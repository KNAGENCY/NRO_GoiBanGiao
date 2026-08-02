#if NRO_E2E || UNITY_EDITOR
using System;
using System.Collections;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using UnityEngine;
#if UNITY_EDITOR
using System.Linq;
using UnityEditor;
using UnityEditor.Build.Reporting;
#endif

namespace NroE2E
{
    public sealed class NroE2EBridge : MonoBehaviour
    {
        private readonly NroE2EMainThreadQueue mainThread =
            new NroE2EMainThreadQueue();
        private readonly NroE2ECommandIds commandIds =
            new NroE2ECommandIds();
        private readonly NroE2ETelemetryProducer telemetry =
            new NroE2ETelemetryProducer();
        private readonly object sessionLock = new object();

        private TcpListener listener;
        private TcpClient client;
        private StreamWriter writer;
        private long writerSessionGeneration;
        private long nextSessionGeneration;
        private long activeSessionGeneration;
        private Thread listenerThread;
        private volatile bool stopping;
        private bool applicationQuitRequested;
        private string expectedToken;
        private string clientName;
        private long sequence;
        private string lastEvent;

#if UNITY_EDITOR
        public static void BuildBatch()
        {
            string output = NroE2EConfig.GetArgumentValue("-e2e-output");
            if (string.IsNullOrWhiteSpace(output) || !Path.IsPathRooted(output))
            {
                throw new InvalidOperationException(
                    "E2E-BUILD-001: -e2e-output phai la duong dan tuyet doi");
            }
            string[] scenes = EditorBuildSettings.scenes
                .Where(scene => scene.enabled)
                .Select(scene => scene.path)
                .ToArray();
            if (scenes.Length == 0)
            {
                throw new InvalidOperationException(
                    "E2E-BUILD-001: project khong co scene build duoc");
            }
            BuildTargetGroup group = BuildTargetGroup.Standalone;
            string previous = PlayerSettings.GetScriptingDefineSymbolsForGroup(group);
            string[] symbols = previous.Split(';');
            if (Array.IndexOf(symbols, "NRO_E2E") < 0)
            {
                PlayerSettings.SetScriptingDefineSymbolsForGroup(
                    group,
                    string.IsNullOrWhiteSpace(previous)
                        ? "NRO_E2E"
                        : previous + ";NRO_E2E");
            }
            PlayerSettings.defaultScreenWidth = 1280;
            PlayerSettings.defaultScreenHeight = 720;
            EmbedAssetProvenance();
            Directory.CreateDirectory(Path.GetDirectoryName(output));
            BuildReport report = BuildPipeline.BuildPlayer(new BuildPlayerOptions
            {
                scenes = scenes,
                locationPathName = output,
                target = BuildTarget.StandaloneWindows64,
                options = BuildOptions.Development
            });
            if (report.summary.result != BuildResult.Succeeded)
            {
                throw new InvalidOperationException(
                    "E2E-BUILD-001: Unity build that bai: " +
                    report.summary.result);
            }
            Debug.Log(
                "E2E-BUILD-OK bytes=" + report.summary.totalSize +
                " duration=" + report.summary.totalTime);
        }

        private static void EmbedAssetProvenance()
        {
            string source = Path.Combine(
                Directory.GetCurrentDirectory(),
                ".nro-asset-lock.json");
            if (!File.Exists(source))
            {
                throw new InvalidOperationException(
                    "E2E-ASSET-006: snapshot thieu .nro-asset-lock.json");
            }
            NroE2EAssetMarker marker;
            try
            {
                marker = JsonUtility.FromJson<NroE2EAssetMarker>(
                    File.ReadAllText(source, Encoding.UTF8));
            }
            catch (Exception error)
            {
                throw new InvalidOperationException(
                    "E2E-ASSET-006: marker provenance khong doc duoc",
                    error);
            }
            if (marker == null || marker.schema_version != 1 ||
                string.IsNullOrWhiteSpace(marker.lock_sha256) ||
                marker.lock_sha256.Length != 64 ||
                string.IsNullOrWhiteSpace(marker.server_client_pair_sha256) ||
                marker.server_client_pair_sha256.Length != 64 ||
                marker.versions == null)
            {
                throw new InvalidOperationException(
                    "E2E-ASSET-006: marker provenance sai schema");
            }
            string destination = Path.Combine(
                Application.dataPath,
                "StreamingAssets",
                "nro-e2e-provenance.json");
            Directory.CreateDirectory(Path.GetDirectoryName(destination));
            File.Copy(source, destination, true);
            AssetDatabase.Refresh();
        }
#endif

