// Chuong trinh nroserver: diem khoi dong may chu game.
//
// Doc cau hinh tu bien moi truong, mo CSDL, roi lang nghe ket noi TCP cho
// toi khi nhan tin hieu dung (Ctrl+C) hoac loi khong the phuc hoi.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

// khoDong la phan toi thieu store.Store can de don dep luc tat may chu — CHI
// Close(), tach thanh interface HEP de test tiem duoc "kho gia" (khong can
// PostgreSQL that), dung khuon voi Kho/KhoLuu (internal/game). *store.Store
// thoa man interface nay tu nhien.
type khoDong interface {
	Close()
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
//     Task 6 nay — nhung PHAI dung no TRUOC buoc 1, vi neu khong, ticker 60
//     giay cua no co the tu goi LuuTatCa->Xep() TREN MOT GOROUTINE DOC LAP,
//     KHONG dong bo gi voi luong tat may chu chinh — neu lan Xep do roi
//     dung luc SAU khi thoGhi.DongVaChoXong() (ben trong buoc 1 duoi day) da
//     tra ve, ban ghi do se nam lai trong hang doi VINH VIEN, khong ai flush nua
//     (dung CHINH XAC canh bao cua brief Task 6: "gọi Xep() SAU khi
//     DongVaChoXong() trả về sẽ mất dữ liệu âm thầm... kể cả KhiDong() của
//     phiên đứt lúc đó" — boLuu la mot nguon Xep() DOC LAP giong het
//     KhiDong(), can chan truoc TIEN cho chac).
//  1. game.TatMayChu(ctx, soKhuVuc, thoGhi) — GOM CA BA buoc con lai cua dac
//     ta §4 diem C: LuuTatCa (xin anh chup TOAN BO nguoi choi o TAT CA khu
//     vuc) -> DongVaChoXong (cho tho ghi ghi HET xuong dia, co han gio qua
//     ctx) -> RỒI MỚI DungTatCa() (dung actor tung khu vuc). Xem tatmaychu.go
//     de biet chi tiet vi sao thu tu nay BAT BUOC.
//  2. kho.Close() — dong pool CSDL, CHI SAU KHI tho ghi da ghi xong va actor
//     da dung (ca hai nam trong buoc 1) — khong con ai dung toi kho nua luc nay.
//
// Loi tu TatMayChu (thuong la ctx.Err() khi vuot han gio) duoc TRA VE cho
// noi goi (main()) log canh bao — donDepKhiTat KHONG dung giua chung vi loi
// do: DungTatCa() (nam trong buoc 1) va kho.Close() (buoc 2) VAN PHAI chay
// du TatMayChu co loi hay khong, vi treo may chu vo han con te hon mot lan
// tat "danh hong".
//
// VE "khong duong nao goi Xep sau khi DongVaChoXong tra ve, ke ca KhiDong()
// cua phien dut luc do" (canh bao rieng cua brief Task 6): mot phien client
// dut ket noi trong LUC soKhuVuc.DungTatCa() dang chay (phan cuoi buoc 1,
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
func donDepKhiTat(ctx context.Context, kho khoDong, soKhuVuc *game.SoKhuVuc, thoGhi *game.ThoGhi, boLuu *game.BoLuuDinhKy) error {
	boLuu.Dung()
	err := game.TatMayChu(ctx, soKhuVuc, thoGhi)
	kho.Close()
	return err
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
		slog.Error("Chuoi ket noi dang dung (kiem tra lai user/mat khau/host/cong/ten CSDL):", "NRO_DSN", dsn)
		os.Exit(1)
	}
	slog.Info("da noi CSDL thanh cong.")

	// Chan tin hieu dung TRUOC khi vao vong lap lang nghe, de dam bao khong
	// bao gio bo lot Ctrl+C dung luc dang khoi dong.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Ke hoach 5: tao 5 khu vuc actor (mapId 1-5) MOT lan, dung chung cho moi phien.
	soKhuVuc := game.MoSoKhuVuc()

	// Ke hoach 11 Task 6 diem B "NOI DAY THAT": TRUOC dong nay, khong noi nao
	// trong ca du an goi MoThoGhi / MoBoLuuDinhKy(...).BatDau() / DatThoGhi —
	// toan bo ha tang luu tru cua Task 1-5 (da qua review) NAM DO KHONG AI
	// GOI TOI, may chu THAT SU KHONG LUU GI CA. MoLuuTru (internal/game) lap
	// rap ca ba manh mot lan duy nhat: mot ThoGhi CHUNG, mot BoLuuDinhKy DA
	// BATDAU (tu dong luu moi 60 giay — dac ta §9 muc 1), va taoXuLy tra ve
	// MOT BoXuLy MOI da DatThoGhi(thoGhi) cho MOI ket noi (dung ben duoi thay
	// game.NewBoXuLy truc tiep).
	thoGhi, boLuu, taoXuLy := game.MoLuuTru(kho, soKhuVuc, soThoGhiMacDinh)

	// LangNghe chay vo han (Accept trong vong lap), nen phai chay tren
	// goroutine rieng de goroutine chinh con ranh tay cho tin hieu dung.
	loiMayChu := make(chan error, 1)
	go func() {
		loiMayChu <- mynet.LangNghe(diaChi, func() mynet.XuLyGoiTin {
			return taoXuLy()
		})
	}()

	slog.Info("may chu san sang. Nhan Ctrl+C de dung may chu an toan.")

	select {
	case sig := <-sigCh:
		slog.Info("da nhan tin hieu dung, dang dong may chu...", "tin_hieu", sig.String())
		ctxTat, huy := context.WithTimeout(context.Background(), hanGioTatMayChu)
		if errTat := donDepKhiTat(ctxTat, kho, soKhuVuc, thoGhi, boLuu); errTat != nil {
			slog.Error("don dep luc tat may chu gap loi (co the mat vai thay doi GAN NHAT chua kip ghi)",
				"loi", errTat)
		}
		huy()
		slog.Info("da dong CSDL. Thoat may chu.")
		os.Exit(0)

	case err := <-loiMayChu:
		// LangNghe() chi tra ve khi khong the mo cong nghe (loi khoi dong):
		// than vong lap Accept cua no tu bat loi tung ket noi va khong bao
		// gio tu thoat, nen err o day luon khac nil. Chua ai ket noi duoc
		// (Accept chua bao gio chay) nen donDepKhiTat o day chu yeu la dong
		// bo luu dinh ky + pool CSDL cho sach, khong mat du lieu nguoi choi
		// nao (chua co ai de mat).
		ctxTat, huy := context.WithTimeout(context.Background(), hanGioTatMayChu)
		if errTat := donDepKhiTat(ctxTat, kho, soKhuVuc, thoGhi, boLuu); errTat != nil {
			slog.Error("don dep luc tat may chu gap loi", "loi", errTat)
		}
		huy()
		slog.Error("MAY CHU DUNG DO LOI: khong mo duoc cong de lang nghe ket noi.", "loi", err)
		os.Exit(1)
	}
}
