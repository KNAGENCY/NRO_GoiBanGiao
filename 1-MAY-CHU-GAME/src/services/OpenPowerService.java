package services;

import item.Item;
import player.NPoint;
import player.Pet;
import player.Player;

public class OpenPowerService {

    public static final int COST_SPEED_OPEN_LIMIT_POWER = 500000000;

    private static OpenPowerService i;

    private OpenPowerService() {
    }

    public static OpenPowerService gI() {
        if (i == null) {
            i = new OpenPowerService();
        }
        return i;
    }

    private boolean isWearingHuyDietItem(Player player) {
        for (Item item : player.inventory.itemsBody) {
            if (item != null && item.isNotNullItem()) {
                if (item.template.id >= 650 && item.template.id <= 662) {
                    return true;
                }
            }
        }
        return false;
    }

    private boolean hasDaTuTien(Player player) {
        for (Item item : player.inventory.itemsBag) {
            if (item != null && item.isNotNullItem()) {
                if (item.template.id == 674 && item.quantity > 0) {
                    return true;
                }
            }
        }
        return false;
    }

    private void removeDaTuTien(Player player) {
        for (Item item : player.inventory.itemsBag) {
            if (item != null && item.isNotNullItem() && item.template.id == 674 && item.quantity > 0) {
                item.quantity--;
                if (item.quantity <= 0) {
                    player.inventory.itemsBag.remove(item);
                }
                InventoryService.gI().sendItemBag(player);
                return;
            }
        }
    }

    private boolean checkRequire(Player p, byte curLimit) {
        Player master = p.isPet ? ((Pet) p).master : p;
        if (curLimit == 7 || curLimit == 8) {
            boolean hasHuyDietItem = isWearingHuyDietItem(master)
                    || (p.isPet && isWearingHuyDietItem(p));
            if (!hasHuyDietItem) {
                Service.gI().sendThongBaoOK(master, "Cần Mặc Ít Nhất 1 Món Hủy Diệt Để Mở!");
                return false;
            }
        } else if (curLimit >= 9) {
            if (!hasDaTuTien(master)) {
                Service.gI().sendThongBaoOK(p, "Cần 1 Thạch Thần Để Mở!");
                return false;
            }
            removeDaTuTien(master);
        }
        return true;
    }

    private void updatePointAfterOpenPower(Player player) {
        player.nPoint.initPowerLimit();
        player.nPoint.calPoint();
        if (!player.isPet) {
            Service.gI().point(player);
        } else {
            Player master = ((Pet) player).master;
            Service.gI().sendChiSoPetGoc(master);
            Service.gI().showInfoPet(master);
            Service.gI().point(master);
        }
    }

    public boolean openPowerBasic(Player player) {
        byte curLimit = player.nPoint.limitPower;
        if (curLimit < NPoint.MAX_LIMIT) {
            if (!player.itemTime.isOpenPower && player.nPoint.canOpenPower()) {
                if (!checkRequire(player, curLimit)) {
                    return false;
                }
                player.itemTime.isOpenPower = true;
                player.itemTime.lastTimeOpenPower = System.currentTimeMillis();
                ItemTimeService.gI().sendAllItemTime(player);
                return true;
            } else {
                return false;
            }
        } else {
            Service.gI().sendThongBao(player, "Sức mạnh của bạn đã đạt tới mức tối đa");
            return false;
        }
    }

    public boolean openPowerSpeed(Player player) {
        byte curLimit = player.nPoint.limitPower;
        if (curLimit < NPoint.MAX_LIMIT) {
            if (!checkRequire(player, curLimit)) {
                return false;
            }
            player.nPoint.limitPower++;
            if (player.nPoint.limitPower > NPoint.MAX_LIMIT) {
                player.nPoint.limitPower = NPoint.MAX_LIMIT;
            }
            updatePointAfterOpenPower(player);
            if (!player.isPet) {
                Service.gI().sendThongBao(player, "Giới hạn sức mạnh của bạn đã được tăng lên 1 bậc");
            } else {
                Service.gI().sendThongBao(((Pet) player).master, "Giới hạn sức mạnh của đệ tử đã được tăng lên 1 bậc");
            }
            return true;
        } else {
            if (!player.isPet) {
                Service.gI().sendThongBao(player, "Sức mạnh của bạn đã đạt tới mức tối đa");
            } else {
                Service.gI().sendThongBao(((Pet) player).master, "Sức mạnh của đệ tử đã đạt tới mức tối đa");
            }
            return false;
        }
    }
}
