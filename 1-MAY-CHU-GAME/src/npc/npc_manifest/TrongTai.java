package npc.npc_manifest;

import consts.ConstNpc;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.util.ArrayList;
import java.util.List;
import jdbc.DBConnecter;
import models.SuperRank.SuperRankManager;
import models.SuperRank.SuperRankService;
import network.Message;
import npc.Npc;
import player.Player;
import services.NpcService;
import services.func.ChangeMapService;
import utils.Util;

public class TrongTai extends Npc {

    private static final int MENU_TOP_CHI_SO = 53100;

    private static final String TOP_SM = "SELECT name, gender, items_body, CAST(JSON_EXTRACT(data_point, '$[1]') AS UNSIGNED) AS value "
            + "FROM player INNER JOIN account ON account.id = player.account_id "
            + "WHERE account.ban = 0 ORDER BY value DESC LIMIT 10";
    private static final String TOP_NV = "SELECT name, gender, items_body, CAST(JSON_EXTRACT(data_task, '$[0]') AS UNSIGNED) AS nv, "
            + "CAST(JSON_EXTRACT(data_task, '$[1]') AS UNSIGNED) AS subnv, CAST(JSON_EXTRACT(data_task, '$[2]') AS UNSIGNED) AS countnv "
            + "FROM player INNER JOIN account ON account.id = player.account_id "
            + "WHERE account.ban = 0 ORDER BY nv DESC, subnv DESC, countnv DESC LIMIT 10";
    private static final String TOP_HP = "SELECT name, gender, items_body, CAST(split_str(data_point, ',', 6) AS UNSIGNED) AS value "
            + "FROM player INNER JOIN account ON account.id = player.account_id "
            + "WHERE account.is_admin = 0 AND account.ban = 0 ORDER BY value DESC LIMIT 10";
    private static final String TOP_KI = "SELECT name, gender, items_body, CAST(split_str(data_point, ',', 7) AS UNSIGNED) AS value "
            + "FROM player INNER JOIN account ON account.id = player.account_id "
            + "WHERE account.is_admin = 0 AND account.ban = 0 ORDER BY value DESC LIMIT 10";
    private static final String TOP_SD = "SELECT name, gender, items_body, CAST(split_str(data_point, ',', 8) AS UNSIGNED) AS value "
            + "FROM player INNER JOIN account ON account.id = player.account_id "
            + "WHERE account.is_admin = 0 AND account.ban = 0 ORDER BY value DESC LIMIT 10";
    private static final String TOP_SAN_BOSS = "SELECT name, gender, items_body, event_point_boss AS value "
            + "FROM player INNER JOIN account ON account.id = player.account_id "
            + "WHERE account.ban = 0 ORDER BY value DESC LIMIT 10";

