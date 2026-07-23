package services.func;

/*
 *
 *
 * @author GiaCatLuong
 */
import GiaCatLuongManager.Functions;
import java.util.HashMap;
import java.util.Map;
import item.Item;
import consts.ConstNpc;
import consts.ConstPlayer;

import item.Item.ItemOption;
import map.Zone;
import services.NpcService;
import player.Inventory;
import player.Player;
import services.Service;
import services.ItemTimeService;
import utils.Util;
import network.Message;
import services.ItemService;
import services.PlayerService;
import services.InventoryService;
import utils.Logger;
import java.util.List;
import server.Maintenance;

public class SummonDragon {

    public static final byte WISHED = 0;
    public static final byte TIME_UP = 1;

    public static final byte DRAGON_SHENRON = 0;
    public static final byte DRAGON_PORUNGA = 1;

    public static final short NGOC_RONG_1_SAO = 14;
    public static final short NGOC_RONG_2_SAO = 15;
    public static final short NGOC_RONG_3_SAO = 16;
    public static final short NGOC_RONG_4_SAO = 17;
    public static final short NGOC_RONG_5_SAO = 18;
    public static final short NGOC_RONG_6_SAO = 19;
    public static final short NGOC_RONG_7_SAO = 20;

    public static final String SUMMON_SHENRON_TUTORIAL = "Chỉ có thể gọi rồng 1 sao.\n"
            + "Để gọi rồng 1 sao cần đủ ngọc rồng từ 1 sao đến 7 sao.\n"
            + "Ngọc rồng sẽ mất ngay khi gọi rồng dù bạn có ước hay không.\n"
            + "Quá 5 phút nếu không ước rồng thần sẽ bay mất";
    public static final String SHENRON_SAY = "Ta sẽ ban cho ngươi 1 điều ước,\n ngươi có 5 phút, hãy chọn đi:\n1. Cải Trang Vip\n2. Hp + 20% trong 2 tiếng\n3. Ki + 20% trong 2 tiếng\n4. Sức đánh + 15% trong 2 tiếng\n5. Giáp + 3000 trong 2 tiếng\n6. TNSM + 30% trong 2 tiếng\n7. Hàng trang + 50 TV";

    public static final String[] SHENRON_1_STAR_WISHES_1 = new String[] {
           "Đẹp Trai\nNhất",
    "Hp Trâu\nNhất",
    "Ki Lớn\nNhất",
    "Dame To\nNhất",
    "Giáp Cao\nNhất",
    "Up Nhanh\nNhất",
    "Giàu Có\nNhất"
};
    public static final String[] SHENRON_1_STAR_WISHES_2 = SHENRON_1_STAR_WISHES_1;
    public static final String[] SHENRON_2_STARS_WHISHES = new String[] { "Chỉ có thể gọi rồng 1 sao" };
    public static final String[] SHENRON_3_STARS_WHISHES = new String[] { "Chỉ có thể gọi rồng 1 sao" };
    // --------------------------------------------------------------------------
    private static SummonDragon instance;
    private final Map pl_dragonStar;
    private long lastTimeShenronAppeared;
    private long lastTimeShenronWait;
    private final int timeResummonShenron = 600000;
    private boolean isShenronAppear;
    private final int timeShenronWait = 300000;

    private final Thread update;
    private boolean active;

    public boolean isPlayerDisconnect;
    public Player playerSummonShenron;
    private int playerSummonShenronId;
    private Zone mapShenronAppear;
    private byte shenronStar;
    private int menuShenron;
    private byte select;

private SummonDragon() {
    this.pl_dragonStar = new HashMap<>();
    this.update = new Thread(() -> {
        while (active && !Maintenance.isRunning) {
            try {
                updateShenron();
                Thread.sleep(1000);
            } catch (Exception e) {
                Logger.logException(SummonDragon.class, e);
            }
        }
    }, "Update Shenron Dragon");
}

public void start() {
    if (!active) {
        active = true;
        update.start();
    }
}

private void updateShenron() {
    if (!isShenronAppear) return;

    if (isPlayerDisconnect && mapShenronAppear != null) {
        for (Player plMap : mapShenronAppear.getPlayers()) {
            if (plMap.isPl() && plMap.id == playerSummonShenronId) {
                playerSummonShenron = plMap;
                reSummonShenron();
                isPlayerDisconnect = false;
                break;
            }
        }
    }

    if (Util.canDoWithTime(lastTimeShenronWait, timeShenronWait)) {
        shenronLeave(playerSummonShenron, TIME_UP);
    }
}

