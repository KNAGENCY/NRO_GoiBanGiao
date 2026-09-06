// Giu phim NHAY (GameScr.cs:2708 doc keyHold[Main.isPC ? 21 : 2]) - tha bang tha-phim.cs. Nhay xong moi giu trai/phai
// (chuoi else-if: dang giu nhay thi phim trai/phai KHONG duoc xu ly).
Game1.GameCanvas.keyHold[21] = true;
Game1.GameCanvas.keyHold[2] = true;
return "giu phim nhay; cx=" + Game1.Char.myCharz().cx + " cy=" + Game1.Char.myCharz().cy;
