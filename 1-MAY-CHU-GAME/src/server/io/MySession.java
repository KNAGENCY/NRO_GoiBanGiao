// package server.io;

// import java.net.Socket;

// import player.Player;
// import server.Controller;
// import data.DataGame;
// import jdbc.daos.NDVSqlFetcher;
// import item.Item;
// import java.io.IOException;
// import network.Session;
// import network.SessionManager;
// import network.inetwork.ISession;
// import network.Message;
// import server.Client;
// import server.Maintenance;
// import server.Manager;
// import models.AntiLogin;
// import services.Service;
// import utils.Logger;

// import java.util.HashMap;
// import java.util.List;
// import java.util.Map;
// import utils.TimeUtil;
// import utils.Util;

// public class MySession extends Session {

//     private static final Map<String, AntiLogin> ANTILOGIN = new HashMap<>();
//     private static final java.util.concurrent.ConcurrentHashMap<String, Object> LOGIN_LOCKS = new java.util.concurrent.ConcurrentHashMap<>();

//     public Player player;

//     public byte timeWait = 100;
//     public boolean sentKey;
//     public int packetCount = 0;
//     public long lastPacketTime = 0;
//     public static final byte[] KEYS = { 0 };
//     public byte curR, curW;

//     public String ipAddress;
//     public boolean isAdmin;
//     public int userId;
//     public String uu;
//     public String pp;

//     public int typeClient;
//     public byte zoomLevel;

//     public long lastTimeLogout;
//     public boolean joinedGame;

//     public long lastTimeReadMessage;

//     public boolean actived;

//     public boolean check;

//     public int goldBar;
//     public long gold;
//     public int eventPoint;
//     public List<Item> itemsReward;
//     public String dataReward;
//     public boolean is_gift_box;
//     public double bdPlayer;

//     public int version;
//     public int cash;
//     public int phaobong;
//     public int lixi;
//     public int danap;
//     public int vip;
//     public int luotquay;
//     public int diemdanh;
//     public boolean finishUpdate;

//     public long firstDataRequestTime;
//     public int dataRequestCount;

//     public MySession(Socket socket) {
//         super(socket);
//         ipAddress = socket.getInetAddress().getHostAddress();
//     }

//     @Override
//     public void sendKey() throws Exception {
//         super.sendKey();
//         this.startSend();
//     }

//     public void sendSessionKey() {
//         Message msg = new Message(-27);
//         try {
//             msg.writer().writeByte(KEYS.length);
//             msg.writer().writeByte(KEYS[0]);
//             for (int i = 1; i < KEYS.length; i++) {
//                 msg.writer().writeByte(KEYS[i] ^ KEYS[i - 1]);
//             }
//             this.sendMessage(msg);
//             msg.cleanup();
//             sentKey = true;
//         } catch (IOException e) {
//         }
//     }

//     public void login(String username, String password) {
//         long startTime = System.currentTimeMillis();
//         Logger.warning("[LOGIN-ATTEMPT] user=" + username + " ip=" + this.ipAddress);

//         AntiLogin al = ANTILOGIN.get(this.ipAddress);
//         if (al == null) {
//             al = new AntiLogin();
//             ANTILOGIN.put(this.ipAddress, al);
//         }

//         if (!al.canLogin()) {
//             Logger.warning("[LOGIN-BLOCK] reason=AntiLogin user=" + username + " ip=" + this.ipAddress);
//             Service.gI().sendThongBaoOK(this, al.getNotifyCannotLogin());
//             return;
//         }

//         if (Manager.LOCAL) {
//             Logger.warning("[LOGIN-BLOCK] reason=LOCAL user=" + username + " ip=" + this.ipAddress);
//             Service.gI().sendThongBaoOK(this,
//                     "Server này chỉ để lưu dữ liệu\nVui lòng qua server khác");
//             return;
//         }

