package game

import (
	"context"
	"testing"

	"nroserver/internal/auth"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

type khoGia struct {
	tk  *store.TaiKhoan
	err error
}

func (k *khoGia) LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) {
	if k.err != nil {
		return nil, k.err
	}
	return k.tk, nil
}

// dungGoiDangNhap dung goi dang nhap dung y het client gui:
// khung lenh -29, ben trong la byte lenh con 0 roi hai chuoi.
func dungGoiDangNhap(ten, matKhau string) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(0) // lenh con: dang nhap
	w.UTF(ten)
	w.UTF(matKhau)
	return protocol.NewMessage(-29, w)
}

// Goi dung dinh dang phai duoc boc tach ra dung ten va mat khau.
func TestBocTachGoiDangNhap(t *testing.T) {
	m := dungGoiDangNhap("nguoi", "matkhau")
	r := protocol.NewReader(m.Data)

	lenhCon, err := r.Byte()
	if err != nil || lenhCon != 0 {
		t.Fatalf("lenh con = %d, err = %v, mong doi 0", lenhCon, err)
	}
	ten, _ := r.UTF()
	mk, _ := r.UTF()
	if ten != "nguoi" || mk != "matkhau" {
		t.Fatalf("boc tach ra %q / %q, mong doi \"nguoi\" / \"matkhau\"", ten, mk)
	}
}

func TestDangNhapDungMatKhau(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, TenDangNhap: "nguoi", MatKhauBam: bam}})

	kq, err := b.KiemTraDangNhap(context.Background(), "nguoi", "matkhau")
	if err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	if !kq.ThanhCong {
		t.Fatalf("phai dang nhap duoc, ly do tu choi: %q", kq.LyDo)
	}
	if kq.TaiKhoan.ID != 7 {
		t.Fatalf("ID = %d, mong doi 7", kq.TaiKhoan.ID)
	}
}

func TestDangNhapSaiMatKhau(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "sai")
	if kq.ThanhCong {
		t.Fatal("mat khau sai ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, mong doi %q", kq.LyDo, LyDoSaiThongTin)
	}
}

// Tai khoan khong ton tai va sai mat khau phai tra ve CUNG MOT cau bao,
// de ke xau khong do duoc tai khoan nao co that.
func TestTaiKhoanKhongTonTaiBaoGiongSaiMatKhau(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: store.ErrKhongThayTaiKhoan})

	kq, _ := b.KiemTraDangNhap(context.Background(), "khong-co", "gi-do")
	if kq.ThanhCong {
		t.Fatal("tai khoan khong ton tai ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, phai giong het truong hop sai mat khau", kq.LyDo)
	}
}

func TestTaiKhoanBiKhoa(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam, BiKhoa: true}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "matkhau")
	if kq.ThanhCong {
		t.Fatal("tai khoan bi khoa ma van cho vao")
	}
	if kq.LyDo != LyDoBiKhoa {
		t.Fatalf("ly do = %q, mong doi %q", kq.LyDo, LyDoBiKhoa)
	}
}

// Ten dang nhap rong khong duoc cham toi CSDL.
func TestTenRongBiTuChoiNgay(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: context.Canceled}) // kho se loi neu bi goi
	kq, err := b.KiemTraDangNhap(context.Background(), "", "gi-do")
	if err != nil {
		t.Fatalf("khong duoc goi toi kho: %v", err)
	}
	if kq.ThanhCong {
		t.Fatal("ten rong ma van cho vao")
	}
}
