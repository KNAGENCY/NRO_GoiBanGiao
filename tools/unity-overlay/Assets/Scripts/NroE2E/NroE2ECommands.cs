#if NRO_E2E || UNITY_EDITOR
using System;
using System.Collections.Generic;
using System.IO;
using System.Security.Cryptography;
using System.Text;
using System.Threading;
using Game1;
using UnityEngine;
using GameChar = Game1.Char;

namespace NroE2E
{
    public sealed class NroE2ECommandIds
    {
        private readonly HashSet<string> accepted =
            new HashSet<string>(StringComparer.Ordinal);

        public bool TryAccept(string commandId, out string errorCode)
        {
            if (string.IsNullOrWhiteSpace(commandId))
            {
                errorCode = "E2E-COMMAND-002";
                return false;
            }
            lock (accepted)
            {
                if (!accepted.Add(commandId))
                {
                    errorCode = "E2E-COMMAND-003";
                    return false;
                }
            }
            errorCode = null;
            return true;
        }
    }

    public sealed class NroE2ESequence
    {
        private long last;

        public bool TryAccept(long current, out string errorCode)
        {
            if (current <= last)
            {
                errorCode = "E2E-SEQUENCE-003";
                return false;
            }
            if (current != last + 1)
            {
                errorCode = "E2E-SEQUENCE-004";
                return false;
            }
            last = current;
            errorCode = null;
            return true;
        }
    }

    public sealed class NroE2EMainThreadQueue
    {
        private readonly int mainThreadId;
        private readonly Queue<Action> pending = new Queue<Action>();

        public NroE2EMainThreadQueue()
        {
            mainThreadId = Thread.CurrentThread.ManagedThreadId;
        }

        public void Enqueue(Action action)
        {
            if (action == null)
            {
                throw new ArgumentNullException(nameof(action));
            }
            lock (pending)
            {
                pending.Enqueue(action);
            }
        }

        public void Drain()
        {
            if (Thread.CurrentThread.ManagedThreadId != mainThreadId)
            {
                throw new InvalidOperationException(
                    "E2E-THREAD-005: command queue chi duoc drain tren Unity main thread");
            }
            while (true)
            {
                Action action;
                lock (pending)
                {
                    if (pending.Count == 0)
                    {
                        return;
                    }
                    action = pending.Dequeue();
                }
                action();
            }
        }
    }

    public static class NroE2EEvidence
    {
        public static string ComputeSha256(string path)
        {
            using (FileStream input = new FileStream(
                path,
                FileMode.Open,
                FileAccess.Read,
                FileShare.Read))
            using (SHA256 sha = SHA256.Create())
            {
                byte[] digest = sha.ComputeHash(input);
                StringBuilder text = new StringBuilder(digest.Length * 2);
                foreach (byte value in digest)
                {
                    text.Append(value.ToString("x2"));
                }
                return text.ToString();
            }
        }

        public static bool TryValidateCapture(
            string path,
            string expectedSha256,
            out string errorCode)
        {
            if (string.IsNullOrWhiteSpace(path) ||
                string.IsNullOrWhiteSpace(expectedSha256) ||
                expectedSha256.Length != 64 ||
                !File.Exists(path))
            {
                errorCode = "E2E-SHA-006";
                return false;
            }
            string actual;
            try
            {
                actual = ComputeSha256(path);
            }
            catch
            {
                errorCode = "E2E-SHA-006";
                return false;
            }
            if (!NroE2ECommands.ConstantTimeTokenEquals(
                actual.ToLowerInvariant(),
                expectedSha256.ToLowerInvariant()))
            {
                errorCode = "E2E-SHA-006";
                return false;
            }
            errorCode = null;
            return true;
        }
    }

    [Serializable]
    public sealed class NroE2EAssetMarker
    {
        public int schema_version;
        public string lock_sha256;
        public string server_client_pair_sha256;
        public NroE2EAssetVersions versions;
    }

    [Serializable]
    public sealed class NroE2EAssetVersions
    {
        public int data;
        public int map;
        public int resource;
        public int runtime_slice;
    }

