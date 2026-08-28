package game

import (
	"bytes"
	"testing"
	"time"

	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 8 — dinh tuyen goi len -60 (nguoi danh nguoi) LEN BOSS.
//
// Boss vao the gioi bang chinh goi -5 cua nguoi choi (goiXuatHienBoss) nen voi
// client no la mot "char" trong map: nut danh gui LEN dung goi -60 nay, mang
// id boss (AM). File nay khoa NAM dieu, moi dieu la mot cach hong am tham:
//
//  1. Don danh vao boss THAT SU tru mau boss VA ghi so dong gop — tuc di qua
//     bossNhanSatThuong (Task 12), khong phai mot duong `b.Hp -= sat` rieng.
//     Duong rieng do lam nguoi choi MAT SACH phan thuong ma khong bao gio bao
//     loi.
//  2. Goi phat ra la -60 (KHONG phai -9/-12 cua quai) va con so trong goi la
//     sat thuong THAT SU da ap len boss.
//  3. Chot id AM: mot id am chi duoc di duong boss, mot id duong chi duoc di
//     duong nguoi choi — khong duong nao nuot duong nao.
//  4. Boss KHONG dang tren man hinh (NGHI / RA_DI) hoac da la cai xac (TROI,
//     hp=0) thi khong danh duoc, va id bia dat khong duoc lam sap actor.
//  5. Nhung chot chung cua ca nhanh -60 (nguoi danh da chet, chieu hoi mau,
//     bat tu, vuot tam) van ap DUNG NGUYEN cho huong danh boss.
//
// Moi test goi k.xuLy(...) TRUC TIEP (actor CHUA BatDau) — dung khuon
// khuvuc_boss_test.go / khuvuc_thuongboss_test.go: khong goroutine, khong
// sleep, dong ho ghim.
// ─────────────────────────────────────────────────────────────────────────────

// mocDanhBoss la moc thoi gian ghim cho ca file. Ghim cung (khong time.Now) de
// moi khang dinh ve bat tu / hoi chieu la tat dinh.
var mocDanhBoss = time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)

// bossDanhDuoc dung mot boss MOT PHA nam DUNG CHO NGUOI CHOI DUNG (0,0) nen
// moi don deu trong tam — test nao muon kiem vuot tam thi tu doi ViTriX.
//
// Lang 3 (Da = He Tho) doi voi nguoi danh Lang 0 (La = He Hoa) la quan he
// TRUNG TINH (100%) — xem heNguyenToTheoLang (khuvuc.go): Hoa khac Phong,
// Thuy khac Hoa; Hoa va Tho khong khac che nhau. Nho vay sat thuong tho trong
// ca file nay bang DUNG con so cua nhanh PvP: satThuongGoc(lang0)=12 x
// chieu can chien cap 1 (SatThuong=20) / 100 = 2.
func bossDanhDuoc(trangThai TrangThaiBoss, hp int32, diemThuong int64) *Boss {
	return &Boss{
		Id:         -7,
		PhaHienTai: 0,
		DuLieu: []PhaBoss{{
			Ten:        "Cửu Vĩ",
			Lang:       3,
			Cap:        90,
			TypePk:     TypePkBoss,
			DiemThuong: diemThuong,
		}},
		Hp:        hp,
		HpToiDa:   1000,
		TrangThai: trangThai,
	}
}

// khuVucDanhBoss dung mot khu vuc CHUA BatDau() voi dong ho ghim, hai nguoi
// choi (A se danh, B chi dung xem — de kiem "phat cho CA khu vuc") va mot con
// boss da vao khu vuc. A da chon chieu can chien (mauChieuId=1, skillId=1,
// tam 100).
//
// Tra ve sau khi DA DON SACH goi -5 phat sinh tu buoc chuan bi, de moi test
// chi con thay dung goi tin do chinh kich ban cua no tao ra.
func khuVucDanhBoss(t *testing.T, b *Boss) (k *KhuVuc, pA, pB *phienGiaKhuVuc, nvA *store.NhanVat) {
	t.Helper()
	k = MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return mocDanhBoss }
	nvA = &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 150, HpToiDa: 150}
	nvB := &store.NhanVat{ID: 2, Ten: "B", Lang: 0, ChakraHienTai: 100, HpHienTai: 150, HpToiDa: 150}
	pA, pB = &phienGiaKhuVuc{}, &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	k.xuLy(suKienVao{nv: nvB, phien: pB})
	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 1}) // can chien, tam 100, skillId=1
	if b != nil {
		k.themBoss(b)
	}
	pA.nhan, pB.nhan = nil, nil
	return k, pA, pB, nvA
}

