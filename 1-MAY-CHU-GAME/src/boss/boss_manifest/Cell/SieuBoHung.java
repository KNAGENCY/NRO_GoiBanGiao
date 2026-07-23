package boss.boss_manifest.Cell;

import GiaCatLuongManager.Functions;
import consts.ConstPlayer;
import boss.Boss;
import boss.BossID;
import boss.BossManager;
import boss.BossStatus;
import boss.BossesData;
import map.ItemMap;
import player.Player;
import server.Manager;
import services.*;
import utils.Util;

public class SieuBoHung extends Boss {

    private long st;
    public boolean callCellCon;
    private long lastTimeHapThu;
    private int timeHapThu;

    private final String text[] = { "Thưa quý vị và các bạn, đây đúng là trận đấu trời long đất lở",
            "Vượt xa mọi dự đoán của chúng tôi",
            "Eo ơi toàn thân lão Xên bốc cháy kìa" };
    private long lastTimeChat;
    private long lastTimeMove;
    private int indexChat = 0;

    public SieuBoHung() throws Exception {
        super(BossID.SIEU_BO_HUNG, false, true, BossesData.SIEU_BO_HUNG_1, BossesData.SIEU_BO_HUNG_2);
    }

    @Override
    protected void resetBase() {
        super.resetBase();
        this.callCellCon = false;
    }

    public void callCellCon() {
        new Thread(() -> {
            try {
                this.changeStatus(BossStatus.AFK);
                this.changeToTypeNonPK();
                this.recoverHP();
                this.callCellCon = true;
                this.chat("Hãy đấu với 7 đứa con của ta, chúng đều là siêu cao thủ");
                Functions.sleep(2000);
                this.chat("Cứ chưởng tiếp đi haha");
                Functions.sleep(2000);
                this.chat("Liệu mà giữ mạng đấy");
                Functions.sleep(2000);
                if (this.bossAppearTogether != null && this.currentLevel >= 0
                        && this.currentLevel < this.bossAppearTogether.length) {
                    Boss[] arr = this.bossAppearTogether[this.currentLevel];
                    if (arr != null) {
                        for (Boss boss : arr) {
                            if (boss == null)
                                continue;
                            int id = (int) boss.id;
                            if (id == BossID.XEN_CON_1 || id == BossID.XEN_CON_2 || id == BossID.XEN_CON_3 ||
                                    id == BossID.XEN_CON_4 || id == BossID.XEN_CON_5 || id == BossID.XEN_CON_6 ||
                                    id == BossID.XEN_CON_7) {
                                boss.changeStatus(BossStatus.RESPAWN);
                            }
                        }
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();
    }

    public void recoverHP() {
        PlayerService.gI().hoiPhuc(this, this.nPoint.hpMax, 0);
    }

    @Override
    public void reward(Player plKill) {
        this.pointBoss(plKill);
        TaskService.gI().checkDoneTaskKillBoss(plKill, this);

        int rd = Util.nextInt(1, 100);
        int itemId;

        if (rd <= 30) {
            itemId = 14;
        } else if (rd <= 60) {
            itemId = 15;
        } else {
            itemId = 16;
        }

        ItemMap it = new ItemMap(
                this.zone,
                itemId,
                1,
                this.location.x + Util.nextInt(-15, 15),
                this.zone.map.yPhysicInTop(this.location.x, this.location.y - 24),
                plKill.id);
        Service.gI().dropItemMap(this.zone, it);
    }

    @Override
    public void active() {
        if (this.typePk == ConstPlayer.NON_PK) {
            this.changeToTypePK();
        }
        this.attack();
    }

    @Override
    public synchronized long injured(Player plAtt, long damage, boolean piercing, boolean isMobAttack) {
        if (prepareBom) {
            return 0;
        }
        if (!this.callCellCon && damage >= this.nPoint.hp) {
            if (this.currentLevel == 0) {
                this.callCellCon();
                return 0;
            }
        }
        if (!this.isDie()) {
            if (!piercing && Util.isTrue(this.nPoint.tlNeDon, 1000)) {
                this.chat("Xí hụt");
                return 0;
            }
            damage = this.nPoint.subDameInjureWithDeff(damage / 3);
            if (!piercing && effectSkill.isShielding) {
                if (damage > nPoint.hpMax) {
                    EffectSkillService.gI().breakShield(this);
                }
                damage = damage / 4;
            }
            this.nPoint.subHP(damage);
            if (isDie()) {
                if (this.data != null && this.currentLevel >= this.data.length - 1) {
                    try {
                        reward(plAtt);
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                    this.leaveMapNew();
                    new Thread(() -> {
                        try {
                            Thread.sleep(30 * 60 * 1000L);
                            try {
                                BossManager.gI().addBoss(new SieuBoHung());
                            } catch (Exception ex) {
                                ex.printStackTrace();
                            }
                        } catch (InterruptedException ignored) {
                        }
                    }).start();
                } else {
                    setBom(plAtt);
                }
                return 0;
            }
            return damage;
        } else {
            return 0;
        }
    }

    @Override
    public void joinMap() {
        super.joinMap();
        st = System.currentTimeMillis();
    }

    @Override
    public void autoLeaveMap() {
        this.mc();
        if (this.bossStatus == BossStatus.AFK) {
            this.changeStatus(BossStatus.ACTIVE);
        }
        if (Util.canDoWithTime(st, 900000)) {
            this.leaveMapNew();
        }
        if (this.zone != null && this.zone.getNumOfPlayers() > 0) {
            st = System.currentTimeMillis();
        }
    }

    public void mc() {
        Player mc = zone.getNpc();
        if (mc != null) {
            if (Util.canDoWithTime(lastTimeChat, 3000)) {
                String textchat = text[indexChat];
                Service.gI().chat(mc, textchat);
                indexChat++;
                if (indexChat == text.length) {
                    indexChat = 0;
                    lastTimeChat = System.currentTimeMillis() + 7000;
                } else {
                    lastTimeChat = System.currentTimeMillis();
                }
            }

            if (Util.canDoWithTime(lastTimeMove, 15000)) {
                if (Util.isTrue(2, 3)) {
                    int x = this.location.x + Util.nextInt(-100, 100);
                    int y = x > 156 && x < 611 ? 288 : 312;
                    PlayerService.gI().playerMove(mc, x, y);
                }
                lastTimeMove = System.currentTimeMillis();
            }
        }
    }

}
