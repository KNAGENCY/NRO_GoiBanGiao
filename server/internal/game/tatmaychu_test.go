package game

import (
	"context"
	"testing"
	"time"

	"nroserver/internal/store"
)

// TestTatMayChuLuuTruocKhiDungKhuVuc la test CHINH cua Task 6 diem A (brief
// Step 1): co nguoi trong khu vuc -> goi TatMayChu -> kho gia (thong qua
// ThoGhi that, khong phai gia lap tat) phai nhan DU ban ghi cua nguoi do
// TRUOC KHI ham tra ve, VA khu vuc phai THUC SU da dung (actor da Dung())
// SAU KHI ham tra ve — kiem ca hai nua cua dac ta §4 diem C trong MOT test.
//
// Day CUNG la test se bat dot bien Step 5 cua brief (chuyen DungTatCa() len
// TRUOC buoc luu): neu TatMayChu goi DungTatCa() TRUOC LuuTatCa/
// DongVaChoXong, actor se dung TRUOC KHI AnhChupNguoiChoi kip hoi -> kho gia
// se KHONG nhan duoc ban ghi nao (assert (1) duoi day do), lam test nay DO.
func TestTatMayChuLuuTruocKhiDungKhuVuc(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() }) // phong ho: neu mot assert Fatal giua chung, actor van duoc dung sach

	nv := &store.NhanVat{
		ID: 777, Ten: "Obito", MapID: 1,
		ViTriX: 12345, ViTriY: 6789,
		HpToiDa: 100, HpHienTai: 42,
		ChakraToiDa: 50, ChakraHienTai: 30,
	}
	k.Gui(suKienVao{nv: nv, phien: &phienGiaKhuVuc{}})

	s := &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}
	kg := &khoLuuGia{}
	tg := MoThoGhi(kg, 2)

	ctx, huy := context.WithTimeout(context.Background(), 5*time.Second)
	defer huy()
	if err := TatMayChu(ctx, s, tg); err != nil {
		t.Fatalf("TatMayChu tra loi: %v", err)
	}

	// (1) kho gia PHAI da nhan dung ban ghi cua nguoi choi — day chinh la
	// bang chung "da luu TRUOC khi tra ve": TatMayChu chi tra ve SAU khi
	// DongVaChoXong xac nhan hang doi da duoc ghi HET xuong kho gia.
	kg.mu.Lock()
	ban := append([]store.NhanVat(nil), kg.goiGhi...)
	kg.mu.Unlock()
	if len(ban) != 1 {
		t.Fatalf("kho gia nhan %d ban ghi, mong doi DUNG 1 (nguoi choi duy nhat trong khu vuc)", len(ban))
	}
	if ban[0].ID != nv.ID || ban[0].HpHienTai != nv.HpHienTai || ban[0].ViTriX != nv.ViTriX {
		t.Fatalf("ban ghi kho gia nhan = %+v, mong doi khop voi nguoi choi %+v", ban[0], nv)
	}

	// (2) khu vuc PHAI da dung THAT SU sau khi TatMayChu tra ve: goi
	// AnhChupNguoiChoi lan nua — actor da Dung() nen se KHONG con ai tra loi,
	// tra ve rong sau 100ms (xem AnhChupNguoiChoi, khuvuc.go). nv van dang
	// dang ky trong khu vuc (chua tung Gui(suKienRa)) nen neu actor VAN CON
	// SONG, ham nay se van thay nv — tra ve rong o day la bang chung actor
	// DA CHET, khong phai "tinh co khong co ai".
	if conSong := k.AnhChupNguoiChoi(0); len(conSong) != 0 {
		t.Fatalf("sau TatMayChu, AnhChupNguoiChoi van tra ve %d nguoi — khu vuc CHUA thuc su dung "+
			"(DungTatCa co the chua duoc goi, hoac goi truoc buoc luu)", len(conSong))
	}
}

// TestTatMayChuVanDungKhuVucKhiHetGioLuu xac nhan TatMayChu KHONG bo qua
// buoc DungTatCa() du DongVaChoXong loi (vd CSDL treo qua han gio, dac ta
// §6.2) — treo ca tien trinh cho mai mai con te hon mot lan tat "danh hong".
//
// Dung context DA HET GIO TU TRUOC (sleep qua han truoc khi goi TatMayChu)
// de ep DongVaChoXong tra ve loi ngay, roi kiem DungTatCa VAN duoc goi bang
// each gui mot nguoi choi MOI sau khi TatMayChu tra ve: neu actor CON SONG,
// no se xu ly suKienVao va AnhChupNguoiChoi se thay nguoi do; actor DA CHET
// thi ca hai su kien nam im trong hop thu, khong ai tra loi -> tra ve rong.
func TestTatMayChuVanDungKhuVucKhiHetGioLuu(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })
	s := &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}

	kg := &khoLuuGia{doTre: 2 * time.Second} // kho gia CHAM hon nhieu so voi han gio ctx duoi day
	tg := MoThoGhi(kg, 1)

	ctx, huy := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer huy()
	time.Sleep(5 * time.Millisecond) // dam bao ctx CHAC CHAN da het gio truoc khi goi TatMayChu

	err := TatMayChu(ctx, s, tg)
	if err == nil {
		t.Fatal("TatMayChu khong tra loi du ctx da het gio tu truoc khi goi — mong doi ctx.Err()")
	}

	nvMoi := &store.NhanVat{ID: 999, Ten: "SauKhiTat"}
	k.Gui(suKienVao{nv: nvMoi, phien: &phienGiaKhuVuc{}})
	if conSong := k.AnhChupNguoiChoi(nvMoi.ID); len(conSong) != 0 {
		t.Fatalf("TatMayChu gap loi DongVaChoXong nhung KHONG dung actor (van xu ly duoc su kien moi sau do) — " +
			"phai LUON goi DungTatCa() du buoc luu co loi hay khong")
	}
}
