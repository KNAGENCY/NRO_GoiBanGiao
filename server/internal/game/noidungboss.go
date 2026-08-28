// noidungboss.go — NOI DUNG boss: hai con boss mau, du lieu THUAN Go
// (KH6b Task 15, nua mot).
//
// PHAM VI FILE NAY: chi DU LIEU va hai ham tra cuu. KHONG co may trang thai,
// KHONG chon ngau nhien, KHONG cham dong ho — "khi nao sinh, o map nao" la
// viec cua bo dieu phoi lich (lichboss.go, task 13) va may trang thai trong
// actor, dung theo CHAN-7 cua ke hoach. File nay chi tra loi cau hoi "con boss
// do gom nhung gi".
//
// VI SAO DU LIEU NAM TRONG CODE Go CHU KHONG PHAI YAML: thiet ke §8 chot
// chuyen kho noi dung sang YAML o KH8. Truoc do, de trong Go co mot cai loi
// that su: trinh bien dich va bo test kiem duoc tung con so ngay luc build.
// Bo test di kem (noidungboss_test.go) chinh la thu ngan mot con boss hong
// CAM lot vao game — loai loi khong lam sap gi ca, chi lam boss dung im, vo
// hinh, hoac sinh trong long dat, va chi lo ra khi co nguoi that vao nhin.
//
// ⚠️ MOI CON SO O DAY DEU TRUY DUOC VE MOT NGUON TRONG CODE. Khi sua, sua
// kem cai nguon do — dung dien mot con so "cho dep".
package game

// ---------------------------------------------------------------------------
// 1. Hai kieu bao ngoai — cho nhung gi PhaBoss khong cho duoc
// ---------------------------------------------------------------------------

// ViTriSinh la MOT cho dat chan hop le cua boss: map nao, toa do nao.
//
// VI SAO KHONG NAM TRONG PhaBoss: PhaBoss.MapVao chi noi boss hien o MAP nao,
// khong noi hien o CHO nao trong map. Ma cho dat chan lai la thu duy nhat
// trong ca bo du lieu nay co the kiem bang mot nguon KHACH QUAN (luoi va cham
// cua map), nen no dang duoc mot kieu rieng de bo test bat duoc.
//
// Toa do la "chan nhan vat dung o dau", cung quy uoc voi DiemDi.DenX/DenY
// (banDo.go) va voi luoi va cham: o chua toa do phai DI DUOC, va o ngay duoi
// phai DAC (co mat dat). Sinh giua troi thi boss roi/lo lung tuy client.
type ViTriSinh struct {
	Map  int32
	X, Y int16
}

// NoiDungBoss la MOT con boss trong kho noi dung: cac pha cua no + cac cho no
// duoc phep dat chan. Ma la khoa tra cuu on dinh (bo dieu phoi lich goi boss
// theo ma, khong theo chi so mang — chen mot con vao giua danh sach khong
// duoc lam doi y nghia cua lich).
type NoiDungBoss struct {
	Ma    string
	Pha   []PhaBoss
	ViTri []ViTriSinh
}

// ---------------------------------------------------------------------------
// 2. Nhung con so dung chung — ghi ro lay tu dau
// ---------------------------------------------------------------------------

// Ba mieng do da co ANH NARUTO THAT trong client. Nguon: closure.part_ids
// trong server/internal/game/assets/asset-lock.json = [30, 57, 58, 2189, 2190,
// 2191]. Dien mot id ngoai tap do vao Outfit la con boss VO HINH — client
// khong tim thay anh thi khong ve gi, va may chu khong co mot dong log nao.
//
// Bo "hoc vien": 30/57/58 (PHASE_ONE_REPLACEMENTS trong assetpack/part.py) —
// ba mieng dau/than/chan cua bo do hoc vien nhan gia.
const (
	dauHocVien  int16 = 30
	thanHocVien int16 = 57
	chanHocVien int16 = 58
)

