"use server";

import { redirect } from "next/navigation";
import { sql } from "@/lib/db";
import { canBangThoiGian, kiemMatKhau } from "@/lib/matkhau";
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
    {
      id: string;
      ten_dang_nhap: string;
      mat_khau_bam: string;
      bi_khoa: boolean;
      la_admin: boolean;
      the_he_phien: number;
    }[]
  >`
    SELECT id, ten_dang_nhap, mat_khau_bam, bi_khoa, la_admin, the_he_phien
      FROM tai_khoan
     WHERE lower(ten_dang_nhap) = lower(${ten})
     LIMIT 1
  `;

  const taiKhoan = hang[0];

  // Thong bao loi CHUNG cho ca hai truong hop "khong ton tai" va "sai mat khau"
  // — tranh lo cho ke tan cong biet ten dang nhap nao da ton tai (do do tai
  // khoan / user enumeration).
  const THONG_BAO_SAI = "Tên đăng nhập hoặc mật khẩu không đúng.";

  // Ro ri THOI GIAN: khong duoc viet `!taiKhoan || !kiemMatKhau(...)` — toan
  // tu `||` doan mach (short-circuit) se BO QUA han kiemMatKhau() (~13ms PBKDF2)
  // khi khong tim thay tai khoan, khien nhanh "khong ton tai" tra loi NHANH
  // HON han nhanh "sai mat khau". Ke xau chi can bam gio dang nhap la biet
  // duoc ten dang nhap nao co that (user enumeration qua kenh thoi gian).
  // Vi vay: luon danh gia CA HAI nhanh — that (kiemMatKhau) hoac gia
  // (canBangThoiGian, ton dung tung do bang), roi moi quyet dinh tu choi.
  const matKhauDung = taiKhoan
    ? kiemMatKhau(matKhau, taiKhoan.mat_khau_bam)
    : canBangThoiGian(matKhau);

  if (!taiKhoan || !matKhauDung) {
    redirect(duongDanLoi(THONG_BAO_SAI, ten));
  }

  if (taiKhoan.bi_khoa) {
    redirect(duongDanLoi("Tài khoản này đã bị khoá. Vui lòng liên hệ quản trị viên.", ten));
  }

  // `theHe` di kem token de sau nay thu hoi duoc phien nay tu xa (dang xuat,
  // doi mat khau) — xem lib/xacthucphien.ts.
  await taoPhien({
    id: Number(taiKhoan.id),
    ten: taiKhoan.ten_dang_nhap,
    laAdmin: taiKhoan.la_admin,
    theHe: taiKhoan.the_he_phien,
  });
  redirect("/");
}
