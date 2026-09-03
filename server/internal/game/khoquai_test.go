// khoquai_test.go — quai doc tu KHO NOI DUNG (KH8 task 8).
//
// Bon cau hoi tep nay tra loi, moi cau la mot dot bien co that:
//
//  1. Doi `sat_thuong_co_ban` trong YAML co that su doi sat thuong quai gay ra
//     khong? (dot bien: tra hang so 5 ve quai.go)
//  2. Sinh quai co that su doc anh chup khong? (dot bien: doc lai mot bang cung)
//  3. `anh_id` co that su duoc doi chieu voi catalog anh DA KY khong?
//  4. Duong sinh quai co lam hong hop dong "mot lan dau nhip" cua task 3 khong?
package game

import (
	"bytes"
	"testing"
	"time"

	"nroserver/internal/noidung"
)

// khoQuaiSuaSatThuong dung mot *noidung.Kho THAT (ban nhung) roi doi DUY NHAT
// truong SatThuongCoBan cua mau quai dau. Khong bia du lieu: moi so khac van la
// so that, nen neu test do thi do dung mot nguyen nhan.
func khoQuaiSuaSatThuong(t *testing.T, moi int32) *noidung.Kho {
	t.Helper()
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	if len(du.MauQuai) == 0 {
		t.Fatal("kho nhung khong co mau quai nao")
	}
	du.MauQuai[0].SatThuongCoBan = moi
	return noidung.MoKho(du)
}

// TestQuaiSatThuongDenTuYAMLKhongPhaiHangSo la test TRUNG TAM cua task 8.
//
// DOT BIEN NHAM TOI: tra `SatThuongCoBan: 5` (hay bat ky hang so nao) vao duong
// sinh quai, thay vi doc tu mau. Voi dot bien do, sat thuong trong goi -11 se
// van la 5 du YAML noi 41 — va test nay do.
//
// Dung 41 chu khong phai mot so tron: mot so tron de trung ngau nhien voi mot
// gia tri mac dinh nao do va che mat loi.
func TestQuaiSatThuongDenTuYAMLKhongPhaiHangSo(t *testing.T) {
	const satThuongMoi int32 = 41

	t0 := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.nguoiTrong[1].nv.HpHienTai = 1000 // du mau de sat thuong khong bi kep §1.2
	k.nguoiTrong[1].nv.ViTriX = 0
	k.nguoiTrong[1].nv.ViTriY = 0

	// Sinh quai TU ANH CHUP da sua — dung duong ma MoKhuVuc di.
	kho := khoQuaiSuaSatThuong(t, satThuongMoi)
	k.quai = taoQuaiTuAnhChup(kho.AnhChup(), TimBanDo(1))
	if len(k.quai) == 0 {
		t.Fatal("khong sinh duoc con quai nao tu anh chup")
	}
	// Keo con quai dau ve dung cho nguoi choi (tam tang 2 = 100) va cho no moc
	// hoi mau hop le, y het cac test aiQuai khac trong goi nay.
	k.quai = k.quai[:1]
	k.quai[0].X, k.quai[0].Y = 0, 0
	k.quai[0].LanHoiMauCuoi = t0

	luc = t0.Add(100 * time.Millisecond)
	k.nhipTim()

	if len(p.nhan) != 1 || p.nhan[0].Command != -11 {
		t.Fatalf("mong doi DUNG mot goi -11, nhan %+v", p.nhan)
	}
	muon := []byte{
		0,           // mobId
		0, 0, 0, 41, // satThuong = sat_thuong_co_ban tu YAML
	}
	if !bytes.Equal(p.nhan[0].Data, muon) {
		t.Fatalf("sat thuong quai KHONG theo YAML:\n got  = %v\n muon = %v\n"+
			"(neu got la 5 thi hang so cu da quay lai duong sinh quai)", p.nhan[0].Data, muon)
	}
}

// TestQuaiHpDenTuYAML: hp cua mau cung phai di tu YAML vao ca Hp lan HpToiDa.
func TestQuaiHpDenTuYAML(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	du.MauQuai[0].Hp = 777
	kho := noidung.MoKho(du)

	quai := taoQuaiTuAnhChup(kho.AnhChup(), TimBanDo(1))
	if len(quai) == 0 {
		t.Fatal("khong sinh duoc con quai nao")
	}
	for i, q := range quai {
		if q.Hp != 777 || q.HpToiDa != 777 {
			t.Fatalf("quai[%d]: Hp=%d HpToiDa=%d, mong doi ca hai = 777 (hp tu YAML)", i, q.Hp, q.HpToiDa)
		}
	}
}

// TestSinhQuaiDocAnhChupChuKhongPhaiBangCung — DOT BIEN NHAM TOI: ai do dat lai
// mot bang quai cung trong Go va cho duong sinh doc no. Voi mot anh chup KHONG
// CO diem sinh nao, duong sinh phai tra ve RONG. Neu no van tra ra hai con Soi
// Rung thi no dang doc mot nguon khac — chinh la cai bang cung ma task 8 xoa.
func TestSinhQuaiDocAnhChupChuKhongPhaiBangCung(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	du.SinhQuai = nil
	kho := noidung.MoKho(du)

	if quai := taoQuaiTuAnhChup(kho.AnhChup(), TimBanDo(1)); len(quai) != 0 {
		t.Fatalf("anh chup khong co diem sinh nao ma van sinh ra %d con quai — "+
			"duong sinh dang doc mot bang cung nao do, khong doc anh chup", len(quai))
	}
}

