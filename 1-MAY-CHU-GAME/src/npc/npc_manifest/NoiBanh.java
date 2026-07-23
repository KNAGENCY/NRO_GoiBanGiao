package npc.npc_manifest;

/**
 * @author GiaCatLuong
 */
import consts.ConstItem;
import consts.ConstNpc;

import item.Item;

import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;

public class NoiBanh extends Npc {

    public NoiBanh(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Xin chào " + player.name + "\n1. Thịt Heo Úp ở mọi Map (yêu cầu trên 5tỷ Sức Mạnh)\n2. Thúng Nếp Úp ở mọi Map (yêu cầu trên 10tỷ Sức Mạnh)\n3. Thùng Đậu Xanh Úp ở Tương Lai\n4. Lá Dong Úp ở bên Cold\nNgoài ra pem Boss Nhiệm vụ cũng rơi rất nhiều vật phẩm",
                    "Làm 1\nCái Bánh", "Refuse");
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.iDMark.isBaseMenu()) {
                switch (select) {
                    case 0:
                        this.createOtherMenu(player, ConstNpc.MENU_NAU_BANH,
                                "|0|Bạn Làm Chóa Gì Đã Có Người Yêu Đúng Không\n",
                                "Làm Bánh\nTest", "Làm Bánh\nChưng", "Refuse");
                        break;
                }
            } else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_NAU_BANH) {
                switch (select) {
                    case 0: {
                        Item thitheo = InventoryService.gI().findItemBag(player, ConstItem.THIT_HEO);
                        Item thungNep = InventoryService.gI().findItemBag(player, ConstItem.THUNG_NEP);
                        Item thungDauXanh = InventoryService.gI().findItemBag(player, ConstItem.THUNG_DAU_XANH);
                        Item laDong = InventoryService.gI().findItemBag(player, ConstItem.LA_DONG);

                        int soLuongThitHeo = (thitheo != null) ? thitheo.quantity : 0;
                        int soLuongThungNep = (thungNep != null) ? thungNep.quantity : 0;
                        int soLuongThungDauXanh = (thungDauXanh != null) ? thungDauXanh.quantity : 0;
                        int soLuongLaDong = (laDong != null) ? laDong.quantity : 0;
                        long gold = player.inventory.gold;
                        if (soLuongThitHeo >= 10 && soLuongThungNep >= 10
                                && soLuongThungDauXanh >= 10 && soLuongLaDong >= 10 && gold >= 5000000) {
                            this.createOtherMenu(player, ConstNpc.MENU_NAU_BANH_1,
                                    "|0|Bạn Muốn Làm Bánh Test?\n"
                                            + "Thịt heo: " + soLuongThitHeo + "/10\n"
                                            + "Thúng nếp: " + soLuongThungNep + "/10\n"
                                            + "Thúng đậu xanh: " + soLuongThungDauXanh + "/10\n"
                                            + "Lá dong: " + soLuongLaDong + "/10\n"
                                            + "Gold: " + gold + "/5,000,000",
                                    "Làm Bánh\nTest", "Refuse");
                        } else {
                            this.createOtherMenu(player, ConstNpc.MENU_NAU_BANH_TU_CHOI_1,
                                    "|0|Bạn Không Đủ nguyên Liệu Để Làm Bánh Test!\n"
                                            + "Thịt heo: " + soLuongThitHeo + "/10\n"
                                            + "Thúng nếp: " + soLuongThungNep + "/10\n"
                                            + "Thúng đậu xanh: " + soLuongThungDauXanh + "/10\n"
                                            + "Lá dong: " + soLuongLaDong + "/10\n"
                                            + "Gold: " + gold + "/5,000,000",
                                    "Refuse");
                        }
                        break;
                    }
                    case 1: {
                        Item thitheo = InventoryService.gI().findItemBag(player, ConstItem.THIT_HEO);
                        Item thungNep = InventoryService.gI().findItemBag(player, ConstItem.THUNG_NEP);
                        Item thungDauXanh = InventoryService.gI().findItemBag(player, ConstItem.THUNG_DAU_XANH);
                        Item laDong = InventoryService.gI().findItemBag(player, ConstItem.LA_DONG);
                        Item trungVitMuoi = InventoryService.gI().findItemBag(player, ConstItem.TRUNG_VIT_MUOI);

                        int soLuongThitHeo = (thitheo != null) ? thitheo.quantity : 0;
                        int soLuongThungNep = (thungNep != null) ? thungNep.quantity : 0;
                        int soLuongThungDauXanh = (thungDauXanh != null) ? thungDauXanh.quantity : 0;
                        int soLuongLaDong = (laDong != null) ? laDong.quantity : 0;
                        int soLuongTrungVitMuoi = (trungVitMuoi != null) ? trungVitMuoi.quantity : 0;
                        long gold = player.inventory.gold;

                        if (soLuongThitHeo >= 10 && soLuongThungNep >= 10
                                && soLuongThungDauXanh >= 10 && soLuongLaDong >= 10
                                && soLuongTrungVitMuoi >= 1 && gold >= 5000000) {
                            this.createOtherMenu(player, ConstNpc.MENU_NAU_BANH_2,
                                    "|0|Bạn Muốn Làm Bánh Chưng ?\n"
                                            + "Thịt heo: " + soLuongThitHeo + "/10\n"
                                            + "Thúng nếp: " + soLuongThungNep + "/10\n"
                                            + "Thúng đậu xanh: " + soLuongThungDauXanh + "/10\n"
                                            + "Lá dong: " + soLuongLaDong + "/10\n"
                                            + "Trứng vịt muối: " + soLuongTrungVitMuoi + "/1\n"
                                            + "Gold: " + gold + "/5,000,000",
                                    "Làm Bánh\nChưng", "Refuse");
                        } else {
                            this.createOtherMenu(player, ConstNpc.MENU_NAU_BANH_TU_CHOI_2,
                                    "|0|Bạn Không Đủ nguyên Liệu Để Làm Bánh Chưng!\n"
                                            + "Thịt heo: " + soLuongThitHeo + "/10\n"
                                            + "Thúng nếp: " + soLuongThungNep + "/10\n"
                                            + "Thúng đậu xanh: " + soLuongThungDauXanh + "/10\n"
                                            + "Lá dong: " + soLuongLaDong + "/10\n"
                                            + "Trứng vịt muối: " + soLuongTrungVitMuoi + "/1\n"
                                            + "Gold: " + gold + "/5,000,000",
                                    "Refuse");
                        }
                        break;
                    }
                }
            } else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_NAU_BANH_1) {
                switch (select) {
                    case 0: {
                        Item thitheo = InventoryService.gI().findItemBag(player, ConstItem.THIT_HEO);
                        Item thungNep = InventoryService.gI().findItemBag(player, ConstItem.THUNG_NEP);
                        Item thungDauXanh = InventoryService.gI().findItemBag(player, ConstItem.THUNG_DAU_XANH);
                        Item laDong = InventoryService.gI().findItemBag(player, ConstItem.LA_DONG);

                        if (thitheo != null && thungNep != null && thungDauXanh != null && laDong != null
                                && thitheo.quantity >= 10 && thungNep.quantity >= 10
                                && thungDauXanh.quantity >= 10 && laDong.quantity >= 10
                                && player.inventory.gold >= 5000000) {
                            InventoryService.gI().subQuantityItemsBag(player, thitheo, 10);
                            InventoryService.gI().subQuantityItemsBag(player, thungNep, 10);
                            InventoryService.gI().subQuantityItemsBag(player, thungDauXanh, 10);
                            InventoryService.gI().subQuantityItemsBag(player, laDong, 10);
                            player.inventory.gold -= 5000000;
                            Item banhTet = ItemService.gI().createNewItem((short) ConstItem.BANH_TET);
                            InventoryService.gI().addItemBag(player, banhTet);
                            InventoryService.gI().sendItemBag(player);
                            Service.gI().sendThongBaoOK(player, "Bạn Không Có Người Yêu Và Đã Làm Thành Công Bánh Test!");
                        } else {
                            Service.gI().sendThongBao(player, "Bạn Không đủ Nguyên Liệu Để Làm Bánh Test!");
                        }
                    }
                        break;
                }
            } else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_NAU_BANH_2) {
                switch (select) {
                    case 0: {
                        Item thitheo = InventoryService.gI().findItemBag(player, ConstItem.THIT_HEO);
                        Item thungNep = InventoryService.gI().findItemBag(player, ConstItem.THUNG_NEP);
                        Item thungDauXanh = InventoryService.gI().findItemBag(player, ConstItem.THUNG_DAU_XANH);
                        Item laDong = InventoryService.gI().findItemBag(player, ConstItem.LA_DONG);
                        Item trungVitMuoi = InventoryService.gI().findItemBag(player, ConstItem.TRUNG_VIT_MUOI);

                        if (thitheo != null && thungNep != null && thungDauXanh != null && laDong != null
                                && trungVitMuoi != null
                                && thitheo.quantity >= 10 && thungNep.quantity >= 10
                                && thungDauXanh.quantity >= 10 && laDong.quantity >= 10
                                && trungVitMuoi.quantity >= 1 && player.inventory.gold >= 5000000) {

                            InventoryService.gI().subQuantityItemsBag(player, thitheo, 10);
                            InventoryService.gI().subQuantityItemsBag(player, thungNep, 10);
                            InventoryService.gI().subQuantityItemsBag(player, thungDauXanh, 10);
                            InventoryService.gI().subQuantityItemsBag(player, laDong, 10);
                            InventoryService.gI().subQuantityItemsBag(player, trungVitMuoi, 1);
                            player.inventory.gold -= 5000000;

                            Item banhChung = ItemService.gI().createNewItem((short) ConstItem.BANH_CHUNG);
                            InventoryService.gI().addItemBag(player, banhChung);
                            InventoryService.gI().sendItemBag(player);

                            Service.gI().sendThongBaoOK(player, "Bạn Không Có Người Yêu Và Đã Làm Thành Công Bánh Chưng!");
                        } else {

                            Service.gI().sendThongBao(player, "Bạn Không Đủ Nguyên Liệu Để Làm Bánh Chưng!");
                        }
                    }
                        break;
                }
            }
        }
    }
}
