# Cầu nối Tiled — để làm map **dần dần** thay vì vẽ lại từ đầu mỗi lần

## Vì sao có thư mục này

Trước đây bố cục map sống trong `ve-bo-cuc-map.py`: **93 con số cứng trong 167 dòng Python**. Muốn dịch một căn nhà sang trái 2 ô thì sửa một con số rồi chạy script — mà script **xoá sạch cả 1248 ô rồi vẽ lại từ đầu**.

Đó không phải là *sửa*. Đó là *sinh lại*. Và vì mỗi lần sinh lại đều đổi vài thứ ngoài ý muốn, kết quả đúng là cảm giác "chắp vá, xếp hình".

Ba thứ thiếu, giờ đã có:

| Thiếu gì | Giờ có gì |
|---|---|
| Bố cục nằm trong code, sửa = sinh lại toàn bộ | Bố cục nằm trong **`map-1-lang-la.tmx`**, sửa từng ô một |
| Không nhìn thấy cái mình đang làm | **`xem-truoc.py`** vẽ ra ảnh trong **0,6 giây** |
| Concept một nơi, bố cục một nơi, không gì buộc chúng khớp | Concept là **lớp giấy can nằm dưới** ngay trong file `.tmx` |

---

## Cần cài gì

**[Tiled](https://www.mapeditor.org/)** — miễn phí, mã nguồn mở, chạy Windows. *Chưa cài trên máy này.*

Không cài cũng dùng được `xem-truoc.py`, nhưng để **sửa bằng chuột** thì cần Tiled.

---

## Bốn lệnh

```bash
python tools/map-lang-la/tiled/tao-bo-tile.py       # 1 lần: sinh bộ tile cho Tiled
python tools/map-lang-la/tiled/map-sang-tiled.py    # map hiện có -> .tmx
# ... mở map-1-lang-la.tmx bằng Tiled, kéo thả, lưu ...
python tools/map-lang-la/tiled/tiled-sang-map.py    # .tmx -> map game
python tools/map-lang-la/sinh-va-cham.py            # va chạm suy lại từ bố cục
```

Rồi **tăng `vsData`** bên client. Quên bước này thì client dùng bản cache cũ và **trong game không đổi gì cả** — đã mắc bẫy này hai lần.

Xem trước bất cứ lúc nào, không cần Unity:

```bash
python tools/map-lang-la/tiled/xem-truoc.py --khung
```

---

## Quy ước neo — chỗ dễ sai nhất

Vật thể to (nhà, cây, cổng) là **một tile id** với ảnh cao nhiều ô. Câu hỏi: ô lưới ghi số đó nằm ở đâu của cái ảnh?

- **Client game**: góc **TRÊN**-trái. Nhà cao 6 ô đặt ở hàng 12 → chiếm hàng 12…17.
- **Tiled**: góc **DƯỚI**-trái. Cùng số đó ở hàng 12 → chiếm hàng 7…12.

Lệch đúng `chiều_cao_ô − 1` hàng. Cầu nối tự dịch qua lại, chi tiết ở `chung.py`.

> **Đã tự mắc bẫy này.** Bản `xem-truoc.py` đầu tiên vẽ theo quy ước dưới-trái, nhìn thấy **16 vật thể đều treo cách mặt đất 1–7 ô**, suýt báo là map hỏng. Đo lại bằng quy ước trên-trái thì **16/16 chạm đất chính xác**. Map vẫn luôn đúng — công cụ mới sai.

---

## Hai điều công cụ này **không** làm được

**1. Thứ tự vẽ khi vật thể chồng nhau có thể khác game — đúng 10 ô trên map hiện tại (0,065% pixel).**

Client vẽ theo thứ tự hàng trong map; Tiled vẽ theo thứ tự hàng của nó. Sau khi dịch toạ độ, hai thứ tự này không còn trùng nhau ở chỗ có vật thể chồng lấn. Ví dụ thật: cột đèn ở cột 7 và căn nhà ở cột 8 chồng nhau — **game vẽ đèn đè lên nhà, Tiled vẽ nhà đè lên đèn**.

Nên khi cần chắc chắn, dùng bản vẽ theo đúng quy ước client:

```bash
python tools/map-lang-la/tiled/xem-truoc.py --game
```

**2. Chỉ có một lớp.** Client chỉ đọc một lưới tile. Không thêm được lớp bóng đổ hay lớp tiền cảnh — muốn vậy phải sửa cả client lẫn server.

---

## Vì sao tin được cầu nối này

Ba phép kiểm, chạy lại được bất cứ lúc nào:

| Kiểm gì | Kết quả |
|---|---|
| `map → .tmx → map` có khớp **từng byte** không | **Khớp từng byte** (1250/1250) |
| Vẽ từ `.tmx` và vẽ từ map game có ra cùng ảnh không | Giống, trừ **10 ô** chồng lấn nói ở trên |
| Đổi mà chồng lấn/tràn map | Báo lỗi và **dừng**, không ghi đè im lặng |

`tiled-sang-map.py` luôn sao lưu ra `1.bytes.bak-tiled` trước khi ghi đè.
