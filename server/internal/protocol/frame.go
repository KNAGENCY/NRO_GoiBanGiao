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
// Quan trong: 3 byte do dai la bat buoc cho MOI goi mang lenh trong
// lenhDoDaiBaByte, KE CA khi than tin rong (n == 0). Client (Session_ME.cs,
// ham readMessage) chi nhin byte lenh de quyet dinh goi readMessage2 — ham nay
// LUON doc du 3 byte do dai, khong he xet than tin dai bao nhieu. Server ma
// rut xuong 2 byte khi than tin rong se lam client doc lem mot byte cua goi ke
// tiep sang truong do dai, lech dong byte tu do tro di — client treo im lang,
// khong bao loi gi.
//
// Ban Java ghi "XOR(x) - 128" con client doc "readKey(b) + 128". Tren byte 8 bit,
// cong 128 va tru 128 cho cung ket qua, nen o day dung cong cho de doc.
func (f *Framer) WriteMessage(m *Message) error {
	out := make([]byte, 0, len(m.Data)+8)

	if f.coKhoa() {
		out = append(out, f.ks.Write(byte(m.Command)))
	} else {
		out = append(out, byte(m.Command))
	}

	n := len(m.Data)
	switch {
	case lenhDoDaiBaByte[m.Command]:
		// Client (Session_ME.cs, readMessage) quyet dinh doc 3 byte do dai CHI
		// dua vao byte lenh, khong xet than tin co rong hay khong (readMessage2
		// luon chay cho 7 lenh nay). Vi vay server bat buoc phai LUON ghi du 3
		// byte cho cac lenh nay, ke ca khi n == 0 — rut xuong 2 byte se lam
		// client doc lem sang goi ke tiep, lech dong byte tu do tro di va client
		// treo im lang khong bao loi gi. KHONG duoc them dieu kien n > 0 o day.
		out = append(out, f.ks.Write(byte(n))+128)
		out = append(out, f.ks.Write(byte(n>>8))+128)
		out = append(out, f.ks.Write(byte(n>>16))+128)
	case n > 0 && f.coKhoa():
		out = append(out, f.ks.Write(byte(n>>8)))
		out = append(out, f.ks.Write(byte(n&0xFF)))
	default:
		// Chua co khoa, hoac than tin rong: ghi 2 byte big-endian khong ma hoa.
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
