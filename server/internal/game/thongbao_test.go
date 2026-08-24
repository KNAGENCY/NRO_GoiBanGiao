package game

import (
	"bytes"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// Phan 1: khuon byte goi 93 (goiThongBaoToanServer)
// ─────────────────────────────────────────────────────────────────────────────

// TestGoiThongBaoToanServerByteDung ghim khuon byte cua goi 93 theo NGUON
// 1-MAY-CHU-GAME/src/server/ServerNotify.java:74-76 (`new Message(93)` +
// `writeUTF(text)`): than goi CHI CO writeUTF, tuc 2 byte do dai (big-endian)
// roi den noi dung ma hoa modified UTF-8 — KHONG co byte loai, KHONG co mau,
// KHONG co truong dem nao khac. Byte mong doi viet TAY o day (khong goi lai
// protocol.Writer de tu sinh) — neu dung Writer de tinh ky vong thi test se
// dong tinh voi bat ky loi nao cua Writer, thanh oracle mu.
func TestGoiThongBaoToanServerByteDung(t *testing.T) {
	for _, tc := range []struct {
		ten     string
		noiDung string
		muon    []byte
	}{
		{
			ten:     "chuoi ASCII",
			noiDung: "Chao lang La",
			muon: []byte{
				0, 12, // do dai sau ma hoa = 12 byte (2 byte big-endian)
				'C', 'h', 'a', 'o', ' ', 'l', 'a', 'n', 'g', ' ', 'L', 'a',
			},
		},
		{
			ten:     "chuoi co dau (modified UTF-8, 2 byte moi ky tu ngoai ASCII)",
			noiDung: "Cá",
			muon: []byte{
				0, 3, // 'C' 1 byte + 'á' 2 byte = 3
				0x43,       // 'C'
				0xC3, 0xA1, // 'á' = U+00E1
			},
		},
		{
			ten:     "chuoi rong van la goi hop le (2 byte do dai = 0)",
			noiDung: "",
			muon:    []byte{0, 0},
		},
	} {
		t.Run(tc.ten, func(t *testing.T) {
			m := goiThongBaoToanServer(tc.noiDung)
			if m == nil {
				t.Fatalf("goiThongBaoToanServer(%q) tra nil, mong doi mot goi hop le", tc.noiDung)
			}
			// 93 ghim CUNG o day (khong doc lai hang so lenhThongBaoToanServer
			// tu code) — con so nay den tu ServerNotify.java:74, doc lai hang so
			// se lam test dong tinh neu ai do sua nham ma lenh.
			if m.Command != 93 {
				t.Fatalf("Command = %d, muon 93 (ServerNotify.java:74)", m.Command)
			}
			if !bytes.Equal(m.Data, tc.muon) {
				t.Fatalf("than goi 93 sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
					m.Data, len(m.Data), tc.muon, len(tc.muon))
			}
		})
	}
}

// TestGoiThongBaoToanServerQuaDaiTraNil: writeUTF ghi do dai bang 2 byte khong
// dau -> toi da 65535 byte. Dai hon the thi KHONG duoc dung goi thieu/hong (ca
// luong tin phia sau se doc lech) — ham phai tra nil de nguoi goi bo qua.
func TestGoiThongBaoToanServerQuaDaiTraNil(t *testing.T) {
	if m := goiThongBaoToanServer(strings.Repeat("a", 65536)); m != nil {
		t.Fatalf("chuoi 65536 byte phai tra nil, nhung nhan duoc goi %d byte than", len(m.Data))
	}
	// Ngay duoi nguong thi VAN phai dung duoc — chan tren dung cho, khong bi
	// cat non mot byte.
	if m := goiThongBaoToanServer(strings.Repeat("a", 65535)); m == nil {
		t.Fatal("chuoi 65535 byte (dung nguong) phai dung duoc goi, nhung tra nil")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2: hang doi KHONG CHAN (CHAN-2 — ly do ton tai cua ca file)
// ─────────────────────────────────────────────────────────────────────────────

// TestBoPhatThanhXepKhongChanKhiHangDoiDay la test QUAN TRONG NHAT file nay:
// Xep() duoc goi TU TRONG actor KhuVuc, nen no TUYET DOI khong duoc chan — du
// hang doi da day cung va goroutine phat chua he chay (o day co tinh KHONG goi
// BatDau de mo phong truong hop xau nhat: khong ai rut hang doi ca).
//
// Ghim CUNG suc chua 64 (khong doc lai hang so tranHangDoiThongBao tu code —
// bay oracle mu, xem cung ly do o luudinhky_test.go: "Ghim CUNG 59 giay"):
// xep 74 lan thi 64 lan dau vao hang doi, 10 lan sau bi BO va dem soLanBo.
//
// Cach do "khong chan": chay toan bo 74 lan Xep trong mot goroutine rieng va
// dat han gio — neu Xep chan that (vd doi ai co ai do sua kenh thanh khong
// tran), test bao loi RO RANG thay vi treo den khi go test tu giet.
func TestBoPhatThanhXepKhongChanKhiHangDoiDay(t *testing.T) {
	const sucChua = 64
	const soXep = 74

	b := MoBoPhatThanh()
	defer b.Dung()

	xongHet := make(chan struct{})
	go func() {
		defer close(xongHet)
		for i := 0; i < soXep; i++ {
			b.Xep(YeuCauThongBao{NoiDung: "thu", MocTao: time.Now()})
		}
	}()

	select {
	case <-xongHet:
	case <-time.After(2 * time.Second):
		t.Fatalf("Xep() BI CHAN: %d lan xep vao hang doi suc chua %d chua xong sau 2 giay "+
			"(hang doi PHAI co tran va Xep PHAI dung select+default)", soXep, sucChua)
	}

	if soBo := b.SoLanBo(); soBo != soXep-sucChua {
		t.Fatalf("SoLanBo() = %d, muon %d (xep %d lan vao hang doi suc chua %d)",
			soBo, soXep-sucChua, soXep, sucChua)
	}
	if con := len(b.hangDoi); con != sucChua {
		t.Fatalf("hang doi giu %d yeu cau, muon dung %d (khong duoc phinh qua tran)", con, sucChua)
	}
}

// TestBoPhatThanhXepNhieuGoroutineKhongChan: nhieu actor KhuVuc chay song song
// (moi khu vuc mot goroutine) co the cung Xep mot luc — dem soLanBo phai la
// atomic, va khong goroutine nao duoc chan. Chay voi -race se bat duoc neu ai
// do doi soLanBo thanh int thuong.
func TestBoPhatThanhXepNhieuGoroutineKhongChan(t *testing.T) {
	const soGoroutine = 8
	const moiGoroutineXep = 50 // 8*50 = 400 >> suc chua 64, chac chan co bo

	b := MoBoPhatThanh()
	defer b.Dung()

	var cho sync.WaitGroup
	for i := 0; i < soGoroutine; i++ {
		cho.Add(1)
		go func() {
			defer cho.Done()
			for j := 0; j < moiGoroutineXep; j++ {
				b.Xep(YeuCauThongBao{NoiDung: "bao", MocTao: time.Now()})
			}
		}()
	}

	xongHet := make(chan struct{})
	go func() { cho.Wait(); close(xongHet) }()

	select {
	case <-xongHet:
	case <-time.After(3 * time.Second):
		t.Fatal("Xep() tu nhieu goroutine BI CHAN (chua xong sau 3 giay)")
	}

	tongXep := int64(soGoroutine * moiGoroutineXep)
	if soBo, con := b.SoLanBo(), int64(len(b.hangDoi)); soBo+con != tongXep {
		t.Fatalf("soBo(%d) + conTrongHangDoi(%d) = %d, muon dung tong so lan Xep = %d "+
			"(dem soLanBo phai atomic, khong duoc mat lan dem nao)", soBo, con, soBo+con, tongXep)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 3: nhip chong ngap 1,5 giay (dong ho tiem duoc)
// ─────────────────────────────────────────────────────────────────────────────

// TestBoPhatThanhTonTrongNhipToiThieu la MOC KIEM THU DOT BIEN cua task: xoa
// khoi phan chong ngap trong nhip() (thongbao.go) thi test nay PHAI DO.
//
// Ghim CUNG 1,5 giay = 1500ms (khong doc lai hang so nhipToiThieuThongBao —
// con so den tu ServerNotify.java:65 `Functions.sleep(1500)`). Dung DONG HO GIA
// (goi nhip() truc tiep voi bayGio tiem san) nen test chay xong trong vai
// micro-giay, khong Sleep that giay nao.
//
// Kich ban: hai thong bao A, B xep san. Lan luot goi nhip() tai cac moc
// 0 / 1400 / 1499 / 1500 ms roi ghi lai NHUNG MOC NAO thuc su phat ra tin.
//   - moc 0    -> A phat NGAY (lanPhatCuoi con zero: thong bao dau tien khong
//     phai cho, dung nhu ServerNotify.run() gui roi moi sleep).
//   - 1400/1499 -> KHONG duoc phat gi (chua du 1,5 giay), va B KHONG DUOC MAT.
//   - 1500      -> B phat (dung nguong).
func TestBoPhatThanhTonTrongNhipToiThieu(t *testing.T) {
	b := MoBoPhatThanh()
	t0 := time.Now()

	b.Xep(YeuCauThongBao{NoiDung: "A", MocTao: t0})
	b.Xep(YeuCauThongBao{NoiDung: "B", MocTao: t0})

	var daPhat []string
	var mocPhat []time.Duration
	phat := func(s string) { daPhat = append(daPhat, s) }

	for _, lech := range []time.Duration{
		0,
		1400 * time.Millisecond,
		1499 * time.Millisecond,
		1500 * time.Millisecond,
	} {
		truoc := len(daPhat)
		b.nhip(t0.Add(lech), phat)
		if len(daPhat) > truoc {
			mocPhat = append(mocPhat, lech)
		}
		// Sau moc 1400ms: B bi day som PHAI con nam trong hang doi (xep lai,
		// khong bi vut bo) — day la nua thu hai cua hop dong chong ngap.
		if lech == 1400*time.Millisecond {
			if con := len(b.hangDoi); con != 1 {
				t.Fatalf("tai +1400ms (chua du nhip), hang doi con %d yeu cau, muon 1 "+
					"(thong bao den som phai duoc XEP LAI cho toi luot, khong duoc mat)", con)
			}
			if soBo := b.SoLanBo(); soBo != 0 {
				t.Fatalf("tai +1400ms, SoLanBo() = %d, muon 0 (hang doi con cho, khong duoc bo)", soBo)
			}
		}
	}

	muonMoc := []time.Duration{0, 1500 * time.Millisecond}
	if len(mocPhat) != len(muonMoc) || mocPhat[0] != muonMoc[0] || mocPhat[1] != muonMoc[1] {
		t.Fatalf("cac moc PHAT ra tin = %v, muon %v "+
			"(hai thong bao lien tiep phai cach nhau it nhat 1,5 giay)", mocPhat, muonMoc)
	}
	if len(daPhat) != 2 || daPhat[0] != "A" || daPhat[1] != "B" {
		t.Fatalf("noi dung da phat = %v, muon [A B] (khong duoc mat, khong duoc dao)", daPhat)
	}
}

// TestBoPhatThanhHangDoiRongThiKhongPhatGi: mot nhip khi khong co gi trong hang
// doi phai la viec khong (khong goi phat, khong cham lanPhatCuoi). Neu nhip lo
// tay cap nhat lanPhatCuoi khi hang doi rong, thong bao dau tien den sau do se
// bi bat cho oan 1,5 giay.
func TestBoPhatThanhHangDoiRongThiKhongPhatGi(t *testing.T) {
	b := MoBoPhatThanh()
	t0 := time.Now()

	soLanGoiPhat := 0
	phat := func(string) { soLanGoiPhat++ }

	b.nhip(t0, phat)
	b.nhip(t0.Add(time.Second), phat)
	if soLanGoiPhat != 0 {
		t.Fatalf("hang doi rong ma goi phat() %d lan, muon 0", soLanGoiPhat)
	}

	// Xep MOT thong bao roi nhip NGAY (chi +1ms sau lan nhip rong cuoi): phai
	// phat luon, chung to hai nhip rong o tren khong lam ban lanPhatCuoi.
	b.Xep(YeuCauThongBao{NoiDung: "dau tien", MocTao: t0})
	b.nhip(t0.Add(time.Second+time.Millisecond), phat)
	if soLanGoiPhat != 1 {
		t.Fatalf("thong bao DAU TIEN khong duoc phat ngay (soLanGoiPhat = %d, muon 1) — "+
			"nhip rong khong duoc cap nhat lanPhatCuoi", soLanGoiPhat)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 4: vong doi goroutine (BatDau / Dung / khong ro ri)
// ─────────────────────────────────────────────────────────────────────────────

// TestBoPhatThanhDungHaiLanKhongPanic: Dung() phai idempotent (sync.Once) —
// luc tat may chu co the co nhieu duong cung goi Dung. Kem luon truong hop
// Dung() khi CHUA BatDau(): phai tra ve ngay, khong treo cho mot goroutine
// khong ton tai.
func TestBoPhatThanhDungHaiLanKhongPanic(t *testing.T) {
	t.Run("da BatDau", func(t *testing.T) {
		b := MoBoPhatThanh()
		b.BatDau(func(string) {})
		b.Dung()
		b.Dung() // lan hai: khong duoc panic (close kenh da dong)
	})

	t.Run("chua BatDau", func(t *testing.T) {
		b := MoBoPhatThanh()
		xong := make(chan struct{})
		go func() {
			defer close(xong)
			b.Dung()
			b.Dung()
		}()
		select {
		case <-xong:
		case <-time.After(2 * time.Second):
			t.Fatal("Dung() khi chua BatDau() bi TREO — phai tra ve ngay vi khong co goroutine nao de cho")
		}
	})
}

// TestBoPhatThanhKhongRoRiGoroutine: mo/dong bo phat thanh nhieu lan roi doi so
// goroutine ve dung cho cu. Ro ri goroutine o day la loi am tham nguy hiem —
// may chu chay hang thang, moi lan mo lai mot bo phat thanh ma khong don sach
// se tich luy dan cho toi khi het bo nho.
func TestBoPhatThanhKhongRoRiGoroutine(t *testing.T) {
	truoc := soGoroutineOnDinh(t)

	for i := 0; i < 5; i++ {
		b := MoBoPhatThanh()
		b.BatDau(func(string) {})
		b.Xep(YeuCauThongBao{NoiDung: "thu", MocTao: time.Now()})
		b.Dung()
	}

	sau := soGoroutineOnDinh(t)
	if sau > truoc {
		t.Fatalf("so goroutine sau 5 lan mo/dong = %d, truoc do = %d — RO RI %d goroutine "+
			"(Dung() phai cho vongLap thoat han)", sau, truoc, sau-truoc)
	}
}

// soGoroutineOnDinh doc runtime.NumGoroutine() lap lai cho toi khi con so DUNG
// YEN mot luc — can thiet vi goroutine cua cac test truoc (actor khu vuc, tho
// ghi...) co the con dang thu don, va vi chinh goroutine vua Dung() cung can
// vai nhip de bien mat khoi bo dem cua runtime. Doc mot phat roi so sanh ngay
// se cho ket qua chap chon (test "nhap nhay" — te hon la khong co test).
func soGoroutineOnDinh(t *testing.T) int {
	t.Helper()

	const soLanGiongNhauLaOnDinh = 5
	moc := runtime.NumGoroutine()
	giongNhau := 0
	hetGio := time.After(3 * time.Second)

	for {
		select {
		case <-hetGio:
			t.Fatalf("so goroutine khong on dinh sau 3 giay (dang %d, moc %d)",
				runtime.NumGoroutine(), moc)
		case <-time.After(10 * time.Millisecond):
		}

		hienTai := runtime.NumGoroutine()
		if hienTai == moc {
			giongNhau++
			if giongNhau >= soLanGiongNhauLaOnDinh {
				return hienTai
			}
			continue
		}
		moc = hienTai
		giongNhau = 0
	}
}

// TestBoPhatThanhVongLapThatPhatQuaBatDau la test WIRING: cac test tren goi
// nhip() TRUC TIEP (bo qua goroutine that), test nay chay QUA BatDau() + ticker
// THAT de xac nhan day du duong day "ticker -> nhip(dongHo()) -> phat()" that
// su noi dung, khong chi logic nhip() bi co lap dung. Dung dong ho THAT
// (time.Now mac dinh) vi thong bao dau tien khong phai cho nhip nao ca.
func TestBoPhatThanhVongLapThatPhatQuaBatDau(t *testing.T) {
	b := MoBoPhatThanh()
	daPhat := make(chan string, 4)

	b.Xep(YeuCauThongBao{NoiDung: "Lang La mo cong", MocTao: time.Now()})
	b.BatDau(func(s string) { daPhat <- s })
	defer b.Dung()

	select {
	case s := <-daPhat:
		if s != "Lang La mo cong" {
			t.Fatalf("vongLap that phat ra %q, muon %q", s, "Lang La mo cong")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("sau 3 giay that, vongLap that (qua BatDau + ticker) van chua phat thong bao nao")
	}
}
