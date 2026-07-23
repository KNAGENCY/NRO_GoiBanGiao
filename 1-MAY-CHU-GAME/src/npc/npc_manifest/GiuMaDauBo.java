package npc.npc_manifest;

/**
 * NPC Gộp:
 * - Nhận quà nhiệm vụ
 * - Nhận quà sức mạnh
 * - Nhận quà mốc nạp (miễn phí)
 * - Bán đồ thần linh
 * - Nhận người yêu
 *
 * @author GiaCatLuong
 */

import consts.ConstNpc;
import item.Item;
import jdbc.DBConnecter;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;

import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

public class GiuMaDauBo extends Npc {

    /* ================== MENU ID ================== */
    private static final int MENU_NHIEM_VU = 1001;
    private static final int MENU_SUC_MANH = 1002;
    private static final int MENU_NHAN_QUA_NAP = 1003;

    /* ================== MỐC ================== */
    private static final int[] MOC_NHIEM_VU = { 21, 22, 23, 28 };
    private static final int[] MOC_SUC_MANH = { 11, 22, 33, 44 };
    private static final int[] MOC_NAP = { 50, 100, 300, 500, 1000 };

    public GiuMaDauBo(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    /* ================== MENU CHÍNH ================== */
    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player))
            return;

     createOtherMenu(player, ConstNpc.BASE_MENU,
        "Này nhóc:\r\n" +
        "Biết tính đạo hàm tích phân rồi\r\n" +
        "Thì biết tính cho hàm ẩn\r\n" +
        "Hàm logalepe, Hàm lượng giác chưa ? \r\n" +
        "Biết đếm vân sáng, vân tối\r\n" +
        "Trong giao Thoa ánh sáng k ?\r\n" +
        "Cho mấy chục viên bi xanh vàng đoả\r\n" +
        "Bốc bừa 3 viên thì có biết tính xác suất k ?\r\n",

        "Quà\nNhiệm vụ",
        "Quà\nSức mạnh",
        "Nhận\nNgười Yêu",
        "Sự Kiện\nTắm Biển"
);
    }

    /* ================== XỬ LÝ MENU ================== */
    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player))
            return;

        if (player.iDMark.isBaseMenu()) {
            switch (select) {
                case 0 -> {
                    createOtherMenu(player, ConstNpc.BASE_MENU,
                            "Nhận quà nhiệm vụ\n(Mỗi mốc 1 lần)\nPhải xong nhiệm vụ mới được nhận",
                            "Nappa x22\nHộp Quà",
                            "KuKu x11\nBánh Test",
                            "TĐST x1\nThạch Thần",
                            "Xên 8 x111\nThỏi Vàng"
                            );
                    player.iDMark.setIndexMenu(MENU_NHIEM_VU);
                }
                case 1 -> {
                    createOtherMenu(player, ConstNpc.BASE_MENU,
                            "Nhận quà sức mạnh\n(Mỗi mốc 1 lần)",
                            "11tỷ x22\nPhiếu Ngon",
                            "22tỷ x11\nBánh Chưng",
                            "33tỷ x1\nTrứng Mabư",
                            "44tỷ x222\nThức Ăn"
                          
                        );
                    player.iDMark.setIndexMenu(MENU_SUC_MANH);
                }
                // case 2 -> {
                //     createOtherMenu(player, ConstNpc.BASE_MENU,
                //             "Nhận quà mốc nạp\nKhông bị trừ tiền\n(Mỗi mốc 1 lần)",
                //             "50k = x1\nĐệ Mabư",
                //             "100k = x10\n3 Sao",
                //             "300k = x20\nThạch Thần",
                //             "500k = x1\nĐệ Black",
                //             "1000k= x1\nĐệ Beerus"
                //             );
                //     player.iDMark.setIndexMenu(MENU_NHAN_QUA_NAP);
                // }
            //     case 3 -> Service.gI().sendThongBaoOK(player,
            //             "Liên hệ Admin để bán đồ Thần Linh\n1 món = 10k ATM hoặc 30k Game");
              case 2 -> Service.gI().sendThongBaoOK(player, "Đang phát triển");

case 3 -> Service.gI().sendThongBaoOK(player, "Sẽ update sau 7 ngày Open");
             }
            return;
        }

        if (player.iDMark.getIndexMenu() == MENU_NHIEM_VU && select < MOC_NHIEM_VU.length) {
            nhanQuaNhiemVu(player, select);
            return;
        }

        if (player.iDMark.getIndexMenu() == MENU_SUC_MANH && select < MOC_SUC_MANH.length) {
            nhanQuaSucManh(player, select);
            return;
        }

        if (player.iDMark.getIndexMenu() == MENU_NHAN_QUA_NAP && select < MOC_NAP.length) {
            nhanQuaNap(player, MOC_NAP[select]);
        }
    }

    /* ================== NHIỆM VỤ ================== */
    private void nhanQuaNhiemVu(Player player, int index) {
        int moc = MOC_NHIEM_VU[index];
        int[] itemIds = { 648, 752, 674, 457 };
        int[] quantities = { 22, 11, 1, 111 };

        try (Connection con = DBConnecter.getConnectionServer()) {
            PreparedStatement ps = con.prepareStatement(
                    "SELECT reward_moc FROM player WHERE id = ?");
            ps.setLong(1, player.id);
            ResultSet rs = ps.executeQuery();
            if (!rs.next())
                return;

            JSONObject json = rs.getString("reward_moc") == null
                    ? new JSONObject()
                    : (JSONObject) JSONValue.parse(rs.getString("reward_moc"));

            String key = "task_" + moc;

            if (player.playerTask.taskMain.id < moc) {
                Service.gI().sendThongBao(player, "Chưa xong nhiệm vụ " + moc);
                return;
            }

            if (json.containsKey(key)) {
                Service.gI().sendThongBao(player, "Đã nhận mốc này rồi");
                return;
            }

            if (InventoryService.gI().getCountEmptyBag(player) < 1) {
                Service.gI().sendThongBao(player, "Hành trang đầy");
                return;
            }

            giveItem(player, itemIds[index], quantities[index]);

            json.put(key, true);
            updatePlayerJson(con, json, player.id);
            Service.gI().sendThongBao(player, "Nhận thành công mốc NV " + moc);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /* ================== SỨC MẠNH ================== */
    private void nhanQuaSucManh(Player player, int index) {
        int moc = MOC_SUC_MANH[index];
        int[] itemIds = { 1194, 753, 568, 663 };
        int[] quantities = { 22, 11, 1, 222 };

        long need = (long) moc * 1_000_000_000L;

        try (Connection con = DBConnecter.getConnectionServer()) {
            PreparedStatement ps = con.prepareStatement(
                    "SELECT reward_moc FROM player WHERE id = ?");
            ps.setLong(1, player.id);
            ResultSet rs = ps.executeQuery();
            if (!rs.next())
                return;

            JSONObject json = rs.getString("reward_moc") == null
                    ? new JSONObject()
                    : (JSONObject) JSONValue.parse(rs.getString("reward_moc"));

            String key = "power_" + moc;

            if (player.nPoint.power < need) {
                Service.gI().sendThongBao(player, "Chưa đủ " + moc + " tỷ sức mạnh");
                return;
            }

            if (json.containsKey(key)) {
                Service.gI().sendThongBao(player, "Đã nhận mốc này rồi");
                return;
            }

            if (InventoryService.gI().getCountEmptyBag(player) < 1) {
                Service.gI().sendThongBao(player, "Hành trang đầy");
                return;
            }

            giveItem(player, itemIds[index], quantities[index]);

            json.put(key, true);
            updatePlayerJson(con, json, player.id);
            Service.gI().sendThongBao(player, "Nhận thành công mốc " + moc + " tỷ");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /* ================== NẠP ================== */
    private void nhanQuaNap(Player player, int moc) {
        try (Connection con = DBConnecter.getConnectionServer()) {
            PreparedStatement ps = con.prepareStatement(
                    "SELECT danap, reward_nap_moc FROM account WHERE id = ?");
            ps.setInt(1, player.getSession().userId);
            ResultSet rs = ps.executeQuery();
            if (!rs.next())
                return;

            if (rs.getLong("danap") < (long) moc * 1000) {
                Service.gI().sendThongBao(player, "Chưa đạt mốc " + moc + "k");
                return;
            }

            JSONObject json = rs.getString("reward_nap_moc") == null
                    ? new JSONObject()
                    : (JSONObject) JSONValue.parse(rs.getString("reward_nap_moc"));

            if (json.containsKey(String.valueOf(moc))) {
                Service.gI().sendThongBao(player, "Đã nhận mốc này rồi");
                return;
            }

        switch (moc) {
    case 50 -> giveItem(player, 568, 1);

    case 100 -> giveItem(player, 16, 10);

    case 300 -> 
        giveItem(player, 674, 20);     
       

    case 500 -> giveItem(player, 1717, 1);

    case 1000 -> giveItem(player, 1716, 1);

}

            json.put(String.valueOf(moc), true);

            PreparedStatement ps2 = con.prepareStatement(
                    "UPDATE account SET reward_nap_moc = ? WHERE id = ?");
            ps2.setString(1, json.toJSONString());
            ps2.setInt(2, player.getSession().userId);
            ps2.executeUpdate();

            InventoryService.gI().sendItemBag(player);
            Service.gI().sendThongBao(player, "Nhận thành công mốc " + moc + "k");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /* ================== HỖ TRỢ ================== */
    private void giveItem(Player player, int tempId, int quantity) {
        Item item = ItemService.gI().createNewItem((short) tempId, quantity);
        item.itemOptions.add(new Item.ItemOption(30, 1));
        InventoryService.gI().addItemBag(player, item);
        InventoryService.gI().sendItemBag(player);
    }

    private void updatePlayerJson(Connection con, JSONObject json, long playerId) throws Exception {
        PreparedStatement ps = con.prepareStatement(
                "UPDATE player SET reward_moc = ? WHERE id = ?");
        ps.setString(1, json.toJSONString());
        ps.setLong(2, playerId);
        ps.executeUpdate();
    }
}
