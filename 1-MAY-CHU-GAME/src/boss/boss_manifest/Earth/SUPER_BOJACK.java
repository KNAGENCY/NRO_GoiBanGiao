package boss.boss_manifest.Earth;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossesData;
import map.ItemMap;
import player.Player;
import services.Service;
import services.TaskService;
import utils.Util;

public class SUPER_BOJACK extends Boss {

    private long st;

    public SUPER_BOJACK() throws Exception {
        super(BossID.SUPER_BOJACK, false, true, BossesData.SUPER_BOJACK_2);
    }

    @Override
    public void moveTo(int x, int y) {
        if (this.currentLevel == 1) {
            return;
        }
        super.moveTo(x, y);
    }
@Override
public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
    // Giới hạn sát thương tối đa = 2999
    if (damage > 69999) {
        damage = 69999;
    }
    return super.injured(plAtt, damage, piercing, isMobAttack);
}
@Override
public void reward(Player plKill) {
      TaskService.gI().checkDoneTaskKillBoss(plKill, this);
    if (this.zone == null) {
        return;
    }

    int totalDrop = 2;
    int radius = 80; // bán kính vòng tròn

    for (int i = 0; i < totalDrop; i++) {
              int itemId = Util.isTrue(99, 100) ? 708 : 702;

        // Tính vị trí theo hình tròn
        double angle = 2 * Math.PI * i / totalDrop;
        int x = this.location.x + (int) (Math.cos(angle) * radius);
        int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);

        ItemMap itemMap = new ItemMap(
                this.zone,
                itemId,
                1,   // số lượng = 1
                x,
                y,
                -1   // ai cũng nhặt được
        );

        Service.gI().dropItemMap(this.zone, itemMap);
    }
}



    @Override
    protected void notifyJoinMap() {
        if (this.currentLevel == 1) {
            return;
        }
        super.notifyJoinMap();
    }

    @Override
    public void joinMap() {
        super.joinMap();
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
