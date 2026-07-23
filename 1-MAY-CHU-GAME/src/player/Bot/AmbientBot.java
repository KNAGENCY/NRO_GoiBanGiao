package player.Bot;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import map.Map;
import map.WayPoint;
import map.Zone;
import npc.Npc;
import server.Manager;
import services.func.ChangeMapService;
import utils.Util;

public final class AmbientBot {
    private static final int ROAMING_COUNT = 30;
    private static final int[] NPC_MAPS = {0, 7, 14, 5, 19, 154, 42, 43, 44, 20, 13, 24, 25, 26};
    private static final long VISIBLE_TIME = 120_000L;
    private static final long HIDDEN_TIME = 60_000L;

    private final Bot bot;
    private final boolean npcSitter;
    private Zone homeZone;
    private WayPoint targetGate;
    private long nextMove;
    private long stateTime;
    private boolean visible = true;
    private final int moveDelay;
    private final int moveStep;

    private AmbientBot(Bot bot, boolean npcSitter) {
        this.bot = bot;
        this.npcSitter = npcSitter;
        this.moveDelay = Util.nextInt(400, 550);
        this.moveStep = Util.nextInt(40, 60);
        this.stateTime = System.currentTimeMillis();
    }

    public static void spawnAll(Set<String> usedNames) {
        List<Map> roamingMaps = eligibleRoamingMaps();
        for (int i = 0; i < ROAMING_COUNT && !roamingMaps.isEmpty(); i++) {
            Map map = roamingMaps.get(Util.nextInt(0, roamingMaps.size() - 1));
            spawn(map, false, usedNames);
        }
        for (int mapId : NPC_MAPS) {
            Map map = findMap(mapId);
            if (map != null) {
                spawn(map, true, usedNames);
            }
        }
    }

    private static void spawn(Map map, boolean npcSitter, Set<String> usedNames) {
        if (map.zones == null || map.zones.isEmpty() || map.zones.get(0) == null) return;
        Zone zone = map.zones.get(0);
        Bot bot = NewBot.gI().createAmbientBot(zone, usedNames);
        if (bot == null) return;
        AmbientBot behavior = new AmbientBot(bot, npcSitter);
        bot.ambient = behavior;
        BotManager.gI().bot.add(bot);
        if (npcSitter) behavior.placeNearRandomNpc(map);
        bot.joinAssignedZone(zone);
        if (!npcSitter) behavior.chooseGate();
    }

    public void update() {
        if (npcSitter) updateNpcSitter();
        else updateRoaming();
    }

    private void updateRoaming() {
        long now = System.currentTimeMillis();
        if (bot.zone == null || now < nextMove) return;
        nextMove = now + moveDelay;
        bot.nPoint.mp = bot.nPoint.mpMax;
        if (targetGate == null || !isEligibleMap(targetGate.goMap)) {
            chooseGate();
            if (targetGate == null) return;
        }
        int tx = (targetGate.minX + targetGate.maxX) / 2;
        int ty = targetGate.maxY;
        double dx = tx - bot.location.x;
        if (Math.abs(dx) <= moveStep) {
            Map next = findMap(targetGate.goMap);
            if (next != null && next.zones != null && !next.zones.isEmpty()) {
                bot.moveTo(tx, ty);
                bot.joinAssignedZone(next.zones.get(0));
                bot.moveTo(targetGate.goX, targetGate.goY);
            }
            targetGate = null;
            chooseGate();
            return;
        }
        bot.moveTo(bot.location.x + (dx > 0 ? moveStep : -moveStep), ty);
    }

    private void chooseGate() {
        if (bot.zone == null || bot.zone.map.wayPoints == null) return;
        List<WayPoint> gates = new ArrayList<>();
        for (WayPoint gate : bot.zone.map.wayPoints) {
            if (gate != null && isEligibleMap(gate.goMap) && gate.goMap != bot.zone.map.mapId
                    && hasEligibleExit(findMap(gate.goMap))) gates.add(gate);
        }
        if (!gates.isEmpty()) targetGate = gates.get(Util.nextInt(0, gates.size() - 1));
    }

    private void updateNpcSitter() {
        long now = System.currentTimeMillis();
        long elapsed = now - stateTime;
        if (visible && elapsed >= VISIBLE_TIME) {
            homeZone = bot.zone;
            ChangeMapService.gI().exitMap(bot);
            visible = false;
            stateTime = now;
        } else if (!visible && elapsed >= HIDDEN_TIME && homeZone != null) {
            placeNearRandomNpc(homeZone.map);
            bot.joinAssignedZone(homeZone);
            visible = true;
            stateTime = now;
        }
    }

    private void placeNearRandomNpc(Map map) {
        if (map == null || map.npcs == null) return;
        List<Npc> npcs = new ArrayList<>();
        for (Npc npc : map.npcs) if (npc != null && npc.tempId != 6) npcs.add(npc);
        if (npcs.isEmpty()) return;
        Npc npc = npcs.get(Util.nextInt(0, npcs.size() - 1));
        bot.location.x = npc.cx + Util.nextInt(-45, 45);
        bot.location.y = npc.cy;
    }

    private static boolean isEligibleMap(int id) {
        Map map = findMap(id);
        return map != null && map.type == 0 && (id <= 110 || id == 122 || id == 123 || id == 124)
                && id != 78 && id != 51 && id != 49;
    }

    private static List<Map> eligibleRoamingMaps() {
        List<Map> result = new ArrayList<>();
        for (Map map : Manager.MAPS) if (map != null && isEligibleMap(map.mapId) && hasEligibleExit(map)) result.add(map);
        return result;
    }

    private static boolean hasEligibleExit(Map map) {
        if (map == null || map.wayPoints == null) return false;
        for (WayPoint gate : map.wayPoints) {
            if (gate != null && gate.goMap != map.mapId && isEligibleMap(gate.goMap)) return true;
        }
        return false;
    }

    private static Map findMap(int id) {
        for (Map map : Manager.MAPS) if (map != null && map.mapId == id) return map;
        return null;
    }
}
