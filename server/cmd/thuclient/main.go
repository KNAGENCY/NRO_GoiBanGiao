// Chuong trinh thuclient: mot "client gia" doc lap, lam DUNG NHUNG GI client
// Unity that lam khi ket noi va dang nhap — de kiem chung may chu Go noi
// dung giao thuc ma khong can mo game that.
//
// Chuong trinh nay CO Y khong dung "nroserver/internal/net" hay
// "nroserver/internal/game" — do la ma phia may chu. No chi dung
// "nroserver/internal/protocol" (bo ma hoa/khung tin da duoc kiem chung khop
// voi ban Java/Unity cu), dung y het cach mot client that (Session_ME.cs) tu
// cai dat logic tuong tu trong ngon ngu cua no — khong "goi lai" ma nguon may
// chu.
//
// Cac hang so giao thuc duoi day lay tu ban Java/Unity cu, DUNG DOI:
//
//	-27  goi bat tay dau tien: client gui rong, may chu tra ve khoa tho
//	-29  khung "chua dang nhap", than tin: 1 byte lenh con roi du lieu
//	  0  lenh con: dang nhap
//	-28  hop thu "khong tren ban do", than tin: 1 byte lenh con roi du lieu
//	  2  lenh con cua khung -28: tao nhan vat (ten, lang, kieu toc)
//	  4  lenh con cua khung -28: versionGame (may chu -> client, Ke hoach 3):
//	     vsData/vsMap/vsSkill/vsItem, bao hieu "da vao the gioi"
//	  6  lenh con cua khung -28: client xin catalog map / may chu tra loi
//	  7  lenh con cua khung -28: client xin catalog skill / may chu tra loi
//	     (chua DANH SACH LANG — Ke hoach 3 thay 3 lop Dragon Ball cu bang 5
//	     lang Naruto, day la noi kiem duoc so lang/ten lang)
//	  8  lenh con cua khung -28: client xin catalog item / may chu tra loi
//	 13  lenh con cua khung -28: client bao da nhan du 4 catalog, san sang
//	-87  goi data (dart/arrow/effect/image/part/skill): client xin bang than
//	     rong, may chu tra loi day du — LENH DAC BIET, do dai BA byte (khac
//	     moi lenh khac o tren, du dai HAI byte) vi goi tra loi ~87 KB vuot qua
//	     gioi han cua truong do dai hai byte (65535) — xem lenhDoDaiBaByte
//	     trong protocol/frame.go
//	-26  may chu hien hop thoai co nut OK (dang nhap/tao nhan vat that bai se kem)
//	-102 may chu bao client thoat khoi trang thai "dang cho" (chi khung -29)
//	   2 may chu bao chuyen sang man tao nhan vat (= dang nhap thanh cong, chua co nhan vat)
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"nroserver/internal/protocol"
)

const (
	lenhBatTay          int8 = -27
	lenhChuaDangNhap    int8 = -29
	lenhConDangNhap     byte = 0
	lenhKhongTrenBanDo  int8 = -28 // hop thu "khong tren ban do" (Controller.java messageNotMap)
	lenhConTaoNhanVat   byte = 2   // lenh con cua khung -28: tao nhan vat
	lenhThongBaoOK      int8 = -26
	lenhDangNhapThatBai int8 = -102
	lenhTaoNhanVat      int8 = 2

	// Ke hoach 3: dong bo du lieu tinh, sau khi vaoTheGioi (dang nhap voi
	// nhan vat co san, hoac vua tao nhan vat xong) may chu gui -28/4
	// (versionGame) — xem cac lenh con duoi day va dongbo.go phia may chu
	// (internal/game/dongbo.go, cung ten hang so, DUNG DOI).
	lenhConVersionGame  byte = 4   // may chu -> client: bao 4 phien ban du lieu tinh
	lenhConCatalogMap   byte = 6   // client xin catalog map / may chu tra
	lenhConCatalogSkill byte = 7   // client xin catalog skill / may chu tra (DANH SACH LANG)
	lenhConCatalogItem  byte = 8   // client xin catalog item / may chu tra
	lenhConClientOk     byte = 13  // client bao da nhan du 4 catalog, san sang
	lenhDuLieu          int8 = -87 // goi data: LENH DAC BIET do dai BA byte (xem chu thich dau file)

	// hanChoPhanHoiTaoNhanVat: sau khi gui yeu cau tao nhan vat, doi toi da
	// ngan nay de xem may chu co tra loi gi khong. Ke hoach 2: tao nhan vat
	// THANH CONG khong gui goi nao ca (xem vaoTheGioi trong
	// internal/game/nhanvat.go) — cong cu chay tay nay khong the tu khang
	// dinh "thanh cong", chi phat hien duoc LOI ro rang (-26); het han ma
	// khong nhan duoc gi la dau hieu (khong chac chan) cua thanh cong, nguoi
	// dung phai tu doi chieu bang log may chu / psql.
	hanChoPhanHoiTaoNhanVat = 5 * time.Second

	// gioiHanThan chan than tin bao do dai khong lo, tranh cap phat het bo
	// nho neu noi nham vao mot thu khong phai may chu that.
	//
	// Truoc day chu thich nay so sanh voi mot hang so cung ten (doDaiThanToiDa)
	// ben protocol/frame.go — hang so do DA BI BO ben phia may chu vi la MA
	// CHET (truong do dai server doc chi co 2 byte, khong bao gio vuot 65535,
	// nen dieu kien so sanh voi 4 MB khong bao gio dung). Ghi chu lai o day de
	// tranh nham lan tuong tu: gioi han 4 MB o DAY (cong cu client gia) cung
	// dang trong tinh trang tuong tu — khong xoa vi ngoai pham vi ke hoach
	// nay (chi sua server/internal/protocol/frame.go), nhung nen luu y neu sau
	// nay dong dung cong cu nay.
	gioiHanThan = 4 << 20
)

