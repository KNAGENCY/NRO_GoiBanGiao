// lichboss.go — BO DIEU PHOI LICH BOSS, chay NGOAI actor (KH6b Task 13,
// ke hoach 2026-08-22-kh6b-boss-ke-hoach.md muc CHAN-3 + CHAN-7).
//
// LY DO TON TAI: cau hoi "khi nao sinh boss" va "sinh o map nao" KHONG duoc
// tra loi ben trong actor khu vuc. Hai ly do rat cu the:
//
//  1. (CHAN-7) Actor khong duoc tu quyet spawn. Neu may trang thai trong
//     nhipTim tu doc dong ho roi tu hoi sinh, ta co HAI co che sinh boss chong
//     nhau (lich + tu hoi sinh) va khong ai biet con boss dang song la do ai
//     sinh ra. Actor CHI lo chuyen pha SAU KHI da nhan lenh sinh.
//  2. (CHAN-3) MoSoKhuVuc() phai GIU NGUYEN hanh vi "5 khu vuc im lang".
//     Hang loat test tich hop dung no de dung the gioi roi khang dinh khong co
//     goi nao ngoai du kien; mot con boss tu nhien hien ra giua kich ban se
//     lam do het ma khong ai hieu vi sao. Vi vay lich boss CHI duoc bat tu
//     cmd/nroserver/main.go (Task 14) — TUYET DOI khong tu MoSoKhuVuc.
//
// CHIEU DI CUA DU LIEU — mot chieu, khong ai cho ai:
//
//	BoLichBoss.vongLap (goroutine rieng)        Actor khu vuc
//	────────────────────────────────────        ─────────────
//	den gio -> dung san *Boss (chon map,   ──►  hop thu DIEU KHIEN
//	  chon HP, cap id am) -> k.SinhBoss()       -> themBoss() -> phat -5
//	CHUYEN GIAO QUYEN SO HUU con tro do,
//	khong bao gio cham lai nua
//
// Chieu nguoc lai (actor bao "con boss nay di roi") CHI duoc phep la MOT phep
// ghi nguyen tu: TheSongBoss.Tat(). Khong kenh, khong mutex, khong cho — vi
// bay §6.10b: bat ky loi goi CHAN nao tu trong actor deu co the treo ca khu vuc.
package game

import (
	"errors"
	"log/slog"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"nroserver/internal/antoan"
)

// nhipKiemLichBoss la tan suat DRIVER THAT (ticker cua vongLap) — KHONG PHAI
// chu ky sinh boss (cai do la GiayNghi trong du lieu tung pha). Chi la "bao lau
// nhin dong ho mot lan", dung y het nhipKiemLuuDinhKy (luudinhky.go) va
// nhipKiemPhatThanh (thongbao.go): mot nhip khi chua den gio chi la vai phep so
// sanh thoi gian, chi phi khong dang ke.
const nhipKiemLichBoss = 100 * time.Millisecond

// khoangHoanKhiGuiHong la thoi gian NGHI truoc khi thu gui lai mot lenh sinh
// da that bai tam thoi (khu vuc dang ban -> ErrControlQuaHan, hoac chua
// BatDau() -> ErrActorChuaBatDau).
//
// Vi sao PHAI co con so nay chu khong thu lai ngay nhip sau: guiDieuKhien CHO
// TOI 100ms moi lan that bai vi qua han. Thu lai moi nhip 100ms nghia la bo
// dieu phoi dung gan het thoi gian cua no de go cua mot khu vuc dang qua tai —
// va neu ca 5 khu vuc cung ket thi mot "nhip" cua bo dieu phoi keo dai nua
// giay. 1 giay la du dai de khu vuc kip tho, du ngan de khong ai thay boss ra
// tre.
const khoangHoanKhiGuiHong = 1 * time.Second

