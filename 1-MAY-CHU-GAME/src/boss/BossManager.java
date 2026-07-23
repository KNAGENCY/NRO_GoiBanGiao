package boss;

/*
 *
 *
 * @author GiaCatLuong
 */
import GiaCatLuongManager.Functions;
// import boss.boss_manifest.AnTrom.AnTrom;
// import boss.boss_manifest.AnTrom.AnTromTV;
import boss.boss_manifest.Black.BlackGoku;
import boss.boss_manifest.Black.BG;
import boss.boss_manifest.Black.BA;
import boss.boss_manifest.Black.BB;
import boss.boss_manifest.Black.BC;
import boss.boss_manifest.Black.BD;

import server.GameExecutor;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.io.OutputStreamWriter;
import java.io.Writer;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.TimeUnit;
import boss.boss_manifest.Nappa.Rambo;
import boss.boss_manifest.Nappa.MapDauDinh;
import boss.boss_manifest.Nappa.Kuku;
import boss.boss_manifest.Android.Android19;
import boss.boss_manifest.Android.Pic;
import boss.boss_manifest.Android.Android14;
import boss.boss_manifest.Android.Poc;
import boss.boss_manifest.Android.Android13;
import boss.boss_manifest.Android.KingKong;
import boss.boss_manifest.Android.DrKore;
import java.util.concurrent.CopyOnWriteArrayList;

import boss.boss_manifest.Android.Android15;
import boss.boss_manifest.Cooler.Cooler;
import boss.boss_manifest.Cell.SieuBoHung;
import boss.boss_manifest.Cell.XenBoHung;
import boss.boss_manifest.Bossvip.Broly;
import boss.boss_manifest.Bossvip.SuperBroly;
import boss.boss_manifest.Bossvip.TestDameBoss;
import boss.boss_manifest.Bossvip.Mabuz;
import boss.boss_manifest.Bossvip.OngGiaNoel;
import boss.boss_manifest.TaoPaiPai.TaoPaiPai;
import boss.boss_manifest.Frieza.Fide;
import boss.boss_manifest.Bossvip.SonTinh;
import boss.boss_manifest.Bossvip.ThuyTinh;
// import boss.boss_manifest.Bossvip.BiMa;
// import boss.boss_manifest.Bossvip.Doi;
// import boss.boss_manifest.Bossvip.MaTroi;
import boss.boss_manifest.Bossvip.NguyetThan;
import boss.boss_manifest.Bossvip.NhatThan;
import boss.boss_manifest.MajinBuu12H.Mabu;
import boss.boss_manifest.MajinBuu12H.BuiBui;
import boss.boss_manifest.MajinBuu12H.BuiBui2;
import boss.boss_manifest.MajinBuu12H.Cadic;
import boss.boss_manifest.MajinBuu12H.Drabura;
import boss.boss_manifest.MajinBuu12H.Drabura2;
import boss.boss_manifest.MajinBuu12H.Drabura3;
import boss.boss_manifest.MajinBuu12H.Goku;
import boss.boss_manifest.MajinBuu12H.Yacon;
import boss.boss_manifest.MajinBuu14H.Mabu2H;
import boss.boss_manifest.MajinBuu14H.G;
import boss.boss_manifest.MajinBuu14H.SuperBu;
import boss.boss_manifest.GinyuForce.SO1;
import boss.boss_manifest.GinyuForce.SO2;
import boss.boss_manifest.GinyuForce.SO3;
import boss.boss_manifest.GinyuForce.SO4;
import boss.boss_manifest.GinyuForce.TDT;
import boss.boss_manifest.NamekGinyuForce.SO1_NM;
import boss.boss_manifest.NamekGinyuForce.SO2_NM;
import boss.boss_manifest.NamekGinyuForce.SO3_NM;
import boss.boss_manifest.NamekGinyuForce.SO4_NM;
import boss.boss_manifest.NamekGinyuForce.TDT_NM;
import boss.boss_manifest.Earth.BIDO;
import boss.boss_manifest.Earth.BOJACK;
import boss.boss_manifest.Earth.BUJIN;
import boss.boss_manifest.Earth.KOGU;
import boss.boss_manifest.Earth.SUPER_BOJACK;
import boss.boss_manifest.Earth.ZANGYA;
import boss.boss_manifest.Cell.XENCON1;
import boss.boss_manifest.Cell.XENCON2;
import boss.boss_manifest.Cell.XENCON3;
import boss.boss_manifest.Cell.XENCON4;
import boss.boss_manifest.Cell.XENCON5;
import boss.boss_manifest.Cell.XENCON6;
import boss.boss_manifest.Cell.XENCON7;
import boss.boss_manifest.Cumber.Cumber;
import boss.boss_manifest.DaiTuongBroly.DaiTuongBroly;
import boss.boss_manifest.Hatchiyac.Hatchiyac;
// import boss.boss_manifest.Doraemon.*;
import boss.boss_manifest.BrolySsj.*;
import boss.boss_manifest.Wolves.*;
import boss.boss_manifest.Bossvip.LanCon;
import player.Player;
import network.Message;
import services.MapService;
import services.Service;
import services.func.ChangeMapService;