func main() {
	diaChi := flag.String("dia-chi", "127.0.0.1:14445", "dia chi may chu (host:port)")
	tenDangNhap := flag.String("ten", "", "ten dang nhap (bat buoc)")
	matKhau := flag.String("mat-khau", "", "mat khau (bat buoc)")
	coTaoNhanVat := flag.Bool("tao-nhan-vat", false,
		"neu dang nhap thanh cong va may chu tra ve lenh 2 (chua co nhan vat), tu dong gui tiep goi tao nhan vat voi -ten-nhan-vat/-lang/-kieu-toc")
	tenNhanVat := flag.String("ten-nhan-vat", "Naruto01",
		"ten nhan vat se thu tao (5-10 ky tu, chi chu cai va chu so) — chi dung khi -tao-nhan-vat")
	langNhanVat := flag.Int("lang", 0,
		"id lang cho nhan vat moi (0..4, xem internal/game/lang.go) — chi dung khi -tao-nhan-vat")
	kieuTocNhanVat := flag.Int("kieu-toc", 5,
		"id kieu toc cho nhan vat moi — chi dung khi -tao-nhan-vat")
	flag.Parse()

	if *tenDangNhap == "" || *matKhau == "" {
		fmt.Fprintln(os.Stderr, "LOI: phai truyen ten dang nhap va mat khau.")
		fmt.Fprintln(os.Stderr, "Vi du:  thuclient -ten thutest -mat-khau 1")
		fmt.Fprintln(os.Stderr, "Xem tat ca tham so:  thuclient -h")
		os.Exit(1)
	}
	// lang/kieu-toc di xuong 1 byte tren day (xem guiGoiTaoNhanVat) — chan
	// som gia tri ngoai [0,255] thay vi de ep kieu am tham sai lech.
	if *coTaoNhanVat && (*langNhanVat < 0 || *langNhanVat > 255) {
		fmt.Fprintln(os.Stderr, "LOI: -lang phai trong khoang 0..255 (mot byte).")
		os.Exit(1)
	}
	if *coTaoNhanVat && (*kieuTocNhanVat < 0 || *kieuTocNhanVat > 255) {
		fmt.Fprintln(os.Stderr, "LOI: -kieu-toc phai trong khoang 0..255 (mot byte).")
		os.Exit(1)
	}

	fmt.Printf("Dang noi toi may chu tai %s ...\n", *diaChi)
	conn, err := net.Dial("tcp", *diaChi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "KHONG NOI DUOC TOI MAY CHU: %v\n", err)
		fmt.Fprintln(os.Stderr, "Kiem tra: may chu (nroserver) da chay chua? Dia chi/cong co dung khong?")
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Da noi duoc toi may chu. Bat dau bat tay trao khoa...")

	r := bufio.NewReader(conn)

	// Buoc 1: gui goi -27 RONG. Truoc khi co khoa, client that ghi do dai
	// theo LITTLE-ENDIAN (BinaryWriter.Write(ushort) cua .NET) — nhung than
	// tin rong nen ca hai byte deu la 0, thu tu khong anh huong o day.
	//
	// cmdBatTay di qua bien (khong convert thang tren hang so) vi Go chan
	// convert hang so am sang byte luc bien dich ("constant overflows byte");
	// convert luc chay thi hop le, day chinh la quy uoc byte(int8) ma ca
	// protocol/frame.go lan phan con lai cua file nay dang dung.
	cmdBatTay := lenhBatTay
	if _, err := conn.Write([]byte{byte(cmdBatTay), 0x00, 0x00}); err != nil {
		fmt.Fprintf(os.Stderr, "LOI khi gui goi bat tay: %v\n", err)
		os.Exit(1)
	}

	// Buoc 2: doc phan hoi cua may chu. Goi nay CHUA MA HOA (khoa chua duoc
	// dat), va may chu ghi do dai theo BIG-ENDIAN (kieu java.io.DataOutputStream).
	cmd, than, err := docGoiChuaMaHoa(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LOI khi doc goi bat tay tu may chu: %v\n", err)
		os.Exit(1)
	}
	if cmd != lenhBatTay {
		fmt.Fprintf(os.Stderr, "MAY CHU TRA VE LENH %d, MONG DOI %d — bat tay khong dung giao thuc.\n", cmd, lenhBatTay)
		os.Exit(1)
	}

	rd := protocol.NewReader(than)
	doDaiKhoa, err := rd.Byte()
	if err != nil {
		fmt.Fprintf(os.Stderr, "LOI: goi bat tay thieu do dai khoa: %v\n", err)
		os.Exit(1)
	}
	khoaTho := make([]byte, doDaiKhoa)
	for i := range khoaTho {
		b, err := rd.Byte()
		if err != nil {
			fmt.Fprintf(os.Stderr, "LOI: goi bat tay thieu byte khoa (thu %d/%d): %v\n", i+1, len(khoaTho), err)
			os.Exit(1)
		}
		khoaTho[i] = b
	}
	// 3 truong con lai trong goi (dia chi may chu phu) khong can dung o day —
	// du an nay chi chay mot may chu duy nhat.

	khoa := protocol.DeriveKey(khoaTho)
	ks := protocol.NewKeyStream(khoa)
	fmt.Printf("Da nhan khoa tho (%d byte) tu may chu va tao khoa lam viec. Bat tay xong.\n", len(khoaTho))

	// Buoc 3: gui goi dang nhap. Khung -29, than tin: 1 byte lenh con (0 =
	// dang nhap) roi 2 chuoi kieu UTF cua Java (ten, mat khau).
	w := protocol.NewWriter()
	w.Byte(lenhConDangNhap)
	if err := w.UTF(*tenDangNhap); err != nil {
		fmt.Fprintf(os.Stderr, "LOI: ten dang nhap khong hop le: %v\n", err)
		os.Exit(1)
	}
	if err := w.UTF(*matKhau); err != nil {
		fmt.Fprintf(os.Stderr, "LOI: mat khau khong hop le: %v\n", err)
		os.Exit(1)
	}
	if err := guiGoiDaMaHoa(conn, ks, lenhChuaDangNhap, w.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "LOI khi gui goi dang nhap: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Da gui goi dang nhap voi ten dang nhap: %q\n", *tenDangNhap)
	fmt.Println("Dang cho may chu phan hoi...")
	fmt.Println()

	// Buoc 4: doc va dien giai cac goi may chu tra ve.
	//
	// Neu -tao-nhan-vat duoc bat: khi thay lenh 2 (chua co nhan vat), THAY VI
	// dung lai nhu binh thuong, gui tiep goi -28/2 tao nhan vat roi doc THEM
	// mot vong nua de xem may chu tra loi gi:
	//   - -26 (loi)     -> dienGiaiGoi da in thong bao, dung lai luon (da ro).
	//   - -28/4 (versionGame) -> nhanh -28 duoi day tu xu ly va dung ca
	//     chuong trinh sau khi xong version dance (xem xuLyVersionDance).
	//   - het han doc   -> KHONG con la dau hieu thanh cong tu Ke hoach 3 (vaoTheGioi
	//     luon gui -28/4 khi thanh cong, xem internal/game/login.go) — van giu
	//     lai nhanh nay cho phong hong (vd may chu cu chua nang cap), nguoi
	//     dung phai tu doi chieu bang log may chu / psql neu gap truong hop nay.
	daGuiTaoNhanVat := false
	dangChoPhanHoiTaoNhanVat := false
	for {
		cmd, than, err := docGoiDaMaHoa(r, ks)
		if err != nil {
			var netErr net.Error
			if dangChoPhanHoiTaoNhanVat && errors.As(err, &netErr) && netErr.Timeout() {
				fmt.Println("Khong nhan duoc goi nao trong vong cho phan hoi tao nhan vat.")
				fmt.Println("Ke hoach 3: tao nhan vat THANH CONG phai tra ve -28/4 (versionGame); het han o day la BAT THUONG.")
				fmt.Println("=> Kiem tra log may chu (dong \"nhan vat vao the gioi\") va CSDL (psql) de xac nhan.")
				return
			}
			if err == io.EOF {
				fmt.Println("May chu da dong ket noi.")
			} else {
				fmt.Fprintf(os.Stderr, "LOI khi doc goi tu may chu: %v\n", err)
			}
			return
		}

		// Ke hoach 3: khung -28 (khong tren ban do) gio co them cac lenh con
		// lien quan dong bo tinh. Chi can xu ly rieng lenh con 4 (versionGame)
		// o day: no la tin hieu "da vao the gioi" (ca hai duong dang nhap co
		// san nhan vat LAN vua tao nhan vat xong deu di qua no, xem vaoTheGioi
		// trong internal/game/login.go) — tu do cong cu nay lam DUNG Y HET
		// client that: tu dong xin het 4 catalog roi bao san sang, de nguoi
		// dung doi chieu bang mat. Dung lai chuong trinh sau khi xong (Ke hoach
		// 3 dung o buoc gui -28/13 — Ke hoach 4 moi doc them ME_INFO/MAP INFO).
		if cmd == lenhKhongTrenBanDo {
			if len(than) == 0 {
				fmt.Println("[Lenh -28] goi thieu byte lenh con, bo qua.")
				continue
			}
			lenhCon, thanConLai := than[0], than[1:]
			if lenhCon == lenhConVersionGame {
				if err := xuLyVersionDance(conn, r, ks, thanConLai); err != nil {
					fmt.Fprintf(os.Stderr, "LOI trong version dance: %v\n", err)
				}
				return
			}
			fmt.Printf("[Lenh -28/%d] goi -28 lenh con chua duoc dien giai o day, dai %d byte: % X\n",
				lenhCon, len(thanConLai), thanConLai)
			continue
		}

		dungLai := dienGiaiGoi(cmd, than)

		// Da co phan hoi RO RANG (vd -26) cho yeu cau tao nhan vat vua gui ->
		// dienGiaiGoi da in thong bao, khong can doi them, dung lai o day.
		if dangChoPhanHoiTaoNhanVat && cmd == lenhThongBaoOK {
			return
		}

		if cmd == lenhTaoNhanVat && *coTaoNhanVat && !daGuiTaoNhanVat {
			daGuiTaoNhanVat = true
			fmt.Println()
			fmt.Printf("Thu tao nhan vat: ten=%q lang=%d kieu-toc=%d\n", *tenNhanVat, *langNhanVat, *kieuTocNhanVat)
			if err := guiGoiTaoNhanVat(conn, ks, *tenNhanVat, byte(*langNhanVat), byte(*kieuTocNhanVat)); err != nil {
				fmt.Fprintf(os.Stderr, "LOI khi gui goi tao nhan vat: %v\n", err)
				return
			}
			fmt.Printf("Da gui goi tao nhan vat. Dang cho may chu phan hoi (toi da %s)...\n", hanChoPhanHoiTaoNhanVat)
			_ = conn.SetReadDeadline(time.Now().Add(hanChoPhanHoiTaoNhanVat))
			dangChoPhanHoiTaoNhanVat = true
			continue
		}

		if dungLai {
			return
		}
	}
}

