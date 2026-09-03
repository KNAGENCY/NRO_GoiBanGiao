package noidung

import (
	"reflect"
	"testing"
)

// Ghim các con số KHUNG của boss.yaml nhúng theo bảng Go cũ (noidungboss.go
// trước KH8 task 7). Lúc duyệt task 7, đổi `dame: 60 -> 61` trong YAML không
// test nào đỏ — bộ kiểm chỉ chặn miền (M7/M8/M10/M11), không ghim giá trị.
// Ai chỉnh số boss (việc của task 13) thì đổi cả đây, có chủ đích.
func TestBossNhung_GhimSoKhoiDauTheoBangGoCu(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatal(err)
	}
	type khung struct {
		Ten        string
		Cap        int
		Dame       int64
		NhipDanhMs int32
		HpChon     []int64
		GiayNghi   int32
		LaBossLon  bool
		MocNoiGian []MocNoiGianBoss
	}
	co := map[string][]khung{}
	for _, b := range du.Boss {
		for _, p := range b.Pha {
			moc := make([]MocNoiGianBoss, len(p.MocNoiGian))
			for i, m := range p.MocNoiGian {
				m.Chat = "" // thoại không ghim ở đây
				moc[i] = m
			}
			co[b.Ma] = append(co[b.Ma], khung{p.Ten, p.Cap, p.Dame, p.NhipDanhMs, p.HpChon, p.GiayNghi, p.LaBossLon, moc})
		}
	}
	muon := map[string][]khung{
		"nhan-gia-mat-na": {
			// KH8 task 13 (vòng cân bằng), có chủ đích: dame 60 → 15 và
			// hp_chon [2400,3000,3600] → [3600,4500,5400]. Lý do đầy đủ ghi cạnh
			// hai số đó trong boss.yaml; tóm tắt: 5 người chỉ số khởi đầu có
			// tổng 1000 máu và KHÔNG có đường hồi máu nào trong trận, nên dame
			// 60 hạ sạch cả năm người trong 55 giây mà boss vẫn còn 2708/3000.
			// Số mới: hạ được boss trong 2 phút 2,5 giây, còn 2/5 người sống,
			// cả hai mốc nổi giận đều kịp chạy.
			{"Kẻ Lạ Mặt", 50, 15, 1500, []int64{3600, 4500, 5400}, 1800, true, []MocNoiGianBoss{
				{PhanTramHP: 50, HeSoDame: 1.30, HeSoDef: 0.85, HeSoNhip: 0.80},
				{PhanTramHP: 20, HeSoDame: 1.60, HeSoDef: 0.70, HeSoNhip: 0.75},
			}},
			{"Mặt Nạ Vỡ", 70, 90, 1000, []int64{4000, 5000, 6000}, 5, true, []MocNoiGianBoss{
				{PhanTramHP: 50, HeSoDame: 1.35, HeSoDef: 0.80, HeSoNhip: 0.80},
				{PhanTramHP: 25, HeSoDame: 1.75, HeSoDef: 0.65, HeSoNhip: 0.75},
			}},
		},
		"hoc-vien-dao-tau": {
			{"Học Viên Đào Tẩu", 25, 35, 1200, []int64{400, 600, 800}, 600, false, []MocNoiGianBoss{
				{PhanTramHP: 40, HeSoDame: 1.50, HeSoDef: 0.80, HeSoNhip: 0.75},
			}},
		},
	}
	if !reflect.DeepEqual(co, muon) {
		t.Fatalf("boss.yaml lech bang Go cu:\n co: %+v\nmuon: %+v", co, muon)
	}
}