// goi60DanhBoss dung than goi -60 mong doi cho "A (id 1) danh boss id -7".
// Dung ham thay vi hang so de tung test chi phai noi ba thu no quan tam.
func goi60DanhBoss(satThuong byte, daChet, chiMang byte) []byte {
	return []byte{
		0, 0, 0, 1, // idNguoiDanh = 1 (A)
		1,                      // skillId = 1 (lang 0, slot 0, cap 1)
		1,                      // soMucTieu = 1
		0xFF, 0xFF, 0xFF, 0xF9, // idNguoiTrung = -7 (id boss, AM)
		1,                  // docTiep = 1 (dang don danh)
		0,                  // loaiChieu = 0
		0, 0, 0, satThuong, // satThuong
		daChet,
		chiMang,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 1 — duong chinh: danh trung boss
// ─────────────────────────────────────────────────────────────────────────────

// TestDanhBossTruMauVaPhatGoi60ChoCaKhuVuc la kich ban chinh: A danh boss dang
// HOAT_DONG, boss con song sau don.
//
// Khang dinh DUNG MOT goi -60 cho MOI NGUOI trong khu vuc (ke ca chinh A) —
// va TUYET DOI khong co goi thu hai. VI SAO khong co -30/9 nhu nhanh PvP: goi
// -5 cua NGUOI CHOI (goiXuatHien) ghi hp = HpToiDa (nguoi vao giua tran thay
// nguoi khac full mau) nen nhanh PvP buoc phai gui -30/9 sua lai; goi -5 cua
// BOSS (goiXuatHienBoss) ghi b.Hp THAT, va con so trong -60 la sat thuong that
// da ap — nen moi client tu tru dung bang may chu, khong co gi de sua.
func TestDanhBossTruMauVaPhatGoi60ChoCaKhuVuc(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, pA, pB, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 998 {
		t.Fatalf("Hp boss sau don = %d, mong doi 998 (1000 - satThuong 2)", b.Hp)
	}
	muon := goi60DanhBoss(2, 0, 0)
	for ten, p := range map[string]*phienGiaKhuVuc{"A (nguoi danh)": pA, "B (nguoi xem)": pB} {
		if len(p.nhan) != 1 {
			t.Fatalf("%s nhan %d goi, mong doi DUNG 1 (chi -60, KHONG co -30/9)", ten, len(p.nhan))
		}
		if p.nhan[0].Command != lenhNguoiDanhNguoi {
			t.Fatalf("%s: command = %d, mong doi %d (-60)", ten, p.nhan[0].Command, lenhNguoiDanhNguoi)
		}
		if !bytes.Equal(p.nhan[0].Data, muon) {
			t.Fatalf("%s: than goi -60 sai:\n got  = %v\n muon = %v", ten, p.nhan[0].Data, muon)
		}
	}
}

// TestDanhBossGhiSoDongGop la KIEM THU DOT BIEN cua rang buoc so mot: neu ai
// do thay k.bossNhanSatThuong bang `b.Hp -= sat` truc tiep thi mau boss VAN
// giam dung (test tren van xanh) nhung so dong gop RONG -> den luc boss chet
// khong ai duoc dong nao. Test nay la cai bay do.
func TestDanhBossGhiSoDongGop(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, _, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if got := b.SatThuongTheoNguoi[nvA.ID]; got != 2 {
		t.Fatalf("so dong gop cua A = %d, mong doi 2 — don danh PHAI di qua bossNhanSatThuong, khong duoc tu tru b.Hp",
			got)
	}
}

// TestDanhBossChetChiaThuongVaChuyenTroi la KIEM THU DOT BIEN thu hai cho cung
// rang buoc: don ket lieu phai keo theo TOAN BO nghi thuc boss chet (chia
// thuong vao nv.KinhNghiem, xoa so dong gop, chuyen trang thai TROI) — tat ca
// nam trong bossNhanSatThuong -> bossChet. Mot duong `b.Hp -= sat` rieng se de
// boss nam o hp=0 ma VAN o trang thai HOAT_DONG, chia thuong khong bao gio
// chay, va con so daChet trong goi -60 khong bao gio bat len.
func TestDanhBossChetChiaThuongVaChuyenTroi(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 2, 100) // dung 2 mau = dung mot don
	k, pA, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 0 {
		t.Fatalf("Hp boss = %d, mong doi 0 (chet boi don nay)", b.Hp)
	}
	if b.TrangThai != BossTroi {
		t.Fatalf("trang thai boss = %s, mong doi TROI (nghi thuc chet chua chay)", b.TrangThai)
	}
	if b.SatThuongTheoNguoi != nil {
		t.Fatalf("so dong gop = %v, mong doi nil (chia xong PHAI xoa sach)", b.SatThuongTheoNguoi)
	}
	if nvA.KinhNghiem != 100 {
		t.Fatalf("KinhNghiem cua A = %d, mong doi 100 (an tron DiemThuong vi chi minh A danh)", nvA.KinhNghiem)
	}
	muon := goi60DanhBoss(2, 1, 0) // daChet = 1
	if len(pA.nhan) != 1 || !bytes.Equal(pA.nhan[0].Data, muon) {
		t.Fatalf("goi -60 cua don ket lieu sai:\n got  = %+v\n muon = %v", pA.nhan, muon)
	}
}

// TestDanhBossSatThuongKhongVuotMauConLai khoa dac ta §1.2 ("sat thuong khong
// bao giờ vuot mau con lai") cho huong danh boss: boss con 1 mau, don danh
// tho la 2 -> goi -60 phai bao DUNG 1, khong phai 2. Client tu tru con so
// trong goi (Char.doInjure) nen bao thua la thanh mau cua moi nguoi lech khoi
// may chu ngay lap tuc.
func TestDanhBossSatThuongKhongVuotMauConLai(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1, 0)
	k, pA, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 0 {
		t.Fatalf("Hp boss = %d, mong doi 0", b.Hp)
	}
	muon := goi60DanhBoss(1, 1, 0) // satThuong KEP ve 1, daChet = 1
	if len(pA.nhan) != 1 || !bytes.Equal(pA.nhan[0].Data, muon) {
		t.Fatalf("goi -60 sai:\n got  = %+v\n muon = %v", pA.nhan, muon)
	}
}