    public sealed class NroE2ERuntimeSample
    {
        public bool ready;
        public bool connected;
        public string account;
        public int loginAttempt;
        public int characterId;
        public int mapId;
        public int zoneId;
        public string mapName;
        public int hp;
        public int hpMax;
        public int selectedSkill;
        public long selectedSkillLastUse;
        public long selectedSkillReadyAt;
        public long monotonicMs;
        public string damageSourceKind;
        public int damageSourceId = -1;
        public int serverDataVersion;
        public int serverMapVersion;
        public int serverSkillVersion;
        public int serverItemVersion;
        public int clientDataVersion;
        public int clientMapVersion;
        public int clientSkillVersion;
        public int clientItemVersion;
        public bool negotiationActive;
        public string appliedLockSha256;
        public string pairSha256;
        public bool dragonBallFallback;
    }

    public sealed class NroE2ETelemetryProducer
    {
        private bool initialized;
        private int previousHp;
        private int previousCharacterId = -1;
        private int previousMapId = -1;
        private int previousZoneId = -1;
        private int previousSelectedSkill = -1;
        private long previousSkillLastUse;
        private long readyEventAt;
        private bool skillReadyEmitted = true;
        private int previousLoginAttempt;
        private bool negotiationObserved;
        private bool appliedLockEmitted;
        private int readyTransitions;

        public void Observe(
            NroE2ERuntimeSample sample,
            Action<string, NroE2EEventPayload> emit)
        {
            if (sample == null || emit == null)
            {
                throw new ArgumentNullException();
            }
            if (!sample.ready)
            {
                // Resource negotiation dien ra truoc khi nhan vat vao map.
                // Ghi nhan canh len ngay luc do; applied chi phat khi ready.
                if (sample.negotiationActive)
                {
                    ObserveNegotiation(sample, emit);
                }
                initialized = false;
                return;
            }

            if (!initialized)
            {
                initialized = true;
                readyTransitions++;
                previousHp = sample.hp;
                previousCharacterId = sample.characterId;
                previousMapId = sample.mapId;
                previousZoneId = sample.zoneId;
                previousSelectedSkill = sample.selectedSkill;
                previousSkillLastUse = sample.selectedSkillLastUse;
                readyEventAt = sample.selectedSkillReadyAt;
                skillReadyEmitted = sample.monotonicMs >= readyEventAt;
                previousLoginAttempt = sample.loginAttempt;
                emit("login_completed", LoginPayload(sample));
                emit("map_resource_state", ResourcePayload(sample));
                ObserveNegotiation(sample, emit);
                return;
            }

            if (sample.loginAttempt != previousLoginAttempt ||
                sample.characterId != previousCharacterId)
            {
                readyTransitions++;
                previousLoginAttempt = sample.loginAttempt;
                previousCharacterId = sample.characterId;
                emit("login_completed", LoginPayload(sample));
            }

            if (sample.mapId != previousMapId || sample.zoneId != previousZoneId)
            {
                previousMapId = sample.mapId;
                previousZoneId = sample.zoneId;
                emit("map_resource_state", ResourcePayload(sample));
            }

            if (sample.hp != previousHp)
            {
                int before = previousHp;
                int after = sample.hp;
                previousHp = after;
                emit("hp_transition", new NroE2EEventPayload
                {
                    hp_before = before,
                    hp_after = after,
                    amount = System.Math.Abs(before - after),
                    hp = after,
                    hp_max = sample.hpMax
                });
                if (after < before)
                {
                    emit("damage_received", new NroE2EEventPayload
                    {
                        amount = before - after,
                        hp_before = before,
                        hp_after = after,
                        hp = after,
                        hp_max = sample.hpMax,
                        source_kind = sample.damageSourceKind ?? "unknown",
                        source_id = sample.damageSourceId
                    });
                }
            }

            if (sample.selectedSkill != previousSelectedSkill)
            {
                previousSelectedSkill = sample.selectedSkill;
                emit("skill_state", new NroE2EEventPayload
                {
                    action = "selected",
                    selected_skill = sample.selectedSkill
                });
            }
            if (sample.selectedSkillLastUse > previousSkillLastUse)
            {
                previousSkillLastUse = sample.selectedSkillLastUse;
                readyEventAt = sample.selectedSkillReadyAt;
                skillReadyEmitted = false;
                emit("skill_state", new NroE2EEventPayload
                {
                    action = "used",
                    selected_skill = sample.selectedSkill,
                    selected_skill_ready_at_ms = readyEventAt
                });
            }
            if (!skillReadyEmitted && sample.monotonicMs >= readyEventAt)
            {
                skillReadyEmitted = true;
                emit("skill_ready", new NroE2EEventPayload
                {
                    action = "ready",
                    selected_skill = sample.selectedSkill,
                    selected_skill_ready_at_ms = readyEventAt
                });
            }
            ObserveNegotiation(sample, emit);
        }

