// Kiem thu nhanh cho lib/matkhau.ts. Chay bang: npm run test:matkhau
//
// Day la kiem thu NOI BO (chi Node, khong dung Go). Kiem tuong thich HAI
// CHIEU that su voi may chu Go nam o script rieng — xem README.md phan
// "Kiem tuong thich mat khau voi may chu game".
import assert from "node:assert/strict";
import { bamMatKhau, kiemMatKhau, laDaBam } from "./matkhau";

// 1) Bam roi doi chieu lai phai dung.
{
  const h = bamMatKhau("matkhau123");
  assert.equal(kiemMatKhau("matkhau123", h), true, "bam roi doi chieu that bai");
  assert.equal(kiemMatKhau("sai-mat-khau", h), false, "mat khau sai ma van cho qua");
}

// 2) Dinh dang chuoi bam phai dung: pbkdf2$100000$<muoi 16 byte b64>$<bam 32 byte b64>
{
  const h = bamMatKhau("gi-do");
  const phan = h.split("$");
  assert.equal(phan.length, 4, "chuoi bam khong dung 4 phan");
  assert.equal(phan[0], "pbkdf2", "thieu tien to pbkdf2");
  assert.equal(phan[1], "100000", "so vong phai la 100000");
  assert.equal(Buffer.from(phan[2], "base64").length, 16, "muoi phai 16 byte");
  assert.equal(Buffer.from(phan[3], "base64").length, 32, "bam phai 32 byte (256 bit)");
  assert.equal(laDaBam(h), true, "laDaBam phai nhan dien duoc chuoi vua bam");
}

// 3) Muoi ngau nhien: bam hai lan cung mat khau phai ra hai chuoi khac nhau.
{
  const a = bamMatKhau("abc");
  const b = bamMatKhau("abc");
  assert.notEqual(a, b, "hai lan bam giong nhau — muoi khong ngau nhien");
}

// 4) Chuoi bam hong khong duoc lam sap (phai tra false, khong duoc nem loi).
{
  for (const x of ["pbkdf2$", "pbkdf2$abc$x$y", "pbkdf2$100000$!!!$!!!", ""]) {
    assert.equal(kiemMatKhau("bat ky", x), false, `chuoi hong "${x}" ma van cho qua`);
  }
}

// 5) Doi chieu duoc chuoi bam do CHINH may chu Go sinh ra (lay tu
// server/internal/auth/password_test.go — mat khau la "1"). Day la bang
// chung tinh, khong thay the buoc kiem hai chieu that su chay Go, nhung bat
// duoc sai lech dinh dang som, ngay tren may khong co Go.
{
  const bamTuGo =
    "pbkdf2$100000$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU=";
  assert.equal(kiemMatKhau("1", bamTuGo), true, "khong doi chieu duoc chuoi bam co san tu Go");
  assert.equal(kiemMatKhau("2", bamTuGo), false, "mat khau sai ma van cho qua (chuoi tu Go)");
}

console.log("OK — tat ca kiem thu lib/matkhau.ts deu qua.");
