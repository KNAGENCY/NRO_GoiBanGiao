// kiemtra.go — BO KIEM DAY DU cua kho noi dung (KH8 task 1).
//
// nap.go giu bo kiem TOI THIEU cua task 0 (du de cau "YAML sinh ra dung byte" co
// nghia). File nay MO RONG no thanh ba bat bien K1/K2/K3 cua thiet ke §2.4, cong
// them doi chieu `khoa-id.yaml` (§2.6). Ca hai bo cung chay trong Nap; ho cung
// mot ky luat FAIL-CLOSED: mot loi la KHONG NAP GI CA.
//
// # KIEM NAO CHAN DUONG NHAN BAN / HONG NGAM NAO
//
// Doc bang nay khi ban dinh "don dep cho gon" mot file YAML.
//
//	Kiem                          | Chan duong hong nao                      | Neu KHONG co kiem thi trieu chung la gi
//	------------------------------|------------------------------------------|-------------------------------------------------
//	K1 id = vi tri (nap.go)       | Hoan vi / chen giua / doi id mau vat pham | Ao ninja cua MOI nguoi choi bien thanh thuoc
//	                              | (id dang nam trong cot vat_pham.mau_id)  | hoi mau. Khong log, khong loi.
//	K1 khoa-id.yaml (file nay)    | XOA mot mau, hay doi thu tu roi sua lai   | Nhu tren, nhung bo kiem "id = vi tri" VAN
//	                              | id cho lien tuc                          | XANH — day la lo hong duy nhat no khong thay.
//	K1 khoa day du                | Them mau moi ma QUEN ghi so              | Mau moi khong bao gio duoc doi chieu; lan sau
//	                              |                                          | co nguoi xoa no thi khong ai bat duoc.
//	K2 cua hang -> mau (nap.go)   | Ban mot mau khong ton tai                | Client nem NullReference vao mot catch RONG:
//	                              |                                          | cua hang hien do dang, khong dau vet.
//	K2 mon co option -> tran_chong| Ban mon CO OPTION nhung mau cho chong     | Hai mon khac option chong chung mot o -> MOT
//	                              | nhieu cai mot o                          | trong hai MAT SACH option khi luu CSDL.
//	K2 chieu -> lang (nap.go)     | Thieu lang / them lang thu 6              | Client cu chi gui duoc lang 0/1/2: lang la =
//	                              |                                          | nguoi choi ket o man chon lang.
//	K2 ma option 6/7              | Doi Y NGHIA cua option id 6 hoac 7        | Client bam cung 22->6, 23->7 (param x1000):
//	                              |                                          | mon do hien sai chi so, khong bao gi.
//	K3 tran client (nap.go+day)   | So vuot suc doc cua client                | Client doc lech dong byte -> catalog RONG hoac
//	                              |                                          | sai tu do tro di. Loi im lang tuyet doi.
//	K3 than goi -28/8             | Danh muc phinh qua 100000 byte            | Client goi d.mark(100000) truoc khi doc: dai
//	                              |                                          | hon la doc lech, KHONG bao loi.
//	K3 hoi_chieu_ms boi 100ms     | Con so khong khop nhip actor              | Chieu hoi som/muon nua nhip mot cach ngau
//	                              |                                          | nhien tuy luc bam — khong ai tai hien duoc.
//	K3 ten tab <= 2 dong          | Ten tab ba dong                          | Tab hien ten RONG, khong bao loi.
//	K3 ten option phai co '#'     | Ten option quen '#'                       | Nguoi choi thay dong chu tinh, khong thay con
//	                              |                                          | so — khong loi, chi vo nghia.
//
// # CON NO — GHI RO DE KHONG AI TUONG DA KIEM
//
// Thiet ke §2.5 doi kiem `GiayNghi` cua pha boss la boi so 100ms (dot bien M7/M8).
// KHONG lam duoc o task nay: kho noi dung ban 1 CHUA CO `boss.yaml` — boss van
// nam trong game/noidungboss.go va duoc kiem boi test cua chinh no. Khi task 7
// chuyen boss sang YAML thi kiem do phai them VAO DAY, canh kiemK3NhipChieu.
//
// Thiet ke §2.6 con noi toi co `ngung_dung: true` de "ngung ban ma khong xoa
// dong khoa". Co do can MOT TRUONG MOI trong kieu MauVatPham (noidung.go) — dung
// pham vi task 4. Ban 1 ngung ban bang cach BO MON KHOI CUA HANG, danh muc giu
// nguyen. Day la cach dung; khong ai duoc xoa dong nao trong danh muc.
package noidung

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// khoa-id.yaml — so cai da cap phat
// ---------------------------------------------------------------------------

