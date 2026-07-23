package models.ShenronEvent;

/*
 *
 *
 * @author GiaCatLuong
 */

import network.Message;
import consts.ConstNpc;
import item.Item;
import player.Player;
import lombok.Getter;
import lombok.Setter;
import map.Zone;
import server.Client;
import services.InventoryService;
import services.ItemTimeService;
import services.NpcService;
import services.Service;
import utils.SkillUtil;
import utils.Util;

public class ShenronEvent {

    @Setter
    @Getter
    private Player player;

    @Setter
    @Getter
    private Zone zone;

    public long playerId;
    public boolean isPlayerDisconnect;
    public byte select;
    public int shenronType;
    public boolean leaveMap;

    public static final byte WISHED = 0;
    public static final byte TIME_UP = 1;

    public static final byte DRAGON_EVENT = 1;

    public long lastTimeShenronWait;
    public static int timeResummonShenron = 300000;
    public static int timeShenronWait = 300000;

    public static final String SHENRON_SAY = "Ta sẽ ban cho người 1 điều ước, \nngươi có 5 phút,\n hãy chọn đi:\n1) Đổi skill 2, 3 đệ tử (có thể trùng skill trước đó).\n2) Đổi skill 3, 4 đệ tử (có thể trùng skill trước đó).\n3) Tăng 10% HP, KI, SD trong 30p.\n4) Tăng 15% Giáp và 30% TNSM trong 30p.";

    public static final String[] SHENRON_WISHES = new String[] { "Điều ước 1", "Điều ước 2", "Điều ước 3",
            "Điều ước 4" };

    public boolean shenronLeave;

