package event.event_manifest;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.BossID;
import event.Event;

public class TrungThu extends Event {

    @Override
    public void boss() {
        // createBoss(BossID.KHIDOT, 10);
        createBoss(BossID.NGUYETTHAN, 1000000);
      
    }
}
