
package player.Bot;

import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import server.ServerManager;
import services.ChatGlobalService;
import utils.Util;

/**
 *
 * @author SADBOI
 */
public class BotManager implements Runnable {

    private static final long GLOBAL_CHAT_INTERVAL = 240_000L;
    private static final String[] GLOBAL_CHAT_TEXTS = {
            "Rừng Đá k5 ttv Up Tnsm",
            "BDKB k2 ttv Up Tnsm (có trói)",
            "Núi Khỉ Đỏ k8 ttv Up Tnsm",
            "Xin slot Up TNSM ké",    "Newbie Up đâu cho nhanh thế mọi người",
            "Mua Vàng X cao",    "Mua Acc đã qua tương lai sll",
            "Mua 3 Sao qua ATM",    "Mua x2 Giầy Thần linh Trái Đất gấp",
            "Bố mày thách bọn mày qua được TĐST",    "Thuê người làm TĐST uy tín",
            "Súc vật phá TĐST cmm à",    "Phá cc, tao cũng nhiệm vụ",
            "Tìm pt săn Boss xuyên đêm"
    };

    public static BotManager i;

    public List<Bot> bot = new CopyOnWriteArrayList<>();
    private long lastGlobalChatTime;
    private long lastGlobalChatBotId = -1;

    public static BotManager gI() {
        if (i == null) {
            i = new BotManager();
        }
        return i;
    }

    @Override
    public void run() {
        while (ServerManager.isRunning) {
            long st = System.currentTimeMillis();
            try {
                for (Bot bot : this.bot) {
                    if (bot != null) {
                        bot.update();
                    }
                }
                this.chatGlobalForTrainingBot();
            } catch (Exception ignored) {
            }

            long sleepTime = Math.max(150 - (System.currentTimeMillis() - st), 10);
            try {
                Thread.sleep(sleepTime);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            }
        }
    }

    private void chatGlobalForTrainingBot() {
        if (this.bot.isEmpty() || !Util.canDoWithTime(this.lastGlobalChatTime, GLOBAL_CHAT_INTERVAL)) {
            return;
        }
        Bot chosenBot = null;
        if (this.bot.size() == 1) {
            chosenBot = this.bot.get(0);
        } else {
            for (int i = 0; i < 10; i++) {
                Bot randomBot = this.bot.get(Util.nextInt(0, this.bot.size() - 1));
                if (randomBot != null && randomBot.id != this.lastGlobalChatBotId) {
                    chosenBot = randomBot;
                    break;
                }
            }
            if (chosenBot == null) {
                chosenBot = this.bot.get(Util.nextInt(0, this.bot.size() - 1));
            }
        }
        if (chosenBot == null) {
            return;
        }
        String text = GLOBAL_CHAT_TEXTS[Util.nextInt(0, GLOBAL_CHAT_TEXTS.length - 1)];
        ChatGlobalService.gI().chat1(chosenBot, text);
        this.lastGlobalChatBotId = chosenBot.id;
        this.lastGlobalChatTime = System.currentTimeMillis();
    }
}
