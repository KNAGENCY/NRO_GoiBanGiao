package game

import (
	"log/slog"

	"nroserver/internal/store"
)

// hinhdang.go — NGUOI CHOI TRONG NHU THE NAO tren client, va vi sao may chu
// phai quyet thay cho client.
//
// Boi canh (05/09/2026, phat hien luc choi thu that bang client Unity qua MCP):
//
//  1. Client Unity la ban port cua game Dragon Ball: no biet DUNG 3 "gender"
//     (Trai Dat / Namek / Xayda) va tu chon than/chan mac dinh theo gender
//     (Char.setDefaultBody/setDefaultLeg: 0 -> 57/58, 1 -> 59/60, 2 -> 57/58,
//     con lai -> KHONG DAT, giu 0). Game Naruto co 5 lang; may chu tu truoc ghi
//     "gender = id lang" (KH3) nen lang 1 ra bo do Namek cu (59/60 van con
//     trong part.bin), lang 3/4 ra part 0 — mot part DAU chi co 3 khung, ma
//     than can toi 17 khung -> client nem IndexOutOfRange khi ve.
//  2. Goi -5 (xuat hien) truoc day ghi than = 0, chan = 0 cho NGUOI CHOI (chu
//     thich cu: "chua co he trang bi nen outfit deu 0"). Client KHAC nhin thay
//     ta bang goi nay -> cung dinh part 0 nhu tren.
//  3. Man tao nhan vat cua client gui id TOC Dragon Ball (CreateCharScr.
//     hairIDOld: 64/30/31, 9/29/32, 6/27/28) lam KieuToc; may chu luu nguyen.
//     Bo hinh Naruto hom nay chi co MOT dau (part 30). Chon toc khac 30 la
//     dau Dragon Ball nho gan len than Naruto to.
//  4. Nhan vat chet (HpHienTai = 0) duoc luu xuong CSDL (KH11 luu 3 thoi diem)
//     va dang nhap lai VAN chet: goi -42 kep mau ve day de HUD khong 0, nen
//     client tuong song, con actor thi coi la chet (khuvuc.go: nguoi chet
//     khong danh duoc) — hai ben noi hai su that.
//
// Bon ham duoi day la MOT cho quyet dinh cho ca bon chuyen, de -5 / ME_INFO /
// -37 / tao nhan vat khong bao gio noi khac nhau.

// Ba part cua bo do "hoc vien" — bo hinh nguoi choi DUY NHAT hien co
// (docs/my-thuat/mapping-id.json parts[], asset-lock.json closure.part_ids).
// HeadNhanVat (thegioi.go) la dau; hai hang nay la than va chan.
//
// Khi co bo do thu hai (KH7 ban 2 "mac do doi hinh"): than/chan lay tu
// MauVatPham.Phan cua mon dang mac, hai hang nay chi con la MAC DINH khi
// khong mac gi.
const (
	ThanNhanVat int16 = 57
	ChanNhanVat int16 = 58
)

// gioiTinhClient la byte "gender" ma client dung cho MAC DINH THAN/CHAN, mat
// nhap nhay, hieu ung theo hanh tinh (Char.cgender). Luon 0 = "Trai Dat":
// bo do Naruto duy nhat nam o cac slot mac dinh cua gender 0 (57/58), va cac
// bang tra theo gender ben client (GameScr.info1.charId[gender]) chi dai 3.
//
// KHONG NHAM voi byte "lop/lang" (chi so nClasss[], quyet ten lang + bo chieu
// trong catalog -28/7): byte do VAN la id lang 0..4. ME_INFO va -5 ghi CA HAI
// byte — mot cai la gioi tinh (ham nay), mot cai la lang (nv.Lang).
func gioiTinhClient(lang int16) byte {
	_ = lang // giu tham so de cho goi doc ro "tu lang nao ra", du hom nay luon 0
	return 0
}

// dauHopLe tra part dau se VE cho KieuToc da luu. Bo hinh Naruto hom nay chi
// co dau 30 (HeadNhanVat); moi gia tri khac (toc Dragon Ball tu man tao nhan
// vat, hang cu trong CSDL, gia tri rac) deu quy ve dau do thay vi ve mot cai
// dau khac bo. Khi co them dau, mo rong bang o day — MOT cho.
func dauHopLe(kieuToc int16) int16 {
	if kieuToc == HeadNhanVat {
		return kieuToc
	}
	return HeadNhanVat
}

// hinhDangNguoiChoi dung HinhDangHien cho goi -5 cua MOT nguoi choi: dau qua
// dauHopLe, than/chan la bo mac dinh (chua co he mac do doi hinh), cac o
// flagBag/aura/effFront = 0 nhu KH5. goiXuatHien va moi cho khac muon "nguoi
// choi nay trong the nao" PHAI di qua day.
func hinhDangNguoiChoi(nv *store.NhanVat) HinhDangHien {
	return HinhDangHien{
		Lang:   nv.Lang,
		Outfit: [6]int16{dauHopLe(nv.KieuToc), ThanNhanVat, ChanNhanVat, 0, 0, 0},
		TypePk: 0,
	}
}

// hoiSinhLucVao: nhan vat dang CHET trong CSDL (HpHienTai <= 0) thi song lai
// ngay luc vao the gioi — day mau va chakra, GIU vi tri — y het case
// suKienHoiSinh (khuvuc.go) lam khi nguoi choi bam "hoi sinh", chi khac la
// khong can client xin. Tra true neu co doi.
//
// Goi tren LUONG PHIEN, TRUOC khi dung ME_INFO va TRUOC khi -39 giao nv cho
// actor (xuLyClientOk, login.go) — luc do chi luong phien cham nv, khong co
// tranh chap. Khong ghi CSDL o day: ThoGhi (KH11) se luu o nhip 60s / khi roi.
func hoiSinhLucVao(nv *store.NhanVat) bool {
	if nv == nil || nv.HpHienTai > 0 {
		return false
	}
	nv.HpHienTai = nv.HpToiDa
	nv.ChakraHienTai = nv.ChakraToiDa
	slog.Info("nhan vat chet trong CSDL, hoi sinh luc vao the gioi",
		"ten", nv.Ten, "id", nv.ID, "hp", nv.HpHienTai, "chakra", nv.ChakraHienTai)
	return true
}