//         if (Maintenance.isRunning) {
//             Logger.warning("[LOGIN-BLOCK] reason=Maintenance user=" + username + " ip=" + this.ipAddress);
//             Service.gI().sendThongBaoOK(this,
//                     "Server đang trong thời gian bảo trì, vui lòng quay lại sau");
//             return;
//         }

//         if (!this.isAdmin && Client.gI().getPlayers().size() >= Manager.MAX_PLAYER) {
//             Logger.warning("[LOGIN-BLOCK] reason=ServerFull user=" + username
//                     + " players=" + Client.gI().getPlayers().size() + "/" + Manager.MAX_PLAYER);
//             Service.gI().sendThongBaoOK(this,
//                     "Máy chủ hiện đang quá tải, cư dân vui lòng di chuyển sang máy chủ khác.");
//             return;
//         }

//         String userKey = username.toLowerCase();
//         // BUG-FIX #2: never remove the lock from the map. Removing it allows two
//         // threads
//         // to wind up holding different lock objects for the same username.
//         Object lock = LOGIN_LOCKS.computeIfAbsent(userKey, k -> new Object());

//         synchronized (lock) {
//             // ----- Sweep #1: active in-game player with this username -----
//             // BUG-FIX #1: kick stale, then FALL THROUGH (no return) so the new session can
//             // log in.
//             // BUG-FIX #3: skip `this` defensively.
//             // for (Player p : Client.gI().getPlayers()) {
//             //     if (p == null || p.getSession() == null || p.getSession().uu == null)
//             //         continue;
//             //     if (p.getSession() == this)
//             //         continue;
//             //     if (!p.getSession().uu.equalsIgnoreCase(username))
//             //         continue;



// List<Player> playersCopy = new java.util.ArrayList<>(Client.gI().getPlayers());
// for (Player p : playersCopy) {
//     if (p == null || p.getSession() == null || p.getSession().uu == null)
//         continue;
//     if (p.getSession() == this)
//         continue;
//     if (!p.getSession().uu.equalsIgnoreCase(username))
//         continue;






//                 MySession old = (MySession) p.getSession();
//                 Logger.warning("[LOGIN-DUP-PLAYER] user=" + username
//                         + " newIp=" + this.ipAddress
//                         + " oldIp=" + old.ipAddress
//                         + " -> kicking old, then proceeding");
//                 Service.gI().sendThongBaoOK(p, "Mất kết nối với máy chủ!");
//                 Client.gI().kickSession(old);
//                 waitForKickCleanup(username);
//                 break; // only one duplicate is possible
//             }

//             // ----- Sweep #2: any other matching session in SessionManager -----
//             // BUG-FIX #3: dropped the `timeWait > 0` filter — kick all matching dupes.
//             // for (ISession s : SessionManager.gI().getSessions()) {
//             //     MySession ms = (MySession) s;
//             //     if (ms == this)
//             //         continue;
//             //     if (ms.uu == null || !ms.uu.equalsIgnoreCase(username))
//             //         continue;


// List<network.inetwork.ISession> sessionsCopy;
// try {
//     sessionsCopy = new java.util.ArrayList<>(SessionManager.gI().getSessions());
// } catch (Exception e) {
//     sessionsCopy = new java.util.ArrayList<>();
// }
// for (ISession s : sessionsCopy) {
//     MySession ms = (MySession) s;
//     if (ms == this)
//         continue;
//     if (ms.uu == null || !ms.uu.equalsIgnoreCase(username))
//         continue;



//                 Logger.warning("[LOGIN-DUP-SESSION] user=" + username
//                         + " newIp=" + this.ipAddress
//                         + " oldIp=" + ms.ipAddress
//                         + " timeWait=" + ms.timeWait
//                         + " -> kicking, then proceeding");
//                 Service.gI().sendThongBaoOK(ms, "Mất kết nối với máy chủ!");
//                 Client.gI().kickSession(ms);
//             }

