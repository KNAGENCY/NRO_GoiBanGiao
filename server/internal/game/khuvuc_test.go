package game

import (
	"testing"

	"nroserver/internal/protocol"
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

// phienGiaKhuVuc la phien mang gia (thoa man interface guiTin) dung rieng cho
// test xuLy cua KhuVuc: chi ghi lai moi goi tin nhan duoc vao slice, khong gui
// di dau ca. Ten khac voi ham `phienGia` (nhanvat_test.go, phien mynet THAT
// qua net.Pipe) de tranh dung ten trong cung package.
type phienGiaKhuVuc struct{ nhan []*protocol.Message }

func (p *phienGiaKhuVuc) Gui(m *protocol.Message) error {
	p.nhan = append(p.nhan, m)
	return nil
}

// TestKhuVucXuLySuKien kiem tra logic xuLy cho ca 3 su kien theo dung THU TU
// phat tin cua Ke hoach 5 Task 5. Goi k.xuLy(sk) TRUC TIEP (actor CHUA BatDau)
// de tach logic xu ly khoi vong lap channel (da test o Task 1).
func TestKhuVucXuLySuKien(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))

	nvA := &store.NhanVat{ID: 1, Ten: "A", HpToiDa: 100}
	nvB := &store.NhanVat{ID: 2, Ten: "B", HpToiDa: 100}
	pA := &phienGiaKhuVuc{}
	pB := &phienGiaKhuVuc{}

	// Buoc 1: A vao khu vuc — chua co ai khac trong khu vuc nen pA khong
	// nhan goi nao (khong tu thay minh).
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	if len(pA.nhan) != 0 {
		t.Fatalf("buoc 1: pA nhan %d goi, mong doi 0", len(pA.nhan))
	}
	if _, co := k.nguoiTrong[nvA.ID]; !co {
		t.Fatal("buoc 1: nvA phai co trong khu vuc sau khi vao")
	}

	// Buoc 2: B vao khu vuc — pB phai nhan DUNG 1 goi -5 cua A; pA phai
	// nhan DUNG 1 goi -5 cua B (moi ben chi thay nguoi kia, khong thay minh).
	k.xuLy(suKienVao{nv: nvB, phien: pB})
	if len(pB.nhan) != 1 {
		t.Fatalf("buoc 2: pB nhan %d goi, mong doi 1", len(pB.nhan))
	}
	if pB.nhan[0].Command != -5 {
		t.Fatalf("buoc 2: pB goi lenh %d, mong doi -5", pB.nhan[0].Command)
	}
	khangDinhIdTrongGoi(t, pB.nhan[0], nvA.ID)

	if len(pA.nhan) != 1 {
		t.Fatalf("buoc 2: pA nhan %d goi, mong doi 1", len(pA.nhan))
	}
	if pA.nhan[0].Command != -5 {
		t.Fatalf("buoc 2: pA goi lenh %d, mong doi -5", pA.nhan[0].Command)
	}
	khangDinhIdTrongGoi(t, pA.nhan[0], nvB.ID)

	// Buoc 3: A di chuyen toi (250,300) — pB nhan them 1 goi -7 [idA][250][300];
	// pA KHONG nhan gi them (client tu ve chuyen dong cua chinh minh).
	k.xuLy(suKienDiChuyen{idNhanVat: nvA.ID, x: 250, y: 300})
	if len(pA.nhan) != 1 {
		t.Fatalf("buoc 3: pA nhan %d goi, mong doi van la 1 (khong nhan lai chinh minh)", len(pA.nhan))
	}
	if len(pB.nhan) != 2 {
		t.Fatalf("buoc 3: pB nhan %d goi, mong doi 2", len(pB.nhan))
	}
	goiDiChuyenB := pB.nhan[1]
	if goiDiChuyenB.Command != -7 {
		t.Fatalf("buoc 3: pB goi lenh %d, mong doi -7", goiDiChuyenB.Command)
	}
	r := protocol.NewReader(goiDiChuyenB.Data)
	if id, err := r.Int(); err != nil || id != int32(nvA.ID) {
		t.Fatalf("buoc 3: id trong goi -7 = %d, err=%v, mong doi %d", id, err, nvA.ID)
	}
	if x, err := r.Short(); err != nil || x != 250 {
		t.Fatalf("buoc 3: x trong goi -7 = %d, err=%v, mong doi 250", x, err)
	}
	if y, err := r.Short(); err != nil || y != 300 {
		t.Fatalf("buoc 3: y trong goi -7 = %d, err=%v, mong doi 300", y, err)
	}
	if nvA.ViTriX != 250 || nvA.ViTriY != 300 {
		t.Fatalf("buoc 3: nvA vi tri = (%d,%d), mong doi (250,300) — phai la authoritative", nvA.ViTriX, nvA.ViTriY)
	}

	// Buoc 4: A roi khu vuc — pB nhan them 1 goi -6 [idA]; nvA bi xoa khoi map.
	k.xuLy(suKienRa{idNhanVat: nvA.ID})
	if len(pB.nhan) != 3 {
		t.Fatalf("buoc 4: pB nhan %d goi, mong doi 3", len(pB.nhan))
	}
	goiRaB := pB.nhan[2]
	if goiRaB.Command != -6 {
		t.Fatalf("buoc 4: pB goi lenh %d, mong doi -6", goiRaB.Command)
	}
	r2 := protocol.NewReader(goiRaB.Data)
	if id, err := r2.Int(); err != nil || id != int32(nvA.ID) {
		t.Fatalf("buoc 4: id trong goi -6 = %d, err=%v, mong doi %d", id, err, nvA.ID)
	}
	if len(k.nguoiTrong) != 1 {
		t.Fatalf("buoc 4: con %d nguoi trong khu vuc, mong doi 1", len(k.nguoiTrong))
	}
	if _, co := k.nguoiTrong[nvA.ID]; co {
		t.Fatal("buoc 4: nvA phai bi xoa khoi khu vuc sau khi ra")
	}
}

// khangDinhIdTrongGoi doc truong dau tien (playerId, writeInt) cua mot goi -5
// va khang dinh no khop id mong doi.
func khangDinhIdTrongGoi(t *testing.T, m *protocol.Message, idMongDoi int64) {
	t.Helper()
	r := protocol.NewReader(m.Data)
	id, err := r.Int()
	if err != nil {
		t.Fatalf("doc id trong goi -5 loi: %v", err)
	}
	if id != int32(idMongDoi) {
		t.Fatalf("id trong goi -5 = %d, mong doi %d", id, idMongDoi)
	}
}
