-- Bang tai khoan. Co tinh dat toi thieu o buoc nay: ke hoach 2 se mo rong
-- day du va chuyen du lieu cu sang.
CREATE TABLE IF NOT EXISTS tai_khoan (
    id            BIGSERIAL PRIMARY KEY,
    ten_dang_nhap TEXT        NOT NULL,
    mat_khau_bam  TEXT        NOT NULL,
    bi_khoa       BOOLEAN     NOT NULL DEFAULT FALSE,
    la_admin      BOOLEAN     NOT NULL DEFAULT FALSE,
    tao_luc       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Ten dang nhap khong phan biet hoa thuong: nguoi choi go "Abc" hay "abc"
-- deu phai vao dung mot tai khoan. Chi muc duy nhat chan viec tao trung.
CREATE UNIQUE INDEX IF NOT EXISTS tai_khoan_ten_duy_nhat
    ON tai_khoan (lower(ten_dang_nhap));

-- Ten dang nhap chi cho phep chu cai tieng Anh, chu so va gach duoi, dai 4-20.
--
-- Vi sao rang buoc o CSDL chu khong chi o tang tren: CSDL nay tao voi locale C,
-- ma duoi locale do ham lower() KHONG doi duoc ky tu ngoai bang chu cai tieng Anh.
-- Da kiem chung: 'Anguoi' va 'anguoi' viet co dau se cung lot vao, thanh hai tai
-- khoan rieng biet, pha vo quy tac "khong phan biet hoa thuong".
--
-- Gioi han nay khong phai luat moi: server Java cu da ap dung dung quy tac nay
-- (Service.java dong 374). Dua no xuong CSDL de web va trang quan tri chen thang
-- vao bang cung khong lach duoc.
--
-- Loi ich kem theo: chan luon kieu gia mao dung chu Kirin/Hy Lap trong nhu chu
-- Latin de dat ten giong het tai khoan quan tri.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tai_khoan_ten_hop_le'
    ) THEN
        ALTER TABLE tai_khoan
            ADD CONSTRAINT tai_khoan_ten_hop_le
            CHECK (ten_dang_nhap ~ '^[A-Za-z0-9_]{4,20}$');
    END IF;
END $$;
