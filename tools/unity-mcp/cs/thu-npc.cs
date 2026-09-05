// Tai hien tay hai buoc cua loadInfoMap de xem buoc nao nem ngoai le: tao Mob thu 2 va tao Npc Iruka.
var sb = new System.Text.StringBuilder();
try {
    var mob = new Game1.Mob(1, false, false, false, false, false, 0, 0, 100, (sbyte)1, 100, (short)40, (short)384, (sbyte)0, (sbyte)1);
    sb.Append("Mob(1) OK x=").Append(mob.x).Append(" tplType=").Append(Game1.Mob.arrMobTemplate[0].type);
} catch (System.Exception e) { sb.Append(" MOB EX=").Append(e.GetType().Name).Append(":").Append(e.Message).Append(" @ ").Append(e.StackTrace == null ? "" : e.StackTrace.Split('\n')[0]); }
try {
    var npc = new Game1.Npc(0, 0, 172, 387, 0, 27078);
    sb.Append(" | Npc OK name=").Append(npc.template == null ? "null" : npc.template.name);
    Game1.GameScr.vNpc.addElement(npc);
    sb.Append(" vNpc=").Append(Game1.GameScr.vNpc.size());
} catch (System.Exception e) { sb.Append(" NPC EX=").Append(e.GetType().Name).Append(":").Append(e.Message).Append(" @ ").Append(e.StackTrace == null ? "" : e.StackTrace.Split('\n')[0]); }
sb.Append(" | vMob=").Append(Game1.GameScr.vMob.size());
for (int i = 0; i < Game1.GameScr.vMob.size(); i++) { var m = (Game1.Mob)Game1.GameScr.vMob.elementAt(i); sb.Append(" [mob id=").Append(m.mobId).Append(" x=").Append(m.x).Append(" hp=").Append(m.hp).Append(" isDie=").Append(m.isDie).Append("]"); }
return sb.ToString();
