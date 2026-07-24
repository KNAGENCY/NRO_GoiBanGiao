package protocol

import (
	"encoding/binary"
	"fmt"
)

// Writer ghi kieu du lieu theo dung cach java.io.DataOutputStream ghi,
// vi client duoc chuyen the tu ban Java ME va van doc theo quy uoc do.
// Moi so nguyen deu big-endian.
type Writer struct {
	buf []byte
}

func NewWriter() *Writer { return &Writer{buf: make([]byte, 0, 64)} }

func (w *Writer) Byte(v byte) { w.buf = append(w.buf, v) }

func (w *Writer) Bool(v bool) {
	if v {
		w.buf = append(w.buf, 1)
	} else {
		w.buf = append(w.buf, 0)
	}
}

func (w *Writer) Short(v int16) {
	w.buf = binary.BigEndian.AppendUint16(w.buf, uint16(v))
}

func (w *Writer) Int(v int32) {
	w.buf = binary.BigEndian.AppendUint32(w.buf, uint32(v))
}

func (w *Writer) Long(v int64) {
	w.buf = binary.BigEndian.AppendUint64(w.buf, uint64(v))
}

// UTF ghi chuoi theo dinh dang writeUTF cua java.io.DataOutputStream: 2 byte
// do dai (so byte SAU KHI ma hoa) roi toi noi dung ma hoa theo "modified UTF-8".
//
// Modified UTF-8 KHONG giong UTF-8 chuan cua Go, o hai diem:
//  1. Ky tu NUL (rune 0) duoc ma hoa thanh 2 byte 0xC0 0x80 thay vi 1 byte 0x00
//     (de chuoi C-style khong bi nham voi ky tu ket thuc chuoi).
//  2. Ky tu ngoai BMP (vd emoji, rune > 0xFFFF) duoc tach thanh cap surrogate
//     UTF-16 roi ma hoa MOI NUA thanh 3 byte rieng (tong 6 byte), thay vi
//     ma hoa thang thanh 1 chuoi 4 byte nhu UTF-8 chuan.
//
// Client game (chuyen the tu Java ME) doc bang java.io.DataInputStream.readUTF,
// ham nay CHI hieu modified UTF-8. Neu ghi emoji theo UTF-8 chuan 4 byte,
// readUTF phia Java se nem UTFDataFormatException (goi tin bi coi la hong),
// nen bat buoc phai tu ma hoa dung quy uoc nay, khong the ghi thang []byte(s).
func (w *Writer) UTF(s string) error {
	encoded := encodeModifiedUTF8(s)
	if len(encoded) > 65535 {
		return fmt.Errorf("chuoi sau khi ma hoa modified UTF-8 dai %d byte, vuot gioi han 65535 cua writeUTF", len(encoded))
	}
	w.Short(int16(len(encoded)))
	w.buf = append(w.buf, encoded...)
	return nil
}

// encodeModifiedUTF8 ma hoa 1 chuoi Go (UTF-8 chuan) thanh modified UTF-8
// dung quy uoc cua java.io.DataOutputStream.writeUTF.
func encodeModifiedUTF8(s string) []byte {
	buf := make([]byte, 0, len(s)+8)
	for _, r := range s {
		switch {
		case r == 0:
			buf = append(buf, 0xC0, 0x80)
		case r >= 1 && r <= 0x7F:
			buf = append(buf, byte(r))
		case r >= 0x80 && r <= 0x7FF:
			buf = append(buf,
				byte(0xC0|(r>>6)),
				byte(0x80|(r&0x3F)),
			)
		case r >= 0x800 && r <= 0xFFFF:
			buf = append(buf,
				byte(0xE0|(r>>12)),
				byte(0x80|((r>>6)&0x3F)),
				byte(0x80|(r&0x3F)),
			)
		default:
			// Ngoai BMP: tach thanh cap surrogate UTF-16 roi ma hoa moi nua
			// thanh 3 byte, dung nhu java.io.DataOutputStream.writeUTF.
			v := r - 0x10000
			hi := rune(0xD800 + (v >> 10))
			lo := rune(0xDC00 + (v & 0x3FF))
			buf = append(buf,
				byte(0xE0|(hi>>12)),
				byte(0x80|((hi>>6)&0x3F)),
				byte(0x80|(hi&0x3F)),
				byte(0xE0|(lo>>12)),
				byte(0x80|((lo>>6)&0x3F)),
				byte(0x80|(lo&0x3F)),
			)
		}
	}
	return buf
}

