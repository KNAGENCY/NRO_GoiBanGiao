package models.Combine.manifest;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.Service;
import utils.Util;

public class GiamDinhSach {

    public static void showInfoCombine(Player player) {
        if (player.combine.itemsCombine.size() != 2) {
            Service.gI().sendDialogMessage(player, "Cần Sách Tuyệt Kỹ và bùa giám định.");
            return;
        }
        Item sachTuyetKy = null;
        Item buaGiamDinh = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isSachTuyetKy() || item.isSachTuyetKy2()) {
                sachTuyetKy = item;
            } else if (item.template.id == 1284) {
                buaGiamDinh = item;
            }
        }
        if (sachTuyetKy == null || buaGiamDinh == null) {
            Service.gI().sendDialogMessage(player, "Cần Sách Tuyệt Kỹ Và Bùa Giám Định.");
            return;
        }
        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_GREEN).append("Giám định ").append(sachTuyetKy.template.name).append(" ?\n");
        text.append(ConstFont.BOLD_BLUE).append("Bùa giám định ").append(buaGiamDinh.quantity).append("/1");
        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(), "Giám Định",
                "Refuse");
    }

    public static void giamDinhSach(Player player) {
        if (player.combine.itemsCombine.size() != 2) {
            return;
        }
        Item sachTuyetKy = null;
        Item buaGiamDinh = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isSachTuyetKy() || item.isSachTuyetKy2()) {
                sachTuyetKy = item;
            } else if (item.template.id == 1284) {
                buaGiamDinh = item;
            }
        }
        if (sachTuyetKy == null || buaGiamDinh == null) {
            return;
        }
        if (!sachTuyetKy.isHaveOption(217)) {
            Service.gI().sendServerMessage(player, "Giám Định Kiểu Chóa Gì");
            return;
        }

        for (int i = 0; i < sachTuyetKy.itemOptions.size(); i++) {
            Item.ItemOption io = sachTuyetKy.itemOptions.get(i);
            if (io.optionTemplate.id == 217) {
                sachTuyetKy.itemOptions.remove(i);

                // Dòng 1
                int[] opt1 = { 77, 103, 50, 47 };
                int id1 = opt1[Util.nextInt(opt1.length)];
                int param1;
                if (id1 == 47) {
                    param1 = Util.nextInt(100, 1000);
                } else {
                    param1 = Util.nextInt(1, 10);
                }
                sachTuyetKy.itemOptions.add(new Item.ItemOption(id1, param1));

                // Dòng 2
                int[] opt2 = { 0, 47, 6, 7 };
                int id2 = opt2[Util.nextInt(opt2.length)];
                int param2;
                if (id2 == 0 || id2 == 47) {
                    param2 = Util.nextInt(100, 1500);
                } else {
                    param2 = Util.nextInt(10000, 30000);
                }
                sachTuyetKy.itemOptions.add(new Item.ItemOption(id2, param2));

                break;
            }
        }

        CombineService.gI().sendEffectSuccessCombine(player);
        InventoryService.gI().subQuantityItemsBag(player, buaGiamDinh, 1);
        InventoryService.gI().sendItemBag(player);
        CombineService.gI().reOpenItemCombine(player);
    }
}
