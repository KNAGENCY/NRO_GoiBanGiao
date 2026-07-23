package server;

/*
 *
 *
 * @author GiaCatLuong
 */
import GiaCatLuongManager.FileRunner;
import GiaCatLuongManager.Functions;
import models.Consign.ConsignShopManager;
import jdbc.daos.HistoryTransactionDAO;
import jdbc.daos.NDVSqlFetcher;
import boss.BossManager;
import boss.OtherBossManager;
import boss.FinalBossManager;
import boss.SkillSummonedManager;
import consts.ConstDataEventNAP;
import consts.ConstDataEventSM;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.io.*;
import java.lang.management.ManagementFactory;
import java.net.URLDecoder;
import jdbc.daos.PlayerDAO;
import player.Player;
import services.Service;
import network.inetwork.ISession;
import network.SessionManager;
import network.Network;
import server.io.MyKeyHandler;
import server.io.MySession;
import services.ClanService;
import services.ChatGlobalService;
import services.NgocRongNamecService;
import utils.Logger;
import utils.TimeUtil;
import java.util.*;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import utils.Util;
import models.DeathOrAliveArena.DeathOrAliveArenaManager;
import event.EventManager;
import java.io.DataOutputStream;
import item.Item;
import item.Item.ItemOption;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.JSONValue;
import services.ItemService;
import jdbc.daos.ShopDAO;

import jdbc.DBConnecter;
import jdbc.daos.EventDAO;
import network.MessageSendCollect;
import models.ShenronEvent.ShenronEventManager;
import models.SuperRank.SuperRankManager;
import network.Message;
import network.inetwork.ISessionAcceptHandler;
import player.Bot.BotManager;
import player.Bot.NewBot;

public class ServerManager {

    public static String timeStart;

    // public static final Map<String, Integer> CLIENTS = new HashMap<>();

    public static final Map<String, Integer> CLIENTS = new java.util.concurrent.ConcurrentHashMap<>();
    public static String NAME = "Local";
    public static String IP = "127.0.0.1";
    public static int PORT = 14445;

    private static ServerManager instance;

    public static boolean isRunning;
    private static final String DEFAULT_GAME_ROOT = "C:\\Users\\Administrator\\Music\\zz1";
    private static final String GAME_ROOT = resolveGameRoot();
    private static final String PANEL_COMMAND_FILE = new File(GAME_ROOT, "panel_cmd.txt").getPath();
    private static final String PANEL_EXP_STATUS_FILE = new File(GAME_ROOT, "panel_exp_status.txt").getPath();
    private static final String PANEL_DROP_STATUS_FILE = new File(GAME_ROOT, "panel_drop_status.txt").getPath();
    private static final String PANEL_ADMIN_MODE_STATUS_FILE = new File(GAME_ROOT, "panel_admin_mode_status.txt").getPath();
    private static final String PANEL_BOT_STATUS_FILE = new File(GAME_ROOT, "panel_bot_status.txt").getPath();
    // Tên file cũ của tính năng này (thời chủ cũ đặt là "XNXX"). Chỉ còn dùng để
    // đọc lại trạng thái bật/tắt đúng một lần sau khi đổi tên, rồi thôi.
    private static final String PANEL_BOT_STATUS_FILE_LEGACY = new File(GAME_ROOT, "panel_xnxx_status.txt").getPath();
    private static final String PANEL_BOSS_PERSIST_STATUS_FILE = new File(GAME_ROOT, "panel_boss_persist_status.txt").getPath();
    private static final String PANEL_SERVER_STATUS_FILE = new File(GAME_ROOT, "panel_server_status.json").getPath();
    private static final Map<Integer, Long> PANEL_KICK_LOCKS = new java.util.concurrent.ConcurrentHashMap<>();
    private static final int PANEL_KICK_SECONDS_DEFAULT = 300;
    private static volatile boolean panelAdminMode = false;

    private static String resolveGameRoot() {
        String root = System.getProperty("game.root");
        if (root == null || root.trim().isEmpty()) {
            root = System.getenv("GAME_ROOT");
        }
        if (root != null && !root.trim().isEmpty()) {
            return root.trim();
        }

        String userDir = System.getProperty("user.dir");
        if (isGameRoot(userDir)) {
            return new File(userDir).getAbsolutePath();
        }
        if (userDir != null) {
            File parent = new File(userDir).getAbsoluteFile().getParentFile();
            if (parent != null && isGameRoot(parent.getPath())) {
                return parent.getAbsolutePath();
            }
        }
        for (String candidate : defaultGameRootCandidates()) {
            if (isGameRoot(candidate)) {
                return new File(candidate).getAbsolutePath();
            }
        }
        return DEFAULT_GAME_ROOT;
    }

