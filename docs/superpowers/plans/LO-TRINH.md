# Lộ trình các kế hoạch thi công

**Thiết kế gốc:** [`../specs/2026-07-23-naruto-tren-loi-nro-design.md`](../specs/2026-07-23-naruto-tren-loi-nro-design.md)

Mỗi kế hoạch là một tài liệu riêng, viết theo skill `superpowers:writing-plans`,
thi công theo skill `superpowers:subagent-driven-development`.

**Mức độ chắc chắn của lộ trình này:** Kế hoạch 1 đã xong. Kế hoạch 2, 3, 9 đã
chốt phạm vi. Các kế hoạch còn lại là bản phác — thứ tự và ranh giới sẽ chỉnh
khi tới nơi, vì làm xong cái trước mới biết cái sau cần gì.

---

## Bảng lộ trình

| # | Tên | Trạng thái | Phụ thuộc |
|---|---|---|---|
| 1 | Lõi giao thức + Đăng nhập | ✅ **Xong** (gộp `6529d35`) | — |
| 2 | **Nhân vật: tạo, lưu, tải** | Chưa bắt đầu | 1 |
| 3 | Mô hình actor cho khu vực + di chuyển | Chưa bắt đầu | 1, 2 |
| 4 | Chiến đấu, chiêu thức, quái, boss | Chưa bắt đầu | 3 |
| 5 | Vật phẩm, hành trang, cửa hàng, giao dịch | Chưa bắt đầu | 3 |
| 6 | Kho nội dung YAML + nạp nóng | Chưa bắt đầu | 4, 5 |
| 7 | Nhiệm vụ, cốt truyện, NPC | Chưa bắt đầu | 6 |
| 8 | Web + Panel quản trị (Next.js) | Chưa bắt đầu | 2 |
| 9 | **Cảm giác chơi** — hoạt ảnh, chuyển màn, hướng dẫn người mới | Chưa bắt đầu | 3, 4 |
| B | Bảo mật: chống gian lận, chống dò mật khẩu, chống DDoS | Chưa bắt đầu | 3, 4 |
| V | Triển khai VPS, giám sát, phát hành | Chưa bắt đầu | tất cả |

Kế hoạch **B** và **V** không đánh số vì chúng cắt ngang mọi kế hoạch khác,
không nằm gọn ở một chỗ trong thứ tự.

> **Ghi chú Kế hoạch 2 (24/07/2026):** ban đầu định là "chuyển dữ liệu cũ".
> Nhưng CSDL cũ chỉ có **3 tài khoản + 3 nhân vật thử nghiệm** (đếm thật từ
> `3-CO-SO-DU-LIEU/z-moi-nhat.sql`) — không có người chơi thật. Quy trình
> chuyển dữ liệu ở thiết kế §6.3 sinh ra để bảo vệ tài sản người chơi thật,
> ở đây không cần. Đây là game Naruto mới, khởi đầu sạch. Nên Kế hoạch 2 đổi
> thành **xây hệ thống nhân vật** (tạo/lưu/tải), nối tiếp chỗ Kế hoạch 1 dừng
> (đăng nhập xong → màn hình tạo nhân vật). Ba tài khoản thử vi phạm ràng buộc
> tên đăng nhập mới → bỏ, không mang sang.

---

## Kế hoạch 9 — Cảm giác chơi

### Tại sao tách riêng

Ba thứ dưới đây quyết định game "đã tay" hay "ức chế", nhưng không thuộc về
bất kỳ kế hoạch nào ở trên. Nếu không tách ra, mỗi kế hoạch sẽ làm một kiểu
rồi cuối cùng game chắp vá.

### Phạm vi

**A. Hoạt ảnh**
- Hệ thống hiệu ứng: định dạng nhị phân đã giải mã từ `EffectData.cs`, nên
  thêm chiêu mới là việc **thêm dữ liệu**, không phải sửa code client
- Hoạt ảnh nhân vật: đứng, chạy, đánh, trúng đòn, gục, hồi sinh
- Hoạt ảnh chiêu thức: khởi chiêu → bay → nổ, mỗi giai đoạn khớp với mốc
  thời gian máy chủ gửi xuống
- Phản hồi khi trúng đòn: giật hình, nhấp nháy, số sát thương bay lên

**B. Chuyển màn hình**
- Vào/ra khu vực: hiện đang cắt phụt, cần làm mờ dần
- Mở/đóng bảng (hành trang, cửa hàng, kỹ năng)
- Màn hình chờ khi tải: phải có thứ để nhìn, không được đứng hình
- Mất kết nối và nối lại: người chơi phải **hiểu chuyện gì đang xảy ra**,
  không phải đứng nhìn màn hình đơ như lỗi "Xin chờ" của bản Java cũ

**C. Hướng dẫn người mới (onboarding)**
- Mười phút đầu quyết định người chơi ở lại hay bỏ đi
- Tạo nhân vật → chọn làng → nhiệm vụ đầu tiên, không được có bước nào
  bắt đọc chữ mà không có gì để bấm
- Dạy bằng cách cho làm, không phải bằng bảng chữ dài

### ⚠️ Phần phải quyết SỚM, không đợi tới Kế hoạch 9

Đây là lý do thật sự khiến Kế hoạch 9 nguy hiểm nếu để cuối. Ba con số dưới
đây là **ràng buộc của máy chủ**, phải chốt và nhét vào **Kế hoạch 3**. Chốt
sai thì tới Kế hoạch 9 mới phát hiện, lúc đó sửa là viết lại tầng khu vực.

| Thứ phải chốt | Vì sao không đợi được |
|---|---|
| **Nhịp máy chủ** (bao nhiêu lần cập nhật mỗi giây) | Nhịp thấp thì di chuyển giật, không hoạt ảnh nào cứu được. Nhịp cao thì tốn băng thông và CPU. |
| **Tần suất gửi vị trí + ngân sách nội suy phía client** | Client phải "đoán" vị trí giữa hai lần máy chủ báo. Không thiết kế từ đầu thì nhân vật sẽ nhảy cóc. |
| **Client dự đoán trước, máy chủ xác nhận sau** | Người chơi bấm đánh phải thấy phản hồi **ngay khung hình đó**, không đợi máy chủ trả lời. Nhưng máy chủ vẫn phải là bên quyết định thật (xem §5.1 của thiết kế: lỗ hổng "máy chủ tin tuyệt đối client"). Hai yêu cầu này mâu thuẫn nhau — phải thiết kế cách hoà giải **trước khi** viết tầng khu vực, không phải vá sau. |

**✅ Đã giải quyết** — xem [`../specs/2026-07-23-hoa-giai-muot-va-chong-gian-lan.md`](../specs/2026-07-23-hoa-giai-muot-va-chong-gian-lan.md).
Bảng số chốt ở §4.3 của tài liệu đó, đưa nguyên vào Global Constraints của
Kế hoạch 3. Ba sửa đổi client ở §4.1 làm cùng đợt dựng lại client.

### Chưa viết chi tiết vì

Kế hoạch 9 đụng vào client Unity và phụ thuộc máy chủ gửi xuống những gì —
mà điều đó do Kế hoạch 3 và 4 quyết định. Viết 50 bước chi tiết bây giờ là
viết ra thứ sẽ phải bỏ đi. Sẽ viết đầy đủ khi Kế hoạch 4 xong.
