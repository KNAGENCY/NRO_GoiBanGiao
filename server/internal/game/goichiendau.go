// Package game — dung goi XUONG cua vong chien dau nguoi-vs-quai (Ke hoach 6
// Task 6), nguoi-vs-nguoi (Task 10), quai-vs-nguoi (Task 9) VA chet/hoi sinh
// nguoi choi (Task 13).
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
// Task 13 them: -17/-8 (chet, dac ta §3.4/§7 — hai loi THAT cua ban Java da
// XAC MINH BANG CLIENT, xem chu thich goiNguoiChet/goiAiDoChet), -16 (hoi
// sinh, cung ma lenh ca hai chieu), -30/15 (ai do hoi sinh, dung lai ma lenh
// -30 cua ME_INFO voi byte lenh con khac).
//
// Task 11 them: dieu phoi goi len -60 (nguoi danh nguoi, xem case
// suKienDanhNguoiChoi trong khuvuc.go) VA goi XUONG -30/9 (cap nhat mau cho
// nguoi khac thay, dac ta §3.4) — goi -60 (dang don danh) KHONG mang truong
// hpConLai nao nen can mot goi rieng de nhung nguoi KHONG PHAI muc tieu ve lai
// dung thanh mau cua muc tieu.
//
// Ca file chi dung goi (byte-exact), KHONG tinh sat thuong/AI/kiem tra dieu
// kien — do la viec cua cac Task/file khac (satthuong.go, chotkiem.go, quai.go).
package game

import "nroserver/internal/protocol"

// Nam ma lenh dung trong file nay: bon ma nhom "nguoi danh quai" (Task 6, dac
// ta §3.3), mot ma nhom "nguoi danh nguoi" (Task 10, dac ta §2.1/§3.3), va hai
// ma nhom "quai danh nguoi" (Task 9, dac ta §3.3).
const (
	lenhNguoiDanhQuai  int8 = 54  // Nguoi danh quai — CHU Y: khong co truong dem
	lenhQuaiConSong    int8 = -9  // Quai con song sau don danh
	lenhQuaiChet       int8 = -12 // Quai chet, kem danh sach do roi (KH6 luon rong)
	lenhQuaiHoiSinh    int8 = -13 // Quai hoi sinh
	lenhNguoiDanhNguoi int8 = -60 // Nguoi danh nguoi — CO HAI DANG DO DAI, xem goiDanhNguoiChoi/goiHoiMauNguoiChoi

	lenhQuaiDanhTrungMinh      int8 = -11 // Quai danh TRUNG chinh nan nhan — gui RIENG nan nhan
	lenhQuaiDanhTrungNguoiKhac int8 = -10 // Quai danh trung ai do — gui NGUOI KHAC trong khu vuc (khong gui nan nhan)
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

// ---------------------------------------------------------------------------
// Task 9 them: -11 (quai danh trung minh) va -10 (quai danh trung ai do) —
// Nhom D "Quai danh nguoi", dac ta §3.3, AI hai tang uu tien o quai.go.
// ---------------------------------------------------------------------------

// goiQuaiDanhTrungMinh dung goi XUONG -11: bao RIENG cho nan nhan mot con
// quai vua danh TRUNG chinh minh (dac ta §3.3). Than gui CHI 2 truong
// (mobId, satThuong) — dac ta ghi ro con mot truong tuy chon o cuoi
// (int satThuongChakra) ma CA ban Java LAN ban Go deu KHONG gui (client boc
// try/catch), nen ham nay KHONG ghi them gi sau satThuong.
//
// ⚠️ Goi cho AI QUAI (khuvuc.go, aiQuai): -11 phai gui THANG qua phien cua
// CHINH nan nhan, KHONG dung phatChoNguoiKhac/phatChoCaKhuVuc — ca hai deu
// SAI doi tuong ("gui cho: Rieng nan nhan" khac han "ca khu vuc" hay "nguoi
// khac", xem bang §3.3 cua dac ta).
func goiQuaiDanhTrungMinh(mobId byte, satThuong int32) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(mobId)
	w.Int(satThuong)
	return protocol.NewMessage(lenhQuaiDanhTrungMinh, w)
}

