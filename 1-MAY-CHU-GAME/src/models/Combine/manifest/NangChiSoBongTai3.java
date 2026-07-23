package models.Combine.manifest;

import consts.ConstFont;
import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.Service;
import utils.Util;

public class NangChiSoBongTai3 {

    private static final int BONG_TAI_CAP_3 = 1499;
    private static final int MANH_HON_BONG_TAI = 934;
    private static final int DA_XANH_LAM = 935;
    private static final int REQUIRED_MANH_HON = 99;
    private static final int REQUIRED_DA_XANH_LAM = 1;
    private static final int GEM_COST = 250;
    private static final int SUCCESS_RATE = 50;

    public static void showInfoCombine(Player player) {
        if (player.combine.itemsCombine.size() != 3) {
            Service.gI().sendDialogMessage(player, "Cần 1 Bông tai cấp 2, 99 Mảnh hồn porata và 1 Đá xanh lam.");
            return;
        }
        Item bongTai = null;
        Item manhHonBongTai = null;
        Item daXanhLam = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case BONG_TAI_CAP_3 ->
                        bongTai = item;
                    case MANH_HON_BONG_TAI ->
                        manhHonBongTai = item;
                    case DA_XANH_LAM ->
                        daXanhLam = item;
                }
            }
        }

        if (bongTai == null || manhHonBongTai == null || daXanhLam == null) {
            Service.gI().sendDialogMessage(player, "Cần 1 Bông tai cấp 2, 99 Mảnh hồn porata và 1 Đá xanh lam.");
            return;
        }

        StringBuilder text = new StringBuilder();
        text.append(ConstFont.BOLD_BLUE).append("Bông tai Porata [+3]\n\n");
        text.append(ConstFont.BOLD_BLUE).append("Tỉ lệ thành công: ").append(SUCCESS_RATE).append("%\n");
        text.append(manhHonBongTai.quantity >= REQUIRED_MANH_HON ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED)
                .append("Cần 99 Mảnh hồn bông tai\n");
        text.append(daXanhLam.quantity >= REQUIRED_DA_XANH_LAM ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED)
                .append("Cần 1 Đá xanh lam\n");
        text.append(player.inventory.getGemAndRuby() >= GEM_COST ? ConstFont.BOLD_BLUE : ConstFont.BOLD_RED)
                .append("Cần ").append(GEM_COST).append(" ngọc\n");
        text.append(ConstFont.BOLD_GREEN).append("+2 Chỉ số ngẫu nhiên mạnh hơn\n");

        Npc npc = getNpc(player);
        if (player.inventory.getGemAndRuby() < GEM_COST) {
            npc.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(),
                    "Còn thiếu\n" + Util.numberToMoney(GEM_COST - player.inventory.getGemAndRuby()) + " ngọc");
            return;
        }
        if (daXanhLam.quantity < REQUIRED_DA_XANH_LAM) {
            npc.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(), "Còn thiếu\nĐá xanh lam");
            return;
        }
        if (manhHonBongTai.quantity < REQUIRED_MANH_HON) {
            npc.createOtherMenu(player, ConstNpc.IGNORE_MENU, text.toString(),
                    "Còn thiếu\n" + (REQUIRED_MANH_HON - manhHonBongTai.quantity) + " Mảnh hồn bông tai");
            return;
        }
        npc.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, text.toString(), "Nâng cấp\n250 ngọc", "Từ chối");
    }

    public static void nangChiSoBongTai3(Player player) {
        if (player.combine.itemsCombine.size() != 3) {
            return;
        }
        Item bongTai = null;
        Item manhHonBongTai = null;
        Item daXanhLam = null;
        for (Item item : player.combine.itemsCombine) {
            if (item.isNotNullItem()) {
                switch (item.template.id) {
                    case BONG_TAI_CAP_3 ->
                        bongTai = item;
                    case MANH_HON_BONG_TAI ->
                        manhHonBongTai = item;
                    case DA_XANH_LAM ->
                        daXanhLam = item;
                }
            }
        }

        if (bongTai == null || manhHonBongTai == null || daXanhLam == null
                || player.inventory.getGemAndRuby() < GEM_COST
                || daXanhLam.quantity < REQUIRED_DA_XANH_LAM
                || manhHonBongTai.quantity < REQUIRED_MANH_HON) {
            return;
        }

        player.inventory.subGemAndRuby(GEM_COST);
        Service.gI().sendMoney(player);
        if (Util.isTrue(SUCCESS_RATE, 100)) {
            int[] options = {77, 103, 50, 108, 94, 14, 101, 5};
            int option1 = options[Util.nextInt(options.length)];
            int param1 = option1 == 94 || option1 == 14 ? Util.nextInt(5, 15) : Util.nextInt(5, 20);
            int option2 = options[Util.nextInt(options.length)];
            int param2 = option2 == 94 || option2 == 14 ? Util.nextInt(5, 20) : Util.nextInt(5, 25);

            bongTai.itemOptions.clear();
            bongTai.itemOptions.add(new Item.ItemOption(option1, param1));
            bongTai.itemOptions.add(new Item.ItemOption(option2, param2));
            CombineService.gI().sendEffectSuccessCombine(player);
            getNpc(player).npcChat(player, "chúc mừng con nhé");
        } else {
            CombineService.gI().sendEffectFailCombine(player);
            getNpc(player).npcChat(player, "chúc mừng Admin nhé");
        }
        InventoryService.gI().subQuantityItemsBag(player, manhHonBongTai, REQUIRED_MANH_HON);
        InventoryService.gI().subQuantityItemsBag(player, daXanhLam, REQUIRED_DA_XANH_LAM);
        InventoryService.gI().sendItemBag(player);
        CombineService.gI().reOpenItemCombine(player);
    }

    private static Npc getNpc(Player player) {
        return player.iDMark.getNpcChose() != null ? player.iDMark.getNpcChose() : CombineService.gI().baHatMit;
    }
}
