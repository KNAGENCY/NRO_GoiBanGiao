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
//	-26  may chu hien hop thoai co nut OK (dang nhap that bai se kem theo)
//	-102 may chu bao client thoat khoi trang thai "dang cho"
//	   2 may chu bao chuyen sang man tao nhan vat (= dang nhap thanh cong)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"nroserver/internal/protocol"
)

const (
	lenhBatTay          int8 = -27
	lenhChuaDangNhap    int8 = -29
	lenhConDangNhap     byte = 0
	lenhThongBaoOK      int8 = -26
	lenhDangNhapThatBai int8 = -102
	lenhTaoNhanVat      int8 = 2

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
	flag.Parse()

	if *tenDangNhap == "" || *matKhau == "" {
		fmt.Fprintln(os.Stderr, "LOI: phai truyen ten dang nhap va mat khau.")
		fmt.Fprintln(os.Stderr, "Vi du:  thuclient -ten thutest -mat-khau 1")
		fmt.Fprintln(os.Stderr, "Xem tat ca tham so:  thuclient -h")
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

	// Buoc 4: doc va dien giai cac goi may chu tra ve, cho toi khi ro ket qua
	// dang nhap (thanh cong hoac that bai) hoac may chu dong ket noi.
	for {
		cmd, than, err := docGoiDaMaHoa(r, ks)
		if err != nil {
			if err == io.EOF {
				fmt.Println("May chu da dong ket noi.")
			} else {
				fmt.Fprintf(os.Stderr, "LOI khi doc goi tu may chu: %v\n", err)
			}
			return
		}
		if dungLai := dienGiaiGoi(cmd, than); dungLai {
			return
		}
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

// docGoiDaMaHoa doc mot goi tin DA MA HOA (sau khi bat tay xong), dung dung
// thu tu byte ma client that dung cho LENH THUONG (khong phai lenh dac biet
// do dai 3 byte — luong dang nhap trong cong cu nay khong dung lenh dac
// biet nao, xem danh sach lenhDoDaiBaByte trong protocol/frame.go).
func docGoiDaMaHoa(r io.Reader, ks *protocol.KeyStream) (int8, []byte, error) {
	var dau [3]byte
	if _, err := io.ReadFull(r, dau[:]); err != nil {
		return 0, nil, err
	}
	cmd := int8(ks.Read(dau[0]))
	hi := ks.Read(dau[1])
	lo := ks.Read(dau[2])
	n := int(hi)<<8 | int(lo)
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
