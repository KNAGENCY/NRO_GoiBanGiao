package boss.boss_manifest.NamekGinyuForce;

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

public class SO1_NM extends Boss {

    private long st;

    public SO1_NM() throws Exception {
        super(BossID.SO_1_NM, false, true, BossesData.SO_1_NM);
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
  if (damage > 99999999) {
        damage = 99999999;
    }
    return super.injured(plAtt, damage, piercing, isMobAttack);
}
@Override
public void reward(Player plKill) {
       TaskService.gI().checkDoneTaskKillBoss(plKill, this);
        this.pointBoss(plKill);
    if (Util.isTrue(30, 100)) {
        // --- 20% rơi x1 1083 hoặc 1078 (Chỉ người hạ boss nhặt được) ---
        int[] rareItems = {1083, 1078,1084, 1085, 1086};
        int itemId = rareItems[Util.nextInt(0, rareItems.length - 1)];

        ItemMap it = new ItemMap(this.zone, itemId, 1,
                this.location.x,
                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                plKill.id);
        Service.gI().dropItemMap(this.zone, it);
        
    } else {
        // --- 80% chọn ngẫu nhiên DUY NHẤT 1 loại từ 1066->1070 (Rơi 10 cái loại đó) ---
        int[] commonItems = {1066, 1067, 1068, 1069, 1070};
        // Chọn ngẫu nhiên 1 ID từ mảng trên
        int randomItemId = commonItems[Util.nextInt(0, commonItems.length - 1)];
        
        for (int i = 0; i < 10; i++) {
            ItemMap it = new ItemMap(this.zone, randomItemId, 1,
                    this.location.x + Util.nextInt(-80, 80), 
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    -1); // Mọi người đều nhặt được
            Service.gI().dropItemMap(this.zone, it);
        }
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
    public void doneChatS() {
        this.changeStatus(BossStatus.AFK);
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
            if (boss.id == BossID.SO_2_NM && !boss.isDie()) {
                return;
            }

        }
        this.parentBoss.changeStatus(BossStatus.ACTIVE);
    }

}
