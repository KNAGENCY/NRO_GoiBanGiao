// Doc part Iruka (2189/2190/2191) client dang nap + NPC trong map — kiem client da nhan part.bin moi (dx/dy) chua.
var sb = new System.Text.StringBuilder();
try {
    var parts = Game1.GameScr.parts;
    sb.Append("parts=").Append(parts == null ? -1 : parts.Length);
    int[] ids = new int[] { 30, 57, 58, 2189, 2190, 2191 };
    foreach (int id in ids) {
        if (parts == null || id >= parts.Length || parts[id] == null) { sb.Append(" [").Append(id).Append("=null]"); continue; }
        var p = parts[id];
        sb.Append(" [").Append(id).Append(" type=").Append(p.type).Append(" n=").Append(p.pi.Length);
        for (int k = 0; k < p.pi.Length && k < 3; k++) sb.Append(" (").Append(p.pi[k].id).Append(",").Append(p.pi[k].dx).Append(",").Append(p.pi[k].dy).Append(")");
        sb.Append("]");
    }
    sb.Append(" | npcs=");
    for (int i = 0; i < Game1.GameScr.vNpc.size(); i++) {
        var n = (Game1.Npc)Game1.GameScr.vNpc.elementAt(i);
        sb.Append("[tpl=").Append(n.template == null ? -1 : n.template.npcTemplateId)
          .Append(" head=").Append(n.template == null ? -1 : n.template.headId)
          .Append(" body=").Append(n.template == null ? -1 : n.template.bodyId)
          .Append(" leg=").Append(n.template == null ? -1 : n.template.legId)
          .Append(" cx=").Append(n.cx).Append(" cy=").Append(n.cy).Append(" cf=").Append(n.cf).Append("] ");
    }
    var me = Game1.Char.myCharz();
    sb.Append(" | me cx=").Append(me.cx).Append(" cy=").Append(me.cy).Append(" cf=").Append(me.cf).Append(" cdir=").Append(me.cdir);
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