// TenFileKhoaId — file thu tu cua kho noi dung. Khong phai noi dung game: no la
// SO GHI "id nao da tung duoc cap phat, cho cai gi".
const TenFileKhoaId = "khoa-id.yaml"

// khoaIdMuc la MOT dong so: mot id da cap phat va ma cua no.
type khoaIdMuc struct {
	Id int    `yaml:"id"`
	Ma string `yaml:"ma"`
}

type fileKhoaId struct {
	MauOption  []khoaIdMuc `yaml:"mau_option"`
	MauVatPham []khoaIdMuc `yaml:"mau_vat_pham"`
}

// khoaIdMaTu bien mot `ten` hien thi thanh MA on dinh: bo dau tieng Viet, ha chu
// thuong, moi cum ky tu khong phai chu/so thanh mot dau '-'.
//
// VI SAO SINH MA TU `ten` chu khong them mot truong `ma:` vao danh muc: them
// truong la sua kieu du lieu dung chung (noidung.go) — ngoai pham vi task nay,
// va la dung mot khoa canh cua thu hai cho cung mot su that. Doi lai, DOI TEN
// HIEN THI cua mot mau da khoa se lam bo kiem DO. Do la CO Y: ten la thu duy
// nhat cho phep con nguoi doc mot dong khoa va biet id do la cai gi. Doi ten mot
// mau da phat hanh la mot viec CO CHU DINH — sua kem dong khoa tuong ung, va
// nguoi review doi chieu git log cua file khoa (§2.6) se thay ro mot dong bi
// SUA chu khong phai THEM, va hoi vi sao. Dung cai ta muon.
//
// KhoaIdMaTu la ban XUAT KHAU cua khoaIdMaTu (KH8 task 4). Goi `game` can no de
// doi chieu danh muc DANG CHAY voi so khoa mot cach DOC LAP — neu test ben do
// tu viet lai phep sinh ma thi no chi chung minh ban sao dung, va hai ban se roi
// nhau ngay lan dau ai sua bang dau.
func KhoaIdMaTu(ten string) string { return khoaIdMaTu(ten) }

