package player;
import services.func.PhucLoiService;

/*
 *
 *
 * @author GiaCatLuong
 */
import GiaCatLuongManager.Functions;
import boss.BossID;
import boss.BossManager;
import consts.ConstDailyGift;
// import minigame.cost.LuckyNumberCost;
// import minigame.LuckyNumber.LuckyNumberService;
import models.Training.TrainingService;
import models.War.War;
import npc.NonInteractiveNPC;
import models.Card.Card;
import models.Card.RadarCard;
import models.Card.RadarService;
import models.MajinBuu.MajinBuuService;

import player.dailyGift.DailyGiftData;
import player.dailyGift.DailyGiftService;
import services.*;
import skill.PlayerSkill;

import java.util.Iterator;
import java.util.List;

import clan.Clan;
import intrinsic.IntrinsicPlayer;
import item.Item;
import item.ItemTime;
import npc.specialnpc.MagicTree;
import consts.ConstPlayer;
import consts.ConstTask;
import npc.specialnpc.MabuEgg;
import mob.MobMe;
import data.DataGame;
import clan.ClanMember;
import consts.ConstAchievement;
import java.time.LocalDate;
import map.Zone;
import matches.IPVP;
import matches.TYPE_LOSE_PVP;
import skill.Skill;
import server.io.MySession;

import task.TaskPlayer;
import network.Message;
import server.Client;
import services.func.ChangeMapService;
import models.Combine.Combine;
import utils.Logger;
import utils.Util;

import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashSet;
import java.util.Set;
import java.util.Timer;

import lombok.Getter;
import lombok.Setter;
import models.Template.AchievementQuest;
import models.BlackBallWar.BlackBallWarService;
import map.ItemMap;
import map.MaBuHold;
import models.MajinBuu.MajinBuu14H;
import models.ShenronEvent.ShenronEvent;
import server.Maintenance;
import utils.SkillUtil;

public class Player implements Runnable {

    public long lastTimeEatPea;

    @Setter
    @Getter


    
    public MySession session;
    public long lastTimeTXChat = 0;
    public long lastBetTime = 0;
    public long id;
    public String name;
    public byte gender;
    public boolean isNewMember;
    public short head;
    public int deltaTime;
    public int xBeforeLogout;
    public int yBeforeLogout;
    public byte typePk;
    public byte cFlag;
    public boolean haveTennisSpaceShip;
    public boolean isCopy;
    public boolean beforeDispose;
    public int mbv = 0;
    public boolean baovetaikhoan;
    public long mbvtime;
    public int timeGohome;
    public long lastUpdateGohomeTime;
    public boolean goHome;
    public long lastPkCommesonTime;
    public boolean callBossPocolo;
    public Zone zoneSieuThanhThuy;
    public boolean winSTT;
    public long lastTimeWinSTT;
    public long lastTimeUpdateSTT;
    public MajinBuu14H maBu2H;
    public boolean isMabuHold;
    public MaBuHold maBuHold;
    public int precentMabuHold;
    public boolean isPhuHoMapMabu;
    public boolean danhanthoivang;
    public long lastRewardGoldBarTime;
    public int timesPerDayBDKB = 0;
    public long lastTimeJoinBDKB;
    public boolean joinCDRD;
    public long lastTimeJoinCDRD;
    public boolean talkToThuongDe;
    public boolean talkToThanMeo;
    public long timeChangeMap144;
    public long lastTimeJoinDT;
    public int typeChibi;
    public long lastTimeChibi;
    public long lastTimeUpdateChibi;
    public String captcha = "";
    public int spamcaptcha = 0;
    public long lasttimebotchat;
    public boolean doesNotAttack;
    public long lastTimePlayerNotAttack;
    public int timeNotAttack = 1800000;
    public boolean isPet;
    public boolean isNewPet;
    public boolean isNewPet1;
    public boolean isBoss;
    public boolean isPlayer;
    public IPVP pvp;
    public byte maxTime = 30;
    public byte type = 0;
    public boolean isOffline = false;
    public String notify = null;
    public int mapIdBeforeLogout;
    public List<Zone> mapBlackBall;
    public List<Zone> mapMaBu;
    public List<Player> temporaryEnemies = new ArrayList<>();
    public Date firstTimeLogin;
    public Zone zone;
    public Zone mapBeforeCapsule;
    public List<Zone> mapCapsule;
    public Pet pet;
    public NewPet newPet;
    public MobMe mobMe;
    public Location location;
    public SetClothes setClothes;
    public EffectSkill effectSkill;
    public MabuEgg mabuEgg;
    public TaskPlayer playerTask;
    public ItemTime itemTime;
    public Fusion fusion;
    public MagicTree magicTree;
    public IntrinsicPlayer playerIntrinsic;
    public Inventory inventory;
    public PlayerSkill playerSkill;
    public Combine combine;
    public IDMark iDMark;
    public Charms charms;
    public EffectSkin effectSkin;
    public NPoint nPoint;
    public RewardBlackBall rewardBlackBall;
    public FightMabu fightMabu;
    public NewSkill newSkill;
    public Satellite satellite;
    public Achievement achievement;
    public GiftCode giftCode;
    public Traning traning;
    public int hoisinh;
    public int lastMenuId = -1;

    public int lastReviveResetDay;
    public long lastTimeLanConFollow = 0;

    public Clan clan;
    public ClanMember clanMember;
    public List<Friend> friends;
    public List<Enemy> enemies;
    public boolean justRevived;
    public long lastTimeRevived;
    public long timeChangeZone;
    public long lastUseOptionTime;
    public short idNRNM = -1;
    public short idGo = -1;
    public long lastTimePickNRNM;
    public List<Card> Cards = new ArrayList<>();
    public int levelWoodChest;
    public long goldChallenge;
    public long rubyChallenge;
    public long lastTimeRewardWoodChest;
    public List<Item> itemsWoodChest = new ArrayList<>();
    public int indexWoodChest;
    public long lastTimePKDHVT23;
    public boolean lostByDeath;
    public boolean isPKDHVT;
    public int xSend;
    public int ySend;
    public boolean isFly;
    public long lastTimeDietQuy;
    // shenron event
    public long lastTimeShenronAppeared;
    public boolean isShenronAppear;
    public ShenronEvent shenronEvent;
    // vo dai sinh tu
    public long lastTimePKVoDaiSinhTu;
    public boolean haveRewardVDST;
    public int thoiVangVoDaiSinhTu;
    public long timePKVDST;
    public int binhChonHatMit;
    public int binhChonPlayer;
    public Zone zoneBinhChon;
    public ItemEvent itemEvent;
    public int levelLuyenTap;
    public boolean isThachDau;
    public int tnsmLuyenTap;
    public boolean dangKyTapTuDong;
    public long lastTimeOffline;
    public int mapIdDangTapTuDong;
    public int lastMapOffline;
    public int lastZoneOffline;
    public int lastXOffline;
    public String thongBaoTapTuDong;
    public boolean teleTapTuDong;
    public int timesPerDayCuuSat;
    public long lastTimeCuuSat;
    public boolean nhanVangNangVIP;
    public boolean nhanDeTuNangVIP;
    public boolean nhanSKHVIP;
    public long totalDamageTaken;
    public boolean thongBaoChangeMap;
    public String textThongBaoChangeMap;
    public boolean thongBaoThua;
    public String textThongBaoThua;
    public SuperRank superRank;
    public boolean canReward;
    public boolean changeMapVIP;
    public boolean haveReward;
    public int tayThong;
    public List<Item> itemsTradeWVP = new ArrayList<>();
    public boolean tradeWVP;
    private DropItem dropItem;
    public boolean isBattu = false;
    public int nangDong;
public int npcPackage = 0;
public int npcDayReceived = 0;
public long npcLastReward = 0;
public long timeExchangeTN = 0;

    public List<DailyGiftData> dailyGiftData = new ArrayList<>();
    public int numUseSkill = 0;
    public int typeRecvieArchiment = 0;
    public PlayerEvent event;
    public PointFusion pointfusion;
    public List<Archivement> archivementList = new ArrayList<>();
    public List<ArchivementSucManh> archivementListSM = new ArrayList<>();

    public boolean isHoldNamecBallTranhDoat;
    public int tempIdNamecBallHoldTranhDoat = -1;
    public long lastTimePickItem;
    public boolean receivedAchievementReward = false;

    public long lastTimeUpdateBallWar;

