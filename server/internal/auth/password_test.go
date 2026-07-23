package auth

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"regexp"
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

// ---------------------------------------------------------------------------
// PHAN B, Test 2: dinh dang chuoi bam phai khop CHINH XAC ban Java
// (PasswordUtil.java) va ban PHP (sv1/config.php). Kiem thu dot bien da chi ra
// BON dot bien deu SONG SOT (khong test nao truoc day bat duoc):
//   1. So vong 100.000 -> 10.000
//   2. Muoi 16 -> 4 byte
//   3. Khoa 256 -> 64 bit
//   4. Bo doc so vong tu chuoi luu tru (dung hang so cung thay vi phan[1])
// Day la noi rang buoc du an ghi ro "phai khop ban Java va ban PHP" — sai
// lech BAT KY dieu nao trong bon dieu tren deu lam TOAN BO nguoi choi mat
// dang nhap (chuoi bam moi khong con khop cac ban con lai).
// ---------------------------------------------------------------------------

// reDinhDangBam khop dung dinh dang Hash() phai sinh ra:
// pbkdf2$<so_vong>$<muoi base64>$<bam base64>
var reDinhDangBam = regexp.MustCompile(`^pbkdf2\$(\d+)\$([A-Za-z0-9+/]+=*)\$([A-Za-z0-9+/]+=*)$`)

// TestHashDungDinhDangKhopJavaVaPhp bat dung ca 3 dot bien dau (so vong,
// do dai muoi, do dai khoa/bam) bang cach kiem tung phan cua chuoi Hash() sinh
// ra, khong chi kiem "Verify(Hash(x)) == true" (kieu kiem do se KHONG bat
// duoc neu ca hai ham cung doi lech nhau MOT CACH DONG BO).
func TestHashDungDinhDangKhopJavaVaPhp(t *testing.T) {
	h, err := Hash("gi-do")
	if err != nil {
		t.Fatalf("loi bam: %v", err)
	}

	phan := reDinhDangBam.FindStringSubmatch(h)
	if phan == nil {
		t.Fatalf("chuoi bam %q khong khop dinh dang pbkdf2$<so_vong>$<muoi>$<bam>", h)
	}

	if phan[1] != "100000" {
		t.Fatalf("so vong = %s, mong doi 100000 (phai khop ban Java/PHP)", phan[1])
	}

	muoi, err := base64.StdEncoding.DecodeString(phan[2])
	if err != nil {
		t.Fatalf("khong giai ma duoc phan muoi (base64): %v", err)
	}
	if len(muoi) != 16 {
		t.Fatalf("muoi dai %d byte, mong doi 16 byte", len(muoi))
	}

	bamRa, err := base64.StdEncoding.DecodeString(phan[3])
	if err != nil {
		t.Fatalf("khong giai ma duoc phan bam (base64): %v", err)
	}
	if len(bamRa) != 32 {
		t.Fatalf("bam dai %d byte, mong doi 32 byte (256 bit)", len(bamRa))
	}
}

// TestVerifyDocSoVongTuChuoiLuuTru bat dung dot bien thu tu: "bo doc so vong
// tu chuoi luu tru" (dung hang so soVong=100000 cung thay vi phan[1] doc
// duoc). Dot bien nay song sot truoc day vi khong test nao tung doi chieu mot
// chuoi bam co so vong KHAC 100000.
//
// Cach kiem: tu dung mot chuoi bam voi so vong = 1000 bang CHINH crypto/pbkdf2
// (khong qua Hash() cua package nay, vi Hash() luon dung cung soVong=100000),
// roi khang dinh Verify van doi chieu dung — chi co the dung neu Verify THAT
// SU doc so vong tu chuoi luu tru thay vi dung hang so cung.
func TestVerifyDocSoVongTuChuoiLuuTru(t *testing.T) {
	const soVongKhac = 1000
	muoi := []byte("0123456789abcdef") // 16 byte, gia lap muoi that
	matKhau := "mat-khau-thu-nghiem-so-vong-khac"

	bamTay, err := pbkdf2.Key(sha256.New, matKhau, muoi, soVongKhac, 32)
	if err != nil {
		t.Fatalf("loi bam thu cong: %v", err)
	}

	luuTru := fmt.Sprintf("%s%d$%s$%s",
		tienTo, soVongKhac,
		base64.StdEncoding.EncodeToString(muoi),
		base64.StdEncoding.EncodeToString(bamTay))

	if !Verify(matKhau, luuTru) {
		t.Fatal("Verify khong doi chieu duoc chuoi bam co so vong KHAC 100000 " +
			"— co the Verify dang dung hang so cung (soVong) thay vi doc so vong tu chuoi luu tru")
	}
	if Verify("mat-khau-sai", luuTru) {
		t.Fatal("mat khau sai ma van cho qua")
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
