package boss.boss_manifest.Wolves;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossesData;
import item.Item;
import item.Item.ItemOption;
import map.ItemMap;
import player.Player;
import services.ItemService;
import services.Service;
import services.TaskService;
import utils.Util;

public class S extends Boss {

    private long st;

    public S() throws Exception {
        super(BossID.S, false, true, BossesData.S1, BossesData.S2, BossesData.S3);
    }

 @Override
    public void reward(Player plKill) {
        this.pointBoss(plKill);

        // ✅ Cập nhật nhiệm vụ giết boss
        TaskService.gI().checkDoneTaskKillBoss(plKill, this);

        if (Util.isTrue(20, 100)) {
            // 40%: rơi đồ TL, chỉ người hạ boss nhặt được
            ItemMap it = ItemService.gI().randDoTL(
                    this.zone,
                    1,
                    this.location.x,
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    plKill.id // chỉ người hạ boss nhặt
            );
            Service.gI().dropItemMap(this.zone, it);
        } else {
            // 60% còn lại: rơi 4 item id = 457, ai cũng nhặt được
            for (int i = 0; i < 3; i++) {
                int x = this.location.x + Util.nextInt(-60, 60); // rải xung quanh
                int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);
                ItemMap it = new ItemMap(this.zone, (short) 674, 1, x, y, -1);
                // -1 => không gán id người chơi => ai cũng nhặt được
                Service.gI().dropItemMap(this.zone, it);
            }
        }
    }
    @Override
    public void joinMap() {
        super.joinMap(); // To change body of generated methods, choose Tools | Templates.
        st = System.currentTimeMillis();
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

}
