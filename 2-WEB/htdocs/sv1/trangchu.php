<?php
/* ==========================================================================
   NGỌC RỒNG KAME — TRANG CHỦ
   Chỉ gọi hàm của kame/layout.php và kame/data.php, không tự viết SQL.
   Số liệu hiển thị đều là số thật; đọc không được thì ghi rõ "Chưa có dữ liệu".
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';

/* --- Lấy dữ liệu một lần, dùng lại cho cả trang ------------------------- */

$user       = kame_current_user();
$status     = kame_server_status();
$allBosses  = kame_bosses();          // đã sắp xếp: ACTIVE trước, rồi maxHp giảm dần
$bossMeta   = kame_boss_meta();
$topPower   = kame_top_power(8);
$latestNews = kame_news(4);

/** Tổng số boss và số boss đang sống — đếm thật từ dữ liệu, không bịa. */
$bossTotal  = count($allBosses);
$bossAlive  = 0;
foreach ($allBosses as $b) {
    if ($b['status'] === 'ACTIVE') {
        $bossAlive++;
    }
}
$bossCards = array_slice($allBosses, 0, 6);

/* --- Vài hàm nhỏ chỉ dùng trong trang này ------------------------------- */

/** Số có dấu chấm ngăn nghìn theo cách viết Việt Nam. */
function tc_so(int $n): string
{
    return number_format($n, 0, ',', '.');
}

/** Phần trăm máu còn lại, chặn trong khoảng 0–100 (dữ liệu game có thể vượt maxHp). */
function tc_phan_tram(int $hp, int $maxHp): float
{
    if ($maxHp <= 0) {
        return 0.0;
    }
    return max(0.0, min(100.0, $hp / $maxHp * 100));
}

/**
 * Nhãn trạng thái của một boss.
 *
 * @return array{0:string,1:bool}  [chữ hiển thị, có phải đang nghỉ hay không]
 */
function tc_trang_thai(array $boss): array
{
    if (($boss['status'] ?? '') === 'ACTIVE') {
        return ['Đang sống', false];
    }

    $con = (int)($boss['respawnIn'] ?? 0);
    if ($con > 0) {
        // Server Java ghi bằng giây; nếu quá lớn thì đó là mili-giây.
        if ($con > 86400) {
            $con = (int)round($con / 1000);
        }
        return ['Hồi sinh sau ' . max(1, (int)ceil($con / 60)) . ' phút', true];
    }

    return ['Đang nghỉ', true];
}

/* --- Liên kết cộng đồng: lấy từ config.php, không hardcode lại ---------- */

$linkZalo     = (isset($SITE_ZALO_URL) && is_string($SITE_ZALO_URL) && $SITE_ZALO_URL !== '') ? $SITE_ZALO_URL : '/sv1/ho-tro';
$linkFacebook = (isset($SITE_FACEBOOK_URL) && is_string($SITE_FACEBOOK_URL) && $SITE_FACEBOOK_URL !== '') ? $SITE_FACEBOOK_URL : '/sv1/ho-tro';

kame_head(
    'Trang chủ',
    'Ngọc Rồng Kame — máy chủ Ngọc Rồng Online ba hành tinh Trái Đất, Namếc, Xayda. Xem boss đang sống, bảng xếp hạng và tải game.'
);
kame_header('trangchu');
?>

<style>
/* --- Chỉ dùng riêng cho trang chủ. Màu lấy từ biến của ui-kit.css. ------- */

/* Nhãn nhỏ phía trên tiêu đề hero */
.tc-hero-tag {
  display: inline-block;
  margin-bottom: var(--kame-space-4);
  padding: 6px 14px;
  border: 2px solid var(--kame-orange-600);
  border-radius: 999px;
  color: var(--kame-white);
  background: rgb(240 96 48 / 88%);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: .12em;
  text-transform: uppercase;
}

/* Bảy viên ngọc: nhãn "1 SAO"… nằm dưới mỗi viên */
.tc-ngoc { text-align: center; }
.tc-ngoc__label {
  margin-top: 10px;
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 15px;
  font-style: italic;
  font-weight: 900;
  letter-spacing: .06em;
  text-transform: uppercase;
}
.tc-ngoc .kame-dragonball {
  display: grid;
  place-items: center;
  padding: 0 8px;
  color: #a02000;
  font-size: 12px;
  line-height: 1;
  letter-spacing: -1px;
  overflow-wrap: anywhere;
}

