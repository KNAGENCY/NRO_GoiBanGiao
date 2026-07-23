package models.Combine;

import consts.ConstNpc;
import item.Item;

import java.io.IOException;

import models.Combine.manifest.CheTaoTrangBiThienSu;
import models.Combine.manifest.CuongHoaLoSaoPhaLe;
import models.Combine.manifest.DanhBongSaoPhaLe;
import models.Combine.manifest.DapDoAoHoa;
import models.Combine.manifest.EpSaoTrangBi;
import models.Combine.manifest.GiaHanVatPham;
import models.Combine.manifest.GiamDinhSach;
import models.Combine.manifest.HoiPhucSach;
import models.Combine.manifest.MoKhoaItem;
import models.Combine.manifest.NangCapBongTai;
import models.Combine.manifest.NangCapBongTai3;
import models.Combine.manifest.NangCapBongTai4;
import models.Combine.manifest.NangCapChanMenh;
import models.Combine.manifest.NangCapDeMabu;
import models.Combine.manifest.NangCapKichHoat;
import models.Combine.manifest.NangCapKichHoatVip;
import models.Combine.manifest.NangCapSachTuyetKy;
import models.Combine.manifest.NangCapSaoPhaLe;
import models.Combine.manifest.NangCapVatPham;
import models.Combine.manifest.NangChiSoBongTai;
import models.Combine.manifest.NangChiSoBongTai3;
import models.Combine.manifest.NangGiapLuyenTap;
import models.Combine.manifest.NhapNgocRong;
import models.Combine.manifest.PhaLeHoaTrangBi;
import models.Combine.manifest.PhanRaSach;
import models.Combine.manifest.PhapSuHoa;
import models.Combine.manifest.RemoveOptionItem;
import models.Combine.manifest.SieuHoaCaiTrang;
import models.Combine.manifest.TaoDaHematite;
import models.Combine.manifest.TaySach;
import models.Combine.manifest.TinhAnTrangBi;
import models.Combine.manifest.TinhThachHoa;
import player.Player;
import network.Message;
import npc.Npc;
import npc.NpcManager;
import services.InventoryService;

public class CombineService {

    private static final int COST = 500000000;
    private static final int TIME_COMBINE = 1500;
    public static final byte MAX_STAR_ITEM = 8;
    public static final byte MAX_LEVEL_ITEM = 8;
    private static final byte OPEN_TAB_COMBINE = 0;
    private static final byte REOPEN_TAB_COMBINE = 1;
    private static final byte combineSUCCESS = 2;
    private static final byte combineFAIL = 3;
    private static final byte combineCHANGE_OPTION = 4;
    private static final byte combineDRAGON_BALL = 5;
    public static final byte OPEN_ITEM = 6;
    public static final int EP_SAO_TRANG_BI = 500;
    public static final int PHA_LE_HOA_TRANG_BI = 501;
    public static final int CHUYEN_HOA_TRANG_BI_DUNG_VANG = 502;
    public static final int CHUYEN_HOA_TRANG_BI_DUNG_NGOC = 503;
    public static final int NHAP_DA = 504;
    public static final int NANG_CAP_SAO_PHA_LE = 100;
    public static final int DANH_BONG_SAO_PHA_LE = 101;
    public static final int CUONG_HOA_LO_SAO_PHA_LE = 102;
    public static final int TAO_DA_HEMATITE = 103;
    public static final int GIAM_DINH_SACH = 104;
    public static final int TAY_SACH = 105;
    public static final int NANG_CAP_SACH_TUYET_KY = 106;
    public static final int HOI_PHUC_SACH = 107;
    public static final int PHAN_RA_SACH = 108;
    public static final int CHE_TAO_TRANG_BI_THIEN_SU = 109;
    public static final int NANG_CAP_VAT_PHAM = 510;
    public static final int NANG_CAP_BONG_TAI = 511;
    public static final int NANG_CAP_BONG_TAI_3 = 531;
    public static final int NANG_CAP_BONG_TAI_4 = 533;
    public static final int LAM_PHEP_NHAP_DA = 512;
    public static final int NHAP_NGOC_RONG = 513;
    public static final int NANG_CHI_SO_BONG_TAI = 517;
    public static final int NANG_CHI_SO_BONG_TAI_3 = 532;
    public static final int NANG_CAP_KICH_HOAT = 518;
    public static final int NANG_CAP_KICH_HOAT_VIP = 519;
    public static final int NANG_CAP_DO_TL = 530;
    public static final int NANG_CAP_DE_MABU = 4921;

