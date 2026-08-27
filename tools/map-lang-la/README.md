# Công cụ dựng map

Chạy bằng `python` (Windows). Đích tác động: `D:\NRO_Project\Assets\Resources\res\`.

> **Đọc [`SO-TAY-LAM-MAP.md`](../../docs/my-thuat/SO-TAY-LAM-MAP.md) trước** — quy trình 7 bước, quy tắc cỡ ảnh, và những bẫy đã vấp.

## Dùng thường xuyên nhất

| Lệnh | Làm gì |
|---|---|
| `do-map.py <ảnh...>` | **Thước đo map** — 7 chỉ số so với concept + **ảnh chẩn đoán chỉ đúng chỗ hỏng**, chạy 1 giây. Dùng nó trước, agent sau |
| `kiem-anh-dot-2.py` | Kiểm ảnh mới **trước khi lắp**: cỡ, viền bẩn, đáy nền, ghép đôi, chân chạm mép |
| `lap-anh-dot-2.py` | Lắp ảnh x4 vào client, tự thu nhỏ ra 4 mức zoom, chép sang bản `b1x` |
| [`tiled/`](tiled/README.md) | **Cách xếp bố cục hiện nay** — sửa map bằng chuột trong Tiled, `xem-truoc.py` vẽ ra ảnh trong 0,6 giây |
| `sinh-va-cham.py` | Sinh `map-1-collision.bin` **thẳng từ bố cục** + vá asset-lock |
| `vao-game-chup.py <tên>` | Dừng → nạp asset → chạy game → vào map → chụp 2 ảnh |

### `do-map.py` xuất thêm ảnh chẩn đoán

Bảng số nói **có** bệnh, ảnh chẩn đoán nói bệnh **nằm ở đâu**:

```bash
python tools/map-lang-la/do-map.py D:/NRO_Project/hoantacA.png
# -> tools/map-lang-la/chan-doan/hoantacA-chan-doan.png
```

Tám ô: ảnh gốc + bảy chỉ số, mỗi ô tô màu thẳng lên ảnh đúng chỗ làm hỏng chỉ số đó — vùng đất bị phẳng tô cam, cột trời lọt xuống dải nhà tô đỏ, hàng pixel có đường kẻ ngang vạch đỏ, khối núi tô lục. Mỗi ô kèm con số của ảnh này, của concept, và ngưỡng trong sổ tay.

Không cần ảnh thì thêm `--khong-anh`. Đổi chỗ để ảnh thì `--ra <thư mục>`.

### Đã nghỉ hưu

| Lệnh | Chuyện gì |
|---|---|
| `ve-bo-cuc-map.py` | Cách **cũ** xếp bố cục: 93 con số cứng trong code, chạy một cái là **xoá sạch 1248 ô rồi vẽ lại từ đầu**. Đã thay bằng `tiled/`. Chạy lên bây giờ nó **in cảnh báo rồi dừng**, không đụng vào map — muốn ép chạy phải gõ cờ `--toi-that-su-muon-ve-lai-tu-dau`. **Giữ lại để đọc**: nó là hồ sơ gốc giải thích vì sao từng công trình nằm ở đúng cột đó |

## Sửa ảnh

| Lệnh | Làm gì |
|---|---|
| `do-vien-tile.py` / `xoa-vien-tile.py` | Dò và xoá viền trắng/đen nằm trong ảnh tile |
| `lam-mem-cho-giap-nen.py` | Nhoà mép nối giữa 4 lớp nền cho hết đường kẻ ngang |
| `cat-troi-thua-lang-xa.py` | Cắt phần trời thừa ở đỉnh ảnh làng xa để nó che được chân núi |
| `lap-khe-ho-lop-gan.py` | Lấp khoảng hở giữa lớp nền thấp nhất và mặt cỏ |
| `tao-o-chuyen-tiep.py` | Trộn hai tile thành ô chuyển tiếp, ranh giới lởm chởm |
| `lam-dat-sau-co-chieu-sau.py` | Tạo tầng đất tối dần theo độ sâu |
| `don-vien-tim-sprite.py` | Dọn pixel magenta sót ở viền sprite nhân vật |

## Sửa dữ liệu game

| Lệnh | Làm gì |
|---|---|
| `them-tile-vao-tileset.py` | Khai báo tile mới là sàn / chặn ngang / trần trong `tile_set_info.bin` |
| `ghep-lai-nhan-vat.py` | Đặt độ lệch ba mảnh đầu/thân/chân nhân vật |
| `tao-chan-dung-nguoi-choi.py` | Cắt chân dung HUD từ ảnh đầu nhân vật |
| `lam-tuong-bien.py` | Dựng tường vô hình ở mép map |
| `nang-nhan-vat.py` | Nâng nhân vật cho chân chạm mặt đất |

> **Đổi `part.bin` thì phải tăng `vsData`** trong `server/internal/game/dongbo.go` — không tăng thì client cứ dùng bản cũ trong bộ đệm, sửa xong chụp lại thấy y hệt. Bẫy này đã vấp hai lần.