import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import map.Zone;
import server.Maintenance;
import utils.Logger;

public class BossManager {

    private static BossManager instance;
    public static byte ratioReward = 10;
    private static final long PANEL_MAX_BOSS_HP = 2_000_000_000L;
    private static final String PANEL_GAME_ROOT = resolvePanelGameRoot();
    private static final String PANEL_BOSS_PERSIST_STATUS_FILE = new File(PANEL_GAME_ROOT, "panel_boss_persist_status.txt").getPath();
    private static final String PANEL_BOSS_OVERRIDES_FILE = new File(PANEL_GAME_ROOT, "panel_boss_overrides.properties").getPath();
    private volatile boolean started = false;
// ===== DANH SÁCH BOSS VIP / CHÍNH =====
    private boolean isMainBoss(Boss boss) {
    return boss instanceof Broly    
     || boss instanceof XenBoHung  || boss instanceof Baconsoi  || boss instanceof SieuBoHung  || boss instanceof XENCON1|| boss instanceof S
     || boss instanceof BrolySsj
     || boss instanceof NguyetThan
        || boss instanceof BlackGoku  || boss instanceof BG  || boss instanceof BA || boss instanceof BB || boss instanceof BC || boss instanceof BD
        || boss instanceof Cumber
        || boss instanceof DaiTuongBroly
        || boss instanceof OngGiaNoel
        || boss instanceof ThuyTinh
        || boss instanceof BUJIN      || boss instanceof SUPER_BOJACK
        || boss instanceof Mabuz
             || boss instanceof TDT      || boss instanceof TDT_NM

               || boss instanceof SO4_NM 
             || boss instanceof SO4
       
     
          || boss instanceof Android19                   || boss instanceof Pic 
                  || boss instanceof Cooler
                       || boss instanceof Kuku
                            || boss instanceof MapDauDinh
     || boss instanceof Rambo
          || boss instanceof Fide
              || boss instanceof Android15

        || boss instanceof Hatchiyac;
}

    private static String resolvePanelGameRoot() {
        String root = System.getProperty("game.root");
        if (root == null || root.trim().isEmpty()) {
            root = System.getenv("GAME_ROOT");
        }
        if (root != null && !root.trim().isEmpty()) {
            return root.trim();
        }
        String userDir = System.getProperty("user.dir");
        if (userDir != null && new File(userDir, "data\\config\\config.properties").isFile()) {
            return new File(userDir).getAbsolutePath();
        }
        String userHome = System.getProperty("user.home");
        if (userHome != null && !userHome.trim().isEmpty()) {
            File homeRoot = new File(userHome, "Music\\zz1");
            if (new File(homeRoot, "data\\config\\config.properties").isFile()) {
                return homeRoot.getAbsolutePath();
            }
        }
        File AdminRoot = new File("C:\\Users\\Admin\\Music\\zz1");
        if (new File(AdminRoot, "data\\config\\config.properties").isFile()) {
            return AdminRoot.getAbsolutePath();
        }
        return "C:\\Users\\Administrator\\Music\\zz1";
    }

