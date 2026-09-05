// Sau openMenu: menu client co hien khong; chon dong 0 (Cua hang) roi doc panel.
var sb = new System.Text.StringBuilder();
try {
    sb.Append("menu.showMenu=").Append(Game1.GameCanvas.menu.showMenu);
    sb.Append(" chatPopup=").Append(Game1.ChatPopup.currChatPopup == null ? "null" : "co");
    Game1.Service.gI().confirmMenu(0, 0);
    sb.Append(" | da gui confirmMenu(0,0)");
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
