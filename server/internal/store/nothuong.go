package store

import (
	"context"
	"fmt"
)

// LoaiThuongKinhNghiem la loai no DUY NHAT hien co. Phai khop y het danh sach
// trong rang buoc thuong_cho_nhan_loai_hop_le (migrations/0004) — them loai moi
// thi sua ca hai cho VA them nhanh tra tuong ung trong NhanNoThuong, neu khong
// khoan no do se nam lai vinh vien khong ai tra.
const LoaiThuongKinhNghiem = "kinh_nghiem"

// tranKinhNghiem la muc bao hoa khi cong thuong — bang math.MaxInt64, ghim
// thang con so thay vi import math (cung le voi hpToiDa trong game/boss.go:
// test doc thang, khong phai doi chieu hang so cua thu vien).
const tranKinhNghiem int64 = 9_223_372_036_854_775_807

// NoThuong la MOT khoan thuong dang cho tra cho mot nhan vat.
//
// Khong co truong ThoiDiem: moc_tao do CSDL tu dat bang now() — cung ly do voi
// dang_nhap_cuoi trong CapNhatNhanVat, tranh lech dong ho giua cac tien trinh.
type NoThuong struct {
	NhanVatID int64
	Loai      string
	SoDiem    int64
}

// GhiNoThuong ghi mot LO khoan no xuong CSDL trong MOT cau lenh.
//
// Ca lo an hoac trot cung nhau: mot dong vi pham rang buoc (so diem <= 0, loai
// la, nhan vat khong ton tai) lam ca lo bi tu choi. Day la lua chon co y —
// ghi mot phan roi bao loi se de lai mot trang thai khong ai biet la da ghi
// duoc bao nhieu, va nguoi goi (hang doi trong game/sonothuong.go) se thu lai
// ca lo, sinh ra thuong nhan doi.
//
// Lo RONG khong lam gi va khong loi: nguoi goi la mot hang doi chay theo nhip,
// phan lon cac lan xa hang doi deu rong.
func (s *Store) GhiNoThuong(ctx context.Context, ds []NoThuong) error {
	if len(ds) == 0 {
		return nil
	}

	// Ba mang song song + unnest: mot vong mang qua day thay vi len(ds) lan
	// khu hoi CSDL. Mot dot boss dong nguoi co the sinh hang chuc khoan no
	// cung luc.
	idNhanVat := make([]int64, len(ds))
	loai := make([]string, len(ds))
	soDiem := make([]int64, len(ds))
	for i, no := range ds {
		idNhanVat[i] = no.NhanVatID
		loai[i] = no.Loai
		soDiem[i] = no.SoDiem
	}

	if _, err := s.pool.Exec(ctx,
		`INSERT INTO thuong_cho_nhan (nhan_vat_id, loai_thuong, so_diem)
		 SELECT * FROM unnest($1::bigint[], $2::text[], $3::bigint[])`,
		idNhanVat, loai, soDiem,
	); err != nil {
		return fmt.Errorf("ghi no thuong that bai: %w", err)
	}
	return nil
}

// NhanNoThuong tra HET no kinh nghiem dang cho cua mot nhan vat: cong vao cot
// kinh_nghiem, danh dau cac dong no la da nhan, va tra ve tong vua cong (0 neu
// khong co no).
//
// TAT CA TRONG MOT GIAO DICH, va day la diem quan trong nhat cua ca ham: neu
// buoc danh dau va buoc cong nam o hai giao dich khac nhau thi mot cu sap giua
// hai buoc se hoac AN MAT phan thuong (danh dau xong, chua cong) hoac PHAT HAI
// LAN (cong xong, chua danh dau — lan dang nhap sau cong tiep). Goi ham nay
// hai lan lien tiep phai ra 500 roi 0, khong bao gio 500 roi 500.
//
// PHAI GOI TRUOC KHI NAP NHAN VAT VAO BO NHO, hoac phai cong ket qua vao ban
// trong bo nho ngay sau do (login.go lam viec thu hai). Ly do: ThoGhi ghi DE CA
// DONG — mot ban NhanVat trong RAM con mang kinh_nghiem cu se xoa sach phan vua
// cong o lan luu dinh ky dau tien. Day dung la cai bay da khien khoan thuong
// nay khong tra duoc ngay tu dau; no van rinh o day.
//
// Loc theo loai_thuong = kinh_nghiem (khong lay het) de mot loai no moi duoc
// them vao CSDL trong tuong lai KHONG bi ham nay danh dau "da nhan" roi vut di
// ma khong tra gi ca. Tha de no nam lai cho toi khi co nguoi viet duong tra.
func (s *Store) NhanNoThuong(ctx context.Context, nhanVatID int64) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("mo giao dich nhan no that bai: %w", err)
	}
	// Rollback sau Commit la mot phep khong-lam-gi trong pgx, nen defer nay an
	// toan tuyet doi o ca hai duong ra.
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx,
		`UPDATE thuong_cho_nhan
		    SET da_nhan = true, nhan_luc = now()
		  WHERE nhan_vat_id = $1 AND da_nhan = false AND loai_thuong = $2
		 RETURNING so_diem`,
		nhanVatID, LoaiThuongKinhNghiem,
	)
	if err != nil {
		return 0, fmt.Errorf("danh dau no da nhan that bai: %w", err)
	}
	var tong int64
	for rows.Next() {
		var diem int64
		if err := rows.Scan(&diem); err != nil {
			rows.Close()
			return 0, fmt.Errorf("doc so diem no that bai: %w", err)
		}
		tong = congBaoHoaKinhNghiem(tong, diem)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("duyet no that bai: %w", err)
	}

	if tong == 0 {
		// Khong co no: dung o day, KHONG cham vao bang nhan_vat. Duong nay chay
		// o moi lan dang nhap nen phai re nhat co the.
		if err := tx.Commit(ctx); err != nil {
			return 0, fmt.Errorf("chot giao dich nhan no that bai: %w", err)
		}
		return 0, nil
	}

	// LEAST(...::numeric, tran) — bao hoa thay vi tran. Cong tren kieu numeric
	// (khong gioi han chu so) roi moi kep ve bigint: neu cong thang tren bigint
	// thi PostgreSQL nem loi 22003 va HUY ca giao dich, nghia la nguoi choi gan
	// tran kinh nghiem se khong bao gio nhan duoc no nao nua.
	ct, err := tx.Exec(ctx,
		`UPDATE nhan_vat
		    SET kinh_nghiem = LEAST(kinh_nghiem::numeric + $1::numeric, $2::numeric)::bigint
		  WHERE id = $3`,
		tong, tranKinhNghiem, nhanVatID,
	)
	if err != nil {
		return 0, fmt.Errorf("cong thuong vao nhan vat that bai: %w", err)
	}
	if ct.RowsAffected() == 0 {
		// Khoa ngoai khong cho no ton tai ma nhan vat thi khong — nen den day
		// duoc la du lieu da hong. Tra loi de giao dich bi huy: tha khong tra
		// thuong con hon danh dau da tra ma khong cong cho ai.
		return 0, ErrKhongThayNhanVat
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("chot giao dich nhan no that bai: %w", err)
	}
	return tong, nil
}

// congBaoHoaKinhNghiem cong hai so khong am, dung o tran thay vi cuon vong ve
// so am. Tong cac khoan no cua mot nguoi ve ly thuyet co the vuot int64 (moi
// khoan da co the toi MaxInt64 sau bao hoa o tang game).
func congBaoHoaKinhNghiem(a, b int64) int64 {
	if b > 0 && a > tranKinhNghiem-b {
		return tranKinhNghiem
	}
	return a + b
}
