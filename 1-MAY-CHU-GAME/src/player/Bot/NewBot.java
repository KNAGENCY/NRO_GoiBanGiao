package player.Bot;

import clan.Clan;
import item.Item;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import jdbc.DBConnecter;
import map.Map;
import map.Zone;
import services.ItemService;
import services.MapService;
import services.func.ChangeMapService;
import utils.Util;

public class NewBot {
    public static NewBot i;

    private static final int[] TRAINING_MAPS = {
            5, 77, 68, 79, 83, 92, 100, 99, 93, 109,
            2, 3, 9, 11, 16, 17, 70, 69, 71, 76
    };
    private static final int BOTS_PER_MAP = 5;
    private static final long REQUIRED_PLAYER_POWER = 5_000_000_000L;
    private static final long POWER_CHECK_INTERVAL = 5_000L;
    private volatile long lastPowerCheck;
    private volatile boolean qualifiedPlayerExists;
    private static final short[] COSTUME_IDS = {
            525, 550, 549, 461, 1319, 1697, 1654, 1587, 1330,
            1333, 1652, 711, 710, 1320, 1321, 1322, 1323, 1324, 1325, 1328, 879, 1265, 860, 1698, 884, 1700
    };
    private static final short[] ACCESSORY_IDS = { 1254, 1680, 1679, 1231, 1022, 1681 };
    private static final String[] CLAN_NAMES = {
            "Top 1", "Cay Chay", "Dam Me NRO", "2 Anh Em Song Sat", "Dai Vuong",
            "Nguc To", "Tu Tien", "Trum Namek", "Yasuo"
    };

    public static NewBot gI() {
        if (i == null) {
            i = new NewBot();
        }
        return i;
    }

    public synchronized void spawnTrainingBots() {
        if (!hasQualifiedDatabasePlayer()) {
            resetTrainingBots();
            return;
        }
        resetTrainingBots();
        Set<String> usedNames = new HashSet<>();
        for (int mapId : TRAINING_MAPS) {
            List<Zone> spawnZones = randomizeZonesForMap(mapId);
            for (Zone zone : spawnZones) {
                Bot bot = createTrainingBot(zone, usedNames);
                if (bot == null) {
                    continue;
                }
                BotManager.gI().bot.add(bot);
                bot.joinAssignedZone(zone);
                moveBotToMobArea(bot, zone);
            }
        }
        AmbientBot.spawnAll(usedNames);
    }

    public synchronized void stopTrainingBots() {
        resetTrainingBots();
    }

    public boolean isTrainingBotsActive() {
        for (Bot bot : BotManager.gI().bot) {
            if (bot != null && (bot.mo1 != null || bot.ambient != null)) {
                return true;
            }
        }
        return false;
    }

    public synchronized void syncAvailability(boolean panelEnabled) {
        boolean shouldRun = panelEnabled && hasQualifiedDatabasePlayer();
        if (shouldRun && !isTrainingBotsActive()) {
            spawnTrainingBots();
        } else if (!shouldRun && isTrainingBotsActive()) {
            stopTrainingBots();
        }
    }

    public boolean hasQualifiedDatabasePlayer() {
        long now = System.currentTimeMillis();
        if (now - lastPowerCheck < POWER_CHECK_INTERVAL) {
            return qualifiedPlayerExists;
        }
        synchronized (this) {
            now = System.currentTimeMillis();
            if (now - lastPowerCheck < POWER_CHECK_INTERVAL) {
                return qualifiedPlayerExists;
            }
            String sql = "SELECT 1 FROM player p INNER JOIN account a ON a.id = p.account_id "
                    + "WHERE a.is_admin = 0 AND CAST(JSON_EXTRACT(p.data_point, '$[1]') AS UNSIGNED) >= ? LIMIT 1";
            try (Connection connection = DBConnecter.getConnectionServer();
                    PreparedStatement statement = connection.prepareStatement(sql)) {
                statement.setLong(1, REQUIRED_PLAYER_POWER);
                try (ResultSet result = statement.executeQuery()) {
                    qualifiedPlayerExists = result.next();
                }
            } catch (Exception ignored) {
            }
            lastPowerCheck = now;
            return qualifiedPlayerExists;
        }
    }

