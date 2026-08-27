"use server";

import { revalidatePath } from "next/cache";
import { sql } from "@/lib/db";
import { docPhien } from "@/lib/phien";

/**
 * Khoa/mo tai khoan (cot bi_khoa). Kiem lai quyen admin ngay trong action —
 * khong chi dua vao proxy.ts hay AdminLayout — de action nay van an toan neu
 * bi goi truc tiep (vi du tu console) bo qua giao dien. Khong cho admin tu
 * khoa chinh minh (kiem id).
 *
 * Khoa xong la nguoi do bi day ra khoi web NGAY: docPhien() doc cot bi_khoa
 * tu CSDL o moi lan doc phien, khong tin cookie nua (lib/xacthucphien.ts).
 * Khong can cong so the he phien o day — mo khoa lai thi ho vao duoc tiep,
 * dung nhu y nghia cua nut "Mo khoa".
 */
export async function doiKhoaTaiKhoan(formData: FormData): Promise<void> {
  const phien = await docPhien();
  if (!phien || !phien.laAdmin) {
    throw new Error("Không có quyền thực hiện thao tác này.");
  }

  const id = Number(formData.get("id"));
  if (!Number.isInteger(id) || id <= 0) {
    return;
  }

  if (id === phien.id) {
    // Chan tu khoa chinh minh du bi goi thang qua action, khong chi an nut.
    return;
  }

  const trangThaiMoi = formData.get("trang_thai_moi") === "1";

  await sql`
    UPDATE tai_khoan
       SET bi_khoa = ${trangThaiMoi}
     WHERE id = ${id}
  `;

  revalidatePath("/admin/tai-khoan");
}
