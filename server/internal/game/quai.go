// Package game — mo hinh quai (mob) va cach ghi danh sach quai vao goi -24 MAP
// INFO (Ke hoach 6 Task 5). Quai la trang thai THUOC actor KhuVuc: tao luc
// MoKhuVuc, chi luong vongLap cua khu vuc duoc doc/ghi (xem khuvuc.go, cung quy
// tac da ap dung cho nguoiTrong o Ke hoach 5).
package game

import (
	"time"

	"nroserver/internal/protocol"
)

// Quai la mot con quai dang song trong PHAM VI mot KhuVuc. Cac Task sau cua Ke
// hoach 6 (6-9, 13) se doc/ghi them cac truong nay tu trong vongLap (AI, chien
// dau, hoi sinh) — Task 5 chi dung mo hinh + phat trong -24.
type Quai struct {
	// Id la tempId — id TAM THOI trong PHAM VI khu vuc nay (khong phai id toan
	// cuc), dung de cac goi chien dau sau nay (54/-9/-10/-11/-12/-13, Task 6+)
	// tham chieu dung con quai. Client doc byte co dau -> Id/soQuai <= 127
	// (dac ta §1.4).
	Id byte

	// MauId la id mau quai (mob template) — KHONG xuat hien trong than goi -24
	// (dac ta §3.6 khong co truong nay), nhung can giu noi bo de cac Task sau
	// (AI, loot, hien thi ten/hinh) tra bang du lieu mau quai.
	MauId byte

	// Cap la truong "level" trong dac ta §3.6 (thu tu 9, dung SAU hp, TRUOC
	// hpMax) — cap dung de tinh sat thuong/song sot.
	Cap byte

	// CapTinhAnh la truong "lvMob" trong dac ta §3.6 (thu tu 14, o cuoi than
	// goi) — cap HIEN THI canh ten quai o client, TACH RIENG voi Cap vi dac
	// ta dinh nghia hai o rieng du thuong bang nhau.
	CapTinhAnh byte

	// Hp/HpToiDa la mau int32 4 byte big-endian (dac ta §1.1: readInt phia
	// client, kep tran 2_147_483_647, kep san 0).
	Hp      int32
	HpToiDa int32

	// X/Y la toa do quai trong ban do (short, dac ta §3.6 truong 11-12).
	X int16
	Y int16

	// TrangThai la truong "status" (dac ta §3.6 truong 13) — 0 = binh thuong;
	// Task sau (choang, buff...) se dung gia tri khac.
	TrangThai byte

	// --- Ba truong duoi day la Ke hoach 6 Task 8 (hoi sinh 5s + hoi mau 30s,
	// dac ta §5.2) — trang thai NOI BO cua actor, KHONG xuat hien trong than
	// goi -24 (ghiQuai o duoi khong dung toi). Chi luong vongLap cua KhuVuc
	// duoc doc/ghi (dac ta §1.6), giong moi truong khac cua Quai.

	// LucChet la thoi diem Hp xuong <= 0 GAN NHAT — nhipTim (khuvuc.go) dung
	// truong nay de tinh "da troi qua bao lau tu luc chet" cho nguong 5 giay
	// hoi sinh. Duoc gan trong case suKienDanhQuai (khuvuc.go) ngay khi phat
	// hien Hp <= 0, TRUOC khi phat goi -12. Zero value (chua tung chet) khong
	// gay hai: quai moi taoQuaiMacDinh luon Hp > 0 nen nhanh "dang chet" cua
	// nhipTim khong bao gio doc truong nay khi con zero.
	LucChet time.Time

	// LucHoiSinh la thoi diem quai HOI SINH GAN NHAT (nhipTim gan khi qua du
	// 5 giay ke tu LucChet). Dung lam ViTriCoMau.LucHoiSinh (chotkiem.go) khi
	// nguoi choi xin danh quai nay — kich hoat bat tu 1500ms sau hoi sinh (dac
	// ta §1.7), CUNG CO CHE voi nguoiTrongKhuVuc ma Task 13 se dung cho nguoi
	// choi. Zero value (quai moi spawn, CHUA TUNG hoi sinh lan nao — khac voi
	// "vua hoi sinh") khien bayGio.Sub(zero) ra khoang RAT LON -> kiem bat tu
	// tu nhien CHO PHEP ngay, dung hanh vi mong doi: quai moi mo ban do khong
	// bi coi la "vua hoi sinh" nen khong duoc mien nhiem oan.
	LucHoiSinh time.Time

	// LanHoiMauCuoi la thoi diem GAN NHAT nhipTim DA KIEM (khong nhat thiet
	// DA HOI — neu dang day mau thi kiem xong khong hoi gi) chu ky hoi mau 30
	// giay (dac ta §5.2). Khoi tao = thoi diem spawn (MoKhuVuc gan sau khi goi
	// taoQuaiMacDinh) de lan kiem DAU TIEN dung 30 giay sau khi ban do mo, chu
	// khong hoi ngay lap tuc do zero value nam 1 tao ra khoang cach "da troi
	// qua" khong lo. Sau moi lan hoi sinh (nhipTim), truong nay CUNG duoc dat
	// lai ve thoi diem hoi sinh — tranh mot con quai vua hoi (da day mau) van
	// bi tinh "no du 30s roi" tu truoc luc chet va hoi "thua" ngay tick ke tiep.
	LanHoiMauCuoi time.Time
}