        [RuntimeInitializeOnLoadMethod(RuntimeInitializeLoadType.AfterSceneLoad)]
        private static void Bootstrap()
        {
#if NRO_E2E
            string portText = NroE2EConfig.GetArgumentValue("-telemetry-port");
            string token = NroE2EConfig.GetArgumentValue("-telemetry-token");
            if (string.IsNullOrWhiteSpace(portText))
            {
                return;
            }
            int port;
            if (!int.TryParse(portText, out port) || port < 1 || port > 65535 ||
                string.IsNullOrWhiteSpace(token) || token.Length < 32)
            {
                Debug.LogError("E2E-TELEMETRY-001: cau hinh telemetry khong hop le");
                return;
            }
            GameObject host = new GameObject("NroE2EBridge");
            DontDestroyOnLoad(host);
            NroE2EBridge bridge = host.AddComponent<NroE2EBridge>();
            bridge.StartListener(
                port,
                token,
                NroE2EConfig.GetArgumentValue("-client") ?? "unknown");
#endif
        }

        public void StartListener(int port, string token, string name)
        {
            if (listener != null)
            {
                throw new InvalidOperationException(
                    "E2E-TELEMETRY-001: listener da chay");
            }
            if (port < 1 || port > 65535 || string.IsNullOrWhiteSpace(token) ||
                token.Length < 32)
            {
                throw new ArgumentException(
                    "E2E-TELEMETRY-001: port/token khong hop le");
            }
            expectedToken = token;
            clientName = string.IsNullOrWhiteSpace(name) ? "unknown" : name;
            listener = new TcpListener(IPAddress.Loopback, port);
            listener.Start(1);
            listenerThread = new Thread(ListenLoop)
            {
                IsBackground = true,
                Name = "NroE2E-loopback-jsonl"
            };
            listenerThread.Start();
        }

        private void ListenLoop()
        {
            while (!stopping)
            {
                TcpClient accepted = null;
                long sessionGeneration = 0;
                try
                {
                    accepted = listener.AcceptTcpClient();
                    accepted.NoDelay = true;
                    sessionGeneration = BeginSession(accepted);
                    if (sessionGeneration == 0)
                    {
                        return;
                    }
                    RunSession(accepted, sessionGeneration);
                }
                catch (SocketException)
                {
                    if (!stopping)
                    {
                        Debug.LogError(
                            "E2E-TELEMETRY-002: socket listener dung bat thuong");
                    }
                    return;
                }
                catch (IOException)
                {
                    // Reset/EOF cua mot session khong duoc giet listener.
                    // finally se dispose het tai nguyen roi vong lap accept lai.
                    if (stopping)
                    {
                        return;
                    }
                }
                catch (Exception)
                {
                    if (!stopping)
                    {
                        // Khong log request/exception text vi message dau co token.
                        Debug.LogError("E2E-TELEMETRY-999: session dung bat thuong");
                    }
                }
                finally
                {
                    InvalidateSession(sessionGeneration);
                    if (accepted != null)
                    {
                        try
                        {
                            accepted.Close();
                        }
                        catch (IOException)
                        {
                        }
                    }
                }
            }
        }

        private long BeginSession(TcpClient accepted)
        {
            lock (sessionLock)
            {
                if (stopping)
                {
                    return 0;
                }
                nextSessionGeneration++;
                activeSessionGeneration = nextSessionGeneration;
                client = accepted;
                writer = null;
                writerSessionGeneration = 0;
                return activeSessionGeneration;
            }
        }

        private bool IsSessionActive(long sessionGeneration)
        {
            lock (sessionLock)
            {
                return !stopping &&
                    sessionGeneration != 0 &&
                    activeSessionGeneration == sessionGeneration;
            }
        }

        private void InvalidateSession(long sessionGeneration)
        {
            lock (sessionLock)
            {
                InvalidateSessionWhileLocked(sessionGeneration);
            }
        }

        private TcpClient InvalidateSessionWhileLocked(long sessionGeneration)
        {
            if (sessionGeneration == 0 ||
                activeSessionGeneration != sessionGeneration)
            {
                return null;
            }
            TcpClient invalidated = client;
            activeSessionGeneration = 0;
            client = null;
            if (writerSessionGeneration == sessionGeneration)
            {
                writer = null;
                writerSessionGeneration = 0;
            }
            return invalidated;
        }

