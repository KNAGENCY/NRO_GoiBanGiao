package npc.npc_manifest;

import consts.ConstNpc;
import item.Item;
import java.util.ArrayList;
import java.util.List;
import jdbc.daos.PlayerDAO;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.NpcService;
import services.Service;
import services.TaskService;
import services.func.Input;
import shop.ShopService;
import utils.Util;

public class OngGohan extends Npc {

    public OngGohan(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    int costNapVang = 1;
    int[][] napVang = { 
        { 20000, 30 }, { 50000, 100 }, { 100000, 250 }, 
        { 200000, 600 }, { 500000, 1750 }, { 1000000, 4000 } 
    };

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
                this.createOtherMenu(player, ConstNpc.BASE_MENU,
                        "Con hãy cố gắng theo đuổi cô ấy,\nĐừng lo lắng cho ta,\n"
                        + "Cô ấy mới đem lại hạnh phúc cho con,\nHàng tháng gửi tiền cho ta là được ?",
                        "Gift Code",
                        "Quy Đổi",
                        "Hòm Thư",
                        "MTV Free",
                        "Giới Thiệu\nLối Chơi" // Thêm menu mới ở vị trí case 4
                );
            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) {
            return;
        }

        if (player.iDMark.isBaseMenu()) {
            switch (select) {
                case 0 -> Input.gI().createFormGiftCode(player);
                case 1 -> {
                    String npcSay = "Số Dư Của Con Là: " + Util.mumberToLouis(player.getSession().cash / 1)
                            + " Cá\nCon không nên xài chúng\nHoặc có thể đổi qua vàng\n ♥  ♥\n"
                            + "Ta đang giữ giúp con " + Util.mumberToLouis(player.getSession().goldBar)
                            + " Thỏi Vàng";
                    createOtherMenu(player, ConstNpc.NAP_TIEN, npcSay, "Quy Đổi\nThỏi Vàng", "Nhận\nThỏi Vàng", "Close");
                }
                case 2 -> this.createOtherMenu(player, ConstNpc.MAIL_BOX,
                        "|0|Website: Chúc bạn chơi game vui vẻ\n",
                        "Hòm Thư\n(" + (player.inventory.itemsMailBox.size()
                        - InventoryService.gI().getCountEmptyListItem(player.inventory.itemsMailBox))
                        + " món)", "Xóa Hết\nHòm Thư", "Close");
              case 3 -> {
    if (player.getSession().actived) {
        Service.gI().sendThongBao(player, "Con đã là thành viên rồi!");
        return;
    }

    this.createOtherMenu(player, ConstNpc.MO_THANH_VIEN,
            "Chỉ Có Thể MTV Free sau nhiệm vụ Fide,\nNạp tiền cũng không mở được ",
            
            "Mở\nFree");
}
                case 4 -> { // Mở menu Giới Thiệu Lối Chơi
                    this.createOtherMenu(player, ConstNpc.GIOI_THIEU_GAME, 
                            "Con muốn tìm hiểu về điều gì?", 
                            "Sao phê la,\nNgọc rồng", "Up TNSM", "Phụ Kiện", "SKH\nThạch Thần","Hướng Dẫn\nMAP WAR","Hướng Dẫn\nCày Vàng");
                }
                default -> Service.gI().sendThongBao(player, "Lựa chọn không hợp lệ");
            }
        } 


else if (player.iDMark.getIndexMenu() == ConstNpc.MO_THANH_VIEN) {
    switch (select) {

    
        case 0 -> {
            if (player.getSession().actived) {
                Service.gI().sendThongBao(player, "Đã là thành viên rồi!");
                return;
            }

            if (player.playerTask.taskMain.id < 24) {
                Service.gI().sendThongBao(player,
                        "Cần hoàn thành nhiệm vụ Fide!");
                return;
            }

            if (PlayerDAO.updateActive(player, 1)) {
                Service.gI().sendThongBao(player,
                        "Chúc mừng con đã mở thành viên miễn phí!");
            } else {
                Service.gI().sendThongBao(player, "Lỗi kích hoạt!");
            }
        }
    }
}


