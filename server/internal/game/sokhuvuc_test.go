package game

import "testing"

func TestMoSoKhuVuc(t *testing.T) {
	s := MoSoKhuVuc()
	if s == nil {
		t.Fatal("MoSoKhuVuc should not return nil")
	}

	// Test: dung co 5 khu vuc (mapId 1-5)
	for mapId := int32(1); mapId <= 5; mapId++ {
		kv := s.Tra(mapId)
		if kv == nil {
			t.Errorf("Tra(%d) should not be nil", mapId)
		}
		if kv.BanDo() == nil {
			t.Errorf("Tra(%d).BanDo() should not be nil", mapId)
		}
		if kv.BanDo().MapID != mapId {
			t.Errorf("Tra(%d).BanDo().MapID = %d, want %d", mapId, kv.BanDo().MapID, mapId)
		}
	}

	// Test: Tra(99) return nil
	if kv := s.Tra(99); kv != nil {
		t.Errorf("Tra(99) should be nil, got %v", kv)
	}

	// Test: Tra(0) return nil
	if kv := s.Tra(0); kv != nil {
		t.Errorf("Tra(0) should be nil, got %v", kv)
	}
}

func TestDungTatCa(t *testing.T) {
	s := MoSoKhuVuc()

	// DungTatCa should not hang. Goi xong va kiem thu ket thuc binh thuong.
	s.DungTatCa()

	// Test ket thuc khong treo (neu treo thi test se timeout).
	// Neu toi day dung tuc kiem test PASSED.
}
