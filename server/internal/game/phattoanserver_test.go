package game

import (
	"bytes"
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 6 — noi LOA toan may chu vao 5 khu vuc (CHAN-2).
//
// Toan bo file nay bao ve MOT cau: "actor KHONG duoc goi ra ngoai". Duong day
// day du sau khi noi xong:
//
//	actor khu vuc            loa (goroutine rieng)        actor tung khu vuc
//	─────────────────        ─────────────────────        ──────────────────
//	k.xepThongBao(s)   ──►   Xep() -> hang doi      ──►   Gui(suKienPhatGoi)
//	di tiep NGAY             nhip 1,5s roi goi            (hop thu GAMEPLAY,
//	(khong cho ai)           SoKhuVuc.PhatToanServer      khong chan) ->
//	                                                      phatChoCaKhuVuc
//
// Diem chet nguoi o giua: PhatToanServer chay tren goroutine LOA chu KHONG
// phai tren actor, nen no TUYET DOI khong duoc dung guiDieuKhien (khuvuc.go:
// cho toi thoiGianChoControl = 100ms tren hop thu ma vongLap co the dang ban).
// Vi vay no di bang k.Gui — hop thu gameplay, khong bao gio chan.
// ─────────────────────────────────────────────────────────────────────────────

// lenhThongBao93 ghim CUNG ma lenh 93 trong test (khong doc lai hang so
// lenhThongBaoToanServer tu code): con so den tu ServerNotify.java:74, doc lai
// hang so se lam test dong tinh neu ai do sua nham ma lenh — dung ly do da ghi
// o TestGoiThongBaoToanServerByteDung.
const lenhThongBao93 int8 = 93

// thanGoi93 dung than goi 93 mong doi cho mot chuoi ASCII: 2 byte do dai
// big-endian roi den cac byte ky tu. Viet TAY (khong goi protocol.Writer) de
// khong thanh oracle mu.
func thanGoi93(noiDung string) []byte {
	than := []byte{byte(len(noiDung) >> 8), byte(len(noiDung))}
	return append(than, []byte(noiDung)...)
}

// choGoi93 doi cho toi khi phien gia nhan duoc MOT goi 93, tra ve goi do. Cac
// goi khac (-11/-10 do quai AI mac dinh danh nguoi choi, -5 luc vao khu vuc)
// bi BO QUA — day la nhieu hop le tu tinh nang khac, dung tinh than
// khongConGoiGiTruQuaiDanh o khuvuc_tichhop_test.go.
func choGoi93(t *testing.T, p *phienGiaKhuVuc, hanGio time.Duration) *protocol.Message {
	t.Helper()
	hetGio := time.After(hanGio)
	for {
		for _, m := range p.banTin() {
			// Mot goi NIL toi duoc tay nguoi choi nghia la ai do da day
			// suKienPhatGoi{m: nil} vao actor (goiThongBaoToanServer tra nil
			// khi noi dung qua dai). Bat o day de loi hien ra RO RANG thay vi
			// thanh mot cu nil-deref kho hieu giua vong lap.
			if m == nil {
				t.Fatal("khu vuc phat mot goi NIL cho nguoi choi — PhatToanServer PHAI " +
					"kiem goiThongBaoToanServer tra nil truoc khi gui vao hop thu actor")
			}
			if m.Command == lenhThongBao93 {
				return m
			}
		}
		select {
		case <-hetGio:
			return nil
		case <-time.After(5 * time.Millisecond):
		}
	}
}

// khongCoGoi93 khang dinh phien gia KHONG nhan goi 93 nao trong suot han cho.
func khongCoGoi93(t *testing.T, p *phienGiaKhuVuc, hanGio time.Duration) bool {
	t.Helper()
	return choGoi93(t, p, hanGio) == nil
}

// datNguoiVaoKhuVuc dua mot nhan vat + phien gia vao khu vuc QUA HOP THU
// (khong ghi thang k.nguoiTrong — do la trang thai cua actor, luat §1.6) roi
// DONG BO bang mot vong control (AnhChupNguoiChoiCoLoi): control mailbox duoc
// uu tien nhung xuLyGameplayDaXepTruoc chay het gameplay da xep truoc do, nen
// khi ham nay tra ve, su kien vao CHAC CHAN da duoc xu ly xong.
func datNguoiVaoKhuVuc(t *testing.T, k *KhuVuc, id int64, ten string) *phienGiaKhuVuc {
	t.Helper()
	p := &phienGiaKhuVuc{}
	nv := &store.NhanVat{ID: id, Ten: ten, HpToiDa: 100, HpHienTai: 100, ViTriX: 100, ViTriY: 384}
	k.Gui(suKienVao{nv: nv, phien: p})

	ban, err := k.AnhChupNguoiChoiCoLoi(id)
	if err != nil {
		t.Fatalf("dong bo sau khi dua nguoi %d vao khu vuc that bai: %v", id, err)
	}
	if len(ban) != 1 {
		t.Fatalf("sau khi vao, anh chup nguoi %d tra %d ban ghi, muon 1 "+
			"(su kien vao chua duoc actor xu ly?)", id, len(ban))
	}
	return p
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 1: xepThongBao — duong DUY NHAT actor duoc dung de bao tin ra ngoai
// ─────────────────────────────────────────────────────────────────────────────

// TestActorKhongGoiGuiDieuKhienKhiXepThongBao la test QUAN TRONG NHAT file
// nay (CHAN-2). No chay xepThongBao TU BEN TRONG actor — bang suKienChay,
// diem chen chay tren dung goroutine vongLap — roi khang dinh HAI dieu:
//
//  1. KHONG TREO: co han gio 2 giay. Neu ai do sua xepThongBao thanh mot vong
//     "gui roi cho tra loi" (guiDieuKhien / RoiKhuVuc / DanhSachQuaiSnapshot...)
//     thi actor se tu cho chinh minh -> het gio -> test DO ro rang thay vi treo
//     den khi go test tu giet.
//
//  2. KHONG CHAM HOP THU DIEU KHIEN: dem len(k.hopThuDieuKhien) NGAY TRUOC va
//     NGAY SAU loi goi, va dem TU TRONG actor. Day la meo lam cho phep do tro
//     nen TAT DINH: actor la nguoi DUY NHAT rut hop thu dieu khien, ma luc nay
//     chinh no dang ban chay ham cua ta — nen khong ai co the rut bot; bat ky
//     su kien dieu khien nao xepThongBao lo gui vao deu CON NGUYEN o do va lam
//     con dem tang. Chi do "khong treo" thoi la chua du: guiDieuKhien khi hop
//     thu con cho trong se THANH CONG ngay (khong treo) ma van pha vo ky luat
//     actor va se tu khoa khi hop thu day.
func TestActorKhongGoiGuiDieuKhienKhiXepThongBao(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := MoKhuVuc(TimBanDo(1))
	k.loa.Store(loa)
	k.BatDau()
	defer k.Dung()

	type doDuoc struct{ truoc, sau int }
	ketQua := make(chan doDuoc, 1)

	k.Gui(suKienChay{f: func() {
		truoc := len(k.hopThuDieuKhien)
		k.xepThongBao("Boss Cuu Vi da xuat hien")
		ketQua <- doDuoc{truoc: truoc, sau: len(k.hopThuDieuKhien)}
	}})

	select {
	case d := <-ketQua:
		if d.sau != d.truoc {
			t.Fatalf("xepThongBao lam hop thu DIEU KHIEN tang tu %d len %d su kien — "+
				"actor KHONG duoc goi guiDieuKhien (cho toi 100ms tren hop thu ma chinh "+
				"no dang ban -> tu khoa). Phai di bang loa.Xep (khong bao gio chan)",
				d.truoc, d.sau)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("xepThongBao goi TU TRONG actor bi TREO qua 2 giay — dau hieu tu khoa: " +
			"no dang cho mot cau tra loi ma chinh actor nay moi tra duoc")
	}

	if con := len(loa.hangDoi); con != 1 {
		t.Fatalf("hang doi loa giu %d yeu cau sau mot lan xepThongBao, muon dung 1", con)
	}
}

// TestXepThongBaoKhiLoaNil: rat nhieu test (va ca duong chay that truoc khi
// main.go goi GanLoa) dung MoKhuVuc TRUC TIEP, khong he co loa. xepThongBao
// khi do phai la viec KHONG — im lang tuyet doi, khong panic. Neu no panic,
// mot con boss bao tin se giet nhip actor (luoi panic bat duoc nhung van bo
// mat nhip do va ngap log).
func TestXepThongBaoKhiLoaNil(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	if k.loa.Load() != nil {
		t.Fatal("MoKhuVuc khong duoc tu gan loa — loa chi den tu GanLoa (main.go)")
	}

	k.xepThongBao("khong co loa thi khong ai nghe") // khong duoc panic

	if k.loa.Load() != nil {
		t.Fatal("xepThongBao khong duoc tu tao loa cho khu vuc")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2: PhatToanServer — mot cau chu toi MOI nguoi trong MOI khu vuc
// ─────────────────────────────────────────────────────────────────────────────

// TestPhatToanServerToiMoiKhuVuc dung SoKhuVuc THAT (MoSoKhuVuc, 5 actor dang
// chay) va dat hai nguoi o HAI khu vuc KHAC NHAU — diem mau chot: mot thong
// bao toan server khong duoc chi toi khu vuc dau tien hay khu vuc cua nguoi
// gay ra no, ma phai toi TAT CA.
func TestPhatToanServerToiMoiKhuVuc(t *testing.T) {
	s := MoSoKhuVuc()
	defer s.DungTatCa()

	pA := datNguoiVaoKhuVuc(t, s.Tra(1), 8001, "A o map 1")
	pB := datNguoiVaoKhuVuc(t, s.Tra(3), 8003, "B o map 3")

	const noiDung = "Boss Cuu Vi da xuat hien tai Lang La"
	s.PhatToanServer(noiDung)

	for _, tc := range []struct {
		ten string
		p   *phienGiaKhuVuc
	}{
		{"nguoi o khu vuc 1", pA},
		{"nguoi o khu vuc 3", pB},
	} {
		m := choGoi93(t, tc.p, 3*time.Second)
		if m == nil {
			t.Fatalf("%s KHONG nhan duoc goi 93 sau 3 giay — thong bao toan server "+
				"phai toi MOI nguoi trong MOI khu vuc", tc.ten)
		}
		if muon := thanGoi93(noiDung); !bytes.Equal(m.Data, muon) {
			t.Fatalf("%s nhan goi 93 sai than:\n got  = %v\n muon = %v", tc.ten, m.Data, muon)
		}
	}
}

// TestPhatToanServerNoiDungQuaDaiKhongPhatGi: goiThongBaoToanServer tra nil khi
// chuoi vuot 65535 byte (writeUTF ghi do dai bang 2 byte khong dau). Luc do
// PhatToanServer phai IM — gui mot goi nil vao hop thu actor se thanh mot loi
// nil pointer trong vongLap (luoi panic bat duoc, nhung do la mat mot nhip vo
// co) hoac te hon la mot goi hong xuong client.
func TestPhatToanServerNoiDungQuaDaiKhongPhatGi(t *testing.T) {
	s := MoSoKhuVuc()
	defer s.DungTatCa()

	p := datNguoiVaoKhuVuc(t, s.Tra(1), 8009, "nguoi thu")

	qua := make([]byte, 65536)
	for i := range qua {
		qua[i] = 'a'
	}
	s.PhatToanServer(string(qua))

	if !khongCoGoi93(t, p, 300*time.Millisecond) {
		t.Fatal("noi dung qua dai (65536 byte) van duoc phat — goiThongBaoToanServer tra nil, " +
			"PhatToanServer PHAI kiem nil truoc khi gui")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 3: GanLoa — noi day, va LUAT 4 (MoSoKhuVuc giu nguyen)
// ─────────────────────────────────────────────────────────────────────────────

// TestMoSoKhuVucKhongTuKhoiDongLoa KHOA LUAT 4 cua du an: MoSoKhuVuc() khong
// duoc tu khoi dong THEM bat cu thu gi. Rat nhieu test tich hop dung no de
// dung 5 khu vuc va kiem "im lang" (helper chi tha goi -11/-10); mot cai loa
// tu bat o day se lam vo hang loat test do — va nguy hiem hon, se khong ai
// tat duoc no khi tat may chu vi khong ai giu tay cam.
func TestMoSoKhuVucKhongTuKhoiDongLoa(t *testing.T) {
	s := MoSoKhuVuc()
	defer s.DungTatCa()

	for mapId := int32(1); mapId <= 5; mapId++ {
		if l := s.Tra(mapId).loa.Load(); l != nil {
			t.Fatalf("khu vuc %d da co loa ngay sau MoSoKhuVuc — loa CHI duoc gan bang "+
				"GanLoa (main.go goi rieng, Task 14)", mapId)
		}
	}

	// Va khong co goi 93 nao tu dong chay ra: mot nguoi ngoi im trong khu vuc
	// khong duoc nghe thong bao nao ca.
	p := datNguoiVaoKhuVuc(t, s.Tra(1), 8011, "nguoi ngoi im")
	if !khongCoGoi93(t, p, 300*time.Millisecond) {
		t.Fatal("MoSoKhuVuc tu phat goi 93 — no phai hoan toan im lang")
	}
}

// TestGanLoaNoiDayTuActorToiMoiNguoi la test WIRING day du: chay CA CHUOI
// that — actor khu vuc 1 goi xepThongBao (tu trong suKienChay, dung goroutine
// actor) -> loa rut hang doi tren goroutine rieng -> goi PhatToanServer ->
// Gui(suKienPhatGoi) vao 5 khu vuc -> phatChoCaKhuVuc. Nguoi nhan nam o khu
// vuc 2, tuc thong bao that su di XUYEN khu vuc, khong quanh quan trong cai
// khu vuc da phat.
func TestGanLoaNoiDayTuActorToiMoiNguoi(t *testing.T) {
	s := MoSoKhuVuc()
	loa := MoBoPhatThanh()
	s.GanLoa(loa)
	defer func() {
		loa.Dung() // tat loa TRUOC: khong con ai Gui vao khu vuc dang dong
		s.DungTatCa()
	}()

	p := datNguoiVaoKhuVuc(t, s.Tra(2), 8021, "nguoi o map 2")

	const noiDung = "Su kien Lang La bat dau"
	k1 := s.Tra(1)
	k1.Gui(suKienChay{f: func() { k1.xepThongBao(noiDung) }})

	m := choGoi93(t, p, 3*time.Second)
	if m == nil {
		t.Fatal("sau 3 giay, thong bao xep tu actor khu vuc 1 van chua toi nguoi o khu vuc 2 — " +
			"duong day GanLoa -> Xep -> PhatToanServer -> Gui bi dut o dau do")
	}
	if muon := thanGoi93(noiDung); !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 93 nhan duoc sai than:\n got  = %v\n muon = %v", m.Data, muon)
	}
}
