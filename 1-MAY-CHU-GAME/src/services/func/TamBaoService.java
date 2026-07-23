package services.func;

import item.Item;
import java.io.IOException;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ThreadLocalRandom;
import jdbc.DBConnecter;
import network.Message;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import server.Client;
import utils.Logger;

public final class TamBaoService {
    private static final int PRICE_PER_SPIN = 2_000_000_000;
    private static final int[] ITEM_IDS = {15, 16, 18, 19, 20, 457, 579, 935, 1628, 987, 1150, 1151, 1152, 1153};
    private static final int[] MILESTONES = {100, 200, 500, 1000, 2000};
    private static final int[] CASH_REWARDS = {10_000, 22_000, 60_000, 130_000, 280_000};
    private static final int COUNT_MASK = (1 << 20) - 1;
    private static final TamBaoService INSTANCE = new TamBaoService();

    public static TamBaoService gI() { return INSTANCE; }
    private int count(Player p) { return p.getSession().luotquay & COUNT_MASK; }
    private int claims(Player p) { return (p.getSession().luotquay >>> 20) & 31; }
    private int packed(int count, int claims) { return (count & COUNT_MASK) | ((claims & 31) << 20); }

    public void sendData(Player p) {
        sendItems(p);
        sendMilestones(p);
        sendClaimedSlots(p, new int[ITEM_IDS.length]);
    }

    private void sendItems(Player p) {
        Message m = null;
        try {
            m = new Message(106); m.writer().writeByte(0); m.writer().writeByte(1);
            m.writer().writeShort(-1); m.writer().writeShort(-1); m.writer().writeLong(PRICE_PER_SPIN); m.writer().writeByte(ITEM_IDS.length);
            for (int id : ITEM_IDS) {
                m.writer().writeByte(0); m.writer().writeShort(id); m.writer().writeInt(1);
                m.writer().writeUTF(""); m.writer().writeUTF(""); m.writer().writeByte(0);
            }
            p.sendMessage(m);
        } catch (IOException e) { Logger.logException(TamBaoService.class, e); }
        finally { if (m != null) m.cleanup(); }
    }

    public void sendMilestones(Player p) {
        Message m = null;
        try {
            m = new Message(106); m.writer().writeByte(1); m.writer().writeInt(count(p));
            m.writer().writeByte(MILESTONES.length);
            int claimed = claims(p);
            for (int i = 0; i < MILESTONES.length; i++) {
                m.writer().writeInt((claimed >>> i) & 1); m.writer().writeShort(457);
                m.writer().writeInt(CASH_REWARDS[i]); m.writer().writeUTF(CASH_REWARDS[i] + " Cá");
                m.writer().writeUTF(""); m.writer().writeInt(i); m.writer().writeInt(MILESTONES[i]); m.writer().writeByte(0);
            }
            p.sendMessage(m);
        } catch (IOException e) { Logger.logException(TamBaoService.class, e); }
        finally { if (m != null) m.cleanup(); }
    }

    private void sendClaimedSlots(Player p, int[] received) {
        Message m = null;
        try {
            m = new Message(106); m.writer().writeByte(2); m.writer().writeByte(ITEM_IDS.length);
            for (int i = 0; i < ITEM_IDS.length; i++) m.writer().writeInt(received[i]);
            p.sendMessage(m);
        } catch (IOException e) { Logger.logException(TamBaoService.class, e); }
        finally { if (m != null) m.cleanup(); }
    }

    public synchronized void handle(Player p, Message request) throws IOException {
        int action = request.reader().readByte();
        if (action == 0) claimMilestone(p, request.reader().readInt());
        else if (action == 1) spin(p, request.reader().readInt()); else if (action == 2) sendTop(p); else if (action == 3) resetTop(p);
    }

    private void spin(Player p, int times) {
        if (times != 1 && times != 10) return;
        long cost = (long) PRICE_PER_SPIN * times;
        if (p.inventory.gold < cost) { Service.gI().sendThongBao(p, "Không đủ vàng để quay"); return; }
        if (InventoryService.gI().getCountEmptyBag(p) < times) { Service.gI().sendThongBao(p, "Hành trang không đủ " + times + " ô trống"); return; }
        p.inventory.gold -= cost;
        int[] received = new int[ITEM_IDS.length];
        StringBuilder rewardNames = new StringBuilder("Bạn nhận được ");
        for (int i = 0; i < times; i++) {
            int rewardIndex = ThreadLocalRandom.current().nextInt(ITEM_IDS.length);
            int id = ITEM_IDS[rewardIndex];
            received[rewardIndex]++;
            Item item = ItemService.gI().createNewItem((short) id, 1);
            InventoryService.gI().addItemBag(p, item);
            if (rewardNames.length() > "Bạn nhận được ".length()) rewardNames.append(", ");
            rewardNames.append(item.template.name);
        }
        p.getSession().luotquay = packed(count(p) + times, claims(p));
        saveSpinState(p);
        InventoryService.gI().sendItemBag(p); Service.gI().sendMoney(p);
        Service.gI().sendThongBao(p, rewardNames.toString());
        sendMilestones(p);
        sendClaimedSlots(p, received);
    }

