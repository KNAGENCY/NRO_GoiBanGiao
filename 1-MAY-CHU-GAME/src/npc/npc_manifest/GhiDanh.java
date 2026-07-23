// package npc.npc_manifest;

// /**
//  *
//  * @author GiaCatLuong
//  */
// import consts.ConstNpc;
// import item.Item;
// // import models.The23rdMartialArtCongress.The23rdMartialArtCongressService;
// import npc.Npc;
// import player.Player;
// import services.InventoryService;
// import services.ItemService;
// import services.NpcService;
// import services.PlayerService;
// import services.Service;
// import services.func.ChangeMapService;
// import utils.Util;

// public class GhiDanh extends Npc {

//     String[] menuselect = new String[] {};

//     public GhiDanh(int mapId, int status, int cx, int cy, int tempId, int avartar) {
//         super(mapId, status, cx, cy, tempId, avartar);
//     }

//     @Override
//     public void openBaseMenu(Player pl) {
//         if (canOpenNpc(pl)) {
//             switch (this.mapId) {

//                 // ❌ ĐÃ XÓA MAP 52 (World Martial Arts Tournament)

//                 case 129 -> {
//                     if (Util.isAfterMidnight(pl.lastTimePKDHVT23)) {
//                         pl.goldChallenge = 500_000_000;
//                         pl.rubyChallenge = 100;
//                         pl.levelWoodChest = 0;
//                     }

//                     long goldchallenge = pl.goldChallenge;
//                     long rubychallenge = pl.rubyChallenge;

//                     if (pl.levelWoodChest == 0) {
//                         menuselect = new String[] {
//                             "Hướng\ndẫn\nthêm",
//                             "Thi đấu\n" + Util.numberToMoney(rubychallenge) + " hồng ngọc",
//                             "Thi đấu\n" + Util.numberToMoney(goldchallenge) + " vàng",
//                             "Về\nĐại Hội\nVõ Thuật"
//                         };
//                     } else {
//                         menuselect = new String[] {
//                             "Hướng\ndẫn\nthêm",
//                             "Thi đấu\n" + Util.numberToMoney(rubychallenge) + " hồng ngọc",
//                             "Thi đấu\n" + Util.numberToMoney(goldchallenge) + " vàng",
//                             "Nhận\nthưởng\nRương Cấp\n" + pl.levelWoodChest,
//                             "Về\nĐại Hội\nVõ Thuật"
//                         };
//                     }

//                     this.createOtherMenu(pl, ConstNpc.BASE_MENU,
//                             "Đại hội võ thuật lần thứ 23\n"
//                             + "Diễn ra bất kể ngày đêm, ngày nghỉ, ngày lễ\n"
//                             + "Phần thưởng vô cùng quý giá\n"
//                             + "Nhanh chóng tham gia nào",
//                             menuselect, "Từ chối");
//                 }

//                 default -> super.openBaseMenu(pl);
//             }
//         }
//     }

//     @Override
//     public void confirmMenu(Player player, int select) {
//         if (!canOpenNpc(player)) return;

//         // ❌ ĐÃ XÓA xử lý mapId == 52

//         if (this.mapId == 129) {
//             switch (player.iDMark.getIndexMenu()) {

//                 case ConstNpc.BASE_MENU -> {
//                     long goldchallenge = player.goldChallenge;
//                     long rubychallenge = player.rubyChallenge;

//                     if (player.levelWoodChest == 0) {
//                         switch (select) {
//                             case 0 ->
//                                 NpcService.gI().createTutorial(player, tempId, this.avartar,
//                                         ConstNpc.NPC_DHVT23);

//                             case 1, 2 -> {
//                                 if (player.levelWoodChest != 12) {
//                                     if (InventoryService.gI().finditemWoodChest(player)) {

//                                         if (select == 1) {
//                                             if (player.inventory.ruby >= rubychallenge) {
//                                                 The23rdMartialArtCongressService.gI()
//                                                         .startChallenge(player);
//                                                 player.inventory.ruby -= rubychallenge;
//                                             } else {
//                                                 Service.gI().sendThongBao(player,
//                                                         "Bạn không đủ hồng ngọc, còn thiếu "
//                                                         + Util.numberToMoney(
//                                                                 rubychallenge - player.inventory.ruby));
//                                                 return;
//                                             }
//                                         } else {
//                                             if (player.inventory.gold >= goldchallenge) {
//                                                 The23rdMartialArtCongressService.gI()
//                                                         .startChallenge(player);
//                                                 player.inventory.gold -= goldchallenge;
//                                             } else {
//                                                 Service.gI().sendThongBao(player,
//                                                         "Bạn không đủ vàng, còn thiếu "
//                                                         + Util.numberToMoney(
//                                                                 goldchallenge - player.inventory.gold));
//                                                 return;
//                                             }
//                                         }