// TestDanhBossHaiNguoiChiaThuongTheoDongGop: hai nguoi cung danh, nguoi danh
// nhieu hon an nhieu hon. Kiem ca duong dong gop nhieu nguoi lan viec don cua
// nguoi thu hai khong xoa mat dong gop cua nguoi dau.
func TestDanhBossHaiNguoiChiaThuongTheoDongGop(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 6, 300)
	k, _, _, nvA := khuVucDanhBoss(t, b)
	// B cung phai co chieu moi danh duoc.
	k.xuLy(suKienChonChieu{idNhanVat: 2, mauChieuId: 1})
	nvB := k.nguoiTrong[2].nv

	// A danh 2 don (2+2=4 mau), B danh 1 don (2 mau) -> 6 mau, boss chet.
	// Moi don cach nhau du xa de qua kiem hoi chieu (1000ms).
	moc := mocDanhBoss
	k.dongHo = func() time.Time { return moc }
	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})
	moc = moc.Add(2 * time.Second)
	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})
	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvB.ID, idMucTieu: b.Id})

	if b.Hp != 0 {
		t.Fatalf("Hp boss = %d, mong doi 0", b.Hp)
	}
	if nvA.KinhNghiem != 200 || nvB.KinhNghiem != 100 {
		t.Fatalf("KinhNghiem A/B = %d/%d, mong doi 200/100 (dong gop 4/6 va 2/6 cua 300 diem)",
			nvA.KinhNghiem, nvB.KinhNghiem)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2 — chot id am: hai duong khong duoc nuot nhau
// ─────────────────────────────────────────────────────────────────────────────

// TestDanhIdDuongVanDiDuongNguoiChoi: co boss trong khu vuc, nhung A nham vao
// NGUOI (id duong) -> nguoi mat mau, boss KHONG suy suyen. Bao ve chieu nguoc
// lai cua chot id am: nhanh boss khong duoc nuot mat PvP.
func TestDanhIdDuongVanDiDuongNguoiChoi(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, _, _, nvA := khuVucDanhBoss(t, b)
	nvB := k.nguoiTrong[2].nv

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: nvB.ID})

	if nvB.HpHienTai != 148 {
		t.Fatalf("HpHienTai cua B = %d, mong doi 148 (150 - 2) — don vao id DUONG phai di duong nguoi choi", nvB.HpHienTai)
	}
	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (khong dinh dang gi toi don PvP)", b.Hp)
	}
}

