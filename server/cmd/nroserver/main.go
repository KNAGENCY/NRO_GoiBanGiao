// Chuong trinh nroserver: diem khoi dong may chu game.
//
// Doc cau hinh tu bien moi truong, mo CSDL, roi lang nghe ket noi TCP cho
// toi khi nhan tin hieu dung (Ctrl+C) hoac loi khong the phuc hoi.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"nroserver/internal/chanip"
	"nroserver/internal/game"
	mynet "nroserver/internal/net"
	"nroserver/internal/store"
)

// soThoGhiMacDinh la so goroutine ghi CSDL song song toi da trong MOI dot xa
// hang doi cua ThoGhi (Ke hoach 11 Task 3, xem MoThoGhi). Chon 4, khong phai
// 1 (qua cham khi hang doi phinh) hay hang chuc (qua tai pool ket noi
// PostgreSQL cua mot VPS nho): dac ta §9 muc 1 lay VI DU 100 nguoi choi dong
// thoi voi nhip luu 60 giay chi ra trung binh ~1,7 lan ghi/GIAY — 4 goroutine
// ghi song song la du du de xa het mot dot hang doi (vd bung no cuc bo luc
// nhip luu dinh ky roi vao dung 100 nguoi cung luc) ma khong can nhieu hon.
const soThoGhiMacDinh = 4

// hanGioTatMayChu la han gio cho toan bo qua trinh "luu roi tat" (dac ta §4
// diem C de xuat 30 giay o task-6-brief.md Step 3): du CSDL treo that su luc
// dang tat may chu, tien trinh van PHAI thoat duoc trong khoang thoi gian co
// the du doan, khong treo vo han (dung tinh than context.WithTimeout da giai
// thich o hanGioGhi, thogi.go).
const hanGioTatMayChu = 30 * time.Second

