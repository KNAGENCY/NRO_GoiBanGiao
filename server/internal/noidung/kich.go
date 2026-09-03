// kich.go — KENH KICH NAP NONG: mot diem cuoi HTTP NOI BO, chi nghe 127.0.0.1.
//
// (KH8 task 10. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// §4.1-§4.3 "kenh kich", §4 "nap nong", §7 bao mat.)
//
// # VI SAO HTTP CHU KHONG PHAI SIGHUP HAY FILE LENH
//
// SIGHUP: may chu dich la VPS Windows — Windows KHONG CO SIGHUP. Loai.
//
// File lenh kieu panel_cmd.txt (ban Java cu): poll 2 giay, xoa file, khong co
// danh tinh nguoi kich, va — nang nhat — KHONG CO PHAN HOI. Nguoi kich khong
// biet "nap thanh cong hay hong o dong nao", phai di mo log. Bai hoc §6.32:
// gui cho ai do mot thu ma khong xu ly phan tra loi la TE HON KHONG GUI.
//
// HTTP: nguoi kich nhan NGAY than tra loi — "the he 7, 20 chieu, 5 mau vat
// pham" hoac "cua-hang.yaml dong 14: mau_id 9 khong ton tai". Do la ca ly do
// chon no.
//
// # BA RUI RO BAO MAT VA BA CHOT (§4.3)
//
//  1. LO RA INTERNET. Chot: BIND CHI 127.0.0.1 (hoac ::1 / localhost). Go nham
//     mot chu — "0.0.0.0:9090" — la mo quyen doi luat choi ra doi. Vi vay
//     MoKenhNapLai TU CHOI KHOI DONG neu dia chi khong phai loopback, va con
//     kiem lai dia chi THAT SU cua listener sau khi bind (mot ten mien co the
//     phan giai ra dia chi ngoai). Xem KiemDiaChiNapLai + test.
//
//  2. BAT KY TIEN TRINH NAO TREN CUNG MAY CUNG GOI DUOC. Tren VPS Windows chay
//     AutoAdminLogon, "cung may" khong phai mot bien gioi manh. Chot: mot KHOA
//     DUNG CHUNG doc tu bien moi truong, gui trong header; thieu hoac sai ->
//     403 + GHI LOG. Noi thang: day CHUA PHAI xac thuc that (khong co danh
//     tinh, khong het han, khong quay vong) — no la mot cai chot cua, du cho
//     "khong ai vo tinh kich duoc", khong du cho "may bi chiem van an toan".
//
//  3. AI KICH, LUC NAO, THE HE NAO. Ban Java cu khong tra loi duoc. Chot: moi
//     lan nap ghi MOT dong log co the he moi + tom tat. Nhat ky sau hon (ai doi
//     so gi) la viec cua `git log` tren chinh cac file YAML.
//
// # CAI GI KHONG NAP NONG DUOC
//
// Kenh nay chi trao ANH CHUP NOI DUNG (chieu / vat pham / cua hang). No KHONG
// dung lai ai, KHONG doi ban do, KHONG doi giao thuc, KHONG doi lua chon da
// luu cua nguoi choi. Cu the ba thu nguoi van hanh hay tuong nham:
//
//   - chieuDangChon cua mot nguoi dang choi KHONG doi theo — no la du lieu CUA
//     NGUOI CHOI trong CSDL, khong phai noi dung. Co y.
//   - goi catalog -28/7 da gui cho nguoi DANG online khong duoc gui lai; ho
//     thay so moi sau khi dang nhap lai. Nguoi dang nhap MOI thay ngay.
//   - so luong lang (SoLangCoDinh) va khuon byte cua goi la CODE, khong phai
//     du lieu — doi chung phai bien dich lai.
//
// # HAI LOP CHONG "HAI THE HE NHAY"
//
// Kho.Thay (kho.go) da dung CompareAndSwap nen hai nguoi nap cung luc khong
// bao gio tao ra hai anh chup cung so the he. Nhung do la lop thu HAI. Lop thu
// NHAT o day: mot mutex bao ca doan "doc dia -> kiem -> Thay", de hai request
// dong thoi khong doc cung mot thu muc dang duoc sua do dang, va de dong log
// the he khong dan xen. Xem TestNapLaiHaiRequestSongSong.
package noidung

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Duong dan hai diem cuoi. De thanh hang de mot cho go nham duong dan khong
// tro thanh "404 im lang" ma nguoi van hanh tuong la kenh chet.
const (
	DuongNapLai = "/admin/nap-lai"
	DuongTheHe  = "/admin/the-he"
)

