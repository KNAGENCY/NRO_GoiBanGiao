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

public class BOJACK extends Boss {

    private long st;

    public BOJACK() throws Exception {
        super(BossID.BOJACK, false, true, BossesData.BOJACK, BossesData.SUPER_BOJACK);
    }

 @Override
public void reward(Player plKill) {
    this.pointBoss(plKill);
    if (this.zone == null) {
        return;
    }

    int totalDrop = 2;
    int radius = 80; // bán kính vòng tròn

    for (int i = 0; i < totalDrop; i++) {
              int itemId = Util.isTrue(99, 100) ? 704 : 702;

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

    @Override
    public void doneChatS() {
        if (this.currentLevel == 1) {
            return;
        }
        this.changeStatus(BossStatus.AFK);
    }
}
