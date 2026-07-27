package game

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"nroserver/internal/store"
)

// khoLuuGia gia lap KhoLuu de kiem thu ThoGhi ma KHONG can CSDL that (dac ta
// Task 3: "Task nay dung kho GIA nen KHONG can CSDL"). Day la mot type RIENG
// voi khoGia (login_test.go): KhoLuu chi co MOT ham CapNhatNhanVat, khac han
// 4 ham cua Kho — gop chung se phai gia lap thua 3 ham khong lien quan.
type khoLuuGia struct {
	mu sync.Mutex

	goiGhi   []store.NhanVat // MOI lan CapNhatNhanVat duoc goi, THEO DUNG THU TU goi toi
	coHanGio []bool          // ctx.Deadline() co ton tai khong tai lan goi do — CUNG chi so voi goiGhi

	doTre time.Duration // >0: sleep NGAN NAY truoc khi tra ve — mo phong kho/dia cham (test b)

	// soLoiDauTien: so LAN GOI DAU TIEN se tra loi; tu lan goi thu (soLoiDauTien+1)
	// tro di LUON thanh cong. Dung cho test (c): dat =1 de mo phong "kho hong
	// mot lan roi tu hoi phuc" — dung tinh huong CSDL chap chon thoang qua.
	soLoiDauTien int
	daGoi        int // dem so lan da goi — CHI dung/sua duoi mu.Lock(), khong can atomic rieng
}

func (k *khoLuuGia) CapNhatNhanVat(ctx context.Context, nv store.NhanVat) error {
	// Sleep TRUOC khi dung k.mu — kho that (PostgreSQL) cham o tang mang/dia,
	// khong phai o tang khoa noi bo cua chinh no; mo phong dung vi tri do giup
	// test (b) đo đúng cái test muốn đo: "Xep goi kho cham" chu khong phai
	// "hai goroutine tranh mu.Lock() cua kho gia".
	if k.doTre > 0 {
		time.Sleep(k.doTre)
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.daGoi++
	_, coHanGio := ctx.Deadline()
	k.goiGhi = append(k.goiGhi, nv)
	k.coHanGio = append(k.coHanGio, coHanGio)
	if k.daGoi <= k.soLoiDauTien {
		return errors.New("khoLuuGia: loi mo phong (chua toi luot thanh cong)")
	}
	return nil
}

func (k *khoLuuGia) soLanGoi() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return len(k.goiGhi)
}

// dongVaChoXongVoiHan la tro giup dung chung cho nhieu test: goi
// DongVaChoXong voi mot han gio hop ly (khong de test treo vo han neu code
// co loi thuc su lam ThoGhi khong bao gio dong duoc).
func dongVaChoXongVoiHan(t *testing.T, tg *ThoGhi, han time.Duration) {
	t.Helper()
	ctx, huy := context.WithTimeout(context.Background(), han)
	defer huy()
	if err := tg.DongVaChoXong(ctx); err != nil {
		t.Fatalf("DongVaChoXong loi: %v", err)
	}
}

