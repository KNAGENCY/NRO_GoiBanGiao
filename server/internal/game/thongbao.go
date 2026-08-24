// Package game — duong PHAT THANH toan may chu (goi 93).
//
// CHAN-2 ("actor KHONG duoc goi ra ngoai"): mot thong bao toan server phai toi
// tay MOI nguoi dang online, tuc la cham vao rat nhieu phien mot luc. Neu actor
// KhuVuc (khuvuc.go, nhip 100ms) tu minh lam viec do — duyet danh sach phien,
// dung goi, day xuong socket — thi mot client ngung doc la ca khu vuc dung tim
// (dung bay §6.10b BAT-DAU-TU-DAY.md: actor tung bi chan toi 30 giay). Vi vay
// dat MOT HANG DOI nam NGOAI actor, y het khuon ThoGhi (thogi.go):
//
//	Actor khu vuc              Xep() — CHI cham kenh      vongLap (goroutine rieng)
//	─────────────              ───────────────────────    ─────────────────────────
//	quyet dinh "can bao"  ──►  hangDoi <- yeuCau     ──►  rut hang doi, goi phat()
//	di tiep NGAY               (select+default, KHONG     (phat() duoc phep cham
//	KHONG BAO GIO CHO           BAO GIO CHAN — day thi     mang/o dia — khong actor
//	                            BO va dem soLanBo)         nao dung o day)
//
// Khac ThoGhi mot diem quan trong: ThoGhi GOP theo id nhan vat (map, ban moi
// ghi de ban cu) vi "ban moi luon tot hon"; thong bao thi NGUOC LAI — moi cai
// la mot su kien RIENG, khong cai nao thay the cai nao — nen dung KENH CO TRAN
// (FIFO co gioi han) chu khong phai map.
package game

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"nroserver/internal/antoan"
	"nroserver/internal/protocol"
)

// lenhThongBaoToanServer la ma lenh goi thong bao chay ngang toan may chu.
//
// NGUON (bang chung, khong phai suy doan): server Java goc
// 1-MAY-CHU-GAME/src/server/ServerNotify.java:74-76 —
//
//	msg = new Message(93);
//	msg.writer().writeUTF(text);
//	Service.gI().sendMessAllPlayer(msg);
//
// Do la NOI DUY NHAT trong ca 489 file .java dung `new Message(93)` (da grep
// toan bo 1-MAY-CHU-GAME/src/; cac lan xuat hien "93" con lai deu la id option
// vat pham hoac mapId, khong lien quan ma lenh). Than goi vi vay CHI CO MOT
// truong: mot chuoi writeUTF — khong co byte loai, khong co mau, khong co dem.
const lenhThongBaoToanServer int8 = 93

const (
	// tranHangDoiThongBao la SUC CHUA hang doi. Co tran la dieu BAT BUOC: hang
	// doi khong tran = Xep() co the phai cho khi day = actor bi chan = dung
	// nguyen cai loi CHAN-2 dang muon chua. 64 la con so "du rong cho mot con
	// bao thong bao that (boss ra, su kien, admin ho), du hep de neu duong phat
	// thanh tac that thi ta BIET (soLanBo tang) thay vi am tham phinh bo nho".
	tranHangDoiThongBao = 64

	// nhipToiThieuThongBao la KHOANG CACH TOI THIEU giua hai thong bao lien
	// tiep — lop CHONG NGAP: khong co no, mot vong lap loi (hoac mot admin go
	// nham) co the ban hang tram dong chu dap kin man hinh moi nguoi choi.
	//
	// NGUON 1,5 giay: ServerNotify.java:65 — vong `run()` cua server Java goc
	// rut MOT thong bao khoi hang doi (dong 52 `sendChatVip(notifies.remove(0))`)
	// roi `Functions.sleep(1500)`. Tuc la ban Java that su chay dung nhip nay
	// suot nhieu nam — day la con so DA KIEM CHUNG TRONG SAN XUAT, khong phai
	// so ta tu nghi ra.
	nhipToiThieuThongBao = 1500 * time.Millisecond

	// nhipKiemPhatThanh la tan suat DRIVER THAT (ticker cua vongLap) — KHONG
	// PHAI khoang cach giua hai thong bao (dung nhipToiThieuThongBao cho viec
	// do). Chi la "bao lau nhin dong ho + ngo hang doi mot lan", dung y het
	// nhipKiemLuuDinhKy (luudinhky.go) va khoangCachThuLai (thogi.go): mot nhip
	// khi hang doi rong chi la mot select-default, chi phi khong dang ke.
	nhipKiemPhatThanh = 50 * time.Millisecond
)

