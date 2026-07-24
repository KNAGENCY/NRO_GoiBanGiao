"use server";

import { redirect } from "next/navigation";
import { sql } from "@/lib/db";
import { kiemMatKhau } from "@/lib/matkhau";
import { taoPhien } from "@/lib/phien";

function duongDanLoi(loi: string, ten?: string): string {
  const qs = new URLSearchParams({ loi });
  if (ten) qs.set("ten", ten);
  return `/dang-nhap?${qs.toString()}`;
}

export async function dangNhap(formData: FormData): Promise<void> {
  const ten = String(formData.get("ten_dang_nhap") ?? "").trim();
  const matKhau = String(formData.get("mat_khau") ?? "");

  if (!ten || !matKhau) {
    redirect(duongDanLoi("Vui lòng nhập đầy đủ tên đăng nhập và mật khẩu.", ten));
  }

  // Chu y: cot id la BIGSERIAL — postgres.js tra ve dang CHUOI (khong phai
  // number) de tranh mat do chinh xac voi so 64-bit. Phai ep kieu ve number
  // truoc khi dua vao phien (xem taoPhien ben duoi).
  const hang = await sql<
    { id: string; ten_dang_nhap: string; mat_khau_bam: string; bi_khoa: boolean; la_admin: boolean }[]
  >`
    SELECT id, ten_dang_nhap, mat_khau_bam, bi_khoa, la_admin
      FROM tai_khoan
     WHERE lower(ten_dang_nhap) = lower(${ten})
     LIMIT 1
  `;

  const taiKhoan = hang[0];

  // Thong bao loi CHUNG cho ca hai truong hop "khong ton tai" va "sai mat khau"
  // — tranh lo cho ke tan cong biet ten dang nhap nao da ton tai (do do tai
  // khoan / user enumeration).
  const THONG_BAO_SAI = "Tên đăng nhập hoặc mật khẩu không đúng.";

  if (!taiKhoan || !kiemMatKhau(matKhau, taiKhoan.mat_khau_bam)) {
    redirect(duongDanLoi(THONG_BAO_SAI, ten));
  }

  if (taiKhoan.bi_khoa) {
    redirect(duongDanLoi("Tài khoản này đã bị khoá. Vui lòng liên hệ quản trị viên.", ten));
  }

  await taoPhien({
    id: Number(taiKhoan.id),
    ten: taiKhoan.ten_dang_nhap,
    laAdmin: taiKhoan.la_admin,
  });
  redirect("/");
}
