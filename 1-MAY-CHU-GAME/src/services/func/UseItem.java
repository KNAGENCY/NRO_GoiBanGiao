package services.func;

/*
 *
 *
 * @author GiaCatLuong
 */
import consts.ConstItem;
import models.Combine.CombineService;
import models.ShenronEvent.ShenronEventService;
import models.Card.Card;
import models.Card.RadarService;
import models.Card.RadarCard;
import consts.ConstMap;
import item.Item;
import item.ItemTime;
import consts.ConstNpc;
import jdbc.DBConnecter;
import utils.Util;

import consts.ConstPlayer;
import item.Item.ItemOption;
import map.Zone;
import player.Inventory;
import services.*;
import player.Player;
import skill.Skill;
import network.Message;
import npc.specialnpc.MabuEgg;
import utils.SkillUtil;
import utils.TimeUtil;
import server.io.MySession;
import utils.Logger;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Random;

import boss.BossManager;

public class UseItem {

    private static final int ITEM_BOX_TO_BODY_OR_BAG = 0;
    private static final int ITEM_BAG_TO_BOX = 1;
    private static final int ITEM_BODY_TO_BOX = 3;
    private static final int ITEM_BAG_TO_BODY = 4;
    private static final int ITEM_BODY_TO_BAG = 5;
    private static final int ITEM_BAG_TO_PET_BODY = 6;
    private static final int ITEM_BODY_PET_TO_BAG = 7;

    private static final byte DO_USE_ITEM = 0;
    private static final byte DO_THROW_ITEM = 1;
    private static final byte ACCEPT_THROW_ITEM = 2;
    private static final byte ACCEPT_USE_ITEM = 3;

    private static UseItem instance;

    private int randClothes(int level) {
        return ConstItem.LIST_ITEM_CLOTHES[Util.nextInt(0, 2)][Util.nextInt(0, 4)][level - 1];
    }

    private UseItem() {

    }

    public static UseItem gI() {
        if (instance == null) {
            instance = new UseItem();
        }
        return instance;
    }

    public void getItem(MySession session, Message msg) {
        Player player = session.player;
        if (player == null) {
            return;
        }
        TransactionService.gI().cancelTrade(player);
        try {
            int type = msg.reader().readByte();
            int index = msg.reader().readByte();
            if (index == -1) {
                return;
            }
            switch (type) {
                case ITEM_BOX_TO_BODY_OR_BAG:
                    InventoryService.gI().itemBoxToBodyOrBag(player, index);
                    TaskService.gI().checkDoneTaskGetItemBox(player);
                    break;
                case ITEM_BAG_TO_BOX:
                    InventoryService.gI().itemBagToBox(player, index);
                    break;
                case ITEM_BODY_TO_BOX:
                    InventoryService.gI().itemBodyToBox(player, index);
                    break;
                case ITEM_BAG_TO_BODY:
                    InventoryService.gI().itemBagToBody(player, index);
                    break;
                case ITEM_BODY_TO_BAG:
                    InventoryService.gI().itemBodyToBag(player, index);
                    break;
                case ITEM_BAG_TO_PET_BODY:
                    InventoryService.gI().itemBagToPetBody(player, index);
                    break;
                case ITEM_BODY_PET_TO_BAG:
                    InventoryService.gI().itemPetBodyToBag(player, index);
                    break;
            }
            if (player.setClothes != null) {
                player.setClothes.setup();
            }
            if (player.pet != null) {
                player.pet.setClothes.setup();
            }
            player.setClanMember();
            Service.gI().sendFlagBag(player);
            Service.gI().point(player);
            Service.gI().sendSpeedPlayer(player, -1);
        } catch (Exception e) {
            Logger.logException(UseItem.class, e);

        }
    }

    public Item finditem(Player player, int iditem) {
        for (Item item : player.inventory.itemsBag) {
            if (item.isNotNullItem() && item.template.id == iditem) {
                return item;
            }
        }
        return null;
    }

    public void doItem(Player player, Message _msg) {
        TransactionService.gI().cancelTrade(player);
        Message msg = null;
        byte type;
        try {
            type = _msg.reader().readByte();
            int where = _msg.reader().readByte();
            int index = _msg.reader().readByte();
            switch (type) {
                case DO_USE_ITEM:
                    if (player != null && player.inventory != null) {
                        if (index != -1) {
                            if (index < 0) {
                                return;
                            }
                            Item item = player.inventory.itemsBag.get(index);
                            if (item.isNotNullItem()) {
                                if (item.template.type == 7) {
                                    msg = new Message(-43);
                                    msg.writer().writeByte(type);
                                    msg.writer().writeByte(where);
                                    msg.writer().writeByte(index);
                                    msg.writer().writeUTF("Bạn chắc chắn học "
                                            + player.inventory.itemsBag.get(index).template.name + "?");
                                    player.sendMessage(msg);
                                } else if (item.template.id == 570) {
                                    // if (!Util.isAfterMidnight(player.lastTimeRewardWoodChest)) {
                                    // Service.gI().sendThongBao(player, "Hãy chờ đến ngày mai");
                                    // return;
                                    // }
                                    msg = new Message(-43);
                                    msg.writer().writeByte(type);
                                    msg.writer().writeByte(where);
                                    msg.writer().writeByte(index);
                                    msg.writer().writeUTF("Bạn chắc muốn mở\n"
                                            + player.inventory.itemsBag.get(index).template.name + " ?");
                                    player.sendMessage(msg);
                                } else if (item.template.type == 22) {
                                    if (player.zone.items.stream()
                                            .filter(it -> it != null && it.itemTemplate.type == 22).count() > 2) {
                                        Service.gI().sendThongBaoOK(player, "Mỗi map chỉ đặt được 3 Vệ Tinh");
                                        return;
                                    }
                                    msg = new Message(-43);
                                    msg.writer().writeByte(type);
                                    msg.writer().writeByte(where);
                                    msg.writer().writeByte(index);
                                    msg.writer().writeUTF("Bạn chắc muốn dùng\n"
                                            + player.inventory.itemsBag.get(index).template.name + " ?");
                                    player.sendMessage(msg);
                                } else {
                                    UseItem.gI().useItem(player, item, index);
                                }
                            }
                        } else {
                            int iditem = _msg.reader().readShort();
                            Item item = finditem(player, iditem);
                            UseItem.gI().useItem(player, item, index);
                        }
                    }
                    break;
                case DO_THROW_ITEM:
                    if (!(player.zone.map.mapId == 21 || player.zone.map.mapId == 22 || player.zone.map.mapId == 23)) {
                        Item item = null;
                        if (index < 0) {
                            return;
                        }
                        if (where == 0) {
                            item = player.inventory.itemsBody.get(index);
                        } else {
                            item = player.inventory.itemsBag.get(index);
                        }

                        if (item.isNotNullItem() && item.template.id == 570) {
                            Service.gI().sendThongBao(player, "Không thể bỏ vật phẩm này.");
                            return;
                        }
                        if (!item.isNotNullItem()) {
                            return;
                        }
                        msg = new Message(-43);
                        msg.writer().writeByte(type);
                        msg.writer().writeByte(where);
                        msg.writer().writeByte(index);
                        msg.writer().writeUTF("Bạn chắc chắn muốn vứt " + item.template.name + "?");
                        player.sendMessage(msg);
                    } else {
                        Service.gI().sendThongBao(player, "Không thể thực hiện");
                    }
                    break;
                case ACCEPT_THROW_ITEM:
                    InventoryService.gI().throwItem(player, where, index);
                    Service.gI().point(player);
                    InventoryService.gI().sendItemBag(player);
                    break;
                case ACCEPT_USE_ITEM:
                    UseItem.gI().useItem(player, player.inventory.itemsBag.get(index), index);
                    break;
            }
        } catch (Exception e) {
            Logger.logException(UseItem.class, e);
        } finally {
            if (msg != null) {
                msg.cleanup();
            }
        }
    }

