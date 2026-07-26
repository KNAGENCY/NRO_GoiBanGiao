package game

import (
	"bytes"
	"testing"
	"time"

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

// TestKhuVucChonChieuCapNhatChieuDangChon: gui suKienChonChieu voi templateId
// HOP LE (thuoc lang cua nhan vat) -> actor phai cap nhat chieuDangChon dung
// CapChieu cap 1 cua chieu do. Kiem tra bang gia tri GHIM CUNG (khong doi
// module gian tiep), tinh tu cong thuc taoCapChieu/xayDungChieuTheoLang trong
// chieuthuc.go: lang=0 (La), templateId=2 (slot 1, tam xa) ->
// skillIdGoc = 1 + (0*4+1)*3 = 4, cap 1 -> SkillId=4, Diem=1.
func TestKhuVucChonChieuCapNhatChieuDangChon(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})

	if k.nguoiTrong[nvA.ID].chieuDangChon != nil {
		t.Fatal("truoc khi chon chieu, chieuDangChon phai la nil")
	}

	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 2})

	cap := k.nguoiTrong[nvA.ID].chieuDangChon
	if cap == nil {
		t.Fatal("sau khi chon chieu HOP LE, chieuDangChon khong duoc la nil")
	}
	if cap.SkillId != 4 || cap.Diem != 1 {
		t.Fatalf("chieuDangChon = {SkillId:%d Diem:%d}, mong doi {SkillId:4 Diem:1} (gia tri ghim cung)",
			cap.SkillId, cap.Diem)
	}
}

// TestKhuVucChonChieuKhongThuocLangKhongDoiGi: nhan vat co Lang khong hop le
// (5, ngoai khoang 0..4) chon templateId=1 (ton tai o MOI lang hop le, nhung
// KHONG co lang nao cho nhan vat nay) -> ChieuCuaLang tra nil -> khong tim
// thay -> chieuDangChon PHAI GIU NGUYEN (khong doi gi), dung dac ta §3.2.
func TestKhuVucChonChieuKhongThuocLangKhongDoiGi(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 5}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})

	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 1})

	if k.nguoiTrong[nvA.ID].chieuDangChon != nil {
		t.Fatal("chon chieu khong thuoc lang: chieuDangChon phai VAN LA nil (khong doi gi)")
	}
}

// TestKhuVucChonChieuNguoiKhongTrongKhuVucKhongPanic: idNhanVat khong co
// trong khu vuc (chua vao, hoac da ra) -> phai la no-op an toan, khong panic.
func TestKhuVucChonChieuNguoiKhongTrongKhuVucKhongPanic(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.xuLy(suKienChonChieu{idNhanVat: 999, mauChieuId: 1})
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 13 — chet va hoi sinh nguoi choi (dac ta §3.4, §6.4).
// nguoiChetHopLeMau/ chieu can chien tam 100 dung chung nhieu test o day —
// khong tach helper rieng nhu chotkiem_test.go vi moi test can mot bo trang
// thai khac nhau (co nguoi thu hai hay khong, quai o dau...).
// ---------------------------------------------------------------------------

// TestSuKienDanhQuaiNguoiDaChetKhongDuocDanh la TEST BAT BUOC cua brief Task
// 13: nguoi choi HpHienTai<=0 (da chet) gui suKienDanhQuai (goi len 54) PHAI
// bi tu choi HOAN TOAN — khong phat goi gi, khong tru mau quai. Day CUNG la
// muc tieu cua kiem thu dot bien (brief buoc 5): bo chot nay o khuvuc.go ->
// test nay phai do (xem bao cao Task 13 phan "Kiem thu dot bien").
func TestSuKienDanhQuaiNguoiDaChetKhongDuocDanh(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 0} // da chet
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	// templateId=1: chieu can chien lang 0 (chieuthuc.go), tam 100 — moi kiem
	// KHAC (chakra/hoi chieu/tam) deu HOP LE, de chung minh chi rieng "da
	// chet" la ly do tu choi (loai tru kha nang bi chan boi ly do khac).
	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 1})
	pA.nhan = nil // don goi (neu co) tu hai buoc chuan bi tren

	k.quai = []*Quai{{Id: 0, Hp: 100, HpToiDa: 150}} // Hp != HpToiDa: tranh nhanh luoi an toan §1.3 xen vao

	k.xuLy(suKienDanhQuai{idNhanVat: nvA.ID, mobId: 0})

	if len(pA.nhan) != 0 {
		t.Fatalf("nguoi choi DA CHET (HpHienTai=0) van gui goi khi danh quai — nhan %d goi, mong doi 0 (test bat buoc cua brief Task 13)", len(pA.nhan))
	}
	if k.quai[0].Hp != 100 {
		t.Fatalf("Hp quai = %d, mong doi VAN LA 100 — don danh cua nguoi da chet phai bi tu choi HOAN TOAN, khong duoc tru mau quai", k.quai[0].Hp)
	}
}