    private void claimMilestone(Player p, int id) {
        if (id < 0 || id >= MILESTONES.length) return;
        int claimed = claims(p);
        if (((claimed >>> id) & 1) != 0) { Service.gI().sendThongBao(p, "Mốc này đã nhận"); return; }
        if (count(p) < MILESTONES[id]) { Service.gI().sendThongBao(p, "Cần " + MILESTONES[id] + " Điểm"); return; }
        int newClaims = claimed | (1 << id);
        int newPacked = packed(count(p), newClaims);
        try (Connection c = DBConnecter.getConnectionServer(); PreparedStatement ps = c.prepareStatement(
                "UPDATE account SET luotquay=?, cash=cash+?, danap=danap+? WHERE id=? AND luotquay=?")) {
            ps.setInt(1, newPacked); ps.setInt(2, CASH_REWARDS[id]); ps.setInt(3, CASH_REWARDS[id]);
            ps.setInt(4, p.getSession().userId); ps.setInt(5, p.getSession().luotquay);
            if (ps.executeUpdate() == 1) {
                p.getSession().luotquay = newPacked; p.getSession().cash += CASH_REWARDS[id]; p.getSession().danap += CASH_REWARDS[id];
                Service.gI().sendThongBao(p, "Nhận " + CASH_REWARDS[id] + " Cá thành công"); sendMilestones(p);
            }
        } catch (Exception e) { Logger.logException(TamBaoService.class, e); }
    }

    private void saveSpinState(Player p) {
        try (Connection c = DBConnecter.getConnectionServer(); PreparedStatement ps = c.prepareStatement("UPDATE account SET luotquay=? WHERE id=?")) {
            ps.setInt(1, p.getSession().luotquay); ps.setInt(2, p.getSession().userId); ps.executeUpdate();
        } catch (Exception e) { Logger.logException(TamBaoService.class, e); }
    }
    private boolean isTopAdmin(Player p) {
        try (Connection c = DBConnecter.getConnectionServer(); PreparedStatement ps = c.prepareStatement("SELECT admin, is_admin FROM account WHERE id=?")) {
            ps.setInt(1, p.getSession().userId);
            try (ResultSet rs = ps.executeQuery()) { return rs.next() && (rs.getInt(1) == 1 || rs.getInt(2) == 1); }
        } catch (Exception e) { Logger.logException(TamBaoService.class, e); return false; }
    }

    private void sendTop(Player p) {
        List<String> names = new ArrayList<>();
        List<Integer> points = new ArrayList<>();
        try (Connection c = DBConnecter.getConnectionServer(); PreparedStatement ps = c.prepareStatement(
                "SELECT player.name, (account.luotquay & 1048575) AS spin_points FROM account INNER JOIN player ON player.account_id=account.id ORDER BY spin_points DESC, player.name ASC LIMIT 20");
             ResultSet rs = ps.executeQuery()) {
            while (rs.next()) { names.add(rs.getString(1)); points.add(rs.getInt(2)); }
        } catch (Exception e) { Logger.logException(TamBaoService.class, e); }
        Message m = null;
        try {
            m = new Message(106); m.writer().writeByte(3); m.writer().writeByte(isTopAdmin(p) ? 1 : 0); m.writer().writeByte(names.size());
            for (int i = 0; i < names.size(); i++) { m.writer().writeUTF(names.get(i)); m.writer().writeInt(points.get(i)); }
            p.sendMessage(m);
        } catch (IOException e) { Logger.logException(TamBaoService.class, e); }
        finally { if (m != null) m.cleanup(); }
    }

    private void resetTop(Player p) {
        if (!isTopAdmin(p)) { Service.gI().sendThongBao(p, "Chi tai khoan admin moi duoc reset diem"); return; }
        try (Connection c = DBConnecter.getConnectionServer(); PreparedStatement ps = c.prepareStatement("UPDATE account SET luotquay=0")) {
            ps.executeUpdate();
            for (Player online : Client.gI().getPlayers()) {
                if (online != null && online.getSession() != null) { online.getSession().luotquay = 0; sendMilestones(online); }
            }
            Service.gI().sendThongBao(p, "Da reset diem vong quay");
            sendMilestones(p); sendTop(p);
        } catch (Exception e) { Logger.logException(TamBaoService.class, e); }
    }
}
