// package npc.npc_manifest;

// /**
//  * NPC Tổ sư Kaio
//  * Chỉ dùng để xem BẢNG XẾP HẠNG
//  *
//  * @author GiaCatLuong
//  */

// import consts.ConstNpc;
// import jdbc.DBConnecter;
// import npc.Npc;
// import player.Player;
// import services.Service;
// import utils.Util;

// import java.sql.Connection;
// import java.sql.PreparedStatement;
// import java.sql.ResultSet;

// public class ToSuKaio extends Npc {

//     private static final int MENU_TOP = 3000;

//     public ToSuKaio(int mapId, int status, int cx, int cy, int tempId, int avartar) {
//         super(mapId, status, cx, cy, tempId, avartar);
//     }

//     @Override
//     public void openBaseMenu(Player player) {
//         if (!canOpenNpc(player)) return;

//         createOtherMenu(player, ConstNpc.BASE_MENU,
//                 "TỔ SƯ KAIO\n\nBảng xếp hạng cao thủ toàn server",
//                 "Top\nSức Đánh",
//                 "Top\nHP",
//                 "Top\nKI",
//                 "Đóng");
//     }

//     @Override
//     public void confirmMenu(Player player, int select) {
//         if (!canOpenNpc(player)) return;

//         if (player.iDMark.isBaseMenu()) {
//             switch (select) {
//                 case 0 -> showTop(player, "dame_point_fusion", "TOP 10 SỨC ĐÁNH");
//                 case 1 -> showTop(player, "hp_point_fusion", "TOP 10 HP");
//                 case 2 -> showTop(player, "mp_point_fusion", "TOP 10 KI");
//             }
//         }
//     }

//     /* ================= BẢNG XẾP HẠNG ================= */
//     private void showTop(Player player, String column, String title) {
//         StringBuilder sb = new StringBuilder("|7|" + title + "\n");

//         try (Connection con = DBConnecter.getConnectionServer()) {
//             PreparedStatement ps = con.prepareStatement(
//                     "SELECT name, " + column +
//                     " FROM player ORDER BY " + column + " DESC LIMIT 10");
//             ResultSet rs = ps.executeQuery();

//             int rank = 1;
//             while (rs.next()) {
//                 sb.append("|0|")
//                         .append(rank++).append(". ")
//                         .append(rs.getString("name"))
//                         .append(": |6|")
//                         .append(Util.numberToMoney(rs.getLong(column)))
//                         .append("\n");
//             }
//         } catch (Exception e) {
//             e.printStackTrace();
//             sb.append("|7|Không thể tải dữ liệu BXH");
//         }

//         Service.gI().sendThongBaoOK(player, sb.toString());
//     }
// }
