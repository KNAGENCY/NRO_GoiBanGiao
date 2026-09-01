package game

import (
	"runtime"
	"testing"
	"time"

	"nroserver/internal/protocol"
)

// ---------------------------------------------------------------------------
// KH6b Task 13 — bo dieu phoi lich boss chay NGOAI actor (lichboss.go).
//
// CACH QUAN SAT: khong test nao duoc doc thang k.boss khi actor DANG CHAY — do
// la data race that (go test -race bat ngay). Thay vao do dat MOT NGUOI QUAN
// SAT (phien gia) trong khu vuc roi DEM GOI -5 nguoi do nhan duoc: boss vao khu
// vuc thi themBoss phat -5 cho ca khu vuc. Day dung la thu nguoi choi that
// thay, nen test khoa dung cai dang duoc hua chu khong khoa mot chi tiet noi bo.
// ---------------------------------------------------------------------------

// datNguoiQuanSat dua mot phien gia vao khu vuc rồi CHO TOI KHI actor xu ly
// xong. Tra ve phien do de dem goi.
//
// Chot dong bo bang AnhChupNguoiChoiCoLoi (khong phai Sleep): no di qua hop thu
// DIEU KHIEN, ma xuLyDieuKhien rut het gameplay da xep truoc do (xem
// xuLyGameplayDaXepTruoc, khuvuc.go) — nen khi no tra ve la chac chan suKienVao
// da chay xong.
func datNguoiQuanSat(t *testing.T, k *KhuVuc, id int64) *phienGiaKhuVuc {
	t.Helper()
	p := &phienGiaKhuVuc{}
	k.Gui(suKienVao{nv: nvMau(id, "QuanSat"), phien: p})
	chotActor(t, k)
	if n := demGoi5(p); n != 0 {
		t.Fatalf("nguoi quan sat vua vao khu vuc da nhan %d goi -5 — khu vuc phai trong truoc khi test bat dau", n)
	}
	return p
}

// chotActor cho actor xu ly het moi thu da gui truoc do (xem datNguoiQuanSat).
func chotActor(t *testing.T, k *KhuVuc) {
	t.Helper()
	if _, err := k.AnhChupNguoiChoiCoLoi(0); err != nil {
		t.Fatalf("cho actor xu ly xong that bai: %v", err)
	}
}

// demGoi5 dem so goi -5 (xuat hien) ma phien gia da nhan.
func demGoi5(p *phienGiaKhuVuc) int {
	n := 0
	for _, m := range p.banTin() {
		if m.Command == lenhXuatHien {
			n++
		}
	}
	return n
}

// goi5DauTien tra goi -5 dau tien phien gia nhan duoc (nil neu chua co).
func goi5DauTien(p *phienGiaKhuVuc) *protocol.Message {
	for _, m := range p.banTin() {
		if m.Command == lenhXuatHien {
			return m
		}
	}
	return nil
}

// phaMau la mot pha boss HOP LE toi thieu: du truong de goiXuatHienBoss dung
// duoc goi -5 (Cap >= 1, TypePk != 0) va du truong de bo dieu phoi lam viec
// (GiayNghi, MapVao, HpChon).
func phaMau(giayNghi int32, mapVao []int32, hpChon []int64) PhaBoss {
	return PhaBoss{
		Ten:       "CuuVi",
		Cap:       50,
		TypePk:    TypePkBoss,
		GiayNghi:  giayNghi,
		MapVao:    mapVao,
		HpChon:    hpChon,
		LaBossLon: false,
	}
}

// mucMau goi tat: mot muc lich mot pha.
func mucMau(giayNghi int32, mapVao []int32, hpChon []int64) MucLichBoss {
	return MucLichBoss{DuLieu: []PhaBoss{phaMau(giayNghi, mapVao, hpChon)}}
}

// moSoKhuVucCoDon mo the gioi that va tu dang ky don dep.
func moSoKhuVucCoDon(t *testing.T) *SoKhuVuc {
	t.Helper()
	s := MoSoKhuVuc()
	t.Cleanup(func() { _ = s.DungTatCa() })
	return s
}

// ---------------------------------------------------------------------------
// 1. TEST CHOT CHAN-3 — ly do ton tai cua ca Task 13.
// ---------------------------------------------------------------------------

