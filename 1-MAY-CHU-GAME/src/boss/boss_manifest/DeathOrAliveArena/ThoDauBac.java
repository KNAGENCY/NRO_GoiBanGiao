package boss.boss_manifest.DeathOrAliveArena;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.BossID;
import boss.BossManager;
import boss.BossesData;
import static boss.BossType.PHOBAN;
import player.Player;

public class ThoDauBac extends DeathOrAliveArena {

    public ThoDauBac(Player player) throws Exception {
        super(PHOBAN, BossID.THO_DAU_BAC, BossesData.THO_DAU_BAC);
        BossManager.gI().addBoss(this);
        this.playerAtt = player;
    }
}
