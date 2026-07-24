import Link from "next/link";
import { dangKy } from "./actions";

export default async function TrangDangKy({
  searchParams,
}: {
  searchParams: Promise<{ loi?: string; ten?: string }>;
}) {
  const { loi, ten } = await searchParams;

  return (
    <main className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-6 py-16">
      <h1 className="mb-1 text-2xl font-bold text-white">Đăng ký tài khoản</h1>
      <p className="mb-8 text-sm text-slate-400">
        Tài khoản dùng chung cho web và game — đăng ký ở đây rồi vào game chơi luôn.
      </p>

      {loi && (
        <p className="mb-4 rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
          {loi}
        </p>
      )}

      <form action={dangKy} className="flex flex-col gap-4">
        <label className="flex flex-col gap-1.5 text-sm">
          <span className="font-medium text-slate-200">Tên đăng nhập</span>
          <input
            type="text"
            name="ten_dang_nhap"
            required
            minLength={4}
            maxLength={20}
            pattern="[A-Za-z0-9_]{4,20}"
            defaultValue={ten}
            placeholder="4-20 ký tự: chữ, số, gạch dưới"
            className="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
          />
        </label>

        <label className="flex flex-col gap-1.5 text-sm">
          <span className="font-medium text-slate-200">Mật khẩu</span>
          <input
            type="password"
            name="mat_khau"
            required
            minLength={6}
            placeholder="Tối thiểu 6 ký tự"
            className="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
          />
        </label>

        <label className="flex flex-col gap-1.5 text-sm">
          <span className="font-medium text-slate-200">Xác nhận mật khẩu</span>
          <input
            type="password"
            name="xac_nhan_mat_khau"
            required
            minLength={6}
            className="rounded-lg border border-slate-700 bg-slate-900 px-4 py-2.5 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
          />
        </label>

        <button
          type="submit"
          className="mt-2 rounded-lg bg-orange-500 px-6 py-3 font-semibold text-white transition hover:bg-orange-400"
        >
          Đăng ký
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-slate-400">
        Đã có tài khoản?{" "}
        <Link href="/dang-nhap" className="font-medium text-orange-400 hover:underline">
          Đăng nhập
        </Link>
      </p>
    </main>
  );
}
