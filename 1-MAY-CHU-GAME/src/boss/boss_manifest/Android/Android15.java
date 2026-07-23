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
import services.PlayerService;
import services.Service;
import services.TaskService;
import utils.Util;

public class Android15 extends Boss {

    public boolean callApk13;

    public Android15() throws Exception {
        super(BossID.ANDROID_15, false, true, BossesData.ANDROID_15);
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
        this.attack();
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (!this.callApk13 && damage >= this.nPoint.hp) {
            if (this.parentBoss != null) {
                ((Android14) this.parentBoss).callApk13();
            }
            return 0;
        }
        return super.injured(plAtt, damage, piercing, isMobAttack);
    }

    public void recoverHP() {
        PlayerService.gI().hoiPhuc(this, this.nPoint.hpMax, 0);
    }
}