    public Player() {
        lastUseOptionTime = System.currentTimeMillis();
        lastTimeDietQuy = 0;
        location = new Location();
        nPoint = new NPoint(this);
        inventory = new Inventory();
        playerSkill = new PlayerSkill(this);
        setClothes = new SetClothes(this);
        effectSkill = new EffectSkill(this);
        fusion = new Fusion(this);
        playerIntrinsic = new IntrinsicPlayer();
        rewardBlackBall = new RewardBlackBall(this);
        fightMabu = new FightMabu(this);
        // ----------------------------------------------------------------------
        iDMark = new IDMark();
        combine = new Combine();
        playerTask = new TaskPlayer();
        friends = new ArrayList<>();
        enemies = new ArrayList<>();
        itemTime = new ItemTime(this);
        charms = new Charms();
        effectSkin = new EffectSkin(this);
        newSkill = new NewSkill(this);
        satellite = new Satellite();
        achievement = new Achievement(this);
        giftCode = new GiftCode();
        traning = new Traning();
        itemEvent = new ItemEvent(this);
        superRank = new SuperRank(this);
        dropItem = new DropItem(this);

        event = new PlayerEvent(this);
        pointfusion = new PointFusion(this);
        archivementList = new ArrayList<>();
        archivementListSM = new ArrayList<>();
    }

    public PointFusion getPointfusion() {
        return this.pointfusion;
    }

    public boolean isBot;

    // --------------------------------------------------------------------------
    private Set<Integer> boughtItems = new HashSet<>(); // Dùng Set để tránh trùng lặp

    // Getter và Setter cho boughtItems
    public Set<Integer> getBoughtItems() {
        return boughtItems;
    }

    public void addBoughtItem(int itemTempId) {
        boughtItems.add(itemTempId);
    }
    // --------------------------------------------------------------------------

    public String getName() {
        return name;
    }

    // --------------------------------------------------------------------------
    public boolean isDie() {
        if (this.nPoint != null) {
            return this.nPoint.hp <= 0;
        }
        return true;
    }

    public void sendMessage(Message msg) {
        if (this.session != null) {
            session.sendMessage(msg);
        }
    }

    public boolean isPl() {
        return isPlayer && !isPet && !isBoss && !isNewPet && !isNewPet1 && !isBot
                && !(this instanceof NonInteractiveNPC);
    }

    @Override
    public void run() {
        Functions.sleep(500);
        while (!Maintenance.isRunning && session != null && session.isConnected() && this.name != null) {
            long st = System.currentTimeMillis();
            update();
            long time = 1000 - (System.currentTimeMillis() - st);
            if (time > 0) {
                Functions.sleep(time);
            }
        }
    }

    public void start() {
        new Thread(this, "Update player " + this.name).start();
    }

    public void update() {
        if (!this.beforeDispose && !this.isBot) {
            try {
                PhucLoiService.gI().tick(this);
                if (this.zone != null || (!this.isPl() && this.zone == null)) {
                    if (itemTime != null) {
                        itemTime.update();
                    }
                    if (magicTree != null) {
                        magicTree.update();
                    }
                    if (this.isPl() && this.zone != null && this.zone.map.mapId == this.gender + 21
                            && (TaskService.gI().getIdTask(this) == ConstTask.TASK_0_0
                                    || TaskService.gI().getIdTask(this) == ConstTask.TASK_0_1)) {
                        this.playerTask.taskMain.index = 2;
                        TaskService.gI().sendTaskMain(this);
                    }
                }
                if ((this.zone != null && !MapService.gI().isHome(this.zone.map.mapId))
                        || (!this.isPl() && this.zone == null)) {
                    if (isPl() && iDMark != null && iDMark.isBan()
                            && Util.canDoWithTime(iDMark.getLastTimeBan(), 5000)) {
                        Client.gI().kickSession(session);
                        return;
                    }
                    if (nPoint != null) {
                        nPoint.update();
                    }
                    if (fusion != null) {
                        fusion.update();
                    }
                    if (effectSkill != null) {
                        effectSkill.update();
                    }
                    if (mobMe != null) {
                        mobMe.update();
                    }
                    if (effectSkin != null) {
                        effectSkin.update();
                    }
                    if (pet != null) {
                        pet.update();
                    }
                    if (newPet != null) {
                        newPet.update();
                    }
                    // if (itemTime != null) {
                    // itemTime.update();
                    // }
                    if (satellite != null) {
                        satellite.update();
                    }

                    // Chibi
                    if (this.isPl() && !this.isDie() && this.effectSkill != null && !this.effectSkill.isChibi
                            && Util.canDoWithTime(lastTimeChibi, 300000)) {
                        if (Util.isTrue(1, 10) && !MapService.gI().isMapBlackBallWar(this.zone.map.mapId)) {
                            EffectSkillService.gI().setChibi(this, 600000);
                        }
                        lastTimeChibi = System.currentTimeMillis();
                    }
                    if (this.isPl() && !this.isDie() && this.effectSkill != null && this.effectSkill.isChibi
                            && Util.canDoWithTime(lastTimeUpdateChibi, 1000)) {
                        if (this.typeChibi == 1) {
                            if (this.nPoint.mp < this.nPoint.mpMax) {
                                if (this.nPoint.mpMax - this.nPoint.mp < this.nPoint.mpMax / 10) {
                                    this.nPoint.mp = this.nPoint.mpMax;
                                } else {
                                    this.nPoint.mp += this.nPoint.mpMax / 10;
                                }
                            }
                            PlayerService.gI().sendInfoMp(this);
                        } else if (this.typeChibi == 3) {
                            if (this.nPoint.hp < this.nPoint.hpMax) {
                                if (this.nPoint.hpMax - this.nPoint.hp < this.nPoint.hpMax / 10) {
                                    this.nPoint.hp = this.nPoint.hpMax;
                                } else {
                                    this.nPoint.hp += this.nPoint.hpMax / 10;
                                }
                            }
                            PlayerService.gI().sendInfoHp(this);
                        }
                        lastTimeUpdateChibi = System.currentTimeMillis();
                    }

                    if (this.isPl() && this.achievement != null) {
                        this.achievement.done(ConstAchievement.HOAT_DONG_CHAM_CHI, 1000);
                    }
                    if (this.isPl()) {
                        Calendar calendar = Calendar.getInstance();
                        int hour = calendar.get(Calendar.HOUR_OF_DAY);
                        if (!(hour >= 22 && hour <= 23)) {
                            if (zone == null) {
                                // System.out.println("Zone is null at Player.update");
                                return;
                            }
                            if (zone.map == null) {
                                // System.out.println("Zone.map is null at Player.update");
                                return;
                            }
                            if (zone.map.mapId == 126) {
                                ChangeMapService.gI().changeMapNonSpaceship(this, 19, 1000 + Util.nextInt(-100, 100),
                                        360);
                            }
                        }
                        // updateCSMM();
                        // TaskService.gI().sendUpdateCountSubTask(this);

                        // sendTextTimeDaiLyGift();
                        // TranhNgoc.gI().update(this);
                        // if (clan != null) {
                        //     ClanService.gI().checkDoneTaskJoinClan(clan);
                        // }
                        // thuhoivp();
                    }
                    if (this.isPl() && this.effectSkill != null && this.effectSkill.isMabuHold) {
                        this.nPoint.subHP(this.nPoint.hpMax / 100);
                        if (Util.isTrue(1, 10)) {
                            Service.gI().chat(this, "Cứu tôi với");
                        }
                        PlayerService.gI().sendInfoHp(this);
                        if (this.precentMabuHold > 15) {
                            EffectSkillService.gI().removeMabuHold(this);
                        }
                        if (this.nPoint.hp <= 0) {
                            EffectSkillService.gI().removeMabuHold(this);
                            setDie();
                        }
                    }

                    if (this.zone != null && this.effectSkin != null && this.effectSkin.xHPKI > 1
                            && !MapService.gI().isMapBlackBallWar(this.zone.map.mapId)) {
                        this.effectSkin.xHPKI = 1;
                        this.nPoint.calPoint();
                        Service.gI().point(this);
                    }

                    if (this.zone != null && this.effectSkin != null && this.effectSkin.xDame > 1
                            && !MapService.gI().isMapBlackBallWar(this.zone.map.mapId)) {
                        this.effectSkin.xDame = 1;
                        this.nPoint.calPoint();
                        Service.gI().point(this);
                    }

                    if (this.isPl() && this.zone != null) {
                        fixBlackBallWar();
                    }

                    if (this.zone != null && this.zone.map.mapId == (21 + this.gender)) {
                        if (this.mabuEgg != null) {
                            this.mabuEgg.sendMabuEgg();
                        }
                    }

                    if (this.isPhuHoMapMabu && this.zone != null && !MapService.gI().isMapMabu2H(this.zone.map.mapId)) {
                        this.isPhuHoMapMabu = false;
                        this.nPoint.calPoint();
                        Service.gI().point(this);
                        Service.gI().Send_Info_NV(this);
                        Service.gI().Send_Caitrang(this);
                    }
                    // Change Map 144 CDRD
                    // if (this.isPl() && this.clan != null && this.clan.ConDuongRanDoc != null
                    //         && this.joinCDRD && this.clan.ConDuongRanDoc.allMobsDead
                    //         && this.talkToThanMeo && this.zone.map.mapId == 47
                    //         && Util.canDoWithTime(timeChangeMap144, 5000)) {
                    //     ChangeMapService.gI().changeMapYardrat(this, this.clan.ConDuongRanDoc.getMapById(144),
                    //             300 + Util.nextInt(-100, 100), 312);
                    //     this.timeChangeMap144 = System.currentTimeMillis();
                    // }
                    // Auto tắt cờ khi rời map Mabu
                    if (this.isPl() && this.zone != null && !MapService.gI().isMapMaBu(this.zone.map.mapId)
                            && (this.cFlag == 9 || this.cFlag == 10)) {
                        Service.gI().changeFlag(this, 0);
                    }

                    if (this.isPl() && this.zone != null && MapService.gI().isMapMaBu(this.zone.map.mapId)
                            && this.cFlag != 9 && this.cFlag != 10) {
                        Service.gI().changeFlag(this, Util.nextInt(9, 10));
                    }
                    if (dropItem != null) {
                        dropItem.update();
                    }
                    MajinBuuService.gI().update(this);
    
                    if (!isBoss && this.iDMark != null && this.iDMark.isGotoFuture()
                            && Util.canDoWithTime(this.iDMark.getLastTimeGoToFuture(), 60000)) {
                        ChangeMapService.gI().changeMapBySpaceShip(this, 102, -1, Util.nextInt(60, 200));
                        this.iDMark.setGotoFuture(false);
                    }
                    if (this.isPl() && location != null
                            && location.lastTimeplayerMove < System.currentTimeMillis() - 5 * 60 * 60 * 1000) {
                        Client.gI().kickSession(session);
                    }
                }
            } catch (Exception e) {
                Logger.logException(Player.class, e, "Lỗi tại player: " + this.name);
            }
        }
    }

