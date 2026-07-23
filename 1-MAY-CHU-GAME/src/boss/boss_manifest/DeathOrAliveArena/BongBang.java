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

public class BongBang extends DeathOrAliveArena {

    public BongBang(Player player) throws Exception {
        super(PHOBAN, BossID.BONG_BANG, BossesData.BONG_BANG);
        BossManager.gI().addBoss(this);
        this.playerAtt = player;
    }
}
