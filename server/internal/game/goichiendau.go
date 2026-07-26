// Package game — dung goi XUONG cua vong chien dau nguoi-vs-quai (Ke hoach 6
// Task 6) VA nguoi-vs-nguoi (Task 10).
//
// Task 6 (giu nguyen, KHONG sua trong Task 10): 54 (nguoi danh quai), -9
// (quai con song), -12 (quai chet), -13 (quai hoi sinh). Ca 4 gui CHO CA KHU
// VUC (dac ta §3.3).
//
// Task 10 them: -60 (nguoi danh nguoi, Nhom F, dac ta §2.1/§3.3) — CO HAI
// DANG DO DAI tuy truong docTiep, dung HAI HAM TACH RIENG goiDanhNguoiChoi
// (dang don danh, 18 byte) va goiHoiMauNguoiChoi (dang hoi mau, 11 byte). Xem
// chu thich hai ham do de biet ly do KHONG duoc gop thanh mot ham co co.
//
// Ca file chi dung goi (byte-exact), KHONG tinh sat thuong/AI/kiem tra dieu
// kien — do la viec cua cac Task/file khac (satthuong.go, chotkiem.go, quai.go).
package game

import "nroserver/internal/protocol"

// Nam ma lenh dung trong file nay: bon ma nhom "nguoi danh quai" (Task 6, dac
// ta §3.3) va mot ma nhom "nguoi danh nguoi" (Task 10, dac ta §2.1/§3.3).
const (
	lenhNguoiDanhQuai  int8 = 54  // Nguoi danh quai — CHU Y: khong co truong dem
	lenhQuaiConSong    int8 = -9  // Quai con song sau don danh
	lenhQuaiChet       int8 = -12 // Quai chet, kem danh sach do roi (KH6 luon rong)
	lenhQuaiHoiSinh    int8 = -13 // Quai hoi sinh
	lenhNguoiDanhNguoi int8 = -60 // Nguoi danh nguoi — CO HAI DANG DO DAI, xem goiDanhNguoiChoi/goiHoiMauNguoiChoi
)

// DoRoi la mot vat pham roi ra khi quai chet (dac ta §3.3, truong lap trong
// goi -12). KH6 KHONG lam he vat pham (xem §0 "Ngoai KH6" cua dac ta) nen
// truong nay chi dinh nghia san cho Task sau (KH7) dung — Task 6 chi goi
// goiQuaiChet voi doRoi rong (soDoRoi = 0).
type DoRoi struct {
	ItemMapId   int16
	TemplateId  int16
	X           int16
	Y           int16
	IdNguoiNhat int32
}

// goiNguoiDanhQuai dung goi XUONG 54: bao ca khu vuc mot nguoi vua danh mot
// con quai. ⚠️ Than gui CHI CO 3 TRUONG (idNguoiDanh, skillId, mobId) — dac ta
// §2.0 noi ro gui 54 KHONG co truong dem, client doc mobId cho toi HET THAN
// GOI (toi da 10 lan). Ghi THEM BAT KY BYTE NAO sau mobId se lam client hieu
// nham thanh "nguoi nay danh 2 con quai" -> quay sai huong, nham sai muc tieu.
// Vi vay ham nay KHONG nhan tham so danh sach quai — chi 1 mobId duy nhat,
// dung dac ta ma khong the vo tinh ghi thua.
func goiNguoiDanhQuai(idNguoiDanh int32, skillId byte, mobId byte) *protocol.Message {
	w := protocol.NewWriter()
	w.Int(idNguoiDanh)
	w.Byte(skillId)
	w.Byte(mobId)
	// KHONG ghi gi them sau day — xem chu thich ham va dac ta §2.0/§3.3.
	return protocol.NewMessage(lenhNguoiDanhQuai, w)
}

