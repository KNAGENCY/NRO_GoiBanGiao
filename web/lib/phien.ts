// Phien dang nhap luu trong cookie ky JWT (thu vien jose), httpOnly de
// JavaScript phia trinh duyet khong doc duoc — chong danh cap phien qua XSS.
import { cookies } from "next/headers";
import { jwtVerify, SignJWT } from "jose";

const TEN_COOKIE = "phien";
const THOI_HAN = "7d"; // 7 ngay

export type DuLieuPhien = {
  id: number;
  ten: string;
  laAdmin: boolean;
};

function layKhoaBiMat(): Uint8Array {
  const bi_mat = process.env.AUTH_SECRET;
  if (!bi_mat) {
    throw new Error(
      "Thieu bien moi truong AUTH_SECRET. Tao file web/.env.local tu " +
        "web/.env.local.example roi dien mot chuoi ngau nhien dai."
    );
  }
  return new TextEncoder().encode(bi_mat);
}

/** Tao cookie phien sau khi dang nhap/dang ky thanh cong. */
export async function taoPhien(du_lieu: DuLieuPhien): Promise<void> {
  const token = await new SignJWT({ ...du_lieu })
    .setProtectedHeader({ alg: "HS256" })
    .setIssuedAt()
    .setExpirationTime(THOI_HAN)
    .sign(layKhoaBiMat());

  const kho_cookie = await cookies();
  kho_cookie.set(TEN_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: 60 * 60 * 24 * 7, // 7 ngay, tinh bang giay
  });
}

/** Doc va xac thuc phien hien tai. Tra null neu chua dang nhap hoac token hong/het han. */
export async function docPhien(): Promise<DuLieuPhien | null> {
  const kho_cookie = await cookies();
  const token = kho_cookie.get(TEN_COOKIE)?.value;
  if (!token) return null;

  try {
    const { payload } = await jwtVerify(token, layKhoaBiMat());
    if (
      typeof payload.id !== "number" ||
      typeof payload.ten !== "string" ||
      typeof payload.laAdmin !== "boolean"
    ) {
      return null;
    }
    return { id: payload.id, ten: payload.ten, laAdmin: payload.laAdmin };
  } catch {
    // Token het han, sai chu ky, hoac bi sua — coi nhu chua dang nhap.
    return null;
  }
}

/** Xoa cookie phien (dang xuat). */
export async function xoaPhien(): Promise<void> {
  const kho_cookie = await cookies();
  kho_cookie.delete(TEN_COOKIE);
}
