package npc.npc_manifest;

/**
 *
 * @author GiaCatLuong
 */

public class OngMoori extends OngGohan {

    public OngMoori(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

}

// package npc.npc_manifest;

// import consts.ConstNpc;
// import consts.ConstTask;
// import item.Item;
// import java.util.ArrayList;
// import java.util.List;
// import jdbc.daos.PlayerDAO;
// import npc.Npc;
// import player.Player;
// import services.InventoryService;
// import services.ItemService;
// import services.NpcService;
// import services.Service;
// import services.TaskService;
// import services.func.Input;
// import shop.ShopService;
// import utils.Util;

// public class OngGohan extends Npc {

//     public OngGohan(int mapId, int status, int cx, int cy, int tempId, int avartar) {
//         super(mapId, status, cx, cy, tempId, avartar);
//     }

//     int costNapVang = 1;
//     int[][] napVang = { { 20000, 30 }, { 50000, 100 }, { 100000, 250 }, { 200000, 600 }, { 500000, 1750 },
//             { 1000000, 4000 } };

//     @Override
//     public void openBaseMenu(Player player) {
//         if (canOpenNpc(player)) {
//             if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
//                 this.createOtherMenu(player, ConstNpc.BASE_MENU,
//                         "Con hãy cố gắng theo đuổi cô ấy,\nĐừng lo lắng cho ta,\nCô ấy mới đem lại hạnh phúc cho con,\nHàng tháng gửi tiền cho ta là được ?",
//                         "Gift-Code",
//                         "Quy Đổi",
//                         // "Nhận 100K\nNgọc xanh",
                     
//                         "Hộp Thư",
//                          "Mở\nThành Viên"
//                      //   ,
//                         // "Refuse"
//                     );
//             }
//         }
//     }

//     @Override
//     public void confirmMenu(Player player, int select) {
//         if (canOpenNpc(player)) {
//             if (player.iDMark.isBaseMenu()) {
//                 switch (select) {
//                     case 0 -> Input.gI().createFormGiftCode(player);
//                     case 1 -> {
//                         String npcSay = "Số Dư Của Con Là: " + Util.mumberToLouis(player.getSession().cash)
//                                 + " VND Dùng Để Nạp Qua Đơn Vị Khác\n"
//                                 + "Ta đang giữ giúp con " + Util.mumberToLouis(player.getSession().goldBar)
//                                 + " Thỏi Vàng";
//                         createOtherMenu(player, ConstNpc.NAP_TIEN, npcSay,
//                                 "Quy Đổi\nThỏi Vàng",
//                                 "Nhận\nThỏi Vàng",
//                                 "Close");
//                     }
//                     // case 2 -> {
//                     //     if (player.inventory.gem >= 500000) {
//                     //         Service.gI().sendThongBao(player, "Hãy Dùng Bớt Ngọc Xanh!");
//                     //         return;
//                     //     }
//                     //     int soLuongNgocXanh = 500_000;
//                     //     player.inventory.gem += soLuongNgocXanh;
//                     //     Service.gI().sendMoney(player);
//                     //     Service.gI().sendThongBao(player,
//                     //             "Bạn Vừa Nhận " + Util.mumberToLouis(soLuongNgocXanh) + " Ngọc Xanh");
//                     // }
//                     // case 3 -> {
//                     //     if (TaskService.gI().getIdTask(player) >= ConstTask.TASK_9_0
//                     //             && TaskService.gI().getIdTask(player) < ConstTask.TASK_11_0) {
//                     //         player.playerTask.taskMain.id = 10;
//                     //         player.playerTask.taskMain.index = 0;
//                     //         TaskService.gI().sendNextTaskMain(player);
//                     //         Service.gI().sendThongBao(player, "Bạn đã được hỗ trợ nhiệm vụ thành công");
//                     //     // } else if (TaskService.gI().getIdTask(player) >= ConstTask.TASK_16_0
//                     //     //         && TaskService.gI().getIdTask(player) < ConstTask.TASK_17_0) {
//                     //     //     player.playerTask.taskMain.id = 17;
//                     //     //     player.playerTask.taskMain.index = 0;
//                     //     //     TaskService.gI().sendNextTaskMain(player);
//                     //     //    Service.gI().sendThongBao(player, "Bạn đã được hỗ trợ nhiệm vụ thành công");
//                     //     } else if (TaskService.gI().getIdTask(player) >= ConstTask.TASK_18_0
//                     //             && TaskService.gI().getIdTask(player) < ConstTask.TASK_20_0) {
//                     //         player.playerTask.taskMain.id = 19;
//                     //         player.playerTask.taskMain.index = 0;
//                     //         TaskService.gI().sendNextTaskMain(player);
//                     //         Service.gI().sendThongBao(player, "Bạn đã được hỗ trợ nhiệm vụ thành công");
//                     //     } else {
//                     //         Service.gI().sendThongBao(player,
//                     //                 "Chỉ hỗ trợ nhiệm vụ Tàu pảy pảy, DHVT và Trung úy trắng");
//                     //     }
//                     // }
//                     case 2 -> this.createOtherMenu(player, ConstNpc.MAIL_BOX,
//                             "|0|Website: Chúc bạn chơi game vui vẻ\n",
//                             "Hòm Thư\n(" + (player.inventory.itemsMailBox.size()
//                                     - InventoryService.gI().getCountEmptyListItem(player.inventory.itemsMailBox))
//                                     + " món)",
//                             "Xóa Hết\nHòm Thư", "Close");
//                          case 3 -> { // Xử lý Mở Thành Viên
//                     // 1. Kiểm tra xem đã là thành viên chưa (dựa trên biến actived trong Session)
//                     if (player.getSession().actived) {
//                         Service.gI().sendThongBao(player, "Con đã là thành viên rồi, không cần mở nữa!");
//                         return;
//                     }

//                     // 2. Kiểm tra số dư cash (10.000)
//                     if (player.getSession().cash >= 10000) {
//                         // 3. Thực hiện trừ tiền trong DB và Session
//                         if (PlayerDAO.subcash(player, 10000)) {
//                             // 4. Gọi hàm updateActive để chuyển active từ 0 -> 1 trong DB
//                             if (PlayerDAO.updateActive(player, 1)) {
//                                 Service.gI().sendThongBao(player, "Chúc mừng con đã trở thành thành viên!");
//                             } else {
//                                 Service.gI().sendThongBao(player, "Có lỗi xảy ra khi kích hoạt, vui lòng báo Admin.");
//                             }
//                         } else {
//                             Service.gI().sendThongBao(player, "Lỗi trừ tiền, vui lòng thử lại.");
//                         }
//                     } else {
//                         Service.gI().sendThongBao(player, "Cần 10.000 VNĐ để Mở thành viên.");
//                     }
//                 }
//                     default -> Service.gI().sendThongBao(player, "Lựa chọn không hợp lệ");
//                 }
//             } else if (player.iDMark.getIndexMenu() == ConstNpc.MAIL_BOX) {
//                 switch (select) {
//                     case 0 -> ShopService.gI().opendShop(player, "ITEMS_MAIL_BOX", true);
//                     case 1 -> NpcService.gI().createMenuConMeo(player,
//                             ConstNpc.CONFIRM_REMOVE_ALL_ITEM_MAIL_BOX, this.avartar,
//                             "|3|Bạn chắc muốn xóa hết vật phẩm trong hòm thư?\n"
//                                     + "|7|Sau khi xóa sẽ không thể khôi phục!",
//                             "Đồng ý", "Hủy bỏ");
//                 }
//             } else if (player.iDMark.getIndexMenu() == ConstNpc.NAP_TIEN) {
//                 switch (select) {
//                     case 0 -> {
//                         List<String> menu = new ArrayList<>();
//                         for (int i = 0; i < napVang.length; i++) {
//                             menu.add(i, Util.mumberToLouis(napVang[i][0]) + "\n"
//                                     + Util.mumberToLouis(napVang[i][1] * costNapVang) + " Thỏi\nvàng");
//                         }
//                         String[] menus = menu.toArray(String[]::new);
//                         createOtherMenu(player, ConstNpc.NAP_VANG, "Ta sẽ giữ giúp con\n"
//                                 + "Nếu con cần dùng tới hãy quay lại đây gặp ta!", menus);
//                     }
//                     case 1 -> {
//                         if (player.getSession().goldBar > 0) {
//                             List<Item> listItem = new ArrayList<>();
//                             Item thoiVang = ItemService.gI().createNewItem((short) 457, player.getSession().goldBar);
//                             listItem.add(thoiVang);
//                             if (InventoryService.gI().getCountEmptyBag(player) < listItem.size()) {
//                                 Service.gI().sendThongBao(player,
//                                         "Cần ít nhất " + listItem.size() + " ô trống trong hành trang");
//                             }
//                             for (Item it : listItem) {
//                                 InventoryService.gI().addItemBag(player, it);
//                                 InventoryService.gI().sendItemBag(player);
//                             }
//                             Service.gI().sendThongBao(player,
//                                     "Bạn đã nhận được " + player.getSession().goldBar + " thỏi vàng");
//                             PlayerDAO.subGoldBar(player, -(napVang[select][1] * costNapVang));
//                             PlayerDAO.subGoldBar(player, player.getSession().goldBar);
//                         }
//                     }
//                 }
//             } else if (player.iDMark.getIndexMenu() == ConstNpc.DOIXUOCEAN) {
//                 switch (select) {
//                     case 0 -> exchangeXukrai(player, 10000, 110);
//                     case 1 -> exchangeXukrai(player, 20000, 250);
//                     case 2 -> exchangeXukrai(player, 30000, 550);
//                     case 3 -> exchangeXukrai(player, 50000, 1100);
//                     case 4 -> exchangeXukrai(player, 100000, 2550);
//                     case 5 -> exchangeXukrai(player, 200000, 5200);
//                     case 6 -> exchangeXukrai(player, 500000, 12000);
//                     case 7 -> exchangeXukrai(player, 1000000, 25000);
//                     default -> Service.gI().sendThongBao(player, "Lựa chọn không hợp lệ");
//                 }
//             } else if (player.iDMark.getIndexMenu() == ConstNpc.NAP_VANG) {
//                 if (player.getSession().cash >= napVang[select][0]) {
//                     List<Item> listItem = new ArrayList<>();
//                     if (InventoryService.gI().getCountEmptyBag(player) < listItem.size()) {
//                         Service.gI().sendThongBao(player,
//                                 "Cần ít nhất " + listItem.size() + " ô trống trong hành trang");
//                     }
//                     for (Item it : listItem) {
//                         InventoryService.gI().addItemBag(player, it);
//                         InventoryService.gI().sendItemBag(player);
//                     }
//                     PlayerDAO.subcash(player, napVang[select][0]);
//                     PlayerDAO.subGoldBar(player, -(napVang[select][1] * costNapVang));
//                     Service.gI().sendThongBao(player,
//                             "Bạn có thêm " + Util.mumberToLouis(napVang[select][1] * costNapVang) + " thỏi vàng");
//                 } else {
//                     Service.gI().sendThongBao(player, "Không đủ số dư");
//                 }
//             }
//         }
//     }

//     private void exchangeXukrai(Player player, int i, int i0) {
//         if (player.getSession().cash >= i) {
//             PlayerDAO.subcash(player, i);
//             Item xu = ItemService.gI().createNewItem((short) 1705, i0);
//             InventoryService.gI().addItemBag(player, xu);
//             Service.gI().sendMoney(player);
//             InventoryService.gI().sendItemBag(player);
//             Service.gI().sendThongBao(player, "Bạn có vừa có thêm " + i0 + "xu");
//         } else {
//             Service.gI().sendThongBao(player, "Bạn không đủ số dư");
//         }
//     }
// }