// Bo "jounin": 2189/2190/2191 — bo cua NPC Iruka (assets/catalog.json, va
// dongbo_tichhop_test.go:297 ghim dung ba so nay xuong day). Day la bo do
// nguoi lon duy nhat da co anh that.
//
// ⚠️ NO KHACH: hai bo tren la TAT CA nhung gi client co. Chua co bo do rieng
// cho boss — mot con boss mac bo jounin trong y het NPC Iruka. Doi la chuyen
// cua kho my thuat (KH7/KH8), khong phai chuyen cua file nay; o day chon
// "nhin thay duoc" thay vi "dep ma vo hinh".
const (
	dauJounin  int16 = 2189
	thanJounin int16 = 2190
	chanJounin int16 = 2191
)

// Hai templateId chieu ma boss duoc dung. Nguon: chieuthuc.go — moi lang co
// dung 4 chieu, templateId 1..4 theo slot:
//
//	1 = can chien (loai DANH, tam 100)
//	2 = tam xa theo he (loai DANH, tam 300)
//	3 = hoi mau dong doi (loai BUFF)  -> boss KHONG dung: khong sinh goi tan
//	    cong nao, va duong hoi mau dang co y hoan (BAT-DAU-TU-DAY.md §8)
//	4 = tu than (loai TU_THAN)        -> boss KHONG dung: tamTuThan = 0 nen
//	    moi don deu vuot tam, tuc LUON truot
//
// Boss cam chieu 3 hoac 4 se "vung tay vao khong khi": trong nhu dang danh
// nhau ma khong bao gio gay mot diem sat thuong nao.
const (
	chieuCanChien int32 = 1
	chieuTamXa    int32 = 2
)

// Map 1 la Lang La — map DUY NHAT chon lam san khau cho ca hai con boss mau.
// Ba ly do, theo thu tu quan trong:
//
//  1. No la map duy nhat co LUOI VA CHAM THAT (assets/map-1-collision.bin).
//     Map 2-5 khong co artifact va cham nao, dichuyen.go phai dung bien du
//     phong cho chung — dat boss o do thi KHONG AI kiem duoc no co dung trong
//     vach da khong, va loi kieu do khong bao gio bao.
//  2. No la map duy nhat da co my thuat Naruto that; map 2-5 con dung tile
//     Dragon Ball cu, tha boss vao do trong lac cho.
//  3. No la map khoi dau cua Lang La (lang.go), tuc noi dong nguoi nhat —
//     boss mau sinh ra de NGHIEM THU, can nguoi nhin thay.
const mapLangLa int32 = 1

// ---------------------------------------------------------------------------
// 3. Can bang so lieu — day la GIA TRI KHOI DAU, khong phai so da chot
// ---------------------------------------------------------------------------
//
// Chua co he thong cap do/chi so nao chot con so nao ca (thiet ke §8: "can
// bang chi tiet — chot khi co content"). Nhung cac so duoi day KHONG bia: moi
// so deu suy ra tu con so co that trong code, va phep tinh ghi ngay canh de
// nguoi sau kiem lai duoc.
//
// BA CON SO GOC:
//
//	mau nguoi choi luc moi tao : 100..250   (lang.go, HpKhoiDau tung lang)
//	sat thuong goc nguoi choi  : 10..16     (lang.go, SatThuongKhoiDau)
//	mau con Soi Rung           : 100        (assets/catalog.json, mobs[0].hp)
//
// CACH SAT THUONG DUOC TINH (satthuong.go): satThuong = Dame × CapChieu.SatThuong / 100.
// CapChieu.SatThuong sinh boi taoCapChieu (chieuthuc.go) = dameGoc + (cap-1)×20,
// voi dameGoc = 20 cho can chien va 30 cho tam xa. Tuc:
//
//	can chien cap 1/2/3 -> 20% / 40% / 60% cua Dame
//	tam xa    cap 1/2/3 -> 30% / 50% / 70% cua Dame
//
// TU DO SUY RA:
//
//	Dame boss 35  -> don 10..21   : nguoi choi 150 mau chiu ~10 don. Boss le.
//	Dame boss 60  -> don 24..30   : chiu ~5 don. Boss lon pha 1.
//	Dame boss 90  -> don 54..63   : chiu ~3 don. Boss lon pha 2 — phai danh
//	                                theo nhom, dung y do.
//
// MAU BOSS suy nguoc tu suc ra don cua nguoi choi: sat thuong goc 10..16 ×
// chieu cap 3 (60-70%) ≈ 7..11 mot don, hoi chieu 1,6 giay -> ~5 mau/giay moi
// nguoi. Vi vay:
//
//	400..800    ≈ 1,5..3 phut cho MOT nguoi   -> boss le, danh solo duoc
//	2400..3600  ≈ 2..3 phut cho NAM nguoi     -> boss lon pha 1
//	4000..6000  ≈ 3..4 phut cho NAM nguoi     -> boss lon pha 2
//
// Ba day mau nay deu nam thoai mai trong mien int32 cua Boss.Hp/HpToiDa
// (tran 2.147.483.647) — bo test ghim dieu do bang chinh KepHp.

