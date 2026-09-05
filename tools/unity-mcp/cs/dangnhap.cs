// Tu man ServerListScreen -> mo LoginScr, dien tai khoan naruto/1 vao RMS (doLogin doc tu RMS), roi doLogin().
var sb = new System.Text.StringBuilder();
try {
    Game1.Rms.saveRMSString("acc", "naruto");
    Game1.Rms.saveRMSString("pass", "1");
    if (Game1.GameCanvas.loginScr == null) Game1.GameCanvas.loginScr = new Game1.LoginScr();
    Game1.GameCanvas.loginScr.switchToMe();
    Game1.GameCanvas.loginScr.tfUser.setText("naruto");
    Game1.GameCanvas.loginScr.tfPass.setText("1");
    Game1.GameCanvas.loginScr.doLogin();
    sb.Append("da goi doLogin; connected=").Append(Game1.Session_ME.gI().isConnected());
    sb.Append(" screen=").Append(Game1.GameCanvas.currentScreen == null ? "null" : Game1.GameCanvas.currentScreen.GetType().Name);
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message).Append(" ").Append(e.StackTrace); }
return sb.ToString();
