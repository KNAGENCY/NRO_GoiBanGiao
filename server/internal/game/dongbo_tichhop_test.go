package game

// Kiem thu dau-cuoi (Task 6, Ke hoach 3): mo phong MOT client gia di TRON
// version dance — dang nhap -> tao nhan vat -> nhan -28/4 (versionGame) ->
// xin ca 4 catalog du lieu tinh (-87 data, -28/6 map, -28/7 skill, -28/8
// item) -> nhan tra loi BYTE-DUNG cho tung cai -> gui -28/13 (client bao san
// sang) -> may chu con phuc vu tiep.
//
// Day la lop kiem thu KHAC HAN cac test dieu phoi rieng le trong dongbo_test.go
// (goiVersionGame()/goiCatalogData()/goiCatalogSkill()/... goi TRUC TIEP, khong
// qua mang) va login_test.go/tichhop_test.go (moi file chi phu MOT doan cua
// chuoi): o day di HET ca chuoi qua mang that (net.Pipe + mynet.PhucVu that)
// va CSDL that (t.Skip neu thieu NRO_TEST_DSN, giong moi file _tichhop_test.go
// khac), dung DUNG mot client gia tu dau chi cuoi — day chinh la kich ban ma
// mot client Unity that se trai qua khi noi lan dau.
//
// Vi sao goi -87 can cach doc/ghi RIENG (khong dung cl.f.WriteMessage/
// ReadMessage nhu -28/x): -87 la lenh dac biet do dai BA byte (goiCatalogData
// ~87 KB vuot qua gioi han HAI byte 65535 — xem lenhDoDaiBaByte trong
// protocol/frame.go), va dinh dang BA byte CHI dung cho MOT chieu that su
// (may-chu-gui-xuong); client LUON gui yeu cau bang HAI byte du la lenh gi.
// Framer.WriteMessage/ReadMessage cua chinh cl.f khong the dung THANG cho
// chieu nguoc (xem chu thich chi tiet o guiGoiXinDataThuCong/
// docGoiLenhDacBietThuCong trong login_test.go cho ban dau tien cua van de
// nay, va o Framer.KeyStream trong protocol/frame.go cho ly do vi sao ham do
// duoc them). O DAY khac han cach dung cua login_test.go (TestPhucVuGoiXinData)
// mot diem quan trong: -87 KHONG con la giao dich DAU TIEN tren ket noi nua —
// truoc do da co dang nhap + tao nhan vat + doc -28/4 di qua cl.f roi, nen
// KHONG the dung mot KeyStream MOI bat dau tu con tro 0 (se lech khoi con tro
// that ma may chu dang dung, giai ma ra rac ngay goi dau). Phai dung DUNG con
// tro cl.f.KeyStream() (cung mot doi tuong KeyStream ma cl.f dang dung cho
// MOI thao tac WriteMessage/ReadMessage khac) de con tro tich luy dung tiep
// tuc — xem xinData/doiGoiLenh/docMotGoiThuCong duoi day.

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// hanDoiGoi la han gio cho MOI lan doi mot goi cu the qua doiGoi/doiGoiLenh —
// du dai cho net.Pipe (khong do tre mang that) + vai truy van CSDL that da
// lam nong tu truoc (xem hanChoImLang trong tichhop_test.go cho lap luan
// tuong tu chi tiet hon).
const hanDoiGoi = 5 * time.Second

// doiGoi doc CAC goi qua cl.f (Framer, dinh dang HAI byte thuong — dung cho
// moi lenh -28/x, khong phai lenh dac biet) TOI KHI gap goi mang dung
// Command == cmd VA byte lenh con dau tien (trong Data) == sub, BO QUA goi
// khac (khong khop cmd, hoac khop cmd nhung khac sub). Dat han gio hanDoiGoi
// cho CA vong lap (khong phai tung goi rieng) de khong treo vinh vien neu may
// chu khong bao gio tra loi dung goi mong doi.
func (cl *clientGia) doiGoi(t *testing.T, cmd int8, sub byte) *protocol.Message {
	t.Helper()
	_ = cl.conn.SetReadDeadline(time.Now().Add(hanDoiGoi))
	defer cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // khoi phuc han "binh thuong" cho buoc sau

	for {
		m, err := cl.f.ReadMessage()
		if err != nil {
			t.Fatalf("khong tim thay goi %d/%d truoc khi het han %s: %v", cmd, sub, hanDoiGoi, err)
		}
		if m.Command != cmd {
			continue
		}
		r := protocol.NewReader(m.Data)
		subNhan, errSub := r.Byte()
		if errSub != nil || subNhan != sub {
			continue
		}
		return m
	}
}