// ---------------------------------------------------------------------------
// 4. Hai con boss mau
// ---------------------------------------------------------------------------

// MaBossNhanGiaMatNa / MaBossHocVienDaoTau: khoa tra cuu cua hai con boss mau.
// Hang so (khong phai chuoi tran) de bo dieu phoi lich va bo test cung tro toi
// MOT cho — go nham mot chu trong chuoi tran se im lang tra nil.
const (
	MaBossNhanGiaMatNa  = "nhan-gia-mat-na"
	MaBossHocVienDaoTau = "hoc-vien-dao-tau"
)

// danhSachNoiDungBoss la kho noi dung boss. Hai con, chon co chu dich de giua
// chung chay HET ba duong dac trung cua he boss — neu khong thi buoi nghiem
// thu bang client that (nua hai cua task 15) khong chung minh duoc gi:
//
//	nhieu pha  -> Nhan Gia Mat Na co 2 pha (duong bien hinh, giu nguyen id)
//	noi gian   -> ca hai con deu co MocNoiGian theo %mau
//	boss lon   -> Nhan Gia Mat Na dat LaBossLon (duong thong bao chay chu 93)
var danhSachNoiDungBoss = []NoiDungBoss{
	{
		// ── BOSS LON: Nhan Gia Mat Na ──────────────────────────────────────
		//
		// Mot ke doi lot hoc vien tron trong Lang La ca mua; danh gan chet
		// thi mat na vo, chan than buoc ra. Hai pha nay chinh la ly do con
		// boss nay ton tai: no la con DUY NHAT chay duong bien hinh.
		//
		// VI SAO CAU CHUYEN LAI LA "LOT XAC" chu khong phai mot con quai vat
		// bien hinh: client chi co DUNG HAI bo do (hoc vien va jounin, xem
		// hang so o tren). Voi hai bo do do, "hoc vien gia -> nguoi lon that"
		// la cau chuyen DUY NHAT ma my thuat hien co ke duoc that. Neu viet
		// mot con ho chin duoi thi pha 2 se hien ra... mot bac thay Iruka.
		Ma: MaBossNhanGiaMatNa,
		ViTri: []ViTriSinh{
			// Giua be da duoi cong torii, ngay tim lang — cho ai cung di qua.
			// Cot 23 hang 17 cua luoi va cham: dung tren MAT be, hang 18 ngay
			// duoi la than be (dac). Cung mat be ma cong Lang La dung
			// (banDo.go, diemDiMap1).
			{Map: mapLangLa, X: 552, Y: 408},
			// Bai dat phia dong, ngoai be da — cot 37 hang 18, mat dat thuong,
			// hang 19 ben duoi la long dat (dac). Co hai cho de con boss
			// khong phai lan nao cung hien dung mot diem.
			{Map: mapLangLa, X: 888, Y: 432},
		},
		Pha: []PhaBoss{
			{
				Ten:  "Kẻ Lạ Mặt",
				Lang: 0, // Lang La — he Hoa. Chieu lay tu catalog cua lang 0.
				Cap:  50,
				// Bo do hoc vien: dang tro lot cua pha 1. Nguoi choi nhin vao
				// tuong mot hoc vien binh thuong — dung y do.
				Outfit: [6]int16{dauHocVien, thanHocVien, chanHocVien, 0, 0, 0},
				Dame:   60, // -> don 24..30, nguoi choi 150 mau chiu ~5 don
				TypePk: TypePkBoss,

				// 1500ms = 15 nhip actor. Boi so cua 100 la BAT BUOC chu khong
				// phai cho tron: boss cap nhat trong nhip 100ms cua khu vuc,
				// nen 1450ms tren thuc te van la 1500ms — phan le la con so
				// noi doi.
				NhipDanhMs: 1500,

				HpChon:  []int64{2400, 3000, 3600},
				MapVao:  []int32{mapLangLa},
				BoChieu: [][2]int32{{chieuCanChien, 2}, {chieuTamXa, 2}},

				// 1800 giay = 30 phut giua hai lan boss lon xuat hien. Doi
				// chieu: boss le duoi kia la 10 phut.
				GiayNghi: 1800,

				// Hai moc, %HP GIAM DAN va khong trung nguong — MocNoiGianDaQua
				// (boss.go) la con tro chay MOT CHIEU, sap sai thu tu thi moc
				// sau khong bao gio toi luot.
				//
				// HeSoNhip nhan DON len NhipDanhHienTai (cong don):
				// 1500 -> ×0,8 = 1200 -> ×0,75 = 900. Ca ba deu la boi so cua
				// 100ms, tuc con so trong du lieu dung bang con so xay ra that.
				MocNoiGian: []MocNoiGian{
					{PhanTramHP: 50, HeSoDame: 1.30, HeSoDef: 0.85, HeSoNhip: 0.80,
						Chat: "Đủ rồi. Thôi giấu tay."},
					{PhanTramHP: 20, HeSoDame: 1.60, HeSoDef: 0.70, HeSoNhip: 0.75,
						Chat: "Ngươi ép ta cởi lớp này ra..."},
				},

				// Boss lon: phat goi 93 (chay chu toan server) luc hien va luc
				// bi ha. Dat o CA HAI pha vi nguoi choi dang o map khac cung
				// nen biet con nay da lot xac.
				LaBossLon: true,

				ThoaiHien: []string{
					"Làng Lá... vẫn ngủ say như cái đêm đó.",
					"Ta đứng giữa các ngươi suốt mùa này. Không ai ngoảnh lại nhìn.",
				},
				ThoaiDanh: []string{
					"Chakra của ngươi run rồi đấy.",
					"Đừng nhắm mắt. Ta muốn ngươi nhìn cho rõ.",
					"Học viện dạy các ngươi kết ấn, có ai dạy các ngươi chết đâu.",
				},
				ThoaiChet: []string{
					"Đó mới chỉ là lớp áo... còn ta thì chưa.",
				},
			},
			{
				Ten:  "Mặt Nạ Vỡ",
				Lang: 0,
				Cap:  70, // cao hon pha 1 — nguoi choi thay ngay con so tren dau boss doi
				// Bo jounin: chan than sau khi mat na vo.
				Outfit: [6]int16{dauJounin, thanJounin, chanJounin, 0, 0, 0},
				Dame:   90, // -> don 54..63: mot nguoi don doc khong tru noi
				TypePk: TypePkBoss,

				NhipDanhMs: 1000, // 10 nhip actor — nhanh hon pha 1 (15 nhip)

				HpChon:  []int64{4000, 5000, 6000},
				MapVao:  []int32{mapLangLa},
				BoChieu: [][2]int32{{chieuCanChien, 3}, {chieuTamXa, 3}}, // cap 3 = cap cao nhat catalog co

				// 5 giay lang giua hai pha: mat na roi xuong dat, im mot nhip
				// roi chan than buoc ra. Khong de 0 — 0 nghia la pha 2 no ra
				// ngay trong cung mot khoanh khac pha 1 chet, nguoi choi khong
				// kip doc gi.
				GiayNghi: 5,

				// 1000 -> ×0,8 = 800 -> ×0,75 = 600. Van boi so cua 100ms.
				MocNoiGian: []MocNoiGian{
					{PhanTramHP: 50, HeSoDame: 1.35, HeSoDef: 0.80, HeSoNhip: 0.80,
						Chat: "Chakra này ta để dành cho đúng hôm nay."},
					{PhanTramHP: 25, HeSoDame: 1.75, HeSoDef: 0.65, HeSoNhip: 0.75,
						Chat: "Cháy đi. Cả cái làng này cùng cháy."},
				},

				LaBossLon: true,

				ThoaiHien: []string{
					"Mặt nạ rơi rồi. Giờ thì nhìn cho kỹ.",
					"Ta từng đeo băng trán của làng này. Trước khi làng này bỏ ta lại.",
				},
				ThoaiDanh: []string{
					"Hoả Độn đâu phải của riêng Làng Lá.",
					"Bao nhiêu người nữa thì các ngươi mới chịu tin?",
					"Đứng dậy đi. Ta chưa muốn xong sớm thế đâu.",
				},
				ThoaiChet: []string{
					"Cuối cùng... cũng có người nhớ mặt ta.",
				},
			},
		},
	},
	{
		// ── BOSS LE: Hoc Vien Dao Tau ──────────────────────────────────────
		//
		// Mot con boss NHO, MOT pha, KHONG chay chu toan server — co y de doi
		// chieu voi con tren. No ton tai de tra loi hai cau hoi khac nhau:
		// "duong noi gian co chay khi chi co MOT moc khong" va "boss khong
		// phai boss lon co im lang dung khong" (LaBossLon = false, mac dinh).
		//
		// Lang 2 (Suong Mu — he Thuy) nhung dat chan o map 1 (Lang La): mot ke
		// dao tau khoi ky thi cua Suong Mu, tron sang dat Lang La. Chon lang
		// khac lang cua map la CO CHU DICH — no bat bo chieu phai tra cuu dung
		// catalog cua lang 2 chu khong phai lang cua map, va bo test khoa dieu
		// do lai.
		Ma: MaBossHocVienDaoTau,
		ViTri: []ViTriSinh{
			// Cot 13 hang 18 — bai dat phia tay, cach xa NPC Iruka (172,384)
			// va hai con Soi Rung de nguoi moi khong dam thang vao boss.
			{Map: mapLangLa, X: 312, Y: 432},
			// Cot 42 hang 18 — bia dong cua map.
			{Map: mapLangLa, X: 1008, Y: 432},
		},
		Pha: []PhaBoss{
			{
				Ten:  "Học Viên Đào Tẩu",
				Lang: 2, // Lang Suong Mu — he Thuy
				Cap:  25,
				// Van bo do hoc vien: han bo tron giua ky thi, chua kip doi do.
				// Day la lan DUY NHAT bo do hoc vien duoc dung dung nghia den.
				Outfit: [6]int16{dauHocVien, thanHocVien, chanHocVien, 0, 0, 0},
				Dame:   35, // -> don 10..21: nguoi choi 150 mau chiu ~10 don
				TypePk: TypePkBoss,

				NhipDanhMs: 1200, // 12 nhip actor

				HpChon: []int64{400, 600, 800},
				MapVao: []int32{mapLangLa},
				// Can chien cap 2 + tam xa cap 1: yeu hon han bo chieu cua boss
				// lon, dung cho mot con danh solo duoc.
				BoChieu: [][2]int32{{chieuCanChien, 2}, {chieuTamXa, 1}},

				GiayNghi: 600, // 10 phut — boss le nen ra thuong xuyen hon

				// MOT moc duy nhat: 1200 -> ×0,75 = 900ms.
				MocNoiGian: []MocNoiGian{
					{PhanTramHP: 40, HeSoDame: 1.50, HeSoDef: 0.80, HeSoNhip: 0.75,
						Chat: "Đừng ép ta nhớ lại kỳ thi đó."},
				},

				// LaBossLon de mac dinh false — con nay ra vao im lang, khong
				// lam ngap kenh chay chu toan server.

				ThoaiHien: []string{
					"Đừng gọi tên cũ của ta. Đứa đó chết trong kỳ thi rồi.",
					"Ta bỏ Sương Mù để khỏi phải giết bạn học. Rồi ta vẫn giết.",
				},
				ThoaiDanh: []string{
					"Ngươi cũng đeo băng trán à? Tháo ra đi, đỡ đau.",
					"Ở Sương Mù, ai ngã xuống là nằm lại đó.",
					"Nước không cần mạnh. Nước chỉ cần lâu.",
				},
				ThoaiChet: []string{
					"Vậy là... ta được về nhà theo cách này.",
				},
			},
		},
	},
}