    private boolean isPanelBossPersistEnabled() {
        File statusFile = new File(PANEL_BOSS_PERSIST_STATUS_FILE);
        if (!statusFile.isFile()) {
            return false;
        }
        try (BufferedReader reader = new BufferedReader(new FileReader(statusFile))) {
            String value = reader.readLine();
            return value != null && value.trim().equals("1");
        } catch (IOException e) {
            Logger.logException(BossManager.class, e, "Failed to read boss persist status");
            return false;
        }
    }

    private Properties loadPanelBossOverrides() {
        Properties properties = new Properties();
        File file = new File(PANEL_BOSS_OVERRIDES_FILE);
        if (!file.isFile()) {
            return properties;
        }
        try (FileInputStream input = new FileInputStream(file)) {
            properties.load(input);
        } catch (IOException e) {
            Logger.logException(BossManager.class, e, "Failed to read boss overrides");
        }
        return properties;
    }

    private void savePanelBossOverride(Boss boss, long hp, long dame, String name) {
        Properties properties = loadPanelBossOverrides();
        String key = String.valueOf((int) boss.id);
        properties.setProperty(key + ".hp", String.valueOf(hp));
        properties.setProperty(key + ".dame", String.valueOf(dame));
        properties.remove(key + ".secondsRest");
        if (name != null && !name.trim().isEmpty()) {
            properties.setProperty(key + ".name", name.trim());
        }
        try (FileOutputStream output = new FileOutputStream(PANEL_BOSS_OVERRIDES_FILE)) {
            properties.store(output, "Panel boss overrides");
        } catch (IOException e) {
            Logger.logException(BossManager.class, e, "Failed to save boss override");
        }
    }

    public void applyPanelBossOverride(Boss boss) {
        if (boss == null || boss.nPoint == null || !isPanelBossPersistEnabled()) {
            return;
        }
        Properties properties = loadPanelBossOverrides();
        String key = String.valueOf((int) boss.id);
        String hpText = properties.getProperty(key + ".hp", "").trim();
        if (hpText.isEmpty()) {
            return;
        }
        try {
            long hp = Long.parseLong(hpText);
            if (hp <= 0 || hp > PANEL_MAX_BOSS_HP) {
                return;
            }
            String name = properties.getProperty(key + ".name", "").trim();
            if (!name.isEmpty()) {
                boss.name = name;
            }
            boss.nPoint.hpg = hp;
            boss.nPoint.hpMax = hp;
            boss.nPoint.hp = hp;
            String dameText = properties.getProperty(key + ".dame", "").trim();
            if (!dameText.isEmpty()) {
                long dame = Long.parseLong(dameText);
                if (dame > 0 && dame <= PANEL_MAX_BOSS_HP) boss.nPoint.dameg = dame;
            }
        } catch (NumberFormatException e) {
            Logger.logException(BossManager.class, e, "Invalid boss override hp");
        }
    }

    public static BossManager gI() {
        if (instance == null) {
            instance = new BossManager();
        }
        return instance;
    }

    private static final String BOSS_STATUS_OUTPUT_FILE = resolveBossStatusOutputFile();

    private static String resolveBossStatusOutputFile() {
        String statusFile = System.getProperty("boss.status.file");
        if (statusFile == null || statusFile.trim().isEmpty()) {
            statusFile = System.getenv("BOSS_STATUS_FILE");
        }
        if (statusFile != null && !statusFile.trim().isEmpty()) {
            return statusFile.trim();
        }

        String webRoot = System.getProperty("web.root");
        if (webRoot == null || webRoot.trim().isEmpty()) {
            webRoot = System.getenv("WEB_ROOT");
        }
        if (webRoot != null && !webRoot.trim().isEmpty()) {
            return new File(new File(webRoot.trim(), "sv1"), "boss_status.json").getPath();
        }

        return "C:\\xampp\\htdocs\\sv1\\boss_status.json";
    }

