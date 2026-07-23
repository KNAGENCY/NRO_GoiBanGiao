package npc.npc_manifest;

/**
 * @author GiaCatLuong
 */
import consts.ConstMenu;
import consts.ConstNpc;
import item.Item;

import models.Combine.CombineService;
import models.Combine.manifest.CheTaoCuonSachCu;
import models.Combine.manifest.DoiSachTuyetKy;
import models.Combine.manifest.NangCapVatPham;
import models.DeathOrAliveArena.DeathOrAliveArena;
import models.DeathOrAliveArena.DeathOrAliveArenaManager;
import models.DeathOrAliveArena.DeathOrAliveArenaService;
import npc.Npc;
import player.Player;
import services.InventoryService;
import services.ItemService;
import services.Service;
import services.func.ChangeMapService;
import services.func.TopService;
import shop.ShopService;
import utils.Util;

public class BaHatMit extends Npc {

    public BaHatMit(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            switch (this.mapId) {
                case 5 ->
                    this.createOtherMenu(player, ConstNpc.BASE_MENU, "Ngươi tìm ta có việc gì?", "Chức năng\nPhê la",
                            "Phụ kiện &\nĐập ấn", "Võ đài\nSinh Tử", "Nâng cấp\nĐệ Mabư", "Nâng cấp\nHD - SKH");
               case 112 -> {
    if (Util.isAfterMidnight(player.lastTimePKVoDaiSinhTu)) {
        player.haveRewardVDST = false;
        player.thoiVangVoDaiSinhTu = 0;
    }

    if (player.haveRewardVDST) {
        this.createOtherMenu(player, ConstNpc.BASE_MENU,
                "Đây là phần thưởng cho con.",
                "Bí ngô\nngẫu nhiên",
                "Sức mạnh\nbí kíp");
        return;
    }

    if (DeathOrAliveArenaManager.gI().getVDST(player.zone) != null) {
        if (DeathOrAliveArenaManager.gI().getVDST(player.zone).getPlayer().equals(player)) {
            this.createOtherMenu(player, ConstNpc.BASE_MENU,
                    "Ngươi muốn hủy đăng ký thi đấu võ đài?",
                    "Đồng ý\n" + player.thoiVangVoDaiSinhTu + " thỏi vàng",
                    "Từ chối",
                    "Về\nđảo rùa");
            return;
        }

        this.createOtherMenu(player, ConstNpc.BASE_MENU,
                "Ngươi muốn đăng ký thi đấu võ đài?\nnhiều phần thưởng giá trị đang đợi ngươi đó\nThắng võ đài sẽ nhận được:\n - Vật phẩm nâng cấp Sách tuyệt kỹ\n - Hoặc Ngọc rồng xương Vip",
                "Đồng ý\n" + player.thoiVangVoDaiSinhTu + " thỏi vàng",
                "Từ chối",
                "Về\nđảo rùa");
        return;
    }

    this.createOtherMenu(player, ConstNpc.BASE_MENU,
            "Ngươi muốn đăng ký thi đấu võ đài?\nnhiều phần thưởng giá trị đang đợi ngươi đó\nThắng võ đài sẽ nhận được:\n - Vật phẩm nâng cấp Sách tuyệt kỹ\n - Hoặc Ngọc rồng xương Vip",
            "Đồng ý\n" + player.thoiVangVoDaiSinhTu + " thỏi vàng",
            "Từ chối",
            "Về\nđảo rùa");
}
                case 174 ->
                    this.createOtherMenu(player, ConstNpc.BASE_MENU, "Ngươi tìm ta có việc gì?", "Quay về", "Từ chối");
                case 181 ->
                    this.createOtherMenu(player, ConstNpc.BASE_MENU, "Ngươi tìm ta có việc gì?", "Quay về", "Từ chối");
                default -> {
                    this.createOtherMenu(player, ConstNpc.BASE_MENU, "M\u00e0y T\u00ecm Tao L\u00e0m C\u00e1i Ch\u00f3a G\u00ec?",
                            "S\u00e1ch\nTuy\u1ec7t K\u1ef9",
                            "C\u1eeda H\u00e0ng\nB\u00f9a",
                            "N\u00e2ng C\u1ea5p\nV\u1eadt Ph\u1ea9m",
                            "N\u00e2ng c\u1ea5p\nb\u00f4ng tai 2",
                            "M\u1edf ch\u1ec9 s\u1ed1\nb\u00f4ng tai 2",
                            "Nh\u1eadp\nNg\u1ecdc R\u1ed3ng");
                }
            }
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            switch (this.mapId) {
                case 5 -> {
                    switch (player.iDMark.getIndexMenu()) {
                        case ConstNpc.BASE_MENU -> {
                            switch (select) {
                                case 0 ->
                                    createOtherMenu(player, ConstMenu.MENU_PHA_LE, "Mày nhìn cái đéo gì",
                                            "Ép Sao\nTrang Bị", "Phê La\nHoá\nTrang Bị");
                                // , "Nâng cấp\nSao pha lê",
                                // "Đánh bóng\nSao pha lê", "Cường hoá\nLỗ sao\npha lê", "Tạo đá Hematite");

                                case 1 ->
                                    this.createOtherMenu(player, ConstMenu.MENU_CHUYEN_HOA_SKH,
                                            "Ta sẽ nâng trang bị của ngươi\nlên một tầm cao mới hoàn toàn khác",
                                            "Pháp Sư\nHóa", "Tẩy\nPháp Sư", "Đập Ấn\nTrang Bị",  "Nâng Cấp\nChân Mệnh");
                                // ,"Đập đồ\nẢo hóa",
                                // "Đập ấn\nTrang bị");
                                case 2 -> {
                                    if (player.nPoint.power >= 9_000_000_000L) { // check >= 40 tỷ
                                        ChangeMapService.gI().changeMapNonSpaceship(player, 112,
                                                200 + Util.nextInt(-100, 100), 408);
                                    } else {
                                        Service.gI().sendThongBaoOK(player,
                                                "Trong Này Có Đồ Ngon, Nhưng Cần Trên 10 Tỷ Sức Mạnh!");
                                    }
                                }

                                // case 3 ->
                                // CombineService.gI().openTabCombine(player,
                                // CombineService.NANG_CAP_CHAN_MENH);
                                case 3 ->
                                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_DE_MABU);
                                                              case 4 ->
            this.createOtherMenu(player, ConstMenu.MENU_NANG_CAP_TRANG_BI,
                    "Ta có thể giúp ngươi rình gái tắm\nTuyệt đối không được quay phim lại",
                    "Chế Tạo\nSKH", "Chế Tạo\nHủy Diệt", 
                    "Phân rã\nTrang bị");
    
                            
                                    }
                                }

                    
                                // case 4 ->
                                // this.createOtherMenu(player, ConstMenu.CHUC_NANG_BHM_KHAC,
                                // "Ngươi tìm ta có việc gì?\n",
                                // "Đập đồ\nẢo Hóa",
                                // "Build Đồ",
                                // "Pháp sư hoá",
                                // "Siêu hóa\n Cải trang",
                                // "Tinh ấn\ntrang bị",
                                // "Tinh thạch\ntrang bị",
                                // "Nâng cấp\nGiáp LT");
                            
                        
                        case ConstMenu.MENU_CHUYEN_HOA_SKH -> {
                            switch (select) {
                                case 0: // PHÁP SƯ HÓA
                                    CombineService.gI().openTabCombine(player, CombineService.PS_HOA_TRANG_BI);
                                    break;
                                case 1: // TẨY
                                    CombineService.gI().openTabCombine(player, CombineService.TAY_PS_HOA_TRANG_BI);
                                    break;
                                case 2: // ẤN
                                    CombineService.gI().openTabCombine(player, CombineService.AN_TRANG_BI);
                                    break;
                                      case 3:
                                CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_CHAN_MENH);
                                break;
                           

                                // case 2: // Đập đồ ảo hóa
                                // CombineService.gI().openTabCombine(player, CombineService.DAP_DO_AO_HOA);
                                // break;
                                // case 3:
                                // CombineService.gI().openTabCombine(player, CombineService.AN_TRANG_BI);
                                // break;// Name: GiaCatLuong
                            }
                        }
                        case ConstMenu.CHUC_NANG_BHM_KHAC -> {
                            switch (select) {
                                case 0: // Đập đồ ảo hóa
                                    CombineService.gI().openTabCombine(player, CombineService.DAP_DO_AO_HOA);
                                    break;
                                case 1: // Build Đồ
                                    this.createOtherMenu(player, ConstMenu.BUILD_DO_BHM, "|7|\bChi Tiết:\n|0|"
                                            + "\nMở Khóa GD: tỉ lệ thành công 30%"
                                            + "\nGia hạn Vật Phẩm: tỉ lệ thành công 30% + 3 - 7 ngày, 70% + 1 ngày"
                                            + "\nTẩy đồ: tẩy sao pha lê, chỉ số đặc biệt một số trang bị",
                                            "Mở Khóa GD", "Gia hạn\n Vật Phẩm", "Tẩy Đồ");
                                    break;
                                case 2: // Pháp sư hóa trang bị
                                    CombineService.gI().openTabCombine(player, CombineService.PS_HOA_TRANG_BI);
                                    break;
                                case 3: // Siêu Hóa
                                    CombineService.gI().openTabCombine(player, CombineService.SIEU_HOA);
                                    break;
                                case 4:
                                    CombineService.gI().openTabCombine(player, CombineService.AN_TRANG_BI);
                                    break;// Name: GiaCatLuong
                                case 5:
                                    CombineService.gI().openTabCombine(player, CombineService.TINH_THACH_HOA);
                                    break;// Name: GiaCatLuong
                                case 6:
                                    CombineService.gI().openTabCombine(player, CombineService.NANG_GIAP_LUYEN_TAP);
                                    break;// Name: GiaCatLuong
                            }
                        }
                        case ConstMenu.BUILD_DO_BHM -> {
                            switch (select) {
                                case 0: // Mở khóa Item
                                    CombineService.gI().openTabCombine(player, CombineService.MO_KHOA_ITEM);
                                    break;
                                case 1: // Gia hạn vật phẩm
                                    CombineService.gI().openTabCombine(player, CombineService.GIA_HAN_VAT_PHAM);
                                    break;
                                case 2: // Tẩy đồ
                                    CombineService.gI().openTabCombine(player, CombineService.TAY_PS_HOA_TRANG_BI);
                                    break;
                            }
                        }

                        case ConstMenu.MENU_PHA_LE -> {
                            switch (select) {
                                case 0: // Ép sao trang bị
                                    CombineService.gI().openTabCombine(player, CombineService.EP_SAO_TRANG_BI);
                                    break;
                                case 1: // Pha lê hoá trang bị
                                    createOtherMenu(player, ConstMenu.MENU_PHA_LE_HOA_TRANG_BI,
                                            "Ngươi muốn pha lê hoá trang bị bằng cách nào?", "Bằng ngọc", "Từ chối");
                                    break;
                                // case 2: // Nâng cấp sao pha lê
                                // CombineService.gI().openTabCombine(player,
                                // CombineService.NANG_CAP_SAO_PHA_LE);
                                // break;
                                // case 3: // Đánh bóng sao pha lê
                                // CombineService.gI().openTabCombine(player,
                                // CombineService.DANH_BONG_SAO_PHA_LE);
                                // break;
                                // case 4: // Cường hoá lỗ sao pha lê
                                // CombineService.gI().openTabCombine(player,
                                // CombineService.CUONG_HOA_LO_SAO_PHA_LE);
                                // break;
                                // case 5: // Tạo đá Hematite
                                // CombineService.gI().openTabCombine(player, CombineService.TAO_DA_HEMATITE);
                                // break;
                            }
                        }
                        case ConstMenu.MENU_NANG_CAP_TRANG_BI -> {
                            switch (select) {
                                case 0:
                                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_KICH_HOAT);
                                    break;
                                case 1:
                                    CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_KICH_HOAT_VIP);
                                    break;
                                      case 2: // Tạo đá Hematite
                                CombineService.gI().openTabCombine(player, CombineService.TAO_DA_HEMATITE);
                                break;
                            }
                        }
                        case ConstMenu.MENU_CHUYEN_HOA_TRANG_BI -> {
                            switch (select) {
                                case 0 -> {
                                    CombineService.gI().openTabCombine(player,
                                            CombineService.CHUYEN_HOA_TRANG_BI_DUNG_VANG);
                                }
                                case 1 -> {
                                    CombineService.gI().openTabCombine(player,
                                            CombineService.CHUYEN_HOA_TRANG_BI_DUNG_NGOC);
                                }
                            }
                        }
                        case ConstMenu.MENU_PHA_LE_HOA_TRANG_BI -> {
                            if (select == 0) {
                                CombineService.gI().openTabCombine(player, CombineService.PHA_LE_HOA_TRANG_BI);
                            }
                        }
                        case ConstNpc.MENU_START_COMBINE -> {
                            switch (player.combine.typeCombine) {
                                case CombineService.PHA_LE_HOA_TRANG_BI -> {
                                    switch (select) {
                                        case 0 ->
                                            CombineService.gI().startCombine(player, 100);
                                        case 1 ->
                                            CombineService.gI().startCombine(player, 10);
                                        case 2 ->
                                            CombineService.gI().startCombine(player);
                                    }
                                }

                                case CombineService.NANG_CAP_DO_TL,
                                        CombineService.NANG_CAP_KICH_HOAT_VIP,
                                        CombineService.NANG_CAP_KICH_HOAT,
                                        CombineService.NANG_CAP_SAO_PHA_LE,
                                        CombineService.DANH_BONG_SAO_PHA_LE,
                                        CombineService.CUONG_HOA_LO_SAO_PHA_LE,
                                        CombineService.TAO_DA_HEMATITE,
                                        CombineService.EP_SAO_TRANG_BI,
                                        CombineService.DAP_DO_AO_HOA,
                                        CombineService.PS_HOA_TRANG_BI,
                                        CombineService.TAY_PS_HOA_TRANG_BI,
                                        CombineService.SIEU_HOA,
                                        CombineService.AN_TRANG_BI,
                                        CombineService.TINH_THACH_HOA,
                                        CombineService.NANG_GIAP_LUYEN_TAP,
                                        CombineService.NANG_CAP_DE_MABU,
                                        CombineService.MO_KHOA_ITEM,
                                        CombineService.GIA_HAN_VAT_PHAM,
                                        CombineService.NANG_CAP_CHAN_MENH -> {
                                    switch (select) {
                                        case 0:
                                            CombineService.gI().startCombine(player);
                                            break;
                                    }
                                }
                            }
                        }
                    }
                }
               case 112 -> {
    if (player.iDMark.isBaseMenu()) {

        if (player.haveRewardVDST) {
            switch (select) {
                case 0 -> {
                    if (InventoryService.gI().getCountEmptyBag(player) > 0) {
                        Item item = ItemService.gI().createNewItem((short) (Util.nextInt(702, 704)));
                        item.itemOptions.add(new Item.ItemOption(30, 1));
                        InventoryService.gI().addItemBag(player, item);
                        InventoryService.gI().sendItemBag(player);
                        Service.gI().sendThongBao(player, "Bạn nhận được " + item.template.name);
                        player.haveRewardVDST = false;
                    } else {
                        Service.gI().sendThongBao(player,
                                "Hành trang không còn chỗ trống, không thể nhặt thêm");
                    }
                }

                case 1 -> {
                    if (InventoryService.gI().getCountEmptyBag(player) > 0) {
                        int rand = Util.nextInt(0, 1);
                        short itemId = (rand == 0 ? (short) 1284 : (short) 1285);
                        Item item = ItemService.gI().createNewItem(itemId, 3);
                        InventoryService.gI().addItemBag(player, item);
                        InventoryService.gI().sendItemBag(player);
                        Service.gI().sendThongBao(player, "Bạn nhận được x5 " + item.template.name);
                        player.haveRewardVDST = false;
                    } else {
                        Service.gI().sendThongBao(player,
                                "Hành trang không còn chỗ trống, không thể nhặt thêm");
                    }
                }
            }
            return;
        }

        if (DeathOrAliveArenaManager.gI().getVDST(player.zone) != null) {

            if (DeathOrAliveArenaManager.gI().getVDST(player.zone).getPlayer().equals(player)) {

                switch (select) {
                    case 0 ->
                        this.npcChat("Không thể thực hiện");
                    case 1 -> {
                    }
                    case 2 ->
                        ChangeMapService.gI().changeMapBySpaceShip(player, 5, -1, 1156);
                }
                return;
            }

            switch (select) {
                case 0 ->
                    DeathOrAliveArenaService.gI().startChallenge(player);
                case 1 -> {
                }
                case 2 ->
                    ChangeMapService.gI().changeMapBySpaceShip(player, 5, -1, 1156);
            }
            return;
        }

        switch (select) {
            case 0 ->
                DeathOrAliveArenaService.gI().startChallenge(player);
            case 1 -> {
            }
            case 2 ->
                ChangeMapService.gI().changeMapBySpaceShip(player, 5, -1, 1156);
        }
    }
                }
                case 174, 181 -> {
                    if (player.iDMark.isBaseMenu()) {
                        switch (select) {
                            case 0 ->
                                ChangeMapService.gI().changeMapBySpaceShip(player, 5, -1, 1156);
                        }
                    }
                }
                case 42, 43, 44, 84 -> {
                    if (player.iDMark.isBaseMenu()) {
                        switch (select) {
                            case 0:
                                createOtherMenu(player, ConstNpc.MENU_SACH_TUYET_KY, "Ta có thể giúp gì cho ngươi ?",
                                        // "Đóng thành\nSách cũ",
                                        "Đổi Sách\nTuyệt kỹ",
                                        "Giám Định\nSách",
                                        "Tẩy\nSách");
                                // "Nâng cấp\nSách\nTuyệt kỹ",
                                // "Hồi phục\nSách",
                                // "Phân rã\nSách");
                                break;
                            case 1:
                                createOtherMenu(player, ConstNpc.MENU_OPTION_SHOP_BUA,
                                        "Bùa của ta rất lợi hại, nhìn ngươi yếu đuối thế này, chắc muốn mua bùa để "
                                                + "mạnh mẽ à, mua không ta bán cho, xài rồi lại thích cho mà xem.",
                                        // "Bùa\n1 giờ",
                                        // "Bùa\n8 giờ",
                                         "Bùa\n1 tháng",
                                        "Close");
                                break;
                            case 2:
                                CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_VAT_PHAM);
                                break;
                            case 3:
                                CombineService.gI().openTabCombine(player, CombineService.NANG_CAP_BONG_TAI);
                                break;
                            case 4:
                                CombineService.gI().openTabCombine(player, CombineService.NANG_CHI_SO_BONG_TAI);
                                break;
                            // case 6:
                            // CombineService.gI().openTabCombine(player, CombineService.LAM_PHEP_NHAP_DA);
                            // break;
                            case 5:
                                CombineService.gI().openTabCombine(player, CombineService.NHAP_NGOC_RONG);
                                break;
                        }
                    } else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_SACH_TUYET_KY) {
                        switch (select) {
                            // case 0:
                            // CheTaoCuonSachCu.showCombine(player);
                            // break;
                            case 0:
                                DoiSachTuyetKy.showCombine(player);
                                break;
                            case 1:
                                CombineService.gI().openTabCombine(player, CombineService.GIAM_DINH_SACH);
                                break;
                            case 2:
                                CombineService.gI().openTabCombine(player, CombineService.TAY_SACH);
                                break;
                            // case 4:
                            // CombineService.gI().openTabCombine(player,
                            // CombineService.NANG_CAP_SACH_TUYET_KY);
                            // break;
                            // case 5:
                            // CombineService.gI().openTabCombine(player, CombineService.HOI_PHUC_SACH);
                            // break;
                            // case 6:
                            // CombineService.gI().openTabCombine(player, CombineService.PHAN_RA_SACH);
                            // break;
                        }
                    } else if (player.iDMark.getIndexMenu() == ConstNpc.DONG_THANH_SACH_CU) {
                        CheTaoCuonSachCu.cheTaoCuonSachCu(player);
                    } else if (player.iDMark.getIndexMenu() == ConstNpc.DOI_SACH_TUYET_KY) {
                        DoiSachTuyetKy.doiSachTuyetKy(player);
                    } else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_OPTION_SHOP_BUA) {
                        switch (select) {
                            // case 0 ->
                            // ShopService.gI().opendShop(player, "BUA_1H", true);
                            // case 0 ->
                            //     ShopService.gI().opendShop(player, "BUA_8H", true);
                            case 0 ->
                            ShopService.gI().opendShop(player, "BUA_1M", true);
                        }
                    } else if (player.iDMark.getIndexMenu() == ConstNpc.MENU_START_COMBINE) {
                        switch (player.combine.typeCombine) {
                            case CombineService.NANG_CAP_BONG_TAI, CombineService.NANG_CHI_SO_BONG_TAI,
                                    CombineService.LAM_PHEP_NHAP_DA, CombineService.NHAP_NGOC_RONG,
                                    CombineService.GIAM_DINH_SACH, CombineService.TAY_SACH,
                                    CombineService.NANG_CAP_SACH_TUYET_KY, CombineService.HOI_PHUC_SACH,
                                    CombineService.PHAN_RA_SACH -> {
                                if (select == 0) {
                                    CombineService.gI().startCombine(player);
                                }
                            }
                            case CombineService.NANG_CAP_VAT_PHAM -> {
                                if (select == 0) {
                                    CombineService.gI().startCombine(player);
                                } else if (select == 1) {
                                    NangCapVatPham.nangCapVatPham(player, true);
                                }
                            }
                        }
                    }
                }
                default -> {
                }
            }
        }
    }
}
