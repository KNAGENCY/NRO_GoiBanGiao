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

public class BrolyManager extends BossManager {

    private static BrolyManager instance;
    private volatile boolean started = false;

    public static BrolyManager gI() {
        if (instance == null) {
            instance = new BrolyManager();
        }
        return instance;
    }

    /**
     * Start Broly boss update
     */
    public synchronized void start() {
        if (started) {
            Logger.warning("BrolyManager already started – skip");
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
                        Logger.error("Broly boss error: " + e.getMessage());
                        removeBoss(boss);
                    }
                }
            } catch (Exception e) {
                Logger.logException(
                        BrolyManager.class,
                        e,
                        "BrolyManager update error"
                );
            }

        }, 0, 150, TimeUnit.MILLISECONDS);

        Logger.success("BrolyManager scheduled");
    }
}