// TheSongBoss la SOI DAY DUY NHAT noi actor nguoc ve bo dieu phoi: "con boss
// cua muc lich nay con trong the gioi khong?".
//
// Vi sao la atomic.Bool chu khong phai mot kenh hay mot callback: ham Tat() se
// duoc goi TU TRONG actor (may trang thai boss, Task 9, luc boss het pha va roi
// khoi the gioi). Moi thu actor goi ra ngoai deu phai la thu KHONG BAO GIO CHAN
// — mot phep ghi nguyen tu la muc do an toan cao nhat co the (bay #1 cua brief:
// actor khong duoc goi nguoc vao bo dieu phoi roi cho).
//
// TRANG THAI BAN DAU la TAT (zero value = false = khong co boss nao song) — dung
// truc giac: muc lich vua tao thi chua sinh con nao.
type TheSongBoss struct{ song atomic.Bool }

// Bat danh dau "muc nay dang co mot con boss trong the gioi". CHI bo dieu phoi
// goi, ngay sau khi gui lenh sinh THANH CONG.
func (t *TheSongBoss) Bat() { t.song.Store(true) }

// Tat danh dau "con boss cua muc nay da roi khoi the gioi" — bo dieu phoi duoc
// phep sinh con moi khi den chu ky ke tiep.
//
// AN TOAN goi tu TRONG actor: mot phep ghi nguyen tu, khong cap phat, khong
// cham I/O, khong bao gio chan. nil-safe de noi goi khong phai kiem truoc.
func (t *TheSongBoss) Tat() {
	if t == nil {
		return
	}
	t.song.Store(false)
}

// ConSong tra ve co dang bat hay khong. nil-safe: khong co the thi coi nhu
// khong co boss nao song.
func (t *TheSongBoss) ConSong() bool {
	if t == nil {
		return false
	}
	return t.song.Load()
}

// MucLichBoss la MOT dong trong lich: mot con boss + noi/luc no xuat hien.
//
// KHAI BAO TOI THIEU CO CHU DICH (Task 15 moi la noi nap du lieu boss that):
// mot muc lich CHI la []PhaBoss — dung y kieu du lieu boss da co san (boss.go),
// khong dinh nghia them mot lop cau hinh song song nao. Ba thu bo dieu phoi can
// deu doc tu PHA DAU TIEN:
//
//   - GiayNghi -> chu ky sinh (bao lau ke tu lan sinh truoc thi sinh lai)
//   - MapVao   -> danh sach khu vuc, chon NGAU NHIEN mot cai
//   - HpChon   -> danh sach mau, chon NGAU NHIEN mot gia tri
//
// Vi sao doc tu pha DAU chu khong phai pha dang chay: pha sau chi ton tai khi
// boss da o trong the gioi (bien hinh), luc do bo dieu phoi khong con quyen gi
// voi no nua. "Khi nao va o dau" chi duoc hoi MOT LAN, luc con boss ra doi.
type MucLichBoss struct {
	DuLieu []PhaBoss
}

// mucDangChay la mot MucLichBoss kem trang thai chay cua no. Tach khoi
// MucLichBoss (kieu CONG KHAI, do Task 15 khai bao du lieu) de nguoi viet du
// lieu boss khong bao gio phai — va khong bao gio co the — dat tay cac truong
// dieu phoi ben trong.
//
// CHI goroutine vongLap (hoac test goi truc tiep nhip(), khi chua BatDau) duoc
// doc/ghi cac truong duoi day. Ngoai le DUY NHAT la the: no la atomic vi ca
// actor cung cham.
type mucDangChay struct {
	MucLichBoss

	the *TheSongBoss

	// lanSinhCuoi la moc cua lan gui lenh sinh THANH CONG gan nhat (hoac luc
	// mo bo dieu phoi, neu chua sinh lan nao) — goc de dem GiayNghi. Gui HONG
	// KHONG dong vao truong nay: neu dong, mot khu vuc ban trong choc lat se
	// day con boss lui ca mot chu ky.
	lanSinhCuoi time.Time

	// hoanToi: truoc moc nay thi khong thu gui lai (xem khoangHoanKhiGuiHong).
	// Zero value = khong hoan gi ca.
	hoanToi time.Time

	// ngungHan bat khi khu vuc dich DA DUNG han. Khac hoanToi (tam thoi): mot
	// actor da dung thi khong bao gio chay lai (Dung() dong kenh xong bang
	// sync.Once), nen thu lai la vo ich — va neu khong ngung, luc TAT MAY CHU
	// ca 5 khu vuc cung dung se lam bo dieu phoi ban ra hang chuc dong log
	// moi giay cho toi khi chinh no bi dung.
	ngungHan bool
}

