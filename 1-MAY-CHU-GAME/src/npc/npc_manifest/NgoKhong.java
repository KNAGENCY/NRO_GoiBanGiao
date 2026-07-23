package npc.npc_manifest;

/**
 *
 * @author GiaCatLuong
 */
import consts.ConstItem;
import consts.ConstNpc;
import item.Item;
import item.Item.ItemOption;
import java.util.Random;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;

public class NgoKhong extends Npc {

    public NgoKhong(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            switch (mapId) {
                case 122 -> {
                    createOtherMenu(player, ConstNpc.BASE_MENU, "|0|Ăn 1 Quả, Trả 1 Cục Vàng, Để Trống Hành Trang, Không Thì Nhót\n",
                            "Tặng\nHồng Đào", "Tặng\nHồng Đào\nchín");
                }
                default ->
                    super.openBaseMenu(player);
            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.iDMark.isBaseMenu()) {
                switch (select) {
                  case 0: {
    Item daohong = InventoryService.gI().findItemBag(player, ConstItem.QUA_HONG_DAO);

    if (daohong != null && daohong.quantity >= 1) {
        // Trừ 1 quả hồng đào
        InventoryService.gI().subQuantityItemsBag(player, daohong, 1);

        // Tạo các sao pha lê có option riêng
        Item item1 = ItemService.gI().createNewItem((short) ConstItem.SAO_PHA_LE);
        item1.itemOptions.add(new ItemOption(95, 5));

        Item item2 = ItemService.gI().createNewItem((short) ConstItem.SAO_PHA_LE_442);
        item2.itemOptions.add(new ItemOption(96, 5));

        Item item3 = ItemService.gI().createNewItem((short) ConstItem.SAO_PHA_LE_443);
        item3.itemOptions.add(new ItemOption(97, 5));

        Item item4 = ItemService.gI().createNewItem((short) ConstItem.SAO_PHA_LE_447);
        item4.itemOptions.add(new ItemOption(101, 5));

        // Đưa vào mảng chọn ngẫu nhiên
        Item[] chuItems = new Item[] { item1, item2, item3, item4 };

        Random rand = new Random();
        Item selectedItem = chuItems[rand.nextInt(chuItems.length)];

        // Thêm vào túi
        InventoryService.gI().addItemBag(player, selectedItem);
        InventoryService.gI().sendItemBag(player);

        // Thông báo
        Service.gI().sendThongBaoOK(player, "Bạn nhận được: " + selectedItem.template.name);
    } else {
        Service.gI().sendThongBao(player, "Cần 1 Quả Hồng Đào!");
    }
    break;
}

                    case 1: {
                        Item daohong = InventoryService.gI().findItemBag(player, ConstItem.QUA_HONG_DAO_CHIN);
                        if (daohong != null && daohong.quantity >= 1) {
                            InventoryService.gI().subQuantityItemsBag(player, daohong, 1);
                            Item[] chuItems = new Item[] {
                                    ItemService.gI().createNewItem((short) ConstItem.NGOC_RONG_7_SAO),
                                    ItemService.gI().createNewItem((short) ConstItem.NGOC_RONG_6_SAO),
                                    ItemService.gI().createNewItem((short) ConstItem.NGOC_RONG_5_SAO),
                                    ItemService.gI().createNewItem((short) ConstItem.NGOC_RONG_3_SAO)
                            };
                            Random rand = new Random();
                            Item selectedItem = chuItems[rand.nextInt(chuItems.length)];
                            InventoryService.gI().addItemBag(player, selectedItem);
                            InventoryService.gI().sendItemBag(player);
                            Service.gI().sendThongBaoOK(player, "Bạn Nhận Được: " + selectedItem.template.name);
                        } else {
                            Service.gI().sendThongBao(player, "Cần 1 Quả Hồng Đào!");
                        }
                        break;
                    }
                }
            }
        }
    }

}
