
package GiaCatLuongManager;

/**
 *
 * @author SADBOI
 */
import java.util.Random;
import java.util.regex.Pattern;
import player.Player;

public class Functions {
    private static final String REGEX = "\\b(dkm|đkm|đbrr|địt|đĩ|đỹ|cm|cmm|lồn|buồi|cc|ôm cl|mẹ mày|cặc|đụ|fuck|damn|clmm|dcmm|cl|tml|đ*t|c*c|dit|d*t|c.a.c|l.o.n|c.ặ.c|l.ồ.n|b.u.ồ.i|bu*i|đặc cầu|đồn lầu|bú cu|buscu|đm|cc|đb|db|lol|nhu lon|nhu cac|vc|vl|vãi|đéo|đờ mờ|đờ cờ mờ|clgt|dell|mẹ|cứt|shit|idiot|khốn|xiên chết|cụ|giao phối|thiểu năng|ngáo|chó|dog|đcmm|vcl|vkl|đ!t|d!t|đỵt|dyt|ngu|óc|.com|.net|.online|.vn|.pw|.pro|.org|.info|.ml|.ga|.gq|.cf|.fun|.xyz|.io|.club)\\b";
    private static Pattern pattern = Pattern.compile(
            "\\b(dkm|đkm|đbrr|địt|đĩ|đỹ|cm|cmm|lồn|buồi|cc|ôm cl|mẹ mày|cặc|đụ|fuck|damn|clmm|dcmm|cl|tml|đ*t|c*c|dit|d*t|c.a.c|l.o.n|c.ặ.c|l.ồ.n|b.u.ồ.i|bu*i|đặc cầu|đồn lầu|bú cu|buscu|đm|cc|đb|db|lol|nhu lon|nhu cac|vc|vl|vãi|đéo|đờ mờ|đờ cờ mờ|clgt|dell|mẹ|cứt|shit|idiot|khốn|xiên chết|cụ|giao phối|thiểu năng|ngáo|chó|dog|đcmm|vcl|vkl|đ!t|d!t|đỵt|dyt|ngu|óc|.com|.net|.online|.vn|.pw|.pro|.org|.info|.ml|.ga|.gq|.cf|.fun|.xyz|.io|.club)\\b",
            2);

    public static boolean checkspam(String text) {
        return pattern.matcher(text.toLowerCase()).find();
    }

    public static boolean checkspam(Player player, String text) {
        return pattern.matcher(text.toLowerCase()).find();
    }

    public static int maxint(long n) {
        return (int) (n > Integer.MAX_VALUE ? Integer.MAX_VALUE : n);
    }

    public static String generateRandomCharacters(int quantity) {
        StringBuilder sb = new StringBuilder();
        Random random = new Random();
        for (int i = 0; i < quantity; ++i) {
            int type = random.nextInt(2);
            char generatedChar = type == 0 ? (char) (random.nextInt(10) + 48) : (char) (random.nextInt(26) + 65);
            sb.append(generatedChar);
        }
        return sb.toString();
    }

    public static void sleep(long milis) {
        try {
            Thread.sleep(milis);
        } catch (InterruptedException interruptedException) {
            // empty catch block
        }
    }

    public static void sleep1Milis() {
        try {
            Thread.sleep(1L);
        } catch (InterruptedException interruptedException) {
            // empty catch block
        }
    }
}
