// Package antoan giu luoi do (recover) dung chung cho MOI goroutine cua may
// chu. Dat rieng thay vi nhet vao game/net vi CA HAI package do deu can luoi:
// neu dat o game thi net phai import game (vong import). antoan CHI phu thuoc
// stdlib (log/slog, runtime/debug) nen ai import cung duoc, khong the sinh vong.
package antoan

import (
	"log/slog"
	"runtime/debug"
)

// ChayCoLuoi chay f; neu f panic thi CHAN lai, ghi log kem dau vet (stack) roi
// tra ve binh thuong — KHONG cho panic lan len lam sap ca tien trinh, da van
// TOAN BO nguoi choi. Tra ve true neu da chan mot panic (de caller tu quyet
// dinh lam gi sau do: dem, dung rieng khu vuc, dong ket noi...), false neu f
// chay tron.
//
// VI SAO khong nuot loi im lang: bay §6 (ban Java) da tung nuot loi trong
// catch rong khien su co khong the chan doan. Moi panic bi chan o day PHAI de
// lai log muc Error kem stack de nguoi van hanh thay va lan ra duoc goc.
//
// `cho` la nhan mo ta NOI dat luoi (vd "khuvuc.xuLy", "net.phucvu:1.2.3.4",
// "thogi.ghi") — doc log thay truong nay la biet panic phat ra tu goroutine nao.
func ChayCoLuoi(cho string, f func()) (daChanPanic bool) {
	defer func() {
		if r := recover(); r != nil {
			daChanPanic = true
			slog.Error("panic bi chan boi luoi",
				"cho", cho, "loi", r, "stack", string(debug.Stack()))
		}
	}()
	f()
	return false
}
