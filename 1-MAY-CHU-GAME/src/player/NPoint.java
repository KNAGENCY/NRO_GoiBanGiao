package player;

/*
 *
 *
 * @author GiaCatLuong
 */
import models.Card.Card;
import models.Card.OptionCard;
import consts.ConstPlayer;
import consts.ConstRatio;
import intrinsic.Intrinsic;
import item.Item;
import item.Item.ItemOption;
import skill.Skill;
import server.Manager;
import services.EffectSkillService;
import services.ItemService;
import services.MapService;
import services.PlayerService;
import services.Service;
import services.TaskService;
import utils.Logger;
import utils.SkillUtil;
import utils.Util;

import java.util.ArrayList;
import java.util.List;

import jdbc.daos.EventDAO;
import lombok.Setter;
import mob.Mob;
import power.PowerLimit;
import power.PowerLimitManager;
import utils.TimeUtil;

public class NPoint {

    public static final byte MAX_LIMIT = 13;

    @Setter
    private Player player;

    public NPoint(Player player) {
        this.player = player;
        this.tlHp = new ArrayList<>();
        this.tlMp = new ArrayList<>();
        this.tlDef = new ArrayList<>();
        this.tlDame = new ArrayList<>();
        this.tlDameAttMob = new ArrayList<>();
        this.tlTNSM = new ArrayList<>();
        this.tlDameCrit = new ArrayList<>();
    }

    public boolean isCrit;
    public boolean isCrit100;
    public boolean isCritTele;

    private Intrinsic intrinsic;
    private int percentDameIntrinsic;
    public long dameAfter;
    private PowerLimit powerLimit;
    /*-----------------------Chỉ số cơ bản------------------------------------*/
    public byte numAttack;
    public short stamina, maxStamina;

    public byte limitPower;
    public long power;
    public long tiemNang;

    public long hp, hpMax, hpg;
    public long mp, mpMax, mpg;
    public long dame, dameg;
    public int def, defg;
    public int crit, critg;
    public byte speed = 5;

    public boolean teleport;

    public boolean khangTDHS;

    public void initPowerLimit() {
        powerLimit = PowerLimitManager.getInstance().get(limitPower);
    }

    /**
     * Chỉ số cộng thêm
     */
    public int hpAdd, mpAdd, dameAdd, defAdd, critAdd, hpHoiAdd, mpHoiAdd;

    /**
     * //+#% sức đánh chí mạng
     */
    public List<Integer> tlDameCrit;
    public int tlSDCM;

    /**
     * Tỉ lệ hp, mp cộng thêm
     */
    public List<Integer> tlHp, tlMp;

    /**
     * Tỉ lệ giáp cộng thêm
     */
    public List<Integer> tlDef;

    /**
     * Tỉ lệ sức đánh/ sức đánh khi đánh quái
     */
    public List<Integer> tlDame, tlDameAttMob;

    /**
     * Lượng hp, mp hồi mỗi 30s, mp hồi cho người khác
     */
    public long hpHoi, mpHoi, mpHoiCute;

    /**
     * Tỉ lệ hp, mp hồi cộng thêm
     */
    public short tlHpHoi, tlMpHoi;

    /**
     * Tỉ lệ hp, mp hồi bản thân và đồng đội cộng thêm
     */
    public short tlHpHoiBanThanVaDongDoi, tlMpHoiBanThanVaDongDoi;

    /**
     * Tỉ lệ hút hp, mp khi đánh, hp khi đánh quái
     */
    public short tlHutHp, tlHutMp, tlHutHpMob;

    /**
     * Tỉ lệ hút hp, mp xung quanh mỗi 5s
     */
    public short tlHutHpMpXQ;
    // giảm dame xung quanh
    public short tlGiamDameXQ;
    public short tlPhanTamXQ;

    /**
     * Tỉ lệ phản sát thương
     */
    public short tlPST;

    /**
     * Tỉ lệ tiềm năng sức mạnh
     */
    public List<Integer> tlTNSM;

    /**
     * Tỉ lệ vàng cộng thêm
     */
    public short tlGold;

    /**
     * Tỉ lệ né đòn
     */
    public short tlNeDon;

    public short tlBom;

    public short tlGiap;

    public short tlxgcc;

    public short tlxgc;

    public short tlchinhxac;

    public short tlTNSMPet;
    public short xChuong;

    public short setltdb;
    public short setTinhAn;
    public short setNhatAn;
    public short setNguyetAn;

    /**
     * Tỉ lệ sức đánh đẹp cộng thêm cho bản thân và người xung quanh
     */
    public int tlSexyDame;

    /**
     * Tỉ lệ giảm sức đánh
     */
    public short tlSubSD;
    public short tlSubSDCurrent;

    public int voHieuChuong;

    /*------------------------Effect skin-------------------------------------*/
    public Item trainArmor;
    public boolean wearingTrainArmor;

    public boolean wearingVoHinh;
    public boolean isKhongLanh;
    public boolean isKhongLaoHoa;
    public boolean islinhthuydanhbac;
    public boolean isTinhAn;
    public boolean isNhatAn;
    public boolean isNguyetAn;
    public boolean isTanHinh;
    public boolean isHoaDa;
    public boolean isLamCham;
    public boolean isDoSPL;
    public boolean isThoBulma;
    public boolean isDietQuy;
    public boolean isBunmaTocMau;
    public boolean isTiecBaiBien;
    public short tlHpGiamODo;

    public boolean isGogeta;

    public int tlSpeed;

    public int levelBT;

    /*-------------------------------------------------------------------------*/
    /**
     * Tính toán mọi chỉ số sau khi có thay đổi
     */
    public void calPoint() {
        if (this.player.pet != null) {
            this.player.pet.nPoint.setPointWhenWearClothes();
        }
        this.setPointWhenWearClothes();
    }