/* Tin tức: ảnh bên trái, chữ bên phải */
.tc-tin {
  display: grid;
  grid-template-columns: 150px 1fr;
  gap: var(--kame-space-4);
  align-items: center;
  padding: var(--kame-space-4);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 92%);
  box-shadow: var(--kame-shadow-sm);
  color: inherit;
  text-decoration: none;
  transition: transform .18s ease, box-shadow .18s ease;
}
.tc-tin:hover { transform: translateY(-2px); box-shadow: var(--kame-shadow-md); }
.tc-tin__anh { overflow: hidden; border-radius: var(--kame-radius-sm); }
.tc-tin__anh img { display: block; width: 100%; aspect-ratio: 800/450; object-fit: cover; }
.tc-tin__ngay { font-size: 12px; font-weight: 800; color: rgb(40 21 15 / 58%); }
.tc-tin__ten {
  margin: 6px 0;
  color: var(--kame-maroon-700);
  font-size: 17px;
  font-weight: 800;
  line-height: 1.3;
}
.tc-tin__tom { margin: 0; font-size: 14px; color: rgb(40 21 15 / 70%); }

/* Ba thế giới */
.tc-tg { display: flex; height: 100%; flex-direction: column; }
.tc-tg__anh img { display: block; width: 100%; aspect-ratio: 800/500; object-fit: cover; }
.tc-tg__than { display: flex; flex: 1; flex-direction: column; gap: var(--kame-space-3); padding: var(--kame-space-5); }
.tc-tg__ky { margin: 0; padding-left: 18px; font-size: 14px; color: rgb(40 21 15 / 72%); }
.tc-tg p.tc-tg__ky { padding-left: 0; }
.tc-tg__ky li { margin-bottom: 4px; }
.tc-tg__nut { margin-top: auto; padding-top: var(--kame-space-2); }

/* Tải game */
.tc-tai { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: var(--kame-space-4); }
.tc-tai a {
  display: grid;
  gap: 4px;
  padding: var(--kame-space-5);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 92%);
  box-shadow: var(--kame-shadow-sm);
  color: var(--kame-maroon-700);
  text-align: center;
  text-decoration: none;
  transition: transform .18s ease, background-color .18s ease;
}
.tc-tai a:hover { transform: translateY(-2px); background: var(--kame-paper-300); }
.tc-tai strong {
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 22px;
  font-style: italic;
  font-weight: 900;
  text-transform: uppercase;
}
.tc-tai span { font-size: 13px; color: rgb(40 21 15 / 66%); }

/* Cộng đồng — chỉ tên nền tảng, không kèm con số */
.tc-cd { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: var(--kame-space-4); }
.tc-cd a {
  display: block;
  padding: var(--kame-space-4);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 92%);
  box-shadow: var(--kame-shadow-sm);
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 20px;
  font-style: italic;
  font-weight: 900;
  text-align: center;
  text-decoration: none;
  text-transform: uppercase;
}
.tc-cd a:hover { background: var(--kame-paper-300); }

/* Dòng trống khi chưa có dữ liệu */
.tc-trong {
  padding: var(--kame-space-6);
  border: 3px dashed rgb(96 0 24 / 40%);
  border-radius: var(--kame-radius-md);
  color: rgb(40 21 15 / 66%);
  text-align: center;
}

/* Nút "Xem tất cả boss" canh giữa */
.tc-giua { display: flex; justify-content: center; margin-top: var(--kame-space-5); }

