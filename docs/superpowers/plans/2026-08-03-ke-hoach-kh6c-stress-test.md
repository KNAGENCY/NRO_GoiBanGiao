# KH6c — Stress test 20 client — Kế hoạch thi công

> **Cho agent thi công:** SUB-SKILL BẮT BUỘC: dùng `superpowers:subagent-driven-development`
> (khuyến nghị) hoặc `superpowers:executing-plans` để làm từng task. Các bước dùng ô tick
> (`- [ ]`) để theo dõi.

**Goal:** Chứng minh dưới tải 20 người cùng khu vực dồn broadcast, actor KHÔNG treo, và thêm 2 chỉ
số quan sát (trễ nhịp) cho server.

**Architecture:** Test in-process trên chính `internal/game` (không mạng thật, không CSDL). Mảnh A
thêm đo trễ nhịp vào actor (`vongLap`), tách hàm thuần `capNhatTreNhip` để test tất định. Mảnh B/C
là 1 stress test dựng 20 phiên giả, 20 goroutine dồn `suKienDiChuyen` (đường broadcast đã được
`TestActorLoaiClientCham...` chứng minh là hit đúng phatChoNguoiKhac + gạch client chậm), 1 prober
chứng minh actor luôn trả lời.

**Tech Stack:** Go (`D:\Go\bin\go.exe`), `sync/atomic`, `testing`. Không thêm thư viện.

## Global Constraints (áp cho MỌI task)

- Định danh + chú thích **tiếng Việt không dấu**, theo đúng phong cách file quanh nó.
- Hai test này **không cần CSDL** (in-process) — chạy được không cần `NRO_TEST_DSN`.
- Trước khi commit mỗi task: `go test ./internal/game/ -count=5` xanh, `go vet ./...` sạch,
  `gofmt -l internal/game/` không in tên file mình sửa.
- **Kiểm thử đột biến bắt buộc** (§6.11): mỗi task có bước phá code → test phải đỏ → khôi phục.
- **KHÔNG** đổi hành vi actor (trần hộp thư 256, gạch client chậm) — chỉ THÊM đo.
- File tiếng Việt sửa bằng Edit/Write, **không** dùng shell (§6.15).
- Sửa lên nhánh `kh6c-stress-test` (đã có 2 commit: doc fix + spec).

---

### Task 1: Mảnh A — Đo trễ nhịp actor (hàm thuần + gắn vào vòng lặp)

**Files:**
- Modify: `server/internal/game/khuvuc.go` (struct `KhuVuc`: +3 trường; +hằng `nguongNhipTre`;
  +method `capNhatTreNhip`; sửa nhánh `case <-ticker.C` trong `vongLap`)
- Test: `server/internal/game/trenhip_test.go` (tạo mới)

**Interfaces:**
- Produces (Task 2 dùng): `k.soNhipTre atomic.Int64`, `k.treNhipLauNhatMs atomic.Int64`,
  `func (k *KhuVuc) capNhatTreNhip(khoangCach time.Duration)`.

- [ ] **Step 1: Viết test thuần (đỏ trước)** — tạo `server/internal/game/trenhip_test.go`:

```go
package game

import (
	"testing"
	"time"
)

// capNhatTreNhip la ham THUAN (khong dung dong ho): nhan khoang cach giua hai
// nhip, quyet dinh co dem la "tre" khong + cap nhat max. Test truc tiep.
func TestCapNhatTreNhip(t *testing.T) {
	k := &KhuVuc{}

	// Duoi/bang nguong 200ms: KHONG tinh la tre.
	k.capNhatTreNhip(150 * time.Millisecond)
	k.capNhatTreNhip(200 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 0 {
		t.Fatalf("150ms va 200ms deu <= nguong, soNhipTre=%d, mong 0", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 0 {
		t.Fatalf("chua co nhip tre, treNhipLauNhatMs=%d, mong 0", got)
	}

	// Vuot nguong: dem + ghi max.
	k.capNhatTreNhip(300 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 1 {
		t.Fatalf("300ms vuot nguong, soNhipTre=%d, mong 1", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 300 {
		t.Fatalf("treNhipLauNhatMs=%d, mong 300", got)
	}

	// Tre hon: max len 500.
	k.capNhatTreNhip(500 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 2 {
		t.Fatalf("soNhipTre=%d, mong 2", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 500 {
		t.Fatalf("treNhipLauNhatMs=%d, mong 500", got)
	}

	// Tre nhung NHO hon max cu: dem tang, max GIU 500.
	k.capNhatTreNhip(250 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 3 {
		t.Fatalf("soNhipTre=%d, mong 3", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 500 {
		t.Fatalf("treNhipLauNhatMs=%d, mong giu 500", got)
	}
}
```