// TestChotCHAN3_MoSoKhuVucKhongTuSinhBossNao la CHOT KIEM cua ke hoach (muc
// CHAN-3): MoSoKhuVuc() THUAN phai giu nguyen hanh vi "5 khu vuc im lang".
//
// Neu ai do nhet khoi dong lich boss vao MoSoKhuVuc(), hang chuc test tich hop
// se do voi thong bao kieu "nhan duoc goi NGOAI DU KIEN: lenh=-5" — mot dau vet
// khong he goi y nguyen nhan that. Test nay dung o day de LOI DO NOI THANG ra
// ngay tai cho, bang ten.
//
// Cho THAT (khong phai dong ho tiem) la CO Y: cai dang duoc khang dinh la
// "khong co goroutine nen nao dang chay" — dieu do chi chung minh duoc bang
// cach de thoi gian that troi qua.
func TestChotCHAN3_MoSoKhuVucKhongTuSinhBossNao(t *testing.T) {
	s := moSoKhuVucCoDon(t)

	quanSat := make(map[int32]*phienGiaKhuVuc, 5)
	for mapId := int32(1); mapId <= 5; mapId++ {
		quanSat[mapId] = datNguoiQuanSat(t, s.Tra(mapId), int64(mapId))
	}

	// 400ms = 4 nhip 100ms cua ca 5 actor — thua du de mot bo dieu phoi bi nhet
	// nham vao MoSoKhuVuc kip sinh lua boss dau tien.
	time.Sleep(400 * time.Millisecond)

	for mapId := int32(1); mapId <= 5; mapId++ {
		chotActor(t, s.Tra(mapId))
		if n := demGoi5(quanSat[mapId]); n != 0 {
			t.Fatalf("khu vuc %d phat %d goi -5 sau MoSoKhuVuc() thuan — MoSoKhuVuc TUYET DOI khong duoc khoi dong lich boss (CHAN-3)", mapId, n)
		}
	}
}

// ---------------------------------------------------------------------------
// 2. Hai buoc Mo/BatDau — dong ho tiem duoc chi co nghia neu goroutine chua chay
// ---------------------------------------------------------------------------

// TestMoBoLichBossChuaChayGiCho ToiKhiBatDau: MoBoLichBoss KHONG duoc tu khoi
// dong goroutine nao. Do la ly do ta tach lam hai buoc — de test kip tiem dong
// ho/bo chon TRUOC khi co ai doc chung.
func TestMoBoLichBossChuaChayGiChoToiKhiBatDau(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	// GiayNghi 0 = "den han ngay lap tuc" — neu MoBoLichBoss tu chay, con boss
	// dau tien se ra trong vong mot nhip.
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(0, []int32{1}, []int64{500})})
	time.Sleep(300 * time.Millisecond)
	chotActor(t, k)
	if n := demGoi5(p); n != 0 {
		t.Fatalf("MoBoLichBoss (chua BatDau) da sinh %d boss — no phai tao xong ma KHONG chay goroutine nao", n)
	}

	l.BatDau()
	defer l.Dung()
	if !choToiKhi(func() bool {
		chotActor(t, k)
		return demGoi5(p) >= 1
	}) {
		t.Fatal("sau BatDau() bo dieu phoi khong sinh boss nao trong 2 giay")
	}
}

