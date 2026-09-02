package game

import (
	"context"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// KH7 Task 5 — kiem tich hop DAU-CUOI tren CSDL THAT:
// tao nhan vat -> cap 2 mon vao o 0 va o 2 cua tui -> dang nhap lai -> vao the
// gioi -> ba goi hanh trang phai mang dung hai mon do, o 1 phai RONG, va bang
// tra o -> id trong BoXuLy phai tra dung id CSDL.
//
// Can NRO_TEST_DSN (moKhoThatChoTest tu t.Skip neu chua dat).
// ---------------------------------------------------------------------------

// docBangHanhTrangTuGoi doc than goi -36 hoac -35 (khuon: [cAction][soO][o]*)
// va tra ve mau cua tung o (-1 = rong) cung so luong tung o.
func docBangHanhTrangTuGoi(t *testing.T, m *protocol.Message, ten string) (mau []int16, soLuong []int32) {
	t.Helper()
	r := protocol.NewReader(m.Data)
	if a, err := r.Byte(); err != nil || a != hanhTrangNapCaBang {
		t.Fatalf("%s: cAction = %d (loi %v), mong doi %d", ten, a, err, hanhTrangNapCaBang)
	}
	soO, err := r.Byte()
	if err != nil {
		t.Fatalf("%s: doc soO that bai: %v", ten, err)
	}
	for i := 0; i < int(soO); i++ {
		id, err := r.Short()
		if err != nil {
			t.Fatalf("%s: doc mau o %d/%d that bai: %v", ten, i, soO, err)
		}
		if id == -1 {
			mau = append(mau, -1)
			soLuong = append(soLuong, 0)
			continue
		}
		sl, err := r.Int()
		if err != nil {
			t.Fatalf("%s o %d: doc so luong that bai: %v", ten, i, err)
		}
		if _, err := r.UTF(); err != nil {
			t.Fatalf("%s o %d: doc info that bai: %v", ten, i, err)
		}
		if _, err := r.UTF(); err != nil {
			t.Fatalf("%s o %d: doc content that bai: %v", ten, i, err)
		}
		soOpt, err := r.Byte()
		if err != nil {
			t.Fatalf("%s o %d: doc so option that bai: %v", ten, i, err)
		}
		for j := 0; j < int(soOpt); j++ {
			if _, err := r.Byte(); err != nil {
				t.Fatalf("%s o %d: doc optionId that bai: %v", ten, i, err)
			}
			if _, err := r.Short(); err != nil {
				t.Fatalf("%s o %d: doc param that bai: %v", ten, i, err)
			}
		}
		mau = append(mau, id)
		soLuong = append(soLuong, sl)
	}
	return mau, soLuong
}

func TestTichHop_VaoTheGioiNapHanhTrangThat(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	tenTK := "nguoi" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taiKhoanID := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Naru" + hauToNganNgau(t)

	// --- Buoc 1: tao nhan vat (phien 1, dong ngay sau do). ---
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK, matKhau)
		if !cl.choDungMotLenh(t, 2) {
			t.Fatal("lan dau dang nhap phai nhan DUNG MOT lenh 2 (man tao nhan vat)")
		}
		cl.taoNhanVat(tenNV, 0, 5)
		if !cl.choVaoTheGioiRoiImLang(t) {
			t.Fatal("tao nhan vat hop le phai nhan DUNG MOT goi -28/4 roi im lang")
		}
	}

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}

	// Don dep vat pham TRUOC khi taoTaiKhoanCoMatKhau xoa nhan_vat: vat_pham va
	// vat_pham_lich_su deu REFERENCES nhan_vat(id), khong xoa truoc thi lenh
	// DELETE nhan_vat cua ham do that bai vi khoa ngoai. Cleanup chay LIFO nen
	// dang ky o day (SAU ham do) la chay TRUOC — dung y do.
	t.Cleanup(func() {
		pool := moPoolThuChoTest(t)
		if _, err := pool.Exec(context.Background(),
			`DELETE FROM vat_pham_lich_su WHERE vat_pham_id IN
			     (SELECT id FROM vat_pham WHERE chu_so_huu_id = $1)`, nv.ID); err != nil {
			t.Errorf("don dep vat_pham_lich_su cua nhan vat %d that bai: %v", nv.ID, err)
		}
		if _, err := pool.Exec(context.Background(),
			`DELETE FROM vat_pham WHERE chu_so_huu_id = $1`, nv.ID); err != nil {
			t.Errorf("don dep vat_pham cua nhan vat %d that bai: %v", nv.ID, err)
		}
	})

	// --- Buoc 2: cap 2 mon vao O 0 va O 2 cua tui. O 1 CO CHU DICH de rong:
	// day la o bat loi L1 (bo qua o rong ma khong ghi short -1). ---
	idO0, err := khoThat.CapVatPham(ctx, nv.ID, store.NoiChuaTui, 0, 0, 7,
		store.ChiSoPhu{}, "SINH_RA")
	if err != nil {
		t.Fatalf("CapVatPham o 0 that bai: %v", err)
	}
	idO2, err := khoThat.CapVatPham(ctx, nv.ID, store.NoiChuaTui, 2, 1, 3,
		store.ChiSoPhu{}, "SINH_RA")
	if err != nil {
		t.Fatalf("CapVatPham o 2 that bai: %v", err)
	}

	// --- Buoc 3: dang nhap lai, di het chuoi vao the gioi. ---
	cl := moPhienClientGia(t, khoThat)
	defer cl.dong()
	cl.dangNhap(tenTK, matKhau)
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

	// --- Buoc 4: THU TU BAT BUOC (pha A). ME_INFO cap lai ca ba mang vat pham
	// (Controller.cs:5860/:5907/:5942), nen -36/-37/-35 den TRUOC no la bi xoa
	// sach khong bao. Doc TUNG goi mot theo dung thu tu de dot bien "gui -36
	// truoc ME_INFO" lam test do.
	thuTu := []int8{-82, lenhMeInfo, -24, lenhChiSoNhanVat, lenhDoDangMac, lenhTui, lenhRuong}
	goiTheoLenh := make(map[int8]*protocol.Message, len(thuTu))
	for i, mong := range thuTu {
		m, errDoc := cl.doiGoiLenhKhongFatal(mong)
		if errDoc != nil {
			t.Fatalf("goi thu %d cua chuoi vao the gioi: khong nhan duoc lenh %d: %v", i, mong, errDoc)
		}
		goiTheoLenh[mong] = m
	}

	// --- Buoc 5: ME_INFO mang dung hai mon, o 1 rong. ---
	rMe := doDenBaMangVatPham(t, nv, goiTheoLenh[lenhMeInfo])
	if soO, _ := docMotBang(t, rMe, "itemsBody"); soO != SoOMacChuan {
		t.Fatalf("ME_INFO itemsBody.size = %d, mong doi %d", soO, SoOMacChuan)
	}
	soOTui, mauTuiMe := docMotBang(t, rMe, "itemsBag")
	if soOTui != SoOTuiChuan {
		t.Fatalf("ME_INFO itemsBag.size = %d, mong doi %d (SO O)", soOTui, SoOTuiChuan)
	}
	if mauTuiMe[0] != 0 || mauTuiMe[1] != -1 || mauTuiMe[2] != 1 {
		t.Fatalf("ME_INFO tui: o0/o1/o2 = %d/%d/%d, mong doi 0/-1/1",
			mauTuiMe[0], mauTuiMe[1], mauTuiMe[2])
	}

	// --- Buoc 6: goi -36 mang dung hai mon, o 1 rong. ---
	mauTui, slTui := docBangHanhTrangTuGoi(t, goiTheoLenh[lenhTui], "-36 (tui)")
	if len(mauTui) != SoOTuiChuan {
		t.Fatalf("-36 co %d o, mong doi %d", len(mauTui), SoOTuiChuan)
	}
	if mauTui[0] != 0 || slTui[0] != 7 {
		t.Errorf("-36 o 0 = mau %d so luong %d, mong doi mau 0 so luong 7", mauTui[0], slTui[0])
	}
	if mauTui[1] != -1 {
		t.Errorf("-36 o 1 = mau %d, mong doi -1 (o RONG giua hai mon)", mauTui[1])
	}
	if mauTui[2] != 1 || slTui[2] != 3 {
		t.Errorf("-36 o 2 = mau %d so luong %d, mong doi mau 1 so luong 3", mauTui[2], slTui[2])
	}
	demMon := 0
	for _, id := range mauTui {
		if id != -1 {
			demMon++
		}
	}
	if demMon != 2 {
		t.Errorf("-36 co %d o co mon, mong doi dung 2", demMon)
	}

	// --- Buoc 7: bang tra o -> vat_pham.id trong BoXuLy phai dung. ---
	if cl.b.hanhTrang == nil {
		t.Fatal("BoXuLy khong giu bang hanh trang sau khi vao the gioi (b.hanhTrang == nil)")
	}
	if got, co := cl.b.hanhTrang.TraId(store.NoiChuaTui, 0); !co || got != idO0 {
		t.Errorf("TraId(TUI,0) = %d,%v; mong doi %d,true", got, co, idO0)
	}
	if got, co := cl.b.hanhTrang.TraId(store.NoiChuaTui, 2); !co || got != idO2 {
		t.Errorf("TraId(TUI,2) = %d,%v; mong doi %d,true", got, co, idO2)
	}
	if got, co := cl.b.hanhTrang.TraId(store.NoiChuaTui, 1); co {
		t.Errorf("TraId(TUI,1) = %d,true; o 1 RONG nen phai tra co=false", got)
	}
	if got, co := cl.b.hanhTrang.TraId(store.NoiChuaRuong, 0); co {
		t.Errorf("TraId(RUONG,0) = %d,true; ruong rong nen phai tra co=false (o so 0 cua TUI va cua RUONG la hai mon khac nhau)", got)
	}
}
