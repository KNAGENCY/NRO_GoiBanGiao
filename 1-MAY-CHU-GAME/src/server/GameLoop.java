package server;

import map.Zone;
import utils.Logger;

public class GameLoop {

    private static final int TICK = 100; // ms

    public static void start() {
        Logger.success("GameLoop scheduled");

        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            if (!ServerManager.isRunning) return;

            long now = System.currentTimeMillis();

            try {
                if (Manager.MAPS == null) return;

                for (map.Map map : Manager.MAPS) {
                    if (map == null || map.zones == null) continue;

                    for (Zone zone : map.zones) {
                        if (zone != null) {
                            zone.update(now);
                        }
                    }
                }
            } catch (Exception e) {
                Logger.logException(GameLoop.class, e, "GameLoop error");
            }

        }, 0, TICK, java.util.concurrent.TimeUnit.MILLISECONDS);
    }
}