// TestDanhBossIdAmKhongCoTrongDanhSachKhongPanic: id AM nhung khong con boss
// nao mang id do (boss vua bi don di, hoac client bia so) -> im lang bo qua,
// khong goi tin nao, KHONG panic. Mot id bia dat khong duoc giet actor.
func TestDanhBossIdAmKhongCoTrongDanhSachKhongPanic(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, pA, pB, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: -999})

	if len(pA.nhan) != 0 || len(pB.nhan) != 0 {
		t.Fatalf("id am bia dat: pA=%d pB=%d goi, mong doi ca hai deu 0", len(pA.nhan), len(pB.nhan))
	}
	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (id khong khop thi khong con boss nao duoc dinh toi)", b.Hp)
	}
}

// TestDanhBossKhuVucKhongCoBossNaoKhongPanic: khu vuc rong tron (k.boss nil) —
// vong tim khong duoc panic hay danh nham ai.
func TestDanhBossKhuVucKhongCoBossNaoKhongPanic(t *testing.T) {
	k, pA, _, nvA := khuVucDanhBoss(t, nil)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: -1})

	if len(pA.nhan) != 0 {
		t.Fatalf("A nhan %d goi khi khu vuc khong co boss nao, mong doi 0", len(pA.nhan))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 3 — trang thai boss khong danh duoc
// ─────────────────────────────────────────────────────────────────────────────

// TestDanhBossKhongTrenManHinhBiTuChoi la KIEM THU DOT BIEN cua chot trang
// thai: boss NGHI (chua hien, client khong he biet id do) va RA_DI (da don di)
// KHONG duoc an don. Bo chot nay thi mot client giu id cu se ban duoc vao con
// boss dang dem nguoc — no chet truoc ca khi kip xuat hien.
func TestDanhBossKhongTrenManHinhBiTuChoi(t *testing.T) {
	for _, tt := range []TrangThaiBoss{BossNghi, BossRaDi} {
		t.Run(tt.String(), func(t *testing.T) {
			b := bossDanhDuoc(tt, 1000, 0)
			k, pA, pB, nvA := khuVucDanhBoss(t, b)

			k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

			if b.Hp != 1000 {
				t.Fatalf("boss %s: Hp = %d, mong doi 1000 (khong tren man hinh thi khong danh duoc)", tt, b.Hp)
			}
			if len(pA.nhan) != 0 || len(pB.nhan) != 0 {
				t.Fatalf("boss %s: pA=%d pB=%d goi, mong doi ca hai deu 0 (im lang, dung nhu id khong khop ai)",
					tt, len(pA.nhan), len(pB.nhan))
			}
		})
	}
}

// TestDanhBossDangHienVaChaoVanDanhDuoc: ba trang thai "dang tro len va con
// song" (HIEN, CHAO, HOAT_DONG) deu phai an don. Thieu test nay thi ai do co
// the thu hep dieu kien thanh "chi HOAT_DONG" — nguoi choi bam danh trong luc
// boss dang chao se thay may chu im lang.
func TestDanhBossDangHienVaChaoVanDanhDuoc(t *testing.T) {
	for _, tt := range []TrangThaiBoss{BossHien, BossChao, BossHoatDong} {
		t.Run(tt.String(), func(t *testing.T) {
			b := bossDanhDuoc(tt, 1000, 0)
			k, pA, _, nvA := khuVucDanhBoss(t, b)

			k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

			if b.Hp != 998 {
				t.Fatalf("boss %s: Hp = %d, mong doi 998", tt, b.Hp)
			}
			if len(pA.nhan) != 1 {
				t.Fatalf("boss %s: A nhan %d goi, mong doi 1 (-60)", tt, len(pA.nhan))
			}
		})
	}
}

// TestDanhCaiXacBossBiTuChoi: boss da chet (TROI, hp=0) van con tren man hinh
// (DangHien) nhung KHONG duoc an them don nao — day chinh la chot chong chia
// thuong hai lan cua Task 12. Im lang, dung y het cach mot NGUOI CHOI da chet
// bi tu choi (kiem #4 cua KiemDanhHopLe).
func TestDanhCaiXacBossBiTuChoi(t *testing.T) {
	b := bossDanhDuoc(BossTroi, 0, 100)
	k, pA, pB, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if len(pA.nhan) != 0 || len(pB.nhan) != 0 {
		t.Fatalf("danh cai xac boss: pA=%d pB=%d goi, mong doi ca hai deu 0", len(pA.nhan), len(pB.nhan))
	}
	if nvA.KinhNghiem != 0 {
		t.Fatalf("KinhNghiem cua A = %d, mong doi 0 (danh cai xac KHONG duoc chia thuong lan hai)", nvA.KinhNghiem)
	}
}

// TestDanhBossChetRoiThiDonSauKhongAnThem chay DUNG kich ban that cua chot
// chong chia thuong hai lan: don thu nhat ket lieu boss, don thu hai (sau khi
// hoi chieu) roi vao cai xac -> khong goi tin, khong thuong them.
func TestDanhBossChetRoiThiDonSauKhongAnThem(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 2, 100)
	k, pA, _, nvA := khuVucDanhBoss(t, b)
	moc := mocDanhBoss
	k.dongHo = func() time.Time { return moc }

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})
	if nvA.KinhNghiem != 100 {
		t.Fatalf("chuan bi test sai: KinhNghiem sau don ket lieu = %d, mong doi 100", nvA.KinhNghiem)
	}
	soGoiSauDonMot := len(pA.nhan)

	moc = moc.Add(2 * time.Second) // qua hoi chieu, don thu hai hop le ve moi mat khac
	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if nvA.KinhNghiem != 100 {
		t.Fatalf("KinhNghiem sau don thu hai = %d, mong doi VAN 100 (khong chia thuong hai lan)", nvA.KinhNghiem)
	}
	if len(pA.nhan) != soGoiSauDonMot {
		t.Fatalf("A nhan them %d goi cho don vao cai xac, mong doi 0", len(pA.nhan)-soGoiSauDonMot)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 4 — cac chot chung cua nhanh -60 phai ap dung y het cho huong boss
// ─────────────────────────────────────────────────────────────────────────────

// TestDanhBossNguoiDanhDaChetBiTuChoi: A dang cho hoi sinh (HpHienTai<=0)
// khong duoc "tranh thu" danh boss — cung luat da ap cho danh quai va PvP.
func TestDanhBossNguoiDanhDaChetBiTuChoi(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, pA, _, nvA := khuVucDanhBoss(t, b)
	nvA.HpHienTai = 0

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (nguoi danh da chet)", b.Hp)
	}
	if len(pA.nhan) != 0 {
		t.Fatalf("A nhan %d goi, mong doi 0", len(pA.nhan))
	}
}

