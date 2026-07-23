package npc.npc_manifest;

/**
 *
 * @author GiaCatLuong
 */
import models.Consign.ConsignShopService;
import npc.Npc;
import player.Player;
import services.Service;

public class KyGui extends Npc {

    public KyGui(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (!canOpenNpc(player)) {
            return;
        }

        // Chưa hoàn thành nhiệm vụ Fide
        if (player.playerTask.taskMain.id < 24) {
            createOtherMenu(player, 0,
                    "Con cần hoàn thành nhiệm vụ Fide\nmới có thể sử dụng chức năng Ký Gửi.",
                    "Đóng");
            return;
        }

        createOtherMenu(player, 0,
                "Toàn bộ Shop đều dùng thỏi vàng\n"
                + "Ký gửi ít nhất 10 THỎI VÀNG\n"
                + "Phí ký gửi là 2 THỎI VÀNG\n"
                + "Phí uptop là 1 THỎI VÀNG\n"
                + "Phí Phương Anh là ca sĩ tôi yêu\n"
                + "CHÚNG TÔI CHỈ ĐỚP 10% HOA HỒNG\n"
                + "(Vàng và vật phẩm bị khóa đều không thể sử dụng)",
                "Mua bán\nKý gửi", "Từ chối", "Cướp\ntiệm vàng");
    }

    @Override
    public void confirmMenu(Player pl, int select) {
        if (!canOpenNpc(pl)) {
            return;
        }

        // Chặn sử dụng nếu chưa hoàn thành nhiệm vụ Fide
        if (pl.playerTask.taskMain.id < 24) {
            Service.gI().sendThongBao(pl,
                    "Cần hoàn thành nhiệm vụ Fide mới có thể sử dụng Ký Gửi!");
            return;
        }

        switch (select) {
            case 0 -> ConsignShopService.gI().openShopKyGui(pl);

            case 1 -> {
                // Chức năng Ký gửi (nếu có)
            }

            case 2 -> Service.gI().sendThongBaoOK(pl, "Chúng tôi chỉ có Vàng Nổi");
        }
    }
}