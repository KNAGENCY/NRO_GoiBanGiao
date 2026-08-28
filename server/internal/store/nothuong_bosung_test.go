package store

import (
	"context"
	"strings"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Bo sung sau mot vong KIEM THU DOT BIEN doi khang.
//
// Bo test dau tien khoa chac moi duong THANH CONG, nhung de trong gan nhu toan
// bo duong LOI, duong PHONG THU, va cac bat bien chi duoc ghi bang chu thich —
// nen nhieu dot bien nguy hiem song sot. Moi test o day sinh ra tu MOT dot bien
// cu the da song sot, va chu thich noi ro do la dot bien nao.
// ─────────────────────────────────────────────────────────────────────────────

// TestGhiNoThuong_GhiLaiCungLoKhongNhanDoi khoa duong chong NHAN DOI phan
// thuong — kich ban that: PostgreSQL commit xong roi moi tra loi cho client
// (het han gio, dut ket noi ngay sau commit). Ben goi khong phan biet duoc
// "chua ghi" voi "ghi roi ma khong biet" nen no thu lai ca lo. Khong co gi chan
// thi nguoi choi nhan thuong gap doi.
func TestGhiNoThuong_GhiLaiCungLoKhongNhanDoi(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notesthh")

	lo := []NoThuong{{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 250, KhoaLo: "lo-lap-lai"}}
	for lan := 1; lan <= 3; lan++ {
		if err := s.GhiNoThuong(ctx, lo); err != nil {
			t.Fatalf("ghi lo lan %d: %v", lan, err)
		}
	}

	if got := demNoChuaNhan(t, s, nv.ID); got != 1 {
		t.Fatalf("ghi CUNG mot lo 3 lan sinh ra %d dong no, mong doi 1", got)
	}
	tong, _, err := s.NhanNoThuong(ctx, nv.ID)
	if err != nil {
		t.Fatalf("nhan no: %v", err)
	}
	if tong != 250 {
		t.Fatalf("nhan duoc %d diem, mong doi 250 (khong nhan doi)", tong)
	}
}

// TestGhiNoThuong_HaiLoKhacNhauDeuVaoDuoc: chan nhan doi khong duoc chan nham
// hai khoan THAT cua cung mot nguoi den tu hai lo khac nhau (vd hai lan boss
// chet cach nhau vai giay).
func TestGhiNoThuong_HaiLoKhacNhauDeuVaoDuoc(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestii")

	for _, khoa := range []string{"lo-mot", "lo-hai"} {
		if err := s.GhiNoThuong(ctx, []NoThuong{
			{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 60, KhoaLo: khoa},
		}); err != nil {
			t.Fatalf("ghi lo %s: %v", khoa, err)
		}
	}

	if got := demNoChuaNhan(t, s, nv.ID); got != 2 {
		t.Fatalf("hai lo khac nhau sinh ra %d dong, mong doi 2", got)
	}
	tong, _, err := s.NhanNoThuong(ctx, nv.ID)
	if err != nil {
		t.Fatalf("nhan no: %v", err)
	}
	if tong != 120 {
		t.Fatalf("nhan duoc %d diem, mong doi 120", tong)
	}
}

// TestNhanNoThuong_NhanVatKhongTonTaiThiKhongLoi — sinh tu dot bien "bo nhanh
// tat if tong == 0", da SONG SOT. Bo nhanh do lam ham chay tiep xuong UPDATE
// nhan_vat, khop 0 dong, roi tra ErrKhongThayNhanVat — doi hanh vi that.
//
// Duong nay chay o MOI lan dang nhap, nen mot id khong co trong bang phai ra
// (0, nil): khong co no thi khong co gi de lam, khong phai loi.
func TestNhanNoThuong_NhanVatKhongTonTaiThiKhongLoi(t *testing.T) {
	s := moKhoCoBangNo(t)
	tong, _, err := s.NhanNoThuong(context.Background(), 987654321)
	if err != nil {
		t.Fatalf("nhan no cua nhan vat khong ton tai: %v", err)
	}
	if tong != 0 {
		t.Fatalf("tra %d, mong doi 0", tong)
	}
}

// TestNhanNoThuong_HaiKhoanLonCongLaiKhongTranSoAm — sinh tu dot bien "bo bao
// hoa trong congBaoHoaKinhNghiem", da SONG SOT.
//
// Bai kiem bao hoa cu chi co MOT dong no, nen bo cong don PHIA GO khong bao gio
// tran — thu no thuc su kiem la LEAST cua SQL. Hai dong lon moi cham toi duoc
// bo cong don do. Tran o day cho ra so AM, tuc nguoi choi dang nhap xong bi TRU
// kinh nghiem.
func TestNhanNoThuong_HaiKhoanLonCongLaiKhongTranSoAm(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestjj")

	const ratLon = int64(5_000_000_000_000_000_000) // hai cai nay cong lai vuot MaxInt64
	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: ratLon, KhoaLo: "lo-j1"},
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: ratLon, KhoaLo: "lo-j2"},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}

	const maxInt64 = int64(9_223_372_036_854_775_807)
	tong, _, err := s.NhanNoThuong(ctx, nv.ID)
	if err != nil {
		t.Fatalf("nhan no: %v", err)
	}
	if tong != maxInt64 {
		t.Fatalf("NhanNoThuong tra %d, mong doi bao hoa o %d (khong duoc cuon ve so am)", tong, maxInt64)
	}
	if got := docKinhNghiem(t, s, nv.ID); got != maxInt64 {
		t.Fatalf("kinh_nghiem = %d, mong doi %d", got, maxInt64)
	}
}

