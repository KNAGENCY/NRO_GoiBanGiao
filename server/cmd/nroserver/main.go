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

	"nroserver/internal/game"
	mynet "nroserver/internal/net"
	"nroserver/internal/store"
)

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

	// LangNghe chay vo han (Accept trong vong lap), nen phai chay tren
	// goroutine rieng de goroutine chinh con ranh tay cho tin hieu dung.
	loiMayChu := make(chan error, 1)
	go func() {
		loiMayChu <- mynet.LangNghe(diaChi, func() mynet.XuLyGoiTin {
			return game.NewBoXuLy(kho)
		})
	}()

	slog.Info("may chu san sang. Nhan Ctrl+C de dung may chu an toan.")

	select {
	case sig := <-sigCh:
		slog.Info("da nhan tin hieu dung, dang dong may chu...", "tin_hieu", sig.String())
		kho.Close()
		slog.Info("da dong CSDL. Thoat may chu.")
		os.Exit(0)

	case err := <-loiMayChu:
		// LangNghe() chi tra ve khi khong the mo cong nghe (loi khoi dong):
		// than vong lap Accept cua no tu bat loi tung ket noi va khong bao
		// gio tu thoat, nen err o day luon khac nil.
		kho.Close()
		slog.Error("MAY CHU DUNG DO LOI: khong mo duoc cong de lang nghe ket noi.", "loi", err)
		os.Exit(1)
	}
}
