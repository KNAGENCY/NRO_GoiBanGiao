var sb = new System.Text.StringBuilder();
try {
    sb.Append("isPlaying=").Append(UnityEditor.EditorApplication.isPlaying);
    sb.Append(" connected=").Append(Game1.Session_ME.gI().isConnected());
    sb.Append(" screen=").Append(Game1.GameCanvas.currentScreen == null ? "null" : Game1.GameCanvas.currentScreen.GetType().Name);
    sb.Append(" dlg=").Append(Game1.GameCanvas.currentDialog == null ? "null" : Game1.GameCanvas.currentDialog.GetType().Name);
    sb.Append(" loginScrNull=").Append(Game1.GameCanvas.loginScr == null);
    var me = Game1.Char.myCharz();
    sb.Append(" me=").Append(me == null ? "null" : (me.cName + " hp=" + me.cHP + "/" + me.cHPFull + " head=" + me.head + " body=" + me.body + " leg=" + me.leg));
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