// guiGoiTaoNhanVat dong goi va gui goi -28/2 (tao nhan vat), dung dung dinh
// dang than tin ma xuLyTaoNhanVat (server, internal/game/nhanvat.go) mong
// doi: byte lenh con (lenhConTaoNhanVat), roi chuoi ten kieu UTF, roi 1 byte
// lang, roi 1 byte kieu toc.
func guiGoiTaoNhanVat(w io.Writer, ks *protocol.KeyStream, ten string, lang, kieuToc byte) error {
	dat := protocol.NewWriter()
	dat.Byte(lenhConTaoNhanVat)
	if err := dat.UTF(ten); err != nil {
		return fmt.Errorf("ten nhan vat khong hop le: %w", err)
	}
	dat.Byte(lang)
	dat.Byte(kieuToc)
	return guiGoiDaMaHoa(w, ks, lenhKhongTrenBanDo, dat.Bytes())
}

// guiGoiXinCatalog dong goi va gui mot goi -28/<lenhCon> than RONG (ngoai
// chinh byte lenh con) — dung cho ca ba loai client xin catalog (map=6,
// skill=7, item=8) LAN client bao san sang (13): ca bon deu chi can dung mot
// byte lenh con, khong co truong nao khac trong than tin client gui len.
func guiGoiXinCatalog(w io.Writer, ks *protocol.KeyStream, lenhCon byte) error {
	dat := protocol.NewWriter()
	dat.Byte(lenhCon)
	return guiGoiDaMaHoa(w, ks, lenhKhongTrenBanDo, dat.Bytes())
}