    public List<Boss> getBosses() {

        return this.bosses;

    }

    private String escapeJson(String text) {
        if (text == null) {
            return "";
        }
        return text.replace("\\", "\\\\")
                .replace("\"", "\\\"")
                .replace("\b", "\\b")
                .replace("\f", "\\f")
                .replace("\n", "\\n")
                .replace("\r", "\\r")
                .replace("\t", "\\t");
    }

    private void exportBossStatus() {
        long now = System.currentTimeMillis();
        StringBuilder json = new StringBuilder();
        json.append("{\"generatedAt\":").append(now).append(",\"bosses\":[");
        boolean first = true;
        for (Boss boss : bosses) {
            if (boss == null || boss.data == null || boss.data.length == 0
                    || boss.zone == null || boss.isDie() || boss.nPoint == null || boss.nPoint.hp <= 0) {
                continue;
            }
            int activeLevel = boss.currentLevel;
            if (activeLevel < 0 || activeLevel >= boss.data.length) {
                continue;
            }
            BossData activeData = boss.data[activeLevel];
            if (!first) {
                json.append(',');
            }
            first = false;
            String state = "S\u1ed1ng";
            long respawnIn = 0;
            if (boss.bossStatus == BossStatus.REST || boss.bossStatus == BossStatus.RESPAWN) {
                respawnIn = Math.max(0, boss.secondsRest * 1000L - (now - boss.lastTimeRest));
            }
            String mapName = boss.zone.map.mapName + " (khu " + boss.zone.zoneId + ")";
            int mapId = boss.zone.map.mapId;
            int zoneId = boss.zone.zoneId;
            long hp = boss.nPoint.hp;
            long maxHp = boss.nPoint.hpg;
            String displayName = boss.name != null && !boss.name.trim().isEmpty() ? boss.name : activeData.getName();
            json.append('{')
                    .append("\"id\":").append(boss.id)
                    .append(",\"name\":\"").append(escapeJson(displayName)).append("\"")
                    .append(",\"status\":\"").append(boss.bossStatus.toString()).append("\"")
                    .append(",\"state\":\"").append(escapeJson(state)).append("\"")
                    .append(",\"hp\":").append(Math.max(0L, hp))
                    .append(",\"maxHp\":").append(Math.max(0L, maxHp))
                    .append(",\"dame\":").append(Math.max(0L, boss.nPoint.dameg))
                    .append(",\"mapName\":\"").append(escapeJson(mapName)).append("\"")
                    .append(",\"mapId\":").append(mapId)
                    .append(",\"zoneId\":").append(zoneId)
                    .append(",\"respawnIn\":").append(respawnIn)
                    .append(",\"outfit0\":").append(activeData.getOutfit()[0])
                    .append(",\"outfit1\":").append(activeData.getOutfit()[1])
                    .append(",\"outfit2\":").append(activeData.getOutfit()[2])
                    .append('}');
        }
        json.append("]}");
        File file = new File(BOSS_STATUS_OUTPUT_FILE);
        File parent = file.getParentFile();
        if (parent != null && !parent.isDirectory()) {
            parent.mkdirs();
        }
        try (Writer writer = new OutputStreamWriter(new FileOutputStream(file, false), StandardCharsets.UTF_8)) {
            writer.write(json.toString());
        } catch (IOException e) {
            Logger.logException(BossManager.class, e, "Failed to export boss status");
        }
    }

    protected final List<Boss> bosses = new CopyOnWriteArrayList<>();

    public void addBoss(Boss boss) {
        this.bosses.add(boss);
    }

    public void removeBoss(Boss boss) {
        this.bosses.remove(boss);
    }

