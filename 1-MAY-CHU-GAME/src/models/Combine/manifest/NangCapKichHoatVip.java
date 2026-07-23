package models.Combine.manifest;

import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class NangCapKichHoatVip {

    // check đồ thần linh
    public static boolean isDoThanLinh(Item item) {
        return item.template.id >= 555 && item.template.id <= 567;
    }

    // random khoảng
    public static int randRange(int min, int max) {
        return Util.nextInt(min, max + 1);
    }

    // random option theo % chỉ số
 public static void addRandomOption(Item item, int id, int base) {
    int min = (int) (base * 0.99);
    int max = (int) (base * 1.15);
    item.itemOptions.add(new Item.ItemOption(id, randRange(min, max)));
}

    public static void showInfoCombine(Player player) {
        if (player.combine != null && player.combine.itemsCombine != null
                && player.combine.itemsCombine.size() == 2) {

            boolean allThanLinh = true;

            for (Item item : player.combine.itemsCombine) {
                if (!isDoThanLinh(item)) {
                    allThanLinh = false;
                    break;
                }
            }

            player.combine.goldCombine = 500_000_000;

            if (allThanLinh) {
                CombineService.gI().baHatMit.createOtherMenu(player,
                        ConstNpc.MENU_START_COMBINE,
                        "Dùng 2 đồ Thần Linh\nđể nâng cấp thành đồ Hủy Diệt\n"
                        + "Đồ đầu tiên sẽ quyết định tất cả",
                        "Nâng cấp\n" + Util.numberToMoney(player.combine.goldCombine) + " vàng",
                        "Từ chối");
            } else {
                Service.gI().sendThongBaoOK(player, "Cần 2 trang bị Thần Linh");
            }

        } else {
            Service.gI().sendThongBaoOK(player, "Cần 2 trang bị Thần Linh");
        }
    }

    public static void startCombine(Player player) {

        if (player.combine.itemsCombine.size() != 2) {
            return;
        }

        int gold = player.combine.goldCombine;

        if (player.inventory.gold < gold) {
            Service.gI().sendThongBao(player,
                    "Bạn không đủ vàng, còn thiếu "
                    + Util.numberToMoney(gold - player.inventory.gold));
            return;
        }

        Item doChinh = player.combine.itemsCombine.get(0);

        if (!isDoThanLinh(doChinh)) {
            Service.gI().sendThongBao(player, "Đồ đầu tiên phải là Thần Linh");
            return;
        }

        short idHuyDiet = (short) (doChinh.template.id + 95);

        Item newItem = ItemService.gI().createNewItem(idHuyDiet);

        newItem.itemOptions.clear();

        // random option đặc biệt
        int[] specialOps = {50, 77, 103};
        int specialId = specialOps[Util.nextInt(specialOps.length)];
        int specialParam = randRange(1, 5);

        switch (idHuyDiet) {

            case 650 -> {
                addRandomOption(newItem, 47, 1600);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 42);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 651 -> {
                addRandomOption(newItem, 22, 104);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 27, 20000);
                addRandomOption(newItem, 21, 40);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 652 -> {
                addRandomOption(newItem, 47, 1700);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 42);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 653 -> {
                addRandomOption(newItem, 22, 100);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 27, 16000);
                addRandomOption(newItem, 21, 40);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 654 -> {
                addRandomOption(newItem, 47, 1800);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 42);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 655 -> {
                addRandomOption(newItem, 22, 96);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 27, 12000);
                addRandomOption(newItem, 21, 40);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 656 -> {
                addRandomOption(newItem, 14, 16);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 46);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 657 -> {
                addRandomOption(newItem, 0, 8800);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 44);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 658 -> {
                addRandomOption(newItem, 23, 96);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 28, 11000);
                addRandomOption(newItem, 21, 42);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 659 -> {
                addRandomOption(newItem, 0, 8600);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 44);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 660 -> {
                addRandomOption(newItem, 23, 100);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 28, 12000);
                addRandomOption(newItem, 21, 42);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 661 -> {
                addRandomOption(newItem, 0, 9000);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 21, 44);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }

            case 662 -> {
                addRandomOption(newItem, 23, 92);
                newItem.itemOptions.add(new Item.ItemOption(specialId, specialParam));
                addRandomOption(newItem, 28, 10000);
                addRandomOption(newItem, 21, 42);
                newItem.itemOptions.add(new Item.ItemOption(30, 1));
            }
        }

        // trừ đồ
        for (Item item : player.combine.itemsCombine) {
            InventoryService.gI().subQuantityItemsBag(player, item, 1);
        }

        // trừ vàng
        player.inventory.gold -= gold;

        // add item
        InventoryService.gI().addItemBag(player, newItem);

        CombineService.gI().sendEffectSuccessCombine(player);

        InventoryService.gI().sendItemBag(player);
        Service.gI().sendMoney(player);

        CombineService.gI().reOpenItemCombine(player);
    }
}