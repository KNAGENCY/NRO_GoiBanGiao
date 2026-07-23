package npc.npc_manifest;

import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import shop.ShopService;
import services.Service;
import item.Item;

public class ThanVuTru extends Npc {

    private static final int MENU_MAIN = 0;
    private static final int MENU_EXCHANGE = 1;
    private static final long TIME_WAIT = 30 * 60 * 1000; // 30 phút

    public ThanVuTru(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) return;

        long tn = player.nPoint.tiemNang;
        String tnText = String.format("%,d", tn);

        String option0 = "Chuyển hóa\nTiềm Năng";

        if (player.timeExchangeTN > 0) {
            long now = System.currentTimeMillis();
            long diff = now - player.timeExchangeTN;

            if (diff >= TIME_WAIT) {
                option0 = "Nhận Đá";
            } else {
                long remain = (TIME_WAIT - diff) / 1000;
                long minutes = remain / 60;
                long seconds = remain % 60;

                option0 = "Còn Lại\n" + minutes + "p " + seconds + "s";
            }
        }

        createOtherMenu(player, MENU_MAIN,
                "Ta là Thần Vú Trụ,\n"
                + "Ta sẽ giúp người học Skill,\n"
                + "♥  ♥  ♥  ♥  ♥\n"
                + "Tiềm năng hiện tại: " + tnText + "\n"
                + "Ngươi Cần: 500,000,000 Tiềm năng để Chuyển Hóa\n"
                + "Rồi dùng Đá Quyền Năng để đổi Skill",
                option0,
                "Học Skill",
                "Từ chối");
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) return;

        switch (player.iDMark.getIndexMenu()) {

            case MENU_MAIN -> {
                switch (select) {

                    case 0 -> {
                        long now = System.currentTimeMillis();

                        if (player.timeExchangeTN > 0) {
                            long diff = now - player.timeExchangeTN;

                            if (diff >= TIME_WAIT) {
                                player.iDMark.setIndexMenu(MENU_EXCHANGE);
                                createOtherMenu(player, MENU_EXCHANGE,
                                        "Chuyển hóa đã hoàn tất\nNgươi muốn nhận Đá Quyền Năng?",
                                        "Nhận",
                                        "Từ chối");
                            } else {
                                long remain = (TIME_WAIT - diff) / 1000;
                                long minutes = remain / 60;
                                long seconds = remain % 60;

                                Service.gI().sendThongBao(player,
                                        "Còn " + minutes + " phút " + seconds + " giây");
                            }
                            return;
                        }

                        // ===== CHƯA CHUYỂN HÓA =====
                        player.iDMark.setIndexMenu(MENU_EXCHANGE);

                        createOtherMenu(player, MENU_EXCHANGE,
                                "Đổi 500.000.000 Tiềm Năng\nlấy 100 Đá Quyền Năng?\n(Thời gian: 30 phút)",
                                "Đồng ý",
                                "Từ chối");
                    }

                    case 1 -> ShopService.gI().opendShop(player, "SHOP_TV", false);
                }
            }

            case MENU_EXCHANGE -> {
                if (select != 0) return;

                long need = 500_000_000L;

                // ===== NHẬN ĐÁ =====
                if (player.timeExchangeTN > 0) {
                    long now = System.currentTimeMillis();

                    if (now - player.timeExchangeTN < TIME_WAIT) {
                        Service.gI().sendThongBao(player,
                                "Chưa hoàn tất chuyển hóa!");
                        return;
                    }

                    if (InventoryService.gI().getCountEmptyBag(player) < 1) {
                        Service.gI().sendThongBao(player,
                                "Hành trang không đủ chỗ trống!");
                        return;
                    }

                    Item item = ItemService.gI().createNewItem((short) 1505, 100);
                    // item.itemOptions.add(new Item.ItemOption(30, 1));

                    InventoryService.gI().addItemBag(player, item);
                    InventoryService.gI().sendItemBag(player);

                    // RESET trạng thái
                    player.timeExchangeTN = 0;

                    Service.gI().sendThongBao(player,
                            "Nhận 100 Đá Quyền Năng!");
                    return;
                }

                // ===== BẮT ĐẦU CHUYỂN HÓA =====
                if (player.nPoint.tiemNang < need) {
                    Service.gI().sendThongBao(player,
                            "Không đủ Tiềm Năng!");
                    return;
                }

                if (InventoryService.gI().getCountEmptyBag(player) < 1) {
                    Service.gI().sendThongBao(player,
                            "Hành trang không đủ chỗ trống!");
                    return;
                }

    
                player.nPoint.tiemNang -= need;
                Service.gI().point(player);


                player.timeExchangeTN = System.currentTimeMillis();


                Service.gI().chat(player,
                        "[Thần Vũ Trụ]: Ta sẽ chuyển hóa sức mạnh của ngươi...");
                Service.gI().chat(player,
                        "[Thần Vũ Trụ]: Quá trình sẽ mất 30 phút...");

                Service.gI().sendThongBao(player,
                        "Hãy quay lại sau 30 phút!");
            }
        }
    }
}