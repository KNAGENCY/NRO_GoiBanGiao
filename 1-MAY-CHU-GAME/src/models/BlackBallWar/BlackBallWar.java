package models.BlackBallWar;

import jdbc.daos.NDVSqlFetcher;
import jdbc.daos.PlayerDAO;
import player.Player;
import services.MapService;
import services.Service;
import services.func.ChangeMapService;
import utils.TimeUtil;
import utils.Util;

import map.Zone;
import server.Client;

public class BlackBallWar {

    public static final int TIME_CAN_PICK_BLACK_BALL_AFTER_DROP = 5000;

    public static final byte X3 = 3;
    public static final byte X5 = 5;
    public static final byte X7 = 7;

    public static final int COST_X3 = 100000000;
    public static final int COST_X5 = 300000000;
    public static final int COST_X7 = 500000000;

    public static final int AVAILABLE = 5;
    private static final int TIME_WIN = 300000;
//    private static final int TIME_WIN = 30;
    public final Zone zone;

    public BlackBallWar(Zone zone) {
        this.zone = zone;
    }

    // ⚠️ ĐƯỢC GAMELOOP GỌI
    public void update(long now) {

        if (!TimeUtil.isBlackBallWarOpen()) {
            zone.finishBlackBallWar = false;
        }

        for (int i = zone.getNumOfPlayers() - 1; i >= 0; i--) {
            updatePlayer(zone.getPlayers().get(i), now);
        }
    }

    private void updatePlayer(Player player, long now) {
        if (player.zone == null
                || !MapService.gI().isMapBlackBallWar(player.zone.map.mapId)) {
            return;
        }

        if (!TimeUtil.isBlackBallWarOpen()) {
            kickOutOfMap(player);
            return;
        }

        if (player.iDMark.isHoldBlackBall()) {
            if (now - player.iDMark.getLastTimeHoldBlackBall() >= TIME_WIN) {
                win(player);
            } else if (Util.canDoWithTime(
                    player.iDMark.getLastTimeNotifyTimeHoldBlackBall(), 10000)) {

                Service.gI().sendThongBao(player,
                        "Cố giữ ngọc thêm "
                        + TimeUtil.getSecondLeft(
                                player.iDMark.getLastTimeHoldBlackBall(),
                                TIME_WIN / 1000)
                        + " giây nữa sẽ thắng");

                player.iDMark.setLastTimeNotifyTimeHoldBlackBall(now);
            }
        }
    }

    private void win(Player player) {
        zone.finishBlackBallWar = true;
        int star = player.iDMark.getTempIdBlackBallHold() - 371;

        player.rewardBlackBall.reward((byte) star);
        Service.gI().sendThongBao(player,
                "Chúc mừng bạn đã giành được Ngọc rồng "
                + star + " sao đen cho bang");

        if (player.clan != null) {
            player.clan.members.forEach(m -> {
                Player p = Client.gI().getPlayer(m.id);
                if (p != null) {
                    p.rewardBlackBall.reward((byte) star);
                } else {
                    Player off = NDVSqlFetcher.loadById(m.id);
                    if (off != null) {
                        off.rewardBlackBall.reward((byte) star);
                        PlayerDAO.updatePlayer(off);
                    }
                }
            });
        }

        kickAllPlayersOutOfMap();
    }

    private void kickOutOfMap(Player player) {
        if (player.cFlag == 8) {
            Service.gI().changeFlag(player, Util.nextInt(1, 7));
        }

        Service.gI().sendThongBao(
                player,
                "Trò chơi tìm ngọc hôm nay đã kết thúc, hẹn gặp lại vào 20h ngày mai"
        );

        ChangeMapService.gI().changeMapBySpaceShip(
                player, player.gender + 24, -1, 250
        );
    }

    private void kickAllPlayersOutOfMap() {
        for (Player p : zone.getPlayers()) {
            kickOutOfMap(p);
        }
    }
}
