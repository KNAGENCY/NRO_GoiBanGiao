-- Bang nhan vat. Theo nguyen tac §6.2 cua thiet ke: CAC COT RIENG BIET, tra cuu
-- duoc — KHONG gop nhieu thu vao mot chuoi JSON nhu bang `player` cu (bang cu
-- co ~40 cot text JSON giai ma bang vi tri mang, rat gion, doi thu tu ghi la
-- hong het du lieu). Cot mo rong (tui do, ky nang, nhiem vu) se la bang rieng o
-- cac ke hoach sau, khong nhet vao day.
CREATE TABLE IF NOT EXISTS nhan_vat (
    id              BIGSERIAL PRIMARY KEY,
    tai_khoan_id    BIGINT      NOT NULL REFERENCES tai_khoan(id),
    ten             TEXT        NOT NULL,
    lang            SMALLINT    NOT NULL,   -- id lang/chung toc (KHONG phai gioi tinh)
    kieu_toc        SMALLINT    NOT NULL,   -- id sprite toc
    -- vi tri trong the gioi
    map_id          INTEGER     NOT NULL,
    vi_tri_x        INTEGER     NOT NULL,
    vi_tri_y        INTEGER     NOT NULL,
    -- chi so cot loi
    cap             INTEGER     NOT NULL DEFAULT 1,
    kinh_nghiem     BIGINT      NOT NULL DEFAULT 0,
    suc_manh        BIGINT      NOT NULL DEFAULT 0,   -- nhan luc
    hp_toi_da       INTEGER     NOT NULL,
    hp_hien_tai     INTEGER     NOT NULL,
    chakra_toi_da   INTEGER     NOT NULL,
    chakra_hien_tai INTEGER     NOT NULL,
    -- tien
    vang            BIGINT      NOT NULL DEFAULT 0,
    ngoc            BIGINT      NOT NULL DEFAULT 0,
    -- thoi gian
    tao_luc         TIMESTAMPTZ NOT NULL DEFAULT now(),
    dang_nhap_cuoi  TIMESTAMPTZ
);

-- Ten nhan vat khong phan biet hoa thuong: khong cho hai nguoi ten "Naruto" va
-- "naruto" cung ton tai (chong gia mao). Cung ly do dung lower() nhu ten dang
-- nhap o Ke hoach 1 — va vi charset chi la [A-Za-z0-9] (kiem o tang Go) nen
-- lower() duoi locale C hoat dong dung, khong dinh loi dau tieng Viet.
CREATE UNIQUE INDEX IF NOT EXISTS nhan_vat_ten_duy_nhat
    ON nhan_vat (lower(ten));

-- Mot tai khoan chi mot nhan vat: dung voi mo hinh server cu (client khong co
-- man chon nhan vat that su). Rang buoc o CSDL de khong the tao hai nhan vat
-- cho cung tai khoan du co dua tranh (race) o tang tren.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'nhan_vat_mot_tai_khoan'
    ) THEN
        ALTER TABLE nhan_vat
            ADD CONSTRAINT nhan_vat_mot_tai_khoan UNIQUE (tai_khoan_id);
    END IF;
END $$;

-- Tra nhan vat theo tai khoan la truy van nong nhat (moi lan dang nhap). Chi
-- muc nay khong that su can vi da co UNIQUE(tai_khoan_id) o tren tu dong tao
-- chi muc — giu lai lam tai lieu y dinh, IF NOT EXISTS nen khong trung.
CREATE INDEX IF NOT EXISTS nhan_vat_theo_tai_khoan
    ON nhan_vat (tai_khoan_id);

-- Ten nhan vat chi cho phep chu cai tieng Anh va chu so, dai 5-10. Cung ly do
-- nhu rang buoc ten dang nhap o 0001: duoi locale C, lower() KHONG case-fold
-- duoc ky tu ngoai ASCII, nen neu de lot ky tu co dau vao cot co
-- UNIQUE(lower(ten)) thi co che chong trung ten khong phan biet hoa thuong bi
-- vo hieu. CHECK o CSDL de web/panel/script import chen thang cung khong lach
-- duoc. Khop regex reTenNhanVat phia Go (5-10 ky tu [A-Za-z0-9], Task 6).
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'nhan_vat_ten_hop_le'
    ) THEN
        ALTER TABLE nhan_vat
            ADD CONSTRAINT nhan_vat_ten_hop_le
            CHECK (ten ~ '^[A-Za-z0-9]{5,10}$');
    END IF;
END $$;