// BoLichBoss la GOROUTINE NEN quyet dinh khi nao + o dau sinh boss.
//
// Dung khuon BoLuuDinhKy (luudinhky.go): Mo... tao nhung CHUA chay, BatDau()
// rieng, dong ho tiem duoc, Dung() sach. Xem chu thich cac truong duoi de biet
// vi sao tung mieng phai nhu vay.
type BoLichBoss struct {
	soKhuVuc *SoKhuVuc
	muc      []*mucDangChay

	// dongHo la DONG HO TIEM DUOC (khuon k.dongHo o khuvuc.go,
	// BoLuuDinhKy.dongHo o luudinhky.go, BoPhatThanh.dongHo o thongbao.go):
	// khong co no thi moi test lich phai Sleep THAT du GiayNghi — cham va chap
	// chon. MoBoLichBoss gan mac dinh time.Now. Test cung package gan thang
	// truong nay, CHI TRUOC khi goi BatDau() (sau do vongLap doc no moi nhip,
	// ghi tu ben ngoai luc do la data race).
	dongHo func() time.Time

	// chon la BO CHON NGAU NHIEN tiem duoc: tra mot so trong [0, n). Cung ly do
	// voi dongHo — test phai khang dinh duoc "sinh o DUNG map nay, DUNG mau
	// nay" ma khong phu thuoc may man. Mac dinh rand.IntN.
	chon func(n int) int

	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once

	// daBatDau chan BatDau() goi lan thu hai. Ngay 28/08/2026 da tra gia dung
	// loi nay o BoPhatThanh: hai goroutine cung `defer close(b.xong)` -> cai
	// thu hai panic "close of closed channel" TREN GOROUTINE NEN luc tat may
	// chu, khong ai bat duoc, chet ca tien trinh. CompareAndSwap (khong phai
	// Store) moi lam duoc viec nay: nguoi DAU TIEN doi duoc false->true moi la
	// nguoi chay that.
	//
	// No con mot vai tro thu hai: cho Dung() biet co goroutine nao de ma cho
	// hay khong — thieu no thi Dung() goi truoc BatDau() se treo vinh vien o
	// <-l.xong (khong ai dong kenh do ca).
	daBatDau atomic.Bool
}

// MoBoLichBoss tao mot bo dieu phoi CHUA CHAY. Goi BatDau() rieng de khoi dong.
//
// Tach hai buoc (khuon MoKhuVuc/MoBoLuuDinhKy/MoBoPhatThanh, KHAC MoThoGhi tu
// chay luon) de test gan duoc dongHo/chon gia TRUOC khi bat ky goroutine nao co
// the doc chung. Day khong phai su cau ky: dong ho that chi tien theo mot chieu
// va khong cho phep "tua nhanh 10 phut", nen neu goroutine da chay truoc khi
// tiem duoc dong ho thi moi test lich deu buoc phai Sleep that.
//
// lanSinhCuoi cua MOI muc khoi tao = dongHo() TAI LUC GOI (dung tinh than
// MoBoLuuDinhKy va MoKhuVuc gan LanHoiMauCuoi luc spawn): tranh zero value (nam
// 1) khien nhip DAU TIEN sau BatDau() hieu nham "da qua GiayNghi tu lau" va bung
// ca dan boss ra ngay giay khoi dong may chu, chua ai kip vao game.
//
// Kenh dung/xong tao NGAY o day (khong doi BatDau, khac BoLuuDinhKy) de Dung()
// luon an toan ke ca khi khong ai goi BatDau(): close(nil chan) la panic.
func MoBoLichBoss(soKhuVuc *SoKhuVuc, lich []MucLichBoss) *BoLichBoss {
	l := &BoLichBoss{
		soKhuVuc: soKhuVuc,
		dongHo:   time.Now,
		chon:     rand.IntN,
		dung:     make(chan struct{}),
		xong:     make(chan struct{}),
	}
	bayGio := l.dongHo()
	l.muc = make([]*mucDangChay, 0, len(lich))
	for _, m := range lich {
		l.muc = append(l.muc, &mucDangChay{
			MucLichBoss: m,
			the:         &TheSongBoss{},
			lanSinhCuoi: bayGio,
		})
	}
	return l
}