    private void active() {
        if (!active) {
            active = true;
            this.update.start();
        }
    }

    public static SummonDragon gI() {
        if (instance == null) {
            instance = new SummonDragon();
        }
        return instance;
    }

    public void openMenuSummonShenron(Player pl, byte dragonBallStar) {
        if (dragonBallStar != 1) {
            Service.gI().sendThongBao(pl, "Chỉ có thể gọi rồng 1 sao");
            return;
        }
        this.pl_dragonStar.put(pl, dragonBallStar);
        NpcService.gI().createMenuConMeo(pl, ConstNpc.SUMMON_SHENRON, -1, "Bạn muốn gọi rồng thần 1 sao ?",
                "Hướng\ndẫn thêm", "Gọi\nRồng Thần\n1 Sao");
    }
    public void summonShenron(Player pl) {
        Byte dragonStarCall = (Byte) this.pl_dragonStar.get(pl);
        if (dragonStarCall == null || dragonStarCall != 1) {
            Service.gI().sendThongBao(pl, "Chỉ có thể gọi rồng 1 sao");
            return;
        }
        if (pl.zone.map.mapId == 0 || pl.zone.map.mapId == 7 || pl.zone.map.mapId == 14) {
            if (checkShenronBall(pl)) {
                if (isShenronAppear) {
                    Service.gI().sendThongBao(pl, "Không thể thực hiện");
                    return;
                }

                if (Util.canDoWithTime(lastTimeShenronAppeared, timeResummonShenron)) {
                    // gọi rồng
                    playerSummonShenron = pl;
                    playerSummonShenronId = (int) pl.id;
                    mapShenronAppear = pl.zone;
                    byte dragonStar = (byte) pl_dragonStar.get(playerSummonShenron);
                    int begin = NGOC_RONG_1_SAO;
                    switch (dragonStar) {
                        case 2:
                            begin = NGOC_RONG_2_SAO;
                            break;
                        case 3:
                            begin = NGOC_RONG_3_SAO;
                            break;
                    }
                    for (int i = begin; i <= NGOC_RONG_7_SAO; i++) {
                        try {
                            InventoryService.gI().subQuantityItemsBag(pl, InventoryService.gI().findItemBag(pl, i), 1);
                        } catch (Exception ex) {
                        }
                    }
                    InventoryService.gI().sendItemBag(pl);
                    sendNotifyShenronAppear();
                    activeShenron(pl, true, SummonDragon.DRAGON_SHENRON);
                    sendWhishesShenron(pl);
                } else {
                    int timeLeft = (int) ((timeResummonShenron - (System.currentTimeMillis() - lastTimeShenronAppeared))
                            / 1000);
                    Service.gI().sendThongBao(pl, "Vui lòng đợi "
                            + (timeLeft < 7200 ? timeLeft + " giây" : timeLeft / 60 + " phút") + " nữa");
                }
            }
        } else {
            Service.gI().sendThongBao(pl, "Chỉ được gọi rồng thần ở ngôi làng trước nhà");
        }
    }

    private void reSummonShenron() {
        activeShenron(playerSummonShenron, true, SummonDragon.DRAGON_SHENRON);
        sendWhishesShenron(playerSummonShenron);
    }

