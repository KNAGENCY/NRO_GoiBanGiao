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
import utils.Util;

public class NangCapBongTai3 {

    private static final int BONG_TAI_CAP_2 = 921;
    private static final int BONG_TAI_CAP_3 = 1499;
    private static final int DA_NGU_SAC = 674;
    private static final int MANH_VO_BONG_TAI = 933;
    private static final int OPTION_MANH_VO = 31;
    private static final int REQUIRED_MANH_VO = 9999;
    private static final int FAIL_SUB_MANH_VO = 999;
    private static final int SUCCESS_RATE = 10;
    private static final long GOLD_COST = 5_000_000;
    private static final int GEM_COST = 20;

    public static void showInfoCombine(Player player) {
        if (player.combine.itemsCombine.size() != 3) {
            Service.gI().sendDialogMessage(player, "Cần 1 Bông tai cấp 2, 1 Thạch Thần và 9999 Mảnh vỡ bông tai.");
            return;
        }
        Item bongTai = null;
        Item daNguSac = null;
        Item manhVo = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case BONG_TAI_CAP_2 ->
                        bongTai = item;
                    case DA_NGU_SAC ->
                        daNguSac = item;
                    case MANH_VO_BONG_TAI ->
                        manhVo = item;
                }
            }
        }
        if (bongTai == null || daNguSac == null || manhVo == null) {
            Service.gI().sendDialogMessage(player, "Cần 1 Bông tai cấp 2, 1 Thạch Thần và 9999 Mảnh vỡ bông tai.");
            return;
        }

        int quantityManhVo = manhVo.getOptionParam(OPTION_MANH_VO);
        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_BLUE).append("Bông tai Porata [+3]\n\n");
        text.append(ConstFont.BOLD_BLUE).append("Tỉ lệ thành công: ").append(SUCCESS_RATE).append("%\n");
        text.append(quantityManhVo >= REQUIRED_MANH_VO ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED)
                .append("Cần 9999 Mảnh vỡ bông tai\n");
        text.append(daNguSac.quantity >= 1 ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED).append("Cần 1 ")
                .append(daNguSac.template.name).append("\n");
        text.append(player.inventory.gold >= GOLD_COST ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED).append("Cần ")
                .append(Util.numberToMoney(GOLD_COST)).append(" vàng\n");
        text.append(player.inventory.getGemAndRuby() >= GEM_COST ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED)
                .append("Cần ").append(GEM_COST).append(" ngọc\n");
        text.append(ConstFont.BOLD_RED).append("Thất bại mất 999 Mảnh vỡ bông tai và 1 ").append(daNguSac.template.name)
                .append("\n");

        Npc npc = getNpc(player);
        if (player.inventory.getGemAndRuby() < GEM_COST) {
            npc.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(),
                    "Còn thiếu\n" + Util.numberToMoney(GEM_COST - player.inventory.getGemAndRuby()) + " ngọc");
            return;
        }
        if (player.inventory.gold < GOLD_COST) {
            npc.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(),
                    "Còn thiếu\n" + Util.numberToMoney(GOLD_COST - player.inventory.gold) + " vàng");
            return;
        }
        if (quantityManhVo < REQUIRED_MANH_VO) {
            npc.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(),
                    "Còn thiếu\n" + (REQUIRED_MANH_VO - quantityManhVo) + " Mảnh vỡ bông tai");
            return;
        }
        npc.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(),
                "Nâng cấp\n" + Util.numberToMoney(GOLD_COST) + " vàng\n" + GEM_COST + " ngọc", "Từ chối");
    }

    public static void nangCapBongTai3(Player player) {
        if (player.combine.itemsCombine.size() != 3) {
            return;
        }
        Item bongTai = null;
        Item daNguSac = null;
        Item manhVo = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case BONG_TAI_CAP_2 ->
                        bongTai = item;
                    case DA_NGU_SAC ->
                        daNguSac = item;
                    case MANH_VO_BONG_TAI ->
                        manhVo = item;
                }
            }
        }
        if (bongTai == null || daNguSac == null || manhVo == null) {
            return;
        }
        int quantityManhVo = manhVo.getOptionParam(OPTION_MANH_VO);
        if (quantityManhVo < REQUIRED_MANH_VO || daNguSac.quantity < 1 || player.inventory.gold < GOLD_COST
                || player.inventory.getGemAndRuby() < GEM_COST) {
            return;
        }
        int bongTaiBagIndex = InventoryService.gI().getIndexBag(player, bongTai);
        if (bongTaiBagIndex < 0) {
            Service.gI().sendDialogMessage(player, "Không tìm thấy Bông tai cấp 2 trong hành trang.");
            return;
        }

        player.inventory.gold -= GOLD_COST;
        player.inventory.subGemAndRuby(GEM_COST);
        InventoryService.gI().subQuantityItemsBag(player, daNguSac, 1);
        if (Util.isTrue(SUCCESS_RATE, 100)) {
            Item btc3 = ItemService.gI().createNewItem((short) BONG_TAI_CAP_3);
            btc3.itemOptions.add(new Item.ItemOption(72, 2));
            player.inventory.itemsBag.set(bongTaiBagIndex, btc3);
            InventoryService.gI().subParamItemsBag(player, MANH_VO_BONG_TAI, OPTION_MANH_VO, REQUIRED_MANH_VO);
            CombineService.gI().sendEffectSuccessCombine(player);
            getNpc(player).npcChat(player, "chúc mừng con nhé");
        } else {
            InventoryService.gI().subParamItemsBag(player, MANH_VO_BONG_TAI, OPTION_MANH_VO, FAIL_SUB_MANH_VO);
            CombineService.gI().sendEffectFailCombine(player);
            getNpc(player).npcChat(player, "chúc mừng Admin nhé");
        }
        InventoryService.gI().sendItemBag(player);
        Service.gI().sendMoney(player);
        CombineService.gI().reOpenItemCombine(player);
    }

    private static Npc getNpc(Player player) {
        return player.iDMark.getNpcChose() != null ? player.iDMark.getNpcChose() : CombineService.gI().baHatMit;
    }
}