    public long lastTimeSendTextTime;

    // public void sendTextTimeDaiLyGift() {
    // if (Util.canDoWithTime(lastTimeSendTextTime, 300000)) {
    // if (DailyGiftService.checkDailyGift(this, ConstDailyGift.NHAN_BUA_MIEN_PHI))
    // {
    // ItemTimeService.gI().sendTextTime(this, itemTime.TEXT_NHAN_BUA_MIEN_PHI,
    // "Nhận ngẫu nhiên bùa 1h mỗi ngày tại Bà Hạt Mít ở vách núi", 30);
    // }
    // lastTimeSendTextTime = System.currentTimeMillis();
    // }
    // }
    // public void updateCSMM() {
    //     minigame.LuckyNumber.LuckyNumber.players.forEach((g) -> {
    //         if (this.id == g.id) {
    //             LuckyNumberService.showNumberPlayer(this, LuckyNumberService.strNumber(this.id));
    //             ItemTimeService.gI().sendItemTime(this, 2295, LuckyNumberCost.timeGame);
    //         }
    //     });
    // }

    public void upSkill(byte id, short cost) {

        int tempId = this.playerSkill.skills.get(id).template.id;
        int level = this.playerSkill.skills.get(id).point + 1;
        if (level > 7) {
            Service.gI().sendThongBao((Player) this, "Kĩ năng đã đạt cấp tối đa!");
        } else if (((Player) this).inventory.gem < cost) {
            Service.gI().sendThongBao((Player) this, "Bạn không đủ ngọc để nâng cấp!");
        } else {
            Skill skill = null;
            try {
                skill = SkillUtil.nClassTD.getSkillTemplate(tempId).skillss.get(level - 1);
            } catch (Exception e) {
                try {
                    skill = SkillUtil.nClassNM.getSkillTemplate(tempId).skillss.get(level - 1);
                } catch (Exception ex) {
                    skill = SkillUtil.nClassXD.getSkillTemplate(tempId).skillss.get(level - 1);
                }
            }
            skill = new Skill(skill);
            if (id == 1) {
                skill.coolDown = 1000;
            }
            this.playerSkill.skills.set(id, skill);
            ((Player) this).inventory.gem -= cost;
            Service.gI().sendMoney((Player) this);
            Service.gI().player((Player) this);

        }
    }

    // --------------------------------------------------------------------------
    /*
     * {380, 381, 382}: ht lưỡng long nhất thể xayda trái đất
     * {383, 384, 385}: ht porata xayda trái đất
     * {391, 392, 393}: ht namếc
     * {870, 871, 872}: ht c2 trái đất
     * {873, 874, 875}: ht c2 namếc
     * {867, 878, 869}: ht c2 xayda
     */
    private static final short[][] idOutfitFusion = {
            // trái đất - namek - xayda

            { 962, 963, 964 }, { 391, 392, 393 }, { 1610, 1611, 1612 }, // hợp thể zin
            { 380, 381, 382 }, {383, 384, 385 }, { 873, 874, 875 }, // cấp 1
// {380, 381, 382}, {391, 392, 393}, {380, 381, 382},
// {383, 384, 385}, {391, 392, 393},{383, 384, 385}, 
           { 1706, 1707, 1708 }, { 1000, 1001, 1002 },  { 1700, 1701, 1702 }, // cấp 2
            { 1643, 1644, 1645 }, { 1604, 1605, 1606 }, { 1646, 1647, 1648 }, // hợp thể cấp 3
        //    { 1700, 1701, 1702 }, { 1703, 1704, 1705 }, { 1706, 1707, 1708 },
            { 1709, 1710, 1711 }, { 1712, 1713, 1714 }, { 1715, 1716, 1717 }, };

    public static final short[][] idOutfitGod = {
            { -1, 472, 473 }, { -1, 476, 477 }, { -1, 474, 475 }
    };

    public static final short[][][] idOutfitHalloween = {
            {
                    { 545, 548, 549 }, { 547, 548, 549 }, { 546, 548, 549 }
            },
            {
                    { 760, 761, 762 }, { 760, 761, 762 }, { 760, 761, 762 }
            },
            {
                    { 654, 655, 656 }, { 654, 655, 656 }, { 654, 655, 656 }
            },
            {
                    { 651, 652, 653 }, { 651, 652, 653 }, { 651, 652, 653 }
            },
            {
                    { 651, 652, 653 }, { 651, 652, 653 }, { 651, 652, 653 }
            }
    };

    public static final short[][] idOutfitMafuba = {
            { 1221, 1222, 1223 }, { -1, -1, -1 }, { 1218, 1219, 1220 }
    };

    public int getHat() {
        return -1;
    }