    private void useItem(Player pl, Item item, int indexBag) {
        if (item != null && item.isNotNullItem()) {
            if (item.template.id == 570) {
                int time = (int) TimeUtil.diffDate(new Date(), new Date(item.createTime), TimeUtil.DAY);
                if (time == 0) {
                    Service.gI().sendThongBao(pl, "Hãy chờ đến ngày mai");
                } else {
                    openRuongGo(pl, item);
                }
                return;
            }
            if (item.template.strRequire <= pl.nPoint.power) {
                switch (item.template.type) {

                    case 21:
                        InventoryService.gI().itemBagToBody(pl, indexBag);
                        PetService.Pet2(pl, pl.getHeadThuCung(), pl.getBodyThuCung(), pl.getLegThuCung());
                        Service.gI().point(pl);
                        break;
                    case 33: // card
                        UseCard(pl, item);
                        break;
                    case 7: // sách học, nâng skill
                        learnSkill(pl, item);
                        break;
                    case 6: // đậu thần
                        this.eatPea(pl);
                        break;
                    case 12: // ngọc rồng các loại
                        controllerCallRongThan(pl, item);
                        break;
                    case 23: // thú cưỡi mới
                    case 24: // thú cưỡi cũ
                        InventoryService.gI().itemBagToBody(pl, indexBag);
                        break;
                    case 11: // item bag
                        InventoryService.gI().itemBagToBody(pl, indexBag);

                        Service.gI().sendFlagBag(pl);

                        break;
                    case 75:
                        InventoryService.gI().itemBagToBody(pl, indexBag);
                        Service.gI().sendchienlinh(pl, (short) (item.template.iconID - 1));
                        break;
                    case 72: {
                        InventoryService.gI().itemBagToBody(pl, indexBag);
                        Service.gI().sendPetFollow(pl, (short) (item.template.iconID - 1));
                        break;
                    }
                    case 98: {
                        InventoryService.gI().itemBagToBody(pl, indexBag);
                        Service.gI().sendEffPlayer(pl);
                        break;
                    }
                    case 99: {
                        InventoryService.gI().itemBagToBody(pl, indexBag);
                        Service.gI().sendEffPlayer(pl);
                        break;
                    }
                    default:
                        switch (item.template.id) {
                            case 992: // Nhan thoi khong
                                pl.type = 2;
                                pl.maxTime = 5;
                                Service.gI().Transport(pl);
                                break;
                            case 361:
                                pl.idGo = (short) Util.nextInt(0, 6);
                                NgocRongNamecService.gI().menuCheckTeleNamekBall(pl);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                InventoryService.gI().sendItemBag(pl);
                                break;
                            // case 942:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 966, 967, 968);
                            // Service.gI().point(pl);
                            // break;
                            // case 943:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 969, 970, 971);
                            // Service.gI().point(pl);
                            // break;
                            // case 944:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 972, 973, 974);
                            // Service.gI().point(pl);
                            // break;
                            // case 967:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 1050, 1051, 1052);
                            // Service.gI().point(pl);
                            // break;
                            // case 1107:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 1183, 1184, 1185);
                            // Service.gI().point(pl);
                            // break;
                            // case 1686:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 1572, 1573, 1574);
                            // Service.gI().point(pl);
                            // break;
                            // case 1683:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 1561, 1562, 1563);
                            // Service.gI().point(pl);
                            // break;
                            // case 1682:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, 1558, 1559, 1560);
                            // Service.gI().point(pl);
                            // break;
                            // case 942:
                            // case 943:
                            // case 944:
                            // case 967:
                            // case 1107:
                            // case 1686:
                            // case 1683:
                            // case 1682:
                            // InventoryService.gI().itemBagToBody(pl, indexBag);
                            // PetService.Pet2(pl, item.template.head, item.template.body,
                            // item.template.leg);
                            // Service.gI().point(pl);
                            // break;
                            case 211: // nho tím
                            case 212: // nho xanh
                                eatGrapes(pl, item);
                                break;
                            case 342:
                            case 343:
                            case 344:
                            case 345:
                                if (pl.zone.items.stream().filter(it -> it != null && it.itemTemplate.type == 22)
                                        .count() < 3) {
                                    Service.gI().dropSatellite(pl, item, pl.zone, pl.location.x, pl.location.y);
                                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                } else {
                                    Service.gI().sendThongBaoOK(pl, "Mỗi map chỉ đặt được 3 Vệ Tinh");
                                }
                                break;
                            case 380: // cskb
                                openCSKB(pl, item);
                                break;
                            case 381: // cuồng nộ
                            case 382: // bổ huyết
                            case 383: // bổ khí
                            case 384: // giáp xên
                            case 385: // ẩn danh
                            case 379: // máy dò capsule
                            case 638: // commeson
                            case 2075: // rocket
                            case 2160: // Nồi cơm điện
                            case 579:
                            case 1045: // đuôi khỉ
                            case 663: // bánh pudding
                            case 664: // xúc xíc
                            case 665: // kem dâu
                            case 666: // mì ly
                            case 667: // sushi
                            case 1150:
                            case 1151:
                            case 1152:
                            case 1153:
                            case 1154:
                            case 1628:
                            case 764:
                            case 1731:
                            case 1727:
                            case 1728:
                            case 1729:
                            case 1730:
                            case 1195:
                            case 1196:
                            case 752:
                            case 753:
                                useItemTime(pl, item);
                                break;
                            case 1560:
                                openRuongNgocRong(pl, item);
                                break;
                            case 1770:
                                ruongnro(pl, item);
                                break;
                            case 880:
                            case 881:
                            case 882:
                                if (pl.itemTime.isEatMeal2) {
                                    Service.gI().sendThongBao(pl, "Chỉ được sử dụng 1 cái");
                                    break;
                                }
                                useItemTime(pl, item);
                                break;
                            case 899:
                            case 900:
                            case 902:
                            case 903:
                                if (pl.itemTime.isEatMeal3) {
                                    Service.gI().sendThongBao(pl, "Chỉ được sử dụng 1 cái");
                                    break;
                                }
                                useItemTime(pl, item);
                                break;
                            case 1732:
                                openxukrai(pl, item, 100);
                                break;
                            case 521: // tdlt
                                useTDLT(pl, item);
                                break;
                            case 454: // bông tai
                                if (pl.zone != null && (pl.zone.map.mapId == 182 || pl.zone.map.mapId == 168)) {
                                    Service.gI().sendThongBao(pl, "Không thể sử dụng vật phẩm này trong map War");
                                    break;
                                }
                                UseItem.gI().usePorata(pl);
                                break;
                            case 569:
                                Service.gI().sendThongBao(pl,
                                        "Haha, th\u1eb1ng ngu, vui l\u00f2ng gi\u1eef b\u00ed m\u1eadt");
                                break;
                            case 921: // bông tai
                                if (pl.zone != null && (pl.zone.map.mapId == 182 || pl.zone.map.mapId == 168)) {
                                    Service.gI().sendThongBao(pl, "Không thể sử dụng vật phẩm này trong map War");
                                    break;
                                }
                                UseItem.gI().usePorata2(pl);
                                break;



                                     case 1498:
                                if (pl.zone != null && (pl.zone.map.mapId == 182 || pl.zone.map.mapId == 168)) {
                                    Service.gI().sendThongBao(pl, "Không thể sử dụng vật phẩm này trong map War");
                                    break;
                                }
                                UseItem.gI().usePorata(pl, (byte) 10);
                                break;
                            case 1499:
                                if (pl.zone != null && (pl.zone.map.mapId == 182 || pl.zone.map.mapId == 168)) {
                                    Service.gI().sendThongBao(pl, "Không thể sử dụng vật phẩm này trong map War");
                                    break;
                                }
                                UseItem.gI().usePorata(pl, (byte) 9);
                                break;
                          case 193: // gói 10 viên capsule

    // 🔒 CHẶN Ở MAP WAR
    if (pl.zone != null) {
        int mapId = pl.zone.map.mapId;
        if (mapId == 168 || mapId == 182) {
            Service.gI().sendThongBao(pl,
                    "Không thể sử dụng Capsule trong khu vực này!");
            return;
        }
    }

    openCapsuleUI(pl);
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    

                          case 194: // capsule đặc biệt

    // 🔒 CHẶN Ở MAP WAR
    if (pl.zone != null) {
        int mapId = pl.zone.map.mapId;
        if (mapId == 168 || mapId == 182) {
            Service.gI().sendThongBao(pl,
                    "Không thể sử dụng Capsule đặc biệt trong khu vực này!");
            return;
        }
    }

    openCapsuleUI(pl);
    break;

case 1715:

    if (pl.pet != null) {
        // đã có đệ -> đổi pet
        PetService.gI().createPetFideNhi(pl, true, pl.gender);
    } else {
        // chưa có -> tạo mới
        PetService.gI().createPetFideNhi(pl, false, pl.gender);
    }

    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    Service.gI().sendThongBao(pl, "Bạn đã thất bại!\nChúc may mắn lần sau");

    break;

case 1714:

    if (pl.pet != null) {
        // đã có đệ -> đổi pet
        PetService.gI().createPetCellNhi(pl, true, pl.gender);
    } else {
        // chưa có -> tạo mới
        PetService.gI().createPetCellNhi(pl, false, pl.gender);
    }

    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    Service.gI().sendThongBao(pl, "Bạn đã thất bại!\nChúc may mắn lần sau");

    break;



 case 1717:
case 1716:
    if (pl.pet == null) {
        Service.gI().sendThongBao(pl, "Yêu cầu phải có đệ tử");
        break;
    }

    int menuId = item.template.id; // 1716 hoặc 1717

 NpcService.gI().createMenuConMeo(
    pl,
    menuId,
    16132, // ✅ avatar
    "Chọn hành tinh cho Nghịch tử",
    "Trái Đất", "Namek", "Xayda", "Hủy"
);


    break;

case 1721:
    NpcService.gI().createMenuConMeo(
            pl,
            1721,
            16132,
            "Chọn hành tinh cho Đệ tử đi con gà:",
            "Trái Đất", "Namek", "Xayda", "Hủy"
    );
    break;

                            case 568: // Trứng Mabu
                                if (pl.mabuEgg == null) {
                                    MabuEgg.createMabuEgg(pl);
                                    Service.gI().sendThongBao(pl, "Trứng đang được ấp ở nhà rồi kìa");
                                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                    InventoryService.gI().sendItemBag(pl);
                                    pl.mabuEgg.sendMabuEgg();
                                } else {
                                    Service.gI().sendThongBao(pl, "Mi có trứng rồi mà??");
                                }
                                break;


    //      case 1722: // Vé cộng cash account
    // int[] rewards = {2000, 5000, 10000};
    // int cashAdd = rewards[Util.nextInt(rewards.length)];

    // try {
    //     // Cộng cash vào database
    //     DBConnecter.executeUpdate(
    //         "UPDATE account SET cash = cash + " + cashAdd +
    //         " WHERE id = " + pl.getSession().userId
    //     );

    //     // Cộng cash trong session (khỏi relog)
    //     pl.getSession().cash += cashAdd;

    //     // Trừ item
    //     InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    //     InventoryService.gI().sendItemBag(pl);

    //     // Thông báo
    //     Service.gI().sendThongBao(
    //         pl,
    //         "Bạn nhận được +" + Util.numberToMoney(cashAdd) + " VNĐ"
    //     );

    // } catch (Exception e) {
    //     // Nếu lỗi DB → KHÔNG trừ item
    //     Service.gI().sendThongBao(
    //         pl,
    //         "Có lỗi khi cộng VNĐ, vui lòng thử lại"
    //     );
    //     e.printStackTrace();
    // }
    // break;



                            case 401: // đổi đệ tử
                                changePet(pl, item);
                                break;
                            case 402: // sách nâng chiêu 1 đệ tử
                            case 403: // sách nâng chiêu 2 đệ tử
                            case 404: // sách nâng chiêu 3 đệ tử
                            case 759: // sách nâng chiêu 4 đệ tử
                                upSkillPet(pl, item);
                                break;
                            /// thay skill 2 đệ tử
                            case 1744:
                                replaceSkill2(pl, item, Skill.KAMEJOKO, "Kamejoko");
                                break;
                            case 1745:
                                replaceSkill2(pl, item, Skill.ANTOMIC, "Antomic");
                                break;
                            case 1746:
                                replaceSkill2(pl, item, Skill.MASENKO, "Masenko");
                                break;
                            /// thay skill 3 đệ tử
                           case 1747:
    replaceSkill3(pl, item, Skill.THAI_DUONG_HA_SAN, "Thái dương hạ sang");
    break; // ✅ thêm break
case 1748:
    replaceSkill3(pl, item, Skill.TAI_TAO_NANG_LUONG, "Tái tạo năng lượng");
    break; // ✅ thêm break
case 1749:
    replaceSkill3(pl, item, Skill.KAIOKEN, "Kaioken");
    break; // ✅ THÊM BREAK — đây là lỗi chính gây ra fall-through
case 1750:
    replaceSkill4(pl, item, Skill.BIEN_KHI, "Biến hình");
    break;
                            case 1751:
                                replaceSkill4(pl, item, Skill.TROI, "Trói");
                                break;
                            case 1752:
                                replaceSkill4(pl, item, Skill.KHIEN_NANG_LUONG, "Khiên năng lượng");
                                break;
                            case 1753:
                                int[] skillItems = { 1750, 1751, 1752 };
                                if (InventoryService.gI().getCountEmptyBag(pl) == 0) {
                                    Service.gI().sendThongBaoOK(pl, "Cần 1 ô hành trang để mở");
                                    return;
                                }
                                Item sachthayskill = ItemService.gI()
                                        .createNewItem((short) skillItems[Util.nextInt(skillItems.length)]);
                                sachthayskill.itemOptions.add(new Item.ItemOption(30, 1));

                                InventoryService.gI().addItemBag(pl, sachthayskill);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                InventoryService.gI().sendItemBag(pl);
                                Service.gI().sendThongBao(pl, "Bạn đã nhận được " + sachthayskill.template.name);
                                break;
                              case 1743: {
    if (InventoryService.gI().getCountEmptyBag(pl) < 1) {
        Service.gI().sendThongBaoOK(pl, "Cần 1 ô hành trang để mở");
        return;
    }

    int rd = Util.nextInt(100);
    int itemId;

    if (rd < 20) {
        itemId = 1744; // 20%
    } else if (rd < 50) {
        itemId = 1745; // 30%
    } else {
        itemId = 1746; // 50%
    }

    Item itemNhan = ItemService.gI().createNewItem((short) itemId);

    // Option 30 param 1
    itemNhan.itemOptions.add(new Item.ItemOption(30, 1));

    InventoryService.gI().addItemBag(pl, itemNhan);
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    InventoryService.gI().sendItemBag(pl);

    Service.gI().sendThongBao(pl,
            "Bạn nhận được " + itemNhan.template.name);

    break;
}

 case 1755: {
    if (InventoryService.gI().getCountEmptyBag(pl) < 1) {
        Service.gI().sendThongBaoOK(pl, "Cần 1 ô hành trang để mở");
        return;
    }

    int rd = Util.nextInt(100);
    int itemId;

    if (rd < 20) {
        itemId = 1747; // 20%
    } else if (rd < 50) {
        itemId = 1748; // 30%
    } else {
        itemId = 1749; // 50%
    }

    Item itemNhan = ItemService.gI().createNewItem((short) itemId);

    // Option 30 param 1
    itemNhan.itemOptions.add(new Item.ItemOption(30, 1));

    InventoryService.gI().addItemBag(pl, itemNhan);
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    InventoryService.gI().sendItemBag(pl);

    Service.gI().sendThongBao(pl,
            "Bạn nhận được " + itemNhan.template.name);

    break;
}






                           case 726:
    if (pl.nPoint != null && pl.nPoint.power < 1_500_000_000L) {
        Service.gI().sendThongBao(pl, "Cần 1tỷ5 sức mạnh");
        return;
    }
    UseItem.gI().ItemManhGiay(pl, item);
    break;
                            case 727:
                            case 728:
                                UseItem.gI().ItemSieuThanThuy(pl, item);
                                break;
                            case 648:
                                ItemService.gI().OpenItem648(pl, item);
                                break;
                            case 1759:
                                ItemService.gI().OpenItem1759(pl, item);
                                break;
                            case 1757:
                                Service.gI().releaseCooldownSkill(pl);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                InventoryService.gI().sendItemBag(pl);
                                Service.gI().sendThongBao(pl, "Đã hồi lại toàn bộ Skill.");
                                break;
                            case 1576:
                                ItemService.gI().OpenItem1576(pl, item);
                                break;
                            case 1758:
                                NpcService.gI().createMenuConMeo(pl, ConstNpc.MENU_OPTION_USE_ITEM1758, -1,
                                        "Hãy chọn Set Thần Linh ngươi cần:",
                                        "Trái Đất", "Namek", "Xayda", "Từ chối");
                                break;
                            case 736:
                                ItemService.gI().OpenItem736(pl, item);
                                break;
                            case 987:
                                Service.gI().sendThongBao(pl, "Bảo vệ trang bị không bị rớt cấp"); // đá bảo vệ
                                break;
                            case 2006:
                                Input.gI().createFormChangeNameByItem(pl);
                                break;
                            case 1623:
                                TaskService.gI().sendNextTaskMain(pl);
                                break;
                            case 725:
                                Item item725 = InventoryService.gI().findItemBag(pl, 725);
                                if (item725 != null && item725.quantity > 0) {
                                    InventoryService.gI().subQuantityItemsBag(pl, item725, 1);
                                    InventoryService.gI().sendItemBag(pl);
                                    Service.gI().addSMTN(pl, (byte) 2, 1_000_000L, false);
                                    Service.gI().sendThongBao(pl, "Mày đã bị lừa");
                                } else {
                                    Service.gI().sendThongBao(pl, "Bạn không có vật phẩm cần để sử dụng!");
                                }
                                break;
                            // case 1228:
                            //     NpcService.gI().createMenuConMeo(pl, ConstNpc.HOP_QUA_THAN_LINH, -1,
                            //             "Chọn hành tinh của đồ thần linh muốn nhận.",
                            //             "Trái đất", "Namek", "Xayda");
                            //     break;
    case 1239: { // Item xem danh sách Boss (yêu cầu SM >= 20 tỷ)
    long requireSM = 20_000_000_000L;

    if (pl.nPoint.power < requireSM) {
        Service.gI().sendThongBao(pl,
                "Yêu cầu sức mạnh tối thiểu 20 tỷ để sử dụng vật phẩm này!");
        break;
    }

    BossManager.gI().showListBoss(pl);
    break; // KHÔNG trừ item
}


                                 case 1238: { // ✅ ITEM SỬA / NEXT NHIỆM VỤ
    if (pl.playerTask == null || pl.playerTask.taskMain == null) {
        Service.gI().sendThongBao(pl, "Không có nhiệm vụ để hỗ trợ!");
        break;
    }

    // Reset bước hiện tại
    pl.playerTask.taskMain.index = 0;

    // NEXT nhiệm vụ theo logic hệ thống
    TaskService.gI().sendNextTaskMain(pl);

    // Trừ item
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    InventoryService.gI().sendItemBag(pl);

    Service.gI().sendThongBao(pl, "Bạn đã được hỗ trợ nhiệm vụ thành công!");
    break;
}     
                            case 1626: {
                                int[] listItem = { 856, 943, 942 };
                                if (InventoryService.gI().getCountEmptyBag(pl) == 0) {
                                    Service.gI().sendThongBaoOK(pl, "Cần 1 ô hành trang để mở");
                                    return;
                                }
                                Item phuKien = ItemService.gI().createNewItem((short) listItem[Util.nextInt(2)]);
                                if (phuKien.template.id == 856) {
                                    phuKien.itemOptions.add(new Item.ItemOption(50, 10));
                                    phuKien.itemOptions.add(new Item.ItemOption(77, 10));
                                    phuKien.itemOptions.add(new Item.ItemOption(103, 10));
                                } else if (phuKien.template.id == 943) {
                                    phuKien.itemOptions.add(new Item.ItemOption(50, 10));
                                } else if (phuKien.template.id == 942) {
                                    phuKien.itemOptions.add(new Item.ItemOption(77, 10));
                                    phuKien.itemOptions.add(new Item.ItemOption(103, 10));
                                }
                                if (Util.isTrue(95, 100)) {
                                    phuKien.itemOptions.add(new Item.ItemOption(93, Util.nextInt(1, 5)));
                                }
                                InventoryService.gI().addItemBag(pl, phuKien);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                InventoryService.gI().sendItemBag(pl);
                                Service.gI().sendThongBao(pl, "Bạn đã nhận được " + phuKien.template.name);
                            }
                                break;
                            // case 1628: {
                            // Player player = pl;
                            // if (player.pet != null) {
                            // if (player.pet.playerSkill.skills.get(1).skillId != -1) {
                            // player.pet.openSkill2();
                            // } else {
                            // Service.gI().sendThongBao(player, "Ít nhất đệ tử ngươi phải có chiêu 2
                            // chứ!");
                            // return;
                            // }
                            // } else {
                            // Service.gI().sendThongBao(player, "Ngươi làm gì có đệ tử?");
                            // return;
                            // }
                            // }
                            // break;
                            case 1629: {
                                Player player = pl;
                                if (player.pet != null) {
                                    if (player.pet.playerSkill.skills.get(2).skillId != -1) {
                                        player.pet.openSkill3();
                                    } else {
                                        Service.gI().sendThongBao(player, "Ít nhất đệ tử ngươi phải có chiêu 3 chứ!");
                                        return;
                                    }
                                } else {
                                    Service.gI().sendThongBao(player, "Ngươi làm gì có đệ tử?");
                                    return;
                                }
                            }
                                break;
                            case 1630: {
                                Player player = pl;
                                if (player.pet != null) {
                                    if (player.pet.playerSkill.skills.get(3).skillId != -1) {
                                        player.pet.openSkill4();
                                    } else {
                                        Service.gI().sendThongBao(player, "Ít nhất đệ tử ngươi phải có chiêu 4 chứ!");
                                        return;
                                    }
                                } else {
                                    Service.gI().sendThongBao(player, "Ngươi làm gì có đệ tử?");
                                    return;
                                }
                            }
                                break;
                            case 628: {
                                int ct = Util.nextInt(618, 626);
                                Item caiTrangHaiTac = ItemService.gI().createNewItem((short) ct);
                                caiTrangHaiTac.itemOptions.add(new Item.ItemOption(93, 30));
                                caiTrangHaiTac.itemOptions.add(new Item.ItemOption(50, 15));
                                caiTrangHaiTac.itemOptions.add(new Item.ItemOption(77, 15));
                                caiTrangHaiTac.itemOptions.add(new Item.ItemOption(103, 15));
                                caiTrangHaiTac.itemOptions.add(new Item.ItemOption(149, 1));
                                InventoryService.gI().addItemBag(pl, caiTrangHaiTac);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                Service.gI().sendThongBao(pl,
                                        "Bạn đã nhận được cải trang " + caiTrangHaiTac.template.name);
                            }
                                break;
                            case 1440: {
                                int ct = Util.nextInt(441, 447);
                                Item caiTrangHaiTac = ItemService.gI().createNewItem((short) ct);
                                caiTrangHaiTac.itemOptions.add(new Item.ItemOption(93, 30));
                                InventoryService.gI().addItemBag(pl, caiTrangHaiTac);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                Service.gI().sendThongBao(pl, "Bạn đã nhận được " + caiTrangHaiTac.template.name);
                            }
                                break;
                            case 1453: {
                                int tile = Util.nextInt(1, 100);
                                int itemId = -1;
                                if (tile <= 50) {
                                    itemId = 1283;
                                } else if (tile <= 70) {
                                    itemId = 579;
                                } else if (tile <= 85) {
                                    itemId = 457;
                                } else if (tile <= 90) {
                                    itemId = 674;
                                } else {
                                    itemId = 1718;
                                }
                                if (itemId != -1) {
                                    Item newItem = ItemService.gI().createNewItem((short) itemId);
                                    InventoryService.gI().addItemBag(pl, newItem);
                                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                    Service.gI().sendThongBao(pl, "Bạn đã nhận được " + newItem.template.name);
                                }
                            }
                                break;
                            case 869: {
                                pl.inventory.gold += 500_000_000L;
                                InventoryService.gI().sendItemBag(pl);
                                Service.gI().sendMoney(pl);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                Service.gI().sendThongBaoOK(pl,
                                        "Đừng bao giờ tin lời mấy thằng tên Đạt\n+1 bài học quý giá");
                            }
                                break;
                            case 1173: {
                                int tile = Util.nextInt(1, 100);
                                int itemId = -1;
                                if (tile <= 30) {
                                    itemId = 579;
                                } else if (tile <= 40) {
                                    itemId = 987;
                                } else if (tile <= 50) {
                                    itemId = 1194;
                                } else if (tile <= 60) {
                                    itemId = 753;
                                } else if (tile <= 70) {
                                    itemId = 752;
                                } else if (tile <= 80) {
                                    itemId = 935;
                                } else if (tile <= 90) {
                                    itemId = 343;
                                } else {
                                    itemId = 345;
                                }
                                if (itemId != -1) {
                                    Item newItem = ItemService.gI().createNewItem((short) itemId);
                                    InventoryService.gI().addItemBag(pl, newItem);
                                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                    Service.gI().sendThongBao(pl, "Bạn đã nhận được " + newItem.template.name);
                                }
                            }
                                break;

                            case 1536: {
                            }
                                break;
                            case 1703:// set tl kh

                                UseItem.gI().Hopdothanlinh(pl, item);
                                break; // Name: GiaCatLuong
                            case 1704:// set hd kh

                                UseItem.gI().Hopdohuydiet(pl, item);
                                break; // Name: GiaCatLuong
                        }
                        break;
                }
                TaskService.gI().checkDoneTaskUseItem(pl, item);
                InventoryService.gI().sendItemBag(pl);
            } else {
                Service.gI().sendThongBaoOK(pl, "Sức mạnh không đủ yêu cầu");
            }
        }
    }

