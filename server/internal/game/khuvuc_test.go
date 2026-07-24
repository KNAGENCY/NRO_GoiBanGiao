package game

import (
	"testing"

	"nroserver/internal/store"
)

func TestMoKhuVucGiuBanDo(t *testing.T) {
	bd := TimBanDo(1)
	k := MoKhuVuc(bd)
	if k.BanDo() != bd {
		t.Error("MoKhuVuc phai giu dung con tro ban do da cho")
	}
	if k.SoNguoiChoi() != 0 {
		t.Errorf("khu vuc moi mo co %d nguoi, mong doi 0", k.SoNguoiChoi())
	}
}

func TestThemNguoiChoiVaoKhuVuc(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.ThemNguoiChoi(&store.NhanVat{ID: 1, Ten: "Naruto"})
	if k.SoNguoiChoi() != 1 {
		t.Fatalf("sau khi them 1 nguoi, SoNguoiChoi = %d, mong doi 1", k.SoNguoiChoi())
	}
	k.ThemNguoiChoi(&store.NhanVat{ID: 2, Ten: "Sasuke"})
	if k.SoNguoiChoi() != 2 {
		t.Errorf("sau khi them 2 nguoi, SoNguoiChoi = %d, mong doi 2", k.SoNguoiChoi())
	}
}

func TestKhuVucActorKhoiDongVaDungSach(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	k.Gui(suKienRa{idNhanVat: 999}) // no-op: khong co ai
	k.Dung()                        // phai tra ve, khong treo
	// gui sau khi dung khong duoc panic:
	k.Gui(suKienRa{idNhanVat: 1})
}