// TestThoGhiXepGopBanMoiNhatCungMotNguoi la test (a) cua brief Task 3: xep 3
// ban CUNG MOT nguoi (cung ID) lien tiep -> hang doi phai GOP lai, ban GIUA
// (nv2) khong bao gio duoc ghi RIENG xuong kho, va ban CUOI CUNG duoc ghi
// phai la ban MOI NHAT (nv3).
//
// VI SAO khong assert "dung 1 lan goi": neu tho ghi (goroutine nen) tinh co
// kip lay ban nv1 ra khoi hang doi TRUOC khi Xep(nv2)/Xep(nv3) chay xong (mot
// kha nang co that, du cuc hiem, tren may nhieu loi CPU that su song song),
// ket qua VAN DUNG kien truc gop — chi la co THEM mot lan ghi nv1 (du lieu cu
// nhung khong SAI, chi la "chua kip moi") truoc lan ghi nv3. Dieu KHONG DUOC
// PHEP xay ra, du trong bat ky tinh huong xap lich nao, la ghi RIENG ban
// GIUA (nv2) — day la bat bien that su cua "gop theo map", nen la dieu duy
// nhat test nay ghim cung, tranh flaky ma van bat dung dot bien §Step 5.
func TestThoGhiXepGopBanMoiNhatCungMotNguoi(t *testing.T) {
	kg := &khoLuuGia{}
	tg := MoThoGhi(kg, 1)

	id := int64(42)
	nv1 := store.NhanVat{ID: id, ViTriX: 10, ViTriY: 10}
	nv2 := store.NhanVat{ID: id, ViTriX: 20, ViTriY: 20}
	nv3 := store.NhanVat{ID: id, ViTriX: 30, ViTriY: 30}
	tg.Xep(nv1)
	tg.Xep(nv2)
	tg.Xep(nv3)

	// DongVaChoXong (test d dam bao hanh vi nay) buoc hang doi phai duoc ghi
	// HET truoc khi tra ve -> sau dong nay, moi ban con trong hang doi (chi
	// co the la ban gop cuoi cung, nv3) chac chan da xuong kho gia.
	dongVaChoXongVoiHan(t, tg, 2*time.Second)

	kg.mu.Lock()
	defer kg.mu.Unlock()
	if len(kg.goiGhi) == 0 {
		t.Fatal("khong co lan ghi nao ca — muon it nhat 1 (ban moi nhat)")
	}
	for _, nv := range kg.goiGhi {
		if nv == nv2 {
			t.Fatalf("ban GIUA (nv2=%+v) bi ghi RIENG xuong kho — hang doi khong gop theo id, "+
				"mat dung y nghia 'map id -> ban moi nhat'", nv2)
		}
	}
	cuoi := kg.goiGhi[len(kg.goiGhi)-1]
	if cuoi != nv3 {
		t.Fatalf("ban ghi CUOI CUNG = %+v, muon ban MOI NHAT nv3 = %+v", cuoi, nv3)
	}
}

