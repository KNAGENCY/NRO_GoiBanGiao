package boss.boss_manifest.Cell;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossManager;
import boss.BossStatus;
import boss.BossesData;
import map.ItemMap;
import player.Player;
import services.Service;
import services.TaskService;
import services.func.ChangeMapService;
import utils.Util;

public class XENCON4 extends Boss {

    private long st;

    public XENCON4() throws Exception {
        super(BossID.XEN_CON_4, false, true, BossesData.XEN_CON_4);
    }

  @Override
public void reward(Player plKill) {
 //   this.pointBoss(plKill);

      int rd = Util.nextInt(1, 100);
    int itemId;
    if (rd <= 20) {
        itemId = 17;
    } else if (rd <= 40) {
        itemId = 18;
    } else if (rd <= 60) {
        itemId = 19;
    } else if (rd <= 80) {
        itemId = 20;
    } else {
        itemId = 702;
    }


    // thả item ra map tại vị trí boss chết
    ItemMap it = new ItemMap(
        this.zone,
        itemId,
        2,
        this.location.x,
        this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
        plKill.id
    );
    Service.gI().dropItemMap(this.zone, it);

    TaskService.gI().checkDoneTaskKillBoss(plKill, this);
}

    @Override
    public void joinMap() {
        st = System.currentTimeMillis();
        this.zone = this.parentBoss.zone;
        ChangeMapService.gI().changeMap(this, this.zone,
                this.parentBoss.location.x + Util.nextInt(-100, 100), this.parentBoss.location.y);
        Service.gI().sendFlagBag(this);
        this.notifyJoinMap();
        this.changeStatus(BossStatus.CHAT_S);
    }

    @Override
    public void doneChatE() {
        if (this.parentBoss == null || this.parentBoss.bossAppearTogether == null
                || this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel] == null) {
            return;
        }
        for (Boss boss : this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel]) {
            if ((boss.id == BossID.XEN_CON_2 || boss.id == BossID.XEN_CON_3 || boss.id == BossID.XEN_CON_1
                    || boss.id == BossID.XEN_CON_5 || boss.id == BossID.XEN_CON_6 || boss.id == BossID.XEN_CON_7)
                    && !boss.isDie()) {
                return;
            }
        }
        this.parentBoss.changeStatus(BossStatus.ACTIVE);
    }

    @Override
    public void leaveMap() {
        ChangeMapService.gI().exitMap(this);
        this.lastZone = null;
        this.lastTimeRest = System.currentTimeMillis();
        this.changeStatus(BossStatus.REST);
        // BossManager.gI().removeBoss(this);
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
