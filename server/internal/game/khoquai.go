// khoquai.go — CUA DUY NHAT cua goi `game` di vao KHO NOI DUNG cho phan QUAI.
//
// (KH8 task 8. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// hang "8" bang Task, §2 can cuoc, §4 nap nong, §5.3 do can bang, §6 anh da ky.)
//
// # CAI GI DOI O TASK NAY
//
// Truoc task nay, quai duoc sinh ra tu MOT BANG CUNG trong game/quai.go
// (taoQuaiMacDinh): hp lay tu catalog.json, con SatThuongCoBan la HANG SO 5 go
// thang vao dong code do, va danh sach diem sinh cung la cua catalog.json dung
// chung cho MOI ban do.
//
// Tu task nay, ca ba thu do nam trong kho-noi-dung/quai.yaml. Bang cung DA BI
// XOA — dot bien "sinh quai doc bang cung" khong con viet duoc nua.
//
// # KHUON GIONG HET khochieu.go (task 3), CO Y
//
// Cung mot nguonNoiDung, cung mot sync.Once cho ban NHUNG mac dinh, cung mot
// quy tac "trong actor chi dung k.acNhip". Khong dat them mot co che thu hai
// canh mot co che da co: mot nguoi doc khochieu.go xong la doc duoc file nay.
// Cu the, interface khoChieuAnhChup (khochieu.go) duoc MO RONG them ba phuong
// thuc quai, thay vi de ra mot interface thu hai — vi *noidung.anhChup thoa man
// ca hai, va mot actor chi giu MOT anh chup cho ca nhip.
//
// # ANH_ID: KIEM O `game`, KHONG KIEM O `noidung`
//
// `anh_id` trong quai.yaml phai tro toi mot id anh quai co that trong goi asset
// DA KY (assets/asset-lock.json, closure.mob_ids). Danh sach do doc boi goi
// `game` (catalog.go); goi `noidung` KHONG duoc import `game` (huong import da
// chot: game -> noidung, khong bao gio nguoc lai). Vi vay bo kiem duoc CHIA:
//
//	noidung/nap.go   kiemQuai        — ma duy nhat, hp>0, sat_thuong>0, id<=127...
//	noidung/kiemtra.go KiemIdAnhQuai — ham THUAN, nhan danh sach id lam THAM SO
//	game/khoquai.go  (file nay)      — nguoi goi ham tren, voi danh sach THAT
//
// # QUAI DANG SONG KHONG NAP NONG DUOC — DA CHOT
//
// Mot con quai lay hp/sat thuong/cap MOT LAN, luc MoKhuVuc goi taoQuaiTuAnhChup.
// Nap nong quai.yaml KHONG dung toi nhung con quai dang dung tren ban do: chung
// giu so cua THE HE LUC SINH. Chi quai SINH RA SAU do (mot khu vuc moi mo) moi
// lay so moi.
//
// Do la lua chon CO CHU DICH:
//
//   - HpToiDa cua tung con quai DA DUOC GUI XUONG CLIENT trong goi -24 MAP INFO
//     luc nguoi choi vao ban do. Doi HpToiDa giua chung ma khong gui lai goi -24
//     la lam thanh mau o client sai vinh vien, khong mot dong log.
//   - Doi hp giua tran la lam mau con quai nhay so — nguoi choi danh 3 don thay
//     mau tang.
//
// Muon so moi co hieu luc: khoi dong lai may chu. Ghi ro trong dau quai.yaml de
// nguoi sua so doc thay truoc khi hoi.
package game

import (
	"sync"

	"nroserver/internal/noidung"
)

var (
	khoQuaiKiemMotLan sync.Once
	khoQuaiLoiKiem    error
)

// kiemIdAnhQuaiMotLan doi chieu anh_id cua MOI mau quai trong anh chup voi danh
// sach id anh quai da ky trong asset-lock. Chay MOT LAN cho ca doi may chu:
// ban NHUNG khong doi giua chung, va lan nap nong (task 10) se tu goi lai qua
// duong cua no.
//
// PANIC khi lech: file YAML nam TRONG file thuc thi (go:embed) nen "anh_id tro
// vao khong khi" khong phai su co moi truong ma la mot ban build hong — im lang
// cho qua la de may chu chay voi mot con quai client khong ve duoc. Cung ky luat
// voi khoChieuMacDinh (khochieu.go).
func kiemIdAnhQuaiMotLan(ac khoChieuAnhChup) {
	khoQuaiKiemMotLan.Do(func() {
		khoQuaiLoiKiem = noidung.KiemIdAnhQuai(ac.DanhSachQuai(), catalogLangLa.MobIDs)
	})
	if khoQuaiLoiKiem != nil {
		panic("game: kho noi dung quai lech catalog anh da ky: " + khoQuaiLoiKiem.Error())
	}
}

// taoQuaiTuAnhChup sinh danh sach quai KHOI TAO cho mot ban do, doc TU MOT ANH
// CHUP da cam san. Day la duong DUY NHAT quai duoc tao ra.
//
// Ham THUAN doi voi anh chup: khong doc dong ho, khong doc bien toan cuc nao
// ngoai catalogLangLa (chi de kiem anh_id). Ban do khong co khoi `sinh:` -> nil,
// khong panic: mot ban do yen binh la mot lua chon noi dung hop le.
//
// Diem sinh tro toi mot `ma` khong co trong `mau_quai` la KHONG THE xay ra —
// noidung.kiemQuai da tu choi nap tu truoc. Neu van xay ra (ai do goi ham nay
// voi mot anh chup tu tay khong qua Nap) thi BO QUA con quai do thay vi panic:
// mot khu vuc thieu mot con quai con chay duoc, mot may chu panic thi khong.
func taoQuaiTuAnhChup(ac khoChieuAnhChup, bd *BanDo) []*Quai {
	if bd == nil {
		return nil
	}
	kiemIdAnhQuaiMotLan(ac)
	dat := ac.SinhQuaiCuaMap(bd.MapID)
	if len(dat) == 0 {
		return nil
	}
	quai := make([]*Quai, 0, len(dat))
	for _, d := range dat {
		m, co := ac.Quai(d.Ma)
		if !co {
			continue
		}
		quai = append(quai, &Quai{
			Id:             d.Id,
			MauId:          byte(m.AnhId),
			Cap:            m.Cap,
			CapTinhAnh:     m.CapTinhAnh,
			Hp:             m.Hp,
			HpToiDa:        m.Hp,
			X:              d.X,
			Y:              d.Y,
			SatThuongCoBan: m.SatThuongCoBan,
		})
	}
	return quai
}
