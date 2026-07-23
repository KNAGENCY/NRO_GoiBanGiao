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
