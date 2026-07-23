<?php
declare(strict_types=1);

/*
 * Lightweight application-layer DDoS guard for PHP/XAMPP behind Cloudflare.
 * This complements Cloudflare; it cannot protect an already-saturated network link.
 */

if (defined('DDOS_GUARD_LOADED')) {
    return;
}
define('DDOS_GUARD_LOADED', true);

const DDOS_GUARD_ENFORCE_CLOUDFLARE = false; // Enable only after the origin firewall allows Cloudflare IPs exclusively.
const DDOS_GUARD_MAX_URI_BYTES = 2048;
const DDOS_GUARD_MAX_QUERY_BYTES = 1024;
const DDOS_GUARD_MAX_QUERY_PARAMS = 24;
const DDOS_GUARD_MAX_HEADER_BYTES = 16384;
const DDOS_GUARD_MAX_HEADERS = 60;
const DDOS_GUARD_DATA_DIR = 'C:/xampp/tmp';
const DDOS_GUARD_SITE_ID = 'z18';
const DDOS_GUARD_COOKIE = 'ddg_z18_v1';
const DDOS_GUARD_COOKIE_TTL = 1800;

function ddos_guard_data_path(string $filename): string
{
    $base = is_dir(DDOS_GUARD_DATA_DIR) && is_writable(DDOS_GUARD_DATA_DIR)
        ? DDOS_GUARD_DATA_DIR
        : sys_get_temp_dir();
    return rtrim($base, '\\/') . DIRECTORY_SEPARATOR . $filename;
}

function ddos_guard_base64url_encode(string $value): string
{
    return rtrim(strtr(base64_encode($value), '+/', '-_'), '=');
}

function ddos_guard_base64url_decode(string $value): string|false
{
    $padding = (4 - strlen($value) % 4) % 4;
    return base64_decode(strtr($value . str_repeat('=', $padding), '-_', '+/'), true);
}

function ddos_guard_secret(): string
{
    static $secret = null;
    if (is_string($secret)) {
        return $secret;
    }
    $path = ddos_guard_data_path('ddos_guard_' . DDOS_GUARD_SITE_ID . '.secret');
    $loaded = @file_get_contents($path);
    if (is_string($loaded) && strlen($loaded) >= 32) {
        return $secret = $loaded;
    }
    $generated = random_bytes(32);
    @file_put_contents($path, $generated, LOCK_EX);
    $loaded = @file_get_contents($path);
    return $secret = (is_string($loaded) && strlen($loaded) >= 32) ? $loaded : $generated;
}

function ddos_guard_ip_prefix(string $ip): string
{
    $packed = @inet_pton($ip);
    if ($packed === false) {
        return 'invalid';
    }
    $bytes = strlen($packed) === 4 ? 3 : 7;
    return bin2hex(substr($packed, 0, $bytes));
}

function ddos_guard_browser_binding(string $clientIp): string
{
    $ua = substr((string)($_SERVER['HTTP_USER_AGENT'] ?? ''), 0, 512);
    return hash('sha256', DDOS_GUARD_SITE_ID . '|' . ddos_guard_ip_prefix($clientIp) . '|' . $ua);
}

function ddos_guard_make_token(string $clientIp, ?int $issuedAt = null): string
{
    $issuedAt ??= time();
    $payload = $issuedAt . '.' . ddos_guard_browser_binding($clientIp);
    $signature = hash_hmac('sha256', $payload, ddos_guard_secret(), true);
    return ddos_guard_base64url_encode($payload) . '.' . ddos_guard_base64url_encode($signature);
}

function ddos_guard_valid_token(string $clientIp): bool
{
    $token = (string)($_COOKIE[DDOS_GUARD_COOKIE] ?? '');
    if ($token === '' || strlen($token) > 256 || substr_count($token, '.') !== 1) {
        return false;
    }
    [$encodedPayload, $encodedSignature] = explode('.', $token, 2);
    $payload = ddos_guard_base64url_decode($encodedPayload);
    $signature = ddos_guard_base64url_decode($encodedSignature);
    if (!is_string($payload) || !is_string($signature) || substr_count($payload, '.') !== 1) {
        return false;
    }
    [$issuedText, $binding] = explode('.', $payload, 2);
    $issuedAt = ctype_digit($issuedText) ? (int)$issuedText : 0;
    if ($issuedAt < time() - DDOS_GUARD_COOKIE_TTL || $issuedAt > time() + 60) {
        return false;
    }
    if (!hash_equals(ddos_guard_browser_binding($clientIp), $binding)) {
        return false;
    }
    return hash_equals(hash_hmac('sha256', $payload, ddos_guard_secret(), true), $signature);
}