//                                         PlayerService.gI().sendInfoHpMpMoney(player);
//                                         player.goldChallenge += 500_000_000;
//                                         player.rubyChallenge += 100;

//                                     } else {
//                                         Service.gI().sendThongBao(player,
//                                                 "Hãy mở rương báu vật trước");
//                                     }
//                                 } else {
//                                     Service.gI().sendThongBao(player,
//                                             "Bạn đã vô địch giải. Vui lòng chờ đến ngày mai");
//                                 }
//                             }

//                             case 3 ->
//                                 ChangeMapService.gI().changeMapNonSpaceship(
//                                         player, 52, player.location.x, 336);
//                         }

//                     } else {
//                         switch (select) {
//                             case 0 ->
//                                 NpcService.gI().createTutorial(player, tempId, this.avartar,
//                                         ConstNpc.NPC_DHVT23);

//                             case 1, 2 -> {
//                                 if (player.levelWoodChest != 12) {
//                                     if (InventoryService.gI().finditemWoodChest(player)) {

//                                         if (select == 1) {
//                                             if (player.inventory.ruby >= rubychallenge) {
//                                                 The23rdMartialArtCongressService.gI()
//                                                         .startChallenge(player);
//                                                 player.inventory.ruby -= rubychallenge;
//                                             } else {
//                                                 Service.gI().sendThongBao(player,
//                                                         "Bạn không đủ hồng ngọc");
//                                                 return;
//                                             }
//                                         } else {
//                                             if (player.inventory.gold >= goldchallenge) {
//                                                 The23rdMartialArtCongressService.gI()
//                                                         .startChallenge(player);
//                                                 player.inventory.gold -= goldchallenge;
//                                             } else {
//                                                 Service.gI().sendThongBao(player,
//                                                         "Bạn không đủ vàng");
//                                                 return;
//                                             }
//                                         }

//                                         PlayerService.gI().sendInfoHpMpMoney(player);
//                                         player.goldChallenge += 500_000_000;
//                                         player.rubyChallenge += 100;

//                                     } else {
//                                         Service.gI().sendThongBao(player,
//                                                 "Hãy mở rương báu vật trước");
//                                     }
//                                 } else {
//                                     Service.gI().sendThongBao(player,
//                                             "Bạn đã vô địch giải. Vui lòng chờ đến ngày mai");
//                                 }
//                             }

//                             case 3 ->
//                                 this.createOtherMenu(player, 1,
//                                         "Phần thưởng của bạn đang ở cấp "
//                                         + player.levelWoodChest + " / 12\n"
//                                         + "Mỗi ngày chỉ được nhận 1 lần\n"
//                                         + "Bạn có chắc muốn nhận ngay?",
//                                         "OK", "Từ chối");

//                             case 4 ->
//                                 ChangeMapService.gI().changeMapNonSpaceship(
//                                         player, 52, player.location.x, 336);
//                         }
//                     }
//                 }

//                 case 1 -> {
//                     if (select == 0) {
//                         if (InventoryService.gI().finditemWoodChest(player)) {
//                             if (InventoryService.gI().getCountEmptyBag(player) > 0) {
//                                 Item it = ItemService.gI().createNewItem((short) 570);
//                                 it.itemOptions.add(
//                                         new Item.ItemOption(72, player.levelWoodChest));
//                                 it.itemOptions.add(
//                                         new Item.ItemOption(30, 0));
//                                 it.createTime = System.currentTimeMillis();

//                                 InventoryService.gI().addItemBag(player, it);
//                                 InventoryService.gI().sendItemBag(player);

//                                 player.levelWoodChest = 0;
//                                 player.lastTimeRewardWoodChest = System.currentTimeMillis();

//                                 NpcService.gI().createMenuConMeo(player, -1, -1,
//                                         "Bạn nhận được\n|1|Rương Gỗ\n"
//                                         + "|2|Giấu bên trong nhiều vật phẩm quý giá",
//                                         "OK");
//                             } else {
//                                 this.npcChat(player,
//                                         "Hành trang đã đầy, cần 1 ô trống");
//                             }
//                         } else {
//                             Service.gI().sendThongBao(player,
//                                     "Hãy mở rương báu vật trước");
//                         }
//                     }
//                 }
//             }
//         }
//     }
// }
