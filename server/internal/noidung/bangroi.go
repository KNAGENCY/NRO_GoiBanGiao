// bangroi.go — BANG ROI DO: quai/boss chet thi rot ra nhung mon gi, voi ti le nao.
//
// (KH8 task 9. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// hang "9" bang Task, §2.3 tham chieu cheo bang MA CHUOI, §4 anh chup mot lan
// dau nhip.)
//
// # CAI GI DOI O TASK NAY
//
// Truoc task nay, bang roi la BON BIEN GOI trong game/roido.go
// (roiDoTheoQuai / roiDoQuaiMacDinh / roiDoTheoBoss / roiDoTheoBossLon) — con so
// can bang nam trong ma nguon Go, va bang quai duoc tra bang `Quai.MauId`, ma
// tu task 8 truong do MANG ANH_ID (id anh trong catalog) chu khong con la can
// cuoc cua mau quai nua. Hai mau quai khac han nhau dung chung mot anh la dung
// chung mot bang roi — mot loi im lang.
//
// Tu task nay: MOT ban duy nhat, kho-noi-dung/bang-roi.yaml, tra cuu bang `ma`
// CHUOI cua mau quai (dung khoa voi khoi `sinh:` cua quai.yaml). Bang Go da bi
// XOA — dot bien "actor doc bang Go cu" khong con viet duoc nua.
//
// # BA BO KIEM, VI SAO TUNG CAI
//
//	ma quai co that      — `ma` go nham la mot con quai IM LANG roi bang mac dinh
//	                       thay vi bang cua no; khong mot dong log, chi la "sao
//	                       con nay it rot do the".
//	mau_id ton tai + DA KHOA — mau_id tro vao khoang khong la mot mon do khong
//	                       ten trong tui nguoi choi (client NullReference roi
//	                       NUOT). "Da khoa" (khoa-id.yaml) chat hon "ton tai":
//	                       no chan ca truong hop bang roi tro toi mot mau vua
//	                       duoc them ma chua ghi so.
//	so_luong <= tran_chong — client chong do theo TranChong cua mau; rot 50 cai
//	                       vao mot o chong toi da 10 la mat 40 cai, im lang.
//
// # MA BOSS: KIEM O `game`, KHONG KIEM O DAY (cung khuon KiemIdAnhQuai)
//
// Danh sach ma boss hop le la cac hang MaBoss* trong goi `game`
// (noidungboss.go), va `noidung` KHONG duoc import `game` (huong import da chot:
// game -> noidung, khong bao gio nguoc lai). Vi vay bo kiem duoc CHIA y het
// khuon `anh_id` cua task 8:
//
//	noidung/bangroi.go kiemBangRoi        — ma khong rong, khong trung, mau_id...
//	noidung/bangroi.go KiemMaBangRoiBoss  — ham THUAN, nhan danh sach ma lam THAM SO
//	game/roido.go      (nguoi goi)        — goi ham tren voi danh sach THAT
//
// # VI SAO CO MOT MUC `ma: boss-chung`
//
// Ngay hom nay `game.Boss` KHONG mang `ma` cua con boss: no chi mang []PhaBoss
// (lichboss.MucLichTuNoiDungBoss bo truong Ma di khi dung lich). Nghia la o cho
// boss chet, actor KHONG BIET no dang giet con boss nao — no chi biet PHA hien
// tai. Cho tra `ma` do vao Boss nam trong lichboss.go/noidungboss.go, hai file
// mot task khac dang so huu; sua chung o day la cam chac dut build luc gop.
//
// Vi vay bang boss hom nay co DUNG MOT muc, ma `boss-chung`, dung cho MOI con
// boss — dung nguyen hanh vi cu (game/roido.go truoc task nay cung dung MOT bang
// cho moi boss). Muc nao co `chi_boss_lon: true` chi rot khi pha dang chet la
// pha LaBossLon — dung vai cua bang roiDoTheoBossLon cu.
//
// KHI NAO BO `boss-chung`: khi `game.Boss` mang duoc `Ma`, them muc theo ma that
// (`nhan-gia-mat-na`, `hoc-vien-dao-tau`) vao YAML va doi cho goi trong
// khuvuc.go tu hang MaBangRoiBossChung sang `b.Ma`. Bo kiem KiemMaBangRoiBoss da
// don san viec do: no CHAP NHAN ca ma that lan ma chung.
package noidung

