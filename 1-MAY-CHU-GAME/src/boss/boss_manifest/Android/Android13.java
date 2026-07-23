package boss.boss_manifest.Android;

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
import services.ItemService;
import services.Service;
import services.TaskService;
import utils.Util;

public class Android13 extends Boss {

    public Android13() throws Exception {
        super(BossID.ANDROID_13, false, true, BossesData.ANDROID_13);
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
    public void doneChatS() {
        if (this.parentBoss == null) {
            return;
        }
        if (this.parentBoss.bossAppearTogether == null
                || this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel] == null) {
            return;
        }
        for (Boss boss : this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel]) {
            if (boss.id == BossID.ANDROID_15 && !boss.isDie()) {
                boss.changeToTypePK();
                break;
            }
        }
        this.parentBoss.changeToTypePK();
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (damage >= this.nPoint.hp) {
            boolean flag = true;
            if (this.parentBoss != null) {
                if (this.parentBoss.bossAppearTogether != null
                        && this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel] != null) {
                    for (Boss boss : this.parentBoss.bossAppearTogether[this.parentBoss.currentLevel]) {
                        if (boss.id == BossID.ANDROID_15 && !boss.isDie()) {
                            flag = false;
                            break;
                        }
                    }
                }
                if (flag && !this.parentBoss.isDie()) {
                    flag = false;
                }
            }
            if (!flag) {
                return 0;
            }
        }
        return super.injured(plAtt, damage, piercing, isMobAttack);
    }
}
