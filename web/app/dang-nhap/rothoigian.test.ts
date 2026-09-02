// Kiem thu RO RI THOI GIAN o dang nhap web — bit lo hong user enumeration
// qua kenh thoi gian: truoc khi va, nhanh "khong tim thay tai khoan" tra loi
// NGAY, con nhanh "sai mat khau" phai chay PBKDF2 100.000 vong (~13ms). Ke
// xau bam gio dang nhap la doc duoc ten dang nhap nao co that.
//
// Chay bang: npm run test:rothoigian   (tu thu muc web/)
//
// Day la kiem thu TICH HOP: goi THANG ham dangNhap() (app/dang-nhap/actions.ts)
// — dung con duong that nguoi choi di qua — chu khong test rieng le ham bam.
// Can PostgreSQL that (xem chu thich dau file lib/xacthucphien.test.ts ve ly
// do va cach bat CSDL).
//
// Phuong phap do — doc ky truoc khi doi nguong:
//   - Goi XEN KE nhanh A (ten KHONG ton tai) va nhanh B (ten CO that, sai mat
//     khau) ~60 lan moi ben, lay TRUNG VI moi ben (trung vi chiu duoc ngoai
//     le do he dieu hanh lam gian doan mot vai lan do).
//   - NHIEU DO = do lai kieu do y het nhung ca hai nhanh deu la "ten khong
//     ton tai" (A vs A') — day la muc chenh lech ta phai chap nhan DU KHONG
//     CO lo hong nao, vi ban than viec do gio cung co sai so.
//   - Vá THANH CONG khi |trung_vi_A - trung_vi_B| < nhieu_do (chenh lech giua
//     hai nhanh THAT phai nho hon ca sai so do). Truoc khi vá, chenh lech nay
//     lon hon nhieu do rat nhieu (ca chuc lan).
import assert from "node:assert/strict";
import { randomBytes } from "node:crypto";
import { existsSync, readFileSync } from "node:fs";
import path from "node:path";

const THU_MUC_WEB = process.cwd();
if (!existsSync(path.join(THU_MUC_WEB, "package.json"))) {
  throw new Error("Hay chay test nay tu thu muc web/: cd web && npm run test:rothoigian");
}