// ---------------------------------------------------------------------------
// 5. Hai cua tra cuu
// ---------------------------------------------------------------------------

// saoNoiDungBoss tra ban sao SAU cua mot NoiDungBoss.
//
// VI SAO PHAI SAO SAU: sao mot struct trong Go chi sao CON TRO cua lat cat,
// nen nguoi goi van ghi de duoc du lieu goc qua ban "sao chep" — con boss ma
// dong tren dinh cam, dung bay ma banDo.go da tung phai va (xem TimBanDo).
// Voi boss thi hau qua nang hon: du lieu bi sua chi lo ra o LAN SINH SAU, o
// mot khu vuc khac, rat lau sau khi cho sua tay xay ra.
//
// KHONG sao sau [6]int16 va cac truong so: mang co do dai co dinh trong Go la
// GIA TRI, copy struct da chep tron.
func saoNoiDungBoss(goc NoiDungBoss) NoiDungBoss {
	cp := goc
	cp.ViTri = append([]ViTriSinh(nil), goc.ViTri...)
	cp.Pha = make([]PhaBoss, len(goc.Pha))
	for i, pha := range goc.Pha {
		p := pha
		p.HpChon = append([]int64(nil), pha.HpChon...)
		p.MapVao = append([]int32(nil), pha.MapVao...)
		p.BoChieu = append([][2]int32(nil), pha.BoChieu...)
		p.ThoaiHien = append([]string(nil), pha.ThoaiHien...)
		p.ThoaiDanh = append([]string(nil), pha.ThoaiDanh...)
		p.ThoaiChet = append([]string(nil), pha.ThoaiChet...)
		p.MocNoiGian = append([]MocNoiGian(nil), pha.MocNoiGian...)
		cp.Pha[i] = p
	}
	return cp
}

// TimNoiDungBoss tra ban sao cua con boss mang ma nay, hoac nil neu khong co.
//
// Fail-closed (tra nil thay vi mot NoiDungBoss rong): mot NoiDungBoss rong se
// di tiep vao may trang thai va sinh ra con boss khong pha, khong mau, khong
// ten — hong CAM. nil thi noi goi buoc phai kiem.
func TimNoiDungBoss(ma string) *NoiDungBoss {
	for i := range danhSachNoiDungBoss {
		if danhSachNoiDungBoss[i].Ma == ma {
			cp := saoNoiDungBoss(danhSachNoiDungBoss[i])
			return &cp
		}
	}
	return nil
}

// DanhSachNoiDungBoss tra ban sao SAU toan bo kho noi dung boss (chi doc).
func DanhSachNoiDungBoss() []NoiDungBoss {
	ds := make([]NoiDungBoss, len(danhSachNoiDungBoss))
	for i := range danhSachNoiDungBoss {
		ds[i] = saoNoiDungBoss(danhSachNoiDungBoss[i])
	}
	return ds
}
