// hanhtrang.go — KH7 Task 4: bang phien dich (noi chua, o so) -> vat_pham.id.
//
// Dac ta: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md
// §4.1 (id mon KHONG BAO GIO xuong client), §2.3 (kiem-roi-ghi), §1.6/§6.9
// (ky luat actor: ai so huu du lieu nao), §5.3 (vat pham ghi thang bang giao
// dich CSDL tren goroutine PHIEN).
//
// Vi sao can bang nay:
//
// Ba goi -36/-37/-35 (goihanhtrang.go) chi gui XUONG client mau_id + so luong +
// option. `vat_pham.id` — dinh danh that cua MOT MON — khong bao gio roi khoi
// may chu (§4.1). Cho nen khi nguoi choi bam "dung mon o o 3", client chi gui
// duoc con so 3. May chu phai tu tra o 3 ra id. Bang nay LA cai tra do.
//
// Vi sao phai dung cung luc voi goi:
//
// Neu bang dung tu mot lat cat du lieu KHAC voi lat cat da ghi ba goi, client
// se thay o 3 la mon A trong khi may chu tra o 3 ra mon B — nguoi choi bam ban
// cai nay lai mat cai kia, khong loi, khong log. Vi vay ham tien ich
// DungBangVaGoi dung bang VA ba goi tu DUNG MOT lat cat, mot lan goi.
package game

import (
	"fmt"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// khoaO la khoa tra cuu: cap (noi chua, o so). PHAI co ca noi chua — o so 3 cua
// TUI va o so 3 cua RUONG la hai mon khac nhau; khoa chi theo o so la tra nham.
type khoaO struct {
	noiChua string
	oSo     int16
}

// BangHanhTrang la ban phien dich (noi chua, o so) -> vat_pham.id cua MOT nguoi
// choi, cho ba noi chua co o: TUI, MAC, RUONG.
//
// AI SO HUU (§1.6/§6.9): goroutine PHIEN cua chinh nguoi choi do — cu the la
// *BoXuLy trong login.go. Khong con tro nao toi bang nay duoc chia se ra ngoai
// goroutine ay, nen khong can khoa.
//
// VI SAO KHONG DE TRONG ACTOR KHU VUC (§6.9): actor khu vuc so huu nhung gi
// NHIEU nguoi cung nhin thay (vi tri, mau, quai, boss). Hanh trang chi mot
// nguoi nhin thay, va vat pham duoc ghi THANG bang giao dich CSDL tren goroutine
// phien, khong qua ThoGhi va khong qua actor khu vuc (§5.3). Nhet bang nay vao
// actor khu vuc se bat moi thao tac tui do phai xep hang sau ca luot mo phong
// cua ban do — cham va sai tang so huu.
//
// KHI NAO PHAI DUNG LAI: sau MOI lan bo o thay doi — doi o, mua, ban, nhat,
// vut, mac/thao do (task 6/7/9 goi). Cach an toan la luon dung lai bang CUNG
// LUC voi ba goi bang DungBangVaGoi, dung sua bang tai cho.
type BangHanhTrang struct {
	m map[khoaO]int64
}

// noiChuaCoO la ba noi chua co khai niem "o". DUOI_DAT / DA_HUY khong co o
// (store.VatPham.OSo = nil) nen khong nam trong bang; KY_GUI danh cho ban 3.
var noiChuaCoO = map[string]bool{
	store.NoiChuaTui:   true,
	store.NoiChuaMac:   true,
	store.NoiChuaRuong: true,
}

// DungBangHanhTrang dung bang tu mot lat cat vat pham.
//
// KIEM-ROI-GHI (§2.3): tra loi va KHONG tra bang nua vo neu gap bat cu dieu gi
// sai — hai mon cung mot o, o ngoai [0,tranSoO), mon khong co o_so, noi chua
// khong phai ba loai tren. Mot bang nua vo con te hon khong co bang.
//
// Mon o DUOI_DAT / DA_HUY / KY_GUI trong lat cat duoc BO QUA lang le (chung
// khong thuoc hanh trang), khong phai loi.
func DungBangHanhTrang(dsMon []store.VatPham) (*BangHanhTrang, error) {
	m := make(map[khoaO]int64, len(dsMon))
	for i := range dsMon {
		v := &dsMon[i]
		if !noiChuaCoO[v.NoiChua] {
			if v.OSo != nil {
				return nil, fmt.Errorf("bang hanh trang: mon id=%d o noi chua %q khong co o nhung o_so=%d",
					v.ID, v.NoiChua, *v.OSo)
			}
			continue
		}
		if v.OSo == nil {
			return nil, fmt.Errorf("bang hanh trang: mon id=%d o %s khong co o_so", v.ID, v.NoiChua)
		}
		o := *v.OSo
		if o < 0 || int(o) >= tranSoO {
			return nil, fmt.Errorf("bang hanh trang: mon id=%d co o_so=%d, ngoai bien [0,%d)",
				v.ID, o, tranSoO)
		}
		k := khoaO{noiChua: v.NoiChua, oSo: o}
		if cu, trung := m[k]; trung {
			return nil, fmt.Errorf("bang hanh trang: hai mon cung o %d cua %s (id=%d va id=%d)",
				o, v.NoiChua, cu, v.ID)
		}
		m[k] = v.ID
	}
	return &BangHanhTrang{m: m}, nil
}

// TraId tra vat_pham.id cua mot o. co = false nghia la o rong (hoac noi chua
// khong co o) — nguoi goi PHAI coi day la "khong lam gi", khong duoc doan.
func (b *BangHanhTrang) TraId(noiChua string, oSo int16) (int64, bool) {
	if b == nil {
		return 0, false
	}
	id, co := b.m[khoaO{noiChua: noiChua, oSo: oSo}]
	return id, co
}

// SoMon tra so o dang co mon trong bang (dung cho test va log).
func (b *BangHanhTrang) SoMon() int {
	if b == nil {
		return 0
	}
	return len(b.m)
}

// HanhTrangDayDu gom bang phien dich VA ba goi gui xuong client — ket qua cua
// DungBangVaGoi.
type HanhTrangDayDu struct {
	Bang     *BangHanhTrang
	GoiTui   *protocol.Message // -36
	GoiMac   *protocol.Message // -37
	GoiRuong *protocol.Message // -35
}

// DungBangVaGoi dung bang + ba goi tu DUNG MOT lat cat dsMon, trong MOT lan
// goi. Day la ham ma task 5 (noi vao login.go / ME_INFO) phai dung — no lam
// cho viec "dung lech" tro thanh khong the: khong co duong nao dung bang tu lat
// cat nay roi dung goi tu lat cat khac.
//
// Loi o bat cu buoc nao => tra loi va KHONG tra gi ca (§2.3).
func DungBangVaGoi(dsMon []store.VatPham, head int16, soOTui, soOMac, soORuong int) (*HanhTrangDayDu, error) {
	bang, err := DungBangHanhTrang(dsMon)
	if err != nil {
		return nil, err
	}
	gTui, err := GoiTui(dsMon, soOTui)
	if err != nil {
		return nil, err
	}
	gMac, err := GoiDoDangMac(dsMon, head, soOMac)
	if err != nil {
		return nil, err
	}
	gRuong, err := GoiRuong(dsMon, soORuong)
	if err != nil {
		return nil, err
	}
	return &HanhTrangDayDu{Bang: bang, GoiTui: gTui, GoiMac: gMac, GoiRuong: gRuong}, nil
}
