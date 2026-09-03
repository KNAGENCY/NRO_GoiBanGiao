// kich_test.go — kiem kenh nap nong (KH8 task 10).
//
// Ba dieu duoc khoa o day, moi dieu la mot dot bien co that neu bo:
//  1. dia chi KHONG loopback phai bi TU CHOI (§4.3 rui ro 1);
//  2. nap that bai phai GIU NGUYEN ban dang chay (fail-closed, §4);
//  3. hai request song song khong duoc lam so the he nhay/trung.
package noidung

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"nroserver/kho-noi-dung"
)

const khoaTest = "khoa-test-123"

// thuMucYAMLTest do ban NHUNG ra mot thu muc tam de test sua duoc mot con so.
func thuMucYAMLTest(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	for _, ten := range TenCacFile() { // nguon duy nhat — xem chu thich TenCacFile
		b, err := fs.ReadFile(khonoidung.FS, ten)
		if err != nil {
			t.Fatalf("khong doc duoc %s: %v", ten, err)
		}
		if err := os.WriteFile(filepath.Join(d, ten), b, 0o600); err != nil {
			t.Fatalf("khong ghi duoc %s: %v", ten, err)
		}
	}
	return d
}

func doiTrongFile(t *testing.T, thuMuc, ten, cu, moi string) {
	t.Helper()
	p := filepath.Join(thuMuc, ten)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("doc %s: %v", ten, err)
	}
	if !strings.Contains(string(b), cu) {
		t.Fatalf("khong tim thay %q trong %s — dot bien khong va vao dau", cu, ten)
	}
	moiVan := strings.Replace(string(b), cu, moi, 1)
	if err := os.WriteFile(p, []byte(moiVan), 0o600); err != nil {
		t.Fatalf("ghi %s: %v", ten, err)
	}
}

func khoTest(t *testing.T) *Kho {
	t.Helper()
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("NapNhung: %v", err)
	}
	return MoKho(du)
}

func logCam() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func moKenhTest(t *testing.T, kho *Kho, thuMuc string) *KenhNapLai {
	t.Helper()
	k, err := MoKenhNapLai("127.0.0.1:0", kho, thuMuc, khoaTest, logCam())
	if err != nil {
		t.Fatalf("MoKenhNapLai: %v", err)
	}
	t.Cleanup(func() {
		ctx, huy := context.WithTimeout(context.Background(), 3*time.Second)
		defer huy()
		_ = k.Dung(ctx)
	})
	return k
}

