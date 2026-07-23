package server;

/**
 * @author GiaCatLuong
 */
// import minigame.TaiXiu.TaiXiu;
import boss.BossManager;
// import boss.GasDestroyManager;
import boss.OtherBossManager;

import consts.ConstNpc;
import data.ItemData;
import item.Item;
import jdbc.DBConnecter;
import jdbc.daos.ShopDAO;

import java.util.List;
import java.util.Timer;
import java.util.TimerTask;
import utils.Logger;
import network.SessionManager;

import GiaCatLuongManager.SystemMetrics;
// import minigame.LuckyNumber.LuckyNumber;
import models.GiftCode.GiftCodeManager;
import models.ShenronEvent.ShenronEvent;
import models.ShenronEvent.ShenronEventManager;
import player.Bot.BotManager;
import player.Bot.NewBot;
import player.Pet;
import player.Player;
import java.sql.Connection;

import services.InventoryService;
import services.ItemService;
import services.ChatGlobalService;
import services.NpcService;
import services.PetService;
import services.PlayerService;
import services.Service;
import services.SkillService;
import services.TaskService;
import services.func.ChangeMapService;
import services.func.Input;
import services.func.SummonDragonNamek;
import services.func.TopService;
import skill.Skill;
import utils.Util;

public class Command {

    private static Command instance;

    public static Command gI() {
        if (instance == null) {
            instance = new Command();
        }
        return instance;
    }

    public void chat(Player player, String text) {
        if (!check(player, text)) {
            Service.gI().chat(player, text);
        }
    }

