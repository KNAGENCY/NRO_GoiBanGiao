package auth

import (
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// CanBangThoiGian sinh ra de bit lo hong "bam gio la biet tai khoan nao co
// that" (xem internal/game/rothoigian_test.go). No chi lam DUNG mot viec: ton
// dung bang khoi luong tinh toan cua mot lan doi chieu mat khau that. Ba test
// duoi day khoa ba cach no co the am tham hong.

// Neu ai do doi so vong bam (soVong) ma quen sinh lai chuoi bam gia, chuoi gia
// se re hon chuoi that va lo hong quay lai NGAY — khong co trieu chung nao
// khac ngoai viec tai khoan lai bi liet ke duoc. Test nay bat dung truong hop do.
func TestBamGiaDungSoVongVoiBamThat(t *testing.T) {
	phan := strings.Split(bamGia, "$")
	if len(phan) != 4 {
		t.Fatalf("chuoi bam gia sai dinh dang: %q", bamGia)
	}
	vong, err := strconv.Atoi(phan[1])
	if err != nil {
		t.Fatalf("khong doc duoc so vong trong chuoi bam gia %q: %v", bamGia, err)
	}
	if vong != soVong {
		t.Fatalf("chuoi bam gia khai %d vong nhung ham bam that dang dung %d vong. "+
			"Bam gia re hon bam that thi hai nhanh dang nhap lai chenh nhau ve thoi gian "+
			"va ke xau liet ke duoc tai khoan nao co that. Cach sua: sinh lai hang so bamGia "+
			"bang auth.Hash() roi dan vao password.go.", vong, soVong)
	}
	if !IsHashed(bamGia) {
		t.Fatalf("chuoi bam gia %q khong o dang da bam — Verify se roi sang nhanh so sanh "+
			"mat khau tho, nhanh do gan nhu khong ton thoi gian nao", bamGia)
	}
}

// Chuoi bam gia khong duoc tinh co khop voi mot mat khau nao — neu khop thi
// mot nguoi go dung chuoi do se "dang nhap thanh cong" vao mot tai khoan
// khong ton tai.
func TestBamGiaKhongKhopMatKhauNao(t *testing.T) {
	for _, mk := range []string{"", "123456", "matkhau", "password", "admin", bamGia} {
		if Verify(mk, bamGia) {
			t.Fatalf("mat khau %q khop voi chuoi bam gia", mk)
		}
	}
	if CanBangThoiGian("matkhau") {
		t.Fatal("CanBangThoiGian tra ve true — no phai LUON tra false")
	}
}

func trungViThoiGian(d []time.Duration) time.Duration {
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	return d[len(d)/2]
}

// Test dot bien quan trong nhat cua goi nay: neu ai do "toi uu" CanBangThoiGian
// thanh mot ham rong, hoac trinh bien dich cat bo loi goi vi ket qua khong
// duoc dung, thi ham van chay dung nhung khong con ton thoi gian — va lo hong
// quay lai trong im lang. Do bang so la cach duy nhat bat duoc chuyen do.
func TestCanBangThoiGianTonBangMotLanDoiChieuThat(t *testing.T) {
	bamThat, err := Hash("matkhaunguoidung")
	if err != nil {
		t.Fatalf("khong bam duoc mat khau mau: %v", err)
	}

	const soLan = 40
	that := make([]time.Duration, 0, soLan)
	gia := make([]time.Duration, 0, soLan)
	// Xen ke hai phep do de moi thay doi cham cua may (CPU ha xung nhip, bo
	// thu gom rac chay mot dot) roi deu len ca hai ben thay vi don het vao mot ben.
	for i := 0; i < soLan; i++ {
		t0 := time.Now()
		Verify("saibet", bamThat)
		that = append(that, time.Since(t0))

		t1 := time.Now()
		CanBangThoiGian("saibet")
		gia = append(gia, time.Since(t1))
	}

	medThat := trungViThoiGian(that)
	medGia := trungViThoiGian(gia)
	t.Logf("doi chieu that = %v | can bang thoi gian = %v", medThat, medGia)

	// Doi hoi hai ben nam trong khoang 25% cua nhau. Rong rai du de khong chop
	// tat vi nhieu do tren may Windows dang lam viec khac, nhung chat hon rat
	// nhieu so voi kieu hong that (ham rong -> gia gan bang 0, tuc lech 100%).
	if medGia*4 < medThat*3 || medGia*3 > medThat*4 {
		t.Fatalf("can bang thoi gian ton %v trong khi mot lan doi chieu that ton %v — "+
			"hai ben phai xap xi nhau, neu khong thi nhanh 'khong co tai khoan' van nhanh hon "+
			"nhanh 'sai mat khau' va ke xau bam gio la liet ke duoc tai khoan", medGia, medThat)
	}
}