    public synchronized boolean respawnBossFromPanel(int bossId, int mapId, int zoneId) {
        Boss oldBoss = findBossForPanelRespawn(bossId, mapId, zoneId);
        int createBossId = bossId;
        if (oldBoss != null) {
            Boss rootBoss = oldBoss.parentBoss != null ? oldBoss.parentBoss : oldBoss;
            createBossId = (int) rootBoss.id;
            removeBossClusterFromPanel(rootBoss);
        }

        Boss newBoss = createBoss(createBossId);
        if (newBoss == null) {
            return false;
        }
        newBoss.lastTimeRest = 0;
        newBoss.changeStatus(BossStatus.RESPAWN);
        exportBossStatus();
        return true;
    }

    public synchronized boolean editBossFromPanel(int bossId, int mapId, int zoneId, long hp, long dame, String name) {
        Boss boss = findBossForPanelRespawn(bossId, mapId, zoneId);
        if (boss == null || boss.zone == null || boss.isDie() || boss.nPoint == null) {
            return false;
        }
        if (hp <= 0 || hp > PANEL_MAX_BOSS_HP || dame <= 0 || dame > PANEL_MAX_BOSS_HP) {
            return false;
        }
        String cleanName = name == null ? "" : name.trim();
        savePanelBossOverride(boss, hp, dame, cleanName);
        if (!isPanelBossPersistEnabled()) {
            exportBossStatus();
            return true;
        }
        if (!cleanName.isEmpty()) {
            boss.name = cleanName;
        }
        boss.nPoint.hpg = hp;
        boss.nPoint.hpMax = hp;
        boss.nPoint.hp = hp;
        boss.nPoint.dameg = dame;
        Service.gI().Send_Info_NV(boss);
        exportBossStatus();
        return true;
    }

    private void removeBossClusterFromPanel(Boss rootBoss) {
        List<Boss> cluster = new ArrayList<>();
        collectBossCluster(rootBoss, cluster);
        for (Boss boss : cluster) {
            removeBossFromPanel(boss);
        }
    }

    private void collectBossCluster(Boss boss, List<Boss> cluster) {
        if (boss == null || cluster.contains(boss)) {
            return;
        }
        cluster.add(boss);
        if (boss.bossAppearTogether == null) {
            return;
        }
        for (Boss[] group : boss.bossAppearTogether) {
            if (group == null) {
                continue;
            }
            for (Boss child : group) {
                collectBossCluster(child, cluster);
            }
        }
    }

    private Boss findBossForPanelRespawn(int bossId, int mapId, int zoneId) {
        if (mapId >= 0 && zoneId >= 0) {
            for (Boss boss : bosses) {
                if (boss != null && boss.id == bossId && boss.zone != null
                        && boss.zone.map.mapId == mapId && boss.zone.zoneId == zoneId) {
                    return boss;
                }
            }
        }
        for (Boss boss : bosses) {
            if (boss != null && boss.id == bossId) {
                return boss;
            }
        }
        return null;
    }

    private void removeBossFromPanel(Boss boss) {
        try {
            if (boss.zone != null) {
                ChangeMapService.gI().exitMap(boss);
            }
        } catch (Exception e) {
            Logger.logException(BossManager.class, e, "Failed to exit old boss from map");
        }
        removeBoss(boss);
        try {
            boss.dispose();
        } catch (Exception e) {
            Logger.logException(BossManager.class, e, "Failed to dispose old boss");
        }
    }

