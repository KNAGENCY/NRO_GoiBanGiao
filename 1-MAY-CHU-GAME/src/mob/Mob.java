package mob;

/*
 *
 *
 * @author GiaCatLuong
 */
import consts.ConstItem;
import consts.ConstMap;
import consts.ConstMob;
import consts.ConstTask;
import consts.cn;
import item.Item;
import item.Item.ItemOption;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import map.ItemMap;
import map.Zone;
import models.Achievement.AchievementService;
import models.Training.TrainingService;
import network.Message;
import player.Location;
import player.Pet;
import player.Player;
import server.Maintenance;
import server.Manager;
import server.ServerNotify;
import services.InventoryService;
import services.ItemMapService;
import services.ItemService;
import services.MapService;
import services.Service;
import services.TaskService;
import skill.Skill;
import utils.TimeUtil;
import utils.Util;

public class Mob {

    public int id;
    public Zone zone;
    public int tempId;
    public String name;
    public byte level;

    public List<Player> temporaryEnemies = new ArrayList<>();

    public MobPoint point;
    public MobEffectSkill effectSkill;
    public Location location;

    public byte pDame;
    public int pTiemNang;
    private long maxTiemNang;

    public long lastTimeDie;
    public int lvMob = 0;
    public int status = 5;
    public int type = 1;

    private long lastTimeAttackPlayer;
    private long timeAttack = 2000;
    public long lastTimePhucHoi = System.currentTimeMillis();
    public long lastTimeSendEffect = System.currentTimeMillis();

    public Mob(Mob mob) {
        this.point = new MobPoint(this);
        this.effectSkill = new MobEffectSkill(this);
        this.location = new Location();
        this.id = mob.id;
        this.tempId = mob.tempId;
        this.level = mob.level;
        this.point.setHpFull(mob.point.getHpFull());
        this.point.sethp((this.point.getHpFull()));
        this.location.x = mob.location.x;
        this.location.y = mob.location.y;
        this.pDame = mob.pDame;
        this.pTiemNang = mob.pTiemNang;
        this.type = mob.type;
        this.setTiemNang();
    }

    private void sendMobAttackMe(Player player, int dame) {
        if (!player.isPet && !player.isNewPet && !player.isBot) {
            Message msg;
            try {// =)) khong biet fix kkk
                msg = new Message(-11);
                msg.writer().writeByte(this.id);
                msg.writer().writeInt(dame); // dame
                player.sendMessage(msg);
                msg.cleanup();
            } catch (Exception e) {
            }
        }
    }

    private void hutItem(Player player, List<ItemMap> items) {
        if (!player.isPet && !player.isNewPet && !player.isBot) {
            if (player.charms.tdThuHut > System.currentTimeMillis()) {
                for (ItemMap item : items) {
                    ItemMapService.gI().pickItem(player, item.itemMapId, true);
                }
            }
        } else {
            if (((Pet) player).master.charms.tdThuHut > System.currentTimeMillis()) {
                for (ItemMap item : items) {
                    ItemMapService.gI().pickItem(((Pet) player).master, item.itemMapId, true);
                }
            }
        }
    }

    public Mob() {
        this.point = new MobPoint(this);
        this.effectSkill = new MobEffectSkill(this);
        this.location = new Location();
    }

    public void setTiemNang() {
        this.maxTiemNang = (int) (long) this.point.getHpFull() * (long) (this.pTiemNang + Util.nextInt(-2, 2)) / 100L;
    }

    public boolean isDie() {
        return this.point.gethp() <= 0;
    }

    public void setDie() {
        this.lastTimePhucHoi = System.currentTimeMillis();
        this.lastTimeDie = System.currentTimeMillis();
    }

    public void addTemporaryEnemies(Player pl) {
        if (pl != null && !temporaryEnemies.contains(pl)) {
            temporaryEnemies.add(pl);
        }
    }

    public void injured(Player plAtt, long damage, boolean dieWhenHpFull) {
        if (!this.isDie()) {
            if (damage >= this.point.hp) {
                damage = this.point.hp;
            }
            if (!dieWhenHpFull) {
                if (this.point.hp == this.point.maxHp && damage >= this.point.hp) {
                    damage = this.point.hp - 1;
                }
                if ((this.tempId == ConstMob.MOC_NHAN || this.tempId == ConstMob.BU_NHIN_MA_QUAI)
                        && damage > this.point.maxHp / 10) {
                    damage = this.point.maxHp / 10;
                }
            }
            if (MapService.gI().isMapKhiGasHuyDiet(this.zone.map.mapId)) {
                boolean mob76Die = true;
                for (Mob mob : this.zone.mobs) {
                    if (!mob.isDie() && mob.tempId == ConstMob.CO_MAY_HUY_DIET) {
                        mob76Die = false;
                        break;
                    }
                }
                if (!mob76Die && plAtt != null && plAtt.playerSkill != null && plAtt.playerSkill.skillSelect != null) {
                    switch (plAtt.playerSkill.skillSelect.template.id) {
                        case Skill.LIEN_HOAN, Skill.ANTOMIC, Skill.MASENKO, Skill.KAMEJOKO ->
                            damage = 1;
                    }
                }
            }
            // if (!dieWhenHpFull && !isBigBoss() &&
            // !MapService.gI().isMapPhoBan(this.zone.map.mapId) && this.lvMob > 0
            // && plAtt != null && plAtt.charms.tdOaiHung < System.currentTimeMillis()) {
            // damage = (int) ((this.point.maxHp <= 20000000 ? this.point.maxHp * 10 :
            // 2000000000) * (10.0 / 100));
            // this.mobAttackPlayer(plAtt);
            // }
            if (plAtt != null && plAtt.isBoss && this.tempId > 0 && Util.isTrue(1, 2)
                    && Util.canDoWithTime(lastTimeAttackPlayer, 2500)) {
                this.mobAttackPlayer(plAtt);
                lastTimeAttackPlayer = System.currentTimeMillis();
            }
            if (plAtt != null && (plAtt.zone.map.mapId == 156 || plAtt.zone.map.mapId == 157) && !plAtt.isPet) {
                damage = 1;
            }
            if (plAtt != null && plAtt.zone != null && plAtt.zone.map.mapId == 155 && !plAtt.isPet) {
                damage = 1;
            }
            if (plAtt != null && plAtt.zone != null && plAtt.zone.map != null
                    && MapService.gI().isMapNguHanhSon(plAtt.zone.map.mapId)) {
                damage = 80_000;
            }

            // if (damage > 2_000_000_000) {
            // damage = 2_000_000_000;
            // }
            this.point.hp -= damage;
            addTemporaryEnemies(plAtt);
            if (this.isDie()) {
                this.status = 0;
                this.setDie();
                this.temporaryEnemies.clear();
                if (plAtt != null) {
                    this.sendMobDieAffterAttacked(plAtt, (int) damage);
                    TaskService.gI().checkDoneTaskKillMob(plAtt, this);
                    TaskService.gI().checkDoneSideTaskKillMob(plAtt, this);
                    TaskService.gI().checkDoneClanTaskKillMob(plAtt, this);
                    AchievementService.gI().checkDoneTaskKillMob(plAtt, this);
                }
                if (this.id == 13) {
                    this.zone.isbulon1Alive = false;
                }
                if (this.id == 14) {
                    this.zone.isbulon2Alive = false;
                }
            } else {
                this.sendMobStillAliveAffterAttacked(damage,
                        plAtt != null ? (plAtt.nPoint != null && plAtt.nPoint.isCrit) : false);
            }
            if (plAtt != null) {
                if (plAtt.isPl() && plAtt.satellite != null && plAtt.satellite.isDefend) {
                    plAtt.satellite.isDefend = false;
                }
                Player tnsmPlayer = plAtt;
                long tnsmDamage = damage;
                long tnsm = isPrisonPetAttack(tnsmPlayer) ? Math.max(0, tnsmDamage)
                        : getTiemNangForPlayer(tnsmPlayer, tnsmDamage);
                if (tnsm > 0) {
                    Service.gI().addSMTN(tnsmPlayer, (byte) 2, tnsm, true);
                    TrainingService.gI().tangTnsmLuyenTap(tnsmPlayer, tnsm);
                }
            }
        }
    }

