package server;

import GiaCatLuongManager.Functions;
import jdbc.DBConnecter;
import jdbc.daos.PlayerDAO;
import map.ItemMap;
import player.Player;
import network.SessionManager;
import network.inetwork.ISession;
import server.io.MySession;
import services.Service;
import services.func.ChangeMapService;
import services.func.SummonDragon;
import services.func.TransactionService;
import services.NgocRongNamecService;
import utils.Logger;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
// import models.DragonNamecWar.TranhNgoc;
import models.SuperRank.SuperRankManager;
import models.War.War;
import services.func.SummonDragonNamek;
import utils.Util;

// public class Client implements Runnable {

//     private static Client instance;
    
//     private final Map<Long, Player> players_id = new HashMap<>();
//     private final Map<Integer, Player> players_userId = new HashMap<>();
//     private final Map<String, Player> players_name = new HashMap<>();
//     private final List<Player> players = new ArrayList<>();
public class Client implements Runnable {

    private static Client instance;
    
    private final Map<Long, Player> players_id = new java.util.concurrent.ConcurrentHashMap<>();
    private final Map<Integer, Player> players_userId = new java.util.concurrent.ConcurrentHashMap<>();
    private final Map<String, Player> players_name = new java.util.concurrent.ConcurrentHashMap<>();
    // ← THAY ArrayList thành CopyOnWriteArrayList
    private final List<Player> players = new java.util.concurrent.CopyOnWriteArrayList<>();



    private Client() {
        new Thread(this, "Update Client").start();
    }

    public List<Player> getPlayers() {
        return this.players;
    }

    public Player getPlayerByName(String name) {
        for (Player player : players) {
            if (player.getName().equals(name)) {
                return player;
            }
        }
        return null;
    }

    public static Client gI() {
        if (instance == null) {
            instance = new Client();
        }
        return instance;
    }

    public boolean put(Player player) {
        if (hasDuplicateActivePlayer(player)) {
            banDuplicateAccount(player);
            return false;
        }
        if (!players_id.containsKey(player.id)) {
            this.players_id.put(player.id, player);
        }
        if (!players_name.containsValue(player)) {
            this.players_name.put(player.name, player);
        }
        if (!players_userId.containsValue(player)) {
            this.players_userId.put(player.getSession().userId, player);
        }
        if (!players.contains(player)) {
            this.players.add(player);
        }
        SuperRankManager.gI().rewardIfMissedDaily(player);
        return true;
    }

    private boolean hasDuplicateActivePlayer(Player player) {
        if (player == null || player.getSession() == null) {
            return false;
        }
        MySession session = (MySession) player.getSession();
        if (session.userId <= 0) {
            return false;
        }
        for (Player online : players) {
            if (online == null || online == player || online.isPet || online.isNewPet || online.getSession() == null) {
                continue;
            }
            MySession onlineSession = (MySession) online.getSession();
            if (onlineSession.userId == session.userId) {
                return true;
            }
        }
        return false;
    }

    private void banDuplicateAccount(Player player) {
        if (player == null || player.getSession() == null) {
            return;
        }
        MySession session = (MySession) player.getSession();
        banDuplicateAccount(session.userId, session.uu);
        List<Player> playersCopy = new ArrayList<>(players);
        for (Player online : playersCopy) {
            if (online == null || online.getSession() == null) {
                continue;
            }
            MySession onlineSession = (MySession) online.getSession();
            if (onlineSession.userId == session.userId) {
                Service.gI().sendThongBaoOK(online, "Tai khoan bi khoa do dang nhap trung nhan vat.");
                kickSession(onlineSession);
            }
        }
        Service.gI().sendThongBaoOK(session, "Tai khoan bi khoa do dang nhap trung nhan vat.");
        session.disconnect();
        remove(player);
    }

    private void banDuplicateAccount(int userId, String username) {
        if (userId <= 0) {
            return;
        }
        try {
            DBConnecter.executeUpdate("update account set ban = 1 where id = ?", userId);
            Logger.warning("[DUPLICATE-PLAYER-BAN] userId=" + userId + " username=" + username);
        } catch (Exception e) {
            Logger.logException(Client.class, e);
        }
    }

