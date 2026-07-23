package services.func;

import boss.Boss;
import item.Item;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.util.Properties;
import java.sql.*;
import java.time.LocalDate;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import jdbc.DBConnecter;
import network.Message;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import utils.Logger;

public final class PhucLoiService {
    private static final PhucLoiService INSTANCE = new PhucLoiService();
    private static final int TOTAL = 15;
    private final Map<Integer, Long> ticks = new ConcurrentHashMap<>();

    // ------------------------------------------------------------------
    // Bộ đếm giờ online: đếm trong bộ nhớ, định kỳ mới ghi xuống CSDL.
    //
    // Bản cũ ghi thẳng xuống CSDL mỗi giây cho MỖI người chơi, và mỗi lần
    // là 3 câu lệnh (INSERT IGNORE + UPDATE đổi ngày + UPDATE cộng giây).
    // 300 người online = ~900 câu lệnh mỗi giây, chỉ để đếm thời gian ngồi
    // chơi — trong khi dữ liệu nhân vật thật thì cả buổi mới ghi một lần.
    // Đúng khung giờ đông nhất (20h-23h) đây là chỗ nghẽn đầu tiên.
    //
    // Nay: cộng dồn trong bộ nhớ, 60 giây mới ghi một lần → giảm khoảng
    // 180 lần số câu lệnh. Người chơi không thấy khác biệt vì mọi chỗ đọc
    // số liệu (mở bảng Phúc Lợi, bấm nhận quà) đều đẩy phần đang chờ
    // xuống CSDL trước khi đọc.
    // ------------------------------------------------------------------
    private static final long FLUSH_INTERVAL_MS = 60_000L;
    private static final long CONFIG_CHECK_INTERVAL_MS = 5_000L;
    /** Số giây online đã đếm nhưng chưa kịp ghi xuống CSDL. */
    private final Map<Integer, Integer> pendingSeconds = new ConcurrentHashMap<>();
    /** Lần ghi xuống CSDL gần nhất của từng tài khoản. */
    private final Map<Integer, Long> lastFlush = new ConcurrentHashMap<>();
    /** Ngày mà dòng dữ liệu của tài khoản đã được bảo đảm tồn tại + reset. */
    private final Map<Integer, String> ensuredDate = new ConcurrentHashMap<>();
    private volatile long configLastCheck = Long.MIN_VALUE;
    private static final String[] TAB_NAMES = {"1. Quà Online", "2. Quà Cày Chay", "3. Quà Đặc Biệt", "4. Quà Sức Mạnh", "5. Quà Nhiệm Vụ"};
    private static final int[][] LIMITS = {{300,900,3600},{50000,100000,500000},{5,5,10},{10,20,40},{11,22,27}};
    private static int[][][] ITEM_IDS = {
        {{447,441,442,443},{381,382,383,384},{457}},
        {{20,19,457},{16,457},{1499,457}},
        {{457},{457},{457}}, {{1152},{1151},{1150}}, {{579},{1628},{457}}
    };
    private static int[][][] QTY = {
        {{5,5,5,5},{5,5,5,5},{10}}, {{10,10,10},{15,15},{1,100}},
        {{10},{20},{50}}, {{10},{10},{10}}, {{5},{5},{20}}
    };
    private static final int[][][] SOURCE_ITEM_IDS = ITEM_IDS;
    private static final int[][][] SOURCE_QTY = QTY;

    private static int[][][][][] CUSTOM_OPTIONS;
    private long configLastModified = Long.MIN_VALUE;
    private long configLastLength = Long.MIN_VALUE;
    public static PhucLoiService gI() { return INSTANCE; }
    private PhucLoiService() { ensureTable(); }
    private Connection connection() throws SQLException { return DBConnecter.getConnectionServer(); }

