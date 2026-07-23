package boss.boss_manifest.Black;

/*
 *
 *
 * @author GiaCatLuong
 */
import boss.*;
import consts.ConstPlayer;
import consts.ConstTask;

import map.ItemMap;
import player.Player;
import server.Manager;
import services.*;
import utils.Util;

import java.util.Random;

import utils.SkillUtil;

public class BA extends Boss {

    private long st;
    private int timeLeaveMap;

    public BA() throws Exception {
        super(BossID.BA, false, true, BossesData.BA, BossesData.SBA);
    }

  @Override
public void reward(Player plKill) {
    int rate = Util.nextInt(1, 100);

    if (rate <= 20) {
        // 20% rơi id 549 x2, rải đều xung quanh
        for (int i = 0; i < 2; i++) {
            int x = this.location.x + Util.nextInt(-60, 60);
            int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);
            ItemMap it = new ItemMap(this.zone, (short) 549, 1, x, y, -1);
            Service.gI().dropItemMap(this.zone, it);
        }

    } else if (rate <= 40) {
        // 20% rơi id 1628 x2, rải đều xung quanh
        for (int i = 0; i < 2; i++) {
            int x = this.location.x + Util.nextInt(-60, 60);
            int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);
            ItemMap it = new ItemMap(this.zone, (short) 1628, 1, x, y, -1);
            Service.gI().dropItemMap(this.zone, it);
        }

    } else {
        // 60% còn lại: rơi 457 x5, ai cũng nhặt được
        for (int i = 0; i < 5; i++) {
            int x = this.location.x + Util.nextInt(-60, 60);
            int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);
            ItemMap it = new ItemMap(this.zone, (short) 457, 1, x, y, -1);
            Service.gI().dropItemMap(this.zone, it);
        }
    }
}
   @Override
public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
    if (!this.isDie()) {
        if (!piercing && Util.isTrue(this.nPoint.tlNeDon, 1000)) {
            this.chat("Xí hụt");
            return 0;
        }

        if (this.currentLevel != 0) {
            damage /= 2;
        }

        damage = this.nPoint.subDameInjureWithDeff(damage - Util.nextInt(100000));

        // Giới hạn dame pem boss: tối đa 5000 / 1 hit
        if (damage > 5000) {
            damage = 5000;
        }
        if (damage < 1) {
            damage = 1;
        }

        if (!piercing && effectSkill.isShielding) {
            if (damage > nPoint.hpMax) {
                EffectSkillService.gI().breakShield(this);
            }
            damage = 1;
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
        if (Util.canDoWithTime(st, timeLeaveMap)) {
            if (Util.isTrue(1, 2)) {
                this.leaveMap();
            } else {
                this.leaveMapNew();
            }
        }
        if (this.zone != null && this.zone.getNumOfPlayers() > 0) {
            st = System.currentTimeMillis();
            timeLeaveMap = Util.nextInt(300000, 900000);
        }
    }

    @Override
    public void joinMap() {
        this.name = this.data[this.currentLevel].getName() + " " + Util.nextInt(1, 100);
        super.joinMap();
        st = System.currentTimeMillis();
        timeLeaveMap = Util.nextInt(600000, 900000);
    }

    @Override
    public void attack() {
        if (Util.canDoWithTime(this.lastTimeAttack, 100) && this.typePk == ConstPlayer.PK_ALL) {
            this.lastTimeAttack = System.currentTimeMillis();
            try {
                Player pl = getPlayerAttack();
                if (pl == null || pl.isDie()) {
                    return;
                }
                this.playerSkill.skillSelect = this.playerSkill.skills
                        .get(Util.nextInt(0, this.playerSkill.skills.size() - 1));
                int dis = Util.getDistance(this, pl);
                if (dis > 450) {
                    move(pl.location.x - 24, pl.location.y);
                } else if (dis > 100) {
                    int dir = (this.location.x - pl.location.x < 0 ? 1 : -1);
                    int move = Util.nextInt(50, 100);
                    move(this.location.x + (dir == 1 ? move : -move), pl.location.y);
                } else {
                    if (Util.isTrue(30, 100)) {
                        int move = Util.nextInt(50);
                        move(pl.location.x + (Util.nextInt(0, 1) == 1 ? move : -move), this.location.y);
                    }
                    SkillService.gI().useSkill(this, pl, null, -1, null);
                    checkPlayerDie(pl);
                }
            } catch (Exception ex) {
                ex.printStackTrace();
            }
        }
    }
}
