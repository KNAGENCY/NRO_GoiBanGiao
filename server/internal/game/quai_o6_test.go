package game

import (
	"testing"

	"nroserver/internal/protocol"
)

// TestGhiQuaiO6LaMauKhongPhaiId tai hien loi 05/09/2026: o thu 6 cua moi con quai
// trong -24 la templateId (client: Mob.arrMobTemplate[o6]) — truoc do ghi Id (id
// tam 0,1,...). Voi 2 con soi cung mau 0, con thu hai co Id=1 -> client tra mau 1
// khong ton tai -> IndexOutOfRange -> bo NGUYEN phan NPC/do roi/bong do cua -24.
// Dot bien "ghi lai q.Id" -> test nay do.
func TestGhiQuaiO6LaMauKhongPhaiId(t *testing.T) {
	q := &Quai{Id: 9, MauId: 3, Cap: 1, Hp: 10, HpToiDa: 10}
	w := protocol.NewWriter()
	ghiQuai(w, q)
	b := w.Bytes()
	if b[5] != 3 {
		t.Fatalf("o 6 = %d, mong doi MauId 3 (khong phai Id 9)", b[5])
	}
}

// TestGhiDanhSachQuaiThuTuIdLaChiSo: client lay mobId = CHI SO trong danh sach
// (Controller.cs:5068 `new Mob(b, ...)`), nen danh sach quai cua moi khu vuc phai
// xep dung Id 0,1,2,... — neu khong, goi -9/-10/-11/-12 mang Id may chu se tro
// nhầm con. Kiem tren bo quai that cua map 1 sinh tu kho noi dung.
func TestGhiDanhSachQuaiThuTuIdLaChiSo(t *testing.T) {
	quai := taoQuaiTuAnhChup(khoChieuMacDinh(), TimBanDo(1))
	for i, q := range quai {
		if int(q.Id) != i {
			t.Fatalf("quai thu %d co Id %d — client se gan mobId=%d, lech voi may chu", i, q.Id, i)
		}
	}
	// Va o 6 cua moi con phai la mau co that trong catalog -28/6: goiCatalogMap
	// (dongbo.go) ghi DUNG MOT mau quai (catalogLangLa.Mob) -> client
	// arrMobTemplate dai 1 -> chi templateId 0 hop le.
	soMau := 1
	w := protocol.NewWriter()
	ghiDanhSachQuai(w, quai)
	b := w.Bytes()
	off := 1
	for i := range quai {
		tpl := b[off+5]
		if int(tpl) >= soMau {
			t.Fatalf("quai thu %d ghi templateId %d >= so mau %d — client IndexOutOfRange, bo phan con lai cua -24", i, tpl, soMau)
		}
		off += 23
	}
}