import "fmt"

// MucRoi la MOT dong cua bang roi: "mon mau_id, so_luong cai, ti le bao nhieu %".
// Moi dong gieo DOC LAP — mot lan quai chet co the rot 0, 1 hay tat ca cac dong.
//
// Ti le de o don vi PHAN TRAM (nguoi doc file hieu ngay) chu khong phai phan van
// nhu ben trong may chu: viec doi don vi la viec cua game/roido.go, khong phai
// viec cua nguoi go so.
type MucRoi struct {
	MauId   int `yaml:"mau_id"`
	SoLuong int `yaml:"so_luong"`

	// TiLePhanTram: 0..100. 0 = KHONG BAO GIO rot (khai mot dong ti le 0 la cach
	// tat mot mon ma van giu lai dong so + chu thich vi sao no tung co).
	TiLePhanTram float64 `yaml:"ti_le_phan_tram"`

	// ChacChan = 100%, khai TUONG MINH thay vi go `ti_le_phan_tram: 100`. Hai ly
	// do: doc "chac chan" ro hon doc mot con so, va no chan mot loi that — go
	// 100 vao mot truong so thuc roi mai sau doi cach lam tron la mot mon "gan
	// nhu chac chan" ma khong ai phat hien.
	//
	// Khai CA HAI (chac_chan + ti_le_phan_tram khac 0) la LOI: hai nguon cho mot
	// su that.
	ChacChan bool `yaml:"chac_chan"`

	// ChiBossLon: dong nay chi tinh khi pha dang chet la pha LaBossLon. Vo nghia
	// trong khoi `quai:` — bo kiem TU CHOI neu khai o do.
	ChiBossLon bool `yaml:"chi_boss_lon"`
}

// BangRoiTheoMa la mot khoi `ma:` + `roi:` trong bang-roi.yaml.
type BangRoiTheoMa struct {
	Ma  string   `yaml:"ma"`
	Roi []MucRoi `yaml:"roi"`
}

type fileBangRoi struct {
	QuaiMacDinh []MucRoi        `yaml:"quai_mac_dinh"`
	Quai        []BangRoiTheoMa `yaml:"quai"`
	Boss        []BangRoiTheoMa `yaml:"boss"`
}

// BangRoi la phan bang roi cua MOT anh chup noi dung. Map duoc dung o day (khac
// voi slice o cho khac) vi tra cuu xay ra trong ACTOR, moi lan quai chet.
type BangRoi struct {
	QuaiMacDinh []MucRoi
	Quai        map[string][]MucRoi
	Boss        map[string][]MucRoi
}

// MaBangRoiBossChung — xem khoi chu thich dau tep, muc "VI SAO CO MOT MUC".
const MaBangRoiBossChung = "boss-chung"

// TranTiLePhanTram / TranSoLuongRoi — hai tran hien nhien, de thanh hang de thong
// bao loi tro toi mot cho.
const TranTiLePhanTram = 100

// ---------------------------------------------------------------------------
// Cua tra cuu tren anh chup — BAN SAO SAU (ky luat kho.go)
// ---------------------------------------------------------------------------

// BangRoiCuaQuai tra ban sao bang roi cua mot mau quai theo MA.
//
// Ma khong co bang rieng -> BANG MAC DINH (`quai_mac_dinh`), KHONG phai nil. Co
// mac dinh (thay vi "khong co bang thi khong rot gi") vi them mot mau quai moi
// ma quen bang roi la loi im lang kho thay nhat trong ca he thong nay — giu
// nguyen ly do da ghi o game/roido.go truoc task 9.
func (a *anhChup) BangRoiCuaQuai(ma string) []MucRoi {
	if ds, co := a.duLieu.BangRoi.Quai[ma]; co {
		return bangRoiSao(ds)
	}
	return bangRoiSao(a.duLieu.BangRoi.QuaiMacDinh)
}