    public void openRuongGo(Player pl, Item item) {
        List<String> textRuongGo = new ArrayList<>();
        int time = (int) TimeUtil.diffDate(new Date(), new Date(item.createTime), TimeUtil.DAY);
        if (time != 0) {
            Item itemReward = null;
            int param = item.itemOptions.get(0).param;
            int gold = 0;
            int[] listItem = { 441, 442, 443, 444, 445, 446, 447, 220, 221, 222, 223, 224, 225 };
            int[] listClothesReward;
            int[] listItemReward;
            String text = "Bạn nhận được\n";
            if (param < 8) {
                gold = 100000 * param;
                listClothesReward = new int[] { randClothes(param) };
                listItemReward = Util.pickNRandInArr(listItem, 3);
            } else if (param < 10) {
                gold = 250000 * param;
                listClothesReward = new int[] { randClothes(param), randClothes(param) };
                listItemReward = Util.pickNRandInArr(listItem, 4);
            } else {
                gold = 500000 * param;
                listClothesReward = new int[] { randClothes(param), randClothes(param), randClothes(param) };
                listItemReward = Util.pickNRandInArr(listItem, 5);
                int ruby = Util.nextInt(1, 5);
                pl.inventory.ruby += ruby;
                textRuongGo.add(text + "|1| " + ruby + " Hồng Ngọc");
            }
            for (var i : listClothesReward) {
                itemReward = ItemService.gI().createNewItem((short) i);
                RewardService.gI().initBaseOptionClothes(itemReward.template.id, itemReward.template.type,
                        itemReward.itemOptions);
                RewardService.gI().initStarOption(itemReward, new RewardService.RatioStar[] {
                        new RewardService.RatioStar((byte) 1, 1, 2), new RewardService.RatioStar((byte) 2, 1, 3),
                        new RewardService.RatioStar((byte) 3, 1, 4), new RewardService.RatioStar((byte) 4, 1, 5), });
                InventoryService.gI().addItemBag(pl, itemReward);
                textRuongGo.add(text + itemReward.info);
            }
            for (var i : listItemReward) {
                itemReward = ItemService.gI().createNewItem((short) i);
                RewardService.gI().initBaseOptionSaoPhaLe(itemReward);
                itemReward.quantity = Util.nextInt(1, 5);
                InventoryService.gI().addItemBag(pl, itemReward);
                textRuongGo.add(text + itemReward.info);
            }
            if (param == 11) {
                itemReward = ItemService.gI().createNewItem((short) ConstItem.MANH_NHAN);
                itemReward.quantity = Util.nextInt(1, 3);
                InventoryService.gI().addItemBag(pl, itemReward);
                textRuongGo.add(text + itemReward.info);
            }
            NpcService.gI().createMenuConMeo(pl, ConstNpc.RUONG_GO, -1,
                    "Bạn nhận được\n|1|+" + Util.numberToMoney(gold) + " vàng", "OK [" + textRuongGo.size() + "]");
            InventoryService.gI().subQuantityItemsBag(pl, item, 1);
            pl.inventory.addGold(gold);
            InventoryService.gI().sendItemBag(pl);
            PlayerService.gI().sendInfoHpMpMoney(pl);
        }
    }