    public void runBot(int type, ShopBot shop, int slot) {
        if (type == 0) {
            spawnTrainingBots();
        }
    }

    public Zone getAssignedZone(int mapId, int zoneId) {
        Map map = MapService.gI().getMapById(mapId);
        if (map == null || zoneId < 0 || zoneId >= map.zones.size()) {
            return null;
        }
        return map.zones.get(zoneId);
    }

    public Zone getRandomEligibleZone(int mapId, int excludeZoneId) {
        Map map = MapService.gI().getMapById(mapId);
        if (map == null || map.zones == null || map.zones.isEmpty()) {
            return null;
        }
        List<Zone> eligibleZones = new ArrayList<>();
        for (Zone zone : map.zones) {
            if (zone == null) {
                continue;
            }
            if (mapId == 5 && zone.zoneId == 0) {
                continue;
            }
            if (zone.zoneId == excludeZoneId && map.zones.size() > 1) {
                continue;
            }
            eligibleZones.add(zone);
        }
        if (eligibleZones.isEmpty()) {
            return getAssignedZone(mapId, excludeZoneId);
        }
        return eligibleZones.get(Util.nextInt(0, eligibleZones.size() - 1));
    }

    private void resetTrainingBots() {
        List<Bot> bots = new ArrayList<>(BotManager.gI().bot);
        for (Bot bot : bots) {
            if (bot == null || (bot.mo1 == null && bot.ambient == null)) {
                continue;
            }
            BotManager.gI().bot.remove(bot);
            try {
                ChangeMapService.gI().exitMap(bot);
            } catch (Exception ignored) {
            }
            try {
                bot.dispose();
            } catch (Exception ignored) {
            }
        }
    }

    private Bot createTrainingBot(Zone zone, Set<String> usedNames) {
        if (zone == null) {
            return null;
        }
        byte gender = (byte) Util.nextInt(0, 2);
        Item costume = ItemService.gI().createNewItem(COSTUME_IDS[Util.nextInt(0, COSTUME_IDS.length - 1)]);
        Item accessory = ItemService.gI().createNewItem(ACCESSORY_IDS[Util.nextInt(0, ACCESSORY_IDS.length - 1)]);
        String name = BotNameRepository.randomTrainingBotName(usedNames);
        Bot bot = new Bot(name, gender, (short) costume.template.head, (short) costume.template.body,
                (short) costume.template.leg, (short) accessory.template.part, zone.map.mapId, zone.zoneId);
        bot.mo1 = new Mobb(bot);
        bot.nPoint.limitPower = 8;
        if (isLowTrainingMap(zone.map.mapId)) {
            bot.nPoint.power = Util.nextLong(50_000_000L, 100_000_000L);
        } else if (isMidTrainingMap(zone.map.mapId)) {
            bot.nPoint.power = Util.nextLong(500_000_000L, 20_000_000_000L);
        } else {
            bot.nPoint.power = Util.nextLong(30_000_000_000L, 60_000_000_000L);
        }
        bot.nPoint.tiemNang = bot.nPoint.power;
        if (isLowTrainingMap(zone.map.mapId)) {
            bot.nPoint.hpg = Util.nextInt(5_000, 10_000);
            bot.nPoint.dameg = Util.nextInt(500, 3_000);
        } else if (isMidTrainingMap(zone.map.mapId)) {
            bot.nPoint.hpg = Util.nextInt(50_000, 200_000);
            bot.nPoint.dameg = Util.nextInt(5_000, 15_000);
        } else {
            bot.nPoint.hpg = Util.nextInt(150_000, 500_000);
            bot.nPoint.dameg = Util.nextInt(10_000, 20_000);
        }
        bot.nPoint.mpg = 2_000_000_000L;
        bot.nPoint.mpMax = 2_000_000_000L;
        bot.nPoint.mp = bot.nPoint.mpMax;
        bot.nPoint.hpMax = bot.nPoint.hpg;
        bot.nPoint.hp = bot.nPoint.hpMax;
        bot.nPoint.maxStamina = 20_000;
        bot.nPoint.stamina = 20_000;
        bot.nPoint.critg = 10;
        bot.nPoint.defg = isHighTrainingMap(zone.map.mapId) ? 10_000 : 500;
        bot.clan = createFakeClan();
        ensureBodySlots(bot);
        bot.leakSkill();
        return bot;
    }