    public byte getAura() {
        // if (!isPl() || this.Cards.isEmpty()) {
        // return -1;
        // }
        if (this.effectSkill != null && this.effectSkill.isSuper) {
            return idAuraSuper[gender][(playerSkill.getSkillbyId(gender == 0 ? 27 : gender == 1 ? 28 : 29).point - 1)
                    - numUseSkill];
        }
        if (isPl() && this.inventory != null && this.inventory.itemsBody.size() > 5) {
            Item item = this.inventory.itemsBody.get(5);
            if (!item.isNotNullItem()) {
                return auraPower();
            }
            switch (item.template.id) {
//   case 1652:
//                     return 12;
//   case 1636:
//                     return 34;
//   case 1656:
//                     return 1;
//   case 1654:
//                     return 58;
//   case 1653:
//                     return 14;
//                     case 1330:
//                     return 32;

//                 case 1735:
//                     return 83;
//                 case 1087:
//                     return 26;
//                 case 1088:
//                     return 25;
//                 case 1089:
//                     return 25;
//                 case 1090:
//                     return 27;
//                 case 1091:
//                     return 24;
//                 case 1208:
//                     return 6;
//                 case 1209:
//                     return 7;
//                 case 1210:
//                     return 13;
//                 case 1234:
//                     return 29;
//                 case 1235:
//                     return 29;
//                 case 1236:
//                     return 29;
//                 case 1255:
//                     return 19;
//                 case 1257:
//                     return 19;
//                 case 1265:
//                     return 19;
//                 case 1274:
//                     return 10;
//                 case 1275:
//                     return 10;
//                 case 1276:
//                     return 60;
//                 case 1277:
//                     return 61;
//                 case 1302:
//                     return 5;
//                       case 1754:
//                     return 77;
//                      case 1674:
//                     return 9;
//                 case 1309:
//                     return 2;
//                 case 1371:
//                     return 3;
//                 case 1378:
//                     return 3;
//                 case 1380:
//                     return 3;
//                 case 1381:
//                     return 3;
//                 case 1383:
//                     return 57;
//                 case 1384:
//                     return 57;
//                 case 1385:
//                     return 57;
//                 case 1391:
//                     return 57;
//                 case 1396:
//                     return 36;
//                 case 1424:
//                     return 37;
//                 case 1436:
//                     return 37;
//                 case 1437:
//                     return 37;
//                 case 1442:
//                     return 31;
//                 case 1450:
//                     return 59;
//                 case 1454:
//                     return 37;
//                 case 1632:
//                     return 56;
//                 case 1657:
//                     return 56;
//                 case 1587:
//                     return 18;
//                 case 1590:
//                     return 4;
//                 case 1593:
//                     return 25;
//                 case 1595:
//                     return 57;
//                 case 1658:
//                     return 88;
//                 case 1659:
//                     return 88;
//                 case 1660:
//                     return 88;
//                 case 1661:
//                     return 88;
//                 case 1662:
//                     return 88;
//                 case 1663:
//                     return 88;
//                 case 1667:
//                     return 88;
               
//                 case 1684:
//                     return 86;
//                 case 1685:
//                     return 85;
//                 case 1693:
//                     return 84;
//                 case 1553:
//                     return 84;
//                 case 1697:
//                     return 91;
//                 case 1698:
//                     return 90;
//                 case 1700:
//                     return 90;
//                 case 1761:
//                     return 28;
//                 case 1762:
//                     return 63;
            }
        }
        for (Card card : this.Cards) {
            if (card != null // && (card.Id == 956 || card.Id == 1142) && card.Level > 1
            ) {
                RadarCard radarTemplate = RadarService.gI().RADAR_TEMPLATE.stream().filter(r -> r.Id == card.Id)
                        .findFirst().orElse(null);
                if (radarTemplate != null) {
                    return (byte) radarTemplate.AuraId;
                }
            }
        }
        return auraPower();
    }

    public byte auraPower() {
        if (this.nPoint == null) {
            return 0;
        }
        if (this.nPoint.power >= 70_000_000_000L) {
            return 32;
             } else if (this.nPoint.power >= 60_000_000_000L) {
            return 59;
        } else if (this.nPoint.power >= 50_000_000_000L) {
            return 56;
        } else if (this.nPoint.power >= 40_000_000_000L) {
            return 20;
        } else if (this.nPoint.power >= 30_000_000_000L) {
            return 18;
        } else if (this.nPoint.power >= 20_000_000_000L) {
            return 17; } 
            else if (this.nPoint.power >= 1_000_000_000L) {
            return 7;
        }
        return -1;
    }

    public byte getEffFront() {
        if (this.inventory == null) {
            return -1;
        }
        if (this.inventory.itemsBody.isEmpty() || this.inventory.itemsBody.size() < 10) {
            return -1;
        }
        int levelAo = 0;
        Item.ItemOption optionLevelAo = null;
        int levelQuan = 0;
        Item.ItemOption optionLevelQuan = null;
        int levelGang = 0;
        Item.ItemOption optionLevelGang = null;
        int levelGiay = 0;
        Item.ItemOption optionLevelGiay = null;
        int levelNhan = 0;
        Item.ItemOption optionLevelNhan = null;
        Item itemAo = this.inventory.itemsBody.get(0);
        Item itemQuan = this.inventory.itemsBody.get(1);
        Item itemGang = this.inventory.itemsBody.get(2);
        Item itemGiay = this.inventory.itemsBody.get(3);
        Item itemNhan = this.inventory.itemsBody.get(4);
        for (Item.ItemOption io : itemAo.itemOptions) {
            if (io.optionTemplate.id == 72) {
                levelAo = io.param;
                optionLevelAo = io;
                break;
            }
        }
        for (Item.ItemOption io : itemQuan.itemOptions) {
            if (io.optionTemplate.id == 72) {
                levelQuan = io.param;
                optionLevelQuan = io;
                break;
            }
        }
        for (Item.ItemOption io : itemGang.itemOptions) {
            if (io.optionTemplate.id == 72) {
                levelGang = io.param;
                optionLevelGang = io;
                break;
            }
        }
        for (Item.ItemOption io : itemGiay.itemOptions) {
            if (io.optionTemplate.id == 72) {
                levelGiay = io.param;
                optionLevelGiay = io;
                break;
            }
        }
        for (Item.ItemOption io : itemNhan.itemOptions) {
            if (io.optionTemplate.id == 72) {
                levelNhan = io.param;
                optionLevelNhan = io;
                break;
            }
        }
        if (optionLevelAo != null && optionLevelQuan != null && optionLevelGang != null && optionLevelGiay != null
                && optionLevelNhan != null
                && levelAo >= 8 && levelQuan >= 8 && levelGang >= 8 && levelGiay >= 8 && levelNhan >= 8) {
            return 8;
        } else if (optionLevelAo != null && optionLevelQuan != null && optionLevelGang != null
                && optionLevelGiay != null && optionLevelNhan != null
                && levelAo >= 7 && levelQuan >= 7 && levelGang >= 7 && levelGiay >= 7 && levelNhan >= 7) {
            return 7;
        } else if (optionLevelAo != null && optionLevelQuan != null && optionLevelGang != null
                && optionLevelGiay != null && optionLevelNhan != null
                && levelAo >= 6 && levelQuan >= 6 && levelGang >= 6 && levelGiay >= 6 && levelNhan >= 6) {
            return 6;
        } else if (optionLevelAo != null && optionLevelQuan != null && optionLevelGang != null
                && optionLevelGiay != null && optionLevelNhan != null
                && levelAo >= 5 && levelQuan >= 5 && levelGang >= 5 && levelGiay >= 5 && levelNhan >= 5) {
            return 5;
        } else if (optionLevelAo != null && optionLevelQuan != null && optionLevelGang != null
                && optionLevelGiay != null && optionLevelNhan != null
                && levelAo >= 4 && levelQuan >= 4 && levelGang >= 4 && levelGiay >= 4 && levelNhan >= 4) {
            return 4;
        } else {
            return -1;
        }
    }

    private static final short[][] idOutFitSuperEarth = {
            { 1436, 1437, 1438 }, // level 1

            { 1436, 1437, 1438 }, // level 2

            { 1442, 1437, 1438 }, // level 3

            { 1440, 1437, 1438 }, // level 4

            { 1439, 1437, 1438 }, // level 5

            { 1441, 1437, 1438 }, // level 6
    };

    private static final short[][] idOutFitSuperNamec = {
            { 1430, 1431, 1432 }, // level 1

            { 1443, 1431, 1432 }, // level 2

            { 1444, 1431, 1432 }, // level 3

            { 1445, 1431, 1432 }, // level 4

            { 1446, 1431, 1432 }, // level 5

            { 1447, 1431, 1432 }, // level 6
    };

    private static final short[][] idOutFitSuperSaiyan = {
            { 1433, 1434, 1435 }, // level 1

            { 1433, 1434, 1435 }, // level 2

            { 1448, 1434, 1435 }, // level 3

            { 1449, 1434, 1435 }, // level 4

            { 1450, 1434, 1435 }, // level 5

            { 1451, 1434, 1435 }, // level 6
    };

    private static final byte[][] idAuraSuper = {
            { 20, 21, 22, 23, 24, 25 }, // Trái đất

            { 26, 27, 28, 29, 30, 31 }, // namec

            { 32, 33, 34, 35, 36, 37 },// xayda
    };

    public short getHeadThuCung() {
        if (this.isPl() && this.inventory != null && this.inventory.itemsBody.size() > 7
                && this.inventory.itemsBody.get(7).isNotNullItem()) {
            return (short) (this.inventory.itemsBody.get(7).template.head);
        }
        return -1;
    }