// TenHeaderKhoaNapLai la header mang khoa dung chung. Dung header rieng chu
// khong phai Authorization: de khong ai nham tuong day la xac thuc that (xem
// rui ro 2 o dau tep).
const TenHeaderKhoaNapLai = "X-Nro-Khoa"

// hanGioDocHeader chan mot ket noi mo ra roi im lang chiem cho mai mai.
const hanGioDocHeader = 5 * time.Second

// KenhNapLai la diem cuoi HTTP noi bo dang chay. Tao qua MoKenhNapLai, tat qua
// Dung. Zero value KHONG dung duoc.
type KenhNapLai struct {
	kho    *Kho
	thuMuc string // "" = chi co ban nhung, nap lai bi tu choi 501
	khoa   string
	log    *slog.Logger

	mayChu   *http.Server
	langNghe net.Listener

	// muNap bao ca doan "doc dia -> kiem -> Thay" (lop thu nhat, xem dau tep).
	muNap sync.Mutex
}

// KiemDiaChiNapLai tra loi neu dia chi KHONG phai loopback.
//
// Chap nhan: "127.0.0.1:9090", "[::1]:9090", "localhost:9090".
// Tu choi: "0.0.0.0:9090", ":9090" (bind moi card mang), ten mien ngoai.
//
// ":9090" bi tu choi CO Y du no rat quen tay trong Go: `net.Listen("tcp",
// ":9090")` nghe tren MOI dia chi, tuc lo ra internet — dung dieu §4.3 rui ro 1
// cam. Muon kich tu xa thi di qua SSH tunnel, khong noi bind.
func KiemDiaChiNapLai(diaChi string) error {
	host, cong, err := net.SplitHostPort(strings.TrimSpace(diaChi))
	if err != nil {
		return fmt.Errorf("dia chi kenh nap lai khong hop le (can dang 127.0.0.1:<cong>): %w", err)
	}
	if cong == "" {
		return errors.New("dia chi kenh nap lai thieu so cong (can dang 127.0.0.1:<cong>)")
	}
	host = strings.Trim(host, "[]")
	if host == "" {
		return errors.New(`dia chi kenh nap lai rong phan host (":<cong>" nghe tren MOI card mang = lo ra internet) — hay ghi ro 127.0.0.1:<cong>`)
	}
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("host %q khong phai dia chi IP loopback — kenh nap lai CHI duoc bind 127.0.0.1", host)
	}
	if !ip.IsLoopback() {
		return fmt.Errorf("dia chi %q khong phai loopback — kenh nap lai CHI duoc bind 127.0.0.1 (§4.3 rui ro 1)", host)
	}
	return nil
}