    private void checkDuplicateActivePlayers() {
        Map<Integer, Player> seen = new HashMap<>();
        List<Player> playersCopy = new ArrayList<>(players);
        for (Player player : playersCopy) {
            if (player == null || player.isPet || player.isNewPet || player.getSession() == null) {
                continue;
            }
            MySession session = (MySession) player.getSession();
            if (session.userId <= 0) {
                continue;
            }
            Player first = seen.putIfAbsent(session.userId, player);
            if (first != null && first != player) {
                banDuplicateAccount(session.userId, session.uu);
                for (Player duplicate : playersCopy) {
                    if (duplicate == null || duplicate.getSession() == null) {
                        continue;
                    }
                    MySession duplicateSession = (MySession) duplicate.getSession();
                    if (duplicateSession.userId == session.userId) {
                        Service.gI().sendThongBaoOK(duplicate, "Tai khoan bi khoa do dang nhap trung nhan vat.");
                        kickSession(duplicateSession);
                    }
                }
                seen.remove(session.userId);
            }
        }
    }
private void remove(MySession session) {
        if (session.player != null) {
            this.remove(session.player);
            session.player.dispose();
        }
        if (session.joinedGame) {
            session.joinedGame = false;
            try {
                DBConnecter.executeUpdate("update account set last_time_logout = ? where id = ?",
                        new Timestamp(System.currentTimeMillis()), session.userId);
            } catch (Exception e) {
                Logger.logException(Client.class, e);
            }
        }
        try {
            if (session.ipAddress != null) {
                Object o = ServerManager.CLIENTS.get(session.ipAddress);
                if (o != null) {
                    int n = Integer.parseInt(String.valueOf(o));
                    if (n <= 1) {
                        ServerManager.CLIENTS.remove(session.ipAddress);
                    } else {
                        ServerManager.CLIENTS.put(session.ipAddress, n - 1);
                    }
                }
            }
        } catch (Exception e) {
            Logger.logException(Client.class, e);
        }
        // ← XÓA DÒNG NÀY: ServerManager.gI().disconnect(session);
    }
    // private void remove(MySession session) {
    //     if (session.player != null) {
    //         this.remove(session.player);
    //         session.player.dispose();
    //     }
    //     if (session.joinedGame) {
    //         session.joinedGame = false;
    //         try {
    //             DBConnecter.executeUpdate("update account set last_time_logout = ? where id = ?",
    //                     new Timestamp(System.currentTimeMillis()), session.userId);
    //         } catch (Exception e) {
    //             Logger.logException(Client.class, e);
    //         }
    //     }
    //     try {
    //         if (session.ipAddress != null) {
    //             Object o = ServerManager.CLIENTS.get(session.ipAddress);
    //             if (o != null) {
    //                 int n = Integer.parseInt(String.valueOf(o));
    //                 if (n <= 1) {
    //                     ServerManager.CLIENTS.remove(session.ipAddress);
    //                 } else {
    //                     ServerManager.CLIENTS.put(session.ipAddress, n - 1);
    //                 }
    //             }
    //         }
    //     } catch (Exception e) {
    //         Logger.logException(Client.class, e);
    //     }
    //     ServerManager.gI().disconnect(session);
    // }

