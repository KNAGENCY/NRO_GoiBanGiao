package boss;

/*
 *
 *
 * @author GiaCatLuong
 */

import server.GameExecutor;
import server.Maintenance;
import utils.Logger;

import java.util.concurrent.TimeUnit;

public class OtherBossManager extends BossManager {

    private static OtherBossManager instance;
    private volatile boolean started = false;

    public static OtherBossManager gI() {
        if (instance == null) {
            instance = new OtherBossManager();
        }
        return instance;
    }

    /**
     * Start update bằng GameExecutor
     */
    public synchronized void start() {
        if (started) {
            Logger.warning("OtherBossManager already started – skip");
            return;
        }
        started = true;

        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            if (Maintenance.isRunning) return;

            try {
                for (int i = bosses.size() - 1; i >= 0; i--) {
                    Boss boss = bosses.get(i);
                    if (boss == null) continue;

                    try {
                        boss.update();
                    } catch (Exception e) {
                        Logger.error("OtherBoss boss error: " + e.getMessage());
                        removeBoss(boss);
                    }
                }
            } catch (Exception e) {
                Logger.logException(OtherBossManager.class, e, "OtherBoss update error");
            }

        }, 0, 150, TimeUnit.MILLISECONDS);

        Logger.success("OtherBossManager scheduled");
    }
}
