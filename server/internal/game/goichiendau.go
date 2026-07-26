// Package game — dung 4 goi XUONG cua vong "nguoi danh quai" (Ke hoach 6 Task
// 6): 54 (nguoi danh quai), -9 (quai con song), -12 (quai chet), -13 (quai hoi
// sinh). Ca 4 gui CHO CA KHU VUC (dac ta §3.3). Chi dung gio (byte-exact),
// KHONG tinh sat thuong/AI — do la viec cua Task 7-9.
package game

import "nroserver/internal/protocol"

// Bon ma lenh cua nhom "nguoi danh quai", dac ta §3.3.
const (
	lenhNguoiDanhQuai int8 = 54  // Nguoi danh quai — CHU Y: khong co truong dem
	lenhQuaiConSong   int8 = -9  // Quai con song sau don danh
	lenhQuaiChet      int8 = -12 // Quai chet, kem danh sach do roi (KH6 luon rong)
	lenhQuaiHoiSinh   int8 = -13 // Quai hoi sinh
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
