package boss.boss_manifest.GinyuForce;

import java.util.Calendar;

/*
 *
 *
 * @author GiaCatLuong
 */

import boss.Boss;
import boss.BossID;
import boss.BossStatus;
import boss.BossesData;
import map.ItemMap;
import player.Player;
import services.EffectSkillService;
import services.Service;
import services.TaskService;
import utils.Util;

public class TDT extends Boss {

    private long st;

    private long lastBodyChangeTime;

    public TDT() throws Exception {
        super(BossID.TIEU_DOI_TRUONG, false, true, BossesData.TIEU_DOI_TRUONG);
    }

    private void bodyChangePlayerInMap() {
        if (this.zone != null) {
            for (Player pl : this.zone.getPlayers()) {
                if (Util.isTrue(5, 10) && pl.effectSkill != null && !pl.effectSkill.isBodyChangeTechnique) {
                    EffectSkillService.gI().setIsBodyChangeTechnique(pl);
                }
            }
        }
    }

    @Override
    public void moveTo(int x, int y) {
        if (this.currentLevel == 1) {
            return;
        }
        super.moveTo(x, y);
    }
@Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {

        // Kiểm tra nhiệm vụ 22
        boolean isTask22 = plAtt != null
                && plAtt.playerTask != null
                && plAtt.playerTask.taskMain != null
                && plAtt.playerTask.taskMain.id == 22;

        // Không nhiệm vụ 22 + trong khung giờ => không gây damage
        if (!isTask22 && isInTime()) {
            Service.gI().sendThongBao(plAtt,
                    "Bạn chỉ có thể đánh boss ngoài khung giờ hoặc khi đang làm nhiệm vụ!");
            return 0;
        }

        // Giới hạn sát thương
        if (damage > 250_000_000) {
            damage = 250_000_000;
        }

        return super.injured(plAtt, damage, piercing, isMobAttack);
    }

 @Override
public void reward(Player plKill) {
    // Tặng điểm cho người hạ boss
  //  this.pointBoss(plKill);

    TaskService.gI().checkDoneTaskKillBoss(plKill, this);

    // 50% rơi item 1194, chia thành 2 vị trí, mỗi vị trí 1 item
    if (Util.isTrue(50, 100)) {
        for (int i = 0; i < 2; i++) {
            ItemMap it = new ItemMap(this.zone, 1194, 1,
                    this.location.x + Util.nextInt(-30, 30),
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    plKill.id);
            Service.gI().dropItemMap(this.zone, it);
        }
    } 
    // 50% còn lại rơi item 751, 10 item, mọi người nhặt được
    else {
        for (int i = 0; i < 10; i++) {
            ItemMap it = new ItemMap(this.zone, 750, 1,
                    this.location.x + Util.nextInt(-60, 60),
                    this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                    -1);
            Service.gI().dropItemMap(this.zone, it);
        }
    }
}



    @Override
    protected void notifyJoinMap() {
        if (this.currentLevel == 1) {
            return;
        }
        super.notifyJoinMap();
    }

    @Override
    public void attack() {
        if (Util.canDoWithTime(lastBodyChangeTime, 10000)) {
            bodyChangePlayerInMap();
            this.chat("Úm ba la xì bùa");
            this.lastBodyChangeTime = System.currentTimeMillis();
        }
        super.attack();
    }

    @Override
    public void joinMap() {
        super.joinMap();
        st = System.currentTimeMillis();
    }

    @Override
    public void doneChatS() {
        this.changeStatus(BossStatus.AFK);
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
    
    // ================== KHUNG GIỜ ==================
    private boolean isInTime() {
        Calendar c = Calendar.getInstance();
        int hour = c.get(Calendar.HOUR_OF_DAY);
        int minute = c.get(Calendar.MINUTE);

        int time = hour * 60 + minute;

        // Sáng 08:30 - 10:00
        if (time >= (8 * 60 + 30) && time <= (10 * 60)) return true;

        // Chiều 14:30 - 17:00
        if (time >= (14 * 60 + 30) && time <= (17 * 60)) return true;

        // Tối 20:30 - 22:00
        if (time >= (20 * 60 + 30) && time <= (22 * 60)) return true;

        return false;
    }
}
