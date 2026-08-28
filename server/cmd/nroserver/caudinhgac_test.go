package main

import "testing"

// Ba con so cua hai cai gac theo dia chi mang (chan do mat khau, gioi han ket
// noi) doc tu bien moi truong de nguoi van hanh noi long duoc NGAY trong ngay
// mo cua neu nguoi choi cua mot nha mang bi va lay. Chinh vi vay cach doc
// chung phai chac: mot gia tri go nham KHONG duoc lang le bien thanh "khong
// gioi han gi ca" — do la kieu hong khong co trieu chung nao cho toi dung hom
// bi tan cong.
func TestDocSoNguyenMoiTruong(t *testing.T) {
	const macDinh = 10

	truongHop := []struct {
		ten       string
		datGiaTri bool
		giaTri    string
		mongDoi   int
	}{
		{"chua dat thi lay mac dinh", false, "", macDinh},
		{"chuoi rong thi lay mac dinh", true, "", macDinh},
		{"chi co dau cach thi lay mac dinh", true, "   ", macDinh},
		{"so hop le", true, "25", 25},
		{"so hop le co dau cach thua", true, " 25 ", 25},
		{"so 0 KHONG duoc coi la khong gioi han", true, "0", macDinh},
		{"so am KHONG duoc coi la khong gioi han", true, "-5", macDinh},
		{"chu thi lay mac dinh", true, "nhieu", macDinh},
		{"so thap phan thi lay mac dinh", true, "12.5", macDinh},
	}

	for _, tc := range truongHop {
		t.Run(tc.ten, func(t *testing.T) {
			const bien = "NRO_THU_DOC_SO"
			if tc.datGiaTri {
				t.Setenv(bien, tc.giaTri)
			}
			if duoc := docSoNguyenMoiTruong(bien, macDinh); duoc != tc.mongDoi {
				t.Fatalf("doc %q duoc %d, mong doi %d", tc.giaTri, duoc, tc.mongDoi)
			}
		})
	}
}
