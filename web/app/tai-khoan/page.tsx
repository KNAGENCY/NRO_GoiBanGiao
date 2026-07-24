import Link from "next/link";
import { redirect } from "next/navigation";
import { sql } from "@/lib/db";
import { docPhien } from "@/lib/phien";
import { tenLang } from "@/lib/lang";
import { dangXuat } from "./actions";

export default async function TrangTaiKhoan({
  searchParams,
}: {
  searchParams: Promise<{ doi_mat_khau?: string }>;
}) {
  const phien = await docPhien();
  if (!phien) {
    redirect("/dang-nhap");
  }

  const { doi_mat_khau } = await searchParams;

  // Cot suc_manh, vang la BIGINT — postgres.js tra ve dang CHUOI (giong id
  // BIGSERIAL o W1) de khong mat do chinh xac voi so 64-bit.
  const hang = await sql<
    {
      id: string;
      ten: string;
      lang: number;
      cap: number;
      suc_manh: string;
      vang: string;
    }[]
  >`
    SELECT id, ten, lang, cap, suc_manh, vang
      FROM nhan_vat
     WHERE tai_khoan_id = ${phien.id}
     LIMIT 1
  `;
  const nhanVat = hang[0];

  return (
    <main className="mx-auto flex min-h-screen max-w-2xl flex-col gap-8 px-6 py-16">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-white">Tài khoản của tôi</h1>
        <form action={dangXuat}>
          <button
            type="submit"
            className="rounded-lg border border-slate-600 px-4 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-400 hover:bg-slate-800/50"
          >
            Đăng xuất
          </button>
        </form>
      </div>

      {doi_mat_khau && (
        <p className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-300">
          Đổi mật khẩu thành công.
        </p>
      )}

      <section className="rounded-xl border border-slate-700 bg-slate-900/50 p-6">
        <h2 className="mb-4 text-lg font-semibold text-white">
          Thông tin đăng nhập
        </h2>
        <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
          <dt className="text-slate-400">Tên đăng nhập</dt>
          <dd className="text-white">{phien.ten}</dd>
          {phien.laAdmin && (
            <>
              <dt className="text-slate-400">Vai trò</dt>
              <dd className="text-orange-300">Quản trị viên</dd>
            </>
          )}
        </dl>
        <Link
          href="/doi-mat-khau"
          className="mt-4 inline-block text-sm font-medium text-orange-400 hover:underline"
        >
          Đổi mật khẩu
        </Link>
      </section>

      <section className="rounded-xl border border-slate-700 bg-slate-900/50 p-6">
        <h2 className="mb-4 text-lg font-semibold text-white">Nhân vật</h2>
        {nhanVat ? (
          <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
            <dt className="text-slate-400">Tên</dt>
            <dd className="text-white">{nhanVat.ten}</dd>
            <dt className="text-slate-400">Làng</dt>
            <dd className="text-white">{tenLang(nhanVat.lang)}</dd>
            <dt className="text-slate-400">Cấp</dt>
            <dd className="text-white">{nhanVat.cap}</dd>
            <dt className="text-slate-400">Sức mạnh</dt>
            <dd className="text-white">
              {Number(nhanVat.suc_manh).toLocaleString("vi-VN")}
            </dd>
            <dt className="text-slate-400">Vàng</dt>
            <dd className="text-white">
              {Number(nhanVat.vang).toLocaleString("vi-VN")}
            </dd>
          </dl>
        ) : (
          <p className="text-slate-400">Chưa có nhân vật — vào game tạo.</p>
        )}
      </section>
    </main>
  );
}