    private static List<String> defaultGameRootCandidates() {
        List<String> candidates = new ArrayList<>();
        String userHome = System.getProperty("user.home");
        if (userHome != null && !userHome.trim().isEmpty()) {
            candidates.add(new File(userHome, "Music\\zz1").getPath());
        }
        candidates.add(DEFAULT_GAME_ROOT);
        candidates.add("C:\\Users\\Admin\\Music\\zz1");
        return candidates;
    }

    public static int getPanelKickRemainingSeconds(int userId) {
        Long lockedUntil = PANEL_KICK_LOCKS.get(userId);
        if (lockedUntil == null) {
            return 0;
        }
        long remainingMs = lockedUntil - System.currentTimeMillis();
        if (remainingMs <= 0) {
            PANEL_KICK_LOCKS.remove(userId);
            return 0;
        }
        return (int) Math.ceil(remainingMs / 1000.0);
    }

    public static boolean isPanelAdminMode() {
        return panelAdminMode;
    }

    private static boolean isGameRoot(String path) {
        if (path == null || path.trim().isEmpty()) {
            return false;
        }
        return new File(path, "data\\config\\config.properties").isFile();
    }

    public void init() {
        Manager.gI();
        ConsignShopManager.gI().load();
        HistoryTransactionDAO.deleteHistory();
    }

    public static ServerManager gI() {
        if (instance == null) {
            instance = new ServerManager();
            instance.init();
        }
        return instance;
    } 

    public static void main(String[] args) {
        timeStart = TimeUtil.getTimeNow("dd/MM/yyyy HH:mm:ss");
        new server.ServerManagerUI().setVisible(true);
    }

