// // package npc.npc_manifest;

// // /**
// //  *
// //  * @author GiaCatLuong
// //  */

// // import clan.Clan;
// // import clan.ClanMember;
// // import consts.ConstNpc;
// // import models.DestronGas.DestronGas;
// // import models.DestronGas.DestronGasService;
// // import npc.Npc;
// // import static npc.NpcFactory.PLAYERID_OBJECT;
// // import player.Player;
// // import services.NpcService;
// // import services.func.Input;
// // import utils.TimeUtil;

// // public class MrPoPo extends Npc {

// //     public MrPoPo(int mapId, int status, int cx, int cy, int tempId, int avartar) {
// //         super(mapId, status, cx, cy, tempId, avartar);
// //     }

// //     @Override
// //     public void openBaseMenu(Player player) {
// //         if (canOpenNpc(player)) {
// //             if (this.mapId == 0) {
// //                 if (player.clan != null) {
// //                     this.createOtherMenu(player, ConstNpc.BASE_MENU,
// //                             "Trong Này Chứa Cải Trang Mạnh Nhất Game\nTrong Này Chứa Cải Trang Mạnh Nhất Game\nTrong Này Chứa Cải Trang Mạnh Nhất Game\nTa sẽ đưa các cậu đến nơi ấy, các cậu đã sẵn sàng chưa?",
// //                             "Thông tin\nChi tiết", "Top 100\nBang hội", "Thành tích\nBang", "OK", "Từ chối");
// //                 } else {
// //                     this.createOtherMenu(player, ConstNpc.BASE_MENU,
// //                             "Trong Này Chứa Cải Trang Mạnh Nhất Game\nTrong Này Chứa Cải Trang Mạnh Nhất Game\nTrong Này Chứa Cải Trang Mạnh Nhất Game\nTa sẽ đưa các cậu đến nơi ấy, các cậu đã sẵn sàng chưa?",
// //                             "Thông tin\nChi tiết", "Top 100\nBang hội", "OK", "Từ chối");
// //                 }
// //             }
// //         }
// //     }

// //     @Override
// //     public void confirmMenu(Player player, int select) {
// //         if (canOpenNpc(player)) {
// //             if (this.mapId == 0) {
// //                 if (player.iDMark.isBaseMenu()) {
// //                     switch (select) {
// //                         case 0 ->
// //                             NpcService.gI().createTutorial(player, tempId, this.avartar,
// //                                     ConstNpc.HUONG_DAN_KHI_GAS_HUY_DIET);
// //                         case 1 -> {
// //                         }
// //                         case 2 -> {
// //                         }
// //                         case 3 -> {
// //                             Clan clan = player.clan;
// //                             if (clan != null) {
// //                                 ClanMember cm = clan.getClanMember((int) player.id);
// //                                 if (cm != null) {
// //                                     // if (player.clanMember.getNumDateFromJoinTimeToToday() < 1) {
// //                                     // NpcService.gI().createTutorial(player, tempId, this.avartar,
// //                                     // "Gia nhập bang hội trên 2 ngày mới được tham gia");
// //                                     // return;
// //                                     // }
// //                                     if (player.clan.KhiGasHuyDiet != null) {
// //                                         createOtherMenu(player, 2,
// //                                                 "Bang hội của cậu đang tham gia Destron Gas cấp độ "
// //                                                         + player.clan.KhiGasHuyDiet.level + "\n"
// //                                                         + "cậu có muốn đi cùng họ không ? ("
// //                                                         + TimeUtil.convertTimeNow(
// //                                                                 player.clan.KhiGasHuyDiet.getLastTimeOpen())
// //                                                         + " trước)",
// //                                                 "Đồng ý", "Từ chối");
// //                                         return;
// //                                     }
// //                                     if (!clan.isLeader(player)) {
// //                                         NpcService.gI().createTutorial(player, tempId, this.avartar,
// //                                                 "Chức năng chỉ dành cho bang chủ");
// //                                         return;
// //                                     }
// //                                     if (clan.members.size() < DestronGas.N_PLAYER_CLAN) {
// //                                         NpcService.gI().createTutorial(player, tempId, this.avartar,
// //                                                 "Bang hội phải có ít nhất 5 thành viên mới có thể tham gia");
// //                                         return;
// //                                     }
// //                                     Input.gI().createFormChooseLevelKGHD(player);
// //                                 }
// //                             }
// //                         }
// //                     }
// //                 } else if (player.iDMark.getIndexMenu() == 2) {
// //                     if (select == 0) {
// //                         if (player.clan.KhiGasHuyDiet == null) {
// //                             DestronGasService.gI().openKhiGasHuyDiet(player,
// //                                     Byte.parseByte(String.valueOf(PLAYERID_OBJECT.get(player.id))));
// //                         } else {
// //                             DestronGasService.gI().openKhiGasHuyDiet(player, (byte) 0);
// //                         }
// //                     }
// //                 }
// //             }
// //         }
// //     }
// // }

