// package npc.npc_manifest;

// /**
//  *
//  * @author GiaCatLuong
//  */
// import clan.Clan;
// import consts.ConstNpc;
// import item.Item;
// import java.sql.Connection;
// import java.sql.PreparedStatement;
// import java.sql.ResultSet;
// import java.sql.SQLException;
// import java.util.ArrayList;
// import java.util.logging.Level;
// import java.util.logging.Logger;
// import jdbc.DBConnecter;
// import npc.Npc;
// import models.War.War;
// import org.json.simple.JSONArray;
// import org.json.simple.JSONObject;
// import org.json.simple.JSONValue;
// import player.Archivement;
// import player.Player;
// import services.InventoryService;
// import services.ItemService;
// import services.NpcService;
// import services.RewardService;
// import services.Service;
// import services.TaskService;
// import services.func.ChangeMapService;
// import services.func.Input;
// import services.func.TopService;
// import shop.ShopService;
// import utils.Util;

// public class QuyLaoKame extends Npc {

//     public QuyLaoKame(int mapId, int status, int cx, int cy, int tempId, int avartar) {
//         super(mapId, status, cx, cy, tempId, avartar);
//     }

//     @Override
//     public void openBaseMenu(Player player) {
//         Item ruacon = InventoryService.gI().findItemBag(player, 874);
//         if (canOpenNpc(player)) {
//             ArrayList<String> menu = new ArrayList<>();
//             if (!player.canReward) {
//                 menu.add("Nói\nChuyện");
//                 menu.add("Map 02\nCày Vàng");
//                 menu.add("Đổi Phiếu\nBé Ngon");
//                    menu.add("Đổi Quà");
//                 if (ruacon != null && ruacon.quantity >= 1) {
//                     menu.add("Con Lân");
//                 }
//             } else {
//                 menu.add("Con Lân");
//             }
//             String[] menus = menu.toArray(String[]::new);
//             if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
//                 this.createOtherMenu(player, ConstNpc.BASE_MENU,
//                         "21H30 - 21H50 Hãy Đi Cày Vàng?", menus);
//             }
//         }
//     }

//     @Override
//     public void confirmMenu(Player player, int select) {
//         if (!canOpenNpc(player)) {
//             return;
//         }

//         if (player.canReward) {
//             RewardService.gI().rewardLancon(player);
//             return;
//         }

//         switch (player.iDMark.getIndexMenu()) {

//             case ConstNpc.BASE_MENU -> {
//                 switch (select) {

//                  case 0 -> {
//     ArrayList<String> menu = new ArrayList<>();

//     // Luôn có
//     menu.add("Nhiệm vụ");
//     menu.add("Go To\nNHS");

//     // Chỉ leader mới thấy giải tán
//     Clan clan = player.clan;
//     if (clan != null && clan.isLeader(player)) {
//         menu.add("Giải tán\nBang hội");
//     }

//     String[] menus = menu.toArray(String[]::new);
//     this.createOtherMenu(player, 0,
//             "Có Thể Tìm NHS Tối Thượng Ở Đây\nTỉ Lệ Là 50% Thành Công",
//             menus);
// }


//                case 1 -> {

    
//     if (player.fusion != null && player.fusion.typeFusion != 0) {
//         Service.gI().sendThongBao(player,
//                 "Không thể tham gia khi đang hợp thể!");
//         return;
//     }


//     if (player.nPoint.power < 17_000_000_000L) {
//         Service.gI().sendThongBao(player,
//                 "Cần tối thiểu 17 Tỷ Sức Mạnh để tham gia!");
//         return;
//     }

//     if (!War.gI().isOpenTime()) {
//         Service.gI().sendThongBao(player,
//                 "Quay lại vào lúc 21H30 - 21H45");
//         return;
//     }

    
//     boolean empty = true;
//     for (int i = 0; i < 14; i++) {
//         if (i < player.inventory.itemsBody.size()) {
//             Item item = player.inventory.itemsBody.get(i);
//             if (item != null && item.isNotNullItem()) {
//                 empty = false;
//                 break;
//             }
//         }
//     }

//     if (!empty) {
//         Service.gI().sendThongBao(player,
//                 "Hãy tháo hết trang bị để tham gia!");
//         return;
//     }

   
//     ChangeMapService.gI()
//             .changeMapNonSpaceship(player,
//                     War.MAP_PREPARE_ID, 1093, 600);
// }



//                     case 2 -> {
//                         Item phieu = InventoryService.gI()
//                                 .findItemBag(player, 1194);
//                         int soLuong = phieu != null ? phieu.quantity : 0;

//                         this.createOtherMenu(player,
//                                 ConstNpc.ĐỔI_PHIẾU_BÉ_NGOAN,
//                                 "Con đang có " + soLuong + " phiếu bé ngon\n"
//                                         + "Con muốn đổi quà gì nào\n"
//                                         + "Pem Boss Nhiệm Vụ là có ...",
//                                 "Up Skill 2\n1 Phiếu",
//                                 "Đổi Skill 2\n20 Phiếu",
//                                 "Đổi Skill 3\n30 Phiếu",
//                                 "Trứng\n50 Phiếu",
//                                 "Đóng");
//                     }

