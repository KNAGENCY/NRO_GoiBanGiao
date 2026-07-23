package npc.npc_manifest;

/**
 *
 * @author GiaCatLuong
 */
import clan.Clan;
import consts.ConstNpc;
import consts.ConstPlayer;
import item.Item;
import java.util.ArrayList;
// import models.DragonNamecWar.TranhNgoc;
// import models.DragonNamecWar.TranhNgocService;
import npc.Npc;
import player.Player;
import server.Manager;
import services.InventoryService;
import services.NpcService;
import services.RewardService;
import services.Service;
import services.TaskService;
import services.func.ChangeMapService;
import services.func.Input;
import shop.ShopService;
import utils.Util;

public class TruongLaoGuru extends Npc {

    public TruongLaoGuru(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            if (!TaskService.gI().checkDoneTaskTalkNpc(player, this)) {
                if (player.gender != ConstPlayer.NAMEC) {
                    NpcService.gI().createTutorial(player, tempId, avartar,
                            "Con hãy về hành tinh của mình mà thể hiện");
                    return;
                }
                ArrayList<String> menu = new ArrayList<>();
                if (!player.canReward) {
                    menu.add("Nhiệm vụ");
                    // menu.add("Học\nKỹ năng");
                    Clan clan = player.clan;
                    if (clan != null) {
                        menu.add("Go To\nNHS");
                        if (clan.isLeader(player)) {
                            menu.add("Giải tán\nBang hội");
                        }
                    }
                } else {
                    menu.add("Giao\nLân con");
                }
                String[] menus = menu.toArray(String[]::new);
                createOtherMenu(player, ConstNpc.BASE_MENU,
                        "Có Thể Tìm NHS Tối Thượng Ở Đây\nTỉ Lệ Là 50% Thành Công", menus);

            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            if (player.canReward) {
                RewardService.gI().rewardLancon(player);
                return;
            }
            if (player.iDMark.isBaseMenu()) {

                switch (select) {
                    case 0 ->
                        NpcService.gI().createTutorial(player, tempId, avartar,
                                player.playerTask.taskMain.subTasks.get(player.playerTask.taskMain.index).name);
                  
                    case 1 -> {
                        long phi = 5_000_000L;
                    

                        if (player.playerTask.taskMain.id < 23) {
                            Service.gI().sendThongBao(player,
                                    "Yêu Cầu Xong NV TĐST!");
                            return;
                        }

                        if (player.inventory.gold < phi) {
                            Service.gI().sendThongBao(player,
                                    "Không đủ 50 triệu vàng!");
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
                        if (clan != null) {
                            if (clan.isLeader(player)) {
                                createOtherMenu(player, 3, "Con có chắc muốn giải tán bang hội không?", "Đồng ý",
                                        "Refuse");
                            }
                        }
                    }
                }
            } else if (player.iDMark.getIndexMenu() == 3) {
                Clan clan = player.clan;
                if (clan != null) {
                    if (clan.isLeader(player)) {
                        if (select == 0) {
                            Input.gI().createFormGiaiTanBangHoi(player);
                        }
                    }
                }
            }
        }
    }
}