        private NroE2EEventPayload LoginPayload(NroE2ERuntimeSample sample)
        {
            return new NroE2EEventPayload
            {
                account = sample.account,
                character_id = sample.characterId,
                map_id = sample.mapId,
                zone_id = sample.zoneId,
                hp = sample.hp,
                hp_max = sample.hpMax,
                login_mode = readyTransitions <= 1 ? "login" : "relogin"
            };
        }

        private static NroE2EEventPayload ResourcePayload(
            NroE2ERuntimeSample sample)
        {
            return new NroE2EEventPayload
            {
                map_id = sample.mapId,
                zone_id = sample.zoneId,
                map_name = sample.mapName,
                applied_lock_sha256 = sample.appliedLockSha256,
                pair_sha256 = sample.pairSha256,
                server_data_version = sample.serverDataVersion,
                server_map_version = sample.serverMapVersion,
                server_skill_version = sample.serverSkillVersion,
                server_item_version = sample.serverItemVersion,
                client_data_version = sample.clientDataVersion,
                client_map_version = sample.clientMapVersion,
                client_skill_version = sample.clientSkillVersion,
                client_item_version = sample.clientItemVersion,
                update_required = sample.negotiationActive,
                dragon_ball_fallback = sample.dragonBallFallback
            };
        }

        private void ObserveNegotiation(
            NroE2ERuntimeSample sample,
            Action<string, NroE2EEventPayload> emit)
        {
            if (sample.negotiationActive && !negotiationObserved)
            {
                negotiationObserved = true;
                appliedLockEmitted = false;
                NroE2EEventPayload payload = ResourcePayload(sample);
                payload.negotiation_detected = true;
                emit("warm_negotiation", payload);
            }
            if (negotiationObserved && !sample.negotiationActive &&
                !appliedLockEmitted)
            {
                appliedLockEmitted = true;
                NroE2EEventPayload payload = ResourcePayload(sample);
                payload.negotiation_detected = true;
                emit("resource_lock_applied", payload);
            }
        }
    }

    public static class NroE2ERuntimeReader
    {
        private static NroE2EAssetMarker marker;
        private static bool markerLoaded;