    public boolean check(Player player, String text) {
        if (text.equalsIgnoreCase("hd1")) {
            Service.gI().sendThongBaoOK(player, ChatGlobalService.SYSTEM_NOTICE_DETAIL_1);
            return true;
        } else if (text.equalsIgnoreCase("hd2")) {
            Service.gI().sendThongBaoOK(player, ChatGlobalService.SYSTEM_NOTICE_DETAIL_2);
            return true;
        } else if (text.equalsIgnoreCase("hd")) {
            Service.gI().sendThongBaoOK(player,
                    ChatGlobalService.SYSTEM_NOTICE_DETAIL_1 + "\n\n" + ChatGlobalService.SYSTEM_NOTICE_DETAIL_2);
            return true;
        }
        if (player.isAdmin()) {
            if (text.equals("code")) {
                GiftCodeManager.gI().checkInfomationGiftCode(player);
                return true;
            } else if (text.equals("a")) {
                BossManager.gI().showListBoss(player);
                return true;
            } else if (text.equals("bot")) {
                boolean enabled = !NewBot.gI().isTrainingBotsActive();
                ServerManager.setPanelBotMode(enabled, false);
                Service.gI().sendThongBaoOK(player, enabled ? "Đã bật bot đánh quái" : "Đã tắt bot đánh quái");
                return true;
            } else if (text.equals("mapboss2")) {
                BossManager.gI().showListBoss(player);
                return true;
            }else if (text.equals("ls")) {
            try (Connection con = DBConnecter.getConnectionServer()) {
                Manager.SHOPS = ShopDAO.getShops(con);
                Service.gI().sendThongBaoOK(player, "Update Shop thành công!");
                System.out.println("Update Shop: " + Manager.SHOPS.size());
            } catch (Exception e) {
                e.printStackTrace();
                Service.gI().sendThongBaoOK(player, "Update Shop thất bại!");
            }
            return true;

            } else if (text.equals("hsk")) {
                Service.gI().releaseCooldownSkill(player);
                return true;
            }else if (text.equals("lg")) {
                GiftCodeManager.gI().loadGiftCodeFromDB();
                Service.gI().sendThongBaoOK(player, "Reload giftcode thành công!");
                return true;

            }else if (text.startsWith("exp")) {
                try {
                    int rate = Integer.parseInt(text.replace("exp", "").trim());

                    if (rate <= 0) {
                        Service.gI().sendThongBao(player, "Tỉ lệ phải > 0!");
                        return true;
                    }
                    Manager.RATE_EXP_SERVER = rate;
                    Service.gI().sendThongBaoOK(player,
                            "EXP server: x" + rate);

                    Service.gI().sendThongBaoAllPlayer(
                            "EXP server hiện tại: x" + rate);

                } catch (NumberFormatException e) {
                    Service.gI().sendThongBao(player, "Sai định dạng! Ví dụ: exp2 hoặc exp 2");
                } catch (Exception e) {
                    e.printStackTrace();
                    Service.gI().sendThongBao(player, "Lỗi!");
                }
                return true;
            }
             else if (text.startsWith("upp")) {
                try {
                    long power = Long.parseLong(text.substring(3));
                    Service.gI().addSMTN(player.pet, (byte) 2, power, false);
                    return true;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            } else if (text.startsWith("up")) {
                long power = Long.parseLong(text.substring(2));
                Service.gI().addSMTN(player, (byte) 2, power, false);
                return true;
            } else if (text.equals("battu")) {
                if (player.isBattu) {
                    player.isBattu = false;
                } else {
                    player.isBattu = true;
                }
                Service.gI().sendThongBao(player, "Bất tử" + (player.isBattu ? ": ON" : ": OFF"));
                return true;
            } else if (text.startsWith("dt")) {
                try {
                    long power = Long.parseLong(text.replaceAll("upp", ""));
                    Service.gI().addSMTN(player.pet, (byte) 2, power, false);
                    return true;
                } catch (Exception e) {
                }
         
            } else if (text.equals("nm")) {
                SummonDragonNamek.gI().summonNamec(player);
                return true;
        
            }else if (text.equals("c")) {
    Input.gI().createFormBuffVND(player);
    return true;

            } else if (text.equals("b")) {
                if (player.isAdmin()) {
                  NpcService.gI().createMenuConMeo(player, ConstNpc.MENU_ADMIN, 16103,
        "|7|Ngọc Rồng Kame\n"
                + "|0|Time start: " + ServerManager.timeStart + "\nCó: "
                + Client.gI().getPlayers().size()
                + " Người chơi\n Sessions: " + SessionManager.gI().getNumSession()
                + ", Thread: " + Thread.activeCount()
                + "\n" + SystemMetrics.ToString(),
        "Đệ tử", "Bảo trì", "Tìm kiếm"
        // , "Gọi Boss",
        // "Buff VND",
      ,  "Buff\nhộp thư");
                } else {
                    PlayerService.gI().banPlayer(player);
                    Logger.warning("Player " + player.name + " có dấu hiệu bug admin tiến hành band acc");
                }
                return true;
            } else if (text.startsWith("m")) {
                String mapIdStr = text.replace("m", "").trim();
                int mapId = Integer.parseInt(mapIdStr);
                ChangeMapService.gI().changeMapInYard(player, mapId, -1, -1);
                return true;
            }
            if (text.startsWith("dm")) {
                try {
                    long dameg = Integer.parseInt(text.replaceAll("dm", ""));
                    player.nPoint.dameg = dameg;
                    Service.gI().point(player);
                    return true;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
            if (text.startsWith("hp")) {
                try {
                    long hpg = Integer.parseInt(text.replaceAll("hp", ""));
                    player.nPoint.hpg = hpg;
                    Service.gI().point(player);
                    return true;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
            if (text.startsWith("ki")) {
                long mpg = Long.parseLong(text.replaceAll("ki", ""));
                player.nPoint.mpg = mpg;
                Service.gI().point(player);
                return true;
            }
            if (text.startsWith("def")) {
                try {
                    int defg = Integer.parseInt(text.replaceAll("def", ""));
                    player.nPoint.defg = defg;
                    Service.gI().point(player);
                    return true;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
           
            if (text.startsWith("cr")) {
                try {
                    int critg = Integer.parseInt(text.replaceAll("cr", ""));
                    player.nPoint.critg = critg;
                    Service.gI().point(player);
                    return true;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
            if (text.startsWith("n")) {
                try {
                    int idTask = Integer.parseInt(text.replaceAll("n", ""));
                    player.playerTask.taskMain.id = idTask - 1;
                    player.playerTask.taskMain.index = 0;
                    TaskService.gI().sendNextTaskMain(player);
                    return true;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
          
            if (text.startsWith("gender_")) {
                byte idGender = Byte.parseByte(text.replaceAll("gender_", ""));
                player.gender = idGender;
                return true;
            }
            if (text.startsWith("i")) {
                String[] parts = text.split(" ");
                if (parts.length >= 3) {
                    short id = Short.parseShort(parts[1]);
                    int quantity = Integer.parseInt(parts[2]);
                    Item item = ItemService.gI().createNewItem(id, quantity);
                    List<Item.ItemOption> ops = ItemService.gI().getListOptionItemShop((short) id);
                    if (!ops.isEmpty()) {
                        item.itemOptions = ops;
                    }
                    InventoryService.gI().addItemBag(player, item);
                    InventoryService.gI().sendItemBag(player);
                    Service.gI().sendThongBao(player,
                            "Lụm " + item.template.name + " [" + item.template.id + "] Thành Công !");
                    return true;
                } else {
                    Service.gI().sendThongBao(player, "Lỗi");
                    return true;
                }
            } 
            else if (text.equals("item")) {
                Input.gI().createFormGiveItem(player);
                return true;
            } else if (text.equals("getitem")) {
                Input.gI().createFormGetItem(player);
                return true;
            } else if (text.equals("d")) {
                Service.gI().setPos(player, player.location.x, player.location.y + 10);
                return true;
            }
        }
        if (text.startsWith("ten con la "))

        {
            PetService.gI().changeNamePet(player, text.replaceAll("ten con la ", ""));
        }

        if (player.pet != null) {
            switch (text) {
                case "di theo", "follow" ->
                    player.pet.changeStatus(Pet.FOLLOW);
                case "bao ve", "protect" ->
                    player.pet.changeStatus(Pet.PROTECT);
                case "tan cong", "attack" ->
                    player.pet.changeStatus(Pet.ATTACK);
                case "ve nha", "go home" ->
                    player.pet.changeStatus(Pet.GOHOME);
                case "bien hinh" ->
                    player.pet.transform();
            }
        }
        return false;
    }

}