    public void loadBoss() {
        this.createBoss(BossID.ONG_GIA_NOEL);
        this.createBoss(BossID.DAITUONGBROLY);
        this.createBoss(BossID.Hatchiyac);
        this.createBoss(BossID.TIEU_DOI_TRUONG);
        this.createBoss(BossID.TIEU_DOI_TRUONG_NM);
        this.createBoss(BossID.BOJACK);
        this.createBoss(BossID.SUPER_BOJACK);
        this.createBoss(BossID.KING_KONG);
        this.createBoss(BossID.TAU_PAY_PAY_DONG_NAM_KARIN, 5);
        this.createBoss(BossID.XEN_BO_HUNG);
        this.createBoss(BossID.SIEU_BO_HUNG);
        this.createBoss(BossID.KUKU, 2);
        this.createBoss(BossID.MAP_DAU_DINH);
        this.createBoss(BossID.RAMBO);
        this.createBoss(BossID.FIDE);
        this.createBoss(BossID.ANDROID_14);
        this.createBoss(BossID.DR_KORE);
        this.createBoss(BossID.COOLER);
    // this.createBoss(BossID.G);
        this.createBoss(BossID.BLACK_GOKU, 3);
this.createBoss(BossID.BG);
        this.createBoss(BossID.BA);
        this.createBoss(BossID.BB);
        this.createBoss(BossID.BC);
              this.createBoss(BossID.BD);
        // this.createBoss(BossID.BLACK_GOKU);
        // this.createBoss(BossID.AN_TROM, 3);
        // this.createBoss(BossID.AN_TROM_TV, 3);
        this.createBoss(BossID.BROLY, 5);
        // this.createBoss(BossID.SUPER_BROLY, 3);
        this.createBoss(BossID.CUMBER);
       this.createBoss(BossID.CUMBER);
        this.createBoss(BossID.BROLY_SSJ);
        this.createBoss(BossID.BA_CON_SOI);
    this.createBoss(BossID.S);
        this.createBoss(BossID.SON_TINH);
        this.createBoss(BossID.THUY_TINH);
        // this.createBoss(BossID.CHILL_1);
        // this.createBoss(BossID.MATROI);
        // this.createBoss(BossID.BIMA);
        // this.createBoss(BossID.DOI);
        this.createBoss(BossID.LAN_CON);
        // this.createBoss(BossID.ONG_GIA_NOEL);
        this.createBoss(BossID.NHATTHAN,6);
        this.createBoss(BossID.NGUYETTHAN,6);
        this.createBoss(BossID.MABUZ);
        this.createBoss(BossID.TEST_DAME);
    }

    public void createBoss(int bossID, int total) {
        for (int i = 0; i < total; i++) {
            createBoss(bossID);
        }
    }

