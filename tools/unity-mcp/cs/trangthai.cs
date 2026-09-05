var sb = new System.Text.StringBuilder();
sb.Append("isPlaying=").Append(UnityEditor.EditorApplication.isPlaying);
sb.Append(" scene=").Append(UnityEngine.SceneManagement.SceneManager.GetActiveScene().name);
sb.Append(" screen=").Append(UnityEngine.Screen.width).Append("x").Append(UnityEngine.Screen.height);
sb.Append(" focused=").Append(UnityEditorInternal.InternalEditorUtility.isApplicationActive);
return sb.ToString();
