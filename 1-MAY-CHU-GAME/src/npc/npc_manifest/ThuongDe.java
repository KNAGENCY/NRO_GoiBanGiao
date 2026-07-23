package npc.npc_manifest;

import consts.ConstNpc;
import npc.Npc;
import player.Player;
import services.Service;
import services.func.ChangeMapService;
import models.Combine.CombineService;
import services.InventoryService;
import item.Item;
import services.ItemService;

public class ThuongDe extends Npc {

    private static final long REQUIRED_POWER = 20000000000L; // 10 tỷ sức mạnh

    public ThuongDe(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) {
            return;
        }

        // MAP 45
        if (this.mapId == 45) {
            createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Ta sẽ đưa con đến Phòng tập thời gian\nTàu 77 đang trốn trong đó\nHãy đổi khu để tìm hắn!",
                    "Đi tìm\nTàu 77");
        }

        // MAP 154
        if (this.mapId == 154) {
            createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Ta là Thượng Đế.\nCon cần ta giúp gì lào?\nNgục Tù Là Nơi Dành Riêng Cho Up Đệ Tử\nĐệ Pem Bao Nhiêu Dame Sẽ Nhận Bấy Nhiêu TNSM",
                    "Đến\nNgục Tù",
                    "Trang Bị\nThiên Sứ"
                    // ,
                    // "Up\nThiên Sứ",
                    // "Đổi\nĐệ Cumber"
                );
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) {
            return;
        }

        // ===== MAP 45 =====
        if (this.mapId == 45 && player.iDMark.isBaseMenu()) {
            if (select == 0) {
                ChangeMapService.gI().changeMapBySpaceShip(player, 49, -1, 300);
            }
        }

        // ===== MAP 154 =====
        if (this.mapId == 154 && player.iDMark.isBaseMenu()) {

            switch (select) {

                case 0: // Ngục tù
                    if (player.nPoint.power < REQUIRED_POWER) {
                        Service.gI().sendThongBao(player, "Bạn cần 40 tỷ sức mạnh");
                        return;
                    }
                    ChangeMapService.gI().changeMapBySpaceShip(player, 155, -1, 300);
                    break;

                case 1: // Chế đồ Thiên Sứ
                    player.iDMark.setNpcChose(this);
                    player.combine.clearItemCombine();
                    CombineService.gI().openTabCombine(player, CombineService.CHE_TAO_TRANG_BI_THIEN_SU);
                    break;
            }}

//                 case 2: // Map Thiên Sứ
//                     ChangeMapService.gI().changeMapBySpaceShip(player, 174, -1, 300);
//                     break;

//                 case 3: // Đổi đệ Cumber
//                     Item item571 = InventoryService.gI().findItemBag(player, 571);
//                     Item item572 = InventoryService.gI().findItemBag(player, 572);

//                     if (item571 == null || item571.quantity < 49 ||
//                         item572 == null || item572.quantity < 49) {

//                         Service.gI().sendThongBao(player,
//                                 "Cần Rương Bạc và Rương Vàng");
//                         return;
//                     }

//                  // trừ nguyên liệu
// InventoryService.gI().subQuantityItemsBag(player, item571, 49);
// InventoryService.gI().subQuantityItemsBag(player, item572, 49);

// // cập nhật hành trang
// InventoryService.gI().sendItemBag(player);

// // tạo item 1711
// Item item = ItemService.gI().createNewItem((short) 1711);
// item.quantity = 1;

// InventoryService.gI().addItemBag(player, item);

// // cập nhật hành trang lần nữa
// InventoryService.gI().sendItemBag(player);

// Service.gI().sendThongBao(player,
//         "Bạn đã đổi thành công Trứng đệ Cumber!");

//                     break;
//             }
//         }

        // START COMBINE
        else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_START_COMBINE) {
            CombineService.gI().startCombine(player);
        }
    }
}