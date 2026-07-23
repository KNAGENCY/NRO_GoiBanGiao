package npc.npc_manifest;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.BossManager;
import consts.ConstNpc;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.NpcService;
import services.Service;
import utils.Util;

public class Potage extends Npc {

    public Potage(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            if (this.mapId == 140) {
                Player BossClone = BossManager.gI().findBossClone(player);
                if (BossClone != null) {
                    this.createOtherMenu(player, ConstNpc.BASE_MENU,
                            "Đang có 1 nhân bản của " + BossClone.name + " hãy chờ kết quả trận đấu",
                            "OK");
                } else {
                    this.createOtherMenu(player, ConstNpc.BASE_MENU,
                            "Hãy giúp ta đánh bại bản sao\nNgươi chỉ có 5 phút để hạ hắn\nPhần thưởng cho ngươi là 1 Chân Mệnh Và 1 bình Commeson",
                            "Hướng\ndẫn\nthêm", "OK\n(xin 30TV)", "Refuse");
                }
            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (this.mapId == 140) {
                if (player.iDMark.isBaseMenu()) {
                    Player BossClone = BossManager.gI().findBossClone(player);
                    if (BossClone == null) {
                        switch (select) {
                            case 0 ->
                                NpcService.gI().createTutorial(player, tempId, this.avartar,
                                        "Thứ bị phong ấn tại đây là vũ khí có tên Commesonđược tạo ra nhằm bảo vệ cho hành tinh PotaufeuTuy nhiên nó đã tàn phá mọi thứ trong quá khứ\nKhiến cư dân Potaufeu niêm phong nó với cái giáphải trả là mạng sống của họTa, Potage là người duy nhất sống sótvà ta đã bảo vệ phong ấn hơn một trăm năm.\nTuy nhiên bọn xâm lượt Gryll đã đến và giải thoát CommesonHãy giúp ta tiêu diệt bản sao do Commeson tạo ravà niêm phong Commeson một lần và mãi mãi");
                            case 1 -> {
                                if (!Util.isAfterMidnight(player.lastPkCommesonTime) && !player.isAdmin()) {
                                    Service.gI().sendThongBaoOK(player, "Hãy chờ đến ngày mai");
                                } else {
                                    var goldBar = InventoryService.gI().findItem(player.inventory.itemsBag, 457);
                                    if (goldBar != null && goldBar.quantity >= 30) {
                                        InventoryService.gI().subQuantityItemsBag(player, goldBar, 30);
                                        InventoryService.gI().sendItemBag(player);
                                        Service.gI().callNhanBan(player);
                                    } else {
                                        Service.gI().sendThongBaoOK(player, "Cần 30 TV Mới Có Thể Bắt Đầu!");
                                    }
                                }
                            }

                        }
                    }
                }
            }
        }
    }

}
