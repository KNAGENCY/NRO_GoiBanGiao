// Vi sao vNpc rong? Doc arrNpcTemplate, taskMaint, waypoint, bg item.
var sb = new System.Text.StringBuilder();
try {
    var arr = Game1.Npc.arrNpcTemplate;
    sb.Append("arrNpcTemplate=").Append(arr == null ? "null" : arr.Length.ToString());
    if (arr != null) for (int i = 0; i < arr.Length; i++) { var t = arr[i]; sb.Append(" [").Append(i).Append(":").Append(t == null ? "null" : (t.npcTemplateId + " " + t.name + " head=" + t.headId + " body=" + t.bodyId + " leg=" + t.legId)).Append("]"); }
    var me = Game1.Char.myCharz();
    sb.Append(" | taskMaint=").Append(me.taskMaint == null ? "null" : ("id=" + me.taskMaint.taskId + " idx=" + me.taskMaint.index));
    sb.Append(" | vNpc=").Append(Game1.GameScr.vNpc == null ? -1 : Game1.GameScr.vNpc.size());
    sb.Append(" | vGo(waypoint)=").Append(Game1.TileMap.vGo == null ? -1 : Game1.TileMap.vGo.size());
    sb.Append(" | vCurrItem=").Append(Game1.TileMap.vCurrItem == null ? -1 : Game1.TileMap.vCurrItem.size());
    sb.Append(" | mapName=").Append(Game1.TileMap.mapName).Append(" mapID=").Append(Game1.TileMap.mapID);
} catch (System.Exception e) { sb.Append(" EX=").Append(e.GetType().Name).Append(":").Append(e.Message).Append(" ").Append(e.StackTrace); }
return sb.ToString();
