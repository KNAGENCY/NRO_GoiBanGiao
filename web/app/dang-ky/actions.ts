"use server";

import { redirect } from "next/navigation";
import { sql } from "@/lib/db";
import { bamMatKhau } from "@/lib/matkhau";
import { taoPhien } from "@/lib/phien";

const RE_TEN_HOP_LE = /^[A-Za-z0-9_]{4,20}$/;

function duongDanLoi(loi: string, ten?: string): string {
  const qs = new URLSearchParams({ loi });
  if (ten) qs.set("ten", ten);
  return `/dang-ky?${qs.toString()}`;
}

export async function dangKy(formData: FormData): Promise<void> {
  const ten = String(formData.get("ten_dang_nhap") ?? "").trim();
  const matKhau = String(formData.get("mat_khau") ?? "");
  const xacNhan = String(formData.get("xac_nhan_mat_khau") ?? "");

  if (!RE_TEN_HOP_LE.test(ten)) {
    redirect(
      duongDanLoi(
        "Tên đăng nhập phải dài 4-20 ký tự, chỉ gồm chữ cái không dấu, số và gạch dưới (_).",
        ten
      )
    );
  }

  if (matKhau.length < 6) {
    redirect(duongDanLoi("Mật khẩu phải có ít nhất 6 ký tự.", ten));
  }

  if (matKhau !== xacNhan) {
    redirect(duongDanLoi("Xác nhận mật khẩu không khớp.", ten));
  }

  const matKhauBam = bamMatKhau(matKhau);

  let idTaiKhoan: number;
  let theHePhien: number;
  try {
    // Chu y: cot id la BIGSERIAL — postgres.js tra ve dang CHUOI (khong phai
    // number) de tranh mat do chinh xac voi so 64-bit. Phai ep kieu ve number
    // truoc khi dua vao phien, neu khong docPhien() se coi payload khong hop le.
    //
    // Lay luon `the_he_phien` bang RETURNING thay vi SELECT lai: token phai
    // mang theo so the he thi sau nay moi thu hoi duoc phien (xem
    // lib/xacthucphien.ts). Tai khoan moi thi so nay bang 0, nhung doc that ra
    // van hon ghi cung — mai sau doi mac dinh trong migration khong vo cho nay.
    const ketQua = await sql<{ id: string; the_he_phien: number }[]>`
      INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam)
      VALUES (${ten}, ${matKhauBam})
      RETURNING id, the_he_phien
    `;
    idTaiKhoan = Number(ketQua[0].id);
    theHePhien = ketQua[0].the_he_phien;
  } catch (err: unknown) {
    const ma = (err as { code?: string } | undefined)?.code;
    if (ma === "23505") {
      // Vi pham chi muc duy nhat tai_khoan_ten_duy_nhat (lower(ten_dang_nhap)).
      redirect(duongDanLoi("Tên đăng nhập đã được sử dụng, vui lòng chọn tên khác.", ten));
    }
    if (ma === "23514") {
      // Vi pham CHECK tai_khoan_ten_hop_le o tang CSDL (phong hach du da kiem o tren).
      redirect(duongDanLoi("Tên đăng nhập không hợp lệ.", ten));
    }
    console.error("Loi dang ky tai khoan:", err);
    redirect(duongDanLoi("Có lỗi xảy ra, vui lòng thử lại sau.", ten));
  }

  await taoPhien({ id: idTaiKhoan, ten, laAdmin: false, theHe: theHePhien });
  redirect("/");
}
