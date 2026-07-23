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
import consts.ConstPlayer;
import map.ItemMap;
import player.Player;
import services.ItemService;
import services.PlayerService;
import services.Service;
import services.TaskService;
import utils.Util;

public class Android14 extends Boss {

    public boolean callApk13;

    public Android14() throws Exception {
        super(BossID.ANDROID_14, false, true, BossesData.ANDROID_14);
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
    protected void resetBase() {
        super.resetBase();
        this.callApk13 = false;
    }

    @Override
    public void active() {
        if (this.typePk == ConstPlayer.NON_PK && !this.callApk13) {
            this.changeToTypePK();
        }
        this.attack();
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (!this.callApk13 && damage >= this.nPoint.hp) {
            this.callApk13();
            return 0;
        }
        return super.injured(plAtt, damage, piercing, isMobAttack);
    }

    public void callApk13() {
        if (this.bossAppearTogether == null || this.bossAppearTogether[this.currentLevel] == null) {
            return;
        }
        for (Boss boss : this.bossAppearTogether[this.currentLevel]) {
            if (boss.id == BossID.ANDROID_13) {
                boss.changeStatus(BossStatus.RESPAWN);
            } else if (boss.id == BossID.ANDROID_15) {
                boss.changeToTypeNonPK();
                ((Android15) boss).callApk13 = true;
                ((Android15) boss).recoverHP();
            }
        }
        this.changeToTypeNonPK();
        this.recoverHP();
        this.callApk13 = true;
    }

    public void recoverHP() {
        PlayerService.gI().hoiPhuc(this, this.nPoint.hpMax, 0);
    }

    @Override
    public void doneChatS() {
        if (this.bossAppearTogether == null || this.bossAppearTogether[this.currentLevel] == null) {
            return;
        }
        for (Boss boss : this.bossAppearTogether[this.currentLevel]) {
            if (boss.id == BossID.ANDROID_15) {
                boss.changeToTypePK();
                break;
            }
        }
    }

}
