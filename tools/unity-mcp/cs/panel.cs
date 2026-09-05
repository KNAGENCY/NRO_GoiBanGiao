// Panel cua hang co mo khong? (sau khi confirmMenu(0,0) -> -44)
var sb = new System.Text.StringBuilder();
try {
    var p = Game1.GameCanvas.panel;
    sb.Append("panel.isShow=").Append(p.isShow).Append(" type=").Append(p.type);
    sb.Append(" currentTab=").Append(p.currentTabIndex);
    try { sb.Append(" tabNames=").Append(p.tabName != null && p.tabName[p.type] != null ? string.Join("|", p.tabName[p.type][0]) : "?"); } catch (System.Exception) { sb.Append(" tabNames=?"); }
    sb.Append(" menu.showMenu=").Append(Game1.GameCanvas.menu.showMenu);
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
