#if UNITY_EDITOR
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Reflection;
using System.Text;
using System.Threading;
using UnityEditor;
using UnityEngine;

namespace NroE2E
{
    // Project legacy nay khong cai Unity Test Framework/NUnit. Vi vay contract
    // EditMode dung cung khuon batch assertion nhu NroResourcePreflightTests:
    // chay trong Editor, fail bang ma on dinh va exit khac 0.
    public static class NroE2EContractTests
    {
        private sealed class ContractException : Exception
        {
            public readonly string Code;

            public ContractException(string code, string message)
                : base(message)
            {
                Code = code;
            }
        }

        private sealed class AuthenticatedSession : IDisposable
        {
            public readonly TcpClient Socket;
            public readonly NetworkStream Stream;
            public readonly StreamReader Reader;
            public readonly StreamWriter Writer;

            public AuthenticatedSession(TcpClient socket)
            {
                Socket = socket;
                Stream = socket.GetStream();
                Reader = new StreamReader(
                    Stream,
                    new UTF8Encoding(false),
                    false,
                    4096,
                    true);
                Writer = new StreamWriter(
                    Stream,
                    new UTF8Encoding(false),
                    4096,
                    true)
                {
                    AutoFlush = true,
                    NewLine = "\n"
                };
            }

            public void Dispose()
            {
                Writer.Dispose();
                Reader.Dispose();
                Stream.Dispose();
                Socket.Dispose();
            }
        }

        public static void RunBatch()
        {
            try
            {
                TestToken();
                TestCommandIds();
                TestSequence();
                TestMainThread();
                TestScreenshotSha();
                TestReconnectAfterRejectedSession();
                TestQueuedCommandCannotCrossSession();
                TestSchema();
                Debug.Log("NRO-E2E-CONTRACT-OK assertions=50");
                EditorApplication.Exit(0);
            }
            catch (ContractException error)
            {
                Debug.LogError(error.Code + ": " + error.Message);
                EditorApplication.Exit(1);
            }
            catch (Exception error)
            {
                Debug.LogError("E2E-CONTRACT-999: " + error);
                EditorApplication.Exit(1);
            }
        }

        private static void Require(bool condition, string code, string message)
        {
            if (!condition)
            {
                throw new ContractException(code, message);
            }
        }

        private static void TestToken()
        {
            const string token = "token-test-32-byte-khong-secret-0001";
            Require(
                !NroE2ECommands.ConstantTimeTokenEquals(token, null),
                "E2E-AUTH-001",
                "missing token was accepted");
            Require(
                !NroE2ECommands.ConstantTimeTokenEquals(token, "token-sai"),
                "E2E-AUTH-001",
                "wrong token was accepted");
            Require(
                NroE2ECommands.ConstantTimeTokenEquals(token, token),
                "E2E-AUTH-001",
                "matching token was rejected");
        }

        private static void TestCommandIds()
        {
            NroE2ECommandIds ids = new NroE2ECommandIds();
            string error;
            Require(
                !ids.TryAccept(null, out error) && error == "E2E-COMMAND-002",
                "E2E-COMMAND-002",
                "missing command_id was accepted");
            Require(
                ids.TryAccept("cmd-1", out error),
                "E2E-COMMAND-003",
                "first command_id was rejected");
            Require(
                !ids.TryAccept("cmd-1", out error) &&
                error == "E2E-COMMAND-003",
                "E2E-COMMAND-003",
                "duplicate command_id was accepted");
        }

        private static void TestSequence()
        {
            string error;
            NroE2ESequence sequence = new NroE2ESequence();
            Require(
                sequence.TryAccept(1, out error),
                "E2E-SEQUENCE-004",
                "sequence 1 was rejected");
            Require(
                !sequence.TryAccept(0, out error) &&
                error == "E2E-SEQUENCE-003",
                "E2E-SEQUENCE-003",
                "backward sequence was accepted");

            sequence = new NroE2ESequence();
            sequence.TryAccept(1, out error);
            Require(
                !sequence.TryAccept(3, out error) &&
                error == "E2E-SEQUENCE-004",
                "E2E-SEQUENCE-004",
                "sequence gap was accepted");
        }

