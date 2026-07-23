package models.Combine.manifest;

import consts.ConstNpc;
import item.Item;
import item.Item.ItemOption;
import models.Combine.CombineService;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class NangCapChanMenh {

    private static final int ID_DA_NANG_CAP = 1718;

    private static final int ID_CHAN_MENH_START = 1733;
    private static final int ID_CHAN_MENH_MAX = 1741;

    private static final int SO_DA_CAN = 10;
    private static final long VANG_CAN = 500_000_000L;

    private static boolean isNguyenLieu(Item item) {
        return item != null
                && item.isNotNullItem()
                && item.template.id == ID_DA_NANG_CAP;
    }

    private static boolean isChanMenh(Item item) {
        return item != null
                && item.isNotNullItem()
                && item.template.id >= ID_CHAN_MENH_START
                && item.template.id <= ID_CHAN_MENH_MAX;
    }

    private static float getRatioCombine(Item item) {
        int level = item.template.id - ID_CHAN_MENH_START;

        float base = 80f;

        // Mỗi cấp giảm tỉ lệ đi một nửa
        float ratio = base / (float) Math.pow(2, level);

        // Chặn tỉ lệ thấp nhất là 1%
        if (ratio < 1) {
            ratio = 1;
        }

        return ratio;
    }

    public static void showInfoCombine(Player player) {
        if (player.combine.itemsCombine.size() < 2) {
            Service.gI().sendThongBao(player, "Thiếu nguyên liệu");
            return;
        }

        Item chanMenh = null;
        Item nguyenLieu = null;

        for (Item item : player.combine.itemsCombine) {
            if (isChanMenh(item)) {
                chanMenh = item;
            } else if (isNguyenLieu(item)) {
                nguyenLieu = item;
            }
        }

        if (chanMenh == null) {
            Service.gI().sendThongBao(player, "Thiếu Chân mệnh");
            return;
        }

        if (nguyenLieu == null) {
            Service.gI().sendThongBao(player, "Thiếu Đá nâng cấp");
            return;
        }

        if (chanMenh.template.id >= ID_CHAN_MENH_MAX) {
            Service.gI().sendThongBao(player, "Chân mệnh của bạn đã đạt cấp tối đa!");
            return;
        }

        if (InventoryService.gI().getCountEmptyBag(player) <= 0) {
            Service.gI().sendThongBao(player, "Hành trang phải có ít nhất 1 ô trống");
            return;
        }

        if (nguyenLieu.quantity < SO_DA_CAN) {
            Service.gI().sendThongBaoOK(player,
                    "Không đủ " + SO_DA_CAN + " " + nguyenLieu.template.name
                    + "\nHiện có: " + nguyenLieu.quantity + "/" + SO_DA_CAN);
            return;
        }

        if (player.inventory.gold < VANG_CAN) {
            Service.gI().sendThongBaoOK(player, "Không đủ 500tr vàng");
            return;
        }

        String npcSay = "Bạn có muốn nâng cấp " + chanMenh.template.name + " không ?\n"
                + "|0|Tỉ lệ thành công là : " + getRatioCombine(chanMenh) + "%\n"
                + "|7|Cần " + SO_DA_CAN + " " + nguyenLieu.template.name + " và 500tr vàng";

        CombineService.gI().baHatMit.createOtherMenu(
                player,
                ConstNpc.MENU_START_COMBINE,
                npcSay,
                "Nâng cấp",
                "Từ chối"
        );
    }

    public static void startCombine(Player player) {
        Item chanMenh = null;
        Item nguyenLieu = null;

        for (Item item : player.combine.itemsCombine) {
            if (isChanMenh(item)) {
                chanMenh = item;
            } else if (isNguyenLieu(item)) {
                nguyenLieu = item;
            }
        }

        if (chanMenh == null) {
            Service.gI().sendThongBaoOK(player, "Thiếu Chân mệnh");
            return;
        }

        if (nguyenLieu == null) {
            Service.gI().sendThongBaoOK(player, "Thiếu Đá nâng cấp");
            return;
        }

        if (chanMenh.template.id >= ID_CHAN_MENH_MAX) {
            Service.gI().sendThongBaoOK(player, "Chân mệnh đã đạt cấp tối đa");
            return;
        }

        if (InventoryService.gI().getCountEmptyBag(player) <= 0) {
            Service.gI().sendThongBao(player, "Hành trang phải có ít nhất 1 ô trống");
            return;
        }

        if (nguyenLieu.quantity < SO_DA_CAN) {
            Service.gI().sendThongBaoOK(player,
                    "Không đủ " + SO_DA_CAN + " " + nguyenLieu.template.name
                    + "\nHiện có: " + nguyenLieu.quantity + "/" + SO_DA_CAN);
            return;
        }

        if (player.inventory.gold < VANG_CAN) {
            Service.gI().sendThongBaoOK(player, "Không đủ 500tr vàng");
            return;
        }

        // Trừ nguyên liệu và vàng
        InventoryService.gI().subQuantityItemsBag(player, nguyenLieu, SO_DA_CAN);
        player.inventory.gold -= VANG_CAN;
        Service.gI().sendMoney(player);

        if (Util.isTrue(getRatioCombine(chanMenh), 100)) {
            // Tạo Chân mệnh cấp kế tiếp
            Item chanMenhMoi = ItemService.gI().createNewItem((short) (chanMenh.template.id + 1));
            int capChanMenh = chanMenhMoi.template.id - ID_CHAN_MENH_START;

            // Chỉ số theo cấp
            chanMenhMoi.itemOptions.add(new ItemOption(50, 5 + capChanMenh));       // % Sức đánh
            chanMenhMoi.itemOptions.add(new ItemOption(77, 5 + capChanMenh * 2));   // % HP
            chanMenhMoi.itemOptions.add(new ItemOption(103, 5 + capChanMenh * 2));  // % KI
            chanMenhMoi.itemOptions.add(new ItemOption(5, 5 + capChanMenh));        // % Chí mạng

            // Thêm Chân mệnh mới
            InventoryService.gI().addItemBag(player, chanMenhMoi);

            // Xóa Chân mệnh cũ
            InventoryService.gI().subQuantityItemsBag(player, chanMenh, 1);

            Service.gI().sendThongBao(player, "|0|Nâng cấp thành công " + chanMenhMoi.template.name);
            CombineService.gI().sendEffectSuccessCombine(player);
        } else {
            Service.gI().sendThongBao(player, "|0|Nâng cấp thất bại!");
            CombineService.gI().sendEffectFailCombine(player);
        }

        InventoryService.gI().sendItemBag(player);
        CombineService.gI().reOpenItemCombine(player);
    }
}