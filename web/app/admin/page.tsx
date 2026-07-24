import { sql } from "@/lib/db";

export default async function TrangAdmin() {
  const [tongTaiKhoan] = await sql<{ tong: string }[]>`
    SELECT COUNT(*)::text AS tong FROM tai_khoan
  `;
  const [tongNhanVat] = await sql<{ tong: string }[]>`
    SELECT COUNT(*)::text AS tong FROM nhan_vat
  `;
  const [tongBiKhoa] = await sql<{ tong: string }[]>`
    SELECT COUNT(*)::text AS tong FROM tai_khoan WHERE bi_khoa = true
  `;

  const theSo = [
    { nhan: "Tổng tài khoản", giaTri: tongTaiKhoan.tong },
    { nhan: "Tổng nhân vật", giaTri: tongNhanVat.tong },
    { nhan: "Tài khoản bị khoá", giaTri: tongBiKhoa.tong },
  ];

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold text-white">Bảng điều khiển</h1>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        {theSo.map((t) => (
          <div
            key={t.nhan}
            className="rounded-xl border border-slate-700 bg-slate-900/50 p-6"
          >
            <p className="text-sm text-slate-400">{t.nhan}</p>
            <p className="mt-2 text-3xl font-bold text-white">
              {Number(t.giaTri).toLocaleString("vi-VN")}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
