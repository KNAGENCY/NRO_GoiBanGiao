package npc.npc_manifest;

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

public class Bardock extends Npc {

    private static final int MENU_GOI  = 10001;
    private static final int MENU_NHAN = 10002;

    private static final int GOI_50K  = 1;
    private static final int GOI_100K = 2;
    private static final int GOI_250K = 4;

    private static final int ONCE_TRUNG = 1;
    private static final int ONCE_SKIP  = 2;

    private static final int CASH_50K  = 50000;
    private static final int CASH_100K = 100000;
    private static final int CASH_250K = 250000;

    private static final long COOLDOWN = 20L * 60 * 60 * 1000;
    private static final int MAX_RECEIVE = 7;

    public Bardock(int mapId, int status, int cx, int cy, int tempId, int avatar) {
        super(mapId, status, cx, cy, tempId, avatar);
    }
@Override
public void openBaseMenu(Player p) {
    createOtherMenu(p, ConstNpc.BASE_MENU,
            "Ê,Ê,Ê \nTui Có Bán Cái Này Hay Lắm\n"
          + "Đây Là 3 Gói Quà, Mỗi Gói Có Thể Nhận 7 Lần Sau Khi Mua\n"
          + "Vàng, Bạc Đá Quý Tao Không Thiếu\nTao Có Cả Trứng Mabư Và Chibi\nChibi Của Tao Có Thể Giúp Mày Skip 1 Nhiệm Vụ Bất Kỳ\nThấy Tao Lợi Hại Chưa\n",
            "Mua Quà",
            "Nhận Quà",
            "Đóng");
}

  @Override
public void confirmMenu(Player p, int select) {
    if (p.iDMark.isBaseMenu()) {
        if (select == 0) {
            createOtherMenu(p, ConstNpc.BASE_MENU,
                    "Chọn Gói Quà Đi Nào Hảo Hán\nMỗi Gói Có Thể Nhận Tối Đa Phần Thưởng 7 Lần\nLần Đầu Tiên Vật Phẩm Sẽ Bị Khóa\nTừ Lần Sau Sẽ Không Bị Khóa\n(Chỉ Nhận Thêm Được Vàng)\n",
                    "50K\n+20 TV\n+0 TV",
                    "100K\n+50 TV\n+ Trứng",
                    "250K\n120 TV\n+1 Chibi",
                    "Đóng");
            p.iDMark.setIndexMenu(MENU_GOI);
        } else if (select == 1) {
            openNhanMenu(p);
        }
        return;
    }


        if (p.iDMark.getIndexMenu() == MENU_GOI) {
            if (select == 0) muaGoi(p, GOI_50K);
            if (select == 1) muaGoi(p, GOI_100K);
            if (select == 2) muaGoi(p, GOI_250K);
            return;
        }

        if (p.iDMark.getIndexMenu() == MENU_NHAN) {
            if (select == 0) nhanQua(p, GOI_50K);
            if (select == 1) nhanQua(p, GOI_100K);
            if (select == 2) nhanQua(p, GOI_250K);
        }
    }

    private void openNhanMenu(Player p) {
        createOtherMenu(p, ConstNpc.BASE_MENU,
                "Nhận Quà Nào Bé Iu\nMỗi Mốc Có Thể Nhận 7 Là Nhé Cưng\nĐầu Tư Vào Đây Là HDPE Luôn",
                "Nhận 50k",
                "Nhận 100k",
                "Nhận 250k",
                "Đóng");
        p.iDMark.setIndexMenu(MENU_NHAN);
    }

