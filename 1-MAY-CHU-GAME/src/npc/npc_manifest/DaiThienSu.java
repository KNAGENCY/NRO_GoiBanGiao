package npc.npc_manifest;

import consts.ConstNpc;
import item.Item;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import shop.ShopService;

public class DaiThienSu extends Npc {

    /* ================= MENU ID ================= */
    private static final int MENU_GOC = 1000;
    private static final int MENU_DOT_TIEN = 1001;

    public DaiThienSu(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    /* ================= OPEN MENU ================= */
    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) return;

        createOtherMenu(player, ConstNpc.BASE_MENU,
                "Tâm tao không ác\nSao số tao lại khổ...\nLòng sống vị tha nhưng đời cứ đưa tao vào thế khó\nTình cảm của tao thật lòng sao trả lại tao toàn điều dối trá?",
                "Thực Đơn\nCủa\nCậu Vàng",
                "Bán\nCậu Vàng",
                "Đốt tiệm",
                "Đóng");

        player.iDMark.setIndexMenu(MENU_GOC);
    }

    /* ================= CONFIRM ================= */
    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) return;

        /* ===== MENU GỐC ===== */
        if (player.iDMark.getIndexMenu() == MENU_GOC) {
            switch (select) {
                case 0 -> { // Thực đơn
                  ShopService.gI().opendShop(player, "OSIN", false);
                }

                case 1 -> { // Giao cậu vàng
                    giaoCauVang(player);
                }

                case 2 -> { // Đốt tiệm
                    createOtherMenu(player, ConstNpc.BASE_MENU,
                            "Ê, ê\nÊ, ê, Nịt pẹ mày nha\nMày có cút nhanh không\ntao mà gọi Cậu Vàng ra là mày xong đời?\nĐốt tiệm nhận ngay rương Thần Linh và Hủy Diệt",
                            "Đốt tiệm",
                            "Đóng");
                    player.iDMark.setIndexMenu(MENU_DOT_TIEN);
                }
            }
            return;
        }

        /* ===== MENU ĐỐT TIỆM ===== */
        if (player.iDMark.getIndexMenu() == MENU_DOT_TIEN) {
            if (select == 0) {
                dotTiem(player);
            }
        }
    }

    /* ================= LOGIC ================= */

    private void giaoCauVang(Player player) {
    Item cauVang = InventoryService.gI().findItemBag(player, 1742);

    if (cauVang == null || cauVang.quantity < 1) {
        Service.gI().sendThongBao(player, "Không tìm thấy\nCậu Vàng");
        return;
    }

    if (InventoryService.gI().getCountEmptyBag(player) < 1) {
        Service.gI().sendThongBao(player, "Hành trang không đủ chỗ trống");
        return;
    }

    InventoryService.gI().subQuantityItemsBag(player, cauVang, 1);

    giveItem(player, 568, 1);

    Service.gI().sendThongBaoOK(player, "Mày nhận được x1 Trứng Mabu béo");
}

private void dotTiem(Player player) {
    Item queDot = InventoryService.gI().findItemBag(player, 1237);
    Item canhKho = InventoryService.gI().findItemBag(player, 1240);

    if (queDot == null || queDot.quantity < 10
            || canhKho == null || canhKho.quantity < 10) {
        Service.gI().sendThongBao(player,
                "Cần x10 que đốt và x10 cành khô");
        return;
    }

    if (InventoryService.gI().getCountEmptyBag(player) < 1) {
        Service.gI().sendThongBao(player, "Hành trang không đủ chỗ trống");
        return;
    }

    InventoryService.gI().subQuantityItemsBag(player, queDot, 10);
    InventoryService.gI().subQuantityItemsBag(player, canhKho, 10);

   Service.gI().sendThongBaoOK(player,
        "Đốt mả mẹ mày à, thù oán đéo gì?\nBạn đíu nhận được gì");

giveItem(player, 1453, 1);
giveItem(player, 1173, 1);

}


    /* ================= SUPPORT ================= */
    private void giveItem(Player player, int tempId, int quantity) {
        Item item = ItemService.gI().createNewItem((short) tempId, quantity);
        InventoryService.gI().addItemBag(player, item);
        InventoryService.gI().sendItemBag(player);
    }
}
