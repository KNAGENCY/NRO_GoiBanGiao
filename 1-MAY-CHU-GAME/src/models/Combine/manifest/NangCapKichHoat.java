package models.Combine.manifest;

import consts.ConstNpc;
import item.Item;
import models.Combine.CombineService;
import player.Player;
import server.Manager;
import server.ServerNotify;
import services.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class NangCapKichHoat {

    // ================= LIST ITEM =================

    public static final int[] AO_TD = {0,3,33,34,136,137,138,139,230,231,232,233,555};
    public static final int[] QUAN_TD = {6,9,35,36,140,141,142,143,242,243,244,245,556};
    public static final int[] GANG_TD = {21,24,37,38,144,145,146,147,254,255,256,257,562};
    public static final int[] GIAY_TD = {27,30,39,40,148,149,150,151,266,267,268,269,563};
    public static final int[] RADA_TD = {12,57,58,59,184,185,186,187,278,279,280,281,561};

    public static final int[] AO_NAMEK = {1,4,41,42,152,153,154,155,234,235,236,237,557};
    public static final int[] QUAN_NAMEK = {7,10,43,44,156,157,158,159,246,247,248,249,558};
    public static final int[] GANG_NAMEK = {22,25,46,45,160,161,162,163,258,259,260,261,564};
    public static final int[] GIAY_NAMEK = {28,31,47,48,164,165,166,167,270,271,272,273,565};
    public static final int[] RADA_NAMEK = {12,57,58,59,184,185,186,187,278,279,280,281,561};

    public static final int[] AO_XAYDA = {2,5,49,50,168,169,170,171,238,239,240,241,559};
    public static final int[] QUAN_XAYDA = {8,11,51,52,172,173,174,175,250,251,252,253,560};
    public static final int[] GANG_XAYDA = {23,26,53,54,176,177,178,179,262,263,264,265,566};
    public static final int[] GIAY_XAYDA = {29,32,55,56,180,181,182,183,274,275,276,277,567};
    public static final int[] RADA_XAYDA = {12,57,58,59,184,185,186,187,278,279,280,281,561};

    // ================= CHECK =================

    public static boolean isDoHuyDiet(Item item){
        return item.template.id >= 650 && item.template.id <= 662;
    }

    public static boolean isDoThan(int id){
        return id >= 555 && id <= 567;
    }

    // ================= RANDOM PARAM =================

    private static int randomParam(int base){
        int min = (int)(base * 0.9);
        int max = (int)(base * 1.2);
        return Util.nextInt(min, max);
    }

    // ================= ADD OPTION =================

    private static void add(Item item,int id,int base){
        item.itemOptions.add(new Item.ItemOption(id,randomParam(base)));
    }

    private static void add2(Item item,int id1,int base1,int id2,int base2){
        add(item,id1,base1);
        add(item,id2,base2);
    }

    // ================= OPTION MẶC ĐỊNH =================

    public static void setOptionMacDinh(Item item){

        switch(item.template.id){

            // ===== ÁO =====

            case 0:case 1:case 2: add(item,47,2); break;
            case 33:case 41:case 49: add(item,47,4); break;
            case 3:case 4:case 5: add(item,47,8); break;
            case 34:case 42:case 50: add(item,47,16); break;
            case 136:case 152:case 168: add(item,47,24); break;
            case 137:case 153:case 169: add(item,47,40); break;
            case 138:case 154:case 170: add(item,47,60); break;
            case 139:case 155:case 171: add(item,47,90); break;
            case 230:case 234:case 238: add(item,47,200); break;
            case 231:case 235:case 239: add(item,47,250); break;
            case 232:case 236:case 240: add(item,47,300); break;
            case 233:case 237:case 241: add(item,47,400); break;
            case 555:case 557:case 559: add(item,47,800); break;

            // ===== QUẦN =====

            case 6: add(item,6,30); break;
            case 7:case 8: add(item,6,20); break;

            case 35:case 43:case 51: add2(item,6,150,27,12); break;
            case 9:case 10:case 11: add2(item,6,300,27,40); break;
            case 36:case 44:case 52: add2(item,6,600,27,120); break;

            case 140:case 156:case 172: add2(item,6,1400,27,280); break;
            case 141:case 157:case 173: add2(item,6,3000,27,600); break;
            case 142:case 158:case 174: add2(item,6,6000,27,1200); break;
            case 143:case 159:case 175: add2(item,6,10000,27,2000); break;

            case 242:case 246:case 250: add2(item,6,14000,27,2500); break;
            case 243:case 247:case 251: add2(item,6,18000,27,3000); break;
            case 244:case 248:case 252: add2(item,6,22000,27,3500); break;
            case 245:case 249:case 253: add2(item,6,26000,27,4000); break;

            case 556:case 558:case 560: add2(item,6,50000,27,10000); break;

            // ===== GĂNG =====

            case 21:case 22:case 23: add(item,0,4); break;
            case 24:case 46:case 53: add(item,0,7); break;
            case 37:case 25:case 26: add(item,0,14); break;
            case 38:case 45:case 54: add(item,0,28); break;

            case 144:case 160:case 176: add(item,0,55); break;
            case 145:case 161:case 177: add(item,0,110); break;
            case 146:case 162:case 178: add(item,0,220); break;
            case 147:case 163:case 179: add(item,0,530); break;

            case 254:case 258:case 262: add(item,0,680); break;
            case 255:case 259:case 263: add(item,0,1000); break;
            case 256:case 260:case 264: add(item,0,1500); break;
            case 257:case 261:case 265: add(item,0,2200); break;

            case 562:case 564:case 566: add(item,0,5000); break;

            // ===== GIÀY =====

            case 27:case 28:case 29: add(item,7,10); break;
            case 30:case 47:case 55: add2(item,7,25,28,5); break;
            case 39:case 31:case 32: add2(item,7,120,28,24); break;
            case 40:case 48:case 56: add2(item,7,250,28,50); break;

            case 148:case 164:case 180: add2(item,7,500,28,100); break;
            case 149:case 165:case 181: add2(item,7,1200,28,240); break;
            case 150:case 166:case 182: add2(item,7,2400,28,480); break;
            case 151:case 167:case 183: add2(item,7,5000,28,1000); break;

            case 266:case 270:case 274: add2(item,7,9000,28,1500); break;
            case 267:case 271:case 275: add2(item,7,14000,28,2000); break;
            case 268:case 272:case 276: add2(item,7,19000,28,2500); break;
            case 269:case 273:case 277: add2(item,7,24000,28,3000); break;

            case 563:case 565:case 567: add2(item,7,60000,28,12000); break;

            // ===== RADA =====

            case 12:add(item,14,1);break;
            case 57:add(item,14,2);break;
            case 58:add(item,14,3);break;
            case 59:add(item,14,4);break;

            case 184:add(item,14,5);break;
            case 185:add(item,14,6);break;
            case 186:add(item,14,7);break;
            case 187:add(item,14,8);break;

            case 278:add(item,14,9);break;
            case 279:add(item,14,10);break;
            case 280:add(item,14,11);break;
            case 281:add(item,14,12);break;

            case 561:add(item,14,16);break;
        }
           if(isDoThan(item.template.id)){
            item.itemOptions.add(new Item.ItemOption(21,Util.nextInt(16,19)));
        }

    }
public static void showInfoCombine(Player player){

    if(player.combine == null || player.combine.itemsCombine == null){
        return;
    }

    int huyDiet = 0;
    int da1692 = 0;

    for(Item item : player.combine.itemsCombine){
        if(isDoHuyDiet(item)){
            huyDiet++;
        }
        if(item.template.id == 674){
            da1692 += item.quantity;
        }
    }

    if(huyDiet != 2){
        Service.gI().sendThongBaoOK(player,"Cần 2 trang bị Hủy Diệt");
        return;
    }

    int tiLe = 1;

    if(da1692 >= 5){
        tiLe = 10;
    }else if(da1692 >= 3){
        tiLe = 5;
    }

    player.combine.goldCombine = 1_000_000_000;

    CombineService.gI().baHatMit.createOtherMenu(
        player,
        ConstNpc.MENU_START_COMBINE,
        "Tỉ lệ ra đồ Thần: " + tiLe + "%\nCần 500Tr vàng\nCó thể thêm Thạch Thần tăng tỉ lệ\nBỏ thêm x5 Đá tăng 5%\nBỏ thêm 9 Đá tăng 10%",
        "Nâng cấp\n"+Util.numberToMoney(player.combine.goldCombine)+" vàng",
        "Từ chối"
    );
}
    // ================= COMBINE =================

    public static void startCombine(Player player){

if(player.combine.itemsCombine == null || player.combine.itemsCombine.size() < 2){
            Service.gI().sendThongBao(player,"Cần 2 trang bị Hủy Diệt");
            return;
        }

      Item item1 = null;
int huyDiet = 0;

for(Item item : player.combine.itemsCombine){

    if(item == null || item.template == null){
        continue;
    }

    if(isDoHuyDiet(item)){
        huyDiet++;
        if(item1 == null){
            item1 = item;
        }
    }
}

if(huyDiet != 2){
            Service.gI().sendThongBao(player,"Cần 2 trang bị Hủy Diệt");
            return;
        }

        int gold = player.combine.goldCombine;

        if(player.inventory.gold < gold){
            Service.gI().sendThongBao(player,"Bạn không đủ vàng");
            return;
        }
if(InventoryService.gI().getCountEmptyBag(player) <= 0){
    Service.gI().sendThongBao(player,"Hành trang cần ít nhất 1 ô trống");
    return;
}
   int type = Manager.ITEM_TEMPLATES.get(item1.template.id).type;

// gender theo logic bạn yêu cầu
int gender;
if(type == 4){
    gender = player.gender; // rada theo player
}else{
    gender = Manager.ITEM_TEMPLATES.get(item1.template.id).gender; // đồ thường theo item
}

int[] list = null;

switch(gender){
    case 0:
        switch(type){
            case 0:list=AO_TD;break;
            case 1:list=QUAN_TD;break;
            case 2:list=GANG_TD;break;
            case 3:list=GIAY_TD;break;
            case 4:list=RADA_TD;break;
        }
        break;

    case 1:
        switch(type){
            case 0:list=AO_NAMEK;break;
            case 1:list=QUAN_NAMEK;break;
            case 2:list=GANG_NAMEK;break;
            case 3:list=GIAY_NAMEK;break;
            case 4:list=RADA_NAMEK;break;
        }
        break;

    case 2:
        switch(type){
            case 0:list=AO_XAYDA;break;
            case 1:list=QUAN_XAYDA;break;
            case 2:list=GANG_XAYDA;break;
            case 3:list=GIAY_XAYDA;break;
            case 4:list=RADA_XAYDA;break;
        }
        break;

    case 3: // đồ all class
        switch(type){
            case 0:list=AO_TD;break;
            case 1:list=QUAN_TD;break;
            case 2:list=GANG_TD;break;
            case 3:list=GIAY_TD;break;
            case 4:list=RADA_TD;break;
        }
        break;
}

// chống crash
if(list == null){
    Service.gI().sendThongBao(player,"Trang bị không hỗ trợ nâng cấp");
    return;
}
int tileThan = 1; // mặc định 1%

Item daTangTiLe = null;

for(Item item : player.combine.itemsCombine){

    if(item == null || item.template == null){
        continue;
    }

    if(item.template.id == 674){
        daTangTiLe = item;
        if(item.quantity >= 5){
            tileThan = 10;
        }else if(item.quantity >= 3){
            tileThan = 5;
        }
        break;
    }
}

boolean raThan = Util.isTrue(tileThan,100);

short newId;

if(raThan){
    // ID cuối list là đồ thần
    newId = (short) list[list.length - 1];
}else{
    // random đồ thường (bỏ ID cuối)
    int random = Util.nextInt(list.length - 1);
    newId = (short) list[random];
}

Item newItem = ItemService.gI().createNewItem(newId);

        newItem.itemOptions.clear();
if(raThan){
    ServerNotify.gI().notify(
        "Thằng Chó " + player.name + " vừa đớp "
        + newItem.template.name + " Kích Hoạt!"
    );
}
        setOptionMacDinh(newItem);

 

        int[][] op;

        if(gender==0) op = new int[][]{{127,139},{128,140},{129,141}};
        else if(gender==1) op = new int[][]{{130,143},{132,144},{251,254}};
        else op = new int[][]{{133,136},{134,137},{135,138}};

        int rd = Util.nextInt(op.length);

        newItem.itemOptions.add(new Item.ItemOption(op[rd][0],0));
        newItem.itemOptions.add(new Item.ItemOption(op[rd][1],0));

        newItem.itemOptions.add(new Item.ItemOption(30,0));

        player.inventory.gold -= gold;

     for(Item item:player.combine.itemsCombine){

    if(item == null || item.template == null){
        continue;
    }

    if(isDoHuyDiet(item)){
        InventoryService.gI().subQuantityItemsBag(player,item,1);
    }
}
for(Item item:player.combine.itemsCombine){
     if(item == null || item.template == null){
        continue;
    }
    if(item.template.id == 674){
        if(item.quantity >= 5){
            InventoryService.gI().subQuantityItemsBag(player,item,5);
        }else if(item.quantity >= 3){
            InventoryService.gI().subQuantityItemsBag(player,item,3);
        }
        break;
    }
}
        InventoryService.gI().addItemBag(player,newItem);

        CombineService.gI().sendEffectSuccessCombine(player);

        InventoryService.gI().sendItemBag(player);
        Service.gI().sendMoney(player);

        CombineService.gI().reOpenItemCombine(player);
    }
}