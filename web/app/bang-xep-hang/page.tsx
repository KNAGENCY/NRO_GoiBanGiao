import { sql } from "@/lib/db";
import { tenLang } from "@/lib/lang";

// Bang xep hang phai doc du lieu SONG moi lan mo (suc manh nhan vat doi lien tuc
// trong game). Khong co export nay, Next.js prerender trang thanh HTML tinh luc
// build -> bang xep hang dong bang mai o snapshot luc build. force-dynamic bat
// render dong moi request.
export const dynamic = "force-dynamic";

export default async function TrangBangXepHang() {
  // Top 50 nhan vat manh nhat. Khong tham so nguoi dung nao trong cau truy
  // van nay (LIMIT la hang so co dinh) nen khong can tham so hoa gi them.
  const hang = await sql<
    { id: string; ten: string; lang: number; cap: number; suc_manh: string }[]
  >`
    SELECT id, ten, lang, cap, suc_manh
      FROM nhan_vat
     ORDER BY suc_manh DESC
     LIMIT 50
  `;

  return (
    <main className="mx-auto flex min-h-screen max-w-3xl flex-col gap-6 px-6 py-16">
      <h1 className="text-2xl font-bold text-white">Bảng xếp hạng</h1>
      <p className="text-sm text-slate-400">
        Top 50 nhân vật mạnh nhất, xếp theo sức mạnh.
      </p>

      <div className="overflow-x-auto rounded-xl border border-slate-700">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-800/60 text-slate-300">
            <tr>
              <th className="px-4 py-3">Hạng</th>
              <th className="px-4 py-3">Tên</th>
              <th className="px-4 py-3">Làng</th>
              <th className="px-4 py-3">Cấp</th>
              <th className="px-4 py-3">Sức mạnh</th>
            </tr>
          </thead>
          <tbody>
            {hang.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-6 text-center text-slate-400">
                  Chưa có nhân vật nào.
                </td>
              </tr>
            )}
            {hang.map((nv, chiSo) => (
              <tr key={nv.id} className="border-t border-slate-800">
                <td className="px-4 py-3 text-slate-300">{chiSo + 1}</td>
                <td className="px-4 py-3 font-medium text-white">{nv.ten}</td>
                <td className="px-4 py-3 text-slate-300">{tenLang(nv.lang)}</td>
                <td className="px-4 py-3 text-slate-300">{nv.cap}</td>
                <td className="px-4 py-3 text-orange-300">
                  {Number(nv.suc_manh).toLocaleString("vi-VN")}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </main>
  );
}