// YeuCauThongBao la MOT yeu cau phat thanh dang xep hang.
//
// MocTao ghi lai luc yeu cau duoc SINH RA (khong phai luc duoc phat) — dung de
// do do tre hang doi trong log: thong bao nam cho bao lau moi len song. Day la
// so lieu chan doan duy nhat cho biet duong phat thanh dang thong hay dang nghen.
type YeuCauThongBao struct {
	NoiDung string
	MocTao  time.Time
}

// BoPhatThanh la hang doi thong bao + goroutine phat rieng (CHAN-2).
//
// Bat bien phan chia trach nhiem, doc ky truoc khi sua:
//   - Xep() chay tren goroutine BAT KY (ke ca actor KhuVuc) — CHI duoc lam viec
//     khong chan va khong cham I/O: mot select+default vao kenh, mot atomic.Add.
//   - nhip()/vongLap() chay tren goroutine RIENG cua chinh BoPhatThanh — day la
//     noi DUY NHAT duoc doc/ghi lanPhatCuoi va goi phat().
type BoPhatThanh struct {
	// hangDoi la kenh co tran (FIFO). Kenh — khong phai slice + mutex — vi day
	// dung la thu Go lam san: gui khong chan bang select+default, rut trong
	// select cung voi dung/ticker ma khong can dieu kien bien (sync.Cond).
	hangDoi chan YeuCauThongBao

	// soLanBo dem TONG so yeu cau bi BO vi hang doi day. atomic vi nhieu
	// goroutine actor cung Xep() song song. KHONG duoc nuot im lang (bay §6
	// ban Java: catch rong khien su co khong the chan doan) — nhung cung KHONG
	// duoc log ngay trong Xep(): slog ghi ra stderr la I/O, ma Xep() chay tren
	// actor. Vi vay dem o day, con VIEC LOG do goroutine phat thanh lam (xem
	// baoCaoSoBo) — dung tinh than "actor khong cham I/O".
	soLanBo atomic.Int64

	// soBoDaBaoCao la moc soLanBo cua lan log gan nhat — de chi log KHI CON SO
	// TANG THEM, khong lap lai mai mot con so cu moi lan phat. CHI goroutine
	// phat thanh cham (qua baoCaoSoBo), nen khong can atomic.
	soBoDaBaoCao int64

	// dongHo la DONG HO TIEM DUOC (dung khuon k.dongHo o khuvuc.go va
	// BoLuuDinhKy.dongHo o luudinhky.go): test kiem duoc nhip 1,5 giay MA KHONG
	// phai Sleep that 1,5 giay. MoBoPhatThanh gan mac dinh time.Now. Test cung
	// package gan thang truong nay — CHI TRUOC khi goi BatDau() (sau do vongLap
	// doc no moi nhip, ghi tu ben ngoai luc do la data race).
	dongHo func() time.Time

	// lanPhatCuoi la moc thoi gian cua thong bao PHAT GAN NHAT. CHI vongLap
	// (hoac test goi truc tiep nhip()) duoc cham — mot luong duy nhat, dung
	// tinh than lanLuuCuoi cua BoLuuDinhKy.
	//
	// Khoi tao = ZERO VALUE co chu dich, KHAC HAN BoLuuDinhKy (gan dongHo()
	// luc tao de tranh mot lan luu THUA ngay khi khoi dong). O day ta MUON dieu
	// nguoc lai: thong bao DAU TIEN phai len song NGAY, khong bat nguoi choi
	// cho 1,5 giay vo co — dung nhu ban Java (ServerNotify.run() gui ngay roi
	// MOI sleep(1500)). Khong co "lan phat thua" nao de lo: hang doi rong thi
	// nhip() khong lam gi ca.
	lanPhatCuoi time.Time

	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once

	// daBatDau cho Dung() biet co goroutine nao de ma cho hay khong. Thieu co
	// nay thi Dung() goi khi CHUA BatDau() se treo vinh vien o <-b.xong (khong
	// ai dong kenh do ca).
	daBatDau atomic.Bool
}