    public static final int DAP_DO_AO_HOA = 520;
    public static final int PS_HOA_TRANG_BI = 521;
    public static final int TAY_PS_HOA_TRANG_BI = 522;
    public static final int MO_KHOA_ITEM = 523;
    public static final int NANG_CAP_CHAN_MENH = 524;
    public static final int AN_TRANG_BI = 525;
    public static final int GIA_HAN_VAT_PHAM = 526;
    public static final int SIEU_HOA = 527;
    public static final int TINH_THACH_HOA = 528;
    public static final int NANG_GIAP_LUYEN_TAP = 529;

    private static CombineService instance;

    public final Npc baHatMit;
    public final Npc whis;

    private CombineService() {
        this.baHatMit = NpcManager.getNpc(ConstNpc.BA_HAT_MIT);
        this.whis = NpcManager.getNpc(ConstNpc.WHIS);
    }

    public static CombineService gI() {
        if (instance == null) {
            instance = new CombineService();
        }
        return instance;
    }

    /**
     * Hiển thị thông tin đập đồ
     *
     * @param player
     * @param index
     */
    public void showInfoCombine(Player player, int[] index) {
        if (player.combine == null) {
            return;
        }
        player.combine.clearItemCombine();
        if (index.length > 0) {
            for (int i = 0; i < index.length; i++) {
                player.combine.itemsCombine.add(player.inventory.itemsBag.get(index[i]));
            }
        }
        switch (player.combine.typeCombine) {
            case EP_SAO_TRANG_BI ->
                EpSaoTrangBi.showInfoCombine(player);
            case PHA_LE_HOA_TRANG_BI ->
                PhaLeHoaTrangBi.showInfoCombine(player);
            case NHAP_NGOC_RONG ->
                NhapNgocRong.showInfoCombine(player);
            case NANG_CAP_VAT_PHAM ->
                NangCapVatPham.showInfoCombine(player);
            case NANG_CAP_BONG_TAI ->
                NangCapBongTai.showInfoCombine(player);
            case NANG_CAP_BONG_TAI_3 ->
                NangCapBongTai3.showInfoCombine(player);
            case NANG_CAP_BONG_TAI_4 ->
                NangCapBongTai4.showInfoCombine(player);
            case NANG_CHI_SO_BONG_TAI ->
                NangChiSoBongTai.showInfoCombine(player);
            case NANG_CHI_SO_BONG_TAI_3 ->
                NangChiSoBongTai3.showInfoCombine(player);
            case NANG_CAP_SAO_PHA_LE ->
                NangCapSaoPhaLe.showInfoCombine(player);
            case DANH_BONG_SAO_PHA_LE ->
                DanhBongSaoPhaLe.showInfoCombine(player);
            case CUONG_HOA_LO_SAO_PHA_LE ->
                CuongHoaLoSaoPhaLe.showInfoCombine(player);
            case TAO_DA_HEMATITE ->
                TaoDaHematite.showInfoCombine(player);
            case GIAM_DINH_SACH ->
                GiamDinhSach.showInfoCombine(player);
            case TAY_SACH ->
                TaySach.showInfoCombine(player);
            case NANG_CAP_SACH_TUYET_KY ->
                NangCapSachTuyetKy.showInfoCombine(player);
            case HOI_PHUC_SACH ->
                HoiPhucSach.showInfoCombine(player);
            case PHAN_RA_SACH ->
                PhanRaSach.showInfoCombine(player);
            case CHE_TAO_TRANG_BI_THIEN_SU ->
                CheTaoTrangBiThienSu.showInfoCombine(player);
            case NANG_CAP_KICH_HOAT ->
                NangCapKichHoat.showInfoCombine(player);
            case NANG_CAP_KICH_HOAT_VIP ->
                NangCapKichHoatVip.showInfoCombine(player);
            case NANG_CAP_DE_MABU ->
                NangCapDeMabu.showInfoCombine(player);
            case DAP_DO_AO_HOA ->
                DapDoAoHoa.showInfoCombine(player);
            case PS_HOA_TRANG_BI ->
                PhapSuHoa.showInfoCombine(player);
            case TAY_PS_HOA_TRANG_BI ->
                RemoveOptionItem.showInfoCombine(player);
            case MO_KHOA_ITEM ->
                MoKhoaItem.showInfoCombine(player);
            case NANG_CAP_CHAN_MENH ->
                NangCapChanMenh.showInfoCombine(player);
            case AN_TRANG_BI ->
                TinhAnTrangBi.showInfoCombine(player);
            case GIA_HAN_VAT_PHAM ->
                GiaHanVatPham.showInfoCombine(player);
            case SIEU_HOA ->
                SieuHoaCaiTrang.showInfoCombine(player);
            case TINH_THACH_HOA ->
                TinhThachHoa.showInfoCombine(player);
            case NANG_GIAP_LUYEN_TAP ->
                NangGiapLuyenTap.showInfoCombine(player);

        }
    }

