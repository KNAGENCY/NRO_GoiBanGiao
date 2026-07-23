package protocol

// Message la mot goi tin tren day: mot byte lenh va phan than tin.
// Than tin co the rong (nil hoac do dai 0).
type Message struct {
	Command int8
	Data    []byte
}

// NewMessage tao goi tin voi than tin duoc dung san bang Writer.
func NewMessage(cmd int8, w *Writer) *Message {
	return &Message{Command: cmd, Data: w.Bytes()}
}
