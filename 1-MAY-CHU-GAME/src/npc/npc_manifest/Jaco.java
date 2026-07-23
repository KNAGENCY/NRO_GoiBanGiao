package npc.npc_manifest;

import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.Service;
import services.func.ChangeMapService;

public class Jaco extends Npc {

    public Jaco(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) {
            return;
        }

        if (this.mapId == 42 || this.mapId == 43 || this.mapId == 44) {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Đây là nơi bí mật, ngươi có muốn đến Vùng Đất Săn Boss?",
                    "Đến ngay", "Nâng cấp\nBông tai 3", "Mở chỉ số\nBông tai 3",
                    "Hủy Diệt\nBông Tai 4", "Từ chối");
        } else if (this.mapId == 178) {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Ngươi muốn quay về hành tinh nào?",
                    "Trái Đất", "Namếc", "Xayda");
        } else {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Nơi này rất nguy hiểm, nhưng có thể tìm thấy đá Ánh Dương. Ngươi có muốn đến đó không?",
                    "Vùng Đất\nÁnh Dương", "Nâng cấp\nBông tai 3", "Mở chỉ số\nBông tai 3",
                    "Nâng Cấp\nBông Tai 4", "Từ chối");
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) {
            return;
        }

        if (player.iDMark.getIndexMenu() == ConstNpc.MENU_START_COMBINE) {
            if (select == 0 && (player.combine.typeCombine == CombineService.NANG_CAP_BONG_TAI_3
                    || player.combine.typeCombine == CombineService.NANG_CHI_SO_BONG_TAI_3
                    || player.combine.typeCombine == CombineService.NANG_CAP_BONG_TAI_4)) {
                npcChat(player, "um ba la xì bùa");
                CombineService.gI().startCombine(player);
            }
            return;
        }
        if (!player.iDMark.isBaseMenu()) {
            return;
        }

        if (this.mapId == 42 || this.mapId == 43 || this.mapId == 44) {
            switch (select) {
                case 0:
                    ChangeMapService.gI().changeMap(player, 178, -1, -1, 0);
                    Service.gI().sendThongBao(player, "Chúc ngươi may mắn trên chiến trường!");
                    break;
                case 1:
                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_BONG_TAI_3);
                    break;
                case 2:
                    CombineService.gI().openTabCombine(player, CombineService.NANG_CHI_SO_BONG_TAI_3);
                    break;
                case 3:
                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_BONG_TAI_4);
                    break;
                default:
                    break;
            }
        } else if (this.mapId == 178) {
            switch (select) {
                case 0:
                    ChangeMapService.gI().changeMap(player, 42, -1, -1, 0);
                    Service.gI().sendThongBao(player, "Quay về Trái Đất!");
                    break;
                case 1:
                    ChangeMapService.gI().changeMap(player, 43, -1, -1, 0);
                    Service.gI().sendThongBao(player, "Quay về Namếc!");
                    break;
                case 2:
                    ChangeMapService.gI().changeMap(player, 44, -1, -1, 0);
                    Service.gI().sendThongBao(player, "Quay về Xayda!");
                    break;
                default:
                    break;
            }
        } else {
            switch (select) {
                case 0:
                    Item phieuBeNgon = InventoryService.gI().findItemBagByTemp(player, 1194);
                    if (phieuBeNgon != null && phieuBeNgon.quantity >= 1) {
                        InventoryService.gI().subQuantityItemsBag(player, phieuBeNgon, 1);
                        InventoryService.gI().sendItemBag(player);
                        ChangeMapService.gI().changeMap(player, 181, -1, -1, 0);
                        Service.gI().sendThongBao(player, "Chúc ngươi may mắn!");
                    } else {
                        Service.gI().sendThongBao(player, "Ngươi cần 1 Phiếu Bé Ngoan!");
                    }
                    break;
                case 1:
                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_BONG_TAI_3);
                    break;
                case 2:
                    CombineService.gI().openTabCombine(player, CombineService.NANG_CHI_SO_BONG_TAI_3);
                    break;
                case 3:
                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_BONG_TAI_4);
                    break;
                default:
                    break;
            }
        }
    }
}
