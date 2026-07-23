package services;

/*
 *
 *
 * @author GiaCatLuong
 */
import GiaCatLuongManager.Functions;
import player.Player;
import network.Message;
import utils.Logger;
import utils.TimeUtil;
import utils.Util;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import server.Maintenance;

public class ChatGlobalService implements Runnable {

    private static final int COUNT_CHAT = 100;
    private static final int COUNT_WAIT = 100;
    private static ChatGlobalService i;
    private static final String SYSTEM_CHAT_NAME = "Hệ Thống";
    private static final short SYSTEM_CHAT_HEAD = 1685;
    private static final short SYSTEM_CHAT_BODY = 1686;
    private static final short SYSTEM_CHAT_LEG = 1687;
    private static final String[] SYSTEM_NOTICES = {
            "1. Spl mặc Ct Ghost Rider để up.\n"
                    + "2. Nr up ở Núi Khỉ Đỏ, Rừng Đá + Săn Boss + Đổi Điểm.\n"
                    + "3. SKH cơ bản up 3 Map đầu & Đảo Kame.\n"
                    + "4. SKH VIP nâng đồ Thần Linh.\n"
                    + "5. Nhận Cá Ib Key Vàng Trong Zalo",
            "6. Dưới 150tr thì vào Ngũ Hành Sơn.\n"
                    + "7. Trên 150tr thì vào BĐKB ở Quy Lão.\n"
                    + "8. 12H Leo tháp ở Ôsin Nhận vàng ở ĐHVT.\n"
                    + "9. 21h30 vào Map War Cày vàng.\n"
                    + "10. 22H pem Mabư ở Ôsin tại Đảo Kame"
    };
    public static final String SYSTEM_NOTICE_DETAIL_1 = "1. Spl mặc Ct Sói Héc Quyn để up.\n"
            + "2. Nr up ở Núi Khỉ Đỏ, Rừng Đá + Săn Boss + Đổi Điểm.\n"
            + "3. SKH cơ bản up 3 Map đầu & Đảo Kame.\n"
            + "4. SKH VIP nâng đồ Thần Linh";
    public static final String SYSTEM_NOTICE_DETAIL_2 = "5. Dưới 150tr thì vào Ngũ Hành Sơn.\n"
            + "6. Trên 150tr thì vào BĐKB ở Quy Lão.\n"
            + "7. 12H Leo tháp ở Ôsin Nhận vàng ở ĐHVT.\n"
            + "8. 21h30 vào Map War Cày vàng.\n"
            + "9. 22H pem Mabư ở Ôsin tại Đảo Kame";
    private final List<ChatGlobal> listChatting;
    private final List<ChatGlobal> waitingChat;

    private ChatGlobalService() {
        this.listChatting = new ArrayList<>();
        this.waitingChat = new LinkedList<>();
        this.start();
    }

    private void start() {
        new Thread(this, "**Chat global").start();
    }

    public static ChatGlobalService gI() {
        if (i == null) {
            i = new ChatGlobalService();
        }
        return i;
    }

    public void chatVip(Player player, String text) {
        waitingChat.add(new ChatGlobal(player, text.length() > 100 ? text.substring(0, 100) : text));
    }

    public void chat1(Player player, String text) {
        player.iDMark.setLastTimeChatGlobal(System.currentTimeMillis());
        waitingChat.add(new ChatGlobal(player, text.length() > 100 ? text.substring(0, 100) : text));
    }

    public void chatSystemNotice() {
        for (String text : SYSTEM_NOTICES) {
            waitingChat.add(new ChatGlobal(SYSTEM_CHAT_NAME, text));
        }
    }

    public void chat(Player player, String text) {
        if (waitingChat.size() >= COUNT_WAIT) {
            Service.gI().sendThongBao(player, "Kênh thế giới hiện đang quá tải, không thể chat lúc này");
            return;
        }
        boolean haveInChatting = false;
        for (ChatGlobal chat : listChatting) {
            if (chat.text.equals(text)) {
                haveInChatting = true;
                break;
            }
        }
        if (haveInChatting) {
            return;
        }
        if (player.inventory.gem >= 5) {
            if (player.isAdmin() || Util.canDoWithTime(player.iDMark.getLastTimeChatGlobal(), 30000)) {
                if (player.isAdmin() || player.nPoint.power > 2000000000) {
                    player.inventory.subGemAndRuby(5);
                    Service.gI().sendMoney(player);
                    player.iDMark.setLastTimeChatGlobal(System.currentTimeMillis());
                    waitingChat.add(new ChatGlobal(player, text.length() > 100 ? text.substring(0, 100) : text));
                } else {
                    Service.gI().sendThongBao(player, "Sức mạnh phải ít nhất 2tỷ mới có thể chat thế giới");
                }
            } else {
                Service.gI().sendThongBao(player, "Không thể chat thế giới lúc này, vui lòng đợi "
                        + TimeUtil.getTimeLeft(player.iDMark.getLastTimeChatGlobal(), 30));
            }
        } else {
            Service.gI().sendThongBao(player, "Không đủ ngọc chat thế giới");
        }
    }

