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

// UTF ghi chuoi theo dinh dang writeUTF cua Java: 2 byte do dai roi noi dung.
// Go luu chuoi san o dang UTF-8 nen khong can chuyen doi, chi can kiem do dai.
func (w *Writer) UTF(s string) error {
	if len(s) > 65535 {
		return fmt.Errorf("chuoi dai %d byte, vuot gioi han 65535 cua writeUTF", len(s))
	}
	w.Short(int16(len(s)))
	w.buf = append(w.buf, s...)
	return nil
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

func (r *Reader) UTF() (string, error) {
	n, err := r.Short()
	if err != nil {
		return "", err
	}
	b, err := r.con(int(uint16(n)))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