    public long getTiemNangForPlayer(Player pl, long dame) {
        if (pl.zone != null && pl.zone.map != null && MapService.gI().isMapNguHanhSon(pl.zone.map.mapId)) {
            return Util.nextInt(5_000, 30_000);
        }
        int levelPlayer = Service.gI().getCurrLevel(pl);
        int n = levelPlayer - this.level;
        if (pl.zone != null && MapService.gI().isMapBanDoKhoBau(pl.zone.map.mapId)) {
            n = 0;
        }
        if (pl.nPoint != null && pl.nPoint.power < 10_000_000_000L) {
            n = 0;
        }
        long hpFull = Math.max(1, point.getHpFull());
        long tiemNang = Math.abs(dame) * Math.max(1, maxTiemNang) / hpFull;
        tiemNang = Math.max(tiemNang, Math.abs(dame) / 1000);
        if (tiemNang <= 0) {
            tiemNang = 1;
        }
        long randomRange = Math.max(1, tiemNang * 15 / 100);
        tiemNang = Util.nextLong(Math.max(1, tiemNang - randomRange), tiemNang + randomRange);
        if (tiemNang >= 1000000000) {
            tiemNang = 999999999;
        }
        if (n >= 0) {
            for (int i = 0; i < n; i++) {
                long sub = tiemNang * 10 / 100;
                if (sub <= 0) {
                    sub = 1;
                }
                tiemNang -= sub;
            }
        } else {
            for (int i = 0; i < -n; i++) {
                long add = tiemNang * 10 / 100;
                if (add <= 0) {
                    add = 1;
                }
                tiemNang += add;
            }
        }
        if (tiemNang <= 0) {
            tiemNang = 1;
        }
        // Map đặc biệt luôn nhận 100 TNSM
if (pl.zone.map.mapId == 80) {
    tiemNang = 100;
}
           if (pl.zone.map.mapId == 135) {
        long dmg = Math.abs(dame);
        if (dmg >= 60000 && dmg <= 80000) {
            tiemNang = Util.nextInt(300000, 450000);
        } else {
            tiemNang = 1;
        }
    }
        if (pl.nPoint != null) {
            tiemNang = pl.nPoint.calSucManhTiemNang(tiemNang);
        } else {
            return 0;
        }
        if (pl.zone.map.mapId == 122 || pl.zone.map.mapId == 123 || pl.zone.map.mapId == 124) {
            tiemNang *= 0.2;
        }
        if (pl.zone.map.mapId == 77) {
            tiemNang *= 1.5;
        }
     
        return tiemNang;
    }

    private boolean isPrisonPetAttack(Player player) {
        return player != null && player.isPet && player.zone != null && player.zone.map != null
                && player.zone.map.mapId == 155;
    }

    public void update() {
        if (zone.isGoldenFriezaAlive && TimeUtil.is21H()) {
            if (!isDie()) {
                startDie();
                return;
            }
        }
        if (!this.isDie() && this.tempId == ConstMob.CO_MAY_HUY_DIET && Util.canDoWithTime(lastTimeSendEffect, 1000)) {
            sendEffect(55);
            lastTimeSendEffect = System.currentTimeMillis();
        }

        if (this.isDie() && !Maintenance.isRunning && !isBigBoss()) {
            switch (zone.map.type) {
                case ConstMap.MAP_DOANH_TRAI:
                    if (this.tempId == ConstMob.BULON && this.zone.isTUTAlive
                            && Util.canDoWithTime(lastTimeDie, 10000)) {
                        this.hoiSinh();
                        this.hoiSinhMobPhoBan();
                        if (this.id == 13) {
                            this.zone.isbulon1Alive = true;
                        }
                        if (this.id == 14) {
                            this.zone.isbulon2Alive = true;
                        }
                    }
                    break;
                case ConstMap.MAP_BAN_DO_KHO_BAU:
                    break;
                case ConstMap.MAP_CON_DUONG_RAN_DOC:
                    break;
                case ConstMap.MAP_KHI_GAS_HUY_DIET:
                    break;
                case ConstMap.MAP_TAY_KARIN:
                    break;
                default:
                    if (this.zone.isGoldenFriezaAlive && TimeUtil.is21H()) {
                        return;
                    }
                    if (Util.canDoWithTime(lastTimeDie, 5000)) {
                        this.hoiSinh();
                        this.sendMobHoiSinh();
                    }
                    if (Util.canDoWithTime(lastTimePhucHoi, 30000) && !isDie()) {
                        lastTimePhucHoi = System.currentTimeMillis();
                        long hpMax = this.point.maxHp;
                        if (this.point.hp < hpMax) {
                            hoi_hp(hpMax / 10);
                        } else {
                            this.sendMobHoiSinh();
                        }
                    }
            }
        }

        effectSkill.update();
        attack();
    }

    public boolean isBigBoss() {
        return (this.tempId == ConstMob.HIRUDEGARN
                || this.tempId == ConstMob.VUA_BACH_TUOC
                || this.tempId == ConstMob.ROBOT_BAO_VE
                || this.tempId == ConstMob.GAU_TUONG_CUOP
                || this.tempId == ConstMob.VOI_CHIN_NGA
                || this.tempId == ConstMob.GA_CHIN_CUA
                || this.tempId == ConstMob.NGUA_CHIN_LMAO
                || this.tempId == ConstMob.PIANO
                || this.tempId == ConstMob.KONG
                || this.tempId == ConstMob.GOZILLA);
    }

    public void attack() {
        Player player = getPlayerCanAttack();
        if (!isDie() && !effectSkill.isHaveEffectSkill() && tempId != ConstMob.MOC_NHAN
                && tempId != ConstMob.BU_NHIN_MA_QUAI && tempId != ConstMob.CO_MAY_HUY_DIET && !this.isBigBoss()
                && (this.lvMob < 1 || MapService.gI().isMapPhoBan(this.zone.map.mapId))
                && Util.canDoWithTime(lastTimeAttackPlayer, timeAttack)) {
            if (player != null) {
                this.mobAttackPlayer(player);
            }
            this.lastTimeAttackPlayer = System.currentTimeMillis();
        }
    }

    public Player getPlayerCanAttack() {
        Player plAttack = getFirstPlayerCanAttack();
        if (plAttack != null) {
            return plAttack;
        }
        int distance = 100;
        try {
            List<Player> players = this.zone.getNotBosses();
            for (Player pl : players) {
                if (!pl.isDie() && !pl.isBoss && !pl.isNewPet && (pl.satellite == null || !pl.satellite.isDefend)
                        && (pl.effectSkin == null || !pl.effectSkin.isVoHinh)
                        && (this.tempId > 18 || (this.tempId > 9 && this.type == 4)) || isBigBoss()) {
                    int dis = Util.getDistance(pl, this);
                    if (dis <= distance || isBigBoss()) {
                        plAttack = pl;
                        distance = dis;
                    }
                }
            }
            this.timeAttack = 2000;
        } catch (Exception e) {

        }
        return plAttack;
    }