        public static NroE2ERuntimeSample Read()
        {
            GameChar me = GameChar.myCharz();
            NroE2ERuntimeSample sample = new NroE2ERuntimeSample
            {
                ready = me != null && me.charID >= 0 && TileMap.mapID >= 0,
                connected = Session_ME.gI().isConnected(),
                account = NroE2ECommands.CurrentAccount,
                loginAttempt = NroE2ECommands.LoginAttempt,
                characterId = me == null ? -1 : me.charID,
                mapId = TileMap.mapID,
                zoneId = TileMap.zoneID,
                hp = me == null ? 0 : (int)me.cHP,
                hpMax = me == null ? 0 : (int)me.cHPFull,
                selectedSkill = me == null || me.myskill == null ||
                    me.myskill.template == null ? -1 : me.myskill.template.id,
                selectedSkillLastUse = me == null || me.myskill == null
                    ? 0
                    : me.myskill.lastTimeUseThisSkill,
                selectedSkillReadyAt = me == null || me.myskill == null
                    ? 0
                    : me.myskill.lastTimeUseThisSkill + me.myskill.coolDown,
                monotonicMs = NroE2ECommands.MonotonicMilliseconds,
                serverDataVersion = GameScr.vsData,
                serverMapVersion = GameScr.vsMap,
                serverSkillVersion = GameScr.vsSkill,
                serverItemVersion = GameScr.vsItem,
                clientDataVersion = GameScr.vcData,
                clientMapVersion = GameScr.vcMap,
                clientSkillVersion = GameScr.vcSkill,
                clientItemVersion = GameScr.vcItem
            };
            sample.negotiationActive =
                sample.serverDataVersion != sample.clientDataVersion ||
                sample.serverMapVersion != sample.clientMapVersion ||
                sample.serverSkillVersion != sample.clientSkillVersion ||
                sample.serverItemVersion != sample.clientItemVersion ||
                LoginScr.isUpdateData || LoginScr.isUpdateMap ||
                LoginScr.isUpdateSkill || LoginScr.isUpdateItem;
            if (TileMap.mapNames != null && sample.mapId >= 0 &&
                sample.mapId < TileMap.mapNames.Length)
            {
                sample.mapName = TileMap.mapNames[sample.mapId];
            }
            sample.dragonBallFallback = sample.mapId == 1 &&
                !string.Equals(sample.mapName, "Làng Lá", StringComparison.Ordinal);

            NroE2EAssetMarker active = LoadMarker();
            if (active != null)
            {
                sample.appliedLockSha256 = active.lock_sha256;
                sample.pairSha256 = active.server_client_pair_sha256;
            }
            FindDamageSource(me, sample);
            return sample;
        }

        private static void FindDamageSource(
            GameChar me,
            NroE2ERuntimeSample sample)
        {
            if (me == null)
            {
                return;
            }
            Mob closest = null;
            int distance = int.MaxValue;
            for (int index = 0; index < GameScr.vMob.size(); index++)
            {
                Mob mob = (Mob)GameScr.vMob.elementAt(index);
                if (mob == null || mob.cFocus != me || mob.dame <= 0)
                {
                    continue;
                }
                int currentDistance = System.Math.Abs(mob.x - me.cx) +
                    System.Math.Abs(mob.y - me.cy);
                if (closest == null || currentDistance < distance)
                {
                    closest = mob;
                    distance = currentDistance;
                }
            }
            if (closest != null)
            {
                sample.damageSourceKind = "mob";
                sample.damageSourceId = closest.mobId;
                return;
            }
            for (int index = 0; index < GameScr.vCharInMap.size(); index++)
            {
                GameChar player = (GameChar)GameScr.vCharInMap.elementAt(index);
                if (player == null || player == me)
                {
                    continue;
                }
                bool targetsMe = player.charFocus == me;
                if (!targetsMe && player.attChars != null)
                {
                    for (int target = 0; target < player.attChars.Length; target++)
                    {
                        targetsMe |= player.attChars[target] == me;
                    }
                }
                if (targetsMe)
                {
                    sample.damageSourceKind = "player";
                    sample.damageSourceId = player.charID;
                    return;
                }
            }
        }

        private static NroE2EAssetMarker LoadMarker()
        {
            if (markerLoaded)
            {
                return marker;
            }
            markerLoaded = true;
            string path = Path.Combine(
                Application.streamingAssetsPath,
                "nro-e2e-provenance.json");
            if (!File.Exists(path))
            {
                return null;
            }
            try
            {
                marker = JsonUtility.FromJson<NroE2EAssetMarker>(
                    File.ReadAllText(path, Encoding.UTF8));
            }
            catch
            {
                marker = null;
            }
            return marker;
        }
    }

    public static class NroE2ECommands
    {
        private static readonly string[] Supported =
        {
            "login", "wait_ready", "move_to", "select_skill", "target_mob",
            "target_player", "attack", "respawn", "open_npc",
            "choose_npc_menu", "capture", "get_state", "shutdown"
        };

        private static string account;
        private static int loginAttempt;
        private static long selectedSkillReadyAtMs;

        public static string CurrentAccount => account;

        public static int LoginAttempt => loginAttempt;

        public static string[] SupportedCommandNames()
        {
            return (string[])Supported.Clone();
        }