// ghiQuai ghi MOT con quai vao than goi -24 MAP INFO, DUNG THEO THU TU dac ta
// §3.6. KH6 chua co quai bi vo hieu hoa/khong di chuyen duoc/nguyen to lua-bang-
// gio hay la BIG BOSS nen 5 co dau va isBigBoss LUON ghi false; sys LUON = 0.
//
// ⚠️ BAY LON NHAT (§3.6, ly do co test rieng TestGhiQuaiThuTuHpCapHpToiDa): thu
// tu la hp -> level -> hpMax, KHONG PHAI hp -> hpMax -> level nhu truc giac hay
// nghi. Doi thu tu hai truong nay van cho ra tong do dai dung (client khong bao
// loi) nhung gia tri doc duoc SAI HOAN TOAN — dung dac loai "sai thu tu, tong do
// dai van khop" o bang §2.0 cua dac ta: mau am, cap ra tuong hp, khong log nao
// bao ca.
func ghiQuai(w *protocol.Writer, q *Quai) {
	w.Bool(false) // 1 isDisable
	w.Bool(false) // 2 isDontMove
	w.Bool(false) // 3 isFire
	w.Bool(false) // 4 isIce
	w.Bool(false) // 5 isWind

	w.Byte(q.Id) // 6 tempId
	w.Byte(0)    // 7 sys = 0

	w.Int(q.Hp)      // 8  hp    — TRUOC level, DUNG doi cho hpMax
	w.Byte(q.Cap)    // 9  level — GIUA hp va hpMax
	w.Int(q.HpToiDa) // 10 hpMax — SAU level

	w.Short(q.X) // 11 x
	w.Short(q.Y) // 12 y

	w.Byte(q.TrangThai)  // 13 status
	w.Byte(q.CapTinhAnh) // 14 lvMob

	w.Bool(false) // 15 isBigBoss
}

// ghiDanhSachQuai ghi soQuai (byte) roi tung con quai theo dac ta §3.6, dung
// trong goi -24 MAP INFO. soQuai la byte CO DAU o phia client (dac ta §1.4) nen
// PHAI <= 127; qua han se bi CAT BOT (khong ghi qua 127 con, khong ghi so am)
// de khong lam hong khung goi — day la loi cau hinh ban do (qua nhieu quai),
// khong phai loi runtime nen khong panic, chi cat cho an toan.
func ghiDanhSachQuai(w *protocol.Writer, ds []*Quai) {
	n := len(ds)
	if n > 127 {
		n = 127
	}
	w.Byte(byte(n))
	for i := 0; i < n; i++ {
		ghiQuai(w, ds[i])
	}
}

// taoQuaiMacDinh sinh danh sach quai KHOI TAO cho mot khu vuc moi mo (goi tu
// MoKhuVuc). CHUA CO catalog quai theo lang/theo map (ngoai pham vi Task 5 —
// khong task nao trong Ke hoach 6 giao viec do), nen tam thoi moi map mo ra
// voi 2 con quai co so lieu CO DINH, quanh diem SpawnX/SpawnY cua ban do, de
// Task 6-9/13/14 (chien dau, AI, hoi sinh, test dau-cuoi) co doi tuong that
// de thao tac. Danh dau ro dat de coordinator biet can thay bang catalog that
// khi thiet ke noi dung (spawn theo lang, do kho, vi tri ban do that) o mot
// ke hoach sau.
func taoQuaiMacDinh(bd *BanDo) []*Quai {
	if bd == nil {
		return nil
	}
	return []*Quai{
		{
			Id: 0, MauId: 1, Cap: 1, CapTinhAnh: 1,
			Hp: 100, HpToiDa: 100,
			X: bd.SpawnX + 60, Y: bd.SpawnY,
		},
		{
			Id: 1, MauId: 1, Cap: 1, CapTinhAnh: 1,
			Hp: 100, HpToiDa: 100,
			X: bd.SpawnX - 60, Y: bd.SpawnY,
		},
	}
}
