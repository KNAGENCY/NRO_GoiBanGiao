-- 0005_vat_pham.sql — KH7 Task 1: bang vat pham + lich su + giao dich.
--
-- Thiet ke day du: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md §5.1.
-- Moi rang buoc o day la mot lop phong thu CHONG NHAN BAN o tang CSDL — dung ngay
-- ca khi toan bo tang Go co bug (§5.2). Danh so 0005 (0004 da la thuong_cho_nhan).
--
-- Vi sao rang buoc nam o CSDL chu khong chi o Go: web panel, script chuyen du lieu,
-- va tay nguoi van hanh deu cham thang bang nay — giong bai hoc `nhan_vat` CHECK ten
-- o 0002. `vat_pham_mot_o_mot_mon` la tuyen phong thu CUOI CUNG cua BB1.

-- Moi MON (instance) mot dong. id BIGSERIAL, KHONG BAO GIO tai su dung (mon bi tieu
-- chi doi noi_chua='DA_HUY', khong xoa dong — de nhat ky doi chu khong noi doi).
CREATE TABLE IF NOT EXISTS vat_pham (
    id              BIGSERIAL   PRIMARY KEY,
    mau_id          INTEGER     NOT NULL,                 -- ItemTemplate.id (short xuong client)
    chu_so_huu_id   BIGINT      REFERENCES nhan_vat(id),  -- NULL khi duoi dat / da huy
    noi_chua        TEXT        NOT NULL,                 -- TUI | MAC | RUONG | DUOI_DAT | KY_GUI | DA_HUY
    o_so            SMALLINT,                             -- vi tri hien thi; NULL khi khong nam trong o nao
    so_luong        INTEGER     NOT NULL DEFAULT 1,
    chi_so_phu      JSONB       NOT NULL DEFAULT '{"opt":[]}'::jsonb,
    tao_luc         TIMESTAMPTZ NOT NULL DEFAULT now(),
    sua_luc         TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT vat_pham_noi_chua_hop_le
        CHECK (noi_chua IN ('TUI','MAC','RUONG','DUOI_DAT','KY_GUI','DA_HUY')),
    -- BB2: so luong khong bao gio am/khong; va mon co option thi luon la 1 cai (§1.3).
    CONSTRAINT vat_pham_so_luong_duong CHECK (so_luong > 0),
    CONSTRAINT vat_pham_co_option_thi_don_chiec
        CHECK (jsonb_array_length(chi_so_phu->'opt') = 0 OR so_luong = 1),
    -- Mon co chu thi phai o mot noi cua nguoi do; mon khong chu thi phai duoi dat/da huy.
    CONSTRAINT vat_pham_chu_va_noi_khop CHECK (
        (chu_so_huu_id IS NOT NULL AND noi_chua IN ('TUI','MAC','RUONG','KY_GUI'))
     OR (chu_so_huu_id IS NULL     AND noi_chua IN ('DUOI_DAT','DA_HUY'))
    )
);

-- BB1: mot o cua mot nguoi chi chua DUNG MOT mon. Rang buoc CHONG NHAN BAN manh nhat
-- ca KH7 — lam "hai mon cung o" KHONG THE ton tai. Partial index vi mon duoi dat/da
-- huy khong co o (o_so NULL).
CREATE UNIQUE INDEX IF NOT EXISTS vat_pham_mot_o_mot_mon
    ON vat_pham (chu_so_huu_id, noi_chua, o_so)
    WHERE chu_so_huu_id IS NOT NULL AND o_so IS NOT NULL;

-- Truy van nong nhat: "lay ca tui cua nguoi X" moi lan dang nhap / doi tui.
CREATE INDEX IF NOT EXISTS vat_pham_theo_chu
    ON vat_pham (chu_so_huu_id, noi_chua)
    WHERE chu_so_huu_id IS NOT NULL;

-- BB3: moi lan doi chu mot dong. CHI THEM, khong bao gio sua/xoa.
CREATE TABLE IF NOT EXISTS vat_pham_lich_su (
    id            BIGSERIAL   PRIMARY KEY,
    vat_pham_id   BIGINT      NOT NULL REFERENCES vat_pham(id),
    tu_ai         BIGINT      REFERENCES nhan_vat(id),   -- NULL = tu he thong sinh ra
    sang_ai       BIGINT      REFERENCES nhan_vat(id),   -- NULL = bi tieu huy
    ly_do         TEXT        NOT NULL,                  -- SINH_RA|MUA_NPC|BAN_NPC|NHAT|DUNG|BO|GIAO_DICH|ADMIN
    so_luong      INTEGER     NOT NULL,
    luc           TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS vat_pham_lich_su_theo_mon ON vat_pham_lich_su (vat_pham_id);

-- D3: chong lam hai lan cung mot hanh dong kinh te (idempotency key).
CREATE TABLE IF NOT EXISTS vat_pham_giao_dich (
    id          BIGSERIAL   PRIMARY KEY,
    khoa_idem   TEXT        NOT NULL UNIQUE,   -- "<nhanVatId>:<loai>:<soThuTu>"
    ket_qua     JSONB       NOT NULL,          -- de tra lai DUNG ket qua lan dau khi goi lap
    luc         TIMESTAMPTZ NOT NULL DEFAULT now()
);
