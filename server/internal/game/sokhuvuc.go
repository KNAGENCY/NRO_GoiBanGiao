package game

// SoKhuVuc giu 5 khu vuc actor co dinh (mapId 1-5) — tao san luc khoi dong may
// chu, song suot doi tien trinh. Don gian hon tao/huy dong; 5 goroutine ngu la
// chi phi khong dang ke. Tao/huy dong de danh khi co map tam (to doi/su kien).
type SoKhuVuc struct {
	theoMap map[int32]*KhuVuc
}

// MoSoKhuVuc tao 5 khu vuc mapId 1-5 va KHOI DONG actor cua tung cai.
func MoSoKhuVuc() *SoKhuVuc {
	s := &SoKhuVuc{theoMap: make(map[int32]*KhuVuc)}
	for mapId := int32(1); mapId <= 5; mapId++ {
		k := MoKhuVuc(TimBanDo(mapId))
		k.BatDau()
		s.theoMap[mapId] = k
	}
	return s
}

// Tra tra khu vuc theo mapId, nil neu khong co.
func (s *SoKhuVuc) Tra(mapId int32) *KhuVuc {
	return s.theoMap[mapId]
}

// DungTatCa dung sach moi khu vuc (goi luc tat may chu).
func (s *SoKhuVuc) DungTatCa() {
	for _, k := range s.theoMap {
		k.Dung()
	}
}
