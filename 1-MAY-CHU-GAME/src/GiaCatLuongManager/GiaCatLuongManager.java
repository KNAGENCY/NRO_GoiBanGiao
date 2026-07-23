// // /*
// //  * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
// //  * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
// //  */
// // package GiaCatLuongManager;

// // import java.util.concurrent.Executors;
// // import java.util.concurrent.ScheduledExecutorService;
// // import java.util.concurrent.TimeUnit;
// // import java.util.logging.Logger;

// // import jdbc.daos.PlayerDAO;

// // import server.Client;

// // /**
// //  *
// //  * @author Lucy An Trom
// //  */
// // public class GiaCatLuongManager {

// //     private static GiaCatLuongManager instance = null;

// //     // Static method
// //     // Static method to create instance of Singleton class
// //     public static synchronized GiaCatLuongManager getInstance() {
// //         if (instance == null) {
// //             instance = new GiaCatLuongManager();
// //         }
// //         return instance;
// //     }

// //     private ScheduledExecutorService scheduler;

// //     public void startAutoSave() {
// //         // int plcount = Client.gI().getPlayers().size();
// //         // int sscount = GiaCatLuongSessionManager.gI().getSessions().size();
// //         // if (sscount >= plcount + 10) {
// //         // Maintenance.gI().start(60);
// //         // }
// //         scheduler = Executors.newSingleThreadScheduledExecutor();

// //         scheduler.scheduleAtFixedRate(() -> {
// //             try {
// //                 handleAutoSave();
// //             } catch (Exception e) {
// //                 System.out.println("[AutoSaveManager] start autosave error: " + e.getLocalizedMessage());
// //             }
// //         }, 60, 90, TimeUnit.SECONDS);
// //     }

// //     public void handleAutoSave() {
// //         // System.out.println("[AutoSaveManager] start autosave sucessfully !!");
// //         //
// //         Client.gI().getPlayers().forEach(player -> {
// //             // long st = System.currentTimeMillis();
// //             System.out.println("Save player"+ player.name);
// //             PlayerDAO.updatePlayer(player);
// //             // Logger.success(TimeUtil.getCurrHour() + "h" + TimeUtil.getCurrMin() + "m: Tự
// //             // động lưu dữ liệu người chơi thành công! " + (System.currentTimeMillis() - st)
// //             // + "ms\n");

// //         });
// //     }

// // }
// package GiaCatLuongManager;

// import java.util.ArrayList;
// import java.util.List;
// import java.util.concurrent.Executors;
// import java.util.concurrent.ScheduledExecutorService;
// import java.util.concurrent.TimeUnit;

// import jdbc.daos.PlayerDAO;
// import player.Player;
// import server.Client;

// public class GiaCatLuongManager {

//     // [FIX-1] volatile: đảm bảo mọi thread thấy instance mới nhất
//     // Bản gốc dùng synchronized trên getInstance() — đúng nhưng chậm hơn cần thiết
//     // Double-checked locking với volatile nhanh hơn sau khi instance đã tạo xong
//     private static volatile GiaCatLuongManager instance = null;

//     public static GiaCatLuongManager getInstance() {
//         if (instance == null) {
//             synchronized (GiaCatLuongManager.class) {
//                 if (instance == null) {
//                     instance = new GiaCatLuongManager();
//                 }
//             }
//         }
//         return instance;
//     }

//     private ScheduledExecutorService scheduler;

//     public void startAutoSave() {
//         // [FIX-2] Nếu gọi startAutoSave() nhiều lần (ví dụ restart nhẹ),
//         // scheduler cũ không được shutdown → thread cũ vẫn chạy song song
//         // → cùng lúc có 2 AutoSave task → deadlock hoặc double-write DB
//         if (scheduler != null && !scheduler.isShutdown()) {
//             scheduler.shutdown();
//         }

//         // [FIX-3] newSingleThreadScheduledExecutor với tên thread rõ ràng
//         // giúp debug dễ hơn khi xem thread dump
//         scheduler = Executors.newSingleThreadScheduledExecutor(r -> {
//             Thread t = new Thread(r, "AutoSave-Thread");
//             t.setDaemon(true); // Tự tắt khi JVM shutdown, không giữ process sống
//             return t;
//         });

//         scheduler.scheduleAtFixedRate(() -> {
//             try {
//                 handleAutoSave();
//             } catch (Exception e) {
//                 System.out.println("[AutoSave] Error: " + e.getMessage());
//             }
//         }, 360, 180, TimeUnit.SECONDS);

//         System.out.println("[AutoSave] Scheduled every 90s, first run after 60s");
//     }

