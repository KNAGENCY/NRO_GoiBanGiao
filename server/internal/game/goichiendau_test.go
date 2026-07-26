package game

import (
	"bytes"
	"testing"
)

// TestGoiNguoiDanhQuaiByteDung: goi 54, dac ta §3.3 — CHI 3 truong, KHONG co
// truong dem sau mobId (dac ta §2.0: client doc mobId cho toi het than goi).
func TestGoiNguoiDanhQuaiByteDung(t *testing.T) {
	m := goiNguoiDanhQuai(1001, 5, 3)

	if m.Command != 54 {
		t.Fatalf("Command = %d, muon 54", m.Command)
	}

	muon := []byte{
		0, 0, 3, 233, // idNguoiDanh = 1001 (int32 BE, 0x3E9)
		5, // skillId
		3, // mobId — PHAI la byte CUOI CUNG cua than goi
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 54 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiConSongByteDung: goi -9, dac ta §3.3. Hieu ung -1 (khong hieu
// ung) PHAI la 1 byte 0xFF, KHONG phai 4 byte writeInt(-1) nhu ban Java.
func TestGoiQuaiConSongByteDung(t *testing.T) {
	m := goiQuaiConSong(3, 80, 20, false, 255) // 255 = byte cua -1 (khong hieu ung)

	if m.Command != -9 {
		t.Fatalf("Command = %d, muon -9", m.Command)
	}

	muon := []byte{
		3,           // mobId
		0, 0, 0, 80, // hpConLai = 80 (int32)
		0, 0, 0, 20, // satThuong = 20 (int32)
		0,   // chiMang = false
		255, // hieuUngId = -1 (1 BYTE, khong phai 4 byte)
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -9 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiConSongChiMangTrueByteDung: kiem gia tri chiMang=true va
// hieuUngId khac -1 de khong trung tinh co voi test truoc.
func TestGoiQuaiConSongChiMangTrueByteDung(t *testing.T) {
	m := goiQuaiConSong(7, 500, 999, true, 4)

	muon := []byte{
		7,            // mobId
		0, 0, 1, 244, // hpConLai = 500
		0, 0, 3, 231, // satThuong = 999
		1, // chiMang = true
		4, // hieuUngId = 4
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -9 (chiMang=true) byte sai:\n got  = %v\n muon = %v", m.Data, muon)
	}
}

// TestGoiQuaiChetByteDung_SoDoRoiRong: goi -12 voi soDoRoi = 0 (KH6 chua lam
// he vat pham, dac ta §0 "Ngoai KH6"), theo dung yeu cau cua brief Task 6.
func TestGoiQuaiChetByteDung_SoDoRoiRong(t *testing.T) {
	m := goiQuaiChet(2, 150, true, nil)

	if m.Command != -12 {
		t.Fatalf("Command = %d, muon -12", m.Command)
	}

	muon := []byte{
		2,            // mobId
		0, 0, 0, 150, // satThuong = 150
		1, // chiMang = true
		0, // soDoRoi = 0
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -12 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiChetByteDung_CoDoRoi: goi -12 voi 2 do roi, kiem TOAN BO 5 truong
// cua moi phan tu lap (itemMapId, templateId, x, y, idNguoiNhat) dung thu tu
// dac ta §3.3.
func TestGoiQuaiChetByteDung_CoDoRoi(t *testing.T) {
	doRoi := []DoRoi{
		{ItemMapId: 10, TemplateId: 20, X: 100, Y: 200, IdNguoiNhat: 1001},
		{ItemMapId: 11, TemplateId: 21, X: 101, Y: 201, IdNguoiNhat: 1002},
	}
	m := goiQuaiChet(5, 999, false, doRoi)

	muon := []byte{
		5,            // mobId
		0, 0, 3, 231, // satThuong = 999
		0, // chiMang = false
		2, // soDoRoi = 2

		0, 10, // itemMapId = 10
		0, 20, // templateId = 20
		0, 100, // x = 100
		0, 200, // y = 200
		0, 0, 3, 233, // idNguoiNhat = 1001

		0, 11, // itemMapId = 11
		0, 21, // templateId = 21
		0, 101, // x = 101
		0, 201, // y = 201
		0, 0, 3, 234, // idNguoiNhat = 1002
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -12 (co do roi) byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiHoiSinhByteDung: goi -13, dac ta §3.3. sys va levelBoss LUON 0
// (dac ta: client khong doi loai quai khi hoi sinh, levelBoss != 0 bat hieu
// ung sieu cap ngau nhien khong mong muon).
func TestGoiQuaiHoiSinhByteDung(t *testing.T) {
	m := goiQuaiHoiSinh(4, 250)

	if m.Command != -13 {
		t.Fatalf("Command = %d, muon -13", m.Command)
	}

	muon := []byte{
		4,            // mobId
		0,            // sys = 0
		0,            // levelBoss = 0
		0, 0, 0, 250, // hp = 250
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -13 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// ---------------------------------------------------------------------------
// Task 10 — goi -60 "nguoi danh nguoi", dac ta §2.1/§3.3. HAI DANG DO DAI,
// dung HAI HAM TACH RIENG (goiDanhNguoiChoi / goiHoiMauNguoiChoi) — KHONG mot
// ham co tham so co "docTiep". Xem chu thich hai ham trong goichiendau.go de
// biet ly do: mot ham co co thi nguoi sua sau rat de vo tinh ghi thua 4
// truong cuoi vao nhanh hoi mau, ma ghi thua la loi IM LANG TUYET DOI (dac ta
// §2.0: client vut byte du, khong loi, khong trieu chung, chay thu KHONG PHAT
// HIEN DUOC).
// ---------------------------------------------------------------------------

// TestGoiDanhNguoiChoiByteDung: goi -60 DANG "DON DANH" (docTiep=1, 18 byte
// than), dac ta §2.1/§3.3. daChet va chiMang co tinh CHU Y dung hai gia tri
// KHAC NHAU (true/false) de test bat duoc loi doi cho giua hai truong bool
// lien nhau (vi du lo viet nguoc thu tu Bool(daChet)/Bool(chiMang)).
func TestGoiDanhNguoiChoiByteDung(t *testing.T) {
	m := goiDanhNguoiChoi(2001, 9, 3002, 500, true, false)

	if m.Command != -60 {
		t.Fatalf("Command = %d, muon -60", m.Command)
	}

	muon := []byte{
		0, 0, 7, 209, // idNguoiDanh = 2001 (int32 BE)
		9,             // skillId
		1,             // soMucTieu = 1 (luon 1 o KH6, dac ta §3.3)
		0, 0, 11, 186, // idNguoiTrung = 3002 (int32 BE)
		1,            // docTiep = 1 -> danh dau DANG DON DANH, dac ta §2.1
		0,            // loaiChieu = 0 (hang so ngay trong ban Java goc)
		0, 0, 1, 244, // satThuong = 500 (int32 BE)
		1, // daChet = true
		0, // chiMang = false
	}
	if len(muon) != 18 {
		t.Fatalf("mang 'muon' cua test tu viet sai, dai %d byte, phai 18 — sua lai test truoc", len(muon))
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -60 (don danh) byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiHoiMauNguoiChoiByteDung: goi -60 DANG "HOI MAU" (docTiep=0, 11 byte
// than, HET GOI NGAY SAU docTiep — khong con truong nao nua), dac ta §2.1/§3.3.
// Day la test QUAN TRONG NHAT cua Task 10 (theo nguyen van brief) — hai test
// dot bien o cuoi file nay (chay tay, khong nam trong bo test tu dong) xac
// nhan test nay THAT SU bat duoc loi ghi thua truong vao dang hoi mau.
func TestGoiHoiMauNguoiChoiByteDung(t *testing.T) {
	m := goiHoiMauNguoiChoi(4004, 12, 5005)

	if m.Command != -60 {
		t.Fatalf("Command = %d, muon -60", m.Command)
	}

	muon := []byte{
		0, 0, 15, 164, // idNguoiDanh = 4004 (int32 BE)
		12,            // skillId
		1,             // soMucTieu = 1 (luon 1 o KH6, dac ta §3.3)
		0, 0, 19, 141, // idNguoiTrung = 5005 (int32 BE)
		0, // docTiep = 0 -> danh dau DANG HOI MAU, HET GOI — dac ta §2.1
	}
	if len(muon) != 11 {
		t.Fatalf("mang 'muon' cua test tu viet sai, dai %d byte, phai 11 — sua lai test truoc", len(muon))
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -60 (hoi mau) byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiNguoiDanhNguoiChenhLech7Byte: doi chieu truc tiep dac ta §2.1 — hai
// dang cua goi -60 phai chenh nhau DUNG 7 byte (loaiChieu 1 + satThuong 4 +
// daChet 1 + chiMang 1 = 7). Test nay chi doc DO DAI than goi (khong doc gia
// tri byte), bo sung cho hai test byte-chinh-xac o tren — huu ich rieng cho
// truong hop sau nay co ai doi thu tu truong ben trong ma quen cap nhat tong
// do dai (hai test kia se bat duoc ngay, nhung test nay noi thang vao dung
// con so ma dac ta nhac toi).
func TestGoiNguoiDanhNguoiChenhLech7Byte(t *testing.T) {
	donDanh := goiDanhNguoiChoi(1, 1, 2, 1, false, false)
	hoiMau := goiHoiMauNguoiChoi(1, 1, 2)

	if len(donDanh.Data) != 18 {
		t.Fatalf("than goi -60 (don danh) dai %d byte, dac ta §2.1 doi hoi 18", len(donDanh.Data))
	}
	if len(hoiMau.Data) != 11 {
		t.Fatalf("than goi -60 (hoi mau) dai %d byte, dac ta §2.1 doi hoi 11", len(hoiMau.Data))
	}
	if chenh := len(donDanh.Data) - len(hoiMau.Data); chenh != 7 {
		t.Fatalf("hai dang goi -60 chenh %d byte, dac ta §2.1 doi hoi DUNG 7 (loaiChieu 1 + satThuong 4 + daChet 1 + chiMang 1)", chenh)
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 9 — goi -11 (quai danh trung minh) va -10 (quai danh trung
// nguoi khac), dac ta §3.3 Nhom D "Quai danh nguoi". Ca hai goi KHONG co
// truong chiMang/hieuUng — chi cac truong lay thang tu bang byte dac ta.
// ---------------------------------------------------------------------------

// TestGoiQuaiDanhTrungMinhByteDung: goi -11, dac ta §3.3 — CHI 2 truong
// (mobId, satThuong), KHONG co truong nao khac (truong satThuongChakra tuy
// chon ma dac ta nhac toi la KHONG duoc ban Java lan ban Go gui).
func TestGoiQuaiDanhTrungMinhByteDung(t *testing.T) {
	m := goiQuaiDanhTrungMinh(3, 42)

	if m.Command != -11 {
		t.Fatalf("Command = %d, muon -11", m.Command)
	}

	muon := []byte{
		3,           // mobId
		0, 0, 0, 42, // satThuong (int32 BE)
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -11 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiDanhTrungNguoiKhacByteDung: goi -10, dac ta §3.3 — 3 truong
// (mobId, idNanNhan, hpConLai). idNanNhan la int32 BE (id nguoi choi la int,
// dac ta §3 quy uoc).
func TestGoiQuaiDanhTrungNguoiKhacByteDung(t *testing.T) {
	m := goiQuaiDanhTrungNguoiKhac(7, 1001, 58)

	if m.Command != -10 {
		t.Fatalf("Command = %d, muon -10", m.Command)
	}

	muon := []byte{
		7,            // mobId
		0, 0, 3, 233, // idNanNhan = 1001 (int32 BE, 0x3E9)
		0, 0, 0, 58, // hpConLai
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -10 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 13 — goi CHET va HOI SINH nguoi choi (dac ta §3.4, §6.4).
// -17/-8 la HAI TEST BAT BUOC cua brief (danh dau (a)/(b) trong chu thich
// tung ham) — CHONG CHEP dung hai loi THAT cua ban Java da xac minh bang
// client (dac ta §7).
// ---------------------------------------------------------------------------

// TestGoiNguoiChetByteDung la TEST BAT BUOC (a) cua brief Task 13: than goi
// -17 phai DUNG 5 BYTE (byte cPk + short x + short y) — TUYET DOI KHONG co
// truong id. Ban Java (services/Service.java:1147) nhet id nhan vat vao DUNG
// o cPk la SAI (dac ta §7.2) — khong lam lech DO DAI (van 5 byte) nen mot
// test chi kiem do dai la CHUA DU; test nay con kiem GIA TRI tung truong de
// bat duoc ca loi "dat dung o, sai y nghia".
func TestGoiNguoiChetByteDung(t *testing.T) {
	m := goiNguoiChet(0, 100, 384)

	if m.Command != -17 {
		t.Fatalf("Command = %d, muon -17", m.Command)
	}
	if len(m.Data) != 5 {
		t.Fatalf("goi -17 dai %d byte, dac ta §7.2/§3.4 doi hoi DUNG 5 (byte cPk + short x + short y, KHONG co id)", len(m.Data))
	}
	muon := []byte{
		0,      // cPk = 0 (chua co he PK)
		0, 100, // x = 100
		1, 128, // y = 384 (0x0180)
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -17 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiAiDoChetByteDung la TEST BAT BUOC (b) cua brief Task 13: than goi -8
// phai DUNG 9 BYTE, TRUONG DAU la INT (4 byte). Ban Java co mot cho ghi SHORT
// (services/Service.java:1157) la LOI THAT dang chay — client doc bang
// readInt() (dac ta §7.1) nen tra id hut, ca goi bi bo qua, nguoi choi khac
// khong thay ai vua chet. Dung id > 65535 (khong vua trong 2 byte) de neu lo
// quay ve short thi gia tri doc duoc SE SAI HOAN TOAN, khong chi lech byte.
func TestGoiAiDoChetByteDung(t *testing.T) {
	m := goiAiDoChet(70000, 0, 100, 384)

	if m.Command != -8 {
		t.Fatalf("Command = %d, muon -8", m.Command)
	}
	if len(m.Data) != 9 {
		t.Fatalf("goi -8 dai %d byte, dac ta §7.1/§3.4 doi hoi DUNG 9 (int idNguoiChoi + byte cPk + short x + short y)", len(m.Data))
	}
	muon := []byte{
		0, 1, 17, 112, // idNguoiChoi = 70000 (int32 BE — CHI dung neu la INT thuc su)
		0,      // cPk = 0
		0, 100, // x = 100
		1, 128, // y = 384
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -8 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiHoiSinhChinhChuByteDung: goi -16 XUONG (xac nhan chinh chu vua hoi
// sinh) phai co than RONG (dac ta §3.4).
func TestGoiHoiSinhChinhChuByteDung(t *testing.T) {
	m := goiHoiSinhChinhChu()

	if m.Command != -16 {
		t.Fatalf("Command = %d, muon -16", m.Command)
	}
	if len(m.Data) != 0 {
		t.Fatalf("goi -16 dai %d byte, dac ta §3.4 doi hoi RONG", len(m.Data))
	}
}

// TestGoiAiDoHoiSinhByteDung: goi -30/15 (dac ta §3.4) — dung LAI ma lenh -30
// cua ME_INFO (thegioi.go) voi byte lenh con 15. Kiem CA byte lenh con LAN 5
// truong sau do, dung thu tu.
func TestGoiAiDoHoiSinhByteDung(t *testing.T) {
	m := goiAiDoHoiSinh(1001, 150, 100, 200, 300)

	if m.Command != -30 {
		t.Fatalf("Command = %d, muon -30 (dung lai ma lenh cua ME_INFO, khac byte lenh con)", m.Command)
	}
	muon := []byte{
		15,           // byte lenh con = 15 (ai do hoi sinh)
		0, 0, 3, 233, // idNguoiChoi = 1001
		0, 0, 0, 150, // hp = 150
		0, 0, 0, 100, // mp = 100
		0, 200, // x = 200
		1, 44, // y = 300 (0x012C)
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -30/15 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 11 — goi -30/9 "cap nhat mau nguoi choi" (dac ta §3.4),
// bao NHUNG NGUOI KHAC (KHAC voi -30/15 o tren la "ca khu vuc") mot nguoi
// choi vua doi mau sau mot don PvP (goi -60 khong tu mang hpConLai).
// ---------------------------------------------------------------------------

// TestGoiCapNhatMauByteDung: goi -30/9, dac ta §3.4 — 3 truong
// (idNguoiChoi, hp, hpMax), dung LAI ma lenh -30 voi byte lenh con 9 (KHAC 15
// cua -30/15 "ai do hoi sinh" — hai gia tri nay KHONG duoc trung nhau, neu
// khong client se hieu nham hai su kien khac han nhau thanh mot).
func TestGoiCapNhatMauByteDung(t *testing.T) {
	m := goiCapNhatMau(1001, 80, 150)

	if m.Command != -30 {
		t.Fatalf("Command = %d, muon -30 (dung lai ma lenh cua ME_INFO, khac byte lenh con)", m.Command)
	}
	muon := []byte{
		9,            // byte lenh con = 9 (cap nhat mau) — KHAC 15 (ai do hoi sinh)
		0, 0, 3, 233, // idNguoiChoi = 1001
		0, 0, 0, 80, // hp = 80
		0, 0, 0, 150, // hpMax = 150
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -30/9 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiCapNhatMauSubKhacAiDoHoiSinh: doi chieu truc tiep hai hang so byte
// lenh con dung chung ma lenh -30 (subCapNhatMau va subAiDoHoiSinh,
// goichiendau.go) — neu ai lo gan trung gia tri (vd copy-paste doi 15 thanh
// 9 nhung quen sua o cho khac), test nay bat duoc NGAY o muc hang so, truoc
// ca khi can so byte goi that.
func TestGoiCapNhatMauSubKhacAiDoHoiSinh(t *testing.T) {
	if subCapNhatMau == subAiDoHoiSinh {
		t.Fatalf("subCapNhatMau (%d) == subAiDoHoiSinh (%d) — hai su kien PHAI co byte lenh con KHAC NHAU, neu khong client se hieu nham 'cap nhat mau' thanh 'ai do hoi sinh' hoac nguoc lai", subCapNhatMau, subAiDoHoiSinh)
	}
}
