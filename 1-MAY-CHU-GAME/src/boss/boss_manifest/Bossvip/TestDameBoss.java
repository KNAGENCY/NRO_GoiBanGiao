package boss.boss_manifest.Bossvip;

import boss.Boss;
import boss.BossData;
import boss.BossID;
import consts.ConstPlayer;
import map.Map;
import map.Zone;
import player.Player;
import services.MapService;
import services.Service;
import services.func.ChangeMapService;
import skill.Skill;
import utils.Util;

public class TestDameBoss extends Boss {

    private static final long FULL_HP = 1_999_999_999L;
    private static final int MAP_ID = 29;
    private static final int ZONE_ID = 0;
    private static final int SPAWN_X = 1188;
    private static final int SPAWN_Y = 360;
    private static final long HEAL_WHEN_BELOW_HP = 500_000_000L;
    private static final long CHAT_INTERVAL = 5_000L;
    private static final String[] CHAT_TEXTS = {
            "Anh hẹn em Pickleball",
            "Ta vờn nhau Pickleball",
            "Em focus vào tâm trí anh chưa ?",
            "vượt mức Pickleball",
            "Anh đập banh, đập banh"
    };

    private long lastChatTime;

    public TestDameBoss() throws Exception {
        super(BossID.TEST_DAME, false, false, new BossData(
                "Đánh Là Chó",
                ConstPlayer.XAYDA,
                new short[] { 400, 401, 402, -1, -1, -1 },
                1,
                new long[] { FULL_HP },
                new int[] { MAP_ID },
                new int[][] { { Skill.TAI_TAO_NANG_LUONG, 1, 10000 } },
                new String[] {},
                new String[] {},
                new String[] {},
                1));
    }

    @Override
    public void joinMap() {
        Map map = MapService.gI().getMapById(MAP_ID);
        if (map == null || map.zones.size() <= ZONE_ID) {
            this.changeStatus(boss.BossStatus.REST);
            return;
        }
        Zone targetZone = map.zones.get(ZONE_ID);
        this.zone = targetZone;
        ChangeMapService.gI().changeMap(this, targetZone, SPAWN_X, SPAWN_Y);
        Service.gI().sendFlagBag(this);
        this.changeStatus(boss.BossStatus.ACTIVE);
        this.lastChatTime = System.currentTimeMillis();
    }

    @Override
    public void active() {
        if (this.typePk == ConstPlayer.NON_PK) {
            this.changeToTypePK();
        }
        if (this.nPoint.hp < HEAL_WHEN_BELOW_HP) {
            this.nPoint.hp = this.nPoint.hpMax;
            Service.gI().Send_Info_NV(this);
        }
        if (Util.canDoWithTime(this.lastChatTime, CHAT_INTERVAL)) {
            this.lastChatTime = System.currentTimeMillis();
            this.chat(CHAT_TEXTS[Util.nextInt(0, CHAT_TEXTS.length - 1)]);
        }
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (plAtt == null || this.isDie()) {
            return 0;
        }
        this.nPoint.subHP(damage);
        if (this.nPoint.hp <= 0) {
            this.nPoint.hp = 1;
        }
        return damage;
    }

    @Override
    public void reward(Player plKill) {
    }
}