function ddos_guard_set_browser_cookie(string $clientIp): void
{
    setcookie(DDOS_GUARD_COOKIE, ddos_guard_make_token($clientIp), [
        'expires' => time() + DDOS_GUARD_COOKIE_TTL,
        'path' => '/',
        'secure' => (!empty($_SERVER['HTTPS']) && strtolower((string)$_SERVER['HTTPS']) !== 'off')
            || strtolower((string)($_SERVER['HTTP_X_FORWARDED_PROTO'] ?? '')) === 'https',
        'httponly' => true,
        'samesite' => 'Lax',
    ]);
}

function ddos_guard_browser_challenge(string $clientIp, string $peerIp, string $reason): never
{
    ddos_guard_log($reason, $clientIp, $peerIp);
    ddos_guard_set_browser_cookie($clientIp);
    http_response_code(503);
    header('Content-Type: text/html; charset=utf-8');
    header('Cache-Control: no-store');
    header('Retry-After: 1');
    header("Content-Security-Policy: default-src 'none'; script-src 'unsafe-inline'; style-src 'unsafe-inline'");
    echo '<!doctype html><meta charset="utf-8"><meta name="robots" content="noindex">'
        . '<title>Đang xác minh trình duyệt</title><style>body{font:16px system-ui;background:#111827;color:#fff;display:grid;place-items:center;height:100vh;margin:0}.b{text-align:center}</style>'
        . '<div class="b"><h1>Đang bảo vệ máy chủ…</h1><p>Trang sẽ tự tải lại.</p></div>'
        . '<script>setTimeout(function(){location.reload()},700)</script>';
    exit;
}

function ddos_guard_cloudflare_ranges(): array
{
    // Official ranges checked against https://www.cloudflare.com/ips/ on 2026-07-17.
    return [
        '103.21.244.0/22', '103.22.200.0/22', '103.31.4.0/22',
        '104.16.0.0/13', '104.24.0.0/14', '108.162.192.0/18',
        '131.0.72.0/22', '141.101.64.0/18', '162.158.0.0/15',
        '172.64.0.0/13', '173.245.48.0/20', '188.114.96.0/20',
        '190.93.240.0/20', '197.234.240.0/22', '198.41.128.0/17',
        '2400:cb00::/32', '2606:4700::/32', '2803:f800::/32',
        '2405:b500::/32', '2405:8100::/32', '2a06:98c0::/29',
        '2c0f:f248::/32',
    ];
}

function ddos_guard_ip_in_cidr(string $ip, string $cidr): bool
{
    [$network, $prefixText] = explode('/', $cidr, 2);
    $packedIp = @inet_pton($ip);
    $packedNetwork = @inet_pton($network);
    if ($packedIp === false || $packedNetwork === false || strlen($packedIp) !== strlen($packedNetwork)) {
        return false;
    }

    $prefix = (int)$prefixText;
    $maxBits = strlen($packedIp) * 8;
    if ($prefix < 0 || $prefix > $maxBits) {
        return false;
    }

    $wholeBytes = intdiv($prefix, 8);
    $remainingBits = $prefix % 8;
    if ($wholeBytes > 0 && substr($packedIp, 0, $wholeBytes) !== substr($packedNetwork, 0, $wholeBytes)) {
        return false;
    }
    if ($remainingBits === 0) {
        return true;
    }

    $mask = (0xff << (8 - $remainingBits)) & 0xff;
    return (ord($packedIp[$wholeBytes]) & $mask) === (ord($packedNetwork[$wholeBytes]) & $mask);
}

function ddos_guard_is_cloudflare_ip(string $ip): bool
{
    foreach (ddos_guard_cloudflare_ranges() as $cidr) {
        if (ddos_guard_ip_in_cidr($ip, $cidr)) {
            return true;
        }
    }
    return false;
}

function ddos_guard_client_context(): array
{
    $peerIp = trim((string)($_SERVER['REMOTE_ADDR'] ?? ''));
    $fromCloudflare = ddos_guard_is_cloudflare_ip($peerIp);
    $clientIp = $peerIp;

    if ($fromCloudflare) {
        $candidate = trim((string)($_SERVER['HTTP_CF_CONNECTING_IP'] ?? ''));
        if (filter_var($candidate, FILTER_VALIDATE_IP) !== false) {
            $clientIp = $candidate;
        }
    }

    if (filter_var($clientIp, FILTER_VALIDATE_IP) === false) {
        $clientIp = '0.0.0.0';
    }

    return [$clientIp, $peerIp, $fromCloudflare];
}