    public Boss createBoss(int bossID) {
        try {
            return switch (bossID) {
                case BossID.BROLY_SSJ ->
                    new BrolySsj();
    case BossID.S ->
                    new S();
                case BossID.BA_CON_SOI ->
                    new Baconsoi();
                case BossID.DAITUONGBROLY ->
                    new DaiTuongBroly();
                case BossID.Hatchiyac ->
                    new Hatchiyac();
                case BossID.SO_4 ->
                    new SO4();
                case BossID.SO_3 ->
                    new SO3();
                case BossID.SO_2 ->
                    new SO2();
                case BossID.SO_1 ->
                    new SO1();
                case BossID.TIEU_DOI_TRUONG ->
                    new TDT();
                case BossID.SO_4_NM ->
                    new SO4_NM();
                case BossID.SO_3_NM ->
                    new SO3_NM();
                case BossID.SO_2_NM ->
                    new SO2_NM();
                case BossID.SO_1_NM ->
                    new SO1_NM();
                case BossID.TIEU_DOI_TRUONG_NM ->
                    new TDT_NM();
                case BossID.BUJIN ->
                    new BUJIN();
                case BossID.KOGU ->
                    new KOGU();
                case BossID.ZANGYA ->
                    new ZANGYA();
                case BossID.BIDO ->
                    new BIDO();
                case BossID.BOJACK ->
                    new BOJACK();
                case BossID.SUPER_BOJACK ->
                    new SUPER_BOJACK();
                case BossID.KUKU ->
                    new Kuku();
                case BossID.MAP_DAU_DINH ->
                    new MapDauDinh();
                case BossID.RAMBO ->
                    new Rambo();
                case BossID.TAU_PAY_PAY_DONG_NAM_KARIN ->
                    new TaoPaiPai();
                case BossID.DRABURA ->
                    new Drabura();
                case BossID.BUI_BUI ->
                    new BuiBui();
                case BossID.BUI_BUI_2 ->
                    new BuiBui2();
                case BossID.YA_CON ->
                    new Yacon();
                case BossID.DRABURA_2 ->
                    new Drabura2();
                case BossID.GOKU ->
                    new Goku();
                case BossID.CADIC ->
                    new Cadic();
                case BossID.MABU_12H ->
                    new Mabu();
                case BossID.DRABURA_3 ->
                    new Drabura3();
                case BossID.MABU ->
                    new Mabu2H();
                    //    case BossID.G ->
                    // new G();
                case BossID.SUPERBU ->
                    new SuperBu();
                case BossID.FIDE ->
                    new Fide();
                case BossID.DR_KORE ->
                    new DrKore();
                case BossID.ANDROID_19 ->
                    new Android19();
                case BossID.ANDROID_13 ->
                    new Android13();
                case BossID.ANDROID_14 ->
                    new Android14();
                case BossID.ANDROID_15 ->
                    new Android15();
                case BossID.PIC ->
                    new Pic();
                case BossID.POC ->
                    new Poc();
                case BossID.KING_KONG ->
                    new KingKong();
                case BossID.XEN_BO_HUNG ->
                    new XenBoHung();
                case BossID.SIEU_BO_HUNG ->
                    new SieuBoHung();
                case BossID.XEN_CON_1 ->
                    new XENCON1();
                case BossID.XEN_CON_2 ->
                    new XENCON2();
                case BossID.XEN_CON_3 ->
                    new XENCON3();
                case BossID.XEN_CON_4 ->
                    new XENCON4();
                case BossID.XEN_CON_5 ->
                    new XENCON5();
                case BossID.XEN_CON_6 ->
                    new XENCON6();
                case BossID.XEN_CON_7 ->
                    new XENCON7();
                case BossID.COOLER ->
                    new Cooler();
                case BossID.BROLY ->
                    new Broly();
                case BossID.SUPER_BROLY ->
                    new SuperBroly();
                // case BossID.AN_TROM ->
                //     new AnTrom();
                // case BossID.AN_TROM_TV ->
                //     new AnTromTV();
                // case BossID.KHIDOT ->
                // new KhiDot();
                case BossID.NGUYETTHAN ->
                    new NguyetThan();
                case BossID.NHATTHAN ->
                    new NhatThan();

                // case BossID.BIMA ->
                //     new BiMa();
                // case BossID.MATROI ->
                //     new MaTroi();
                // case BossID.DOI ->
                //     new Doi();
                case BossID.ONG_GIA_NOEL ->
                    new OngGiaNoel();
                case BossID.SON_TINH ->
                    new SonTinh();
                case BossID.THUY_TINH ->
                    new ThuyTinh();
                case BossID.LAN_CON ->
                    new LanCon();
                case BossID.BLACK_GOKU ->
                    new BlackGoku();
   case BossID.BA ->
                    new BA();
                       case BossID.BB ->
                    new BB();
                       case BossID.BC ->
                    new BC();


       case BossID.BD ->
                    new BD();

                    case BossID.BG ->
                    new BG();
                case BossID.MABUZ ->
                    new Mabuz();
                case BossID.TEST_DAME ->
                    new TestDameBoss();
                case BossID.CUMBER ->
                    new Cumber();
                default ->
                    null;
            };
        } catch (Exception e) {
            Logger.error(e + "\n");
            return null;
        }
    }

    public void showListBoss(Player player) {
        showListBoss(player, false);
    }