    public Bot createAmbientBot(Zone zone, Set<String> usedNames) {
        Bot bot = createTrainingBot(zone, usedNames);
        if (bot != null) {
            bot.mo1 = null;
            bot.nPoint.mp = bot.nPoint.mpMax;
        }
        return bot;
    }

    public boolean isLowTrainingMap(int mapId) {
        return mapId == 5 || mapId == 2 || mapId == 3 || mapId == 9 || mapId == 11
                || mapId == 16 || mapId == 17;
    }

    public boolean isMidTrainingMap(int mapId) {
        return mapId == 68 || mapId == 77 || mapId == 79 || mapId == 83
                || mapId == 70 || mapId == 69 || mapId == 71 || mapId == 76;
    }

    public boolean isHighTrainingMap(int mapId) {
        return mapId == 92 || mapId == 100 || mapId == 99 || mapId == 93 || mapId == 109;
    }

    private void ensureBodySlots(Bot bot) {
        while (bot.inventory.itemsBody.size() < 12) {
            bot.inventory.itemsBody.add(ItemService.gI().createItemNull());
        }
    }

    private void moveBotToMobArea(Bot bot, Zone zone) {
        if (zone == null || zone.mobs == null || zone.mobs.isEmpty()) {
            return;
        }
        int index = Util.nextInt(0, zone.mobs.size() - 1);
        bot.moveTo(zone.mobs.get(index).location.x, zone.mobs.get(index).location.y);
    }

    private List<Zone> randomizeZonesForMap(int mapId) {
        Map map = MapService.gI().getMapById(mapId);
        List<Zone> result = new ArrayList<>();
        if (map == null || map.zones == null || map.zones.isEmpty()) {
            return result;
        }
        List<Zone> eligibleZones = new ArrayList<>();
        for (Zone zone : map.zones) {
            if (zone == null) {
                continue;
            }
            if (mapId == 5 && zone.zoneId == 0) {
                continue;
            }
            eligibleZones.add(zone);
        }
        if (eligibleZones.isEmpty()) {
            return result;
        }
        int[] perZone = new int[eligibleZones.size()];
        int assigned = 0;
        while (assigned < BOTS_PER_MAP) {
            List<Integer> candidates = new ArrayList<>();
            for (int i = 0; i < eligibleZones.size(); i++) {
                if (perZone[i] < 2) {
                    candidates.add(i);
                }
            }
            if (candidates.isEmpty()) {
                for (int i = 0; i < eligibleZones.size(); i++) {
                    candidates.add(i);
                }
            }
            int pickedIndex = candidates.get(Util.nextInt(0, candidates.size() - 1));
            perZone[pickedIndex]++;
            assigned++;
        }
        for (int i = 0; i < eligibleZones.size(); i++) {
            for (int count = 0; count < perZone[i]; count++) {
                result.add(eligibleZones.get(i));
            }
        }
        return result;
    }

    private Clan createFakeClan() {
        Clan clan = new Clan();
        clan.name = CLAN_NAMES[Util.nextInt(0, CLAN_NAMES.length - 1)];
        clan.name2 = "";
        clan.slogan = "";
        clan.imgId = -1;
        clan.active = false;
        return clan;
    }
}