// package npc.npc_manifest;

// /**
//  * NPC Nhận quà nhiệm vụ + sức mạnh + bán đồ thần linh
//  * Logic lưu JSON giống GiuMaDauBo
//  * @author GiaCatLuong
//  */

// import consts.ConstNpc;
// import item.Item;
// import jdbc.DBConnecter;
// import npc.Npc;
// import player.Player;
// import services.InventoryService;
// import services.ItemService;
// import services.Service;

// import java.sql.Connection;
// import java.sql.PreparedStatement;
// import java.sql.ResultSet;

// import org.json.simple.JSONObject;
// import org.json.simple.JSONValue;

// public class MrPoPo extends Npc {

//     private static final int MENU_NHIEM_VU = 1001;
//     private static final int MENU_SUC_MANH = 1002;

//     // Mốc nhiệm vụ
//     private static final int[] MOC_NHIEM_VU = {21, 22, 23, 28};

//     // Mốc sức mạnh (tỷ)
//     private static final int[] MOC_SUC_MANH = {11, 22, 33, 44};

//     public MrPoPo(int mapId, int status, int cx, int cy, int tempId, int avartar) {
//         super(mapId, status, cx, cy, tempId, avartar);
//     }
// @Override
// public void openBaseMenu(Player player) {
//     if (!canOpenNpc(player)) return;

//     // CHỈ CHO PHÉP Ở MAP 0
//     if (this.mapId != 0) {
//         Service.gI().sendThongBao(player,
//                 "NPC này chỉ hoạt động tại Map 0");
//         return;
//     }

//     createOtherMenu(player, ConstNpc.BASE_MENU,
//             "1. Kẻ thù lớn nhất của đời người là chính mình\r\n" + //
                                
//                                 "2. Thất bại lớn nhất của đời người là tự đại\r\n" + //
                       
//                                 "3. Ngu dốt lớn nhất của đời người là dối trá\r\n" + //
                    
//                                 "4. Bi ai lớn nhất của đời người là ghen tị\r\n" + //
                    
//                                 "5. Sai lầm lớn nhất của đời người là đánh mất mình\r\n" + //
           
//                                 "6. Tội lỗi lớn nhất của đời người là bất hiếu\r\n" + //
               
//                                 "7. Đáng thương lớn nhất của đời người là tự ti\r\n" + //
              
//                                 "8. Đáng khâm phục nhất của đời người là vươn lên sau khi vấp ngã\r\n" + //                         
// "...\r\n" +
//                                 "14. An ủi lớn nhất của đời người là bố thí.",
//             "Nhận quà\nNhiệm Vụ",
//             "Nhận quà\nSức Mạnh",
//             "Bán đồ\nThần Kinh",
//             "Đóng");
// }

//     /* ================= XỬ LÝ MENU ================= */
//     @Override
//     public void confirmMenu(Player player, int select) {
//         if (!canOpenNpc(player)) return;

//         if (player.iDMark.isBaseMenu()) {
//             switch (select) {
//                 case 0 -> {
//                     createOtherMenu(player, ConstNpc.BASE_MENU,
//                             "Nhận quà nhiệm vụ\n(Mỗi mốc nhận 1 lần)",
//                             "Xong Nappa",
//                             "Xong KuKu",
//                             "Xong TĐST",
//                             "Xong Xên 8",
//                             "Đóng");
//                     player.iDMark.setIndexMenu(MENU_NHIEM_VU);
//                 }
//                 case 1 -> {
//                     createOtherMenu(player, ConstNpc.BASE_MENU,
//                             "Nhận quà sức mạnh\n(Mỗi mốc nhận 1 lần)",
//                             "11 Tỷ",
//                             "22 Tỷ",
//                             "33 Tỷ",
//                             "44 Tỷ",
//                             "Đóng");
//                     player.iDMark.setIndexMenu(MENU_SUC_MANH);
//                 }
//                 case 2 -> Service.gI().sendThongBaoOK(player,
//                         "Liên hệ Admin để bán lại đồ Thần Kinh\n"
//                                 + "1 Món = 20K ATM\n"
//                                 + "Hoặc 50K Game");
//             }
//             return;
//         }