func khoaIdMaTu(ten string) string {
	var b strings.Builder
	truocLaGach := true
	for _, r := range strings.ToLower(ten) {
		r = khoaIdBoDau(r)
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			truocLaGach = false
			continue
		}
		if !truocLaGach {
			b.WriteByte('-')
			truocLaGach = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// khoaIdNhomDau: moi dong la "chu ASCII" + tat ca bien the co dau cua no (da ha
// chu thuong). Bang tra cuu tay — KHONG dung golang.org/x/text: du an giu ky luat
// "cang it phu thuoc cang tot" (§8.3 thiet ke), va bang nay khong bao gio doi.
var khoaIdNhomDau = []string{
	"aàáạảãâầấậẩẫăằắặẳẵ",
	"eèéẹẻẽêềếệểễ",
	"iìíịỉĩ",
	"oòóọỏõôồốộổỗơờớợởỡ",
	"uùúụủũưừứựửữ",
	"yỳýỵỷỹ",
	"dđ",
}

// khoaIdBangDau duoc dung mot lan luc khoi dong goi, tu khoaIdNhomDau.
var khoaIdBangDau = func() map[rune]rune {
	m := map[rune]rune{}
	for _, nhom := range khoaIdNhomDau {
		r := []rune(nhom)
		for _, coDau := range r[1:] {
			m[coDau] = r[0]
		}
	}
	return m
}()

// khoaIdBoDau tra ve chu cai ASCII tuong ung cua mot chu cai tieng Viet co dau.
func khoaIdBoDau(r rune) rune {
	if khong, co := khoaIdBangDau[r]; co {
		return khong
	}
	return r
}

// ---------------------------------------------------------------------------
// Cua vao: kiemDayDu
// ---------------------------------------------------------------------------

// kiemDayDu chay TOAN BO K1/K2/K3 con lai + doi chieu so khoa. Goi sau bo kiem
// toi thieu trong Nap. Tra loi DAU TIEN gap phai — fail-closed, khong gom loi:
// mot file YAML sai mot cho thi cho thu hai thuong la he qua, bao ca chum chi
// lam nguoi doc rooi.
func kiemDayDu(du *DuLieu, khoa *fileKhoaId) error {
	if err := kiemK1SoKhoa(khoa); err != nil {
		return fmt.Errorf("%s: %w", TenFileKhoaId, err)
	}
	if err := kiemK1DoiChieuKhoa("mau option", maCuaOption(du.MauOption), khoa.MauOption); err != nil {
		return err
	}
	if err := kiemK1DoiChieuKhoa("mau vat pham", maCuaVatPham(du.MauVatPham), khoa.MauVatPham); err != nil {
		return err
	}
	if err := kiemK3TenOption(du.MauOption); err != nil {
		return fmt.Errorf("%s: %w", TenFileVatPham, err)
	}
	if err := kiemK3ThanGoiItem(du.MauVatPham); err != nil {
		return fmt.Errorf("%s: %w", TenFileVatPham, err)
	}
	if err := kiemK2MonCoOptionPhaiDonChiec(du.CuaHang, du.MauVatPham); err != nil {
		return fmt.Errorf("%s: %w", TenFileCuaHang, err)
	}
	if err := kiemK3TenTab(du.CuaHang); err != nil {
		return fmt.Errorf("%s: %w", TenFileCuaHang, err)
	}
	if err := kiemK1SkillIdToanCuc(du); err != nil {
		return fmt.Errorf("%s: %w", TenFileChieu, err)
	}
	if err := kiemK3NhipChieu(du); err != nil {
		return fmt.Errorf("%s: %w", TenFileChieu, err)
	}
	// Boss cam chieu -> can CA boss.yaml lan chieu.yaml, nen chay o day (KH8
	// task 7, bo kiem o noidung/boss.go).
	if err := kiemBossThamChieuChieu(du); err != nil {
		return fmt.Errorf("%s: %w", TenFileBoss, err)
	}
	return nil
}

func maCuaOption(ds []MauOption) []string {
	ma := make([]string, len(ds))
	for i, o := range ds {
		ma[i] = khoaIdMaTu(o.Ten)
	}
	return ma
}

func maCuaVatPham(ds []MauVatPham) []string {
	ma := make([]string, len(ds))
	for i, m := range ds {
		ma[i] = khoaIdMaTu(m.Ten)
	}
	return ma
}

// ---------------------------------------------------------------------------
// K1 — can cuoc + so khoa (dot bien M1, M2, M3)
// ---------------------------------------------------------------------------

// kiemK1SoKhoa kiem chinh SO KHOA truoc khi dung no lam thuoc do: id lien tuc tu
// 0, ma khong rong, ma khong trung. Mot so khoa hong thi moi ket luan rut ra tu
// no deu vo nghia.
func kiemK1SoKhoa(khoa *fileKhoaId) error {
	if err := kiemK1MotBangKhoa("mau_option", khoa.MauOption); err != nil {
		return err
	}
	return kiemK1MotBangKhoa("mau_vat_pham", khoa.MauVatPham)
}

func kiemK1MotBangKhoa(ten string, ds []khoaIdMuc) error {
	if len(ds) == 0 {
		return fmt.Errorf("bang %s trong so khoa RONG — so khoa la file CHI-THEM, "+
			"mot bang rong nghia la co nguoi da xoa het (xem git log cua chinh file nay)", ten)
	}
	daThayMa := map[string]int{}
	for i, m := range ds {
		if m.Id != i {
			return fmt.Errorf("bang %s dong thu %d khai id = %d — so khoa phai lien tuc tu 0 va "+
				"CHI DUOC THEM VAO CUOI; mot lo hong o day nghia la co dong bi xoa", ten, i, m.Id)
		}
		if m.Ma == "" {
			return fmt.Errorf("bang %s id %d: ma rong — moi dong khoa phai co ma de doi chieu", ten, i)
		}
		if cu, co := daThayMa[m.Ma]; co {
			return fmt.Errorf("bang %s: ma %q khai o ca id %d va id %d — mot ma chi thuoc ve MOT id, "+
				"vinh vien (thiet ke §2.6: id KHONG BAO GIO duoc tai su dung)", ten, m.Ma, cu, i)
		}
		daThayMa[m.Ma] = i
	}
	return nil
}

// kiemK1DoiChieuKhoa la kiem GIET DOT BIEN M3 — cai duy nhat bat duoc viec XOA
// mot muc, hay doi cho hai muc roi sua lai id cho lien tuc. Ca hai thao tac deu
// qua duoc kiem "id = vi tri" cua nap.go.
func kiemK1DoiChieuKhoa(loai string, maDanhMuc []string, khoa []khoaIdMuc) error {
	if len(maDanhMuc) < len(khoa) {
		thieu := khoa[len(maDanhMuc)]
		return fmt.Errorf("%s: danh muc chi con %d muc nhung %s da khoa %d — id %d (%q) DA BIEN MAT. "+
			"Danh muc CHI DUOC DAI THEM: id da cap phat dang nam trong CSDL nguoi choi, "+
			"xoa mot muc la doi mon do cua nguoi khac. Muon ngung ban thi bo mon khoi %s, "+
			"KHONG xoa khoi danh muc",
			loai, len(maDanhMuc), TenFileKhoaId, len(khoa), thieu.Id, thieu.Ma, TenFileCuaHang)
	}
	for i, m := range khoa {
		if maDanhMuc[i] != m.Ma {
			return fmt.Errorf("%s id %d: danh muc dang la %q nhung %s khoa id nay cho %q. "+
				"Hoac hai muc bi doi cho, hoac mot muc bi xoa, hoac mot muc vua bi DOI TEN. "+
				"Neu that su co y doi ten thi sua dong khoa tuong ung — va noi ro trong commit, "+
				"vi nguoi review doi chieu file khoa voi git log cua chinh no (thiet ke §2.6)",
				loai, m.Id, maDanhMuc[i], TenFileKhoaId, m.Ma)
		}
	}
	if len(maDanhMuc) > len(khoa) {
		i := len(khoa)
		return fmt.Errorf("%s id %d (%q) co trong danh muc nhung KHONG co dong nao trong %s. "+
			"Them mot muc moi thi phai them MOT DONG KHOA tuong ung (id: %d, ma: %q) — "+
			"khong co dong do thi lan sau co nguoi xoa muc nay se khong ai bat duoc",
			loai, i, maDanhMuc[i], TenFileKhoaId, i, maDanhMuc[i])
	}
	return nil
}

// kiemK1SkillIdToanCuc: skill_id la can cuoc TOAN CUC cua mot cap chieu, khong
// phai can cuoc trong mot lang. nap.go da kiem trung trong CUNG mot lang; day mo
// rong ra ca nam lang, vi client tra cuu chieu bang skillId khong kem theo lang.
func kiemK1SkillIdToanCuc(du *DuLieu) error {
	type noiKhai struct {
		lang  int
		chieu string
		cap   int
	}
	daThay := map[int16]noiKhai{}
	for iLang := range du.ChieuTheoLang {
		for _, c := range du.ChieuTheoLang[iLang] {
			for k, cap := range c.CacCap {
				if cu, co := daThay[cap.SkillId]; co {
					return fmt.Errorf("skill_id %d khai hai lan: lang %d chieu %q cap %d, va lang %d chieu %q cap %d "+
						"— skill_id la can cuoc TOAN CUC, client tra cuu no khong kem theo lang",
						cap.SkillId, cu.lang, cu.chieu, cu.cap, iLang, c.Ten, k+1)
				}
				daThay[cap.SkillId] = noiKhai{lang: iLang, chieu: c.Ten, cap: k + 1}
			}
		}
	}
	// template_id: duy nhat TRONG MOT LANG (no la so slot cua lang do).
	for iLang := range du.ChieuTheoLang {
		daThayTpl := map[byte]string{}
		for _, c := range du.ChieuTheoLang[iLang] {
			if cu, co := daThayTpl[c.TemplateId]; co {
				return fmt.Errorf("lang %d: template_id %d dung cho ca chieu %q va %q — "+
					"template_id la so SLOT trong mot lang, hai chieu cung slot thi mot cai khong bao gio hien",
					iLang, c.TemplateId, cu, c.Ten)
			}
			daThayTpl[c.TemplateId] = c.Ten
		}
	}
	return nil
}

// KiemIdAnhQuai doi chieu `anh_id` cua MOI mau quai voi danh sach id anh quai
// HOP LE (KH8 task 8, thiet ke §6).
//
// VI SAO HAM NAY XUAT KHAU va nhan danh sach lam THAM SO, thay vi nam trong Nap
// nhu moi bo kiem khac: nguon su that cua id anh la `internal/game/assets/
// catalog.json`, mot file bi GHIM SHA trong `asset-lock.json`, va no duoc doc
// boi goi `game`. Huong import da chot cua du an la game -> noidung, TUYET DOI
// khong nguoc lai (noidung.go, dau tep) — nen goi nay khong the tu di lay danh
// sach do. Nguoi TIEM du lieu (game/khoquai.go) goi ham nay ngay sau khi nap.
//
// Ham THUAN, khong doc file, khong doc dong ho — nen test duoc bang mot danh
// sach id gia, khong can dung toi catalog that.
func KiemIdAnhQuai(mau []MauQuai, idHopLe []int) error {
	co := make(map[int]bool, len(idHopLe))
	for _, id := range idHopLe {
		co[id] = true
	}
	for _, m := range mau {
		if !co[m.AnhId] {
			return fmt.Errorf("mau quai %q: anh_id = %d KHONG co trong catalog anh da ky "+
				"(asset-lock.json, closure.mob_ids) — client khong co anh de ve con quai nay va "+
				"KHONG bao mot dong loi nao. Muon them anh quai moi thi phai lam lai goi asset va "+
				"ky lai lock, khong sua duoc bang %s", m.Ma, m.AnhId, TenFileQuai)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// K2 — tham chieu cheo (dot bien M9)
// ---------------------------------------------------------------------------

// kiemK2MonCoOptionPhaiDonChiec: mot mon BAY BAN kem option chi dung khi mau cua
// no chi chong duoc MOT cai trong mot o tui. Hai mon khac option ma chong chung
// mot o thi khi luu CSDL mot trong hai MAT SACH option — va nguoi choi thay mon
// do "tu nhien mat chi so".
func kiemK2MonCoOptionPhaiDonChiec(ds []CuaHang, mau []MauVatPham) error {
	for _, ch := range ds {
		for iTab, tab := range ch.Tabs {
			for iMon, mon := range tab.Mon {
				if len(mon.Options) == 0 {
					continue
				}
				// nap.go da kiem mau_id nam trong bien truoc khi toi day.
				m := mau[mon.MauId]
				if m.TranChong != 1 {
					return fmt.Errorf("cua hang %d tab %d mon %d: ban mau %d (%q) KEM %d option, "+
						"nhung mau do khai tran_chong = %d. Mon co option BAT BUOC don chiec "+
						"(tran_chong: 1): hai mon khac option chong chung mot o tui thi mot trong hai "+
						"mat sach option, khong mot dong log",
						ch.Id, iTab, iMon, mon.MauId, m.Ten, len(mon.Options), m.TranChong)
				}
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// K3 — mien gia tri (dot bien M4, M5, M6, M14)
// ---------------------------------------------------------------------------

// kiemK3TenOption: ten mau option PHAI chua '#'. Client thay '#' bang gia tri
// param khi ve (ItemOption.getOptionString). Ten khong co '#' thi nguoi choi
// thay mot dong chu tinh — khong loi, chi vo nghia.
func kiemK3TenOption(ds []MauOption) error {
	for i, o := range ds {
		if !strings.Contains(o.Ten, "#") {
			return fmt.Errorf("mau option %d (%q): ten khong chua '#'. Client thay '#' bang con so "+
				"param khi ve; thieu no thi nguoi choi thay mot dong chu tinh, khong thay chi so", i, o.Ten)
		}
	}
	return nil
}

// kiemK3ThanGoiItem: client goi d.mark(100000) TRUOC khi doc goi catalog -28/8
// (Controller.cs:6468). Than dai hon la client doc lech ma KHONG bao loi. Kiem o
// day, luc NAP, chu khong phai luc gui goi — dung tinh than K3.
//
// Do la mot UOC LUONG CAN TREN theo dung khuon byte cua goiCatalogItemMau
// (game/mauvatpham.go): moi mau = loai 1 + gioi_tinh 1 + ten UTF + mo_ta UTF +
// cap 1 + yeu_cau_suc 4 + icon 2 + phan 2 + nang_cap 1.
func kiemK3ThanGoiItem(ds []MauVatPham) error {
	tong := 4 // [lenh con][vsItem][subtype][soMau: short]
	for _, m := range ds {
		tong += 11 + soByteUTF(m.Ten) + soByteUTF(m.MoTa)
	}
	if tong > TranThanGoiItem {
		return fmt.Errorf("danh muc mau vat pham dung ra than goi ~%d byte, client chi mark %d "+
			"truoc khi doc (Controller.cs:6468) — dai hon la client doc LECH tu do tro di ma "+
			"KHONG bao loi. Cat bot `mo_ta`, hoac chia goi (viet code, khong phai sua so)",
			tong, TranThanGoiItem)
	}
	return nil
}

// soByteUTF dem so byte mot chuoi chiem trong mot o UTF cua giao thuc: 2 byte do
// dai + than theo Modified UTF-8 (1 byte cho U+0001..U+007F, 2 cho U+0080..U+07FF
// va cho U+0000, 3 cho phan con lai trong BMP, 6 cho tren BMP).
func soByteUTF(s string) int {
	n := 2
	for _, r := range s {
		switch {
		case r >= 0x0001 && r <= 0x007F:
			n++
		case r <= 0x07FF:
			n += 2
		case r <= 0xFFFF:
			n += 3
		default:
			n += 6
		}
	}
	return n
}

// kiemK3TenTab: client tach ten tab bang '\n' va CHI nhan 1 hoac 2 dong. Ba dong
// tro len: tab hien ten RONG, khong bao loi.
func kiemK3TenTab(ds []CuaHang) error {
	for _, ch := range ds {
		for iTab, tab := range ch.Tabs {
			if soDong := strings.Count(tab.Ten, "\n") + 1; soDong > 2 {
				return fmt.Errorf("cua hang %d tab %d (%q): ten %d dong, client chi ve duoc 1 hoac 2 "+
					"— tu 3 dong tro len tab hien ten RONG va khong bao loi", ch.Id, iTab, tab.Ten, soDong)
			}
		}
	}
	return nil
}

// NhipActorMs — nhip vong lap actor. KHONG phai lua chon noi dung: no la rang
// buoc kien truc (thiet ke §1.2). Moi con so thoi gian trong YAML phai la boi so
// cua no, neu khong thi hanh vi phu thuoc vao "roi vao nhip nao" — khong tai
// hien duoc.
const NhipActorMs = 100

func kiemK3NhipChieu(du *DuLieu) error {
	for iLang := range du.ChieuTheoLang {
		for _, c := range du.ChieuTheoLang[iLang] {
			for k, cap := range c.CacCap {
				if cap.HoiChieuMs <= 0 {
					return fmt.Errorf("lang %d chieu %q cap %d: hoi_chieu_ms = %d, phai > 0 "+
						"— hoi chieu 0 la chieu ban duoc moi nhip actor (10 lan mot giay)",
						iLang, c.Ten, k+1, cap.HoiChieuMs)
				}
				if cap.HoiChieuMs%NhipActorMs != 0 {
					return fmt.Errorf("lang %d chieu %q cap %d: hoi_chieu_ms = %d, phai la boi so cua "+
						"nhip actor %dms — so le nhip lam chieu hoi som hay muon nua nhip tuy luc bam, "+
						"khong ai tai hien duoc",
						iLang, c.Ten, k+1, cap.HoiChieuMs, NhipActorMs)
				}
				if cap.ChakraDung < 0 {
					return fmt.Errorf("lang %d chieu %q cap %d: chakra_dung am (%d)",
						iLang, c.Ten, k+1, cap.ChakraDung)
				}
				if cap.NhanLucCan < 0 {
					return fmt.Errorf("lang %d chieu %q cap %d: nhan_luc_can am (%d)",
						iLang, c.Ten, k+1, cap.NhanLucCan)
				}
				if cap.Gia < 0 {
					return fmt.Errorf("lang %d chieu %q cap %d: gia am (%d)",
						iLang, c.Ten, k+1, cap.Gia)
				}
			}
		}
	}
	return nil
}