// choToiKhi hoi dieu kien moi 20ms cho toi 2 giay. Dung cho cac khang dinh
// PHAI cho goroutine nen that chay (khong the tiem dong ho o day) — van tra ve
// ngay khi dieu kien dat, khong ngu du 2 giay.
func choToiKhi(dieuKien func() bool) bool {
	han := time.Now().Add(2 * time.Second)
	for time.Now().Before(han) {
		if dieuKien() {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return dieuKien()
}

// ---------------------------------------------------------------------------
// 3. Nguong GiayNghi — goi nhip() truc tiep voi bayGio tiem san, khong Sleep
// ---------------------------------------------------------------------------

// TestBoLichBossChiSinhKhiDuGiayNghi: 29 giay CHUA sinh, 30 giay MOI sinh.
// Khong co dong ho tiem duoc thi test nay phai ngu 30 giay that.
func TestBoLichBossChiSinhKhiDuGiayNghi(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	l.muc[0].lanSinhCuoi = goc // moc goc ro rang, khong phu thuoc time.Now luc Mo

	l.nhip(goc.Add(29 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 0 {
		t.Fatalf("moi 29/30 giay ma da sinh %d boss", n)
	}

	l.nhip(goc.Add(30 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("du 30 giay: sinh %d boss, mong doi dung 1", n)
	}
}

// TestBoLichBossKhongSinhTrungKhiBossCuConSong tra loi cau hoi "boss dang song
// ma toi gio sinh tiep thi sao": BO QUA. Mot muc lich chi co TOI DA MOT con
// boss trong the gioi — sinh trung se lam nguoi choi thay hai con Cuu Vi dung
// canh nhau, va khong con cach nao noi con nao la con nao.
//
// Cua so mo lai la TheSongBoss.Tat() — soi day ma may trang thai (Task 9) se
// keo khi boss roi khoi the gioi.
func TestBoLichBossKhongSinhTrungKhiBossCuConSong(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	l.muc[0].lanSinhCuoi = goc

	l.nhip(goc.Add(30 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("lan sinh dau: %d boss, mong doi 1", n)
	}

	// Qua HAI chu ky nua ma con cu van song -> van chi mot con.
	l.nhip(goc.Add(60 * time.Second))
	l.nhip(goc.Add(90 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("boss cu con song ma da sinh them: tong %d goi -5, mong doi van 1", n)
	}

	// Boss cu roi khoi the gioi -> chu ky sau sinh lai.
	l.TheSong(0).Tat()
	l.nhip(goc.Add(120 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 2 {
		t.Fatalf("sau khi boss cu roi di, chu ky moi sinh khong ra: tong %d goi -5, mong doi 2", n)
	}
}

// ---------------------------------------------------------------------------
// 4. "O map nao" — bo chon ngau nhien tiem duoc
// ---------------------------------------------------------------------------

// TestBoLichBossChonMapVaMauTheoBoChonTiem: bo dieu phoi phai gui boss DUNG khu
// vuc da chon (khong phai khu vuc dau danh sach, khong phai tat ca), va nap
// dung gia tri mau da chon.
func TestBoLichBossChonMapVaMauTheoBoChonTiem(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	p := make(map[int32]*phienGiaKhuVuc, 3)
	for mapId := int32(1); mapId <= 3; mapId++ {
		p[mapId] = datNguoiQuanSat(t, s.Tra(mapId), int64(mapId))
	}

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1, 2, 3}, []int64{111, 222})})
	l.muc[0].lanSinhCuoi = goc
	// Kich ban chon: lan hoi thu nhat (map, n=3) tra 2 -> map thu ba = mapId 3;
	// lan thu hai (mau, n=2) tra 1 -> 222.
	kichBan := []int{2, 1}
	buoc := 0
	l.chon = func(n int) int {
		if buoc >= len(kichBan) {
			t.Errorf("bo chon bi hoi %d lan, kich ban chi co %d", buoc+1, len(kichBan))
			return 0
		}
		v := kichBan[buoc]
		buoc++
		if v >= n {
			t.Errorf("kich ban tra %d nhung mien hop le la [0,%d)", v, n)
			return 0
		}
		return v
	}

	l.nhip(goc.Add(30 * time.Second))
	for mapId := int32(1); mapId <= 3; mapId++ {
		chotActor(t, s.Tra(mapId))
	}
	if n := demGoi5(p[1]); n != 0 {
		t.Errorf("map 1 khong duoc chon nhung nhan %d goi -5", n)
	}
	if n := demGoi5(p[2]); n != 0 {
		t.Errorf("map 2 khong duoc chon nhung nhan %d goi -5", n)
	}
	if n := demGoi5(p[3]); n != 1 {
		t.Fatalf("map 3 (da chon) nhan %d goi -5, mong doi 1", n)
	}

	g := docGoi5(t, goi5DauTien(p[3]))
	if g.hp != 222 || g.hpToiDa != 222 {
		t.Errorf("mau boss = %d/%d, mong doi 222/222 (gia tri bo chon da chon)", g.hp, g.hpToiDa)
	}
	if g.id >= 0 {
		t.Errorf("id boss = %d, phai la so AM (bay #2 thiet ke §7)", g.id)
	}
	if g.ten != "CuuVi" || g.cap != 50 || g.typePk != TypePkBoss {
		t.Errorf("goi -5 sai: ten=%q cap=%d typePk=%d", g.ten, g.cap, g.typePk)
	}
}

// TestBoLichBossSinhTaiDiemSinhCuaBanDo: boss phai dung tren nen dat (diem sinh
// mac dinh cua ban do), khong phai (0,0) — cot 0 la tuong bien, con boss se
// nam trong da.
func TestBoLichBossSinhTaiDiemSinhCuaBanDo(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(2)
	p := datNguoiQuanSat(t, k, 1)

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(0, []int32{2}, []int64{500})})
	l.muc[0].lanSinhCuoi = goc
	l.nhip(goc)
	chotActor(t, k)

	g := docGoi5(t, goi5DauTien(p))
	bd := TimBanDo(2)
	if g.x != bd.SpawnX || g.y != bd.SpawnY {
		t.Errorf("boss sinh o (%d,%d), mong doi diem sinh cua ban do (%d,%d)", g.x, g.y, bd.SpawnX, bd.SpawnY)
	}
}

// ---------------------------------------------------------------------------
// 5. Gui lenh THAT BAI — hai chinh sach khac han nhau
// ---------------------------------------------------------------------------

// TestBoLichBossGuiHongTamThoiThiHoanRoiThuLai: khu vuc chua BatDau() ->
// guiDieuKhien tra ErrActorChuaBatDau. Ba dieu phai dung:
//
//  1. KHONG coi nhu da sinh (neu coi la sinh thi con boss do bien mat vinh vien).
//  2. KHONG dong vao lanSinhCuoi (neu dong, mot lan hong day boss lui ca chu ky).
//  3. HOAN khoangHoanKhiGuiHong roi moi thu lai — khong dap cua 10 lan/giay.
func TestBoLichBossGuiHongTamThoiThiHoanRoiThuLai(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	s := &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}} // CHUA BatDau -> gui se hong

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	l.muc[0].lanSinhCuoi = goc

	lucHong := goc.Add(30 * time.Second)
	l.nhip(lucHong)
	if l.muc[0].the.ConSong() {
		t.Fatal("gui lenh sinh THAT BAI ma the song van duoc bat — con boss do se khong bao gio ra")
	}
	if !l.muc[0].lanSinhCuoi.Equal(goc) {
		t.Fatalf("gui hong ma lanSinhCuoi bi doi thanh %v — mot lan hong khong duoc day boss lui ca chu ky", l.muc[0].lanSinhCuoi)
	}
	if l.muc[0].ngungHan {
		t.Fatal("khu vuc chua BatDau la loi TAM THOI, khong duoc ngung han")
	}
	hoanCu := l.muc[0].hoanToi
	if !hoanCu.Equal(lucHong.Add(khoangHoanKhiGuiHong)) {
		t.Fatalf("hoanToi = %v, mong doi %v", hoanCu, lucHong.Add(khoangHoanKhiGuiHong))
	}

	// Trong thoi gian hoan: KHONG duoc thu lai (neu thu, hoanToi se nhich len).
	l.nhip(lucHong.Add(500 * time.Millisecond))
	if !l.muc[0].hoanToi.Equal(hoanCu) {
		t.Fatalf("dang trong thoi gian hoan ma van thu gui lai (hoanToi nhich tu %v sang %v)", hoanCu, l.muc[0].hoanToi)
	}

	// Het hoan + khu vuc da song lai -> boss ra, khong mat chu ky nao.
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })
	p := datNguoiQuanSat(t, k, 1)
	l.nhip(lucHong.Add(khoangHoanKhiGuiHong))
	chotActor(t, k)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("sau khi het hoan va khu vuc da chay, sinh %d boss, mong doi 1", n)
	}
}

