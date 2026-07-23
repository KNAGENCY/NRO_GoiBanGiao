package boss.boss_manifest.Cooler;

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
import services.EffectSkillService;
import services.ItemService;
import services.Service;
import services.TaskService;
import utils.Util;

public class Cooler extends Boss {

    private long st;

    public Cooler() throws Exception {
        super(BossID.COOLER, false, true, BossesData.COOLER, BossesData.COOLER_2);
    }

    @Override
    public void reward(Player plKill) {
        this.pointBoss(plKill);

        // ✅ Cập nhật nhiệm vụ giết boss
        TaskService.gI().checkDoneTaskKillBoss(plKill, this);

        if (Util.isTrue(20, 100)) {
            // 40%: rơi đồ TL, chỉ người hạ boss nhặt được
            ItemMap it = ItemService.gI().randDoTL(
                    this.zone,
                    1,
                    this.location.x,
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    plKill.id // chỉ người hạ boss nhặt
            );
            Service.gI().dropItemMap(this.zone, it);
        } else {
            // 60% còn lại: rơi 4 item id = 457, ai cũng nhặt được
            for (int i = 0; i < 2; i++) {
                int x = this.location.x + Util.nextInt(-60, 60); // rải xung quanh
                int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);
                ItemMap it = new ItemMap(this.zone, (short) 674, 1, x, y, -1);
                // -1 => không gán id người chơi => ai cũng nhặt được
                Service.gI().dropItemMap(this.zone, it);
            }
        }
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (!this.isDie()) {
            if (piercing) {
                damage /= 100;
            }
            if (Util.isTrue(200, 1000)) {
                this.chat("Xí hụt");
                return 0;
            }
            damage = this.nPoint.subDameInjureWithDeff(damage);
            if (effectSkill.isShielding) {
                if (damage > nPoint.hpMax) {
                    EffectSkillService.gI().breakShield(this);
                }
                damage = 1;
            }
            if (damage > 2_000_000_000) {
                damage = Util.nextInt(1_900_000_000, 2_000_000_000);
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
