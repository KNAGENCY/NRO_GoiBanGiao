package models.Consign;

import jdbc.DBConnecter;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.util.ArrayList;
import java.util.List;

import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

import item.Item.ItemOption;

public class ConsignShopManager {

    private static ConsignShopManager instance;

    public static ConsignShopManager gI() {
        if (instance == null) {
            instance = new ConsignShopManager();
        }
        return instance;
    }

    public long lastTimeUpdate;

    public String[] tabName = {
        "Trang bị",
        "Phụ kiện",
        "Hỗ trợ",
        "Linh tinh",
        ""
    };

    public final List<ConsignItem> listItem = new ArrayList<>();

    // ===================== LOAD =====================

    public void load() {

        final String SELECT_SQL =
                "SELECT id, player_id, tab, item_id, gold, gem, quantity, " +
                "itemOption, lasttime, isBuy FROM shop_ky_gui";

        List<ConsignItem> loaded = new ArrayList<>();

        try (
                Connection con = DBConnecter.getConnectionServer();
                PreparedStatement ps = con.prepareStatement(SELECT_SQL);
                ResultSet rs = ps.executeQuery()
        ) {

            while (rs.next()) {

                ConsignItem it = new ConsignItem();

                it.id = rs.getInt("id");
                it.player_sell = rs.getInt("player_id");
                it.tab = rs.getByte("tab");
                it.itemId = rs.getShort("item_id");

                it.goldSell = rs.getInt("gold");
                it.gemSell = rs.getInt("gem");

                it.quantity = rs.getInt("quantity");

                it.lasttime = rs.getLong("lasttime");

                it.isBuy = rs.getInt("isBuy") == 1;

                it.options = parseOptions(rs.getString("itemOption"));

                loaded.add(it);
            }

            synchronized (listItem) {
                listItem.clear();
                listItem.addAll(loaded);
            }

            System.out.println("[CONSIGN-LOAD] loaded " + loaded.size() + " items");

        } catch (Exception e) {

            System.out.println(
                    "[CONSIGN-LOAD-ERR] "
                    + e.getClass().getSimpleName()
                    + ": "
                    + e.getMessage()
            );

            e.printStackTrace();
        }
    }

    // ===================== SAVE =====================

    /**
     * SAVE KIỂU MỚI:
     * - XÓA TOÀN BỘ TABLE
     * - INSERT LẠI TOÀN BỘ listItem hiện tại
     *
     * => FIX TRIỆT ĐỂ lỗi:
     * item đã mua nhưng reboot vẫn hiện lại
     */
    public void save() {

        final String TRUNCATE_SQL = "TRUNCATE TABLE shop_ky_gui";

        final String INSERT_SQL =
                "INSERT INTO shop_ky_gui " +
                "(id, player_id, tab, item_id, gold, gem, quantity, itemOption, lasttime, isBuy) " +
                "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)";

        synchronized (listItem) {

            try (Connection con = DBConnecter.getConnectionServer()) {

                con.setAutoCommit(false);

                try {

                    // XÓA HẾT DB CŨ
                    try (PreparedStatement ps = con.prepareStatement(TRUNCATE_SQL)) {
                        ps.executeUpdate();
                    }

                    // INSERT TOÀN BỘ LIST HIỆN TẠI
                    int rows = insertAll(con, INSERT_SQL);

                    con.commit();

                    System.out.println(
                            "[CONSIGN-SAVE] committed "
                            + rows
                            + " rows"
                    );

                } catch (Exception inner) {

                    con.rollback();

                    throw inner;
                }

            } catch (Exception e) {

                System.out.println(
                        "[CONSIGN-SAVE-ERR] "
                        + e.getClass().getSimpleName()
                        + ": "
                        + e.getMessage()
                );

                e.printStackTrace();
            }
        }
    }

    // ===================== INSERT =====================

    private int insertAll(Connection con, String sql) throws Exception {

        int count = 0;

        try (PreparedStatement ps = con.prepareStatement(sql)) {

            for (ConsignItem it : listItem) {

                if (it == null) {
                    continue;
                }

                String options = normalizeOptions(it.options);

                int isBuy = it.isBuy ? 1 : 0;

                ps.setInt(1, it.id);

                ps.setInt(2, it.player_sell);

                ps.setByte(3, it.tab);

                ps.setShort(4, it.itemId);

                ps.setInt(5, it.goldSell);

                ps.setInt(6, it.gemSell);

                ps.setInt(7, it.quantity);

                ps.setString(8, options);

                ps.setLong(9, it.lasttime);

                ps.setInt(10, isBuy);

                ps.addBatch();

                count++;
            }

            ps.executeBatch();
        }

        return count;
    }

    // ===================== REMOVE ITEM =====================

    /**
     * XÓA ITEM KHỎI LIST + DB NGAY LẬP TỨC
     */
    public void removeItem(ConsignItem item) {

        if (item == null) {
            return;
        }

        synchronized (listItem) {

            listItem.remove(item);

            try (
                    Connection con = DBConnecter.getConnectionServer();
                    PreparedStatement ps =
                            con.prepareStatement(
                                    "DELETE FROM shop_ky_gui WHERE id = ?"
                            )
            ) {

                ps.setInt(1, item.id);

                ps.executeUpdate();

            } catch (Exception e) {

                e.printStackTrace();
            }
        }
    }

    // ===================== UPDATE BUY =====================

    /**
     * UPDATE trạng thái mua NGAY xuống DB
     */
    public void updateItem(ConsignItem item) {

        if (item == null) {
            return;
        }

        try (
                Connection con = DBConnecter.getConnectionServer();
                PreparedStatement ps =
                        con.prepareStatement(
                                "UPDATE shop_ky_gui SET isBuy = ?, lasttime = ? WHERE id = ?"
                        )
        ) {

            ps.setInt(1, item.isBuy ? 1 : 0);

            ps.setLong(2, item.lasttime);

            ps.setInt(3, item.id);

            ps.executeUpdate();

        } catch (Exception e) {

            e.printStackTrace();
        }
    }

    // ===================== HELPERS =====================

    private static String normalizeOptions(Object options) {

        String s = JSONValue.toJSONString(options);

        return (s == null || s.equals("null"))
                ? "[]"
                : s;
    }

    private static List<ItemOption> parseOptions(String json) {

        List<ItemOption> opts = new ArrayList<>();

        if (json == null
                || json.isEmpty()
                || json.equals("[]")
                || json.equals("null")) {

            return opts;
        }

        try {

            Object parsed = JSONValue.parse(json);

            if (!(parsed instanceof JSONArray)) {
                return opts;
            }

            for (Object o : (JSONArray) parsed) {

                if (!(o instanceof JSONObject)) {
                    continue;
                }

                JSONObject obj = (JSONObject) o;

                int id = Integer.parseInt(
                        String.valueOf(obj.get("id"))
                );

                int param = Integer.parseInt(
                        String.valueOf(obj.get("param"))
                );

                opts.add(new ItemOption(id, param));
            }

        } catch (Exception e) {

            System.out.println(
                    "[CONSIGN-LOAD-OPT-ERR] bad json="
                    + json
                    + " err="
                    + e.getMessage()
            );
        }

        return opts;
    }
}