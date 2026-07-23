// Package auth lo viec bam va doi chieu mat khau.
//
// Dinh dang PHAI khop tuyet doi voi ban Java cu (auth/PasswordUtil.java) va
// ban PHP (sv1/config.php). Sua mot ben ma quen ben kia se khien toan bo
// nguoi choi khong dang nhap duoc.
package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	tienTo      = "pbkdf2$"
	soVong      = 100000
	soVongToiDa = 1_000_000
	soByteMu    = 16
	soBitKey    = 256
)

// soVongToiDa chan tren so vong lap doc tu chuoi bam.
//
// Vi sao can: neu cot mat khau trong CSDL bi hong hoac bi chen sai mot so
// vong khong lo (vd 2 ty), Verify se chay that va treo hang phut - mot luong
// xu ly dang nhap bi treo cung nguy hiem ngang sap. Ban Java tinh co an toan
// nho so 32 bit bi tran, con Go dung so 64 bit nen khong tran.
//
// Muc 1 trieu la gap 10 lan gia tri dang dung (100.000), du du dat de sau nay
// tang so vong theo phan cung moi, ma van chan duoc truong hop du lieu hong.

// IsHashed cho biet chuoi luu trong CSDL da o dang bam hay con la mat khau tho
// cua he thong cu.
func IsHashed(luuTru string) bool {
	return strings.HasPrefix(luuTru, tienTo)
}

// Hash bam mat khau moi voi muoi ngau nhien.
func Hash(matKhau string) (string, error) {
	muoi := make([]byte, soByteMu)
	if _, err := rand.Read(muoi); err != nil {
		return "", fmt.Errorf("khong sinh duoc muoi: %w", err)
	}
	// Chu y chu ky ham cua thu vien CHUAN khac voi golang.org/x/crypto:
	//   pbkdf2.Key(ham_bam, mat_khau_string, muoi, so_vong, so_byte) (bytes, error)
	bam, err := pbkdf2.Key(sha256.New, matKhau, muoi, soVong, soBitKey/8)
	if err != nil {
		return "", fmt.Errorf("bam that bai: %w", err)
	}
	return fmt.Sprintf("%s%d$%s$%s",
		tienTo, soVong,
		base64.StdEncoding.EncodeToString(muoi),
		base64.StdEncoding.EncodeToString(bam)), nil
}

// Verify doi chieu mat khau nguoi dung nhap voi gia tri luu trong CSDL.
// Ho tro ca tai khoan cu con luu mat khau tho — nhung van so sanh bang
// thuat toan hang thoi gian de khong lo thong tin qua do tre.
func Verify(matKhau, luuTru string) bool {
	if luuTru == "" {
		return false
	}
	if !IsHashed(luuTru) {
		return subtle.ConstantTimeCompare([]byte(luuTru), []byte(matKhau)) == 1
	}

	phan := strings.Split(luuTru, "$")
	if len(phan) != 4 {
		return false
	}
	vong, err := strconv.Atoi(phan[1])
	if err != nil || vong < 1 || vong > soVongToiDa {
		return false
	}
	muoi, err := base64.StdEncoding.DecodeString(phan[2])
	if err != nil {
		return false
	}
	muon, err := base64.StdEncoding.DecodeString(phan[3])
	if err != nil || len(muon) == 0 {
		return false
	}

	duoc, err := pbkdf2.Key(sha256.New, matKhau, muoi, vong, len(muon))
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(muon, duoc) == 1
}