// TestThoGhiXepKhongBaoGioChan la test (b) cua brief Task 3 — quan trong
// NHAT: LUAT SO MOT "actor khong bao gio cho o dia". Kho gia co tinh cham 2
// giay; goi Xep phai tra ve DUOI 50ms (ghim cung con so nay theo dung brief)
// DU kho cham bao nhieu — vi Xep chi cham hang doi trong bo nho (map+mutex),
// khong he goi toi kho.
//
// Test nay CUNG kiem hanh vi "moi lan ghi boc context.WithTimeout 5 giay"
// (dac ta §6.2) bang cach kiem tra ctx MA khoLuuGia nhan duoc co Deadline()
// hay khong — day la ly do dot bien "bo han gio ghi" (Step 5 brief) lam DUNG
// test nay do: bo context.WithTimeout la ctx.Deadline() se tra ok=false.
func TestThoGhiXepKhongBaoGioChan(t *testing.T) {
	kg := &khoLuuGia{doTre: 2 * time.Second}
	tg := MoThoGhi(kg, 1)
	defer dongVaChoXongVoiHan(t, tg, 4*time.Second)

	bd := time.Now()
	tg.Xep(store.NhanVat{ID: 1, ViTriX: 99})
	mat := time.Since(bd)
	if mat >= 50*time.Millisecond {
		t.Fatalf("Xep mat %v de tra ve — PHAI duoi 50ms du kho gia cham 2s "+
			"(vi pham LUAT SO MOT: actor khong bao gio cho o dia)", mat)
	}

	// Lan ghi NEN (khong lien quan gi toi phep do o tren) can toi ~2s de
	// hoan tat vi kho gia cham 2s — cho toi da 3s roi kiem tra context nhan
	// duoc co han gio hay khong.
	hetGio := time.After(3 * time.Second)
	for {
		kg.mu.Lock()
		soLan := len(kg.coHanGio)
		var coHanGio bool
		if soLan > 0 {
			coHanGio = kg.coHanGio[0]
		}
		kg.mu.Unlock()
		if soLan > 0 {
			if !coHanGio {
				t.Fatal("CapNhatNhanVat nhan ctx KHONG co han gio (Deadline() ok=false) — " +
					"thieu context.WithTimeout khi tho ghi goi kho")
			}
			return
		}
		select {
		case <-hetGio:
			t.Fatal("sau 3s van khong thay kho gia duoc goi lan nao")
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// TestThoGhiThuLaiKhiKhoLoiVaDemLoi la test (c) cua brief Task 3: kho tra loi
// o lan ghi DAU TIEN -> tho ghi phai (1) DEM loi (khong nuot), va (2) TU
// DONG thu lai o luot sau (KHONG can ai goi Xep them lan nao nua — day la
// diem phan biet "tu co nhip thu lai rieng" voi thiet ke "cho ben ngoai kich
// hoat lai"), va cuoi cung ghi THANH CONG.
func TestThoGhiThuLaiKhiKhoLoiVaDemLoi(t *testing.T) {
	kg := &khoLuuGia{soLoiDauTien: 1} // lan goi DAU tra loi; tu lan 2 tro di thanh cong
	tg := MoThoGhi(kg, 1)
	defer dongVaChoXongVoiHan(t, tg, 2*time.Second)

	tg.Xep(store.NhanVat{ID: 7, ViTriX: 5})

	hetGio := time.After(2 * time.Second)
	for {
		if kg.soLanGoi() >= 2 {
			break
		}
		select {
		case <-hetGio:
			t.Fatalf("sau 2s chi thay %d lan goi CapNhatNhanVat, muon >= 2 "+
				"(1 lan loi + tu dong thu lai thanh cong o luot sau)", kg.soLanGoi())
		case <-time.After(10 * time.Millisecond):
		}
	}

	if got := tg.soLoiGhi.Load(); got != 1 {
		t.Fatalf("soLoiGhi = %d, muon dung 1 (dung MOT lan loi duoc dem, khong bi nuot)", got)
	}

	// Doi them mot nhip nua de chac chan KHONG co lan goi thu 3 phat sinh
	// (vd bug lam ban DA GHI THANH CONG van bi xep lai vao hang doi) — chot
	// lai rang hang doi that su da RONG sau khi ghi thanh cong.
	time.Sleep(150 * time.Millisecond)
	if soLanCuoi := kg.soLanGoi(); soLanCuoi != 2 {
		t.Fatalf("sau khi da ghi thanh cong con phat sinh them %d lan goi nua (tong %d) — "+
			"ban da ghi thanh cong KHONG duoc nam lai trong hang doi", soLanCuoi-2, soLanCuoi)
	}
}

// TestThoGhiXepVaChoXoaBanCuTrongHangDoi la kiem chung Finding 2 cua soi
// toan nhanh KH11 (fix-cuoi-report.md): XepVaCho PHAI xoa ban CU cua CUNG id
// khoi hangDoi TRUOC khi ghi dong bo — neu khong, mot Xep(v_cu) vua xep ngay
// truoc do (con nam trong hangDoi, chua kip vongLap xa) se bi vongLap GHI DE
// LEN ban MOI (v_moi) ma XepVaCho vua ghi xong, dung bay dac ta §7.1 "ghi de
// ban moi bang ban cu".
//
// De dung bay nay MOT CACH TAT DINH (khong phu thuoc may man ve lich
// goroutine, tranh test flaky): dung mot "ban CHIEM CHO" (id 999, khong lien
// quan) de giu vongLap BAN trong toan bo phan con lai cua test — vongLap goi
// ghiHetHangDoi() MOT CACH DONG BO (khong phai goroutine rieng) va
// ghiHetHangDoi cho.Wait() cho toi khi ban chiem cho ghi xong (kg.doTre =
// 500ms). Trong luc do, ghiHetHangDoi() DA hoan doi t.hangDoi thanh map RONG
// MOI (xem chu thich ghiHetHangDoi) — nen Xep(vCu) goi SAU do CHAC CHAN roi
// vao map RONG MOI nay, khong the nao bi gop lan vao dot dang xu ly. Doi
// 100ms (gap doi khoangCachThuLai=50ms) sau khi Xep ban chiem cho de chac
// chan vongLap da kip nhan tin hieu va bat dau cho ghi (van con thua rat
// nhieu trong 500ms doTre).
func TestThoGhiXepVaChoXoaBanCuTrongHangDoi(t *testing.T) {
	kg := &khoLuuGia{doTre: 500 * time.Millisecond}
	tg := MoThoGhi(kg, 1)
	defer func() {
		// Luoi an toan don dep (dung tinh than dangnhaptrung_tichhop_test.go):
		// dam bao goroutine nen luon duoc dung han du test that bai som.
		ctx, huy := context.WithTimeout(context.Background(), 3*time.Second)
		defer huy()
		_ = tg.DongVaChoXong(ctx)
	}()

	// Ban chiem cho: giu vongLap ban (dang trong ghiHetHangDoi/cho.Wait())
	// trong 500ms.
	tg.Xep(store.NhanVat{ID: 999, ViTriX: 1})
	time.Sleep(100 * time.Millisecond)

	id := int64(77)
	vCu := store.NhanVat{ID: id, ViTriX: 10, ViTriY: 10}
	vMoi := store.NhanVat{ID: id, ViTriX: 99, ViTriY: 99}

	// Mo phong "nhip luu dinh ky vua xep xong, chua kip xa" (§7.1): vCu nam
	// trong hangDoi MOI (chac chan roi vao day, xem chu thich tren ham) trong
	// luc vongLap van con dang ban voi ban chiem cho.
	tg.Xep(vCu)

	// Cung luc do, phien BI DA cua CUNG nhan vat nay ghi DONG BO ban MOI HON
	// (vMoi) — day la duong XepVaCho that su, dung y het KhiDong goi khi
	// phien bi da (login.go, ghiDongBoKhiBiDa).
	if err := tg.XepVaCho(context.Background(), vMoi); err != nil {
		t.Fatalf("XepVaCho loi: %v", err)
	}

	// Xa het moi thu con lai (ban chiem cho neu chua xong, va — NEU thieu
	// fix — ban vCu con sot trong hangDoi) truoc khi kiem ket qua sau cung.
	dongVaChoXongVoiHan(t, tg, 3*time.Second)

	kg.mu.Lock()
	defer kg.mu.Unlock()
	var banCuoiCungChoId store.NhanVat
	var coBanNaoChoId bool
	for _, nv := range kg.goiGhi {
		if nv.ID == id {
			banCuoiCungChoId = nv
			coBanNaoChoId = true
		}
	}
	if !coBanNaoChoId {
		t.Fatalf("khong co lan ghi nao cho id=%d ca — mong doi it nhat ban MOI (tu XepVaCho)", id)
	}
	if banCuoiCungChoId != vMoi {
		t.Fatalf("ban CUOI CUNG duoc ghi xuong kho cho id=%d la %+v, mong doi VAN LA ban MOI %+v — "+
			"neu day la ban CU %+v: vongLap da GHI DE ban MOI (vua XepVaCho dong bo) bang ban CU con sot lai "+
			"trong hangDoi, dung bay dac ta §7.1 ma XepVaCho phai chan (xoa ban cu khoi hangDoi TRUOC khi ghi)",
			id, banCuoiCungChoId, vMoi, vCu)
	}
}

// TestThoGhiDongVaChoXongGhiHetHangDoi la test (d) cua brief Task 3: xep 3
// nguoi KHAC NHAU (khac ID, khong ai gop voi ai) roi goi DongVaChoXong ngay —
// ham nay phai ghi HET nhung gi con trong hang doi TRUOC KHI tra ve (dung
// cho luc tat may chu, dac ta §4 diem C: khong duoc mat tien do ai).
func TestThoGhiDongVaChoXongGhiHetHangDoi(t *testing.T) {
	kg := &khoLuuGia{}
	tg := MoThoGhi(kg, 2) // soThoGhi=2: kiem ca duong ghi song song nhieu tho

	tg.Xep(store.NhanVat{ID: 1, ViTriX: 1})
	tg.Xep(store.NhanVat{ID: 2, ViTriX: 2})
	tg.Xep(store.NhanVat{ID: 3, ViTriX: 3})

	dongVaChoXongVoiHan(t, tg, 2*time.Second)

	kg.mu.Lock()
	defer kg.mu.Unlock()
	if len(kg.goiGhi) != 3 {
		t.Fatalf("so lan CapNhatNhanVat = %d, muon dung 3 — hang doi phai duoc ghi HET "+
			"truoc khi DongVaChoXong tra ve", len(kg.goiGhi))
	}
	thay := map[int64]bool{}
	for _, nv := range kg.goiGhi {
		thay[nv.ID] = true
	}
	for _, id := range []int64{1, 2, 3} {
		if !thay[id] {
			t.Fatalf("thieu ban ghi cho id=%d — hang doi CHUA duoc ghi het khi DongVaChoXong tra ve", id)
		}
	}
}
