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

public class SkillSummonedManager extends BossManager {

    private static SkillSummonedManager instance;
    private volatile boolean started = false;

    public static SkillSummonedManager gI() {
        if (instance == null) {
            instance = new SkillSummonedManager();
        }
        return instance;
    }

    /**
     * Start update bằng GameExecutor
     */
    public synchronized void start() {
        if (started) {
            Logger.warning("SkillSummonedManager already started – skip");
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
                        Logger.error("SkillSummoned boss error: " + e.getMessage());
                        removeBoss(boss);
                    }
                }
            } catch (Exception e) {
                Logger.logException(SkillSummonedManager.class, e, "SkillSummoned update error");
            }

        }, 0, 150, TimeUnit.MILLISECONDS);

        Logger.success("SkillSummonedManager scheduled");
    }
}