    public void showListBoss(Player player, boolean readOnly) {
        // if (!player.isAdmin()) {
        //     return;
        // }

        player.iDMark.setMenuType(3); // MENU LIST BOSS (ADMIN)

        Message msg;
        try {
            msg = new Message(-96);
            msg.writer().writeByte(readOnly ? 2 : 0);
            msg.writer().writeUTF("Boss");

       msg.writer().writeByte((int) bosses.stream()
        .filter(boss ->
                isMainBoss(boss)
                && !MapService.gI().isMapBossFinal(boss.data[0].getMapJoin()[0])
                && !MapService.gI().isMapHuyDiet(boss.data[0].getMapJoin()[0])
                && !MapService.gI().isMapYardart(boss.data[0].getMapJoin()[0])
                && !MapService.gI().isMapMaBu(boss.data[0].getMapJoin()[0])
                && !MapService.gI().isMapBlackBallWar(boss.data[0].getMapJoin()[0])
        ).count());

            for (int i = 0; i < bosses.size(); i++) {
                Boss boss = bosses.get(i);

         if (!isMainBoss(boss)
        || MapService.gI().isMapBossFinal(boss.data[0].getMapJoin()[0])
        || MapService.gI().isMapYardart(boss.data[0].getMapJoin()[0])
        || MapService.gI().isMapHuyDiet(boss.data[0].getMapJoin()[0])
        || MapService.gI().isMapMaBu(boss.data[0].getMapJoin()[0])
        || MapService.gI().isMapBlackBallWar(boss.data[0].getMapJoin()[0])) {
    continue;
}


                msg.writer().writeInt(i); // ID
                msg.writer().writeInt(i); // INDEX -> CLIENT TRẢ LẠI KHI BẤM "THÁCH ĐẤU"
                msg.writer().writeShort(boss.data[0].getOutfit()[0]);

                if (player.session != null && player.session.version >= 214) {
                    msg.writer().writeShort(-1);
                }

                msg.writer().writeShort(boss.data[0].getOutfit()[1]);
                msg.writer().writeShort(boss.data[0].getOutfit()[2]);
                msg.writer().writeUTF(boss.data[0].getName());

                if (boss.zone != null) {
                    msg.writer().writeUTF(boss.bossStatus.toString());
                    msg.writer().writeUTF(
                            boss.zone.map.mapName);
                } else {
                    msg.writer().writeUTF(boss.bossStatus.toString());
                    msg.writer().writeUTF("Chết rồi");
                }
            }

            player.sendMessage(msg);
            msg.cleanup();
        } catch (Exception ignored) {
        }
    }

    public Boss getBossById(int bossId) {
        return this.bosses.stream().filter(boss -> boss.id == bossId && !boss.isDie()).findFirst().orElse(null);
    }

    public boolean checkBosses(Zone zone, int BossID) {
        return this.bosses.stream()
                .filter(boss -> boss.id == BossID && boss.zone != null && boss.zone.equals(zone) && !boss.isDie())
                .findFirst().orElse(null) != null;
    }

    public Player findBossClone(Player player) {
        return player.zone.getBosses().stream().filter(boss -> boss.id < -100_000_000 && !boss.isDie()).findFirst()
                .orElse(null);
    }

    public Boss getBossById(int bossId, int mapId, int zoneId) {
        return this.bosses.stream().filter(boss -> boss.id == bossId && boss.zone != null
                && boss.zone.map.mapId == mapId && boss.zone.zoneId == zoneId && !boss.isDie()).findFirst()
                .orElse(null);
    }

    public synchronized void start() {
        if (started) {
            Logger.warning("BossManager already started – skip");
            return;
        }
        started = true;

        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            if (Maintenance.isRunning)
                return;

            try {
                for (Boss boss : bosses) {
                    if (boss != null) {
                        boss.update();
                    }
                }
            } catch (Exception e) {
                Logger.logException(BossManager.class, e, "Boss update error");
            }

        }, 0, 150, TimeUnit.MILLISECONDS);

        GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
            if (Maintenance.isRunning)
                return;
            exportBossStatus();
        }, 0, 1, TimeUnit.SECONDS);

        Logger.success("BossManager scheduled");
    }

    public Boss getBossByIndex(int index) {
        if (index < 0 || index >= this.bosses.size()) {
            return null;
        }
        return this.bosses.get(index);
    }

}
