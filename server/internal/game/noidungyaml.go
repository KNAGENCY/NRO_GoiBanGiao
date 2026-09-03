// noidungyaml.go — CAU noi giua kho noi dung YAML (goi noidung) va cac kieu du
// lieu cua goi game.
//
// KH8 task 0. O task nay cac cua tra cuu (ChieuCuaLang, DanhSachMauVatPham,
// DanhSachCuaHang...) VAN doc du lieu Go cung nhu cu — khong doi mot hanh vi
// nao. Cac ham duoi day chi ton tai de CONG TUONG DUONG BYTE (kh8gate_test.go)
// lay duoc du lieu YAML ve dung kieu Go roi cho qua CHINH cac bo dung goi that.
// Task 3/4/5 se dao chieu: cac cua tra cuu doc tu anh chup, du lieu Go bien mat.
//
// VI SAO CHUYEN DOI NAM O GOI `game` CHU KHONG PHAI GOI `noidung`: huong import
// da chot la game -> noidung, KHONG nguoc lai. Goi noidung khong duoc biet gi ve
// khuon goi tin; no chi biet cac con so cua noi dung.
package game

import (
	"fmt"

	"nroserver/internal/noidung"
)

// chieuTuYAML doi bang chieu YAML sang [5][]MauChieu — dung khuon ma
// goiCatalogSkillTu doi.
func chieuTuYAML(du *noidung.DuLieu) [5][]MauChieu {
	var ketQua [5][]MauChieu
	for lang := range du.ChieuTheoLang {
		ds := make([]MauChieu, 0, len(du.ChieuTheoLang[lang]))
		for _, c := range du.ChieuTheoLang[lang] {
			cac := make([]CapChieu, 0, len(c.CacCap))
			for _, cap := range c.CacCap {
				cac = append(cac, CapChieu{
					SkillId:      cap.SkillId,
					Diem:         cap.Diem,
					NhanLucCan:   cap.NhanLucCan,
					ChakraDung:   cap.ChakraDung,
					HoiChieuMs:   cap.HoiChieuMs,
					Dx:           cap.Dx,
					Dy:           cap.Dy,
					MaxFight:     cap.MaxFight,
					SatThuong:    cap.SatThuong,
					Gia:          cap.Gia,
					ThemThongTin: cap.ThemThongTin,
				})
			}
			ds = append(ds, MauChieu{
				Id:            c.TemplateId,
				Ten:           c.Ten,
				CapToiDa:      c.CapToiDa,
				KieuTonChakra: c.KieuTonChakra,
				Loai:          c.Loai,
				IconId:        c.IconId,
				MoTa:          c.MoTa,
				CacCap:        cac,
			})
		}
		ketQua[lang] = ds
	}
	return ketQua
}

// mauOptionTuYAML doi bang mau option YAML sang []MauOption. Chi so cua slice
// tra ve CHINH LA id option — bo kiem cua goi noidung da khang dinh `id` bang
// vi tri truoc khi den day.
func mauOptionTuYAML(du *noidung.DuLieu) []MauOption {
	ds := make([]MauOption, 0, len(du.MauOption))
	for _, o := range du.MauOption {
		ds = append(ds, mauOptionTuNoiDung(o))
	}
	return ds
}

// mauOptionTuNoiDung doi MOT dong mau option. `Id` KHONG chep sang: kieu cua goi
// `game` khong co truong id vi id = CHI SO trong slice (client tu dat the,
// Controller.cs:6479) — giu mot ban thu hai la tao co hoi cho hai ban lech nhau.
func mauOptionTuNoiDung(o noidung.MauOption) MauOption {
	return MauOption{Ten: o.Ten, Loai: o.Loai}
}

// mauVatPhamTuYAML doi danh muc mau vat pham YAML sang []MauVatPham.
func mauVatPhamTuYAML(du *noidung.DuLieu) []MauVatPham {
	ds := make([]MauVatPham, 0, len(du.MauVatPham))
	for _, m := range du.MauVatPham {
		ds = append(ds, mauVatPhamTuNoiDung(m))
	}
	return ds
}

// mauVatPhamTuNoiDung doi MOT mau vat pham. Cung ly do tren: `Id` khong chep.
func mauVatPhamTuNoiDung(m noidung.MauVatPham) MauVatPham {
	return MauVatPham{
		Loai:        m.Loai,
		GioiTinh:    m.GioiTinh,
		Ten:         m.Ten,
		MoTa:        m.MoTa,
		Cap:         m.Cap,
		YeuCauSuc:   m.YeuCauSuc,
		IconID:      m.IconId,
		Phan:        m.Phan,
		NangCapDuoc: m.NangCapDuoc,
		TranChong:   m.TranChong,
	}
}

// cuaHangTuYAML doi danh sach cua hang YAML sang []CuaHang.
//
// Tra loi (khong am tham bo qua) khi `tien` la mot chuoi la: bo kiem cua goi
// noidung da chan roi, nhung neu ai do goi ham nay voi du lieu chua qua bo kiem
// thi mot loai tien khong nhan ra PHAI la loi, KHONG duoc mac dinh thanh vang —
// mac dinh im lang la dung cach bien mot cai gia bang ngoc thanh mot cai gia
// bang vang.
func cuaHangTuYAML(du *noidung.DuLieu) ([]CuaHang, error) {
	ds := make([]CuaHang, 0, len(du.CuaHang))
	for _, ch := range du.CuaHang {
		tabs := make([]TabCuaHang, 0, len(ch.Tabs))
		for iTab, tab := range ch.Tabs {
			mon := make([]MonHang, 0, len(tab.Mon))
			for iMon, m := range tab.Mon {
				var tien LoaiTien
				switch m.Tien {
				case noidung.TienVang:
					tien = TienVang
				case noidung.TienNgoc:
					tien = TienNgoc
				default:
					return nil, fmt.Errorf("cua hang %d tab %d mon %d: loai tien %q khong nhan ra",
						ch.Id, iTab, iMon, m.Tien)
				}
				var opts []OptionMon
				if len(m.Options) > 0 {
					opts = make([]OptionMon, 0, len(m.Options))
					for _, o := range m.Options {
						opts = append(opts, OptionMon{OptionID: o.OptionId, Param: o.Param})
					}
				}
				mon = append(mon, MonHang{
					MauID:   m.MauId,
					Gia:     m.Gia,
					Tien:    tien,
					Options: opts,
					LaMoi:   m.LaMoi,
				})
			}
			tabs = append(tabs, TabCuaHang{Ten: tab.Ten, Mon: mon})
		}
		ds = append(ds, CuaHang{ID: ch.Id, Ten: ch.Ten, Tabs: tabs})
	}
	return ds, nil
}
