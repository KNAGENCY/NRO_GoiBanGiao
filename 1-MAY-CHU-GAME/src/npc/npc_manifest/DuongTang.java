package npc.npc_manifest;

import consts.ConstNpc;
import npc.Npc;
import player.Player;
import services.Service;
import services.func.ChangeMapService;
import shop.ShopService;
import utils.Util;

public class DuongTang extends Npc {

    public DuongTang(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }
@Override
public void openBaseMenu(Player player) {
    if (canOpenNpc(player)) {
        switch (mapId) {
            case 0 -> {
                this.createOtherMenu(player, ConstNpc.BASE_MENU,
                        "Ngũ Hành Sơn tuổi thơ\n- Pem quái luôn gây 80.000 sát thương\n- Luôn nhận được 5.000 - 30.000 TNSM mỗi Hit\nCó cả Mảnh Vỡ & Mảnh Hồn Bông Tai\nCó 20% May Mắn để vào NHS Tối Thượng, ở Quy Lão tỉ lệ là 50%",
                        "Thử\n1 Lần", "Shop\nTrùm Cuối");
            }
            case 123, 122 -> {
                this.createOtherMenu(player, ConstNpc.BASE_MENU,
                        "Ngươi Có Muốn Quay Về Thử Lại Không?",
                        "Quay Về", "Refuse");
            }
            default -> super.openBaseMenu(player);
        }
    }
}

@Override
public void confirmMenu(Player player, int select) {
    if (canOpenNpc(player)) {
        if (player.iDMark.isBaseMenu()) {
            switch (mapId) {
               case 0 -> {
    if (select == 0) {
        if (player.playerTask.taskMain.id < 18) {
            Service.gI().sendThongBaoOK(player, "Yêu Cầu Đến Nhiệm Vụ Bojack, hẹ hẹ!");
            return;
        }
        long phi = 5_000_000L;
        if (player.inventory.gold >= phi) {
            player.inventory.gold -= phi;
            Service.gI().sendMoney(player);
            int rd = Util.nextInt(100);
            if (rd < 80) {
                ChangeMapService.gI().changeMap(player, 123, -1, 200, 336);
            } else {
                ChangeMapService.gI().changeMap(player, 122, -1, 200, 336);
            }
        } else {
            Service.gI().sendThongBao(player, "Ngươi Không Đủ 5 Triệu Vàng!");
        }
    } else if (select == 1) {
        ShopService.gI().opendShop(player, "TRUONG_LAO", false);
    }
}
                
                case 123, 122 -> {
                    if (select == 0) {
                        ChangeMapService.gI().changeMap(player, 0, -1, 200, 336);
                    }
                }
            }
        }
    }
}
}