// MoBoPhatThanh tao mot BoPhatThanh CHUA CHAY — goi BatDau() rieng de khoi
// dong goroutine phat. Tach hai buoc (khuon MoKhuVuc/MoBoLuuDinhKy, KHAC
// MoThoGhi tu chay luon) de test gan duoc dongHo gia TRUOC khi bat ky goroutine
// nao co the doc no.
//
// Kenh dung/xong tao NGAY o day (khong doi BatDau nhu BoLuuDinhKy) de Dung()
// luon an toan, ke ca khi khong ai goi BatDau(): close(nil chan) la panic.
func MoBoPhatThanh() *BoPhatThanh {
	return &BoPhatThanh{
		hangDoi: make(chan YeuCauThongBao, tranHangDoiThongBao),
		dongHo:  time.Now,
		dung:    make(chan struct{}),
		xong:    make(chan struct{}),
	}
}

// Xep dua mot yeu cau thong bao vao hang doi — KHONG BAO GIO CHAN, du hang doi
// day, du goroutine phat dang ban hay da dung han. Day la CA LY DO TON TAI cua
// file nay (CHAN-2): ham nay AN TOAN goi tu trong actor KhuVuc.
//
// Hang doi day -> BO yeu cau va dem soLanBo. Bo la lua chon CO CHU DICH, khong
// phai thieu sot: mot thong bao mat di lam nguoi choi lo mot dong chu; mot
// actor bi chan lam ca khu vuc dung tim. Doi voi du lieu KHONG duoc mat (tien,
// chi so nhan vat) da co ThoGhi lo — hai loai du lieu, hai chinh sach.
func (b *BoPhatThanh) Xep(yc YeuCauThongBao) {
	select {
	case b.hangDoi <- yc:
	default:
		// KHONG log o day: Xep() chay tren actor, slog cham stderr la I/O.
		// Goroutine phat se bao cao con so nay ho (xem baoCaoSoBo).
		b.soLanBo.Add(1)
	}
}

// SoLanBo tra ve tong so yeu cau da bi bo vi hang doi day — so lieu van hanh
// (va la moc kiem cua test). Doc duoc tu bat ky goroutine nao.
func (b *BoPhatThanh) SoLanBo() int64 { return b.soLanBo.Load() }

// BatDau khoi dong goroutine phat thanh rieng. phat la ham THAT SU dua noi
// dung toi nguoi choi — vd dung goiThongBaoToanServer(noiDung) roi phien.XepGui
// cho tung nguoi dang online.
//
// phat() chay tren goroutine CUA BO PHAT THANH, KHONG PHAI tren actor nao —
// nen no DUOC PHEP cham mang/o dia/log thoai mai. Dung goi BatDau hai lan (se
// co hai goroutine cung rut mot hang doi, va Dung() chi cho mot cai thoat).
func (b *BoPhatThanh) BatDau(phat func(string)) {
	b.daBatDau.Store(true)
	go b.vongLap(phat)
}

// vongLap la than bo phat thanh: moi nhipKiemPhatThanh thi ngo hang doi mot
// lan, hoac thoat sach khi co lenh dung. Cau truc select giong het
// KhuVuc.vongLap (khuvuc.go), ThoGhi.vongLap (thogi.go) va BoLuuDinhKy.vongLap
// (luudinhky.go) — nhat quan phong cach goroutine-nen toan du an.
//
// LUOI PANIC boc THAN moi vong, KHONG boc ngoai: boc ngoai thi mot panic giet
// ca vong lap, may chu mat han tinh nang thong bao ma khong ai hay. Mot nhip
// hong -> bo nhip do, vong van quay nhip sau.
func (b *BoPhatThanh) vongLap(phat func(string)) {
	ticker := time.NewTicker(nhipKiemPhatThanh)
	defer ticker.Stop()
	defer close(b.xong)

	for {
		select {
		case <-b.dung:
			return
		case <-ticker.C:
			antoan.ChayCoLuoi("thongbao.nhip", func() { b.nhip(b.dongHo(), phat) })
		}
	}
}

