// Kiem thu nhanh cho lib/matkhau.ts. Chay bang: npm run test:matkhau
//
// Day la kiem thu NOI BO (chi Node, khong dung Go). Kiem tuong thich HAI
// CHIEU that su voi may chu Go nam o script rieng — xem README.md phan
// "Kiem tuong thich mat khau voi may chu game".
import assert from "node:assert/strict";
import { BAM_GIA, bamMatKhau, canBangThoiGian, kiemMatKhau, laDaBam } from "./matkhau";

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

// === CAN BANG THOI GIAN (bit ro ri user enumeration qua kenh thoi gian) ===
// Ba test duoi day khoa ba cach BAM_GIA / canBangThoiGian co the am tham hong
// — y het ben Go (server/internal/auth/canbangthoigian_test.go).

// 6) Neu ai do doi SO_VONG ma quen sinh lai BAM_GIA, chuoi gia se re hon
// chuoi that va lo hong quay lai NGAY, khong co trieu chung nao khac.
{
  const phan = BAM_GIA.split("$");
  assert.equal(phan.length, 4, "BAM_GIA sai dinh dang");
  assert.equal(
    phan[1],
    "100000",
    "BAM_GIA khai so vong khac voi SO_VONG dang dung — bam gia re hon bam that " +
      "thi hai nhanh dang nhap lai chenh nhau thoi gian. Sinh lai BAM_GIA bang bamMatKhau()."
  );
  assert.equal(laDaBam(BAM_GIA), true, "BAM_GIA khong o dang da bam");
}

// 7) BAM_GIA khong duoc khop mat khau nao — neu khop, ai go dung chuoi do se
// "dang nhap thanh cong" vao mot tai khoan khong ton tai.
{
  for (const mk of ["", "123456", "matkhau", "password", "admin", BAM_GIA]) {
    assert.equal(kiemMatKhau(mk, BAM_GIA), false, `mat khau "${mk}" khop voi BAM_GIA`);
  }
  assert.equal(canBangThoiGian("matkhau"), false, "canBangThoiGian phai LUON tra ve false");
}

// 8) Test dot bien quan trong nhat: neu canBangThoiGian bi "toi uu" thanh ham
// rong (hoac khong con goi kiemMatKhau ben trong), no van tra ve dung nhung
// khong con ton thoi gian — lo hong quay lai trong im lang. Do bang so la
// cach duy nhat bat duoc chuyen do.
{
  const bamThat = bamMatKhau("matkhaunguoidung");
  const soLan = 40;
  const that: number[] = [];
  const gia: number[] = [];
  // Xen ke hai phep do de deu nhieu do (CPU ha xung nhip, GC...) len ca hai
  // ben thay vi don het vao mot ben.
  for (let i = 0; i < soLan; i++) {
    const t0 = process.hrtime.bigint();
    kiemMatKhau("saibet", bamThat);
    that.push(Number(process.hrtime.bigint() - t0));

    const t1 = process.hrtime.bigint();
    canBangThoiGian("saibet");
    gia.push(Number(process.hrtime.bigint() - t1));
  }
  const trungVi = (mang: number[]): number => {
    const s = [...mang].sort((a, b) => a - b);
    return s[Math.floor(s.length / 2)];
  };
  const medThat = trungVi(that);
  const medGia = trungVi(gia);
  console.log(
    `  (doi chieu that = ${(medThat / 1e6).toFixed(3)}ms | can bang thoi gian = ${(
      medGia / 1e6
    ).toFixed(3)}ms)`
  );
  // Doi hoi hai ben nam trong khoang 25% cua nhau — rong rai du de khong chop
  // tat vi may Windows dang lam viec khac, nhung chat hon nhieu so voi kieu
  // hong that (ham rong -> gia gan bang 0, tuc lech ~100%).
  assert.ok(
    medGia * 4 >= medThat * 3 && medGia * 3 <= medThat * 4,
    `canBangThoiGian ton ${medGia}ns trong khi mot lan doi chieu that ton ${medThat}ns — ` +
      "hai ben phai xap xi nhau, neu khong nhanh 'khong co tai khoan' van nhanh hon nhanh " +
      "'sai mat khau' va ke xau bam gio la liet ke duoc tai khoan"
  );
}

console.log("OK — tat ca kiem thu lib/matkhau.ts deu qua.");