//             if (this.player != null) {
//                 Logger.warning("[LOGIN-SKIP] reason=SessionAlreadyHasPlayer user=" + username
//                         + " ip=" + this.ipAddress
//                         + " existing=" + this.player.name);
//                 return;
//             }

//             Player pl = null;
//             try {
//                 long st = System.currentTimeMillis();
//                 this.uu = username;
//                 this.pp = password;

//                 pl = NDVSqlFetcher.login(this, al);

//                 if (pl == null) {
//                     Logger.warning("[LOGIN-FAIL] reason=SqlReturnedNull user=" + username
//                             + " ip=" + this.ipAddress);
//                     return;
//                 }

//                 DataGame.sendSmallVersion(this);
//                 DataGame.sendBgItemVersion(this);

//                 this.timeWait = 0;
//                 this.joinedGame = true;

//                 pl.nPoint.calPoint();
//                 pl.nPoint.setHp(pl.nPoint.hp);
//                 pl.nPoint.setMp(pl.nPoint.mp);

//                 pl.zone.addPlayer(pl);

//                 if (pl.pet != null) {
//                     pl.pet.nPoint.calPoint();
//                     pl.pet.nPoint.setHp(pl.pet.nPoint.hp);
//                     pl.pet.nPoint.setMp(pl.pet.nPoint.mp);
//                 }

//                 pl.setSession(this);
//                 Client.gI().put(pl);
//                 this.player = pl;

//                 DataGame.sendVersionGame(this);
//                 DataGame.sendDataItemBG(this);
//                 Controller.gI().sendInfo(this);

//                 Logger.warning("[LOGIN-OK] user=" + this.player.name
//                         + " ip=" + this.ipAddress
//                         + " sqlMs=" + (System.currentTimeMillis() - st)
//                         + " totalMs=" + (System.currentTimeMillis() - startTime));

//                 if (this.player.notify != null
//                         && !this.player.notify.equals("null")
//                         && !this.player.notify.isEmpty()) {
//                     Service.gI().sendThongBao(this.player, this.player.notify);
//                     this.player.notify = null;
//                 }
//             } catch (Exception e) {
//                 Logger.warning("[LOGIN-ERR] user=" + username + " ip=" + this.ipAddress
//                         + " ex=" + e.getClass().getSimpleName() + " msg=" + e.getMessage());
//                 e.printStackTrace();
//                 if (pl != null) {
//                     pl.dispose();
//                 }
//             }
//         }
//     }

//     /** Poll until the kicked session is gone from Client, or time out. */
//     private static void waitForKickCleanup(String username) {
//         long deadline = System.currentTimeMillis() + 1500;
//         while (System.currentTimeMillis() < deadline) {
//             boolean stillThere = false;
//             for (Player pp : Client.gI().getPlayers()) {
//                 if (pp != null && pp.getSession() != null && pp.getSession().uu != null
//                         && pp.getSession().uu.equalsIgnoreCase(username)) {
//                     stillThere = true;
//                     break;
//                 }
//             }
//             if (!stillThere)
//                 return;
//             try {
//                 Thread.sleep(50);
//             } catch (InterruptedException e) {
//                 Thread.currentThread().interrupt();
//                 return;
//             }
//         }
//         Logger.warning("[LOGIN-CLEANUP-TIMEOUT] user=" + username
//                 + " — old entry didn't disappear within 1500ms, proceeding anyway");
//     }
// }
package server.io;

import java.net.Socket;

import player.Player;
import server.Controller;
import data.DataGame;
import jdbc.daos.NDVSqlFetcher;
import item.Item;
import java.io.IOException;
import network.Session;
import network.SessionManager;
import network.inetwork.ISession;
import network.Message;
import server.Client;
import server.Maintenance;
import server.Manager;
import models.AntiLogin;
import services.Service;
import utils.Logger;

