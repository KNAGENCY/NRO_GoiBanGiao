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
	if err := w.UTF(c.NPC.Dialog); err != nil {
		return nil, err
	}
	w.Byte(byte(len(c.NPC.Menu)))
	for _, menu := range c.NPC.Menu {
		if err := w.UTF(menu); err != nil {
			return nil, err
		}
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
	w.Short(c.X)
	w.Short(c.Y)
	w.Byte(c.ID)
	w.Short(c.Portrait)
}
