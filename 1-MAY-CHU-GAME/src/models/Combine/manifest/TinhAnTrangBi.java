package models.Combine.manifest;

import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.Service;
import java.util.Random;
import utils.Util;

public class TinhAnTrangBi {

    private static final int NGOC_TINH_AN = 674;
    private static final long PHI_VANG = 500_000_000L;

    private static boolean isTrangBiAn(Item item) {
        return item != null && item.isNotNullItem() && item.template.type >= 0 && item.template.type <= 5;
    }

    public static void showInfoCombine(Player player) {
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.combine.itemsCombine.size() == 2) {
                Item item1 = player.combine.itemsCombine.get(0);
                Item item2 = player.combine.itemsCombine.get(1);

                Item trangBi = null;
                Item ngoc = null;

                if (isTrangBiAn(item1) && item2.template.id == NGOC_TINH_AN) {
                    trangBi = item1;
                    ngoc = item2;
                } else if (isTrangBiAn(item2) && item1.template.id == NGOC_TINH_AN) {
                    trangBi = item2;
                    ngoc = item1;
                }

                if (trangBi != null && ngoc != null) {
                    if (player.inventory.gold >= PHI_VANG) {
                        String npcSay = trangBi.template.name + "\n|2|";
                        for (Item.ItemOption io : trangBi.itemOptions) {
                            npcSay += io.getOptionString() + "\n";
                        }
                        npcSay += "|1|Ra 1 ấn bất kỳ \n"
                                + "|7|Cần 1 " + ngoc.template.name + " + 500tr vàng";
                        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE,
                                npcSay, "Lụm", "Từ chối");
                    } else {
                        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                                "Bạn không đủ 500tr vàng", "Đóng");
                    }
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                            "Vật phẩm không hợp lệ", "Đóng");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần bỏ đủ vật phẩm", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Hành trang cần ít nhất 1 ô trống", "Đóng");
        }
    }

    public static void startCombine(Player player) {
        if (InventoryService.gI().getCountEmptyBag(player) > 0 && player.combine.itemsCombine.size() == 2) {
            Item item1 = player.combine.itemsCombine.get(0);
            Item item2 = player.combine.itemsCombine.get(1);

            Item trangBi = null;
            Item ngoc = null;

            if (isTrangBiAn(item1) && item2.template.id == NGOC_TINH_AN) {
                trangBi = item1;
                ngoc = item2;
            } else if (isTrangBiAn(item2) && item1.template.id == NGOC_TINH_AN) {
                trangBi = item2;
                ngoc = item1;
            }

            if (trangBi != null && ngoc != null && player.inventory.gold >= PHI_VANG && ngoc.quantity > 0) {
                player.inventory.gold -= PHI_VANG;
                Service.gI().sendMoney(player);
                InventoryService.gI().subQuantityItemsBag(player, ngoc, 1);

                if (Util.isTrue(50, 100)) {
                    int[] optIds = { 34, 35, 36 };
                    int randomId = optIds[new Random().nextInt(optIds.length)];
                    trangBi.itemOptions.removeIf(io -> io.optionTemplate.id == 34
                            || io.optionTemplate.id == 35
                            || io.optionTemplate.id == 36);
                    trangBi.itemOptions.add(new Item.ItemOption(randomId, 1));
                    InventoryService.gI().sendItemBag(player);
                    CombineService.gI().sendEffectSuccessCombine(player);
                    Service.gI().sendThongBao(player, "Khai mở Ấn thành công!");
                } else {
                    InventoryService.gI().sendItemBag(player);
                    CombineService.gI().sendEffectFailCombine(player);
                    Service.gI().sendThongBao(player, "Khai mở Ấn thất bại!");
                }

                CombineService.gI().reOpenItemCombine(player);
            } else {
                Service.gI().sendThongBao(player, "Cần 500tr vàng và 1 Thạch thần để thực hiện");
            }
        }
    }

}
