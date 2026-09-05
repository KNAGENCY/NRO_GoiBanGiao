package game

import (
	"testing"

	"nroserver/internal/protocol"
)

// TestTichHop_NPCIrukaMoCuaHangQuaMenu: duong nguoi choi THAT se di de mua ban
// (KH7 nghiem thu tay): vao the gioi -> bam Iruka (33) -> nhan hoi thoai 32 co
// dong "Cửa hàng" -> chon dong 0 (32 LEN) -> nhan -44 va phien da nho cua hang
// dang mo (goi 6/7 mua/ban khong mang id cua hang). Truoc 05/09 khong co duong
// nay: menu chi co "Đóng", DatCuaHangDangMo chua ai goi -> mua/ban khong bao gio
// chay tu client that.
//
// Cung khang dinh goi 40 (nhiem vu giu cho) den TRUOC -24: doi 40 truoc, roi doi
// -24 — neu 40 den sau -24 thi buoc doi -24 het gio (goi -24 da bi bo qua).
func TestTichHop_NPCIrukaMoCuaHangQuaMenu(t *testing.T) {
	kho := moKhoThatChoTest(t)
	hauTo := hauToNganNgau(t)
	so := moSoKhuVucKhongQuaiMacDinh(t)
	cl := moPhienClientGiaKhuVuc(t, kho, so)
	defer cl.dong()

	tenTK := "npc" + hauTo
	taoTaiKhoanCoMatKhau(t, tenTK, "1")
	cl.dangNhap(tenTK, "1")
	if !cl.choLenh(t, 2) {
		t.Fatal("tai khoan moi phai nhan lenh 2 (tao nhan vat)")
	}
	cl.taoNhanVat("Npc"+hauTo, 0, 5)
	kiemVersionGame(t, cl.doiGoi(t, lenhKhongTrenBanDo, lenhConVersionGame))
	cl.xinData()
	cl.doiGoiLenh(t, lenhDuLieu)
	cl.xinCatalog(lenhConCatalogMap)
	cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogMap)
	cl.xinCatalog(lenhConCatalogSkill)
	cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogSkill)
	cl.xinCatalog(lenhConCatalogItem)
	doiChuoiCatalogItem(t, cl)
	cl.guiClientOk()
	cl.doiGoiLenh(t, -82)
	cl.doiGoiLenh(t, lenhMeInfo)
	// 40 PHAI truoc -24 (client doc NPC trong -24 can taskMaint != null).
	m40 := cl.doiGoiLenh(t, lenhNhiemVuChinh)
	if r := protocol.NewReader(m40.Data); true {
		if id, _ := r.Short(); id < 7 {
			t.Fatalf("goi 40 taskId = %d, phai >= 7", id)
		}
	}
	cl.doiGoiLenh(t, -24)
	cl.doiGoiLenh(t, lenhChiSoNhanVat)
	cl.doiGoiLenh(t, lenhDoDangMac)
	cl.doiGoiLenh(t, lenhTui)
	cl.doiGoiLenh(t, lenhRuong)
	cl.baoVaoKhuVuc()

	// Bam Iruka: 33 [npcId short 0].
	w := protocol.NewWriter()
	w.Short(0)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhMoNPC, w)); err != nil {
		t.Fatal(err)
	}
	m32 := cl.doiGoiLenh(t, 32)
	r := protocol.NewReader(m32.Data)
	r.Short()
	r.UTF()
	n, _ := r.Byte()
	dong0, _ := r.UTF()
	if n != 2 || dong0 != "Cửa hàng" {
		t.Fatalf("hoi thoai Iruka: %d dong, dong 0 = %q; mong 2 dong, dong 0 'Cửa hàng'", n, dong0)
	}

	// Chon dong 0: 32 LEN [npcID short 0][chon byte 0] -> -44.
	w = protocol.NewWriter()
	w.Short(0)
	w.Byte(0)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhChonMenuNPC, w)); err != nil {
		t.Fatal(err)
	}
	m44 := cl.doiGoiLenh(t, lenhMoCuaHang)
	r = protocol.NewReader(m44.Data)
	loai, _ := r.Byte()
	soTab, _ := r.Byte()
	if loai != loaiShopThuong || soTab == 0 {
		t.Fatalf("-44: loaiShop=%d soTab=%d", loai, soTab)
	}
	// Phien da nho cua hang dang mo (net.Pipe dong bo: -44 den sau khi handler
	// gan truong nay).
	if cl.b.muaBanCuaHangDangMo == nil || cl.b.muaBanCuaHangDangMo.ID != cuaHangIruka {
		t.Fatalf("phien chua nho cua hang dang mo: %+v", cl.b.muaBanCuaHangDangMo)
	}

	// Chon "Đóng" (dong 1) -> khong gui gi.
	w = protocol.NewWriter()
	w.Short(0)
	w.Byte(1)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhChonMenuNPC, w)); err != nil {
		t.Fatal(err)
	}
	if !cl.khongConGoiGi(t) {
		t.Fatal("chon 'Đóng' ma may chu van gui goi")
	}
}
