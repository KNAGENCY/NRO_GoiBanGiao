package boss.boss_manifest.DaiTuongBroly;

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

public class DaiTuongBroly extends Boss {

    private long st;

    public DaiTuongBroly() throws Exception {
        super(BossID.DAITUONGBROLY, true, true, BossesData.DAITUONGBROLY);
    }

    @Override
    public void joinMap() {
        super.joinMap();
        st = System.currentTimeMillis();
    }

    @Override
public void reward(Player plKill) {
    this.pointBoss(plKill);
   

    // 50% rơi item 1284 hoặc 1285, mỗi loại 1 item, chỉ người hạ nhận
    if (Util.isTrue(50, 100)) {
        // item 1284
        ItemMap it1 = new ItemMap(this.zone, 1284, 3,
                this.location.x + Util.nextInt(-10, 10),
                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                plKill.id);
        Service.gI().dropItemMap(this.zone, it1);

        // item 1285
        ItemMap it2 = new ItemMap(this.zone, 1285, 3,
                this.location.x + Util.nextInt(-10, 10),
                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                plKill.id);
        Service.gI().dropItemMap(this.zone, it2);
    } 
    // 50% còn lại rơi item 1283, 30 item, mọi người nhận được
    else {
        for (int i = 0; i < 30; i++) {
            ItemMap it = new ItemMap(this.zone, 1283, 1,
                    this.location.x + Util.nextInt(-30, 30),
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    -1); // -1 nghĩa là tất cả người chơi đều nhặt được
            Service.gI().dropItemMap(this.zone, it);
        }
    }
}


    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (plAtt == null) {
            return 0;
        }
        if (!this.isDie()) {

            if (damage > 100_000_000) {
                damage = 100_000_000;
            }
            this.nPoint.subHP(damage);
            if (isDie()) {
                this.setDie(plAtt);
                die(plAtt);
            }
            return damage;
        } else {
            return 0;
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
