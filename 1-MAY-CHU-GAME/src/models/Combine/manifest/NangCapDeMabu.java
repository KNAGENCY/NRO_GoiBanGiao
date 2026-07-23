package models.Combine.manifest;

import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class NangCapDeMabu {

    private static boolean isNguyenLieu674(Item item) {
        return item != null && item.isNotNullItem() && item.template.id == 674;
    }

    private static boolean isTrungMabu(Item item) {
        return item != null && item.isNotNullItem() && item.template.id == 568;
    }

    public static void showInfoCombine(Player player) {
        Item nguyenLieu674 = null;
        int trungCount = 0;

        for (Item item : player.combine.itemsCombine) {
            if (isNguyenLieu674(item)) {
                nguyenLieu674 = item;
            } else if (isTrungMabu(item)) {
                trungCount += item.quantity;
            }
        }

        if (nguyenLieu674 == null || trungCount < 2) {
            Service.gI().sendThongBaoOK(player, "Thiếu nguyên liệu");
            return;
        }
        if (player.pet == null || player.pet.typePet != 1) {
            Service.gI().sendThongBaoOK(player, "Bạn chưa có đệ Mabu");
            return;
        }
        if (player.pet.levelUpgradeMabu >= 7) {
            Service.gI().sendThongBaoOK(player, "Đệ Mabu đã đạt cấp tối đa!");
            return;
        }

        String npcSay = "Bạn có muốn nâng cấp Đệ Mabu không?\n"
                + "|0|Tỉ lệ thành CONG: 50%\n"
                + "|7|Cần 1 " + nguyenLieu674.template.name + " và 2 "
                + ItemService.gI().getTemplate(568).name;

        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE,
                npcSay, "Nâng cấp", "Từ chối");
    }

    public static void startCombine(Player player) {
        Item nguyenLieu674 = null;
        Item trungItem = null;
        int trungCount = 0;

        for (Item item : player.combine.itemsCombine) {
            if (isNguyenLieu674(item)) {
                nguyenLieu674 = item;
            } else if (isTrungMabu(item)) {
                trungItem = item;
                trungCount += item.quantity;
            }
        }

        if (nguyenLieu674 == null || trungCount < 2) {
            Service.gI().sendThongBaoOK(player, "Thiếu nguyên liệu");
            return;
        }
        if (trungItem == null || trungItem.quantity < 2) {
            Service.gI().sendThongBaoOK(player, "Không đủ Trứng Mabu");
            return;
        }
        if (player.pet == null || player.pet.typePet != 1) {
            Service.gI().sendThongBaoOK(player, "Bạn chưa có đệ Mabu");
            return;
        }
        if (player.pet.levelUpgradeMabu >= 7) {
            Service.gI().sendThongBaoOK(player, "Đệ Mabu đã đạt cấp tối đa!");
            return;
        }
        if (player.inventory.gold < 500_000_000) {
            Service.gI().sendThongBaoOK(player, "Cần 500 triệu vàng để nâng cấp");
            return;
        }

        player.inventory.gold -= 500_000_000;
        Service.gI().sendMoney(player);

        InventoryService.gI().subQuantityItemsBag(player, nguyenLieu674, 1);
        InventoryService.gI().subQuantityItemsBag(player, trungItem, 2);

        if (Util.isTrue(50, 100)) {
            player.pet.levelUpgradeMabu++;
            Service.gI().sendThongBao(player,
                    "Nâng cấp thành công! Đệ Mabu đạt cấp " + player.pet.levelUpgradeMabu);
            CombineService.gI().sendEffectSuccessCombine(player);
        } else {
            Service.gI().sendThongBao(player, "Nâng cấp thất bại!");
            CombineService.gI().sendEffectFailCombine(player);
        }

        InventoryService.gI().sendItemBag(player);
        CombineService.gI().reOpenItemCombine(player);
    }

}
