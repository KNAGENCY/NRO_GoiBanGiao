package models.Combine.manifest;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import server.ServerNotify;
import services.InventoryService;
import services.Service;
import utils.Util;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.Properties;

public class PhaLeHoaTrangBi {

    private static final File PANEL_STAR_UPGRADE_FILE = new File(resolvePanelGameRoot(), "panel_star_upgrade_rates.properties");
    private static final Properties PANEL_STAR_UPGRADE_RATES = new Properties();
    private static long panelStarUpgradeModified = Long.MIN_VALUE;

    private static String resolvePanelGameRoot() {
        String root = System.getProperty("game.root");
        if (root == null || root.trim().isEmpty()) {
            root = System.getenv("GAME_ROOT");
        }
        if (root != null && !root.trim().isEmpty()) {
            return root.trim();
        }
        String userDir = System.getProperty("user.dir");
        if (userDir != null && new File(userDir, "data\\config\\config.properties").isFile()) {
            return new File(userDir).getAbsolutePath();
        }
        String userHome = System.getProperty("user.home");
        if (userHome != null && !userHome.trim().isEmpty()) {
            File homeRoot = new File(userHome, "Music\\zz1");
            if (new File(homeRoot, "data\\config\\config.properties").isFile()) {
                return homeRoot.getAbsolutePath();
            }
        }
        File AdminRoot = new File("C:\\Users\\Admin\\Music\\zz1");
        if (new File(AdminRoot, "data\\config\\config.properties").isFile()) {
            return AdminRoot.getAbsolutePath();
        }
        return "C:\\Users\\Administrator\\Music\\zz1";
    }

    private static synchronized void reloadPanelStarUpgradeRates() {
        long modified = PANEL_STAR_UPGRADE_FILE.isFile() ? PANEL_STAR_UPGRADE_FILE.lastModified() : -1L;
        if (modified == panelStarUpgradeModified) {
            return;
        }
        PANEL_STAR_UPGRADE_RATES.clear();
        if (modified >= 0) {
            try (FileInputStream input = new FileInputStream(PANEL_STAR_UPGRADE_FILE)) {
                PANEL_STAR_UPGRADE_RATES.load(input);
            } catch (IOException ignored) {
                PANEL_STAR_UPGRADE_RATES.clear();
            }
        }
        panelStarUpgradeModified = modified;
    }

    private static float getDefaultRatioWeight(int star) {
        return switch (star) {
            case 0 ->
                40;
            case 1 ->
                20;
            case 2 ->
                10;
            case 3 ->
                5;
            case 4 ->
                4;
            case 5 ->
                3;
            case 6 ->
                2;
            case 7 ->
                1;
            case 8 ->
                1;
            default ->
                0;
        };
    }

    private static float getDefaultRatio(int star) {
        return getDefaultRatioWeight(star) * 100F / 3000F;
    }

    private static synchronized float getRatio(int star) {
        reloadPanelStarUpgradeRates();
        float fallback = getDefaultRatio(star);
        try {
            float ratio = Float.parseFloat(PANEL_STAR_UPGRADE_RATES.getProperty(
                    "star." + (star + 1) + ".ratio", String.valueOf(fallback)));
            return ratio >= 0 && ratio <= 100 ? ratio : fallback;
        } catch (NumberFormatException e) {
            return fallback;
        }
    }

    private static String getRatioStr(int star) {
        return BigDecimal.valueOf(getRatio(star))
                .setScale(6, RoundingMode.HALF_UP)
                .stripTrailingZeros()
                .toPlainString();
    }

    private static int getDefaultGold(int star) {
        return switch (star) {
            case 0 ->
                5_000_000;
            case 1 ->
                10_000_000;
            case 2 ->
                20_000_000;
            case 3 ->
                40_000_000;
            case 4 ->
                80_000_000;
            case 5 ->
                120_000_000;
            case 6 ->
                160_000_000;
            case 7 ->
                200_000_000;
            case 8 ->
                256_000_000;
            default ->
                0;
        };
    }

    private static synchronized int getGold(int star) {
        reloadPanelStarUpgradeRates();
        int fallback = getDefaultGold(star);
        try {
            long gold = Long.parseLong(PANEL_STAR_UPGRADE_RATES.getProperty(
                    "star." + (star + 1) + ".gold", String.valueOf(fallback)));
            return gold >= 0 && gold <= Integer.MAX_VALUE ? (int) gold : fallback;
        } catch (NumberFormatException e) {
            return fallback;
        }
    }

    private static boolean isUpgradeSuccess(int star) {
        long ratio = (long) (getRatio(star) * 1_000_000);
        return Util.isTrue(ratio, 100_000_000);
    }

    private static int getGem(int star) {
        return switch (star) {
            case 0 ->
                1;
            case 1 ->
                1;
            case 2 ->
                1;
            case 3 ->
                1;
            case 4 ->
                1;
            case 5 ->
                1;
            case 6 ->
                1;
            case 7 ->
                1;
            case 8 ->
                1;
            default ->
                0;
        };
    }