// docMotGoiThuCong doc MOT goi tu conn, dung DUNG thuat toan client that (can
// cu VAO LENH, khong can cu than tin, de quyet dinh doc hai hay BA byte do
// dai — xem docGoiLenhDacBietThuCong trong login_test.go cho ban dau tien),
// nhung khac ham do o HAI diem:
//  1. Dung CHINH con tro k truyen vao (khong tu tao con tro moi) — xem chu
//     thich dau file va o Framer.KeyStream (protocol/frame.go) cho ly do.
//  2. Tra loi qua error thay vi t.Fatal, de ham goi (doiGoiLenh) tu quyet
//     dinh xu ly loi (timeout thi bo qua/dung lai co chu dich, loi khac thi
//     that bai ro rang) thay vi bi dung ngang tai day.
func docMotGoiThuCong(conn net.Conn, k *protocol.KeyStream) (*protocol.Message, error) {
	var cmdByte [1]byte
	if _, err := io.ReadFull(conn, cmdByte[:]); err != nil {
		return nil, err
	}
	cmd := int8(k.Read(cmdByte[0]))

	var n int
	if cmd == lenhDuLieu {
		var doDai [3]byte
		if _, err := io.ReadFull(conn, doDai[:]); err != nil {
			return nil, err
		}
		b0 := k.Read(doDai[0]) + 128
		b1 := k.Read(doDai[1]) + 128
		b2 := k.Read(doDai[2]) + 128
		n = int(b0) | int(b1)<<8 | int(b2)<<16
	} else {
		var doDai [2]byte
		if _, err := io.ReadFull(conn, doDai[:]); err != nil {
			return nil, err
		}
		hi := k.Read(doDai[0])
		lo := k.Read(doDai[1])
		n = int(hi)<<8 | int(lo)
	}

	data := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(conn, data); err != nil {
			return nil, err
		}
		for i := range data {
			data[i] = k.Read(data[i])
		}
	}
	return &protocol.Message{Command: cmd, Data: data}, nil
}

// doiGoiLenh giong het doiGoi ve tinh than (doc toi khi gap dung lenh mong
// doi, bo qua goi khac, co han gio tranh treo), nhung dung docMotGoiThuCong
// (biet ca dinh dang BA byte) thay vi cl.f.ReadMessage() — bat buoc cho lenh
// dac biet (o day chi -87/lenhDuLieu can, nhung ham viet chung cho ro rang).
func (cl *clientGia) doiGoiLenh(t *testing.T, lenh int8) *protocol.Message {
	t.Helper()
	k := cl.f.KeyStream()
	_ = cl.conn.SetReadDeadline(time.Now().Add(hanDoiGoi))
	defer cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		m, err := docMotGoiThuCong(cl.conn, k)
		if err != nil {
			t.Fatalf("khong tim thay goi lenh %d truoc khi het han %s: %v", lenh, hanDoiGoi, err)
		}
		if m.Command == lenh {
			return m
		}
	}
}

// xinData gui THU CONG goi yeu cau -87 than rong (client xin data), dung DUNG
// dinh dang HAI byte do dai ma client that gui (Session_ME.cs, doSendMessage
// — LUON hai byte cho MOI goi client gui LEN, bat ke lenh gi, xem
// guiGoiXinDataThuCong trong login_test.go cho ban dau tien cua logic nay).
// Khac ham do o cho dung CHINH con tro cl.f.KeyStream() (xem chu thich dau
// file) thay vi mot con tro moi bat dau tu 0.
func (cl *clientGia) xinData() {
	cl.t.Helper()
	k := cl.f.KeyStream()
	cmd := lenhDuLieu // bien (khong phai hang), de byte(cmd) la chuyen doi kieu binh thuong
	// (Go chan convert hang so am sang byte luc bien dich: "constant overflows
	// byte" — xem chu thich cung ly do o guiGoiXinDataThuCong, login_test.go).
	goi := []byte{
		k.Write(byte(cmd)),
		k.Write(0), k.Write(0), // do dai = 0 (than rong), HAI byte thuong
	}
	if _, err := cl.conn.Write(goi); err != nil {
		cl.t.Fatalf("gui goi xin data (-87) that bai: %v", err)
	}
}

