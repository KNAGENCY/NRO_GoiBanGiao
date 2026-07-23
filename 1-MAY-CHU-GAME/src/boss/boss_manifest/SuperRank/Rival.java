package boss.boss_manifest.SuperRank;

import boss.BossData;
import boss.BossManager;
import boss.BossStatus;
import static boss.BossType.PHOBAN;
import java.util.ArrayList;
import java.util.List;
import player.Player;
import skill.Skill;
import utils.Util;

public class Rival extends SuperRank {

    public Rival(Player player, Player rival) throws Exception {
        super(PHOBAN, Util.createIdBossClone((int) rival.id), bossData(rival));

        this.player = rival; // hình dạng
        this.playerAtt = player; // mục tiêu

        this.currentLevel = 0;
        this.isBoss = true;
        this.typePk = 3;

        // 🔥 FIX CỐT LÕI CHO CORE NÀY
        this.changeStatus(BossStatus.ACTIVE);
        this.initSkill(); // load skill từ BossData
        BossManager.gI().addBoss(this);

    }

    // 🔥 ép AI luôn đánh player thật
    @Override
    public Player getPlayerAttack() {
        return playerAtt;
    }

    // 🔥 outfit lấy từ player
    @Override
    public short getHead() {
        return player.getHead();
    }

    @Override
    public short getBody() {
        return player.getBody();
    }

    @Override
    public short getLeg() {
        return player.getLeg();
    }

    // ===================== BossData =====================
    public static BossData bossData(Player player) {
        List<Skill> skillList = new ArrayList<>();

        for (Skill skill : player.playerSkill.skills) {
            if (skill.point > 0
                    && skill.template.id != Skill.TU_SAT
                    && skill.template.id != Skill.QUA_CAU_KENH_KHI
                    && skill.template.id != Skill.MAKANKOSAPPO
                    && skill.template.id != Skill.TRI_THUONG) {
                skillList.add(skill);
            }
        }

        int[][] skillTemp = new int[skillList.size()][3];
        for (int i = 0; i < skillList.size(); i++) {
            Skill s = skillList.get(i);
            skillTemp[i][0] = s.template.id;
            skillTemp[i][1] = s.point;
            skillTemp[i][2] = s.coolDown;
        }

        return new BossData(
                player.name,
                player.gender,
                new short[] {
                        player.getHead(),
                        player.getBody(),
                        player.getLeg(),
                        player.getFlagBag(),
                        player.getAura(),
                        player.getEffFront()
                },
                player.nPoint.dameg,
                new long[] { player.nPoint.hpg },
                new int[] { 113 },
                skillTemp,
                new String[] {},
                new String[] {},
                new String[] {},
                60);
    }
}