// mucNhatKyTuChuoi dich gia tri bien moi truong NRO_MUC_NHAT_KY thanh slog.Level.
// hopLe == false nghia la chuoi khong nhan dang duoc; ham van tra ve
// slog.LevelInfo trong truong hop do de goi noi khong can xu ly rieng, nhung
// nen canh bao nguoi van hanh biet gia tri ho dat bi bo qua.
func mucNhatKyTuChuoi(s string) (muc slog.Level, hopLe bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "":
		return slog.LevelInfo, true
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

// docSoNguyenMoiTruong doc mot so nguyen duong tu bien moi truong, roi ve
// macDinh neu bien chua dat hoac dat gia tri khong dung.
//
// VI SAO ba con so cua hai cai gac theo dia chi mang phai chinh duoc bang bien
// moi truong, khac voi moi hang so khac cua may chu: chung la cho DUY NHAT ma
// mot con so dat qua chat se DUOI NGUOI CHOI THAT ra ngoai. Nha mang Viet Nam
// dung NAT lon — hang nghin nguoi choi co the cung mot dia chi IP cong cong.
// Neu ngay mo cua thay nguoi choi cua mot nha mang bao "khong vao duoc" hang
// loat, nguoi van hanh phai noi so ra duoc NGAY LUC DO, khong the ngoi doi mot
// ban bien dich moi.
//
// Gia tri sai (chu, so am, so 0) KHONG duoc lang le bien thanh "khong gioi han
// gi ca" — kieu hong do khong co trieu chung nao cho toi dung hom bi tan cong.
// Roi ve mac dinh va ghi mot dong canh bao de nguoi van hanh biet gia tri ho
// dat da bi bo qua.
func docSoNguyenMoiTruong(ten string, macDinh int) int {
	tho := strings.TrimSpace(os.Getenv(ten))
	if tho == "" {
		return macDinh
	}
	n, err := strconv.Atoi(tho)
	if err != nil || n <= 0 {
		slog.Warn("bien moi truong khong doc duoc thanh so nguyen duong, dung gia tri mac dinh",
			"bien", ten, "gia_tri_nhan_duoc", tho, "mac_dinh", macDinh)
		return macDinh
	}
	return n
}

// khoDong la phan toi thieu store.Store can de don dep luc tat may chu — CHI
// Close(), tach thanh interface HEP de test tiem duoc "kho gia" (khong can
// PostgreSQL that), dung khuon voi Kho/KhoLuu (internal/game). *store.Store
// thoa man interface nay tu nhien.
type khoDong interface {
	Close()
}

type mayChuCoTheTat interface {
	Shutdown(context.Context) error
}

// loaCoTheDung la phan toi thieu *game.BoPhatThanh ma trinh tu tat may chu
// can — CHI Dung(), tach thanh interface HEP cung ly do voi khoDong o tren:
// de test tiem duoc mot cai loa gia BIET GHI LAI no bi dung vao luc nao trong
// day don dep. *game.BoPhatThanh thoa man interface nay tu nhien.
type loaCoTheDung interface {
	Dung()
}

// lichBossCoTheDung la phan toi thieu *game.BoLichBoss ma trinh tu tat may
// chu can — CHI Dung(), cung ly do va cung khuon voi loaCoTheDung o tren: de
// test tiem duoc mot bo dieu phoi gia BIET GHI LAI no bi dung vao luc nao
// trong day don dep (xem lichboss_test.go). *game.BoLichBoss thoa man
// interface nay tu nhien.
type lichBossCoTheDung interface {
	Dung()
}

// noiLichBoss dung bo dieu phoi lich boss CHUNG cua may chu tu kho noi dung
// boss (noidungboss.go, Task 15) roi BatDau() no ngay — day la "diem noi day"
// cua KH6b Task 14 phan hai.
//
// VI SAO PHAI CO HAM NAY, cung ly do voi noiBoPhatThanh o tren: truoc phan
// hai nay, khong noi nao trong ca du an goi game.MoBoLichBoss/BatDau(). Toan
// bo bo dieu phoi (Task 13, da qua review) nam do KHONG AI GOI TOI — may chu
// khong bao gio sinh mot con boss nao, ma khong co lay MOT dong log bao hieu
// dieu do (SoKhuVuc.MoSoKhuVuc() im lang theo dung CHAN-3, xem lichboss.go).
//
// game.MucLichTuNoiDungBoss chuyen kho noi dung (NoiDungBoss, mang ca Ma va
// ViTri chua dung toi o Task 14) sang dinh dang bo dieu phoi doc duoc
// ([]MucLichBoss, chi con DuLieu). Xem chu thich ham do de biet vi sao hai
// kieu tach rieng.
//
// GOI DUNG MOT LAN luc khoi dong, NGAY SAU noiBoPhatThanh (xem khung chu
// thich "CHO NOI BO DIEU PHOI LICH BOSS" trong main() — bo dieu phoi la
// NGUON sinh thong bao boss, phai co loa som hon no thi thong bao dau tien
// moi khong roi vao hu khong).
func noiLichBoss(soKhuVuc *game.SoKhuVuc) *game.BoLichBoss {
	lich := game.MucLichTuNoiDungBoss(game.DanhSachNoiDungBoss())
	lichBoss := game.MoBoLichBoss(soKhuVuc, lich)
	lichBoss.BatDau()
	return lichBoss
}

// noiBoPhatThanh tao bo phat thanh CHUNG cua may chu roi cam no vao ca 5 khu
// vuc — day la "diem noi day" cua KH6b Task 14.
//
// VI SAO PHAI CO HAM NAY, va vi sao no la MOT ham rieng chu khong phai hai
// dong nam giua main(): truoc Task 14, khong noi nao trong ca du an goi
// GanLoa. Ha tang phat thanh (hang doi, goroutine, goi 93, nhip chong ngap)
// da xay xong tu truoc va DA QUA REVIEW, nhung KhuVuc.xepThongBao doc
// k.loa.Load() thay nil thi LANG LE `return` — nen moi thong bao boss roi vao
// hu khong ma khong de lai MOT DONG LOG nao. Dung lai vet xe cua Task 6
// (MoLuuTru): thu gi khong duoc noi day o mot cho GOI DUOC TU TEST thi som
// muon cung nam do khong ai goi toi, test tu che van xanh ma may chu that van
// tho.
//
// GanLoa (sokhuvuc.go) lam TRON ca hai dau day trong mot loi goi — cam loa vao
// tung khu vuc VA BatDau() goroutine phat voi duong phat la
// SoKhuVuc.PhatToanServer. Vi vay o day khong duoc tu goi BatDau them lan nua:
// BatDau hai lan la hai goroutine cung `defer close(b.xong)`, va cai thu hai
// se panic "close of closed channel" tren goroutine nen dung luc tat may chu
// (thongbao.go da chan bang CompareAndSwap tu 28/08, nhung dung dua vao lop
// chan do de goi bua).
//
// GOI DUNG MOT LAN luc khoi dong. Goi lan hai se tao cai loa THU HAI, de len
// cai cu trong tung khu vuc, va cai cu thanh goroutine mo coi khong ai tat.
func noiBoPhatThanh(soKhuVuc *game.SoKhuVuc) *game.BoPhatThanh {
	loa := game.MoBoPhatThanh()
	soKhuVuc.GanLoa(loa)
	return loa
}

// donDepKhiTat gom TOAN BO trinh tu tat may chu an toan (Ke hoach 11 Task 6,
// dac ta §4 diem C) thanh MOT ham GOI DUOC TU TEST — tach khoi main() va
// KHONG dinh chat os.Exit (brief Task 6: "Tách phần dọn dẹp khi tắt thành
// hàm gọi được từ test, đừng dính chặt os.Exit").
//
// Thu tu BAT BUOC:
//
//  0. boLuu.Dung() — dung TRUOC TIEN bo luu dinh ky (Ke hoach 11 Task 4).
//     boLuu KHONG nam trong 4 buoc dac ta §4 diem C liet ke (LuuTatCa/
//     DongVaChoXong/DungTatCa/kho.Close) vi day la manh MOI duoc noi day o
//     Task 6 nay — nhung PHAI dung no TRUOC buoc 2, vi neu khong, ticker 60
//     giay cua no co the tu goi LuuTatCa->Xep() TREN MOT GOROUTINE DOC LAP,
//     KHONG dong bo gi voi luong tat may chu chinh — neu lan Xep do roi
//     dung luc SAU khi thoGhi.DongVaChoXong() (ben trong buoc 2 duoi day) da
//     tra ve, ban ghi do se nam lai trong hang doi VINH VIEN, khong ai flush nua
//     (dung CHINH XAC canh bao cua brief Task 6: "gọi Xep() SAU khi
//     DongVaChoXong() trả về sẽ mất dữ liệu âm thầm... kể cả KhiDong() của
//     phiên đứt lúc đó" — boLuu la mot nguon Xep() DOC LAP giong het
//     KhiDong(), can chan truoc TIEN cho chac).
//
//  1. lichBoss.Dung() — dung bo dieu phoi lich boss (KH6b Task 14 phan hai),
//     NGAY TRUOC loa.Dung(). Cung logic voi (a)/(b) o duoi cho loa, ap dung
//     mot buoc som hon: lichBoss la NGUON SINH ra thong bao (boss xuat hien
//     goi k.Gui(suKienPhatGoi) qua loa), loa CHI LA duong ONG cho thong bao do
//     chay qua. Tat ong truoc ma nguon van con bom thi thong bao moi chi con
//     biet nam chet trong hang doi cua loa; tat nguon truoc, ong sau moi dam
//     bao khong con gi bom vao mot cai ong da nguoi.
//
//  2. loa.Dung() — dung bo phat thanh toan may chu (KH6b Task 14 phan mot).
//     Vi sao dung O DAY, khong som hon va khong muon hon: xem doan "VE VI TRI
//     CUA loa.Dung()" ngay duoi danh sach nay — ly do (a)/(b)/(c) o do ap dung
//     dung y het cho lichBoss.Dung(), sai khac duy nhat la lichBoss phai dung
//     TRUOC loa vi no la nguon con loa la ong dan (xem buoc 1). Ca hai (buoc
//     1 va 2) DEU PHAI xong TRUOC buoc 3 (game.TatMayChu): ca hai deu day su
//     kien vao hop thu GAMEPLAY cua khu vuc (lichBoss qua k.SinhBoss(), loa
//     qua k.Gui(suKienPhatGoi)) — thu chen giua LuuTatCa (dau buoc 3) va han
//     100ms cua no se lam anh chup nguoi choi bi rong, y het truong hop (a)
//     duoi day.
//
//  3. game.TatMayChu(ctx, soKhuVuc, thoGhi) — GOM CA BA buoc con lai cua dac
//     ta §4 diem C: LuuTatCa (xin anh chup TOAN BO nguoi choi o TAT CA khu
//     vuc) -> DongVaChoXong (cho tho ghi ghi HET xuong dia, co han gio qua
//     ctx) -> RỒI MỚI DungTatCa() (dung actor tung khu vuc). Xem tatmaychu.go
//     de biet chi tiet vi sao thu tu nay BAT BUOC.
//
//  4. kho.Close() — dong pool CSDL, CHI SAU KHI tho ghi da ghi xong va actor
//     da dung (ca hai nam trong buoc 3) — khong con ai dung toi kho nua luc nay.
//
// VE VI TRI CUA lichBoss.Dung()/loa.Dung() (buoc 1-2) — day la mot lua chon
// can nhac, khong phai tien tay nhet vao cho trong. Ba ly do, xep theo do
// nang (viet cho loa, ap dung dung y het cho lichBoss vi ly do da giai thich
// o buoc 1):
//
// (a) BAO VE LAN LUU CUOI CUNG — ly do NANG NHAT. LuuTatCa (dau buoc 3) hoi
// tung khu vuc bang AnhChupNguoiChoi: di qua hop thu DIEU KHIEN va chi cho
// toi thoiGianChoControl = 100ms, qua han la tra ve RONG, tuc CA KHU VUC DO
// KHONG DUOC LUU. Ma truoc khi tra loi anh chup, actor con chay
// xuLyGameplayDaXepTruoc() — no XA HET nhung su kien gameplay dang xep san.
// Loa phat qua DUNG cai hop thu gameplay do (SoKhuVuc.PhatToanServer ->
// k.Gui(suKienPhatGoi)), va moi su kien phat la mot vong ghi socket toi TAT
// CA nguoi trong khu vuc. Neu loa con song luc LuuTatCa chay, mot con bao
// thong bao (10 boss chet cung luc, mot vong lap loi, mot admin go nham)
// hoan toan co the day cai hang xa do vuot 100ms -> anh chup rong -> mat
// tien do CUA CA MOT KHU VUC NGUOI CHOI. Tat loa truoc la bit han duong do.
//
// (b) KHONG PHAT CHO NGUOI DA DI. Neu de loa (hoac lichBoss) song toi sau
// DungTatCa (cuoi buoc 3), no van cham chi bom su kien vao cac khu vuc dang
// bi dong — khong panic (k.Gui kiem k.xong roi bo goi, khuvuc.go), nhung la
// cong vo ich va lam nhieu chan doan. Rang buoc "loa.Dung() TRUOC DungTatCa"
// nay da duoc ghi san tu truoc o hai cho: GanLoa (sokhuvuc.go) va
// BoPhatThanh.Dung (thongbao.go) — buoc 2 chi la cho THI HANH no. Voi
// lichBoss thi buoc 1 dong vai tro tuong duong.
//
// (c) RE, NEN KHONG LAM CHAM VIEC TAT MAY. Dung() chi dong mot kenh roi cho
// vongLap thoat; no CO Y KHONG xa not hang doi (thongbao.go: xa not la ton
// trong nhip 1,5 giay cho TUNG cai, hang doi day 64 se giu tien trinh lai
// toi 96 GIAY trong khi ca han gio tat may chu chi 30). Buoc nay ton nhieu
// nhat la mot lan phat dang do — mili-giay.
//
// CAI PHAI DANH DOI, noi thang: thong bao SINH RA trong luc dang tat may (vd
// mot con boss chet dung khoanh khac do) se nam lai trong hang doi va KHONG
// BAO GIO len song. Chap nhan, vi hai loai du lieu nay khong cung gia: mat
// mot dong chu chay ngang man hinh la nguoi choi lo mot tin; mat mot ban ghi
// nhan vat la mat cap do/tien that cua ho, khong lay lai duoc. Khi phai
// chon, luon chon giu ban ghi — dung tinh than "hang doi day thi BO" cua
// Xep() va "khong xa not" cua Dung().
//
// VI SAO SAU boLuu.Dung() chu khong phai TRUOC no: boLuu la manh DUY NHAT
// trong day nay co the gay MAT DU LIEU AM THAM neu bi tri hoan (xem buoc 0),
// nen no phai la thu dau tien bi chan — khong chen bat cu gi len truoc, ke ca
// mot buoc re. Doi lai, loa/lichBoss song them vai mili-giay thi khong ai mat
// gi: lan luu dinh ky ma boLuu dang lam do neu co hong thi buoc 3 luu lai het.
//
// Loi tu TatMayChu (thuong la ctx.Err() khi vuot han gio) duoc TRA VE cho
// noi goi (main()) log canh bao — donDepKhiTat KHONG dung giua chung vi loi
// do: DungTatCa() (nam trong buoc 3) va kho.Close() (buoc 4) VAN PHAI chay
// du TatMayChu co loi hay khong, vi treo may chu vo han con te hon mot lan
// tat "danh hong".
//
// VE "khong duong nao goi Xep sau khi DongVaChoXong tra ve, ke ca KhiDong()
// cua phien dut luc do" (canh bao rieng cua brief Task 6): mot phien client
// dut ket noi trong LUC soKhuVuc.DungTatCa() dang chay (phan cuoi buoc 3,
// vong lap dung TUNG khu vuc mot) se goi KhiDong() (internal/game/login.go)
// -> AnhChupNguoiChoi. VOI cac khu vuc DA dung, AnhChupNguoiChoi tra rong
// ngay (guiTin cua actor khong con ai doc) nen KhiDong tu bo qua Xep (guard
// "len(ban) > 0"). VOI cac khu vuc CHUA kip dung toi (con lai trong vong
// lap DungTatCa), AnhChupNguoiChoi VAN tra du lieu that -> KhiDong CO THE
// goi Xep() SAU KHI thoGhi da dung — day la MOT KHE HO RIENG BIET, hep, con
// lai (KHONG phai loai da chan o buoc 0): DungTatCa() lap qua tung khu vuc
// CHI dong mot kenh + cho goroutine actor (khong co I/O nao) thoat, nen tong thoi gian ca vong lap thuong la
// mili-giay — xac suat mot client rot dung khoanh khac do la cuc nho, tuong
// duong cac gioi han "best-effort" khac cua graceful shutdown thong thuong
// (khong co co che "rut het session dang mo" — se can them ha tang theo doi
// TOAN BO session dang song, ngoai pham vi Task 6). Ghi lai o day de nguoi
// sau khong tuong nham day la khe ho CHUA duoc xem xet.
func donDepKhiTat(ctx context.Context, kho khoDong, soKhuVuc *game.SoKhuVuc, thoGhi *game.ThoGhi, boLuu *game.BoLuuDinhKy, lichBoss lichBossCoTheDung, loa loaCoTheDung) error {
	boLuu.Dung()
	lichBoss.Dung()
	loa.Dung()
	err := game.TatMayChu(ctx, soKhuVuc, thoGhi)
	kho.Close()
	return err
}

// tatMayChuChung la mot cua duy nhat cho moi nguon dung: signal, E2E,
// listener game/control loi va control khong mo duoc. Ham nhan callback don
// dep de test noi duoc control that toi donDepKhiTat that ma khong can chay
// main/os.Exit.
func tatMayChuChung(
	control mayChuCoTheTat,
	donDep func(context.Context) error,
) error {
	var loiControl error
	if control != nil {
		ctxControl, huyControl := context.WithTimeout(context.Background(), 2*time.Second)
		loiControl = control.Shutdown(ctxControl)
		huyControl()
	}
	ctxTat, huyTat := context.WithTimeout(context.Background(), hanGioTatMayChu)
	loiDonDep := donDep(ctxTat)
	huyTat()
	return errors.Join(loiControl, loiDonDep)
}

func choVaTatMayChu(
	sigCh <-chan os.Signal,
	yeuCauTatE2E <-chan struct{},
	loiMayChu <-chan error,
	loiControl <-chan error,
	control mayChuCoTheTat,
	donDep func(context.Context) error,
) (maThoat int, loiNguon error, loiDonDep error) {
	select {
	case sig := <-sigCh:
		slog.Info("da nhan tin hieu dung, dang dong may chu...", "tin_hieu", sig.String())
	case <-yeuCauTatE2E:
		slog.Info("da nhan yeu cau shutdown E2E, dang dong may chu an toan")
	case err := <-loiMayChu:
		maThoat = 1
		loiNguon = err
	case err := <-loiControl:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			maThoat = 1
			loiNguon = err
		}
	}
	loiDonDep = tatMayChuChung(control, donDep)
	return maThoat, loiNguon, loiDonDep
}

