string dataPath = UnityEngine.Application.dataPath.Replace("\\", "/").Replace("\"", "\\\"");
string json = "{\"data_path\":\"" + dataPath
    + "\",\"is_playing\":" + UnityEditor.EditorApplication.isPlaying.ToString().ToLowerInvariant()
    + ",\"is_playing_or_will_change\":" + UnityEditor.EditorApplication.isPlayingOrWillChangePlaymode.ToString().ToLowerInvariant()
    + ",\"is_compiling\":" + UnityEditor.EditorApplication.isCompiling.ToString().ToLowerInvariant()
    + ",\"is_updating\":" + UnityEditor.EditorApplication.isUpdating.ToString().ToLowerInvariant()
    + "}";
return "UI_DOT4_STATE=" + json;