        private void RunSession(TcpClient accepted, long sessionGeneration)
        {
            using (NetworkStream stream = accepted.GetStream())
            using (StreamReader reader = new StreamReader(
                stream,
                new UTF8Encoding(false),
                false,
                4096,
                true))
            using (StreamWriter sessionWriter = new StreamWriter(
                stream,
                new UTF8Encoding(false),
                4096,
                true))
            {
                sessionWriter.AutoFlush = true;
                sessionWriter.NewLine = "\n";
                lock (sessionLock)
                {
                    if (activeSessionGeneration != sessionGeneration ||
                        stopping)
                    {
                        return;
                    }
                    writer = sessionWriter;
                    writerSessionGeneration = sessionGeneration;
                }
                try
                {
                    string first = reader.ReadLine();
                    NroE2EAuthMessage auth = Parse<NroE2EAuthMessage>(first);
                    if (auth == null ||
                        !NroE2ECommands.ConstantTimeTokenEquals(
                            expectedToken,
                            auth.token))
                    {
                        Send(sessionGeneration, new NroE2EAuthResponse
                        {
                            ok = false,
                            error_code = "E2E-AUTH-001"
                        });
                        return;
                    }
                    Send(
                        sessionGeneration,
                        new NroE2EAuthResponse { ok = true });

                    while (!stopping && IsSessionActive(sessionGeneration))
                    {
                        string line = reader.ReadLine();
                        if (line == null)
                        {
                            return;
                        }
                        NroE2ECommandRequest request =
                            Parse<NroE2ECommandRequest>(line);
                        if (request == null)
                        {
                            SendError(
                                sessionGeneration,
                                null,
                                "E2E-COMMAND-001",
                                "invalid JSON request");
                            continue;
                        }
                        string idError;
                        if (!commandIds.TryAccept(
                            request.command_id,
                            out idError))
                        {
                            SendError(
                                sessionGeneration,
                                request.command_id,
                                idError,
                                "command_id rejected");
                            continue;
                        }
                        // Luong socket chi xep viec. Moi API Game1/Unity nam sau
                        // action nay va chi duoc goi khi Update drain tren main thread.
                        mainThread.Enqueue(
                            () => ExecuteOnMainThread(
                                request,
                                sessionGeneration));
                    }
                }
                finally
                {
                    InvalidateSession(sessionGeneration);
                }
            }
        }

        private void Update()
        {
            mainThread.Drain();
            long sessionGeneration;
            lock (sessionLock)
            {
                sessionGeneration = activeSessionGeneration;
            }
            if (sessionGeneration != 0 && IsSessionActive(sessionGeneration))
            {
                telemetry.Observe(
                    NroE2ERuntimeReader.Read(),
                    (type, payload) => EmitEvent(
                        sessionGeneration,
                        type,
                        payload));
            }
        }

        private void ExecuteOnMainThread(
            NroE2ECommandRequest request,
            long sessionGeneration)
        {
            if (!IsSessionActive(sessionGeneration))
            {
                return;
            }
            string payload;
            string errorCode;
            if (!NroE2ECommands.TryExecute(request, out payload, out errorCode))
            {
                SendError(
                    sessionGeneration,
                    request.command_id,
                    errorCode,
                    "command rejected");
                return;
            }

            switch (request.command)
            {
                case "get_state":
                    payload = JsonUtility.ToJson(
                        NroE2ECommands.ReadState(lastEvent, sequence));
                    SendSuccess(
                        sessionGeneration,
                        request.command_id,
                        payload);
                    return;

                case "capture":
                    StartCoroutine(
                        CaptureAndRespond(request, sessionGeneration));
                    return;

                case "shutdown":
                    EmitEvent(
                        sessionGeneration,
                        "telemetry_shutdown",
                        new NroE2EEventPayload());
                    SendSuccess(
                        sessionGeneration,
                        request.command_id,
                        "{}");
                    RequestApplicationQuit();
#if UNITY_EDITOR
                    StopBridge();
#endif
                    return;

                default:
                    EmitEvent(
                        sessionGeneration,
                        "command_completed",
                        new NroE2EEventPayload
                        {
                            command = request.command
                        });
                    SendSuccess(
                        sessionGeneration,
                        request.command_id,
                        payload ?? "{}");
                    return;
            }
        }

