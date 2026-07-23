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

import java.util.Arrays;
import java.util.List;

public class NangCapVatPham {

    private static int getGold(Item item) {
        int levelItem = item.getOptionParam(72);
        int levelOption = item.template.level;
        int level = levelItem + levelOption;
        if (levelOption >= 7) {
            return 10_000_000;
        } else {
            return level * 10000;
        }
    }

    private static int getDa(Item item) {
        int levelItem = item.getOptionParam(72);
        int levelOption = item.template.level;
        int level = levelItem + levelOption;
        return level + 1;
    }

    private static int getRatio(int level) {
        return switch (level) {
            case 0 -> 80;
            case 1 -> 40;
            case 2 -> 20;
            case 3 -> 10;
            case 4 -> 5;
            case 5 -> 4;
            case 6 -> 3;
            case 7 -> 2;
            default -> 0;
        };
    }

    public static void showInfoCombine(Player player) {
        if (player.combine.itemsCombine.size() != 2) {
            Service.gI().sendDialogMessage(player, "Cần 1 trang bị và đúng loại đá nâng cấp");
            return;
        }

        Item trangBi = null;
        Item daNangCap = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.template.type < 5) {
                trangBi = item;
            } else if (item.isDaNangCap()) {
                daNangCap = item;
            }
        }
        
        if (trangBi == null || daNangCap == null || !trangBi.canNangCapWithNDC(daNangCap)) {
            Service.gI().sendDialogMessage(player, "Cần 1 trang bị và đúng loại đá nâng cấp");
            return;
        }

        Item daBaoVe = InventoryService.gI().findItemBag(player, 987);
        Item daBaoVeKhoa = InventoryService.gI().findItemBag(player, 1143);
        int level = trangBi.getOptionParam(72);
        int gold = getGold(trangBi);
        int da = getDa(trangBi);

        if (level >= CombineService.MAX_LEVEL_ITEM) {
            Service.gI().sendDialogMessage(player, "Vật phẩm đã đạt cấp độ tối đa, không thể nâng cấp nữa");
            return;
        }

        boolean canUseDBV = level == 2 || level == 4 || level == 6 || level == 7;

        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_BLUE).append("Hiện tại ").append(trangBi.template.name);
        if (level > 0) {
            text.append(" [+").append(level).append("]\n");
        } else {
            text.append("\n");
        }
        text.append(ConstFont.BOLD_DARK).append(trangBi.getOptionInfoUpgrade()).append("\n");
        text.append(ConstFont.BOLD_BLUE).append("Sau khi nâng cấp [+").append(level + 1).append("]\n");
        text.append(ConstFont.BOLD_GREEN).append(trangBi.getOptionInfoUpgradeFinal()).append("\n");
        
        // Hiển thị thông tin bonus nếu sắp lên +3 hoặc +5
        if (level + 1 == 3 || level + 1 == 5) {
            text.append(ConstFont.BOLD_RED).append("LƯU Ý: Có thể nhận thêm 1 dòng chỉ số ngẫu nhiên!\n");
        }

        text.append(ConstFont.BOLD_BLUE).append("Tỉ lệ thành công: ").append(getRatio(level)).append("%\n");
        text.append(daNangCap.quantity >= da ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED).append("Cần ").append(da).append(" ").append(daNangCap.template.name).append("\n");
        text.append(player.inventory.gold >= gold ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED).append("Cần ").append(Util.numberToMoney(gold)).append(" vàng");

        if (canUseDBV) {
            text.append("\n");
            text.append(ConstFont.BOLD_BLUE).append("Nếu thất bại sẽ rớt xuống [+").append(level - 1).append("]\n");
            text.append(ConstFont.BOLD_RED).append("Nếu dùng đá bảo vệ sẽ không bị rớt cấp.");
        }

        if (daNangCap.quantity < da) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    text.toString(), "Còn thiếu\n" + Util.numberToMoney(da - daNangCap.quantity) + " " + daNangCap.template.name);
            return;
        }

        if (player.inventory.gold < gold) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    text.toString(), "Còn thiếu\n" + Util.numberToMoney(gold - player.inventory.gold) + " vàng");
            return;
        }

        if (canUseDBV && (daBaoVe != null || daBaoVeKhoa != null)) {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE,
                    text.toString(), "Nâng cấp", "Nâng cấp\ndùng đá\nbảo vệ", "Đóng");
            return;
        }

        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE,
                text.toString(), "Nâng cấp", "Đóng");
    }

    public static void nangCapVatPham(Player player, boolean useDBV) {
        if (player.combine.itemsCombine.size() != 2) {
            return;
        }

        Item trangBi = null;
        Item daNangCap = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.template.type < 5) {
                trangBi = item;
            } else if (item.isDaNangCap()) {
                daNangCap = item;
            }
        }

        if (trangBi == null || daNangCap == null || !trangBi.canNangCapWithNDC(daNangCap)) {
            return;
        }

        Item daBaoVe = InventoryService.gI().findItemBag(player, 987);
        Item daBaoVeKhoa = InventoryService.gI().findItemBag(player, 1143);
        int level = trangBi.getOptionParam(72);
        int gold = getGold(trangBi);
        int da = getDa(trangBi);

        boolean canUseDBV = level == 2 || level == 4 || level == 6 || level == 7;
        
        if (daNangCap.quantity < da || player.inventory.gold < gold || level >= CombineService.MAX_LEVEL_ITEM) {
            return;
        }
        
        if (canUseDBV && useDBV && daBaoVe == null && daBaoVeKhoa == null) {
            return;
        }

        // --- XỬ LÝ KHI THÀNH CÔNG ---
        if (Util.isTrue(getRatio(level), 100)) {
            // 1. Tăng 10% các chỉ số cũ có thể nâng cấp
            for (Item.ItemOption io : trangBi.itemOptions) {
                if (io.isOptionCanUpgrade()) {
                    io.param += (io.param * 10 / 100);
                }
            }
            
            // 2. Tăng cấp độ vật phẩm (Option 72)
            trangBi.addOptionParam(72, 1);
            int newLevel = trangBi.getOptionParam(72);

            // 3. Logic tặng thêm dòng chỉ số khi lên +3 hoặc +5
            if (newLevel == 3 || newLevel == 5) {
                int[] listOptionId = {50, 77, 103, 5};
                int randomOptionId = listOptionId[Util.nextInt(0, listOptionId.length - 1)];
                int randomParam = Util.nextInt(1, 5);
                
                // Thêm dòng option mới
                trangBi.itemOptions.add(new Item.ItemOption(randomOptionId, randomParam));
            }

            CombineService.gI().sendEffectSuccessCombine(player);
            
            if (newLevel > 1) {
                Service.gI().chatJustForMe(player, player, player.name + ": Vừa nâng cấp thành công " + trangBi.template.name + " lên +" + newLevel);
            }
            if (newLevel > 5) {
                ServerNotify.gI().notify(player.name + ": Vừa nâng cấp thành công " + trangBi.template.name + " lên +" + newLevel);
            }
           
        } else {
            // --- XỬ LÝ KHI THẤT BẠI ---
            if (canUseDBV) {
                if (!useDBV || (daBaoVe == null && daBaoVeKhoa == null)) {
                    if (trangBi.template.type < 5) {
                        for (Item.ItemOption io : trangBi.itemOptions) {
                            if (io.isOptionCanUpgrade()) {
                                io.param -= (io.param * 11 / 100);
                            }
                        }
                    }
                    trangBi.subOptionParam(72, 1);
                    trangBi.addOptionParam(209, 1); // Option ghi chú thất bại nếu có
                }
            }
            CombineService.gI().sendEffectFailCombine(player);
        }

        // --- TRỪ CHI PHÍ ---
        if (canUseDBV && useDBV) {
            if (daBaoVe != null) {
                InventoryService.gI().subQuantityItemsBag(player, daBaoVe, 1);
            } else if (daBaoVeKhoa != null) {
                InventoryService.gI().subQuantityItemsBag(player, daBaoVeKhoa, 1);
            }
        }
        
        player.inventory.gold -= gold;
        InventoryService.gI().subQuantityItemsBag(player, daNangCap, da);
        InventoryService.gI().sendItemBag(player);
        Service.gI().sendMoney(player);
        CombineService.gI().reOpenItemCombine(player);
    }
}