var sb = new System.Text.StringBuilder();
try {
    var me = Game1.Char.myCharz();
    sb.Append("screen=").Append(Game1.GameCanvas.currentScreen == null ? "null" : Game1.GameCanvas.currentScreen.GetType().Name);
    sb.Append(" id=").Append(me.charID).Append(" name=").Append(me.cName);
    sb.Append(" gender=").Append(me.cgender).Append(" head=").Append(me.head).Append(" body=").Append(me.body).Append(" leg=").Append(me.leg);
    sb.Append(" hp=").Append(me.cHP).Append("/").Append(me.cHPFull).Append(" cx=").Append(me.cx).Append(" cy=").Append(me.cy).Append(" cf=").Append(me.cf);
    sb.Append(" parts=").Append(Game1.GameScr.parts == null ? -1 : Game1.GameScr.parts.Length);
    if (Game1.GameScr.parts != null && Game1.GameScr.parts.Length > 58) {
        var ph = Game1.GameScr.parts[me.head]; var pb = Game1.GameScr.parts[me.body]; var pl = Game1.GameScr.parts[me.leg];
        sb.Append(" ph=").Append(ph == null ? "null" : ph.type + ":" + ph.pi[0].id + "," + ph.pi[0].dy);
        sb.Append(" pb=").Append(pb == null ? "null" : pb.type + ":" + pb.pi[0].id + "," + pb.pi[0].dy);
        sb.Append(" pl=").Append(pl == null ? "null" : pl.type + ":" + pl.pi[0].id + "," + pl.pi[0].dy);
    }
    sb.Append(" mapName=").Append(Game1.TileMap.mapName).Append(" tmw=").Append(Game1.TileMap.tmw).Append(" tmh=").Append(Game1.TileMap.tmh);
    sb.Append(" vChar=").Append(Game1.GameScr.vCharInMap == null ? -1 : Game1.GameScr.vCharInMap.size());
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message); }
return sb.ToString();