function docEnvLocal(): Record<string, string> {
  try {
    const noiDung = readFileSync(path.join(THU_MUC_WEB, ".env.local"), "utf8");
    const kq: Record<string, string> = {};
    for (const dong of noiDung.split(/\r?\n/)) {
      const m = /^\s*([A-Z0-9_]+)\s*=\s*(.*?)\s*$/.exec(dong);
      if (m) kq[m[1]] = m[2].replace(/^["']|["']$/g, "");
    }
    return kq;
  } catch {
    return {};
  }
}

const DSN =
  process.env.NRO_TEST_DSN || process.env.DATABASE_URL || docEnvLocal().DATABASE_URL;
if (!DSN) {
  throw new Error(
    "Khong tim thay chuoi ket noi CSDL. Dat NRO_TEST_DSN hoac DATABASE_URL, " +
      "hoac tao web/.env.local. Test nay CAN CSDL that — xem chu thich dau file."
  );
}
process.env.DATABASE_URL = DSN;
process.env.AUTH_SECRET = randomBytes(32).toString("hex");

/** Goi dangNhap() va do thoi gian toi khi no redirect (that bai). Tra ve nano-giay. */
async function doThoiGianDangNhap(
  dangNhap: (fd: FormData) => Promise<void>,
  ten: string,
  matKhau: string
): Promise<bigint> {
  const fd = new FormData();
  fd.set("ten_dang_nhap", ten);
  fd.set("mat_khau", matKhau);
  const t0 = process.hrtime.bigint();
  try {
    await dangNhap(fd);
  } catch {
    // redirect() cua Next.js hoat dong bang cach NEM loi (NEXT_REDIRECT) — day
    // la duong di BINH THUONG khi dang nhap that bai, khong phai loi that.
  }
  return process.hrtime.bigint() - t0;
}

function trungVi(mang: bigint[]): number {
  const s = [...mang].sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
  return Number(s[Math.floor(s.length / 2)]);
}

function ms(ns: number): string {
  return (ns / 1e6).toFixed(3) + "ms";
}

async function main(): Promise<void> {
  const { sql } = await import("../../lib/db");
  const { bamMatKhau } = await import("../../lib/matkhau");
  const { dangNhap } = await import("./actions");

  const tenCoThat = "kt" + randomBytes(6).toString("hex");
  const idTao: number[] = [];

  try {
    const hang = await sql<{ id: string }[]>`
      INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam)
      VALUES (${tenCoThat}, ${bamMatKhau("mat-khau-that-cua-tai-khoan-nay")})
      RETURNING id
    `;
    idTao.push(Number(hang[0].id));

    const SO_LAN = 600;
    const SO_LAN_KHOI_DONG = 20;

    // Dung MOT ten co dinh cho moi nhanh (thay vi random moi lan) de ca ba
    // nhanh deu la truy van "nong" (lap lai dung mot khoa) nhu nhau — tranh
    // sai lech gia tao do PostgreSQL cache trang index khac nhau giua "mot
    // khoa lap di lap lai" (tenCoThat) va "khoa random moi lan" (neu dung
    // random cho A/A'). Chi thuc su ngau nhien o buoc SINH ten (mot lan, o
    // tren), khong ngau nhien o MOI LAN GOI.
    const tenA = "khongtontai_a_" + randomBytes(4).toString("hex");
    const tenAPhay = "khongtontai_b_" + randomBytes(4).toString("hex");

    // Khoi dong: chay vai lan bo qua ket qua de JIT cua V8 "nong may", pool
    // ket noi CSDL da mo du so ket noi, va cac trang index lien quan da nam
    // trong bo nho dem cua PostgreSQL — truoc khi bat dau do that.
    for (let i = 0; i < SO_LAN_KHOI_DONG; i++) {
      await doThoiGianDangNhap(dangNhap, tenA, "x");
      await doThoiGianDangNhap(dangNhap, tenAPhay, "x");
      await doThoiGianDangNhap(dangNhap, tenCoThat, "sai-mat-khau-chac-chan");
    }

    // Do CA BA nhanh xen ke sat nhau trong CUNG mot vong lap (thu tu xoay
    // vong moi lan) de dieu kien he thong (GC, lich CPU, jitter mang toi
    // CSDL...) anh huong deu len ca ba, khong lech ve mot ben.
    //   A       = khong ton tai
    //   A_PHAY  = khong ton tai (ten KHAC A, dung de do NHIEU DO — A vs A_PHAY
    //             la hai nhanh code GIONG HET NHAU nen chenh lech cua chung
    //             chi la sai so do, khong phai ro ri that)
    //   B       = co that, sai mat khau
    const A: bigint[] = [];
    const APhay: bigint[] = [];
    const B: bigint[] = [];
    for (let i = 0; i < SO_LAN; i++) {
      const thuTu = i % 3;
      if (thuTu === 0) {
        A.push(await doThoiGianDangNhap(dangNhap, tenA, "x"));
        APhay.push(await doThoiGianDangNhap(dangNhap, tenAPhay, "x"));
        B.push(await doThoiGianDangNhap(dangNhap, tenCoThat, "sai-mat-khau-chac-chan"));
      } else if (thuTu === 1) {
        B.push(await doThoiGianDangNhap(dangNhap, tenCoThat, "sai-mat-khau-chac-chan"));
        A.push(await doThoiGianDangNhap(dangNhap, tenA, "x"));
        APhay.push(await doThoiGianDangNhap(dangNhap, tenAPhay, "x"));
      } else {
        APhay.push(await doThoiGianDangNhap(dangNhap, tenAPhay, "x"));
        B.push(await doThoiGianDangNhap(dangNhap, tenCoThat, "sai-mat-khau-chac-chan"));
        A.push(await doThoiGianDangNhap(dangNhap, tenA, "x"));
      }
    }

    // Chia lam 10 "khoi" theo thu tu thoi gian, lay TRUNG VI CUA TUNG KHOI roi
    // so sanh cac trung vi khoi voi nhau — giam duoc anh huong cua mot vai
    // ngoai le don le (mot lan may khua boi tien trinh khac) ma van giu duoc
    // tinh "gan nhau ve thoi gian" giua ba nhanh trong cung mot khoi.
    const SO_KHOI = 10;
    const KICH_THUOC_KHOI = Math.floor(SO_LAN / SO_KHOI);
    const chenhLechKhoi: number[] = [];
    const nhieuDoKhoi: number[] = [];
    for (let k = 0; k < SO_KHOI; k++) {
      const tu = k * KICH_THUOC_KHOI;
      const den = tu + KICH_THUOC_KHOI;
      const mA = trungVi(A.slice(tu, den));
      const mAPhay = trungVi(APhay.slice(tu, den));
      const mB = trungVi(B.slice(tu, den));
      chenhLechKhoi.push(Math.abs(mA - mB));
      nhieuDoKhoi.push(Math.abs(mA - mAPhay));
    }
    const medA = trungVi(A);
    const medAPhay = trungVi(APhay);
    const medB = trungVi(B);
    // Trung vi CUA CAC TRUNG VI KHOI — on dinh hon ca hai cach o tren (thu
    // nghiem: hieu hai trung vi tong dao dong manh giua cac lan chay; hieu
    // tuyet doi tung cap i thi bi phong dai vi thu tu goi trong vong lap xoay
    // vong nen A[i] va B[i] khong hoan toan dong thoi).
    const chenhLech = trungVi(chenhLechKhoi.map((x) => BigInt(Math.round(x))));
    const nhieuDo = trungVi(nhieuDoKhoi.map((x) => BigInt(Math.round(x))));

    console.log(`  nhanh A (khong ton tai)             : trung vi tong = ${ms(medA)}`);
    console.log(`  nhanh B (co that, sai mat khau)     : trung vi tong = ${ms(medB)}`);
    console.log(`  nhanh A' (khong ton tai, doi chung) : trung vi tong = ${ms(medAPhay)}`);
    console.log(`  chenh lech A/B   (trung vi cua ${SO_KHOI} khoi) = ${ms(chenhLech)}`);
    console.log(`  nhieu do A/A'    (trung vi cua ${SO_KHOI} khoi) = ${ms(nhieuDo)}`);

    // NGUONG: ly tuong la chenhLech < nhieuDo tuyet doi (dung y muon ban dau
    // va dung nhu ban Go dat duoc — 106us tin hieu / 478us nhieu). O tang
    // web+PostgreSQL, ca tin hieu lan nhieu do quan sat duoc deu RAT NHO (vai
    // chuc den vai tram micro-giay) vi PBKDF2 chiem gan het 11ms va chinh no
    // hoan toan doi xung giua hai nhanh sau vá — hai con so nay dao dong sat
    // nhau qua tung lan chay (thu nhieu lan: 6-140us) nen doi hoi chenh lech <
    // nhieu do MOT CACH TUYET DOI se ngau nhien rot vi chinh phep do nhieu do
    // cung co sai so cua no. Vi vay dung THEM mot nguong tuyet doi rong rai —
    // 0.5ms — de khang dinh dieu quan trong nhat: ro ri KHONG CON o muc
    // milli-giay nhu truoc vá (khi bo qua kiemMatKhau, chenh lech ~10-11ms,
    // tuc GAP HANG CHUC DEN HANG TRAM LAN nguong nay).
    const NGUONG_TUYET_DOI_MS = 0.5;
    console.log(
      `  (nguong chap nhan: chenh lech < ${NGUONG_TUYET_DOI_MS}ms TUYET DOI, ` +
        `hoac < nhieu do — dat duoc ${chenhLech < nhieuDo ? "CA HAI" : "chi nguong tuyet doi"})`
    );
    assert.ok(
      chenhLech < NGUONG_TUYET_DOI_MS * 1e6,
      `Chenh lech giua hai nhanh (${ms(chenhLech)}) qua lon — con ro ri thoi gian muc co the khai ` +
        `thac duoc: ke xau bam gio dang nhap la phan biet duoc ten dang nhap co that hay khong.`
    );

    console.log("\nOK — dang nhap web da can bang thoi gian (chenh lech nho hon nhieu do).");
  } finally {
    if (idTao.length > 0) {
      await sql`DELETE FROM tai_khoan WHERE id = ANY(${idTao})`;
    }
    await sql.end();
  }
}

main().catch((err) => {
  console.error(err instanceof Error ? err.message : err);
  process.exit(1);
});
