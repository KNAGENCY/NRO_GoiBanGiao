package npc.npc_manifest;

/**
 * @author GiaCatLuong
 */
import npc.Npc;
import player.Player;
import services.Service;
import services.func.ChangeMapService;

public class Tapion extends Npc {

    public Tapion(int mapId, int status, int cx, int cy, int tempId, int avartar) {
        super(mapId, status, cx, cy, tempId, avartar);
    }

    @Override
    public void openBaseMenu(Player player) {
        if (canOpenNpc(player)) {
            this.createOtherMenu(player, 0,
                    "Ở đây,\n"
                    + "- Nếu ngươi có đủ Sét Hủy Diệt có thể up mảnh Thiên Sứ,\n"
                    + "- Nếu ngươi hạ Boss thì sẽ có tỉ lệ rơi Đá, Công Thức, Mảnh Thiên Sứ,\n"
                    + "- Nếu ngươi đẹp trai thì tỉ lệ cao hơn so với bình thường.",
                    "Chúc \nmay mắn", "Trở về");
        }
    }

    @Override
    public void confirmMenu(Player player, int select) {
        if (canOpenNpc(player)) {
            switch (select) {

                // Chúc may mắn -> không làm gì
                case 0 -> {

                }

                // Trở về -> về map 5
                case 1 -> {
                    ChangeMapService.gI().changeMapNonSpaceship(player, 5, 1000, 408);
                }
            }
        }
    }
}