// TestDanhBossBangChieuHoiMauKhongGaySatThuong: chieu Loai=loaiBuff (Y Lieu
// Nhan Thuat) tuyet doi khong duoc xu nhu don danh — dung ly do va dung cach
// xu ly voi nhanh PvP (dac ta hoi mau con thieu, xem chu thich trong khuvuc.go).
func TestDanhBossBangChieuHoiMauKhongGaySatThuong(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, pA, _, nvA := khuVucDanhBoss(t, b)
	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 3}) // slot 2: hoi mau dong doi
	if k.nguoiTrong[nvA.ID].loaiChieuDangChon != loaiBuff {
		t.Fatalf("chuan bi test sai: loaiChieuDangChon = %d, mong doi %d (loaiBuff)",
			k.nguoiTrong[nvA.ID].loaiChieuDangChon, loaiBuff)
	}
	pA.nhan = nil

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (chieu hoi mau khong duoc lam boss mat mau)", b.Hp)
	}
	if len(pA.nhan) != 0 {
		t.Fatalf("A nhan %d goi, mong doi 0", len(pA.nhan))
	}
}

// TestDanhBossVuaBienHinhDangBatTu la KIEM THU DOT BIEN cua truong LucHienPha:
// boss vua vao pha moi duoc mot khoang bat tu 1500ms (dac ta §1.7), giong het
// nguoi choi/quai vua hoi sinh. Bo truong nay khoi ViTriCoMau thi boss bi ban
// nat ngay giay dau tien cua moi lan bien hinh.
func TestDanhBossVuaBienHinhDangBatTu(t *testing.T) {
	b := bossDanhDuoc(BossHien, 1000, 0)
	b.LucHienPha = mocDanhBoss // vua vao pha DUNG luc nay -> con nguyen 1500ms bat tu
	k, pA, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (dang bat tu sau khi vao pha)", b.Hp)
	}
	if len(pA.nhan) != 0 {
		t.Fatalf("A nhan %d goi, mong doi 0 (4 kiem dau hong -> im lang, khong phat goi)", len(pA.nhan))
	}
}

