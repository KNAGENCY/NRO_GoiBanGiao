package boss.boss_manifest.Android;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossStatus;
import boss.BossesData;
import map.ItemMap;
import player.Player;
import services.Service;
import services.TaskService;
import utils.Util;

public class KingKong extends Boss {

    public KingKong() throws Exception {
        super(BossID.KING_KONG, false, true, BossesData.KING_KONG);
    }

 @Override
public void reward(Player plKill) {

    TaskService.gI().checkDoneTaskKillBoss(plKill, this);

    // 50% rơi 1194
    if (Util.isTrue(50, 100)) {
        for (int i = 0; i < 2; i++) {
            ItemMap it = new ItemMap(
                    this.zone,
                    1194,
                    1,
                    this.location.x + Util.nextInt(-30, 30),
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    plKill.id);
            Service.gI().dropItemMap(this.zone, it);
        }
    }
    // 50% rơi 1150, 1151, 1152, 1153
    else {
        int[] items = {1150, 1151, 1152, 1153};

        for (int itemId : items) {
            int x = this.location.x + Util.nextInt(-80, 80);
            int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);

            ItemMap it = new ItemMap(this.zone, itemId, 1, x, y, -1);
            Service.gI().dropItemMap(this.zone, it);
        }
    }

    // Luôn rơi 10 lá dong (751), ai cũng nhặt được
    for (int i = 0; i < 10; i++) {
        int x = this.location.x + Util.nextInt(-120, 120);
        int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);

        ItemMap it = new ItemMap(this.zone, 751, 1, x, y, -1);
        Service.gI().dropItemMap(this.zone, it);
    }
}


    @Override
    public void autoLeaveMap() {
        if (Util.canDoWithTime(st, 900000)) {
            this.leaveMapNew();
        }
        if (this.zone != null && this.zone.getNumOfPlayers() > 0) {
            st = System.currentTimeMillis();
        }
    }

    @Override
    public void joinMap() {
        super.joinMap(); // To change body of generated methods, choose Tools | Templates.
        st = System.currentTimeMillis();
    }

    private long st;

    @Override
    public void doneChatS() {
        this.changeStatus(BossStatus.AFK);
    }
}