// nhip la logic THUAN cua mot nhip phat thanh — tach rieng khoi vongLap (dung
// ly do KH6 Task 9 tach hoiSinhQuai/hoiMauQuai khoi nhipTim, va kiemVaLuu khoi
// vongLap o luudinhky.go): test goi TRUC TIEP voi bayGio tiem san, khong phu
// thuoc ticker that.
//
// Ba nhanh:
//  1. Hang doi rong -> khong lam gi.
//  2. CHUA du nhipToiThieuThongBao ke tu lan phat truoc -> XEP LAI yeu cau vao
//     hang doi cho toi luot. Yeu cau KHONG BI MAT — day la khac biet cot loi so
//     voi "chong ngap kieu vut bo": nguoi choi van doc duoc thong bao, chi la
//     doc cham hon nhip toi thieu. Ngoai le DUY NHAT: hang doi trong luc do da
//     bi lap day boi Xep() khac -> luc do moi dem soLanBo (van khong duoc chan).
//  3. Du nhip -> ghi lanPhatCuoi roi goi phat().
//
// lanPhatCuoi dat TRUOC khi goi phat() (dung tinh than lanLuuCuoi cua
// BoLuuDinhKy va q.LanDanhCuoi cua khuvuc.go): moc phan anh luc BAT DAU phat,
// va chac chan duoc cap nhat du phat() co cham bao lau di nua — neu dat sau,
// mot phat() cham 3 giay se khien thong bao ke tiep di ngay lap tuc, thung
// dung lop chong ngap.
func (b *BoPhatThanh) nhip(bayGio time.Time, phat func(string)) {
	var yc YeuCauThongBao
	select {
	case yc = <-b.hangDoi:
	default:
		return // hang doi rong — khong co gi de phat
	}

	if !b.lanPhatCuoi.IsZero() && bayGio.Sub(b.lanPhatCuoi) < nhipToiThieuThongBao {
		select {
		case b.hangDoi <- yc:
		default:
			b.soLanBo.Add(1)
		}
		return
	}

	b.lanPhatCuoi = bayGio
	b.baoCaoSoBo()
	slog.Info("phat thanh toan may chu",
		"noiDung", yc.NoiDung, "doTreHangDoi", bayGio.Sub(yc.MocTao))
	phat(yc.NoiDung)
}

// baoCaoSoBo log so yeu cau bi bo KHI con so tang them ke tu lan bao truoc.
// Goi tu goroutine phat thanh (KHONG phai tu Xep/actor — xem chu thich soLanBo).
func (b *BoPhatThanh) baoCaoSoBo() {
	hienTai := b.soLanBo.Load()
	if hienTai <= b.soBoDaBaoCao {
		return
	}
	slog.Warn("hang doi thong bao day — da BO bot yeu cau phat thanh",
		"tongSoBo", hienTai, "boThemKeTuLanBao", hienTai-b.soBoDaBaoCao,
		"tranHangDoi", tranHangDoiThongBao)
	b.soBoDaBaoCao = hienTai
}

// Dung dung goroutine phat thanh mot cach sach: dong kenh dung (MOT LAN, nho
// sync.Once nen goi bao nhieu lan cung khong panic) roi cho vongLap thoat han
// truoc khi tra ve — dam bao khong ro ri goroutine, dung khuon KhuVuc.Dung
// (khuvuc.go) va BoLuuDinhKy.Dung (luudinhky.go).
//
// Goi khi CHUA BatDau() cung an toan (tra ve ngay, khong treo) — nho co
// daBatDau: khong co goroutine nao thi khong co gi de cho.
func (b *BoPhatThanh) Dung() {
	b.dungMotLan.Do(func() {
		close(b.dung)
	})
	if b.daBatDau.Load() {
		<-b.xong
	}
}

// goiThongBaoToanServer dung goi XUONG 93: mot dong thong bao chay ngang cho
// TOAN BO nguoi choi dang online.
//
// KHUON BYTE — nguon: 1-MAY-CHU-GAME/src/server/ServerNotify.java:74-76
// (`new Message(93)` + `msg.writer().writeUTF(text)` + `sendMessAllPlayer`).
// Than goi CHI CO MOT truong writeUTF, KHONG co gi khac: khong byte loai,
// khong mau chu, khong so luong. Ghi them bat ky byte nao la client doc lech.
// Da grep toan bo 489 file .java — day la lan xuat hien DUY NHAT cua ma lenh
// 93, nen khong co dang goi 93 thu hai nao de phai chieu theo.
//
// Tra ve nil neu noiDung qua dai: writeUTF ghi do dai bang 2 byte khong dau nen
// toi da 65535 byte SAU KHI ma hoa modified UTF-8 (xem Writer.UTF,
// protocol/datastream.go). Chu ky ham khong tra loi duoc error, ma mot goi
// THIEU BYTE con te hon khong goi (client doc lech ca luong tin phia sau) —
// nen chon nil, va nguoi goi PHAI kiem nil truoc khi gui.
func goiThongBaoToanServer(noiDung string) *protocol.Message {
	w := protocol.NewWriter()
	if err := w.UTF(noiDung); err != nil {
		slog.Error("thong bao qua dai, khong dung duoc goi 93 — bo qua",
			"loi", err, "soKyTu", len(noiDung))
		return nil
	}
	return protocol.NewMessage(lenhThongBaoToanServer, w)
}