// goiQuaiDanhTrungNguoiKhac dung goi XUONG -10: bao NHUNG NGUOI KHAC (khong
// phai nan nhan) trong khu vuc mot con quai vua danh trung ai do (dac ta
// §3.3). Than gui 3 truong (mobId, idNanNhan, hpConLai) — cung nhu -11,
// truong tuy chon cuoi (int mpConLai) ca hai ban deu khong gui, KHONG ghi
// them gi sau hpConLai.
//
// ⚠️ Goi cho AI QUAI (khuvuc.go, aiQuai): dung k.phatChoNguoiKhac(idNanNhan, m)
// (co san tu Ke hoach 5) — dung nghia "nguoi khac" cua dac ta: loai tru CHINH
// nan nhan (nan nhan da nhan -11 rieng o tren, khong can/khong duoc nhan lai
// -10 cua chinh minh).
func goiQuaiDanhTrungNguoiKhac(mobId byte, idNanNhan int32, hpConLai int32) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(mobId)
	w.Int(idNanNhan)
	w.Int(hpConLai)
	return protocol.NewMessage(lenhQuaiDanhTrungNguoiKhac, w)
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 13 them: CHET va HOI SINH nguoi choi (dac ta §3.4, §6.4).
//
// -17/-8 la HAI GOI DA DUOC XAC MINH BANG CLIENT (dac ta §7 — doc tan noi
// Controller.cs, KHONG suy doan tu ban Java): ban Java dang chay co HAI LOI
// THAT o day. Hai ham goiNguoiChet/goiAiDoChet duoi day CHONG CHEP dung hai
// loi do, tung ham co chu thich rieng giai thich VI SAO.
//
// -16 dung CHUNG mot ma lenh CHO CA HAI CHIEU (dac ta §3.1 "goi len" va §3.4
// "goi xuong" cung ghi -16, than deu RONG): client gui LEN -16 xin hoi sinh
// (xuLyHoiSinh, login.go), may chu gui XUONG -16 RONG xac nhan (ham
// goiHoiSinhChinhChu duoi day) — hai chieu khac nhau nhung TRUNG so hieu lenh
// la binh thuong o giao thuc nay (khong phai loi).
//
// -30/15 DUNG LAI ma lenh -30 (lenhMeInfo, khai bao o thegioi.go — cung file
// voi goiMeInfo/-30/0) VOI BYTE LENH CON KHAC (15 thay vi 0) — giong cach
// -28 co nhieu lenh con (2/6/7/8/13). KHONG dinh nghia hang so lenh moi cho
// -30, chi them subAiDoHoiSinh.
// ---------------------------------------------------------------------------

const (
	lenhNguoiChet int8 = -17 // Chinh chu da chet — RIENG chinh chu, KHONG co truong id (dac ta §7.2)
	lenhAiDoChet  int8 = -8  // Ai do da chet — NGUOI KHAC trong khu vuc, id la INT (dac ta §7.1)
	lenhHoiSinh   int8 = -16 // Hoi sinh — CUNG ma lenh ca hai chieu, than LUON rong

	subAiDoHoiSinh byte = 15 // byte lenh con cua -30 bao "ai do vua hoi sinh" (dac ta §3.4)
	subCapNhatMau  byte = 9  // byte lenh con cua -30 bao "cap nhat mau nguoi choi" (dac ta §3.4, Ke hoach 6 Task 11)
)

// goiNguoiChet dung goi XUONG -17: bao RIENG chinh chu vua chet. Than gui
// DUNG 5 BYTE (byte cPk, short x, short y) — TUYET DOI KHONG co truong id
// (dac ta §7.2/§3.4, TEST BAT BUOC (a) cua brief Task 13).
//
// ⚠️ Goi nay CHI gui cho chinh chu (client dung thang Char.myCharz(), khong
// can tra id trong danh sach nguoi trong map — Controller.cs:4533-4547) nen
// KHONG CO ly do gi de ghi id ca. Ban Java (services/Service.java:1147) ghi
// `writeByte((int) pl.id)` VAO DUNG O `cPk` — khong lam lech tong do dai (van
// 1 byte) nhung GHI SAI GIA TRI cPk cua nguoi vua chet. Day chinh la kieu loi
// "dat dung o, sai y nghia" ma dac ta canh bao o §3.3 diem ve -13 (sys/
// levelBoss) — TUYET DOI KHONG chep. cPk truyen 0 khi chua co he PK (KH6
// chua co, xem dac ta §0 "Ngoai KH6": "PvP co luat -> KH9").
//
// x/y: dac ta ghi ro truong "power" (long, cPower) la TUY CHON o client
// (boc try/catch, Controller.cs:4536-4541) — ban Go BO HAN, khop hanh vi ban
// Java (dong comment "// msg.writer().writeLong(-1); send Power" bi comment
// o Service.java:1150, dac ta §7.2).
func goiNguoiChet(cPk byte, x, y int16) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(cPk)
	w.Short(x)
	w.Short(y)
	return protocol.NewMessage(lenhNguoiChet, w)
}