@media (max-width: 900px) {
  .tc-tai { grid-template-columns: 1fr; }
  .tc-cd { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 600px) {
  .tc-tin { grid-template-columns: 1fr; }
  .tc-tin__anh img { aspect-ratio: 16/9; }
  .tc-ngoc__label { font-size: 12px; }
}
</style>

<!-- ===================== 1. HERO ===================== -->
<section class="kame-hero">
  <div class="kame-hero__inner">
    <div class="kame-container">
      <div class="kame-hero__content">
        <span class="tc-hero-tag">Máy chủ chính thức</span>
        <h1 class="kame-heading kame-heading--xl kame-hero__title">
          Tụ Ngọc<br><em>Khai Chiến</em>
        </h1>
        <p class="kame-hero__desc">
          Ba hành tinh, bảy viên ngọc và những trận săn boss không hồi kết — hành trình của bạn bắt đầu tại đây.
        </p>
        <div class="kame-hero__actions">
          <?php if ($user['logged']): ?>
            <a class="kame-button" href="/sv1/tai-game">Chơi ngay</a>
          <?php else: ?>
            <button type="button" class="kame-button" data-kame-open="dangky">Chơi ngay</button>
          <?php endif; ?>
          <a class="kame-button kame-button--secondary" href="/sv1/vu-tru">Khám phá</a>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- ===================== 2. SỐ LIỆU THẬT ===================== -->
<section class="kame-section kame-section--tight">
  <div class="kame-container">
    <div class="kame-grid">

      <div class="kame-col-4">
        <?php if ($status['ok']): ?>
          <div class="kame-stat kame-stat--live">
            <div class="kame-stat__value"><span class="kame-live-dot"></span><?= tc_so($status['online']) ?></div>
            <div class="kame-stat__label">Đang chơi</div>
          </div>
        <?php else: ?>
          <div class="kame-stat kame-stat--empty">
            <div class="kame-stat__value"><span class="kame-live-dot kame-live-dot--off"></span>Máy chủ đang bảo trì</div>
            <div class="kame-stat__label">Đang chơi</div>
          </div>
        <?php endif; ?>
      </div>

      <div class="kame-col-4">
        <?php if ($status['accounts'] > 0): ?>
          <div class="kame-stat">
            <div class="kame-stat__value"><?= tc_so($status['accounts']) ?></div>
            <div class="kame-stat__label">Tài khoản</div>
          </div>
        <?php else: ?>
          <div class="kame-stat kame-stat--empty">
            <div class="kame-stat__value">Chưa có dữ liệu</div>
            <div class="kame-stat__label">Tài khoản</div>
          </div>
        <?php endif; ?>
      </div>

      <div class="kame-col-4">
        <?php if ($bossTotal > 0): ?>
          <div class="kame-stat">
            <div class="kame-stat__value" data-tc-boss-alive><?= tc_so($bossAlive) ?></div>
            <div class="kame-stat__label">Boss đang sống</div>
          </div>
        <?php else: ?>
          <div class="kame-stat kame-stat--empty">
            <div class="kame-stat__value">Chưa có dữ liệu</div>
            <div class="kame-stat__label">Boss đang sống</div>
          </div>
        <?php endif; ?>
      </div>

    </div>
  </div>
</section>

<!-- ===================== 3. BẢY VIÊN NGỌC RỒNG ===================== -->
<section class="kame-section kame-section--tight">
  <div class="kame-container">
    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">Thu thập đủ bảy viên ngọc rồng</h2>
      <span class="kame-muted">Gom đủ bảy viên để gọi Rồng Thần và ước một điều</span>
    </div>
    <div class="kame-dragonballs">
      <?php for ($sao = 1; $sao <= 7; $sao++): ?>
        <div class="tc-ngoc">
          <div class="kame-dragonball" role="img" aria-label="Ngọc <?= $sao ?> sao"><?= str_repeat('★', $sao) ?></div>
          <div class="tc-ngoc__label"><?= $sao ?> Sao</div>
        </div>
      <?php endfor; ?>
    </div>
  </div>
</section>

<!-- ===================== 4. BOSS ĐANG SỐNG ===================== -->
<section class="kame-section">
  <div class="kame-container">

    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">
        <span class="kame-live-dot" style="background:var(--kame-red)"></span>Boss đang sống
      </h2>
      <?php if ($bossMeta['stale']): ?>
        <span class="kame-badge">Dữ liệu chưa cập nhật</span>
      <?php else: ?>
        <span class="kame-muted">Tự làm mới mỗi 10 giây</span>
      <?php endif; ?>
    </div>

    <?php if ($bossCards): ?>
      <div class="kame-boss-grid" data-tc-boss-grid>
        <?php foreach ($bossCards as $boss): ?>
          <?php
          [$nhan, $dangNghi] = tc_trang_thai($boss);
          $phanTram = tc_phan_tram($boss['hp'], $boss['maxHp']);
          ?>
          <article class="kame-boss-card<?= $dangNghi ? ' kame-boss-card--afk' : '' ?>">
            <div class="kame-boss-card__avatar">
              <img src="<?= kame_e(kame_icon_url($boss['avatar_id'])) ?>" alt="" loading="lazy" width="48" height="48">
            </div>
            <div>
              <div class="kame-boss-card__top">
                <h3 class="kame-boss-card__name"><?= kame_e($boss['name']) ?></h3>
                <span class="kame-badge<?= $dangNghi ? '' : ' kame-badge--hot' ?>"><?= kame_e($nhan) ?></span>
              </div>
              <div class="kame-boss-card__map"><?= kame_e($boss['mapName']) ?></div>
              <div class="kame-boss-hp">
                <div class="kame-boss-hp__fill" style="width:<?= number_format($phanTram, 2, '.', '') ?>%"></div>
              </div>
              <div class="kame-boss-hp__text"><?= tc_so($boss['hp']) ?> / <?= tc_so($boss['maxHp']) ?> HP</div>
            </div>
          </article>
        <?php endforeach; ?>
      </div>

      <div class="tc-giua">
        <a class="kame-button kame-button--secondary" href="/sv1/vu-tru">Xem tất cả <?= tc_so($bossTotal) ?> boss</a>
      </div>
    <?php else: ?>
      <div class="tc-trong">Chưa có dữ liệu boss. Máy chủ có thể đang bảo trì.</div>
    <?php endif; ?>

  </div>
</section>

<!-- ===================== 5. TIN TỨC MỚI NHẤT ===================== -->
<section class="kame-section">
  <div class="kame-container">

    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">Tin tức mới nhất</h2>
      <a class="kame-button kame-button--secondary kame-button--sm" href="/sv1/tin-tuc">Xem tất cả</a>
    </div>

    <?php if ($latestNews): ?>
      <div class="kame-grid">
        <?php foreach ($latestNews as $bai): ?>
          <div class="kame-col-6">
            <a class="tc-tin" href="/sv1/tin-tuc?bai=<?= kame_e(rawurlencode((string)$bai['slug'])) ?>">
              <div class="tc-tin__anh">
                <?php if (!empty($bai['image'])): ?>
                  <img src="<?= kame_e((string)$bai['image']) ?>" alt="" loading="lazy">
                <?php else: ?>
                  <?= kame_placeholder('tin-mac-dinh.jpg', '800/450') ?>
                <?php endif; ?>
              </div>
              <div>
                <?php $luc = strtotime((string)$bai['created_at']); ?>
                <div class="tc-tin__ngay"><?= $luc ? kame_e(date('d.m.Y', $luc)) : '' ?></div>
                <div class="tc-tin__ten"><?= kame_e((string)$bai['title']) ?></div>
                <p class="tc-tin__tom"><?= kame_e((string)$bai['excerpt']) ?></p>
              </div>
            </a>
          </div>
        <?php endforeach; ?>
      </div>
    <?php else: ?>
      <div class="tc-trong">Chưa có dữ liệu</div>
    <?php endif; ?>

  </div>
</section>

<!-- ===================== 6. BA THẾ GIỚI ===================== -->
<section class="kame-section">
  <div class="kame-container">

    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">Khám phá ba thế giới</h2>
      <span class="kame-muted">Mỗi hành tinh một lối chơi riêng</span>
    </div>

    <div class="kame-grid">

      <div class="kame-col-4">
        <article class="kame-card tc-tg">
          <div class="tc-tg__anh"><?= kame_placeholder('the-gioi-traidat.jpg', '800/500') ?></div>
          <div class="tc-tg__than">
            <h3 class="kame-heading" style="font-size:26px">Trái Đất</h3>
            <p class="tc-tg__ky">Lối chơi cân bằng nhất: sức đánh, máu và khí đều ở mức khá, dễ làm quen cho người mới.</p>
            <ul class="tc-tg__ky">
              <li><strong>Dịch chuyển tức thời</strong> — đi lại giữa các bản đồ nhanh hơn hẳn.</li>
              <li>Kỹ năng thiên về đánh thường và chưởng năng lượng ổn định.</li>
            </ul>
            <div class="tc-tg__nut"><a class="kame-button kame-button--sm" href="/sv1/vu-tru">Khám phá ngay</a></div>
          </div>
        </article>
      </div>

      <div class="kame-col-4">
        <article class="kame-card tc-tg">
          <div class="tc-tg__anh"><?= kame_placeholder('the-gioi-namec.jpg', '800/500') ?></div>
          <div class="tc-tg__than">
            <h3 class="kame-heading" style="font-size:26px">Namếc</h3>
            <p class="tc-tg__ky">Trâu và dai sức: hồi máu tốt nên đánh lâu không sợ hết trụ, hợp người thích cày dài hơi.</p>
            <ul class="tc-tg__ky">
              <li><strong>Trị thương</strong> — hồi máu cho mình và cho đồng đội.</li>
              <li><strong>Đẻ trứng</strong> — tạo quả trứng chịu đòn thay mình.</li>
            </ul>
            <div class="tc-tg__nut"><a class="kame-button kame-button--sm" href="/sv1/vu-tru">Khám phá ngay</a></div>
          </div>
        </article>
      </div>

      <div class="kame-col-4">
        <article class="kame-card tc-tg">
          <div class="tc-tg__anh"><?= kame_placeholder('the-gioi-xayda.jpg', '800/500') ?></div>
          <div class="tc-tg__than">
            <h3 class="kame-heading" style="font-size:26px">Xayda</h3>
            <p class="tc-tg__ky">Sát thương dồn cực mạnh trong thời gian ngắn, đổi lại máu mỏng hơn — đánh nhanh thắng nhanh.</p>
            <ul class="tc-tg__ky">
              <li><strong>Biến hình</strong> — tăng vọt sức đánh trong một khoảng thời gian.</li>
              <li><strong>Tự phát nổ</strong> — dồn toàn lực vào một đòn quyết định.</li>
            </ul>
            <div class="tc-tg__nut"><a class="kame-button kame-button--sm" href="/sv1/vu-tru">Khám phá ngay</a></div>
          </div>
        </article>
      </div>

    </div>
  </div>
</section>

<!-- ===================== 7. TOP NGƯỜI CHƠI ===================== -->
<section class="kame-section">
  <div class="kame-container">

    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">Top sức mạnh</h2>
      <a class="kame-button kame-button--secondary kame-button--sm" href="/sv1/bxh01">Xem bảng xếp hạng</a>
    </div>

    <?php if ($topPower): ?>
      <div class="kame-rank">
        <?php foreach ($topPower as $i => $nguoi): ?>
          <?php $hang = $i + 1; ?>
          <div class="kame-rank-row<?= $hang <= 3 ? ' kame-rank-row--top' : '' ?>">
            <div class="kame-rank-row__pos"><?= $hang ?></div>
            <div class="kame-rank-row__avatar">
              <img src="<?= kame_e(kame_icon_url($nguoi['avatar_id'])) ?>" alt="" loading="lazy" width="40" height="40">
            </div>
            <div>
              <div class="kame-rank-row__name"><?= kame_e($nguoi['name']) ?></div>
              <div class="kame-rank-row__meta"><?= kame_e($nguoi['planet']) ?></div>
            </div>
            <div class="kame-rank-row__value"><?= tc_so($nguoi['power']) ?></div>
          </div>
        <?php endforeach; ?>
      </div>
    <?php else: ?>
      <div class="tc-trong">Chưa có dữ liệu — hãy là người đầu tiên</div>
    <?php endif; ?>

  </div>
</section>

<!-- ===================== 8. TẢI GAME ===================== -->
<section class="kame-section">
  <div class="kame-container">

    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">Tải game</h2>
      <span class="kame-muted">Một tài khoản chơi được trên cả ba nền tảng</span>
    </div>

    <div class="tc-tai">
      <a href="/sv1/tai-game"><strong>Windows</strong><span>Bản cho máy tính</span></a>
      <a href="/sv1/tai-game"><strong>Android</strong><span>Bản cài đặt APK</span></a>
      <a href="/sv1/tai-game"><strong>iOS</strong><span>Bản cho iPhone, iPad</span></a>
    </div>

  </div>
</section>

<!-- ===================== 9. CỘNG ĐỒNG ===================== -->
<section class="kame-section">
  <div class="kame-container">

    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--md">Cộng đồng Ngọc Rồng Kame</h2>
      <span class="kame-muted">Hỏi đáp, nhận giftcode và theo dõi sự kiện</span>
    </div>

    <div class="tc-cd">
      <a href="<?= kame_e($linkZalo) ?>">Zalo</a>
      <a href="<?= kame_e($linkFacebook) ?>">Facebook</a>
      <a href="/sv1/ho-tro">YouTube</a>
      <a href="/sv1/ho-tro">TikTok</a>
    </div>

  </div>
</section>

<script>
/* --- Tự làm mới danh sách boss mỗi 10 giây -----------------------------
   Gọi lại đúng endpoint sẵn có /sv1/boss_status, không tạo endpoint mới.
   Sắp xếp và dựng thẻ y hệt phía PHP để hai bên nhìn giống nhau.
   ---------------------------------------------------------------------- */
(function () {
  "use strict";

  var luoi  = document.querySelector("[data-tc-boss-grid]");
  var oSong = document.querySelector("[data-tc-boss-alive]");
  if (!luoi) return;

  var SO_THE = <?= (int)count($bossCards) ?>;

  function esc(s) {
    return String(s === null || s === undefined ? "" : s)
      .replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;").replace(/'/g, "&#039;");
  }

  function so(n) {
    return String(Math.round(Number(n) || 0)).replace(/\B(?=(\d{3})+(?!\d))/g, ".");
  }

  function trangThai(b) {
    if (b.status === "ACTIVE") return { chu: "Đang sống", nghi: false };
    var con = Number(b.respawnIn) || 0;
    if (con > 0) {
      if (con > 86400) con = Math.round(con / 1000);
      return { chu: "Hồi sinh sau " + Math.max(1, Math.ceil(con / 60)) + " phút", nghi: true };
    }
    return { chu: "Đang nghỉ", nghi: true };
  }

  function dungThe(b) {
    var tt = trangThai(b);
    var maxHp = Math.max(1, Number(b.maxHp) || 1);
    var hp = Number(b.hp) || 0;
    var pt = Math.max(0, Math.min(100, hp / maxHp * 100));

    return '<article class="kame-boss-card' + (tt.nghi ? " kame-boss-card--afk" : "") + '">'
      + '<div class="kame-boss-card__avatar">'
      + '<img src="/icon_asset.php?file=' + encodeURIComponent(String(b.avatar_id || 1)) + '.png" alt="" loading="lazy" width="48" height="48">'
      + '</div><div>'
      + '<div class="kame-boss-card__top">'
      + '<h3 class="kame-boss-card__name">' + esc(b.name) + '</h3>'
      + '<span class="kame-badge' + (tt.nghi ? "" : " kame-badge--hot") + '">' + esc(tt.chu) + '</span>'
      + '</div>'
      + '<div class="kame-boss-card__map">' + esc(b.mapName) + '</div>'
      + '<div class="kame-boss-hp"><div class="kame-boss-hp__fill" style="width:' + pt.toFixed(2) + '%"></div></div>'
      + '<div class="kame-boss-hp__text">' + so(hp) + " / " + so(maxHp) + ' HP</div>'
      + '</div></article>';
  }

  function lamMoi() {
    fetch("/sv1/boss_status", {
      headers: { "Accept": "application/json" },
      credentials: "same-origin",
      cache: "no-store"
    })
      .then(function (r) { return r.ok ? r.json() : null; })
      .then(function (data) {
        if (!data || !Array.isArray(data.bosses)) return;

        var ds = data.bosses.slice();
        // Cùng thứ tự với kame_bosses(): đang sống lên trước, rồi máu tối đa giảm dần
        ds.sort(function (a, b) {
          var sa = a.status === "ACTIVE" ? 0 : 1;
          var sb = b.status === "ACTIVE" ? 0 : 1;
          if (sa !== sb) return sa - sb;
          return (Number(b.maxHp) || 0) - (Number(a.maxHp) || 0);
        });

        if (oSong) {
          var song = ds.filter(function (b) { return b.status === "ACTIVE"; }).length;
          oSong.textContent = so(song);
        }

        luoi.innerHTML = ds.slice(0, SO_THE).map(dungThe).join("");
      })
      .catch(function () {
        /* Mất mạng hoặc máy chủ tắt: giữ nguyên nội dung đang hiển thị. */
      });
  }

  setInterval(lamMoi, 10000);
})();
</script>

<?php
kame_auth_modals();
kame_footer();