//         if (player.iDMark.getIndexMenu() == MENU_NHIEM_VU) {
//             if (select >= 0 && select < MOC_NHIEM_VU.length) {
//                 nhanQuaNhiemVu(player, MOC_NHIEM_VU[select]);
//             }
//             return;
//         }

//         if (player.iDMark.getIndexMenu() == MENU_SUC_MANH) {
//             if (select >= 0 && select < MOC_SUC_MANH.length) {
//                 nhanQuaSucManh(player, MOC_SUC_MANH[select]);
//             }
//         }
//     }

//     /* ================= NHẬN QUÀ NHIỆM VỤ ================= */
//     private void nhanQuaNhiemVu(Player player, int moc) {
//         try (Connection con = DBConnecter.getConnectionServer()) {

//             PreparedStatement ps = con.prepareStatement(
//                     "SELECT reward_moc FROM player WHERE id = ?");
//             ps.setLong(1, player.id);

//             ResultSet rs = ps.executeQuery();
//             if (!rs.next()) return;

//             JSONObject json = rs.getString("reward_moc") == null
//                     ? new JSONObject()
//                     : (JSONObject) JSONValue.parse(rs.getString("reward_moc"));

//             String key = "task_" + moc;

//             if (player.playerTask.taskMain.id < moc) {
//                 Service.gI().sendThongBao(player,
//                         "Bạn chưa đến nhiệm vụ " + moc);
//                 return;
//             }

//             if (json.containsKey(key)) {
//                 Service.gI().sendThongBao(player,
//                         "Bạn đã nhận quà mốc này rồi");
//                 return;
//             }

//             if (InventoryService.gI().getCountEmptyBag(player) < 1) {
//                 Service.gI().sendThongBao(player,
//                         "Hành trang không đủ chỗ trống");
//                 return;
//             }

//             // TRAO QUÀ (TẠM)
//             giveItem(player, 1, 1);

//             json.put(key, true);
//             updateJson(con, json, player.id);

//             Service.gI().sendThongBao(player,
//                     "Nhận thành công quà nhiệm vụ " + moc);

//         } catch (Exception e) {
//             e.printStackTrace();
//         }
//     }

//     /* ================= NHẬN QUÀ SỨC MẠNH ================= */
//     private void nhanQuaSucManh(Player player, int mocTy) {
//         try (Connection con = DBConnecter.getConnectionServer()) {

//             PreparedStatement ps = con.prepareStatement(
//                     "SELECT reward_moc FROM player WHERE id = ?");
//             ps.setLong(1, player.id);

//             ResultSet rs = ps.executeQuery();
//             if (!rs.next()) return;

//             JSONObject json = rs.getString("reward_moc") == null
//                     ? new JSONObject()
//                     : (JSONObject) JSONValue.parse(rs.getString("reward_moc"));

//             String key = "power_" + mocTy;
//             long powerNeed = (long) mocTy * 1_000_000_000L;

//             if (player.nPoint.power < powerNeed) {
//                 Service.gI().sendThongBao(player,
//                         "Bạn chưa đạt " + mocTy + " tỷ sức mạnh");
//                 return;
//             }

//             if (json.containsKey(key)) {
//                 Service.gI().sendThongBao(player,
//                         "Bạn đã nhận quà mốc này rồi");
//                 return;
//             }

//             if (InventoryService.gI().getCountEmptyBag(player) < 1) {
//                 Service.gI().sendThongBao(player,
//                         "Hành trang không đủ chỗ trống");
//                 return;
//             }

//             // TRAO QUÀ (TẠM)
//             giveItem(player, 1, 1);

//             json.put(key, true);
//             updateJson(con, json, player.id);

//             Service.gI().sendThongBao(player,
//                     "Nhận thành công quà " + mocTy + " tỷ sức mạnh");

//         } catch (Exception e) {
//             e.printStackTrace();
//         }
//     }

//     /* ================= GIVE ITEM ================= */
//     private void giveItem(Player player, int tempId, int quantity) {
//         Item item = ItemService.gI().createNewItem((short) tempId, quantity);
//         item.itemOptions.add(new Item.ItemOption(30, 1));
//         InventoryService.gI().addItemBag(player, item);
//         InventoryService.gI().sendItemBag(player);
//     }

//     /* ================= UPDATE JSON ================= */
//     private void updateJson(Connection con, JSONObject json, long playerId) throws Exception {
//         PreparedStatement ps = con.prepareStatement(
//                 "UPDATE player SET reward_moc = ? WHERE id = ?");
//         ps.setString(1, json.toJSONString());
//         ps.setLong(2, playerId);
//         ps.executeUpdate();
//     }
// }