// BangRoiCuaBoss tra ban sao bang roi cua mot con boss theo MA. Khong co ->
// nil (boss khong khai bang thi khong rot gi): khac han quai o cho boss la noi
// dung dat tay tung con, khong co "dan boss vo danh" nao de lo.
func (a *anhChup) BangRoiCuaBoss(ma string) []MucRoi {
	ds, co := a.duLieu.BangRoi.Boss[ma]
	if !co {
		return nil
	}
	return bangRoiSao(ds)
}

// bangRoiSao sao mot dong bang roi. MucRoi toan truong vo huong nen mot lan
// copy slice la du sau.
func bangRoiSao(ds []MucRoi) []MucRoi {
	if len(ds) == 0 {
		return nil
	}
	ra := make([]MucRoi, len(ds))
	copy(ra, ds)
	return ra
}

// ---------------------------------------------------------------------------
// Dung BangRoi tu file da phan tich
// ---------------------------------------------------------------------------

func dungBangRoi(f *fileBangRoi) BangRoi {
	br := BangRoi{
		QuaiMacDinh: f.QuaiMacDinh,
		Quai:        make(map[string][]MucRoi, len(f.Quai)),
		Boss:        make(map[string][]MucRoi, len(f.Boss)),
	}
	for _, k := range f.Quai {
		br.Quai[k.Ma] = k.Roi
	}
	for _, k := range f.Boss {
		br.Boss[k.Ma] = k.Roi
	}
	return br
}

// ---------------------------------------------------------------------------
// Bo kiem
// ---------------------------------------------------------------------------

// kiemBangRoi kiem toan bo bang-roi.yaml.
//
//	f        — file da phan tich (doc THANG tu file, khong qua map: chi co o day
//	           moi phat hien duoc "khai hai lan mot ma").
//	mauQuai  — de doi chieu `ma` quai; PHAI la cua CUNG mot DuLieu (§2.3).
//	mauVP    — de doi chieu `mau_id` + tran chong.
//	khoa     — so khoa id (khoa-id.yaml) de doi chieu "mau_id DA KHOA".
func kiemBangRoi(f *fileBangRoi, mauQuai []MauQuai, mauVP []MauVatPham, khoa []khoaIdMuc) error {
	maQuaiCo := make(map[string]bool, len(mauQuai))
	for _, m := range mauQuai {
		maQuaiCo[m.Ma] = true
	}
	daKhoa := make(map[int]bool, len(khoa))
	for _, k := range khoa {
		daKhoa[k.Id] = true
	}

	if err := kiemMotBangRoi("quai_mac_dinh", f.QuaiMacDinh, mauVP, daKhoa, false); err != nil {
		return err
	}

	daThayQuai := map[string]bool{}
	for i, k := range f.Quai {
		if k.Ma == "" {
			return fmt.Errorf("khoi quai thu %d: `ma` rong — `ma` la can cuoc, bang roi tra cuu bang chinh chuoi nay", i)
		}
		if daThayQuai[k.Ma] {
			return fmt.Errorf("bang roi quai %q khai hai lan — khoi thu hai se de len khoi thu nhat im lang", k.Ma)
		}
		daThayQuai[k.Ma] = true
		if !maQuaiCo[k.Ma] {
			return fmt.Errorf("bang roi tro toi ma quai %q khong co trong `mau_quai` cua %s — "+
				"hoac go nham, hoac vua co nguoi doi `ma` cua mot mau ma quen sua o day; "+
				"de nguyen thi con quai do IM LANG rot bang mac dinh", k.Ma, TenFileQuai)
		}
		if err := kiemMotBangRoi("quai "+k.Ma, k.Roi, mauVP, daKhoa, false); err != nil {
			return err
		}
	}

	daThayBoss := map[string]bool{}
	for i, k := range f.Boss {
		if k.Ma == "" {
			return fmt.Errorf("khoi boss thu %d: `ma` rong", i)
		}
		if daThayBoss[k.Ma] {
			return fmt.Errorf("bang roi boss %q khai hai lan", k.Ma)
		}
		daThayBoss[k.Ma] = true
		// `ma` boss KHONG doi chieu duoc o day — xem khoi chu thich dau tep, muc
		// "MA BOSS: KIEM O `game`". KiemMaBangRoiBoss lam viec do.
		if err := kiemMotBangRoi("boss "+k.Ma, k.Roi, mauVP, daKhoa, true); err != nil {
			return err
		}
	}
	return nil
}