// xinCatalog gui goi -28/<lenhCon> (client xin catalog map/skill/item, hoac
// bao san sang voi lenhCon = lenhConClientOk = 13) qua cl.f.WriteMessage —
// an toan dung Framer o day vi -28 KHONG nam trong lenhDoDaiBaByte, nen dinh
// dang HAI byte thuong dung cho ca hai chieu mot khi da co khoa.
func (cl *clientGia) xinCatalog(lenhCon byte) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Byte(lenhCon)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		cl.t.Fatalf("gui goi -28/%d that bai: %v", lenhCon, err)
	}
}

// guiClientOk gui goi -28/13 (client bao da nhan du 4 catalog, san sang).
func (cl *clientGia) guiClientOk() {
	cl.t.Helper()
	cl.xinCatalog(lenhConClientOk)
}

// ---------------------------------------------------------------------------
// Cac ham kiem* khang dinh THAN TIN byte-dung cho tung goi, doi chieu voi tai
// lieu giao thuc §2-§6 va DUNG cong thuc ma dongbo_test.go da kiem tra o muc
// don vi (goiVersionGame/goiCatalogData/goiCatalogSkill/goiCatalogMap/
// goiCatalogItem) — o day chi khac cho DU LIEU DEN QUA MANG THAT (da qua ma
// hoa/giai ma qua net.Pipe that) thay vi goi ham tao goi truc tiep trong bo
// nho.
// ---------------------------------------------------------------------------

// kiemVersionGame khang dinh than goi -28/4: [4][vsData][vsMap][vsSkill][vsItem][0][0].
func kiemVersionGame(t *testing.T, m *protocol.Message) {
	t.Helper()
	want := []byte{lenhConVersionGame, vsData, vsMap, vsSkill, vsItem, 0, 0}
	if !bytes.Equal(m.Data, want) {
		t.Fatalf("than goi -28/4 (versionGame) = %v, mong doi %v (version phải theo asset lock)", m.Data, want)
	}
}

// kiemDataCoPartVaDuoi100KB khang dinh goi -87: lenh dung, tong do dai duoi
// 100000 byte (mốc mark client), vsData đúng, và ĐÚNG 6 blob theo thứ tự
// dart/arrow/effect/image/part/skill; part + skill phải byte-identical artifact —
// đối chiếu với TestGoiCatalogDataDungCauTruc trong
// dongbo_test.go.
func kiemDataCoPartVaDuoi100KB(t *testing.T, m *protocol.Message) {
	t.Helper()
	if m.Command != lenhDuLieu {
		t.Fatalf("lenh goi data = %d, mong doi %d (-87)", m.Command, lenhDuLieu)
	}
	if len(m.Data) >= 100000 {
		t.Fatalf("goi -87 = %d byte, phai < 100000 (moc mark client)", len(m.Data))
	}

	r := protocol.NewReader(m.Data)
	if b, err := r.Byte(); err != nil || b != vsData {
		t.Fatalf("byte dau goi -87 = %d (loi doc: %v), mong doi vsData = %d", b, err, vsData)
	}

	docBlob := func(ten string) []byte {
		n, err := r.Int()
		if err != nil {
			t.Fatalf("doc do dai blob %s that bai: %v", ten, err)
		}
		data := make([]byte, n)
		for i := int32(0); i < n; i++ {
			value, err := r.Byte()
			if err != nil {
				t.Fatalf("doc noi dung blob %s (byte %d/%d) that bai: %v", ten, i, n, err)
			}
			data[i] = value
		}
		return data
	}

	if data := docBlob("dart"); !bytes.Equal(data, []byte{0, 0}) {
		t.Errorf("blob dart = %v, mong doi bang hieu ung trong 128 o", data)
	}
	if data := docBlob("arrow"); !bytes.Equal(data, []byte{0, 0}) {
		t.Errorf("blob arrow = %v, mong doi bang hieu ung trong 128 o", data)
	}
	if data := docBlob("effect"); !bytes.Equal(data, duLieuHieuUngTrong) {
		t.Errorf("blob effect = %v, mong doi bang hieu ung trong 128 o", data)
	}
	if data := docBlob("image"); !bytes.Equal(data, []byte{0, 0}) {
		t.Errorf("blob image = %v, mong doi bang hieu ung trong 128 o", data)
	}
	if data := docBlob("part"); !bytes.Equal(data, duLieuPart) {
		t.Errorf("blob part không byte-identical artifact nhúng (nhận %d, muốn %d)", len(data), len(duLieuPart))
	}
	if data := docBlob("skill"); !bytes.Equal(data, duLieuSkill) {
		t.Errorf("blob skill không byte-identical artifact nhúng (nhận %d, muốn %d)", len(data), len(duLieuSkill))
	}

	if _, err := r.Byte(); err == nil {
		t.Error("goi -87 con byte thua sau 6 blob — phai het buffer")
	}
}

