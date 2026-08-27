import Link from "next/link";
import { redirect } from "next/navigation";
import { docPhien } from "@/lib/phien";

// ĐÂY là chốt kiểm quyền admin thật sự của cả miền /admin/*.
//
// proxy.ts o goc du an chi loc so bo (co phien hop le hay khong) — no chay o
// Edge runtime nen khong doc duoc CSDL, ma tu ban vá "phien thu hoi duoc" thi
// co `laAdmin` trong token khong con dang tin nua. Con docPhien() o day doc
// cot `la_admin` MOI NHAT tu CSDL, nen ai vua bi ha quyen la bi day ra ngay
// lan bam chuot ke tiep, khong phai cho token het han 7 ngay.
//
// Moi server action trong /admin cung tu kiem lai bang docPhien() — de van an
// toan neu bi goi thang, bo qua giao dien. Nho React cache(), ca layout va
// action cung mot request chi ton MOT truy van CSDL.
export default async function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const phien = await docPhien();
  if (!phien || !phien.laAdmin) {
    redirect("/dang-nhap");
  }

  return (
    <div className="min-h-screen">
      <header className="border-b border-slate-800 bg-slate-900/60">
        <div className="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-4 px-6 py-4">
          <div className="flex flex-wrap items-center gap-6">
            <span className="font-bold text-white">Panel quản trị</span>
            <nav className="flex gap-4 text-sm">
              <Link href="/admin" className="text-slate-300 hover:text-white">
                Bảng điều khiển
              </Link>
              <Link
                href="/admin/tai-khoan"
                className="text-slate-300 hover:text-white"
              >
                Tài khoản
              </Link>
            </nav>
          </div>
          <div className="flex items-center gap-4 text-sm">
            <span className="text-slate-400">
              Xin chào, <span className="text-slate-200">{phien.ten}</span>
            </span>
            <Link href="/" className="text-slate-400 hover:text-white">
              Về trang chủ
            </Link>
          </div>
        </div>
      </header>
      <div className="mx-auto max-w-5xl px-6 py-10">{children}</div>
    </div>
  );
}
