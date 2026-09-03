// nap.go — bo nap YAML + bo kiem TOI THIEU cua task 0.
//
// "Toi thieu" la co y: task 0 chi can du kiem de cau "YAML sinh ra dung byte"
// co nghia. Bo kiem DAY DU (K1/K2/K3 + khoa-id.yaml) la task 1; anh chup bat
// bien + atomic.Pointer + so the he la task 2. Dung nhet chung vao day.
//
// KY LUAT FAIL-CLOSED: bat ky loi nao — cu phap, can cuoc, vuot tran — la Nap
// tra loi va KHONG tra ve du lieu mot phan. Khong bao gio co trang thai "nap
// duoc mot nua" (thiet ke §3.5).
package noidung

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"

	"nroserver/kho-noi-dung"
)

// Ten ba file YAML cua ban 1. De thanh hang de mot cho go nham ten file khong
// tro thanh "im lang bo qua mot danh muc".
const (
	TenFileChieu   = "chieu.yaml"
	TenFileVatPham = "vat-pham.yaml"
	TenFileCuaHang = "cua-hang.yaml"
	TenFileQuai    = "quai.yaml"     // KH8 task 8
	TenFileLang    = "lang.yaml"     // KH8 task 6
	TenFileBangRoi = "bang-roi.yaml" // KH8 task 9
)

// TenCacFile tra ve DANH SACH DAY DU cac file YAML ma Nap doc — NGUON DUY NHAT
// cho moi fixture test "chep kho nhung ra thu muc tam". Truoc 03/09 co BA danh
// sach cung o ba test khac nhau; task 8 them quai.yaml vao Nap thi nhanh cua no
// xanh, nhanh task 10 xanh, GOP LAI thi test nap nong do vi thu muc tam thieu
// quai.yaml. Them file moi vao Nap thi PHAI them vao day — va chi vao day.
// (Task 6 them lang.yaml: gop 03/09 dieu phoi them vao day, khong o fixture.)
func TenCacFile() []string {
	return []string{TenFileChieu, TenFileVatPham, TenFileCuaHang, TenFileKhoaId, TenFileQuai, TenFileLang, TenFileBangRoi}
}

// NapNhung nap ban NHUNG trong file thuc thi — duong dung luc KHOI DONG. Khong
// phu thuoc thu muc lam viec, khong the thieu file.
func NapNhung() (*DuLieu, error) { return Nap(khonoidung.FS) }

