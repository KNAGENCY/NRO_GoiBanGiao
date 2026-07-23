package models.Combine.manifest;

import player.Player;
import consts.ConstNpc;
import item.Item;
import item.Item.ItemOption;
import models.Combine.CombineService;
import services.InventoryService;
import services.Service;
import utils.Util;

public class RemoveOptionItem {

    public static void showInfoCombine(Player player) {
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.combine.itemsCombine.size() == 2) {
                Item daHacHoa = null;
                Item itemTarget = null;
                for (Item item_ : player.combine.itemsCombine) {
                    if (item_.template.id == 674) {
                        daHacHoa = item_;
                    } else if (item_.isTrangBiHacHoa()) {
                        itemTarget = item_;
                    }
                }
                if (daHacHoa == null) {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                            "Bạn còn thiếu Đá giải Hắc Hỏa", "Đóng");
                    return;
                }
                if (itemTarget == null) {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                            "Cần trang bị hợp lệ để tẩy (trang bị có pháp sư hóa)", "Đóng");
                    return;
                }

                String npcSay = "|2|Hiện tại " + itemTarget.template.name + "\n|0|";
                for (Item.ItemOption io : itemTarget.itemOptions) {
                    npcSay += io.getOptionString() + "\n";
                }

                npcSay += ("|2|Sau khi tẩy sẽ xoá toàn bộ \n|7|")
                        + "\n|7|Tỉ lệ thành công: 100%\n"
                        + "Cần 1 Thạch Thần\n"
                        + "Cần " + Util.numberToMoney(500_000_000) + " vàng";

                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE,
                        npcSay, "Tẩy nhá", "Từ chối");
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần có trang bị và Đá giải Hắc Hỏa (674)", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Hành trang cần ít nhất 1 chỗ trống", "Đóng");
        }
    }

    public static void startCombine(Player player) {
        if (player.combine.itemsCombine.size() != 2) {
            return;
        }
        if (player.combine.itemsCombine.stream().filter(item -> item.isNotNullItem() && item.isTrangBiHacHoa())
                .count() != 1) {
            return;
        }
        if (player.combine.itemsCombine.stream().filter(item -> item.isNotNullItem() && item.template.id == 674)
                .count() != 1) {
            return;
        }
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.inventory.gold < 500_000_000) {
                Service.gI().sendThongBao(player, "Con cần 500tr vàng để tẩy...");
                return;
            }

            Item daHacHoa = player.combine.itemsCombine.stream().filter(item -> item.template.id == 674)
                    .findFirst().get();
            Item trangBi = player.combine.itemsCombine.stream().filter(Item::isTrangBiHacHoa).findFirst().get();

            if (daHacHoa == null || trangBi == null) {
                Service.gI().sendThongBao(player, "Thiếu nguyên liệu hoặc trang bị cần tẩy");
                return;
            }

            if (Util.isTrue(100, 100)) {
                CombineService.gI().sendEffectSuccessCombine(player);

                int index57 = -1;
                for (int i = 0; i < trangBi.itemOptions.size(); i++) {
                    if (trangBi.itemOptions.get(i).optionTemplate.id == 57) {
                        index57 = i;
                        break;
                    }
                }

                if (index57 != -1) {
                    while (trangBi.itemOptions.size() > index57) {
                        trangBi.itemOptions.remove(trangBi.itemOptions.size() - 1);
                    }
                    Service.gI().sendThongBao(player, "Bạn đã tẩy thành công");
                } else {
                    Service.gI().sendThongBao(player, "k có cái nịt");
                }

                player.combine.ratioCombine = 0;
                InventoryService.gI().sendItemBag(player);
            } else {
                CombineService.gI().sendEffectFailCombine(player);
            }

            InventoryService.gI().subQuantityItemsBag(player, daHacHoa, 1);
            player.inventory.gold -= 500_000_000;
            InventoryService.gI().sendItemBag(player);
            Service.gI().sendMoney(player);
            player.combine.itemsCombine.clear();
            CombineService.gI().reOpenItemCombine(player);
        }
    }
}