//                     case 3 -> { ShopService.gI().opendShop(player, "BUNMA_PHU_KIEN", false);}
//                 }
//             }

//             case ConstNpc.ĐỔI_PHIẾU_BÉ_NGOAN -> {
//                 int soluong = switch (select) {
//                     case 0 -> 1;
//                     case 1 -> 20;
//                     case 2 -> 30;
//                     case 3 -> 50;
//                     default -> 0;
//                 };

//                 if (soluong == 0)
//                     return;

//                 if (InventoryService.gI().getCountEmptyBag(player) < 2) {
//                     Service.gI().sendThongBao(player,
//                             "Cần 2 ô hành trang!");
//                     return;
//                 }

//                 Item phieu = InventoryService.gI()
//                         .findItemBag(player, 1194);
//                 if (phieu == null || phieu.quantity < soluong) {
//                     Service.gI().sendThongBao(player,
//                             "Không đủ Phiếu Bé Ngon!");
//                     return;
//                 }

//                 InventoryService.gI()
//                         .subQuantityItemsBag(player, phieu, soluong);

//                 switch (select) {
//                     case 0 ->
//                         InventoryService.gI().addItemBag(player,
//                                 ItemService.gI().createNewItem((short) 403));
//                     case 1 ->
//                         InventoryService.gI().addItemBag(player,
//                                 ItemService.gI().createNewItem((short) 1744));
//                     case 2 ->
//                         InventoryService.gI().addItemBag(player,
//                                 ItemService.gI().createNewItem((short) 1748));

//                     case 3 -> {
//                         InventoryService.gI().addItemBag(player,
//                                 ItemService.gI().createNewItem((short) 568));
                       
                              
//                     }
//                 }

//                 InventoryService.gI().sendItemBag(player);
//                 Service.gI().sendThongBao(player,
//                         "Đổi thành công!");
//             }

//             case 0 -> {
//                 switch (select) {
//                     case 0 ->
//                         NpcService.gI().createTutorial(player,
//                                 tempId, avartar,
//                                 player.playerTask.taskMain.subTasks
//                                         .get(player.playerTask.taskMain.index).name);

//                     case 1 -> {
//                         long phi = 5_000_000L;
                

//                         if (player.playerTask.taskMain.id < 23) {
//                             Service.gI().sendThongBao(player,
//                                     "Yêu Cầu Xong NV TĐST!");
//                             return;
//                         }

//                         if (player.inventory.gold < phi) {
//                             Service.gI().sendThongBao(player,
//                                     "Không đủ 50 triệu vàng!");
//                             return;
//                         }

//                         player.inventory.gold -= phi;
//                         Service.gI().sendMoney(player);

//                         int rd = Util.nextInt(100);
//                         if (rd < 50) {
//                             ChangeMapService.gI().changeMap(player, 123, -1, 200, 336);
//                         } else {
//                             ChangeMapService.gI().changeMap(player, 122, -1, 200, 336);
//                         }
//                     }
//                     case 2 -> {
//                         Clan clan = player.clan;
//                         if (clan != null && clan.isLeader(player)) {
//                             createOtherMenu(player, 3,
//                                     "Con có chắc muốn giải tán bang hội?",
//                                     "Đồng ý", "Hủy");
//                         }
//                     }
//                 }
//             }

//             case 3 -> {
//                 Clan clan = player.clan;
//                 if (clan != null && clan.isLeader(player) && select == 0) {
//                     Input.gI().createFormGiaiTanBangHoi(player);
//                 }
//             }
//         }
//     }
// }
package npc.npc_manifest;

/**
 *
 * @author GiaCatLuong
 */
import clan.Clan;
import consts.ConstNpc;
import item.Item;
import java.util.ArrayList;
import npc.Npc;
import npc.NpcFactory;
import models.War.War;
import player.Player;
import services.InventoryService;
import services.NpcService;
import services.RewardService;
import services.Service;
import services.TaskService;
import services.func.ChangeMapService;
import services.func.Input;
import shop.ShopService;
import utils.Util;

public class QuyLaoKame extends Npc {