func goi(t *testing.T, k *KenhNapLai, phuongThuc, duong, khoa string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(phuongThuc, "http://"+k.DiaChi()+duong, nil)
	if err != nil {
		t.Fatalf("tao request: %v", err)
	}
	if khoa != "" {
		req.Header.Set(TenHeaderKhoaNapLai, khoa)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("goi %s %s: %v", phuongThuc, duong, err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

// ⚠️ DOT BIEN 1 — §4.3 rui ro 1: doi dieu kien bind thanh "chap nhan het" (bo
// KiemDiaChiNapLai trong MoKenhNapLai) thi test nay DO.
func TestKenhNapLaiTuChoiDiaChiKhongLoopback(t *testing.T) {
	xau := []string{"0.0.0.0:9090", ":9090", "192.168.1.10:9090", "8.8.8.8:9090"}
	for _, dc := range xau {
		if err := KiemDiaChiNapLai(dc); err == nil {
			t.Fatalf("dia chi %q PHAI bi tu choi — kenh nap lai chi duoc bind loopback", dc)
		}
		k, err := MoKenhNapLai(dc, khoTest(t), "", khoaTest, logCam())
		if err == nil {
			ctx, huy := context.WithTimeout(context.Background(), time.Second)
			_ = k.Dung(ctx)
			huy()
			t.Fatalf("MoKenhNapLai(%q) PHAI loi — da mo mot cong ra ngoai", dc)
		}
	}
	for _, tot := range []string{"127.0.0.1:9090", "[::1]:9090", "localhost:9090", "127.0.0.5:1"} {
		if err := KiemDiaChiNapLai(tot); err != nil {
			t.Fatalf("dia chi loopback %q bi tu choi nham: %v", tot, err)
		}
	}
}

// Thieu khoa dung chung = KHONG mo kenh (fail-closed, §4.3 rui ro 2).
func TestKenhNapLaiTuChoiKhiThieuKhoa(t *testing.T) {
	if _, err := MoKenhNapLai("127.0.0.1:0", khoTest(t), "", "  ", logCam()); err == nil {
		t.Fatal("MoKenhNapLai PHAI loi khi khoa dung chung rong")
	}
}

func TestKenhNapLaiSaiKhoaTraVe403(t *testing.T) {
	k := moKenhTest(t, khoTest(t), thuMucYAMLTest(t))
	if ma, _ := goi(t, k, http.MethodPost, DuongNapLai, "sai-khoa"); ma != http.StatusForbidden {
		t.Fatalf("sai khoa phai tra 403, nhan %d", ma)
	}
	if ma, _ := goi(t, k, http.MethodPost, DuongNapLai, ""); ma != http.StatusForbidden {
		t.Fatalf("thieu khoa phai tra 403, nhan %d", ma)
	}
	if ma, _ := goi(t, k, http.MethodGet, DuongNapLai, khoaTest); ma != http.StatusMethodNotAllowed {
		t.Fatalf("GET vao duong nap lai phai tra 405, nhan %d", ma)
	}
}

func TestKenhNapLaiThanhCongTangTheHe(t *testing.T) {
	kho := khoTest(t)
	truoc := kho.TheHe()
	k := moKenhTest(t, kho, thuMucYAMLTest(t))

	ma, than := goi(t, k, http.MethodPost, DuongNapLai, khoaTest)
	if ma != http.StatusOK {
		t.Fatalf("nap lai hop le phai tra 200, nhan %d: %s", ma, than)
	}
	if kho.TheHe() != truoc+1 {
		t.Fatalf("the he phai tang 1 (%d -> %d), nhan %d", truoc, truoc+1, kho.TheHe())
	}
	for _, can := range []string{"the_he", "tom_tat", "so_chieu"} {
		if !strings.Contains(than, can) {
			t.Fatalf("than tra loi thieu %q: %s", can, than)
		}
	}
	ma, than = goi(t, k, http.MethodGet, DuongTheHe, khoaTest)
	if ma != http.StatusOK || !strings.Contains(than, "the_he") {
		t.Fatalf("GET the-he sai: ma=%d than=%s", ma, than)
	}
}

// Thu muc rong = chi co ban nhung -> 501 kem loi giai thich, KHONG phai 500.
func TestKenhNapLaiKhongCoThuMucTraVe501(t *testing.T) {
	k := moKenhTest(t, khoTest(t), "")
	ma, than := goi(t, k, http.MethodPost, DuongNapLai, khoaTest)
	if ma != http.StatusNotImplemented {
		t.Fatalf("khong co thu muc YAML phai tra 501, nhan %d", ma)
	}
	if !strings.Contains(than, "NRO_KHO_NOI_DUNG") {
		t.Fatalf("501 phai chi ro cach sua (NRO_KHO_NOI_DUNG): %s", than)
	}
}

// ⚠️ DOT BIEN 2 — fail-closed (§4): doi xuLyNapLai thanh "Thay(du) roi moi kiem
// loi" (hoac bo nhanh `if err != nil { return }`) thi test nay DO: the he se
// tang du bo so lieu hong.
func TestKenhNapLaiLoiValidatorGiuNguyenBanCu(t *testing.T) {
	kho := khoTest(t)
	truoc := kho.TheHe()
	acTruoc := kho.AnhChup()
	ds := acTruoc.ChieuCuaLang(0)
	if len(ds) == 0 {
		t.Fatal("lang 0 khong co chieu nao — ban goc da hong, test khong chung minh gi")
	}

	thuMuc := thuMucYAMLTest(t)
	// Mot dot bien ma bo kiem PHAI bat: khoi lang thu 0 khai id_lang sai.
	doiTrongFile(t, thuMuc, TenFileChieu, "id_lang: 0", "id_lang: 3")

	k := moKenhTest(t, kho, thuMuc)
	ma, than := goi(t, k, http.MethodPost, DuongNapLai, khoaTest)
	if ma != http.StatusUnprocessableEntity {
		t.Fatalf("nap lai hong phai tra 422, nhan %d: %s", ma, than)
	}
	if !strings.Contains(than, "id_lang") {
		t.Fatalf("than tra loi phai chua THONG DIEP NGUYEN VAN cua validator: %s", than)
	}
	if kho.TheHe() != truoc {
		t.Fatalf("nap that bai ma the he van doi (%d -> %d) — KHONG fail-closed", truoc, kho.TheHe())
	}
	if len(kho.AnhChup().ChieuCuaLang(0)) != len(ds) {
		t.Fatal("anh chup dang chay bi doi sau mot lan nap that bai")
	}
}

// ⚠️ DOT BIEN 3 — bo k.muNap trong xuLyNapLai. Xem ghi chu duoi day: day la LOP
// THU NHAT; Kho.Thay (CompareAndSwap) van la lop thu hai, nen bo mutex chua
// chac lam so the he trung. Cai test nay khoa la: N request song song PHAI cho
// dung N the he moi, khong hon khong kem, va bo so lieu cuoi cung phai nap
// duoc — tuc khong co lan Thay nao roi vao du lieu doc do dang.
func TestNapLaiHaiRequestSongSong(t *testing.T) {
	kho := khoTest(t)
	truoc := kho.TheHe()
	k := moKenhTest(t, kho, thuMucYAMLTest(t))

	const soRequest = 8
	var wg sync.WaitGroup
	maCode := make([]int, soRequest)
	for i := 0; i < soRequest; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, "http://"+k.DiaChi()+DuongNapLai, nil)
			req.Header.Set(TenHeaderKhoaNapLai, khoaTest)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				maCode[i] = -1
				return
			}
			_, _ = io.Copy(io.Discard, res.Body)
			res.Body.Close()
			maCode[i] = res.StatusCode
		}(i)
	}
	wg.Wait()
	for i, ma := range maCode {
		if ma != http.StatusOK {
			t.Fatalf("request %d tra %d, mong 200", i, ma)
		}
	}
	if got := kho.TheHe(); got != truoc+soRequest {
		t.Fatalf("%d lan nap thanh cong phai cho dung %d the he moi: %d -> %d",
			soRequest, soRequest, truoc, got)
	}
}