    private Player getFirstPlayerCanAttack() {
        Player plAtt = null;
        try {
            List<Player> playersMap = zone.getHumanoids();
            int dis = 300;
            if (playersMap != null) {
                for (Player plAttt : playersMap) {
                    if (plAttt.isDie() || plAttt.isBoss || (plAttt.satellite != null && plAttt.satellite.isDefend)
                            || (plAttt.effectSkin != null && plAttt.effectSkin.isVoHinh)
                            || !this.temporaryEnemies.contains(plAttt)) {
                        continue;
                    }
                    int d = Util.getDistance(plAttt, this);
                    if (d <= dis) {
                        dis = d;
                        plAtt = plAttt;
                    }
                }
            }
            this.timeAttack = 1000;
        } catch (Exception e) {

        }
        return plAtt;
    }

    private void mobAttackPlayer(Player player) {
        long dameMob = Util.maxIntValue(this.point.getDameAttack());
        if (player.effectSkill != null && player.effectSkill.isShielding) {
            dameMob = 0; // Nếu đang sử dụng khiên, sát thương là 0
        } else {
            if (player.charms != null && player.charms.tdDaTrau > System.currentTimeMillis()) {
                dameMob /= 2;
            }
            if (player.isPet && ((Pet) player).master.charms != null
                    && ((Pet) player).master.charms.tdDeTu > System.currentTimeMillis()) {
                dameMob /= 2;
            }
            if (this.lvMob > 0 && !MapService.gI().isMapPhoBan(this.zone.map.mapId)) {
                dameMob = (long) (player.nPoint.hpMax * (10.0 / 100));
            }
            if (player.satellite != null && player.satellite.isDefend) {
                dameMob -= dameMob / 5;
            }
            if (player.itemTime != null && player.itemTime.isUseCMS) {
                dameMob = (long) Math.round(dameMob * 0.1);
            }
            if (this.lvMob > 0 && player.charms.tdOaiHung > System.currentTimeMillis()) {
                dameMob = 0;
            }
        }
        long dame = player.injured(null, Util.maxIntValue(dameMob), false, true);

        this.sendMobAttackMe(player, dame);
        this.sendMobAttackPlayer(player);
        this.phanSatThuong(player, dame);
    }

    private void sendMobAttackMe(Player player, long dame) {
        if (!player.isPet && !player.isNewPet) {
            Message msg;
            try {
                msg = new Message(-11);
                msg.writer().writeByte(this.id);
                msg.writeLongByGiaCatLuong(Util.maxIntValue(dame), cn.readInt); // dame
                player.sendMessage(msg);
                msg.cleanup();
            } catch (Exception e) {
            }
        }
    }

    private void sendMobAttackPlayer(Player player) {
        Message msg;
        try {
            msg = new Message(-10);
            msg.writer().writeByte(this.id);
            msg.writer().writeInt((int) player.id);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(player.nPoint.hp), cn.readInt); // hp
            Service.gI().sendMessAnotherNotMeInMap(player, msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    public void hoiSinh() {
        this.status = 5;
        this.point.hp = this.point.maxHp;
        this.setTiemNang();
    }

    public int lvMob() {
        this.lvMob = 0;
        this.point.hp = this.point.maxHp;
        return 0;
    }

    public void sendMobHoiSinh() {
        Message msg = null;
        try {
            msg = new Message(-13);
            msg.writer().writeByte(this.id);
            msg.writer().writeByte(this.tempId);
            msg.writer().writeByte(lvMob());
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            this.sendMobMaxHp(this.point.hp);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    public void hoi_hp(long hp) {
        Message msg = null;
        try {
            this.point.sethp((this.point.gethp() + hp));
            long HP = hp > 0 ? 1 : Math.abs(hp);
            msg = new Message(-9);
            msg.writer().writeByte(this.id);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.gethp()), cn.readInt);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(HP), cn.readInt);
            msg.writer().writeBoolean(false);
            msg.writer().writeByte(-1);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
                msg = null;
            }
        }
    }