// Nap doc toan bo kho noi dung tu mot he thong file bat ky (embed.FS luc khoi
// dong, os.DirFS luc nap nong), kiem het, roi dung MOT anh chup hoan chinh.
//
// Trinh tu bat buoc: doc HET -> phan tich HET -> kiem HET -> moi dung. Khong
// buoc nao trong day duoc dung du lieu nua voi (thiet ke §3.3).
func Nap(fsys fs.FS) (*DuLieu, error) {
	var fChieu fileChieu
	if err := docYAML(fsys, TenFileChieu, &fChieu); err != nil {
		return nil, err
	}
	var fVatPham fileVatPham
	if err := docYAML(fsys, TenFileVatPham, &fVatPham); err != nil {
		return nil, err
	}
	var fCuaHang fileCuaHang
	if err := docYAML(fsys, TenFileCuaHang, &fCuaHang); err != nil {
		return nil, err
	}
	var fLang fileLang
	if err := docYAML(fsys, TenFileLang, &fLang); err != nil {
		return nil, err
	}

	var fQuai fileQuai
	if err := docYAML(fsys, TenFileQuai, &fQuai); err != nil {
		return nil, err
	}

	var fBangRoi fileBangRoi
	if err := docYAML(fsys, TenFileBangRoi, &fBangRoi); err != nil {
		return nil, err
	}

	du := &DuLieu{
		MauOption:  fVatPham.MauOption,
		MauVatPham: fVatPham.MauVatPham,
		CuaHang:    fCuaHang.CuaHang,
		MauQuai:    fQuai.MauQuai,
		SinhQuai:   fQuai.Sinh,
		BangRoi:    dungBangRoi(&fBangRoi),
	}
	if err := kiemChieu(fChieu.Lang); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileChieu, err)
	}
	for _, l := range fChieu.Lang {
		du.ChieuTheoLang[l.IdLang] = l.Chieu
	}
	if err := kiemLang(fLang.Lang); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileLang, err)
	}
	for _, l := range fLang.Lang {
		du.Lang[l.Id] = l
	}
	if err := kiemMauOption(du.MauOption); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileVatPham, err)
	}
	if err := kiemMauVatPham(du.MauVatPham); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileVatPham, err)
	}
	if err := kiemCuaHang(du.CuaHang, len(du.MauOption), len(du.MauVatPham)); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileCuaHang, err)
	}
	if err := kiemQuai(du.MauQuai, du.SinhQuai); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileQuai, err)
	}
	// Bo kiem DAY DU (K1/K2/K3 + so khoa id) — kiemtra.go, task 1. Chay SAU bo
	// kiem toi thieu vi nhieu kiem o do dua vao viec "id = vi tri" da dung.
	var fKhoa fileKhoaId
	if err := docYAML(fsys, TenFileKhoaId, &fKhoa); err != nil {
		return nil, err
	}
	if err := kiemDayDu(du, &fKhoa); err != nil {
		return nil, err
	}
	// KH8 task 9: bang roi kiem SAU CUNG vi no doi chieu CA BA — mau quai (ma),
	// mau vat pham (mau_id + tran chong) va so khoa id (mau_id da khoa chua).
	// Chay sau kiemDayDu la CO Y: neu danh muc vat pham lech so khoa thi loi
	// DUNG la loi lech so, khong phai "bang roi tro vao khoang khong" (he qua).
	if err := kiemBangRoi(&fBangRoi, du.MauQuai, du.MauVatPham, fKhoa.MauVatPham); err != nil {
		return nil, fmt.Errorf("%s: %w", TenFileBangRoi, err)
	}
	return du, nil
}

// NapTuThuMuc nap tu mot thu muc that tren dia — duong danh cho nap nong (task
// 2/10) va cho test muon sua mot con so roi xem bo kiem co do khong.
func NapTuThuMuc(duongDan string) (*DuLieu, error) {
	return Nap(os.DirFS(duongDan))
}