    public TrongTai(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) {
            return;
        }
        switch (mapId) {
            case 113 -> openSuperRankMenu(player);
            case 5 -> this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "BXH Update đồng bộ 3 phút 1 lần liên tục\n"
                            + "Công bằng, minh bạch, thắng làm vua\n"
                            + "Kỹ năng yếu thì không được đổ thừa\n"
                            + "Nếu có buff bẩn, bug, hack ban thẳng tay kể cả người quen\n\n"
                            + "Số điểm săn boss hiện có: " + player.event.getEventPointBHM() + " điểm.",
                    "Top SM", "Top NV", "Top\nChỉ Số", "Top\nSăn Boss", "Top\nTrùm Cuối");
            default -> super.openBaseMenu(player);
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (!canOpenNpc(player)) {
            return;
        }
        if (player.iDMark.isBaseMenu()) {
            switch (mapId) {
                case 113 -> confirmSuperRankMenu(player, select);
                case 5 -> {
                    switch (select) {
                        case 0 ->
                            showTop(player, "Top SM", TOP_SM, "Sức mạnh", false);
                        case 1 ->
                            showTop(player, "Top NV", TOP_NV, "Nhiệm vụ", true);
                        case 2 ->
                            createOtherMenu(player, MENU_TOP_CHI_SO, "Ngươi muốn xem bảng xếp hạng chỉ số nào?",
                                    "Top HP", "Top KI", "Top SĐ", "Từ chối");
                        case 3 ->
                            showTop(player, "Top Săn Boss", TOP_SAN_BOSS, "Điểm săn boss", false);
                        case 4 ->
                            NpcService.gI().createTutorial(player, tempId, avartar,
                                    "Nhận 'Cá' Ib Key Vàng Trong Zalo");
                        default -> {
                        }
                    }
                }
                default -> {
                }
            }
        } else if (player.iDMark.getIndexMenu() == MENU_TOP_CHI_SO) {
            switch (select) {
                case 0 ->
                    showTop(player, "Top HP", TOP_HP, "HP", false);
                case 1 ->
                    showTop(player, "Top KI", TOP_KI, "KI", false);
                case 2 ->
                    showTop(player, "Top SĐ", TOP_SD, "Sức đánh", false);
                default -> {
                }
            }
        }
    }

    private void openSuperRankMenu(Player player) {
        if (SuperRankManager.gI().awaiting(player)) {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Vui lòng chờ, số thứ tự của bạn là " + SuperRankManager.gI().ordinal(player.id),
                    "OK", "Về\nĐại Hội\nVõ Thuật");
            return;
        }
        this.createOtherMenu(player, ConstNpc.BASE_MENU,
                "Đại hội võ thuật Siêu Hạng\ndiễn ra 24/7 kể cả ngày lễ và chủ nhật\n"
                        + "Hãy thi đấu ngay để khẳng định đẳng cấp của mình nhé",
                "Top 100\nCao Thủ", "Hướng\ndẫn\nthêm",
                player.superRank.ticket > 0 ? "Miễn phí\nCòn " + player.superRank.ticket + " vé" : "Thi đấu",
                "Ưu tiên\nđấu ngay", "Về\nĐảo Kame");
    }

    private void confirmSuperRankMenu(Player player, int select) {
        if (SuperRankManager.gI().awaiting(player)) {
            if (select == 1) {
                ChangeMapService.gI().changeMapNonSpaceship(player, 52, player.location.x, 336);
            }
            return;
        }
        switch (select) {
            case 0 ->
                SuperRankService.gI().topList(player, 0);
            case 1 ->
                NpcService.gI().createTutorial(player, tempId, avartar, ConstNpc.THONG_TIN_SIEU_HANG);
            case 2 ->
                SuperRankService.gI().topList(player, 1);
            case 3 ->
                SuperRankService.gI().topList(player, 2);
            case 4 ->
                ChangeMapService.gI().changeMapNonSpaceship(player, 5, player.location.x, 288);
            default -> {
            }
        }
    }

    private void showTop(Player player, String title, String sql, String label, boolean taskTop) {
        Message msg = null;
        try (Connection con = DBConnecter.getConnectionServer();
                PreparedStatement ps = con.prepareStatement(sql);
                ResultSet rs = ps.executeQuery()) {
            List<TopEntry> entries = new ArrayList<>();
            int rank = 1;
            while (rs.next() && rank <= 10) {
                entries.add(new TopEntry(
                        rank,
                        rs.getByte("gender"),
                        rs.getString("name"),
                        taskTop ? 0 : rs.getLong("value"),
                        taskTop ? rs.getInt("nv") : 0,
                        taskTop ? rs.getInt("subnv") : 0,
                        taskTop ? rs.getInt("countnv") : 0));
                rank++;
            }

            msg = new Message(-96);
            msg.writer().writeByte(0);
            msg.writer().writeUTF(title);
            msg.writer().writeByte(entries.size());
            for (TopEntry entry : entries) {
                msg.writer().writeInt(entry.rank);
                msg.writer().writeInt(entry.rank);
                msg.writer().writeShort(Util.getHead(entry.gender));
                if (player.session != null && player.session.version >= 214) {
                    msg.writer().writeShort(-1);
                }
                msg.writer().writeShort(entry.gender == 1 ? 59 : 57);
                msg.writer().writeShort(entry.gender == 1 ? 60 : 58);
                msg.writer().writeUTF(entry.name);
                if (taskTop) {
                    msg.writer().writeUTF(label + ": " + entry.nv + " - " + entry.subnv);
                    msg.writer().writeUTF("Tiến độ: " + entry.countnv);
                } else {
                    msg.writer().writeUTF(label + ": " + Util.numberToMoney(entry.value));
                    msg.writer().writeUTF(String.valueOf(entry.value));
                }
            }
            player.sendMessage(msg);
        } catch (Exception e) {
            NpcService.gI().createTutorial(player, tempId, avartar, "Không thể tải bảng xếp hạng lúc này.");
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    private static class TopEntry {

        private final int rank;
        private final byte gender;
        private final String name;
        private final long value;
        private final int nv;
        private final int subnv;
        private final int countnv;

        private TopEntry(int rank, byte gender, String name, long value, int nv, int subnv, int countnv) {
            this.rank = rank;
            this.gender = gender;
            this.name = name;
            this.value = value;
            this.nv = nv;
            this.subnv = subnv;
            this.countnv = countnv;
        }
    }
}
