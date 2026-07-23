package boss.boss_manifest.Nappa;

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

public class MapDauDinh extends Boss {

    private long st;

    public MapDauDinh() throws Exception {
        super(BossID.MAP_DAU_DINH, false, true, BossesData.MAP_DAU_DINH);
    }

    @Override
    public void joinMap() {
        super.joinMap();
        st = System.currentTimeMillis();
    }
@Override
public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
    // Giới hạn sát thương tối đa = 2999
    if (damage > 59999990) {
        damage = 59999990;
    }
    return super.injured(plAtt, damage, piercing, isMobAttack);
}

  @Override
public void reward(Player plKill) {
    // Tặng điểm cho người hạ boss
  //  this.pointBoss(plKill);

    TaskService.gI().checkDoneTaskKillBoss(plKill, this);

    // 50% rơi item 1194, số lượng 2, mỗi cái 1 vị trí, chỉ người hạ boss nhặt được
    if (Util.isTrue(50, 100)) {
        for (int i = 0; i < 2; i++) {
            ItemMap it = new ItemMap(this.zone, 1194, 1,
                    this.location.x + Util.nextInt(-20, 20), // rải 2 vị trí xung quanh boss
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    plKill.id); // chỉ người hạ boss nhặt được
            Service.gI().dropItemMap(this.zone, it);
        }
    } 
    // 50% còn lại rơi item 751, 10 item, mọi người nhặt được
    else {
        for (int i = 0; i < 10; i++) {
            ItemMap it = new ItemMap(this.zone, 750, 1,
                    this.location.x + Util.nextInt(-60, 60),
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    -1); // -1 nghĩa là tất cả người chơi đều nhặt được
            Service.gI().dropItemMap(this.zone, it);
        }
    }
}



    @Override
    public void autoLeaveMap() {
        if (Util.canDoWithTime(st, 900000)) {
            this.changeStatus(BossStatus.LEAVE_MAP);
        }
        // if (this.zone != null && this.zone.getNumOfPlayers() > 0) {
        // st = System.currentTimeMillis();
        // }
    }
}