        public static long MonotonicMilliseconds
        {
            get
            {
                return (long)(Time.realtimeSinceStartupAsDouble * 1000.0);
            }
        }

        public static bool ConstantTimeTokenEquals(string expected, string actual)
        {
            byte[] expectedBytes = Encoding.UTF8.GetBytes(expected ?? string.Empty);
            byte[] actualBytes = Encoding.UTF8.GetBytes(actual ?? string.Empty);
            byte[] expectedDigest;
            byte[] actualDigest;
            using (SHA256 sha = SHA256.Create())
            {
                expectedDigest = sha.ComputeHash(expectedBytes);
            }
            using (SHA256 sha = SHA256.Create())
            {
                actualDigest = sha.ComputeHash(actualBytes);
            }
            int difference = 0;
            for (int index = 0; index < expectedDigest.Length; index++)
            {
                difference |= expectedDigest[index] ^ actualDigest[index];
            }
            return difference == 0;
        }

        public static bool TryExecute(
            NroE2ECommandRequest request,
            out string payload,
            out string errorCode)
        {
            payload = null;
            errorCode = null;
            if (request == null || Array.IndexOf(Supported, request.command) < 0)
            {
                errorCode = "E2E-COMMAND-004";
                return false;
            }
            NroE2ECommandPayload input = request.payload ?? new NroE2ECommandPayload();
            GameChar me = GameChar.myCharz();
            switch (request.command)
            {
                case "login":
                    if (string.IsNullOrWhiteSpace(input.account) ||
                        string.IsNullOrWhiteSpace(input.password))
                    {
                        errorCode = "E2E-COMMAND-005";
                        return false;
                    }
                    account = input.account;
                    loginAttempt++;
                    GameMidlet.IP = "127.0.0.1";
                    int gamePort;
                    string gamePortText = NroE2EConfig.GetArgumentValue("-game-port");
                    if (!int.TryParse(gamePortText, out gamePort))
                    {
                        gamePort = 14445;
                    }
                    GameMidlet.PORT = gamePort;
                    if (!Session_ME.gI().isConnected())
                    {
                        GameCanvas.connect();
                    }
                    Service.gI().login(
                        input.account,
                        input.password,
                        GameMidlet.VERSION,
                        0);
                    return true;

                case "wait_ready":
                    if (me == null || me.charID < 0 || TileMap.mapID < 0)
                    {
                        errorCode = "E2E-NOT-READY";
                        return false;
                    }
                    return true;

                case "move_to":
                    if (me == null)
                    {
                        errorCode = "E2E-NOT-READY";
                        return false;
                    }
                    me.cx = input.x;
                    me.cy = input.y;
                    Service.gI().charMove();
                    return true;

                case "select_skill":
                    if (me == null || me.vSkillFight == null)
                    {
                        errorCode = "E2E-NOT-READY";
                        return false;
                    }
                    for (int index = 0; index < me.vSkillFight.size(); index++)
                    {
                        Skill skill = (Skill)me.vSkillFight.elementAt(index);
                        if (skill != null && skill.template != null &&
                            skill.template.id == input.skill_id)
                        {
                            me.myskill = skill;
                            Service.gI().selectSkill(input.skill_id);
                            selectedSkillReadyAtMs =
                                MonotonicMilliseconds + skill.coolDown;
                            return true;
                        }
                    }
                    errorCode = "E2E-COMMAND-006";
                    return false;

                case "target_mob":
                    if (me == null)
                    {
                        errorCode = "E2E-NOT-READY";
                        return false;
                    }
                    for (int index = 0; index < GameScr.vMob.size(); index++)
                    {
                        Mob mob = (Mob)GameScr.vMob.elementAt(index);
                        if (mob != null && mob.mobId == input.target_id)
                        {
                            me.mobFocus = mob;
                            me.charFocus = null;
                            me.npcFocus = null;
                            return true;
                        }
                    }
                    errorCode = "E2E-COMMAND-007";
                    return false;

                case "target_player":
                    if (me == null)
                    {
                        errorCode = "E2E-NOT-READY";
                        return false;
                    }
                    for (int index = 0; index < GameScr.vCharInMap.size(); index++)
                    {
                        GameChar player =
                            (GameChar)GameScr.vCharInMap.elementAt(index);
                        if (player != null && player.charID == input.target_id)
                        {
                            me.charFocus = player;
                            me.mobFocus = null;
                            me.npcFocus = null;
                            return true;
                        }
                    }
                    errorCode = "E2E-COMMAND-007";
                    return false;

                case "attack":
                    GameScr.gI().doFire(false, true);
                    return true;

                case "respawn":
                    Service.gI().wakeUpFromDead();
                    return true;

                case "open_npc":
                    if (me == null)
                    {
                        errorCode = "E2E-NOT-READY";
                        return false;
                    }
                    for (int index = 0; index < GameScr.vNpc.size(); index++)
                    {
                        Npc npc = (Npc)GameScr.vNpc.elementAt(index);
                        if (npc != null && npc.template != null &&
                            npc.template.npcTemplateId == input.npc_id)
                        {
                            me.npcFocus = npc;
                            Service.gI().openMenu(input.npc_id);
                            return true;
                        }
                    }
                    errorCode = "E2E-COMMAND-007";
                    return false;

                case "choose_npc_menu":
                    if (me == null || me.npcFocus == null ||
                        me.npcFocus.template == null)
                    {
                        errorCode = "E2E-COMMAND-007";
                        return false;
                    }
                    Service.gI().menu(
                        me.npcFocus.template.npcTemplateId,
                        input.menu_index,
                        0);
                    return true;

                case "capture":
                case "get_state":
                case "shutdown":
                    // Bridge xu ly ba command co lifecycle/response dac biet.
                    return true;
            }
            errorCode = "E2E-COMMAND-004";
            return false;
        }

