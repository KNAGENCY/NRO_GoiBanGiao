package game

import (
	"strings"
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 9 — MAY TRANG THAI VONG DOI BOSS trong actor.
//
//	NGHI -> HIEN -> CHAO -> HOAT_DONG -> (chet) -> TROI -> RA_DI -> (roi han)
//
// File nay khoa NAM dieu, moi dieu la mot cach he boss co the hong ma KHONG
// bao mot loi nao:
//
//  1. Boss chay HET vong doi. Hom nay (truoc task nay) no chet xong nam mai o
//     TROI: cai xac dung giua ban do vinh vien, khong ai don.
//  2. Khi roi han the gioi PHAI phat goi -6. Thieu no thi client giu cai xac
//     tren man hinh cho toi khi nguoi choi doi map — mot con ma khong danh
//     duoc, khong bien mat.
//  3. Khi roi han the gioi PHAI goi TheSongBoss.Tat(). Thieu no thi moi muc
//     lich boss sinh DUNG MOT LAN roi im lang mai mai — khong loi, khong log,
//     chi la "sao lau roi khong thay boss".
//  4. CHAN-7: may trang thai KHONG duoc tu quyet "khi nao" va "o map nao". No
//     khong doc dong ho lich, khong tu hoi sinh. Ai nhet lai nhanh do vao thi
//     ta co HAI co che sinh boss chong nhau.
//  5. Boss roi di phai duoc GO khoi k.boss. Khong go thi moi chu ky lich de
//     lai mot xac trong danh sach, va nhipTim (10 lan/giay) lap qua chung mai.
//
// KY THUAT: moi test goi k.nhipTim() TRUC TIEP tren mot khu vuc CHUA BatDau()
// (dung khuon quai_test.go/khuvuc_boss_test.go) — mot goroutine duy nhat, dong
// ho tiem duoc, khong sleep that. Rieng TestLichSinhDuocConThuHai chay tu
// TRONG actor that vi no phai chung minh duong noi giua HAI he thong.
// ─────────────────────────────────────────────────────────────────────────────

// khuVucBossGia dung mot khu vuc CHUA BatDau() de chay may trang thai boss:
// dong ho tro toi *luc (test doi "bay gio" khong sleep), MOT nguoi xem co phien
// gia de dem goi phat ra, va QUAI BI DON SACH.
//
// Vi sao phai `k.quai = nil`: nhipTim chay AI quai trong cung mot nhip. De
// nguyen dan quai mac dinh thi moi khang dinh "boss phat dung N goi" se bi tron
// them goi cua quai (-11/-10) va do vi mot ly do khong lien quan gi toi boss.
// Cac test quai lam dieu tuong tu theo chieu nguoc lai (gan thang k.quai).
func khuVucBossGia(luc *time.Time) (*KhuVuc, *phienGiaKhuVuc) {
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return *luc }
	k.quai = nil
	p := &phienGiaKhuVuc{}
	k.nguoiTrong[1] = &nguoiTrongKhuVuc{
		nv:    &store.NhanVat{ID: 1, Ten: "NguoiXem", HpToiDa: 100, HpHienTai: 100, ViTriX: 100, ViTriY: 384},
		phien: p,
	}
	return k, p
}

// bossVongDoi dung mot boss MOT PHA co du thoai de chay het vong doi. Trang
// thai + moc thoi gian vao trang thai do nguoi goi dat, vi moi test bat dau o
// mot buoc khac nhau.
func bossVongDoi(trangThai TrangThaiBoss, luc time.Time, thoaiHien, thoaiChet []string) *Boss {
	return &Boss{
		Id: -7,
		DuLieu: []PhaBoss{{
			Ten:       "Cửu Vĩ",
			Cap:       90,
			TypePk:    TypePkBoss,
			ThoaiHien: thoaiHien,
			ThoaiChet: thoaiChet,
		}},
		Hp:          1000,
		HpToiDa:     1000,
		ViTriX:      100,
		ViTriY:      384,
		TrangThai:   trangThai,
		MocThoiGian: luc,
		LucHienPha:  luc,
	}
}