    /**
     * Bắt đầu đập đồ - điều hướng từng loại đập đồ
     *
     * @param player
     * @param n
     */
    public void startCombine(Player player, int... n) {
        int num = 0;
        if (n.length > 0) {
            num = n[0];
        }
        switch (player.combine.typeCombine) {
            case EP_SAO_TRANG_BI ->
                EpSaoTrangBi.epSaoTrangBi(player);
            case PHA_LE_HOA_TRANG_BI ->
                PhaLeHoaTrangBi.phaLeHoa(player, num);
            case NHAP_NGOC_RONG ->
                NhapNgocRong.nhapNgocRong(player, num == 1);
            case NANG_CAP_VAT_PHAM ->
                NangCapVatPham.nangCapVatPham(player, num == 1);
            case NANG_CAP_BONG_TAI ->
                NangCapBongTai.nangCapBongTai(player);
            case NANG_CAP_BONG_TAI_3 ->
                NangCapBongTai3.nangCapBongTai3(player);
            case NANG_CAP_BONG_TAI_4 ->
                NangCapBongTai4.nangCapBongTai4(player);
            case NANG_CHI_SO_BONG_TAI ->
                NangChiSoBongTai.nangChiSoBongTai(player);
            case NANG_CHI_SO_BONG_TAI_3 ->
                NangChiSoBongTai3.nangChiSoBongTai3(player);
            case NANG_CAP_SAO_PHA_LE ->
                NangCapSaoPhaLe.nangCapSaoPhaLe(player);
            case DANH_BONG_SAO_PHA_LE ->
                DanhBongSaoPhaLe.danhBongSaoPhaLe(player);
            case CUONG_HOA_LO_SAO_PHA_LE ->
                CuongHoaLoSaoPhaLe.cuongHoaLoSaoPhaLe(player);
            case TAO_DA_HEMATITE ->
                TaoDaHematite.taoDaHematite(player);
            case GIAM_DINH_SACH ->
                GiamDinhSach.giamDinhSach(player);
            case TAY_SACH ->
                TaySach.taySach(player);
            case NANG_CAP_SACH_TUYET_KY ->
                NangCapSachTuyetKy.nangCapSachTuyetKy(player);
            case HOI_PHUC_SACH ->
                HoiPhucSach.hoiPhucSach(player);
            case PHAN_RA_SACH ->
                PhanRaSach.phanRaSach(player);
            case CHE_TAO_TRANG_BI_THIEN_SU ->
                CheTaoTrangBiThienSu.cheTaoTrangBiThienSu(player);
            case NANG_CAP_KICH_HOAT ->
                NangCapKichHoat.startCombine(player);
            case NANG_CAP_KICH_HOAT_VIP ->
                NangCapKichHoatVip.startCombine(player);
            case NANG_CAP_DE_MABU ->
                NangCapDeMabu.startCombine(player);
            case DAP_DO_AO_HOA ->
                DapDoAoHoa.startCombine(player);
            case PS_HOA_TRANG_BI ->
                PhapSuHoa.startCombine(player);
            case TAY_PS_HOA_TRANG_BI ->
                RemoveOptionItem.startCombine(player);
            case MO_KHOA_ITEM ->
                MoKhoaItem.startCombine(player);
            case NANG_CAP_CHAN_MENH ->
                NangCapChanMenh.startCombine(player);
            case AN_TRANG_BI ->
                TinhAnTrangBi.startCombine(player);
            case GIA_HAN_VAT_PHAM ->
                GiaHanVatPham.startCombine(player);
            case SIEU_HOA ->
                SieuHoaCaiTrang.startCombine(player);
            case TINH_THACH_HOA ->
                TinhThachHoa.startCombine(player);
            case NANG_GIAP_LUYEN_TAP ->
                NangGiapLuyenTap.startCombine(player);
        }

        player.iDMark.setIndexMenu(ConstNpc.IGNORE_MENU);
        player.combine.clearParamCombine();
        player.combine.lastTimeCombine = System.currentTimeMillis();

    }