// kiemCatalogMapRong giu ten de test cu goi lai, nhung doc catalog Lang La that.
func kiemCatalogMapRong(t *testing.T, m *protocol.Message) {
	t.Helper()
	r := protocol.NewReader(m.Data)
	if sub, err := r.Byte(); err != nil || sub != lenhConCatalogMap {
		t.Fatalf("sub=%d err=%v", sub, err)
	}
	if v, _ := r.Byte(); v != vsMap {
		t.Fatalf("vsMap=%d", v)
	}
	n, _ := r.Byte()
	if n != 6 {
		t.Fatalf("map count=%d", n)
	}
	for _, mapCatalog := range catalogLangLa.Maps {
		got, err := r.UTF()
		if err != nil || got != mapCatalog.Name {
			t.Fatalf("map=%q err=%v", got, err)
		}
	}
	n, _ = r.Byte()
	if n != 1 {
		t.Fatalf("npc count=%d", n)
	}
	if name, err := r.UTF(); err != nil || name != "Iruka" {
		t.Fatalf("npc name=%q err=%v", name, err)
	}
	for _, field := range []struct {
		label string
		want  int16
	}{{"head", 2189}, {"body", 2190}, {"leg", 2191}} {
		label, want := field.label, field.want
		got, err := r.Short()
		if err != nil || got != want {
			t.Fatalf("npc %s=%d err=%v", label, got, err)
		}
	}
	if _, err := r.Byte(); err != nil {
		t.Fatal(err)
	}
	n, _ = r.Byte()
	if n != 1 {
		t.Fatalf("mob count=%d", n)
	}
	if typ, err := r.Byte(); err != nil || typ != 0 {
		t.Fatalf("mob type=%d err=%v", typ, err)
	}
	if name, err := r.UTF(); err != nil || name != "Sói Rừng" {
		t.Fatalf("mob name=%q err=%v", name, err)
	}
	if hp, err := r.Int(); err != nil || hp != 100 {
		t.Fatalf("mob hp=%d err=%v", hp, err)
	}
	for _, want := range []byte{100, 0, 0} {
		got, err := r.Byte()
		if err != nil || got != want {
			t.Fatalf("mob runtime=%d want=%d err=%v", got, want, err)
		}
	}
	if _, err := r.Byte(); err == nil {
		t.Fatal("catalog map co byte thua")
	}
}

