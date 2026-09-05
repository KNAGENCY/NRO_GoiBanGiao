package game

import (
	"fmt"
	"log/slog"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
)

// lenhMoNPC: client gui LEN khi bam vao NPC (Service.openMenu, cmd 33: [npcId short]).
const lenhMoNPC int8 = 33

// lenhChonMenuNPC: client gui LEN khi chon mot dong trong menu NPC
// (Service.confirmMenu, cmd 32: [npcID short][chon byte]). Cung so 32 voi goi
// XUONG "hoi thoai NPC" (goiMoNPC) — hai chieu, hai khuon khac nhau.
const lenhChonMenuNPC int8 = 32

// cuaHangIruka: id cua hang trong kho-noi-dung/cua-hang.yaml ma thay Iruka mo
// ("Tạp hoá Làng Lá"). Cua hang thu hai cho NPC khac -> them dong vao menuNPC.
const cuaHangIruka int32 = 1

// menuNPC tra menu hoi thoai theo mau NPC. MENU o day THAY menu trong
// catalog.json (["Đóng"]): catalog bi GHIM SHA trong asset-lock (sua mot chu la
// doi lock 444KB), va menu la NOI DUNG ta quyet dinh — dung cho cua no la
// kho-noi-dung (KH9 npc.yaml). Bang Go nay la cau noi de KH7 nghiem thu duoc
// (mo cua hang qua NPC) truoc khi KH9 bat dau. Thu tu dong = chi so `chon`
// client gui len o cmd 32 — KHONG doi thu tu ma khong sua xuLyChonMenuNPC.
func menuNPC(id byte) []string {
	if id == catalogLangLa.NPC.ID {
		return []string{"Cửa hàng", "Đóng"}
	}
	return nil
}

func goiMoNPC(id byte) (*protocol.Message, error) {
	c := catalogLangLa
	if id != c.NPC.ID {
		return nil, fmt.Errorf("npc %d khong co trong catalog", id)
	}
	menu := menuNPC(id)
	w := protocol.NewWriter()
	w.Short(int16(c.NPC.ID))
	if err := w.UTF(c.NPC.Dialog); err != nil {
		return nil, err
	}
	w.Byte(byte(len(menu)))
	for _, dong := range menu {
		if err := w.UTF(dong); err != nil {
			return nil, err
		}
	}
	w.Short(c.NPC.Portrait)
	return protocol.NewMessage(32, w), nil
}

func ghiNPCMap(w *protocol.Writer, bd *BanDo) {
	if bd == nil || bd.MapID != 1 {
		w.Byte(0)
		return
	}
	c := catalogLangLa.NPC
	w.Byte(1)
	w.Byte(0)
	w.Short(c.X)
	w.Short(c.Y)
	w.Byte(c.ID)
	w.Short(c.Portrait)
}

// xuLyChonMenuNPC: nguoi choi chon mot dong trong menu NPC (cmd 32 LEN).
// Iruka dong 0 "Cửa hàng" -> nho cua hang dang mo cho phien (DatCuaHangDangMo,
// muaban.go — goi 6/7 mua/ban khong mang id cua hang) roi gui -44. Dong "Đóng"
// hay npc/dong la -> khong lam gi, khong log (goi rac lap lai duoc).
//
// Chay tren goroutine PHIEN: TimCuaHang doc anh chup kho noi dung (bang muaban.go
// da ghi), goiCuaHang thuan; khong cham actor.
func (b *BoXuLy) xuLyChonMenuNPC(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil || b.nhanVat.MapID != 1 {
		return nil
	}
	r := protocol.NewReader(m.Data)
	id, err := r.Short()
	if err != nil {
		return nil
	}
	chon, err := r.Byte()
	if err != nil {
		return nil
	}
	if id != int16(catalogLangLa.NPC.ID) || chon != 0 {
		return nil // "Đóng" hoac dong khong co
	}
	ch := TimCuaHang(cuaHangIruka)
	if ch == nil {
		// Kho noi dung khong con cua hang nay: loi DU LIEU, phai thay trong log,
		// khong phai goi rac. Khong gui gi de client khong mo mot bang rong.
		slog.Error("menu NPC 'Cửa hàng': khong tim thay cua hang trong kho noi dung",
			"cua_hang", cuaHangIruka, "id_nhan_vat", b.nhanVat.ID)
		return nil
	}
	goi, err := goiCuaHang(ch)
	if err != nil {
		slog.Error("menu NPC 'Cửa hàng': dung goi -44 that bai", "loi", err)
		return nil
	}
	b.DatCuaHangDangMo(ch)
	return s.XepGui(goi)
}
