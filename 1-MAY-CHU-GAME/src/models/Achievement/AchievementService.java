package models.Achievement;

import consts.ConstAchievement;
import item.Item;
import mob.Mob;
import models.Template.AchievementTemplate;
import network.Message;
import player.Player;
import server.Manager;
import services.InventoryService;
import services.ItemService;
import services.Service;
import skill.Skill;
import utils.Util;

public class AchievementService {

    private static AchievementService instance;

    public static AchievementService gI() {
        if (instance == null) {
            instance = new AchievementService();
        }
        return instance;
    }

    public void openAchievementUI(Player player) {
        Message msg = null;
        try {
            msg = new Message(-76);
            msg.writer().writeByte(0);
            msg.writer().writeByte(Manager.ACHIEVEMENT_TEMPLATE.size());
            for (int i = 0; i < Manager.ACHIEVEMENT_TEMPLATE.size(); i++) {
                AchievementTemplate at = Manager.ACHIEVEMENT_TEMPLATE.get(i);
                msg.writer().writeUTF(at.info1);
                msg.writer().writeUTF(regex(player, at.info2) + " ("
                        + numberToString(player.achievement.getCompleted(i)) + "/"
                        + numberToString(at.maxCount) + ")");
                msg.writer().writeShort(at.money);
                msg.writer().writeBoolean(player.achievement.isFinish(i, at.maxCount));
                msg.writer().writeBoolean(player.achievement.isRecieve(i));
            }
            player.sendMessage(msg);
            player.typeRecvieArchiment = 0;
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    public void confirmAchievement(Player player, byte select) {
        if (player.achievement == null) {
            return;
        }
        if (!player.achievement.canReward(select)) {
            Service.gI().sendThongBao(player, "Không thể thực hiện");
            return;
        }
        AchievementTemplate at = Manager.ACHIEVEMENT_TEMPLATE.get(select);
        player.achievement.reward(select);
        // if (at.money > 0) {
        //     player.inventory.ruby += at.money;
        //     Service.gI().sendMoney(player);
        //     Service.gI().sendThongBao(player, "Bạn vừa nhận được " + at.money + " Hồng ngọc.");
        // }
        if (at.itemId > 0 && at.itemQuantity > 0) {
            Item item = ItemService.gI().createNewItem((short) at.itemId, at.itemQuantity);
            InventoryService.gI().addItemBag(player, item);
            InventoryService.gI().sendItemBag(player);
            Service.gI().sendThongBao(player, "Bạn vừa nhận được " + at.itemQuantity + " " + item.template.name + ".");
        }

        if (select < ConstAchievement.DAC_DAO_THANH_TIEN) {
            player.achievement.done(ConstAchievement.DAC_DAO_THANH_TIEN, 1);
        }

        Message msg = null;
        try {
            msg = new Message(-76);
            msg.writer().writeByte(1);
            msg.writer().writeByte(select);
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    public String numberToString(long num) {
        return num <= 10000 ? num + "" : Util.numberToMoney(num);
    }

    public String regex(Player player, String text) {
        int gen = player.gender;
        return text.replaceAll("%1", gen == 0 ? "Siêu nhân" : gen == 1 ? "Siêu Namếc" : "Siêu Xayda")
                .replaceAll("%2", gen == 0 ? "Bunma" : gen == 1 ? "Dende" : "Appule");
    }

    public void checkDoneTask(Player player, int aId) {
        if (player.isPl() && player.achievement != null) {
            switch (aId) {
                default -> player.achievement.done(aId, 1);
            }
        }
    }

    public void checkDoneTaskKillMob(Player player, Mob mob) {
        try {
            if (mob.lvMob > 0) {
                checkDoneTask(player, ConstAchievement.DANH_BAI_SIEU_QUAI);
            }
            if (mob.type == 4) {
                checkDoneTask(player, ConstAchievement.THO_SAN_THIEN_XA);
            }
            if (mob.tempId == 0) {
                checkDoneTask(player, ConstAchievement.TAP_LUYEN_BAI_BAN);
            }
        } catch (Exception e) {
        }
    }

    public void checkDoneTaskUseSkill(Player player) {
        if (player.isPl()) {
            switch (player.playerSkill.skillSelect.template.id) {
                case Skill.KAMEJOKO, Skill.MASENKO, Skill.ANTOMIC -> {
                    checkDoneTask(player, ConstAchievement.NOI_CONG_CAO_CUONG);
                }
                case Skill.DRAGON, Skill.DEMON, Skill.GALICK, Skill.LIEN_HOAN, Skill.KAIOKEN,
                        Skill.DICH_CHUYEN_TUC_THOI ->
                    {
                    }
                default -> checkDoneTask(player, ConstAchievement.KY_NANG_THANH_THAO);
            }
        }
    }

    public void checkDoneTaskFly(Player player, int length) {
        if (player.isPl() && player.achievement != null) {
            length = Math.abs(length / 10);
            if (length < 10) {
                player.achievement.done(ConstAchievement.KHINH_CONG_THANH_THAO, length);
            }
        }
    }
}
