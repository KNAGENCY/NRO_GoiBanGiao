package player.Bot;

import consts.ConstPlayer;
import map.Map;
import map.Zone;
import player.Player;
import services.MapService;
import services.PlayerService;
import services.Service;
import services.func.ChangeMapService;
import skill.Skill;
import utils.SkillUtil;
import utils.Util;

public class Bot extends Player {

    private final short head_;
    private final short body_;
    private final short leg_;
    private final short flag_;
    private final int assignedMapId;
    private int assignedZoneId;
    private long lastTimeChangeZone;
    public ShopBot shop;
    public Mobb mo1;
    public AmbientBot ambient;

    public Bot(String name, byte gender, short head, short body, short leg, short flag, int assignedMapId, int assignedZoneId) {
        this.name = name;
        this.gender = gender;
        this.head_ = head;
        this.body_ = body;
        this.leg_ = leg;
        this.flag_ = flag;
        this.assignedMapId = assignedMapId;
        this.assignedZoneId = assignedZoneId;
        this.lastTimeChangeZone = System.currentTimeMillis();
        this.id = Util.nextInt(1, 2_000_000_000);
        this.isBot = true;
        this.isPlayer = false;
    }

    public int getAssignedMapId() {
        return assignedMapId;
    }

    public int getAssignedZoneId() {
        return assignedZoneId;
    }

    public void setAssignedZoneId(int assignedZoneId) {
        this.assignedZoneId = assignedZoneId;
    }

    public long getLastTimeChangeZone() {
        return lastTimeChangeZone;
    }

    public void setLastTimeChangeZone(long lastTimeChangeZone) {
        this.lastTimeChangeZone = lastTimeChangeZone;
    }

    public void joinAssignedZone(Zone zone) {
        if (zone != null) {
            ChangeMapService.gI().goToMap(this, zone);
            this.zone.load_Me_To_Another(this);
        }
    }

    public void joinMap() {
        joinAssignedZone(NewBot.gI().getAssignedZone(this.assignedMapId, this.assignedZoneId));
    }

    public Zone getRandomZone(int mapId) {
        Map map = MapService.gI().getMapById(mapId);
        if (map == null || map.zones == null || map.zones.isEmpty()) {
            return null;
        }
        if (mapId == this.assignedMapId && this.assignedZoneId >= 0 && this.assignedZoneId < map.zones.size()) {
            return map.zones.get(this.assignedZoneId);
        }
        return map.zones.get(0);
    }

    @Override
    public short getHead() {
        if (effectSkill.isMonkey) {
            return (short) ConstPlayer.HEADMONKEY[effectSkill.levelMonkey - 1];
        }
        return this.head_;
    }

    @Override
    public short getBody() {
        if (effectSkill.isMonkey) {
            return 193;
        }
        return this.body_;
    }

    @Override
    public short getLeg() {
        if (effectSkill.isMonkey) {
            return 194;
        }
        return this.leg_;
    }

    @Override
    public short getFlagBag() {
        return this.flag_;
    }

    @Override
    public void update() {
        super.update();
        this.nPoint.tlHutHp = (short) (NewBot.gI().isHighTrainingMap(this.assignedMapId) ? 100 : 50);
        this.increasePoint();
        if (this.mo1 != null) {
            this.mo1.update();
        }
        if (this.ambient != null) {
            this.ambient.update();
        }
        if (this.isDie()) {
            Service.gI().hsChar(this, nPoint.hpMax, nPoint.mpMax);
        }
    }

    public void leakSkill() {
        this.playerSkill.skills.clear();
        Skill firstSkill = null;
        Skill secondSkill = null;
        switch (this.gender) {
            case ConstPlayer.TRAI_DAT -> {
                firstSkill = SkillUtil.createSkill(Skill.DRAGON, 7);
                secondSkill = SkillUtil.createSkill(Skill.KAMEJOKO, 7);
            }
            case ConstPlayer.NAMEC -> {
                firstSkill = SkillUtil.createSkill(Skill.DEMON, 7);
                secondSkill = SkillUtil.createSkill(Skill.MASENKO, 7);
            }
            case ConstPlayer.XAYDA -> {
                firstSkill = SkillUtil.createSkill(Skill.GALICK, 7);
                secondSkill = SkillUtil.createSkill(Skill.ANTOMIC, 7);
            }
        }
        if (firstSkill != null) {
            firstSkill.manaUse = 0;
            this.playerSkill.skills.add(firstSkill);
        }
        if (secondSkill != null) {
            secondSkill.manaUse = 0;
            this.playerSkill.skills.add(secondSkill);
        }
        if (!this.playerSkill.skills.isEmpty()) {
            this.playerSkill.skillSelect = this.playerSkill.skills.get(0);
        }
    }

    public boolean UseLastTimeSkill() {
        if (this.playerSkill.skillSelect.lastTimeUseThisSkillbot < (System.currentTimeMillis()
                - this.playerSkill.skillSelect.coolDown)) {
            this.playerSkill.skillSelect.lastTimeUseThisSkillbot = System.currentTimeMillis();
            return true;
        }
        return false;
    }

    private void increasePoint() {
        long tiemNangUse;
        int point;
        if (this.nPoint == null) {
            return;
        }
        if (Util.isTrue(50, 100)) {
            point = 100;
            int pointHp = point * 20;
            tiemNangUse = point * (2 * (this.nPoint.hpg + 1000) + pointHp - 20) / 2;
            if (doUseTiemNang(tiemNangUse)) {
                this.nPoint.hpMax += point;
                this.nPoint.hpg += point;
                Service.gI().point(this);
            }
        } else {
            point = 10;
            tiemNangUse = point * (2 * this.nPoint.dameg + point - 1) / 2 * 100;
            if (doUseTiemNang(tiemNangUse)) {
                this.nPoint.dameg += point;
                Service.gI().point(this);
            }
        }
    }

    private boolean doUseTiemNang(long tiemNang) {
        if (this.nPoint.tiemNang < tiemNang) {
            return false;
        }
        this.nPoint.tiemNang -= tiemNang;
        return true;
    }

    public void moveTo(int x, int y) {
        PlayerService.gI().playerMove(this, x, y);
    }

    public void useSkill(int skillId) {
    }
}