// goiAiDoChet dung goi XUONG -8: bao NHUNG NGUOI KHAC trong khu vuc mot
// nguoi choi vua chet. Than gui DUNG 9 BYTE, TRUONG DAU la INT 4 byte — dac
// ta §7.1/§3.4, TEST BAT BUOC (b) cua brief Task 13.
//
// ⚠️ Client doc id BANG readInt() (Controller.cs:4548-4554, dac ta §7.1) —
// ban Java co MOT CHO ghi SHORT (services/Service.java:1157) la LOI THAT
// dang chay: client tra id hut (readInt doc ca 4 byte, lay nham 2 byte cPk/x
// tiep theo lam mot phan cua id) -> findCharInMap tra null -> CA GOI BI BO
// QUA. Hau qua that: nguoi choi khac KHONG THAY ai vua chet cho toi khi ho
// roi khu vuc roi vao lai (luc do Zone.java lai phat dung ban writeInt).
// TUYET DOI KHONG duoc chep loi nay — id PHAI la int (4 byte).
func goiAiDoChet(idNguoiChoi int32, cPk byte, x, y int16) *protocol.Message {
	w := protocol.NewWriter()
	w.Int(idNguoiChoi)
	w.Byte(cPk)
	w.Short(x)
	w.Short(y)
	return protocol.NewMessage(lenhAiDoChet, w)
}

// goiHoiSinhChinhChu dung goi XUONG -16: xac nhan RIENG chinh chu vua hoi
// sinh thanh cong. Than RONG (dac ta §3.4) — CUNG ma lenh voi goi LEN -16
// (xin hoi sinh, cung rong, xem xuLyHoiSinh trong login.go) nhung o day la
// chieu MAY CHU -> CLIENT.
func goiHoiSinhChinhChu() *protocol.Message {
	return protocol.NewMessage(lenhHoiSinh, protocol.NewWriter())
}

// goiAiDoHoiSinh dung goi XUONG -30/15: bao CA KHU VUC (ke ca chinh nguoi
// vua hoi sinh — dac ta §3.4 cot "Gui cho": "Ca khu vuc", cung tinh than voi
// 54/-9/-12 o Task 6: nguoi gay ra su kien cung nhan lai) mot nguoi choi vua
// hoi sinh xong, kem mau/chakra DAY DU va vi tri hien tai.
//
// Dung LAI ma lenh -30 (lenhMeInfo, thegioi.go) voi byte lenh con 15 thay vi
// 0 — dac ta khong danh rieng ma lenh nao khac cho su kien nay, giong cach
// -28 da co nhieu lenh con. Truong dau tien LUON la subAiDoHoiSinh (15) —
// client phan biet ME_INFO (0) va "ai do hoi sinh" (15) bang chinh byte nay.
func goiAiDoHoiSinh(idNguoiChoi int32, hp, mp int32, x, y int16) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(subAiDoHoiSinh)
	w.Int(idNguoiChoi)
	w.Int(hp)
	w.Int(mp)
	w.Short(x)
	w.Short(y)
	return protocol.NewMessage(lenhMeInfo, w)
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 11 them: -30/9 (cap nhat mau nguoi choi, dac ta §3.4) — mot
// phan cua Nhom F "nguoi danh nguoi" (§2.1/§3.3). Goi -60 (dang don danh,
// goiDanhNguoiChoi o tren) KHONG mang truong hpConLai nao (khac han -9 nguoi
// danh quai hay -10 quai danh nguoi khac, ca hai da co san hpConLai) — thieu
// goi nay thi NHUNG NGUOI KHAC (khong phai chinh muc tieu) se khong biet ve
// lai dung thanh mau cua muc tieu sau mot don PvP.
// ---------------------------------------------------------------------------

// goiCapNhatMau dung goi XUONG -30/9: bao NHUNG NGUOI KHAC (khong phai chinh
// nguoi vua doi mau — dac ta §3.4 cot "Gui cho": "Người khác", KHAC voi -30/15
// o tren la "Ca khu vuc") mot nguoi choi vua doi mau, kem hp/hpMax hien tai.
//
// Dung LAI ma lenh -30 (lenhMeInfo, thegioi.go) voi byte lenh con 9 — cung ky
// thuat voi goiAiDoHoiSinh (-30/15), client phan biet bang chinh byte lenh
// con dau tien. Noi GOI ham nay (case suKienDanhNguoiChoi, khuvuc.go) PHAI
// dung k.phatChoNguoiKhac(idMucTieu, ...), KHONG dung phatChoCaKhuVuc — sai
// ham se gui goi nay CA cho chinh muc tieu, di nguoc dac ta.
func goiCapNhatMau(idNguoiChoi int32, hp, hpMax int32) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(subCapNhatMau)
	w.Int(idNguoiChoi)
	w.Int(hp)
	w.Int(hpMax)
	return protocol.NewMessage(lenhMeInfo, w)
}