// goiQuaiConSong dung goi XUONG -9: quai bi danh nhung con song, bao ca khu
// vuc mau con lai + sat thuong vua nhan.
//
// ⚠️ Ba diem da xac minh bang client (dac ta §3.3), DUNG CHEP theo ban Java:
//  1. Truong cuoi hieuUngId la BYTE (1 byte), khong phai int (4 byte) nhu
//     Mob.java:1369 (writeInt(-1)). Client (Controller.cs:4352) chi doc 1
//     byte — ghi int se day 3 byte thua vao truong SAU (khong con truong nao
//     trong -9 nhung se lam sai KHI co ai them truong moi vao cuoi sau nay,
//     va la sai kieu du lieu ro rang can tranh).
//  2. Goi nguoi goi PHAI tranh satThuong == 1 (client an toan bo qua ve khi
//     dame == 1, Controller.cs:4338) — trach nhiem cua ham tinh sat thuong
//     (Task khac), KHONG kiem tra o day.
//  3. satThuong == 0 la "truot" — CHỈ hien "miss" dung khi chiMang == false
//     (Controller.cs:4356-4363). Goi ham nay PHAI truyen chiMang=false khi
//     truot (rang buoc bo sung §1.9/§3.3).
func goiQuaiConSong(mobId byte, hpConLai int32, satThuong int32, chiMang bool, hieuUngId byte) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(mobId)
	w.Int(hpConLai)
	w.Int(satThuong)
	w.Bool(chiMang)
	w.Byte(hieuUngId) // 1 BYTE — KHONG phai int, xem chu thich ham
	return protocol.NewMessage(lenhQuaiConSong, w)
}

// goiQuaiChet dung goi XUONG -12: quai chet, bao ca khu vuc sat thuong don cuoi
// + danh sach do roi. KH6 chua lam he vat pham (dac ta §0 "Ngoai KH6") nen goi
// tu Task 6 luon truyen doRoi rong/nil -> soDoRoi = 0, KHONG lap qua vong nao.
//
// Dung xong TOAN BO du lieu (bao gom cat len 127 do roi) roi moi ghi byte —
// rang buoc §1.10 cua dac ta: loi giua chung phai bo ca goi, khong ghi goi
// thieu byte. Vi ham nay thuan tuy lap qua slice co san (khong the loi giua
// chung tru viec cat do dai), khong can co che rollback rieng — nhung van cat
// so luong TRUOC khi bat dau vong ghi de dam bao khong bao gio ghi qua han.
func goiQuaiChet(mobId byte, satThuong int32, chiMang bool, doRoi []DoRoi) *protocol.Message {
	n := len(doRoi)
	if n > 127 {
		n = 127 // soDoRoi la byte co dau phia client (dac ta §1.4) -> <= 127
	}

	w := protocol.NewWriter()
	w.Byte(mobId)
	w.Int(satThuong)
	w.Bool(chiMang)
	w.Byte(byte(n))
	for i := 0; i < n; i++ {
		dr := doRoi[i]
		w.Short(dr.ItemMapId)
		w.Short(dr.TemplateId)
		w.Short(dr.X)
		w.Short(dr.Y)
		w.Int(dr.IdNguoiNhat)
	}
	return protocol.NewMessage(lenhQuaiChet, w)
}

// goiQuaiHoiSinh dung goi XUONG -13: quai vua hoi sinh, bao ca khu vuc mau day
// du. KHONG gui goi 87 theo sau — dac ta §0.1 giai thich ky: ma lenh 87 phia
// client la "quai bi gio thoi" (mobId + bool), KHONG PHAI cap nhat maxHp nhu
// ban Java tuong nham; gui 87 sau -13 con lam quai DUNG IM neu hp >= 16.777.216.
//
// ⚠️ sys va levelBoss LUON ghi 0 — dac ta §3.3 xac minh client ghi hai byte
// nay vao mob.sys/mob.levelBoss chu KHONG doi loai/cap quai khi hoi sinh, va
// levelBoss != 0 con BAT hieu ung sieu cap NGAU NHIEN khong mong muon. Vi vay
// ham nay KHONG nhan tham so cho hai truong nay — tranh nguy co truyen nham.
func goiQuaiHoiSinh(mobId byte, hp int32) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(mobId)
	w.Byte(0) // sys — LUON 0, xem chu thich ham
	w.Byte(0) // levelBoss — LUON 0, xem chu thich ham
	w.Int(hp)
	return protocol.NewMessage(lenhQuaiHoiSinh, w)
}

