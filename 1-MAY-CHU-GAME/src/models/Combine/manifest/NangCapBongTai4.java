package models.Combine.manifest;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;

public class NangCapBongTai4 {

    private static final int BONG_TAI_CAP_2 = 921;
    private static final int BONG_TAI_CAP_3 = 1499;
    private static final int BONG_TAI_CAP_4 = 569;

    public static void showInfoCombine(Player player) {
        Item bongTaiCap2 = null;
        Item bongTaiCap3 = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case BONG_TAI_CAP_2 ->
                        bongTaiCap2 = item;
                    case BONG_TAI_CAP_3 ->
                        bongTaiCap3 = item;
                }
            }
        }

        if (player.combine.itemsCombine.size() != 2 || bongTaiCap2 == null || bongTaiCap3 == null) {
            Service.gI().sendDialogMessage(player,
                    "Cần x1 Bông Tai Cấp 3 và x1 Bông Tai Cấp 2.");
            return;
        }

        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_BLUE)
                .append("Đổi x1 Bông Tai Cấp 3 và x1 Bông Tai Cấp 2 để lấy Bông Tai Cấp 4\nChỉ người Đẳng Cấp mới có Bông tai cấp 4\n");
        text.append(ConstFont.BOLD_BLUE).append("Tỉ lệ thành công là 100%");

        getNpc(player).createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(),
                "Đổi", "Từ chối");
    }

    public static void nangCapBongTai4(Player player) {
        if (player.combine.itemsCombine.size() != 2) {
            return;
        }
        Item bongTaiCap2 = null;
        Item bongTaiCap3 = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case BONG_TAI_CAP_2 ->
                        bongTaiCap2 = item;
                    case BONG_TAI_CAP_3 ->
                        bongTaiCap3 = item;
                }
            }
        }
        if (bongTaiCap2 == null || bongTaiCap3 == null) {
            return;
        }

        Item bongTaiCap4 = ItemService.gI().createNewItem((short) BONG_TAI_CAP_4);
        InventoryService.gI().subQuantityItemsBag(player, bongTaiCap3, 1);
        InventoryService.gI().subQuantityItemsBag(player, bongTaiCap2, 1);
        InventoryService.gI().addItemBag(player, bongTaiCap4);
        InventoryService.gI().sendItemBag(player);
        CombineService.gI().sendEffectSuccessCombine(player);
        getNpc(player).npcChat(player, "Chúc mừng con nhé");
        CombineService.gI().reOpenItemCombine(player);
    }

    private static Npc getNpc(Player player) {
        return player.iDMark.getNpcChose() != null ? player.iDMark.getNpcChose() : CombineService.gI().baHatMit;
    }
}
