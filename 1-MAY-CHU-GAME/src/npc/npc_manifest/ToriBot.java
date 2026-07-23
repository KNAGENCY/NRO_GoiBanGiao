package npc.npc_manifest;

import consts.ConstNpc;
import item.Item;
import models.War.War;
import npc.Npc;
import player.Player;
import services.Service;
import services.func.ChangeMapService;
import shop.ShopService;

public class ToriBot extends Npc {

    public ToriBot(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            switch (this.mapId) {
                case 182 -> {
                    if (War.gI().isStartWarTime()) {
                        this.createOtherMenu(player, ConstNpc.BASE_MENU,
                                "Mày Đã Có Thể Đi Vào Map War\nCắn Full Item Các Thứ Vào",
                                "Đồng ý", "Trốn");
                    } else {
this.createOtherMenu(player, ConstNpc.BASE_MENU,
        "Hãy Nhanh Chóng Hạ Quái Để Nhận Đồ Vip\n"
      + "Muốn Vào Map War Thì Đi Theo Ta\n"
      + "Ta Sẽ Dẫn Ngươi Đi Vào Lúc 21h45 - 21h45 - 21h45!",
        "Hạ Quái",
        "Về Đảo");

                    }
                }
                case 183 -> {
                    this.createOtherMenu(player, 1,
                            "Cửa hàng chuyên kì trân dị bảo không ở đâu có được",
                            "Cửa hàng");
                }
                default -> {
                }
            }
        }
    }

   @Override
public void confirmMenu(Player player, int select) {
    if (canOpenNpc(player)) {

        // 📍 MAP 182
        if (this.mapId == 182) {
            if (player.iDMark.isBaseMenu()) {
                switch (select) {

                    // 🔥 GIỮ NGUYÊN CHỨC NĂNG CŨ
                    case 0 -> {
                        if (War.gI().isStartWarTime()) {
                            Service.gI().changeFlag(player, 8);
                            ChangeMapService.gI().changeMap(
                                    player,
                                    War.MAP_WAR_ID,
                                    -1,
                                    111,
                                    288
                            );
                        } else {
                            // ➕ BỔ SUNG: CHƯA TỚI GIỜ WAR
                            Service.gI().sendThongBao(player,
                                    "Quái ở trong map, hãy tiêu diệt!");
                        }
                    }

                    // 🚪 GO OUT
                    case 1 -> {
                        boolean empty = true;
                        for (int i = 0; i < player.inventory.itemsBody.size(); i++) {
                            Item item = player.inventory.itemsBody.get(i);
                            if (item != null && item.isNotNullItem()) {
                                empty = false;
                                break;
                            }
                        }

                        if (!empty) {
                            Service.gI().sendThongBao(player,
                                    "Vui lòng tháo toàn bộ trang bị trước khi di chuyển!");
                            return;
                        }

                        ChangeMapService.gI().changeMap(player, 5, -1, 1070, 408);
                    }
                }
            }
            } else if (this.mapId == 183) {
                switch (select) {
                    case 0 -> ShopService.gI().opendShop(player, "TORI_BOT", true);
                }
            }
        }
    }
}
