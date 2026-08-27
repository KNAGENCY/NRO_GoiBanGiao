// Kiem thu cho lib/xacthucphien.ts — phan LOI cua phien dang nhap.
// Chay bang: npm run test:phien   (tu thu muc web/)
//
// Day la kiem thu TICH HOP: no can PostgreSQL that, vi dieu can chung minh
// chinh la "web co hoi CSDL moi lan doc phien hay khong". Kiem bang do gia
// (mock) thi khong chung minh duoc gi ca.
//
// ⚠️ CO TINH LAM CHO NO CHET AM I neu khong co CSDL, chu KHONG bo qua im lang.
// Ly do: §9 cua BAT-DAU-TU-DAY.md ghi lai mot vu that — test tich hop cua may
// chu Go tu bo qua khi PostgreSQL tat, bang tong ket van xanh, va 15 test do
// nam im nhieu ngay. Mot test bao mat "xanh" vi khong he chay thi con te hon
// la khong co test.
//
// Bat CSDL:  D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start
// Chi CSDL:  dat NRO_TEST_DSN (uu tien) hoac DATABASE_URL. Vi du:
//   $env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:5432/nro_test"
// Khong dat gi thi test doc web/.env.local.
//
// LUU Y ve cong PostgreSQL: no KHONG co dinh giua cac lan khoi dong (§9), nen
// test khong ghim cong o dau ca. Bao "connection refused" thi kiem cong va
// kiem CSDL da bat chua, truoc khi nghi la code sai.
//
// (Toan bo than test nam trong main() vi du an nay khong bat "type": "module",
// tsx bien ra CommonJS nen khong dung duoc `await` o muc cao nhat.)
import assert from "node:assert/strict";
import { randomBytes } from "node:crypto";
import { existsSync, readFileSync } from "node:fs";
import path from "node:path";

// Chay tu web/ (npm script dat cwd o day). Kiem som cho ro rang.
const THU_MUC_WEB = process.cwd();
if (!existsSync(path.join(THU_MUC_WEB, "package.json"))) {
  throw new Error("Hay chay test nay tu thu muc web/: cd web && npm run test:phien");
}
const GOC_REPO = path.resolve(THU_MUC_WEB, "..");

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
  process.env.NRO_TEST_DSN ||
  process.env.DATABASE_URL ||
  docEnvLocal().DATABASE_URL;

if (!DSN) {
  throw new Error(
    "Khong tim thay chuoi ket noi CSDL. Dat NRO_TEST_DSN hoac DATABASE_URL, " +
      "hoac tao web/.env.local. Test nay CAN CSDL that — xem chu thich dau file."
  );
}

// Phai dat TRUOC khi nap module: lib/db.ts nem loi ngay luc import neu thieu
// DATABASE_URL, nen cac module duoi day duoc nap bang import() trong main().
process.env.DATABASE_URL = DSN;
// Khoa ky rieng cho test: khong dung nham khoa that, va moi lan chay mot khoa
// khac nhau nen token cua lan chay truoc chac chan khong con hieu luc.
process.env.AUTH_SECRET = randomBytes(32).toString("hex");

// Bo chay ca kiem thu. CO Y khong dung lai khi mot ca do: chay het roi bao cao
// tat ca — nhin mot lan la thay het cai gi hong, thay vi sua tung cai roi chay
// lai. (Quan trong voi bo test nay: ba lo hong la ba ca doc lap.)
const caDat: string[] = [];
const caHong: { ten: string; loi: unknown }[] = [];

async function ca(ten: string, than: () => Promise<void>): Promise<void> {
  try {
    await than();
    caDat.push(ten);
    console.log(`  \u2713 ${ten}`);
  } catch (loi) {
    caHong.push({ ten, loi });
    console.log(`  \u2717 ${ten}`);
  }
}