// TheSong tra the song cua muc thu i — day la HANDLE ma duong "boss da roi khoi
// the gioi" (may trang thai, Task 9) se dung de bao cho bo dieu phoi biet duoc
// phep sinh con moi. Ngoai mien -> nil (TheSongBoss nil-safe nen noi goi khong
// phai kiem).
//
// ⚠️ NO CHUA DUOC NOI DAY: hom nay chua co ai goi Tat(), nen moi muc lich sinh
// DUNG MOT LAN roi thoi. Do la lua chon CO CHU DICH — hong ve phia AN TOAN: mot
// con boss thieu thi nguoi choi thieu mot muc tieu; mot vong lap sinh boss
// khong phanh thi khu vuc ngap boss va khong co cach nao go. Task 9 noi day
// bang cach goi l.TheSong(i).Tat() khi boss het pha va roi the gioi.
func (l *BoLichBoss) TheSong(i int) *TheSongBoss {
	if i < 0 || i >= len(l.muc) {
		return nil
	}
	return l.muc[i].the
}

// BatDau khoi dong goroutine nen. Goi lan thu hai la VIEC KHONG (bo qua + log
// muc loi) — xem chu thich truong daBatDau ve cai gia da tra cho loi nay.
func (l *BoLichBoss) BatDau() {
	if !l.daBatDau.CompareAndSwap(false, true) {
		slog.Error("BatDau() bo lich boss duoc goi HAI LAN — bo qua lan nay " +
			"(hai goroutine cung mot lich se sinh boss doi va lam panic luc tat may chu)")
		return
	}
	go l.vongLap()
}

// vongLap la than bo dieu phoi: moi nhipKiemLichBoss thi nhin dong ho mot lan,
// hoac thoat sach khi co lenh dung. Cau truc select giong het KhuVuc.vongLap
// (khuvuc.go), ThoGhi.vongLap (thogi.go), BoLuuDinhKy.vongLap (luudinhky.go) va
// BoPhatThanh.vongLap (thongbao.go) — nhat quan phong cach goroutine-nen.
//
// LUOI PANIC boc THAN moi vong, KHONG boc ngoai: boc ngoai thi mot panic giet
// ca vong lap va may chu mat han he thong boss ma khong ai hay. Mot nhip hong ->
// bo nhip do, vong van quay nhip sau.
func (l *BoLichBoss) vongLap() {
	ticker := time.NewTicker(nhipKiemLichBoss)
	defer ticker.Stop()
	defer close(l.xong)

	for {
		select {
		case <-l.dung:
			return
		case <-ticker.C:
			antoan.ChayCoLuoi("lichboss.nhip", func() { l.nhip(l.dongHo()) })
		}
	}
}

// nhip la logic THUAN cua mot nhip dieu phoi — tach rieng khoi vongLap (dung ly
// do kiemVaLuu tach khoi vongLap o luudinhky.go, va hoiSinhQuai/aiQuai tach
// khoi nhipTim o khuvuc.go): test goi TRUC TIEP voi bayGio tiem san, khong phu
// thuoc ticker that.
func (l *BoLichBoss) nhip(bayGio time.Time) {
	for _, m := range l.muc {
		l.thuSinh(m, bayGio)
	}
}