// TestBang_TuChoiLoaiThuongLa — sinh tu dot bien "xoa CHECK loai_thuong", da
// SONG SOT. Rang buoc nay quan trong hon ve ngoai cua no: no la thu DUY NHAT
// dang giu cho menh de loc trong NhanNoThuong tuong duong voi "lay het". Xoa no
// ma quen sua ham thi no loai moi nam lai vinh vien khong ai tra — im lang
// tuyet doi.
func TestBang_TuChoiLoaiThuongLa(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestkk")

	if _, err := s.pool.Exec(ctx,
		`INSERT INTO thuong_cho_nhan (nhan_vat_id, loai_thuong, so_diem, khoa_lo)
		 VALUES ($1, 'vang', 10, 'lo-k1')`, nv.ID); err == nil {
		t.Fatal("CSDL nhan loai_thuong la — rang buoc thuong_cho_nhan_loai_hop_le da mat")
	}
}

// TestBang_TuChoiDaNhanMaKhongCoMocNhan — sinh tu dot bien "xoa CHECK khop
// da_nhan/nhan_luc", da SONG SOT. Mat rang buoc nay thi so sach "ai nhan gi luc
// nao" thanh vo nghia ma khong bao loi o dau ca.
func TestBang_TuChoiDaNhanMaKhongCoMocNhan(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestll")

	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 5, KhoaLo: "lo-l1"},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}
	if _, err := s.pool.Exec(ctx,
		`UPDATE thuong_cho_nhan SET da_nhan = true WHERE nhan_vat_id = $1`, nv.ID); err == nil {
		t.Fatal("CSDL cho danh dau da nhan ma khong co moc nhan")
	}
}

// TestBang_CoIndexMotPhanChuaNhan — sinh tu dot bien "xoa index mot phan", da
// SONG SOT. Day thuan tuy la hieu nang nen khong test HANH VI nao bat duoc;
// khoa thang hop dong luoc do thay vi de no troi. Mat index mot phan thi moi
// lan dang nhap phai quet ca bang chi-tang, cham dan theo doi may chu.
func TestBang_CoIndexMotPhanChuaNhan(t *testing.T) {
	s := moKhoCoBangNo(t)
	var dinhNghia string
	if err := s.pool.QueryRow(context.Background(),
		`SELECT indexdef FROM pg_indexes
		  WHERE tablename = 'thuong_cho_nhan' AND indexname = 'thuong_cho_nhan_chua_nhan'`,
	).Scan(&dinhNghia); err != nil {
		t.Fatalf("khong thay index thuong_cho_nhan_chua_nhan: %v", err)
	}
	if !strings.Contains(dinhNghia, "da_nhan = false") {
		t.Fatalf("index khong con la index MOT PHAN: %s", dinhNghia)
	}
}

// TestNhanNoThuong_KhongRoKetNoiKhiGiuaChung_ThatBai — sinh tu dot bien "bo
// defer tx.Rollback", da SONG SOT vi khong duong LOI nao cua NhanNoThuong duoc
// test cham toi.
//
// pool.Begin muon MOT ket noi khoi pool; thoat ma khong Commit/Rollback la
// khong bao gio tra lai. Het pool = ca may chu khong ai dang nhap duoc nua.
// Dung mot trigger tam lam UPDATE nhan_vat no giua giao dich.
func TestNhanNoThuong_KhongRoKetNoiKhiGiuaChungThatBai(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestmm")

	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 15, KhoaLo: "lo-m1"},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}

	// Trigger lam moi UPDATE len nhan_vat nay that bai — tuc la giao dich cua
	// NhanNoThuong se no SAU khi da danh dau no.
	if _, err := s.pool.Exec(ctx, `
		CREATE OR REPLACE FUNCTION nro_test_chan_cap_nhat() RETURNS trigger AS $$
		BEGIN RAISE EXCEPTION 'chan co y de kiem'; END; $$ LANGUAGE plpgsql;
		CREATE TRIGGER nro_test_chan BEFORE UPDATE ON nhan_vat
		    FOR EACH ROW EXECUTE FUNCTION nro_test_chan_cap_nhat();`); err != nil {
		t.Fatalf("tao trigger thu: %v", err)
	}
	t.Cleanup(func() {
		s.pool.Exec(context.Background(), `DROP TRIGGER IF EXISTS nro_test_chan ON nhan_vat`)
		s.pool.Exec(context.Background(), `DROP FUNCTION IF EXISTS nro_test_chan_cap_nhat()`)
	})

	truoc := s.pool.Stat().AcquiredConns()
	if _, _, err := s.NhanNoThuong(ctx, nv.ID); err == nil {
		t.Fatal("NhanNoThuong khong bao loi du UPDATE nhan_vat bi trigger chan")
	}
	if sau := s.pool.Stat().AcquiredConns(); sau > truoc {
		t.Fatalf("con %d ket noi bi giu sau khi giao dich hong (truoc %d) — thieu Rollback, "+
			"het pool la ca may chu khong ai dang nhap duoc", sau, truoc)
	}

	// Va no PHAI con nguyen (giao dich bi huy tron ven), de lan sau tra duoc.
	if got := demNoChuaNhan(t, s, nv.ID); got != 1 {
		t.Fatalf("con %d dong no chua nhan sau khi giao dich hong, mong doi 1 "+
			"(giao dich phai bi huy tron ven, khong duoc danh dau da nhan)", got)
	}
}
