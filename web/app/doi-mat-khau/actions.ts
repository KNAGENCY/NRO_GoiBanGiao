"use server";

import { redirect } from "next/navigation";
import { sql } from "@/lib/db";
import { bamMatKhau, kiemMatKhau } from "@/lib/matkhau";
import { docPhien, taoPhien } from "@/lib/phien";

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

  // Doi mat khau xong thi day MOI phien cu ra ngoai — day la ly do chinh
  // nguoi ta doi mat khau: nghi co ke khac dang dung tai khoan cua minh. Cong
  // 1 vao `the_he_phien` ngay trong cung lenh UPDATE (mot vong CSDL, va khong
  // co khe cua nao giua hai buoc) lam moi token da phat hanh chet ngay.
  const matKhauBamMoi = bamMatKhau(matKhauMoi);
  const hangMoi = await sql<
    { ten_dang_nhap: string; la_admin: boolean; the_he_phien: number }[]
  >`
    UPDATE tai_khoan
       SET mat_khau_bam = ${matKhauBamMoi},
           the_he_phien = the_he_phien + 1
     WHERE id = ${phien.id}
    RETURNING ten_dang_nhap, la_admin, the_he_phien
  `;

  // ...ke ca token cua CHINH trinh duyet nay. Nen cap lai cookie moi ngay, de
  // nguoi vua doi mat khau khong bi da ra trang dang nhap mot cach kho hieu.
  const tk = hangMoi[0];
  if (tk) {
    await taoPhien({
      id: phien.id,
      ten: tk.ten_dang_nhap,
      laAdmin: tk.la_admin,
      theHe: tk.the_he_phien,
    });
  }

  redirect("/tai-khoan?doi_mat_khau=1");
}