    // public void openRuongGo(Player player) {
    // Item ruongGo = InventoryService.gI().findItemBag(player, 570);
    // if (ruongGo != null) {
    // int level = InventoryService.gI().getParam(player, 72, 570);
    // if (InventoryService.gI().getCountEmptyBag(player) < level) {
    // Service.gI().sendThongBao(player, "Cần ít nhất " + level + " ô trống trong
    // hành trang");
    // } else {
    // player.itemsWoodChest.clear();
    // if (level == 0) {
    // InventoryService.gI().subQuantityItemsBag(player, ruongGo, 1);
    // InventoryService.gI().sendItemBag(player);
    // Item item = ItemService.gI().createNewItem((short) 673);
    // InventoryService.gI().addItemBag(player, item);
    // InventoryService.gI().sendItemBag(player);
    // Service.gI().sendThongBao(player, "Bạn vừa nhận được set Lính Thủy Đánh
    // Bạc!");
    // return;
    // }
    // int rand = Util.nextInt(0, 6);
    // Item level1 = ItemService.gI().createNewItem(((short) (441 + rand)));
    // level1.itemOptions.add(new Item.ItemOption(95 + rand, (rand == 3 || rand ==
    // 4) ? 3 : 5));
    // level1.quantity = Util.nextInt(1, level * 2);
    // player.itemsWoodChest.add(level1);
    // if (level > 1) {
    // rand = Util.nextInt(0, 4);
    // Item level2 = ItemService.gI().createNewItem(((short) (220 + rand)));
    // level2.itemOptions.add(new Item.ItemOption(71 - rand, 0));
    // level2.quantity = Util.nextInt(1, level * 3);
    // player.itemsWoodChest.add(level2);
    // }
    // if (level > 2) {
    // Item level3 = ItemService.gI().createNewItem((short) 20); // 7 sao
    // level3.quantity = Util.nextInt(1, level);
    // player.itemsWoodChest.add(level3);
    // }
    // if (level > 3) {
    // Item level4 = ItemService.gI().createNewItem((short) 19); // 6 sao
    // level4.quantity = Util.nextInt(1, level);
    // player.itemsWoodChest.add(level4);
    // }
    // if (level > 4) {
    // Item level5 = ItemService.gI().createNewItem((short) 18); // 5 sao
    // level5.quantity = Util.nextInt(1, level);
    // player.itemsWoodChest.add(level5);
    // }
    // if (level > 5) {
    // Item level6 = ItemService.gI().createNewItem((short) 17); // 4 sao
    // level6.quantity = Util.nextInt(1, level);
    // player.itemsWoodChest.add(level6);
    // }
    // if (level > 6) {
    // Item level7 = ItemService.gI().createNewItem((short) 16); // 3 sao
    // level7.quantity = Util.nextInt(1, level);
    // player.itemsWoodChest.add(level7);
    // }
    // if (level > 7) {
    // Item level8 = ItemService.gI().createNewItem((short) 457); // thoi vang
    // level8.quantity = Util.nextInt(100, level * 100);
    // player.itemsWoodChest.add(level8);
    // }
    // if (level > 8) {
    // Item level9 = ItemService.gI().createNewItem((short) 1229); // bi kiep tuyet
    // ky
    // level9.quantity = Util.nextInt(15, level + 15);
    // player.itemsWoodChest.add(level9);
    // }
    // if (level > 9) {
    // int[] itemId = {1229, 2026, 2036, 2037, 2038, 2039, 2040, 2019, 2020, 2021,
    // 2022, 2023, 2024};
    // byte[] option = {77, 80, 81, 103, 50, 94, 5};
    // byte[] option_v2 = {14, 16, 17, 19, 27, 28, 5, 47, 87}; //77 %hp // 80 //81
    // //103 //50 //94 //5 % sdcm
    // byte optionid;
    // byte optionid_v2;
    // byte param;
    // Item level10 = ItemService.gI().createNewItem((short) itemId[Util.nextInt(0,
    // 12)]);
    // level10.itemOptions.clear();
    // optionid = option[Util.nextInt(0, 6)];
    // param = (byte) Util.nextInt(5, 10);
    // level10.itemOptions.add(new Item.ItemOption(optionid, param));
    // if (Util.isTrue(20, 100)) {
    // optionid_v2 = option_v2[Util.nextInt(option_v2.length)];
    // level10.itemOptions.add(new Item.ItemOption(optionid_v2, param));
    // }
    // level10.itemOptions.add(new Item.ItemOption(30, 0));
    // level10.itemOptions.add(new Item.ItemOption(93, Util.nextInt(level, 30)));
    // level10.quantity = 1;
    // player.itemsWoodChest.add(level10);
    // }
    // if (level > 10) {
    // byte[] option = {77, 80, 81, 103, 50, 94, 5};
    // byte[] option_v2 = {14, 16, 17, 19, 27, 28, 5, 47, 87}; //77 %hp // 80 //81
    // //103 //50 //94 //5 % sdcm
    // byte optionid;
    // byte optionid_v2;
    // byte param;
    // Item level11 = ItemService.gI().createNewItem((short) Util.nextInt(2112,
    // 2118));
    // level11.itemOptions.clear();
    // optionid = option[Util.nextInt(0, 6)];
    // param = (byte) Util.nextInt(5, 10);
    // level11.itemOptions.add(new Item.ItemOption(optionid, param));
    // if (Util.isTrue(20, 100)) {
    // optionid_v2 = option_v2[Util.nextInt(option_v2.length)];
    // level11.itemOptions.add(new Item.ItemOption(optionid_v2, param));
    // }
    // level11.itemOptions.add(new Item.ItemOption(30, 0));
    // if (Util.isTrue(90, 100)) {
    // level11.itemOptions.add(new Item.ItemOption(93, Util.nextInt(level, 30)));
    // }
    // level11.quantity = 1;
    // player.itemsWoodChest.add(level11);
    // }
    // if (level > 11) {
    // byte[] option = {77, 80, 81, 103, 50, 94, 5};
    // byte[] option_v2 = {14, 16, 17, 19, 27, 28, 5, 47, 87}; //77 %hp // 80 //81
    // //103 //50 //94 //5 % sdcm
    // byte optionid;
    // byte optionid_v2;
    // byte param;
    // Item level12 = ItemService.gI().createNewItem((short) Util.nextInt(2108,
    // 2122));
    // level12.itemOptions.clear();
    // optionid = option[Util.nextInt(0, 6)];
    // param = (byte) Util.nextInt(5, 10);
    // level12.itemOptions.add(new Item.ItemOption(optionid, param));
    // if (Util.isTrue(20, 100)) {
    // optionid_v2 = option_v2[Util.nextInt(option_v2.length)];
    // level12.itemOptions.add(new Item.ItemOption(optionid_v2, param));
    // }
    // level12.itemOptions.add(new Item.ItemOption(30, 0));
    // if (Util.isTrue(50, 100)) {
    // level12.itemOptions.add(new Item.ItemOption(93, Util.nextInt(level, 30)));
    // }
    // level12.quantity = 1;
    // player.itemsWoodChest.add(level12);
    // }
    // InventoryService.gI().subQuantityItemsBag(player, ruongGo, 1);
    // InventoryService.gI().sendItemBag(player);
    // for (Item it : player.itemsWoodChest) {
    // InventoryService.gI().addItemBag(player, it);
    // }
    // InventoryService.gI().sendItemBag(player);
    // player.indexWoodChest = player.itemsWoodChest.size() - 1;
    // int i = player.indexWoodChest;
    // if (i < 0) {
    // return;
    // }
    // Item itemWoodChest = player.itemsWoodChest.get(i);
    // player.indexWoodChest--;
    // String info = "|1|" + itemWoodChest.template.name;
    // String info2 = "\n|2|";
    // if (!itemWoodChest.itemOptions.isEmpty()) {
    // for (Item.ItemOption io : itemWoodChest.itemOptions) {
    // if (io.optionTemplate.id != 102 && io.optionTemplate.id != 73) {
    // info2 += io.getOptionString() + "\n";
    // }
    // }
    // }
    // info = (info2.length() > "\n|2|".length() ? (info + info2).trim() :
    // info.trim()) + "\n|0|" + itemWoodChest.template.description;
    // NpcService.gI().createMenuConMeo(player, ConstNpc.RUONG_GO, -1, "Bạn nhận
    // được\n"
    // + info.trim(), "OK" + (i > 0 ? " [" + i + "]" : ""));
    // }
    // }
    // }
    private void changePet(Player player, Item item) {
        if (player.pet != null) {
            int gender = player.pet.gender + 1;
            if (gender > 2) {
                gender = 0;
            }
            PetService.gI().changeNormalPet(player, gender);
            InventoryService.gI().subQuantityItemsBag(player, item, 1);
        } else {
            Service.gI().sendThongBao(player, "Không thể thực hiện");
        }
    }