    private void setPointWhenWearClothes() {
        resetPoint();
        if (this.player.rewardBlackBall.timeOutOfDateReward[2] > System.currentTimeMillis()) {
            tlHutHp += RewardBlackBall.R4S_1;
        }
        // if (this.player.rewardBlackBall.timeOutOfDateReward[3] >
        // System.currentTimeMillis()) {
        // tlPST += RewardBlackBall.R4S_2;
        // }
        if (this.player.rewardBlackBall.timeOutOfDateReward[4] > System.currentTimeMillis()) {
            tlDameCrit.add(RewardBlackBall.R5S_1);
            tlSDCM += RewardBlackBall.R5S_1;
        }
        // if (this.player.rewardBlackBall.timeOutOfDateReward[6] >
        // System.currentTimeMillis()) {
        // tlNeDon += RewardBlackBall.R7S_1;
        // }

        Card card = player.Cards.stream().filter(r -> r != null && r.Used == 1).findFirst().orElse(null);
        if (card != null) {
            for (OptionCard io : card.Options) {
                if (io.active == card.Level || (card.Level == -1 && io.active == 0)) {
                    switch (io.id) {
                        case 0: // Tấn công +#
                            this.dameAdd += io.param;
                            break;
                        case 2: // HP, KI+#000
                            this.hpAdd += io.param * 1000;
                            this.mpAdd += io.param * 1000;
                            break;
                        case 3:// vô hiệu chưởng
                            this.voHieuChuong += io.param;
                            break;
                        case 5: // +#% sức đánh chí mạng
                            this.tlDameCrit.add(io.param);
                            this.tlSDCM += io.param;
                            break;
                        case 6: // HP+#
                            this.hpAdd += io.param;
                            break;
                        case 7: // KI+#
                            this.mpAdd += io.param;
                            break;
                        case 8: // Hút #% HP, KI xung quanh mỗi 5 giây
                            this.tlHutHpMpXQ += io.param;
                            break;
                        case 10: // giảm dame
                        case 169:
                            this.tlGiamDameXQ += io.param;
                            break;
                        case 111:
                            this.tlPhanTamXQ += io.param > 1 ? io.param : 40;
                            break;
                        case 14: // Chí mạng+#%
                            this.critAdd += io.param;
                            break;
                        case 16: // Speed
                        case 114:
                        case 148:
                            this.tlSpeed += io.param;
                            break;
                        case 18: // Chinh xac
                            this.tlchinhxac += io.param;
                            break;
                        case 19: // Tấn công+#% khi đánh quái
                            this.tlDameAttMob.add(io.param);
                            break;
                        case 22: // HP+#K
                            this.hpAdd += io.param * 1000;
                            break;
                        case 23: // MP+#K
                            this.mpAdd += io.param * 1000;
                            break;
                        case 27: // +# HP/30s
                            this.hpHoiAdd += io.param;
                            break;
                        case 28: // +# KI/30s
                            this.mpHoiAdd += io.param;
                            break;
                        case 33: // dịch chuyển tức thời
                            this.teleport = true;
                            break;
                        case 34:
                            this.setTinhAn += 1;
                            break;
                        case 35:
                            this.setNguyetAn += 1;
                            break;
                        case 36:
                            this.setNhatAn += 1;
                            break;
                        case 47: // Giáp+#
                            this.defAdd += io.param;
                            break;
                        case 48: // HP/KI+#
                            this.hpAdd += io.param;
                            this.mpAdd += io.param;
                            break;
                        case 49: // Tấn công+#%
                        case 50: // Sức đánh+#%
                            this.tlDame.add(io.param);
                            break;
                        case 53: // Không ảnh hưởng bởi lão hóa
                            this.isKhongLaoHoa = true;
                            break;
                        case 77: // HP+#%
                            this.tlHp.add(io.param);
                            break;
                        case 80: // HP+#%/30s
                            this.tlHpHoi += io.param;
                            break;
                        case 81: // MP+#%/30s
                            this.tlMpHoi += io.param;
                            break;
                        case 88: // Cộng #% exp khi đánh quái
                            this.tlTNSM.add(io.param);
                            break;
                        case 94: // Giáp #%
                            this.tlGiap += io.param;
                            break;
                        case 95: // Biến #% tấn công thành HP
                            this.tlHutHp += io.param;
                            break;
                        case 96: // Biến #% tấn công thành MP
                            this.tlHutMp += io.param;
                            break;
                        case 97: // Phản #% sát thương
                            this.tlPST += io.param;
                            break;
                        case 98: // Xuyen giap chuong
                            this.tlxgc += io.param;
                            break;
                        case 99: // Xuyen giap can chien
                            this.tlxgcc += io.param;
                            break;
                        case 100: // +#% vàng từ quái
                            this.tlGold += io.param;
                            break;
                        case 101: // +#% TN,SM
                            this.tlTNSM.add(io.param);
                            break;
                        case 103: // KI +#%
                            this.tlMp.add(io.param);
                            break;
                        case 104: // Biến #% tấn công quái thành HP
                            this.tlHutHpMob += io.param;
                            break;
                        case 105: // Vô hình khi không đánh quái và boss
                            this.wearingVoHinh = true;
                            break;
                        case 106: // Không ảnh hưởng bởi cái lạnh
                            this.isKhongLanh = true;
                            break;
                        case 108: // #% Né đòn
                            this.tlNeDon += io.param;
                            break;
                        case 109: // Hôi, giảm #% HP
                            this.tlHpGiamODo += io.param;
                            break;
                        case 116: // Kháng thái dương hạ san
                            this.khangTDHS = true;
                            break;
                        case 153: // Kháng thái dương hạ san
                            this.tlBom += io.param;
                            break;
                        case 117: // Đẹp +#% SĐ cho mình và người xung quanh
                            if (io.param > this.tlSexyDame) {
                                this.tlSexyDame = io.param;
                            }
                            break;
                        case 147: // +#% sức đánh
                            this.tlDame.add(io.param);
                            break;
                        case 156: // Giảm 50% sức đánh, HP, KI và +#% SM, TN, vàng từ quái
                            this.tlSubSDCurrent += 50;
                            this.tlTNSM.add(io.param);
                            this.tlGold += io.param;
                            break;
                        case 160: // TNSM PET;
                            this.tlTNSMPet += io.param;
                            break;
                        case 162: // Cute hồi #% KI/s bản thân và xung quanh
                            this.mpHoiCute += io.param;
                            break;
                        case 173: // Phục hồi #% HP và KI cho đồng đội
                            this.tlHpHoiBanThanVaDongDoi += io.param;
                            this.tlMpHoiBanThanVaDongDoi += io.param;
                            break;

                    }
                }
            }
        }

        // Bông tai cấp 2
        if (this.player.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                || this.player.fusion.typeFusion == 9) {
            this.player.inventory.itemsBag.stream()
                    .filter(it -> it.isNotNullItem() && (it.template.id == 921 || it.template.id == 1499))
                    .findFirst().ifPresent(btc2 -> {
                        for (ItemOption io : btc2.itemOptions) {
                            addOption(io);
                            if (io.optionTemplate.id == 72) {
                                this.levelBT = io.param;
                            }
                        }
                    });
        }

        this.player.setClothes.worldcup = 0;
        for (Item item : this.player.inventory.itemsBody) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case 966:
                    case 982:
                    case 983:
                    case 883:
                    case 904:
                        player.setClothes.worldcup++;
                }
                if (item.template.id >= 592 && item.template.id <= 594) {
                    teleport = true;
                }
                for (ItemOption io : item.itemOptions) {
                    addOption(io);
                }
            }
        }
        setDameTrainArmor();
        setBasePoint();
        setOutfitFusion();
        setSpeed();
    }

    private void addOption(ItemOption io) {
        switch (io.optionTemplate.id) {
            case 0: // Tấn công +#
                this.dameAdd += io.param;
                break;
            case 2: // HP, KI+#000
                this.hpAdd += io.param * 1000;
                this.mpAdd += io.param * 1000;
                break;
            case 3:// vô hiệu chưởng
                this.voHieuChuong += io.param;
                break;
            case 5: // +#% sức đánh chí mạng
                this.tlDameCrit.add(io.param);
                this.tlSDCM += io.param;
                break;
            case 6: // HP+#
                this.hpAdd += io.param;
                break;
            case 7: // KI+#
                this.mpAdd += io.param;
                break;
            case 8: // Hút #% HP, KI xung quanh mỗi 5 giây
                this.tlHutHpMpXQ += io.param;
                break;
            case 10: // giảm dame
            case 169:
                this.tlGiamDameXQ += io.param;
                break;
            case 111:
                this.tlPhanTamXQ += io.param > 1 ? io.param : 40;
                break;
            case 14: // Chí mạng+#%
                this.critAdd += io.param;
                break;
            case 16: // Speed
            case 114:
            case 148:
                this.tlSpeed += io.param;
                break;
            case 18: // Chinh xac
                this.tlchinhxac += io.param;
                break;
            case 19: // Tấn công+#% khi đánh quái
                this.tlDameAttMob.add(io.param);
                break;
            case 22: // HP+#K
                this.hpAdd += io.param * 1000;
                break;
            case 23: // MP+#K
                this.mpAdd += io.param * 1000;
                break;
            case 24: // Làm chậm
                this.isLamCham = true;
                break;
            case 25: // Tàn hình
                this.isTanHinh = true;
                break;
            case 26: // Hóa đá
                this.isHoaDa = true;
                break;
            case 27: // +# HP/30s
                this.hpHoiAdd += io.param;
                break;
            case 28: // +# KI/30s
                this.mpHoiAdd += io.param;
                break;
            case 33: // dịch chuyển tức thời
                this.teleport = true;
                break;
            case 34:
                this.setTinhAn += 1;
                break;
            case 35:
                this.setNguyetAn += 1;
                break;
            case 36:
                this.setNhatAn += 1;
                break;
            case 47: // Giáp+#
                this.defAdd += io.param;
                break;
            case 48: // HP/KI+#
                this.hpAdd += io.param;
                this.mpAdd += io.param;
                break;
            case 49: // Tấn công+#%
            case 50: // Sức đánh+#%
                this.tlDame.add(io.param);
                break;
            case 53: // Không ảnh hưởng bởi lão hóa
                this.isKhongLaoHoa = true;
                break;
            case 77: // HP+#%
                this.tlHp.add(io.param);
                break;
            case 80: // HP+#%/30s
                this.tlHpHoi += io.param;
                break;
            case 81: // MP+#%/30s
                this.tlMpHoi += io.param;
                break;
            case 88: // Cộng #% exp khi đánh quái
                this.tlTNSM.add(io.param);
                break;
            case 94: // Giáp #%
                this.tlGiap += io.param;
                break;
            case 95: // Biến #% tấn công thành HP
                this.tlHutHp += io.param;
                break;
            case 96: // Biến #% tấn công thành MP
                this.tlHutMp += io.param;
                break;
            case 97: // Phản #% sát thương
                this.tlPST += io.param;
                break;
            case 98: // Xuyen giap chuong
                this.tlxgc += io.param;
                break;
            case 99: // Xuyen giap can chien
                this.tlxgcc += io.param;
                break;
            case 100: // +#% vàng từ quái
                this.tlGold += io.param;
                break;
            case 101: // +#% TN,SM
                this.tlTNSM.add(io.param);
                break;
            case 103: // KI +#%
                this.tlMp.add(io.param);
                break;
            case 104: // Biến #% tấn công quái thành HP
                this.tlHutHpMob += io.param;
                break;
            case 105: // Vô hình khi không đánh quái và boss
                this.wearingVoHinh = true;
                break;
            case 106: // Không ảnh hưởng bởi cái lạnh
                this.isKhongLanh = true;
                break;
            case 108: // #% Né đòn
                this.tlNeDon += io.param;
                break;
            case 109: // Hôi, giảm #% HP
                this.tlHpGiamODo += io.param;
                break;
            case 110: // Do spl
                this.isDoSPL = true;
                break;
            case 116: // Kháng thái dương hạ san
                this.khangTDHS = true;
                break;

            case 117: // Đẹp +#% SĐ cho mình và người xung quanh
                if (io.param > this.tlSexyDame) {
                    this.tlSexyDame = io.param;
                }
                break;
            case 147: // +#% sức đánh
                this.tlDame.add(io.param);
                break;
            case 156: // Giảm 50% sức đánh, HP, KI và +#% SM, TN, vàng từ quái
                this.tlSubSDCurrent += 50;
                this.tlTNSM.add(io.param);
                this.tlGold += io.param;
                break;
            case 162: // Cute hồi #% KI/s bản thân và xung quanh
                this.mpHoiCute += io.param;
                break;
            case 159: // x chưởng
                this.xChuong = (short) io.param;
                break;
            case 160: // TNSM PET;
                this.tlTNSMPet += io.param;
                break;
            case 173: // Phục hồi #% HP và KI cho đồng đội
                this.tlHpHoiBanThanVaDongDoi += io.param;
                this.tlMpHoiBanThanVaDongDoi += io.param;
                break;
            case 176: //
                setInfoOption176();
                break;
            case 211:
                this.setltdb += 1;
                break;
            case 153: // % phát nổ sau khi chết
                this.tlBom += io.param;
                break;

        }
    }

    private void setSpeed() {
        if (player.isPl()) {
            speed = (byte) (5 + 3 * (tlSpeed / 100));
        }
    }

    private void setInfoOption176() {
        if (player.isPl()) {
            this.tlDame.add(10);
            speed = (byte) (5 + 3 * (50 / 100));
        }
    }

    private void setOutfitFusion() {
        if (this.player.inventory.itemsBody.size() < 6 || this.player.pet == null
                || this.player.pet.inventory.itemsBody.size() < 6) {
            return;
        }
        Item skin = this.player.inventory.itemsBody.get(5);
        Item pskin = this.player.pet.inventory.itemsBody.get(5);
        if (skin.isNotNullItem() && pskin.isNotNullItem()) {
            this.isGogeta = skin.template.id == 1693 && pskin.template.id == 1553
                    || skin.template.id == 1553 && pskin.template.id == 1693;
        } else {
            this.isGogeta = false;
        }
    }

    private void setDameTrainArmor() {
        if (!this.player.isPet && !this.player.isBoss) {
            if (this.player.inventory.itemsBody.size() < 7) {
                return;
            }
            try {
                Item gtl = this.player.inventory.itemsBody.get(6);
                if (gtl.isNotNullItem()) {
                    this.wearingTrainArmor = true;
                    this.player.inventory.trainArmor = gtl;
                    this.tlSubSD += ItemService.gI().getPercentTrainArmor(gtl);
                } else {
                    if (this.player.inventory.trainArmor == null) {
                        gtl = this.player.inventory.itemsBag.stream()
                                .filter(item -> item.isNotNullItem() && item.template.type == 32
                                        && item.itemOptions != null
                                        && item.itemOptions.stream()
                                                .filter(io -> io.optionTemplate.id == 9 && io.param > 0).findFirst()
                                                .orElse(null) != null)
                                .findFirst().orElse(null);
                        if (gtl == null) {
                            return;
                        }
                        this.player.inventory.trainArmor = gtl;
                    }
                    this.wearingTrainArmor = false;
                    for (Item.ItemOption io : this.player.inventory.trainArmor.itemOptions) {
                        if (io.optionTemplate.id == 9 && io.param > 0) {
                            this.tlDame.add(ItemService.gI().getPercentTrainArmor(this.player.inventory.trainArmor));
                            break;
                        }
                    }
                }
            } catch (Exception e) {
                Logger.error("Lỗi get giáp tập luyện " + this.player.name + "\n" + e + "\n");
            }
        }
    }

    public void setBasePoint() {
        setHpMax();
        setHp();
        setMpMax();
        setMp();
        setDame();
        setDef();
        setCrit();
        setHpHoi();
        setMpHoi();
        setLtdb();
        setThoBulma();
        setDietQuy();
        setTiecbaiBien();
        setBunmaTocMau();
        setTinhNhatNguyetAn();
    }

    private void setLtdb() {
        this.islinhthuydanhbac = this.setltdb >= 5;
    }

    private void setThoBulma() {
        this.isThoBulma = (this.player.inventory != null && this.player.inventory.itemsBody != null
                && this.player.inventory.itemsBody.size() >= 5 && this.player.inventory.itemsBody.get(5).isNotNullItem()
                && this.player.inventory.itemsBody.get(5).template.id == 584);
    }

    private void setDietQuy() {
        this.isDietQuy = (this.player.inventory != null && this.player.inventory.itemsBody != null
                && this.player.inventory.itemsBody.size() >= 5 && this.player.inventory.itemsBody.get(5).isNotNullItem()
                && this.player.inventory.itemsBody.get(5).template.id >= 1087
                && this.player.inventory.itemsBody.get(5).template.id <= 1091);
    }

    private void setBunmaTocMau() {
        this.isBunmaTocMau = (this.player.inventory != null && this.player.inventory.itemsBody != null
                && this.player.inventory.itemsBody.size() >= 5 && this.player.inventory.itemsBody.get(5).isNotNullItem()
                && this.player.inventory.itemsBody.get(5).template.id >= 1208
                && this.player.inventory.itemsBody.get(5).template.id <= 1210);
    }

    private void setTiecbaiBien() {
        this.isTiecBaiBien = (this.player.inventory != null && this.player.inventory.itemsBody != null
                && this.player.inventory.itemsBody.size() >= 5 && this.player.inventory.itemsBody.get(5).isNotNullItem()
                && this.player.inventory.itemsBody.get(5).template.id >= 1234
                && this.player.inventory.itemsBody.get(5).template.id <= 1236);
    }

    private void setTinhNhatNguyetAn() {
        this.isTinhAn = this.setTinhAn >= 5;
        this.isNhatAn = this.setNhatAn >= 5;
        this.isNguyetAn = this.setNguyetAn >= 5;
    }

    private void setHpHoi() {
        this.hpHoi = this.hpMax / 100;
        this.hpHoi += this.hpHoiAdd;

        // Kiểm tra giá trị tlHpHoi không vượt quá giới hạn
        if (this.tlHpHoi > 100) {
            this.tlHpHoi = 100;
        } else if (this.tlHpHoi < 0) {
            this.tlHpHoi = 0;
        }

        this.hpHoi += ((long) this.hpMax * this.tlHpHoi / 100);

        // Kiểm tra giá trị tlHpHoiBanThanVaDongDoi không vượt quá giới hạn
        if (this.tlHpHoiBanThanVaDongDoi > 100) {
            this.tlHpHoiBanThanVaDongDoi = 100;
        } else if (this.tlHpHoiBanThanVaDongDoi < 0) {
            this.tlHpHoiBanThanVaDongDoi = 0;
        }

        this.hpHoi += ((long) this.hpMax * this.tlHpHoiBanThanVaDongDoi / 100);
    }

    private void setMpHoi() {
        this.mpHoi = this.mpMax / 100;
        this.mpHoi += this.mpHoiAdd;

        // Kiểm tra giá trị tlMpHoi không vượt quá giới hạn
        if (this.tlMpHoi > 100) {
            this.tlMpHoi = 100;
        } else if (this.tlMpHoi < 0) {
            this.tlMpHoi = 0;
        }

        this.mpHoi += ((long) this.mpMax * this.tlMpHoi / 100);

        // Kiểm tra giá trị tlMpHoiBanThanVaDongDoi không vượt quá giới hạn
        if (this.tlMpHoiBanThanVaDongDoi > 100) {
            this.tlMpHoiBanThanVaDongDoi = 100;
        } else if (this.tlMpHoiBanThanVaDongDoi < 0) {
            this.tlMpHoiBanThanVaDongDoi = 0;
        }

        this.mpHoi += ((long) this.mpMax * this.tlMpHoiBanThanVaDongDoi / 100);
    }

    private void setHpMax() {
        // Tính toán giới hạn hpMax
        long hpMax = Util.maxIntValue(this.hpg + this.hpAdd);

        // Áp dụng các yếu tố ảnh hưởng đến hpMax
        for (Integer tl : this.tlHp) {
            hpMax += (hpMax * tl / 100L);
        }

        // Xử lý set nappa
        if (this.player.setClothes.nappa == 5) {
            hpMax += (hpMax * 80L / 100L);
        }

        if (this.player.setClothes.cadicM >= 2) {
            hpMax += (hpMax * 20L / 100L);
        }

        // Xử lý set worldcup
        if (this.player.setClothes.worldcup == 2) {
            hpMax += (hpMax * 10 / 100L);
        }
        // xử lý nrnm
        if (player.itemTime.isUsePorungaHPKI) {
            hpMax += (hpMax * 25L / 100L);
        }
        // Xử lý rồng xương
        if (player.itemTime != null && player.itemTime.isUseRX) {
            hpMax += (hpMax * 10L / 100L);
        }
        if (player.itemTime != null && player.itemTime.isUseShenronHp) {
            hpMax += (hpMax * 20L / 100L);
        }

        // Xử lý set nhật ấn
        if (this.isNhatAn) {
            hpMax += (hpMax * 15L / 100L);
        }

        // Xử lý ngọc rồng đen 2 sao
        if (this.player.rewardBlackBall.timeOutOfDateReward[1] > System.currentTimeMillis()) {
            hpMax += (hpMax * RewardBlackBall.R2S_1 / 100L);
        }

        // Xử lý khỉ
        if (this.player.effectSkill.isMonkey) {
            if (!this.player.isPet || (this.player.isPet && ((Pet) this.player).status != Pet.FUSION)) {
                int percent = SkillUtil.getPercentHpMonkey(player.effectSkill.levelMonkey);
                hpMax += (hpMax * percent / 100L);
            }
        }
        /* Cộng Chỉ Số Hp Khi Hợp Thể Của Đệ Vip */
        /// pet fide nhí tăng 20%
        if (this.player.isPet && ((Pet) this.player).typePet == 5
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 20 / 100L);
        }
        /// pet cell nhí tăng 25%
        if (this.player.isPet && ((Pet) this.player).typePet == 6
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 40 / 100L);
        }
        /// pet bư nhí tăng 30%
        if (this.player.isPet && ((Pet) this.player).typePet == 7
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 30 / 100L);
        }
        /// pet berus nhí tăng 35%
        if (this.player.isPet && ((Pet) this.player).typePet == 8
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 35 / 100L);
        }
        /// pet Adr Bãi biển tăng 40%
        if (this.player.isPet && ((Pet) this.player).typePet == 9
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 40 / 100L);
        }
        /// pet Mabu gầy tăng 50%
        if (this.player.isPet && ((Pet) this.player).typePet == 10
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 50 / 100L);
        }
        /* End Cộng Chỉ Số Hp Khi Hợp Thể Của Đệ Vip */
        // Xử lý pet pic
        if (this.player.isPet && ((Pet) this.player).typePet == 3
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 15 / 100L);
        }

        // Xử lý pet mabư
        if (this.player.isPet && ((Pet) this.player).typePet == 1
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            Pet pet = (Pet) this.player;
            int bonusPercent = 10 + (pet.levelUpgradeMabu * 10); // base 10% + mỗi cấp +5%
            hpMax += (hpMax * bonusPercent / 100L);

        }

        // Xử lý pet berus
        if (this.player.isPet && ((Pet) this.player).typePet == 2
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 80 / 100L);
        }

        // Xử lý pet black
        if (this.player.isPet && ((Pet) this.player).typePet == 4
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            hpMax += (hpMax * 60 / 100L);
        }

        if (this.player.pet != null && this.player.fusion.typeFusion != ConstPlayer.NON_FUSION) {// Zalo:
            //
            // GiaCatLuong
            if (this.player.pet.typePet >= 5) {
                hpMax += this.player.pet.nPoint.hpMax * this.player.getPointfusion().getHpFusion() / 100L;
            }
        }
        // Xử lý phù
        if (this.player.zone != null && MapService.gI().isMapBlackBallWar(this.player.zone.map.mapId)) {
            hpMax *= this.player.effectSkin.xHPKI;
        }

        // Xử lý thức ăn 2
        if (this.player.itemTime != null && this.player.itemTime.isEatMeal2 && this.player.itemTime.iconMeal2 == 8062) {
            hpMax += (hpMax * 5 / 100L);
        }
        // Xử lý thức ăn 3
        if (this.player.itemTime != null && this.player.itemTime.isEatMeal3 && this.player.itemTime.iconMeal3 == 8244) {
            hpMax += (hpMax * 10 / 100L);
        }
        if (this.player.itemTime != null && this.player.itemTime.isUseLoX15) {
            hpMax -= (hpMax * 85 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX10) {
            hpMax -= (hpMax * 65 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX7) {
            hpMax -= (hpMax * 45 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX5) {
            hpMax -= (hpMax * 35 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX2) {
            hpMax -= (hpMax * 25 / 100L);
        }

        // Xử lý gogeta
        if (this.isGogeta) {
            hpMax += (hpMax * 10 / 100L);
        }

        // Phù map mabu
        if (this.player.isPhuHoMapMabu) {
            hpMax += 1_000_000;
        }

        // Xử lý +hp đệ
        if (this.player.fusion.typeFusion != ConstPlayer.NON_FUSION) {
            hpMax += this.player.pet.nPoint.hpMax;
        }

        // Xử lý bổ huyết
        if (this.player.itemTime != null && this.player.itemTime.isUseBoHuyet && !this.player.itemTime.isUseBoHuyet2) {
            hpMax *= 1.55;
        }

        // Xử lý item sieu cap
        if (this.player.itemTime != null && this.player.itemTime.isUseBoHuyet2) {
            hpMax *= 2.2;
        }

        // Xử lý huýt sáo
        if (!this.player.isPet || (this.player.isPet && ((Pet) this.player).status != Pet.FUSION)) {
            if (this.player.effectSkill.tiLeHPHuytSao != 0) {
                hpMax += (hpMax * this.player.effectSkill.tiLeHPHuytSao / 100L);
            }
        }

        // Xử lý chibi
        if (this.player.effectSkill != null && this.player.effectSkill.isChibi && this.player.typeChibi == 3) {
            hpMax *= 2;
        }

        // Xử lý map lạnh
        if (this.player.zone != null && MapService.gI().isMapCold(this.player.zone.map) && !this.isKhongLanh) {
            hpMax /= 2;
        }
        // xử lí map tương lai
        if (this.player.zone != null
                && !this.player.isBoss
                && MapService.gI().isMapTuongLai(this.player.zone.map)
                && !this.isKhongLaoHoa
                && !(this.player.itemTime != null
                        && (this.player.itemTime.isBanhtet || this.player.itemTime.isBanhchung))) {
            hpMax = (hpMax * 2) / 3;
        }

        if (!this.player.isBoss && !this.player.isNewPet
                && TimeUtil.checkTime(EventDAO.getRemainingTimeToIncreaseHP())) {
            hpMax += hpMax / 10;
        }

        // if (hpMax > 2_000_000_000) {
        // hpMax = 2_000_000_000;
        // }
        this.hpMax = hpMax;
    }

    private void setHp() {
        // Giới hạn giá trị hp không vượt quá hpMax
        if (this.hp > this.hpMax) {
            this.hp = this.hpMax;
        }
    }

    private void setMpMax() {
        // Tính toán giới hạn mpMax
        long mpMax = Util.maxIntValue(this.mpg + this.mpAdd);

        // Áp dụng các yếu tố ảnh hưởng đến mpMax
        for (Integer tl : this.tlMp) {
            mpMax += (mpMax * tl / 100L);
        }

        // Xử lý set picolo
        if (this.player.setClothes.ocTieu == 5) {
            mpMax *= 2;
        }

        // Xử lý set nguyệt ấn
        if (this.isNguyetAn) {
            mpMax += (mpMax * 15L / 100L);
        }

        // Xử lý ngọc rồng đen 6 sao
        if (this.player.rewardBlackBall.timeOutOfDateReward[5] > System.currentTimeMillis()) {
            mpMax += (mpMax * RewardBlackBall.R6S_1 / 100L);
        }

        // Xử lý set worldcup
        if (this.player.setClothes.worldcup == 2) {
            mpMax += (mpMax * 10 / 100L);
        }
        /* Cộng Chỉ Số Ki Khi Hợp Thể Của Đệ Vip */
        /// pet fide nhí tăng 20%
        if (this.player.isPet && ((Pet) this.player).typePet == 5
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 20 / 100L);
        }
        /// pet cell nhí tăng 25%
        if (this.player.isPet && ((Pet) this.player).typePet == 6
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 40 / 100L);
        }
        /// pet bư nhí tăng 30%
        if (this.player.isPet && ((Pet) this.player).typePet == 7
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 30 / 100L);
        }
        /// pet berus nhí tăng 35%
        if (this.player.isPet && ((Pet) this.player).typePet == 8
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 35 / 100L);
        }
        /// pet Adr Bãi biển tăng 40%
        if (this.player.isPet && ((Pet) this.player).typePet == 9
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 40 / 100L);
        }
        /// pet Mabu gầy tăng 50%
        if (this.player.isPet && ((Pet) this.player).typePet == 10
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 50 / 100L);
        }
        /* End Cộng Chỉ Số Ki Khi Hợp Thể Của Đệ Vip */
        // Xử lý pet pic
        if (this.player.isPet && ((Pet) this.player).typePet == 3
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 15 / 100L);
        }

        // Xử lý pet mabư
        if (this.player.isPet && ((Pet) this.player).typePet == 1
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            Pet pet = (Pet) this.player;
            int bonusPercent = 10 + (pet.levelUpgradeMabu * 10); // base 10% + mỗi cấp +5%
            mpMax += (mpMax * bonusPercent / 100L);

        }

        // Xử lý pet br
        if (this.player.isPet && ((Pet) this.player).typePet == 2
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            mpMax += (mpMax * 80 / 100L);// MP berus
        }

        // Xử lý pet black
 if (this.player.isPet && ((Pet) this.player).typePet == 4
        && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
    mpMax += (mpMax * 60 / 100L);
}
        if (this.player.pet != null && this.player.fusion.typeFusion != ConstPlayer.NON_FUSION) {// Zalo:
            //
            // GiaCatLuong
            if (this.player.pet.typePet >= 5) {
                mpMax += this.player.pet.nPoint.mpMax * this.player.getPointfusion().getMpFusion() / 100L;
            }
        }
        // Xử lý phù
        if (this.player.zone != null && MapService.gI().isMapBlackBallWar(this.player.zone.map.mapId)) {
            mpMax *= this.player.effectSkin.xHPKI;
        }

        // Xử lý gogeta
        if (this.isGogeta) {
            mpMax += (mpMax * 10 / 100L);
        }

        // Phù map mabu
        if (this.player.isPhuHoMapMabu) {
            mpMax += 1_000_000;
        }
        if (this.player.itemTime != null && this.player.itemTime.isUseLoX15) {
            mpMax -= (mpMax * 85 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX10) {
            mpMax -= (mpMax * 65 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX7) {
            mpMax -= (mpMax * 45 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX5) {
            mpMax -= (mpMax * 35 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX2) {
            mpMax -= (mpMax * 25 / 100L);
        }
        // Xử lý rồng xương
        if (player.itemTime != null && player.itemTime.isUseRX) {
            mpMax += (mpMax * 10L / 100L);
        }
        if (player.itemTime != null && player.itemTime.isUseShenronKi) {
            mpMax += (mpMax * 20L / 100L);
        }
        // xử lý nrnm
        if (player.itemTime.isUsePorungaHPKI) {
            mpMax += (mpMax * 25L / 100L);
        }

        // Xử lý hợp thể
        if (this.player.fusion.typeFusion != 0) {
            mpMax += this.player.pet.nPoint.mpMax;
        }

        // Xử lý bổ khí
        if (this.player.itemTime != null && this.player.itemTime.isUseBoKhi && !this.player.itemTime.isUseBoKhi2) {
            mpMax *= 1.55;
        }
        // Xử lý map lạnh
        if (this.player.zone != null && MapService.gI().isMapCold(this.player.zone.map) && !this.isKhongLanh) {
            mpMax /= 2;
        }
        // xử lí map tương lai
        if (this.player.zone != null
                && !this.player.isBoss
                && MapService.gI().isMapTuongLai(this.player.zone.map)
                && !this.isKhongLaoHoa
                && !(this.player.itemTime != null
                        && (this.player.itemTime.isBanhtet || this.player.itemTime.isBanhchung))) {
            mpMax = (mpMax * 2) / 3;
        }
        // Xử lý item sieu cap
        if (this.player.itemTime != null && this.player.itemTime.isUseBoKhi2) {
            mpMax *= 2.2;
        }

        if (!this.player.isBoss && !this.player.isNewPet
                && TimeUtil.checkTime(EventDAO.getRemainingTimeToIncreaseMP())) {
            mpMax += mpMax / 10;
        }

        // if (mpMax > 2_000_000_000) {
        // mpMax = 2_000_000_000;
        // }
        this.mpMax = mpMax;
    }

    private void setMp() {
        if (this.mp > this.mpMax) {
            this.mp = this.mpMax;
        }
    }

    public long getHP() {
        return this.hp <= this.hpMax ? this.hp : this.hpMax;
    }

    public void setHP(long hp) {
        if (hp > 0) {
            this.hp = (hp <= this.hpMax ? hp : this.hpMax);
        } else {
            player.setDie();
        }
    }

    public long getMP() {
        return this.mp <= this.mpMax ? this.mp : this.mpMax;
    }

    public void setMP(long mp) {
        if (mp > 0) {
            this.mp = (mp <= this.mpMax ? mp : this.mpMax);
        } else {
            this.mp = 0;
        }
    }

    private void setDame() {
        // Tính toán giới hạn dame
        long dame = Util.maxIntValue(this.dameg + this.dameAdd);

        // Áp dụng các yếu tố ảnh hưởng đến dame
        for (Integer tl : this.tlDame) {
            dame += (dame * tl / 100L);
        }

        // for (Integer tl : this.tlSDDep) {
        // dame += (dame * tl / 100L);
        // }
        /* Cộng Chỉ Số Dame Khi Hợp Thể Của Đệ Vip */
        /// pet fide nhí tăng 20%
        if (this.player.isPet && ((Pet) this.player).typePet == 5
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 20 / 100L);
        }
        /// pet cell nhí tăng 25%
        if (this.player.isPet && ((Pet) this.player).typePet == 6
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 40 / 100L);
        }
        /// pet bư nhí tăng 30%
        if (this.player.isPet && ((Pet) this.player).typePet == 7
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 30 / 100L);
        }
        /// pet berus nhí tăng 35%
        if (this.player.isPet && ((Pet) this.player).typePet == 8
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 35 / 100L);
        }
        /// pet Adr Bãi biển tăng 40%
        if (this.player.isPet && ((Pet) this.player).typePet == 9
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 40 / 100L);
        }

        /// pet Mabu gầy tăng 50%
        if (this.player.isPet && ((Pet) this.player).typePet == 10
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 50 / 100L);
        }
        /* End Cộng Chỉ Số Dame Khi Hợp Thể Của Đệ Vip */
        // Xử lý pet pic
        if (this.player.isPet && ((Pet) this.player).typePet == 3
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 15 / 100L);
        }

        // Xử lý pet mabư
        if (this.player.isPet && ((Pet) this.player).typePet == 1
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            Pet pet = (Pet) this.player;
            int bonusPercent = 10 + (pet.levelUpgradeMabu * 10); // base 10% + mỗi cấp +5%
            dame += (dame * bonusPercent / 100L);

        }

        // Xử lý pet br
        if (this.player.isPet && ((Pet) this.player).typePet == 2
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 80 / 100L);
        }

        // Xử lý pet black
        if (this.player.isPet && ((Pet) this.player).typePet == 4
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            dame += (dame * 60 / 100L);
        }
        if (this.player.pet != null && this.player.fusion.typeFusion != ConstPlayer.NON_FUSION) {// Zalo:
            //
            // GiaCatLuong
            if (this.player.pet.typePet >= 5) {
                dame += this.player.pet.nPoint.dame * this.player.getPointfusion().getDameFusion() / 100L;
            }
        }
        // Xử lý set tinh ấn
        if (this.isTinhAn) {
            dame += (dame * 15L / 100L);
        }

        // Xử lý thức ăn
        if (!this.player.isPet && this.player.itemTime != null && this.player.itemTime.isEatMeal
                || this.player.isPet && this.player.itemTime != null && ((Pet) this.player).master.itemTime.isEatMeal) {
            dame += (dame * 10 / 100L);
        }
        if (this.player.itemTime != null && this.player.itemTime.isUseLoX15) {
            dame -= (dame * 85 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX10) {
            dame -= (dame * 65 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX7) {
            dame -= (dame * 45 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX5) {
            dame -= (dame * 35 / 100L);
        } else if (this.player.itemTime != null && this.player.itemTime.isUseLoX2) {
            dame -= (dame * 25 / 100L);
        }
        // Xử lý thức ăn 2
        if (this.player.itemTime != null && this.player.itemTime.isEatMeal2 && this.player.itemTime.iconMeal2 == 8060) {
            dame += (dame * 5 / 100L);
        }

        // Xử lý thức ăn 2
        if (this.player.itemTime != null && this.player.itemTime.isEatMeal2 && this.player.itemTime.iconMeal2 == 8061) {
            this.tlDameCrit.add(5);
            this.tlSDCM += 5;
        }

        // Xử lý cuồng nộ
        if (this.player.itemTime != null && this.player.itemTime.isUseCuongNo && !this.player.itemTime.isUseCuongNo2) {
            dame *= 1.55;
        }
        if (this.player.itemTime != null && this.player.itemTime.isXimuoihoadao) {
            dame += (dame * 20 / 100L);
        }
        if (this.player.itemTime != null && this.player.itemTime.isBanhtet) {
            dame += (dame * 10 / 100L);
        }
        if (this.player.itemTime != null && this.player.itemTime.isBanhchung) {
            dame += (dame * 30 / 100L);
        }
        if (this.player.itemTime != null && this.player.itemTime.isXimuoihoamai) {
            dame += (dame * 20 / 100L);
        }
        if (this.player.itemTime != null && this.player.itemTime.isUseCuongNo2) {
            dame *= 2.2;
        }
        if (this.player.itemTime != null && this.player.itemTime.isUseKhauTrang) {
            this.tlDameAttMob.add(10);
        }

        if (this.player.itemTime != null && this.player.itemTime.isEatMeal3 && this.player.itemTime.iconMeal3 == 8247) {
            dame += (dame * 10 / 100L);
        }

        // Xử lý ngọc rồng đen 1 sao
        if (this.player.rewardBlackBall.timeOutOfDateReward[0] > System.currentTimeMillis()) {
            dame += (dame * RewardBlackBall.R1S_2 / 100L);
        }

        // Xử lý set worldcup
        if (this.player.setClothes.worldcup == 2) {
            dame += (dame * 10 / 100L);
        }
        // Xử lý set nail
        if (this.player.setClothes.nail >= 2) {
            this.tlDameCrit.add(20);
        }
        // Xử lý gogeta
        if (this.isGogeta) {
            dame += (dame * 10 / 100L);
        }

        // Phù map mabu
        if (this.player.isPhuHoMapMabu) {
            dame += 10_000;
        }

        // Xử lý rồng xương
        if (player.itemTime != null && player.itemTime.isUseRX) {
            dame += (dame * 10L / 100L);
        }
        if (player.itemTime != null && player.itemTime.isUseShenronDame) {
            dame += (dame * 15L / 100L);
        }

        // Xử lý phù
        if (this.player.zone != null && MapService.gI().isMapBlackBallWar(this.player.zone.map.mapId)) {
            dame *= this.player.effectSkin.xDame;
        }

        // Xử lý hợp thể
        if (this.player.fusion.typeFusion != 0) {
            dame += this.player.pet.nPoint.dame;
        }

        // Xử lý nrnm
        if (player.itemTime.isUsePorungaSD) {
            dame += (dame * 20L / 100L); // +20% SD
        }

        // Xử lý khỉ
        if (this.player.effectSkill.isMonkey) {
            if (!this.player.isPet || (this.player.isPet && ((Pet) this.player).status != Pet.FUSION)) {
                int percent = SkillUtil.getPercentDameMonkey(player.effectSkill.levelMonkey);
                dame += (dame * percent / 100L);
            }
        }

        // Sức đánh đẹp
        dame += (dame * tlSexyDame / 100L);

        // Xử lý giảm dame
        dame -= (dame * tlSubSD / 100L);

        // Xử lý map cold
        if (this.player.zone != null && MapService.gI().isMapCold(this.player.zone.map) && !this.isKhongLanh) {
            dame /= 2;
        }
        // xử lí map tương lai
        if (this.player.zone != null
                && !this.player.isBoss
                && MapService.gI().isMapTuongLai(this.player.zone.map)
                && !this.isKhongLaoHoa
                && !(this.player.itemTime != null
                        && (this.player.itemTime.isBanhtet || this.player.itemTime.isBanhchung))) {
            dame = (dame * 2) / 3;
        }

        if (!this.player.isBoss && !this.player.isNewPet
                && TimeUtil.checkTime(EventDAO.getRemainingTimeToIncreaseDame())) {
            dame += dame / 10;
        }

        dame -= (dame * tlSubSDCurrent / 100L);

        // if (dame > 2_000_000_000) {
        // dame = 2_000_000_000;
        // }
        this.dame = dame;
    }

    private void setDef() {
        this.def = this.defg * 4;
        this.def += this.defAdd;
        // Xử lý thức ăn 3
        if (this.player.itemTime != null && this.player.itemTime.isEatMeal3 && this.player.itemTime.iconMeal3 == 8246) {
            this.def += (this.def * 10 / 100L);
        }
        if (this.player.rewardBlackBall.timeOutOfDateReward[1] > System.currentTimeMillis()) {
            this.def += RewardBlackBall.R3S_1;
        }
        // xử lý giáp xên
        if (this.player.itemTime != null && this.player.itemTime.isUseGiapXen && !this.player.itemTime.isUseGiapXen2) {
            this.def *= 1.55;
        }
        // xử lý giáp xên 2
        if (this.player.itemTime != null && this.player.itemTime.isUseGiapXen2 && !this.player.itemTime.isUseGiapXen) {
            this.def *= 2.2;
        }
        // rồng xương
        if (player.itemTime != null && player.itemTime.isUseRX2) {
            this.def *= 1.1;
        }
        if (player.itemTime != null && player.itemTime.isUseShenronArmor) {
            this.def += 3000;
        }
        // xử lý nrnm
        if (player.itemTime.isUsePorungaGIAP) {
            this.def += 3000; // +3000 giáp
        }
        // xử lí đệ thường và mabu
        if (this.player.fusion.typeFusion != ConstPlayer.NON_FUSION && this.player.pet != null) {
            this.def += this.player.pet.nPoint.def;
        }
        if (this.player.isPet && ((Pet) this.player).typePet == 1
                && (((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA
                        || ((Pet) this.player).master.fusion.typeFusion == ConstPlayer.HOP_THE_PORATA2
                        || ((Pet) this.player).master.fusion.typeFusion == 9)) {
            Pet pet = (Pet) this.player;
            int bonusPercent = 5 + (pet.levelUpgradeMabu * 10);
            def += (def * bonusPercent / 100L);
        }
        // xử lý map lạnh
        if (this.player.zone != null
                && MapService.gI().isMapCold(this.player.zone.map)
                && !this.isKhongLanh) {
            def /= 2;
        }

        // xử lý map tương lai
        if (this.player.zone != null
                && !this.player.isBoss
                && MapService.gI().isMapTuongLai(this.player.zone.map)
                && !this.isKhongLaoHoa
                && !(this.player.itemTime != null
                        && (this.player.itemTime.isBanhtet || this.player.itemTime.isBanhchung))) {
            def = (def * 2) / 3;
        }
    }

    private void setCrit() {
        this.crit = this.critg;
        this.crit += this.critAdd;
        // biến khỉ
        if (this.player.effectSkill.isMonkey) {
            this.crit = 110;
        }
        if (player.setClothes.thanVuTruKaio >= 2) {
            this.crit += 20;
        }
        // Xử lý thức ăn 3
        if (this.player.itemTime != null && this.player.itemTime.isEatMeal3 && this.player.itemTime.iconMeal3 == 8244) {
            this.crit = this.crit + 5;
        }
        if (this.player.itemTime != null && this.player.itemTime.isBanhtet) {
            this.crit = this.crit + 15;
        }
        if (this.player.itemTime != null && this.player.itemTime.isBanhchung) {
            this.crit = this.crit + 15;
        }
    }

    private void resetPoint() {
        this.voHieuChuong = 0;
        this.hpAdd = 0;
        this.mpAdd = 0;
        this.dameAdd = 0;
        this.defAdd = 0;
        this.critAdd = 0;
        this.tlHp.clear();
        this.tlMp.clear();
        this.tlDef.clear();
        this.tlDame.clear();
        this.tlDameCrit.clear();
        this.tlDameAttMob.clear();
        this.tlSDCM = 0;
        this.tlHpHoiBanThanVaDongDoi = 0;
        this.tlMpHoiBanThanVaDongDoi = 0;
        this.hpHoi = 0;
        this.mpHoi = 0;
        this.mpHoiCute = 0;
        this.tlHpHoi = 0;
        this.tlMpHoi = 0;
        this.tlHutHp = 0;
        this.tlHutMp = 0;
        this.tlHutHpMob = 0;
        this.tlHutHpMpXQ = 0;
        this.tlGiamDameXQ = 0;
        this.tlPhanTamXQ = 0;
        this.tlPST = 0;
        this.tlTNSM.clear();
        this.tlDameAttMob.clear();
        this.tlGold = 0;
        this.tlNeDon = 0;
        this.tlBom = 0;
        this.tlGiap = 0;
        this.tlxgcc = 0;
        this.tlxgc = 0;
        this.tlchinhxac = 0;
        this.tlTNSMPet = 0;
        this.xChuong = 0;
        this.setltdb = 0;
        this.setTinhAn = 0;
        this.setNhatAn = 0;
        this.setNguyetAn = 0;
        this.tlSexyDame = 0;
        this.tlSubSD = 0;
        this.tlSubSDCurrent = 0;
        this.tlHpGiamODo = 0;
        this.tlSpeed = 0;
        this.teleport = false;

        this.wearingVoHinh = false;
        this.isKhongLanh = false;
        this.isKhongLaoHoa = false;
        this.khangTDHS = false;
        this.isTanHinh = false;
        this.isHoaDa = false;
        this.isLamCham = false;
        this.isDoSPL = false;
        this.isThoBulma = false;
        this.isDietQuy = false;
        this.isBunmaTocMau = false;
        this.isTiecBaiBien = false;
    }

    public void addHp(long hp) {
        if (hp > 0) {
            long potentialHp = this.hp + hp;
            if (potentialHp > this.hpMax) {
                this.hp = this.hpMax;
            } else {
                this.hp = potentialHp;
            }
        }
    }

    public void addMp(long mp) {
        long potentialMp = this.mp + mp;

        if (potentialMp > this.mpMax) {
            this.mp = this.mpMax;
        } else if (potentialMp < 0) {
            this.mp = 0;
        } else {
            this.mp = potentialMp;
        }
    }

    public void setHp(long hp) {
        if (hp < 0) {
            this.hp = 0;
        } else {
            this.hp = hp;
        }
    }

    public void setMp(long mp) {
        // if (mp > this.mpMax) {
        // this.mp = this.mpMax;
        // } else
        if (mp < 0) {
            this.mp = 0;
        } else {
            this.mp = mp;
        }
    }

    private void setIsCrit() {
        if (intrinsic != null && intrinsic.id == 25 && this.getCurrPercentHP() <= intrinsic.param1) {
            isCrit = true;
        } else if (isCrit100) {
            isCrit100 = false;
            isCrit = true;
        } else {
            isCrit = Util.isTrue(this.crit, ConstRatio.PER100);
        }
    }

    public long getDameAttack(boolean isAttackMob) {
        if (isAttackMob && this.player.zone != null && this.player.zone.map != null
                && MapService.gI().isMapNguHanhSon(this.player.zone.map.mapId)) {
            return 80_000;
        }
        setIsCrit();
        long dameAttack = this.dame;
        intrinsic = this.player.playerIntrinsic.intrinsic;
        percentDameIntrinsic = 0;
        int percentDameSkill = 0;
        byte percentXDame = 0;
        Skill skillSelect = player.playerSkill.skillSelect;
        if (skillSelect.template.id != Skill.DICH_CHUYEN_TUC_THOI && isCritTele) {
            isCrit = true;
            isCritTele = false;
        }
        switch (skillSelect.template.id) {
            case Skill.DRAGON:
                if (intrinsic.id == 1) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                percentDameSkill = skillSelect.damage;
                break;
            case Skill.KAMEJOKO:
                if (intrinsic.id == 2) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                percentDameSkill = skillSelect.damage;
                if (this.player.setClothes.songoku == 5) {
                    percentXDame = 100;
                }
                break;
            case Skill.GALICK:
                if (intrinsic.id == 16) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                percentDameSkill = skillSelect.damage;
                if (this.player.setClothes.kakarot == 5) {
                    percentXDame = 100;
                }
                break;
            case Skill.ANTOMIC:
                if (intrinsic.id == 17) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                percentDameSkill = skillSelect.damage;
                break;
            case Skill.DEMON:
                if (intrinsic.id == 8) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                percentDameSkill = skillSelect.damage;
                break;
            case Skill.MASENKO:
                if (intrinsic.id == 9) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                if (this.player.setClothes.nail == 5) {
                    percentXDame = 50;
                }
                percentDameSkill = skillSelect.damage;
                break;
            case Skill.LIEN_HOAN:
                if (intrinsic.id == 13) {
                    percentDameIntrinsic = intrinsic.param1;
                }
                percentDameSkill = skillSelect.damage;
                if (this.player.setClothes.lienHoan == 5) {
                    percentXDame = 101;
                }
                break;
            case Skill.KAIOKEN:
                percentDameSkill = skillSelect.damage;
                if (this.player.setClothes.kaioken == 5) {
                    percentXDame = 50;
                } else if (player.setClothes.thanVuTruKaio == 5) {
                    percentXDame = 30;
                }
                break;
            case Skill.TU_SAT:
                percentDameSkill = skillSelect.damage;
                if (this.player.setClothes.cadicM == 4) {
                    percentXDame = 20;
                } else if (this.player.setClothes.cadicM == 5) {
                    percentXDame = 40;
                }
                break;
            case Skill.DICH_CHUYEN_TUC_THOI:
                isCrit = true;
                isCritTele = true;
                dameAttack = Util.nextLong(Util.maxIntValue((dameAttack - (dameAttack / 100 * 5))),
                        Util.maxIntValue((dameAttack + (dameAttack / 100 * 5))));
                break;
            case Skill.MAKANKOSAPPO:
                percentDameSkill = skillSelect.damage;
                long dameSkill = Util.maxIntValue((long) this.mpMax * percentDameSkill / 100);
                if (this.player.setClothes.picolo == 5) {
                    dameSkill += dameSkill * 50 / 100;
                }
                return dameSkill;
            // case Skill.QUA_CAU_KENH_KHI:
            // long hpmob = 0;
            // long hppl = 0;

            // for (Mob mob : this.player.zone.mobs) {
            // if (!mob.isDie() && Util.getDistance(this.player, mob) <= SkillUtil
            // .getRangeQCKK(this.player.playerSkill.skillSelect.point)) {
            // hpmob += mob.point.hp;
            // }
            // }

            // for (Player pl : this.player.zone.getHumanoids()) {
            // if (!pl.isDie() && this.player.id != pl.id && Util.getDistance(this.player,
            // pl) <= SkillUtil
            // .getRangeQCKK(this.player.playerSkill.skillSelect.point)) {
            // hppl += pl.nPoint.hp;
            // }
            // }

            // long dameqckk = (hpmob * 10 / 100) + (hppl * 10 / 100) + this.dame * 10;

            // if (this.player.setClothes.kirin == 5) {
            // dameqckk *= 2;
            // }

            // dameqckk = dameqckk + (Util.nextInt(-5, 5) * dameqckk / 100);
            // // if (dameqckk > 2_000_000_000) {
            // // dameqckk = 2_000_000_000;
            // // }
            // return dameqckk;
            case Skill.QUA_CAU_KENH_KHI:
                long dameqckk = (this.def * 150L);
                if (this.player.setClothes.kirin == 5) {
                    dameqckk *= 2;
                }
                dameqckk = dameqckk + (Util.nextInt(-5, 5) * dameqckk / 100);
                return dameqckk;

            case Skill.DE_TRUNG:
                if (player.setClothes.pikkoroDaimao == 5) {
                    dameAttack *= 3;
                }
                // if (dameAttack > 2_000_000_000) {
                // dameAttack = 2_000_000_000;
                // }
                return dameAttack;
        }

        if (intrinsic.id == 18 && this.player.effectSkill.isMonkey) {
            percentDameIntrinsic = intrinsic.param1;
        }

        if (percentDameSkill != 0) {
            dameAttack = dameAttack * percentDameSkill / 100;
        }

        dameAttack += (dameAttack * percentDameIntrinsic / 100);
        dameAttack += (dameAttack * dameAfter / 100);
        if (this.player.effectSkill != null && this.player.effectSkill.isDameBuff
                && (this.player.effectSkill.tileDameBuff < 0 || tlSexyDame == 0)) {
            int tiLeDame = this.player.effectSkill.tileDameBuff;
            dameAttack += (dameAttack * tiLeDame / 100L);
        }
        if (isAttackMob) {
            for (Integer tl : this.tlDameAttMob) {
                dameAttack += (dameAttack * tl / 100);
            }
            if (this.player.isPet && ((Pet) this.player).master.charms.tdDeTu > System.currentTimeMillis()) {
                dameAttack *= 2;
            }
        }

        dameAfter = 0;

        if (isCrit) {
            dameAttack *= 2;
            dameAttack += (dameAttack * tlSDCM / 100);
        }

        dameAttack += ((long) dameAttack * percentXDame / 100);

        long tempDameAttack = (long) (dameAttack / 100L * 5L);
        if (tempDameAttack <= 0) {
            tempDameAttack = 1;
        }
        dameAttack += (long) (Util.getOne(-1, 1) * Util.nextInt((int) tempDameAttack) + 1);

        if (player.effectSkin != null && player.effectSkin.isXChuong
                && (player.playerSkill.skillSelect.template.id == Skill.KAMEJOKO
                        || player.playerSkill.skillSelect.template.id == Skill.ANTOMIC
                        || player.playerSkill.skillSelect.template.id == Skill.MASENKO)) {
            dameAttack *= xChuong;
            player.effectSkin.isXDame = true;
            player.effectSkin.isXChuong = false;
            player.effectSkin.lastTimeXChuong = System.currentTimeMillis();
        }

        // if (dameAttack > 2_000_000_000) {
        // dameAttack = 2_000_000_000;
        // }
        return dameAttack;
    }

    public int getCurrPercentHP() {
        if (this.hpMax == 0) {
            return 100;
        }
        return (int) ((long) this.hp * 100 / this.hpMax);
    }

    public int getCurrPercentMP() {
        return (int) ((long) this.mp * 100 / this.mpMax);
    }

    public void setFullHpMp() {
        this.hp = this.hpMax;
        this.mp = this.mpMax;
    }

    public void subHP(long sub) {
        this.hp -= sub;
        if (this.hp <= 0) {
            this.hp = 0;
            this.setHp(0);
        }
    }

    public void subMP(long sub) {
        this.mp -= sub;
        if (this.mp <= 0) {
            this.mp = 0;
        }

        // if (this.mp > 2_000_000_000) {
        // this.mp = 2_000_000_000;
        // }
    }

    public long calSucManhTiemNang(long tiemNang) {
        if (player.zone.map.type == 3) {
            return 0;
        }
        if (power < getPowerLimit()) {
            for (Integer tl : this.tlTNSM) {
                tiemNang += ((long) tiemNang * tl / 100);
            }
            if (this.player.cFlag != 0) {
                if (this.player.cFlag == 8) {
                    tiemNang += ((long) tiemNang * 10 / 100);
                } else {
                    tiemNang += ((long) tiemNang * 5 / 100);
                }
            }
            long tn = tiemNang;
            if (this.player.charms.tdTriTue > System.currentTimeMillis()) {
                tiemNang += tn;
            }
            if (this.player.charms.tdTriTue3 > System.currentTimeMillis()) {
                tiemNang += tn * 2;
            }
            if (this.player.charms.tdTriTue4 > System.currentTimeMillis()) {
                tiemNang += tn * 3;
            }
            if (this.player.charms.tdTriTue4 > System.currentTimeMillis()) {
                tiemNang += tn * 3;
            }
            if (this.player.effectSkill.isChibi && this.player.typeChibi == 2) {
                tiemNang += tn * 2;
            }
            if (this.player.getSession() != null && this.player.getSession().vip > 0
                    || this.player.isPet && ((Pet) this.player).master.getSession() != null
                            && ((Pet) this.player).master.getSession().vip > 0) {
                tiemNang += tn * 3;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseDK) {
                tiemNang += tn * 1.5;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseKhauTrang) {
                tiemNang += tn * 5 / 100;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseLoX2) {
                tiemNang += tn * 0.5;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseLoX5) {
                tiemNang += tn * 5;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseLoX7) {
                tiemNang += tn * 7;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseLoX10) {
                tiemNang += tn * 10;
            }
            if (this.player.itemTime != null && this.player.itemTime.isUseLoX15) {
                tiemNang += tn * 15;
            }
            // Xử lý rồng xương
            if (player.itemTime != null && player.itemTime.isUseRX2) {
                tiemNang += tn * 1.15;
            }
            if (player.itemTime != null && player.itemTime.isUseShenronTnsm) {
                tiemNang += tn * 30 / 100;
            }
            if (this.player.satellite != null && this.player.satellite.isIntelligent) {
                tiemNang += tn / 5;
            }
            if (this.intrinsic != null && this.intrinsic.id == 24) {
                tiemNang += ((long) tiemNang * this.intrinsic.param1 / 100);
            }
            if (this.power >= 70000000000L) {
                tiemNang -= ((long) tiemNang * 80 / 100);
            }
            // if (this.player.isPet) {
            // if (((Pet) this.player).master.charms.tdDeTu > System.currentTimeMillis()) {
            // tiemNang += tn * 2;
            // }
            // if (((Pet) this.player).itemTime.lastTimeBuax2DeTu >
            // System.currentTimeMillis()) {
            // tiemNang += tn * 2;
            // }
          if (this.player instanceof Pet) {
    Pet pet = (Pet) this.player;

    if (pet.master != null
            && pet.master.itemTime != null
            && pet.master.itemTime.isUseBuax2DeTu) {
        tiemNang += tn;
    }

    if (pet.master != null
            && pet.master.nPoint != null
            && pet.master.nPoint.tlTNSMPet > 0) {

        tiemNang += (tn * pet.master.nPoint.tlTNSMPet) / 100;
    }
}
            

            if (TimeUtil.checkTime(EventDAO.getRemainingTimeToIncreasePotentialAndPower())) {
                tiemNang *= 2;
            }
            if (MapService.gI().isMapNguHanhSon(this.player.zone.map.mapId)) {
                tiemNang *= 4;
            }
            // if (MapService.gI().isMapBanDoKhoBau(this.player.zone.map.mapId)) {
            //     tiemNang *= 0.05;
            // }if (this.player.zone.map.mapId == 135) {
        
           tiemNang *= Manager.RATE_EXP_SERVER;
tiemNang = calSubTNSM(tiemNang);

// 👉 R7S đặt SAU CÙNG (chuẩn)
if (this.player.rewardBlackBall != null
        && this.player.rewardBlackBall.timeOutOfDateReward[6] > System.currentTimeMillis()) {
    tiemNang += (tiemNang * RewardBlackBall.R7S_1 / 100L);
}

// chống âm
if (tiemNang <= 0) {
    tiemNang = 1;
}
        } else {
            tiemNang = 0;
        }
        return tiemNang;
    }

    public long calSubTNSM(long tiemNang) {
        long power = player.nPoint.power;

        // giảm cơ bản 50%
        tiemNang = (long) (tiemNang * 0.5);

        // dãy các ngưỡng sức mạnh
        long[] thresholds = {
                 5_000_000_000L,
                15_100_000_000L,
                25_100_000_000L,
                35_100_000_000L,
                45_100_000_000L,              
               
              
              
                55_000_000_000L,
             
              
                65_000_000_000L,
               
                75_000_000_000L,
                80_000_000_000L,
                85_000_000_000L,
                90_000_000_000L,
                95_000_000_000L,
                100_000_000_000L
        };

        // mỗi ngưỡng vượt qua → giảm thêm 50%
        for (long threshold : thresholds) {
            if (power >= threshold) {
                tiemNang = tiemNang / 2;
            }
        }

        return tiemNang;
    }

    public short getTileHutHp(boolean isMob) {
        if (isMob) {
            return (short) (this.tlHutHp + this.tlHutHpMob);
        } else {
            return this.tlHutHp;
        }
    }

    public short getTiLeHutMp() {
        return this.tlHutMp;
    }

    public int subDameInjureWithDeff(long dame) {
        long def = this.def;
        dame -= def;
        if (dame < 0) {
            dame = 1;
        }
        return (int) dame;
    }

    /*------------------------------------------------------------------------*/
    public boolean canOpenPower() {
        return this.power >= getPowerLimit();
    }

    public long getPowerLimit() {
        if (powerLimit != null) {
            return powerLimit.getPower();
        }
        return 0;
    }

    public long getPowerNextLimit() {
        PowerLimit powerLimit = PowerLimitManager.getInstance().get(limitPower + 1);
        if (powerLimit != null) {
            return powerLimit.getPower();
        }
        return 0;
    }

    // **************************************************************************
    // POWER - TIEM NANG
    public void powerUp(long power) {
        this.power += power;
        TaskService.gI().checkDoneTaskPower(player, this.power);
    }

    public void tiemNangUp(long tiemNang) {
        this.tiemNang += tiemNang;
    }

    public void increasePoint(byte type, short point, boolean manualForPet) {
        if (powerLimit == null) {
            return;
        }
        if (point <= 0) {
            return;
        }
        boolean updatePoint = false;
        long tiemNangUse = 0;
        if (type == 0) {
            int pointHp = point * 20;
            tiemNangUse = point * (2 * (this.hpg + 1000) + pointHp - 20) / 2;
            if ((this.hpg + pointHp) <= powerLimit.getHp()) {
                if (doUseTiemNang(tiemNangUse)) {
                    hpg += pointHp;
                    updatePoint = true;
                }
            } else {
                Service.gI().sendThongBao(player, "HP của bạn đã đạt mức tối đa");
                Service.gI().sendMoney(player);
                return;
            }
        }
        if (type == 1) {
            int pointMp = point * 20;
            tiemNangUse = point * (2 * (this.mpg + 1000) + pointMp - 20) / 2;
            if ((this.mpg + pointMp) <= powerLimit.getMp()) {
                if (doUseTiemNang(tiemNangUse)) {
                    mpg += pointMp;
                    updatePoint = true;
                }
            } else {
                Service.gI().sendThongBao(player, "KI của bạn đã đạt mức tối đa");
                Service.gI().sendMoney(player);
                return;
            }
        }
        if (type == 2) {
            tiemNangUse = point * (2 * this.dameg + point - 1) / 2 * 100;
            if ((this.dameg + point) <= powerLimit.getDamage()) {
                if (doUseTiemNang(tiemNangUse)) {
                    dameg += point;
                    updatePoint = true;
                }
            } else {
                Service.gI().sendThongBao(player, "Sức đánh của bạn đã đạt mức tối đa");
                Service.gI().sendMoney(player);
                return;
            }
        }
        if (type == 3) {
            tiemNangUse = 2 * (this.defg + 5) / 2 * 100000;
            if ((this.defg + point) <= powerLimit.getDefense()) {
                if (doUseTiemNang(tiemNangUse)) {
                    defg += point;
                    updatePoint = true;
                }
            } else {
                Service.gI().sendThongBao(player, "Giáp của bạn đã đạt mức tối đa");
                Service.gI().sendMoney(player);
                return;
            }
        }
        if (type == 4) {
            tiemNangUse = 50000000L;
            for (int i = 0; i < this.critg; i++) {
                tiemNangUse *= 5L;
            }
            if ((this.critg + point) <= powerLimit.getCritical()) {
                if (doUseTiemNang(tiemNangUse)) {
                    critg += point;
                    updatePoint = true;
                }
            } else {
                Service.gI().sendThongBao(player, "Chí mạng của bạn đã đạt mức tối đa");
                Service.gI().sendMoney(player);
                return;
            }
        }
        if (updatePoint) {
            Service.gI().point(player);
        }
        if (manualForPet) {
            Service.gI().sendChiSoPetGoc(((Pet) player).master);
            Service.gI().showInfoPet(((Pet) player).master);
            Service.gI().point(((Pet) player).master);
        }
    }
    // public void increasePoint(byte type, short point, boolean manualForPet) {
    // if (powerLimit == null) {
    // return;
    // }
    // if (point <= 0) {
    // return;
    // }
    // boolean updatePoint = false;
    // long tiemNangUse = 0;
    // if (type == 0) {
    // int pointHp = point * 20;
    // tiemNangUse = point * (2 * (this.hpg + 1000) + pointHp - 20) / 2;
    // if ((this.hpg + pointHp) <= powerLimit.getHp()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // hpg += pointHp;
    // updatePoint = true;
    // }
    // } else {
    // Service.gI().sendThongBao(player, "HP của bạn đã đạt mức tối đa");
    // Service.gI().sendMoney(player);
    // return;
    // }
    // }
    // if (type == 1) {
    // int pointMp = point * 20;
    // tiemNangUse = point * (2 * (this.mpg + 1000) + pointMp - 20) / 2;
    // if ((this.mpg + pointMp) <= powerLimit.getMp()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // mpg += pointMp;
    // updatePoint = true;
    // }
    // } else {
    // Service.gI().sendThongBao(player, "KI của bạn đã đạt mức tối đa");
    // Service.gI().sendMoney(player);
    // return;
    // }
    // }
    // if (type == 2) {
    // tiemNangUse = point * (2 * this.dameg + point - 1) / 2 * 100;
    // if ((this.dameg + point) <= powerLimit.getDamage()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // dameg += point;
    // updatePoint = true;
    // }
    // } else {
    // Service.gI().sendThongBao(player, "Sức đánh của bạn đã đạt mức tối đa");
    // Service.gI().sendMoney(player);
    // return;
    // }
    // }
    // if (type == 3) {
    // tiemNangUse = 2 * (this.defg + 5) / 2 * 100000;
    // if ((this.defg + point) <= powerLimit.getDefense()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // defg += point;
    // updatePoint = true;
    // }
    // } else {
    // Service.gI().sendThongBao(player, "Giáp của bạn đã đạt mức tối đa");
    // Service.gI().sendMoney(player);
    // return;
    // }
    // }
    // if (type == 4) {
    // tiemNangUse = 50000000L;
    // for (int i = 0; i < this.critg; i++) {
    // tiemNangUse *= 5L;
    // }
    // if ((this.critg + point) <= powerLimit.getCritical()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // critg += point;
    // updatePoint = true;
    // }
    // } else {
    // Service.gI().sendThongBao(player, "Chí mạng của bạn đã đạt mức tối đa");
    // Service.gI().sendMoney(player);
    // return;
    // }
    // }
    // if (updatePoint) {
    // Service.gI().point(player);
    // }
    // if (manualForPet) {
    // Service.gI().sendChiSoPetGoc(((Pet) player).master);
    // Service.gI().showInfoPet(((Pet) player).master);
    // Service.gI().point(((Pet) player).master);
    // }
    // }

    // public void increasePoint(byte type, short point) {
    // if (point <= 0 || point > 100) {
    // return;
    // }
    // long tiemNangUse;
    // if (type == 0) {
    // int pointHp = point * 20;
    // tiemNangUse = point * (2 * (this.hpg + 1000) + pointHp - 20) / 2;
    // if ((this.hpg + pointHp) <= getHpMpLimit()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // hpg += pointHp;
    // }
    // } else {
    // Service.gI().sendThongBaoOK(player, "Vui lòng mở giới hạn sức mạnh");
    // return;
    // }
    // }
    // if (type == 1) {
    // int pointMp = point * 20;
    // tiemNangUse = point * (2 * (this.mpg + 1000) + pointMp - 20) / 2;
    // if ((this.mpg + pointMp) <= getHpMpLimit()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // mpg += pointMp;
    // }
    // } else {
    // Service.gI().sendThongBaoOK(player, "Vui lòng mở giới hạn sức mạnh");
    // return;
    // }
    // }
    // if (type == 2) {
    // TaskService.gI().checkDoneTaskNangCS(player);
    // tiemNangUse = point * (2 * this.dameg + point - 1) / 2 * 100;
    // if ((this.dameg + point) <= getDameLimit()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // dameg += point;
    // }
    // TaskService.gI().checkDoneTaskNangCS(player);
    // } else {
    // Service.gI().sendThongBaoOK(player, "Vui lòng mở giới hạn sức mạnh");
    // return;
    // }
    // }
    // if (type == 3) {
    // tiemNangUse = 2 * (this.defg + 5) / 2 * 100000;
    // if ((this.defg + point) <= getDefLimit()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // defg += point;
    // }
    // } else {
    // Service.gI().sendThongBaoOK(player, "Vui lòng mở giới hạn sức mạnh");
    // return;
    // }
    // }
    // if (type == 4) {
    // tiemNangUse = 50000000L;
    // for (int i = 0; i < this.critg; i++) {
    // tiemNangUse *= 5L;
    // }
    // if ((this.critg + point) <= getCritLimit()) {
    // if (doUseTiemNang(tiemNangUse)) {
    // critg += point;
    // }
    // } else {
    // Service.gI().sendThongBaoOK(player, "Vui lòng mở giới hạn sức mạnh");
    // return;
    // }
    // }
    // Service.gI().point(player);
    // }
    private boolean doUseTiemNang(long tiemNang) {
        if (this.tiemNang < tiemNang) {
            Service.gI().sendThongBaoOK(player, "Bạn không đủ tiềm năng");
            return false;
        }
        if (this.tiemNang >= tiemNang && this.tiemNang - tiemNang >= 0) {
            this.tiemNang -= tiemNang;
            TaskService.gI().checkDoneTaskUseTiemNang(player);
            return true;
        }
        return false;
    }

    public long getFullTN() {
        long tnhp = 0, tnki = 0, tnsd = 0, tng = 0, tncm = 0;

        if (hpg > 0) {
            tnhp = (((hpg / 20L) * (50L + (50L + (hpg / 20L) - 1L)) / 2L) * 20L);
        }
        if (mpg > 0) {
            tnki = (((mpg / 20L) * (50L + (50L + (mpg / 20L) - 1L)) / 2L) * 20L);
        }
        if (dameg > 0) {
            tnsd = ((dameg * (dameg - 1L) * 100L) / 2L);
        }
        if (defg > 0) {
            tng = ((defg * (500000L + (500000L + (defg - 1L) * 100000L))) / 2L);
        }
        if (critg > 0) {
            tncm = ((50L * (((long) Math.pow(5L, critg) - 1L)) / (5L - 1L) * 1000000L));
        }
        return tnhp + tnki + tnsd + tng + tncm;
    }

    // --------------------------------------------------------------------------
    private long lastTimeHoiPhuc;
    private long lastTimeHoiStamina;

    public void update() {
        if (player != null && player.effectSkill != null) {
            if (player.effectSkill.isCharging && player.effectSkill.countCharging < 10) {
                int tiLeHoiPhuc = SkillUtil.getPercentCharge(player.playerSkill.skillSelect.point);
                if (player.effectSkill.isCharging && !player.isDie() && !player.effectSkill.isHaveEffectSkill()
                        && (hp < hpMax || mp < mpMax)) {
                    long hpRecovered = hpMax / 100 * tiLeHoiPhuc;
                    long mpRecovered = mpMax / 100 * tiLeHoiPhuc;

                    // if (hp + hpRecovered > 2_000_000_000) {
                    // hpRecovered = 2_000_000_000 - hp;
                    // }
                    // if (mp + mpRecovered > 2_000_000_000) {
                    // mpRecovered = 2_000_000_000 - mp;
                    // }
                    PlayerService.gI().hoiPhuc(player, Util.maxIntValue(hpRecovered), Util.maxIntValue(mpRecovered));

                    if (player.effectSkill.countCharging % 3 == 0) {
                        Service.gI().chat(player, "Phục hồi năng lượng " + getCurrPercentHP() + "%");
                    }
                } else {
                    EffectSkillService.gI().stopCharge(player);
                }
                if (++player.effectSkill.countCharging >= 10) {
                    EffectSkillService.gI().stopCharge(player);
                }
            }

            if (Util.canDoWithTime(lastTimeHoiPhuc, 30000)) {
                PlayerService.gI().hoiPhuc(this.player, Util.maxIntValue(hpHoi), Util.maxIntValue(mpHoi));
                this.lastTimeHoiPhuc = System.currentTimeMillis();
            }

            if (Util.canDoWithTime(lastTimeHoiStamina, 60000) && this.stamina < this.maxStamina) {
                this.stamina++;
                this.lastTimeHoiStamina = System.currentTimeMillis();

                if (!this.player.isBoss && !this.player.isPet) {
                    PlayerService.gI().sendCurrentStamina(this.player);
                }
            }
        }
    }

    public void dispose() {
        this.intrinsic = null;
        this.player = null;
        this.tlHp = null;
        this.tlMp = null;
        this.tlDef = null;
        this.tlDame = null;
        this.tlDameAttMob = null;
        this.tlTNSM = null;
    }

}