// TestNhipTimAiQuaiHaChetPhatChuoiGoiChetChoCaHaiPhia: quai AI ha CHET mot
// nguoi choi (SatThuongCoBan lon hon HpHienTai con lai, kep boi §1.2) PHAI
// phat DUNG chuoi -11 roi -17 cho nan nhan (dac ta §6.3 roi §6.4 noi tiep
// nhau), VA -10 roi -8 cho nguoi khac trong khu vuc — dung THU TU "dam sat
// thuong truoc, bao chet sau".
func TestNhipTimAiQuaiHaChetPhatChuoiGoiChetChoCaHaiPhia(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, pNanNhan := khuVucNhipTimGia(&luc)
	k.nguoiTrong[1].nv.HpHienTai = 10 // it mau — mot don la chet
	k.nguoiTrong[1].nv.ViTriX = 0
	k.nguoiTrong[1].nv.ViTriY = 0

	pNguoiKhac := &phienGiaKhuVuc{}
	k.nguoiTrong[2] = &nguoiTrongKhuVuc{
		nv:    &store.NhanVat{ID: 2, Ten: "NguoiKhac", HpHienTai: 100, ViTriX: 99_999}, // xa, ngoai CA HAI tam AI -> khong bao gio thanh muc tieu
		phien: pNguoiKhac,
	}

	k.quai = []*Quai{{Id: 7, Hp: 30, HpToiDa: 30, X: 0, Y: 0, SatThuongCoBan: 999_999, LanHoiMauCuoi: t0}}

	k.nhipTim()

	if len(pNanNhan.nhan) != 2 {
		t.Fatalf("nan nhan nhan %d goi, mong doi DUNG 2 (-11 roi -17)", len(pNanNhan.nhan))
	}
	if pNanNhan.nhan[0].Command != -11 {
		t.Fatalf("goi dau cua nan nhan = %d, mong doi -11", pNanNhan.nhan[0].Command)
	}
	if pNanNhan.nhan[1].Command != -17 {
		t.Fatalf("goi thu hai cua nan nhan = %d, mong doi -17 (chinh chu chet, PHAI den SAU -11)", pNanNhan.nhan[1].Command)
	}
	muon17 := []byte{0, 0, 0, 0, 0} // cPk=0, x=0, y=0
	if !bytes.Equal(pNanNhan.nhan[1].Data, muon17) {
		t.Fatalf("goi -17 byte sai:\n got  = %v\n muon = %v", pNanNhan.nhan[1].Data, muon17)
	}

	if len(pNguoiKhac.nhan) != 2 {
		t.Fatalf("nguoi khac nhan %d goi, mong doi DUNG 2 (-10 roi -8)", len(pNguoiKhac.nhan))
	}
	if pNguoiKhac.nhan[0].Command != -10 {
		t.Fatalf("goi dau cua nguoi khac = %d, mong doi -10", pNguoiKhac.nhan[0].Command)
	}
	if pNguoiKhac.nhan[1].Command != -8 {
		t.Fatalf("goi thu hai cua nguoi khac = %d, mong doi -8 (ai do chet, PHAI den SAU -10)", pNguoiKhac.nhan[1].Command)
	}
	muon8 := []byte{
		0, 0, 0, 1, // idNguoiChoi = 1 (int)
		0,    // cPk = 0
		0, 0, // x = 0
		0, 0, // y = 0
	}
	if !bytes.Equal(pNguoiKhac.nhan[1].Data, muon8) {
		t.Fatalf("goi -8 byte sai:\n got  = %v\n muon = %v", pNguoiKhac.nhan[1].Data, muon8)
	}

	if k.nguoiTrong[1].nv.HpHienTai != 0 {
		t.Fatalf("HpHienTai nan nhan = %d, mong doi 0 (chet, kep boi §1.2 — khong am mau)", k.nguoiTrong[1].nv.HpHienTai)
	}
}

