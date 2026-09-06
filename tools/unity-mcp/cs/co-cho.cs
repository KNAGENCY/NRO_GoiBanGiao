// Cac co lam client ve man "Xin cho" (GameCanvas.isWaiting / paint loading) - chan doan client treo o loading.
var sb = new System.Text.StringBuilder();
try {
    sb.Append("GameCanvas.isLoading=").Append(Game1.GameCanvas.isLoading);
    sb.Append(" Char.isLoadingMap=").Append(Game1.Char.isLoadingMap);
    sb.Append(" timeBreakLoading=").Append(Game1.GameCanvas.timeBreakLoading).Append(" now=").Append(Game1.mSystem.currentTimeMillis());
    sb.Append(" LoginScr.isContinueToLogin=").Append(Game1.LoginScr.isContinueToLogin);
    sb.Append(" ServerList.waitToLogin=").Append(Game1.ServerListScreen.waitToLogin).Append(" isWait=").Append(Game1.ServerListScreen.isWait);
    sb.Append(" InfoDlg.isShow=").Append(Game1.InfoDlg.isShow);
    sb.Append(" msgdlg=").Append(Game1.GameCanvas.msgdlg == null ? "null" : (Game1.GameCanvas.msgdlg.info == null ? "?" : string.Join("|", Game1.GameCanvas.msgdlg.info)));
    sb.Append(" screen=").Append(Game1.GameCanvas.currentScreen == null ? "null" : Game1.GameCanvas.currentScreen.GetType().Name);
    sb.Append(" currentDialog=").Append(Game1.GameCanvas.currentDialog == null ? "null" : Game1.GameCanvas.currentDialog.GetType().Name);
    var me = Game1.Char.myCharz();
    sb.Append(" me hp=").Append(me.cHP).Append(" statusMe=").Append(me.statusMe).Append(" cf=").Append(me.cf).Append(" isDie=").Append(me.isDie);
    sb.Append(" tileMap=").Append(Game1.TileMap.mapName).Append(" tmw=").Append(Game1.TileMap.tmw).Append(" maps=").Append(Game1.TileMap.maps == null ? "null" : Game1.TileMap.maps.Length.ToString());
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
