package antoan

import "testing"

func TestChayCoLuoi_ChanPanicVaKhongLanRa(t *testing.T) {
	// f panic -> ChayCoLuoi PHAI chan lai (test khong duoc chet) va tra true.
	daChan := ChayCoLuoi("test.panic", func() { panic("bug gia") })
	if !daChan {
		t.Fatal("panic phai bi chan va tra true")
	}
}

func TestChayCoLuoi_ChayTrotTraFalse(t *testing.T) {
	chay := false
	daChan := ChayCoLuoi("test.ok", func() { chay = true })
	if daChan {
		t.Fatal("chay trot khong duoc bao la chan panic")
	}
	if !chay {
		t.Fatal("f phai duoc chay")
	}
}