    private void nhanQua(Player p, int goi) {
        try (Connection con = DBConnecter.getConnectionServer()) {

            String lastCol = goi == GOI_50K ? "npc_last_50k"
                           : goi == GOI_100K ? "npc_last_100k"
                           : "npc_last_250k";

            String countCol = goi == GOI_50K ? "npc_count_50k"
                            : goi == GOI_100K ? "npc_count_100k"
                            : "npc_count_250k";

            PreparedStatement ps = con.prepareStatement(
                    "SELECT npc_package," + lastCol + "," + countCol + ",npc_item_once FROM account WHERE id=?");
            ps.setInt(1, p.getSession().userId);
            ResultSet rs = ps.executeQuery();
            if (!rs.next()) return;

            int pkg = rs.getInt("npc_package");
            long last = rs.getLong(lastCol);
            int count = rs.getInt(countCol);
            int once = rs.getInt("npc_item_once");

            if ((pkg & goi) == 0) {
                Service.gI().sendThongBao(p, "Chưa mua gói!");
                return;
            }

            if (count >= MAX_RECEIVE) {
                Service.gI().sendThongBao(p, "Đã nhận đủ 7 lượt!");
                return;
            }

            long now = System.currentTimeMillis();
            if (last > 0 && now - last < COOLDOWN) {
                long h = (COOLDOWN - (now - last)) / 3600000;
                Service.gI().sendThongBao(p, "Còn " + h + " giờ");
                return;
            }

            // ===== PHÁT QUÀ =====

            // GÓI 50K
            if (goi == GOI_50K) {
                if (count == 0) {
                    giveThoiVangWithParam(p, 20);
                } else {
                    giveThoiVang(p, 20);
                }
            }

            // GÓI 100K
            if (goi == GOI_100K) {
                if (count == 0) {
                    giveThoiVangWithParam(p, 50);
                } else {
                    giveThoiVang(p, 50);
                }

                if ((once & ONCE_TRUNG) == 0) {
                    giveItem(p, 568, 1);
                    once |= ONCE_TRUNG;
                }
            }

            // GÓI 250K
            if (goi == GOI_250K) {
                if (count == 0) {
                    giveThoiVangWithParam(p, 120);
                } else {
                    giveThoiVang(p, 120);
                }

                if ((once & ONCE_SKIP) == 0) {
                    giveItem(p, 1238, 1);
                    once |= ONCE_SKIP;
                }
            }

            PreparedStatement up = con.prepareStatement(
                    "UPDATE account SET " + lastCol + "=?, " + countCol + "=" + countCol + "+1, npc_item_once=? WHERE id=?");
            up.setLong(1, now);
            up.setInt(2, once);
            up.setInt(3, p.getSession().userId);
            up.executeUpdate();

            Service.gI().sendThongBao(p, "Nhận thành công (" + (count + 1) + "/7)");

        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void muaGoi(Player p, int goi) {
        try (Connection con = DBConnecter.getConnectionServer()) {
            PreparedStatement ps = con.prepareStatement(
                    "SELECT cash,npc_package FROM account WHERE id=?");
            ps.setInt(1, p.getSession().userId);
            ResultSet rs = ps.executeQuery();
            if (!rs.next()) return;

            int cash = rs.getInt("cash");
            int pkg = rs.getInt("npc_package");

            int price = goi == GOI_50K ? CASH_50K :
                        goi == GOI_100K ? CASH_100K : CASH_250K;

            if (cash < price) {
                Service.gI().sendThongBao(p, "Không đủ tiền rồi baby!");
                return;
            }

            if ((pkg & goi) != 0) {
                Service.gI().sendThongBao(p, "Chỉ mua được 1 lần!");
                return;
            }

            pkg |= goi;
            cash -= price;

            PreparedStatement up = con.prepareStatement(
                    "UPDATE account SET cash=?, npc_package=? WHERE id=?");
            up.setInt(1, cash);
            up.setInt(2, pkg);
            up.setInt(3, p.getSession().userId);
            up.executeUpdate();

            Service.gI().sendThongBaoOK(p, "BẠN ĐÃ MUA THÀNH CONG!");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void giveThoiVang(Player p, int q) {
        Item it = ItemService.gI().createNewItem((short) 457, q);
        InventoryService.gI().addItemBag(p, it);
        InventoryService.gI().sendItemBag(p);
    }

    private void giveThoiVangWithParam(Player p, int q) {
        Item it = ItemService.gI().createNewItem((short) 457, q);
        it.itemOptions.add(new Item.ItemOption(30, 1));
        InventoryService.gI().addItemBag(p, it);
        InventoryService.gI().sendItemBag(p);
    }

    private void giveItem(Player p, int id, int q) {
        Item it = ItemService.gI().createNewItem((short) id, q);
        InventoryService.gI().addItemBag(p, it);
        InventoryService.gI().sendItemBag(p);
    }
}
