package server;

import java.util.concurrent.ScheduledThreadPoolExecutor;

public class GameExecutor {

    public static final ScheduledThreadPoolExecutor EXECUTOR =
            new ScheduledThreadPoolExecutor(
                    12, // 10–16 thread là đẹp
                    r -> {
                        Thread t = new Thread(r, "GameScheduler");
                        t.setDaemon(true);
                        return t;
                    }
            );

    static {
        EXECUTOR.setRemoveOnCancelPolicy(true);
    }
}