    private void eatGrapes(Player pl, Item item) {
        int percentCurrentStatima = pl.nPoint.stamina * 100 / pl.nPoint.maxStamina;
        if (percentCurrentStatima > 50) {
            Service.gI().sendThongBao(pl, "Thể lực vẫn còn trên 50%");
            return;
        } else if (item.template.id == 211) {
            pl.nPoint.stamina = pl.nPoint.maxStamina;
            Service.gI().sendThongBao(pl, "Thể lực của bạn đã được hồi phục 100%");
        } else if (item.template.id == 212) {
            pl.nPoint.stamina += (pl.nPoint.maxStamina * 20 / 100);
            Service.gI().sendThongBao(pl, "Thể lực của bạn đã được hồi phục 20%");
        }
        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
        InventoryService.gI().sendItemBag(pl);
        PlayerService.gI().sendCurrentStamina(pl);
    }

    private void openCSKB(Player pl, Item item) {
        if (InventoryService.gI().getCountEmptyBag(pl) > 0) {
            short[] temp = { 76, 188, 189, 190, 381, 382, 383, 384, 385 };
            int[][] gold = { { 5000, 20000 } };
            byte index = (byte) Util.nextInt(0, temp.length - 1);
            short[] icon = new short[2];
            icon[0] = item.template.iconID;
            if (index <= 3) {
                pl.inventory.gold += Util.nextInt(gold[0][0], gold[0][1]);
                if (pl.inventory.gold > Inventory.LIMIT_GOLD) {
                    pl.inventory.gold = Inventory.LIMIT_GOLD;
                }
                PlayerService.gI().sendInfoHpMpMoney(pl);
                icon[1] = 930;
            } else {
                Item it = ItemService.gI().createNewItem(temp[index]);
                it.itemOptions.add(new ItemOption(73, 0));
                InventoryService.gI().addItemBag(pl, it);
                icon[1] = it.template.iconID;
            }
            InventoryService.gI().subQuantityItemsBag(pl, item, 1);
            InventoryService.gI().sendItemBag(pl);

            CombineService.gI().sendEffectOpenItem(pl, icon[0], icon[1]);
        } else {
            Service.gI().sendThongBao(pl, "Hàng trang đã đầy");
        }
    }

    private void useItemTime(Player pl, Item item) {
        switch (item.template.id) {
            case 1731:
                if (pl.itemTime.isUseLoX5 == true || pl.itemTime.isUseLoX7 == true || pl.itemTime.isUseLoX10 == true
                        || pl.itemTime.isUseLoX15 == true) {
                    Service.gI().sendThongBao(pl, "Bạn đang sử dụng nước thánh rồi");
                    return;
                }
                pl.itemTime.lastTimeLoX2 = System.currentTimeMillis();
                pl.itemTime.isUseLoX2 = true;
                break;
            case 1727:
                if (pl.itemTime.isUseLoX2 == true || pl.itemTime.isUseLoX7 == true || pl.itemTime.isUseLoX10 == true
                        || pl.itemTime.isUseLoX15 == true) {
                    Service.gI().sendThongBao(pl, "Bạn đang sử dụng nước thánh rồi");
                    return;
                }
                pl.itemTime.lastTimeLoX5 = System.currentTimeMillis();
                pl.itemTime.isUseLoX5 = true;
                break;
            case 1728:
                if (pl.itemTime.isUseLoX5 == true || pl.itemTime.isUseLoX2 == true || pl.itemTime.isUseLoX10 == true
                        || pl.itemTime.isUseLoX15 == true) {
                    Service.gI().sendThongBao(pl, "Bạn đang sử dụng nước thánh rồi");
                    return;
                }
                pl.itemTime.lastTimeLoX7 = System.currentTimeMillis();
                pl.itemTime.isUseLoX7 = true;
                break;
            case 1729:
                if (pl.itemTime.isUseLoX5 == true || pl.itemTime.isUseLoX7 == true || pl.itemTime.isUseLoX2 == true
                        || pl.itemTime.isUseLoX15 == true) {
                    Service.gI().sendThongBao(pl, "Bạn đang sử dụng nước thánh rồi");
                    return;
                }
                pl.itemTime.lastTimeLoX10 = System.currentTimeMillis();
                pl.itemTime.isUseLoX10 = true;
                break;
            case 1730:
                if (pl.itemTime.isUseLoX5 == true || pl.itemTime.isUseLoX7 == true || pl.itemTime.isUseLoX10 == true
                        || pl.itemTime.isUseLoX2 == true) {
                    Service.gI().sendThongBao(pl, "Bạn đang sử dụng nước thánh rồi");
                    return;
                }
                pl.itemTime.lastTimeLoX15 = System.currentTimeMillis();
                pl.itemTime.isUseLoX15 = true;
                break;
            case 764:
                pl.itemTime.lastTimeKhauTrang = System.currentTimeMillis();
                pl.itemTime.isUseKhauTrang = true;
                break;
            case 1628:
                if (!addTimeBuax2DeTu(pl)) {
                    return;
                }
                break;
            case 382: // bổ huyết
                if (pl.itemTime != null && pl.itemTime.isUseBoHuyet2) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeBoHuyet = System.currentTimeMillis();
                pl.itemTime.isUseBoHuyet = true;
                break;
            case 383: // bổ khí
                if (pl.itemTime != null && pl.itemTime.isUseBoKhi2) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeBoKhi = System.currentTimeMillis();
                pl.itemTime.isUseBoKhi = true;
                break;
            case 384: // giáp xên
                if (pl.itemTime != null && pl.itemTime.isUseGiapXen2) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeGiapXen = System.currentTimeMillis();
                pl.itemTime.isUseGiapXen = true;
                break;
            case 381: // cuồng nộ
                if (pl.itemTime != null && pl.itemTime.isUseCuongNo2) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeCuongNo = System.currentTimeMillis();
                pl.itemTime.isUseCuongNo = true;
                Service.gI().point(pl);
                break;
            case 385: // ẩn danh
                if (pl.itemTime != null && pl.itemTime.isUseAnDanh2) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeAnDanh = System.currentTimeMillis();
                pl.itemTime.isUseAnDanh = true;
                break;
            case 379: // máy dò capsule
                pl.itemTime.lastTimeUseMayDo = System.currentTimeMillis();
                pl.itemTime.isUseMayDo = true;
                break;
            case 1150:// cn
                if (pl.itemTime != null && pl.itemTime.isUseCuongNo) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeCuongNo2 = System.currentTimeMillis();
                pl.itemTime.isUseCuongNo2 = true;
                Service.gI().point(pl);
                break;
            case 752://
                pl.itemTime.lastTimeBanhtet = System.currentTimeMillis();
                pl.itemTime.isBanhtet = true;
                Service.gI().point(pl);
                break;
            case 753://
                pl.itemTime.lastTimeBanhchung = System.currentTimeMillis();
                pl.itemTime.isBanhchung = true;
                Service.gI().point(pl);
                break;
            case 1195://
                pl.itemTime.lastTimeXimuoihoadao = System.currentTimeMillis();
                pl.itemTime.isXimuoihoadao = true;
                Service.gI().point(pl);
                break;
            case 1196:// cn 2
                pl.itemTime.lastTimeXimuoihoamai = System.currentTimeMillis();
                pl.itemTime.isXimuoihoamai = true;
                Service.gI().point(pl);
                break;
            case 1152:// bo huyet 2
                if (pl.itemTime != null && pl.itemTime.isUseBoHuyet) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeBoHuyet2 = System.currentTimeMillis();
                pl.itemTime.isUseBoHuyet2 = true;
                break;
            case 1151:// bo khi
                if (pl.itemTime != null && pl.itemTime.isUseBoKhi) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeBoKhi2 = System.currentTimeMillis();
                pl.itemTime.isUseBoKhi2 = true;
                break;
            case 1153:// gx 2
                if (pl.itemTime != null && pl.itemTime.isUseGiapXen) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeGiapXen2 = System.currentTimeMillis();
                pl.itemTime.isUseGiapXen2 = true;
                break;
            case 1154:// an danh
                if (pl.itemTime != null && pl.itemTime.isUseAnDanh) {
                    Service.gI().sendThongBao(pl, "Tham lam!");
                    return;
                }
                pl.itemTime.lastTimeAnDanh2 = System.currentTimeMillis();
                pl.itemTime.isUseAnDanh2 = true;
                break;
            case 638: // Commeson
                pl.itemTime.lastTimeUseCMS = System.currentTimeMillis();
                pl.itemTime.isUseCMS = true;
                break;
            case 2160: // Nồi cơm điện
                pl.itemTime.lastTimeUseNCD = System.currentTimeMillis();
                pl.itemTime.isUseNCD = true;
                break;
            case 579:
                addTimeUseDK(pl);
                break;
            case 1045: // Đuôi khỉ
                pl.itemTime.lastTimeUseDK = System.currentTimeMillis();
                pl.itemTime.isUseDK = true;
                break;
            case 663: // bánh pudding
            case 664: // xúc xíc
            case 665: // kem dâu
            case 666: // mì ly
            case 667: // sushi
                pl.itemTime.lastTimeEatMeal = System.currentTimeMillis();
                pl.itemTime.isEatMeal = true;
                ItemTimeService.gI().removeItemTime(pl, pl.itemTime.iconMeal);
                pl.itemTime.iconMeal = item.template.iconID;
                break;
            case 880:
            case 881:
            case 882:
                pl.itemTime.lastTimeEatMeal2 = System.currentTimeMillis();
                pl.itemTime.isEatMeal2 = true;
                ItemTimeService.gI().removeItemTime(pl, pl.itemTime.iconMeal2);
                pl.itemTime.iconMeal2 = item.template.iconID;
                break;
            case 889:
            case 900:
            case 902:
            case 903:
                pl.itemTime.lastTimeEatMeal3 = System.currentTimeMillis();
                pl.itemTime.isEatMeal3 = true;
                ItemTimeService.gI().removeItemTime(pl, pl.itemTime.iconMeal3);
                pl.itemTime.iconMeal3 = item.template.iconID;
                break;
            case 1109: // máy dò đồ
                pl.itemTime.lastTimeUseMayDo2 = System.currentTimeMillis();
                pl.itemTime.isUseMayDo2 = true;
                break;
        }
        Service.gI().point(pl);
        ItemTimeService.gI().sendAllItemTime(pl);
        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
        InventoryService.gI().sendItemBag(pl);
    }