        private static void TestMainThread()
        {
            NroE2EMainThreadQueue queue = new NroE2EMainThreadQueue();
            int mainThreadId = Thread.CurrentThread.ManagedThreadId;
            int executedThreadId = -1;
            Exception backgroundError = null;
            Thread background = new Thread(() =>
            {
                try
                {
                    queue.Enqueue(() =>
                    {
                        executedThreadId = Thread.CurrentThread.ManagedThreadId;
                    });
                }
                catch (Exception error)
                {
                    backgroundError = error;
                }
            });
            background.Start();
            Require(
                background.Join(2000),
                "E2E-THREAD-005",
                "background enqueue did not finish");
            Require(
                backgroundError == null,
                "E2E-THREAD-005",
                "background enqueue failed");
            Require(
                executedThreadId == -1,
                "E2E-THREAD-005",
                "command executed directly on background thread");
            queue.Drain();
            Require(
                executedThreadId == mainThreadId,
                "E2E-THREAD-005",
                "command did not execute on Unity main thread");
        }

        private static void TestScreenshotSha()
        {
            string root = Path.Combine(
                Path.GetTempPath(),
                "nro-e2e-contract-" + Guid.NewGuid().ToString("N"));
            Directory.CreateDirectory(root);
            try
            {
                string path = Path.Combine(root, "capture.png");
                File.WriteAllBytes(path, new byte[] { 1, 3, 7, 9 });
                string sha = NroE2EEvidence.ComputeSha256(path);
                string error;
                Require(
                    !NroE2EEvidence.TryValidateCapture(path, null, out error) &&
                    error == "E2E-SHA-006",
                    "E2E-SHA-006",
                    "missing screenshot SHA was accepted");
                Require(
                    !NroE2EEvidence.TryValidateCapture(
                        path,
                        new string('0', 64),
                        out error) &&
                    error == "E2E-SHA-006",
                    "E2E-SHA-006",
                    "wrong screenshot SHA was accepted");
                Require(
                    NroE2EEvidence.TryValidateCapture(path, sha, out error),
                    "E2E-SHA-006",
                    "correct screenshot SHA was rejected");
            }
            finally
            {
                Directory.Delete(root, true);
            }
        }

        private static void TestReconnectAfterRejectedSession()
        {
            const string token = "token-test-32-byte-khong-secret-0001";
            TcpListener probe = new TcpListener(IPAddress.Loopback, 0);
            probe.Start();
            int port = ((IPEndPoint)probe.LocalEndpoint).Port;
            probe.Stop();

            GameObject host = new GameObject("NroE2EContractBridge");
            NroE2EBridge bridge = host.AddComponent<NroE2EBridge>();
            try
            {
                bridge.StartListener(port, token, "contract-client");
                NroE2EAuthResponse rejected = ExchangeAuth(
                    port,
                    "wrong-token-32-byte-not-secret-0000");
                Require(
                    rejected != null && !rejected.ok &&
                    rejected.error_code == "E2E-AUTH-001",
                    "E2E-SESSION-009",
                    "wrong-auth session did not receive rejection");

                NroE2EAuthResponse firstValid = ExchangeAuth(port, token);
                Require(
                    firstValid != null && firstValid.ok,
                    "E2E-SESSION-009",
                    "valid reconnect after wrong auth failed");

                // ExchangeAuth dispose client ngay sau auth OK, tao EOF that.
                // Listener phai don session do va accept tiep mot client khac.
                NroE2EAuthResponse afterEof = ExchangeAuth(port, token);
                Require(
                    afterEof != null && afterEof.ok,
                    "E2E-SESSION-009",
                    "valid reconnect after EOF failed");
            }
            finally
            {
                bridge.StopBridge();
                Thread listenerWorker = (Thread)typeof(NroE2EBridge)
                    .GetField(
                        "listenerThread",
                        BindingFlags.Instance | BindingFlags.NonPublic)
                    .GetValue(bridge);
                Require(
                    listenerWorker == null || !listenerWorker.IsAlive,
                    "E2E-SESSION-010",
                    "listener thread remained alive after StopBridge");
                UnityEngine.Object.DestroyImmediate(host);
            }
        }

