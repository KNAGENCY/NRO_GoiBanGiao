import { sql } from "@/lib/db";
import { docPhien } from "@/lib/phien";
import { doiKhoaTaiKhoan } from "./actions";

export default async function TrangAdminTaiKhoan({
  searchParams,
}: {
  searchParams: Promise<{ q?: string }>;
}) {
  const phien = await docPhien();
  const { q } = await searchParams;
  const tuKhoa = (q ?? "").trim();

  // Top 100 (khong phan trang o W2), loc theo ten_dang_nhap neu co tu khoa.
  // Truy van tham so hoa qua postgres.js template — khong noi chuoi.
  const hang = tuKhoa
    ? await sql<
        {
          id: string;
          ten_dang_nhap: string;
          bi_khoa: boolean;
          la_admin: boolean;
          tao_luc: Date;
        }[]
      >`
        SELECT id, ten_dang_nhap, bi_khoa, la_admin, tao_luc
          FROM tai_khoan
         WHERE ten_dang_nhap ILIKE ${"%" + tuKhoa + "%"}
         ORDER BY tao_luc DESC
         LIMIT 100
      `
    : await sql<
        {
          id: string;
          ten_dang_nhap: string;
          bi_khoa: boolean;
          la_admin: boolean;
          tao_luc: Date;
        }[]
      >`
        SELECT id, ten_dang_nhap, bi_khoa, la_admin, tao_luc
          FROM tai_khoan
         ORDER BY tao_luc DESC
         LIMIT 100
      `;

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold text-white">Danh sách tài khoản</h1>

      <form method="GET" className="flex gap-2">
        <input
          type="text"
          name="q"
          defaultValue={tuKhoa}
          placeholder="Tìm theo tên đăng nhập..."
          className="w-full max-w-sm rounded-lg border border-slate-700 bg-slate-900 px-4 py-2 text-white outline-none placeholder:text-slate-500 focus:border-orange-500"
        />
        <button
          type="submit"
          className="rounded-lg bg-orange-500 px-4 py-2 font-medium text-white transition hover:bg-orange-400"
        >
          Tìm
        </button>
      </form>

      <div className="overflow-x-auto rounded-xl border border-slate-700">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-800/60 text-slate-300">
            <tr>
              <th className="px-4 py-3">ID</th>
              <th className="px-4 py-3">Tên đăng nhập</th>
              <th className="px-4 py-3">Vai trò</th>
              <th className="px-4 py-3">Trạng thái</th>
              <th className="px-4 py-3">Ngày tạo</th>
              <th className="px-4 py-3">Hành động</th>
            </tr>
          </thead>
          <tbody>
            {hang.length === 0 && (
              <tr>
                <td colSpan={6} className="px-4 py-6 text-center text-slate-400">
                  Không tìm thấy tài khoản nào.
                </td>
              </tr>
            )}
            {hang.map((tk) => {
              const laChinhMinh = phien !== null && Number(tk.id) === phien.id;
              return (
                <tr key={tk.id} className="border-t border-slate-800">
                  <td className="px-4 py-3 text-slate-400">{tk.id}</td>
                  <td className="px-4 py-3 font-medium text-white">
                    {tk.ten_dang_nhap}
                  </td>
                  <td className="px-4 py-3 text-slate-300">
                    {tk.la_admin ? "Quản trị" : "Người chơi"}
                  </td>
                  <td className="px-4 py-3">
                    {tk.bi_khoa ? (
                      <span className="rounded-full bg-red-500/10 px-3 py-1 text-xs font-medium text-red-300">
                        Đã khoá
                      </span>
                    ) : (
                      <span className="rounded-full bg-emerald-500/10 px-3 py-1 text-xs font-medium text-emerald-300">
                        Hoạt động
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-3 text-slate-400">
                    {new Date(tk.tao_luc).toLocaleDateString("vi-VN")}
                  </td>
                  <td className="px-4 py-3">
                    {laChinhMinh ? (
                      <span className="text-xs text-slate-500">
                        Tài khoản của bạn
                      </span>
                    ) : (
                      <form action={doiKhoaTaiKhoan}>
                        <input type="hidden" name="id" value={tk.id} />
                        <input
                          type="hidden"
                          name="trang_thai_moi"
                          value={tk.bi_khoa ? "0" : "1"}
                        />
                        <button
                          type="submit"
                          className={
                            tk.bi_khoa
                              ? "rounded-lg border border-emerald-500/40 px-3 py-1.5 text-xs font-medium text-emerald-300 transition hover:bg-emerald-500/10"
                              : "rounded-lg border border-red-500/40 px-3 py-1.5 text-xs font-medium text-red-300 transition hover:bg-red-500/10"
                          }
                        >
                          {tk.bi_khoa ? "Mở khoá" : "Khoá"}
                        </button>
                      </form>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
