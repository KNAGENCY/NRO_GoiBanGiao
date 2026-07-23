package npc.npc_manifest;

import consts.ConstNpc;
import consts.cn;
import item.Item;
import jdbc.daos.PlayerDAO;
// import models.DragonNamecWar.TranhNgocService;
import npc.Npc;
import player.Player;
import services.PetService;
import services.Service;
import services.TaskService;
import services.func.ChangeMapService;
import services.func.TopService;
import services.func.Input;
import shop.ShopService;
import utils.Util;

public class Fide extends Npc {

    public Fide(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            TaskService.gI().checkDoneTaskTalkNpc(player, this);
            if (mapId == 5) {
                createOtherMenu(player, ConstNpc.BASE_MENU,
                        "|0| Ta Nắm Giữa 70% Vật Phẩm Trong Game\nNhưng Sức Mạnh Của Ta Rất Đặc Biệt\n|6|Đen Thì Nhận Được Cái Nịt\n\n"
                                + "Số điểm săn boss hiện có: "
                                + player.event.getEventPointBHM() + " điểm.\nBoss Trên 500tr Máu mới có Điểm",
                        "CẢI TRANG", "PHỤ KIỆN", "ĐỔI ĐIỂM\nSĂN BOSS", "ĐỔI TÊN\nNGƯỜI CHƠI");
            } else if (mapId == 153) {
                if (player.iDMark.getTranhNgoc() == 1) {
                    this.createOtherMenu(player, ConstNpc.BASE_MENU,
                            "Phắn đê! Ta không nói chuyện với sinh vật hạ đẳng", "Đóng");
                    return;
                }
                this.createOtherMenu(player, ConstNpc.BASE_MENU,
                        "Hãy mang ngọc rồng về cho ta", "Đưa ngọc", "Đóng");
            }
        }
    }
    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.iDMark.getIndexMenu() == ConstNpc.BASE_MENU) {
                if (this.mapId == 5) {
                    switch (select) {
                        // case 0 -> {
                        // ShopService.gI().opendShop(player, "SANTA_PHUKIEN", false);
                        // }
                        // case 6 -> {
                        // ShopService.gI().opendShop(player, "SHOP_TV", false);
                        // }
                        case 0 -> {
                            ShopService.gI().opendShop(player, "SHOP_NANGCAP", false);
                        }
                        case 1 -> {
                            ShopService.gI().opendShop(player, "SHOP_XU_KRAI", false);
                        }
//      case 2 -> {
//     if (player.playerTask.taskMain.id > 21) {

//         if (Util.isTrue(90, 100)) { // 90% không cho vào
//             Service.gI().sendThongBao(player, "BĐKB đã đầy, hãy thử lại");
//             return;
//         }

//         // 10% còn lại mới vào map
//         ChangeMapService.gI().changeMapBySpaceShip(player, 135, -1, -1);

//     } else {
//         Service.gI().sendThongBao(player, "Cần hoàn thành nhiệm vụ KuKu.");
//     }
// }
                        case 2-> {
                            ShopService.gI().opendShop(player, "SHOP_DIEMSANBOSS", false);
                        }
                        case 3 -> {
                            Input.gI().createFormChangeNameByCash(player);
                        }
                        // case 4 -> {
                        // if (player.getSession() != null) {
                        // this.createOtherMenu(player, 13000,
                        // "Thông tin từng loại đệ:\n|0|"
                        // + "Fide nhí: 20% chỉ số khi hợp thể\n|0|"
                        // + "Cell nhí: 25% chỉ số khi hợp thể\n|0|"
                        // + "Berrus nhí: 30% chỉ số khi hợp thể\n|0|"
                        // + "Adr bãi biển: 40% chỉ số khi hợp thể\n|0|"
                        // // + "Bưu gầy: 50% chỉ số khi hợp thể\n|6|"
                        // + "Bạn đang có: " + player.getSession().cash + " VND\n",
                        // "Mua đệ\nFide nhí\n[" + cn.de31 + "]", "Mua đệ\nCell nhí\n[" + cn.de32 + "]",
                        // "Mua đệ\nBưu nhí\n[" + cn.de33 + "]", "Mua đệ\nBerrus nhí\n[" + cn.de34
                        // + "]",
                        // "Mua đệ\nAdr bãi biển\n[" + cn.de36 + "]");
                        // // + "]", "Mua đệ\nAdr bãi biển\n[" + cn.de36 + "]", "Mua đệ\nBưu gầy\n[" +
                        // // cn.de35
                        // // + "]");
                        // player.iDMark.setIndexMenu(13000); // Cập nhật IndexMenu
                        // }

                        // break; // Thêm lệnh break để tránh fall-through
                        // }
                        // case 5 -> {
                        // TopService.showListTop(player, 8);
                        // }
                    }
                }
            } else if (player.iDMark.getIndexMenu() == 13000) {
                switch (select) {
                //     case 0:
                //         if (player.pet == null) {
                //             Service.gI().sendThongBao(player, "Ngươi cần phải có đệ mới sử dụng được chức năng này?");
                //             return;
                //         }
                //         for (Item item : player.pet.inventory.itemsBody) {
                //             if (item.isNotNullItem()) {
                //                 Service.gI().sendThongBao(player, "Cần bỏ đồ đệ tử đang mặc để sử dụng chức năng?");
                //                 return;
                //             }
                //         }
                //         if (player.getSession() != null && player.getSession().cash < cn.de31) {
                //             Service.gI().sendThongBao(player, "Bạn không đủ " + cn.de31 + " VND");
                //             return;
                //         }

                //         if (PlayerDAO.subcash(player, cn.de31)) {
                //             PetService.gI().createPetFideNhi(player, player.pet != null, player.gender);
                //             Service.gI().sendThongBao(player, "Bạn đã nhận được đệ Fide Nhí");
                //         } else {
                //             Service.gI().sendThongBao(player, "Đã có lỗi xảy ra !!");
                //         }
                //         break;

                //     case 1:
                //         if (player.pet == null) {
                //             Service.gI().sendThongBao(player, "Ngươi cần phải có đệ mới sử dụng được chức năng này?");
                //             return;
                //         }
                //         for (Item item : player.pet.inventory.itemsBody) {
                //             if (item.isNotNullItem()) {
                //                 Service.gI().sendThongBao(player, "Cần bỏ đồ đệ tử đang mặc để sử dụng chức năng?");
                //                 return;
                //             }
                //         }
                //         if (player.getSession() != null && player.getSession().cash < cn.de32) {
                //             Service.gI().sendThongBao(player, "Bạn không đủ " + cn.de32 + " VND");
                //             return;
                //         }

                //         if (PlayerDAO.subcash(player, cn.de32)) {
                //             PetService.gI().createPetCellNhi(player, player.pet != null, player.gender);
                //             Service.gI().sendThongBao(player, "Bạn đã nhận được đệ Cell Nhí");
                //         } else {
                //             Service.gI().sendThongBao(player, "Đã có lỗi xảy ra !!");
                //         }
                //         break;

                //     case 2:
                //         if (player.pet == null) {
                //             Service.gI().sendThongBao(player, "Ngươi cần phải có đệ mới sử dụng được chức năng này?");
                //             return;
                //         }
                //         for (Item item : player.pet.inventory.itemsBody) {
                //             if (item.isNotNullItem()) {
                //                 Service.gI().sendThongBao(player, "Cần bỏ đồ đệ tử đang mặc để sử dụng chức năng?");
                //                 return;
                //             }
                //         }
                //         if (player.getSession() != null && player.getSession().cash < cn.de33) {
                //             Service.gI().sendThongBao(player, "Bạn không đủ " + cn.de33 + " VND");
                //             return;
                //         }

                //         if (PlayerDAO.subcash(player, cn.de33)) {
                //             PetService.gI().createPetBuuNhi(player, player.pet != null, player.gender);
                //             Service.gI().sendThongBao(player, "Bạn đã nhận được đệ Bưu Nhí");
                //         } else {
                //             Service.gI().sendThongBao(player, "Đã có lỗi xảy ra !!");
                //         }
                //         break;

                //     case 3:
                //         if (player.pet == null) {
                //             Service.gI().sendThongBao(player, "Ngươi cần phải có đệ mới sử dụng được chức năng này?");
                //             return;
                //         }
                //         for (Item item : player.pet.inventory.itemsBody) {
                //             if (item.isNotNullItem()) {
                //                 Service.gI().sendThongBao(player, "Cần bỏ đồ đệ tử đang mặc để sử dụng chức năng?");
                //                 return;
                //             }
                //         }
                //         if (player.getSession() != null && player.getSession().cash < cn.de34) {
                //             Service.gI().sendThongBao(player, "Bạn không đủ " + cn.de34 + " VND");
                //             return;
                //         }

                //         if (PlayerDAO.subcash(player, cn.de34)) {
                //             PetService.gI().createPetBerrusNhi(player, player.pet != null, player.gender);
                //             Service.gI().sendThongBao(player, "Bạn đã nhận được đệ Berrus");

                //         } else {
                //             Service.gI().sendThongBao(player, "Đã có lỗi xảy ra !!");
                //         }
                //         break;

                //     case 4:
                //         if (player.pet == null) {
                //             Service.gI().sendThongBao(player, "Ngươi cần phải có đệ mới sử dụng được chức năng này?");
                //             return;
                //         }
                //         for (Item item : player.pet.inventory.itemsBody) {
                //             if (item.isNotNullItem()) {
                //                 Service.gI().sendThongBao(player, "Cần bỏ đồ đệ tử đang mặc để sử dụng chức năng?");
                //                 return;
                //             }
                //         }
                //         if (player.getSession() != null && player.getSession().cash < cn.de36) {
                //             Service.gI().sendThongBao(player, "Bạn không đủ " + cn.de36 + " VND");
                //             return;
                //         }

                //         if (PlayerDAO.subcash(player, cn.de36)) {
                //             PetService.gI().createPetAdrBeach(player, player.pet != null, player.gender);
                //             Service.gI().sendThongBao(player, "Bạn đã nhận được đệ Adr Bãi biển");
                //         } else {
                //             Service.gI().sendThongBao(player, "Đã có lỗi xảy ra !!");
                //         }
                //         break;
                //     case 5:
                //         if (player.pet == null) {
                //             Service.gI().sendThongBao(player, "Ngươi cần phải có đệ mới sử dụng được chức năng này?");
                //             return;
                //         }
                //         for (Item item : player.pet.inventory.itemsBody) {
                //             if (item.isNotNullItem()) {
                //                 Service.gI().sendThongBao(player, "Cần bỏ đồ đệ tử đang mặc để sử dụng chức năng?");
                //                 return;
                //             }
                //         }
                //         if (player.getSession() != null && player.getSession().cash < cn.de35) {
                //             Service.gI().sendThongBao(player, "Bạn không đủ " + cn.de35 + " VND");
                //             return;
                //         }

                //         if (PlayerDAO.subcash(player, cn.de35)) {
                //             PetService.gI().createPetMabuGay(player, player.pet != null, player.gender);
                //             Service.gI().sendThongBao(player, "Bạn đã nhận được đệ Bưu Gầy");
                //         } else {
                //             Service.gI().sendThongBao(player, "Đã có lỗi xảy ra !!");
                //         }
                //         break;
                // }
        // }
            //  else if (this.mapId == 153) {
            //     switch (select) {
            //         case 0 -> {
            //             if (player.iDMark.getTranhNgoc() == 2 && player.isHoldNamecBallTranhDoat) {
            //                 if (!Util.canDoWithTime(player.lastTimePickItem, 20000)) {
            //                     Service.gI().sendThongBao(player, "Vui lòng đợi "
            //                             + ((player.lastTimePickItem + 20000 - System.currentTimeMillis()) / 1000)
            //                             + " giây để có thể trả");
            //                     return;
            //                 }
                        //     TranhNgocService.getInstance().dropBall(player, (byte) 2);
                        //     player.zone.pointRed++;
                        //     if (player.zone.pointRed > ConstTranhNgocNamek.MAX_POINT) {
                        //         player.zone.pointRed = ConstTranhNgocNamek.MAX_POINT;
                        //     }
                        //     TranhNgocService.getInstance().sendUpdatePoint(player);
                         }
                    }
                }
            }
       // }
    }
//}