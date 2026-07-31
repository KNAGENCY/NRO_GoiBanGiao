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

    public static class NroE2ECommands
    {
        private static readonly string[] Supported =
        {
            "login", "wait_ready", "move_to", "select_skill", "target_mob",
            "target_player", "attack", "respawn", "open_npc",
            "choose_npc_menu", "capture", "get_state", "shutdown"
        };

        private static string account;
        private static long selectedSkillReadyAtMs;

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
                selected_skill_ready_at_ms = selectedSkillReadyAtMs,
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