    public void sendEffect(int Effect) {
        Message msg = null;
        try {
            msg = new Message(-9);
            msg.writer().writeByte(this.id);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.gethp()), cn.readInt);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.gethp()), cn.readInt);
            msg.writer().writeBoolean(false);
            msg.writer().writeByte(Effect);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
                msg = null;
            }
        }
    }

    private void sendMobDieAffterAttacked(Player plKill, long dameHit) {
        Message msg;
        try {
            msg = new Message(-12);
            msg.writer().writeByte(this.id);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(dameHit), cn.readInt);
            msg.writer().writeBoolean(plKill.nPoint.isCrit); // crit
            List<ItemMap> items = mobReward(plKill, this.dropItemTask(plKill), msg);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
            hutItem(plKill, items);
        } catch (Exception e) {
        }
    }

    private List<ItemMap> mobReward(Player player, ItemMap itemTask, Message msg) {
        if (player.isBot) {
            return null;
        }
        List<ItemMap> itemReward = new ArrayList<>();
        try {
            itemReward = this.getItemMobReward(player, this.location.x + Util.nextInt(-10, 10),
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y));
            if (itemTask != null) {
                itemReward.add(itemTask);
            }
            msg.writer().writeByte(itemReward.size()); // sl item roi
            for (ItemMap itemMap : itemReward) {
                msg.writer().writeShort(itemMap.itemMapId);// itemmapid
                msg.writer().writeShort(itemMap.itemTemplate.id); // id item
                msg.writer().writeShort(itemMap.x); // xend item
                msg.writer().writeShort(itemMap.y); // yend item
                msg.writer().writeInt((int) itemMap.playerId); // id nhan vat
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return itemReward;
    }

    private int adjustMobDropRatio(int baseRatio) {
        int rate = Math.max(1, Manager.RATE_DROP_SERVER);
        return Math.max(1, (int) Math.round(baseRatio * 5.0 / rate));
    }

    public List<ItemMap> getItemMobReward(Player player, int x, int yEnd) {
        List<ItemMap> list = new ArrayList<>();

        if (player.isBoss) {
            return list;
        }

        // if (player.isPl() && Util.isTrue(1, 10000) && this.tempId == 0) {
        // short itTemp = (short) ItemService.gI().randTempItemKichHoat(player.gender);
        // ItemMap it = new ItemMap(zone, itTemp, 1, x, yEnd, player.id);
        // List<Item.ItemOption> ops = ItemService.gI().getListOptionItemShop(itTemp);
        // if (!ops.isEmpty()) {
        // it.options = ops;
        // }
        // it.options.add(new Item.ItemOption(210, 0));
        // it.options.add(new Item.ItemOption(216, 0));
        // it.options.add(new Item.ItemOption(30, 0));
        // list.add(it);
        // }
        if (this.tempId == 0) {
            return list;
        }
        int mapid = player.zone.map.mapId;

        // if (EventManager.CHRISTMAS) {
        // Player pl = player;
        // if (pl.isPet) {
        // pl = ((Pet) pl).master;
        // }
        // if (pl.isPet) {
        // pl = ((Pet) pl).master;
        // }
        // if (Util.isTrue(1, adjustMobDropRatio(100))) {
        // if (pl.itemEvent != null && pl.itemEvent.canDropTatVoGiangSinh(100)) {
        // list.add(new ItemMap(zone, 649, 1, x, yEnd, player.id));
        // }
        // }
        // }

        // if (mapid == 5 || mapid == 13) {
        // Player pl = player;
        // if (pl.isPet) {
        // pl = ((Pet) pl).master;
        // }
        // if (pl.isPet) {
        // pl = ((Pet) pl).master;
        // }
        // if (Util.isTrue(1, adjustMobDropRatio(500))) {
        // if (pl.itemEvent != null && pl.itemEvent.canDropBinhNuoc(100)) {
        // list.add(new ItemMap(zone, 456, 1, x, yEnd, pl.id));
        // }
        // }
        // }

        // if (EventManager.INTERNATIONAL_WOMANS_DAY) {
        // Player pl = player;
        // if (pl.isPet) {
        // pl = ((Pet) pl).master;
        // }
        // if (pl.isPet) {
        // pl = ((Pet) pl).master;
        // }
        // if (Util.isTrue(1, adjustMobDropRatio(50))) {
        // if (pl.itemEvent != null && pl.itemEvent.canDropHoaHong(100)) {
        // list.add(new ItemMap(zone, 610, 1, x, yEnd, player.id));
        // }
        // }
        // }
        // if (Util.isTrue(15, 300)) {
        // list.add(new ItemMap(zone, 18, 1, x, yEnd, player.id));
        // }
        // if (Util.isTrue(15, 300)) {
        // list.add(new ItemMap(zone, 19, 1, x, yEnd, player.id));
        // }
        // if (Util.isTrue(15, 300)) {
        // list.add(new ItemMap(zone, 20, 1, x, yEnd, player.id));
        // }

        // if (Util.isTrue(1, adjustMobDropRatio(400))) {
        // ItemMap it = new ItemMap(zone, 1530, 1, x, yEnd, player.id);
        // it.options.add(new Item.ItemOption(30, 0));
        // list.add(it);
        //
        // }

        // if (EventManager.HALLOWEEN) {
        // if (MapService.gI().isMapEventHalloween(mapid)) {
        // if (Util.isTrue(1, adjustMobDropRatio(50))) {
        // list.add(new ItemMap(zone, 707, 1, x, yEnd, player.id));
        // } else if (Util.isTrue(1, adjustMobDropRatio(50))) {
        // list.add(new ItemMap(zone, 708, 1, x, yEnd, player.id));
        // }
        // }
        // }

        // Rơi đồ đặc biệt cho Namek khi đánh quái dưới 53
        // Rơi đồ đặc biệt cho từng tộc khi đánh quái dưới 53
        // if (zone.map.mapId == 5 && player.nPoint.power > 10_000_000_000L) {
        //     // Nếu trên 40 tỷ thì tỉ lệ 1/200, còn lại là 1/1000
        //     int rate = player.nPoint.power > 20_000_000_000L ? 100 : 50;

        //     if (Util.isTrue(1, rate)) {
        //         int randomItem = -1;

        //         if (player.gender == 0) { // Trái đất
        //             int[] itemIdsEarth = {488, 489, 490, 491, 492, 493, 494, 495, 496, 497, 498, 499, 500, 501};
        //             randomItem = itemIdsEarth[Util.nextInt(0, itemIdsEarth.length - 1)];
        //         } else if (player.gender == 1) { // Namek
        //             int[] itemIdsNamek = { 481, 482, 483, 484, 485, 486, 487, 474, 475, 476, 477, 478, 479, 480 };
        //             randomItem = itemIdsNamek[Util.nextInt(0, itemIdsNamek.length - 1)];
        //         } else if (player.gender == 2) { // Xayda
        //             int[] itemIdsSaiyan = { 502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 512, 513, 514, 515 };
        //             randomItem = itemIdsSaiyan[Util.nextInt(0, itemIdsSaiyan.length - 1)];
        //         }

        //         if (randomItem != -1) {
        //             ItemMap item = new ItemMap(zone, randomItem, 1, x, yEnd, player.id);
        //             Service.gI().dropItemMap(zone, item);
        //         }
        //     }
        // }

        // if (this.tempId > 1 && this.tempId < 54) {
        //     int ratio = player.nPoint.power < 2_000_000_000L ? 100 : 200; // dưới 1tỷ 1/200, trên 1tỷ 1/500
        //     if (Util.isTrue(1, adjustMobDropRatio(ratio))) {
        //         int randomItem = -1;
        //         if (player.gender == 0) { // Trái đất
        //             int[] itemIdsEarth = { 69, 70, 71, 116, 117, 118, 119, 120, 308, 309, 310, 311, 312 };
        //             randomItem = itemIdsEarth[Util.nextInt(0, itemIdsEarth.length - 1)];
        //         } else if (player.gender == 1) { // Namek
        //             int[] itemIdsNamek = { 82, 83, 84, 123, 124, 125, 126, 127, 336, 337, 338, 339, 340 };
        //             randomItem = itemIdsNamek[Util.nextInt(0, itemIdsNamek.length - 1)];
        //         } else if (player.gender == 2) { // Xayda
        //             int[] itemIdsSaiyan = { 90, 91, 92, 130, 131, 132, 133, 134, 315, 316, 317, 318, 319 };
        //             randomItem = itemIdsSaiyan[Util.nextInt(0, itemIdsSaiyan.length - 1)];
        //         }

        //         if (randomItem != -1) {
        //             list.add(new ItemMap(zone, randomItem, 1, x, yEnd, player.id));
        //         }
        //     }
        // }

        // if (this.tempId > 1 && this.tempId < 58) {
        //     int ratio = player.nPoint.power < 2_000_000_000L ? 100 : 200; // dưới 10tỷ 1/500, trên 10tỷ 1/1000
        //     if (Util.isTrue(1, adjustMobDropRatio(ratio))) {
        //         int randomItem = -1;
        //         if (player.gender == 0) { // Trái đất
        //             int[] itemIdsEarth = { 96, 97, 98, 99, 302, 303, 304, 305 };
        //             randomItem = itemIdsEarth[Util.nextInt(0, itemIdsEarth.length - 1)];
        //         } else if (player.gender == 1) { // Namek
        //             int[] itemIdsNamek = { 103, 104, 105, 106, 330, 331, 332, 333 };
        //             randomItem = itemIdsNamek[Util.nextInt(0, itemIdsNamek.length - 1)];
        //         } else if (player.gender == 2) { // Xayda
        //             int[] itemIdsSaiyan = { 110, 111, 112, 113, 323, 324, 325, 326 };
        //             randomItem = itemIdsSaiyan[Util.nextInt(0, itemIdsSaiyan.length - 1)];
        //         }

        //         if (randomItem != -1) {
        //             list.add(new ItemMap(zone, randomItem, 1, x, yEnd, player.id));
        //         }
        //     }
        // }

        if (this.tempId > 38 && this.tempId < 54) {
            int ratio = player.nPoint.power < 2_000_000_000L ? 600 : 900; // dưới 10tỷ 1/1000, trên 10tỷ 1/2000
            if (Util.isTrue(1, adjustMobDropRatio(ratio))) {
                // code xử lý rơi đồ ở đây
                int[] dropItems = { 241, 253, 265, 277, 233, 245, 257, 269, 237, 249, 261, 273, 281 };
                int randomItem = dropItems[Util.nextInt(0, dropItems.length - 1)];

                Item item = ItemService.gI().createNewItem((short) randomItem);

                // Gắn option 210 theo yêu cầu
                int[] op210Choices = { 50, 77, 103, 47 };
                int op210 = op210Choices[Util.nextInt(0, op210Choices.length - 1)];
                int value210 = (op210 == 47) ? Util.nextInt(100, 400) : Util.nextInt(1, 4);
                item.itemOptions.add(new ItemOption(op210, value210));

                // Gắn option theo từng nhóm
                switch (randomItem) {
                    case 241:
                    case 233:
                    case 237: // nhóm 1
                        item.itemOptions.add(new ItemOption(47, Util.nextInt(400, 550)));
                        break;
                    case 253:
                    case 245:
                    case 249: // nhóm 2
                        item.itemOptions.add(new ItemOption(6, Util.nextInt(22000, 27000)));
                        item.itemOptions.add(new ItemOption(27, Util.nextInt(3000, 5000)));
                        break;
                    case 265:
                    case 261:
                    case 257: // nhóm 3
                        item.itemOptions.add(new ItemOption(0, Util.nextInt(2100, 2400)));
                        break;
                    case 277:
                    case 269:
                    case 273: // nhóm 4
                        item.itemOptions.add(new ItemOption(7, Util.nextInt(22000, 26000)));
                        item.itemOptions.add(new ItemOption(28, Util.nextInt(4000, 6000)));
                        break;
                    case 281: // nhóm 5
                        item.itemOptions.add(new ItemOption(14, Util.nextInt(11, 13)));
                        break;
                }

                // rơi ra map
                ItemMap it = new ItemMap(zone, item, x, yEnd, player.id);
                list.add(it);
            }
        }

        if (this.tempId > 53) {
            int ratio = player.nPoint.power < 2_000_000_000L ? 400 : 700; // dưới 10 tỷ 1/1000, trên 10 tỷ 1/2000
            if (Util.isTrue(1, adjustMobDropRatio(ratio))) {
                int[] dropItems = { 241, 253, 265, 277, 233, 245, 257, 269, 237, 249, 261, 273, 281 };
                int randomItem = dropItems[Util.nextInt(0, dropItems.length - 1)];

                Item item = ItemService.gI().createNewItem((short) randomItem);

                // Gắn option 210 theo yêu cầu
                int[] op210Choices = { 50, 77, 103, 47 };
                int op210 = op210Choices[Util.nextInt(0, op210Choices.length - 1)];
                int value210 = (op210 == 47) ? Util.nextInt(100, 500) : Util.nextInt(1, 5);
                item.itemOptions.add(new ItemOption(op210, value210));

                // Gắn option theo từng nhóm
                switch (randomItem) {
                    case 241:
                    case 233:
                    case 237: // nhóm 1
                        item.itemOptions.add(new ItemOption(47, Util.nextInt(400, 550)));
                        break;
                    case 253:
                    case 245:
                    case 249: // nhóm 2
                        item.itemOptions.add(new ItemOption(6, Util.nextInt(22000, 27000)));
                        item.itemOptions.add(new ItemOption(27, Util.nextInt(3000, 5000)));
                        break;
                    case 265:
                    case 261:
                    case 257: // nhóm 3
                        item.itemOptions.add(new ItemOption(0, Util.nextInt(2100, 2400)));
                        break;
                    case 277:
                    case 269:
                    case 273: // nhóm 4
                        item.itemOptions.add(new ItemOption(7, Util.nextInt(22000, 26000)));
                        item.itemOptions.add(new ItemOption(28, Util.nextInt(4000, 6000)));
                        break;
                    case 281: // nhóm 5
                        item.itemOptions.add(new ItemOption(14, Util.nextInt(11, 13)));
                        break;
                }

                // rơi ra map
                ItemMap it = new ItemMap(zone, item, x, yEnd, player.id);
                list.add(it);
            }
        }

        if (mapid == 80) {
            if (Util.isTrue(1, adjustMobDropRatio(100))) { // 1% rơi item đặc biệt
                int[] dropItems = { 241, 253, 265, 277, 233, 245, 257, 269, 237, 249, 261, 273, 281 };
                int randomItem = dropItems[Util.nextInt(0, dropItems.length - 1)];

                Item item = ItemService.gI().createNewItem((short) randomItem);

                // Gắn option đầu tiên theo yêu cầu
                int[] opChoices = { 50, 77, 103, 47 };
                int op = opChoices[Util.nextInt(0, opChoices.length - 1)];
                int value;
                if (op == 47) {
                    value = Util.nextInt(100, 500);
                } else {
                    value = Util.nextInt(1, 5);
                }
                item.itemOptions.add(new ItemOption(op, value));

                // Gắn option theo từng nhóm
                switch (randomItem) {
                    case 241:
                    case 233:
                    case 237: // nhóm 1
                        item.itemOptions.add(new ItemOption(47, Util.nextInt(400, 550)));
                        break;
                    case 253:
                    case 245:
                    case 249: // nhóm 2
                        item.itemOptions.add(new ItemOption(6, Util.nextInt(22000, 27000)));
                        item.itemOptions.add(new ItemOption(27, Util.nextInt(3000, 5000)));
                        break;
                    case 265:
                    case 261:
                    case 257: // nhóm 3
                        item.itemOptions.add(new ItemOption(0, Util.nextInt(2100, 2400)));
                        break;
                    case 277:
                    case 269:
                    case 273: // nhóm 4
                        item.itemOptions.add(new ItemOption(7, Util.nextInt(22000, 26000)));
                        item.itemOptions.add(new ItemOption(28, Util.nextInt(4000, 6000)));
                        break;
                    case 281: // nhóm 5
                        item.itemOptions.add(new ItemOption(14, Util.nextInt(11, 13)));
                        break;
                }

                // rơi ra map
                ItemMap it = new ItemMap(zone, item, x, yEnd, player.id);
                list.add(it);
            } else { // 99% còn lại rơi 1283
                Item item = ItemService.gI().createNewItem((short) 1283);
                item.quantity = 1;

                ItemMap it = new ItemMap(zone, item, x, yEnd, player.id);
                list.add(it);
            }
        }

        // if (mapid == 131 || mapid == 132) {
        // if (player != null) {
        // if (InventoryService.gI().findItemBody(player, 1276) != null) {
        // if (Util.isTrue(10, 100)) {
        // Item item = ItemService.gI().createNewItem((short) 1243);
        // item.itemOptions.add(new ItemOption(210, 1));
        // ItemMap it = new ItemMap(zone, item, x, yEnd, player.id);
        // list.add(it);
        // }
        // }
        // if (InventoryService.gI().findItemBody(player, 1277) != null) {
        // if (Util.isTrue(10, 100)) {
        // Item item = ItemService.gI().createNewItem((short) 1742);
        // item.itemOptions.add(new ItemOption(210, 1));
        // ItemMap it = new ItemMap(zone, item, x, yEnd, player.id);
        // list.add(it);
        // }
        // }
        // }
        // }

        if (player.itemTime.isUseMayDo &&
                Util.isTrue(1, adjustMobDropRatio(100)) &&
                this.tempId > 57 && this.tempId < 66) {
            list.add(new ItemMap(zone, 380, 1, x, yEnd, player.id));
        }

        // if (player.itemTime.isUseMayDo2 && Util.isTrue(15, 100) && this.tempId > 80
        // && this.tempId < 81) {
        // list.add(new ItemMap(zone, 1110, 1, x, yEnd, player.id));
        // }

        if (player.isPl() && TaskService.gI().getIdTask(player) == ConstTask.TASK_8_1) {
            if (player.gender == 0 && this.tempId == 11 || player.gender == 1 && this.tempId == 12
                    || player.gender == 1 && this.tempId == 2) {
                list.add(new ItemMap(zone, 194, 1, x, yEnd, player.id));
            }
        }
        if (player.nPoint.power > 5_000_000_000L && Util.isTrue(1, adjustMobDropRatio(50))) {
            ItemMap it = new ItemMap(zone, ConstItem.THIT_HEO, 1, x, yEnd, player.id);
            it.options.add(new Item.ItemOption(86, 0));
            list.add(it);
        }

        if (player.nPoint.power > 10_000_000_000L && Util.isTrue(1, adjustMobDropRatio(100))) {
            ItemMap it = new ItemMap(zone, ConstItem.THUNG_NEP, 1, x, yEnd, player.id);
            it.options.add(new Item.ItemOption(86, 0));
            list.add(it);
        }

        if (Util.isTrue(1, adjustMobDropRatio(50)) && this.tempId != 0 && mapid >= 92 && mapid <= 100) {
            ItemMap it = new ItemMap(zone, ConstItem.THUNG_DAU_XANH, 1, x, yEnd, player.id);
            it.options.add(new Item.ItemOption(86, 0));
            list.add(it);
        }

        if (Util.isTrue(1, adjustMobDropRatio(50)) && this.tempId != 0 && mapid >= 105 && mapid <= 110) {
            ItemMap it = new ItemMap(zone, ConstItem.LA_DONG, 1, x, yEnd, player.id);
            it.options.add(new Item.ItemOption(86, 0));
            list.add(it);
        }

        // if (this.tempId != 0 && mapid >= 122 && mapid <= 124) {
        //     // Rơi QUA_HONG_DAO theo sức mạnh
        //     int chance = player.nPoint.power < 2_000_000_000L ? 50 : 300; // dưới 5 tỷ = 1/100, trên 10 tỷ = 1/500
        //     if (Util.isTrue(1, chance)) {
        //         ItemMap it = new ItemMap(zone, ConstItem.QUA_HONG_DAO, 1, x, yEnd, player.id);
        //         list.add(it);
        //     }

        //     // Rơi item 542 với tỉ lệ 1/1000
        //     if (Util.isTrue(1, 250)) {
        //         ItemMap it542 = new ItemMap(zone, 542, 1, x, yEnd, player.id);
        //         list.add(it542);
        //     }
        // }
// Drop khi mặc cải trang 577
Item caiTrang = InventoryService.gI().findItemBody(player, 461);
if (caiTrang != null) {

    // 443 - tỉ lệ 1/100
    if (Util.isTrue(1, adjustMobDropRatio(50))) {
        Item item = ItemService.gI().createNewItem((short) 443);
        item.itemOptions.add(new ItemOption(97, 5));
        list.add(new ItemMap(zone, item, x, yEnd, player.id));
    }

    // 447 - tỉ lệ 1/200
    if (Util.isTrue(1, adjustMobDropRatio(200))) {
        Item item = ItemService.gI().createNewItem((short) 447);
        item.itemOptions.add(new ItemOption(101, 5));
        list.add(new ItemMap(zone, item, x, yEnd, player.id));
    }

    // 442 - tỉ lệ 1/400
    if (Util.isTrue(1, adjustMobDropRatio(400))) {
        Item item = ItemService.gI().createNewItem((short) 442);
        item.itemOptions.add(new ItemOption(96, 5));
        list.add(new ItemMap(zone, item, x, yEnd, player.id));
    }

    // 441 - tỉ lệ 1/500
    if (Util.isTrue(1, adjustMobDropRatio(500))) {
        Item item = ItemService.gI().createNewItem((short) 441);
        item.itemOptions.add(new ItemOption(95, 5));
        list.add(new ItemMap(zone, item, x, yEnd, player.id));
    }
}
        if (mapid == 182) {
            if (Util.isTrue(1, adjustMobDropRatio(20))) { // 10% tỉ lệ rơi
                int[] itemIds = { 691, 692, 693, 3,4,5,227,228,229 };
                int itemId = itemIds[Util.nextInt(0, itemIds.length - 1)];
                ItemMap it = new ItemMap(zone, itemId, 1, x, yEnd, player.id);

                // Nhóm 1: 0, 6, 7, 47 -> param 1 -> 30000
                int[] group1 = { 0, 6, 7, 47 };
                it.options.add(new Item.ItemOption(group1[Util.nextInt(0, group1.length - 1)], Util.nextInt(1, 30000)));

                // Nhóm 2: 50, 77, 103, 5 -> param 1 -> 15
                int[] group2 = { 50, 77, 103, 47 };
                it.options.add(new Item.ItemOption(group2[Util.nextInt(0, group2.length - 1)], Util.nextInt(1, 15)));

                // Nhóm 3: 95, 96 -> param 1 -> 20
                int[] group3 = { 95, 96, 97 };
                it.options.add(new Item.ItemOption(group3[Util.nextInt(0, group3.length - 1)], Util.nextInt(1, 20)));

  int[] group4 = { 93 };
                it.options.add(new Item.ItemOption(group4[Util.nextInt(0, group4.length - 1)], Util.nextInt(0, 2)));

                list.add(it);
            }
        }

        // if (this.tempId != 0
        // && (mapid == 165 || mapid == 155 || mapid == 167 || mapid == 177)) {

        // if (Util.isTrue(1, adjustMobDropRatio(500))) {
        // ItemMap it1 = new ItemMap(zone, ConstItem.SUP_BI_HAC_AM, 1, x, yEnd,
        // player.id);
        // list.add(it1);
        // }

        // if (Util.isTrue(1, adjustMobDropRatio(500))) {
        // ItemMap it2 = new ItemMap(zone, ConstItem.HAMBURGER_SAU, 1, x, yEnd,
        // player.id);
        // list.add(it2);
        // }
        // }
if (mapid == 5) {

    int ratio;

    if (player.nPoint.power < 5_000_000_000L) {
        ratio = 20;   // 1/50
    } else {
        ratio = 500;  // 1/200
    }

    if (Util.isTrue(1, adjustMobDropRatio(ratio))) {
        ItemMap it = new ItemMap(zone, 728, 1, x, yEnd, player.id);
        list.add(it);
    }
}

   if (mapid == 79 || mapid == 77) {

    // item 20 - 1/500
    if (Util.isTrue(1, adjustMobDropRatio(500))) {
        ItemMap it = new ItemMap(zone, 20, 1, x, yEnd, player.id);
        list.add(it);
    }

    // item 19 - 1/500
    if (Util.isTrue(1, adjustMobDropRatio(500))) {
        ItemMap it = new ItemMap(zone, 19, 1, x, yEnd, player.id);
        list.add(it);
    }

    // item 18 - 1/100
    if (Util.isTrue(1, adjustMobDropRatio(100))) {
        ItemMap it = new ItemMap(zone, 18, 1, x, yEnd, player.id);
        list.add(it);
    }

    // item 728 - 1/500
    if (Util.isTrue(1, adjustMobDropRatio(500))) {
        ItemMap it = new ItemMap(zone, 728, 1, x, yEnd, player.id);
        list.add(it);
    }
}
if (mapid == 181) {
    if (Util.isTrue(1, adjustMobDropRatio(500))) { 
        ItemMap it = new ItemMap(zone, 1718, 1, x, yEnd, player.id);
        list.add(it);
    }
}
        // if (Util.isTrue(10, 100) && MapService.gI().isMapPhoBan(mapid)) {
        // ItemMap it = new ItemMap(zone, ConstItem.DAU_XANH, 1, x, yEnd, player.id);
        // it.options.add(new Item.ItemOption(86, 0));
        // list.add(it);
        // }

        // if (MapService.gI().isMapPhoBan(mapid) && this.tempId != 22) {
        // if (Util.isTrue(10, 100) || (player.actived() && Util.isTrue(20, 50))) {
        // list.add(new ItemMap(zone, 2029, 1, x, yEnd, player.id));
        // }
        // if (Util.isTrue(50, 100) || (player.actived() && Util.isTrue(10, 50))) {
        // list.add(new ItemMap(zone, 2055, 1, x, yEnd, player.id));
        // }
        // }
        // if (MapService.gI().isMapNguHanhSon(mapid) && this.tempId != 22) {
        // if (Util.isTrue(20, 100) || (player.actived() && Util.isTrue(50, 100))) {
        // player.event.addEventPointNHS(1);
        // }
        // }

        if (mapid == 124 || mapid == 123) {
            if (Util.isTrue(1, adjustMobDropRatio(5))) {
                ItemMap it = new ItemMap(zone, 934, Util.nextInt(1, 10), x, yEnd, player.id);
                list.add(it);
            }
        }

        if (mapid == 122) {
            if (Util.isTrue(1, adjustMobDropRatio(10))) {
                ItemMap it = new ItemMap(zone, 933, 1, x, yEnd, player.id);
                it.options.add(new Item.ItemOption(31, Util.nextInt(1, 5)));
                list.add(it);
            }
        }

        // Vang roi
        // if (Util.isTrue(1, 1000) || (Manager.TEST && Util.isTrue(1, 10)) ||
        // (player.actived() && Util.isTrue(1, 10))) {
        //
        // ItemMap it = new ItemMap(zone, 457, 1, x, yEnd, player.id);
        //
        // list.add(it);
        // }
        // Da ngu sac
        // if (mapid >= 155 && mapid <= 155 && Util.isTrue(1, adjustMobDropRatio(20))) {
        // list.add(new ItemMap(zone, 674, 1, x, yEnd, player.id));
        // }
        // Set kich hoat

        // if (((Util.isTrue(1, 1)))
        // && MapService.gI().isMapUpSKH(mapid)) {
        // short itTemp = (short) ItemService.gI().randTempItemKichHoat(player.gender);
        // ItemMap it = new ItemMap(zone, itTemp, 1, x, yEnd, player.id);
        // List<Item.ItemOption> ops = ItemService.gI().getListOptionItemShop(itTemp);
        // if (!ops.isEmpty()) {
        // it.options = ops;
        // }

        // int[] opsrand = ItemService.gI().randOptionItemKichHoat(player.gender);
        // it.options.add(new Item.ItemOption(opsrand[0], 0));
        // it.options.add(new Item.ItemOption(opsrand[1], 0));

        // it.options.add(new Item.ItemOption(30, 0));
        // list.add(it);
        // }

        // Set kich hoat Vip
if (Util.isTrue(1, adjustMobDropRatio(100))
        && (MapService.gI().isMapUpSKH(mapid) || mapid == 5)) {

    short itTemp = (short) ItemService.gI().randTempItemKichHoat(player.gender);
    ItemMap it = new ItemMap(zone, itTemp, 1, x, yEnd, player.id);

    List<Item.ItemOption> ops = ItemService.gI().getListOptionItemShop(itTemp);
    if (!ops.isEmpty()) {
        it.options = ops;
    }

    int[] opsrand = ItemService.gI().randOptionItemKichHoatNew(player.gender);

    it.options.add(new Item.ItemOption(opsrand[0], 0));
    it.options.add(new Item.ItemOption(opsrand[1], 0));
    it.options.add(new Item.ItemOption(opsrand[2], 0));
    it.options.add(new Item.ItemOption(opsrand[3], 0));
    it.options.add(new Item.ItemOption(30, 0));

    list.add(it);
}
        // Sao pha le
        // if (Util.isTrue(1, adjustMobDropRatio(100)) || (player.nPoint.isDoSPL && Util.isTrue(5, 100))
        // || (player.actived() && Util.isTrue(3, 100))) {
        // int rand = Util.nextInt(0, 6);
        // ItemMap it = new ItemMap(zone, 441 + rand, 1, x, yEnd, player.id);
        // it.options.add(new Item.ItemOption(95 + rand, (rand == 3 || rand == 4) ? 3 :
        // 5));
        // list.add(it);
        // }

        // Da nang cap
        if ((MapService.gI().isMapTuongLai(mapid) && Util.isTrue(1, adjustMobDropRatio(50)))
                || (!MapService.gI().isMapTuongLai(mapid) && Util.isTrue(1, adjustMobDropRatio(100)))) {
            int rand = Util.nextInt(0, 4);
            ItemMap it = new ItemMap(zone, 220 + rand, 1, x, yEnd, player.id);
            it.options.add(new Item.ItemOption(71 - rand, 0));
            list.add(it);
        }

        if (MapService.gI().isMapCold(mapid)) {
            if (player.isPet) {
                player = ((Pet) player).master;
            }
            // if (player.isPet) {
            // player = ((Pet) player).master;
            // }
            if (Util.isTrue(1, adjustMobDropRatio(10000000))) {
                ItemMap it = ItemService.gI().randDoTL(this.zone, 1, x, yEnd, player.id);
                list.add(it);
                ServerNotify.gI().notify(player.name + " vừa nhặt được " + it.itemTemplate.name + " tại "
                        + this.zone.map.mapName + " khu " + this.zone.zoneId);
            }
            if (Util.isTrue(1, adjustMobDropRatio(50)) && InventoryService.gI().fullSetThan(player)) {
                ItemMap it = new ItemMap(zone, Util.nextInt(663, 667), 1, x, yEnd, player.id);
                it.options.add(new Item.ItemOption(50, 10));
                list.add(it);
            }
        }

        // if (MapService.gI().isMapTuongLai(mapid)
        // && ((Util.isTrue(1, 1000) || (player.actived() && Util.isTrue(1, 300))))
        // && InventoryService.gI().fullSetThan(player)) {
        // ItemMap it = new ItemMap(zone, Util.nextInt(663, 667), 1, x, yEnd,
        // player.id);
        // it.options.add(new Item.ItemOption(30, 0));
        // list.add(it);
        // }

        // if (Util.isTrue(1, 100000) || (player.actived() && Util.isTrue(1, 10000))) {
        // list.add(new ItemMap(zone, 457, 1, x, yEnd, player.id));
        // }
        // Manh thien su
if (Util.isTrue(1, adjustMobDropRatio(100)) 
    && MapService.gI().isMapHanhTinhThucVat(mapid) 
    && InventoryService.gI().findItemNTK(player)) {
    
    // Sử dụng Util.nextInt(1, 5) để lấy ngẫu nhiên từ 1 đến 5 vật phẩm
    list.add(new ItemMap(zone, Util.nextInt(1066, 1070), Util.nextInt(1, 5), x, yEnd, player.id));
}
        // if (Util.isTrue(1, adjustMobDropRatio(200)) || (player.actived() && Util.isTrue(1, adjustMobDropRatio(100)))) {
        // list.add(new ItemMap(zone, 861, 1, x, yEnd, player.id));
        // }
        // if (player.nPoint.power >= 80000000000L) {
        // if (player.zone.map.mapId == 155) {
        // list.add(new ItemMap(zone, 2055, 1, x, yEnd, player.id));
        // } else {
        // if (Util.isTrue(1, adjustMobDropRatio(500)) || (player.actived() && Util.isTrue(1, adjustMobDropRatio(100)))) {
        // list.add(new ItemMap(zone, 2051, 1, x, yEnd, player.id));
        // }
        // if (Util.isTrue(1, 1000) || (player.actived() && Util.isTrue(1, adjustMobDropRatio(200)))) {
        // list.add(new ItemMap(zone, 2052, 1, x, yEnd, player.id));
        // }
        // }
        // }
        return list;
    }

    private ItemMap dropItemTask(Player player) {
        ItemMap itemMap = null;
        switch (tempId) {
            case ConstMob.KHUNG_LONG:
            case ConstMob.LON_LOI:
            case ConstMob.QUY_DAT:
                if (TaskService.gI().getIdTask(player) == ConstTask.TASK_2_0) {
                    itemMap = new ItemMap(zone, 73, 1, location.x, location.y, player.id);
                }
                break;
            case ConstMob.THAN_LAN_ME:
                if (TaskService.gI().getIdTask(player) == ConstTask.TASK_8_1) {
                    if (Util.isTrue(1, adjustMobDropRatio(2))) {
                        // Tạo item với option
                        Item item = ItemService.gI().createNewItem((short) 194, 1); // 223 là tempId


                        // Rớt ra map
                        itemMap = new ItemMap(zone, item, location.x, location.y, player.id);
                    } else {
                        Service.gI().sendThongBao(player,
                                "Con thằn lằn mẹ này không giữ Viên CapSule, hãy tìm con thằn lằn mẹ khác");
                    }

                }
            case ConstMob.OC_MUON_HON:
                if (TaskService.gI().getIdTask(player) == ConstTask.TASK_14_1) {
                    if (Util.isTrue(1, adjustMobDropRatio(20))) {
                        itemMap = new ItemMap(zone, 85, 1, location.x, location.y, player.id);
                    } else {
                        Service.gI().sendThongBao(player,
                                "Con ốc mượn hồn này không giữ truyện tranh, hãy thử tìm con ốc mượn hồn khác");
                    }
                }
            case ConstMob.HEO_XAYDA_ME:
                if (TaskService.gI().getIdTask(player) == ConstTask.TASK_14_1) {
                    if (Util.isTrue(1, adjustMobDropRatio(20))) {
                        itemMap = new ItemMap(zone, 85, 1, location.x, location.y, player.id);
                    } else {
                        Service.gI().sendThongBao(player,
                                "Con heo xayda mẹ này không giữ truyện tranh, hãy thử tìm con heo xayda mẹ khác");
                    }
                }
            case ConstMob.OC_SEN:
                if (TaskService.gI().getIdTask(player) == ConstTask.TASK_14_1) {
                    if (Util.isTrue(1, adjustMobDropRatio(20))) {
                        itemMap = new ItemMap(zone, 85, 1, location.x, location.y, player.id);
                    } else {
                        Service.gI().sendThongBao(player,
                                "Con ốc xên này không giữ truyện tranh, hãy thử tìm con ốc xên khác");
                    }
                }
        }
        if (itemMap != null) {
            return itemMap;
        }
        return null;
    }

    private void sendMobStillAliveAffterAttacked(long dameHit, boolean crit) {
        Message msg;
        try {
            msg = new Message(-9);
            msg.writer().writeByte(this.id);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.gethp()), cn.readInt);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(dameHit), cn.readInt);
            msg.writer().writeBoolean(crit); // chí mạng
            msg.writer().writeInt(-1);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    public void hoiSinhMobPhoBan() {
        this.point.hp = this.point.maxHp;
        this.setTiemNang();
        Message msg;
        try {
            msg = new Message(-13);
            msg.writer().writeByte(this.id);
            msg.writer().writeByte(this.tempId);
            msg.writer().writeByte(this.lvMob); // level mob
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    public void hoiSinhMobTayKarin() {
        this.point.hp = this.point.maxHp;
        this.maxTiemNang = 1;
        Message msg;
        try {
            msg = new Message(-13);
            msg.writer().writeByte(this.id);
            msg.writer().writeByte(this.tempId);
            msg.writer().writeByte(this.lvMob); // level mob
            msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendSieuQuai(int type) {
        Message msg;
        try {
            msg = new Message(-75);
            msg.writer().writeByte(this.id);
            msg.writer().writeByte(type);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendDisable(boolean bool) {
        Message msg;
        try {
            msg = new Message(81);
            msg.writer().writeByte(this.id);
            msg.writer().writeBoolean(bool);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendDoneMove(boolean bool) {
        Message msg;
        try {
            msg = new Message(82);
            msg.writer().writeByte(this.id);
            msg.writer().writeBoolean(bool);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendFire(boolean bool) {
        Message msg;
        try {
            msg = new Message(85);
            msg.writer().writeByte(this.id);
            msg.writer().writeBoolean(bool);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendIce(boolean bool) {
        Message msg;
        try {
            msg = new Message(86);
            msg.writer().writeByte(this.id);
            msg.writer().writeBoolean(bool);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendWind(boolean bool) {
        Message msg;
        try {
            msg = new Message(87);
            msg.writer().writeByte(this.id);
            msg.writer().writeBoolean(bool);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendMobMaxHp(long maxHp) {
        Message msg;
        try {
            msg = new Message(87);
            msg.writer().writeByte(this.id);
            msg.writeLongByGiaCatLuong(Util.maxIntValue(maxHp), cn.readInt);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    private void phanSatThuong(Player plTarget, long dame) {
        if (plTarget.nPoint == null) {
            return;
        }
        int percentPST = plTarget.nPoint.tlPST;
        if (percentPST != 0) {
            long damePST = Util.maxIntValue(dame * percentPST / 100L);
            Message msg;
            try {
                msg = new Message(-9);
                msg.writer().writeByte(this.id);
                if (damePST >= this.point.hp) {
                    damePST = this.point.hp - 1;
                }
                long hpMob = Util.maxIntValue(this.point.hp);
                injured(null, damePST, true);
                damePST = hpMob - this.point.hp;
                msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);
                msg.writeLongByGiaCatLuong(Util.maxIntValue(damePST), cn.readInt);
                msg.writer().writeBoolean(false);
                msg.writer().writeByte(36);
                Service.gI().sendMessAllPlayerInMap(this.zone, msg);
                msg.cleanup();
            } catch (IOException e) {
            }
        }
    }

    public void startDie() {
        Message msg;
        try {
            setDie();
            this.point.hp = -1;
            this.status = 0;
            msg = new Message(-12);
            msg.writer().writeByte(this.id);
            Service.gI().sendMessAllPlayerInMap(this.zone, msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }
}