// goiDanhNguoiChoi dung goi XUONG -60 DANG "DON DANH" (docTiep=1, 18 byte
// than), dac ta §2.1/§3.3, Ke hoach 6 Task 10: nguoi choi gay sat thuong len
// nguoi choi khac.
//
// ⚠️ HAM NAY VA goiHoiMauNguoiChoi BAT BUOC TACH HAN THANH HAI HAM, KHONG
// DUOC gop lai thanh MOT HAM CO THAM SO CO (vi du mot tham so "docTiep byte"
// hay "coDanhTiep bool" quyet dinh co ghi tiep 4 truong cuoi hay khong). Ly
// do (dac ta §2.1, nhac lai §2.0):
//
//	Neu gop thanh mot ham co co, nguoi sua code sau nay RAT DE vo tinh ghi ca
//	4 truong loaiChieu/satThuong/daChet/chiMang vao NHANH HOI MAU (vi du: sua
//	default cua tham so co, hoac quen nhanh if/else khi them truong moi). Day
//	la loi "ghi thua" cua bang §2.0 dac ta — tuc la IM LANG TUYET DOI: client
//	vut byte du ngay luc msg.cleanup(), KHONG loi, KHONG trieu chung, chay thu
//	trong game KHONG PHAT HIEN DUOC — chi lo ra khi co nguoi ngoi so byte tay.
//	Tach thanh hai ham doc lap (khong dung chung mot ham noi bo co co) thi loi
//	nay KHONG THE TON TAI: goiHoiMauNguoiChoi don gian KHONG CO cau lenh nao
//	de ghi 4 truong do, du nguoi sua sau co can trong hay khong.
//
// soMucTieu va loaiChieu la HANG SO trong pham vi KH6 (dac ta §3.3: soMucTieu
// "luon ghi 1" ngay trong ban Java goc — client ho tro nhieu muc tieu qua ten
// bien nhung KH6 chi danh 1 nguoi/luot; loaiChieu "= 0" cung la hang so ngay
// trong ban Java goc, chua ro y nghia nao khac) — vi vay KHONG lam tham so
// ham, tranh nguy co nguoi goi truyen nham gia tri.
func goiDanhNguoiChoi(idNguoiDanh int32, skillId byte, idNguoiTrung int32, satThuong int32, daChet bool, chiMang bool) *protocol.Message {
	w := protocol.NewWriter()
	w.Int(idNguoiDanh)
	w.Byte(skillId)
	w.Byte(1) // soMucTieu — LUON 1 o KH6, xem chu thich ham
	w.Int(idNguoiTrung)
	w.Byte(1) // docTiep = 1 -> danh dau day la DANG DON DANH (18 byte), dac ta §2.1
	w.Byte(0) // loaiChieu — hang so 0, xem chu thich ham
	w.Int(satThuong)
	w.Bool(daChet)
	w.Bool(chiMang)
	return protocol.NewMessage(lenhNguoiDanhNguoi, w)
}

// goiHoiMauNguoiChoi dung goi XUONG -60 DANG "HOI MAU" (docTiep=0, 11 byte
// than, HET GOI NGAY SAU docTiep), dac ta §2.1/§3.3, Ke hoach 6 Task 10: chieu
// hoi mau/hoi sinh dong doi (vi du Y Lieu Nhan Thuat, chieu type=2 — dac ta
// §4.3 noi ro chieu type=2 la BAT BUOC vi day la duong chay THAT DUY NHAT cua
// ham nay).
//
// Chi ghi 4 truong dau (idNguoiDanh, skillId, soMucTieu, idNguoiTrung) roi
// docTiep=0 la DUNG NGAY — KHONG duoc ghi them loaiChieu/satThuong/daChet/
// chiMang. Xem chu thich goiDanhNguoiChoi de biet day du ly do BAT BUOC tach
// rieng thanh hai ham thay vi mot ham co tham so co.
func goiHoiMauNguoiChoi(idNguoiDanh int32, skillId byte, idNguoiTrung int32) *protocol.Message {
	w := protocol.NewWriter()
	w.Int(idNguoiDanh)
	w.Byte(skillId)
	w.Byte(1) // soMucTieu — LUON 1 o KH6, xem chu thich goiDanhNguoiChoi
	w.Int(idNguoiTrung)
	w.Byte(0) // docTiep = 0 -> danh dau day la DANG HOI MAU (11 byte), dac ta §2.1
	// KHONG GHI GI THEM SAU DAY — dung o day la HET GOI. Xem chu thich ham va
	// dac ta §2.0/§2.1 (ghi thua truong o day la loi im lang tuyet doi).
	return protocol.NewMessage(lenhNguoiDanhNguoi, w)
}
