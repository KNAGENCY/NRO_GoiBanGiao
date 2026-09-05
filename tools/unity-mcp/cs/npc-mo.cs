// Kiem NPC hien va mo cua hang tu client that: doc vNpc, taskMaint; goi openMenu(0) roi confirmMenu(0,0).
var sb = new System.Text.StringBuilder();
try {
    var me = Game1.Char.myCharz();
    sb.Append("taskMaint=").Append(me.taskMaint == null ? "null" : ("id=" + me.taskMaint.taskId));
    sb.Append(" vNpc=").Append(Game1.GameScr.vNpc.size());
    for (int i = 0; i < Game1.GameScr.vNpc.size(); i++) {
        var n = (Game1.Npc)Game1.GameScr.vNpc.elementAt(i);
        sb.Append(" [npc ").Append(n.template == null ? "?" : n.template.name).Append(" cx=").Append(n.cx).Append(" cy=").Append(n.cy).Append("]");
        me.npcFocus = n; // client chi nhan hoi thoai 32 cho NPC dang focus
    }
    Game1.Service.gI().openMenu(0);
    sb.Append(" | da gui openMenu(0)");
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