    private void sendWhishesShenron(Player pl) {
        byte dragonStar;
        try {
            dragonStar = (byte) pl_dragonStar.get(pl);
            this.shenronStar = dragonStar;
        } catch (Exception e) {
            dragonStar = this.shenronStar;
        }
        if (dragonStar != 1) {
            Service.gI().sendThongBao(pl, "Chỉ có thể gọi rồng 1 sao");
            return;
        }
        NpcService.gI().createMenuRongThieng(pl, ConstNpc.SHENRON_1_1, SHENRON_SAY, SHENRON_1_STAR_WISHES_1);
    }
    private void activeShenron(Player pl, boolean appear, byte type) {
        Message msg;
        try {
            msg = new Message(-83);
            msg.writer().writeByte(appear ? 0 : (byte) 1);
            if (appear) {
                msg.writer().writeShort(pl.zone.map.mapId);
                msg.writer().writeShort(pl.zone.map.bgId);
                msg.writer().writeByte(pl.zone.zoneId);
                msg.writer().writeInt((int) pl.id);
                msg.writer().writeUTF("DucVuPro");
                msg.writer().writeShort(pl.location.x);
                msg.writer().writeShort(pl.location.y);
                msg.writer().writeByte(type);
                lastTimeShenronWait = System.currentTimeMillis();
                isShenronAppear = true;
                pl.iDMark.setShenronType(-1);
            }
            Service.gI().sendMessAllPlayer(msg);
        } catch (Exception e) {
        }
    }

