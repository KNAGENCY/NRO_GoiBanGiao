// Chup man hinh Game view ra PNG (ScreenCapture can 1 frame; doc file o luot sau). Ghi de file cu.
string p = "C:/Users/Admin/AppData/Local/Temp/claude/D--NRO-GoiBanGiao/8d082123-bbf0-446c-aa97-6342f0eda3a3/scratchpad/shot.png";
UnityEngine.ScreenCapture.CaptureScreenshot(p);
return "da yeu cau chup: " + p + " isPlaying=" + UnityEditor.EditorApplication.isPlaying + " paused=" + UnityEditor.EditorApplication.isPaused;
