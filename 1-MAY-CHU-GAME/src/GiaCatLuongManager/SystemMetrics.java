package GiaCatLuongManager;

/**
 *
 * @author SADBOI
 */
import com.sun.management.OperatingSystemMXBean;
import java.lang.management.ManagementFactory;
 import com.sun.management.OperatingSystemMXBean;
import java.lang.management.ManagementFactory;
public class SystemMetrics {
    private static String formatMemory(double used, double total) {
        return String.format("(%.1f/%.1f) GB", used, total);
    }

    private static String formatPercentage(double percentage) {
        return String.format("%.0f%%", percentage);
    }

    private static double getUsedMemoryGB() {
        OperatingSystemMXBean osBean = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
        long totalMemory = osBean.getTotalMemorySize();
        long freeMemory = osBean.getFreeMemorySize();
        long usedMemory = totalMemory - freeMemory;
        return (double) usedMemory / 1.073741824E9;
    }

    private static double getTotalMemoryGB() {
        OperatingSystemMXBean osBean = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
        long totalMemory = osBean.getTotalMemorySize();
        return (double) totalMemory / 1.073741824E9;
    }

    private static double getRAMUsagePercentage() {
        OperatingSystemMXBean osBean = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
        long totalMemory = osBean.getTotalMemorySize();
        long freeMemory = osBean.getFreeMemorySize();
        long usedMemory = totalMemory - freeMemory;
        return (double) usedMemory / (double) totalMemory * 100.0;
    }

    private static double getCPUUsagePercentage() {
        OperatingSystemMXBean osBean = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
        double cpuUsage = osBean.getCpuLoad() * 100.0;
        return cpuUsage;
    }



public static String ToString() {
    Runtime runtime = Runtime.getRuntime();

    // ===== RAM JAVA =====
    long used = runtime.totalMemory() - runtime.freeMemory();
    long max = runtime.maxMemory();

    double usedMB = used / 1024.0 / 1024.0;
    double maxMB = max / 1024.0 / 1024.0;
    double ramJavaPercent = (used * 100.0) / max;

    // ===== RAM VPS =====
    OperatingSystemMXBean osBean =
            (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();

    long totalMem = osBean.getTotalMemorySize();
    long freeMem = osBean.getFreeMemorySize();
    long usedMem = totalMem - freeMem;

    double usedMemGB = usedMem / 1024.0 / 1024.0 / 1024.0;
    double totalMemGB = totalMem / 1024.0 / 1024.0 / 1024.0;
    double ramOSPercent = (usedMem * 100.0) / totalMem;

    // ===== CPU =====
    double cpuJava = osBean.getProcessCpuLoad(); // Java
    double cpuSystem = osBean.getCpuLoad();      // VPS

    cpuJava = cpuJava < 0 ? 0 : cpuJava * 100;
    cpuSystem = cpuSystem < 0 ? 0 : cpuSystem * 100;

    int cores = Runtime.getRuntime().availableProcessors();

    // ===== FORMAT =====
    return "RAM: "
            + String.format("%.0fMB (%.0f%%)", usedMB, ramJavaPercent) + " / "
            + String.format("%.2fGB (%.0f%%)", usedMemGB, ramOSPercent) + " / "
            + String.format("%.2fGB", totalMemGB)

            + "\nCPU: "
            + String.format("%.0f%%", cpuJava) + " / "
            + String.format("%.0f%%", cpuSystem) + " / "
            + cores + " cores";
}
}