// TestBoLichBossKhuVucDaDungThiNgungHan: khu vuc DA DUNG thi khong bao gio chay
// lai (Dung() dong kenh xong bang sync.Once) — thu lai mai mai chi lam ngap log
// dung luc tat may chu. Va quan trong hon: mot khu vuc chet KHONG duoc lam chet
// ca bo dieu phoi — cac muc khac van phai sinh binh thuong.
func TestBoLichBossKhuVucDaDungThiNgungHanNhungMucKhacVanChay(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	kSong := s.Tra(1)
	kChet := s.Tra(2)
	p := datNguoiQuanSat(t, kSong, 1)
	if err := kChet.Dung(); err != nil {
		t.Fatalf("dung khu vuc 2 that bai: %v", err)
	}

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, []MucLichBoss{
		mucMau(30, []int32{2}, []int64{500}), // muc 0 -> khu vuc DA CHET
		mucMau(30, []int32{1}, []int64{500}), // muc 1 -> khu vuc con song
	})
	l.muc[0].lanSinhCuoi = goc
	l.muc[1].lanSinhCuoi = goc

	l.nhip(goc.Add(30 * time.Second))
	if !l.muc[0].ngungHan {
		t.Error("khu vuc da dung ma muc lich khong bi ngung han — se ngap log luc tat may chu")
	}
	chotActor(t, kSong)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("mot khu vuc chet da lam chet ca bo dieu phoi: muc con lai sinh %d boss, mong doi 1", n)
	}

	// Da ngung han thi khong dung lai the song nua o cac nhip sau.
	l.nhip(goc.Add(120 * time.Second))
	if l.muc[0].the.ConSong() {
		t.Error("muc da ngung han van bat the song")
	}
}

