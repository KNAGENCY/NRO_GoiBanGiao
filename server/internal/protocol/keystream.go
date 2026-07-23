// Package protocol cai dat giao thuc nhi phan cua client Unity cu.
//
// Toan bo phan "noi chuyen bang byte" nam gon trong package nay. Phan con lai
// cua may chu khong duoc biet gi ve khoa ma hay thu tu byte — nho vay sau nay
// them client moi chi can them mot package tuong tu, khong dung vao luat choi.
package protocol

import (
	"crypto/rand"
	"fmt"
)

// doDaiKhoaToiDa gioi han o 127 vi client dung kieu sbyte lam con tro khoa
// (xem Session_ME.cs, ham readKey/writeKey). Khoa dai hon se lam tran so
// ben phia client va giai ma sai — loi im lang, rat kho tim.
const doDaiKhoaToiDa = 127

// GenerateRawKey sinh khoa tho ngau nhien de gui cho client trong goi bat tay.
func GenerateRawKey(n int) ([]byte, error) {
	if n < 1 || n > doDaiKhoaToiDa {
		return nil, fmt.Errorf("do dai khoa phai trong khoang 1..%d, nhan duoc %d", doDaiKhoaToiDa, n)
	}
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("khong sinh duoc khoa ngau nhien: %w", err)
	}
	return raw, nil
}

// DeriveKey ap dung phep bien doi ma client thuc hien sau khi nhan khoa tho:
//
//	for j := 0; j < len-1; j++ { key[j+1] ^= key[j] }
//
// Hai ben phai ra cung ket qua, neu khong moi byte sau do deu sai.
// Ham tra ve mang moi, khong sua mang dau vao.
func DeriveKey(raw []byte) []byte {
	key := make([]byte, len(raw))
	copy(key, raw)
	for j := 0; j+1 < len(key); j++ {
		key[j+1] ^= key[j]
	}
	return key
}

// KeyStream giu hai con tro doc lap: mot cho chieu doc, mot cho chieu ghi.
// Dung chung mot con tro cho ca hai chieu se lech ngay goi tin dau tien.
type KeyStream struct {
	key  []byte
	curR int
	curW int
}

func NewKeyStream(key []byte) *KeyStream {
	return &KeyStream{key: key}
}

// Read giai ma mot byte nhan tu client.
func (k *KeyStream) Read(b byte) byte {
	out := k.key[k.curR] ^ b
	k.curR++
	if k.curR >= len(k.key) {
		k.curR %= len(k.key)
	}
	return out
}

// Write ma hoa mot byte truoc khi gui cho client.
func (k *KeyStream) Write(b byte) byte {
	out := k.key[k.curW] ^ b
	k.curW++
	if k.curW >= len(k.key) {
		k.curW %= len(k.key)
	}
	return out
}
