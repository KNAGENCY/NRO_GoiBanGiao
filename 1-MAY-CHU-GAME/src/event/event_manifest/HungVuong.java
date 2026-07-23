package event.event_manifest;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.BossID;
import event.Event;

public class HungVuong extends Event {

    @Override
    public void boss() {
        createBoss(BossID.THUY_TINH, 1000000);
    }
}