// docYAML doc mot file va phan tich vao dich. KnownFields(true) la quan trong:
// no lam mot truong go nham ten (`sat_thuog:`) thanh LOI thay vi thanh so 0 im
// lang — dung tinh than "so vuot tran bi chan o cua" cua K3.
func docYAML(fsys fs.FS, ten string, dich any) error {
	data, err := fs.ReadFile(fsys, ten)
	if err != nil {
		return fmt.Errorf("khong doc duoc %s: %w", ten, err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(dich); err != nil {
		return fmt.Errorf("%s sai cu phap YAML: %w", ten, err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Bo kiem toi thieu
// ---------------------------------------------------------------------------

func kiemChieu(langs []ChieuCuaLang) error {
	if len(langs) != SoLangCoDinh {
		return fmt.Errorf("co %d lang, phai dung %d (client cu chi biet tung nay lop)",
			len(langs), SoLangCoDinh)
	}
	daThay := map[byte]bool{}
	for i, l := range langs {
		if int(l.IdLang) != i {
			return fmt.Errorf("khoi lang thu %d khai id_lang = %d — thu tu trong file PHAI bang id",
				i, l.IdLang)
		}
		if daThay[l.IdLang] {
			return fmt.Errorf("id_lang %d khai hai lan", l.IdLang)
		}
		daThay[l.IdLang] = true
		if len(l.Chieu) == 0 {
			return fmt.Errorf("lang %d khong co chieu nao", l.IdLang)
		}
		skillDaThay := map[int16]bool{}
		for j, c := range l.Chieu {
			if c.TemplateId == 0 {
				return fmt.Errorf("lang %d chieu %d (%q): template_id = 0 — cam, client co nhanh nhoi hai cap gia khi templateId == 0",
					l.IdLang, j, c.Ten)
			}
			if c.Loai != LoaiChieuDanh && c.Loai != LoaiChieuBuff && c.Loai != LoaiChieuTuThan {
				return fmt.Errorf("lang %d chieu %q: loai = %d, may chu chi cai dat %d/%d/%d — them kieu moi la VIET CODE, khong phai them du lieu",
					l.IdLang, c.Ten, c.Loai, LoaiChieuDanh, LoaiChieuBuff, LoaiChieuTuThan)
			}
			if int(c.CapToiDa) != len(c.CacCap) {
				return fmt.Errorf("lang %d chieu %q: cap_toi_da = %d nhung khai %d cap",
					l.IdLang, c.Ten, c.CapToiDa, len(c.CacCap))
			}
			if len(c.CacCap) > 255 {
				return fmt.Errorf("lang %d chieu %q: %d cap, client doc so cap bang mot byte",
					l.IdLang, c.Ten, len(c.CacCap))
			}
			for k, cap := range c.CacCap {
				if int(cap.Diem) != k+1 {
					return fmt.Errorf("lang %d chieu %q cap thu %d: diem = %d, phai bang %d (thu tu trong file = so cap)",
						l.IdLang, c.Ten, k, cap.Diem, k+1)
				}
				if skillDaThay[cap.SkillId] {
					return fmt.Errorf("lang %d: skill_id %d trung trong cung mot lang", l.IdLang, cap.SkillId)
				}
				skillDaThay[cap.SkillId] = true
			}
		}
	}
	return nil
}

// TranToaDoBanDo: toa do khoi dau cua mot lang cuoi cung duoc gan vao
// BanDo.SpawnX/SpawnY (game/banDo.go), hai truong do la int16. Am hoac vuot
// gia tri nay la tran khi gan — kiem O DAY, luc nap, dung tinh than K3.
const TranToaDoBanDo = 32767

// kiemLang la bo kiem TOI THIEU cua lang.yaml (KH8 task 6): 5 lang, id lien
// tuc tu 0 (nhu mau_vat_pham/mau_option), ten khong rong, hai chi so chien
// dau khac 0, va toa do khoi dau nam trong mien client/BanDo doc duoc.
func kiemLang(ds []Lang) error {
	if len(ds) != SoLangCoDinh {
		return fmt.Errorf("co %d lang, phai dung %d (client cu chi gui duoc id 0..%d)",
			len(ds), SoLangCoDinh, SoLangCoDinh-1)
	}
	for i, l := range ds {
		if l.Id != i {
			return fmt.Errorf("lang vi tri %d khai id = %d — thu tu trong file PHAI bang id (nhu mau vat pham/mau option)",
				i, l.Id)
		}
		if l.Ten == "" {
			return fmt.Errorf("lang %d: ten rong", i)
		}
		if l.HpKhoiDau <= 0 {
			return fmt.Errorf("lang %d (%q): hp_khoi_dau = %d, phai > 0", i, l.Ten, l.HpKhoiDau)
		}
		if l.SatThuongKhoiDau <= 0 {
			return fmt.Errorf("lang %d (%q): sat_thuong_khoi_dau = %d, phai > 0 — bang 0 lam MOI nhan vat "+
				"lang nay gay 0 sat thuong vinh vien (game/khuvuc.go TinhSatThuong)", i, l.Ten, l.SatThuongKhoiDau)
		}
		if l.ViTriXKhoiDau < 0 || l.ViTriXKhoiDau > TranToaDoBanDo ||
			l.ViTriYKhoiDau < 0 || l.ViTriYKhoiDau > TranToaDoBanDo {
			return fmt.Errorf("lang %d (%q): vi tri khoi dau (%d,%d) ngoai bien [0,%d] — BanDo.SpawnX/SpawnY "+
				"(game/banDo.go) la int16", i, l.Ten, l.ViTriXKhoiDau, l.ViTriYKhoiDau, TranToaDoBanDo)
		}
	}
	return nil
}

func kiemMauOption(ds []MauOption) error {
	if len(ds) > TranMauOption {
		return fmt.Errorf("bang mau option %d muc, client chi doc duoc %d (mot byte khong dau)",
			len(ds), TranMauOption)
	}
	for i, o := range ds {
		if o.Id != i {
			return fmt.Errorf("mau option vi tri %d khai id = %d — tap id phai la {0..n-1} va THU TU TRONG FILE = ID (hoan vi hai dong la doi option cua moi mon do trong CSDL)",
				i, o.Id)
		}
		if o.Ten == "" {
			return fmt.Errorf("mau option %d: ten rong", i)
		}
	}
	return nil
}

func kiemMauVatPham(ds []MauVatPham) error {
	if len(ds) > TranMauVatPham {
		return fmt.Errorf("danh muc %d mau, client chi doc duoc %d (short co dau)",
			len(ds), TranMauVatPham)
	}
	for i, m := range ds {
		if m.Id != i {
			return fmt.Errorf("mau vat pham vi tri %d khai id = %d — id KHONG DI QUA MANG, client dat id = chi so vong lap va id do DANG NAM TRONG CSDL nguoi choi. Tap id phai la {0..n-1}, CHI DUOC THEM VAO CUOI",
				i, m.Id)
		}
		if m.Ten == "" {
			return fmt.Errorf("mau vat pham %d: ten rong", i)
		}
		if m.TranChong < 1 {
			return fmt.Errorf("mau vat pham %d (%q): tran_chong = %d, phai >= 1", i, m.Ten, m.TranChong)
		}
	}
	return nil
}

// kiemQuai — bo kiem cua quai.yaml (KH8 task 8).
//
// KHONG kiem `anh_id` o day: gia tri hop le cua no nam trong catalog anh DA KY
// (assets/asset-lock.json), ma goi nay khong duoc phep import goi `game`. Viec
// do lam boi KiemIdAnhQuai (kiemtra.go), goi tu game/khoquai.go luc tiem.
func kiemQuai(mau []MauQuai, sinh []SinhQuaiTrenMap) error {
	if len(mau) == 0 {
		return fmt.Errorf("khong co mau quai nao — mot ban do khong quai la mot ban do chet, "+
			"neu that su muon vay thi de `%s` co mau nhung khong co diem sinh", TenFileQuai)
	}
	theoMa := map[string]bool{}
	for i, m := range mau {
		if m.Ma == "" {
			return fmt.Errorf("mau quai thu %d: `ma` rong — `ma` la can cuoc, khoi `sinh:` tro toi mau bang chuoi nay", i)
		}
		if theoMa[m.Ma] {
			return fmt.Errorf("ma quai %q khai hai lan — mot ma chi thuoc ve MOT mau, neu khong thi `sinh:` "+
				"tro toi mau nao la chuyen may rui", m.Ma)
		}
		theoMa[m.Ma] = true
		if m.Ten == "" {
			return fmt.Errorf("mau quai %q: ten rong", m.Ma)
		}
		if m.Hp <= 0 {
			return fmt.Errorf("mau quai %q: hp = %d, phai > 0 — quai sinh ra voi hp <= 0 la quai CHET SAN: "+
				"nhipTim thay Hp<=0 la vao nhanh cho hoi sinh, nhung LucChet con la zero value nen no hoi sinh "+
				"NGAY tick dau, roi lai chet, mai mai. Khong mot dong log", m.Ma, m.Hp)
		}
		if m.SatThuongCoBan <= 0 {
			return fmt.Errorf("mau quai %q: sat_thuong_co_ban = %d, phai > 0 — quai danh 0 sat thuong van phat "+
				"goi -11/-10 moi nhip AI: nguoi choi thay hieu ung bi danh ma mau khong doi, khong hieu chuyen gi",
				m.Ma, m.SatThuongCoBan)
		}
		if m.Cap == 0 {
			return fmt.Errorf("mau quai %q: cap = 0, phai >= 1", m.Ma)
		}
		if m.CapTinhAnh == 0 {
			return fmt.Errorf("mau quai %q: cap_tinh_anh = 0, phai >= 1 (cap HIEN THI canh ten quai)", m.Ma)
		}
		if m.AnhId < 0 {
			return fmt.Errorf("mau quai %q: anh_id am (%d)", m.Ma, m.AnhId)
		}
	}
	daThayMap := map[int32]bool{}
	for _, s := range sinh {
		if daThayMap[s.Map] {
			return fmt.Errorf("map %d khai hai khoi `sinh:` — khoi thu hai se bi bo qua im lang", s.Map)
		}
		daThayMap[s.Map] = true
		if len(s.Dat) > TranTempIdQuai {
			return fmt.Errorf("map %d: %d con quai, client doc soQuai bang byte CO DAU nen toi da %d",
				s.Map, len(s.Dat), TranTempIdQuai)
		}
		daThayId := map[byte]string{}
		for i, d := range s.Dat {
			if !theoMa[d.Ma] {
				return fmt.Errorf("map %d con quai thu %d: `ma` = %q khong co trong `mau_quai` — "+
					"hoac go nham, hoac vua co nguoi doi `ma` cua mot mau ma quen sua khoi `sinh:`",
					s.Map, i, d.Ma)
			}
			if int(d.Id) > TranTempIdQuai {
				return fmt.Errorf("map %d con quai thu %d: id = %d, client doc tempId bang byte CO DAU nen phai <= %d",
					s.Map, i, d.Id, TranTempIdQuai)
			}
			if cu, co := daThayId[d.Id]; co {
				return fmt.Errorf("map %d: id %d dung cho ca %q va %q — tempId la can cuoc cua con quai "+
					"TRONG mot khu vuc; hai con trung id thi goi 54 (xin danh) luon trung con dau, con kia bat tu",
					s.Map, d.Id, cu, d.Ma)
			}
			daThayId[d.Id] = d.Ma
		}
	}
	return nil
}

func kiemCuaHang(ds []CuaHang, soMauOption, soMauVatPham int) error {
	daThayId := map[int32]bool{}
	for _, ch := range ds {
		if daThayId[ch.Id] {
			return fmt.Errorf("hai cua hang cung id %d", ch.Id)
		}
		daThayId[ch.Id] = true
		if len(ch.Tabs) > TranSoTab {
			return fmt.Errorf("cua hang %d: %d tab, toi da %d", ch.Id, len(ch.Tabs), TranSoTab)
		}
		for iTab, tab := range ch.Tabs {
			if tab.Ten == "" {
				return fmt.Errorf("cua hang %d tab %d: ten rong", ch.Id, iTab)
			}
			if len(tab.Mon) > TranSoMonMoiTab {
				return fmt.Errorf("cua hang %d tab %d: %d mon, toi da %d",
					ch.Id, iTab, len(tab.Mon), TranSoMonMoiTab)
			}
			for iMon, mon := range tab.Mon {
				if mon.Tien != TienVang && mon.Tien != TienNgoc {
					return fmt.Errorf("cua hang %d tab %d mon %d: tien = %q, chi nhan %q hoac %q",
						ch.Id, iTab, iMon, mon.Tien, TienVang, TienNgoc)
				}
				if mon.Gia < 0 {
					return fmt.Errorf("cua hang %d tab %d mon %d: gia am (%d)", ch.Id, iTab, iMon, mon.Gia)
				}
				if int(mon.MauId) < 0 || int(mon.MauId) >= soMauVatPham {
					return fmt.Errorf("cua hang %d tab %d mon %d: mau_id %d khong ton tai (co %d mau) — client se nem NullReference roi NUOT, cua hang hien do dang",
						ch.Id, iTab, iMon, mon.MauId, soMauVatPham)
				}
				if len(mon.Options) > TranSoOptionMoiMon {
					return fmt.Errorf("cua hang %d tab %d mon %d: %d option, toi da %d",
						ch.Id, iTab, iMon, len(mon.Options), TranSoOptionMoiMon)
				}
				for iOpt, o := range mon.Options {
					if int(o.OptionId) >= soMauOption {
						return fmt.Errorf("cua hang %d tab %d mon %d option %d: option_id %d >= %d mau option (client KHONG kiem bien -> IndexOutOfRange)",
							ch.Id, iTab, iMon, iOpt, o.OptionId, soMauOption)
					}
					if int(o.Param) > TranParamOption {
						return fmt.Errorf("cua hang %d tab %d mon %d option %d: param %d vuot %d",
							ch.Id, iTab, iMon, iOpt, o.Param, TranParamOption)
					}
				}
			}
		}
	}
	return nil
}