    private void remove(Player player) {
        // Ghi nốt số giây online còn đang chờ trong bộ nhớ của hệ thống Phúc Lợi.
        services.func.PhucLoiService.gI().onPlayerLogout(player);
        this.players_id.remove(player.id);
        this.players_name.remove(player.name);
        this.players_userId.remove(player.getSession().userId);
        this.players.remove(player);
        if (!player.beforeDispose) {
            player.beforeDispose = true;
            if (player.zone != null && player.zone.map.mapId == War.MAP_WAR_ID) {
                player.mapIdBeforeLogout = 0;
            } else {
                player.mapIdBeforeLogout = player.zone.map.mapId;
            }
            // TranhNgoc.gI().removePlayersBlue(player);
            // TranhNgoc.gI().removePlayersRed(player);
            if (player.idNRNM != -1) {
                ItemMap itemMap = new ItemMap(player.zone, player.idNRNM, 1, player.location.x, player.location.y, -1);
                Service.gI().dropItemMap(player.zone, itemMap);
                NgocRongNamecService.gI().pNrNamec[player.idNRNM - 353] = "";
                NgocRongNamecService.gI().idpNrNamec[player.idNRNM - 353] = -1;
                player.idNRNM = -1;
            }
            ChangeMapService.gI().exitMap(player);
            TransactionService.gI().cancelTrade(player);
            if (player.clan != null) {
                player.clan.removeMemberOnline(null, player);
            }
            if (SummonDragon.gI().playerSummonShenron != null
                    && SummonDragon.gI().playerSummonShenron.id == player.id) {
                SummonDragon.gI().isPlayerDisconnect = true;
            }
            if (SummonDragonNamek.gI().playerSummonShenron != null
                    && SummonDragonNamek.gI().playerSummonShenron.id == player.id) {
                SummonDragonNamek.gI().isPlayerDisconnect = true;
            }
            if (player.shenronEvent != null) {
                player.shenronEvent.isPlayerDisconnect = true;
            }
            if (player.mobMe != null) {
                player.mobMe.mobMeDie();
            }
            if (player.pet != null) {
                if (player.pet.mobMe != null) {
                    player.pet.mobMe.mobMeDie();
                }
                ChangeMapService.gI().exitMap(player.pet);
            }
        }
        PlayerDAO.updatePlayer(player);
    }

    public void kickSession(MySession session) {
        if (session != null) {
            session.disconnect();
            this.remove(session);
        }
    }

    public Player getPlayer(long playerId) {
        return this.players_id.get(playerId);
    }

    public Player getRandPlayer() {
        if (this.players.isEmpty()) {
            return null;
        }
        return this.players.get(Util.nextInt(players.size()));
    }

    public Player getPlayerByUser(int userId) {
        return this.players_userId.get(userId);
    }

    public Player getPlayer(String name) {
        return this.players_name.get(name);
    }

    public Player getPlayerByID(int playerId) {
        for (int i = 0; i < players.size(); i++) {
            Player player = players.get(i);
            if (player != null && player.id == playerId) {
                return player;
            }
        }
        return null;
    }

    public void close() {
        Logger.log(Logger.YELLOW, "BEGIN KICK OUT SESSION " + players.size() + "\n");
        while (!players.isEmpty()) {
            Player pl = players.remove(0);
            if (pl != null && pl.getSession() != null) {
                this.kickSession((MySession) pl.getSession());
            }
        }
        Logger.success("SUCCESSFUL\n");
    }

 private void update() {
    List<ISession> sessions = SessionManager.gI().getSessions();

    for (int i = sessions.size() - 1; i >= 0; i--) {

        if (i >= sessions.size()) {
            continue; // list đã bị thay đổi
        }

        ISession s = sessions.get(i);
        if (!(s instanceof MySession)) {
            continue;
        }

        MySession session = (MySession) s;

        if (session == null) {
            sessions.remove(i);
            continue;
        }

        if (session.timeWait > 0) {
            session.timeWait--;
            if (session.timeWait == 0) {
                kickSession(session);
            }
        }
    }
    checkDuplicateActivePlayers();
}


    @Override
    public void run() {
        while (ServerManager.isRunning) {
            long st = System.currentTimeMillis();
            try {
                update();
            } catch (Exception e) {
                e.printStackTrace();
            }
            Functions.sleep(Math.max(1000 - (System.currentTimeMillis() - st), 10));
        }
    }

    public void show(Player player) {
        String txt = "";
        txt += "sessions: " + SessionManager.gI().getNumSession() + "\n";
        txt += "players_id: " + players_id.size() + "\n";
        txt += "players_userId: " + players_userId.size() + "\n";
        txt += "players_name: " + players_name.size() + "\n";
        txt += "players: " + players.size() + "\n";
        Service.gI().sendThongBao(player, txt);
    }
}