// MoKenhNapLai mo diem cuoi va BAT DAU phuc vu tren mot goroutine rieng.
//
//	diaChi  — PHAI loopback (xem KiemDiaChiNapLai). Rong = goi nham, tra loi.
//	kho     — kho dang chay; Thay() se duoc goi tren chinh no.
//	thuMuc  — thu muc YAML tren dia. RONG la hop le: kenh van mo (de doc the he)
//	          nhung POST nap lai tra 501 kem loi giai thich.
//	khoa    — khoa dung chung. RONG = tu choi mo (fail-closed, §4.3 rui ro 2).
//	log     — nil thi dung slog.Default().
func MoKenhNapLai(diaChi string, kho *Kho, thuMuc, khoa string, log *slog.Logger) (*KenhNapLai, error) {
	if kho == nil {
		return nil, errors.New("MoKenhNapLai: kho nil")
	}
	if strings.TrimSpace(khoa) == "" {
		return nil, errors.New("MoKenhNapLai: thieu khoa dung chung — dat NRO_NAP_LAI_KHOA truoc khi bat kenh (§4.3 rui ro 2)")
	}
	if err := KiemDiaChiNapLai(diaChi); err != nil {
		return nil, err
	}
	if log == nil {
		log = slog.Default()
	}

	langNghe, err := net.Listen("tcp", strings.TrimSpace(diaChi))
	if err != nil {
		return nil, fmt.Errorf("khong bind duoc kenh nap lai tren %s: %w", diaChi, err)
	}
	// Kiem LAI dia chi THAT SU sau khi bind: "localhost" tren mot may cau hinh
	// la co the phan giai ra dia chi ngoai. Kiem chuoi truoc khi bind la chua du.
	if err := KiemDiaChiNapLai(langNghe.Addr().String()); err != nil {
		_ = langNghe.Close()
		return nil, fmt.Errorf("listener thuc te KHONG phai loopback, dong lai ngay: %w", err)
	}

	k := &KenhNapLai{kho: kho, thuMuc: thuMuc, khoa: khoa, log: log, langNghe: langNghe}
	mux := http.NewServeMux()
	mux.HandleFunc(DuongNapLai, k.xuLyNapLai)
	mux.HandleFunc(DuongTheHe, k.xuLyTheHe)
	k.mayChu = &http.Server{Handler: mux, ReadHeaderTimeout: hanGioDocHeader}

	go func() {
		if err := k.mayChu.Serve(langNghe); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("kenh nap lai dung bat thuong", "loi", err)
		}
	}()
	log.Info("kenh nap lai noi dung dang nghe (CHI loopback)",
		"dia_chi", langNghe.Addr().String(),
		"thu_muc_yaml", thuMuc,
		"the_he_hien_tai", kho.TheHe())
	return k, nil
}

// DiaChi la dia chi THAT SU dang nghe (co ich khi cau hinh cong 0 trong test).
func (k *KenhNapLai) DiaChi() string { return k.langNghe.Addr().String() }

