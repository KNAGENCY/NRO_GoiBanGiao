import Link from "next/link";
import { redirect } from "next/navigation";
import { docPhien } from "@/lib/phien";
import { doiMatKhau } from "./actions";

export default async function TrangDoiMatKhau({
  searchParams,
}: {
  searchParams: Promise<{ loi?: string }>;
}) {
  const phien = await docPhien();
  if (!phien) {
    redirect("/dang-nhap");
  }

  const { loi } = await searchParams;

  return (
    <main className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-6 py-16">
      <h1 className="mb-1 text-2xl font-bold text-white">Đổi mật khẩu</h1>
      <p className="mb-8 text-sm text-slate-400">
        Tài khoản: <span className="text-slate-200">{phien.ten}</span>
      </p>

      {loi && (
        <p className="mb-4 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
          {loi}
        </p>
      )}

      <form action={doiMatKhau} className="flex flex-col gap-4">
        <label className="flex flex-col gap-1.5 text-sm">
          <span className="font-medium text-slate-200">Mật khẩu hiện tại</span>
          <input
            type="password"
            name="mat_khau_cu"
            required
            className="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
          />
        </label>

        <label className="flex flex-col gap-1.5 text-sm">
          <span className="font-medium text-slate-200">Mật khẩu mới</span>
          <input
            type="password"
            name="mat_khau_moi"
            required
            minLength={6}
            placeholder="Tối thiểu 6 ký tự"
            className="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
          />
        </label>

        <label className="flex flex-col gap-1.5 text-sm">
          <span className="font-medium text-slate-200">
            Xác nhận mật khẩu mới
          </span>
          <input
            type="password"
            name="xac_nhan_mat_khau_moi"
            required
            minLength={6}
            className="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
          />
        </label>

        <button
          type="submit"
          className="mt-2 rounded-lg bg-orange-500 px-6 py-3 font-semibold text-white transition hover:bg-orange-400"
        >
          Đổi mật khẩu
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-slate-400">
        <Link
          href="/tai-khoan"
          className="font-medium text-orange-400 hover:underline"
        >
          Quay lại trang cá nhân
        </Link>
      </p>
    </main>
  );
}