    @Override
    public void run() {
        while (!Maintenance.isRunning) {
            try {
                if (!listChatting.isEmpty()) {
                    ChatGlobal chat = listChatting.get(0);
                    if (Util.canDoWithTime(chat.timeSendToPlayer, 1000)) {
                        listChatting.remove(0).dispose();
                    }
                }

                if (!waitingChat.isEmpty()) {
                    ChatGlobal chat = waitingChat.get(0);
                    if (listChatting.size() < COUNT_CHAT) {
                        waitingChat.remove(0);
                        chat.timeSendToPlayer = System.currentTimeMillis();
                        listChatting.add(chat);
                        chatGlobal(chat);
                    }
                }
                Functions.sleep(1000);
            } catch (Exception e) {
                Logger.logException(ChatGlobalService.class, e);
            }
        }
    }

    private void chatGlobal(ChatGlobal chat) {
        Message msg;
        try {
            msg = new Message(92);
            msg.writer().writeUTF(chat.playerName);
            msg.writer().writeUTF("|5|" + chat.text);
            msg.writer().writeInt((int) chat.playerId);
            msg.writer().writeShort(chat.head);
            msg.writer().writeShort(-1);
            msg.writer().writeShort(chat.body);
            msg.writer().writeShort(chat.bag); // bag
            msg.writer().writeShort(chat.leg);
            msg.writer().writeByte(0);
            Service.gI().sendMessAllPlayer(msg);
            msg.cleanup();
        } catch (Exception e) {
        }
    }

    private void transformText(ChatGlobal chat) {
        String text = chat.text;
        text = text.replaceAll("\\.com", "***")
                .replaceAll("\\.net", "***")
                .replaceAll("\\.xyz", "***")
                .replaceAll("\\.me", "***")
                .replaceAll("\\.pro", "***")
                .replaceAll("\\.mobi", "***")
                .replaceAll("\\.online", "***")
                .replaceAll("\\.info", "***")
                .replaceAll("\\.tk", "***")
                .replaceAll("\\.ml", "***")
                .replaceAll("\\.ga", "***")
                .replaceAll("\\.gq", "***")
                .replaceAll("\\.io", "***")
                .replaceAll("\\.club", "***")
                .replaceAll("cltx", "***")
                .replaceAll("cl", "***")
                .replaceAll("địt", "***")
                .replaceAll("lồn", "***")
                .replaceAll("sv admin", "admin đẹp trai")
                .replaceAll("sv lùa gà", "sv vip thế")
                .replaceAll("lùa", "***")
                .replaceAll("lùa gà", "***")
                .replaceAll("server lùa gà", "server vip thế")
                .replaceAll("cặc", "***");
        chat.text = text;
    }

    private class ChatGlobal {

        public String playerName;
        public int playerId;
        public short head;
        public short body;
        public short leg;
        public short bag;
        public String text;
        public long timeSendToPlayer;

        public ChatGlobal(String playerName, String text) {
            this.playerName = playerName;
            this.playerId = 0;
            this.head = SYSTEM_CHAT_HEAD;
            this.body = SYSTEM_CHAT_BODY;
            this.leg = SYSTEM_CHAT_LEG;
            this.bag = -1;
            this.text = text;
        }

        public ChatGlobal(Player player, String text) {
            if (!player.isAdmin()) {
                this.playerName = player.name;
            } else {
                this.playerName = player.name + " - Quản Trị Viên";
            }
            this.playerId = (int) player.id;
            this.head = player.getHead();
            this.body = player.getBody();
            this.leg = player.getLeg();
            this.bag = player.getFlagBag();
            this.text = text;
            transformText(this);
        }

        private void dispose() {
            this.playerName = null;
            this.text = null;
        }

    }

}
