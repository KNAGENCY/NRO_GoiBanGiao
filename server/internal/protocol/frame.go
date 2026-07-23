package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

// doDaiThanToiDa chan ke xau bao mot do dai khong lo de ep may chu cap phat
// het bo nho. Goi lon nhat that su trong game la du lieu vat pham, duoi 1 MB.
const doDaiThanToiDa = 4 << 20 // 4 MB

// lenhDoDaiBaByte la cac lenh dung 3 byte do dai thay vi 2.
// Danh sach lay tu ServerManager.java (ham doSendMessage) va Session_ME.cs
// (ham readMessage). Hai ben phai khop tuyet doi, lech mot lenh la lech ca dong.
var lenhDoDaiBaByte = map[int8]bool{
	-32: true, -66: true, -74: true, 11: true, -67: true, -87: true, 66: true,
}

// Framer doc va ghi goi tin tren mot ket noi.
// Truoc khi goi SetKey thi moi thu di nguyen ban — do la giai doan bat tay.
//
// Framer mo hinh hoa phia MAY CHU cua duong day:
//
//	WriteMessage = dinh dang server -> client
//	ReadMessage  = dinh dang client -> server
//
// KHONG dung Framer lam "client gia" o giai doan CHUA CO KHOA. Doan do hai
// chieu dung hai thu tu byte KHAC NHAU, va do la co y - no phan chieu dung
// hanh vi cua client that (Session_ME.cs):
//   - client DOC do dai theo big-endian:  (b1 & 0xFF00) | (b2 & 0xFF)
//   - client GHI do dai theo little-endian: BinaryWriter.Write(ushort) cua .NET
//
// Sau khi co khoa thi hai chieu doi xung, luc do dung hai Framer rieng cung
// mot khoa la hop le (xem cac test san co).
type Framer struct {
	r  io.Reader
	w  io.Writer
	ks *KeyStream
}

func NewFramer(r io.Reader, w io.Writer) *Framer {
	return &Framer{r: r, w: w}
}

// SetKey bat che do ma hoa. Goi dung mot lan, ngay sau khi gui goi bat tay.
func (f *Framer) SetKey(key []byte) {
	f.ks = NewKeyStream(key)
}

func (f *Framer) coKhoa() bool { return f.ks != nil }