// Dung dong diem cuoi, cho cac request dang do xong trong han cua ctx.
//
// An toan goi nhieu lan. Sau khi tra ve, khong con request nao co the goi
// Kho.Thay qua kenh nay nua — day la ly do no phai chay TRUOC buoc tat actor
// trong trinh tu tat may chu (xem main.go).
func (k *KenhNapLai) Dung(ctx context.Context) error {
	if k == nil || k.mayChu == nil {
		return nil
	}
	err := k.mayChu.Shutdown(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// ---------------------------------------------------------------------------
// Hai diem cuoi
// ---------------------------------------------------------------------------

// TomTatNoiDung la vai con so du de nguoi kich TIN rang bo so lieu vua nap la
// bo ho vua sua — "20 chieu, 5 mau vat pham" chu khong phai mot chu "OK" tron.
type TomTatNoiDung struct {
	SoChieu      int `json:"so_chieu"`
	SoMauVatPham int `json:"so_mau_vat_pham"`
	SoMauOption  int `json:"so_mau_option"`
	SoCuaHang    int `json:"so_cua_hang"`
}

func tomTatTu(du *DuLieu) TomTatNoiDung {
	t := TomTatNoiDung{
		SoMauVatPham: len(du.MauVatPham),
		SoMauOption:  len(du.MauOption),
		SoCuaHang:    len(du.CuaHang),
	}
	for i := range du.ChieuTheoLang {
		t.SoChieu += len(du.ChieuTheoLang[i])
	}
	return t
}

type traLoiNapLai struct {
	TheHe  uint64        `json:"the_he"`
	Luc    string        `json:"luc"`
	TomTat TomTatNoiDung `json:"tom_tat"`
}

// chanCua kiem phuong thuc + khoa dung chung. Tra ve false nghia la da tra loi
// xong cho client, nguoi goi dung lai ngay.
func (k *KenhNapLai) chanCua(w http.ResponseWriter, r *http.Request, phuongThuc string) bool {
	if r.Method != phuongThuc {
		w.Header().Set("Allow", phuongThuc)
		http.Error(w, "phuong thuc khong dung, can "+phuongThuc, http.StatusMethodNotAllowed)
		return false
	}
	// So sanh thuong: khoa nay khong bao ve bi mat nguoi dung, va mot phep so
	// sanh hang so thoi gian o day khong doi duoc su that "day chua phai xac
	// thuc that" (xem rui ro 2). Cai QUAN TRONG la GHI LOG moi lan sai — do
	// moi la thu nguoi van hanh nhin thay khi co ai do mo cua.
	if r.Header.Get(TenHeaderKhoaNapLai) != k.khoa {
		k.log.Warn("kenh nap lai: TU CHOI mot yeu cau sai khoa",
			"tu", r.RemoteAddr, "duong", r.URL.Path, "co_header_khoa", r.Header.Get(TenHeaderKhoaNapLai) != "")
		http.Error(w, "sai khoa dung chung", http.StatusForbidden)
		return false
	}
	return true
}

// xuLyTheHe: GET /admin/the-he — "may chu dang chay bo so nao".
func (k *KenhNapLai) xuLyTheHe(w http.ResponseWriter, r *http.Request) {
	if !k.chanCua(w, r, http.MethodGet) {
		return
	}
	ac := k.kho.AnhChup()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(traLoiNapLai{
		TheHe: ac.TheHe(), Luc: ac.Luc().Format(time.RFC3339),
	})
}

// xuLyNapLai: POST /admin/nap-lai.
//
// Trinh tu BAT BUOC (fail-closed, §4): doc + kiem TOAN BO tu dia -> chi khi
// KHONG loi moi Thay. Loi validator -> giu nguyen ban dang chay va tra 422 kem
// THONG DIEP NGUYEN VAN cua validator: chu du an khong phai lap trinh vien, cau
// "cua-hang.yaml dong 14: mau_id 9 khong ton tai" la thu duy nhat dung duoc.
func (k *KenhNapLai) xuLyNapLai(w http.ResponseWriter, r *http.Request) {
	if !k.chanCua(w, r, http.MethodPost) {
		return
	}
	if strings.TrimSpace(k.thuMuc) == "" {
		http.Error(w,
			"khong nap lai duoc: may chu dang chay ban noi dung NHUNG trong file thuc thi. "+
				"Muon nap nong thi dat NRO_KHO_NOI_DUNG tro toi thu muc YAML tren dia roi khoi dong lai.",
			http.StatusNotImplemented)
		return
	}

	// Lop thu nhat chong "hai the he nhay" — xem dau tep.
	k.muNap.Lock()
	defer k.muNap.Unlock()

	du, err := NapTuThuMuc(k.thuMuc)
	if err != nil {
		cu := k.kho.AnhChup()
		k.log.Warn("nap nong THAT BAI, giu nguyen ban dang chay",
			"the_he_giu_nguyen", cu.TheHe(), "thu_muc", k.thuMuc, "loi", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "NAP LAI THAT BAI — may chu VAN chay the he %d nhu cu.\n\n%v\n",
			cu.TheHe(), err)
		return
	}

	theHe := k.kho.Thay(du)
	ac := k.kho.AnhChup()
	tt := tomTatTu(du)
	k.log.Info("nap nong THANH CONG, da trao the he moi",
		"the_he", theHe, "thu_muc", k.thuMuc,
		"so_chieu", tt.SoChieu, "so_mau_vat_pham", tt.SoMauVatPham,
		"so_mau_option", tt.SoMauOption, "so_cua_hang", tt.SoCuaHang,
		"tu", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(traLoiNapLai{
		TheHe: theHe, Luc: ac.Luc().Format(time.RFC3339), TomTat: tt,
	})
}