    private void ensureTable() {
        String sql = "CREATE TABLE IF NOT EXISTS phuc_loi_progress (account_id INT PRIMARY KEY, online_date DATE NOT NULL, online_seconds INT NOT NULL DEFAULT 0, online_claim INT NOT NULL DEFAULT 0, nap_claim INT NOT NULL DEFAULT 0, broly INT NOT NULL DEFAULT 0, so4 INT NOT NULL DEFAULT 0, xen INT NOT NULL DEFAULT 0, boss_claim INT NOT NULL DEFAULT 0, power_claim INT NOT NULL DEFAULT 0, task_claim INT NOT NULL DEFAULT 0, last_tick BIGINT NOT NULL DEFAULT 0)";
        try (Connection c=connection(); Statement s=c.createStatement()) { s.executeUpdate(sql); }
        catch (Exception e) { Logger.logException(PhucLoiService.class,e); }
    }
    /**
     * Bảo đảm tài khoản đã có dòng dữ liệu của ngày hôm nay.
     * Mỗi tài khoản chỉ chạy thật đúng một lần mỗi ngày, các lần sau bỏ qua ngay.
     */
    private void ensureRow(int userId) throws SQLException {
        String today = LocalDate.now().toString();
        if (today.equals(ensuredDate.get(userId))) return;
        try (Connection c=connection(); PreparedStatement ps=c.prepareStatement("INSERT IGNORE INTO phuc_loi_progress(account_id,online_date,last_tick) VALUES(?,?,?)")) {
            ps.setInt(1,userId); ps.setString(2,today); ps.setLong(3,System.currentTimeMillis()); ps.executeUpdate();
        }
        try (Connection c=connection(); PreparedStatement ps=c.prepareStatement("UPDATE phuc_loi_progress SET online_date=?,online_seconds=0,online_claim=0,last_tick=? WHERE account_id=? AND online_date<>?")) {
            ps.setString(1,today); ps.setLong(2,System.currentTimeMillis()); ps.setInt(3,userId); ps.setString(4,today); ps.executeUpdate();
        }
        ensuredDate.put(userId, today);
    }

    /** Đẩy số giây đang chờ của một tài khoản xuống CSDL. */
    private void flush(int userId) {
        long now = System.currentTimeMillis();
        lastFlush.put(userId, now);
        Integer pending = pendingSeconds.remove(userId);
        if (pending == null || pending <= 0) return;
        try {
            ensureRow(userId);
            try (Connection c=connection(); PreparedStatement ps=c.prepareStatement("UPDATE phuc_loi_progress SET online_seconds=LEAST(online_seconds+?,3600),last_tick=? WHERE account_id=?")) {
                ps.setInt(1,pending); ps.setLong(2,now); ps.setInt(3,userId); ps.executeUpdate();
            }
        } catch (Exception e) {
            // Ghi hỏng thì trả lại số giây vào hàng chờ để lần sau ghi tiếp,
            // người chơi không mất thời gian online đã tích được.
            pendingSeconds.merge(userId, pending, Integer::sum);
            Logger.logException(PhucLoiService.class,e);
        }
    }

    /**
     * Gọi khi người chơi rời game: ghi nốt số giây còn chờ và dọn bộ nhớ.
     * Không có bước này thì mỗi lần thoát sẽ mất tối đa 60 giây online,
     * và các bảng đếm trong bộ nhớ sẽ phình dần theo số tài khoản đã vào.
     */
    public void onPlayerLogout(Player p) {
        try {
            if (p == null || p.getSession() == null) return;
            int userId = p.getSession().userId;
            if (userId <= 0) return;
            flush(userId);
            ticks.remove(userId);
            lastFlush.remove(userId);
            ensuredDate.remove(userId);
        } catch (Exception e) {
            Logger.logException(PhucLoiService.class,e);
        }
    }