// TestBoLichBossDuLieuHongThiNgungHan: du lieu sai (khong pha nao, khong map
// nao, khong mau nao, cap 0, map khong ton tai) la loi NGUOI VIET DU LIEU —
// thu lai mai khong bao gio dung, nen ngung han + log mot lan.
func TestBoLichBossDuLieuHongThiNgungHan(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	capHong := phaMau(0, []int32{1}, []int64{500})
	capHong.Cap = 0

	lich := []MucLichBoss{
		{DuLieu: nil},                        // 0: khong pha nao
		mucMau(0, nil, []int64{500}),         // 1: khong map nao
		mucMau(0, []int32{1}, nil),           // 2: khong mau nao
		{DuLieu: []PhaBoss{capHong}},         // 3: cap 0 -> client hien "cap 0"
		mucMau(0, []int32{99}, []int64{500}), // 4: map khong ton tai
	}
	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	l := MoBoLichBoss(s, lich)
	for i := range lich {
		l.muc[i].lanSinhCuoi = goc // GiayNghi 0 -> den han ngay, ke ca muc 4
	}
	l.nhip(goc)

	for i := range lich {
		if !l.muc[i].ngungHan {
			t.Errorf("muc %d du lieu hong ma khong bi ngung han", i)
		}
	}
	chotActor(t, k)
	if n := demGoi5(p); n != 0 {
		t.Fatalf("du lieu hong ma van sinh %d boss", n)
	}
}

// ---------------------------------------------------------------------------
// 6. Vong doi goroutine — BatDau hai lan, Dung truoc BatDau, Dung hai lan
// ---------------------------------------------------------------------------

// TestBatDauLichBossHaiLanChiChayMOT goroutine. Khong phai chuyen "dep code":
// ngay 28/08/2026 BoPhatThanh da dinh dung loi nay — hai goroutine cung `defer
// close(xong)` lam PANIC "close of closed channel" tren goroutine nen luc tat
// may chu, giet ca tien trinh. Vi vay test DEM goroutine, khong chi bat "khong
// panic".
func TestBatDauLichBossHaiLanChiChayMotGoroutine(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(3600, []int32{1}, []int64{500})})

	l.BatDau()
	time.Sleep(50 * time.Millisecond) // cho goroutine thu nhat vao vong lap
	truoc := runtime.NumGoroutine()

	l.BatDau() // lan hai — phai la viec khong
	time.Sleep(50 * time.Millisecond)
	sau := runtime.NumGoroutine()

	if sau > truoc {
		t.Fatalf("BatDau() lan hai da chay them goroutine: truoc=%d sau=%d", truoc, sau)
	}

	// Dung phai tra ve (khong treo) va khong panic dong kenh hai lan.
	l.Dung()
	l.Dung()
}

