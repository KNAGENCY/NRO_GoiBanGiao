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
