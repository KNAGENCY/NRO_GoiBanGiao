// Package game — goi XUONG 44 "chat bong bong": mot nhan vat (nguoi choi hoac
// boss) noi mot cau, client ve bong bong thoai tren dau nhan vat do.
//
// ⚠️ DAC TA 2026-07-26-giao-thuc-chien-dau.md KHONG MO TA goi 44 (da tim ca
// file, khong co muc nao). Khuon byte duoi day lay tu HAI NGUON DOC LAP, khop
// nhau tuyet doi — KHONG suy doan, KHONG bia.
//
// NGUON A — may chu Java goc, 1-MAY-CHU-GAME/src/services/Service.java. CA BA
// cho gui 44 ghi than goi Y HET NHAU (chi khac nguoi NHAN goi):
//
//	Service.java:596-607  chat(Player, String)               -> sendMessAllPlayerInMap
//	Service.java:609-620  chatToAnotherNotMe(Player, String) -> sendMessAnotherNotMeInMap
//	Service.java:622-632  chatJustForMe(Player, Player, ...) -> me.sendMessage
//
//	msg = new Message(44);
//	msg.writer().writeInt((int) player.id);
//	msg.writer().writeUTF(text);
//
// NGUON B — client Unity, D:/NRO_Project/Assets/Scripts/Game1/Controller.cs
// dong 4567-4573 (ban Game2/Controller.cs giong het, cung dong 4567):
//
//	case 44:
//	    int num280   = msg.reader().readInt();
//	    string text8 = msg.reader().readUTF();
//	    (... findCharInMap(num280) ...)?.addInfo(text8);
//
// addInfo la bong bong thoai tren dau nhan vat — Game1/Char.cs:1508-1515 goi
// chatInfo.addInfo(info, 0, cInfo, isChatServer: false).
//
// KET LUAN: than goi CHI CO DUNG 2 TRUONG — int id + UTF noi dung. KHONG co
// truong "kieu chat", KHONG mau sac, KHONG thoi luong, KHONG gi khac.
package game

import "nroserver/internal/protocol"

// Ma lenh goi chat bong bong (Service.java:599/612/625 `new Message(44)`;
// client Controller.cs:4567 `case 44`).
const lenhChat int8 = 44

// soKyTuToiDaChat la so KY TU (rune) toi da mot cau chat co the mang.
//
// Day KHONG phai luat choi (gioi han spam/do dai cau la viec cua tang tren) ma
// la GIOI HAN VAT LY cua giao thuc: writeUTF chi co 2 byte do dai -> toi da
// 65535 byte sau khi ma hoa, va mot rune ton toi da 6 byte modified UTF-8
// (rune ngoai BMP bi tach thanh cap surrogate, moi nua 3 byte — xem
// protocol.Writer.UTF). 65535/6 = 10922 rune la nguong an toan cho MOI chuoi.
const soKyTuToiDaChat = 65535 / 6

// goiChat44 dung goi XUONG 44: nhan vat idNguoiNoi noi cau noiDung, client ve
// bong bong thoai tren dau nhan vat do. Khuon byte (bang chung o dau file):
//
//	int  idNguoiNoi  — 4 byte big-endian
//	UTF  noiDung     — 2 byte do dai + noi dung modified UTF-8
//
// ⚠️ TUYET DOI KHONG ghi them byte nao sau noiDung. Client doc dung 2 truong
// roi `break` — byte thua bi vut IM LANG o msg.cleanup(): KHONG loi, KHONG
// trieu chung, chay thu trong game KHONG PHAT HIEN DUOC (dung loai loi "ghi
// thua" ma dac ta §2.0 canh bao). Vi vay ham nay KHONG nhan tham so nao khac:
// khong the vo tinh truyen them mot truong "kieu chat" hay tuong tu.
//
// idNguoiNoi la int64 (khop store.NhanVat.ID) nhung ghi xuong day la int32 —
// dung nhu ban Java ep kieu `writeInt((int) player.id)`. Client cung doc bang
// readInt() roi so voi charID (int), nen ghi Long 8 byte se lam lech HET phan
// sau va client tra nham nhan vat.
//
// Ai nhan goi nay la tuy y nghia loi thoai — ca khu vuc (Service.chat), nguoi
// khac (chatToAnotherNotMe), hay rieng mot nguoi (chatJustForMe). Trong actor
// KhuVuc phai phat qua phien.XepGui() (vd k.phatChoCaKhuVuc), KHONG BAO GIO
// Gui().
func goiChat44(idNguoiNoi int64, noiDung string) *protocol.Message {
	// Cat TRUOC khi ghi mot byte nao — cung tinh than voi goiQuaiChet (cat so
	// do roi ve 127 truoc vong lap): khong bao gio de writeUTF that bai giua
	// chung roi phai bo ca goi. So byte cua chuoi Go luon >= so rune, nen chi
	// can dem rune khi chuoi da dai qua nguong.
	if len(noiDung) > soKyTuToiDaChat {
		kyTu := []rune(noiDung)
		if len(kyTu) > soKyTuToiDaChat {
			noiDung = string(kyTu[:soKyTuToiDaChat])
		}
	}

	w := protocol.NewWriter()
	w.Int(int32(idNguoiNoi)) // 1 id nguoi noi — int 4 byte, xem chu thich ham
	_ = w.UTF(noiDung)       // 2 noi dung — loi da bi chan boi buoc cat o tren
	// KHONG GHI GI THEM SAU DAY — het goi. Xem chu thich ham.
	return protocol.NewMessage(lenhChat, w)
}
