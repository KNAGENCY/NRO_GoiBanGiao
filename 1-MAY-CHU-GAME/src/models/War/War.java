package models.War;

import utils.TimeUtil;

public class War {
    private static War i;

    public static final int MAP_PREPARE_ID = 182;
    public static final int MAP_WAR_ID = 168;

    public static final int HOUR_OPEN = 21;
    public static final int MIN_OPEN = 29;

    public static final int HOUR_CLOSE_ENTER = 21;
    public static final int MIN_CLOSE_ENTER = 50;

    public static final int HOUR_START_WAR = 21;
    public static final int MIN_START_WAR = 44;

    public static final int HOUR_END = 21;
    public static final int MIN_END = 50;

    private long TIME_OPEN;
    private long TIME_CLOSE_ENTER;
    private long TIME_START_WAR;
    private long TIME_END;

    private boolean warStarted = false;

    /**
     * Ngày (dd/MM/yyyy) mà bộ mốc giờ hiện tại đang được tính cho.
     *
     * LÝ DO PHẢI CÓ BIẾN NÀY — đây là lỗi khiến Map War chỉ chạy đúng một ngày:
     *
     * resetTime() tính TIME_OPEN/TIME_END theo ngày HÔM NAY, và trước đây nó chỉ
     * chạy đúng một lần duy nhất, lúc đối tượng War được tạo ra (mà đối tượng này
     * lại được tạo lười - mãi tới khi có người chơi đầu tiên nói chuyện với Quy Lão
     * Kame hoặc ToriBot). Từ đó về sau bộ mốc giờ đứng yên.
     *
     * Sang ngày hôm sau, "bây giờ" luôn lớn hơn TIME_END của hôm qua, nên
     * isOpenTime() vĩnh viễn trả về false → không ai vào được Map War nữa,
     * cho tới khi khởi động lại máy chủ.
     *
     * Hàm update() lẽ ra lo việc tính lại này, nhưng nó KHÔNG ĐƯỢC GỌI Ở ĐÂU CẢ
     * (đã rà toàn bộ mã nguồn: chỉ có isOpenTime và isStartWarTime được dùng).
     * Vì vậy cách sửa ở đây là tự chữa ngay trong các hàm kiểm tra, thay vì đi
     * sửa vòng lặp chính của máy chủ - ít rủi ro hơn nhiều.
     */
    private String ngayDangTinh = "";

    public static synchronized War gI() {
        if (i == null) {
            i = new War();
        }
        return i;
    }

    private War() {
        resetTime();
    }

    /**
     * Nếu đã sang ngày mới thì tính lại toàn bộ mốc giờ cho ngày đó.
     * Được gọi ở đầu mọi hàm kiểm tra, nên Map War tự mở lại mỗi ngày
     * mà không cần khởi động lại máy chủ.
     */
    private synchronized void baoDamDungNgay() {
        String homNay = TimeUtil.getTimeNow("dd/MM/yyyy");
        if (!homNay.equals(ngayDangTinh)) {
            resetTime();
            warStarted = false;
        }
    }

    public synchronized void resetTime() {
        try {
            String today = TimeUtil.getTimeNow("dd/MM/yyyy");
            ngayDangTinh = today;
            TIME_OPEN = TimeUtil.getTime(today + " " + HOUR_OPEN + ":" + MIN_OPEN + ":00", "dd/MM/yyyy HH:mm:ss");
            TIME_CLOSE_ENTER = TimeUtil.getTime(today + " " + HOUR_CLOSE_ENTER + ":" + MIN_CLOSE_ENTER + ":00",
                    "dd/MM/yyyy HH:mm:ss");
            TIME_START_WAR = TimeUtil.getTime(today + " " + HOUR_START_WAR + ":" + MIN_START_WAR + ":00",
                    "dd/MM/yyyy HH:mm:ss");
            TIME_END = TimeUtil.getTime(today + " " + HOUR_END + ":" + MIN_END + ":00", "dd/MM/yyyy HH:mm:ss");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * Hiện KHÔNG được vòng lặp chính gọi. Giữ lại và sửa cho đúng, phòng khi
     * sau này có người nối nó vào GameLoop. Không còn gọi resetTime() ở đây nữa
     * vì việc đó đã do baoDamDungNgay() lo, và gọi ở đây là sai: lúc kết thúc
     * trận thì vẫn đang trong ngày cũ, tính lại cũng ra đúng bộ mốc giờ cũ.
     */
    public synchronized void update() {
        baoDamDungNgay();
        long now = System.currentTimeMillis();
        if (now >= TIME_START_WAR && now < TIME_END && !warStarted) {
            warStarted = true;
        }
        if (now >= TIME_END && warStarted) {
            warStarted = false;
        }
    }

    public synchronized boolean isOpenTime() {
        baoDamDungNgay();
        long now = System.currentTimeMillis();
        return now >= TIME_OPEN && now < TIME_CLOSE_ENTER;
    }

    public synchronized boolean isStartWarTime() {
        baoDamDungNgay();
        long now = System.currentTimeMillis();
        return now >= TIME_START_WAR && now < TIME_END;
    }

    public synchronized boolean isEndTime() {
        baoDamDungNgay();
        long now = System.currentTimeMillis();
        return now >= TIME_END;
    }

    public synchronized boolean isWarStarted() {
        baoDamDungNgay();
        return warStarted;
    }
}