// xuLyVersionDance chay TOAN BO version dance ngay sau khi nhan -28/4
// (versionGame): boc tach 4 phien ban, gui lien tiep -87 (xin data) roi
// -28/6, -28/7, -28/8 (xin catalog map/skill/item) DUNG THU TU ma client that
// lam, in ra tung phan hoi (dac biet in SO LANG va TEN TUNG LANG trong -28/7
// — day la noi doi chieu bang mat khang dinh loi cua Ke hoach 3: 5 lang
// Naruto thay 3 lop Dragon Ball cu), roi gui -28/13 (client bao san sang).
//
// Ke hoach 3 dung o day: KHONG doc them ME_INFO/MAP INFO sau -28/13 (may chu
// hien chua gui gi them — do la pham vi Ke hoach 4), nen ham nay tra ve ngay
// sau khi gui xong -28/13, khong doc them goi nao nua.
func xuLyVersionDance(w io.Writer, r io.Reader, ks *protocol.KeyStream, thanVersionGame []byte) error {
	rd := protocol.NewReader(thanVersionGame)
	vsData, _ := rd.Byte()
	vsMap, _ := rd.Byte()
	vsSkill, _ := rd.Byte()
	vsItem, _ := rd.Byte()
	fmt.Println()
	fmt.Printf("[Lenh -28/4] versionGame: vsData=%d vsMap=%d vsSkill=%d vsItem=%d — bat dau xin cap nhat...\n",
		vsData, vsMap, vsSkill, vsItem)

	// xinVaDoc gom chung buoc gui mot goi xin catalog THUONG (-28/<lenhCon>,
	// hai byte do dai) roi doc phan hoi TUONG UNG — dung cho ca ba lenh con
	// 6/7/8 (khac -87, lenh dac biet ba byte, xu ly rieng ngay duoi day).
	xinVaDoc := func(lenhCon byte, ten string) ([]byte, error) {
		if err := guiGoiXinCatalog(w, ks, lenhCon); err != nil {
			return nil, fmt.Errorf("gui -28/%d (xin catalog %s) that bai: %w", lenhCon, ten, err)
		}
		cmd, than, err := docGoiDaMaHoa(r, ks)
		if err != nil {
			return nil, fmt.Errorf("doc phan hoi -28/%d (catalog %s) that bai: %w", lenhCon, ten, err)
		}
		if cmd != lenhKhongTrenBanDo {
			return nil, fmt.Errorf("phan hoi catalog %s: nhan lenh %d, mong doi %d (-28)", ten, cmd, lenhKhongTrenBanDo)
		}
		if len(than) == 0 || than[0] != lenhCon {
			return nil, fmt.Errorf("phan hoi catalog %s: lenh con khong khop (goi nguyen than: % X)", ten, than)
		}
		return than[1:], nil
	}

	// -87 xin data: KHAC ba cai tren o cho la lenh DAC BIET (do dai BA byte),
	// nen gui/doc rieng thay vi qua xinVaDoc — xem chu thich dau file va o
	// docGoiDaMaHoa cho ly do (goi tra ve ~87 KB vuot qua gioi han hai byte).
	fmt.Println("Xin data (-87, than rong)...")
	if err := guiGoiDaMaHoa(w, ks, lenhDuLieu, nil); err != nil {
		return fmt.Errorf("gui -87 (xin data) that bai: %w", err)
	}
	cmd, thanData, err := docGoiDaMaHoa(r, ks)
	if err != nil {
		return fmt.Errorf("doc phan hoi -87 (xin data) that bai: %w", err)
	}
	if cmd != lenhDuLieu {
		return fmt.Errorf("phan hoi xin data: nhan lenh %d, mong doi %d (-87)", cmd, lenhDuLieu)
	}
	fmt.Printf("[Lenh -87] Nhan goi data: %d byte (duoi moc 100000 client: %v)\n", len(thanData), len(thanData) < 100000)

	thanMap, err := xinVaDoc(lenhConCatalogMap, "map")
	if err != nil {
		return err
	}
	fmt.Printf("[Lenh -28/6] catalog map: %d byte: % X\n", len(thanMap), thanMap)

	thanSkill, err := xinVaDoc(lenhConCatalogSkill, "skill")
	if err != nil {
		return err
	}
	inCatalogSkill(thanSkill)

	thanItem, err := xinVaDoc(lenhConCatalogItem, "item")
	if err != nil {
		return err
	}
	fmt.Printf("[Lenh -28/8] catalog item: %d byte: % X\n", len(thanItem), thanItem)

	fmt.Println("Gui -28/13 (client bao da nhan du 4 catalog, san sang)...")
	if err := guiGoiXinCatalog(w, ks, lenhConClientOk); err != nil {
		return fmt.Errorf("gui -28/13 (client san sang) that bai: %w", err)
	}
	fmt.Println("Da gui -28/13. Ke hoach 3 dung o day (Ke hoach 4 se gui/doc them ME_INFO + MAP INFO).")
	return nil
}