// TestSuKienHoiSinhHoiDayMauChakraVaPhatGoi: nguoi choi DA CHET gui
// suKienHoiSinh -> phai HOI DAY DU mau/chakra, dat lucHoiSinh = bayGio, va
// phat DUNG chuoi -16 (rieng chinh chu) roi -30/15 (ca khu vuc) — dung dac ta
// §6.4 "Máy chủ → chính chủ: -16 ... Máy chủ → cả khu vực: -30/15".
func TestSuKienHoiSinhHoiDayMauChakraVaPhatGoi(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return t0 }

	nvA := &store.NhanVat{
		ID: 1, Ten: "A",
		HpToiDa: 150, HpHienTai: 0, // da chet
		ChakraToiDa: 100, ChakraHienTai: 0,
		ViTriX: 200, ViTriY: 300,
	}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})

	nvB := &store.NhanVat{ID: 2, Ten: "B", HpHienTai: 50}
	pB := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvB, phien: pB})
	pA.nhan = nil // don goi -5 tu hai buoc vao khu vuc o tren
	pB.nhan = nil

	k.xuLy(suKienHoiSinh{idNhanVat: nvA.ID})

	if nvA.HpHienTai != 150 {
		t.Fatalf("HpHienTai sau hoi sinh = %d, mong doi 150 (day mau)", nvA.HpHienTai)
	}
	if nvA.ChakraHienTai != 100 {
		t.Fatalf("ChakraHienTai sau hoi sinh = %d, mong doi 100 (day chakra)", nvA.ChakraHienTai)
	}
	if k.nguoiTrong[nvA.ID].lucHoiSinh != t0 {
		t.Fatalf("lucHoiSinh = %v, mong doi %v (dac ta §1.7: chot bat tu 1500ms sau hoi sinh)", k.nguoiTrong[nvA.ID].lucHoiSinh, t0)
	}

	if len(pA.nhan) != 2 {
		t.Fatalf("A (chinh chu) nhan %d goi, mong doi DUNG 2 (-16 roi -30/15)", len(pA.nhan))
	}
	if pA.nhan[0].Command != -16 || len(pA.nhan[0].Data) != 0 {
		t.Fatalf("goi dau cua A = lenh %d (%d byte), mong doi -16 than RONG", pA.nhan[0].Command, len(pA.nhan[0].Data))
	}
	if pA.nhan[1].Command != -30 {
		t.Fatalf("goi thu hai cua A = %d, mong doi -30 (ai do hoi sinh, dung lai ma lenh ME_INFO)", pA.nhan[1].Command)
	}
	muon30 := []byte{
		15,         // byte lenh con = 15
		0, 0, 0, 1, // idNguoiChoi = 1
		0, 0, 0, 150, // hp = 150
		0, 0, 0, 100, // mp = 100
		0, 200, // x = 200
		1, 44, // y = 300
	}
	if !bytes.Equal(pA.nhan[1].Data, muon30) {
		t.Fatalf("goi -30/15 cua A byte sai:\n got  = %v\n muon = %v", pA.nhan[1].Data, muon30)
	}

	if len(pB.nhan) != 1 || pB.nhan[0].Command != -30 {
		t.Fatalf("B (khac, dac ta doi 'ca khu vuc') phai nhan DUNG 1 goi -30/15, nhan %+v", pB.nhan)
	}
	if !bytes.Equal(pB.nhan[0].Data, muon30) {
		t.Fatalf("goi -30/15 cua B (ca khu vuc) byte sai:\n got  = %v\n muon = %v", pB.nhan[0].Data, muon30)
	}
}

// TestSuKienHoiSinhNguoiConSongBiBoQua: gui suKienHoiSinh cho nguoi CON SONG
// (HpHienTai > 0) — phai la no-op HOAN TOAN (khong hoi mau oan, khong phat
// goi gi). Bao ve chong client gui -16 khi khong can (bam nham/bam kep).
func TestSuKienHoiSinhNguoiConSongBiBoQua(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	nvA := &store.NhanVat{ID: 1, HpHienTai: 50, HpToiDa: 150}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	pA.nhan = nil

	k.xuLy(suKienHoiSinh{idNhanVat: nvA.ID})

	if nvA.HpHienTai != 50 {
		t.Fatalf("HpHienTai = %d, mong doi VAN LA 50 (dang song, khong duoc hoi sinh oan)", nvA.HpHienTai)
	}
	if len(pA.nhan) != 0 {
		t.Fatalf("nhan %d goi khi dang song ma van gui suKienHoiSinh, mong doi 0 (im lang bo qua)", len(pA.nhan))
	}
}

// TestSuKienHoiSinhNguoiKhongTrongKhuVucKhongPanic: idNhanVat khong co trong
// khu vuc (chua vao, hoac da ra) -> phai la no-op an toan, khong panic —
// dung tinh than voi TestKhuVucChonChieuNguoiKhongTrongKhuVucKhongPanic.
func TestSuKienHoiSinhNguoiKhongTrongKhuVucKhongPanic(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.xuLy(suKienHoiSinh{idNhanVat: 999})
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