    private void addTimeUseDK(Player pl) {
        long now = System.currentTimeMillis();
        long timeLeft = pl.itemTime.isUseDK ? ItemTime.TIME_DK - (now - pl.itemTime.lastTimeUseDK) : 0;
        long newTimeLeft = Math.max(0, timeLeft) + ItemTime.TIME_DK;
        pl.itemTime.lastTimeUseDK = now - (ItemTime.TIME_DK - newTimeLeft);
        pl.itemTime.isUseDK = true;
    }

    private boolean addTimeBuax2DeTu(Player pl) {
        if (pl.pet == null) {
            Service.gI().sendThongBao(pl, "Có Đệ tử đéo đâu");
            return false;
        }
        long now = System.currentTimeMillis();
        long timeLeft = pl.itemTime.isUseBuax2DeTu ? ItemTime.TIME_30P - (now - pl.itemTime.lastTimeBuax2DeTu) : 0;
        long newTimeLeft = Math.max(0, timeLeft) + ItemTime.TIME_30P;
        pl.itemTime.lastTimeBuax2DeTu = now - (ItemTime.TIME_30P - newTimeLeft);
        pl.itemTime.isUseBuax2DeTu = true;
        Service.gI().sendThongBao(pl, "Đã kích hoạt X2 TNSM Đệ tử");
        return true;
    }

    private void controllerCallRongThan(Player pl, Item item) {
        int tempId = item.template.id;
        if (tempId >= SummonDragon.NGOC_RONG_1_SAO && tempId <= SummonDragon.NGOC_RONG_7_SAO) {
            switch (tempId) {
                case SummonDragon.NGOC_RONG_1_SAO:
                    // case SummonDragon.NGOC_RONG_2_SAO:
                    // case SummonDragon.NGOC_RONG_3_SAO:
                    SummonDragon.gI().openMenuSummonShenron(pl, (byte) (tempId - 13));
                    break;
                default:
                    NpcService.gI().createMenuConMeo(pl, ConstNpc.TUTORIAL_SUMMON_DRAGON,
                            -1, "Bạn chỉ có thể gọi rồng từ ngọc 1 sao", "Hướng\ndẫn thêm\n(mới)", "OK");
                    break;
            }
        } else if (tempId >= ShenronEventService.NGOC_RONG_1_SAO && tempId <= ShenronEventService.NGOC_RONG_7_SAO) {
            ShenronEventService.gI().openMenuSummonShenron(pl, 0);
        }
    }