// TestSinhQuaiMapKhongKhaiThiRong: mot ban do khong co khoi `sinh:` -> khong
// quai, KHONG panic. (Map 99 khong ton tai trong quai.yaml.)
func TestSinhQuaiMapKhongKhaiThiRong(t *testing.T) {
	bd := *TimBanDo(1)
	bd.MapID = 99
	if quai := taoQuaiTuAnhChup(khoChieuMacDinh(), &bd); len(quai) != 0 {
		t.Fatalf("map 99 khong khai `sinh:` ma sinh ra %d con quai", len(quai))
	}
}

// TestAnhIdQuaiPhaiCoTrongCatalogDaKy — DOT BIEN NHAM TOI: doi `anh_id` trong
// quai.yaml sang mot id khong co trong goi asset da ky. Kiem o day dung DANH
// SACH THAT (catalogLangLa.MobIDs), khac voi test thuan trong goi noidung.
func TestAnhIdQuaiPhaiCoTrongCatalogDaKy(t *testing.T) {
	ac := khoChieuMacDinh()
	if err := noidung.KiemIdAnhQuai(ac.DanhSachQuai(), catalogLangLa.MobIDs); err != nil {
		t.Fatalf("quai.yaml lech catalog anh da ky: %v", err)
	}
	// Va bo kiem phai THAT SU tu choi mot id la — neu khong, cau tren xanh vo
	// nghia (vd MobIDs rong thi moi thu deu "hop le").
	la := ac.DanhSachQuai()
	la[0].AnhId = 123456
	if err := noidung.KiemIdAnhQuai(la, catalogLangLa.MobIDs); err == nil {
		t.Fatal("anh_id 123456 duoc chap nhan — danh sach id anh da ky co the dang RONG")
	}
}

// TestQuaiDangSongGiuSoCuaTheHeLucSinh khoa quyet dinh "cai gi khong nap nong
// duoc" (thiet ke §4, khoi chu thich dau khoquai.go): nap nong quai.yaml KHONG
// sua nhung con quai dang dung san tren ban do.
//
// Day KHONG phai mot thieu sot duoc dung tha: no la hanh vi CO CHU DICH (HpToiDa
// da gui xuong client trong goi -24 luc nguoi choi vao ban do). Test nay ton tai
// de mot ngay nao do co nguoi doi y thi ho phai XOA no mot cach co y thuc, chu
// khong lam vo no mot cach tinh co.
func TestQuaiDangSongGiuSoCuaTheHeLucSinh(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	kho := noidung.MoKho(du)

	quai := taoQuaiTuAnhChup(kho.AnhChup(), TimBanDo(1))
	if len(quai) == 0 {
		t.Fatal("khong sinh duoc con quai nao")
	}
	satThuongLucSinh := quai[0].SatThuongCoBan
	hpLucSinh := quai[0].HpToiDa

	// NAP NONG: mot bo so hoan toan khac.
	du2, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap lan hai: %v", err)
	}
	du2.MauQuai[0].SatThuongCoBan = satThuongLucSinh + 100
	du2.MauQuai[0].Hp = hpLucSinh + 100
	kho.Thay(du2)

	if quai[0].SatThuongCoBan != satThuongLucSinh || quai[0].HpToiDa != hpLucSinh {
		t.Fatalf("quai DANG SONG doi so sau nap nong: satThuong %d -> %d, hpToiDa %d -> %d "+
			"(hanh vi da chot la GIU nguyen — xem khoi chu thich dau khoquai.go)",
			satThuongLucSinh, quai[0].SatThuongCoBan, hpLucSinh, quai[0].HpToiDa)
	}

	// Nhung quai SINH RA SAU nap nong thi phai lay so MOI — neu khong thi nap
	// nong hoan toan vo tac dung voi quai, va do moi la hong that.
	moi := taoQuaiTuAnhChup(kho.AnhChup(), TimBanDo(1))
	if len(moi) == 0 || moi[0].SatThuongCoBan != satThuongLucSinh+100 {
		t.Fatalf("quai sinh SAU nap nong khong lay so moi: %+v", moi)
	}
}

// TestSinhQuaiKhongLamHongHopDongMotLanDauNhip: duong sinh quai KHONG duoc goi
// AnhChup() them lan nao trong nhip. MoKhuVuc sinh quai TRUOC khi actor chay
// (va doc ban NHUNG, khong qua nguon da gan), nen so lan dem cua nguon phai
// khong doi qua mot nhip — bo sung cho TestKhoChieuMotLanMoiNhip (task 3) bang
// mot khu vuc CO QUAI THAT.
func TestSinhQuaiKhongLamHongHopDongMotLanDauNhip(t *testing.T) {
	d := moKhoChieuDem(t)
	k := MoKhuVuc(TimBanDo(1))
	k.ganNguonNoiDung(d.nguon())
	if len(k.quai) == 0 {
		t.Fatal("chuan bi test sai: khu vuc khong co quai nao de nhipTim cham toi")
	}

	d.soLan = 0
	k.nhipTim()
	if d.soLan != 1 {
		t.Fatalf("mot nhip (co %d quai) lay anh chup %d lan, hop dong la DUNG 1 lan",
			len(k.quai), d.soLan)
	}
}