    public short getBodyThuCung() {
        if (this.isPl() && this.inventory != null && this.inventory.itemsBody.size() > 7
                && this.inventory.itemsBody.get(7).isNotNullItem()) {
            return (short) (this.inventory.itemsBody.get(7).template.body);
        }
        return -1;
    }

    public short getLegThuCung() {
        if (this.isPl() && this.inventory != null && this.inventory.itemsBody.size() > 7
                && this.inventory.itemsBody.get(7).isNotNullItem()) {
            return (short) (this.inventory.itemsBody.get(7).template.leg);
        }
        return -1;
    }

    public short getHeadSuper() {
        switch (gender) {
            case 0:
                return idOutFitSuperEarth[(playerSkill.getSkillbyId(27).point - 1) - numUseSkill][0];
            case 1:
                return idOutFitSuperNamec[(playerSkill.getSkillbyId(28).point - 1) - numUseSkill][0];
            case 2:
                return idOutFitSuperSaiyan[(playerSkill.getSkillbyId(29).point - 1) - numUseSkill][0];
        }
        return -1;
    }

    public short getBodySuper() {
        switch (gender) {
            case 0:
                return idOutFitSuperEarth[(playerSkill.getSkillbyId(27).point - 1) - numUseSkill][1];
            case 1:
                return idOutFitSuperNamec[(playerSkill.getSkillbyId(28).point - 1) - numUseSkill][1];
            case 2:
                return idOutFitSuperSaiyan[(playerSkill.getSkillbyId(29).point - 1) - numUseSkill][1];
        }
        return -1;
    }

    public short getLegSuper() {
        switch (gender) {
            case 0:
                return idOutFitSuperEarth[(playerSkill.getSkillbyId(27).point - 1) - numUseSkill][2];
            case 1:
                return idOutFitSuperNamec[(playerSkill.getSkillbyId(28).point - 1) - numUseSkill][2];
            case 2:
                return idOutFitSuperSaiyan[(playerSkill.getSkillbyId(29).point - 1) - numUseSkill][2];
        }
        return -1;
    }

    public short getHead() {
        if (this.isPl() && this.pet != null && this.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                || this.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2) {
            Item item = inventory.itemsBody.get(5);
            Item petItem = pet.inventory.itemsBody.get(5);

            boolean hasItem1 = item.isNotNullItem() && (item.template.id == 1693 || item.template.id == 1553);
            boolean hasItem2 = petItem.isNotNullItem() && (petItem.template.id == 1693 || petItem.template.id == 1553);
            boolean sameItem = item.isNotNullItem() && petItem.isNotNullItem()
                    && item.template.id == petItem.template.id;
            if (hasItem1 && hasItem2 && !sameItem) {
                return 1578;
            }
        }

        if (effectSkill != null && effectSkill.isSuper) {
            return getHeadSuper();
        } else if (effectSkill != null && effectSkill.isBinh) {
            return idOutfitMafuba[effectSkill.typeBinh][0];
        }
        if (effectSkill != null && effectSkill.isStone) {
            return 454;
        }
        if (effectSkill != null && effectSkill.isHalloween) {
            return idOutfitHalloween[effectSkill.idOutfitHalloween][this.gender][0];
        } else if (this.itemTime != null && this.itemTime.isUseKhauTrang) {
            if (this.gender == 1) {
                return 777;
            } else if (this.gender == 2) {
                return 776;
            }
            return 775;
        }
        if (effectSkill != null && effectSkill.isMonkey) {
            return (short) ConstPlayer.HEADMONKEY[effectSkill.levelMonkey - 1];
        } else if (effectSkill != null && effectSkill.isSocola) {
            return 412;
        } else if (fusion != null && fusion.typeFusion != ConstPlayer.NON_FUSION) {
            if (fusion.typeFusion == ConstPlayer.LUONG_LONG_NHAT_THE) {
                return idOutfitFusion[this.gender][0];
            } else if (fusion.typeFusion == ConstPlayer.HOP_THE_PORATA) {
                return idOutfitFusion[3 + this.gender][0];
            } else if (fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2) {
                return idOutfitFusion[6 + this.gender][0];
            } else if (fusion.typeFusion == 9) {
                return idOutfitFusion[9 + this.gender][0];
            } else if (fusion.typeFusion == 10) {
                return idOutfitFusion[12 + this.gender][0];
            }
        } else if (inventory != null && inventory.itemsBody.get(5).isNotNullItem()) {
            int headId = inventory.itemsBody.get(5).template.head;
            if (headId != -1) {
                return (short) headId;
            }
        }
        return this.head;
    }

    public short getBody() {
        if (this.isPl() && this.pet != null && this.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                || this.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2) {
            Item item = inventory.itemsBody.get(5);
            Item petItem = pet.inventory.itemsBody.get(5);

            boolean hasItem1 = item.isNotNullItem() && (item.template.id == 1693 || item.template.id == 1553);
            boolean hasItem2 = petItem.isNotNullItem() && (petItem.template.id == 1693 || petItem.template.id == 1553);
            boolean sameItem = item.isNotNullItem() && petItem.isNotNullItem()
                    && item.template.id == petItem.template.id;
            if (hasItem1 && hasItem2 && !sameItem) {
                return 1581;
            }
        }
        if (effectSkill != null && effectSkill.isSuper) {
            return getBodySuper();
        } else if (effectSkill != null && effectSkill.isBinh) {
            return idOutfitMafuba[effectSkill.typeBinh][1];
        }
        if (effectSkill != null && effectSkill.isStone) {
            return 455;
        }
        if (effectSkill != null && effectSkill.isHalloween) {
            return idOutfitHalloween[effectSkill.idOutfitHalloween][this.gender][1];
        }
        if (effectSkill != null && effectSkill.isMonkey) {
            return 193;
        } else if (effectSkill != null && effectSkill.isSocola) {
            return 413;
        } else if (isPhuHoMapMabu && fusion != null && fusion.typeFusion == ConstPlayer.NON_FUSION) {
            return idOutfitGod[this.gender][1];
        } else if (fusion != null && fusion.typeFusion != ConstPlayer.NON_FUSION) {
            if (fusion.typeFusion == ConstPlayer.LUONG_LONG_NHAT_THE) {
                return idOutfitFusion[this.gender][1];
            } else if (fusion.typeFusion == ConstPlayer.HOP_THE_PORATA) {
                return idOutfitFusion[3 + this.gender][1];
            } else if (fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2) {
                return idOutfitFusion[6 + this.gender][1];
            } else if (fusion.typeFusion == 9) {
                return idOutfitFusion[9 + this.gender][1];
            } else if (fusion.typeFusion == 10) {
                return idOutfitFusion[12 + this.gender][1];
            }
        } else if (inventory != null && inventory.itemsBody.get(5).isNotNullItem()) {
            int body = inventory.itemsBody.get(5).template.body;
            if (body != -1) {
                return (short) body;
            }
        }
        if (inventory != null && inventory.itemsBody.get(0).isNotNullItem()) {
            return inventory.itemsBody.get(0).template.part;
        }
        return (short) (gender == ConstPlayer.NAMEC ? 59 : 57);
    }

