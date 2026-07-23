package models.Combine.manifest;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class TaoDaHematite {
    public static void showInfoCombine(Player player) {
    if (InventoryService.gI().getCountEmptyBag(player) == 0) {
        Service.gI().sendThongBao(player, "Cần 1 ô trống trong hành trang.");
        return;
    }

    if (player.combine.itemsCombine.size() != 2) {
        Service.gI().sendDialogMessage(player, "Cần 2 đồ khác nhau.");
        return;
    }

    Item item1 = player.combine.itemsCombine.get(0);
    Item item2 = player.combine.itemsCombine.get(1);

    if (item1 == null || item2 == null || !item1.isNotNullItem() || !item2.isNotNullItem()) {
        return;
    }

    int id1 = item1.template.id;
    int id2 = item2.template.id;

    boolean validA = (id1 >= 555 && id1 <= 567) && (id2 >= 555 && id2 <= 567);
    boolean validB = (id1 >= 650 && id1 <= 662) && (id2 >= 650 && id2 <= 662);

    if (!validA && !validB || id1 == id2) {
        Service.gI().sendDialogMessage(player, "Nguyên liệu không đẳng cấp");
        return;
    }

    StringBuilder text = new StringBuilder();
    text.append(ConstFont.BOLD_BLUE).append("Chế tạo Thạch Thần\n");

    CombineService.gI().baHatMit.createOtherMenu(
            player,
            ConstNpc.MENU_START_COMBINE,
            text.toString(),
            "Tạo",
            "Hủy"
    );
}

   public static void taoDaHematite(Player player) {
    if (InventoryService.gI().getCountEmptyBag(player) == 0) {
        return;
    }

    if (player.combine.itemsCombine.size() != 2) {
        return;
    }

    Item item1 = player.combine.itemsCombine.get(0);
    Item item2 = player.combine.itemsCombine.get(1);

    if (item1 == null || item2 == null || !item1.isNotNullItem() || !item2.isNotNullItem()) {
        return;
    }

    int id1 = item1.template.id;
    int id2 = item2.template.id;

    boolean validA = (id1 >= 555 && id1 <= 567) && (id2 >= 555 && id2 <= 567);
    boolean validB = (id1 >= 650 && id1 <= 662) && (id2 >= 650 && id2 <= 662);

    if (!validA && !validB || id1 == id2) {
        return;
    }

    int rewardQty = validA ? 1 : 2;

    CombineService.gI().baHatMit.npcChat(player, "Bư cô lô, ba cô la...");

    Item daHematite = ItemService.gI().createNewItem((short) 674);

    daHematite.quantity = rewardQty;

    daHematite.itemOptions.add(new Item.ItemOption(30, 0));

    CombineService.gI().sendEffectCombineItem(player, (byte) 7,
            (short) daHematite.template.iconID, (short) -1);

    InventoryService.gI().addItemBag(player, daHematite);

    Service.gI().sendServerMessage(player,
            "Bạn nhận được " + rewardQty + " " + daHematite.template.name);

    player.inventory.gold -= 1_000_000;

    InventoryService.gI().subQuantityItemsBag(player, item1, 1);
    InventoryService.gI().subQuantityItemsBag(player, item2, 1);

    InventoryService.gI().sendItemBag(player);
    Service.gI().sendMoney(player);

    CombineService.gI().reOpenItemCombine(player);
}
}