        private IEnumerator CaptureAndRespond(
            NroE2ECommandRequest request,
            long sessionGeneration)
        {
            if (!IsSessionActive(sessionGeneration))
            {
                yield break;
            }
            string requested = request.payload == null ? null : request.payload.path;
            if (string.IsNullOrWhiteSpace(requested) ||
                !Path.IsPathRooted(requested) ||
                !string.Equals(
                    Path.GetExtension(requested),
                    ".png",
                    StringComparison.OrdinalIgnoreCase))
            {
                SendError(
                    sessionGeneration,
                    request.command_id,
                    "E2E-SHA-006",
                    "capture path rejected");
                yield break;
            }
            if (!IsSessionActive(sessionGeneration))
            {
                yield break;
            }
            string path = Path.GetFullPath(requested);
            string parent = Path.GetDirectoryName(path);
            Directory.CreateDirectory(parent);
            ScreenCapture.CaptureScreenshot(path);

            long deadline = NroE2ECommands.MonotonicMilliseconds + 10000;
            long previousLength = -1;
            int stableFrames = 0;
            while (NroE2ECommands.MonotonicMilliseconds < deadline)
            {
                if (!IsSessionActive(sessionGeneration))
                {
                    yield break;
                }
                yield return new WaitForEndOfFrame();
                if (!IsSessionActive(sessionGeneration))
                {
                    yield break;
                }
                if (!File.Exists(path))
                {
                    continue;
                }
                long length;
                try
                {
                    using (FileStream input = new FileStream(
                        path,
                        FileMode.Open,
                        FileAccess.Read,
                        FileShare.Read))
                    {
                        length = input.Length;
                    }
                }
                catch (IOException)
                {
                    continue;
                }
                if (length > 0 && length == previousLength)
                {
                    stableFrames++;
                    if (stableFrames >= 2)
                    {
                        if (!IsSessionActive(sessionGeneration))
                        {
                            yield break;
                        }
                        string sha = NroE2EEvidence.ComputeSha256(path);
                        NroE2EEventPayload eventPayload =
                            new NroE2EEventPayload
                            {
                                path = path,
                                sha256 = sha
                            };
                        EmitEvent(
                            sessionGeneration,
                            "capture_completed",
                            eventPayload);
                        SendSuccess(
                            sessionGeneration,
                            request.command_id,
                            JsonUtility.ToJson(eventPayload));
                        yield break;
                    }
                }
                else
                {
                    stableFrames = 0;
                    previousLength = length;
                }
            }
            SendError(
                sessionGeneration,
                request.command_id,
                "E2E-SHA-006",
                "capture flush timeout");
        }

        private void EmitEvent(
            long sessionGeneration,
            string type,
            NroE2EEventPayload payload)
        {
            if (!IsSessionActive(sessionGeneration))
            {
                return;
            }
            sequence++;
            lastEvent = type;
            Send(sessionGeneration, new NroE2EEvent
            {
                monotonic_ms = NroE2ECommands.MonotonicMilliseconds,
                utc = DateTime.UtcNow.ToString("O"),
                client = clientName,
                sequence = sequence,
                type = type,
                payload = payload ?? new NroE2EEventPayload()
            });
        }

        private void SendSuccess(
            long sessionGeneration,
            string commandId,
            string payload)
        {
            Send(sessionGeneration, new NroE2EResponse
            {
                command_id = commandId,
                ok = true,
                payload = payload
            });
        }

        private void SendError(
            long sessionGeneration,
            string commandId,
            string code,
            string message)
        {
            Send(sessionGeneration, new NroE2EResponse
            {
                command_id = commandId,
                ok = false,
                error_code = code,
                message = message
            });
        }

        private void Send(long sessionGeneration, object value)
        {
            string json = JsonUtility.ToJson(value);
            TcpClient failedClient = null;
            lock (sessionLock)
            {
                if (stopping ||
                    activeSessionGeneration != sessionGeneration ||
                    writerSessionGeneration != sessionGeneration ||
                    writer == null)
                {
                    return;
                }
                try
                {
                    writer.WriteLine(json);
                }
                catch (IOException)
                {
                    failedClient = InvalidateSessionWhileLocked(
                        sessionGeneration);
                }
                catch (ObjectDisposedException)
                {
                    failedClient = InvalidateSessionWhileLocked(
                        sessionGeneration);
                }
            }
            if (failedClient != null)
            {
                failedClient.Close();
            }
        }

        private static T Parse<T>(string json) where T : class
        {
            if (string.IsNullOrWhiteSpace(json))
            {
                return null;
            }
            try
            {
                return JsonUtility.FromJson<T>(json);
            }
            catch
            {
                return null;
            }
        }

        private static string EscapeJson(string value)
        {
            return (value ?? string.Empty)
                .Replace("\\", "\\\\")
                .Replace("\"", "\\\"");
        }

        public void StopBridge()
        {
            TcpClient activeClient;
            lock (sessionLock)
            {
                stopping = true;
                activeClient = InvalidateSessionWhileLocked(
                    activeSessionGeneration);
            }
            Thread worker = listenerThread;
            try
            {
                listener?.Stop();
            }
            catch (SocketException)
            {
            }
            try
            {
                activeClient?.Close();
            }
            catch (IOException)
            {
            }
            if (worker != null && worker != Thread.CurrentThread)
            {
                worker.Join(2000);
            }
            if (worker == null || !worker.IsAlive)
            {
                listenerThread = null;
                listener = null;
                lock (sessionLock)
                {
                    client = null;
                    writer = null;
                    writerSessionGeneration = 0;
                    activeSessionGeneration = 0;
                }
            }
        }

        private void RequestApplicationQuit()
        {
            applicationQuitRequested = true;
#if !UNITY_EDITOR
            Application.Quit();
#endif
        }

        private void OnDestroy()
        {
            StopBridge();
        }
    }
}
#endif