    public short getLeg() {
        if (this.isPl() && this.pet != null && this.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                || this.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2) {
            Item item = inventory.itemsBody.get(5);
            Item petItem = pet.inventory.itemsBody.get(5);

            boolean hasItem1 = item.isNotNullItem() && (item.template.id == 1693 || item.template.id == 1553);
            boolean hasItem2 = petItem.isNotNullItem() && (petItem.template.id == 1693 || petItem.template.id == 1553);
            boolean sameItem = item.isNotNullItem() && petItem.isNotNullItem()
                    && item.template.id == petItem.template.id;
            if (hasItem1 && hasItem2 && !sameItem) {
                return 1582;
            }
        }
        if (effectSkill != null && effectSkill.isSuper) {
            return getLegSuper();
        } else if (effectSkill != null && effectSkill.isBinh) {
            return idOutfitMafuba[effectSkill.typeBinh][2];
        }
        if (effectSkill != null && effectSkill.isStone) {
            return 456;
        }
        if (effectSkill != null && effectSkill.isHalloween) {
            return idOutfitHalloween[effectSkill.idOutfitHalloween][this.gender][2];
        }
        if (effectSkill != null && effectSkill.isMonkey) {
            return 194;
        } else if (effectSkill != null && effectSkill.isSocola) {
            return 414;
        } else if (isPhuHoMapMabu && fusion != null && fusion.typeFusion == ConstPlayer.NON_FUSION) {
            return idOutfitGod[this.gender][2];
        } else if (fusion != null && fusion.typeFusion != ConstPlayer.NON_FUSION) {
            if (fusion.typeFusion == ConstPlayer.LUONG_LONG_NHAT_THE) {
                return idOutfitFusion[this.gender][2];
            } else if (fusion.typeFusion == ConstPlayer.HOP_THE_PORATA) {
                return idOutfitFusion[3 + this.gender][2];
            } else if (fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2) {
                return idOutfitFusion[6 + this.gender][2];
            } else if (fusion.typeFusion == 9) {
                return idOutfitFusion[9 + this.gender][2];
            } else if (fusion.typeFusion == 10) {
                return idOutfitFusion[12 + this.gender][2];
            }
        } else if (inventory != null && inventory.itemsBody.get(5).isNotNullItem()) {
            int leg = inventory.itemsBody.get(5).template.leg;
            if (leg != -1) {
                return (short) leg;
            }
        }
        if (inventory != null && inventory.itemsBody.get(1).isNotNullItem()) {
            return inventory.itemsBody.get(1).template.part;
        }
        return (short) (gender == 1 ? 60 : 58);
    }

    public short getFlagBag() {
        if (this.iDMark != null && this.iDMark.isHoldBlackBall()) {
            return 31;
        } else if (this.idNRNM >= 353 && this.idNRNM <= 359) {
            return 30;
        }
        if (TaskService.gI().getIdTask(this) == ConstTask.TASK_3_2) {
            return 28;
        }
        if (this.inventory != null && this.inventory.itemsBody.size() >= 11) {
            if (this.inventory.itemsBody.get(8).isNotNullItem()) {
                return this.inventory.itemsBody.get(8).template.part;
            }
        }
        // if (this.isPet && this.inventory.itemsBody.size() >= 8) {
        // if (this.inventory.itemsBody.get(7).isNotNullItem()) {
        // return this.inventory.itemsBody.get(7).template.part;
        // }
        // }
        if (this.clan != null) {
            return (short) this.clan.imgId;
        }
        return -1;
    }

    public short getMount() {
        if (this.inventory.itemsBody.isEmpty() || this.inventory.itemsBody.size() < 10) {
            return -1;
        }
        Item item = this.inventory.itemsBody.get(7);
        if (item == null || !item.isNotNullItem()) {
            return -1;
        }
        if (item.template.type == 23 || item.template.type == 24) {
            if (item.template.gender == 3 || item.template.gender == this.gender) {
                return item.template.id;
            } else {
                return -1;
            }
        } else {
            if (item.template.id < 500) {
                return item.template.id;
            } else {
                return -1;
            }
        }
    }