        public static NroE2EState ReadState(string lastEvent, long sequence)
        {
            GameChar me = GameChar.myCharz();
            NroE2EState state = new NroE2EState
            {
                account = account,
                character_id = me == null ? -1 : me.charID,
                map_id = TileMap.mapID,
                zone_id = TileMap.zoneID,
                x = me == null ? 0 : me.cx,
                y = me == null ? 0 : me.cy,
                hp = me == null ? 0 : (int)me.cHP,
                hp_max = me == null ? 0 : (int)me.cHPFull,
                chakra = me == null ? 0 : (int)me.cMP,
                chakra_max = me == null ? 0 : (int)me.cMPFull,
                alive = me != null && me.cHP > 0 && me.statusMe != 14,
                selected_skill =
                    me == null || me.myskill == null || me.myskill.template == null
                        ? -1
                        : me.myskill.template.id,
                selected_skill_ready_at_ms =
                    me == null || me.myskill == null
                        ? selectedSkillReadyAtMs
                        : me.myskill.lastTimeUseThisSkill + me.myskill.coolDown,
                last_event = lastEvent,
                sequence = sequence,
                background_layers = new[]
                {
                    new NroE2EBackgroundLayer
                    {
                        resource = "b03",
                        offset_x = -GameScr.cmx / 4,
                        offset_y = -GameScr.cmy / 8
                    }
                }
            };

            List<int> players = new List<int>();
            for (int index = 0; index < GameScr.vCharInMap.size(); index++)
            {
                GameChar player =
                    (GameChar)GameScr.vCharInMap.elementAt(index);
                if (player != null)
                {
                    players.Add(player.charID);
                }
            }
            state.visible_player_ids = players.ToArray();

            List<int> mobs = new List<int>();
            for (int index = 0; index < GameScr.vMob.size(); index++)
            {
                Mob mob = (Mob)GameScr.vMob.elementAt(index);
                if (mob != null)
                {
                    mobs.Add(mob.mobId);
                }
            }
            state.visible_mob_ids = mobs.ToArray();
            if (me != null && me.mobFocus != null)
            {
                state.target_id = me.mobFocus.mobId;
            }
            else if (me != null && me.charFocus != null)
            {
                state.target_id = me.charFocus.charID;
            }
            return state;
        }
    }
}
#endif
