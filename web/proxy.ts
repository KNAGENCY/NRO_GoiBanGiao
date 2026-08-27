// Loc so bo mien /admin/*: khong co cookie phien hop le thi da thang ve trang
// dang nhap, khong cho request di xa hon. Chay o tang proxy — truoc moi route
// handler / server component (day la ten moi cua "middleware" tu Next.js 16,
// xem https://nextjs.org/docs/messages/middleware-to-proxy).
//
// ⚠️ DAY KHONG PHAI CHOT KIEM QUYEN ADMIN. Chot that nam o AdminLayout
// (app/admin/layout.tsx) va o tung server action — cho nao goi docPhien(),
// tuc cho nao doc duoc CSDL.
//
// VI SAO BO PHAN KIEM `laAdmin` O DAY (truoc kia co):
// Tang proxy chay o Edge runtime, khong mo duoc ket noi TCP toi PostgreSQL,
// nen thu duy nhat no doc duoc la noi dung token. Ma tu ban vá "phien thu hoi
// duoc", noi dung token KHONG con la su that: mot nguoi vua bi ha quyen admin
// van mang token ghi `laAdmin: true` toi 7 ngay. Giu lai phep kiem do thi:
//   - no khong chan them duoc ai ca — no chi sai ve phia CHO QUA (nguoi da bi
//     ha quyen van lot, va nguoi vua duoc nang quyen thi bi chan oan);
//   - va no tao ra nguon su that thu hai, ben canh CSDL. Hai nguon su that ve
//     cung mot cau hoi la cach kinh dien de mot chot bao mat muc ruong: sau
//     nay sua CSDL ma quen cho nay.
// Nen o day chi con dung mot cau hoi ma tang nay tra loi duoc chac chan:
// "co phai mot phien do CHINH may chu nay ky khong?" — chan duoc khach vang
// lai va token bia dat, re va khong noi doi.
import { NextRequest, NextResponse } from "next/server";
import { jwtVerify } from "jose";

const TEN_COOKIE = "phien"; // PHAI khop TEN_COOKIE trong lib/xacthucphien.ts

function layKhoaBiMat(): Uint8Array {
  const bi_mat = process.env.AUTH_SECRET;
  if (!bi_mat) {
    throw new Error("Thieu bien moi truong AUTH_SECRET.");
  }
  return new TextEncoder().encode(bi_mat);
}

// Dung chinh thu vien jose (tuong thich Edge runtime, dung Web Crypto API) —
// khong dung duoc lib/xacthucphien.ts vi ham do co truy van CSDL.
async function coPhienHopLe(token: string | undefined): Promise<boolean> {
  if (!token) return false;
  try {
    await jwtVerify(token, layKhoaBiMat());
    return true;
  } catch {
    // Token het han, sai chu ky, hoac bi sua — coi nhu khong hop le.
    return false;
  }
}

export async function proxy(req: NextRequest) {
  const token = req.cookies.get(TEN_COOKIE)?.value;

  if (!(await coPhienHopLe(token))) {
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
