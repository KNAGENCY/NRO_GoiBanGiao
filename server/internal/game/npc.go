package game

import (
	"fmt"
	"nroserver/internal/protocol"
)

const lenhMoNPC int8 = 33

func goiMoNPC(id byte) (*protocol.Message, error) {
	c := catalogLangLa
	if id != c.NPC.ID {
		return nil, fmt.Errorf("npc %d khong co trong catalog", id)
	}
	w := protocol.NewWriter()
	w.Short(int16(c.NPC.ID))
	if err := w.UTF("Chào mừng em đến Làng Lá. Hãy vững tin trên con đường ninja của mình."); err != nil {
		return nil, err
	}
	w.Byte(1)
	if err := w.UTF("Đóng"); err != nil {
		return nil, err
	}
	w.Short(c.NPC.Portrait)
	return protocol.NewMessage(32, w), nil
}

func ghiNPCMap(w *protocol.Writer, bd *BanDo) {
	if bd == nil || bd.MapID != 1 {
		w.Byte(0)
		return
	}
	c := catalogLangLa.NPC
	w.Byte(1)
	w.Byte(0)
	w.Short(172)
	w.Short(384)
	w.Byte(c.ID)
	w.Short(c.Portrait)
}