// TestDanhBossHetBatTuThiDanhDuoc la mat con lai cua test tren: qua 1500ms thi
// don danh vao lai binh thuong. Thieu no thi mot cai chot bat tu qua tay
// (chan vinh vien) van xanh.
func TestDanhBossHetBatTuThiDanhDuoc(t *testing.T) {
	b := bossDanhDuoc(BossHien, 1000, 0)
	b.LucHienPha = mocDanhBoss.Add(-1600 * time.Millisecond)
	k, pA, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 998 {
		t.Fatalf("Hp boss = %d, mong doi 998 (da qua 1500ms bat tu)", b.Hp)
	}
	if len(pA.nhan) != 1 {
		t.Fatalf("A nhan %d goi, mong doi 1", len(pA.nhan))
	}
}

// TestDanhBossVuotTamVanPhatGoiMiss khoa dac ta §1.9 cho huong boss: vuot tam
// -> sat thuong 0, chiMang=false, NHUNG VAN PHAT GOI (im lang lam nguoi choi
// tuong lag). Boss KHONG duoc mat mau.
func TestDanhBossVuotTamVanPhatGoiMiss(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	b.ViTriX = 300 // tam can chien chi 100
	k, pA, pB, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (vuot tam khong gay sat thuong)", b.Hp)
	}
	if len(b.SatThuongTheoNguoi) != 0 {
		t.Fatalf("so dong gop = %v, mong doi rong (don truot khong tinh dong gop)", b.SatThuongTheoNguoi)
	}
	muon := goi60DanhBoss(0, 0, 0) // satThuong 0, daChet 0, chiMang 0 (client hien "miss")
	for ten, p := range map[string]*phienGiaKhuVuc{"A": pA, "B": pB} {
		if len(p.nhan) != 1 || !bytes.Equal(p.nhan[0].Data, muon) {
			t.Fatalf("%s: goi -60 vuot tam sai:\n got  = %+v\n muon = %v", ten, p.nhan, muon)
		}
	}
}

// TestDanhBossChuaChonChieuBiTuChoi: A chua bam chon chieu nao -> KiemDanhHopLe
// tu choi ngay o bao ve dau ham (Chieu == nil) va nhanh boss khong duoc
// nil-deref khi doc SkillId.
func TestDanhBossChuaChonChieuBiTuChoi(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return mocDanhBoss }
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 150, HpToiDa: 150}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	k.themBoss(b)
	pA.nhan = nil

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (chua chon chieu)", b.Hp)
	}
	if len(pA.nhan) != 0 {
		t.Fatalf("A nhan %d goi, mong doi 0", len(pA.nhan))
	}
}

// TestDanhBossHoiChieuChuaXongBiTuChoi: hai don lien tiep trong cung mot moc
// thoi gian -> don thu hai bi kiem hoi chieu (#3) tu choi. Bao dam nhanh boss
// KHONG mo mot duong vong qua hoi chieu (danh boss 10 don/giay).
func TestDanhBossHoiChieuChuaXongBiTuChoi(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, pA, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})
	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 998 {
		t.Fatalf("Hp boss = %d, mong doi 998 (chi don DAU trung, don sau bi hoi chieu chan)", b.Hp)
	}
	if len(pA.nhan) != 1 {
		t.Fatalf("A nhan %d goi, mong doi 1 (don thu hai im lang)", len(pA.nhan))
	}
}

