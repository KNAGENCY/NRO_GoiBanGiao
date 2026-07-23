package models.Combine.manifest;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.Service;
import utils.Util;

public class EpSaoTrangBi {

    public static int getGem(int star) {
        return 1; // Chỉnh về 10 ngọc đồng bộ với hiển thị
    }

    public static void showInfoCombine(Player player) {
        // Kiểm tra số lượng item: 7 lên 8 cần 3 món (Trang bị + Đá + Dùi), còn lại cần 2
        int starOfItem = -1;
        for (Item it : player.combine.itemsCombine) {
            if (it.canPhaLeHoa()) {
                starOfItem = it.getOptionParam(102);
                break;
            }
        }

        int sizeRequired = (starOfItem == 7) ? 3 : 2;

        if (player.combine.itemsCombine.size() < sizeRequired) {
            String msg = (starOfItem == 7) 
                ? "Cần thêm 1 Thạch Thần" 
                : "Cần 1 Trang bị có lỗ sao pha lê và 1 loại ngọc để ép vào.";
            Service.gI().sendDialogMessage(player, msg);
            return;
        }

        Item trangBi = null;
        Item daPhaLe = null;
        Item duiDuc = null;

        for (Item item : player.combine.itemsCombine) {
            if (item.canPhaLeHoa()) {
                trangBi = item;
            } else if (item.isDaPhaLeEpSao()) {
                daPhaLe = item;
            } else if (item.template.id == 674) {
                duiDuc = item;
            }
        }

        if (trangBi == null || daPhaLe == null) {
            Service.gI().sendDialogMessage(player, "Thiếu Trang bị hoặc Đá pha lê.");
            return;
        }

        int star = trangBi.getOptionParam(102);
        int starEmpty = trangBi.getOptionParam(107);

        if (star >= starEmpty) {
            Service.gI().sendDialogMessage(player, "Đã Max Setting");
            return;
        }

        // Kiểm tra dùi đục khi lên sao thứ 8
        if (star == 7 && duiDuc == null) {
            Service.gI().sendDialogMessage(player, "Cần thêm 1 Thạch Thần");
            return;
        }

        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_BLUE).append(trangBi.template.name).append("\n");
        text.append(ConstFont.BOLD_DARK).append(star >= 7 ? trangBi.getOptionInfoCuongHoa(daPhaLe) : trangBi.getOptionInfo(daPhaLe)).append("\n");
        
        if (star == 7) {
            text.append(ConstFont.BOLD_RED).append("Mất 1 Thạch Thần\n");
        }

        text.append(player.inventory.getGemAndRuby() < getGem(star) ? ConstFont.BOLD_RED : ConstFont.BOLD_BLUE).append("Mất ").append(getGem(star)).append(" ngọc");

        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(), 
                "Nâng cấp", "Từ chối");
    }

    public static void epSaoTrangBi(Player player) {
        Item trangBi = null;
        Item daPhaLe = null;
        Item duiDuc = null;

        for (Item item : player.combine.itemsCombine) {
            if (item.canPhaLeHoa()) {
                trangBi = item;
            } else if (item.isDaPhaLeEpSao()) {
                daPhaLe = item;
            } else if (item.template.id == 674) {
                duiDuc = item;
            }
        }

        if (trangBi == null || daPhaLe == null) return;

        int star = trangBi.getOptionParam(102);
        int starEmpty = trangBi.getOptionParam(107);

        // Chặn nếu thiếu dùi đục khi ép từ 7 lên 8
        if (star == 7 && duiDuc == null) {
            Service.gI().sendThongBao(player, "Thiếu Thạch Thần!");
            return;
        }

        if (star >= starEmpty || player.inventory.getGemAndRuby() < getGem(star)) {
            return;
        }

        // Thực hiện ép sao
        trangBi.addOptionParam(102, 1);
        
        if (star >= 7) {
            if (star == 7) {
                trangBi.itemOptions.add(new Item.ItemOption(218, 0)); // Option đánh dấu cường hóa/sao đặc biệt
            }
            trangBi.itemOptions.add(new Item.ItemOption(daPhaLe.getOptionDaPhaLe().optionTemplate.id, daPhaLe.getOptionDaPhaLe().param));
        } else {
            trangBi.addOptionParam(daPhaLe.getOptionDaPhaLe().optionTemplate.id, daPhaLe.getOptionDaPhaLe().param);
        }

        // Trừ ngọc và đá pha lê
        player.inventory.subGemAndRuby(getGem(star));
        InventoryService.gI().subQuantityItemsBag(player, daPhaLe, 1);

        // Nếu là từ 7 lên 8, trừ thêm dùi đục
        if (star == 7 && duiDuc != null) {
            InventoryService.gI().subQuantityItemsBag(player, duiDuc, 1);
            Service.gI().sendThongBao(player, "Thành Cong!");
        }

        CombineService.gI().sendEffectSuccessCombine(player);
        InventoryService.gI().sendItemBag(player);
        Service.gI().sendMoney(player);
        CombineService.gI().reOpenItemCombine(player);
    }
}