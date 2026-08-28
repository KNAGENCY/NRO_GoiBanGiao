-- SO NO PHAN THUONG — cho nhung khoan thuong khong the tra NGAY luc phat sinh
-- vi nguoi nhan khong con o day nua.
--
-- VI SAO CAN MOT BANG RIENG thay vi cong thang vao nhan_vat luc chia thuong:
-- actor khu vuc (internal/game/khuvuc.go) chi so huu con tro store.NhanVat cua
-- NHUNG NGUOI DANG O TRONG no. Voi nguoi da roi khu vuc / mat ket noi dung luc
-- boss chet, actor khong so huu gi ca — ghi vao con tro do la tranh chap du
-- lieu that (BAT-DAU-TU-DAY.md §6.9). Con duong vong "day mot ban chup
-- store.NhanVat cu qua ThoGhi" thi te hon nua: ThoGhi ghi DE CA DONG, nen ban
-- chup cu se KEO LUI moi thu nguoi do lam tu luc roi di (vi tri, mau, vang) —
-- mat du lieu that de doi lay mot khoan thuong.
--
-- Nen khoan thuong duoc GHI NO o day (qua mot hang doi ngoai actor, xem
-- internal/game/sonothuong.go) va tra vao LAN DANG NHAP SAU (internal/game/
-- login.go). Bang nay la thu duy nhat giua "cong bang" va "mat du lieu".
CREATE TABLE IF NOT EXISTS thuong_cho_nhan (
    id          BIGSERIAL   PRIMARY KEY,
    -- ON DELETE CASCADE: no cua mot nhan vat da bi xoa la rac, khong ai nhan
    -- duoc nua. Khong CASCADE thi moi lan xoa nhan vat deu vap khoa ngoai.
    nhan_vat_id BIGINT      NOT NULL REFERENCES nhan_vat(id) ON DELETE CASCADE,
    loai_thuong TEXT        NOT NULL,
    so_diem     BIGINT      NOT NULL,
    moc_tao     TIMESTAMPTZ NOT NULL DEFAULT now(),
    da_nhan     BOOLEAN     NOT NULL DEFAULT false,
    nhan_luc    TIMESTAMPTZ,

    -- Danh sach loai DONG: them loai moi phai sua CA o day VA o NhanNoThuong
    -- (nothuong.go) — co y bat nguoi sua phai di qua ca hai cho, vi mot loai
    -- duoc ghi vao ma khong ai biet cach TRA se nam lai vinh vien.
    CONSTRAINT thuong_cho_nhan_loai_hop_le
        CHECK (loai_thuong IN ('kinh_nghiem')),

    -- Chot chan cuoi cung chong mot khoan no AM hoac RONG: no am se TRU kinh
    -- nghiem cua nguoi choi luc nhan. ChiaThuongTheoDongGop (boss.go) da bo
    -- cac phan <= 0 roi, day la lop thu hai — dat o CSDL vi no chan duoc CA
    -- nhung duong ghi tuong lai chua ai viet.
    CONSTRAINT thuong_cho_nhan_so_diem_duong
        CHECK (so_diem > 0),

    -- Hai cot phai ke cung mot cau chuyen: da nhan thi phai co moc nhan, chua
    -- nhan thi khong duoc co. Sai lech giua chung lam moi bao cao "ai nhan gi
    -- luc nao" thanh vo nghia ma khong bao loi o dau ca.
    CONSTRAINT thuong_cho_nhan_nhan_luc_khop
        CHECK ((da_nhan AND nhan_luc IS NOT NULL) OR (NOT da_nhan AND nhan_luc IS NULL))
);

-- Index MOT PHAN (chi cac dong chua nhan). Duong doc duy nhat cua bang nay la
-- "no chua nhan cua nguoi X" chay o MOI lan dang nhap; con cac dong da nhan chi
-- nam do lam so sach. Index mot phan vi vay nho hon han index day du va khong
-- phinh theo lich su — bang nay chi tang, khong bao gio xoa dong.
CREATE INDEX IF NOT EXISTS thuong_cho_nhan_chua_nhan
    ON thuong_cho_nhan (nhan_vat_id) WHERE da_nhan = false;
