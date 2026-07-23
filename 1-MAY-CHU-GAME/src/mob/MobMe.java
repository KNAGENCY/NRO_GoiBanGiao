package mob;

/*
 *
 *
 * @author GiaCatLuong
 */
import consts.cn;
import map.Zone;
import player.Player;
import utils.SkillUtil;
import services.Service;
import utils.Util;
import network.Message;

public final class MobMe extends Mob {

    private Player player;
    private final long lastTimeSpawn;
    private final int timeSurvive;

    public MobMe(Player player) {
        super();
        this.player = player;
        this.id = (int) player.id;
        int level = player.playerSkill.getSkillbyId(12).point;
        this.tempId = SkillUtil.getTempMobMe(level);
        this.point.maxHp = SkillUtil.getHPMobMe(Util.maxIntValue(player.nPoint.hpMax), level);
        this.point.dame = getDameBySkillLevel();
        this.point.hp = this.point.maxHp;
        this.zone = player.zone;
        this.lastTimeSpawn = System.currentTimeMillis();
        this.timeSurvive = SkillUtil.getTimeSurviveMobMe(level);
        spawn();
    }

    @Override
    public void update() {
        if (Util.canDoWithTime(lastTimeSpawn, timeSurvive) && this.player.setClothes.pikkoroDaimao != 5) {
            this.mobMeDie();
            this.dispose();
        }
    }

    public void attack(Player pl, Mob mob, boolean miss) {
        Message msg;
        try {
            this.point.dame = getDameBySkillLevel();
            if (pl != null) {
                long dame = !miss ? this.point.dame : 0;
                if ((pl.nPoint.hp > dame && pl.nPoint.hp > pl.nPoint.hpMax * 0.05)
                        || this.player.setClothes.pikkoroDaimao == 5) {
                    long dameHit = pl.injured(this.player, dame, true, true);
                    msg = new Message(-95);
                    msg.writer().writeByte(2);
                    msg.writer().writeInt(this.id);
                    msg.writer().writeInt((int) pl.id);
                    msg.writeLongByGiaCatLuong(Util.maxIntValue(dameHit), cn.readInt);
                    msg.writeLongByGiaCatLuong(Util.maxIntValue(pl.nPoint.hp), cn.readInt);
                    Service.gI().sendMessAllPlayerInMap(this.player, msg);
                    msg.cleanup();
                }
            }

            if (mob != null) {
                long dameMob = !miss ? this.point.dame : 0;
                if (this.player != null && this.player.zone != null && this.player.zone.map.mapId == 155 && !this.player.isPet) {
                    dameMob = 1;
                }
                if (mob.point.gethp() > dameMob) {
                    long tnsm = mob.getTiemNangForPlayer(this.player, dameMob);
                    msg = new Message(-95);
                    msg.writer().writeByte(3);
                    msg.writer().writeInt(this.id);
                    msg.writer().writeInt((int) mob.id);
                    mob.point.sethp((mob.point.gethp() - dameMob));
                    msg.writeLongByGiaCatLuong(Util.maxIntValue(mob.point.gethp()), cn.readInt);
                    msg.writeLongByGiaCatLuong(Util.maxIntValue(dameMob), cn.readInt);
                    Service.gI().sendMessAllPlayerInMap(this.player, msg);
                    msg.cleanup();
                    Service.gI().addSMTN(player, (byte) 2, tnsm, true);
                }
            }
        } catch (Exception e) {
        }
    }

    // tạo mobme
    public void spawn() {
        Message msg;
        try {
            msg = new Message(-95);
            msg.writer().writeByte(0);// type
            msg.writer().writeInt((int) player.id);
            msg.writer().writeShort(this.tempId);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);// hp mob
            Service.gI().sendMessAllPlayerInMap(this.player, msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    public void spawnFor(Player target) {
        Message msg;
        try {
            if (target == null) {
                return;
            }
            msg = new Message(-95);
            msg.writer().writeByte(0);// type
            msg.writer().writeInt((int) player.id);
            msg.writer().writeShort(this.tempId);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);// hp mob
            target.sendMessage(msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    public void goToMap(Zone zone) {
        if (zone != null) {
            this.removeMobInMap();
            this.zone = zone;
        }
    }

    // xóa mobme khỏi map
    private void removeMobInMap() {
        Message msg;
        try {
            msg = new Message(-95);
            msg.writer().writeByte(7);// type
            msg.writer().writeInt((int) player.id);
            Service.gI().sendMessAllPlayerInMap(this.player, msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    public void mobMeDie() {
        Message msg;
        try {
            msg = new Message(-95);
            msg.writer().writeByte(6);// type
            msg.writer().writeInt((int) player.id);
            Service.gI().sendMessAllPlayerInMap(this.player, msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    @Override
    public synchronized void injured(Player plAtt, long damage, boolean dieWhenHpFull) {
        Message msg;
        try {
            if (damage > point.maxHp / 20) {
                damage = point.maxHp / 20;
            }
            point.hp -= damage;
            msg = new Message(-95);
            msg.writer().writeByte(5);// type
            msg.writer().writeInt((int) plAtt.id);
            msg.writer().writeByte(plAtt.playerSkill.skillSelect.template.id); // id skill
            msg.writer().writeInt(id); // mob id
            msg.writeLongByGiaCatLuong(Util.maxIntValue(damage), cn.readInt);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(point.hp), cn.readInt);
            Service.gI().sendMessAllPlayerInMap(this.player, msg);
            msg.cleanup();
            if (isDie()) {
                mobMeDie();
                dispose();
            }
        } catch (Exception e) {
        }
    }

    private int getSkillLevel() {
        try {
            return player.playerSkill.getSkillbyId(12).point;
        } catch (Exception e) {
            return 1;
        }
    }

    private long getDameBySkillLevel() {
        return SkillUtil.getHPMobMe(Util.maxIntValue(player.nPoint.getDameAttack(false)), getSkillLevel());
    }
    public void dispose() {
        player.mobMe = null;
        this.player = null;
    }
}