    public QuyLaoKame(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        Item ruacon = InventoryService.gI().findItemBag(player, 874);
        if (canOpenNpc(player)) {
            ArrayList<String> menu = new ArrayList<>();
            if (!player.canReward) {
                menu.add("Nói\nChuyện");
                menu.add("Map War");
                menu.add("Đổi\nThưởng");
                menu.add("Siêu Tốc\nBĐKB");
                menu.add("Đổi mật khẩu");
                if (ruacon != null && ruacon.quantity >= 1) {
                    menu.add("Con Lân");
                }
            } else {
                menu.add("Con Lân");
            }
            String[] menus = menu.toArray(String[]::new);
            if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
                this.createOtherMenu(player, ConstNpc.BASE_MENU,
        "|7|21H30 - 21H50 Map War Sẽ Mở\n Đây Là Sự Kiện Thường Niên Mỗi Tối\n"
        + "Nơi Hội Đụ Anh Tài Với Cực Nhiều Phần Quà\n\n"
        + "|7|BĐKB chuẩn Teamobi:\n"
        + "Cần Nén Sức Đánh 60k ~ 80k mới đc TNSM",
        menus);
            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) {
            return;
        }

        if (player.canReward) {
            RewardService.gI().rewardLancon(player);
            return;
        }

        switch (player.iDMark.getIndexMenu()) {

            case ConstNpc.BASE_MENU -> {
                switch (select) {

                    case 0 -> {
                        ArrayList<String> menu = new ArrayList<>();

                        // Luôn có
                        menu.add("Nhiệm vụ");
                        menu.add("Go To\nNHS");

                        // Chỉ leader mới thấy giải tán
                        Clan clan = player.clan;
                        if (clan != null && clan.isLeader(player)) {
                            menu.add("Giải tán\nBang hội");
                        }

                        String[] menus = menu.toArray(String[]::new);
                        this.createOtherMenu(player, 0,
                                "Có Thể Tìm NHS Tối Thượng Ở Đây\nTỉ Lệ Là 50% Thành Công",
                                menus);
                    }

                    case 1 -> {

                        if (player.fusion != null && player.fusion.typeFusion != 0) {
                            Service.gI().sendThongBao(player,
                                    "Không thể tham gia khi đang hợp thể!");
                            return;
                        }

                        if (player.nPoint.power < 17_000_000_000L) {
                            Service.gI().sendThongBao(player,
                                    "Cần tối thiểu 17 Tỷ Sức Mạnh để tham gia!");
                            return;
                        }

                        if (!War.gI().isOpenTime()) {
                            Service.gI().sendThongBao(player,
                                    "Quay lại vào lúc 21H30 - 21H45");
                            return;
                        }

                        boolean empty = true;
                        for (int i = 0; i < 14; i++) {
                            if (i < player.inventory.itemsBody.size()) {
                                Item item = player.inventory.itemsBody.get(i);
                                if (item != null && item.isNotNullItem()) {
                                    empty = false;
                                    break;
                                }
                            }
                        }

                        if (!empty) {
                            Service.gI().sendThongBao(player,
                                    "Hãy tháo hết trang bị để tham gia!");
                            return;
                        }

                        ChangeMapService.gI()
                                .changeMapNonSpaceship(player,
                                        War.MAP_PREPARE_ID, 1093, 600);
                    }

                 case 2 -> {
    ShopService.gI().opendShop(player, "BUNMA_PHU_KIEN", false);
}

case 3 -> {
    if (player.playerTask.taskMain.id > 21) {

        if (Util.isTrue(50, 100)) {
            Service.gI().sendThongBao(player, "BĐKB đã đầy, hãy thử lại sau 5 giây");
            return;
        }

        ChangeMapService.gI().changeMapBySpaceShip(player, 135, -1, -1);

    } else {
        Service.gI().sendThongBao(player, "Cần hoàn thành nhiệm vụ KuKu.");
    }
}
case 4 -> Input.gI().createFormChangePassword(player);
                }
            }

            case 10001 -> {
                Object value = NpcFactory.PLAYERID_OBJECT.remove(player.id);
                if (select == 0 && value instanceof String) {
                    Service.gI().changePassword(player, (String) value);
                }
            }

            case 0 -> {
                switch (select) {
                    case 0 ->
                        NpcService.gI().createTutorial(player,
                                tempId, avartar,
                                player.playerTask.taskMain.subTasks
                                        .get(player.playerTask.taskMain.index).name);

                    case 1 -> {
                        long phi = 5_000_000L;

                        if (player.playerTask.taskMain.id < 23) {
                            Service.gI().sendThongBao(player,
                                    "Yêu Cầu Xong NV TĐST!");
                            return;
                        }

                        if (player.inventory.gold < phi) {
                            Service.gI().sendThongBao(player,
                                    "Không đủ 5 triệu vàng!");
                            return;
                        }

                        player.inventory.gold -= phi;
                        Service.gI().sendMoney(player);

                        int rd = Util.nextInt(100);
                        if (rd < 50) {
                            ChangeMapService.gI().changeMap(player, 123, -1, 200, 336);
                        } else {
                            ChangeMapService.gI().changeMap(player, 122, -1, 200, 336);
                        }
                    }
                    case 2 -> {
                        Clan clan = player.clan;
                        if (clan != null && clan.isLeader(player)) {
                            createOtherMenu(player, 3,
                                    "Con có chắc muốn giải tán bang hội?",
                                    "Đồng ý", "Hủy");
                        }
                    }
                }
            }

            case 3 -> {
                Clan clan = player.clan;
                if (clan != null && clan.isLeader(player) && select == 0) {
                    Input.gI().createFormGiaiTanBangHoi(player);
                }
            }
        }
    }
}