function ddos_guard_log(string $reason, string $clientIp, string $peerIp): void
{
    $path = ddos_guard_data_path('ddos_guard_z18_blocked.log');
    if (is_file($path) && (int)@filesize($path) > 5 * 1024 * 1024) {
        @rename($path, $path . '.1');
    }
    $line = json_encode([
        'time' => gmdate('c'),
        'reason' => $reason,
        'client_ip' => $clientIp,
        'peer_ip' => $peerIp,
        'method' => (string)($_SERVER['REQUEST_METHOD'] ?? ''),
        'uri' => substr((string)($_SERVER['REQUEST_URI'] ?? ''), 0, 500),
        'ray' => substr((string)($_SERVER['HTTP_CF_RAY'] ?? ''), 0, 100),
    ], JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
    if (is_string($line)) {
        @file_put_contents($path, $line . PHP_EOL, FILE_APPEND | LOCK_EX);
    }
}

function ddos_guard_reject(int $status, string $reason, string $clientIp, string $peerIp, int $retryAfter = 0): never
{
    ddos_guard_log($reason, $clientIp, $peerIp);
    http_response_code($status);
    header('Content-Type: application/json; charset=utf-8');
    header('Cache-Control: no-store');
    header('X-Content-Type-Options: nosniff');
    if ($retryAfter > 0) {
        header('Retry-After: ' . $retryAfter);
    }
    echo json_encode([
        'ok' => false,
        'error' => $status === 429 ? 'Quá nhiều yêu cầu. Vui lòng thử lại sau.' : 'Yêu cầu không hợp lệ.',
        'request_id' => substr((string)($_SERVER['HTTP_CF_RAY'] ?? bin2hex(random_bytes(8))), 0, 64),
    ], JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
    exit;
}

function ddos_guard_header_stats(): array
{
    $count = 0;
    $bytes = 0;
    foreach ($_SERVER as $key => $value) {
        if ($key === 'CONTENT_TYPE' || $key === 'CONTENT_LENGTH' || str_starts_with($key, 'HTTP_')) {
            $count++;
            $bytes += strlen((string)$key) + strlen((string)$value) + 4;
        }
    }
    return [$count, $bytes];
}

function ddos_guard_bucket(string $identity, string $scope, int $windowSeconds, int $limit): array
{
    $now = time();
    $windowStart = intdiv($now, $windowSeconds) * $windowSeconds;
    $key = hash('sha256', DDOS_GUARD_SITE_ID . '|' . $identity . '|' . $scope . '|' . $windowStart);
    $dbPath = ddos_guard_data_path('ddos_guard_z18.sqlite');

    try {
        $db = new PDO('sqlite:' . $dbPath, null, null, [
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_TIMEOUT => 1,
        ]);
        $db->exec('PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL; PRAGMA busy_timeout=150;');
        $db->exec('CREATE TABLE IF NOT EXISTS buckets (k TEXT PRIMARY KEY, count INTEGER NOT NULL, expires INTEGER NOT NULL) WITHOUT ROWID');
        $db->beginTransaction();
        $stmt = $db->prepare(
            'INSERT INTO buckets(k, count, expires) VALUES(:k, 1, :expires) '
            . 'ON CONFLICT(k) DO UPDATE SET count = count + 1 RETURNING count'
        );
        $stmt->execute([':k' => $key, ':expires' => $windowStart + ($windowSeconds * 2)]);
        $count = (int)$stmt->fetchColumn();
        $stmt->closeCursor();
        $db->commit();

        if (random_int(1, 500) === 1) {
            $cleanup = $db->prepare('DELETE FROM buckets WHERE expires < :now');
            $cleanup->execute([':now' => $now]);
        }
        return [$count <= $limit, max(1, ($windowStart + $windowSeconds) - $now), $count];
    } catch (Throwable $error) {
        // Fail open if the local counter is temporarily busy; Cloudflare remains the outer protection layer.
        return [true, 0, 0];
    }
}

[$clientIp, $peerIp, $fromCloudflare] = ddos_guard_client_context();
$hasBrowserProof = ddos_guard_valid_token($clientIp);

if (DDOS_GUARD_ENFORCE_CLOUDFLARE && !$fromCloudflare && !in_array($peerIp, ['127.0.0.1', '::1'], true)) {
    ddos_guard_reject(403, 'direct_origin_request', $clientIp, $peerIp);
}

// Make existing application logs/rate limits see the verified visitor IP, not a Cloudflare edge IP.
$_SERVER['DDOS_GUARD_EDGE_IP'] = $peerIp;
$_SERVER['DDOS_GUARD_FROM_CLOUDFLARE'] = $fromCloudflare ? '1' : '0';
$_SERVER['REMOTE_ADDR'] = $clientIp;

$method = strtoupper((string)($_SERVER['REQUEST_METHOD'] ?? 'GET'));
$uri = (string)($_SERVER['REQUEST_URI'] ?? '/');
$path = strtolower((string)(parse_url($uri, PHP_URL_PATH) ?: '/'));
$query = (string)($_SERVER['QUERY_STRING'] ?? '');

if (!in_array($method, ['GET', 'HEAD', 'POST', 'OPTIONS'], true)) {
    ddos_guard_reject(405, 'method_not_allowed', $clientIp, $peerIp);
}
if (strlen($uri) > DDOS_GUARD_MAX_URI_BYTES || strlen($query) > DDOS_GUARD_MAX_QUERY_BYTES) {
    ddos_guard_reject(414, 'request_uri_too_large', $clientIp, $peerIp);
}
if ($query !== '' && substr_count($query, '&') + 1 > DDOS_GUARD_MAX_QUERY_PARAMS) {
    ddos_guard_reject(400, 'too_many_query_parameters', $clientIp, $peerIp);
}
[$headerCount, $headerBytes] = ddos_guard_header_stats();
if ($headerCount > DDOS_GUARD_MAX_HEADERS || $headerBytes > DDOS_GUARD_MAX_HEADER_BYTES) {
    ddos_guard_reject(431, 'request_headers_too_large', $clientIp, $peerIp);
}

$contentLength = (int)($_SERVER['CONTENT_LENGTH'] ?? 0);
if ($contentLength < 0 || $contentLength > 2 * 1024 * 1024) {
    ddos_guard_reject(413, 'request_body_too_large', $clientIp, $peerIp);
}

// General burst control. HEAD/OPTIONS cost more because the observed flooders rotate methods.
$generalCost = in_array($method, ['HEAD', 'OPTIONS'], true) ? 3 : 1;
$generalCount = 0;
for ($i = 0; $i < $generalCost; $i++) {
    [$allowed, $retryAfter, $generalCount] = ddos_guard_bucket($clientIp, 'global', 10, 100);
    if (!$allowed) {
        ddos_guard_reject(429, 'global_rate_limit', $clientIp, $peerIp, $retryAfter);
    }
}

if (!$hasBrowserProof) {
    if ($method === 'POST' || $generalCount > 3 || ($query !== '' && strlen($query) > 64)) {
        ddos_guard_browser_challenge($clientIp, $peerIp, 'browser_proof_required');
    }
    if (in_array($method, ['GET', 'HEAD'], true)) {
        ddos_guard_set_browser_cookie($clientIp);
    }
}

$isAuth = preg_match('~/(dangnhap|dangky|quenmatkhau|doimatkhau)(?:\.php)?/?$~', $path) === 1;
$isAdmin = str_contains($path, '/admin/');
$isExpensive = preg_match('~/(boss_status|bxh0[12])(?:\.php)?/?$~', $path) === 1;

if ($isAuth) {
    $limit = $method === 'POST' ? 8 : 30;
    $window = $method === 'POST' ? 300 : 60;
    [$allowed, $retryAfter] = ddos_guard_bucket($clientIp, 'auth:' . $method, $window, $limit);
    if (!$allowed) {
        ddos_guard_reject(429, 'auth_rate_limit', $clientIp, $peerIp, $retryAfter);
    }
}
if ($isAdmin) {
    [$allowed, $retryAfter] = ddos_guard_bucket($clientIp, 'admin', 60, 90);
    if (!$allowed) {
        ddos_guard_reject(429, 'admin_rate_limit', $clientIp, $peerIp, $retryAfter);
    }
}
if ($isExpensive) {
    [$allowed, $retryAfter] = ddos_guard_bucket($clientIp, 'expensive', 60, 60);
    if (!$allowed) {
        ddos_guard_reject(429, 'expensive_rate_limit', $clientIp, $peerIp, $retryAfter);
    }
}

header('X-Content-Type-Options: nosniff');
header('Referrer-Policy: same-origin');
header('Permissions-Policy: camera=(), microphone=(), geolocation=()');