    public void update() {
        try {
            if (!shenronLeave) {
                if (isPlayerDisconnect) {
                    Player pl = Client.gI().getPlayer(playerId);
                    if (pl != null) {
                        player = pl;
                        if (player.zone != null && player.zone.map.mapId != 0 && player.zone.map.mapId != 7
                                && player.zone.map.mapId != 14
                                && player.zone.map.mapId != 21 && player.zone.map.mapId != 22
                                && player.zone.map.mapId != 23) {
                            player.shenronEvent = this;
                            zone = player.zone;
                            player.iDMark.setShenronType(shenronType);
                            isPlayerDisconnect = false;
                            reSummonShenron();
                        }
                    }
                }
                if (Util.canDoWithTime(lastTimeShenronWait, timeShenronWait)) {
                    leaveMap = true;
                    NpcService.gI().createMenuRongThieng(player, ConstNpc.IGNORE_MENU,
                            "Còn cái nịt =))\nCó không ước mất đừng tìm.", "Xin vĩnh biệt cụ........");
                    shenronLeave();
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void reSummonShenron() {
        activeShenron(true, DRAGON_EVENT);
        sendWhishesShenron();
    }

    public void sendWhishesShenron() {
        NpcService.gI().createMenuRongThieng(player, ConstNpc.SHOW_SHENRON_EVENT_CONFIRM, SHENRON_SAY, SHENRON_WISHES);
    }

    public void showConfirmShenron(byte select) {
        this.select = select;
        String wish = null;
        switch (player.iDMark.getShenronType()) {
            case 0:
                wish = SHENRON_WISHES[select];
                break;
        }
        NpcService.gI().createMenuRongThieng(player, ConstNpc.SHENRON_EVENT_CONFIRM, "Ngươi có chắc muốn ước?", wish,
                "Từ chối");
    }

    public void activeShenron(boolean appear, byte type) {
        Message msg;
        try {
            msg = new Message(-83);
            msg.writer().writeByte(appear ? 0 : (byte) 1);
            if (appear) {
                msg.writer().writeShort(player.zone.map.mapId);
                msg.writer().writeShort(player.zone.map.bgId);
                msg.writer().writeByte(player.zone.zoneId);
                msg.writer().writeInt((int) player.id);
                msg.writer().writeUTF("Ngoc Rong Kame");
                msg.writer().writeShort(player.location.x);
                msg.writer().writeShort(player.location.y);
                msg.writer().writeByte(type);
                playerId = player.id;
                shenronType = player.iDMark.getShenronType();
                zone.shenronType = shenronType;
                lastTimeShenronWait = System.currentTimeMillis();
                player.isShenronAppear = true;
            }
            Service.gI().sendMessAllPlayerInMap(player, msg);
        } catch (Exception e) {
        }
    }

    public void confirmWish() {
        switch (player.iDMark.getShenronType()) {
            case 0:
                switch (this.select) {
                    case 0: // thay chiêu 2-3 đệ tử
                        if (player.pet != null) {
                            if (player.pet.playerSkill.skills.get(1).skillId != -1) {
                                player.pet.openSkill2();
                                if (player.pet.playerSkill.skills.get(2).skillId != -1) {
                                    player.pet.openSkill3();
                                }
                            } else {
                                Service.gI().sendThongBao(player, "Ít nhất đệ tử ngươi phải có chiêu 2 chứ!");
                                sendWhishesShenron();
                                return;
                            }
                        } else {
                            Service.gI().sendThongBao(player, "Ngươi làm gì có đệ tử?");
                            sendWhishesShenron();
                            return;
                        }
                        break;
                    case 1: // thay chiêu 3-4 đệ tử
                        if (player.pet != null) {
                            if (player.pet.playerSkill.skills.get(2).skillId != -1) {
                                player.pet.openSkill3();
                                if (player.pet.playerSkill.skills.get(3).skillId != -1) {
                                    player.pet.openSkill4();
                                }
                            } else {
                                Service.gI().sendThongBao(player, "Ít nhất đệ tử ngươi phải có chiêu 3 chứ!");
                                sendWhishesShenron();
                                return;
                            }
                        } else {
                            Service.gI().sendThongBao(player, "Ngươi làm gì có đệ tử?");
                            sendWhishesShenron();
                            return;
                        }
                        break;
                    case 2:// Tăng hp, ki, sd
                        int timeRX = player.itemTime.timeRX / 1000 + 1800;
                        int maxTimeInSeconds = 32767;
                        if (timeRX >= maxTimeInSeconds) {
                            Service.gI().sendThongBao(player, "Ước ít thôi con :v");
                            sendWhishesShenron();
                            return;
                        }
                        player.itemTime.isUseRX = true;
                        player.itemTime.timeRX = timeRX * 1000;
                        player.itemTime.lastTimeUseRX = System.currentTimeMillis();
                        ItemTimeService.gI().sendItemTime(player, 6581, timeRX);
                        player.nPoint.calPoint();
                        player.nPoint.setHp(Util.maxIntValue(player.nPoint.hpMax));
                        player.nPoint.setMp(Util.maxIntValue(player.nPoint.mpMax));
                        Service.gI().point(player);
                        Service.gI().Send_Info_NV(player);
                        break;
                    case 3:// Tăng giáp, tnsm
                        int timeRX2 = player.itemTime.timeRX / 1000 + 1800;
                        int maxTimeInSeconds2 = 32767;
                        if (timeRX2 >= maxTimeInSeconds2) {
                            Service.gI().sendThongBao(player, "Ước ít thôi con :v");
                            sendWhishesShenron();
                            return;
                        }
                        player.itemTime.isUseRX2 = true;
                        player.itemTime.timeRX = timeRX2 * 1000;
                        player.itemTime.lastTimeUseRX = System.currentTimeMillis();
                        ItemTimeService.gI().sendItemTime(player, 6581, timeRX2);
                        player.nPoint.calPoint();
                        Service.gI().point(player);
                        Service.gI().Send_Info_NV(player);
                        break;
                }
                break;
        }
        shenronLeave();
    }

    public void shenronLeave() {
        if (!shenronLeave) {
            shenronLeave = true;
            if (player != null && player.zone != null) {
                player.shenronEvent = null;
                if (!leaveMap) {
                    NpcService.gI().createTutorial(player, 0, "Điều ước của ngươi đã được thực hiện...tạm biệt");
                }
                activeShenron(false, DRAGON_EVENT);
                player.isShenronAppear = false;
                select = -1;
            }
            zone.shenronType = -1;
            player.lastTimeShenronAppeared = System.currentTimeMillis();
            ShenronEventManager.gI().remove(this);
        }
    }
}