// inCatalogSkill dien giai than tin catalog skill (-28/7, DA BO byte lenh
// con): vsSkill, countOption, roi countClass (SO LANG) va voi moi lang la ten
// (UTF) + countTemplate. In ro so lang va tung ten lang de nguoi dung doi
// chieu bang mat voi 5 lang Naruto (Ke hoach 3) — day la noi de phat hien
// bang mat neu may chu con phuc vu nham danh sach lop Dragon Ball cu.
func inCatalogSkill(than []byte) {
	r := protocol.NewReader(than)
	vsSkill, err := r.Byte()
	if err != nil {
		fmt.Printf("[Lenh -28/7] LOI doc vsSkill: %v\n", err)
		return
	}
	countOption, _ := r.Byte()
	countClass, err := r.Byte()
	if err != nil {
		fmt.Printf("[Lenh -28/7] LOI doc countClass: %v\n", err)
		return
	}
	fmt.Printf("[Lenh -28/7] catalog skill: vsSkill=%d countOption=%d SO LANG=%d\n", vsSkill, countOption, countClass)
	for i := 0; i < int(countClass); i++ {
		ten, err := r.UTF()
		if err != nil {
			fmt.Printf("  LOI doc ten lang thu %d: %v\n", i, err)
			return
		}
		countTemplate, _ := r.Byte()
		fmt.Printf("  Lang %d: %q (countTemplate=%d)\n", i, ten, countTemplate)
	}
}