- [ ] **Step 2: Chạy để chắc nó ĐỎ**

Chạy: `D:\Go\bin\go.exe test ./internal/game/ -run TestCapNhatTreNhip`
Mong: FAIL biên dịch — `k.soNhipTre undefined` / `capNhatTreNhip undefined`.

- [ ] **Step 3: Thêm trường vào struct `KhuVuc`** (trong `khuvuc.go`, đặt cạnh `soLanGuiRoi`):

```go
	// Do tre nhip (KH6c): mocNhipTruoc chi vongLap cham (mot luong) nen la bien
	// thuong; soNhipTre/treNhipLauNhatMs la atomic vi test/giam sat doc song song.
	mocNhipTruoc     time.Time
	soNhipTre        atomic.Int64
	treNhipLauNhatMs atomic.Int64
```

- [ ] **Step 4: Thêm hằng + method `capNhatTreNhip`** (trong `khuvuc.go`, gần `ghiNhanNhip`):

```go
// nguongNhipTre: khoang cach giua hai nhip vuot con so nay thi tinh la "tre"
// (2x chu ky 100ms). Kieu hang so de kiem thu dot bien sua duoc.
const nguongNhipTre = 200 * time.Millisecond

// capNhatTreNhip nhan khoang cach tu nhip truoc toi nhip nay. Vuot nguong thi
// dem (soNhipTre) + cap nhat max (treNhipLauNhatMs). Ham THUAN theo nghia khong
// doc dong ho — vongLap truyen khoang cach vao — nen test tat dinh duoc.
func (k *KhuVuc) capNhatTreNhip(khoangCach time.Duration) {
	if khoangCach <= nguongNhipTre {
		return
	}
	k.soNhipTre.Add(1)
	if ms := khoangCach.Milliseconds(); ms > k.treNhipLauNhatMs.Load() {
		k.treNhipLauNhatMs.Store(ms) // chi vongLap ghi -> khong dua ghi
	}
}
```

- [ ] **Step 5: Chạy lại — test thuần XANH**

Chạy: `D:\Go\bin\go.exe test ./internal/game/ -run TestCapNhatTreNhip`
Mong: PASS.

- [ ] **Step 6: Gắn đo vào `vongLap`** — sửa nhánh `case <-ticker.C` (hiện ~dòng 563):

Từ:
```go
		case <-ticker.C:
			daPanic := antoan.ChayCoLuoi("khuvuc.nhipTim", k.nhipTim)
			if k.ghiNhanNhip(daPanic) {
				return
			}
```
Thành:
```go
		case <-ticker.C:
			// KH6c: do khoang cach thuc te giua hai nhip. Ticker cua Go BO tick
			// neu actor ban khong kip nhan -> khoang cach > 100ms bao actor dang bi
			// don. Nhip dau (mocNhipTruoc zero) bo qua vi chua co moc.
			bayGio := time.Now()
			if !k.mocNhipTruoc.IsZero() {
				k.capNhatTreNhip(bayGio.Sub(k.mocNhipTruoc))
			}
			k.mocNhipTruoc = bayGio
			daPanic := antoan.ChayCoLuoi("khuvuc.nhipTim", k.nhipTim)
			if k.ghiNhanNhip(daPanic) {
				return
			}
```

- [ ] **Step 7: Toàn bộ test gói game vẫn xanh (không hồi quy)**

Chạy: `D:\Go\bin\go.exe test ./internal/game/ -count=5`
Mong: PASS (không cần CSDL cho các test này; test cần CSDL tự SKIP).
Rồi: `D:\Go\bin\go.exe vet ./...` và `D:\Go\bin\gofmt.exe -l internal/game/` (không in `khuvuc.go`).

- [ ] **Step 8: Kiểm thử đột biến** — sửa tạm `capNhatTreNhip` thành no-op:

```go
func (k *KhuVuc) capNhatTreNhip(khoangCach time.Duration) { return }
```
Chạy: `D:\Go\bin\go.exe test ./internal/game/ -run TestCapNhatTreNhip`
Mong: **FAIL** (soNhipTre=0). Nếu vẫn PASS → test vô dụng, viết lại. **Khôi phục** bản đúng, chạy lại PASS.

- [ ] **Step 9: Commit**

