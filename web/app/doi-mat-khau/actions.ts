"use server";

import { redirect } from "next/navigation";
import { sql } from "@/lib/db";
import { bamMatKhau, kiemMatKhau } from "@/lib/matkhau";
import { docPhien } from "@/lib/phien";

function duongDanLoi(loi: string): string {
  const qs = new URLSearchParams({ loi });
  return `/doi-mat-khau?${qs.toString()}`;
}

export async function doiMatKhau(formData: FormData): Promise<void> {
  const phien = await docPhien();
  if (!phien) {
    redirect("/dang-nhap");
  }

  const matKhauCu = String(formData.get("mat_khau_cu") ?? "");
  const matKhauMoi = String(formData.get("mat_khau_moi") ?? "");
  const xacNhan = String(formData.get("xac_nhan_mat_khau_moi") ?? "");

  if (!matKhauCu || !matKhauMoi) {
    redirect(
      duongDanLoi("Vui lòng nhập đầy đủ mật khẩu cũ và mật khẩu mới.")
    );
  }

  if (matKhauMoi.length < 6) {
    redirect(duongDanLoi("Mật khẩu mới phải có ít nhất 6 ký tự."));
  }

  if (matKhauMoi !== xacNhan) {
    redirect(duongDanLoi("Xác nhận mật khẩu mới không khớp."));
  }

  const hang = await sql<{ mat_khau_bam: string }[]>`
    SELECT mat_khau_bam FROM tai_khoan WHERE id = ${phien.id} LIMIT 1
  `;
  const taiKhoan = hang[0];

  if (!taiKhoan || !kiemMatKhau(matKhauCu, taiKhoan.mat_khau_bam)) {
    redirect(duongDanLoi("Mật khẩu cũ không đúng."));
  }

  const matKhauBamMoi = bamMatKhau(matKhauMoi);
  await sql`
    UPDATE tai_khoan
       SET mat_khau_bam = ${matKhauBamMoi}
     WHERE id = ${phien.id}
  `;

  redirect("/tai-khoan?doi_mat_khau=1");
}