// decodeModifiedUTF8 giai ma nguoc lai encodeModifiedUTF8: doc tung ky tu
// modified UTF-8 (1-3 byte), sau do ghep lai cac cap surrogate UTF-16
// (0xD800-0xDBFF theo sau 0xDC00-0xDFFF) thanh 1 rune that (Go string van
// la UTF-8 chuan nen phai ghep truoc khi tra ve).
func decodeModifiedUTF8(b []byte) (string, error) {
	units := make([]rune, 0, len(b))
	i := 0
	for i < len(b) {
		b0 := b[i]
		switch {
		case b0&0x80 == 0:
			units = append(units, rune(b0))
			i++
		case b0&0xE0 == 0xC0:
			if i+1 >= len(b) {
				return "", fmt.Errorf("modified UTF-8 hong: thieu byte thu 2 tai vi tri %d", i)
			}
			b1 := b[i+1]
			if b1&0xC0 != 0x80 {
				return "", fmt.Errorf("modified UTF-8 hong: byte lien theo khong hop le tai vi tri %d", i+1)
			}
			units = append(units, rune(b0&0x1F)<<6|rune(b1&0x3F))
			i += 2
		case b0&0xF0 == 0xE0:
			if i+2 >= len(b) {
				return "", fmt.Errorf("modified UTF-8 hong: thieu byte cho chuoi 3 byte tai vi tri %d", i)
			}
			b1, b2 := b[i+1], b[i+2]
			if b1&0xC0 != 0x80 || b2&0xC0 != 0x80 {
				return "", fmt.Errorf("modified UTF-8 hong: byte lien theo khong hop le tai vi tri %d", i+1)
			}
			units = append(units, rune(b0&0x0F)<<12|rune(b1&0x3F)<<6|rune(b2&0x3F))
			i += 3
		default:
			return "", fmt.Errorf("modified UTF-8 hong: byte dau 0x%02X khong hop le tai vi tri %d", b0, i)
		}
	}

	runes := make([]rune, 0, len(units))
	for idx := 0; idx < len(units); idx++ {
		u := units[idx]
		if u >= 0xD800 && u <= 0xDBFF && idx+1 < len(units) {
			next := units[idx+1]
			if next >= 0xDC00 && next <= 0xDFFF {
				combined := 0x10000 + (u-0xD800)<<10 + (next - 0xDC00)
				runes = append(runes, combined)
				idx++
				continue
			}
		}
		runes = append(runes, u)
	}
	return string(runes), nil
}

// ThemMang ghi ca khoi byte THO vao cuoi bo dem, KHONG kem tien to do dai (khac
// UTF). Dung cho cac blob nhi phan da co san (vi du du lieu part 87KB trong goi
// -87), nơi do dai da duoc ghi rieng bang Int() truoc do.
func (w *Writer) ThemMang(p []byte) {
	w.buf = append(w.buf, p...)
}

func (w *Writer) Bytes() []byte { return w.buf }

// Reader doc nguoc lai nhung gi Writer ghi.
type Reader struct {
	buf []byte
	pos int
}

func NewReader(b []byte) *Reader { return &Reader{buf: b} }

func (r *Reader) con(n int) ([]byte, error) {
	if r.pos+n > len(r.buf) {
		return nil, fmt.Errorf("doc %d byte tai vi tri %d nhung bo dem chi co %d byte", n, r.pos, len(r.buf))
	}
	out := r.buf[r.pos : r.pos+n]
	r.pos += n
	return out, nil
}

func (r *Reader) Byte() (byte, error) {
	b, err := r.con(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *Reader) Short() (int16, error) {
	b, err := r.con(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b)), nil
}

func (r *Reader) Int() (int32, error) {
	b, err := r.con(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

// Long doc 8 byte big-endian thanh int64 — nghich dao cua Writer.Long, khop
// java.io.DataInputStream.readLong. Dung cho cac truong long trong ME_INFO
// (power, gold): client luon readLong() nen ta ghi/doc dung 8 byte.
func (r *Reader) Long() (int64, error) {
	b, err := r.con(8)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

// UTF doc chuoi theo dinh dang readUTF cua java.io.DataInputStream: 2 byte
// do dai (so byte modified UTF-8 tiep theo) roi giai ma modified UTF-8 thanh
// chuoi Go (UTF-8 chuan). Xem chu thich cua Writer.UTF de biet ly do phai
// giai ma theo quy uoc rieng nay thay vi coi noi dung la UTF-8 chuan.
func (r *Reader) UTF() (string, error) {
	n, err := r.Short()
	if err != nil {
		return "", err
	}
	b, err := r.con(int(uint16(n)))
	if err != nil {
		return "", err
	}
	return decodeModifiedUTF8(b)
}
