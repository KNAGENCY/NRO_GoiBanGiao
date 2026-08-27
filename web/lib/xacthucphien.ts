// Loi cua phien dang nhap: ky token, va XAC THUC token dua tren CSDL.
//
// File nay CO Y khong dung `next/headers` (khong dung cookie) — no chi lam
// viec voi chuoi token. Nho vay:
//   - kiem thu duoc bang tsx thuan (xem xacthucphien.test.ts), khong can dung
//     len ca mot request Next.js;
//   - phan doc/ghi cookie nam gon trong lib/phien.ts, mong va de soi.
//
// =========================================================================
// VI SAO PHAI HOI CSDL — va cai gia phai tra
// =========================================================================
// Cookie phien la JWT tu chua, song 7 ngay. Neu chi kiem chu ky + han dung
// (cach lam cu) thi noi dung token la "su that duy nhat" trong suot 7 ngay,
// keo theo ba hau qua da xac minh trong code, khong phai nghi ngo:
//   1. Khoa tai khoan khong co tac dung — `bi_khoa` chi duoc doc luc dang
//      nhap, nguoi dang co phien van dung web binh thuong toi 7 ngay.
//   2. Ha quyen admin khong co tac dung — `laAdmin` nuong cung trong token.
//   3. Dang xuat chi xoa cookie o mot trinh duyet — ban sao token van song.
// Ca ba deu la "quyet dinh cua nguoi quan tri khong an vao he thong". Muon
// chung an ngay thi bat buoc phai doc trang thai MOI NHAT tu CSDL.
//
// >>> HUONG DA CHON: hoi CSDL o MOI lan doc phien. <<<
//
// Cai gia, noi thang: **them dung 1 truy van cho moi request CO cookie phien**
// (tra cuu theo khoa chinh `tai_khoan.id`, tra ve 4 cot — la index lookup,
// duoi 1 ms tren CSDL cung may). Khach vang lai chua dang nhap KHONG ton them
// truy van nao ca (khong co cookie thi thoat som, chua cham CSDL).
// Con "1 truy van cho moi request" chu khong phai "moi lan goi docPhien()" la
// nho docPhien() duoc boc trong React `cache()` — mot request dung ca layout +
// page + server action van chi hoi CSDL mot lan (xem lib/phien.ts).
//
// -------------------------------------------------------------------------
// NHUNG HUONG DA CAN NHAC VA LOAI — ghi lai de phien sau khoi nghi lai tu dau
// -------------------------------------------------------------------------
// (a) CHI KIEM O ROUTE NHAY CAM (vd chi /admin/*), route thuong thi tin token.
//     Loai. Ba ly do: mot la nguoi bi khoa van doi duoc mat khau, van xem duoc
//     trang ca nhan — tuc lo hong 1 chi vá duoc mot nua; hai la "the nao la
//     nhay cam" se muc theo thoi gian, them route moi ma quen them lop kiem la
//     lai thung, va kieu thung nay im lang; ba la mac dinh cua no la MO — sai
//     huong cho mot chot bao mat.
//
// (b) NHO DEM NGAN HAN (vd nho ket qua 30 giay trong bo nho tien trinh).
//     Loai. No bien "khoa ngay lap tuc" thanh "khoa sau vai chuc giay" — dung
//     luc can nhat (dang xu ly ke pha game) thi do la khoang tre te nhat. Va
//     no khong dang tin: chay nhieu tien trinh / nhieu may thi moi ban co bo
//     nho rieng, khoang tre thanh khong doan duoc. Doi lai duoc rat it, vi thu
//     ta tiet kiem la mot index lookup tren CSDL noi bo. Neu sau nay do duoc
//     that su nghen o day thi hay nho dem — con bay gio la toi uu mu.
//
// (c) SO THE HE `the_he_phien` DUNG MOT MINH, khong doc CSDL moi request.
//     Loai vi... khong the: muon biet so the he hien tai la bao nhieu thi VAN
//     phai doc CSDL. `the_he_phien` khong thay the duoc viec hoi CSDL, no la
//     thu THEM VAO de vá duoc lo hong 3 (huy phien tu xa). Nen o day dung ca
//     hai: hoi CSDL moi lan, va trong lan hoi do doi chieu luon so the he.
//
// (d) BANG `phien` RIENG (moi phien dang nhap mot dong, dang xuat thi xoa).
//     Loai. Chinh xac hon that (huy duoc tung thiet bi), nhung phai don rac
//     dinh ky va them mot bang la vao CSDL dung chung voi may chu game Go.
//     Chua dang cong suc. Cai gia cua lua chon nay: dang xuat o dien thoai se
//     dang xuat luon o may tinh. Doi lai duoc mien phi "thoat khoi moi thiet bi".
//
// -------------------------------------------------------------------------
// KHI CSDL CHET THI SAO: xacThucToken CO Y de loi truy van NEM RA NGOAI, khong
// nuot thanh `null`. Nem ra -> trang bao loi, nhin thay ngay. Nuot thanh null
// -> moi nguoi bong dung "bi dang xuat" ma khong ai hieu tai sao. Ca hai deu
// la fail-closed (khong ai duoc coi la da dang nhap), nhung chi mot cai noi
// that. Truy van CO Y dat NGOAI khoi try/catch ben duoi — dung keo vao trong.
import { jwtVerify, SignJWT } from "jose";
import { sql } from "./db";

/** Ten cookie phien. lib/phien.ts va proxy.ts deu phai dung dung ten nay. */
export const TEN_COOKIE = "phien";