        private static NroE2EAuthResponse ExchangeAuth(int port, string token)
        {
            DateTime deadline = DateTime.UtcNow.AddSeconds(2);
            TcpClient socket = null;
            while (socket == null && DateTime.UtcNow < deadline)
            {
                TcpClient candidate = new TcpClient();
                try
                {
                    candidate.Connect(IPAddress.Loopback, port);
                    socket = candidate;
                }
                catch (SocketException)
                {
                    candidate.Dispose();
                    Thread.Sleep(20);
                }
            }
            if (socket == null)
            {
                return null;
            }
            using (socket)
            {
                socket.ReceiveTimeout = 2000;
                socket.SendTimeout = 2000;
                using (NetworkStream stream = socket.GetStream())
                using (StreamReader reader = new StreamReader(
                    stream,
                    new UTF8Encoding(false),
                    false,
                    4096,
                    true))
                using (StreamWriter writer = new StreamWriter(
                    stream,
                    new UTF8Encoding(false),
                    4096,
                    true))
                {
                    writer.AutoFlush = true;
                    writer.NewLine = "\n";
                    writer.WriteLine(JsonUtility.ToJson(
                        new NroE2EAuthMessage { token = token }));
                    try
                    {
                        string line = reader.ReadLine();
                        return string.IsNullOrWhiteSpace(line)
                            ? null
                            : JsonUtility.FromJson<NroE2EAuthResponse>(line);
                    }
                    catch (IOException)
                    {
                        return null;
                    }
                }
            }
        }

        private static void TestQueuedCommandCannotCrossSession()
        {
            const string token = "token-test-32-byte-khong-secret-0001";
            TcpListener probe = new TcpListener(IPAddress.Loopback, 0);
            probe.Start();
            int port = ((IPEndPoint)probe.LocalEndpoint).Port;
            probe.Stop();

            GameObject host = new GameObject("NroE2ESessionIsolationBridge");
            NroE2EBridge bridge = host.AddComponent<NroE2EBridge>();
            AuthenticatedSession sessionB = null;
            try
            {
                bridge.StartListener(port, token, "session-isolation");
                using (AuthenticatedSession sessionA =
                    OpenAuthenticatedSession(port, token))
                {
                    Require(
                        sessionA != null,
                        "E2E-SESSION-011",
                        "session A did not authenticate");
                    sessionA.Writer.WriteLine(JsonUtility.ToJson(
                        new NroE2ECommandRequest
                        {
                            command_id = "stale-shutdown",
                            command = "shutdown",
                            payload = new NroE2ECommandPayload()
                        }));
                    // Dispose tao EOF sau khi command da nam truoc EOF tren TCP.
                    // B chi auth duoc sau khi accept-loop da xu ly xong A.
                }

                sessionB = OpenAuthenticatedSession(port, token);
                Require(
                    sessionB != null,
                    "E2E-SESSION-011",
                    "session B did not reconnect");
                Thread.Sleep(50);
                InvokeBridgeUpdate(bridge);

                Thread.Sleep(200);
                Require(
                    sessionB.Socket.Available == 0,
                    "E2E-SESSION-011",
                    "session B received response/event queued by session A");
                bool stopping = (bool)typeof(NroE2EBridge)
                    .GetField(
                        "stopping",
                        BindingFlags.Instance | BindingFlags.NonPublic)
                    .GetValue(bridge);
                Require(
                    !stopping,
                    "E2E-SESSION-011",
                    "stale shutdown command from session A was executed");

                sessionB.Socket.ReceiveTimeout = 2000;
                sessionB.Writer.WriteLine(JsonUtility.ToJson(
                    new NroE2ECommandRequest
                    {
                        command_id = "live-wait",
                        command = "not_supported",
                        payload = new NroE2ECommandPayload()
                    }));
                DateTime responseDeadline = DateTime.UtcNow.AddSeconds(2);
                NroE2EResponse ownResponse = null;
                while (ownResponse == null &&
                    DateTime.UtcNow < responseDeadline)
                {
                    if (sessionB.Socket.Available == 0)
                    {
                        Thread.Sleep(20);
                        InvokeBridgeUpdate(bridge);
                        continue;
                    }
                    string ownLine = TryReadLine(sessionB.Reader);
                    NroE2EResponse candidate =
                        string.IsNullOrWhiteSpace(ownLine)
                            ? null
                            : JsonUtility.FromJson<NroE2EResponse>(ownLine);
                    if (candidate != null &&
                        candidate.command_id == "live-wait")
                    {
                        ownResponse = candidate;
                    }
                }
                Require(
                    ownResponse != null &&
                    ownResponse.command_id == "live-wait",
                    "E2E-SESSION-011",
                    "session B could not execute its own command");
            }
            finally
            {
                if (sessionB != null)
                {
                    sessionB.Dispose();
                }
                bridge.StopBridge();
                UnityEngine.Object.DestroyImmediate(host);
            }
        }

