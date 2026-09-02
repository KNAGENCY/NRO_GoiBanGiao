// Bam va doi chieu mat khau bang PBKDF2-HMAC-SHA256.
//
// Dinh dang chuoi luu PHAI khop tuyet doi voi may chu game (Go, xem
// server/internal/auth/password.go) de web va game dung chung mot mat khau:
//
//   pbkdf2$<so_vong>$<muoi base64>$<bam base64>
//
// - HMAC-SHA256, 100.000 vong, muoi ngau nhien 16 byte, khoa dan xuat 32 byte.
// - Base64 CHUAN (khong phai base64url) — giong Go base64.StdEncoding.
//
// Sua bam ma quen doi chieu (hoac nguoc lai) se khien nguoi choi dang ky tren
// web khong dang nhap duoc vao game.
import { pbkdf2Sync, randomBytes, timingSafeEqual } from "node:crypto";

const TIEN_TO = "pbkdf2$";
const SO_VONG = 100_000;
const SO_VONG_TOI_DA = 1_000_000; // chan tren, phong chuoi bam hong lam treo may chu
const SO_BYTE_MUOI = 16;
const SO_BYTE_KHOA = 32; // 256 bit

// BAM_GIA la mot chuoi bam THAT nhung cua mot mat khau ngau nhien 32 byte ma
// khong ai tren doi biet — sinh MOT LAN bang chinh bamMatKhau() roi dan cung
// vao day. Y het hang so `bamGia` ben may chu Go
// (server/internal/auth/password.go) — doc chu thich o do de biet day du VI
// SAO can no.
//
// Tom tat: truoc khi co hang so nay, nhanh "khong tim thay tai khoan" (web/
// app/dang-nhap/actions.ts) tra loi NGAY, con nhanh "sai mat khau" phai chay
// PBKDF2 100.000 vong (~13ms). Ke xau bam gio dang nhap la doc duoc ten dang
// nhap nao co that. Cach bit: nhanh "khong tim thay" cung phai chay kiemMatKhau
// voi mot chuoi bam gia (khong khop mat khau nao) de ton dung thoi gian nhu
// nhanh that, roi moi tra loi loi chung.
//
// VI SAO LA HANG SO DAN CUNG chu khong bam luc chay: bam luc khoi dong ton
// them 13ms va phai xu ly loi o cho khong dang co loi; bam luc CHAY thi con
// te hon — moi lan dang nhap sai lai bam them mot lan, gap doi chi phi. Hang
// so thi khong ton gi va luon giong nhau giua cac lan chay.
//
// BAY: neu sau nay doi SO_VONG ma quen sinh lai chuoi nay, chuoi gia se re
// hon chuoi that va lo hong quay lai TRONG IM LANG — xem test "BAM_GIA phai
// dung SO_VONG voi ham bam that" trong matkhau.test.ts.
// Xuat ra (thay vi giu private) de test khoa duoc dung SO_VONG voi ham bam
// that — xem TestBamGiaDungSoVongVoiBamThat ben Go, y het o matkhau.test.ts.
export const BAM_GIA =
  "pbkdf2$100000$ymOxsTrtw2kd1K8yPL4qyA==$rRXayQIM3KcEQgFBA8I7sKOfjsBb4Sl+WU8YXQ1wJWQ=";

/** IsHashed: chuoi luu co phai dang bam pbkdf2$... hay khong. */
export function laDaBam(luuTru: string): boolean {
  return luuTru.startsWith(TIEN_TO);
}

/** Bam mat khau moi voi muoi ngau nhien. */
export function bamMatKhau(matKhau: string): string {
  const muoi = randomBytes(SO_BYTE_MUOI);
  const bam = pbkdf2Sync(matKhau, muoi, SO_VONG, SO_BYTE_KHOA, "sha256");
  return `${TIEN_TO}${SO_VONG}$${muoi.toString("base64")}$${bam.toString(
    "base64"
  )}`;
}

/**
 * Doi chieu mat khau nguoi dung nhap voi chuoi luu trong CSDL. Doc so vong
 * lap tu chinh chuoi luu (khong dung hang so cung) de tuong thich nguoc neu
 * sau nay tang so vong; so sanh bam bang thuat toan hang thoi gian de khong
 * lo thong tin qua do tre xu ly.
 */
export function kiemMatKhau(matKhau: string, luuTru: string): boolean {
  if (!luuTru || !laDaBam(luuTru)) return false;

  const phan = luuTru.split("$");
  if (phan.length !== 4) return false;

  const soVong = Number(phan[1]);
  if (!Number.isInteger(soVong) || soVong < 1 || soVong > SO_VONG_TOI_DA) {
    return false;
  }

  let muoi: Buffer;
  let muon: Buffer;
  try {
    muoi = Buffer.from(phan[2], "base64");
    muon = Buffer.from(phan[3], "base64");
  } catch {
    return false;
  }
  if (muon.length === 0) return false;

  const duoc = pbkdf2Sync(matKhau, muoi, soVong, muon.length, "sha256");
  if (duoc.length !== muon.length) return false;

  return timingSafeEqual(duoc, muon);
}

/**
 * CanBangThoiGian: ton dung bang khoi luong tinh toan cua mot lan kiemMatKhau
 * that (doi chieu voi BAM_GIA), roi LUON tra ve false.
 *
 * Goi no o nhanh "khong tim thay tai khoan" (khi khong co chuoi bam that nao
 * de doi chieu) de nhanh do ton thoi gian xap xi nhanh "sai mat khau" — xem
 * chu thich cua BAM_GIA o tren.
 *
 * Luon tra ve (thay vi khong tra gi) de noi goi dung duoc no trong mot bieu
 * thuc, tranh bi trinh bien dich / linter coi la loi goi thua ma xoa mat.
 */
export function canBangThoiGian(matKhau: string): boolean {
  kiemMatKhau(matKhau, BAM_GIA);
  return false;
}