// TestDungLichBossAnToanKhiChuaBatDau: Dung() goi khi chua BatDau() phai tra ve
// NGAY. Thieu co daBatDau thi no treo vinh vien o <-l.xong — va vi tat may chu
// goi Dung() theo thu tu, ca chuoi tat se dung im khong ai hieu tai sao.
func TestDungLichBossAnToanKhiChuaBatDau(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})

	xong := make(chan struct{})
	go func() {
		l.Dung()
		l.Dung() // goi hai lan cung phai an toan
		close(xong)
	}()
	select {
	case <-xong:
	case <-time.After(2 * time.Second):
		t.Fatal("Dung() goi truoc BatDau() bi treo — tat may chu se dung im o day")
	}
}

// TestLichBossRongKhongLamGi: lich rong la trang thai HOP LE (may chu chua khai
// bao boss nao) — khong duoc panic, khong duoc sinh gi.
func TestLichBossRongKhongLamGi(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	l := MoBoLichBoss(s, nil)
	l.nhip(time.Now())
	l.BatDau()
	time.Sleep(250 * time.Millisecond)
	l.Dung()

	chotActor(t, k)
	if n := demGoi5(p); n != 0 {
		t.Fatalf("lich rong ma sinh %d boss", n)
	}
	if l.TheSong(0) != nil || l.TheSong(-1) != nil {
		t.Error("TheSong() ngoai mien phai tra nil")
	}
}

// ---------------------------------------------------------------------------
// KH6b Task 14 (phan hai) — chuyen NoiDungBoss (kho noi dung, noidungboss.go)
// sang MucLichBoss (dinh dang bo dieu phoi doc duoc). Hai kieu nay tach rieng
// co chu dich (xem chu thich noidungboss.go muc 4 va MucLichBoss o tren): kho
// noi dung con mang theo Ma + ViTri (chua dung toi trong lich), bo dieu phoi
// chi can DuLieu ([]PhaBoss). Ham chuyen doi la CAU noi giua hai file, nen dat
// ngay trong lichboss.go — noi dinh nghia dich (MucLichBoss) song.
// ---------------------------------------------------------------------------

// TestMucLichTuNoiDungBossGiuNguyenDuLieuPha khoa dung dan chuyen doi don gian
// nhat: MOI NoiDungBoss trong danh sach dau vao phai tro thanh MOT MucLichBoss
// cung so luong, DuLieu (cac pha) phai giu nguyen NOI DUNG — sai thu tu hoac
// thieu mot con la boss do se sinh voi pha sai (map sai, HP sai).
func TestMucLichTuNoiDungBossGiuNguyenDuLieuPha(t *testing.T) {
	ds := DanhSachNoiDungBoss()
	if len(ds) == 0 {
		t.Fatal("DanhSachNoiDungBoss() rong — khong con gi de kiem chuyen doi (kho noi dung Task 15 co the da doi)")
	}

	muc := MucLichTuNoiDungBoss(ds)
	if len(muc) != len(ds) {
		t.Fatalf("so muc lich = %d, muon = %d (moi con boss trong kho noi dung phai thanh MOT muc lich)",
			len(muc), len(ds))
	}
	for i, nd := range ds {
		if len(muc[i].DuLieu) != len(nd.Pha) {
			t.Fatalf("muc %d (%s): so pha = %d, muon = %d", i, nd.Ma, len(muc[i].DuLieu), len(nd.Pha))
		}
		for j, pha := range nd.Pha {
			if muc[i].DuLieu[j].Ten != pha.Ten || muc[i].DuLieu[j].Cap != pha.Cap {
				t.Fatalf("muc %d (%s) pha %d: du lieu KHONG khop kho noi dung goc (Ten/Cap doi)",
					i, nd.Ma, j)
			}
		}
	}
}

// TestMucLichTuNoiDungBossRongTraRong: dau vao rong phai tra lat cat rong,
// khong panic, khong tra nil khien noi goi phai kiem them mot lop.
func TestMucLichTuNoiDungBossRongTraRong(t *testing.T) {
	muc := MucLichTuNoiDungBoss(nil)
	if len(muc) != 0 {
		t.Fatalf("dau vao rong nhung tra ve %d muc", len(muc))
	}
}
