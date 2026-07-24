// Chan mien /admin/* neu phien khong co laAdmin === true. Chay o tang proxy
// (truoc khi vao bat ky route handler / server component nao — day la ten
// moi cua "middleware" tu Next.js 16, xem
// https://nextjs.org/docs/messages/middleware-to-proxy) nen khong the bo qua
// bang cach goi thang mot trang admin — day la lop kiem tra server-side,
// khong tin client.
//
// Khong dung duoc lib/phien.ts nguyen ban (ham do dung next/headers
// cookies() — API danh cho route/server component, khac API cookie cua
// NextRequest/NextResponse o day) nen doc + xac thuc JWT lai o day bang
// chinh thu vien jose (tuong thich Edge runtime, dung Web Crypto API).
import { NextRequest, NextResponse } from "next/server";
import { jwtVerify } from "jose";

const TEN_COOKIE = "phien"; // PHAI khop TEN_COOKIE trong lib/phien.ts

function layKhoaBiMat(): Uint8Array {
  const bi_mat = process.env.AUTH_SECRET;
  if (!bi_mat) {
    throw new Error("Thieu bien moi truong AUTH_SECRET.");
  }
  return new TextEncoder().encode(bi_mat);
}

async function laPhienAdmin(token: string | undefined): Promise<boolean> {
  if (!token) return false;
  try {
    const { payload } = await jwtVerify(token, layKhoaBiMat());
    return payload.laAdmin === true;
  } catch {
    // Token het han, sai chu ky, hoac bi sua — coi nhu khong hop le.
    return false;
  }
}

export async function proxy(req: NextRequest) {
  const token = req.cookies.get(TEN_COOKIE)?.value;

  if (!(await laPhienAdmin(token))) {
    const url = req.nextUrl.clone();
    url.pathname = "/dang-nhap";
    url.search = "";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*"],
};
