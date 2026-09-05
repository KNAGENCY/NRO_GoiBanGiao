package game

import (
	"fmt"

	"nroserver/internal/protocol"
)

// nhiemvu.go — NHIEM VU CHINH GIU CHO (placeholder) cho client.
//
// Vi sao phai gui mot "nhiem vu" khi game chua co he nhiem vu (KH9 chua bat dau):
// client Unity la ban port Dragon Ball va COI NHIEM VU CHINH LUON TON TAI.
// `Char.myCharz().taskMaint` chi duoc gan trong case 40 (Controller.cs:3273
// "TASK_GET"); may chu chua bao gio gui goi do nen truong ay la null, va client
// deref no o nhieu cho KHONG kiem null:
//   - Controller.cs:5145 (doc NPC trong MAP INFO -24): `taskMaint.taskId >= 7`
//     -> NullReferenceException -> NPC KHONG BAO GIO duoc them vao vNpc. Day la
//     ly do thay Iruka khong hien du may chu ghi dung 1 NPC vao -24 (phat hien
//     05/09/2026 qua Unity MCP: arrNpcTemplate=1, vNpc=0, taskMaint=null).
//   - GameScr.cs:3876, :5774 (`taskMaint.taskId <= 1`, `== 0`) trong vong ve/
//     phim — nuot ngoai le im lang moi khung hinh.
//
// Khong co NPC thi khong mo duoc cua hang -> KH7 khong nghiem thu tay duoc.
//
// Goi nay gui MOT lan trong chuoi vao the gioi, SAU ME_INFO va TRUOC -24 (vi
// -24 la noi doc NPC). Noi dung la "nhiem vu chao mung" trung tinh; KH9 se thay
// bang he nhiem vu that (va doi nguon sang kho-noi-dung).
//
// Khuon byte (Controller.cs:3273-3316, doc theo dung thu tu):
//
//	[taskId short][index byte][ten UTF][chi tiet UTF]
//	[soMucCon byte]( [tenMuc UTF][npcDich byte][mapDich short][ghiChu UTF] )*
//	[count short]( [dem short] )*soMucCon      -- phan nay client boc try/catch
//
// taskId: 100. Bo loc NPC o -24 la
//
//	b4 != 6 && ((taskId >= 7 && (taskId != 7 || index > 1)) || b4 not in {7,8,9}) && (taskId >= 6 || b4 != 16)
//
// -> taskId >= 7 cho MOI mau NPC di qua (tru 6 va 16 la NPC dac thu Dragon Ball).
// KHONG dung 0: taskId 0 la "nhiem vu mo dau" cua Dragon Ball, client re nhanh
// theo NPC mau 3/4 (GameScr.cs:2151-2157) va chan mot so phim (GameScr.cs:5774).
const lenhNhiemVuChinh int8 = 40

const (
	nhiemVuGiuChoID int16 = 100
	nhiemVuNPCDich  byte  = 0 // mau NPC Iruka (catalog.json npcs[0].id) — radar tro toi
	nhiemVuMapDich  int16 = 1 // Lang La
)

// goiNhiemVuChinh dung goi 40 giu cho. Loi chi co the den tu writeUTF (chuoi
// qua dai) — cac chuoi o day la hang ngan, nen loi la loi lap trinh.
func goiNhiemVuChinh() (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Short(nhiemVuGiuChoID)
	w.Byte(0) // index: dang o muc con dau tien
	for _, s := range []string{
		"Chào mừng đến Làng Lá",
		"Hãy làm quen với làng và trò chuyện với thầy Iruka.",
	} {
		if err := w.UTF(s); err != nil {
			return nil, fmt.Errorf("goi 40: %w", err)
		}
	}
	w.Byte(1) // mot muc con
	if err := w.UTF("Trò chuyện với thầy Iruka"); err != nil {
		return nil, fmt.Errorf("goi 40: %w", err)
	}
	w.Byte(nhiemVuNPCDich)
	w.Short(nhiemVuMapDich)
	if err := w.UTF(""); err != nil { // ghi chu them (contentInfo) — rong
		return nil, fmt.Errorf("goi 40: %w", err)
	}
	w.Short(0) // count
	w.Short(0) // dem cua muc con 0
	return protocol.NewMessage(lenhNhiemVuChinh, w), nil
}