    /**
     * Mở tab đập đồ
     *
     * @param player
     * @param type   kiểu đập đồ
     */
    public void openTabCombine(Player player, int type) {
        player.combine.setTypeCombine(type);
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(OPEN_TAB_COMBINE);
            msg.writer().writeUTF(getTextInfoTabCombine(type));
            msg.writer().writeUTF(getTextTopTabCombine(type));
            if (player.iDMark.getNpcChose() != null) {
                msg.writer().writeShort(player.iDMark.getNpcChose().tempId);
            }
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    /**
     * Hiệu ứng mở item
     *
     * @param player
     * @param icon1
     * @param icon2
     */
    public void sendEffectOpenItem(Player player, short icon1, short icon2) {
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(OPEN_ITEM);
            msg.writer().writeShort(icon1);
            msg.writer().writeShort(icon2);
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    public void sendEffectCombineItem(Player player, byte type, short icon1, short icon2) {
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(type);
            switch (type) {
                case 0:
                    msg.writer().writeUTF("");
                    msg.writer().writeUTF("");
                    break;
                case 1:
                    msg.writer().writeByte(0);
                    msg.writer().writeByte(-1);
                    break;
                case 2: // success 0 eff 0
                case 3: // success 1 eff 0
                    break;
                case 4: // success 0 eff 1
                    msg.writer().writeShort(icon1);
                    break;
                case 5: // success 0 eff 2
                    msg.writer().writeShort(icon1);
                    break;
                case 6: // success 0 eff 3
                    msg.writer().writeShort(icon1);
                    msg.writer().writeShort(icon2);
                    break;
                case 7: // success 0 eff 4
                    msg.writer().writeShort(icon1);
                    break;
                case 8: // success 1 eff 4
                    break;
            }
            msg.writer().writeShort(-1); // id npc
            // msg.writer().writeShort(-1); // x
            // msg.writer().writeShort(-1); // y
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    /**
     * Hiệu ứng đập đồ thành công
     *
     * @param player
     */
    public void sendEffectSuccessCombine(Player player) {
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(combineSUCCESS);
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    /**
     * Hiệu ứng đập đồ thất bại
     *
     * @param player
     */
    public void sendEffectFailCombine(Player player) {
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(combineFAIL);
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    /**
     * Gửi lại danh sách đồ trong tab combine
     *
     * @param player
     */
    public void reOpenItemCombine(Player player) {
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(REOPEN_TAB_COMBINE);
            msg.writer().writeByte(player.combine.itemsCombine.size());
            for (Item it : player.combine.itemsCombine) {
                for (int j = 0; j < player.inventory.itemsBag.size(); j++) {
                    if (it == player.inventory.itemsBag.get(j)) {
                        msg.writer().writeByte(j);
                    }
                }
            }
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    /**
     * Hiệu ứng ghép ngọc rồng
     *
     * @param player
     * @param icon
     */
    public void sendEffectCombineDB(Player player, short icon) {
        Message msg = null;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(combineDRAGON_BALL);
            msg.writer().writeShort(icon);
            player.sendMessage(msg);
        } catch (Exception e) {
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    public void sendAddItemCombine(Player player, int npcId, Item... items) {
        Message msg;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(0);
            msg.writer().writeUTF("By Kame");
            msg.writer().writeUTF("Ngọc Rồng Kame");
            msg.writer().writeShort(npcId);
            player.sendMessage(msg);
            msg.cleanup();
            msg = new Message(-81);
            msg.writer().writeByte(1);
            msg.writer().writeByte(items.length);
            for (Item item : items) {
                msg.writer().writeByte(InventoryService.gI().getIndexItemBag(player, item));
            }
            player.sendMessage(msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendEffSuccessVip(Player player, int iconID) {
        Message msg;
        try {
            msg = new Message(-81);
            msg.writer().writeByte(7);
            msg.writer().writeShort(iconID);
            player.sendMessage(msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    public void sendEffFailVip(Player player) {
        try {
            Message msg;
            msg = new Message(-81);
            msg.writer().writeByte(8);
            player.sendMessage(msg);
            msg.cleanup();
        } catch (IOException e) {
        }
    }

    private String getTextTopTabCombine(int type) {
        return switch (type) {
            case NANG_CAP_BONG_TAI_3 ->
                "Ta sẽ phù phép\ncho Bông tai cấp 2\nthành cấp 3";
            case NANG_CAP_BONG_TAI_4 ->
              "Ta sẽ phù phép\ncho Bông tai cấp 3\nthành cấp 4";
            case NANG_CHI_SO_BONG_TAI_3 ->
                "Ta sẽ phù phép\ncho Bông tai cấp 3\ncó 2 chỉ số ngẫu nhiên";
            case EP_SAO_TRANG_BI ->
                "Ta sẽ phù phép\ncho trang bị của ngươi\ntrở nên mạnh mẽ";
            case PHA_LE_HOA_TRANG_BI ->
                "Ta sẽ phù phép\ncho trang bị của ngươi\ntrở thành trang bị pha lê";
            case CHUYEN_HOA_TRANG_BI_DUNG_VANG, CHUYEN_HOA_TRANG_BI_DUNG_NGOC ->
                "Lưu ý trang bị mới\nphải hơn trang bị gốc\n1 bậc";
            case NHAP_NGOC_RONG ->
                "Ta sẽ phù phép\ncho 7 viên Ngọc Rồng\nthành 1 viên Ngọc Rồng cấp cao";
            case NHAP_DA ->
                "Ta sẽ phù phép\ncho 10 mảnh đá vụn\ntrở thành 1 đá nâng cấp";
            case NANG_CAP_VAT_PHAM ->
                "Ta sẽ phù phép\ncho trang bị của ngươi\ntrở nên mạnh mẽ";
            case NANG_CAP_BONG_TAI ->
                "Ta sẽ phù phép\ncho Bông tai\nthành cấp 2";
            case NANG_CHI_SO_BONG_TAI ->
                "Ta sẽ phù phép\ncho Bông tai cấp 2\ncó 1 chỉ số ngẫu nhiên";
            case NANG_CAP_SAO_PHA_LE ->
                "Ta sẽ phù phép\nnâng cấp Sao Pha Lê\nthành cấp 2";
            case DANH_BONG_SAO_PHA_LE ->
                "Đánh bóng\nSao pha lê cấp 2";
            case CUONG_HOA_LO_SAO_PHA_LE ->
                "Cường hóa\nÔ Sao Pha Lê";
            case TAO_DA_HEMATITE ->
                "Ta sẽ phù phép\ntạo đá hematite";
            case GIAM_DINH_SACH ->
                "Ta sẽ phù phép\ngiám định sách đó cho ngươi";
            case TAY_SACH ->
                "Ta sẽ phù phép\ntẩy sách đó cho ngươi";
            case NANG_CAP_SACH_TUYET_KY ->
                "Ta sẽ phù phép\nnâng cấp Sách Tuyệt Kỹ cho ngươi";
            case HOI_PHUC_SACH ->
                "Ta sẽ phù phép\nphục hồi sách cho ngươi";
            case PHAN_RA_SACH ->
                "Ta sẽ phù phép\nphân rã sách đó cho ngươi";
            case CHE_TAO_TRANG_BI_THIEN_SU ->
                "Chế tạo\ntrang bị thiên sứ";
            case LAM_PHEP_NHAP_DA ->
                "Ta sẽ phù phép\n"
                        + "cho 10 mảnh đá vụn\n"
                        + "trở thành 1 đá nâng cấp";
            case NANG_CAP_KICH_HOAT ->
                "Ta sẽ phù phép\nchế tạo Đồ Thần Kinh\nthành đồ bị Kích Thích";
            case NANG_CAP_KICH_HOAT_VIP ->
                "Ta sẽ phù phép\nchế tạo trang bị Thiên Sứ\nthành trang bị Kích Hoạt Vip";
            case NANG_CAP_DO_TL ->
                "Ta sẽ phù phép\nchế tạo trang bị Thần Linh\nthành trang bị Hủy Diệt";
            case GIA_HAN_VAT_PHAM ->
                "Ta sẽ phù phép\ncho trang bị của ngươi\nthêm hạn sử dụng";
            case SIEU_HOA ->
                "Ta sẽ giúp con siêu hóa\n Cải trang";
            case TINH_THACH_HOA ->
                "Ta sẽ giúp con Tinh Thạch đồ";
            case DAP_DO_AO_HOA ->
                "Ta sẽ giúp ngươi ảo hóa đồ để có thuộc tính cao hơn";
            case NANG_CAP_CHAN_MENH ->
                "Ta sẽ giúp ngươi\nnâng cấp Chân Mệnh\ncủa ngươi lên 1 cấp bậc";
            case PS_HOA_TRANG_BI ->
                "Pháp sư hóa Pet,\nLinh thú và Ván bay";
            case TAY_PS_HOA_TRANG_BI ->
                "Ta sẽ trả lại\nnguyên bản phụ kiện\ncho ngươi";
            case MO_KHOA_ITEM ->
                "Mở Khóa giao dịch Item";
            case NANG_CAP_DE_MABU ->
                "Tăng Chỉ Số\nKhi Hợp Thể\nNhé Ku";
            case AN_TRANG_BI ->
                "Ta sẽ phù phép\ncho trang bị của ngươi\ntrở thành Đồ Vip ";
            default ->
                "";
        };
    }

    private String getTextInfoTabCombine(int type) {
        return switch (type) {
            case NANG_CAP_BONG_TAI_3 ->
                "Vào hành trang\nChọn Bông tai cấp 2\n+ 1 Thạch Thần\n+ Mảnh vỡ bông tai 9999\nSau đó chọn 'Nâng cấp'";
            case NANG_CAP_BONG_TAI_4 ->
                "V\u00e0o h\u00e0nh trang\nCh\u1ecdn x1 B\u00f4ng Tai C\u1ea5p 3\n+ x1 B\u00f4ng Tai C\u1ea5p 2\nSau \u0111\u00f3 ch\u1ecdn '\u0110\u1ed5i'";
            case NANG_CHI_SO_BONG_TAI_3 ->
                "Chọn Bông tai cấp 3\n+ Mảnh hồn bông tai x99\n+ 1 Đá xanh lam.\nSau đó chọn 'Nâng cấp'";
            case EP_SAO_TRANG_BI ->
                "Chọn: Áo, Quần, Găng, Giày, Rađa, Giáp\n+ 1 Sao pha lê hoặc Ngọc Rồng\n--------\n8 sao cần thêm Thạch Thần";
            case PHA_LE_HOA_TRANG_BI ->
                "Chọn: Áo, Quần, Găng, Giày, Rađa, Giáp\n+ 5 Ảnh gái xin để lọ\n--------\nTối đa Max 8 sao\nTỉ lệ thấp dần";
            case CHUYEN_HOA_TRANG_BI_DUNG_VANG, CHUYEN_HOA_TRANG_BI_DUNG_NGOC ->
                "Vào hành trang\nChọn trang bị gốc\n(Áo,quần,găng,giày hoặc rađa)\ntừ cấp [+4] trở lên\nChọn tiếp trang bị mới\nchưa nâng cấp cần nhập thể\nsau đó chọn 'Nâng cấp'";
            case NHAP_NGOC_RONG ->
                "Vào hành trang\nChọn 7 viên ngọc cùng sao\nSau đó chọn 'Làm phép'";
            case NHAP_DA ->
                "Vào hành trang\nChọn 10 mảnh đá vụn\nChọn 1 bình nước phép\n(mua tại Uron ở trạm tàu vũ trụ)\nSau đó chọn 'Làm phép'";
            case NANG_CAP_VAT_PHAM ->
                "Vào hành trang\nChọn trang bị\n(Áo, Quần, Găng, Giày, Rađa)\nChọn loại đá để nâng cấp\nSau đó chọn 'Nâng cấp'";
            case NANG_CAP_BONG_TAI ->
                "Vào hành trang\nChọn Bông tai cấp 1\n+ Mảnh vỡ bông tai 9999\nSau đó chọn 'Nâng cấp'";
            case NANG_CHI_SO_BONG_TAI ->
               "Chọn Bông tai cấp 3\n+ Mảnh hồn bông tai x99\n+ 1 Đá xanh lam.\nSau đó chọn 'Nâng cấp'";
            case NANG_CAP_SAO_PHA_LE ->
                "Vào hành trang\nChọn đá Hematite\nChọn loại sao pha lê (cấp 1)\nSau đó chọn 'Nâng cấp'";
            case DANH_BONG_SAO_PHA_LE ->
                "Vào hành trang\nChọn loại sao pha lê cấp 2 có từ 2 viên trở lên\nChọn 1 đá mài\nSau đó chọn 'Đánh bóng'";
            case CUONG_HOA_LO_SAO_PHA_LE ->
                "Vào hành trang\nChọn trang bị có Ô sao thứ 8 trở lên chưa cường hóa\nChọn đá Hematite\nChọn dùi đục\nSau đó chọn 'Cường hóa'";
            case TAO_DA_HEMATITE ->
                "Chọn 2 đồ cùng loại và #:\n+ 2 đồ Thần Linh\nHoặc: + 2 đồ Hủy Diệt\n--------\nTùy vào chất lượng đồ\nNhận Thạch Thần #";
            case GIAM_DINH_SACH ->
                "1 Sách Có Thể Giám Định\n+1 Bùa Giám Định";
            case TAY_SACH ->
                "Vào hành trang chọn\n1 sách cần tẩy";
            case NANG_CAP_SACH_TUYET_KY ->
                "Vào hành trang chọn\nSách Tuyệt Kỹ 1 cần nâng cấp và 10 Kìm bấm giấy";
            case HOI_PHUC_SACH ->
                "Vào hành trang chọn\nCác Sách Tuyệt Kỹ cần phục hồi";
            case PHAN_RA_SACH ->
                "Vào hành trang chọn\n1 sách cần phân rã";
            case CHE_TAO_TRANG_BI_THIEN_SU ->
                "Cần (4):\n- 1 công thức\n- 999 mảnh Thiên Sứ tương ứng\n- 1 đá nâng cấp\n- 1 đá may mắn";
            case LAM_PHEP_NHAP_DA ->
                "Vào hành trang\n"
                        + "Chọn 10 mảnh đá vụn\n"
                        + "Chọn 1 bình nước phép\n"
                        + "(mua tại Uron ở trạm tàu vũ trụ)\n"
                        + "Sau đó chọn 'Làm phép'";
            case NANG_CAP_KICH_HOAT ->
                "Chọn: 2 đồ Hủy Diệt #\n+ Nhiều Thạch Thần\n--------\nĐồ KH mới = đồ HD đầu";
            case NANG_CAP_KICH_HOAT_VIP ->
                "Chọn: 2 đồ Thần Linh #\n+ Lọ thêm 10 lần\n--------\nĐồ Thần đầu = đồ HD mới";
            case NANG_CAP_DO_TL ->
                "Vào hành trang\nChọn 1 trang bị Thần Linh\nChọn 1 viên đá Kích Hoạt Vip\nSau đó chọn 'Nâng cấp'";
            case DAP_DO_AO_HOA ->
                "vào hành trang\nChọn trang bị\n(Áo, quần, găng, giày hoặc rađa)"
                        + "\nChọn loại đá quý để nâng cấp\n"
                        + "\nCó thể thêm đá bảo vệ để tránh tụt cấp\n"
                        + "Sau đó chọn 'Nâng cấp'";
            case NANG_CAP_CHAN_MENH ->
                "Chọn: Chân Mệnh"
                        + "\n+ 10 đá Ánh Dương\n--------\nTối đa cấp 9"
                        + "\nTăng mạnh chỉ số";
            case PS_HOA_TRANG_BI ->
                "Chọn: Ván Bay hoặc Đeo Lưng\n+1 Thạch Thần\n"
                        + "+ 1 chút cầu nguyện\n--------\n"
                        + "Nhận được 2% HP, KI, SĐ, SĐCM";
            case MO_KHOA_ITEM ->
                "vào hành trang\nChọn 1 trang bị khóa giao dịch ( bông tai, item sự kiện, thỏi vàng,..) và Đá Hoàng Kim \n "
                        + " để mở khóa giao dịch Item"
                        + "Chỉ cần chọn 'Mở Khóa'";

            case TAY_PS_HOA_TRANG_BI ->
                "Chọn: Ván Bay hoặc Đeo Lưng bị lỏ\n +1 Thạch Thần\n--------\n"
                        + "Tẩy xong có thể Pháp Sư Hóa lại\n"
                        + "Nhưng vẫn có thể lỏ tiếp, hẹ hẹ";

            case AN_TRANG_BI ->
                "Chọn: Áo, Quần, Giày, Găng, Rada\n+1 Thạch Thần\n+ Khô máu đến khi ra Ấn phù hợp\n--------\nTinh ấn (5 món +15%SD)\n Nhật ấn (5 món +15%Hp\n Nguyệt ấn (5 món +15%KI)";

            case GIA_HAN_VAT_PHAM ->
                "Vào hành trang\n"
                        + "Chọn 1 trang bị có hạn sử dụng\n"
                        + "Chọn Đá Hoàng Kim\n"
                        + "Sau đó chọn 'Gia hạn'";
            case NANG_CAP_DE_MABU ->
                "Chọn: 2 Quả Trứng Mabư\n+ 1 Thạch Thần\n--------\nTối đa 7 Level\nHợp thể tăng chỉ số";
            case SIEU_HOA ->
                "Vào hành trang\n"
                        + "Chọn 1 Cải trang\n"
                        + "Chọn Đá Siêu Hóa\n"
                        + "Sau đó chọn 'Nâng Cấp'";
            case TINH_THACH_HOA ->
                "Vào hành trang\n"
                        + "Chọn 1 Vật Phẩm (Pet, Linh Thú, VPDL)\n"
                        + "Chọn 1 loại đá Tinh thạch\n"
                        + "Sau đó chọn 'Nâng Cấp'";

            case NANG_GIAP_LUYEN_TAP ->
                "Vào hành trang\n"
                        + "Chọn 1 Giáp luyện tập\n"
                        + "Chọn đá hổ phách\n"
                        + "Sau đó chọn 'Nâng Cấp'";
            default ->
                "";
        };
    }

}