// thuSinh xet MOT muc lich va gui lenh sinh neu den gio.
//
// THU TU CAC BUOC LA CO Y, doc tu tren xuong:
//
//  1. Da ngung han / boss cu con song / dang trong thoi gian hoan -> ve ngay.
//  2. Kiem DU LIEU TINH (khong dung toi bo chon ngau nhien): du lieu hong la
//     loi nguoi viet du lieu, phai la ra NGAY tu nhip dau chu khong phai doi
//     het mot chu ky GiayNghi roi moi bao.
//  3. Chua du GiayNghi -> ve.
//  4. Chon map + chon mau + dung con boss + gui lenh.
func (l *BoLichBoss) thuSinh(m *mucDangChay, bayGio time.Time) {
	if m.ngungHan || m.the.ConSong() || bayGio.Before(m.hoanToi) {
		return
	}

	pha := m.phaDau()
	if lyDo := kiemPhaDauHopLe(pha); lyDo != "" {
		m.ngungHan = true
		slog.Error("muc lich boss co du lieu hong — ngung han muc nay (thu lai bao nhieu lan cung the)",
			"lyDo", lyDo)
		return
	}

	if bayGio.Sub(m.lanSinhCuoi) < time.Duration(pha.GiayNghi)*time.Second {
		return
	}

	// Chon map TRUOC, chon mau SAU — thu tu nay la HOP DONG voi bo chon tiem
	// duoc (test kich ban dua vao no). Doi thu tu la doi y nghia cua moi test
	// tiem bo chon.
	mapId := pha.MapVao[l.chon(len(pha.MapVao))]
	k := l.soKhuVuc.Tra(mapId)
	if k == nil {
		m.ngungHan = true
		slog.Error("muc lich boss tro toi mot khu vuc KHONG TON TAI — ngung han muc nay",
			"ten", pha.Ten, "mapId", mapId)
		return
	}

	b := l.taoBoss(m, pha, k.BanDo(), bayGio)

	// KE TU DONG NAY con tro b THUOC VE ACTOR — khong duoc doc/ghi lai gi cua
	// no nua, ke ca khi gui that bai (actor co the da nhan duoc va dang xu ly).
	switch err := k.SinhBoss(b); {
	case err == nil:
		m.the.Bat()
		m.lanSinhCuoi = bayGio
		m.hoanToi = time.Time{}
		slog.Info("da gui lenh sinh boss", "ten", pha.Ten, "mapId", mapId)

	case errors.Is(err, ErrActorDaDung):
		// Khu vuc da dung thi KHONG BAO GIO chay lai (Dung() dong kenh xong
		// bang sync.Once). Thu lai la vo ich, va luc tat may chu ca 5 khu vuc
		// cung dung -> khong ngung thi bo dieu phoi ban hang chuc dong log moi
		// giay cho toi khi chinh no bi dung.
		m.ngungHan = true
		slog.Warn("khu vuc da dung — ngung han muc lich boss nay",
			"ten", pha.Ten, "mapId", mapId)

	default:
		// Loi TAM THOI (khu vuc dang ban -> ErrControlQuaHan; chua BatDau ->
		// ErrActorChuaBatDau). HOAN roi thu lai, va TUYET DOI khong dong vao
		// lanSinhCuoi: mot khu vuc ban trong choc lat khong duoc quyen day con
		// boss lui ca mot chu ky.
		m.hoanToi = bayGio.Add(khoangHoanKhiGuiHong)
		slog.Warn("gui lenh sinh boss that bai — hoan lai roi thu tiep",
			"ten", pha.Ten, "mapId", mapId, "loi", err, "hoanToi", m.hoanToi)
	}
}

// phaDau tra con tro toi pha DAU TIEN cua muc (nil neu muc khong co pha nao).
// Xem chu thich MucLichBoss ve ly do bo dieu phoi chi doc pha dau.
func (m *mucDangChay) phaDau() *PhaBoss {
	if len(m.DuLieu) == 0 {
		return nil
	}
	return &m.DuLieu[0]
}

