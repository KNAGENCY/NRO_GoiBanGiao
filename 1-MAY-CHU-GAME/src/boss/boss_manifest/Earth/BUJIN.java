package boss.boss_manifest.Earth;

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
import utils.Util;

public class BUJIN extends Boss {

    private long st;

    public BUJIN() throws Exception {
        super(BossID.BUJIN, false, true, BossesData.BUJIN);
    }

    @Override
    public void moveTo(int x, int y) {
        if (this.currentLevel == 1) {
            return;
        }
        super.moveTo(x, y);
    }

 @Override
public void reward(Player plKill) {
    if (this.zone == null) {
        return;
    }

    int totalDrop = 2;
    int radius = 80; // bán kính vòng tròn

    for (int i = 0; i < totalDrop; i++) {
              int itemId = Util.isTrue(99, 100) ? 706 : 703;

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
        super.joinMap();
        st = System.currentTimeMillis();
    }

    @Override
    public void doneChatE() {
        if (this.parentBoss == null || this.parentBoss.bossAppearTogether == null
                || this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel] == null) {
            return;
        }
        for (Boss boss : this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel]) {
            if ((boss.id == BossID.KOGU || boss.id == BossID.ZANGYA || boss.id == BossID.BIDO) && !boss.isDie()) {
                return;
            }
        }
        this.parentBoss.changeStatus(BossStatus.ACTIVE);
    }

}