// dienGiaiGoi in ra bang tieng Viet de nguoi khong lap trinh cung doc duoc.
// Tra ve true neu da ro ket qua dang nhap va co the dung doc tiep.
func dienGiaiGoi(cmd int8, than []byte) (dungLai bool) {
	switch cmd {
	case lenhThongBaoOK:
		rd := protocol.NewReader(than)
		noiDung, err := rd.UTF()
		if err != nil {
			fmt.Printf("[Lenh %d] May chu gui thong bao nhung khong doc duoc noi dung: %v\n", cmd, err)
			return false
		}
		fmt.Printf("[Lenh %d] THONG BAO TU MAY CHU: %s\n", cmd, noiDung)
		return false

	case lenhDangNhapThatBai:
		fmt.Printf("[Lenh %d] Ket qua: DANG NHAP THAT BAI. May chu da bao client thoat khoi man hinh \"Xin cho\".\n", cmd)
		return true

	case lenhTaoNhanVat:
		fmt.Printf("[Lenh %d] Ket qua: DANG NHAP THANH CONG! May chu chuyen sang man hinh tao nhan vat.\n", cmd)
		return true

	default:
		fmt.Printf("[Lenh %d] Goi tin chua ro y nghia trong pham vi cong cu nay, dai %d byte: % X\n", cmd, len(than), than)
		return false
	}
}