// TestDanhBossKhongDuChakraBiTuChoi: boss duoc mien kiem chakra (VoHanChakra)
// nhung NGUOI CHOI thi khong — nguoi danh boss van phai du chakra.
func TestDanhBossKhongDuChakraBiTuChoi(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, pA, _, nvA := khuVucDanhBoss(t, b)
	nvA.ChakraHienTai = 0 // chieu can chien cap 1 ton 10

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 1000 {
		t.Fatalf("Hp boss = %d, mong doi 1000 (nguoi danh het chakra)", b.Hp)
	}
	if len(pA.nhan) != 0 {
		t.Fatalf("A nhan %d goi, mong doi 0", len(pA.nhan))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 5 — he nguyen to va du lieu pha hong
// ─────────────────────────────────────────────────────────────────────────────

// TestDanhBossApHeNguyenToCuaPha: he nguyen to cua MUC TIEU lay tu Lang cua
// PHA dang chay (khong phai mot hang so trung tinh nhu quai). Ba nhanh khac
// che cho ba con so khac nhau — thieu test nay thi ai do gan cung
// heQuaiChuaXacDinh cho boss ma khong test nao do, va ca he ngu hanh bien mat
// khoi tran boss.
//
// Nguoi danh: Lang 2 (Suong Mu = He Thuy), chieu tam xa cap 1 (SatThuong=30,
// skillId=28) -> sat thuong tho = 16 x 30 / 100 = 4.
func TestDanhBossApHeNguyenToCuaPha(t *testing.T) {
	cac := []struct {
		langBoss int16
		mongDoi  int32
		vaiTro   string
	}{
		{0, 5, "Thuy khac Hoa -> 125%"},  // 4 * 125 / 100 = 5
		{3, 3, "Tho khac Thuy -> 80%"},   // 4 * 80 / 100 = 3
		{1, 4, "khong khac che -> 100%"}, // 4 * 100 / 100 = 4
	}
	for _, tt := range cac {
		t.Run(tt.vaiTro, func(t *testing.T) {
			b := bossDanhDuoc(BossHoatDong, 1000, 0)
			b.DuLieu[0].Lang = tt.langBoss
			k := MoKhuVuc(TimBanDo(1))
			k.dongHo = func() time.Time { return mocDanhBoss }
			nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 2, ChakraHienTai: 100, HpHienTai: 150, HpToiDa: 150}
			k.xuLy(suKienVao{nv: nvA, phien: &phienGiaKhuVuc{}})
			k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 2}) // slot 1: tam xa
			k.themBoss(b)

			k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

			if got := 1000 - b.Hp; got != tt.mongDoi {
				t.Fatalf("%s: sat thuong = %d, mong doi %d", tt.vaiTro, got, tt.mongDoi)
			}
		})
	}
}

// TestDanhBossPhaHongVanDanhDuocVaKhongPanic: PhaHienTai nam ngoai mien du
// lieu (du lieu noi dung hong) -> PhaHienTaiData() tra nil. Don danh KHONG
// duoc lam sap actor; he nguyen to roi ve trung tinh (100%) thay vi doan bua.
func TestDanhBossPhaHongVanDanhDuocVaKhongPanic(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	b.PhaHienTai = 99 // ngoai [0, len(DuLieu))
	k, _, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 998 {
		t.Fatalf("Hp boss = %d, mong doi 998 (he trung tinh 100%%: 12 x 20 / 100 = 2)", b.Hp)
	}
}

// TestDanhBossBoQuaOTrongTrongDanhSach: k.boss chua mot o nil (ve ly thuyet
// khong xay ra vi themBoss chan nil, nhung vong tim phai chiu duoc) — khong
// panic, van tim thay con boss that o sau.
func TestDanhBossBoQuaOTrongTrongDanhSach(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	k, _, _, nvA := khuVucDanhBoss(t, b)
	k.boss = append([]*Boss{nil}, k.boss...)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 998 {
		t.Fatalf("Hp boss = %d, mong doi 998 (o nil trong danh sach khong duoc chan duong)", b.Hp)
	}
}