import java.util.List;
import java.util.Map;
// [FIX-1] HashMap → ConcurrentHashMap vì nhiều thread gọi login() từ nhiều IP cùng lúc
import java.util.concurrent.ConcurrentHashMap;

import utils.TimeUtil;
import utils.Util;

public class MySession extends Session {

    // [FIX-1] ConcurrentHashMap: thread-safe, tránh corrupt khi nhiều IP login đồng thời
    private static final Map<String, AntiLogin> ANTILOGIN = new ConcurrentHashMap<>();

    // Không thay đổi — đã đúng, không được remove key khỏi map này
    private static final ConcurrentHashMap<String, Object> LOGIN_LOCKS = new ConcurrentHashMap<>();

    public Player player;

    public byte timeWait = 100;
    public boolean sentKey;
    public int packetCount = 0;
    public long lastPacketTime = 0;
    public static final byte[] KEYS = { 0 };
    public byte curR, curW;

    public String ipAddress;
    public boolean isAdmin;
    public int userId;
    public String uu;
    public String pp;

    public int typeClient;
    public byte zoomLevel;

    public long lastTimeLogout;
    public boolean joinedGame;

    public long lastTimeReadMessage;

    public boolean actived;

    public boolean check;

    public int goldBar;
    public long gold;
    public int eventPoint;
    public List<Item> itemsReward;
    public String dataReward;
    public boolean is_gift_box;
    public double bdPlayer;

    public int version;
    public int cash;
    public int phaobong;
    public int lixi;
    public int danap;
    public int vip;
    public int luotquay;
    public int diemdanh;
    public boolean finishUpdate;

    public long firstDataRequestTime;
    public int dataRequestCount;

    public MySession(Socket socket) {
        super(socket);
        ipAddress = socket.getInetAddress().getHostAddress();
    }

    @Override
    public void sendKey() throws Exception {
        super.sendKey();
        this.startSend();
    }

    public void sendSessionKey() {
        Message msg = new Message(-27);
        try {
            msg.writer().writeByte(KEYS.length);
            msg.writer().writeByte(KEYS[0]);
            for (int i = 1; i < KEYS.length; i++) {
                msg.writer().writeByte(KEYS[i] ^ KEYS[i - 1]);
            }
            this.sendMessage(msg);
            msg.cleanup();
            sentKey = true;
        } catch (IOException e) {
        }
    }