```bash
git add server/internal/game/khuvuc.go server/internal/game/trenhip_test.go
git commit -m "feat(game): do tre nhip actor (KH6c manh A, ham thuan capNhatTreNhip)

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 2: Mảnh B+C — Stress test 20 client dồn broadcast, actor không treo

**Files:**
- Create: `server/internal/game/khuvuc_stress_test.go`
- (Dùng lại `phienGiaKhuVuc` + `batDauLoiXep()` + `AnhChupNguoiChoiCoLoi` đã có trong gói)

**Interfaces:**
- Consumes: `MoKhuVuc`, `TimBanDo`, `k.BatDau/Dung`, `k.Gui(suKienVao/suKienDiChuyen)`,
  `k.AnhChupNguoiChoiCoLoi(0) ([]store.NhanVat, error)`, `phienGiaKhuVuc{}` + `.batDauLoiXep()`,
  `k.soLanGuiRoi/soNhipTre/treNhipLauNhatMs` (Task 1).

- [ ] **Step 1: Viết stress test** — tạo `server/internal/game/khuvuc_stress_test.go`:

```go
package game

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"nroserver/internal/store"
)

// TestStress20ClientDonBroadcastActorKhongTreo: 20 nguoi cung khu vuc, 20
// goroutine don su kien di chuyen (moi cai -> broadcast phatChoNguoiKhac cho 19
// nguoi con lai). 3 nguoi "cham" (hang day) phai bi gach ma actor KHONG treo:
// prober phai lay duoc anh chup (bang chung actor con song), va cuoi cung chi
// con 17 nguoi khoe. Do la bang chung §6.10b khong tai dien duoi tai.
func TestStress20ClientDonBroadcastActorKhongTreo(t *testing.T) {
	const soNguoi = 20
	const soCham = 3
	const thoiLuong = 2 * time.Second

	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })

	phiens := make([]*phienGiaKhuVuc, soNguoi)
	ids := make([]int64, soNguoi)
	for i := 0; i < soNguoi; i++ {
		id := int64(1000 + i)
		ids[i] = id
		p := &phienGiaKhuVuc{}
		phiens[i] = p
		nv := &store.NhanVat{ID: id, Ten: "N", HpToiDa: 100, ViTriX: 100, ViTriY: 384}
		k.Gui(suKienVao{nv: nv, phien: p})
	}

	// Cho actor xu ly xong 20 suKienVao truoc khi do tai.
	doiDuNguoi(t, k, soNguoi)

	// soCham nguoi cuoi thanh "cham": hang day -> se bi gach khi co broadcast.
	for i := soNguoi - soCham; i < soNguoi; i++ {
		phiens[i].batDauLoiXep()
	}

	var dungBom atomic.Bool
	var tongBom atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < soNguoi; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			x := int32(100)
			for !dungBom.Load() {
				x++
				if x > 200 {
					x = 100
				}
				k.Gui(suKienDiChuyen{idNhanVat: id, x: x, y: 384})
				tongBom.Add(1)
			}
		}(ids[i])
	}

	// Prober: lien tuc lay anh chup, do actor con tra loi khong.
	var soProbeOK, soProbeTimeout atomic.Int64
	var maxProbeNs atomic.Int64
	var dungProbe atomic.Bool
	var wgProbe sync.WaitGroup
	wgProbe.Add(1)
	go func() {
		defer wgProbe.Done()
		for !dungProbe.Load() {
			t0 := time.Now()
			_, err := k.AnhChupNguoiChoiCoLoi(0)
			d := time.Since(t0).Nanoseconds()
			if err != nil {
				soProbeTimeout.Add(1) // actor ban qua 100ms luc nay
			} else {
				soProbeOK.Add(1)
				if d > maxProbeNs.Load() {
					maxProbeNs.Store(d)
				}
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()

	time.Sleep(thoiLuong)
	dungBom.Store(true)
	wg.Wait()
	dungProbe.Store(true)
	wgProbe.Wait()

	// --- KHANG DINH (gate cung, chon nguong RONG de khong chap chon) ---

	// 1) Actor KHONG treo kieu §6.10b: neu actor ket 30s tren socket client cham
	//    thi prober se co 0 lan OK trong suot 2s. Yeu cau co lay duoc anh chup.
	if soProbeOK.Load() == 0 {
		t.Fatal("actor KHONG tra loi probe lan nao trong 2s -> nghi treo kieu §6.10b")
	}

	// 2) Sau khi ngung don tai, actor phai hoi phuc va tra dung trang thai cuoi.
	//    Client cham da bi gach -> chi con soNguoi-soCham nguoi khoe.
	var conLai int
	hanHoiPhuc := time.Now().Add(time.Second)
	for {
		ban, err := k.AnhChupNguoiChoiCoLoi(0)
		if err == nil {
			conLai = len(ban)
			break
		}
		if time.Now().After(hanHoiPhuc) {
			t.Fatal("sau khi ngung tai, actor van khong tra loi trong 1s -> treo")
		}
	}
	if conLai != soNguoi-soCham {
		t.Fatalf("con %d nguoi, mong %d (phai gach het %d client cham)", conLai, soNguoi-soCham, soCham)
	}

	// 3) Chan stall benh ly (rong rai, chi bat treo that su dai).
	if tre := k.treNhipLauNhatMs.Load(); tre >= 1000 {
		t.Errorf("tre nhip lau nhat %dms >= 1000ms -> actor bi don qua nang", tre)
	}

	// --- BAO SO (quan sat, khong lam dieu kien dau/rot) ---
	t.Logf("KH6c stress: bom=%d roiHopThu=%d soNhipTre=%d treNhipLauNhat=%dms probeOK=%d probeTimeout=%d maxProbe=%dms",
		tongBom.Load(), k.soLanGuiRoi.Load(), k.soNhipTre.Load(), k.treNhipLauNhatMs.Load(),
		soProbeOK.Load(), soProbeTimeout.Load(), time.Duration(maxProbeNs.Load()).Milliseconds())
}

// doiDuNguoi cho toi khi actor da xu ly du `mong` nguoi (hoac fail sau 2s).
func doiDuNguoi(t *testing.T, k *KhuVuc, mong int) {
	t.Helper()
	han := time.Now().Add(2 * time.Second)
	for {
		ban, err := k.AnhChupNguoiChoiCoLoi(0)
		if err == nil && len(ban) == mong {
			return
		}
		if time.Now().After(han) {
			t.Fatalf("actor chua xu ly du %d suKienVao (hien %d)", mong, len(ban))
		}
		time.Sleep(2 * time.Millisecond)
	}
}
```

- [ ] **Step 2: Chạy stress test — mong XANH**

Chạy: `D:\Go\bin\go.exe test ./internal/game/ -run TestStress20Client -v`
Mong: PASS, và dòng `t.Logf` in ra số liệu. Xem `probeOK` phải **> 0** và `conLai` = 17.
Nếu chập chờn: xem lại §6 spec (nới nhịp probe/thời lượng), KHÔNG nới tới mức che lỗi thật.

- [ ] **Step 3: Chạy lặp chống chập chờn**

Chạy: `D:\Go\bin\go.exe test ./internal/game/ -run TestStress20Client -count=5`
Mong: PASS cả 5 lần. Nếu có lần đỏ ở gate #3 (trễ nhịp) → nới ngưỡng lên (vd 1500ms) và ghi lý do
vào chú thích; gate #1/#2 (không treo, gạch đủ) **không được nới**.

- [ ] **Step 4: Kiểm thử đột biến** — bỏ gạch client chậm trong `phatChoNguoiKhac` (`khuvuc.go`):

Tạm xoá dòng `k.loaiNhieuNguoiLoiGui(loi)` trong `phatChoNguoiKhac`.
Chạy: `D:\Go\bin\go.exe test ./internal/game/ -run TestStress20Client`
Mong: **FAIL** ở gate #2 (`conLai` = 20 chứ không phải 17, vì client chậm không bị gạch).
**Khôi phục** dòng đó, chạy lại PASS.

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/khuvuc_stress_test.go
git commit -m "test(game): stress 20 client don broadcast, actor khong treo (KH6c)

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

## Self-Review (điền khi viết xong, kiểm plan với spec)

- **Phủ spec:** Mảnh A (đo trễ nhịp) → Task 1. Mảnh B (stress 20 người + prober) → Task 2 Step 1.
  Mảnh C (khẳng định + báo số + đột biến) → Task 2 Step 2–4. ✅
- **Lệch spec có chủ đích (ghi rõ):**
  1. Dùng `suKienDiChuyen` (di chuyển) làm sự kiện dồn, KHÔNG phải `suKienDanhQuai` (đánh quái).
     Lý do: cùng đường tới hạn (mailbox → phatChoNguoiKhac → XepGui → gạch client chậm) mà
     `TestActorLoaiClientCham...` đã chứng minh là hit đúng, lại tất định (không phụ thuộc trạng
     thái quái/HP). Đúng mục tiêu stress "dồn broadcast vào actor".
  2. Trễ nhịp + rơi hộp thư → **báo số** (quan sát), gate cứng chỉ đặt ở "actor không treo" +
     "gạch đủ client chậm" + stall benh ly (trễ ≥1000ms). Spec §6 đã lường: tránh chập chờn.
- **Không placeholder:** mọi step có code thật. ✅
- **Nhất quán kiểu:** `soNhipTre`/`treNhipLauNhatMs`/`capNhatTreNhip` dùng cùng tên ở Task 1 và 2. ✅
