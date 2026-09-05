// Doc vi tri quai/NPC/nguoi choi tren client (chi doc) de kiem "lo lung 2 o".
var sb = new System.Text.StringBuilder();
try {
    var me = Game1.Char.myCharz();
    sb.Append("me cx=").Append(me.cx).Append(" cy=").Append(me.cy);
    sb.Append(" | cmx=").Append(Game1.GameScr.cmx).Append(" cmy=").Append(Game1.GameScr.cmy);
    sb.Append(" | mobs=");
    for (int i = 0; i < Game1.GameScr.vMob.size(); i++) {
        var m = (Game1.Mob)Game1.GameScr.vMob.elementAt(i);
        sb.Append("[id=").Append(m.mobId).Append(" tpl=").Append(m.templateId).Append(" x=").Append(m.x).Append(" y=").Append(m.y).Append(" yFirst=").Append(m.yFirst).Append(" hp=").Append(m.hp).Append("] ");
    }
    sb.Append(" | npcs=");
    for (int i = 0; i < Game1.GameScr.vNpc.size(); i++) {
        var n = (Game1.Npc)Game1.GameScr.vNpc.elementAt(i);
        sb.Append("[tpl=").Append(n.template == null ? -1 : n.template.npcTemplateId).Append(" cx=").Append(n.cx).Append(" cy=").Append(n.cy).Append("] ");
    }
    sb.Append(" | chars=");
    for (int i = 0; i < Game1.GameScr.vCharInMap.size(); i++) {
        var c = (Game1.Char)Game1.GameScr.vCharInMap.elementAt(i);
        sb.Append("[").Append(c.cName).Append(" id=").Append(c.charID).Append(" cx=").Append(c.cx).Append(" cy=").Append(c.cy).Append(" hp=").Append(c.cHP).Append("] ");
    }
    sb.Append(" | map tmw=").Append(Game1.TileMap.tmw).Append(" tmh=").Append(Game1.TileMap.tmh).Append(" size=").Append(Game1.TileMap.size);
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