    public void login(String username, String password) {
        long startTime = System.currentTimeMillis();

        // ============================================================
        // [FIX-2] ROOT CAUSE FIX — đây là nguyên nhân chính gây lỗi
        // "Cannot read field uu because ms is null" sau vài tiếng chạy
        //
        // Khi session bị disconnect(), Session.dispose() được gọi và
        // set tất cả field về null (socket, sender, collector...).
        // Nhưng session đó vẫn còn nằm trong SessionManager.getSessions()
        // chưa được dọn ra. Khi player mới login, Sweep #2 duyệt qua
        // list này, cast ra MySession → ms trỏ vào object đã dispose
        // → ms.uu crash NullPointerException.
        //
        // Fix: kiểm tra isConnected() ngay đầu — session đã dispose
        // thì isConnected() = false, thoát ra luôn, không xử lý tiếp.
        // ============================================================
        if (!this.isConnected()) {
            Logger.warning("[LOGIN-SKIP] reason=SessionDisconnected user=" + username
                    + " ip=" + this.ipAddress);
            return;
        }

        Logger.warning("[LOGIN-ATTEMPT] user=" + username + " ip=" + this.ipAddress);

        // [FIX-1] computeIfAbsent thay get+put riêng lẻ — atomic, thread-safe
        AntiLogin al = ANTILOGIN.computeIfAbsent(this.ipAddress, k -> new AntiLogin());

        if (!al.canLogin()) {
            Logger.warning("[LOGIN-BLOCK] reason=AntiLogin user=" + username + " ip=" + this.ipAddress);
            Service.gI().sendThongBaoOK(this, al.getNotifyCannotLogin());
            return;
        }

        if (Manager.LOCAL) {
            Logger.warning("[LOGIN-BLOCK] reason=LOCAL user=" + username + " ip=" + this.ipAddress);
            Service.gI().sendThongBaoOK(this, "Server này chỉ để lưu dữ liệu\nVui lòng qua server khác");
            return;
        }

        if (Maintenance.isRunning) {
            Logger.warning("[LOGIN-BLOCK] reason=Maintenance user=" + username + " ip=" + this.ipAddress);
            Service.gI().sendThongBaoOK(this, "Server đang trong thời gian bảo trì, vui lòng quay lại sau");
            return;
        }

        if (!this.isAdmin && Client.gI().getPlayers().size() >= Manager.MAX_PLAYER) {
            Logger.warning("[LOGIN-BLOCK] reason=ServerFull user=" + username
                    + " players=" + Client.gI().getPlayers().size() + "/" + Manager.MAX_PLAYER);
            Service.gI().sendThongBaoOK(this, "Máy chủ hiện đang quá tải, cư dân vui lòng di chuyển sang máy chủ khác.");
            return;
        }

        String userKey = username.toLowerCase();
        Object lock = LOGIN_LOCKS.computeIfAbsent(userKey, k -> new Object());

        synchronized (lock) {

            // [FIX-2] Kiểm tra lại sau khi chờ lock — client có thể đã ngắt kết nối
            // trong lúc chờ (ví dụ user bấm thoát), tránh xử lý vô ích
            if (!this.isConnected()) {
                Logger.warning("[LOGIN-SKIP] reason=DisconnectedWhileWaitingLock user=" + username
                        + " ip=" + this.ipAddress);
                return;
            }

            // ============================================================
            // Sweep #1: kick player đang in-game với cùng username
            // [FIX-3] Copy list trước khi duyệt — tránh ConcurrentModificationException
            // khi kickSession() xóa phần tử khỏi list trong lúc đang duyệt
            // ============================================================
            List<Player> playersCopy = new java.util.ArrayList<>(Client.gI().getPlayers());
            for (Player p : playersCopy) {
                if (p == null || p.getSession() == null || p.getSession().uu == null)
                    continue;
                if (p.getSession() == this)
                    continue;
                if (!p.getSession().uu.equalsIgnoreCase(username))
                    continue;

                MySession old = (MySession) p.getSession();
                Logger.warning("[LOGIN-DUP-PLAYER] user=" + username
                        + " newIp=" + this.ipAddress
                        + " oldIp=" + old.ipAddress
                        + " -> kicking old");
                Service.gI().sendThongBaoOK(p, "Mất kết nối với máy chủ!");
                Client.gI().kickSession(old);

                // [FIX-4] Bỏ waitForKickCleanup (block tới 1500ms) — thay bằng sleep ngắn
                // waitForKickCleanup block toàn bộ lock trong 1500ms → các login khác
                // xếp hàng chờ → sau vài tiếng tích lũy đủ thread blocked → server đứng
                try { Thread.sleep(100); } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
                break;
            }

            // ============================================================
            // Sweep #2: kick session còn trong SessionManager với cùng username
            // [FIX-2+FIX-3] Copy list + null-check ms trước khi dùng
            // Đây là chỗ crash "ms is null" — session trong list đã dispose
            // ============================================================
            List<ISession> sessionsCopy;
            try {
                sessionsCopy = new java.util.ArrayList<>(SessionManager.gI().getSessions());
            } catch (Exception e) {
                sessionsCopy = new java.util.ArrayList<>();
            }
            for (ISession s : sessionsCopy) {
                // [FIX-2] NULL CHECK: s có thể là session đã dispose → cast ra MySession
                // nhưng các field bên trong đã null hết → phải check trước khi truy cập
                if (s == null) continue;
                MySession ms = (MySession) s;
                if (ms == this) continue;

                // [FIX-2] Check isConnected() thay vì chỉ check uu == null
                // Session đã dispose có thể có uu != null nhưng socket đã đóng
                // → kick sẽ crash hoặc không có tác dụng
                if (!ms.isConnected()) continue;

                if (ms.uu == null || !ms.uu.equalsIgnoreCase(username)) continue;

                Logger.warning("[LOGIN-DUP-SESSION] user=" + username
                        + " newIp=" + this.ipAddress
                        + " oldIp=" + ms.ipAddress
                        + " -> kicking");
                Service.gI().sendThongBaoOK(ms, "Mất kết nối với máy chủ!");
                Client.gI().kickSession(ms);
            }

            // Session này đã có player → thoát
            if (this.player != null) {
                Logger.warning("[LOGIN-SKIP] reason=SessionAlreadyHasPlayer user=" + username
                        + " ip=" + this.ipAddress
                        + " existing=" + this.player.name);
                return;
            }

            Player pl = null;
            try {
                long st = System.currentTimeMillis();
                this.uu = username;
                this.pp = password;

                pl = NDVSqlFetcher.login(this, al);

                if (pl == null) {
                    Logger.warning("[LOGIN-FAIL] reason=SqlReturnedNull user=" + username
                            + " ip=" + this.ipAddress);
                    return;
                }

                // [FIX-2] Kiểm tra lại sau SQL — query mất 20-200ms, trong thời gian đó
                // client có thể đã đóng kết nối → tránh tiếp tục với session chết
                if (!this.isConnected()) {
                    Logger.warning("[LOGIN-SKIP] reason=DisconnectedAfterSql user=" + username
                            + " ip=" + this.ipAddress);
                    pl.dispose();
                    return;
                }

                DataGame.sendSmallVersion(this);
                DataGame.sendBgItemVersion(this);

                this.timeWait = 0;
                this.joinedGame = true;

                pl.nPoint.calPoint();
                pl.nPoint.setHp(pl.nPoint.hp);
                pl.nPoint.setMp(pl.nPoint.mp);

                pl.zone.addPlayer(pl);

                if (pl.pet != null) {
                    pl.pet.nPoint.calPoint();
                    pl.pet.nPoint.setHp(pl.pet.nPoint.hp);
                    pl.pet.nPoint.setMp(pl.pet.nPoint.mp);
                }

                pl.setSession(this);
                if (!Client.gI().put(pl)) {
                    Logger.warning("[LOGIN-BLOCK] reason=DuplicatePlayerBanned user=" + username
                            + " ip=" + this.ipAddress);
                    return;
                }
                this.player = pl;

                DataGame.sendVersionGame(this);
                DataGame.sendDataItemBG(this);
                Controller.gI().sendInfo(this);

                Logger.warning("[LOGIN-OK] user=" + this.player.name
                        + " ip=" + this.ipAddress
                        + " sqlMs=" + (System.currentTimeMillis() - st)
                        + " totalMs=" + (System.currentTimeMillis() - startTime));

                if (this.player.notify != null
                        && !this.player.notify.equals("null")
                        && !this.player.notify.isEmpty()) {
                    Service.gI().sendThongBao(this.player, this.player.notify);
                    this.player.notify = null;
                }
            } catch (Exception e) {
                Logger.warning("[LOGIN-ERR] user=" + username + " ip=" + this.ipAddress
                        + " ex=" + e.getClass().getSimpleName() + " msg=" + e.getMessage());
                e.printStackTrace();
                if (pl != null) {
                    pl.dispose();
                }
            }
        }
    }

    // [FIX-4] Đã xóa waitForKickCleanup
    // Block thread 1500ms trong khi giữ synchronized lock → tất cả login cùng username
    // phải chờ → sau vài tiếng tích lũy đủ blocked thread → server không nhận login mới
}