// docGoiChuaMaHoa doc mot goi CHUA MA HOA (giai doan bat tay, truoc khi dat
// khoa). May chu ghi do dai 2 byte theo BIG-ENDIAN (java.io.DataOutputStream),
// xem WriteMessage trong protocol/frame.go, nhanh "chua co khoa".
func docGoiChuaMaHoa(r io.Reader) (int8, []byte, error) {
	var dau [3]byte
	if _, err := io.ReadFull(r, dau[:]); err != nil {
		return 0, nil, fmt.Errorf("doc dau goi that bai: %w", err)
	}
	cmd := int8(dau[0])
	n := int(dau[1])<<8 | int(dau[2])
	if n > gioiHanThan {
		return 0, nil, fmt.Errorf("may chu bao do dai than tin %d, vuot gioi han an toan %d", n, gioiHanThan)
	}
	than := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(r, than); err != nil {
			return 0, nil, fmt.Errorf("doc than tin (%d byte) that bai: %w", n, err)
		}
	}
	return cmd, than, nil
}

// docGoiDaMaHoa doc mot goi tin DA MA HOA (sau khi bat tay xong). Can cu VAO
// LENH (khong can cu than tin) de quyet dinh doc hai hay BA byte do dai, dung
// dung thuat toan client that (Session_ME.cs, readMessage/readMessage2) —
// Ke hoach 3: them lenh dac biet lenhDuLieu (-87, goi data ~87 KB vuot qua
// gioi han hai byte 65535) can BA byte do dai; MOI lenh khac (ke ca moi lenh
// con cua khung -28) van dung HAI byte nhu truoc. Doc cmd TRUOC roi moi quyet
// dinh doc them may byte do dai — khac ban truoc Ke hoach 3 (doc san 3 byte
// dau roi moi xet), vi gio khong the biet truoc can hai hay ba byte do dai
// khi chua biet lenh.
func docGoiDaMaHoa(r io.Reader, ks *protocol.KeyStream) (int8, []byte, error) {
	var cmdByte [1]byte
	if _, err := io.ReadFull(r, cmdByte[:]); err != nil {
		return 0, nil, err
	}
	cmd := int8(ks.Read(cmdByte[0]))

	var n int
	if cmd == lenhDuLieu {
		// Lenh dac biet: BA byte do dai, moi byte XOR roi CONG 128 (xem chu
		// thich "Ban Java ghi XOR(x)-128, client doc readKey(b)+128" trong
		// protocol/frame.go — o day lam dung y het de giai ma dung).
		var doDai [3]byte
		if _, err := io.ReadFull(r, doDai[:]); err != nil {
			return 0, nil, err
		}
		b0 := ks.Read(doDai[0]) + 128
		b1 := ks.Read(doDai[1]) + 128
		b2 := ks.Read(doDai[2]) + 128
		n = int(b0) | int(b1)<<8 | int(b2)<<16
	} else {
		var doDai [2]byte
		if _, err := io.ReadFull(r, doDai[:]); err != nil {
			return 0, nil, err
		}
		hi := ks.Read(doDai[0])
		lo := ks.Read(doDai[1])
		n = int(hi)<<8 | int(lo)
	}
	if n > gioiHanThan {
		return 0, nil, fmt.Errorf("may chu bao do dai than tin %d, vuot gioi han an toan %d", n, gioiHanThan)
	}
	than := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(r, than); err != nil {
			return 0, nil, fmt.Errorf("doc than tin (%d byte) that bai: %w", n, err)
		}
		for i := range than {
			than[i] = ks.Read(than[i])
		}
	}
	return cmd, than, nil
}

// guiGoiDaMaHoa ghi mot goi tin DA MA HOA (lenh thuong, khong phai lenh dac
// biet do dai 3 byte) xuong may chu.
func guiGoiDaMaHoa(w io.Writer, ks *protocol.KeyStream, cmd int8, than []byte) error {
	out := make([]byte, 0, len(than)+3)
	out = append(out, ks.Write(byte(cmd)))
	n := len(than)
	out = append(out, ks.Write(byte(n>>8)))
	out = append(out, ks.Write(byte(n&0xFF)))
	for _, b := range than {
		out = append(out, ks.Write(b))
	}
	_, err := w.Write(out)
	return err
}
