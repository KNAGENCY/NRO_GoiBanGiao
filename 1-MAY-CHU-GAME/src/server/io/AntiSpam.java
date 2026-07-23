package server.io;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class AntiSpam {

    private static final Map<String, Long> bannedIPs = new ConcurrentHashMap<>();
    private static final Map<String, Integer> unknownCmdCount = new ConcurrentHashMap<>();
    private static final Map<String, Long> windowStart = new ConcurrentHashMap<>();

    private static final int LIMIT = 10;              // 10 CMD lạ
    private static final long WINDOW = 10 * 1000L;    // trong 10 giây
    private static final String RULE_PREFIX = "NRO_AntiSpam_";

    public static boolean isBanned(String ip) {
        return bannedIPs.containsKey(ip);
    }

    public static void handleUnknownCmd(String ip) {
        long now = System.currentTimeMillis();
        long start = windowStart.getOrDefault(ip, 0L);

        if (now - start > WINDOW) {
            windowStart.put(ip, now);
            unknownCmdCount.put(ip, 1);
            return;
        }

        int count = unknownCmdCount.getOrDefault(ip, 0) + 1;
        unknownCmdCount.put(ip, count);

        if (count >= LIMIT) {
            banIP(ip);
            unknownCmdCount.remove(ip);
            windowStart.remove(ip);
        }
    }

    public static synchronized void banIP(String ip) {
        if (ip == null || ip.isEmpty()) return;

        if ("127.0.0.1".equals(ip)
                || "::1".equals(ip)
                || "localhost".equalsIgnoreCase(ip)) {
            return;
        }

        if (bannedIPs.containsKey(ip)) return;

        bannedIPs.put(ip, Long.MAX_VALUE);

        String ruleName = RULE_PREFIX + ip;

        try {
            Runtime.getRuntime().exec(new String[]{
                    "cmd", "/c",
                    "netsh advfirewall firewall add rule "
                            + "name=\"" + ruleName + "\" "
                            + "dir=in "
                            + "action=block "
                            + "remoteip=" + ip
            }).waitFor();

            System.out.println("[BAN IP VINH VIEN WINDOWS] " + ip);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public static synchronized void unbanIP(String ip) {
        bannedIPs.remove(ip);

        String ruleName = RULE_PREFIX + ip;

        try {
            Runtime.getRuntime().exec(new String[]{
                    "cmd", "/c",
                    "netsh advfirewall firewall delete rule "
                            + "name=\"" + ruleName + "\""
            }).waitFor();

            System.out.println("[UNBAN IP WINDOWS] " + ip);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}