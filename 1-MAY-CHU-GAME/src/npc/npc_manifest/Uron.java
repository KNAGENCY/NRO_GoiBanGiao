package npc.npc_manifest;

/**
 *
 * @author GiaCatLuong
 */

import consts.ConstNpc;
import npc.Npc;
import player.Player;
import services.Service;
import shop.ShopService;

public class Uron extends Npc {

    public Uron(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player pl) {
        if (canOpenNpc(pl)) {
            // Hiển thị menu gồm 2 lựa chọn: Shop và Học Skill
            this.createOtherMenu(pl, ConstNpc.BASE_MENU,
                    "Ta có thể giúp gì cho ngươi?",
                    "Shop", "Học Skill");
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.iDMark.isBaseMenu()) {
                switch (select) {
                    case 0: // Chọn Shop
                        ShopService.gI().opendShop(player, "URON", false);
                        break;
                    case 1: // Chọn Học Skill
                        Service.gI().sendThongBaoOK(player, "Đến nam Kame tìm thần Vũ Trụ");
                        break;
                }
            }
        }
    }
}