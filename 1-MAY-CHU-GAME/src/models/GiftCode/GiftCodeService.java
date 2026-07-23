package models.GiftCode;

/*
 *
 *
 * @author GiaCatLuong
 */
import consts.ConstNpc;
import item.Item;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.util.ArrayList;
import java.util.Set;

import jdbc.DBConnecter;
import jdbc.daos.NDVSqlFetcher;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.JSONValue;
import player.Player;
import services.ItemService;
import services.NpcService;
import services.Service;
import utils.Logger;

public class GiftCodeService {

    private static GiftCodeService instance;

    public static GiftCodeService gI() {
        if (instance == null) {
            instance = new GiftCodeService();
        }
        return instance;
    }

    public void updateGiftCode() {
        PreparedStatement ps = null;
        ResultSet rs = null;
        try (Connection con2 = DBConnecter.getConnectionServer();) {
            GiftCodeManager.gI().listGiftCode.clear();
            ps = con2.prepareStatement("SELECT * FROM giftcode");
            rs = ps.executeQuery();
            while (rs.next()) {
                GiftCode giftcode = new GiftCode();
                giftcode.code = rs.getString("code");
                giftcode.id = rs.getInt("id");
                giftcode.countLeft = rs.getInt("count_left");
                if (giftcode.countLeft == -1) {
                    giftcode.countLeft = 999999999;
                }
                giftcode.datecreate = rs.getTimestamp("datecreate");
                giftcode.dateexpired = rs.getTimestamp("expired");
                JSONArray jar = (JSONArray) JSONValue.parse(rs.getString("detail"));
                if (jar != null) {
                    for (int i = 0; i < jar.size(); ++i) {
                        JSONObject jsonObj = (JSONObject) jar.get(i);

                        int id = Integer.parseInt(jsonObj.get("temp_id").toString());
                        int quantity = Integer.parseInt(jsonObj.get("quantity").toString());

                        JSONArray option = (JSONArray) jsonObj.get("options");
                        ArrayList<Item.ItemOption> optionList = new ArrayList<>();

                        if (option != null) {
                            for (int u = 0; u < option.size(); u++) {
                                JSONObject jsonobject = (JSONObject) option.get(u);
                                int optionId = Integer.parseInt(jsonobject.get("id").toString());
                                int param = Integer.parseInt(jsonobject.get("param").toString());
                                optionList.add(new Item.ItemOption(optionId, param));
                            }
                        }
                        giftcode.option.put(id, optionList);
                        giftcode.detail.put(id, quantity);
                    }
                }
                giftcode.type = rs.getInt("type");
                GiftCodeManager.gI().listGiftCode.add(giftcode);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void giftCode(Player player, String code) {
        try {
            // Cập nhật danh sách GiftCode
            updateGiftCode();
            // Kiểm tra và lấy thông tin GiftCode
            GiftCode giftcode = GiftCodeManager.gI().checkUseGiftCode(player, code);
            if (giftcode == null) {
                Service.gI().sendThongBao(player, "GiftCode đã được sử dụng hoặc không tồn tại.");
                return;
            }
            if (giftcode.timeCode()) {
                Service.gI().sendThongBao(player, "Code đã hết hạn");
                return;
            }
            // Xử lý phần thưởng
            Set<Integer> keySet = giftcode.detail.keySet();
            StringBuilder textGift = new StringBuilder("|7|Bạn Nhận Được:\b");
            for (Integer key : keySet) {
                int idItem = key;
                int quantity = giftcode.detail.get(key);

                switch (idItem) {
                    case -1 -> {
                        player.inventory.gold = Math.min(player.inventory.gold + (long) quantity, 2000000000L);
                        textGift.append("|2|").append(quantity).append(" vàng\b");
                    }
                    case -2 -> {
                        player.inventory.gem = Math.min(player.inventory.gem + quantity, 20000);
                        textGift.append("|3|").append(quantity).append(" ngọc\b");
                    }
                    case -3 -> {
                        player.inventory.ruby = Math.min(player.inventory.ruby + quantity, 20000);
                        textGift.append("|4|").append(quantity).append(" ngọc khóa\b");
                    }
                    default -> {
                        Item itemGiftTemplate = ItemService.gI().createNewItem((short) idItem);
                        if (itemGiftTemplate != null) {
                            Item itemGift = new Item((short) idItem);
                            itemGift.itemOptions = giftcode.option.get(key);
                            itemGift.quantity = quantity;
                            player.inventory.itemsMailBox.add(itemGift);
                            if (NDVSqlFetcher.updateMailBox(player)) {
                                textGift.append("|0|").append(quantity).append(" ").append(itemGift.template.name)
                                        .append("\b");
                            } else {
                                Logger.warning("Lỗi khi cập nhật hòm thư cho người chơi: " + player);
                            }
                        }
                    }
                }
            }
            NpcService.gI().createMenuConMeo(player, ConstNpc.IGNORE_MENU, -1, textGift + "Vui lòng kiểm tra hòm thư",
                    "OK");
        } catch (Exception e) {
            Service.gI().sendThongBao(player, "Đã xảy ra lỗi khi xử lý GiftCode. Vui lòng thử lại sau.");
        }
    }

}
