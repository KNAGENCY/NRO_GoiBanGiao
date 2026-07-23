package services.func;

/*
 *
 *
 * @author GiaCatLuong
 */

import GiaCatLuongManager.Functions;
import network.Message;
import consts.ConstNpc;
import jdbc.daos.NDVSqlFetcher;
import jdbc.daos.PlayerDAO;
import item.Item;
import java.util.List;
import map.Zone;
import player.Player;
import server.Client;
import services.InventoryService;
import services.ItemService;
import services.ItemTimeService;
import services.NpcService;
import services.Service;
import utils.Util;

public class SummonDragonNamek {

    public static final byte DRAGON_PORUNGA = 1;
    private static SummonDragonNamek instance;

    public static final byte WISHED = 0;
    public static final byte TIME_UP = 1;
    private boolean isShenronAppear;
    public Player playerSummonShenron;
    private int playerSummonShenronId;
    private Zone mapShenronAppear;
    private int menuShenron;
    private byte select;
    private final Thread update;
    private boolean active;
    public boolean isPlayerDisconnect;
    private long lastTimeShenronWait;
    private final int timeShenronWait = 300000;

    public static SummonDragonNamek gI() {
        if (instance == null) {
            instance = new SummonDragonNamek();
        }
        return instance;
    }

    private SummonDragonNamek() {
        this.update = new Thread(() -> {
            while (active) {
                try {
                    if (isShenronAppear) {
                        if (isPlayerDisconnect) {
                            List<Player> players = mapShenronAppear.getPlayers();
                            for (Player plMap : players) {
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
                    Functions.sleep(1000);
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
        });
        this.active();
    }

    private void active() {
        if (!active) {
            active = true;
            this.update.start();
        }
    }

    public void summonNamec(Player pl) {
        if (pl.zone.map.mapId == 7) {
            playerSummonShenron = pl;
            playerSummonShenronId = (int) pl.id;
            mapShenronAppear = pl.zone;
            lastTimeShenronWait = System.currentTimeMillis();
            sendNotifyShenronNamekAppear();
            activeShenron(pl, true, DRAGON_PORUNGA);
            sendWhishesNamec(pl);
        } else {
            Service.gI().sendThongBao(pl, "Không thể thực hiện");
        }
    }

    private void reSummonShenron() {
        activeShenron(playerSummonShenron, true, DRAGON_PORUNGA);
        sendWhishesNamec(playerSummonShenron);
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
                msg.writer().writeUTF("Ngoc Rong Kame");
                msg.writer().writeShort(pl.location.x);
                msg.writer().writeShort(pl.location.y);
                msg.writer().writeByte(type);
                isShenronAppear = true;
            }
            Service.gI().sendMessAllPlayer(msg);
        } catch (Exception e) {
        }
    }

    private void sendNotifyShenronNamekAppear() {
        Message msg = null;
        try {
            msg = new Message(-25);
            msg.writer().writeUTF(playerSummonShenron.name + " vừa gọi rồng thần namek tại "
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
        switch (this.menuShenron) {
            case ConstNpc.SHOW_SHENRON_NAMEK_CONFIRM:
                try {
                    switch (select) {
                        case 0: // +20% sức đánh
                            if (playerSummonShenron.clan != null) {
                                playerSummonShenron.clan.members.forEach(m -> {
                                    Player p = Client.gI().getPlayer(m.id);
                                    if (p != null) {
                                        p.itemTime.isUsePorungaSD = true;
                                        p.itemTime.timePorungaSD = 6 * 60 * 60 * 1000; // 6h
                                        p.itemTime.lastTimePorungaSD = System.currentTimeMillis();
                                        ItemTimeService.gI().sendItemTime(p, 9999,
                                                (int) (p.itemTime.timePorungaSD / 1000));
                                        p.nPoint.calPoint();
                                        Service.gI().point(p);
                                        Service.gI().Send_Info_NV(p);
                                    }
                                });
                            }
                            break;

                        case 1: // +25% HP & KI
                            if (playerSummonShenron.clan != null) {
                                playerSummonShenron.clan.members.forEach(m -> {
                                    Player p = Client.gI().getPlayer(m.id);
                                    if (p != null) {
                                        p.itemTime.isUsePorungaHPKI = true;
                                        p.itemTime.timePorungaHPKI = 6 * 60 * 60 * 1000;
                                        p.itemTime.lastTimePorungaHPKI = System.currentTimeMillis();
                                        ItemTimeService.gI().sendItemTime(p, 10000,
                                                (int) (p.itemTime.timePorungaHPKI / 1000));
                                        p.nPoint.calPoint();
                                        p.nPoint.setHp(Util.maxIntValue(p.nPoint.hpMax));
                                        p.nPoint.setMp(Util.maxIntValue(p.nPoint.mpMax));
                                        Service.gI().point(p);
                                        Service.gI().Send_Info_NV(p);
                                    }
                                });
                            }
                            break;

                        case 2: // +3000 giáp
                            if (playerSummonShenron.clan != null) {
                                playerSummonShenron.clan.members.forEach(m -> {
                                    Player p = Client.gI().getPlayer(m.id);
                                    if (p != null) {
                                        p.itemTime.isUsePorungaGIAP = true;
                                        p.itemTime.timePorungaGIAP = 6 * 60 * 60 * 1000;
                                        p.itemTime.lastTimePorungaGIAP = System.currentTimeMillis();
                                        ItemTimeService.gI().sendItemTime(p, 10001,
                                                (int) (p.itemTime.timePorungaGIAP / 1000));
                                        p.nPoint.calPoint();
                                        Service.gI().point(p);
                                        Service.gI().Send_Info_NV(p);
                                    }
                                });
                            }
                            break;
                    }
                } catch (Exception e) {
                    e.printStackTrace();
                }
                break;
        }
        shenronLeave(this.playerSummonShenron, WISHED);
    }

    public void showConfirmShenron(Player pl, int menu, byte select) {
        this.menuShenron = menu;
        this.select = select;
        String wish = null;
        switch (menu) {
            case ConstNpc.SHOW_SHENRON_NAMEK_CONFIRM:
                switch (select) {
                    case 0:
                        wish = "20% SĐ\n6 giờ";
                        break;
                    case 1:
                        wish = "25% HP, KI\n6 giờ";
                        break;
                    case 2:
                        wish = "3000 Giáp\n 6 giờ";
                        break;
                }
                break;
        }
        NpcService.gI().createMenuRongThieng(pl, ConstNpc.SHENRON_NAMEK_CONFIRM,
                "Ngươi có chắc muốn ước điều này cho cả bang hội?", wish, "Từ chối");
    }

    public void sendWhishesNamec(Player pl) {
        NpcService.gI().createMenuRongThieng(pl, ConstNpc.SHOW_SHENRON_NAMEK_CONFIRM,
                "Ta sẽ ban cho cả bang hội ngươi 1 điều ước.\nNgươi có 5 phút, hãy suy nghĩ thật kỹ trước khi quyết định:",
                "20% SĐ\\n6 giờ",
                "25% HP, KI\\n6 giờ",
                "3000 Giáp\\n 6 giờ");
    }

    public void shenronLeave(Player pl, byte type) {
        if (type == WISHED) {
          
            NpcService.gI().createTutorial(pl, 0, "Điều ước của ngươi đã được thực hiện...tạm biệt");
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
        this.mapShenronAppear = null;
    }
}
