package boss.boss_manifest.Cell;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossesData;
import consts.ConstPlayer;
import map.ItemMap;
import player.Player;
import services.EffectSkillService;
import services.ItemService;
import services.PlayerService;
import services.Service;
import services.TaskService;
import services.func.ChangeMapService;
import utils.Logger;
import utils.Util;

public class XenBoHung extends Boss {

    private long lastTimeHapThu;
    private int timeHapThu;

    public XenBoHung() throws Exception {
        super(BossID.XEN_BO_HUNG, false, true, BossesData.XEN_BO_HUNG_1, BossesData.XEN_BO_HUNG_2, BossesData.XEN_BO_HUNG_3);
    }

    @Override
    public void reward(Player plKill) {
        services.func.PhucLoiService.gI().recordBossKill(plKill, this);
        try {
            TaskService.gI().checkDoneTaskKillBoss(plKill, this);
        } catch (Exception e) {
            Logger.logException(XenBoHung.class, e, "XenBoHung reward task/point error");
        }

        try {
            if (this.zone == null || this.zone.map == null || this.location == null) {
                return;
            }

            if (Util.isTrue(20, 100)) {
                ItemMap it = ItemService.gI().randDoTL(
                        this.zone,
                        1,
                        this.location.x,
                        this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                        plKill != null ? plKill.id : -1);
                if (it != null) {
                    Service.gI().dropItemMap(this.zone, it);
                }
            } else {
                for (int i = 0; i < 2; i++) {
                    int x = this.location.x + Util.nextInt(-60, 60);
                    int y = this.zone.map.yPhysicInTop(x, this.location.y - 24);
                    ItemMap it = new ItemMap(this.zone, (short) 457, 1, x, y, -1);
                    Service.gI().dropItemMap(this.zone, it);
                }
            }
        } catch (Exception e) {
            Logger.logException(XenBoHung.class, e, "XenBoHung reward drop error");
        }
    }


    @Override
    public void active() {
        if (this.typePk == ConstPlayer.NON_PK) {
            this.changeToTypePK();
        }
        this.hapThu();
        this.attack();
    }

    private void hapThu() {
        if (!Util.canDoWithTime(this.lastTimeHapThu, this.timeHapThu) || !Util.isTrue(1, 100)) {
            return;
        }

        Player pl = this.zone.getRandomPlayerInMap();
        if (pl == null || pl.isDie()) {
            return;
        }
        ChangeMapService.gI().changeMapYardrat(this, this.zone, pl.location.x, pl.location.y);
        this.nPoint.dameg += (pl.nPoint.dame * 5 / 100);
        this.nPoint.hpg += (pl.nPoint.hp * 2 / 100);
        this.nPoint.critg++;
        this.nPoint.calPoint();
        PlayerService.gI().hoiPhuc(this, pl.nPoint.hp, 0);
        pl.injured(null, pl.nPoint.hpMax, true, false);
        Service.gI().sendThongBao(pl, "Bạn vừa bị " + this.name + " hấp thu!");
        this.chat(2, "Ui cha cha, kinh dị quá. " + pl.name + " vừa bị tên " + this.name + " nuốt chửng kìa!!!");
        this.chat("Haha, ngọt lắm đấy " + pl.name + "..");
        this.lastTimeHapThu = System.currentTimeMillis();
        this.timeHapThu = Util.nextInt(10000, 20000);
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (!this.isDie()) {
            if (!piercing && Util.isTrue(this.nPoint.tlNeDon, 1000)) {
                this.chat("Xí hụt");
                return 0;
            }
            damage = this.nPoint.subDameInjureWithDeff(damage / 2);
            if (!piercing && effectSkill.isShielding) {
                if (damage > nPoint.hpMax) {
                    EffectSkillService.gI().breakShield(this);
                }
                damage = damage / 4;
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

}