        // Xử lý menu Giới Thiệu Lối Chơi
        else if (player.iDMark.getIndexMenu() == ConstNpc.GIOI_THIEU_GAME) {
            switch (select) {
                case 0 -> Service.gI().sendThongBaoOK(player, "SPL cần mặc Cải trang Sói Héc Quyn up',\nNgọc Rồng chỉ up ở 'Núi Khỉ Đỏ', 'Rừng Đá' hoặc săn Boss Võ Đài.");
                case 1 -> Service.gI().sendThongBaoOK(player, "Dưới 150tr thì vào Ngũ Hành Sơn, \nTrên 150tr và Khỏe thì vào BĐKB ở Quy Lão.");
                case 2 -> Service.gI().sendThongBaoOK(player, "Săn Boss Ngũ Hành Sơn, \nVõ Đài Sinh Tử (Bà hạt mít), \nUp ở NPC Batoba (Núi khỉ đỏ).");
                // Trước đây câu này ghi "săn Boss có chữ AZ" - AZ là tên server cũ.
                // Nay liệt kê thẳng tên boss cho người chơi dễ tìm.
                // Danh sách này khớp với 5 lớp boss thả vật phẩm 674 (Thạch Thần):
                // BD, BrolySsj, Cooler, Cumber, S - xem thư mục src\boss\boss_manifest.
                case 3 -> Service.gI().sendThongBaoOK(player, "\nSKH Up 3 Map đầu hoặc Đảo Kame = 1 SKH cơ bản\n3 SKH VIP còn lại nâng từ Đồ Thần Linh\n\nThạch Thần rơi từ các Boss:\n- MA THẦN (Vùng Đất Săn Boss)\n- Frost King, Frost KCN, Cooler (vùng băng)\n- Cumber God, Cumber Blue, Cumber Black Rose (Hành tinh ngục tù)\n- Thế Giới, TẬN THẾ (Bản đồ kho báu)\nHoặc Quy đổi khi Nạp.");
                 case 4 -> Service.gI().sendThongBaoOK(player, "Tất cả đều không sử dụng Trang Bị cũ,\nPhải tự tìm Trang Bị mới,\nSau đó Đại Chiến Sinh Tử,\nCàng sống lâu càng nhiều quà,\nHãy tìm Quy Lão để Khám Phá");
                  case 5 -> Service.gI().sendThongBaoOK(player, "Thỏi Vàng: \n- Có thể nhận từ leo tháp ở Ôsin (ĐHVT 12H mỗi ngày),\n- Có thể nhận thưởng từ Top Siêu Hạng,\n- Có thể tham gia Map War để nhặt,\n-Có thể săn Boss (ai cũng nhặt được)");
            }
        }
        else if (player.iDMark.getIndexMenu() == ConstNpc.MAIL_BOX) {
            switch (select) {
                case 0 -> ShopService.gI().opendShop(player, "ITEMS_MAIL_BOX", true);
                case 1 -> NpcService.gI().createMenuConMeo(player,
                        ConstNpc.CONFIRM_REMOVE_ALL_ITEM_MAIL_BOX, this.avartar,
                        "|3|Bạn chắc muốn xóa hết vật phẩm trong hòm thư?\n"
                        + "|7|Sau khi xóa sẽ không thể khôi phục!", "Đồng ý", "Hủy bỏ");
            }
        } 
        else if (player.iDMark.getIndexMenu() == ConstNpc.NAP_TIEN) {
            switch (select) {
                case 0 -> {
                    List<String> menu = new ArrayList<>();
                    for (int i = 0; i < napVang.length; i++) {
                        menu.add(Util.mumberToLouis(napVang[i][0] / 1) + " Cá\n"
                                + Util.mumberToLouis(napVang[i][1] * costNapVang) + " Thỏi");
                    }
                    createOtherMenu(player, ConstNpc.NAP_VANG, "Ta rất thương con!\nTa sẽ tặng thêm con 1 ít Thạch Thần:\n"
                            + "- 20.000 Cá tặng 2 Thạch Thần - 50.000 Cá tặng 6 Thạch Thần\n- 100.000 Cá tặng 13 Thạch Thần - 200.000 Cá tặng 30 Thạch Thần\n- 500.000 Cá tặng 90 Thạch Thần - 1.000.000 Cá tặng 200 Thạch Thần", menu.toArray(new String[0]));
                }
                case 1 -> {
                    if (player.getSession().goldBar > 0) {
                        Item thoiVang = ItemService.gI().createNewItem((short) 457, player.getSession().goldBar);
                        if (InventoryService.gI().getCountEmptyBag(player) < 1) {
                            Service.gI().sendThongBao(player, "Cần ít nhất 1 ô trống trong hành trang");
                            return;
                        }
                        InventoryService.gI().addItemBag(player, thoiVang);
                        InventoryService.gI().sendItemBag(player);
                        Service.gI().sendThongBao(player, "Bạn đã nhận được " + player.getSession().goldBar + " thỏi vàng");
                        PlayerDAO.subGoldBar(player, player.getSession().goldBar);
                    } else {
                        Service.gI().sendThongBao(player, "Con không còn thỏi vàng nào để nhận.");
                    }
                }
            }
        } 
        else if (player.iDMark.getIndexMenu() == ConstNpc.NAP_VANG) {
    if (select >= 0 && select < napVang.length) {
        if (player.getSession().cash >= napVang[select][0]) {

            int tien = napVang[select][0];
            int vang = napVang[select][1] * costNapVang;

            // Trừ tiền + cộng thỏi vàng
            PlayerDAO.subcash(player, tien);
            PlayerDAO.subGoldBar(player, -vang);

            // ===== THƯỞNG ITEM 674 =====
            int soLuongThuong = 0;
            switch (tien) {
                case 20000 -> soLuongThuong = 2;
                case 50000 -> soLuongThuong = 6;
                case 100000 -> soLuongThuong = 13;
                case 200000 -> soLuongThuong = 30;
                case 500000 -> soLuongThuong = 90;
                case 1000000 -> soLuongThuong = 200;
            }

            if (soLuongThuong > 0) {
                if (InventoryService.gI().getCountEmptyBag(player) < 1) {
                    Service.gI().sendThongBao(player, "Cần ít nhất 1 ô trống để nhận quà!");
                    return;
                }

                Item thuong = ItemService.gI().createNewItem((short) 674, soLuongThuong);
                InventoryService.gI().addItemBag(player, thuong);
            }

            InventoryService.gI().sendItemBag(player);

            Service.gI().sendThongBao(player,
                    "Nhận "
                    + Util.mumberToLouis(vang) + " thỏi"
                    + (soLuongThuong > 0 ? "\n+ " + soLuongThuong + " Thạch Thần" : ""));

        } else {
            Service.gI().sendThongBao(player, "Không đủ số dư");
        }
    }

            }
        }
    }
