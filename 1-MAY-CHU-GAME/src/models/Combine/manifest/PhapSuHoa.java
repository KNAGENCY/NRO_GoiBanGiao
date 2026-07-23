
package models.Combine.manifest;

import consts.ConstNpc;
import item.Item;
import item.Item.ItemOption;
import java.util.Arrays;
import java.util.List;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.Service;
import utils.Util;

/**
 *
 * @author Administrator
 */
public class PhapSuHoa {

    public static void showInfoCombine(Player player) {
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.combine.itemsCombine.size() == 2) {
                Item daHacHoa = null;
                Item itemHacHoa = null;
                for (Item item_ : player.combine.itemsCombine) {
                    if (item_.template.id == 674) {
                        daHacHoa = item_;
                    } else if (item_.isTrangBiPSH()) {
                        itemHacHoa = item_;
                    }
                }
                if (daHacHoa == null) {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                            "Bạn còn thiếu Thạch Thần", "Đóng");
                    return;
                }
                if (itemHacHoa == null) {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                            "Bạn còn thiếu trang bị", "Đóng");
                    return;
                }
                for (ItemOption itopt : itemHacHoa.itemOptions) {
                    if (itopt.optionTemplate.id == 57) {
                        if (itopt.param >= CombineService.MAX_LEVEL_ITEM) {
                            Service.gI().sendThongBaoOK(player, "Trang bị đã đạt tới giới hạn pháp sư");
                            return;
                        }
                    }
                }
                String npcSay = "|2|Hiện tại " + itemHacHoa.template.name + "\n|0|";
                for (Item.ItemOption io : itemHacHoa.itemOptions) {
                    if (io.optionTemplate.id != 57) {
                        npcSay += io.getOptionString() + "\n";
                    }
                }
                player.combine.ratioCombine = 100;
                npcSay += "|2|Sau khi nâng cấp sẽ cộng 1 chỉ số pháp sư ngẫu nhiên \n|7|"
                        + "\n|7|Tỉ lệ thành CONG: " + player.combine.ratioCombine + "%\n"
                        + "Cần " + Util.numberToMoney(500000000) + " vàng";

                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE,
                        npcSay, "Nâng cấp\n" + Util.numberToMoney(500000000) + " vàng", "Từ chối");
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần có trang bị có thể pháp sư và Thạch Thần", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Hành trang cần ít nhất 1 chỗ trống", "Đóng");
        }
    }

    public static void startCombine(Player player) {
        if (player.combine.itemsCombine.size() != 2) {
            Service.gI().sendThongBao(player, "Thiếu nguyên liệu");
            return;
        }
        if (player.combine.itemsCombine.stream().filter(item -> item.isNotNullItem() && item.isTrangBiPSH())
                .count() != 1) {
            Service.gI().sendThongBao(player, "Thiếu trang bị pháp sư");
            return;
        }
        if (player.combine.itemsCombine.stream().filter(item -> item.isNotNullItem() && item.template.id == 674)
                .count() != 1) {
            Service.gI().sendThongBao(player, "Thiếu Thạch Thần");
            return;
        }
        if (InventoryService.gI().getCountEmptyBag(player) > 0) {
            if (player.inventory.gold < 500000000) {
                Service.gI().sendThongBao(player, "Con cần 500 triệu vàng để đổi...");
                return;
            }
            player.inventory.gold -= 500000000;
            Item daHacHoa = player.combine.itemsCombine.stream()
                    .filter(item -> item.template.id == 674).findFirst().get();
            Item trangBiHacHoa = player.combine.itemsCombine.stream()
                    .filter(Item::isTrangBiPSH).findFirst().get();

            for (ItemOption itopt : trangBiHacHoa.itemOptions) {
                if (itopt.optionTemplate.id == 57) {
                    if (itopt.param >= CombineService.MAX_LEVEL_ITEM) {
                        Service.gI().sendThongBaoOK(player, "Trang bị đã đạt tới giới hạn pháp sư");
                        return;
                    }
                }
            }

            if (Util.isTrue(player.combine.ratioCombine, 100)) {
                CombineService.gI().sendEffectSuccessCombine(player);

                ItemOption hacHoaOpt = null;
                int indexLastDefault = -1;

                for (int i = 0; i < trangBiHacHoa.itemOptions.size(); i++) {
                    ItemOption io = trangBiHacHoa.itemOptions.get(i);
                    if (io.optionTemplate.id == 57) {
                        hacHoaOpt = io;
                    } else {
                        indexLastDefault = i;
                    }
                }

                if (hacHoaOpt != null) {
                    hacHoaOpt.param += 1;
                } else {
                    ItemOption newOpt = new ItemOption(57, 1);
                    trangBiHacHoa.itemOptions.add(indexLastDefault + 1, newOpt);
                    hacHoaOpt = newOpt;
                }

                List<Integer> idOptionHacHoa = Arrays.asList(50, 77, 103, 5);
                int randomOption = idOptionHacHoa.get(Util.nextInt(0, idOptionHacHoa.size() - 1));

                int indexHacHoa = -1;
                for (int i = 0; i < trangBiHacHoa.itemOptions.size(); i++) {
                    if (trangBiHacHoa.itemOptions.get(i).optionTemplate.id == 57) {
                        indexHacHoa = i;
                        break;
                    }
                }

                boolean merged = false;
                if (indexHacHoa != -1) {
                    for (int i = indexHacHoa + 1; i < trangBiHacHoa.itemOptions.size(); i++) {
                        ItemOption io = trangBiHacHoa.itemOptions.get(i);
                        if (io.optionTemplate.id == randomOption) {
                            io.param += 2;
                            merged = true;
                            break;
                        }
                    }
                }

                if (!merged) {
                    ItemOption newRandom = new ItemOption(randomOption, 2);
                    if (indexHacHoa != -1) {
                        trangBiHacHoa.itemOptions.add(indexHacHoa + 1, newRandom);
                    } else {
                        trangBiHacHoa.itemOptions.add(indexLastDefault + 1, newRandom);
                    }
                }

                player.combine.ratioCombine = 0;
                Service.gI().sendThongBao(player, "Bạn đã nâng cấp thành công");
            } else {
                CombineService.gI().sendEffectFailCombine(player);
            }

            InventoryService.gI().subQuantityItemsBag(player, daHacHoa, 1);
            InventoryService.gI().sendItemBag(player);
            Service.gI().sendMoney(player);

            CombineService.gI().reOpenItemCombine(player);

        } else {
            Service.gI().sendThongBao(player, "Bạn phải có ít nhất 1 ô trống hành trang");
        }
    }
}