//     public void handleAutoSave() {
//         // [FIX-4] ROOT CAUSE của ConcurrentModificationException:
//         // Client.gI().getPlayers() trả về list gốc đang được các thread khác
//         // add/remove player liên tục (login, logout, disconnect...).
//         // forEach trên list gốc trong khi list bị modify → ConcurrentModificationException
//         // → AutoSave crash giữa chừng → một số player không được lưu
//         //
//         // Fix: copy toàn bộ list ra ArrayList riêng trước khi duyệt
//         // ArrayList copy không bị ảnh hưởng bởi thay đổi của list gốc
//         List<Player> snapshot = new ArrayList<>(Client.gI().getPlayers());

//         int saved = 0;
//         int failed = 0;

//         for (Player player : snapshot) {
//             // [FIX-5] Null check từng player — list copy có thể chứa null
//             // nếu list gốc có slot null (ví dụ player vừa disconnect đúng lúc copy)
//             if (player == null) continue;

//             // [FIX-6] Null check player.name và session — player đang trong quá trình
//             // khởi tạo hoặc đang disconnect có thể có name/session null
//             // → PlayerDAO.updatePlayer(player) sẽ crash với NPE
//             if (player.name == null) continue;
//             if (player.getSession() == null) continue;

//             try {
//                 System.out.println("[AutoSave] Saving: " + player.name);
//                 PlayerDAO.updatePlayer(player);
//                 saved++;
//             } catch (Exception e) {
//                 // [FIX-7] Catch từng player riêng lẻ — bản gốc chỉ catch ở scheduleAtFixedRate
//                 // → 1 player lỗi sẽ làm crash toàn bộ vòng lặp, các player sau không được lưu
//                 failed++;
//                 System.out.println("[AutoSave] Failed to save player: " + player.name
//                         + " | Error: " + e.getMessage());
//             }
//         }

//         System.out.println("[AutoSave] Done — saved: " + saved
//                 + ", failed: " + failed
//                 + ", total: " + snapshot.size());
//     }

//     // [NEW] Gọi khi server shutdown để dừng scheduler sạch
//     // Không có cái này thì thread AutoSave vẫn chạy sau khi server "tắt"
//     public void shutdown() {
//         if (scheduler != null && !scheduler.isShutdown()) {
//             scheduler.shutdown();
//             try {
//                 // Chờ tối đa 10s để job hiện tại hoàn thành
//                 if (!scheduler.awaitTermination(10, TimeUnit.SECONDS)) {
//                     scheduler.shutdownNow();
//                 }
//             } catch (InterruptedException e) {
//                 scheduler.shutdownNow();
//                 Thread.currentThread().interrupt();
//             }
//             System.out.println("[AutoSave] Scheduler stopped");
//         }
//     }
// }
package GiaCatLuongManager;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import jdbc.daos.PlayerDAO;
import player.Player;
import server.Client;

public class GiaCatLuongManager {

    private static volatile GiaCatLuongManager instance = null;

    public static GiaCatLuongManager getInstance() {
        if (instance == null) {
            synchronized (GiaCatLuongManager.class) {
                if (instance == null) {
                    instance = new GiaCatLuongManager();
                }
            }
        }
        return instance;
    }

    private ScheduledExecutorService scheduler;

    public void startAutoSave() {
        if (scheduler != null && !scheduler.isShutdown()) {
            scheduler.shutdown();
        }

        scheduler = Executors.newSingleThreadScheduledExecutor(r -> {
            Thread t = new Thread(r, "AutoSave-Thread");
            t.setDaemon(true);
            return t;
        });

        scheduler.scheduleAtFixedRate(() -> {
            try {
                handleAutoSave();
            } catch (Exception e) {
                System.out.println("[AutoSave] Error: " + e.getMessage());
            }
        }, 600, 300, TimeUnit.SECONDS);

        System.out.println("[AutoSave] Tu dong luu sau 10 phut, sau do moi 5 phut.");
    }

    public void handleAutoSave() {
        List<Player> snapshot = new ArrayList<>(Client.gI().getPlayers());

        int saved = 0;
        int failed = 0;

        for (Player player : snapshot) {
            if (player == null) {
                continue;
            }

            if (player.name == null) {
                continue;
            }

            if (player.getSession() == null) {
                continue;
            }

            try {
                System.out.println("[AutoSave] Saving: " + player.name);
                PlayerDAO.updatePlayer(player);
                saved++;
            } catch (Exception e) {
                failed++;
                System.out.println("[AutoSave] Failed to save player: "
                        + player.name + " | Error: " + e.getMessage());
            }
        }

        System.out.println("[AutoSave] Done - saved: " + saved
                + ", failed: " + failed
                + ", total: " + snapshot.size());
    }

    public void shutdown() {
        if (scheduler != null && !scheduler.isShutdown()) {
            scheduler.shutdown();

            try {
                if (!scheduler.awaitTermination(10, TimeUnit.SECONDS)) {
                    scheduler.shutdownNow();
                }
            } catch (InterruptedException e) {
                scheduler.shutdownNow();
                Thread.currentThread().interrupt();
            }

            System.out.println("[AutoSave] Scheduler stopped");
        }
    }
}