    public void tick(Player p) {
        maybeLoadPanelConfig();
        if (p==null || p.getSession()==null) return;
        int userId = p.getSession().userId;
        if (userId <= 0) return;
        long now=System.currentTimeMillis(), last=ticks.getOrDefault(userId,now);
        if (now-last<1000) { ticks.putIfAbsent(userId,now); return; }
        ticks.put(userId,now);
        int add=(int)Math.min(2,(now-last)/1000);
        if (add>0) pendingSeconds.merge(userId, add, Integer::sum);
        if (now - lastFlush.getOrDefault(userId, now) >= FLUSH_INTERVAL_MS) flush(userId);
        else lastFlush.putIfAbsent(userId, now);
    }
    public void recordBossKill(Player p, Boss boss) {
        if(p==null||boss==null||p.getSession()==null)return; String n=boss.getClass().getSimpleName(); String col=n.equals("Broly")?"broly":n.equals("SO4")?"so4":n.equals("XenBoHung")?"xen":null; if(col==null)return;
        try { ensureRow(p.getSession().userId); try(Connection c=connection();Statement s=c.createStatement()){s.executeUpdate("UPDATE phuc_loi_progress SET "+col+"=LEAST("+col+"+1,10) WHERE account_id="+p.getSession().userId);} }catch(Exception e){Logger.logException(PhucLoiService.class,e);}
    }
    private int[] row(Player p) throws SQLException {
        // Đẩy phần đang chờ xuống trước khi đọc, để số liệu hiện ra luôn đúng.
        flush(p.getSession().userId);
        ensureRow(p.getSession().userId); int[] r=new int[9];
        try(Connection c=connection();PreparedStatement ps=c.prepareStatement("SELECT online_seconds,online_claim,nap_claim,broly,so4,xen,boss_claim,power_claim,task_claim FROM phuc_loi_progress WHERE account_id=?")){ps.setInt(1,p.getSession().userId);try(ResultSet rs=ps.executeQuery()){if(rs.next())for(int i=0;i<9;i++)r[i]=rs.getInt(i+1);}} return r;
    }
    private int progress(Player p,int[] r,int tab,int m){ if(tab==0)return r[0]; if(tab==1)return p.getSession().danap; if(tab==2)return r[3+m]; if(tab==3)return (int)Math.min(100,p.nPoint.power/1_000_000_000L); return p.playerTask==null||p.playerTask.taskMain==null?0:p.playerTask.taskMain.id; }
    private int mask(int[]r,int tab){return tab==0?r[1]:tab==1?r[2]:tab==2?r[6]:tab==3?r[7]:r[8];}
    private String maskColumn(int tab){return tab==0?"online_claim":tab==1?"nap_claim":tab==2?"boss_claim":tab==3?"power_claim":"task_claim";}
    private String info(int tab,int m){ if(tab==0)return new String[]{"Online 5 phút","Online 15 phút","Online 60 phút"}[m]; if(tab==1)return "Đạt Mốc "+LIMITS[tab][m]+" Cá"; if(tab==2)return new String[]{"Tiêu diệt 5 Broly","Tiêu diệt 5 Số 4","Tiêu diệt 10 Xên Bọ Hung"}[m]; if(tab==3)return "Đạt "+LIMITS[tab][m]+" tỷ sức mạnh"; return taskName(new int[]{10,21,26}[m]); }
    private String taskName(int id) { for (task.TaskMain task : server.Manager.TASKS) if (task.id == id) return task.name; return "Nhiệm vụ " + id; }
    private int globalId(int tab,int m){int id=m;for(int i=0;i<tab;i++)id+=LIMITS[i].length;return id;}
    private int[][] options(int tab,int m,int item){
        if (CUSTOM_OPTIONS != null && tab < CUSTOM_OPTIONS.length && m < CUSTOM_OPTIONS[tab].length && item < CUSTOM_OPTIONS[tab][m].length && CUSTOM_OPTIONS[tab][m][item] != null) return CUSTOM_OPTIONS[tab][m][item];
        if(tab==0&&m==0)return new int[][]{{new int[]{101,95,96,97}[item],5}};
        if(tab==0)return new int[][]{{30,1}};
        if(tab==1&&m==0)return item<2?new int[][]{{30,1}}:new int[0][0];
        if(tab==1&&m==1)return item==0?new int[][]{{30,1}}:new int[0][0];
        if(tab==1&&m==2)return item==0?new int[][]{{50,25},{77,25},{103,25},{5,25}}:new int[0][0];
        if(tab==2)return new int[0][0];
        if(tab==3)return new int[][]{{30,1}};
        if(tab==4)return m<2?new int[][]{{30,1}}:new int[0][0];
        return new int[0][0];
    }
    private File panelConfigFile() {
        String gameRoot = System.getProperty("game.root");
        if (gameRoot != null && !gameRoot.trim().isEmpty()) {
            File configured = new File(gameRoot.trim(), "phucloi_config.properties");
            if (configured.isFile()) return configured;
        }
        String userDir = System.getProperty("user.dir");
        if (userDir != null && !userDir.trim().isEmpty()) {
            File workingDirectory = new File(userDir, "phucloi_config.properties");
            if (workingDirectory.isFile()) return workingDirectory;
        }
        return new File("C:\\Users\\Admin\\Music\\zz1", "phucloi_config.properties");
    }
    /**
     * Dò lại file cấu hình, nhưng nhiều nhất 5 giây một lần cho cả server.
     *
     * loadPanelConfig() là hàm synchronized. Bản cũ gọi nó mỗi giây từ luồng
     * riêng của TỪNG người chơi, nên mọi luồng phải xếp hàng chờ nhau ở đây
     * mỗi giây, kèm 3 lần hỏi hệ điều hành về file. Càng đông càng tắc.
     */
    private void maybeLoadPanelConfig() {
        long now = System.currentTimeMillis();
        if (now - configLastCheck < CONFIG_CHECK_INTERVAL_MS) return;
        configLastCheck = now;
        loadPanelConfig();
    }
    private synchronized void loadPanelConfig() {
        try {
            File file = panelConfigFile(); long modified = file.isFile() ? file.lastModified() : -1L; long length = file.isFile() ? file.length() : -1L;
            if (modified == configLastModified && length == configLastLength) return;
            if (modified < 0) { useSourceRewards(); configLastModified=modified; configLastLength=length; return; }
            Properties props = new Properties(); try (FileInputStream in = new FileInputStream(file)) { props.load(in); }
            if (!"true".equalsIgnoreCase(props.getProperty("enabled"))) { useSourceRewards(); configLastModified=modified; configLastLength=length; return; }
            int[][][] ids = new int[5][][], qty = new int[5][][]; int[][][][][] ops = new int[5][][][][];
            for (int t=0;t<5;t++) { ids[t]=new int[LIMITS[t].length][]; qty[t]=new int[LIMITS[t].length][]; ops[t]=new int[LIMITS[t].length][][][];
                for (int m=0;m<LIMITS[t].length;m++) { String value=props.getProperty("reward."+t+"."+m, "").trim();
                    if (value.isEmpty()) { ids[t][m]=ITEM_IDS[t][m]; qty[t][m]=QTY[t][m]; ops[t][m]=new int[ITEM_IDS[t][m].length][][]; continue; }
                    String[] entries=value.split(";"); ids[t][m]=new int[entries.length]; qty[t][m]=new int[entries.length]; ops[t][m]=new int[entries.length][][];
                    for(int i=0;i<entries.length;i++) { String[] parts=entries[i].split(",",3); ids[t][m][i]=Short.parseShort(parts[0].trim()); qty[t][m][i]=Math.max(1,Integer.parseInt(parts[1].trim()));
                        if(parts.length<3||parts[2].trim().isEmpty()||"-".equals(parts[2].trim())) {ops[t][m][i]=new int[0][0];continue;}
                        String[] raw=parts[2].split("\\|");ops[t][m][i]=new int[raw.length][2];for(int o=0;o<raw.length;o++){String[] pair=raw[o].split(":",2);ops[t][m][i][o][0]=Integer.parseInt(pair[0].trim());ops[t][m][i][o][1]=Integer.parseInt(pair[1].trim());}
                    }
                }
            }
            ITEM_IDS=ids; QTY=qty; CUSTOM_OPTIONS=ops; configLastModified=modified; configLastLength=length;
        } catch(Exception e) { Logger.logException(PhucLoiService.class,e); }
    }
    private void useSourceRewards() { ITEM_IDS=SOURCE_ITEM_IDS; QTY=SOURCE_QTY; CUSTOM_OPTIONS=null; }
    // Lúc người chơi mở bảng Phúc Lợi thì đọc lại cấu hình ngay, không chờ
    // nhịp 5 giây — để admin vừa sửa phần thưởng ở panel là thấy hiệu lực liền.
    public void sendData(Player p){ loadPanelConfig(); tick(p); Message msg=null; try{int[]r=row(p);msg=new Message(103);msg.writer().writeByte(1);msg.writer().writeByte(5);msg.writer().writeByte(TOTAL);for(int t=0;t<5;t++){msg.writer().writeUTF(TAB_NAMES[t]);msg.writer().writeUTF("Phần thưởng chỉ nhận một lần khi đủ điều kiện");msg.writer().writeUTF(t==4?taskName(progress(p,r,t,0)):t==0?"giây":t==1?"Cá":t==3?"tỷ":"tiến độ");msg.writer().writeInt(progress(p,r,t,0));msg.writer().writeByte(LIMITS[t].length);for(int m=0;m<LIMITS[t].length;m++){msg.writer().writeInt(globalId(t,m));msg.writer().writeUTF(t==2 ? progress(p,r,t,m)+"/"+LIMITS[t][m]+" "+new String[]{"Broly","Số 4","Xên Bọ Hung"}[m] : info(t,m));msg.writer().writeInt(LIMITS[t][m]);msg.writer().writeInt(progress(p,r,t,m));msg.writer().writeInt((mask(r,t)>>m)&1);msg.writer().writeByte(ITEM_IDS[t][m].length);for(int i=0;i<ITEM_IDS[t][m].length;i++){int id=ITEM_IDS[t][m][i];msg.writer().writeShort(id);msg.writer().writeInt(QTY[t][m][i]);msg.writer().writeUTF("");msg.writer().writeUTF("");int[][]ops=options(t,m,i);msg.writer().writeByte(ops.length);for(int[]op:ops){msg.writer().writeByte(op[0]);msg.writer().writeInt(op[1]);}}}}p.sendMessage(msg);}catch(Exception e){Logger.logException(PhucLoiService.class,e);}finally{if(msg!=null)msg.cleanup();}}
    public synchronized void claim(Player p,int id){int tab=0,m=id;while(tab<5&&m>=LIMITS[tab].length)m-=LIMITS[tab++].length;if(tab>=5)return;try{int[]r=row(p);if(((mask(r,tab)>>m)&1)!=0){Service.gI().sendThongBao(p,"Mốc này đã nhận");return;}if(progress(p,r,tab,m)<LIMITS[tab][m]){Service.gI().sendThongBao(p,"Chưa đủ điều kiện nhận");return;}if(InventoryService.gI().getCountEmptyBag(p)<ITEM_IDS[tab][m].length){Service.gI().sendThongBao(p,"Hành trang không đủ chỗ trống");return;}try(Connection c=connection();PreparedStatement ps=c.prepareStatement("UPDATE phuc_loi_progress SET "+maskColumn(tab)+"="+maskColumn(tab)+"|? WHERE account_id=?")){ps.setInt(1,1<<m);ps.setInt(2,p.getSession().userId);ps.executeUpdate();}for(int i=0;i<ITEM_IDS[tab][m].length;i++){Item it=ItemService.gI().createNewItem((short)ITEM_IDS[tab][m][i],QTY[tab][m][i]);for(int[]op:options(tab,m,i))it.itemOptions.add(new Item.ItemOption(op[0],op[1]));InventoryService.gI().addItemBag(p,it);}InventoryService.gI().sendItemBag(p);Service.gI().sendThongBao(p,"Nhận quà thành công");sendData(p);}catch(Exception e){Logger.logException(PhucLoiService.class,e);}}
}
