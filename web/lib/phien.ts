// Phien dang nhap luu trong cookie ky JWT (thu vien jose), httpOnly de
// JavaScript phia trinh duyet khong doc duoc — chong danh cap phien qua XSS.
//
// File nay chi lo phan COOKIE. Toan bo phan quyet dinh "token nay con hieu
// luc khong" nam o lib/xacthucphien.ts — doc chu thich dau file do de biet vi
// sao moi lan doc phien deu hoi CSDL, va nhung huong nao da bi loai.
import { cookies } from "next/headers";
import { cache } from "react";
import {
  TEN_COOKIE,
  THOI_HAN_GIAY,
  docIdTuToken,
  huyToanBoPhien,
  kyToken,
  xacThucToken,
  type DuLieuPhien,
  type NoiDungToken,
} from "./xacthucphien";

export type { DuLieuPhien, NoiDungToken };
export { huyToanBoPhien };

/**
 * Tao cookie phien sau khi dang nhap/dang ky thanh cong.
 *
 * Ben goi phai truyen `theHe` — lay tu cot `the_he_phien` cua chinh tai khoan
 * do (SELECT luc dang nhap, hoac RETURNING luc dang ky). Bat buoc truyen chu
 * khong tu doc o day de khong am tham them mot truy van vao duong dang nhap,
 * va de TypeScript chan viec quen mat no.
 */
export async function taoPhien(noiDung: NoiDungToken): Promise<void> {
  const token = await kyToken(noiDung);

  const kho_cookie = await cookies();
  kho_cookie.set(TEN_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: THOI_HAN_GIAY,
  });
}

/**
 * Doc va xac thuc phien hien tai. Tra null neu chua dang nhap, token hong/het
 * han, tai khoan bi khoa, hoac phien da bi thu hoi.
 *
 * Boc trong React `cache()`: mot request cua Next.js thuong goi docPhien()
 * nhieu lan (layout goi mot lan, page goi mot lan, server action goi mot lan
 * nua). `cache()` nho ket qua trong PHAM VI MOT REQUEST, nen chi phi thuc te
 * la **1 truy van CSDL cho moi request**, khong phai moi lan goi. Het request
 * la bo nho do bien mat — khong co chuyen du lieu cu song sot sang request
 * sau, tuc khoa tai khoan van an vao ngay lap tuc.
 */
export const docPhien = cache(async function docPhien(): Promise<DuLieuPhien | null> {
  const kho_cookie = await cookies();
  return xacThucToken(kho_cookie.get(TEN_COOKIE)?.value);
});

/**
 * Dang xuat. Lam HAI viec, va viec thu hai moi la viec quan trong:
 *  1. Xoa cookie o trinh duyet nay.
 *  2. Cong 1 vao `the_he_phien` -> moi ban sao cua token (o thiet bi khac, hay
 *     trong tay ke da cop duoc chuoi cookie) chet ngay lap tuc.
 *
 * Danh doi da biet truoc: dang xuat o dien thoai se dang xuat luon o may tinh.
 * Doi lai, "dang xuat" moi thuc su co nghia — xem lib/xacthucphien.ts muc (d).
 */
export async function xoaPhien(): Promise<void> {
  const kho_cookie = await cookies();
  const token = kho_cookie.get(TEN_COOKIE)?.value;
  kho_cookie.delete(TEN_COOKIE);

  // docIdTuToken chu khong phai docPhien(): tai khoan co the dang bi khoa,
  // luc do docPhien() tra null nhung ta VAN can huy phien cua ho.
  const id = await docIdTuToken(token);
  if (id !== null) {
    await huyToanBoPhien(id);
  }
}