    private boolean checkShenronBall(Player pl) {
        byte dragonStar = (byte) this.pl_dragonStar.get(pl);
        if (dragonStar == 1) {
            if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_2_SAO)) {
                Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 2 sao");
                return false;
            }
            if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_3_SAO)) {
                Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 3 sao");
                return false;
            }
        } else if (dragonStar == 2) {
            if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_3_SAO)) {
                Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 3 sao");
                return false;
            }
        }
        if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_4_SAO)) {
            Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 4 sao");
            return false;
        }
        if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_5_SAO)) {
            Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 5 sao");
            return false;
        }
        if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_6_SAO)) {
            Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 6 sao");
            return false;
        }
        if (!InventoryService.gI().isExistItemBag(pl, NGOC_RONG_7_SAO)) {
            Service.gI().sendThongBao(pl, "Bạn còn thiếu 1 viên ngọc rồng 7 sao");
            return false;
        }
        return true;
    }

    private void sendNotifyShenronAppear() {
        Message msg = null;
        try {
            msg = new Message(-25);
            msg.writer().writeUTF(playerSummonShenron.name + " vừa gọi rồng thần tại "
                    + playerSummonShenron.zone.map.mapName + " khu vực " + playerSummonShenron.zone.zoneId);
            Service.gI().sendMessAllPlayerIgnoreMe(playerSummonShenron, msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    public void confirmWish() {
        boolean success = false;
        switch (this.menuShenron) {
            case ConstNpc.SHENRON_1_1:
                switch (this.select) {
                    case 0:
                        success = wishOutfit(playerSummonShenron);
                        break;
                    case 1:
                        success = wishTimedBuff(playerSummonShenron, 1);
                        break;
                    case 2:
                        success = wishTimedBuff(playerSummonShenron, 2);
                        break;
                    case 3:
                        success = wishTimedBuff(playerSummonShenron, 3);
                        break;
                    case 4:
                        success = wishTimedBuff(playerSummonShenron, 4);
                        break;
                    case 5:
                        success = wishTimedBuff(playerSummonShenron, 5);
                        break;
                    case 6:
                        success = wishGold(playerSummonShenron);
                        break;
                }
                break;
        }
        if (success) {
            shenronLeave(this.playerSummonShenron, WISHED);
        }
    }

    public void showConfirmShenron(Player pl, int menu, byte select) {
        if (menu != ConstNpc.SHENRON_1_1 || select < 0 || select >= SHENRON_1_STAR_WISHES_1.length) {
            Service.gI().sendThongBao(pl, "Chỉ có thể gọi rồng 1 sao");
            return;
        }
        this.menuShenron = menu;
        this.select = select;
        String wish = SHENRON_1_STAR_WISHES_1[select];
        NpcService.gI().createMenuRongThieng(pl, ConstNpc.SHENRON_CONFIRM,
                "Ngươi có chắc muốn ước:\n" + wish.replace("\n", " ") + "?", "Đồng ý", "Quay lại");
    }
    public void reOpenShenronWishes(Player pl) {
        NpcService.gI().createMenuRongThieng(pl, ConstNpc.SHENRON_1_1, SHENRON_SAY, SHENRON_1_STAR_WISHES_1);
    }

    private boolean wishOutfit(Player player) {
    if (player == null) {
        return false;
    }

    if (InventoryService.gI().getCountEmptyBag(player) < 1) {
        Service.gI().sendThongBao(player, "Cần ít nhất 1 ô trống trong hành trang");
        reOpenShenronWishes(player);
        return false;
    }

    short itemId;
    switch (player.gender) {
        case 0:
            itemId = 604;
            break;
        case 1:
            itemId = 605;
            break;
        default:
            itemId = 606;
            break;
    }

    Item item = ItemService.gI().createNewItem(itemId);

    item.itemOptions.add(new ItemOption(50, Util.nextInt(1, 44)));
    item.itemOptions.add(new ItemOption(77, Util.nextInt(1, 44)));
    item.itemOptions.add(new ItemOption(103, Util.nextInt(1, 44)));
    item.itemOptions.add(new ItemOption(47, Util.nextInt(1, 4444)));
    item.itemOptions.add(new ItemOption(93, Util.nextInt(1, 3)));

    InventoryService.gI().addItemBag(player, item);
    InventoryService.gI().sendItemBag(player);

    Service.gI().sendThongBao(player, "Bạn nhận được " + item.template.name);
    return true;
}
    private boolean wishTimedBuff(Player player, int type) {
        if (player == null) {
            return false;
        }
        int seconds = 7200;
        long now = System.currentTimeMillis();
        switch (type) {
            case 1:
                player.itemTime.isUseShenronHp = true;
                player.itemTime.lastTimeShenronHp = now;
                player.itemTime.timeShenronHp = seconds * 1000;
                ItemTimeService.gI().sendItemTime(player, 16586, seconds);
                Service.gI().sendThongBao(player, "HP của bạn tăng 20% trong 2 tiếng");
                break;
            case 2:
                player.itemTime.isUseShenronKi = true;
                player.itemTime.lastTimeShenronKi = now;
                player.itemTime.timeShenronKi = seconds * 1000;
                ItemTimeService.gI().sendItemTime(player, 16585, seconds);
                Service.gI().sendThongBao(player, "KI của bạn tăng 20% trong 2 tiếng");
                break;
            case 3:
                player.itemTime.isUseShenronDame = true;
                player.itemTime.lastTimeShenronDame = now;
                player.itemTime.timeShenronDame = seconds * 1000;
                ItemTimeService.gI().sendItemTime(player, 16582, seconds);
                Service.gI().sendThongBao(player, "Sức đánh của bạn tăng 15% trong 2 tiếng");
                break;
            case 4:
                player.itemTime.isUseShenronArmor = true;
                player.itemTime.lastTimeShenronArmor = now;
                player.itemTime.timeShenronArmor = seconds * 1000;
                ItemTimeService.gI().sendItemTime(player, 16584, seconds);
                Service.gI().sendThongBao(player, "Giáp của bạn tăng 3000 trong 2 tiếng");
                break;
            case 5:
                player.itemTime.isUseShenronTnsm = true;
                player.itemTime.lastTimeShenronTnsm = now;
                player.itemTime.timeShenronTnsm = seconds * 1000;
                ItemTimeService.gI().sendItemTime(player, 16583, seconds);
                Service.gI().sendThongBao(player, "TNSM nhận được tăng 30% trong 2 tiếng");
                break;
            default:
                return false;
        }
        player.nPoint.calPoint();
        player.nPoint.setHp(Util.maxIntValue(player.nPoint.hpMax));
        player.nPoint.setMp(Util.maxIntValue(player.nPoint.mpMax));
        Service.gI().point(player);
        Service.gI().Send_Info_NV(player);
        return true;
    }

    private boolean wishGold(Player player) {
        if (player == null) {
            return false;
        }
        if (InventoryService.gI().getCountEmptyBag(player) < 1) {
            Service.gI().sendThongBao(player, "Cần ít nhất 1 ô trống trong hành trang");
            reOpenShenronWishes(player);
            return false;
        }
        Item item = ItemService.gI().createNewItem((short) 457, 50);
        InventoryService.gI().addItemBag(player, item);
        InventoryService.gI().sendItemBag(player);
        Service.gI().sendThongBao(player, "Bạn nhận được 50 thỏi vàng");
        return true;
    }
    public void shenronLeave(Player pl, byte type) {
        if (type == WISHED) {
            NpcService.gI().createTutorial(pl, 0,
                    "Điều ước của ngươi đã trở thành sự thật\nHẹn gặp ngươi lần sau, ta đi ngủ đây, bái bai");
        } else {
            NpcService.gI().createMenuRongThieng(pl, ConstNpc.IGNORE_MENU,
                    "Ta buồn ngủ quá rồi\nHẹn gặp ngươi lần sau, ta đi đây, bái bai");
        }
        activeShenron(pl, false, SummonDragon.DRAGON_SHENRON);
        this.isShenronAppear = false;
        this.menuShenron = -1;
        this.select = -1;
        this.playerSummonShenron = null;
        this.playerSummonShenronId = -1;
        this.shenronStar = -1;
        this.mapShenronAppear = null;
        lastTimeShenronAppeared = System.currentTimeMillis();
    }

    // nâng đồ sư phụ
    private boolean tangCapTrangBi(Player player, int index, String tenTrangBi) {
        Item item = player.inventory.itemsBody.get(index);
        if (item.isNotNullItem()) {
            int level = 0;
            for (ItemOption io : item.itemOptions) {
                if (io.optionTemplate.id == 72) {
                    level = io.param;
                    if (level < 7) {
                        io.param++;
                    }
                    break;
                }
            }
            if (level < 7) {
                if (level == 0) {
                    item.itemOptions.add(new ItemOption(72, 1));
                }

                int optionId = switch (index) {
                    case 0 -> 47;
                    case 1 -> 6;
                    case 2 -> 0;
                    case 3 -> 7;
                      case 4 -> 14;
                    default -> -1;
                };
                if (optionId != -1) {
                    for (ItemOption io : item.itemOptions) {
                        if (io.optionTemplate.id == optionId) {
                            io.param += (io.param * 10 / 100);
                            break;
                        }
                    }
                }
                InventoryService.gI().sendItemBody(player);
                Service.gI().sendThongBao(player, "Đêm qua em tuyệt lắm");
                return true;
            } else {
                Service.gI().sendThongBao(player, "Đêm qua em tuyệt lắm");
                reOpenShenronWishes(player);
                return false;
            }
        } else {
            Service.gI().sendThongBao(player, "Có đâu mà nâng mă");
            reOpenShenronWishes(player);
            return false;
        }
    }

    // nâng đồ đệ
    private boolean tangCapTrangBiPet(Player player, int index, String tenTrangBi) {
        if (player.pet != null) {
            Item item = player.pet.inventory.itemsBody.get(index);
            if (item.isNotNullItem()) {
                int level = 0;
                for (ItemOption io : item.itemOptions) {
                    if (io.optionTemplate.id == 72) {
                        level = io.param;
                        if (level < 7) {
                            io.param++;
                        }
                        break;
                    }
                }
                if (level < 7) {
                    if (level == 0) {
                        item.itemOptions.add(new ItemOption(72, 1));
                    }
                    int optionId = switch (index) {
                        case 0 -> 47;
                        case 1 -> 6;
                        case 2 -> 0;
                        case 3 -> 7;
                           case 4 -> 14;
                        default -> -1;
                    };
                    if (optionId != -1) {
                        for (ItemOption io : item.itemOptions) {
                            if (io.optionTemplate.id == optionId) {
                                io.param += (io.param * 10 / 100);
                                break;
                            }
                        }
                    }
                    Service.gI().point(player);
                    Service.gI().sendThongBao(player, "Đêm qua em tuyệt lắm");
                    return true;
                } else {
                    Service.gI().sendThongBao(player, "Đêm qua em tuyệt lắm");
                    reOpenShenronWishes(player);
                    return false;
                }
            } else {
                Service.gI().sendThongBao(player, "Có đâu mà nâng mă");
                reOpenShenronWishes(player);
                return false;
            }
        } else {
            Service.gI().sendThongBao(player, "Có đâu mà nâng mă");
            reOpenShenronWishes(player);
            return false;
        }
    }

}
