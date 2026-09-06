// Chup man hinh Game view ra PNG. Ep Editor ve lai truoc (khi Unity KHONG co focus, Game view khong tu ve lai
// -> ScreenCapture tra ve khung cu, vi du man "Xin cho" - da vap 06/09). Doc file o luot sau. Ghi de file cu.
string p = "C:/Users/Admin/AppData/Local/Temp/claude/D--NRO-GoiBanGiao/8d082123-bbf0-446c-aa97-6342f0eda3a3/scratchpad/shot.png";
UnityEditorInternal.InternalEditorUtility.RepaintAllViews();
UnityEngine.ScreenCapture.CaptureScreenshot(p);
return "da yeu cau chup: " + p + " isPlaying=" + UnityEditor.EditorApplication.isPlaying + " paused=" + UnityEditor.EditorApplication.isPaused;