    private void learnSkill(Player pl, Item item) {
        Message msg;
        try {
            if (item.template.id >= 1334 && item.template.id <= 1351) {
                learnSkillSuperNew(pl, item);
            } else if (item.template.id >= 1356 && item.template.id <= 1376) {
                if (item.template.id == 1356 || item.template.id == 1363 || item.template.id == 1370) {
                    if (item.template.id >= 1363 && item.template.id <= 1369) {
                        SkillService.gI().learSkillSpecial(pl, Skill.PHAN_THAN);
                    }
                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                    return;
                }
                learnSkillNew2(pl, item);
            } else {
                if (item.template.gender == pl.gender || item.template.gender == 3) {
                    String[] subName = item.template.name.split("");
                    byte level = Byte.parseByte(subName[subName.length - 1]);
                    Skill curSkill = SkillUtil.getSkillByItemID(pl, item.template.id);
                    if (curSkill.point == 7) {
                        Service.gI().sendThongBao(pl, "Kỹ năng đã đạt tối đa!");
                    } else {
                        if (curSkill.point == 0) {
                            if (level == 1) {
                                curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id),
                                        level);
                                SkillUtil.setSkill(pl, curSkill);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                msg = Service.gI().messageSubCommand((byte) 23);
                                msg.writer().writeShort(curSkill.skillId);
                                pl.sendMessage(msg);
                                msg.cleanup();
                            } else {
                                Skill skillNeed = SkillUtil
                                        .createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id), level);
                                Service.gI().sendThongBao(pl, "Vui lòng học " + skillNeed.template.name + " cấp "
                                        + skillNeed.point + " trước!");
                            }
                        } else {
                            if (curSkill.point + 1 == level) {
                                curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id),
                                        level);
                                // System.out.println(curSkill.template.name + " - " + curSkill.point);
                                SkillUtil.setSkill(pl, curSkill);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                msg = Service.gI().messageSubCommand((byte) 62);
                                msg.writer().writeShort(curSkill.skillId);
                                pl.sendMessage(msg);
                                msg.cleanup();
                            } else {
                                Service.gI().sendThongBao(pl, "Vui lòng học " + curSkill.template.name + " cấp "
                                        + (curSkill.point + 1) + " trước!");
                            }
                        }
                        InventoryService.gI().sendItemBag(pl);
                    }
                } else {
                    Service.gI().sendThongBao(pl, "Không thể thực hiện");
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void learnSkillNew2(Player pl, Item item) {
        Message msg;
        try {
            if (item.template.gender == pl.gender || item.template.gender == 3) {
                byte level = SkillUtil.getLevelSkillByItemID(item.template.id);
                Skill curSkill = SkillUtil.getSkillByItemID(pl, item.template.id);
                if (curSkill == null) {
                    SkillService.gI().learSkillSpecial(pl,
                            (byte) SkillUtil.getSkillByItemID(pl, item.template.id).skillId);
                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                    return;
                } else {
                    if (curSkill.point == 7) {
                        Service.gI().sendThongBao(pl, "Kỹ năng đã đạt tối đa!");
                    } else {
                        if (curSkill.point == 0) {
                            if (level == 1) {
                                curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id),
                                        level);
                                SkillUtil.setSkill(pl, curSkill);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                msg = Service.gI().messageSubCommand((byte) 23);
                                msg.writer().writeShort(curSkill.skillId);
                                pl.sendMessage(msg);
                                msg.cleanup();
                                if (curSkill.template.id == Skill.SUPER_NAMEC
                                        || curSkill.template.id == Skill.SUPER_SAIYAN
                                        || curSkill.template.id == Skill.SUPER_TRAI_DAT) {
                                    curSkill = SkillUtil.createSkill(Skill.GONG, level);
                                    SkillUtil.setSkill(pl, curSkill);
                                    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                    msg = Service.gI().messageSubCommand((byte) 23);
                                    msg.writer().writeShort(curSkill.skillId);
                                    pl.sendMessage(msg);
                                    msg.cleanup();
                                }
                            } else {
                                Skill skillNeed = SkillUtil
                                        .createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id), level);
                                Service.gI().sendThongBao(pl, "Vui lòng học " + skillNeed.template.name + " cấp "
                                        + skillNeed.point + " trước!");
                            }
                        } else {
                            if (curSkill.point + 1 == level) {
                                curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id),
                                        level);
                                // System.out.println(curSkill.template.name + " - " + curSkill.point);
                                SkillUtil.setSkill(pl, curSkill);
                                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                                msg = Service.gI().messageSubCommand((byte) 62);
                                msg.writer().writeShort(curSkill.skillId);
                                pl.sendMessage(msg);
                                msg.cleanup();
                            } else {
                                Service.gI().sendThongBao(pl, "Vui lòng học " + curSkill.template.name + " cấp "
                                        + (curSkill.point + 1) + " trước!");
                            }
                        }
                        InventoryService.gI().sendItemBag(pl);
                    }
                }
            } else {
                Service.gI().sendThongBao(pl, "Không thể thực hiện");
            }
        } catch (Exception e) {

        }
    }

    private void learnSkillSuperNew(Player pl, Item item) {
        Message msg;
        try {
            if (item.template.gender == pl.gender || item.template.gender == 3) {
                byte level = SkillUtil.getLevelSkillByItemID(item.template.id);
                Skill curSkill = SkillUtil.getSkillByItemID(pl, item.template.id);
                if (curSkill.point == 6) {
                    Service.gI().sendThongBao(pl, "Kỹ năng đã đạt tối đa!");
                } else {
                    if (curSkill.point == 0) {
                        if (level == 1) {
                            curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id),
                                    level);
                            SkillUtil.setSkill(pl, curSkill);
                            InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                            msg = Service.gI().messageSubCommand((byte) 23);
                            msg.writer().writeShort(curSkill.skillId);
                            pl.sendMessage(msg);
                            msg.cleanup();
                            SkillService.gI().learSkillSpecial(pl, (byte) 30);
                        } else {
                            Skill skillNeed = SkillUtil
                                    .createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id), level);
                            if (level > 1) {
                                Item itemNew = ItemService.gI().createNewItem((short) (item.template.id - 1));
                                String name = itemNew.template.name;
                                String desiredName = name.substring(5);
                                Service.gI().sendThongBao(pl, "Vui lòng học " + desiredName + " trước!");
                            } else {
                                Service.gI().sendThongBao(pl, "Vui lòng học " + skillNeed.template.name + " trước!");
                            }
                        }
                    } else {
                        if (curSkill.point + 1 == level) {
                            curSkill = SkillUtil.createSkill(SkillUtil.getTempSkillSkillByItemID(item.template.id),
                                    level);
                            SkillUtil.setSkill(pl, curSkill);
                            InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                            msg = Service.gI().messageSubCommand((byte) 62);
                            msg.writer().writeShort(curSkill.skillId);
                            pl.sendMessage(msg);
                            msg.cleanup();
                        } else {
                            if (level > 1) {
                                Item itemNew = ItemService.gI().createNewItem((short) (item.template.id - 1));
                                String name = itemNew.template.name;
                                String desiredName = name.substring(5);
                                Service.gI().sendThongBao(pl, "Vui lòng học " + desiredName + " trước!");
                            } else {
                                Service.gI().sendThongBao(pl, "Vui lòng học " + curSkill.template.name + " trước!");
                            }
                        }
                    }
                    InventoryService.gI().sendItemBag(pl);
                }
            } else {
                Service.gI().sendThongBao(pl, "Không thể thực hiện");
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void useTDLT(Player pl, Item item) {
        if (pl.itemTime.isUseTDLT) {
            ItemTimeService.gI().turnOffTDLT(pl, item);
        } else {
            ItemTimeService.gI().turnOnTDLT(pl, item);
        }
    }

    // private void usePorataGogeta(Player pl) {
    // if (pl.pet == null || pl.fusion.typeFusion == 4) {
    // Service.gI().sendThongBao(pl, "Không thể thực hiện");
    // } else {
    // if (pl.fusion.typeFusion == ConstPlayer.NON_FUSION) {
    // pl.pet.fusionGogeta(true);
    // } else {
    // pl.pet.unFusion();
    // }
    // }
    // }
    private void usePorata2(Player pl) {
        if (pl.pet == null || pl.fusion.typeFusion == 4) {
            Service.gI().sendThongBao(pl, "Không thể thực hiện");
        } else {
            if (pl.fusion.typeFusion == ConstPlayer.NON_FUSION) {
                pl.pet.fusion2(true);
            } else {
                pl.pet.unFusion();
            }
        }
    }

    private void usePorata(Player pl) {
        if (pl.pet == null || pl.fusion.typeFusion == 4) {
            Service.gI().sendThongBao(pl, "Không thể thực hiện");
        } else {
            if (pl.fusion.typeFusion == ConstPlayer.NON_FUSION) {
                pl.pet.fusion(true);
            } else {
                pl.pet.unFusion();
            }
        }
    }

    private void usePorata(Player pl, byte typeFusion) {
        if (pl.pet == null || pl.fusion.typeFusion == 4) {
            Service.gI().sendThongBao(pl, "Khong the thuc hien");
        } else {
            if (pl.fusion.typeFusion == ConstPlayer.NON_FUSION) {
                pl.pet.fusion2(true);
                if (pl.fusion.typeFusion != ConstPlayer.NON_FUSION) {
                    pl.fusion.typeFusion = typeFusion;
                    pl.nPoint.calPoint();
                    Service.gI().point(pl);
                    Service.gI().Send_Caitrang(pl);
                }
            } else {
                pl.pet.unFusion();
            }
        }
    }

    private void openCapsuleUI(Player pl) {
        pl.iDMark.setTypeChangeMap(ConstMap.CHANGE_CAPSULE);
        ChangeMapService.gI().openChangeMapTab(pl);
    }

    private void openxukrai(Player pl, Item item, int soluong) {
        if (pl == null || item == null) {
            System.out.println("Lỗi: Người chơi hoặc vật phẩm không hợp lệ." + pl + "");
            return;
        }
        if (soluong <= 0) {
            System.out.println("Lỗi: Số lượng xu ocean không hợp lệ.");
            return;
        }
        if (InventoryService.gI().getCountEmptyBag(pl) < 1) {
            Service.gI().sendThongBao(pl, "Cần 1 ô hành trang mới có thể đổi!!!");
            return;
        }
        Item xukrai = ItemService.gI().createNewItem((short) 1705, soluong);
        xukrai.itemOptions.add(new Item.ItemOption(30, 1));
        InventoryService.gI().addItemBag(pl, xukrai);
        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
        Service.gI().sendThongBao(pl, "Bạn đã nhận được " + soluong + " xu ocean");

    }

    public void openRuongNgocRong(Player pl, Item item) {
        int nr = Util.nextInt(14, 20);
        int[] vp = { 220, 221, 222, 223, 224, 225 };
        int[] vpVip = { 1682, 1683, 1688, 1689, 1677, 1678, 1013, 1021, 1022 };
        Item item2 = null;
        if (Util.isTrue(80, 100)) {
            item2 = ItemService.gI().createNewItem((short) nr);
            item2.quantity = 1;
            item2.itemOptions.add(new Item.ItemOption(30, 1));
        } else if (Util.isTrue(80, 100)) {
            item2 = ItemService.gI().createNewItem((short) vp[Util.nextInt(0, vp.length - 1)]);
            item2.quantity = 1;
            item2.itemOptions.add(new Item.ItemOption(30, 1));
        } else {
            item2 = ItemService.gI().createNewItem((short) vpVip[Util.nextInt(0, vpVip.length - 1)]);
            item2.quantity = 1;
            item2.itemOptions.add(new Item.ItemOption(93, 30));
            item2.itemOptions.add(new Item.ItemOption(50, 10));
            item2.itemOptions.add(new Item.ItemOption(77, 10));
            item2.itemOptions.add(new Item.ItemOption(103, 10));
            item2.itemOptions.add(new Item.ItemOption(101, 10));

        }
        InventoryService.gI().addItemBag(pl, item2);
        InventoryService.gI().sendItemBag(pl);
        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
        Service.gI().sendThongBao(pl, "Bạn mở rương nhận được " + item2.template.name);
    }

 public void ruongnro(Player pl, Item item) {
    // Danh sách item có thể ra
    short[] list = {16, 18, 19, 20};

    // Random 1 item
    short nr = list[Util.nextInt(list.length)];

    // Tạo item
    Item item2 = ItemService.gI().createNewItem(nr);
    item2.quantity = 2; // x2 số lượng
    // ❌ KHÔNG thêm option 30

    // Trừ rương trước cho an toàn
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);

    // Thêm item vào hành trang
    InventoryService.gI().addItemBag(pl, item2);
    InventoryService.gI().sendItemBag(pl);

    Service.gI().sendThongBao(
        pl,
        "Bạn mở rương nhận được x2 " + item2.template.name
    );
}


    public void choseMapCapsule(Player pl, int index) {

        if (pl.idNRNM != -1) {
            Service.gI().sendThongBao(pl, "Không thể mang ngọc rồng này lên Phi thuyền");
            Service.gI().hideWaitDialog(pl);
            return;
        }

        int zoneId = -1;
        if (index > pl.mapCapsule.size() - 1 || index < 0) {
            Service.gI().sendThongBao(pl, "Không thể thực hiện");
            Service.gI().hideWaitDialog(pl);
            return;
        }
        Zone zoneChose = pl.mapCapsule.get(index);
        // Kiểm tra số lượng người trong khu

        if (zoneChose.getNumOfPlayers() > 25
                || MapService.gI().isMapDoanhTrai(zoneChose.map.mapId)
                || MapService.gI().isMapMaBu(zoneChose.map.mapId)
                || MapService.gI().isMapHuyDiet(zoneChose.map.mapId)) {
            Service.gI().sendThongBao(pl, "Hiện tại không thể vào được khu!");
            return;
        }
        if (index != 0 || zoneChose.map.mapId == 21
                || zoneChose.map.mapId == 22
                || zoneChose.map.mapId == 23) {
            pl.mapBeforeCapsule = pl.zone;
        } else {
            zoneId = pl.mapBeforeCapsule != null ? pl.mapBeforeCapsule.zoneId : -1;
            pl.mapBeforeCapsule = null;
        }
        pl.changeMapVIP = true;
        ChangeMapService.gI().changeMapBySpaceShip(pl, pl.mapCapsule.get(index).map.mapId, zoneId, -1);
    }

    public void eatPea(Player player) {
        if (!Util.canDoWithTime(player.lastTimeEatPea, 1000)) {
            return;
        }
        player.lastTimeEatPea = System.currentTimeMillis();
        Item pea = null;
        for (Item item : player.inventory.itemsBag) {
            if (item.isNotNullItem() && item.template.type == 6) {
                pea = item;
                break;
            }
        }
        if (pea != null) {
            long hpKiHoiPhuc = 0;
            int lvPea = Integer.parseInt(pea.template.name.substring(13));
            for (Item.ItemOption io : pea.itemOptions) {
                if (io.optionTemplate.id == 2) {
                    hpKiHoiPhuc = io.param * 1000;
                    break;
                }
                if (io.optionTemplate.id == 48) {
                    hpKiHoiPhuc = io.param;
                    break;
                }
            }
            player.nPoint.setHp(Util.maxIntValue(player.nPoint.hp + hpKiHoiPhuc));
            player.nPoint.setMp(Util.maxIntValue(player.nPoint.mp + hpKiHoiPhuc));
            PlayerService.gI().sendInfoHpMp(player);
            Service.gI().sendInfoPlayerEatPea(player);
            if (player.pet != null && player.zone.equals(player.pet.zone) && !player.pet.isDie()) {
                int statima = 100 * lvPea;
                player.pet.nPoint.stamina += statima;
                if (player.pet.nPoint.stamina > player.pet.nPoint.maxStamina) {
                    player.pet.nPoint.stamina = player.pet.nPoint.maxStamina;
                }
                player.pet.nPoint.setHp(Util.maxIntValue(player.pet.nPoint.hp + hpKiHoiPhuc));
                player.pet.nPoint.setMp(Util.maxIntValue(player.pet.nPoint.mp + hpKiHoiPhuc));
                Service.gI().sendInfoPlayerEatPea(player.pet);
                Service.gI().chatJustForMe(player, player.pet, "Cám ơn sư phụ");
            }

            InventoryService.gI().subQuantityItemsBag(player, pea, 1);
            InventoryService.gI().sendItemBag(player);
        }
    }

    private void upSkillPet(Player pl, Item item) {
        if (pl.pet == null) {
            Service.gI().sendThongBao(pl, "Không thể thực hiện");
            return;
        }
        try {
            switch (item.template.id) {
                case 402: // skill 1
                    if (SkillUtil.upSkillPet(pl.pet.playerSkill.skills, 0)) {
                        Service.gI().chatJustForMe(pl, pl.pet, "Cám ơn sư phụ");
                        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                    } else {
                        Service.gI().sendThongBao(pl, "Không thể thực hiện");
                    }
                    break;
                case 403: // skill 2
                    if (SkillUtil.upSkillPet(pl.pet.playerSkill.skills, 1)) {
                        Service.gI().chatJustForMe(pl, pl.pet, "Cám ơn sư phụ");
                        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                    } else {
                        Service.gI().sendThongBao(pl, "Không thể thực hiện");
                    }
                    break;
                case 404: // skill 3
                    if (SkillUtil.upSkillPet(pl.pet.playerSkill.skills, 2)) {
                        Service.gI().chatJustForMe(pl, pl.pet, "Cám ơn sư phụ");
                        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                    } else {
                        Service.gI().sendThongBao(pl, "Không thể thực hiện");
                    }
                    break;
                case 759: // skill 4
                    if (SkillUtil.upSkillPet(pl.pet.playerSkill.skills, 3)) {
                        Service.gI().chatJustForMe(pl, pl.pet, "Cám ơn sư phụ");
                        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                    } else {
                        Service.gI().sendThongBao(pl, "Không thể thực hiện");
                    }
                    break;

            }

        } catch (Exception e) {
            Service.gI().sendThongBao(pl, "Không thể thực hiện");
        }
    }

    // private void ItemManhGiay(Player pl, Item item) {
    //     if (pl.winSTT && !Util.isAfterMidnight(pl.lastTimeWinSTT)) {
    //         Service.gI().sendThongBao(pl, "Hãy gặp thần mèo Karin để sử dụng");
    //         return;
    //     } else if (pl.winSTT && Util.isAfterMidnight(pl.lastTimeWinSTT)) {
    //         pl.winSTT = false;
    //         pl.callBossPocolo = false;
    //         pl.zoneSieuThanhThuy = null;
    //     }
    //     NpcService.gI().createMenuConMeo(pl, item.template.id, 564,
    //             "Đây chính là dấu hiệu riêng của...\nĐại Ma Vương Pôcôlô\nĐó là một tên quỷ dữ đội lốt người, một kẻ đại gian ác\ncó sức mạnh vô địch và lòng tham không đáy...\nĐối phó với hắn không phải dễ\nCon có chắc chắn muốn tìm hắn không?",
    //             "Đồng ý", "Từ chối");
    // }


