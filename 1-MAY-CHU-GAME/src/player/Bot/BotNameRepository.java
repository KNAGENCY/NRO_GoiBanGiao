package player.Bot;

import java.util.Set;
import utils.Util;

public final class BotNameRepository {

    private static final char[] NAME_SUFFIX_CHARS = "0123456789abcdefghijklmnopqrstuvwxyz".toCharArray();
    private static final String[] PLAYER_NAMES = {
            "0nekill", "0thua", "1082004", "1laze", "1tepheroin", "5sonline", "900njejd", "abcasd",
            "admin", "admin1", "adminn", "aevccc", "ahhjhj", "alibaba", "alookk", "amhbar",
            "anhba", "anhbadz", "anhbazz", "babyboy", "baongoc", "beboncute", "belieu", "bills7",
            "billsv7", "black", "blackgoku", "bluerkid", "bom0hp", "bom0s", "botuanday", "boydz",
            "bulma", "bunma", "cadic", "caiditcu", "calic", "carwyn", "caubedior", "cavexanh",
            "chattv", "chienthan", "chimcu", "chioccho", "chumsv", "clone", "clone1", "clonedzz",
            "cloneee", "cloneeeee", "damtutlon", "daxanh", "death", "dinhmy0511", "dinhmy123", "ditcum",
            "dopamin", "drama", "dsdxva", "emling", "emxinhsay2", "entruy", "finnjr", "gamekohay",
            "gamekyniem", "ggggg", "godnamec", "gogeta", "gohan", "gohansama", "gohanz", "green",
            "guku01", "haclong", "hades", "hai2k1", "hakai", "hangxinh", "hdjnjjj", "helios",
            "hemmetbo", "henlinh", "hermes", "heroin", "hhhhdhs", "hhhhh", "hieuthuhai", "hinata",
            "hoang", "honghainhi", "hongjin", "hotboy", "hotgirl", "huynhhcc", "idoltd", "iphone",
            "itachi", "iteam", "julius", "kakalot", "kame100m", "kame1hp", "kame1ty", "kame69",
            "kamejok0", "kamejoko", "kamekame", "kameop07", "kamesenco", "kamesenko", "kazuha", "kenyairing",
            "khanhsky", "kieuthuong", "kingdeath", "kingheaven", "kingkame", "kingko", "kingnamec", "kingpem",
            "kingwar", "kingz", "kovui", "krillin", "kuibap77", "l3igboy", "laning", "laze1b",
            "linhma", "lino0", "llbossll", "luxury", "lysv4", "madara", "manlaze", "manvip",
            "mavuong", "mayngu", "mazuong", "metuanday", "milooo", "mina1", "mkhang", "modeptrai",
            "n7m3k", "namec", "namecgod", "namectop1", "namek", "namephake", "namna", "nappa",
            "nappam", "ngminhung", "ngocrong", "ngulol", "nguyetxinh", "nhincacgi", "nhungne", "nhungxinh",
            "nhunhi", "noctraam", "nokia", "noname1", "novaroyal", "nvtvcl", "ocbuoi", "octiiu",
            "ojisann", "okking", "okvip", "onehit", "onekill", "onemill", "onetop", "osiin",
            "phong03", "phule", "picollo", "picolo", "piggy", "pikuhan", "piulazeee", "protein",
            "pstsuoper", "qtaii", "quanghuy", "quelam", "quockhanh", "quoctuan", "quylaokame", "radic",
            "radic1", "ronaldo", "rongthan", "s1sgk", "saitama", "sakura", "sayingan", "senvo",
            "sexybaby", "shadowss", "shakame", "shiba", "shoppe", "shushi", "sie7vatly", "silver",
            "silver1", "sogohan", "soicha", "sollh", "sollh1", "sollh2", "sollh3", "sollh4",
            "sollh5", "sollh9", "sollh99", "songoku", "spwwww", "ss7brensan", "ssskacaa", "taichuai",
            "tdvip", "thais", "theace", "thieugia", "thuonghn", "tlathangca", "top1pem", "topstkame",
            "traidat", "trgiangday", "troithue", "trumb0ss", "trumbem", "trumks", "trumlaze", "trumnamec",
            "trumnrsd", "trumoccho", "trumpem", "trums", "trumxayda", "truongan", "uchiha", "valentine",
            "vcuong", "vcuongg", "vegeta", "vegeto", "vigito", "violet", "vipxinh", "vjpxinh",
            "vnsun", "vnsun2", "vuasumu", "vutrunro", "xadagod", "xayda", "xaydagod", "xaydanew",
            "xdviolet", "xfdshjhj", "yamiie", "zlienhoan", "zombie"
    };

    private BotNameRepository() {
    }

    public static String randomTrainingBotName(Set<String> usedNames) {
        for (int attempt = 0; attempt < 200; attempt++) {
            String name = randomBaseName() + randomSuffixChar();
            if (usedNames == null || usedNames.add(name)) {
                return name;
            }
        }
        String fallback;
        do {
            fallback = randomBaseName() + randomSuffixChar() + randomSuffixChar();
        } while (usedNames != null && !usedNames.add(fallback));
        return fallback;
    }

    private static String randomBaseName() {
        return PLAYER_NAMES[Util.nextInt(0, PLAYER_NAMES.length - 1)];
    }

    private static char randomSuffixChar() {
        return NAME_SUFFIX_CHARS[Util.nextInt(0, NAME_SUFFIX_CHARS.length - 1)];
    }
}