        private static AuthenticatedSession OpenAuthenticatedSession(
            int port,
            string token)
        {
            DateTime deadline = DateTime.UtcNow.AddSeconds(2);
            while (DateTime.UtcNow < deadline)
            {
                TcpClient socket = new TcpClient
                {
                    ReceiveTimeout = 2000,
                    SendTimeout = 2000
                };
                try
                {
                    socket.Connect(IPAddress.Loopback, port);
                    AuthenticatedSession session =
                        new AuthenticatedSession(socket);
                    session.Writer.WriteLine(JsonUtility.ToJson(
                        new NroE2EAuthMessage { token = token }));
                    string line = TryReadLine(session.Reader);
                    NroE2EAuthResponse response =
                        string.IsNullOrWhiteSpace(line)
                            ? null
                            : JsonUtility.FromJson<NroE2EAuthResponse>(line);
                    if (response != null && response.ok)
                    {
                        return session;
                    }
                    session.Dispose();
                    return null;
                }
                catch (SocketException)
                {
                    socket.Dispose();
                    Thread.Sleep(20);
                }
                catch (IOException)
                {
                    socket.Dispose();
                    return null;
                }
            }
            return null;
        }

        private static string TryReadLine(StreamReader reader)
        {
            try
            {
                return reader.ReadLine();
            }
            catch (IOException)
            {
                return null;
            }
        }

        private static void InvokeBridgeUpdate(NroE2EBridge bridge)
        {
            typeof(NroE2EBridge)
                .GetMethod(
                    "Update",
                    BindingFlags.Instance | BindingFlags.NonPublic)
                .Invoke(bridge, null);
        }

        private static void TestSchema()
        {
            string[] stateFields =
            {
                "account", "character_id", "map_id", "zone_id", "x", "y",
                "hp", "hp_max", "chakra", "chakra_max", "alive",
                "visible_player_ids", "visible_mob_ids", "selected_skill",
                "selected_skill_ready_at_ms", "target_id", "background_layers",
                "last_event", "sequence"
            };
            foreach (string field in stateFields)
            {
                Require(
                    typeof(NroE2EState).GetField(
                        field,
                        BindingFlags.Instance | BindingFlags.Public) != null,
                    "E2E-STATE-007",
                    "missing state field " + field);
            }

            string[] eventFields =
            {
                "monotonic_ms", "utc", "client", "sequence", "type", "payload"
            };
            foreach (string field in eventFields)
            {
                Require(
                    typeof(NroE2EEvent).GetField(
                        field,
                        BindingFlags.Instance | BindingFlags.Public) != null,
                    "E2E-EVENT-008",
                    "missing event field " + field);
            }

            string[] expectedCommands =
            {
                "login", "wait_ready", "move_to", "select_skill", "target_mob",
                "target_player", "attack", "respawn", "open_npc",
                "choose_npc_menu", "capture", "get_state", "shutdown"
            };
            HashSet<string> actual = new HashSet<string>(
                NroE2ECommands.SupportedCommandNames(),
                StringComparer.Ordinal);
            Require(
                actual.SetEquals(expectedCommands),
                "E2E-COMMAND-004",
                "supported command set is incomplete");
        }
    }
}
#endif