// kiemPhaDauHopLe tra ve LY DO hong (chuoi rong = hop le). Tra chuoi thay vi
// bool de dong log noi duoc DUNG cai gi sai — mot muc lich im lang khong sinh
// boss ma khong noi vi sao la loai loi ton nhieu gio nhat de tim.
//
// Bon dieu kien deu la "sai thi khong the sinh duoc, khong phai sai tam thoi":
//   - khong co pha nao         -> khong biet sinh con gi
//   - Cap ngoai 1..255         -> client hien "cap 0" (xem PhaBoss.Cap)
//   - khong co map nao         -> khong biet sinh o dau
//   - khong co mau nao         -> boss ra doi voi 0 mau, chet ngay lap tuc
func kiemPhaDauHopLe(pha *PhaBoss) string {
	switch {
	case pha == nil:
		return "muc lich khong co pha nao (DuLieu rong)"
	case !pha.CapHopLe():
		return "pha dau co Cap = 0 (hop le la 1..255)"
	case len(pha.MapVao) == 0:
		return "pha dau khong khai bao MapVao (khong biet sinh o khu vuc nao)"
	case len(pha.HpChon) == 0:
		return "pha dau khong khai bao HpChon (boss se ra doi voi 0 mau)"
	default:
		return ""
	}
}

// taoBoss dung MOT con boss hoan chinh, san sang ban giao cho actor.
//
// TRANG THAI BAN DAU LA BossHien (khong phai BossNghi): CHAN-7 chia doi trach
// nhiem — "khi nao + o dau" da duoc tra loi XONG o day, nen lenh gui di co
// nghia la "hien ra ngay". Neu ban giao o trang thai NGHI thi actor lai phai
// tu dem nguoc dong ho, dung cai ma CHAN-7 cam. themBoss thay TrangThai.
// DangHien() = true nen phat -5 ngay cho ca khu vuc.
//
// DuLieu duoc SAO CHEP (khong dung chung lat cat voi muc lich): bo dieu phoi
// con doc muc lich moi nhip, con actor se doc — va qua PhaHienTaiData() la co
// the ghi — vao du lieu pha cua con boss. Dung chung mot mang nen se la data
// race that su ngay khi Task 9 sua mot truong nao do cua pha. Vai chuc byte
// moi lan sinh boss la cai gia khong dang ke de xoa han mot lop rui ro.
//
// Vi tri = diem sinh mac dinh cua ban do (BanDo.SpawnX/SpawnY) — dung cho boss
// dung tren nen dat. KHONG dung (0,0): cot 0 cua luoi la tuong bien, con boss
// se nam trong da.
func (l *BoLichBoss) taoBoss(m *mucDangChay, pha *PhaBoss, bd *BanDo, bayGio time.Time) *Boss {
	// KepHp: chan ca gia tri am lan gia tri vuot tran int32 trong du lieu
	// (§6.17 + bay #4 thiet ke §7) — du lieu boss do nguoi go tay, khong the
	// tin la luon nam trong mien.
	hp := KepHp(pha.HpChon[l.chon(len(pha.HpChon))])

	duLieu := make([]PhaBoss, len(m.DuLieu))
	copy(duLieu, m.DuLieu)

	return &Boss{
		Id:              CapIdBoss(),
		PhaHienTai:      0,
		DuLieu:          duLieu,
		Hp:              hp,
		HpToiDa:         hp,
		ViTriX:          bd.SpawnX,
		ViTriY:          bd.SpawnY,
		TrangThai:       BossHien,
		MocThoiGian:     bayGio,
		LucHienPha:      bayGio,
		NhipDanhHienTai: pha.NhipDanhMs,
	}
}

// Dung dung goroutine nen mot cach sach: dong kenh dung (MOT LAN, nho sync.Once
// nen goi bao nhieu lan cung khong panic) roi cho vongLap thoat han truoc khi
// tra ve — khong ro ri goroutine, dung khuon KhuVuc.Dung/BoLuuDinhKy.Dung/
// BoPhatThanh.Dung.
//
// Goi khi CHUA BatDau() cung an toan (tra ve ngay, khong treo) nho co daBatDau:
// khong co goroutine nao thi khong co gi de cho.
func (l *BoLichBoss) Dung() {
	l.dungMotLan.Do(func() {
		close(l.dung)
	})
	if l.daBatDau.Load() {
		<-l.xong
	}
}
