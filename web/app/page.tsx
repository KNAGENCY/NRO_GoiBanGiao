import Link from "next/link";
import { docPhien } from "@/lib/phien";

export default async function TrangChu() {
  const phien = await docPhien();

  return (
    <main className="mx-auto flex min-h-screen max-w-3xl flex-col items-center justify-center gap-8 px-6 py-16 text-center">
      <div className="flex flex-col items-center gap-3">
        <span className="rounded-full border border-orange-500/30 bg-orange-500/10 px-4 py-1 text-sm font-medium text-orange-300">
          Naruto Online
        </span>
        <h1 className="text-4xl font-bold tracking-tight text-white sm:text-5xl">
          NRO — Con Đường Ninja
        </h1>
        <p className="max-w-xl text-balance text-base text-slate-300 sm:text-lg">
          Nhập vai ninja, luyện chakra, chiêu mộ đồng đội và chinh phục thế
          giới nhẫn giả. Tham gia ngay để bắt đầu hành trình trở thành Hokage.
        </p>
      </div>

      {phien ? (
        <div className="flex flex-col items-center gap-3">
          <p className="text-slate-300">
            Xin chào, <span className="font-semibold text-white">{phien.ten}</span>
          </p>
          <div className="flex flex-col gap-3 sm:flex-row">
            <Link
              href="/tai-khoan"
              className="rounded-lg bg-orange-500 px-6 py-3 font-semibold text-white shadow-lg shadow-orange-500/20 transition hover:bg-orange-400"
            >
              Vào trang cá nhân
            </Link>
            {phien.laAdmin && (
              <Link
                href="/admin"
                className="rounded-lg border border-slate-600 px-6 py-3 font-semibold text-slate-100 transition hover:border-slate-400 hover:bg-slate-800/50"
              >
                Panel quản trị
              </Link>
            )}
          </div>
        </div>
      ) : (
        <div className="flex flex-col gap-3 sm:flex-row">
          <Link
            href="/dang-ky"
            className="rounded-lg bg-orange-500 px-6 py-3 font-semibold text-white shadow-lg shadow-orange-500/20 transition hover:bg-orange-400"
          >
            Đăng ký
          </Link>
          <Link
            href="/dang-nhap"
            className="rounded-lg border border-slate-600 px-6 py-3 font-semibold text-slate-100 transition hover:border-slate-400 hover:bg-slate-800/50"
          >
            Đăng nhập
          </Link>
        </div>
      )}

      <Link
        href="/bang-xep-hang"
        className="text-sm font-medium text-slate-400 underline-offset-4 hover:text-orange-300 hover:underline"
      >
        Xem bảng xếp hạng
      </Link>
    </main>
  );
}
