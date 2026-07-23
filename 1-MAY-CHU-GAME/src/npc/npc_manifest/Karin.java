package npc.npc_manifest;

/**
 *
 * NPC Karin - Chỉ còn menu Nhiệm vụ
 *
 */

import consts.ConstNpc;
import npc.Npc;
import player.Player;
import services.TaskService;
import services.func.ChangeMapService;

public class Karin extends Npc {

    public Karin(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) {
            return;
        }

        if (this.mapId == 46) {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Tàu 77 hắn mạnh lắm\nTa không thể đánh bại hắn 1 mình\nHãy vào Phòng tập thời gian\nPhòng Tập Thời Gian ở Thượng Đế (trên Thần Điện)",
                    "Nhiệm vụ\n10.000 SĐG","Đấm\nThần Mèo");
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) {
            return;
        }

        if (this.mapId != 46) {
            return;
        }

      if (player.iDMark.isBaseMenu()) {
    switch (select) {
        case 0:
            TaskService.gI().checkDoneTaskTalkNpc(player, this);
            break;

        case 1:
            // Đấm Thần Mèo -> vào phòng tập thời gian (map 49)
            ChangeMapService.gI().changeMapBySpaceShip(player, 49, -1, 300);
            break;
    }
}
    }
}