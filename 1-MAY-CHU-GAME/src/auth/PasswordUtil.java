package auth;

import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.Base64;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;

/**
 * Mã hoá mật khẩu bằng PBKDF2-HMAC-SHA256.
 *
 * Định dạng lưu trong cột account.password:
 *     pbkdf2$&lt;số vòng&gt;$&lt;muối base64&gt;$&lt;băm base64&gt;
 *
 * PHẢI khớp chính xác với kame_password_hash()/kame_password_verify()
 * trong C:\xampp\htdocs\sv1\config.php. Sửa một bên mà quên bên kia
 * sẽ khiến toàn bộ tài khoản không đăng nhập được.
 *
 * Tài khoản cũ lưu mật khẩu dạng chữ thường vẫn đăng nhập được;
 * NDVSqlFetcher sẽ tự nâng cấp chúng sang dạng băm sau lần đăng nhập đúng.
 */
public final class PasswordUtil {

    private static final String PREFIX = "pbkdf2$";
    private static final int ITERATIONS = 100000;
    private static final int SALT_BYTES = 16;
    private static final int KEY_BITS = 256;
    private static final SecureRandom RNG = new SecureRandom();

    private PasswordUtil() {
    }

    /** Chuỗi đã ở dạng băm hay còn là mật khẩu thô của hệ thống cũ? */
    public static boolean isHashed(String stored) {
        return stored != null && stored.startsWith(PREFIX);
    }

    /** Băm mật khẩu mới, sinh muối ngẫu nhiên. */
    public static String hash(String password) {
        byte[] salt = new byte[SALT_BYTES];
        RNG.nextBytes(salt);
        byte[] dk = pbkdf2(password, salt, ITERATIONS, KEY_BITS);
        return PREFIX + ITERATIONS + "$"
                + Base64.getEncoder().encodeToString(salt) + "$"
                + Base64.getEncoder().encodeToString(dk);
    }

    /**
     * Đối chiếu mật khẩu người dùng nhập với giá trị lưu trong DB.
     * Tự nhận biết dạng băm hay dạng thô cũ.
     */
    public static boolean verify(String password, String stored) {
        if (password == null || stored == null || stored.isEmpty()) {
            return false;
        }
        if (!isHashed(stored)) {
            // Tài khoản cũ: so sánh thô nhưng vẫn dùng so sánh hằng thời gian
            return constantTimeEquals(
                    stored.getBytes(StandardCharsets.UTF_8),
                    password.getBytes(StandardCharsets.UTF_8));
        }
        String[] parts = stored.split("\\$");
        if (parts.length != 4) {
            return false;
        }
        try {
            int iterations = Integer.parseInt(parts[1]);
            byte[] salt = Base64.getDecoder().decode(parts[2]);
            byte[] want = Base64.getDecoder().decode(parts[3]);
            if (iterations < 1 || want.length == 0) {
                return false;
            }
            byte[] got = pbkdf2(password, salt, iterations, want.length * 8);
            return constantTimeEquals(want, got);
        } catch (Exception e) {
            return false;
        }
    }

    private static byte[] pbkdf2(String password, byte[] salt, int iterations, int keyBits) {
        PBEKeySpec spec = null;
        try {
            spec = new PBEKeySpec(password.toCharArray(), salt, iterations, keyBits);
            SecretKeyFactory f = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256");
            return f.generateSecret(spec).getEncoded();
        } catch (Exception e) {
            throw new IllegalStateException("Khong the bam mat khau", e);
        } finally {
            if (spec != null) {
                spec.clearPassword();
            }
        }
    }

    /** So sánh không phụ thuộc thời gian, tránh lộ thông tin qua độ trễ. */
    private static boolean constantTimeEquals(byte[] a, byte[] b) {
        if (a == null || b == null) {
            return false;
        }
        int diff = a.length ^ b.length;
        for (int i = 0; i < a.length && i < b.length; i++) {
            diff |= a[i] ^ b[i];
        }
        return diff == 0;
    }
}