// kiemCatalogSkill5Lang la khang dinh QUAN TRONG NHAT cua Task 6: -28/7 phai
// chua DUNG 5 lang, DUNG TEN. Ten GHIM DOC LAP voi DanhSachLang()/danhSachLang
// (khong lay tu lang.go) — neu ca hai phia cung doc chung mot nguon thi mot
// dot bien sua sai/dao thu tu ten lang trong lang.go se lam CA HAI phia doi
// theo nhau va test khong bat duoc gi (cung ly do ghim gia tri HP/Chakra/vi
// tri trong TestTaoRoiDangNhapLai_ThayNhanVat, tichhop_test.go).
func kiemCatalogSkill5Lang(t *testing.T, m *protocol.Message) {
	t.Helper()
	tenLangMongDoi := []string{"Làng Lá", "Làng Cát", "Làng Sương Mù", "Làng Đá", "Làng Mây"}

	r := protocol.NewReader(m.Data)
	sub, err := r.Byte()
	if err != nil || sub != lenhConCatalogSkill {
		t.Fatalf("lenh con = %d (loi doc: %v), mong doi %d (catalog skill)", sub, err, lenhConCatalogSkill)
	}
	if v, err := r.Byte(); err != nil || v != vsSkill {
		t.Fatalf("vsSkill = %d (loi doc: %v), mong doi %d", v, err, vsSkill)
	}
	if opt, err := r.Byte(); err != nil || opt != 0 {
		t.Fatalf("countOption = %d (loi doc: %v), mong doi 0", opt, err)
	}
	countClass, err := r.Byte()
	if err != nil {
		t.Fatalf("doc countClass that bai: %v", err)
	}
	if int(countClass) != len(tenLangMongDoi) {
		t.Fatalf("countClass = %d, mong doi DUNG %d (5 lang Naruto thay 3 lop Dragon Ball cu)", countClass, len(tenLangMongDoi))
	}
	for i := 0; i < int(countClass); i++ {
		ten, err := r.UTF()
		if err != nil {
			t.Fatalf("doc ten lang %d that bai: %v", i, err)
		}
		if ten != tenLangMongDoi[i] {
			t.Errorf("ten lang %d = %q, mong doi %q (ghim doc lap voi lang.go)", i, ten, tenLangMongDoi[i])
		}
		ct, err := r.Byte()
		if err != nil {
			t.Fatalf("doc countTemplate lang %d that bai: %v", i, err)
		}
		if ct != 4 {
			t.Errorf("countTemplate lang %d = %d, mong doi 4 (4 chieu x 5 lang, Task 1 KH6)", i, ct)
		}
		for j := 0; j < int(ct); j++ {
			if err := boQuaMotMauChieu(r); err != nil {
				t.Fatalf("doc chieu %d cua lang %d that bai: %v", j, i, err)
			}
		}
	}
}

// doiChuoiCatalogItem doc DU BA message -28/8 ma may chu tra cho mot yeu cau
// catalog vat pham va khang dinh chung den DUNG THU TU 0 -> 1 -> 2 (KH7 task 2).
//
// Thu tu la bat buoc: chi subtype 2 moi kich Service.clientOk() phia client
// (Controller.cs:6517-6534). Truoc KH7 task 2 may chu chi gui MOT message
// subtype 2 rong — moi cho goi doiGoi(-28, 8) trong bo test tich hop phai doc
// du ba, khong thi hai message con lai nam lai trong duong ong.
func doiChuoiCatalogItem(t *testing.T, cl *clientGia) {
	t.Helper()
	for i, wantSub := range []byte{0, 1, 2} {
		m := cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogItem)
		if len(m.Data) < 3 {
			t.Fatalf("message catalog item thu %d qua ngan: %v", i, m.Data)
		}
		if m.Data[1] != vsItem {
			t.Fatalf("message catalog item thu %d vsItem = %d, mong doi %d", i, m.Data[1], vsItem)
		}
		if m.Data[2] != wantSub {
			t.Fatalf("message catalog item thu %d subtype = %d, mong doi %d", i, m.Data[2], wantSub)
		}
		if wantSub == 2 {
			want := []byte{lenhConCatalogItem, vsItem, 2, 0, 0, 0, 0}
			if !bytes.Equal(m.Data, want) {
				t.Fatalf("goi dong dai -28/8 = %v, mong doi %v", m.Data, want)
			}
		}
	}
}

// kiemTileSetInfo khang dinh goi -82: lenh dung, than = NGUYEN file tileset.
func kiemTileSetInfo(t *testing.T, m *protocol.Message) {
	t.Helper()
	if m.Command != -82 {
		t.Fatalf("lenh goi tileset = %d, mong doi -82", m.Command)
	}
	if !bytes.Equal(m.Data, duLieuTileSet) {
		t.Fatalf("than -82 khong bang nguyen file tileset (dai %d, mong doi %d)", len(m.Data), len(duLieuTileSet))
	}
}

