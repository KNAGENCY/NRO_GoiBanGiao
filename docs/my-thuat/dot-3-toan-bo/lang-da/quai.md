# QUÁI — Làng Đá (hệ Thổ) 🟡 chờ server chốt danh sách

> **Đọc `../_KHUON.md` trước.** File này chỉ ghi nội dung riêng của quái Làng Đá.
> 🟡 **CHỜ:** server chưa chốt chính xác Làng Đá có bao nhiêu quái / con nào.
> Bộ 8 con dưới là **đề xuất đại diện** (núi đá/hang động) — khi server chốt nội dung sẽ
> thêm/bớt. Agent **cứ vẽ 8 con này trước** (chắc chắn dùng được), đừng tự bịa thêm.

**Số ảnh: 8.** Kích thước: **96×96** (x4), tỉ lệ 1:1. Nền trong suốt.

---

## Phong cách (hệ Thổ — Làng Đá)

Bảng màu: **nâu đất · xám đá · be**, ánh sáng núi cao lạnh, bụi đá mờ. Quái sống ở núi
đá/hẻm núi/hang động. Nhìn nghiêng side-view, mặt xoay sang **PHẢI**,
toàn thân, có khối + đổ bóng (≥40 màu). **Tham chiếu chuẩn:** `../anh-tham-chieu/TOT-quai-soi-rung.png`
(lông/vảy có lớp, cơ có khối, nền cắt sạch). **Tránh:** `XAU-*.png` (phẳng, chân que).

**Subject block (dán kèm master prompt khi gen):**
`a rocky-mountain / cave monster, stone-textured hide, dust-worn, earth-brown and grey palette, side view facing right, full body, floating on neutral background.`

---

## Danh sách 8 quái

| Tên file | Cỡ | Mô tả chi tiết (màu + đặc trưng) |
|---|---|---|
| `quai_da_golem_x4.png` | 96×96 | **Golem đá**: khối đá nén hình người thô kệch, tay to nặng, khớp nứt lộ mạch khoáng, mắt phát sáng nâu-cam (chakra Thổ), rêu bám vai, bụi rơi. Chậm chạp mà uy lực. |
| `quai_da_nhen_da_x4.png` | 96×96 | **Nhện đá hang** tám cẳng khớp đá, thân phủ mai đá xám nứt, mắt cụm đỏ nhỏ, càng gai đá. Tư thế rình trong bóng tối. |
| `quai_da_doi_hang_x4.png` | 96×96 | **Dơi hang khổng lồ**: lông xám-nâu, cánh da căng gân, tai nhọn, răng nanh trắng, mắt đỏ, treo ngược hoặc dang cánh. Bụi hang mờ quanh. |
| `quai_da_gau_nui_x4.png` | 96×96 | **Gấu núi hung**: lông nâu sẫm dày, vai u cơ bắp, vuốt to, miệng nhe răng, sẹo trên mõm. Đứng chồm dữ tợn, bụi đá dưới chân. |
| `quai_da_de_nui_x4.png` | 96×96 | **Dê núi hoang** (sơn dương): lông xám nâu xù, cặp sừng cong dày, móng bám đá, mắt hằn, đầu cúi thủ thế húc. |
| `quai_da_than_lan_x4.png` | 96×96 | **Kỳ đà đá** (thằn lằn núi): vảy nâu-xám như đá, gai lưng chạy dọc sống, lưỡi thè, đuôi dài quật, bám vách đá. Bóng đổ dài. |
| `quai_da_ung_nui_x4.png` | 96×96 | **Ưng núi đá**: lông nâu-xám sạm, mỏ khoằm vàng sắc, móng vuốt quắp, cánh dang rộng lộ lông vũ cứng, đậu trên mỏm đá. |
| `quai_da_linh_gac_x4.png` | 96×96 | **Lính gác Đá** (ninja Iwa cấp thấp): áo choàng nâu đất trùm đầu, ốp phiến đá vai, băng trán khắc hai phiến đá, tay thủ kunai. Người, dáng thủ thế. |

---

## Tự kiểm trước khi giao (bắt buộc)

```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-ảnh-của-bạn> nguon
```
Dán nguyên kết quả vào báo cáo. **Tên mới** (`quai_da_*`) chưa có trong `kich-thuoc-anh.json`
— báo phía dự án thêm vào bảng (đều 96×96) trước khi công cụ tra được. Không tự đổi tên.

**Chưa cần vẽ đợt này:** biến thể tinh anh / boss của Làng Đá (xem `boss.md`), quái 3 làng còn lại.