// kiemMotBangRoi kiem tung dong cua MOT bang. `choPhepChiBossLon` chi bat o khoi
// `boss:` — `chi_boss_lon` trong khoi quai la mot dong KHONG BAO GIO chay, tuc mot
// mon nguoi go tuong la co ma thuc te khong bao gio rot.
func kiemMotBangRoi(cho string, ds []MucRoi, mauVP []MauVatPham, daKhoa map[int]bool, choPhepChiBossLon bool) error {
	for i, m := range ds {
		if m.SoLuong <= 0 {
			return fmt.Errorf("%s dong %d: so_luong = %d, phai > 0 — mot dong so_luong 0 la mot dong "+
				"khong bao gio rot gi, viet ti_le_phan_tram: 0 neu that su muon tat no", cho, i, m.SoLuong)
		}
		if m.TiLePhanTram < 0 || m.TiLePhanTram > TranTiLePhanTram {
			return fmt.Errorf("%s dong %d: ti_le_phan_tram = %g, phai trong [0, %d]",
				cho, i, m.TiLePhanTram, TranTiLePhanTram)
		}
		if m.ChacChan && m.TiLePhanTram != 0 {
			return fmt.Errorf("%s dong %d: khai CA `chac_chan: true` lan `ti_le_phan_tram: %g` — "+
				"hai nguon cho mot su that; chon mot", cho, i, m.TiLePhanTram)
		}
		if m.ChiBossLon && !choPhepChiBossLon {
			return fmt.Errorf("%s dong %d: `chi_boss_lon` chi co nghia trong khoi `boss:` — "+
				"o khoi quai no la mot dong KHONG BAO GIO rot", cho, i)
		}
		mau, co := timMauVatPhamTheoId(mauVP, m.MauId)
		if !co {
			return fmt.Errorf("%s dong %d: mau_id %d khong ton tai (co %d mau trong %s) — "+
				"mon rot ra se khong co ten, client nem NullReference roi NUOT",
				cho, i, m.MauId, len(mauVP), TenFileVatPham)
		}
		if !daKhoa[m.MauId] {
			return fmt.Errorf("%s dong %d: mau_id %d CHUA CO trong %s — bang roi chi duoc tro toi "+
				"mau da ghi so; mot mau chua khoa la mau con co the bi doi id",
				cho, i, m.MauId, TenFileKhoaId)
		}
		if int(mau.TranChong) > 0 && m.SoLuong > int(mau.TranChong) {
			return fmt.Errorf("%s dong %d: so_luong %d vuot tran_chong %d cua mau %q — phan vuot "+
				"bien mat im lang luc mon vao tui", cho, i, m.SoLuong, mau.TranChong, mau.Ten)
		}
	}
	return nil
}

func timMauVatPhamTheoId(ds []MauVatPham, id int) (MauVatPham, bool) {
	for _, m := range ds {
		if m.Id == id {
			return m, true
		}
	}
	return MauVatPham{}, false
}

// KiemMaBangRoiBoss doi chieu MOI ma trong khoi `boss:` voi mot danh sach ma hop
// le. HAM THUAN — nhan danh sach lam THAM SO, dung khuon KiemIdAnhQuai (task 8):
// danh sach that nam trong goi `game` (cac hang MaBoss*), va `noidung` khong
// duoc import `game`.
//
// Ma MaBangRoiBossChung LUON hop le, khong can co trong danh sach.
func KiemMaBangRoiBoss(du *DuLieu, maHopLe []string) error {
	if du == nil {
		return nil
	}
	hopLe := make(map[string]bool, len(maHopLe)+1)
	hopLe[MaBangRoiBossChung] = true
	for _, m := range maHopLe {
		hopLe[m] = true
	}
	for ma := range du.BangRoi.Boss {
		if !hopLe[ma] {
			return fmt.Errorf("bang roi boss %q khong tro toi con boss nao trong kho noi dung boss — "+
				"bang do KHONG BAO GIO duoc dung toi, va khong mot dong log nao noi vay", ma)
		}
	}
	return nil
}