// kiemMeInfo khang dinh goi -30/0 khop nhan vat da tao: byte lenh con 0, id/gender/
// name/gold/gem dung, gender lan 2 < 5. Doi chieu voi nv lay tu CSDL that.
func kiemMeInfo(t *testing.T, m *protocol.Message, nv *store.NhanVat) {
	t.Helper()
	if m.Command != -30 {
		t.Fatalf("lenh ME_INFO = %d, mong doi -30", m.Command)
	}
	r := protocol.NewReader(m.Data)
	if sub, err := r.Byte(); err != nil || sub != 0 {
		t.Fatalf("byte lenh con ME_INFO = %d (loi %v), mong doi 0", sub, err)
	}
	if id, err := r.Int(); err != nil || id != int32(nv.ID) {
		t.Fatalf("id = %d (loi %v), mong doi %d", id, err, nv.ID)
	}
	if _, err := r.Byte(); err != nil { // taskMain
		t.Fatal(err)
	}
	if g, err := r.Byte(); err != nil || g != byte(nv.Lang) {
		t.Fatalf("gender = %d (loi %v), mong doi %d (= lang)", g, err, nv.Lang)
	}
	if _, err := r.Short(); err != nil { // head
		t.Fatal(err)
	}
	if ten, err := r.UTF(); err != nil || ten != nv.Ten {
		t.Fatalf("name = %q (loi %v), mong doi %q", ten, err, nv.Ten)
	}
	r.Byte()  // cPK
	r.Byte()  // typePk
	r.Long()  // power
	r.Short() // buffHp
	r.Short() // buffMp
	g2, err := r.Byte()
	if err != nil || g2 != byte(nv.Lang) {
		t.Fatalf("gender lan 2 = %d (loi %v), mong doi %d (= lang)", g2, err, nv.Lang)
	}
	if g2 >= 5 {
		t.Fatalf("gender lan 2 = %d, phai < countClass (5) — khong client IndexOutOfRange", g2)
	}
	if sc, err := r.Byte(); err != nil || sc != 0 {
		t.Fatalf("skillCount = %d (loi %v), mong doi 0", sc, err)
	}
	if gold, err := r.Long(); err != nil || gold != nv.Vang {
		t.Fatalf("gold = %d (loi %v), mong doi %d (long 8 byte = Vang)", gold, err, nv.Vang)
	}
	r.Int() // ruby
	if gem, err := r.Int(); err != nil || gem != int32(nv.Ngoc) {
		t.Fatalf("gem = %d (loi %v), mong doi %d (= Ngoc)", gem, err, nv.Ngoc)
	}

	// Duoi goi -30: 3 mang item + headAvatar + 3 conMeo + fusion + deltaTime +
	// isNewMember, deu 0 cho nhan vat moi. Khang dinh byte-dung CA duoi de test
	// dau-cuoi chung minh TRON goi -30 (khong dung o gem), roi HET buffer — 3 truong
	// cuoi (aura/effFront/hat) da bo han (§8 note e). Neu thieu doan nay, sai byte o
	// duoi goi se "xanh gia".
	meByte := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	meShort := func(ten string, muon int16) {
		g, err := r.Short()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	// KH7 Task 5: ba mang mang SO O that (suc chua), roi DU so o than o. Nhan
	// vat vua tao khong co mon nao -> moi o la mot short -1.
	meBangRong := func(ten string, soOMongDoi int) {
		soO, err := r.Byte()
		if err != nil || int(soO) != soOMongDoi {
			t.Fatalf("%s.size = %d (loi %v), mong doi %d (SO O, khong phai so mon)", ten, soO, err, soOMongDoi)
		}
		for i := 0; i < int(soO); i++ {
			if id, e := r.Short(); e != nil || id != -1 {
				t.Fatalf("%s o %d = %d (loi %v), mong doi -1 (o rong)", ten, i, id, e)
			}
		}
	}
	meBangRong("itemsBody", SoOMacChuan)
	meBangRong("itemsBag", SoOTuiChuan)
	meBangRong("itemsBox", SoORuongChuan)
	// headAvatar: may chu gui MOT cap (part dau -> sprite chan dung) de client ve
	// duoc o chan dung tren HUD. Truoc day gui 0 cap nen o do trong tron.
	// Ba so nay phai khop HeadNhanVat/AvatarNhanVat trong thegioi.go — lech la
	// client doc lech 4 byte roi ca doan duoi goi sai theo.
	meShort("headAvatar.size", 1)
	meShort("headAvatar.head", HeadNhanVat)
	meShort("headAvatar.avatar", AvatarNhanVat)
	meShort("conMeo0", 0)
	meShort("conMeo1", 0)
	meShort("conMeo2", 0)
	meByte("fusion", 0)
	if dt, err := r.Int(); err != nil || dt != 0 {
		t.Fatalf("deltaTime = %d (loi %v), mong doi 0", dt, err)
	}
	meByte("isNewMember", 0)
	if _, err := r.Byte(); err == nil {
		t.Error("ME_INFO con byte thua sau isNewMember — 3 truong cuoi (aura/effFront/hat) phai bo han (§8 note e)")
	}
}

// kiemMapInfo khang dinh goi -24 khop ban do cua nhan vat: mapId/tileId/name +
// byte du phong (decode het buffer den cuoi khong loi = du phong dung cho).
func kiemMapInfo(t *testing.T, m *protocol.Message, nv *store.NhanVat, bd *BanDo) {
	t.Helper()
	if m.Command != -24 {
		t.Fatalf("lenh MAP INFO = %d, mong doi -24", m.Command)
	}
	r := protocol.NewReader(m.Data)
	must := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	must("mapId", byte(bd.MapID))
	must("planetId", bd.PlanetID)
	must("tileId", bd.TileID)
	must("bgId", bd.BgID)
	must("type", bd.Loai)
	if ten, err := r.UTF(); err != nil || ten != bd.Ten {
		t.Fatalf("mapName = %q (loi %v), mong doi %q", ten, err, bd.Ten)
	}
	must("zoneId", 0)
	if x, err := r.Short(); err != nil || x != int16(nv.ViTriX) {
		t.Fatalf("x = %d (loi %v), mong doi %d", x, err, nv.ViTriX)
	}
	if y, err := r.Short(); err != nil || y != int16(nv.ViTriY) {
		t.Fatalf("y = %d (loi %v), mong doi %d", y, err, nv.ViTriY)
	}
	// Map 1 (Lang La) co DUNG MOT cong roi map; cac map con lai chua co.
	// Byte-chinh-xac cua o cong thuoc ve TestGoiMapInfoMap1CoCongRoiMap
	// (diemdi_test.go) — o day chi doc dung so byte de phan sau khong lech.
	if bd.MapID == 1 {
		must("waypointCount", 1)
		for _, ten := range []string{"minX", "minY", "maxX", "maxY"} {
			if _, err := r.Short(); err != nil {
				t.Fatalf("cong thieu truong %s: %v", ten, err)
			}
		}
		must("wp.laCong", 1)
		must("wp.ngoaiTuyen", 0)
		if _, err := r.UTF(); err != nil {
			t.Fatalf("cong thieu ten: %v", err)
		}
	} else {
		must("waypointCount", 0)
	}
	must("mobCount", 0)
	must("byte du phong", 0) // giua mob va npc
	if bd.MapID == 1 {
		must("npcCount", 1)
		must("npc.status", 0)
		if _, err := r.Short(); err != nil {
			t.Fatal(err)
		}
		if _, err := r.Short(); err != nil {
			t.Fatal(err)
		}
		must("npc.template", 0)
		if _, err := r.Short(); err != nil {
			t.Fatal(err)
		}
	} else {
		must("npcCount", 0)
	}
	must("itemCount", 0)
	// Map 1 co 11 bong do (bg-item); cac map khac chua co my thuat rieng nen
	// van la short 0. Byte-chinh-xac thuoc ve TestGoiMapInfoMap1CoBongDo.
	soBongMongDoi := int16(0)
	if bd.MapID == 1 {
		soBongMongDoi = 11
	}
	if bg, err := r.Short(); err != nil || bg != soBongMongDoi {
		t.Fatalf("bgItem = %d (loi %v), mong doi short %d", bg, err, soBongMongDoi)
	}
	for i := int16(0); i < soBongMongDoi; i++ {
		for _, ten := range []string{"mau", "tileX", "tileY"} {
			if _, err := r.Short(); err != nil {
				t.Fatalf("bong do #%d thieu truong %s: %v", i, ten, err)
			}
		}
	}
	if ef, err := r.Short(); err != nil || ef != 0 {
		t.Fatalf("effMap = %d (loi %v), mong doi short 0", ef, err)
	}
	must("bgType", bd.BgType)
	must("idSpaceShip", 0)
	must("isMapDouble", 0)
	if _, err := r.Byte(); err == nil {
		t.Error("MAP INFO con byte thua sau isMapDouble — byte du phong co the sai vi tri")
	}
}

// ---------------------------------------------------------------------------
// TestVersionDanceDayDu: kich ban dau-cuoi day du.
// ---------------------------------------------------------------------------

// TestVersionDanceDayDu mo phong MOT client gia di TRON version dance: tai
// khoan moi -> tao nhan vat -> nhan -28/4 -> xin ca 4 catalog (-87/-28/6/
// -28/7/-28/8) -> nhan tra loi byte-dung cho tung cai -> gui -28/13 -> may chu
// van con phuc vu (xac nhan bang cach xin lai skill mot lan nua). Day la muc
// tieu loi cua Ke hoach 3 duoc kiem tra o muc CAO NHAT: dung mang that + CSDL
// that, khong goi thang ham tao goi trong bo nho nhu dongbo_test.go.
func TestVersionDanceDayDu(t *testing.T) {
	khoThat := moKhoThatChoTest(t)

	ctx := context.Background()
	tenTK := "vaogame" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taiKhoanID := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)

	cl := moPhienClientGia(t, khoThat)
	defer cl.dong()

	cl.dangNhap(tenTK, matKhau)
	if !cl.choLenh(t, 2) {
		t.Fatal("tai khoan moi chua co nhan vat -> phai nhan DUNG MOT lenh 2 (man tao nhan vat)")
	}

	// Tao nhan vat -> vaoTheGioi -> may chu gui -28/4 (versionGame).
	cl.taoNhanVat("Naru"+hauToNganNgau(t), 0, 5)
	m4 := cl.doiGoi(t, lenhKhongTrenBanDo, lenhConVersionGame)
	kiemVersionGame(t, m4)

	// Client gia xin ca 4 cap nhat theo dung thu tu client that lam, nhan tra
	// loi BYTE-DUNG cho tung cai.
	cl.xinData()
	kiemDataCoPartVaDuoi100KB(t, cl.doiGoiLenh(t, lenhDuLieu))

	cl.xinCatalog(lenhConCatalogMap)
	kiemCatalogMapRong(t, cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogMap))

	cl.xinCatalog(lenhConCatalogSkill)
	kiemCatalogSkill5Lang(t, cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogSkill)) // KHANG DINH 5 lang, dung ten

	cl.xinCatalog(lenhConCatalogItem)
	doiChuoiCatalogItem(t, cl)

	// Gui client ok -> Ke hoach 4: may chu gui NOT 3 goi vao the gioi theo thu tu
	// -82 (tileset) -> -30/0 (ME_INFO) -> -24 (MAP INFO). Doc va khang dinh byte-dung.
	cl.guiClientOk()

	// Lay nhan vat vua tao tu CSDL de doi chieu ME_INFO/MAP INFO byte-dung.
	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan sau khi tao that bai: %v", err)
	}
	bd := TimBanDo(nv.MapID)
	if bd == nil {
		t.Fatalf("nhan vat o mapId %d khong co ban do trong danh muc", nv.MapID)
	}

	// doiGoiLenh khop THEO LENH (bo qua goi khac), dung duoc cho -82/-30/-24 vi ca
	// ba dung do dai 2 byte (khong nam trong lenhDoDaiBaByte). Thu tu: -82 truoc -24.
	kiemTileSetInfo(t, cl.doiGoiLenh(t, -82))
	kiemMeInfo(t, cl.doiGoiLenh(t, -30), nv)
	kiemMapInfo(t, cl.doiGoiLenh(t, -24), nv, bd)
	// -42 CHI SO gui NGAY SAU MAP INFO (login.go xuLyClientOk). Khong doc thi ong
	// pipe day, may chu ket thuc bang loi GHI het gio 30 giay — trong nhu treo,
	// khong phai loi logic. Da lam do 13 test tich hop vi thieu dung dong nay.
	cl.doiGoiLenh(t, lenhChiSoNhanVat)
}
