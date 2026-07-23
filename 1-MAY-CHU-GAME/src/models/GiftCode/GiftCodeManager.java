package models.GiftCode;

/*
 *
 *
 * @author GiaCatLuong
 */
import jdbc.DBConnecter;
import player.Player;
import services.NpcService;
import services.Service;
import services.InventoryService;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.util.ArrayList;
import java.util.HashMap;

import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

import item.Item.ItemOption;

public class GiftCodeManager {

    public String name;
    public final ArrayList<GiftCode> listGiftCode = new ArrayList<>();

    private static GiftCodeManager instance;

    public static GiftCodeManager gI() {
        if (instance == null) {
            instance = new GiftCodeManager();
        }
        return instance;
    }

   public void loadGiftCodeFromDB() {
    listGiftCode.clear();

    String sql = "SELECT * FROM giftcode";

    Connection con = null;
    PreparedStatement ps = null;
    ResultSet rs = null;

    try {
        con = DBConnecter.getConnectionServer();
        if (con == null || con.isClosed()) {
            System.out.println("Connection null hoặc đã đóng!");
            return;
        }

        ps = con.prepareStatement(sql);
        rs = ps.executeQuery();

        while (rs.next()) {
            GiftCode gc = new GiftCode();

            gc.id = rs.getInt("id");
            gc.code = rs.getString("code");
            gc.countLeft = rs.getInt("count_left");
            gc.datecreate = rs.getTimestamp("datecreate");
            gc.dateexpired = rs.getTimestamp("expired");
            gc.type = rs.getInt("type"); 

            gc.detail = new HashMap<>();
            gc.option = new HashMap<>();

            String detailJson = rs.getString("detail");

            if (detailJson != null && !detailJson.trim().isEmpty()) {
                Object parsed = JSONValue.parse(detailJson);

                if (parsed instanceof JSONArray) {
                    JSONArray jar = (JSONArray) parsed;

                    for (Object obj : jar) {
                        if (!(obj instanceof JSONObject)) {
                            continue;
                        }

                        JSONObject json = (JSONObject) obj;

                        Object itemIdObj = json.get("temp_id");
                        Object quantityObj = json.get("quantity");

                        if (itemIdObj == null || quantityObj == null) {
                            System.out.println("GiftCode item lỗi: " + json);
                            continue;
                        }

                        int itemId;
                        int quantity;

                        try {
                            itemId = Integer.parseInt(itemIdObj.toString());
                            quantity = Integer.parseInt(quantityObj.toString());
                        } catch (Exception e) {
                            System.out.println("GiftCode item parse lỗi: " + json);
                            continue;
                        }

                        ArrayList<ItemOption> optionList = new ArrayList<>();
                        Object optionsObj = json.get("options");

                        if (optionsObj instanceof JSONArray) {
                            JSONArray options = (JSONArray) optionsObj;

                            for (Object optObj : options) {
                                if (!(optObj instanceof JSONObject)) {
                                    continue;
                                }

                                JSONObject opt = (JSONObject) optObj;

                                Object optIdObj = opt.get("id");
                                Object paramObj = opt.get("param");

                                if (optIdObj == null || paramObj == null) {
                                    System.out.println("GiftCode option lỗi: " + opt);
                                    continue;
                                }

                                try {
                                    int optId = Integer.parseInt(optIdObj.toString());
                                    int param = Integer.parseInt(paramObj.toString());

                                    ItemOption io = new ItemOption(optId, param);
                                    if (io.optionTemplate != null) {
                                        optionList.add(io);
                                    }
                                } catch (Exception e) {
                                    System.out.println("GiftCode option parse lỗi: " + opt);
                                }
                            }
                        }

                        gc.detail.put(itemId, quantity);
                        gc.option.put(itemId, optionList);
                    }
                } else {
                    System.out.println("detail không phải JSONArray: " + detailJson);
                }
            }

            listGiftCode.add(gc);
        }

        System.out.println("Reload GiftCode: " + listGiftCode.size());

    } catch (Exception e) {
        e.printStackTrace();
    } finally {
        try {
            if (rs != null) rs.close();
        } catch (Exception ignored) {
        }
        try {
            if (ps != null) ps.close();
        } catch (Exception ignored) {
        }
        try {
            if (con != null) con.close();
        } catch (Exception ignored) {
        }
    }
}

    public GiftCode checkUseGiftCode(Player player, String code) {
        for (GiftCode giftCode : listGiftCode) {
            if (giftCode.code.equals(code)) {
                if (giftCode.countLeft <= 0) {
                    Service.gI().sendThongBao(player, "Giftcode đã hết lượt nhập");
                    return null;
                }
                if (giftCode.isUsedGiftCode(player)) {
                    Service.gI().sendThongBao(player, "Bạn đã sử dụng GiftCode này rồi");
                    return null;
                }
                if (giftCode.type == 1 && !player.getSession().actived) {
                    Service.gI().sendThongBao(player, "Bạn cần mở thành viên để sử dụng mã này.");
                    return null;
                }
                if (InventoryService.gI().getCountEmptyBag(player) < giftCode.detail.size()) {
                    Service.gI().sendThongBaoOK(player,
                            "Cần tối thiểu " + giftCode.detail.size() + " ô hành trang trống");
                    return null;
                }

                giftCode.countLeft -= 1;
                player.giftCode.add(code);
                updateGiftCode(giftCode);
                return giftCode;
            }
        }
        return null;
    }

    public void updateGiftCode(GiftCode giftcode) {
        try {
            DBConnecter.executeUpdate(
                    "update giftcode set count_left = ? where id = ?",
                    giftcode.countLeft,
                    giftcode.id);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void checkInfomationGiftCode(Player p) {
        StringBuilder sb = new StringBuilder();
        for (GiftCode giftCode : listGiftCode) {
            sb.append("Code: ")
                    .append(giftCode.code)
                    .append(", Số lượng còn lại: ")
                    .append(giftCode.countLeft)
                    .append("\n")
                    .append("Ngày tạo: ")
                    .append(giftCode.datecreate)
                    .append(", Ngày hết hạn: ")
                    .append(giftCode.dateexpired)
                    .append("\n");
        }

        if (sb.length() > 0) {
            sb.deleteCharAt(sb.length() - 1);
        }

        NpcService.gI().createTutorial(p, 5073, sb.toString());
    }
}