/** Han dung cua token (dinh dang cua thu vien jose). */
export const THOI_HAN = "7d";

/** Cung mot con so nhung tinh bang giay, cho thuoc tinh maxAge cua cookie. */
export const THOI_HAN_GIAY = 60 * 60 * 24 * 7;

/** Nhung gi phan con lai cua web duoc biet ve nguoi dang dang nhap. */
export type DuLieuPhien = {
  id: number;
  ten: string;
  laAdmin: boolean;
};

/**
 * Noi dung ben trong token. Khac `DuLieuPhien` o cho co them `theHe` — so the
 * he phien tai thoi diem phat hanh, dung de doi chieu voi cot `the_he_phien`.
 */
export type NoiDungToken = DuLieuPhien & { theHe: number };

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

/** Ky mot token phien. Ben goi phai tu lay `theHe` tu CSDL (xem taoPhien). */
export async function kyToken(noiDung: NoiDungToken): Promise<string> {
  return new SignJWT({ ...noiDung })
    .setProtectedHeader({ alg: "HS256" })
    .setIssuedAt()
    .setExpirationTime(THOI_HAN)
    .sign(layKhoaBiMat());
}

/**
 * Doc va xac thuc mot token. Tra `null` neu khong duoc phep coi la dang nhap.
 *
 * Bon chot, theo thu tu tu re den dat:
 *   1. Khong co token          -> null (khong cham CSDL)
 *   2. Chu ky sai / het han / thieu truong -> null (khong cham CSDL)
 *   3. Doc CSDL 1 lan:
 *      - khong con tai khoan   -> null
 *      - `bi_khoa`             -> null   (vá lo hong 1)
 *      - so the he lech        -> null   (vá lo hong 3)
 *   4. Tra du lieu LAY TU CSDL, khong lay tu token (vá lo hong 2: `laAdmin`
 *      va `ten` luon la ban moi nhat).
 */
export async function xacThucToken(
  token: string | undefined | null
): Promise<DuLieuPhien | null> {
  if (!token) return null;

  let noiDung: NoiDungToken;
  try {
    const { payload } = await jwtVerify(token, layKhoaBiMat());
    if (
      typeof payload.id !== "number" ||
      typeof payload.ten !== "string" ||
      typeof payload.laAdmin !== "boolean" ||
      // Token phat hanh TRUOC ban vá nay khong co `theHe` -> khong doi chieu
      // duoc so the he -> khong thu hoi duoc -> tu choi. Hau qua: ai dang giu
      // cookie cu phai dang nhap lai dung mot lan. Chap nhan, vi giu chung lai
      // la giu nguyen lo hong.
      typeof payload.theHe !== "number"
    ) {
      return null;
    }
    noiDung = {
      id: payload.id,
      ten: payload.ten,
      laAdmin: payload.laAdmin,
      theHe: payload.theHe,
    };
  } catch {
    // Token het han, sai chu ky, hoac bi sua — coi nhu chua dang nhap.
    return null;
  }

  // CO Y nam NGOAI try/catch o tren: CSDL chet thi phai bao loi ro rang, dung
  // im lang bien thanh "chua dang nhap" (xem ghi chu dau file).
  const hang = await sql<
    {
      ten_dang_nhap: string;
      bi_khoa: boolean;
      la_admin: boolean;
      the_he_phien: number;
    }[]
  >`
    SELECT ten_dang_nhap, bi_khoa, la_admin, the_he_phien
      FROM tai_khoan
     WHERE id = ${noiDung.id}
     LIMIT 1
  `;

  const taiKhoan = hang[0];
  if (!taiKhoan) return null; // tai khoan da bi xoa han
  if (taiKhoan.bi_khoa) return null; // bi khoa giua chung -> ra ngay
  if (taiKhoan.the_he_phien !== noiDung.theHe) return null; // da dang xuat / doi mat khau

  return {
    id: noiDung.id,
    ten: taiKhoan.ten_dang_nhap,
    laAdmin: taiKhoan.la_admin,
  };
}

/**
 * Lay id tai khoan tu token ma KHONG kiem CSDL.
 *
 * Chi dung cho dang xuat: luc do tai khoan co the dang bi khoa hoac so the he
 * da lech, `xacThucToken` se tra null — nhung ta van can biet id de con cong
 * so the he len. Chu ky VAN duoc kiem, nen khong ai bia id nguoi khac duoc.
 */
export async function docIdTuToken(
  token: string | undefined | null
): Promise<number | null> {
  if (!token) return null;
  try {
    const { payload } = await jwtVerify(token, layKhoaBiMat());
    return typeof payload.id === "number" ? payload.id : null;
  } catch {
    return null;
  }
}

/**
 * Huy TAT CA phien cua mot tai khoan bang cach cong 1 vao `the_he_phien`.
 * Moi token da phat hanh truoc do lech so -> chet ngay lap tuc, o moi thiet bi.
 *
 * Goi khi: dang xuat, doi mat khau, hoac nghi token bi lo.
 * Tra ve so the he MOI (de ben goi ky lai token cho chinh trinh duyet nay ma
 * khong phai truy van them lan nua), hoac null neu khong co tai khoan do.
 */
export async function huyToanBoPhien(id: number): Promise<number | null> {
  const hang = await sql<{ the_he_phien: number }[]>`
    UPDATE tai_khoan
       SET the_he_phien = the_he_phien + 1
     WHERE id = ${id}
    RETURNING the_he_phien
  `;
  return hang[0]?.the_he_phien ?? null;
}