    // --------------------------------------------------------------------------
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (!this.isDie()) {
            if (plAtt != null && !plAtt.equals(this)) {
                setTemporaryEnemies(plAtt);
            }
            if (this.isBattu) {
                return 0;
            }
            // if (plAtt != null && this.isPet && ((Pet) this).master.id == plAtt.id) {
            // if (this.effectSkill != null && !this.effectSkill.isHalloween) {
            // EffectSkillService.gI().setIsHalloween(this, -1, 1800000);
            // }
            // }
            if (plAtt != null && plAtt.playerSkill.skillSelect != null && !plAtt.isBoss
                    && MapService.gI().isMapMaBu(this.zone.map.mapId)) {
                switch (plAtt.playerSkill.skillSelect.template.id) {
                    case Skill.KAMEJOKO, Skill.MASENKO, Skill.ANTOMIC, Skill.DRAGON, Skill.DEMON, Skill.GALICK,
                            Skill.LIEN_HOAN, Skill.KAIOKEN ->
                        damage = Util.maxIntValue(damage > this.nPoint.hpMax / 20 ? this.nPoint.hpMax / 20 : damage);
                }
            }
            if (plAtt != null && plAtt.isBoss) {
                this.effectSkin.isVoHinh = false;
                this.effectSkin.lastTimeVoHinh = System.currentTimeMillis();
            }
            if (plAtt != null && plAtt.effectSkill != null && plAtt.effectSkill.isBinh
                    && !Util.canDoWithTime(plAtt.effectSkill.lastTimeUpBinh, 3000)) {
                return 0;
            }
            if (plAtt != null && plAtt.isPl() && this.maBuHold != null && this.zone != null
                    && this.zone.map.mapId == 128) {
                this.precentMabuHold++;
                damage = 1;
            }
            if (plAtt != null && this.nPoint.islinhthuydanhbac) {
                Service.gI().sendThongBao(plAtt, "Không thể tấn công! Vì người chơi này đã nạp lần đầu!");
                return 0;
            }
            if (plAtt != null && plAtt.idNRNM != -1 && (this.isBoss || this.isNewPet)) {
                return 1;
            }
            if (plAtt != null && (plAtt.idNRNM != -1 || this.idNRNM != -1) && plAtt.clan != null && this.clan != null
                    && plAtt.clan == this.clan) {
                Service.gI().chatJustForMe(plAtt, this, "Ê cùng bang mà");
                return 0;
            }
            if (!Util.canDoWithTime(this.lastTimeRevived, 1500)) {
                return 0;
            }

            if (plAtt != null && plAtt.playerSkill.skillSelect != null) {
                switch (plAtt.playerSkill.skillSelect.template.id) {
                    case Skill.KAMEJOKO, Skill.MASENKO, Skill.ANTOMIC -> {
                        if (this.nPoint.voHieuChuong > 0) {
                            services.PlayerService.gI().hoiPhuc(this, 0,
                                    Util.maxIntValue(damage * this.nPoint.voHieuChuong / 100));
                            return 0;
                        }
                    }
                }
            }

            int tlGiap = this.nPoint.tlGiap;
            int tlNeDon = this.nPoint.tlNeDon;

            if (plAtt != null && !isMobAttack && plAtt.playerSkill.skillSelect != null) {
                switch (plAtt.playerSkill.skillSelect.template.id) {
                    case Skill.KAMEJOKO, Skill.MASENKO, Skill.ANTOMIC, Skill.DRAGON, Skill.DEMON, Skill.GALICK,
                            Skill.LIEN_HOAN, Skill.KAIOKEN, Skill.QUA_CAU_KENH_KHI, Skill.MAKANKOSAPPO,
                            Skill.DICH_CHUYEN_TUC_THOI ->
                        tlNeDon -= plAtt.nPoint.tlchinhxac;
                    default ->
                        tlNeDon = 0;
                }

                switch (plAtt.playerSkill.skillSelect.template.id) {
                    case Skill.KAMEJOKO, Skill.MASENKO, Skill.ANTOMIC -> {
                        if (tlGiap - plAtt.nPoint.tlxgc >= 0) {
                            tlGiap -= plAtt.nPoint.tlxgc;
                        } else {
                            tlGiap = 0;
                        }
                    }
                    case Skill.DRAGON, Skill.DEMON, Skill.GALICK, Skill.LIEN_HOAN, Skill.KAIOKEN -> {
                        if (tlGiap - plAtt.nPoint.tlxgcc >= 0) {
                            tlGiap -= plAtt.nPoint.tlxgcc;
                        } else {
                            tlGiap = 0;
                        }
                    }
                }
            }

            if (piercing) {
                tlGiap = 0;
            }

            if (tlNeDon > 90) {
                tlNeDon = 90;
            }
            if (tlGiap > 86) {
                tlGiap = 86;
            }

            if (Util.isTrue(tlNeDon, 100)) {
                return 0;
            }

            damage -= Util.maxIntValue((damage / 100) * tlGiap);

            if (!piercing) {
                damage = this.nPoint.subDameInjureWithDeff(damage);
            }

            boolean isUseGX = false;
            if (!piercing && plAtt != null && plAtt.playerSkill.skillSelect != null) {
                switch (plAtt.playerSkill.skillSelect.template.id) {
                    case Skill.KAMEJOKO, Skill.MASENKO, Skill.ANTOMIC, Skill.DRAGON, Skill.DEMON, Skill.GALICK,
                            Skill.LIEN_HOAN, Skill.KAIOKEN, Skill.QUA_CAU_KENH_KHI, Skill.MAKANKOSAPPO,
                            Skill.DICH_CHUYEN_TUC_THOI ->
                        isUseGX = true;
                }
            }
            if ((isUseGX || isMobAttack) && this.itemTime != null) {
                if (this.itemTime.isUseGiapXen && !this.itemTime.isUseGiapXen2) {
                    damage /= 2;
                }
                if (this.itemTime.isUseGiapXen2) {
                    damage = damage / 100 * 40;
                }
            }

            if (!piercing && effectSkill.isShielding && !isMobAttack) {
                if (this.iDMark != null) {
                    this.iDMark.setDamePST(Util.maxIntValue(damage));
                }
                if (damage > nPoint.hpMax) {
                    EffectSkillService.gI().breakShield(this);
                }
                damage = 1;
                if (MapService.gI().isMapPhoBan(this.zone.map.mapId)) {
                    damage = 10;
                }
            }
            // damage = Math.min(damage, 2_000_000_000L); // Giới hạn damage không vượt quá
            // 2 tỷ
            if (!piercing && plAtt == null && isMobAttack && (this.charms.tdBatTu > System.currentTimeMillis()
                    || this.effectSkill != null && this.effectSkill.isHalloween) && damage >= this.nPoint.hp) {
                damage = this.nPoint.hp - 1;
            }

            // if (this.zone.map.mapId == 129) {
            //     if (damage >= this.nPoint.hp) {
            //         this.lostByDeath = true;
            //         The23rdMartialArtCongress mc = The23rdMartialArtCongressManager.gI().getMC(zone);
            //         if (mc != null) {
            //             mc.die();
            //         }
            //         return 0;
            //     }
            // }
            if (this.zone.map.mapId == 51) {
                this.totalDamageTaken += damage;
            }
            this.nPoint.subHP(Util.maxIntValue(damage));
            if ((plAtt != null || isMobAttack) && isDie() && !isBoss && !isNewPet && !isNewPet1) {
                if (Util.isTrue(this.nPoint.tlBom, 100)) {
                    setBom(plAtt);
                } else {
                    setDie(plAtt);
                }
            }

            return damage;
        } else {
            return 0;
        }
    }

    public void setTemporaryEnemies(Player pl) {
        if (!temporaryEnemies.contains(pl)) {
            temporaryEnemies.add(pl);
        }
    }

    protected void setBom(Player plAtt) {
        setDie(plAtt);
        SkillService.gI().sendPlayerPrepareBom(plAtt, 2000);
        // Service.gI().callClone(this);
    }

    public void kill(Player pl) {
        pl.injured(this, Util.maxIntValue(pl.nPoint.hpMax), false, false);
        PlayerService.gI().sendInfoHpMpMoney(this);
        Service.gI().Send_Info_NV(this);
    }

    public void setDie() {
        this.setDie(null);
    }

    protected void setDie(Player plAtt) {
        TaskService.gI().checkDoneTaskKillPlayer(plAtt);
        if (this.isPl()) {
            if (this.zone != null && this.zone.map != null && this.zone.map.mapId == War.MAP_WAR_ID) {
                int itemDropId = 457;
                int soLuongRoi = 10;
                int tongSoItem = 0;
                Item item457 = null;
                for (Item item : this.inventory.itemsBag) {
                    if (item != null && item.template != null && item.template.id == itemDropId) {
                        tongSoItem += item.quantity;
                        if (item457 == null) {
                            item457 = item;
                        }
                    }
                }
                if (tongSoItem >= soLuongRoi && item457 != null) {
                    InventoryService.gI().subQuantityItemsBag(this, item457, soLuongRoi);
                    InventoryService.gI().sendItemBag(this);
                    for (int i = 0; i < soLuongRoi; i++) {
                        int randX = this.location.x + Util.nextInt(-30, 30);
                        int y = this.zone.map.yPhysicInTop(randX, this.location.y - 24);
                        ItemMap itemMap = new ItemMap(this.zone, itemDropId, 1, randX, y, -1);
                        Service.gI().dropItemMap(this.zone, itemMap);
                    }
                }
            } else {
                long vangtru = this.nPoint.power / 1000000;
                if (vangtru > 32000) {
                    vangtru = 32000;
                }
                int vang = (int) vangtru - Util.nextInt(10, 100);
                if (this.inventory.gold >= vang && vang >= 1) {
                    this.inventory.gold -= vang;
                    Service.gI().sendMoney(this);
                    vang = vang * 95 / 100;
                    if (vang < 10000) {
                        Service.gI().dropItemMap(this.zone, new ItemMap(zone, 189, vang, this.location.x,
                                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24), this.id));
                    } else if (vang < 20000) {
                        Service.gI().dropItemMap(this.zone, new ItemMap(zone, 188, vang, this.location.x,
                                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24), this.id));
                    } else {
                        Service.gI().dropItemMap(this.zone, new ItemMap(zone, 190, vang, this.location.x,
                                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24), this.id));
                    }
                }
            }
        }

        // xóa phù
        if (this.effectSkin.xHPKI > 1) {
            this.effectSkin.xHPKI = 1;
            Service.gI().point(this);
        }
        if (this.effectSkin.xDame > 1) {
            this.effectSkin.xDame = 1;
            Service.gI().point(this);
        }
        // xóa tụ skill đặc biệt
        this.playerSkill.prepareQCKK = false;
        this.playerSkill.prepareLaze = false;
        this.playerSkill.prepareTuSat = false;
        // xóa hiệu ứng skill
        this.effectSkill.removeSkillEffectWhenDie();
        //
        nPoint.setHp(Util.maxIntValue(0));
        nPoint.setMp(Util.maxIntValue(0));
        // xóa trứng
        if (this.mobMe != null) {
            this.mobMe.mobMeDie();
            this.mobMe.dispose();
            this.mobMe = null;
        }
        Service.gI().charDie(this);
        // add kẻ thù
        if (!this.isPet && !this.isNewPet && !this.isNewPet1 && !this.isBoss && plAtt != null && !plAtt.isPet
                && !plAtt.isNewPet && !plAtt.isNewPet1 && !plAtt.isBoss && !plAtt.isBot) {
            if (!plAtt.itemTime.isUseAnDanh) {
                FriendAndEnemyService.gI().addEnemy(this, plAtt);
            }
        }
        // kết thúc pk

        this.typePk = 0;

        if (this.pvp != null && this.zone.map.mapId != 140) {
            this.pvp.lose(this, TYPE_LOSE_PVP.DEAD);
        }

        BlackBallWarService.gI().dropBlackBall(this);
        NgocRongNamecService.gI().dropNamekBall(this);
        // if (isHoldNamecBallTranhDoat) {
        //     TranhNgocService.getInstance().dropBall(this, (byte) -1);
        //     TranhNgocService.getInstance().sendUpdateLift(this);
        // }
    }

    // --------------------------------------------------------------------------
    public void setClanMember() {
        if (this.clanMember != null) {
            this.clanMember.powerPoint = this.nPoint.power;
            this.clanMember.head = this.getHead();
            this.clanMember.body = this.getBody();
            this.clanMember.leg = this.getLeg();
        }
    }

    public boolean isAdmin() {
        return this.session != null && this.session.isAdmin;
    }

    public void thuhoivp() {
        // Xử lý itemsBag
        Iterator<Item> bagIterator = inventory.itemsBag.iterator();
        while (bagIterator.hasNext()) {
            Item item = bagIterator.next();
            if (item.isNotNullItem() && item.itemOptions != null && !item.itemOptions.isEmpty()) {
                for (Item.ItemOption io : item.itemOptions) {
                    if (io.optionTemplate.id == 14 && io.param > 20) {
                        bagIterator.remove();
                        InventoryService.gI().removeItemBag(this, item);
                        InventoryService.gI().sendItemBag(this);
                        Service.gI().sendThongBao(this, "Đã thu hồi Vật phẩm vì gây lỗi game!");
                        break; // Dừng vòng lặp khi đã xóa phần tử
                    }
                }
            }
        }

        // Xử lý itemsBody
        Iterator<Item> bodyIterator = inventory.itemsBody.iterator();
        while (bodyIterator.hasNext()) {
            Item item = bodyIterator.next();
            if (item.isNotNullItem() && item.itemOptions != null && !item.itemOptions.isEmpty()) {
                for (Item.ItemOption io : item.itemOptions) {
                    if (io.optionTemplate.id == 14 && io.param > 20) {
                        bodyIterator.remove();
                        InventoryService.gI().removeItem(inventory.itemsBody, item);
                        InventoryService.gI().sendItemBody(this);
                        Service.gI().sendThongBao(this, "Đã thu hồi Vật phẩm vì gây lỗi game!");
                        break;
                    }
                }
            }
        }

        // Xử lý itemsBox
        Iterator<Item> boxIterator = inventory.itemsBox.iterator();
        while (boxIterator.hasNext()) {
            Item item = boxIterator.next();
            if (item.isNotNullItem() && item.itemOptions != null && !item.itemOptions.isEmpty()) {
                for (Item.ItemOption io : item.itemOptions) {
                    if (io.optionTemplate.id == 14 && io.param > 20) {
                        boxIterator.remove();
                        InventoryService.gI().removeItem(inventory.itemsBox, item);
                        InventoryService.gI().sendItemBox(this);
                        Service.gI().sendThongBao(this, "Đã thu hồi Vật phẩm vì gây lỗi game!");
                        break;
                    }
                }
            }
        }

        // Xử lý pet.itemsBody
        if (pet != null) {
            Iterator<Item> petBodyIterator = pet.inventory.itemsBody.iterator();
            while (petBodyIterator.hasNext()) {
                Item item = petBodyIterator.next();
                if (item.isNotNullItem() && item.itemOptions != null && !item.itemOptions.isEmpty()) {
                    for (Item.ItemOption io : item.itemOptions) {
                        if (io.optionTemplate.id == 14 && io.param > 20) {
                            petBodyIterator.remove();
                            InventoryService.gI().removeItem(pet.inventory.itemsBody, item);
                            InventoryService.gI().sendItemBag(pet);
                            Service.gI().sendThongBao(this, "Đã thu hồi Vật phẩm vì gây lỗi game!");
                            break;
                        }
                    }
                }
            }
        }
    }

    public void setJustRevivaled() {
        this.justRevived = true;
        this.lastTimeRevived = System.currentTimeMillis();
    }

    public boolean actived() {
        return (this.isPl() && this.session != null && this.session.actived)
                || (this.isPet && ((Pet) this).master.session != null && ((Pet) this).master.session.actived);
    }

    // public void sendNewPet() {
    // if (isPl() && inventory != null && inventory.itemsBody.get(7) != null) {
    // Item it = inventory.itemsBody.get(7);
    // if (it != null && it.isNotNullItem() && newPet == null) {
    // switch (it.template.id) {
    // case 942 -> {
    // PetService.Pet2(this, 966, 967, 968);
    // Service.gI().point(this);
    // }
    // case 943 -> {
    // PetService.Pet2(this, 969, 970, 971);
    // Service.gI().point(this);
    // }
    // case 944 -> {
    // PetService.Pet2(this, 972, 973, 974);
    // Service.gI().point(this);
    // }
    // case 967 -> {
    // PetService.Pet2(this, 1050, 1051, 1052);
    // Service.gI().point(this);
    // }
    // case 968 -> {
    // PetService.Pet2(this, 1183, 1184, 1185);
    // Service.gI().point(this);
    // }
    // }
    // }
    // }
    // }
    private void fixBlackBallWar() {
        int x = this.location.x;
        int y = this.location.y;
        switch (this.zone.map.mapId) {
            case 85, 86, 87, 88, 89, 90, 91 -> {
                if (this.isPl()) {
                    if (x < 24 || x > this.zone.map.mapWidth - 24 || y < 0 || y > this.zone.map.mapHeight - 24) {
                        if (MapService.gI().getWaypointPlayerIn(this) == null) {
                            Service.gI().resetPoint(this, x, this.zone.map.yPhysicInTop(this.location.x, 100));
                            this.nPoint.hp -= this.nPoint.hpMax / 10;
                            PlayerService.gI().sendInfoHp(this);
                            return;
                        }
                    }
                    int yTop = this.zone.map.yPhysicInTop(this.location.x, this.location.y);
                    if (yTop >= this.zone.map.mapHeight - 24) {
                        Service.gI().resetPoint(this, x, this.zone.map.yPhysicInTop(this.location.x, 100));
                        this.nPoint.hp -= this.nPoint.hpMax / 10;
                        PlayerService.gI().sendInfoHp(this);
                    }
                }
            }
        }
    }

    public void move(int _toX, int _toY) {
        if (_toX != this.location.x) {
            this.location.x = _toX;
        }
        if (_toY != this.location.y) {
            this.location.y = _toY;
        }
        MapService.gI().sendPlayerMove(this);
    }

    public void dispose() {
        if (itemsTradeWVP != null) {
            if (!itemsTradeWVP.isEmpty()) {
                for (Item item : itemsTradeWVP) {
                    InventoryService.gI().addItemBag(this, item);
                }
            }
            itemsTradeWVP.clear();
            itemsTradeWVP = null;
        }
        if (pet != null) {
            pet.dispose();
            pet = null;
        }
        if (newPet != null) {
            newPet.dispose();
            newPet = null;
        }
        if (mapBlackBall != null) {
            mapBlackBall.clear();
            mapBlackBall = null;
        }
        zone = null;
        mapBeforeCapsule = null;
        if (mapMaBu != null) {
            mapMaBu.clear();
            mapMaBu = null;
        }
        mapBeforeCapsule = null;
        if (mapCapsule != null) {
            mapCapsule.clear();
            mapCapsule = null;
        }
        if (mobMe != null) {
            mobMe.dispose();
            mobMe = null;
        }
        location = null;
        if (setClothes != null) {
            setClothes.dispose();
            setClothes = null;
        }
        if (effectSkill != null) {
            effectSkill.dispose();
            effectSkill = null;
        }
        if (mabuEgg != null) {
            mabuEgg.dispose();
            mabuEgg = null;
        }
        if (playerTask != null) {
            playerTask.dispose();
            playerTask = null;
        }
        if (itemTime != null) {
            itemTime.dispose();
            itemTime = null;
        }
        if (fusion != null) {
            fusion.dispose();
            fusion = null;
        }
        if (magicTree != null) {
            magicTree.dispose();
            magicTree = null;
        }
        if (playerIntrinsic != null) {
            playerIntrinsic.dispose();
            playerIntrinsic = null;
        }
        if (inventory != null) {
            inventory.dispose();
            inventory = null;
        }
        if (playerSkill != null) {
            playerSkill.dispose();
            playerSkill = null;
        }
        if (combine != null) {
            combine.dispose();
            combine = null;
        }
        if (iDMark != null) {
            iDMark.dispose();
            iDMark = null;
        }
        if (charms != null) {
            charms.dispose();
            charms = null;
        }
        if (effectSkin != null) {
            effectSkin.dispose();
            effectSkin = null;
        }
        if (nPoint != null) {
            nPoint.dispose();
            nPoint = null;
        }
        if (rewardBlackBall != null) {
            rewardBlackBall.dispose();
            rewardBlackBall = null;
        }
        if (pvp != null) {
            pvp.dispose();
            pvp = null;
        }
        if (superRank != null) {
            superRank.dispose();
            superRank = null;
        }
        if (dropItem != null) {
            dropItem.dispose();
            dropItem = null;
        }
        if (satellite != null) {
            satellite = null;
        }
        if (achievement != null) {
            achievement.dispose();
            achievement = null;
        }
        if (giftCode != null) {
            giftCode.dispose();
            giftCode = null;
        }
        if (traning != null) {
            traning = null;
        }
        if (mapCapsule != null) {
            mapCapsule.clear();
            mapCapsule = null;
        }
        if (Cards != null) {
            Cards.clear();
            Cards = null;
        }
        if (itemsWoodChest != null) {
            itemsWoodChest.clear();
            itemsWoodChest = null;
        }
        if (friends != null) {
            friends.clear();
            friends = null;
        }
        if (enemies != null) {
            enemies.clear();
            enemies = null;
        }
        if (temporaryEnemies != null) {
            temporaryEnemies.clear();
            temporaryEnemies = null;
        }
        itemsWoodChest = null;
        Cards = null;
        itemEvent = null;
        maBu2H = null;
        maBuHold = null;
        zoneSieuThanhThuy = null;
        thongBaoTapTuDong = null;
        notify = null;
        clan = null;
        clanMember = null;
        friends = null;
        enemies = null;
        session = null;
        newSkill = null;
        name = null;
        textThongBaoChangeMap = null;
        textThongBaoThua = null;
    }

}
