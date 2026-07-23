package boss.boss_manifest.Bossvip;

import boss.Boss;
import boss.BossData;
import boss.BossID;
import boss.BossStatus;
import consts.ConstPlayer;
import map.ItemMap;
import player.Player;
import services.Service;
import skill.Skill;

/**
 * Mabuz boss (placed under boss.boss_manifest.Broly)
 * Khi chết rơi item id 568 cho người hạ.
 */
public class Mabuz extends Boss {

    private long st;
    private long lastTimeAttack;

    public Mabuz() throws Exception {
        super(BossID.MABUZ, false, true, new BossData(
                "Ma Bư", // name
                ConstPlayer.XAYDA, // gender (chỉnh nếu cần)
                new short[]{297, 298, 299, -1, -1, -1}, // outfit {head, body, leg, bag, aura, eff}
                5000, // dame base
                new long[]{150000000}, // hp base (mảng)
                new int[]{29}, // map join
                new int[][]{
                      {Skill.KAMEJOKO, 1, 2000},
                        {Skill.LIEN_HOAN, 1, 1000}, {Skill.MASENKO, 1, 2000}, {Skill.THAI_DUONG_HA_SAN,1,50000}, {Skill.KHIEN_NANG_LUONG, 1, 50000}
                }, // skills
                new String[]{}, // text chat 1
                new String[]{"|-1|Ta sẽ hủy diệt tất cả!", "|-1|Các ngươi không có cửa đâu!"}, // text chat 2
                new String[]{"|-1|Ta sẽ quay lại..."}, // text chat 3
                600 // type appear
        ));
    }

    @Override
    public void active() {
        super.active();
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (!this.isDie()) {
            // giới hạn dame tối đa = 1.000.000
            if (damage > 500_000_000) {
                damage = 500_000_000;
            }
            this.nPoint.subHP(damage);
            if (this.isDie()) {
                this.setDie(plAtt);
                this.die(plAtt);
            }
            return damage;
        }
        return 0;
    }

    @Override
    public void die(Player plKill) {
        this.reward(plKill);
        this.changeStatus(BossStatus.DIE);
    }

    @Override
    public void reward(Player plKill) {
        if (plKill != null && this.zone != null) {
            // rơi item 568, chỉ người hạ mới nhặt được
            ItemMap itemMap = new ItemMap(
                    this.zone,
                    568,
                    1,
                    this.location.x,
                    this.location.y,
                    plKill.id
            );
            Service.gI().dropItemMap(this.zone, itemMap);
            Service.gI().sendThongBaoOK(plKill, "Bạn nhận được Trứng Mabu!");
        }
    }
}