// docChat44 doc nguoc mot goi 44 thanh (id nguoi noi, noi dung).
func docChat44(t *testing.T, m *protocol.Message) (int32, string) {
	t.Helper()
	if m == nil {
		t.Fatal("goi 44 la nil")
	}
	if m.Command != lenhChat {
		t.Fatalf("command = %d, mong doi %d (44 chat)", m.Command, lenhChat)
	}
	r := protocol.NewReader(m.Data)
	id, err := r.Int()
	if err != nil {
		t.Fatalf("doc id nguoi noi loi: %v", err)
	}
	noiDung, err := r.UTF()
	if err != nil {
		t.Fatalf("doc noi dung loi: %v", err)
	}
	return id, noiDung
}

// locGoi tra ve moi goi co ma lenh da cho ma phien gia nhan duoc.
func locGoi(p *phienGiaKhuVuc, lenh int8) []*protocol.Message {
	var ra []*protocol.Message
	for _, m := range p.banTin() {
		if m.Command == lenh {
			ra = append(ra, m)
		}
	}
	return ra
}

// cauDaNoi tra ve danh sach cau boss da noi, theo dung thu tu phat ra.
func cauDaNoi(t *testing.T, p *phienGiaKhuVuc) []string {
	t.Helper()
	var ra []string
	for _, m := range locGoi(p, lenhChat) {
		_, cau := docChat44(t, m)
		ra = append(ra, cau)
	}
	return ra
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 1 — hai ham THUAN ve nhip thoai
// ─────────────────────────────────────────────────────────────────────────────

// TestTreoThoaiTheoDoDaiCau: cau cang dai thi treo tren dau cang lau (100ms mot
// ky tu), nhung KHONG BAO GIO qua 2 giay.
//
// Hai con so nay khong phai em nghi ra — chung la luat cua ban Java goc,
// Boss.java:507-510 (`timeChatS = textChat.length() * 100; if > 2000 -> 2000`).
// Ghim o day de ai doi con biet minh dang doi mot con so co nguon.
//
// DEM THEO KY TU (rune) chu khong theo BYTE: "Đủ rồi." la 7 ky tu nhung 10 byte
// UTF-8. Dem byte thi moi cau tieng Viet treo lau hon cau tieng Anh cung do dai
// khoang 40% — mot su thien vi vo nghia sinh ra tu chi tiet ma hoa.
func TestTreoThoaiTheoDoDaiCau(t *testing.T) {
	if got := treoThoai(""); got != 0 {
		t.Errorf("treoThoai(\"\") = %v, mong doi 0 (khong co gi de doc)", got)
	}
	if got := treoThoai("abcdef"); got != 600*time.Millisecond {
		t.Errorf("treoThoai(6 ky tu) = %v, mong doi 600ms", got)
	}
	if got := treoThoai("Đủ rồi."); got != 700*time.Millisecond {
		t.Errorf("treoThoai(%q) = %v, mong doi 700ms — 7 KY TU, khong phai 10 byte", "Đủ rồi.", got)
	}
	if got := treoThoai(strings.Repeat("a", 500)); got != 2*time.Second {
		t.Errorf("treoThoai(cau rat dai) = %v, mong doi kep o 2s", got)
	}
}

// TestThoaiKeTiepNoiTungCauMotTheoNhip do ham thuan quyet dinh "tick nay co noi
// khong, va noi cau nao".
//
// Ba dieu duoc khoa:
//   - Cau DAU TIEN noi ngay, khong cho (chua co cau nao treo trên đầu).
//   - Cau sau phai doi cau truoc treo du lau.
//   - "Xong" chi dung SAU KHI cau CUOI da treo du — khong phai ngay khi vua
//     noi xong cau cuoi. Sai cho nay thi cau chot cua boss vua hien ra da bi
//     buoc HOAT_DONG dam len.
func TestThoaiKeTiepNoiTungCauMotTheoNhip(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	cau := []string{"abcde", "xyz"} // treo 500ms / 300ms

	if i, xong := thoaiKeTiep(cau, 0, time.Time{}, t0); i != 0 || xong {
		t.Errorf("chiSo 0: (%d,%v), mong doi (0,false) — cau dau noi NGAY", i, xong)
	}
	if i, xong := thoaiKeTiep(cau, 1, t0, t0.Add(400*time.Millisecond)); i != -1 || xong {
		t.Errorf("moi 400ms sau cau 500ms: (%d,%v), mong doi (-1,false) — con dang treo", i, xong)
	}
	if i, xong := thoaiKeTiep(cau, 1, t0, t0.Add(500*time.Millisecond)); i != 1 || xong {
		t.Errorf("du 500ms: (%d,%v), mong doi (1,false)", i, xong)
	}
	if i, xong := thoaiKeTiep(cau, 2, t0, t0.Add(200*time.Millisecond)); i != -1 || xong {
		t.Errorf("cau cuoi con dang treo: (%d,%v), mong doi (-1,false) — CHUA duoc bao xong", i, xong)
	}
	if i, xong := thoaiKeTiep(cau, 2, t0, t0.Add(300*time.Millisecond)); i != -1 || !xong {
		t.Errorf("cau cuoi treo du: (%d,%v), mong doi (-1,true)", i, xong)
	}
	if i, xong := thoaiKeTiep(nil, 0, time.Time{}, t0); i != -1 || !xong {
		t.Errorf("khong co cau nao: (%d,%v), mong doi (-1,true) — xong ngay", i, xong)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2 — HIEN -> CHAO -> HOAT_DONG
// ─────────────────────────────────────────────────────────────────────────────

// TestBossHienBuocSangChao: goi -5 da do themBoss phat truoc do; viec cua may
// trang thai o buoc HIEN chi la buoc sang CHAO. Neu no dung lai o HIEN thi boss
// khong bao gio noi cau nao va khong bao gio danh ai.
func TestBossHienBuocSangChao(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, _ := khuVucBossGia(&luc)
	b := bossVongDoi(BossHien, luc, []string{"Làng Lá vẫn ngủ say."}, nil)
	k.boss = append(k.boss, b)

	k.nhipTim()

	if b.TrangThai != BossChao {
		t.Fatalf("sau mot nhip, boss o trang thai %s, mong doi CHAO", b.TrangThai)
	}
	if !b.MocThoiGian.Equal(luc) {
		t.Errorf("MocThoiGian = %v, mong doi dat lai = luc vao trang thai moi (%v)", b.MocThoiGian, luc)
	}
}

// TestBossChaoNoiTungCauRoiSangHoatDong la test chinh cua buoc CHAO: hai cau
// thoai phai ra LAN LUOT (khong don cuc mot luc), qua goi 44 mang ID BOSS, va
// chi khi cau cuoi da treo du thi boss moi buoc sang HOAT_DONG.
func TestBossChaoNoiTungCauRoiSangHoatDong(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	// "abcde" treo 500ms, "xyz" treo 300ms.
	b := bossVongDoi(BossChao, luc, []string{"abcde", "xyz"}, nil)
	k.boss = append(k.boss, b)

	k.nhipTim() // cau 1 ra ngay
	if cau := cauDaNoi(t, p); len(cau) != 1 || cau[0] != "abcde" {
		t.Fatalf("nhip 1: da noi %v, mong doi [abcde]", cau)
	}
	if id, _ := docChat44(t, locGoi(p, lenhChat)[0]); id != -7 {
		t.Errorf("goi 44 mang id %d, mong doi -7 (ID CUA BOSS, khong phai cua ai khac)", id)
	}

	luc = t0.Add(400 * time.Millisecond)
	k.nhipTim() // chua du 500ms -> chua noi them
	if cau := cauDaNoi(t, p); len(cau) != 1 {
		t.Fatalf("moi 400ms ma da noi %v — cau truoc chua kip doc", cau)
	}
	if b.TrangThai != BossChao {
		t.Fatalf("boss dang chao do ma da sang %s", b.TrangThai)
	}

	luc = t0.Add(500 * time.Millisecond)
	k.nhipTim() // cau 2
	if cau := cauDaNoi(t, p); len(cau) != 2 || cau[1] != "xyz" {
		t.Fatalf("nhip 3: da noi %v, mong doi [abcde xyz]", cau)
	}
	if b.TrangThai != BossChao {
		t.Fatalf("vua noi xong cau cuoi da sang %s — cau chot chua kip doc", b.TrangThai)
	}

	luc = t0.Add(800 * time.Millisecond) // 500 + 300, cau cuoi treo du
	k.nhipTim()
	if b.TrangThai != BossHoatDong {
		t.Fatalf("noi het thoai roi ma van o %s, mong doi HOAT_DONG", b.TrangThai)
	}
	if cau := cauDaNoi(t, p); len(cau) != 2 {
		t.Fatalf("da noi %v — khong duoc noi lai cau nao khi da sang HOAT_DONG", cau)
	}
}

// TestBossChaoKhongCoThoaiSangHoatDongNgay: mot pha khong khai bao ThoaiHien
// khong duoc lam boss ket o buoc CHAO. Day la kieu du lieu that (boss le, im
// lang) chu khong phai truong hop benh hoan.
func TestBossChaoKhongCoThoaiSangHoatDongNgay(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, p := khuVucBossGia(&luc)
	b := bossVongDoi(BossChao, luc, nil, nil)
	k.boss = append(k.boss, b)

	k.nhipTim()

	if b.TrangThai != BossHoatDong {
		t.Fatalf("boss khong co thoai van o %s, mong doi HOAT_DONG ngay", b.TrangThai)
	}
	if cau := cauDaNoi(t, p); len(cau) != 0 {
		t.Fatalf("boss khong co thoai ma noi %v", cau)
	}
}

// TestBossHoatDongDungYen_ChoTask10: o ban nay buoc HOAT_DONG chua lam gi ca
// (AI danh/di chuyen/thoai-danh la Task 10). Test nay ghim NGUYEN TRANG do lai
// de khi Task 10 lam that, no phai SUA test nay — tuc la co mot cho bat buoc
// dung lai doc, thay vi lang le them hanh vi vao mot buoc khong ai kiem.
func TestBossHoatDongDungYen_ChoTask10(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossVongDoi(BossHoatDong, luc, nil, nil)
	b.DuLieu[0].ThoaiDanh = []string{"Chakra của ngươi run rồi đấy."}
	k.boss = append(k.boss, b)

	for i := 0; i < 100; i++ { // 10 giay nhip
		luc = t0.Add(time.Duration(i) * 100 * time.Millisecond)
		k.nhipTim()
	}

	if b.TrangThai != BossHoatDong {
		t.Fatalf("boss dang HOAT_DONG tu doi sang %s — khong ai duoc quyen doi buoc nay tru khi bi danh chet", b.TrangThai)
	}
	if n := len(p.banTin()); n != 0 {
		t.Fatalf("buoc HOAT_DONG phat %d goi — Task 9 chua lam AI, phai im lang (Task 10 se doi test nay)", n)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 3 — TROI (cai xac) -> RA_DI -> roi han the gioi
// ─────────────────────────────────────────────────────────────────────────────

// TestBossTroiNoiLoiTrangTroi: buoc TROI phat ThoaiChet, cung co che nhip voi
// buoc CHAO.
func TestBossTroiNoiLoiTrangTroi(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, p := khuVucBossGia(&luc)
	b := bossVongDoi(BossTroi, luc, nil, []string{"Đó mới chỉ là lớp áo..."})
	b.Hp = 0
	k.boss = append(k.boss, b)

	k.nhipTim()

	if cau := cauDaNoi(t, p); len(cau) != 1 || cau[0] != "Đó mới chỉ là lớp áo..." {
		t.Fatalf("loi trang troi: da noi %v", cau)
	}
}

// TestXacBossNamDuBaGiayMoiRaDi khoa THOI GIAN TOI THIEU cai xac nam lai.
//
// Vi sao PHAI co con so nay chu khong de "noi xong thoai la di": mot pha khong
// khai bao ThoaiChet se bien mat NGAY trong nhip boss chet. Nguoi choi vua ha
// duoc boss, nhin len thi no da khong con — khong kip thay minh vua lam gi. Te
// hon: dong chu "X vua ha guc Boss Y" di qua hang doi loa voi nhip 1,5 giay
// (nhipToiThieuThongBao, thongbao.go) nen no toi SAU khi cai xac da bien mat —
// ca may chu doc tin ve mot con boss khong con dau vet.
//
// 3 giay = hai lan nhip loa, du de dong chu toi noi va nguoi choi con thay xac.
func TestXacBossNamDuBaGiayMoiRaDi(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossVongDoi(BossTroi, luc, nil, nil) // KHONG co thoai chet
	b.Hp = 0
	k.boss = append(k.boss, b)

	luc = t0.Add(2900 * time.Millisecond)
	k.nhipTim()
	if b.TrangThai != BossTroi {
		t.Fatalf("moi 2,9 giay ma xac da sang %s — chua du 3 giay", b.TrangThai)
	}
	if n := len(locGoi(p, lenhRoiKhuVuc)); n != 0 {
		t.Fatalf("moi 2,9 giay ma da phat %d goi -6", n)
	}

	luc = t0.Add(3 * time.Second)
	k.nhipTim()
	if b.TrangThai != BossRaDi {
		t.Fatalf("du 3 giay ma xac van o %s, mong doi RA_DI", b.TrangThai)
	}
	if len(k.boss) != 1 {
		t.Fatalf("boss bi go khoi danh sach ngay o buoc RA_DI — phai doi toi khi phat -6 xong")
	}
}

// TestBossHetPhaRoiHanTheGioi la TEST QUAN TRONG NHAT file nay: mot con boss
// het pha, den buoc RA_DI, phai lam DU BA VIEC va bo sot bat ky cai nao cung
// hong theo mot kieu rieng:
//
//	phat -6      -> thieu: client giu cai xac boss tren man hinh VINH VIEN
//	goi The.Tat() -> thieu: muc lich do khong bao gio sinh con thu hai nua
//	go khoi k.boss -> thieu: moi chu ky lich de lai mot xac, nhipTim lap mai
func TestBossHetPhaRoiHanTheGioi(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, p := khuVucBossGia(&luc)

	the := &TheSongBoss{}
	the.Bat() // bo dieu phoi da danh dau "con nay dang song"

	b := bossVongDoi(BossRaDi, luc, nil, nil)
	b.The = the
	k.boss = append(k.boss, b)

	k.nhipTim()

	goi6 := locGoi(p, lenhRoiKhuVuc)
	if len(goi6) != 1 {
		t.Fatalf("phat %d goi -6, mong doi DUNG 1 — thieu no thi client giu cai xac boss mai mai", len(goi6))
	}
	r := protocol.NewReader(goi6[0].Data)
	id, err := r.Int()
	if err != nil {
		t.Fatalf("doc id trong goi -6 loi: %v", err)
	}
	if id != -7 {
		t.Errorf("goi -6 mang id %d, mong doi -7 (id cua chinh con boss)", id)
	}

	if the.ConSong() {
		t.Fatal("TheSongBoss van BAT sau khi boss roi the gioi — muc lich nay se khong bao gio " +
			"sinh con thu hai nua (khong loi, khong log, chi la 'sao lau roi khong thay boss')")
	}
	if len(k.boss) != 0 {
		t.Fatalf("con %d boss trong k.boss sau khi da roi the gioi, mong doi 0", len(k.boss))
	}
	if b.TrangThai != BossNghi {
		t.Errorf("trang thai cuoi = %s, mong doi NGHI (thiet ke §5: het pha -> -6 -> ve NGHI)", b.TrangThai)
	}
}

// TestBossChuaTroiDuThiChuaPhatGoi6 chan mot dot bien rat de lot: nguoi sua bo
// buoc TROI di, cho boss chet la RA_DI luon. Luc do -6 van duoc phat va The van
// duoc tat — moi test o tren van xanh — nhung cai xac bien mat trong 100ms.
func TestBossChuaTroiDuThiChuaPhatGoi6(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	the := &TheSongBoss{}
	the.Bat()
	b := bossVongDoi(BossTroi, luc, nil, nil)
	b.The = the
	b.Hp = 0
	k.boss = append(k.boss, b)

	for i := 0; i < 29; i++ { // 29 nhip = 2,9 giay, van duoi nguong
		luc = t0.Add(time.Duration(i) * 100 * time.Millisecond)
		k.nhipTim()
	}

	if n := len(locGoi(p, lenhRoiKhuVuc)); n != 0 {
		t.Fatalf("da phat %d goi -6 khi xac moi nam 2,9 giay", n)
	}
	if !the.ConSong() {
		t.Fatal("The da bi tat khi boss VAN CON tren man hinh — bo dieu phoi se sinh con thu hai " +
			"trong luc con thu nhat chua don di")
	}
	if len(k.boss) != 1 {
		t.Fatalf("boss bi go khoi danh sach qua som, con %d", len(k.boss))
	}
}

// TestBossChetGiuaLucChaoVanNoiDuLoiTrangTroi: nguoi choi xong vao danh ngay
// khi boss dang chao. Loi trang troi PHAI bat dau tu cau DAU TIEN — day la ly
// do ChiSoThoaiHien va ChiSoThoaiChet la HAI con tro rieng (boss.go). Dung
// chung mot con tro thi boss chet giua cau chao thu hai se nuot mat cau trang
// troi dau tien.
func TestBossChetGiuaLucChaoVanNoiDuLoiTrangTroi(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossVongDoi(BossChao, luc, []string{"chao1", "chao2"}, []string{"chet1", "chet2"})
	k.boss = append(k.boss, b)

	k.nhipTim() // noi "chao1"
	if cau := cauDaNoi(t, p); len(cau) != 1 {
		t.Fatalf("chuan bi: da noi %v, mong doi [chao1]", cau)
	}

	// Chet ngay giua chung.
	k.bossChet(b, 1, luc)
	if b.TrangThai != BossTroi {
		t.Fatalf("sau bossChet, trang thai = %s, mong doi TROI", b.TrangThai)
	}

	luc = t0.Add(600 * time.Millisecond) // "chao1" (5 ky tu) da treo du tu lau
	k.nhipTim()

	cau := cauDaNoi(t, p)
	if len(cau) != 2 || cau[1] != "chet1" {
		t.Fatalf("da noi %v, mong doi cau thu hai la \"chet1\" — loi trang troi phai bat dau tu CAU DAU", cau)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 4 — CHAN-7: may trang thai KHONG duoc tu quyet "khi nao" va "o map nao"
// ─────────────────────────────────────────────────────────────────────────────

// TestChotCHAN7_MayTrangThaiKhongTuSinhBoss la CHOT KIEM cua ke hoach muc
// CHAN-7. Sau khi mot con boss roi the gioi, khu vuc phai IM LANG MAI MAI —
// bat ke GiayNghi trong du lieu la bao nhieu va dong ho da chay bao lau.
//
// "Khi nao sinh lai" la cau hoi cua bo dieu phoi lich (lichboss.go, NGOAI
// actor). Neu may trang thai cung tra loi no, ta co HAI co che sinh boss chong
// nhau va khong ai biet con boss dang song la do ai sinh ra.
func TestChotCHAN7_MayTrangThaiKhongTuSinhBoss(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)

	b := bossVongDoi(BossRaDi, luc, nil, nil)
	b.DuLieu[0].GiayNghi = 1 // mot giay: neu actor doc con so nay, no se lo ngay
	b.DuLieu[0].MapVao = []int32{1}
	b.DuLieu[0].HpChon = []int64{500}
	k.boss = append(k.boss, b)

	k.nhipTim() // boss roi the gioi
	if len(k.boss) != 0 {
		t.Fatalf("chuan bi: boss chua roi the gioi, con %d", len(k.boss))
	}
	soGoiSauKhiRoi := len(p.banTin())

	// Nap them MOT con dang o buoc NGHI voi GiayNghi = 1 giay. Day la cai bay
	// truc tiep: neu may trang thai co ay may doc GiayNghi de tu danh thuc boss,
	// con nay se bat day chi sau mot giay dong ho.
	nghi := bossVongDoi(BossNghi, luc, nil, nil)
	nghi.Id = -8
	nghi.DuLieu[0].GiayNghi = 1
	nghi.DuLieu[0].MapVao = []int32{1}
	nghi.DuLieu[0].HpChon = []int64{500}
	k.themBoss(nghi)

	// Mot gio dong ho troi qua — 3600 lan qua nguong GiayNghi.
	for i := 1; i <= 200; i++ {
		luc = t0.Add(time.Duration(i) * 20 * time.Second)
		k.nhipTim()
	}

	if len(k.boss) != 0 {
		t.Fatalf("may trang thai TU SINH %d boss sau khi qua GiayNghi — CHAN-7 cam viec do: "+
			"'khi nao + o map nao' la viec cua bo dieu phoi lich, NGOAI actor", len(k.boss))
	}
	if n := len(locGoi(p, lenhXuatHien)); n != 0 {
		t.Fatalf("may trang thai phat %d goi -5 — no da tu sinh boss", n)
	}
	if n := len(p.banTin()); n != soGoiSauKhiRoi {
		t.Fatalf("khu vuc phat them %d goi sau khi boss da roi di, mong doi im lang tuyet doi",
			n-soGoiSauKhiRoi)
	}
}

// TestKhuVucKhongCoBossThiNhipTimImLang: canh gac re tien nhat cua CHAN-7 —
// mot khu vuc chua tung co boss nao khong bao gio tu co.
func TestKhuVucKhongCoBossThiNhipTimImLang(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)

	for i := 1; i <= 100; i++ {
		luc = t0.Add(time.Duration(i) * time.Minute)
		k.nhipTim()
	}

	if len(k.boss) != 0 {
		t.Fatalf("khu vuc trong tu sinh %d boss", len(k.boss))
	}
	if n := len(p.banTin()); n != 0 {
		t.Fatalf("khu vuc trong phat %d goi", n)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 5 — don rac + ben truoc du lieu hong
// ─────────────────────────────────────────────────────────────────────────────

// TestBossONghiBiGoKhoiDanhSach: mot con boss o trang thai NGHI nam trong actor
// la RAC theo dinh nghia — CHAN-7 cam actor dem nguoc GiayNghi de danh thuc no,
// nen no khong co duong nao tro lai the gioi. Giu lai chi ton bo nho va ton mot
// vong lap moi 100ms.
func TestBossONghiBiGoKhoiDanhSach(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, p := khuVucBossGia(&luc)
	k.themBoss(bossVongDoi(BossNghi, luc, nil, nil))
	if len(k.boss) != 1 {
		t.Fatalf("chuan bi: themBoss phai giu con boss NGHI, len=%d", len(k.boss))
	}

	k.nhipTim()

	if len(k.boss) != 0 {
		t.Fatalf("boss NGHI van nam trong k.boss sau mot nhip, len=%d", len(k.boss))
	}
	if n := len(p.banTin()); n != 0 {
		t.Fatalf("go mot con boss chua bao gio hien ra ma phat %d goi — client chua tung thay no", n)
	}
}

// TestBossDuLieuHongDuocDonVaMoLaiCuaLich: mot con boss khong co pha nao dang
// chay (du lieu hong / da qua pha cuoi) khong the lam gi duoc nua. Neu cu de no
// nam do thi The KHONG BAO GIO duoc tat -> muc lich do chet vinh vien vi mot
// dong du lieu go sai.
func TestBossDuLieuHongDuocDonVaMoLaiCuaLich(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, _ := khuVucBossGia(&luc)

	the := &TheSongBoss{}
	the.Bat()
	b := &Boss{Id: -9, DuLieu: nil, TrangThai: BossHoatDong, The: the, MocThoiGian: luc}
	k.boss = append(k.boss, b)

	k.nhipTim()

	if len(k.boss) != 0 {
		t.Fatalf("boss du lieu hong van nam trong danh sach, len=%d", len(k.boss))
	}
	if the.ConSong() {
		t.Fatal("The van bat sau khi don mot con boss hong — muc lich do chet vinh vien")
	}
}

// TestNhipBossKhongPanicVoiConTroNil: mot o nil lac vao k.boss khong duoc lam
// sap nhip cua ca khu vuc (moi thu khac trong khu vuc dung theo).
func TestNhipBossKhongPanicVoiConTroNil(t *testing.T) {
	luc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	k, _ := khuVucBossGia(&luc)
	k.boss = append(k.boss, nil, bossVongDoi(BossHien, luc, nil, nil), nil)

	k.nhipTim() // khong duoc panic

	if len(k.boss) != 1 {
		t.Fatalf("con %d phan tu trong k.boss, mong doi 1 (hai o nil bi don)", len(k.boss))
	}
}

// TestVongDoiBossTronVen chay HET mot vong doi tren cung mot con boss, dung thu
// tu thiet ke §5. Cac test tren khoa tung buoc; test nay khoa DUONG NOI giua
// chung — thieu no thi moi buoc dung rieng le van xanh trong khi ca chuoi ket.
func TestVongDoiBossTronVen(t *testing.T) {
	t0 := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)

	the := &TheSongBoss{}
	the.Bat()
	b := bossVongDoi(BossHien, luc, []string{"chao"}, []string{"tam biet"})
	b.The = the
	k.boss = append(k.boss, b)

	buoc := func(sau time.Duration, mongDoi TrangThaiBoss, ghiChu string) {
		t.Helper()
		luc = t0.Add(sau)
		k.nhipTim()
		if b.TrangThai != mongDoi {
			t.Fatalf("%s: trang thai = %s, mong doi %s", ghiChu, b.TrangThai, mongDoi)
		}
	}

	buoc(0, BossChao, "HIEN -> CHAO")
	buoc(100*time.Millisecond, BossChao, "CHAO: noi cau chao")
	buoc(500*time.Millisecond, BossHoatDong, "CHAO -> HOAT_DONG (cau 4 ky tu treo 400ms)")

	// Nguoi choi ha boss o giay thu 2.
	luc = t0.Add(2 * time.Second)
	k.bossChet(b, 1, luc)
	if b.TrangThai != BossTroi {
		t.Fatalf("sau bossChet: %s, mong doi TROI", b.TrangThai)
	}

	buoc(2*time.Second+100*time.Millisecond, BossTroi, "TROI: noi loi trang troi")
	// Xac nam du 3 giay tinh tu luc CHET (giay thu 2), khong phai tu luc noi
	// xong loi trang troi — nen moc la giay thu 5.
	buoc(4*time.Second, BossTroi, "TROI: giay thu 4 van chua du 3 giay nam xac")
	buoc(5*time.Second, BossRaDi, "TROI -> RA_DI (xac nam du 3 giay)")

	luc = t0.Add(5100 * time.Millisecond)
	k.nhipTim()
	if len(k.boss) != 0 {
		t.Fatalf("RA_DI xong ma boss van con trong danh sach, len=%d", len(k.boss))
	}
	if the.ConSong() {
		t.Fatal("het vong doi ma The van bat")
	}
	if cau := cauDaNoi(t, p); len(cau) != 2 || cau[0] != "chao" || cau[1] != "tam biet" {
		t.Fatalf("ca vong doi noi %v, mong doi [chao \"tam biet\"]", cau)
	}
	if n := len(locGoi(p, lenhRoiKhuVuc)); n != 1 {
		t.Fatalf("phat %d goi -6 trong ca vong doi, mong doi dung 1", n)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 6 — duong noi that voi bo dieu phoi lich (lichboss.go)
// ─────────────────────────────────────────────────────────────────────────────

// TestLichBossGanTheSongVaoConBoss: bo dieu phoi PHAI cam the song cua muc lich
// vao chinh con boss no de ra. Day la ca soi day — thieu no thi may trang thai
// cam mot con tro nil va Tat() tro thanh mot lenh khong lam gi (nil-safe), im
// lang tuyet doi.
func TestLichBossGanTheSongVaoConBoss(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	m := l.muc[0]

	b := l.taoBoss(m, m.phaDau(), s.Tra(1).BanDo(), time.Now())

	if b.The == nil {
		t.Fatal("boss sinh ra khong mang the song cua muc lich — may trang thai se khong co gi de tat")
	}
	if b.The != m.the {
		t.Fatal("boss mang MOT the song khac voi the song cua muc lich da de ra no")
	}
}

// TestLichSinhDuocConThuHaiSauKhiConDauRoiDi la BANG CHUNG DAU-CUOI cua ca
// task: hai he thong (bo dieu phoi NGOAI actor + may trang thai TRONG actor)
// noi duoc vao nhau.
//
// Truoc task nay, moi muc lich sinh DUNG MOT LAN roi thoi (chu thich
// BoLichBoss.TheSong ghi ro dieu do). Test nay do dung cai da doi.
//
// KY THUAT: khu vuc lay tu MoSoKhuVuc() DA CHAY, nen moi thao tac len trang
// thai cua no (gan dong ho, giet boss, goi nhipTim) deu duoc gui vao TRONG actor
// bang suKienChay — khong mot dong nao cham k.boss/k.dongHo tu goroutine test.
func TestLichSinhDuocConThuHaiSauKhiConDauRoiDi(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	goc := time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)
	luc := goc // CHI goroutine actor duoc cham bien nay

	trongActor := func(f func()) {
		t.Helper()
		k.Gui(suKienChay{f: f})
		chotActor(t, k)
	}
	trongActor(func() { k.dongHo = func() time.Time { return luc }; k.quai = nil })

	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	l.muc[0].lanSinhCuoi = goc

	l.nhip(goc.Add(30 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("lan sinh dau: %d boss, mong doi 1", n)
	}

	// Ha con thu nhat roi de may trang thai lam not viec cua no.
	trongActor(func() {
		for _, b := range k.boss {
			k.bossChet(b, 1, k.dongHo())
		}
	})
	trongActor(func() { luc = goc.Add(40 * time.Second) }) // xac nam qua 3 giay
	trongActor(k.nhipTim)                                  // TROI -> RA_DI
	trongActor(k.nhipTim)                                  // RA_DI -> phat -6, tat the, go khoi danh sach

	if l.TheSong(0).ConSong() {
		t.Fatal("boss da roi the gioi ma the song van bat — may trang thai chua goi Tat()")
	}

	l.nhip(goc.Add(120 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 2 {
		t.Fatalf("tong %d goi -5, mong doi 2 — bo dieu phoi phai sinh duoc con THU HAI sau khi "+
			"may trang thai bao con thu nhat da roi di", n)
	}
}
