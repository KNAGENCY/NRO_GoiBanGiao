// Ket noi PostgreSQL dung chung voi may chu game (Go). Dung thu vien
// postgres.js — nhe, khong can ORM. Chi mot pool duy nhat cho ca ung dung
// (Next.js dev mode reload module nhieu lan nen luu vao globalThis de khong
// mo tran lan ket noi).
import postgres from "postgres";

const DATABASE_URL = process.env.DATABASE_URL;

if (!DATABASE_URL) {
  throw new Error(
    "Thieu bien moi truong DATABASE_URL. Tao file web/.env.local tu " +
      "web/.env.local.example roi dien chuoi ket noi PostgreSQL."
  );
}

declare global {
  // eslint-disable-next-line no-var
  var __nroSql: ReturnType<typeof postgres> | undefined;
}

export const sql =
  globalThis.__nroSql ??
  postgres(DATABASE_URL, {
    // Giu it ket noi trong dev; production co the tang qua bien moi truong sau.
    max: 10,
  });

if (process.env.NODE_ENV !== "production") {
  globalThis.__nroSql = sql;
}
