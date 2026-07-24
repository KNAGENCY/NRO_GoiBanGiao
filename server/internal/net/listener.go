package net

import (
	"fmt"
	"log/slog"
	"net"
)

// LangNghe mo cong TCP va phuc vu tung ket noi tren mot goroutine rieng.
// taoXuLy duoc goi mot lan cho MOI ket noi (xem PhucVu).
func LangNghe(diaChi string, taoXuLy func() XuLyGoiTin) error {
	ln, err := net.Listen("tcp", diaChi)
	if err != nil {
		return fmt.Errorf("khong mo duoc cong %s: %w", diaChi, err)
	}
	defer ln.Close()

	slog.Info("may chu game dang lang nghe", "dia_chi", diaChi)

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("nhan ket noi that bai", "loi", err)
			continue
		}
		go PhucVu(conn, taoXuLy)
	}
}