private void ItemManhGiay(Player pl, Item item) {
    if (item == null || item.template.id != 726) {
        return;
    }

    if (pl.zone == null || pl.isDie()) {
        return;
    }

    Service.gI().sendThongBao(pl, "Đang dịch chuyển...");

    InventoryService.gI().subQuantityItemsBag(pl, item, 1);

    ChangeMapService.gI().changeMap(pl, 135, -1, -1, 0);
}



    private void ItemSieuThanThuy(Player pl, Item item) {
        long tnsm = 5_000_000;
        int n = 0;
        switch (item.template.id) {
            case 727:
                n = 2;
                break;
            case 728:
                n = 2;
                break;
        }
        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
        InventoryService.gI().sendItemBag(pl);
        if (Util.isTrue(50, 100)) {
            Service.gI().sendThongBao(pl, "Bạn đã bị chết vì độc của thuốc tăng lực siêu thần thủy.");
            pl.setDie();
        } else {
            for (int i = 0; i < n; i++) {
                Service.gI().addSMTN(pl, (byte) 2, tnsm, true);
            }
        }
    }

    private void Hopdothanlinh(Player pl, Item item) {// hop qua do thần linh
        NpcService.gI().createMenuConMeo(pl, item.template.id, -1, "Chọn hành tinh của Bạn đi", "Set trái đất",
                "Set namec", "Set xayda", "Từ chổi");
    }

    private void Hopdohuydiet(Player pl, Item item) {// hop qua do huy diet
        NpcService.gI().createMenuConMeo(pl, item.template.id, -1, "Chọn hành tinh của Bạn đi", "Set trái đất",
                "Set namec", "Set xayda", "Từ chổi");
    }

    public void UseCard(Player pl, Item item) {
        RadarCard radarTemplate = RadarService.gI().RADAR_TEMPLATE.stream().filter(c -> c.Id == item.template.id)
                .findFirst().orElse(null);
        if (radarTemplate == null) {
            return;
        }
        if (radarTemplate.Require != -1) {
            RadarCard radarRequireTemplate = RadarService.gI().RADAR_TEMPLATE.stream()
                    .filter(r -> r.Id == radarTemplate.Require).findFirst().orElse(null);
            if (radarRequireTemplate == null) {
                return;
            }
            Card cardRequire = pl.Cards.stream().filter(r -> r.Id == radarRequireTemplate.Id).findFirst().orElse(null);
            if (cardRequire == null || cardRequire.Level < radarTemplate.RequireLevel) {
                Service.gI().sendThongBao(pl, "Bạn cần sưu tầm " + radarRequireTemplate.Name + " ở cấp độ "
                        + radarTemplate.RequireLevel + " mới có thể sử dụng thẻ này");
                return;
            }
        }
        Card card = pl.Cards.stream().filter(r -> r.Id == item.template.id).findFirst().orElse(null);
        if (card == null) {
            Card newCard = new Card(item.template.id, (byte) 1, radarTemplate.Max, (byte) -1, radarTemplate.Options);
            if (pl.Cards.add(newCard)) {
                RadarService.gI().RadarSetAmount(pl, newCard.Id, newCard.Amount, newCard.MaxAmount);
                RadarService.gI().RadarSetLevel(pl, newCard.Id, newCard.Level);
                InventoryService.gI().subQuantityItemsBag(pl, item, 1);
                InventoryService.gI().sendItemBag(pl);
            }
        } else {
            if (card.Level >= 2) {
                Service.gI().sendThongBao(pl, "Thẻ này đã đạt cấp tối đa");
                return;
            }
            card.Amount++;
            if (card.Amount >= card.MaxAmount) {
                card.Amount = 0;
                if (card.Level == -1) {
                    card.Level = 1;
                } else {
                    card.Level++;
                }
                Service.gI().point(pl);
            }
            RadarService.gI().RadarSetAmount(pl, card.Id, card.Amount, card.MaxAmount);
            RadarService.gI().RadarSetLevel(pl, card.Id, card.Level);
            InventoryService.gI().subQuantityItemsBag(pl, item, 1);
            InventoryService.gI().sendItemBag(pl);
        }
    }

    private void replaceSkill2(Player pl, Item item, int skillId, String skillName) {
        if (pl.pet != null) {
            if (pl.pet.playerSkill.skills.get(1).skillId != -1) {
                Skill skill = SkillUtil.createSkill(skillId, 1);
                skill.coolDown = 1000;
                pl.pet.playerSkill.skills.set(1, skill);
                Service.gI().sendThongBao(pl, "Thay thành công chiêu 2 đệ tử thành " + skillName + "!");
            } else {
                Service.gI().sendThongBao(pl, "Ít nhất đệ tử ngươi phải có chiêu 2 chứ!");
                return;
            }
        } else {
            Service.gI().sendThongBao(pl, "Ngươi làm gì có đệ tử?");
            return;
        }
        InventoryService.gI().subQuantityItemsBag(pl, item, 1);
    }

   private void replaceSkill3(Player pl, Item item, int skillId, String skillName) {
    if (pl.pet == null) {
        Service.gI().sendThongBao(pl, "Ngươi làm gì có đệ tử?");
        return;
    }

    if (pl.pet.playerSkill.skills.get(2).skillId == -1) {
        Service.gI().sendThongBao(pl, "Ít nhất đệ tử ngươi phải có chiêu 3 chứ!");
        return;
    }

    Skill skill = SkillUtil.createSkill(skillId, 1);
    skill.coolDown = 1000;
    pl.pet.playerSkill.skills.set(2, skill);

    Service.gI().sendThongBao(pl, "Thay thành công chiêu 3 đệ tử thành " + skillName + "!");
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
}
private void replaceSkill4(Player pl, Item item, int skillId, String skillName) {
    if (pl.pet != null) {
        if (pl.pet.playerSkill.skills.get(3).skillId != -1) { // ✅ sửa get(1) → get(3)
            Skill skill = SkillUtil.createSkill(skillId, 1);
            pl.pet.playerSkill.skills.set(3, skill);
            Service.gI().sendThongBao(pl, "Thay thành công chiêu 4 đệ tử thành " + skillName + "!");
        } else {
            Service.gI().sendThongBao(pl, "Ít nhất đệ tử ngươi phải có chiêu 4 chứ!");
            return;
        }
    } else {
        Service.gI().sendThongBao(pl, "Ngươi làm gì có đệ tử?");
        return;
    }
    InventoryService.gI().subQuantityItemsBag(pl, item, 1);
}
}