    public static void showInfoCombine(Player player) {
        if (player.combine.itemsCombine.size() != 1) {
            Service.gI().sendDialogMessage(player, "Trang bị không phù hợp");
            return;
        }
        Item item = player.combine.itemsCombine.get(0);
        if (item == null || !item.isNotNullItem()) {
            return;
        }
        if (item.isHaveOption(93)) {
            Service.gI().sendDialogMessage(player, "Trang bị có hạn sử dụng, không thể thực hiện");
            return;
        }
        if (!item.canPhaLeHoa()) {
            Service.gI().sendDialogMessage(player, "Trang bị không phù hợp");
            return;
        }
        int star = item.getOptionParam(107);
        int gem = getGem(star);
        int gold = getGold(star);
        if (star >= CombineService.MAX_STAR_ITEM) {
            Service.gI().sendDialogMessage(player, "Đã đạt số pha lê tối đa");
            return;
        }
        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_BLUE).append(item.template.name).append("\n");
        text.append(ConstFont.BOLD_DARK).append(item.getOptionInfo()).append("\n");
        text.append(ConstFont.BOLD_GREEN).append(star + 1).append(" ô Sao Pha Lê\n");
        text.append(ConstFont.BOLD_BLUE).append("Tỉ lệ thành CONG: ").append(getRatioStr(star)).append("%\n");
        text.append(player.inventory.gold < gold ? ConstFont.BOLD_RED : ConstFont.BOLD_BLUE).append("Cần ")
                .append(Util.numberToMoney(gold)).append(" vàng");
        if (player.inventory.gold < gold) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(),
                    "Còn thiếu\n" + Util.numberToMoney(gold - player.inventory.gold) + " vàng");
            return;
        }
        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(),
                "Nâng cấp\n" + gem + " ngọc\nx100 lần", "Nâng cấp\n" + gem + " ngọc\nx10 lần",
                "Nâng cấp\n" + gem + " ngọc", "Từ chối");
    }

    public static void phaLeHoa(Player player, int... numm) {
        int n = 1;
        if (numm.length > 0) {
            n = numm[0];
        }
        if (!player.combine.itemsCombine.isEmpty()) {
            Item item = player.combine.itemsCombine.get(0);
            if (item == null || !item.isNotNullItem() || item.isHaveOption(93) || !item.canPhaLeHoa()) {
                return;
            }
            int star = item.getOptionParam(107);
            if (star >= CombineService.MAX_STAR_ITEM) {
                return;
            }

            // Điều kiện nâng từ 4 → 5 sao cần 50 tỷ sức mạnh
            if (star == 4) {
                if (player.nPoint.power < 20_000_000_000L) {
                    Service.gI().sendThongBaoOK(player, "Cần Đạt 20 Tỷ Sức Mạnh Để Nâng Lên 5 Sao!");
                    return;
                }
            }

            // Điều kiện nâng từ 5 → 6 sao cần item đặc biệt
            // if (star == 5) {
            //     boolean hasSpecialItem = false;
            //     for (Item it : player.inventory.itemsBody) {
            //         if (it != null && it.isNotNullItem() && (it.template.id == 1276 || it.template.id == 1277)) {
            //             hasSpecialItem = true;
            //             break;
            //         }
            //     }
            //     if (!hasSpecialItem) {
            //         Service.gI().sendThongBaoOK(player, "Cần Cải Trang Cực VIP Ở Bò Mộng!");
            //         return;
            //     }
            // }

            int gold = getGold(star);
            int gem = getGem(star);
            if (n == 1) {
                if (player.inventory.gold < gold) {
                    return;
                } else if (player.inventory.getGemAndRuby() < gem) {
                    Service.gI().sendServerMessage(player,
                            "Bạn không đủ ngọc, còn thiếu " + (gem - player.inventory.getGemAndRuby()) + " ngọc nữa");
                    return;
                }
            }
            int num = 0;
            boolean success = false;
            for (int i = 0; i < n; i++) {
                num = i + 1;
                if (player.inventory.getGemAndRuby() < gem) {
                    Service.gI().sendServerMessage(player,
                            "Sau " + i + " lần nâng cấp thất bại, bạn không đủ ngọc để tiếp tục.");
                    break;
                }
                if (player.inventory.gold < gold) {
                    Service.gI().sendServerMessage(player,
                            "Sau " + i + " lần nâng cấp thất bại, bạn không đủ vàng để tiếp tục.");
                    break;
                }
                player.inventory.gold -= gold;
                player.inventory.subGemAndRuby(gem);
                if (isUpgradeSuccess(star)) {
                    success = true;
                    break;
                }
            }
            if (success) {
                item.addOptionParam(107, 1);
                if (star > 4) {
                    ServerNotify.gI().notify("Chúc mừng " + player.name + " vừa pha lê hóa "
                            + "thành công " + item.template.name + " lên " + (star + 1) + " sao pha lê");
                }
                if (n > 1) {
                    Service.gI().sendServerMessage(player, "Thành công sau " + num + " lần nâng cấp.");
                }
                CombineService.gI().sendEffectSuccessCombine(player);
            } else {
                CombineService.gI().sendEffectFailCombine(player);
            }
            InventoryService.gI().sendItemBag(player);
            Service.gI().sendMoney(player);
            CombineService.gI().reOpenItemCombine(player);
        }
    }

}
