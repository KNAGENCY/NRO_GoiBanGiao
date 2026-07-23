package npc.npc_manifest;

import consts.ConstNpc;
import item.Item;
import npc.Npc;
import player.NPoint;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.OpenPowerService;
import services.Service;
import services.func.ChangeMapService;
import utils.Util;
import item.Item;

public class QuocVuong extends Npc {

    private static final int CONFIRM_REBORN = 999;
private static final int MENU_MAP_168 = 998;

    public QuocVuong(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public boolean canOpenNpc(Player player) {
        if (player == null || player.isDie()) {
            return false;
        }
        if (player.nPoint.power < 17_000_000_000L) {
            Service.gI().sendThongBaoOK(player, "Lên 17 Tỷ Sức Mạnh Rồi Nói Chuyện");
            return false;
        }
        return super.canOpenNpc(player);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) {
            return;
        }
    if (this.mapId == 168) {
        this.createOtherMenu(player, MENU_MAP_168,
                "Xin chào, Ta Đợi Con Lâu Lắm Rồi\nĐây Là Phần Quà Cho Con\nChỉ Được Chọn 1 Trong 2",
                "Nhận 1\nĐồ Thần",
                "Nhận\n200 TV",
                "Go Out");
        return;
    }
        if (this.mapId == 43) { this.createOtherMenu(player, ConstNpc.BASE_MENU,
                "Mày Muốn 'Phá Vỡ Giới Hạn' Hay 'Tái sinh'?\nTái sinh sẽ nhận được Thạch Thần",
                "Nâng Win\nBản Thân",
                "Nâng Win\nĐệ Tử",
                "Tái sinh",
                "Đóng");
    }
    }
  @Override
public void confirmMenu(Player player, int select) {
    if (!canOpenNpc(player)) {
        return;
    }

    // ================= MAP 172 =================
   if (player.iDMark.getIndexMenu() == MENU_MAP_168) {
    if (select == 0 || select == 1) {
        Service.gI().sendThongBaoOK(player,
                "Có cái nịt\nĐồ sát người chơi để nhận vàng!");
        return;
    }

   if (select == 2) { // Quay về map 5

    boolean empty = true;
    for (int i = 0; i < player.inventory.itemsBody.size(); i++) {
        Item item = player.inventory.itemsBody.get(i);
        if (item != null && item.isNotNullItem()) {
            empty = false;
            break;
        }
    }

    if (!empty) {
        Service.gI().sendThongBao(player,
                "Vui lòng tháo toàn bộ trang bị trước khi di chuyển!");
        return;
    }

    ChangeMapService.gI().changeMap(player, 5, -1, 1070, 408);
    return;
}
   }

        // ================= BASE MENU =================
        if (player.iDMark.isBaseMenu()) {
            switch (select) {
                case 0 -> {
                    if (player.nPoint.limitPower < NPoint.MAX_LIMIT) {
                        this.createOtherMenu(player, ConstNpc.OPEN_POWER_MYSEFT,
                                "Ta sẽ truyền năng lượng giúp con phá bỏ giới hạn của bản thân lên "
                                        + Util.numberToMoney(player.nPoint.getPowerNextLimit()),
                                "Update\n"
                                        + Util.numberToMoney(OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER)
                                        + " vàng",
                                "Đóng");
                        return; // ⭐ CHỐT LẠI TẠI ĐÂY
                    } else {
                        this.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                                "Đã Max Giới Hạn", "Đóng");
                        return;
                    }
                }
                case 1 -> {
                    if (player.pet != null) {
                        if (player.pet.nPoint.limitPower < NPoint.MAX_LIMIT) {
                            this.createOtherMenu(player, ConstNpc.OPEN_POWER_PET,
                                    "Ta sẽ truyền năng lượng giúp con phá bỏ giới hạn của đệ tử lên "
                                            + Util.numberToMoney(player.pet.nPoint.getPowerNextLimit()),
                                    "Update\n"
                                            + Util.numberToMoney(OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER)
                                            + " vàng",
                                    "Đóng");
                            return;
                        } else {
                            this.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                                    "Đã Max Giới Hạn", "Đóng");
                            return;
                        }
                    } else {
                        Service.gI().sendThongBao(player, "Không có đệ tử");
                        return;
                    }
                }
                case 2 -> {
                    if (player.nPoint.power < 19_000_000_000L) {
                        Service.gI().sendThongBaoOK(player,
                                "Cần tối thiểu 19 Tỷ Sức Mạnh để Tái sinh");
                        return;
                    }
                    this.createOtherMenu(player, CONFIRM_REBORN,
                            "Tái sinh về 150 triệu sức mạnh\n"
                                    + "Yêu cầu trống ít nhất 13 ô hành trang",
                            "Xác nhận",
                            "Hủy");
                    return;
                }
            }
        }

        // ================= CONFIRM MENU =================
        switch (player.iDMark.getIndexMenu()) {
            case ConstNpc.OPEN_POWER_MYSEFT -> {
                if (select == 0) {
                    if (player.inventory.gold >= OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER) {
                        if (OpenPowerService.gI().openPowerSpeed(player)) {
                            player.inventory.gold -= OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER;
                            Service.gI().sendMoney(player);
                        }
                    } else {
                        Service.gI().sendThongBao(player,
                                "Không đủ vàng, còn thiếu "
                                        + Util.numberToMoney(
                                                OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER
                                                        - player.inventory.gold));
                    }
                }
            }

            case ConstNpc.OPEN_POWER_PET -> {
                if (select == 0) {
                    if (player.inventory.gold >= OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER) {
                        if (OpenPowerService.gI().openPowerSpeed(player.pet)) {
                            player.inventory.gold -= OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER;
                            Service.gI().sendMoney(player);
                        }
                    } else {
                        Service.gI().sendThongBao(player,
                                "Không đủ vàng, còn thiếu "
                                        + Util.numberToMoney(
                                                OpenPowerService.COST_SPEED_OPEN_LIMIT_POWER
                                                        - player.inventory.gold));
                    }
                }
            }

            case CONFIRM_REBORN -> {
                if (select == 0) {
                    player.nPoint.power = 150_000_000;
                    player.nPoint.calPoint();
                    Service.gI().point(player);

                    for (int i = 0; i < player.inventory.itemsBody.size(); i++) {
                        Item item = player.inventory.itemsBody.get(i);
                        if (item != null && item.isNotNullItem()) {
                            InventoryService.gI().addItemBag(player, item);
                            player.inventory.itemsBody.set(i, ItemService.gI().createItemNull());
                        }
                    }

                    InventoryService.gI().sendItemBody(player);
                    InventoryService.gI().sendItemBag(player);

                    int soDa = Util.nextInt(1, 3);
                    for (int i = 0; i < soDa; i++) {
                        Item da = ItemService.gI().createNewItem((short) 674, 1);
                        InventoryService.gI().addItemBag(player, da);
                    }
                    InventoryService.gI().sendItemBag(player);

                    Service.gI().sendThongBaoOK(player, "Tái sinh thành công!");
                }
            }
        }
    }
}
