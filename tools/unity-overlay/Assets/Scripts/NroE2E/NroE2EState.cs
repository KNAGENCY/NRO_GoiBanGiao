#if NRO_E2E || UNITY_EDITOR
using System;

namespace NroE2E
{
    [Serializable]
    public sealed class NroE2EBackgroundLayer
    {
        public string resource;
        public int offset_x;
        public int offset_y;
    }

    [Serializable]
    public sealed class NroE2EState
    {
        public string account;
        public int character_id;
        public int map_id;
        public int zone_id;
        public int x;
        public int y;
        public int hp;
        public int hp_max;
        public int chakra;
        public int chakra_max;
        public bool alive;
        public int[] visible_player_ids = new int[0];
        public int[] visible_mob_ids = new int[0];
        public int selected_skill;
        public long selected_skill_ready_at_ms;
        public int target_id = -1;
        public NroE2EBackgroundLayer[] background_layers =
            new NroE2EBackgroundLayer[0];
        public string last_event;
        public long sequence;
    }

    [Serializable]
    public sealed class NroE2EEvent
    {
        public long monotonic_ms;
        public string utc;
        public string client;
        public long sequence;
        public string type;
        public NroE2EEventPayload payload = new NroE2EEventPayload();
    }

    [Serializable]
    public sealed class NroE2EEventPayload
    {
        public string command;
        public string path;
        public string sha256;
        public string account;
        public string login_mode;
        public string action;
        public string source_kind;
        public string map_name;
        public string applied_lock_sha256;
        public string pair_sha256;
        public int character_id = -1;
        public int map_id = -1;
        public int zone_id = -1;
        public int source_id = -1;
        public int amount;
        public int hp_before;
        public int hp_after;
        public int hp;
        public int hp_max;
        public int selected_skill = -1;
        public long selected_skill_ready_at_ms;
        public int server_data_version = -1;
        public int server_map_version = -1;
        public int server_skill_version = -1;
        public int server_item_version = -1;
        public int client_data_version = -1;
        public int client_map_version = -1;
        public int client_skill_version = -1;
        public int client_item_version = -1;
        public bool negotiation_detected;
        public bool update_required;
        public bool dragon_ball_fallback;
    }

    [Serializable]
    public sealed class NroE2ECommandPayload
    {
        public string account;
        public string password;
        public int x;
        public int y;
        public int skill_id;
        public int target_id;
        public int npc_id;
        public int menu_index;
        public string path;
        public int timeout_ms;
    }

    [Serializable]
    public sealed class NroE2ECommandRequest
    {
        public string command_id;
        public string command;
        public NroE2ECommandPayload payload;
    }

    [Serializable]
    public sealed class NroE2EResponse
    {
        public string command_id;
        public bool ok;
        public string error_code;
        public string message;
        public string payload;
    }

    [Serializable]
    public sealed class NroE2EAuthMessage
    {
        public string token;
    }

    [Serializable]
    public sealed class NroE2EAuthResponse
    {
        public bool ok;
        public string error_code;
        public string type = "auth";
    }
}
#endif
