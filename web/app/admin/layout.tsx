import Link from "next/link";
import { redirect } from "next/navigation";
import { docPhien } from "@/lib/phien";

// Lop kiem tra thu hai (defense-in-depth) song song voi middleware.ts o goc
// du an: middleware chan o tang request truoc khi vao route, con o day doc
// lai phien tu server component de chac chan khong co nhanh nao trong
// /admin/* render duoc du server action hay cache nao do lam middleware bi
// bo qua. Ca hai lop DEU kiem o server, khong tin client.
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
