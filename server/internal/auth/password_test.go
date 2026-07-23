package auth

import (
	"testing"
	"time"
)

// Chuoi bam do may chu Java cu sinh ra phai doi chieu duoc.
// Gia tri nay lay tu tai khoan that trong CSDL `z`, mat khau la "1".
const bamTuBanJava = "pbkdf2$100000$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU="

func TestDoiChieuBamCuaBanJava(t *testing.T) {
	if !Verify("1", bamTuBanJava) {
		t.Fatal("mat khau dung nhung doi chieu that bai — dinh dang bam khong khop ban Java")
	}
	if Verify("2", bamTuBanJava) {
		t.Fatal("mat khau sai nhung van cho qua")
	}
}

func TestBamRoiDoiChieu(t *testing.T) {
	h, err := Hash("matkhau-thu")
	if err != nil {
		t.Fatalf("loi bam: %v", err)
	}
	if !Verify("matkhau-thu", h) {
		t.Fatal("bam roi doi chieu lai that bai")
	}
	if Verify("matkhau-khac", h) {
		t.Fatal("mat khau khac ma van cho qua")
	}
}

// Hai lan bam cung mot mat khau phai ra hai chuoi khac nhau (muoi ngau nhien).
func TestMuoiNgauNhien(t *testing.T) {
	a, _ := Hash("abc")
	b, _ := Hash("abc")
	if a == b {
		t.Fatal("hai lan bam ra ket qua giong nhau — muoi khong ngau nhien")
	}
}

// Tai khoan cu luu mat khau tho van phai dang nhap duoc.
func TestTaiKhoanCuLuuTho(t *testing.T) {
	if !Verify("123456", "123456") {
		t.Fatal("tai khoan cu luu mat khau tho phai dang nhap duoc")
	}
	if Verify("sai", "123456") {
		t.Fatal("mat khau tho sai ma van cho qua")
	}
}

func TestNhanBietDangBam(t *testing.T) {
	if !IsHashed(bamTuBanJava) {
		t.Fatal("chuoi bam phai duoc nhan biet")
	}
	if IsHashed("123456") {
		t.Fatal("mat khau tho khong phai dang bam")
	}
}

// Chuoi bam hong khong duoc lam sap may chu.
func TestChuoiBamHong(t *testing.T) {
	for _, x := range []string{"pbkdf2$", "pbkdf2$abc$x$y", "pbkdf2$100000$!!!$!!!", ""} {
		if Verify("bat ky", x) {
			t.Fatalf("chuoi hong %q ma van cho qua", x)
		}
	}
}

// Chan so vong qua lon khong duoc treo may chu.
func TestChanSoVongQuaLon(t *testing.T) {
	tests := []struct {
		name   string
		soVong string
	}{
		{"vua vuot nguong", "1000001"},
		{"2 ty vong", "2147483648"},
		{"64-bit max", "9223372036854775807"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Tao chuoi bam voi so vong qua lon
			bamMau := "pbkdf2$" + test.soVong + "$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU="

			// Do thoi gian — phai tra nhanh (duoi 1 giay)
			start := time.Now()
			result := Verify("mat khau gi", bamMau)
			elapsed := time.Since(start)

			// Phai tra false (so vong qua lon)
			if result {
				t.Errorf("so vong %s qua lon ma van cho qua", test.soVong)
			}

			// Phai tra nhanh (duoi 1 giay)
			if elapsed > time.Second {
				t.Errorf("so vong %s lam treo hang phut: %v", test.soVong, elapsed)
			}
		})
	}
}