    public void run() {
        isRunning = true;
        activeServerSocket();
        loadPanelAdminModeStatus();

        // ===== THREAD GAME / EVENT =====
        new Thread(BotManager.gI(), "Thread Bot Game").start();
        NewBot.gI().syncAvailability(loadPanelBotModeStatus());
        new Thread(NgocRongNamecService.gI(), "Update NRNM").start();
        new Thread(SuperRankManager.gI(), "Update Super Rank").start();
        // new Thread(The23rdMartialArtCongressManager.gI(), "Update DHVT23").start();
        new Thread(DeathOrAliveArenaManager.gI(), "Update Võ Đài Sinh Tử").start();
        new Thread(AutoMaintenance.gI(), "Update Bảo Trì Tự Động").start();
        new Thread(ShenronEventManager.gI(), "Update Shenron").start();

        // ===== LOAD BOSS =====
        BossManager.gI().loadBoss();
        Manager.MAPS.forEach(map.Map::initBoss);
        EventManager.gI().init();
        startPanelCommandWatcher();
        startPanelServerStatusExporter();
        // ===== GAME LOOP =====
        GameLoop.start();
        // ===== START SCHEDULER BOSS =====
        BossManager.gI().start();
        // ===== EVENT / SPECIAL =====
        services.func.SummonDragon.gI().start();
        // ===== CẬP NHẬT TRẠNG THÁI SK =====
        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            try {
                ConstDataEventSM.isRunningSK = ConstDataEventSM.isActiveEvent();
                ConstDataEventNAP.isRunningSK = ConstDataEventNAP.isActiveEvent();
            } catch (Exception e) {
                Logger.logException(ServerManager.class, e, "Update SK error");
            }
        }, 0, 500, TimeUnit.MILLISECONDS);

        // ===== RESET HỒI SINH HÀNG NGÀY =====
        startSystemGlobalNoticeScheduler();
        scheduleResetHoiSinh();
    }


    private static boolean loadPanelBotModeStatus() {
        File statusFile = new File(PANEL_BOT_STATUS_FILE);
        if (!statusFile.isFile()) {
            // Lần chạy đầu sau khi đổi tên: lấy lại trạng thái từ file cũ nếu còn,
            // để bot không bị bật/tắt ngược so với lúc trước khi cập nhật.
            File legacy = new File(PANEL_BOT_STATUS_FILE_LEGACY);
            if (legacy.isFile()) {
                boolean enabled = readBotStatusFile(legacy);
                writePanelBotModeStatus(enabled);
                return enabled;
            }
            writePanelBotModeStatus(true);
            return true;
        }
        return readBotStatusFile(statusFile);
    }

    private static boolean readBotStatusFile(File statusFile) {
        try (BufferedReader reader = new BufferedReader(new FileReader(statusFile))) {
            String value = reader.readLine();
            return value == null || !value.trim().equals("0");
        } catch (IOException e) {
            Logger.logException(ServerManager.class, e, "Failed to read bot status");
            return true;
        }
    }

    private static void writePanelBotModeStatus(boolean enabled) {
        try (FileWriter writer = new FileWriter(PANEL_BOT_STATUS_FILE, false)) {
            writer.write(enabled ? "1" : "0");
        } catch (IOException e) {
            Logger.logException(ServerManager.class, e, "Failed to write bot status");
        }
    }

    /**
     * Bật/tắt bot đánh quái tự động.
     *
     * notifyPlayers hiện luôn được gọi với false: thông báo ra toàn server sẽ
     * nói thẳng cho người chơi biết những nhân vật kia là bot, làm mất tác dụng
     * của chính tính năng này. Trạng thái vẫn hiện đầy đủ ở trang quản trị và
     * trên màn hình máy chủ. Muốn bật lại thì đổi tham số ở chỗ gọi.
     */
    public static void setPanelBotMode(boolean enabled, boolean notifyPlayers) {
        NewBot.gI().syncAvailability(enabled);
        if (enabled) {
            if (notifyPlayers) {
                Service.gI().sendThongBaoAllPlayer("Bot đánh quái đã được bật.");
            }
            Logger.log(Logger.YELLOW, "Panel command: BOT DANH QUAI ON\n");
        } else {
            if (NewBot.gI().isTrainingBotsActive()) {
                NewBot.gI().stopTrainingBots();
            }
            if (notifyPlayers) {
                Service.gI().sendThongBaoAllPlayer("Bot đánh quái đã được tắt.");
            }
            Logger.log(Logger.YELLOW, "Panel command: BOT DANH QUAI OFF\n");
        }
        writePanelBotModeStatus(enabled);
    }
    private void startSystemGlobalNoticeScheduler() {
        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            try {
                ChatGlobalService.gI().chatSystemNotice();
            } catch (Exception e) {
                Logger.logException(ServerManager.class, e, "System global notice error");
            }
        }, 10, 10, TimeUnit.MINUTES);
    }

    private void startPanelServerStatusExporter() {
        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            try {
                exportPanelServerStatus();
            } catch (Exception e) {
                Logger.logException(ServerManager.class, e, "Panel server status export error");
            }
        }, 0, 1, TimeUnit.SECONDS);
    }

    private void exportPanelServerStatus() throws IOException {
        Runtime runtime = Runtime.getRuntime();
        long maxMemory = runtime.maxMemory();
        long usedMemory = runtime.totalMemory() - runtime.freeMemory();
        double heapRamPercent = maxMemory > 0 ? (usedMemory * 100.0 / maxMemory) : 0;
        double systemRamPercent = 0;
        double processCpuPercent = 0;
        double systemCpuPercent = 0;

        java.lang.management.OperatingSystemMXBean baseBean = ManagementFactory.getOperatingSystemMXBean();
        if (baseBean instanceof com.sun.management.OperatingSystemMXBean) {
            com.sun.management.OperatingSystemMXBean osBean = (com.sun.management.OperatingSystemMXBean) baseBean;
            double processCpuLoad = osBean.getProcessCpuLoad();
            double systemCpuLoad = osBean.getSystemCpuLoad();
            long totalMemory = osBean.getTotalPhysicalMemorySize();
            long freeMemory = osBean.getFreePhysicalMemorySize();
            processCpuPercent = processCpuLoad >= 0 ? processCpuLoad * 100.0 : 0;
            systemCpuPercent = systemCpuLoad >= 0 ? systemCpuLoad * 100.0 : 0;
            systemRamPercent = totalMemory > 0 ? ((totalMemory - freeMemory) * 100.0 / totalMemory) : 0;
        }

        String json = "{"
                + "\"generatedAt\":" + System.currentTimeMillis() + ","
                + "\"online\":" + Client.gI().getPlayers().size() + ","
                + "\"thread\":" + Thread.activeCount() + ","
                + "\"session\":" + SessionManager.gI().getNumSession() + ","
                + "\"cpu\":" + formatPanelPercent(processCpuPercent) + ","
                + "\"cpu2\":" + formatPanelPercent(systemCpuPercent) + ","
                + "\"ram\":" + formatPanelPercent(heapRamPercent) + ","
                + "\"ram2\":" + formatPanelPercent(systemRamPercent)
                + "}";

        File statusFile = new File(PANEL_SERVER_STATUS_FILE);
        File parent = statusFile.getParentFile();
        if (parent != null && !parent.isDirectory()) {
            parent.mkdirs();
        }
        try (FileWriter writer = new FileWriter(statusFile, false)) {
            writer.write(json);
        }
    }

    private static String formatPanelPercent(double value) {
        if (Double.isNaN(value) || Double.isInfinite(value)) {
            value = 0;
        }
        value = Math.max(0, Math.min(100, value));
        return String.format(java.util.Locale.US, "%.2f", value);
    }

    private void startPanelCommandWatcher() {
        new Thread(() -> {
            File commandFile = new File(PANEL_COMMAND_FILE);
            while (isRunning) {
                try {
                    if (commandFile.exists() && commandFile.isFile()) {
                        List<String> commands = new ArrayList<>();
                        try (BufferedReader reader = new BufferedReader(new FileReader(commandFile))) {
                            String line;
                            while ((line = reader.readLine()) != null) {
                                line = line.trim();
                                if (!line.isEmpty()) {
                                    commands.add(line);
                                }
                            }
                        }
                        if (!commands.isEmpty()) {
                            processPanelCommands(commands);
                        }
                        commandFile.delete();
                    }
                    NewBot.gI().syncAvailability(loadPanelBotModeStatus());
                } catch (Exception ex) {
                    Logger.logException(ServerManager.class, ex, "Panel command watcher error");
                }
                try {
                    Thread.sleep(2000);
                } catch (InterruptedException ignored) {
                }
            }
        }, "PanelCommandWatcher").start();
    }

    private void processPanelCommands(List<String> commands) {
        for (String command : commands) {
            if (command.startsWith("MAINTENANCE_NOW")) {
                Logger.log(Logger.YELLOW, "Panel command: SCHEDULE MAINTENANCE 30s\n");
                Maintenance.gI().start(30);
                return;
            } else if (command.startsWith("ADMIN_MODE|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    boolean enabled = "1".equals(parts[1].trim()) || "ON".equalsIgnoreCase(parts[1].trim());
                    setPanelAdminMode(enabled, true);
                }
            } else if (command.startsWith("BOT_MODE|") || command.startsWith("XNXX_MODE|")) {
                // Vẫn chấp nhận tên lệnh cũ "XNXX_MODE" để bản web cũ trên VPS
                // không bị chết nút trong lúc chưa kịp cập nhật cùng lúc.
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    boolean enabled = "1".equals(parts[1].trim()) || "ON".equalsIgnoreCase(parts[1].trim());
                    setPanelBotMode(enabled, false);
                }
            } else if (command.startsWith("BOSS_PERSIST_MODE|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    boolean enabled = "1".equals(parts[1].trim()) || "ON".equalsIgnoreCase(parts[1].trim());
                    try (FileWriter writer = new FileWriter(PANEL_BOSS_PERSIST_STATUS_FILE, false)) {
                        writer.write(enabled ? "1" : "0");
                    } catch (IOException e) {
                        Logger.logException(ServerManager.class, e, "Failed to write boss persist status");
                    }
                    Logger.log(Logger.YELLOW, "Panel command: BOSS PERSIST " + (enabled ? "ON" : "OFF") + "\n");
                }
            } else if (command.startsWith("BUFF_VND|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 3) {
                    String charName = parts[1].trim();
                    int amount = 0;
                    try {
                        amount = Integer.parseInt(parts[2].trim());
                    } catch (NumberFormatException e) {
                        Logger.logException(ServerManager.class, e, "Invalid buff amount");
                        continue;
                    }
                    if (amount <= 0) {
                        continue;
                    }
                    if (PlayerDAO.addcash(charName, amount)) {
                        Player target = Client.gI().getPlayerByName(charName);
                        if (target != null && target.session != null) {
                            target.session.cash += amount;
                            Service.gI().sendThongBao(target, "Bạn vừa nhận được " + amount + " Cá từ Admin.");
                        }
                        Logger.log(Logger.YELLOW, "Panel command: Buff Cá " + amount + " cho " + charName + "\n");
                    } else {
                        Logger.log(Logger.YELLOW, "Panel command: Không tìm thấy nhân vật " + charName + "\n");
                    }
                }
            } else if (command.startsWith("SET_EXP|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    try {
                        int rate = Integer.parseInt(parts[1].trim());
                        if (rate <= 0) {
                            Logger.log(Logger.YELLOW, "Panel command: EXP rate phai > 0\n");
                            continue;
                        }
                        Manager.RATE_EXP_SERVER = rate;
                        try (FileWriter writer = new FileWriter(PANEL_EXP_STATUS_FILE, false)) {
                            writer.write(String.valueOf(rate));
                        } catch (IOException e) {
                            Logger.logException(ServerManager.class, e, "Failed to write EXP status");
                        }
                        Service.gI().sendThongBaoAllPlayer("EXP server hiện tại: x" + rate);
                        Logger.log(Logger.YELLOW, "Panel command: SET EXP x" + rate + "\n");
                    } catch (NumberFormatException e) {
                        Logger.logException(ServerManager.class, e, "Invalid EXP rate");
                    }
                }
            } else if (command.startsWith("SET_DROP|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    try {
                        int rate = Integer.parseInt(parts[1].trim());
                        if (rate <= 0) {
                            Logger.log(Logger.YELLOW, "Panel command: DROP rate phai > 0\n");
                            continue;
                        }
                        Manager.RATE_DROP_SERVER = rate;
                        try (FileWriter writer = new FileWriter(PANEL_DROP_STATUS_FILE, false)) {
                            writer.write(String.valueOf(rate));
                        } catch (IOException e) {
                            Logger.logException(ServerManager.class, e, "Failed to write DROP status");
                        }
                        Service.gI().sendThongBaoAllPlayer("Tỉ lệ rơi đồ server hiện tại: x" + rate);
                        Logger.log(Logger.YELLOW, "Panel command: SET DROP x" + rate + "\n");
                    } catch (NumberFormatException e) {
                        Logger.logException(ServerManager.class, e, "Invalid DROP rate");
                    }
                }
            } else if (command.equalsIgnoreCase("LS") || command.startsWith("UPDATE_SHOP")) {
                try (Connection con = DBConnecter.getConnectionServer()) {
                    Manager.SHOPS = ShopDAO.getShops(con);
                    Logger.log(Logger.YELLOW, "Panel command: LS UPDATE SHOP " + Manager.SHOPS.size() + "\n");
                } catch (Exception e) {
                    Logger.logException(ServerManager.class, e, "Panel command LS update shop failed");
                }
            } else if (command.startsWith("EDIT_BOSS|")) {
                String[] parts = command.split("\\|", 7);
                if (parts.length >= 7) {
                    try {
                        int bossId = Integer.parseInt(parts[1].trim());
                        int mapId = Integer.parseInt(parts[2].trim());
                        int zoneId = Integer.parseInt(parts[3].trim());
                        long hp = Long.parseLong(parts[4].trim());
                        long dame = Long.parseLong(parts[5].trim());
                        String bossName = URLDecoder.decode(parts[6].trim(), "UTF-8");
                        if (bossId >= 0 || hp <= 0 || hp > 2_000_000_000L || dame <= 0 || dame > 2_000_000_000L) {
                            Logger.log(Logger.YELLOW, "Panel command: EDIT_BOSS du lieu khong hop le\n");
                            continue;
                        }
                        boolean edited = BossManager.gI().editBossFromPanel(bossId, mapId, zoneId, hp, dame, bossName);
                        if (edited) {
                            Logger.log(Logger.YELLOW, "Panel command: EDIT BOSS " + bossId + " hp=" + hp + " dame=" + dame + "\n");
                        } else {
                            Logger.log(Logger.YELLOW, "Panel command: Khong tim thay boss dang song de sua " + bossId + "\n");
                        }
                    } catch (Exception e) {
                        Logger.logException(ServerManager.class, e, "Invalid EDIT_BOSS command");
                    }
                }
            } else if (command.startsWith("SPAWN_BOSS|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    try {
                        int bossId = Integer.parseInt(parts[1].trim());
                        int mapId = parts.length >= 3 ? Integer.parseInt(parts[2].trim()) : -1;
                        int zoneId = parts.length >= 4 ? Integer.parseInt(parts[3].trim()) : -1;
                        if (bossId >= 0) {
                            Logger.log(Logger.YELLOW, "Panel command: Boss ID khong hop le\n");
                            continue;
                        }
                        boolean spawned = BossManager.gI().respawnBossFromPanel(bossId, mapId, zoneId);
                        if (spawned) {
                            Service.gI().sendThongBaoAllPlayer("Admin vua goi boss moi.");
                            Logger.log(Logger.YELLOW, "Panel command: SPAWN BOSS " + bossId + "\n");
                        } else {
                            Logger.log(Logger.YELLOW, "Panel command: Khong the goi boss " + bossId + "\n");
                        }
                    } catch (NumberFormatException e) {
                        Logger.logException(ServerManager.class, e, "Invalid BOSS id");
                    }
                }
            } else if (command.startsWith("BUFF_ITEM|")) {
                String[] parts = command.split("\\|", 3);
                if (parts.length >= 3) {
                    processPanelBuffItem(parts[1].trim(), parts[2].trim());
                }
            } else if (command.startsWith("KICK_PLAYER|")) {
                String[] parts = command.split("\\|");
                if (parts.length >= 2) {
                    int seconds = PANEL_KICK_SECONDS_DEFAULT;
                    if (parts.length >= 3) {
                        try {
                            seconds = Integer.parseInt(parts[2].trim());
                        } catch (NumberFormatException e) {
                            Logger.logException(ServerManager.class, e, "Invalid KICK seconds");
                        }
                    }
                    processPanelKickPlayer(parts[1].trim(), seconds);
                }
            }
        }
    }

    private static void loadPanelAdminModeStatus() {
        File statusFile = new File(PANEL_ADMIN_MODE_STATUS_FILE);
        if (!statusFile.isFile()) {
            panelAdminMode = false;
            return;
        }
        try (BufferedReader reader = new BufferedReader(new FileReader(statusFile))) {
            String value = reader.readLine();
            panelAdminMode = value != null && value.trim().equals("1");
            if (panelAdminMode) {
                Logger.log(Logger.YELLOW, "Panel admin mode: ON\n");
            }
        } catch (IOException e) {
            Logger.logException(ServerManager.class, e, "Failed to read admin mode status");
        }
    }

    private void setPanelAdminMode(boolean enabled, boolean notifyPlayers) {
        panelAdminMode = enabled;
        try (FileWriter writer = new FileWriter(PANEL_ADMIN_MODE_STATUS_FILE, false)) {
            writer.write(enabled ? "1" : "0");
        } catch (IOException e) {
            Logger.logException(ServerManager.class, e, "Failed to write admin mode status");
        }

        if (enabled) {
            if (notifyPlayers) {
                Service.gI().sendThongBaoAllPlayer("Server dang bat che do bao tri admin. Tai khoan thuong se bi ngat ket noi.");
            }
            List<Player> playersCopy = new ArrayList<>(Client.gI().getPlayers());
            for (Player player : playersCopy) {
                if (player != null && !player.isAdmin() && player.getSession() != null) {
                    Service.gI().sendThongBao(player, "Server dang bao tri admin. Vui long quay lai sau.");
                    Client.gI().kickSession(player.getSession());
                }
            }
            Logger.log(Logger.YELLOW, "Panel command: ADMIN MODE ON\n");
        } else {
            if (notifyPlayers) {
                Service.gI().sendThongBaoAllPlayer("Server da tat che do bao tri admin.");
            }
            Logger.log(Logger.YELLOW, "Panel command: ADMIN MODE OFF\n");
        }
    }

    private void processPanelKickPlayer(String charName, int seconds) {
        if (charName == null || charName.isEmpty()) {
            return;
        }
        if (seconds <= 0) {
            seconds = PANEL_KICK_SECONDS_DEFAULT;
        }

        Player target = Client.gI().getPlayerByName(charName);
        int userId = target != null && target.getSession() != null ? target.getSession().userId : findAccountIdByPlayerName(charName);
        if (userId <= 0) {
            Logger.log(Logger.YELLOW, "Panel command: Khong tim thay nhan vat de kick " + charName + "\n");
            return;
        }
        if (target == null) {
            target = Client.gI().getPlayerByUser(userId);
        }

        PANEL_KICK_LOCKS.put(userId, System.currentTimeMillis() + seconds * 1000L);
        if (target != null && target.getSession() != null) {
            Service.gI().sendThongBao(target, "Bạn bị kick khỏi game. Vui lòng đăng nhập lại sau 5 phút.");
            Client.gI().kickSession(target.getSession());
        }
        Logger.log(Logger.YELLOW, "Panel command: KICK " + charName + " khoa " + seconds + " giay\n");
    }

    private int findAccountIdByPlayerName(String charName) {
        try (Connection con = DBConnecter.getConnectionServer();
                PreparedStatement ps = con.prepareStatement("select account_id from player where name = ? limit 1")) {
            ps.setString(1, charName);
            try (ResultSet rs = ps.executeQuery()) {
                if (rs.next()) {
                    return rs.getInt("account_id");
                }
            }
        } catch (Exception e) {
            Logger.logException(ServerManager.class, e, "Panel KICK find player failed");
        }
        return -1;
    }

    private void processPanelBuffItem(String charName, String payload) {
        try {
            Object parsed = JSONValue.parse(payload);
            if (!(parsed instanceof JSONArray)) {
                Logger.log(Logger.RED, "Panel command: BUFF_ITEM payload khong hop le\n");
                return;
            }
            JSONArray items = (JSONArray) parsed;
            if (items.isEmpty()) {
                return;
            }

            if ("all".equalsIgnoreCase(charName)) {
                List<Player> allPlayer = NDVSqlFetcher.getAllPlayer();
                for (Player pBuffItem : allPlayer) {
                    if (pBuffItem != null) {
                        applyPanelBuffItems(pBuffItem, items);
                    }
                }
                Logger.log(Logger.YELLOW, "Panel command: Buff item cho tat ca nhan vat\n");
                return;
            }

            Player target = Client.gI().getPlayerByName(charName);
            if (target == null) {
                target = NDVSqlFetcher.loadPlayerByName(charName);
            }
            if (target == null) {
                Logger.log(Logger.YELLOW, "Panel command: Khong tim thay nhan vat " + charName + "\n");
                return;
            }

            if (applyPanelBuffItems(target, items)) {
                Logger.log(Logger.YELLOW, "Panel command: Buff item cho " + target.name + "\n");
                if (target.session != null) {
                    Service.gI().sendThongBao(target, "Ban vua nhan vat pham tu Admin trong Hom Thu.");
                }
            }
        } catch (Exception e) {
            Logger.logException(ServerManager.class, e, "Panel BUFF_ITEM error");
        }
    }

    private boolean applyPanelBuffItems(Player target, JSONArray items) {
        boolean added = false;
        for (Object obj : items) {
            if (!(obj instanceof JSONObject)) {
                continue;
            }
            JSONObject itemData = (JSONObject) obj;
            int tempId = toInt(itemData.get("temp_id"), -1);
            int quantity = toInt(itemData.get("quantity"), 0);
            if (tempId <= 0 || quantity <= 0 || quantity > 20000) {
                continue;
            }

            Item itemBuff = ItemService.gI().createNewItem((short) tempId, quantity);
            Object optionsObj = itemData.get("options");
            if (optionsObj instanceof JSONArray) {
                JSONArray options = (JSONArray) optionsObj;
                for (Object optObj : options) {
                    if (!(optObj instanceof JSONObject)) {
                        continue;
                    }
                    JSONObject opt = (JSONObject) optObj;
                    int optId = toInt(opt.get("id"), -1);
                    int param = toInt(opt.get("param"), 0);
                    if (optId >= 0) {
                        itemBuff.itemOptions.add(new ItemOption(optId, param));
                    }
                }
            }
            target.inventory.itemsMailBox.add(itemBuff);
            added = true;
        }
        return added && NDVSqlFetcher.updateMailBox(target);
    }

    private int toInt(Object value, int defaultValue) {
        if (value instanceof Number) {
            return ((Number) value).intValue();
        }
        if (value != null) {
            try {
                return Integer.parseInt(value.toString());
            } catch (NumberFormatException ignored) {
            }
        }
        return defaultValue;
    }

    private void activeServerSocket() {
        try {
            Network.gI().init().setAcceptHandler(new ISessionAcceptHandler() {
                @Override
                public void sessionInit(ISession is) {
                    if (!canConnectWithIp(is.getIP())) {
                        is.disconnect();
                        return; 
                    }
                    is.setMessageHandler(Controller.gI())
                            .setSendCollect(new MessageSendCollect() {
                                @Override
                                public void doSendMessage(ISession session, DataOutputStream dos, Message msg)
                                        throws Exception {
                                    byte[] data = msg.getData();
                                    if (session.sentKey()) {
                                        dos.writeByte(this.writeKey(session, msg.command));
                                    } else {
                                        dos.writeByte(msg.command);
                                    }

                                    if (data != null) {
                                        int size = data.length;
                                        if (msg.command == -32 || msg.command == -66 || msg.command == -74
                                                || msg.command == 11 || msg.command == -67 || msg.command == -87
                                                || msg.command == 66) {
                                            dos.writeByte(this.writeKey(session, (byte) size) - 128);
                                            dos.writeByte(this.writeKey(session, (byte) (size >> 8)) - 128);
                                            dos.writeByte(this.writeKey(session, (byte) (size >> 16)) - 128);
                                        } else if (session.sentKey()) {
                                            dos.writeByte(this.writeKey(session, (byte) (size >> 8)));
                                            dos.writeByte(this.writeKey(session, (byte) (size & 0xFF)));
                                        } else {
                                            dos.writeShort(size);
                                        }
                                        if (session.sentKey()) {
                                            for (int i = 0; i < data.length; ++i) {
                                                data[i] = this.writeKey(session, data[i]);
                                            }
                                        }
                                        dos.write(data);
                                    } else {
                                        dos.writeShort(0);
                                    }
                                    dos.flush();
                                    msg.cleanup();
                                }
                            })
                            .setKeyHandler(new MyKeyHandler())
                            .startCollect();
                }

            @Override
public void sessionDisconnect(ISession session) {
    Client.gI().kickSession((MySession) session);
}
            }).setTypeSessioClone(MySession.class)
                    .setDoSomeThingWhenClose(() -> {
                        Logger.error("SERVER CLOSE\n");
                        System.exit(0);
                    })
                    .start(PORT);
        } catch (Exception e) {
            Logger.logException(ServerManager.class, e, "Error starting server socket");
        }
    }

    // private boolean canConnectWithIp(String ipAddress) {
    //     CLIENTS.putIfAbsent(ipAddress, 0);
    //     int n = CLIENTS.get(ipAddress);
    //     if (n < Manager.MAX_PER_IP) {
    //         CLIENTS.put(ipAddress, n + 1);
    //         return true;
    //     }
    //     return false;
    // }
  private boolean canConnectWithIp(String ipAddress) {
    int[] result = {0};
    CLIENTS.compute(ipAddress, (key, val) -> {
        int current = (val == null) ? 0 : val;
        if (current < Manager.MAX_PER_IP) {
            result[0] = 1; // cho phép
            return current + 1;
        }
        result[0] = 0; // chặn
        return current;
    });
    return result[0] == 1;
}
    public void disconnect(MySession session) {
        // String ip = session.getIP();
        // CLIENTS.put(ip, Math.max(CLIENTS.getOrDefault(ip, 0) - 1, 0));
    }

    public void close() {
        isRunning = false;
        try { ClanService.gI().close(); } catch (Exception e) { Logger.error("Lỗi save clan!\n"); }
        try { ConsignShopManager.gI().save(); } catch (Exception e) { Logger.error("Lỗi save shop ký gửi!\n"); }
        Client.gI().close();
        EventDAO.save();
        Logger.success("SUCCESSFULLY MAINTENANCE!\n");

        // Thoát với mã 2 để run.bat biết là cần KHỞI ĐỘNG LẠI.
        // Trước đây chỗ này tự gọi FileRunner.runBatchFile("run.bat") — cách đó mở
        // một cửa sổ console MỚI mỗi lần bảo trì, còn cửa sổ cũ kẹt vĩnh viễn ở lệnh
        // PAUSE. Bảo trì nhiều lần là màn hình đầy cửa sổ chết.
        // Nay run.bat đứng canh và tự chạy lại, luôn chỉ có đúng một cửa sổ.
        //   mã 2 = bảo trì, chạy lại ngay
        //   mã 0 = dừng hẳn (nút "Save data" và "Thoát" trong giao diện dùng mã này)
        //   mã khác = server lỗi, run.bat tự chạy lại sau 10 giây
        System.exit(2);
    }

    private void scheduleResetHoiSinh() {
        ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);
        long initialDelay = Util.getDelayToMidnight();
        long period = 24 * 60 * 60;

        scheduler.scheduleAtFixedRate(() -> {
            try (Connection con = DBConnecter.getConnectionServer();
                 PreparedStatement st = con.prepareStatement("UPDATE player SET hoisinh = 0")) {
                st.executeUpdate();
                System.out.println("Reset số lần hồi sinh thành công lúc 0h");
            } catch (Exception e) {
                e.printStackTrace();
            }
        }, initialDelay, period, TimeUnit.SECONDS);
    }
}