func main() {
	// MAC DINH LA INFO, TUYET DOI KHONG PHAI DEBUG.
	//
	// Ly do: da do duoc tren may chu that mot ket noi don gui goi 3 byte lenh
	// la co the dat 32.452 goi/giay. Truoc day moi goi nhu vay lam
	// internal/game/login.go ghi MOT dong Debug (~67 byte) - voi muc mac dinh
	// la Debug, dieu do khuech dai thanh 9 GB/GIO tren dia CHI TU MOT ket noi,
	// khong gioi han. Da bo han cac dong Debug do (xem login.go), nhung muc
	// mac dinh cua chinh logger cung phai la info: neu sau nay co ai them log
	// Debug moi ma quen doi lai muc mac dinh, may chu it nhat khong tu dong
	// bat no len khi trien khai binh thuong (khong dat NRO_MUC_NHAT_KY).
	// Muon xem log chi tiet luc dieu tra su co thi tu dat
	// NRO_MUC_NHAT_KY=debug, KHONG duoc lam do la mac dinh.
	muc, hopLe := mucNhatKyTuChuoi(os.Getenv("NRO_MUC_NHAT_KY"))
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: muc,
	})))
	if !hopLe {
		slog.Warn("NRO_MUC_NHAT_KY khong nhan dang duoc, dung muc mac dinh info",
			"gia_tri_nhan_duoc", os.Getenv("NRO_MUC_NHAT_KY"))
	}

	cauHinhControl, err := docCauHinhE2EControl(os.Getenv)
	if err != nil {
		slog.Error("KHONG THE KHOI DONG E2E control.", "loi", err)
		os.Exit(1)
	}

	dsn := os.Getenv("NRO_DSN")
	if dsn == "" {
		slog.Error("KHONG THE KHOI DONG: thieu bien moi truong NRO_DSN (chuoi ket noi PostgreSQL).")
		slog.Error("Cach sua: dat bien moi truong NRO_DSN roi chay lai. Vi du (PowerShell):")
		slog.Error(`  $env:NRO_DSN = "postgres://postgres:matkhau@127.0.0.1:5433/nro"`)
		slog.Error("Hoac (cmd.exe):")
		slog.Error(`  set NRO_DSN=postgres://postgres:matkhau@127.0.0.1:5433/nro`)
		os.Exit(1)
	}

	diaChi := os.Getenv("NRO_DIA_CHI")
	if diaChi == "" {
		diaChi = ":14445"
		slog.Info("chua dat NRO_DIA_CHI, dung dia chi mac dinh", "dia_chi", diaChi)
	}

	slog.Info("dang ket noi CSDL...")
	ctx := context.Background()
	kho, err := store.Open(ctx, dsn)
	if err != nil {
		slog.Error("KHONG THE KHOI DONG: khong noi duoc CSDL.", "loi", err)
		slog.Error("Cach kiem tra: PostgreSQL da chay chua? Dung dia chi/cong/mat khau chua?")
		// Khong bao gio dua DSN vao log: no thuong chua mat khau ro.
		slog.Error("Hay kiem tra lai user/mat khau/host/cong/ten CSDL trong NRO_DSN (gia tri duoc an).")
		os.Exit(1)
	}
	slog.Info("da noi CSDL thanh cong.")

	// No ha tang C: tu ap migration CON THIEU truoc khi mo cong. Fail-closed —
	// schema sai thi DUNG, khong mo cong voi CSDL nua voi. Han gio 60s de migration
	// treo (khoa bang, hoac cho advisory lock cua server khac) khong lam khoi dong
	// ket vo han — thoat co thong bao thay vi treo im.
	ctxMig, huyMig := context.WithTimeout(ctx, 60*time.Second)
	soAp, errMig := kho.ApDungMigrationConLai(ctxMig)
	huyMig()
	if errMig != nil {
		slog.Error("KHONG THE KHOI DONG: ap migration CSDL that bai.", "loi", errMig)
		slog.Error("Cach kiem tra: xem migration nao loi o log tren; sua SQL roi chay lai. May chu KHONG mo cong voi schema sai.")
		kho.Close()
		os.Exit(1)
	}
	if soAp > 0 {
		slog.Info("da ap migration CSDL luc khoi dong", "so_ap_moi", soAp)
	}

	// Chan tin hieu dung TRUOC khi vao vong lap lang nghe, de dam bao khong
	// bao gio bo lot Ctrl+C dung luc dang khoi dong.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// Ke hoach 5: tao 5 khu vuc actor (mapId 1-5) MOT lan, dung chung cho moi phien.
	soKhuVuc := game.MoSoKhuVuc()

	// KH6b Task 14: bat bo phat thanh toan may chu (goi 93) NGAY SAU khi co 5
	// khu vuc va TRUOC khi mo cong nghe — de khong co khoang thoi gian nao ma
	// mot khu vuc da song nhung chua co loa (thong bao sinh ra trong khoang do
	// se bien mat khong dau vet, xem noiBoPhatThanh).
	loa := noiBoPhatThanh(soKhuVuc)

	// KH6b Task 14 (phan hai): bat bo dieu phoi lich boss NGAY DUOI dong tren
	// — phai SAU noiBoPhatThanh vi bo dieu phoi la NGUON sinh boss (boss xuat
	// hien la sinh thong bao goi 93), bat no truoc khi co loa thi nhung dong
	// dau tien roi vao hu khong (chi tiet day du: xem noiLichBoss va khoi
	// chu thich "VE VI TRI CUA lichBoss.Dung()/loa.Dung()" trong donDepKhiTat
	// ben duoi — TAT dat NGAY TRUOC loa.Dung() trong ham do, cung ly do doi
	// xung voi thu tu BAT o day).
	lichBoss := noiLichBoss(soKhuVuc)

	// Ke hoach 11 Task 6 diem B "NOI DAY THAT": TRUOC dong nay, khong noi nao
	// trong ca du an goi MoThoGhi / MoBoLuuDinhKy(...).BatDau() / DatThoGhi —
	// toan bo ha tang luu tru cua Task 1-5 (da qua review) NAM DO KHONG AI
	// GOI TOI, may chu THAT SU KHONG LUU GI CA. MoLuuTru (internal/game) lap
	// rap ca ba manh mot lan duy nhat: mot ThoGhi CHUNG, mot BoLuuDinhKy DA
	// BATDAU (tu dong luu moi 60 giay — dac ta §9 muc 1), va taoXuLy tra ve
	// MOT BoXuLy MOI da DatThoGhi(thoGhi) cho MOI ket noi (dung ben duoi thay
	// game.NewBoXuLy truc tiep).
	//
	// No bao mat §7: MoLuuTru cung tao va noi day gac CHAN DO MAT KHAU (mot
	// cai chung cho ca may chu). Hai so cau hinh cua no doc tu bien moi
	// truong — xem docSoNguyenMoiTruong ben duoi de biet vi sao can no long
	// duoc ma khong phai bien dich lai.
	sucChua := docSoNguyenMoiTruong("NRO_CHONGDO_SUCCHUA", chanip.SucChuaMacDinh)
	nhipHoiPhuc := time.Duration(docSoNguyenMoiTruong(
		"NRO_CHONGDO_HOIPHUC_GIAY", int(chanip.NhipHoiPhucMacDinh/time.Second))) * time.Second
	thoGhi, boLuu, taoXuLy := game.MoLuuTru(kho, soKhuVuc, soThoGhiMacDinh, sucChua, nhipHoiPhuc)

	// No bao mat §7 (lo hong 3): gioi han so ket noi cung luc tren MOT dia chi
	// IP. Khong co no thi mot may duy nhat mo bao nhieu ket noi tuy thich va
	// an het bo nho may chu.
	toiDaKetNoi := docSoNguyenMoiTruong("NRO_TOIDA_KETNOI_MOI_IP", chanip.ToiDaKetNoiMoiIPMacDinh)
	gacKetNoi := chanip.MoGacKetNoi(toiDaKetNoi)
	slog.Info("gac theo dia chi mang da bat",
		"chong_do_suc_chua", sucChua,
		"chong_do_hoi_phuc", nhipHoiPhuc,
		"toi_da_ket_noi_moi_ip", toiDaKetNoi)

	// LangNghe chay vo han (Accept trong vong lap), nen phai chay tren
	// goroutine rieng de goroutine chinh con ranh tay cho tin hieu dung.
	loiMayChu := make(chan error, 1)
	go func() {
		loiMayChu <- mynet.LangNghe(diaChi, gacKetNoi, func() mynet.XuLyGoiTin {
			return taoXuLy()
		})
	}()

	yeuCauTatE2E := make(chan struct{}, 1)
	donDep := func(ctx context.Context) error {
		return donDepKhiTat(ctx, kho, soKhuVuc, thoGhi, boLuu, lichBoss, loa)
	}
	control, controlListener, err := moE2EControl(cauHinhControl, yeuCauTatE2E)
	if err != nil {
		if errTat := tatMayChuChung(nil, donDep); errTat != nil {
			slog.Error("don dep sau loi mo E2E control gap loi", "loi", errTat)
		}
		slog.Error("MAY CHU DUNG: khong mo duoc E2E control.", "loi", err)
		os.Exit(1)
	}
	loiControl := make(chan error, 1)
	if control != nil {
		go func() {
			loiControl <- control.Serve(controlListener)
		}()
		slog.Info("E2E control dang nghe tren loopback", "dia_chi", controlListener.Addr().String())
	}

	slog.Info("may chu san sang. Nhan Ctrl+C de dung may chu an toan.")

	maThoat, loiKhoiDong, errTat := choVaTatMayChu(
		sigCh,
		yeuCauTatE2E,
		loiMayChu,
		loiControl,
		control,
		donDep,
	)
	if errTat != nil {
		slog.Error("don dep luc tat may chu gap loi (co the mat vai thay doi GAN NHAT chua kip ghi)",
			"loi", errTat)
	}
	if loiKhoiDong != nil {
		slog.Error("MAY CHU DUNG DO LOI: listener da dung bat thuong.", "loi", loiKhoiDong)
	}
	slog.Info("da dong CSDL. Thoat may chu.", "ma_thoat", maThoat)
	os.Exit(maThoat)
}
