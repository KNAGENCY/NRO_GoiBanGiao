package boss.boss_manifest.Wolves;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossesData;
import item.Item;
import item.Item.ItemOption;
import map.ItemMap;
import player.Player;
import services.ItemService;
import services.Service;
import services.TaskService;
import utils.Util;

public class Baconsoi extends Boss {

    private long st;

    public Baconsoi() throws Exception {
        super(BossID.BA_CON_SOI, false, true, BossesData.SOI_1, BossesData.SOI_2, BossesData.SOI_3);
    }

 @Override
public void reward(Player plKill) {
        TaskService.gI().checkDoneTaskKillBoss(plKill, this);
        this.pointBoss(plKill);
    // Tỉ lệ 10% rơi Chân mệnh (10/100)
    if (Util.isTrue(10, 100)) { 
        // Tạo vật phẩm ID 1733
        Item item = ItemService.gI().createNewItem((short) 1733);

        // Thêm các Option cố định theo yêu cầu
        item.itemOptions.add(new ItemOption(50, 2));   // Sức đánh %
        item.itemOptions.add(new ItemOption(77, 3));   // HP %
        item.itemOptions.add(new ItemOption(103, 3));  // KI %
        item.itemOptions.add(new ItemOption(5, 2));    // Sức đánh % (Option 5)

        // Drop ra map - chỉ plKill (người hạ boss) nhặt được
        ItemMap it = new ItemMap(
                this.zone,
                item,
                this.location.x,
                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                plKill.id); // Gán ID người giết boss
        Service.gI().dropItemMap(this.zone, it);
        
    } else { 
        // 90% còn lại: Rơi 10 viên ID 1718 rải đều
        for (int i = 0; i < 10; i++) {
            // Tính toán vị trí x ngẫu nhiên xung quanh Boss (trong khoảng -50 đến +50 pixel)
            int xDrop = this.location.x + Util.nextInt(-50, 50);
            
            // Tạo vật phẩm ID 1718
            Item itemNl = ItemService.gI().createNewItem((short) 1718);
            
            // Drop ra map - Gán ID người sở hữu là -1 để ai cũng nhặt được
            ItemMap it = new ItemMap(
                    this.zone,
                    itemNl,
                    xDrop,
                    this.zone.map.yPhysicInTop(xDrop, this.location.y - 24),
                    -1); // -1 nghĩa là vật phẩm tự do (Public)
            
            Service.gI().dropItemMap(this.zone, it);
        }
    }
}

    @Override
    public void joinMap() {
        super.joinMap(); // To change body of generated methods, choose Tools | Templates.
        st = System.currentTimeMillis();
    }

    @Override
    public void autoLeaveMap() {
        if (Util.canDoWithTime(st, 900000)) {
            this.leaveMapNew();
        }
        if (this.zone != null && this.zone.getNumOfPlayers() > 0) {
            st = System.currentTimeMillis();
        }
    }

}
