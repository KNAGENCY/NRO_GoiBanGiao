package game

import (
	"io"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// phienStress: phien gia toi thieu cho stress test.
//   - khoe (cham=false): NUOT goi (khong luu) — tranh phinh bo nho khi bi
//     broadcast hang chuc nghin goi trong 2s.
//   - cham (cham=true): luon tra ErrHangGuiDay (mo phong client khong rut hang)
//     -> actor phai gach ten no.
type phienStress struct {
	cham      atomic.Bool
	soLanDong atomic.Int64
}

func (p *phienStress) XepGui(*protocol.Message) error {
	if p.cham.Load() {
		return mynet.ErrHangGuiDay
	}
	return nil
}

func (p *phienStress) DongCham() { p.soLanDong.Add(1) }

// TestStress20ClientDonBroadcastActorKhongTreo: 20 nguoi cung khu vuc, 20
// goroutine don su kien di chuyen (moi cai xu xong -> phatChoNguoiKhac cho 19
// nguoi con lai). 3 nguoi "cham" phai bi gach ma actor KHONG treo: prober phai
// lay duoc anh chup (bang chung actor con song), va cuoi cung chi con 17 nguoi
// khoe. Do la bang chung §6.10b khong tai dien duoi tai.
func TestStress20ClientDonBroadcastActorKhongTreo(t *testing.T) {
	const soNguoi = 20
	const soCham = 3
	const thoiLuong = 2 * time.Second

	// Tat log trong luc stress: hang tram nghin drop se sinh hang tram nghin dong
	// WARN, ngap output. Drop la CO Y (hop thu day duoi tai) — da co soLanGuiRoi
	// dem. Khoi phuc logger sau test.
	logCu := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	t.Cleanup(func() { slog.SetDefault(logCu) })

	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })

	phiens := make([]*phienStress, soNguoi)
	ids := make([]int64, soNguoi)
	for i := 0; i < soNguoi; i++ {
		id := int64(1000 + i)
		ids[i] = id
		p := &phienStress{}
		phiens[i] = p
		nv := &store.NhanVat{ID: id, Ten: "N", HpToiDa: 100, ViTriX: 100, ViTriY: 384}
		k.Gui(suKienVao{nv: nv, phien: p})
	}

	// Cho actor xu ly xong 20 suKienVao truoc khi do tai.
	doiDuNguoi(t, k, soNguoi)

	// soCham nguoi cuoi thanh "cham": hang day -> se bi gach khi co broadcast toi.
	for i := soNguoi - soCham; i < soNguoi; i++ {
		phiens[i].cham.Store(true)
	}

	var dungBom atomic.Bool
	var tongBom atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < soNguoi; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			x := int16(100)
			for !dungBom.Load() {
				x++
				if x > 110 { // giu trong [100,110]: hop o di duoc cua map 1 (footprint 120x408)
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

	// --- KHANG DINH (gate cung, nguong RONG de khong chap chon) ---

	// 1) Actor KHONG treo kieu §6.10b: neu actor ket 30s tren socket client cham
	//    thi prober se co 0 lan OK trong suot 2s. Yeu cau co lay duoc anh chup.
	if soProbeOK.Load() == 0 {
		t.Fatal("actor KHONG tra loi probe lan nao trong 2s -> nghi treo kieu §6.10b")
	}

	// 2) Sau khi ngung don tai, actor phai hoi phuc va tra dung trang thai cuoi:
	//    client cham da bi gach -> chi con soNguoi-soCham nguoi khoe.
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