func (f *Framer) docMotByte() (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(f.r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

// ReadMessage doc mot goi tin do client gui len.
//
// Dinh dang (xem Session_ME.cs, ham doSendMessage):
//
//	[1 byte lenh]  [2 byte do dai]  [than tin]
//
// Khi da co khoa, do dai ghi theo thu tu (len>>8) roi (len & 0xFF), moi byte
// deu XOR. Khi chua co khoa, client dung BinaryWriter.Write(ushort) cua .NET
// — day la little-endian. Truong hop nay chi xay ra o goi bat tay -27 voi do
// dai bang 0, nen hai byte deu la 0 va thu tu khong anh huong.
func (f *Framer) ReadMessage() (*Message, error) {
	rawCmd, err := f.docMotByte()
	if err != nil {
		return nil, err
	}
	if f.coKhoa() {
		rawCmd = f.ks.Read(rawCmd)
	}
	cmd := int8(rawCmd)

	b1, err := f.docMotByte()
	if err != nil {
		return nil, err
	}
	b2, err := f.docMotByte()
	if err != nil {
		return nil, err
	}

	var n int
	if f.coKhoa() {
		hi := f.ks.Read(b1)
		lo := f.ks.Read(b2)
		n = int(hi)<<8 | int(lo)
	} else {
		n = int(binary.LittleEndian.Uint16([]byte{b1, b2}))
	}

	if n < 0 || n > doDaiThanToiDa {
		return nil, fmt.Errorf("goi tin bao do dai %d, vuot gioi han %d", n, doDaiThanToiDa)
	}

	data := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(f.r, data); err != nil {
			return nil, fmt.Errorf("doc than tin %d byte that bai: %w", n, err)
		}
		if f.coKhoa() {
			for i := range data {
				data[i] = f.ks.Read(data[i])
			}
		}
	}
	return &Message{Command: cmd, Data: data}, nil
}

// WriteMessage gui mot goi tin xuong client.
//
// Dinh dang (xem ServerManager.java, ham doSendMessage): lenh, do dai, than tin.
// Do dai dung 3 byte cho cac lenh trong lenhDoDaiBaByte, 2 byte cho phan con lai.
//
// Nguyen tac chung, ap dung nhat quan cho CA BA nhanh ben duoi: quyet dinh CO
// MA HOA hay khong CHI dua vao coKhoa(); quyet dinh 2 byte hay 3 byte CHI dua
// vao byte lenh (co nam trong lenhDoDaiBaByte hay khong). Do dai than tin (n)
// KHONG duoc tham gia vao bat ky quyet dinh nao trong so nay — dung n de gate
// la nguyen nhan cua ca hai loi da tung xay ra:
//
//  1. Nhanh lenh dac biet: 3 byte do dai la bat buoc cho MOI goi mang lenh
//     trong lenhDoDaiBaByte, KE CA khi than tin rong (n == 0). Client
//     (Session_ME.cs, readMessage) chi nhin byte lenh de quyet dinh goi
//     readMessage2 — ham nay LUON doc du 3 byte do dai, khong he xet than tin
//     dai bao nhieu. Server ma rut xuong 2 byte khi than tin rong se lam
//     client doc lem mot byte cua goi ke tiep sang truong do dai, lech dong
//     byte tu do tro di — client treo im lang, khong bao loi gi.
//
//  2. Nhanh lenh THUONG (khong nam trong lenhDoDaiBaByte) khi DA CO KHOA: dieu
//     kien "n > 0" o day la SAI — phai LUON ma hoa 2 byte do dai mien la
//     coKhoa(), bat ke than tin rong hay khong. Ba bang chung:
//     a) Client doc (Session_ME.cs, readMessage): chi xet getKeyComplete,
//     hoan toan khong xet do dai than tin — "if (getKeyComplete) { ...
//     doc 2 byte da XOR }".
//     b) Chinh client GUI goi RONG (Session_ME.cs, doSendMessage, nhanh
//     data == null) van tu ma hoa 2 byte (qua writeKey) khi da co khoa,
//     chu khong ghi tho.
//     c) Java cu (ServerManager.java): data != null va size == 0 van roi
//     vao nhanh "else if (session.sentKey())" -> ghi 2 byte CO ma hoa.
//     Neu server ghi 2 byte RAW cho goi rong luc da co khoa: client XOR chung
//     ra do dai rac, VA con tro cua KeyStream phia ghi khong duoc tien (vi
//     ks.Write khong duoc goi) trong khi con tro phia doc lai tien (vi
//     ReadMessage luon goi ks.Read cho 2 byte do dai khi coKhoa()) — hai con
//     tro lech nhau VINH VIEN tu do, hong ca cac goi sau chu khong rieng goi
//     rong nay. Ma "lenh thuong gui rong sau bat tay" la chuyen xay ra lien
//     tuc (thong bao, xac nhan).
//
// Ban Java ghi "XOR(x) - 128" con client doc "readKey(b) + 128". Tren byte 8 bit,
// cong 128 va tru 128 cho cung ket qua, nen o day dung cong cho de doc.
func (f *Framer) WriteMessage(m *Message) error {
	// LOI panic: nhanh lenh dac biet ben duoi goi f.ks.Write(...) khong dieu
	// kien. Neu chua SetKey thi f.ks == nil, cham vao se panic — va panic
	// khong duoc recover() se lam sap ca tien trinh may chu, da vang TOAN BO
	// nguoi choi dang online chu khong rieng mot ket noi. Chan tu day: chua co
	// khoa thi tu choi ro rang bang error, tuyet doi khong cham vao f.ks.
	if lenhDoDaiBaByte[m.Command] && !f.coKhoa() {
		return fmt.Errorf("lenh dac biet %d can ma hoa do dai 3 byte nhung Framer chua goi SetKey (chua co khoa)", m.Command)
	}

	out := make([]byte, 0, len(m.Data)+8)

	if f.coKhoa() {
		out = append(out, f.ks.Write(byte(m.Command)))
	} else {
		out = append(out, byte(m.Command))
	}

	n := len(m.Data)
	switch {
	case lenhDoDaiBaByte[m.Command]:
		// f.coKhoa() da duoc dam bao la true o day (kiem tra dau ham).
		// KHONG duoc them dieu kien n > 0: xem bang chung (1) o docstring.
		out = append(out, f.ks.Write(byte(n))+128)
		out = append(out, f.ks.Write(byte(n>>8))+128)
		out = append(out, f.ks.Write(byte(n>>16))+128)
	case f.coKhoa():
		// Lenh thuong, da co khoa: LUON ma hoa 2 byte do dai, bat ke n co
		// bang 0 hay khong. KHONG duoc them dieu kien n > 0: xem bang chung
		// (2a, 2b, 2c) o docstring.
		out = append(out, f.ks.Write(byte(n>>8)))
		out = append(out, f.ks.Write(byte(n&0xFF)))
	default:
		// Chua co khoa (giai doan bat tay): ghi 2 byte big-endian khong ma hoa.
		out = binary.BigEndian.AppendUint16(out, uint16(n))
	}

	if n > 0 {
		than := make([]byte, n)
		copy(than, m.Data)
		if f.coKhoa() {
			for i := range than {
				than[i] = f.ks.Write(than[i])
			}
		}
		out = append(out, than...)
	}

	_, err := f.w.Write(out)
	return err
}