async function main(): Promise<void> {
  const { sql } = await import("./db");
  const { kyToken, xacThucToken, huyToanBoPhien, docIdTuToken } = await import(
    "./xacthucphien"
  );

  // Danh sach tai khoan tam da tao, de don sach o cuoi.
  const daTao: number[] = [];

  // -------------------------------------------------------------------------
  // Chuan bi luoc do: bang `tai_khoan` phai co san; cot `the_he_phien` thi test
  // tu ap tu CHINH file migration 0003 neu chua co.
  //
  // Vi sao ap tu file that chu khong go tay lenh ALTER trong test: file
  // migration sai cu phap thi test nay do luon — no kiem ca migration, khong
  // chi kiem code TypeScript.
  // -------------------------------------------------------------------------
  async function chuanBiLuocDo(): Promise<void> {
    const [coBang] = await sql<{ co: boolean }[]>`
      SELECT EXISTS (
        SELECT 1 FROM information_schema.tables
         WHERE table_schema = 'public' AND table_name = 'tai_khoan'
      ) AS co
    `;
    if (!coBang.co) {
      throw new Error(
        `CSDL "${DSN}" chua co bang tai_khoan. Chay may chu Go mot lan de no tu ap ` +
          "migration (server/README.md), hoac tro NRO_TEST_DSN sang CSDL da co luoc do."
      );
    }

    const [coCot] = await sql<{ co: boolean }[]>`
      SELECT EXISTS (
        SELECT 1 FROM information_schema.columns
         WHERE table_name = 'tai_khoan' AND column_name = 'the_he_phien'
      ) AS co
    `;
    if (!coCot.co) {
      const duongDan = path.join(
        GOC_REPO,
        "server",
        "internal",
        "store",
        "migrations",
        "0003_the_he_phien.sql"
      );
      await sql.unsafe(readFileSync(duongDan, "utf8"));
      console.log("  (da ap migration 0003_the_he_phien.sql vao CSDL kiem thu)");
    }
  }

  // Ten ngau nhien de nhieu agent chay cung luc tren `nro_test` khong dam nhau.
  async function taoTaiKhoanTam(laAdmin = false): Promise<{ id: number; ten: string }> {
    const ten = "kt" + randomBytes(6).toString("hex"); // 14 ky tu, khop CHECK ^[A-Za-z0-9_]{4,20}$
    const hang = await sql<{ id: string }[]>`
      INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam, la_admin)
      VALUES (${ten}, 'khong-dung-den', ${laAdmin})
      RETURNING id
    `;
    const id = Number(hang[0].id);
    daTao.push(id);
    return { id, ten };
  }

  /** Ky token theo dung trang thai hien tai trong CSDL — giong het luc dang nhap that. */
  async function kyTokenNhuLucDangNhap(id: number): Promise<string> {
    const [tk] = await sql<
      { ten_dang_nhap: string; la_admin: boolean; the_he_phien: number }[]
    >`
      SELECT ten_dang_nhap, la_admin, the_he_phien FROM tai_khoan WHERE id = ${id}
    `;
    return kyToken({
      id,
      ten: tk.ten_dang_nhap,
      laAdmin: tk.la_admin,
      theHe: tk.the_he_phien,
    });
  }

  try {
    await chuanBiLuocDo();

    // === LO HONG 1 — KHOA TAI KHOAN GIUA CHUNG => PHIEN CHET NGAY =========
    // Truoc khi vá, `bi_khoa` chi duoc kiem luc dang nhap, nen khoa mot nguoi
    // dang co phien la vo tac dung suot 7 ngay.
    await ca("tai khoan bi khoa giua chung -> phien dang co chet NGAY", async () => {
      const tk = await taoTaiKhoanTam();
      const token = await kyTokenNhuLucDangNhap(tk.id);
      assert.notEqual(
        await xacThucToken(token),
        null,
        "chua khoa ma da khong vao duoc — sai tu buoc chuan bi"
      );

      // Admin bam nut "Khoa" trong panel — dung lenh UPDATE y het actions.ts.
      await sql`UPDATE tai_khoan SET bi_khoa = true WHERE id = ${tk.id}`;

      assert.equal(
        await xacThucToken(token),
        null,
        "TAI KHOAN BI KHOA MA PHIEN CU VAN DUNG DUOC — day chinh la lo hong"
      );
    });

    await ca("mo khoa lai -> phien cu dung duoc tro lai", async () => {
      const tk = await taoTaiKhoanTam();
      const token = await kyTokenNhuLucDangNhap(tk.id);
      await sql`UPDATE tai_khoan SET bi_khoa = true WHERE id = ${tk.id}`;
      assert.equal(await xacThucToken(token), null, "chuan bi that bai");
      await sql`UPDATE tai_khoan SET bi_khoa = false WHERE id = ${tk.id}`;
      assert.notEqual(await xacThucToken(token), null, "mo khoa roi ma van khong vao duoc");
    });

    // === LO HONG 2 — HA QUYEN ADMIN GIUA CHUNG => MAT QUYEN NGAY ==========
    // `laAdmin` truoc day bi nuong cung vao token: go quyen ai do thi ho van la
    // admin tiep 7 ngay.
    await ca(
      "admin bi ha quyen giua chung -> mat quyen NGAY, van dang nhap binh thuong",
      async () => {
        const tk = await taoTaiKhoanTam(true);
        const token = await kyTokenNhuLucDangNhap(tk.id);
        assert.equal(
          (await xacThucToken(token))?.laAdmin,
          true,
          "admin that ma doc ra khong phai admin"
        );

        await sql`UPDATE tai_khoan SET la_admin = false WHERE id = ${tk.id}`;

        const sau = await xacThucToken(token);
        assert.notEqual(sau, null, "ha quyen admin ma dang xuat luon nguoi ta — qua tay");
        assert.equal(
          sau?.laAdmin,
          false,
          "ADMIN DA BI HA QUYEN MA TOKEN CU VAN LA ADMIN — day chinh la lo hong"
        );
      }
    );

    await ca("duoc nang quyen admin -> co quyen ngay, khong phai dang nhap lai", async () => {
      const tk = await taoTaiKhoanTam(false);
      const token = await kyTokenNhuLucDangNhap(tk.id);
      assert.equal((await xacThucToken(token))?.laAdmin, false, "chuan bi that bai");
      await sql`UPDATE tai_khoan SET la_admin = true WHERE id = ${tk.id}`;
      assert.equal(
        (await xacThucToken(token))?.laAdmin,
        true,
        "duoc nang quyen admin ma phien cu chua thay"
      );
    });

    // === LO HONG 3 — DANG XUAT / DOI MAT KHAU GIET MOI BAN SAO TOKEN ======
    // Dang xuat truoc day chi xoa cookie o mot trinh duyet; ban sao token van
    // con hieu luc. `huyToanBoPhien` la thu ma xoaPhien() va doiMatKhau() goi.
    await ca("dang xuat / doi mat khau -> moi ban sao token cu chet NGAY", async () => {
      const tk = await taoTaiKhoanTam();
      const tokenBiSaoTrom = await kyTokenNhuLucDangNhap(tk.id);
      assert.notEqual(await xacThucToken(tokenBiSaoTrom), null, "chuan bi that bai");

      await huyToanBoPhien(tk.id); // <- dang xuat, hoac doi mat khau

      assert.equal(
        await xacThucToken(tokenBiSaoTrom),
        null,
        "DANG XUAT ROI MA BAN SAO TOKEN VAN DUNG DUOC — day chinh la lo hong"
      );
    });

    await ca("dang nhap lai sau khi dang xuat -> vao duoc ngay", async () => {
      const tk = await taoTaiKhoanTam();
      await huyToanBoPhien(tk.id);
      const tokenMoi = await kyTokenNhuLucDangNhap(tk.id);
      assert.notEqual(await xacThucToken(tokenMoi), null, "dang nhap lai ma van khong vao duoc");
    });

    // === KHONG DUOC CHUA LON LANH THANH LON QUE ===========================
    // Tai khoan khong bi dong den phai doc duoc phien nhieu lan lien tiep, va
    // du lieu tra ve phai dung.
    await ca(
      "nguoi dung binh thuong: doc phien 5 lan lien tiep deu dung, khong bi day ra",
      async () => {
        const tk = await taoTaiKhoanTam();
        const token = await kyTokenNhuLucDangNhap(tk.id);
        for (let i = 0; i < 5; i++) {
          assert.deepEqual(
            await xacThucToken(token),
            { id: tk.id, ten: tk.ten, laAdmin: false },
            `lan doc thu ${i + 1} tra ve sai du lieu`
          );
        }
      }
    );

    // He qua cua viec thoi tin token: ten trong CSDL doi thi web phai hien ten
    // moi. Cung la bang chung "du lieu tra ve den tu CSDL that".
    await ca("doi ten trong CSDL -> web hien ten moi ngay (du lieu den tu CSDL)", async () => {
      const tk = await taoTaiKhoanTam();
      const token = await kyTokenNhuLucDangNhap(tk.id);
      const tenMoi = "kt" + randomBytes(6).toString("hex");
      await sql`UPDATE tai_khoan SET ten_dang_nhap = ${tenMoi} WHERE id = ${tk.id}`;
      assert.equal(
        (await xacThucToken(token))?.ten,
        tenMoi,
        "ten hien thi van lay tu token cu — chua doc CSDL"
      );
    });

    await ca("tai khoan bi xoa han -> phien chet", async () => {
      const tk = await taoTaiKhoanTam();
      const token = await kyTokenNhuLucDangNhap(tk.id);
      await sql`DELETE FROM tai_khoan WHERE id = ${tk.id}`;
      assert.equal(await xacThucToken(token), null, "tai khoan da bi xoa ma phien van song");
    });

    await ca("token rong / rac / sai chu ky deu bi tu choi", async () => {
      assert.equal(await xacThucToken(undefined), null, "khong co cookie ma van cho qua");
      assert.equal(await xacThucToken(""), null, "cookie rong ma van cho qua");
      assert.equal(await xacThucToken("khong-phai-jwt"), null, "chuoi rac ma van cho qua");

      const tk = await taoTaiKhoanTam();
      const token = await kyTokenNhuLucDangNhap(tk.id);
      const phan = token.split(".");
      const chuKyHong =
        phan[2][0] === "A" ? "B" + phan[2].slice(1) : "A" + phan[2].slice(1);
      assert.equal(
        await xacThucToken(`${phan[0]}.${phan[1]}.${chuKyHong}`),
        null,
        "token bi sua chu ky ma van cho qua"
      );
    });

    // Token kieu CU (chua co truong theHe) phai bi tu choi: nguoi dang giu
    // cookie tu truoc khi vá se phai dang nhap lai mot lan — dung y muon, vi
    // dung token cu thi khong the thu hoi duoc.
    await ca("token kieu cu (khong co so the he) bi tu choi", async () => {
      const tk = await taoTaiKhoanTam();
      const tokenKieuCu = await kyToken({
        id: tk.id,
        ten: tk.ten,
        laAdmin: false,
        // @ts-expect-error co tinh bo truong theHe de gia lap token phat hanh truoc khi vá
        theHe: undefined,
      });
      assert.equal(
        await xacThucToken(tokenKieuCu),
        null,
        "token kieu cu (khong co so the he) ma van cho qua"
      );
    });

    // xoaPhien() can id de cong so the he, nhung luc do tai khoan co the dang
    // bi khoa — nen buoc doc id phai tach roi khoi buoc kiem CSDL.
    await ca("docIdTuToken: van doc duoc id ke ca khi tai khoan bi khoa", async () => {
      const tk = await taoTaiKhoanTam();
      const token = await kyTokenNhuLucDangNhap(tk.id);
      await sql`UPDATE tai_khoan SET bi_khoa = true WHERE id = ${tk.id}`;
      assert.equal(await xacThucToken(token), null, "chuan bi that bai");
      assert.equal(
        await docIdTuToken(token),
        tk.id,
        "khong doc duoc id tu token cua tai khoan bi khoa — dang xuat se khong huy duoc phien"
      );
      assert.equal(await docIdTuToken("rac"), null, "token rac ma doc ra id");
    });
  } finally {
    // Don sach du test do giua chung — CSDL `nro_test` dung chung voi agent khac.
    if (daTao.length > 0) {
      await sql`DELETE FROM tai_khoan WHERE id = ANY(${daTao})`;
    }
    await sql.end();
  }

  if (caHong.length > 0) {
    console.error(`\n${caHong.length} ca KHONG dat:\n`);
    for (const h of caHong) {
      console.error(`--- ${h.ten}`);
      console.error(h.loi instanceof Error ? h.loi.message : h.loi);
      console.error("");
    }
    throw new Error(`${caHong.length}/${caDat.length + caHong.length} ca kiem thu that bai.`);
  }

  console.log(`\nOK — ${caDat.length} ca kiem thu lib/xacthucphien.ts deu qua.`);
}

main().catch((err) => {
  console.error(err instanceof Error ? err.message : err);
  process.exit(1);
});
