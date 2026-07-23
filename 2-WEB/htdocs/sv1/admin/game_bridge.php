<?php
declare(strict_types=1);

function game_root(): string {
    static $root = null;
    if ($root !== null) {
        return $root;
    }

    $candidates = [];
    $webKey = web_server_key();
    $serverEnvName = 'GAME_ROOT_' . strtoupper($webKey);
    $envRoot = getenv($serverEnvName);
    if (is_string($envRoot) && trim($envRoot) !== '') {
        $candidates[] = trim($envRoot);
    }

    if (!empty($_SERVER[$serverEnvName])) {
        $candidates[] = (string)$_SERVER[$serverEnvName];
    }

    $envRoot = getenv('GAME_ROOT');
    if (is_string($envRoot) && trim($envRoot) !== '') {
        $candidates[] = trim($envRoot);
    }

    if (!empty($_SERVER['GAME_ROOT'])) {
        $candidates[] = (string)$_SERVER['GAME_ROOT'];
    }

    foreach (default_game_roots_by_web($webKey) as $defaultRoot) {
        $candidates[] = $defaultRoot;
    }

    $userProfile = getenv('USERPROFILE');
    if (is_string($userProfile) && trim($userProfile) !== '') {
        $candidates[] = rtrim($userProfile, "\\/") . DIRECTORY_SEPARATOR . 'Music' . DIRECTORY_SEPARATOR . $webKey;
    }

    $candidates[] = 'C:\\Users\\Administrator\\Music' . DIRECTORY_SEPARATOR . $webKey;
    $candidates[] = 'C:\\Users\\Admin\\Music' . DIRECTORY_SEPARATOR . $webKey;

    $existingCandidates = [];
    foreach (array_unique($candidates) as $candidate) {
        $candidate = rtrim($candidate, "\\/");
        if ($candidate !== '' && is_dir($candidate)) {
            $existingCandidates[] = $candidate;
            if (is_file($candidate . DIRECTORY_SEPARATOR . 'data' . DIRECTORY_SEPARATOR . 'config' . DIRECTORY_SEPARATOR . 'config.properties') && is_writable($candidate)) {
                $root = $candidate;
                return $root;
            }
        }
    }

    foreach ($existingCandidates as $candidate) {
        if (is_writable($candidate)) {
            $root = $candidate;
            return $root;
        }
    }

    foreach ($existingCandidates as $candidate) {
            $root = $candidate;
            return $root;
    }

    $root = rtrim((string)end($candidates), "\\/");
    return $root;
}

function web_server_key(): string {
    $webRoot = realpath(dirname(__DIR__));
    $key = $webRoot !== false ? basename($webRoot) : basename(dirname(__DIR__));
    $key = preg_replace('/[^A-Za-z0-9_-]/', '', (string)$key);
    return $key !== '' ? $key : 'sv1';
}

function default_game_roots_by_web(string $webKey): array {
    $roots = [
        'sv1' => [
            // Tren VPS game dat o C:\zz1 - duong dan ngan, khong nam trong ho so
            // nguoi dung nen khong hong khi dang nhap bang tai khoan khac.
            'C:\\zz1',
            'C:\\Users\\Admin\\Music\\zz1',
            'C:\\Users\\Administrator\\Music\\zz1',
        ],
    ];
    return $roots[$webKey] ?? [];
}

function game_file_path(string $fileName): string {
    return game_root() . DIRECTORY_SEPARATOR . ltrim($fileName, DIRECTORY_SEPARATOR);
}

function panel_command_file(): string {
    return game_file_path('panel_cmd.txt');
}

function panel_command_targets(): array {
    return [panel_command_file()];
}

function panel_exp_status_file(): string {
    return game_file_path('panel_exp_status.txt');
}

function panel_drop_status_file(): string {
    return game_file_path('panel_drop_status.txt');
}

function panel_admin_mode_status_file(): string {
    return game_file_path('panel_admin_mode_status.txt');
}

function panel_bot_status_file(): string {
    return game_file_path('panel_bot_status.txt');
}

// Tên file cũ của tính năng bot (thời chủ cũ đặt là "XNXX").
// Chỉ dùng để đọc lại trạng thái nếu file mới chưa được máy chủ game tạo ra.
function panel_bot_status_file_legacy(): string {
    return game_file_path('panel_xnxx_status.txt');
}

function panel_boss_persist_status_file(): string {
    return game_file_path('panel_boss_persist_status.txt');
}

function panel_boss_overrides_file(): string {
    return game_file_path('panel_boss_overrides.properties');
}

function panel_star_upgrade_config_file(): string {
    return game_file_path('panel_star_upgrade_rates.properties');
}

function panel_phucloi_config_targets(): array {
    $targets = [game_file_path('phucloi_config.properties')];
    // The Phúc Lợi feature is shared by the two local server sources.
    foreach (['C:\\Users\\Admin\\Music\\zz1', 'C:\\Users\\Admin\\Music\\j'] as $root) {
        if (is_dir($root)) {
            $targets[] = $root . DIRECTORY_SEPARATOR . 'phucloi_config.properties';
        }
    }
    return array_values(array_unique($targets));
}

function panel_read_phucloi_config(): array {
    foreach (panel_phucloi_config_targets() as $file) {
        if (!is_file($file)) continue;
        $out = [];
        foreach (file($file, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES) ?: [] as $line) {
            if ($line === '' || $line[0] === '#') continue;
            $pair = explode('=', $line, 2);
            if (count($pair) === 2) $out[trim($pair[0])] = trim($pair[1]);
        }
        return $out;
    }
    return [];
}

function panel_write_phucloi_config(array $rewards, ?bool $enabled = null): bool {
    if ($enabled === null) {
        $current = panel_read_phucloi_config();
        $enabled = filter_var($current['enabled'] ?? true, FILTER_VALIDATE_BOOLEAN);
    }
    $lines = ['# Phuc Loi rewards written by panel', 'enabled=' . ($enabled ? 'true' : 'false')];
    foreach ($rewards as $key => $value) {
        if (!preg_match('/^reward\.[0-4]\.[0-2]$/', (string)$key)) continue;
        $lines[] = $key . '=' . str_replace(["\r", "\n"], '', (string)$value);
    }
    $data = implode(PHP_EOL, $lines) . PHP_EOL;
    $ok = true;
    foreach (panel_phucloi_config_targets() as $target) {
        if (@file_put_contents($target, $data, LOCK_EX) === false) $ok = false;
    }
    return $ok;
}

function panel_boss_status_file(): string {
    return dirname(__DIR__) . DIRECTORY_SEPARATOR . 'boss_status.json';
}

function panel_server_status_file(): string {
    return game_file_path('panel_server_status.json');
}

function panel_server_status_candidates(): array {
    $root = game_root();
    return [
        panel_server_status_file(),
        $root . DIRECTORY_SEPARATOR . 'server_status.json',
        $root . DIRECTORY_SEPARATOR . 'panel_server_status.txt',
        $root . DIRECTORY_SEPARATOR . 'server_status.txt',
        $root . DIRECTORY_SEPARATOR . 'status.txt',
    ];
}

function panel_config_file(): string {
    return game_file_path('data' . DIRECTORY_SEPARATOR . 'config' . DIRECTORY_SEPARATOR . 'config.properties');
}

function panel_java_source_file(string $relativePath): string {
    return game_root() . DIRECTORY_SEPARATOR . 'src' . DIRECTORY_SEPARATOR . str_replace(['/', '\\'], DIRECTORY_SEPARATOR, $relativePath);
}

function panel_parse_java_switch_numbers(string $source, string $methodName): array {
    $values = [];
    $pattern = '/(?:private|public|protected)\s+static\s+(?:float|int|long|double)\s+' . preg_quote($methodName, '/') . '\s*\([^)]*\)\s*\{.*?return\s+switch\s*\([^)]*\)\s*\{(.*?)\};/s';
    if (!preg_match($pattern, $source, $match)) {
        return $values;
    }
    if (preg_match_all('/case\s+(\d+)\s*->\s*([0-9_]+(?:\.[0-9_]+)?)/', $match[1], $cases, PREG_SET_ORDER)) {
        foreach ($cases as $case) {
            $values[(int)$case[1]] = (float)str_replace('_', '', $case[2]);
        }
    }
    return $values;
}

function panel_star_upgrade_rates(): array {
    $maxStar = 8;
    $combineService = panel_java_source_file('models/Combine/CombineService.java');
    if (is_file($combineService)) {
        $source = (string)@file_get_contents($combineService);
        if (preg_match('/MAX_STAR_ITEM\s*=\s*(\d+)/', $source, $m)) {
            $maxStar = max(1, (int)$m[1]);
        }
    }

    $defaults = [];
    $phaLeSourceFile = panel_java_source_file('models/Combine/manifest/PhaLeHoaTrangBi.java');
    if (is_file($phaLeSourceFile)) {
        $source = (string)@file_get_contents($phaLeSourceFile);
        $ratioWeights = panel_parse_java_switch_numbers($source, 'getDefaultRatioWeight');
        $golds = panel_parse_java_switch_numbers($source, 'getDefaultGold');
        for ($star = 0; $star < $maxStar; $star++) {
            $defaults[$star + 1] = [
                'star' => $star + 1,
                'gold' => (int)($golds[$star] ?? 0),
                'ratio' => round((float)($ratioWeights[$star] ?? 0) * 100 / 3000, 6),
            ];
        }
    }

    for ($star = 1; $star <= $maxStar; $star++) {
        $defaults[$star] ??= ['star' => $star, 'gold' => 0, 'ratio' => 0.0];
    }

    $configFile = panel_star_upgrade_config_file();
    if (is_file($configFile)) {
        $props = [];
        foreach (file($configFile, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES) ?: [] as $line) {
            if (preg_match('/^\s*(star\.\d+\.(?:gold|ratio))\s*=\s*([^#;]*)/', $line, $match)) {
                $props[$match[1]] = trim($match[2]);
            }
        }
        for ($star = 1; $star <= $maxStar; $star++) {
            if (isset($props["star.$star.gold"])) {
                $defaults[$star]['gold'] = max(0, (int)str_replace(['.', ','], '', $props["star.$star.gold"]));
            }
            if (isset($props["star.$star.ratio"])) {
                $defaults[$star]['ratio'] = min(100, max(0, (float)str_replace(',', '.', $props["star.$star.ratio"])));
            }
        }
    }

    ksort($defaults);
    return array_values($defaults);
}

function panel_write_star_upgrade_rates(array $rates): bool {
    $lines = [
        '# star.N means upgrade to N Sao Pha Le',
        '# gold is vang, ratio is percent',
    ];
    foreach ($rates as $rate) {
        $star = (int)$rate['star'];
        $gold = max(0, (int)$rate['gold']);
        $ratio = max(0, (float)$rate['ratio']);
        $lines[] = "star.$star.gold=$gold";
        $lines[] = "star.$star.ratio=" . rtrim(rtrim(number_format($ratio, 6, '.', ''), '0'), '.');
    }
    return @file_put_contents(panel_star_upgrade_config_file(), implode(PHP_EOL, $lines) . PHP_EOL, LOCK_EX) !== false;
}

function panel_number_or_null($value): ?float {
    if ($value === null || $value === '') {
        return null;
    }
    if (is_string($value)) {
        $value = str_replace(',', '.', trim($value));
        if (preg_match('/-?\d+(?:\.\d+)?/', $value, $m)) {
            return (float)$m[0];
        }
        return null;
    }
    return is_numeric($value) ? (float)$value : null;
}

function panel_first_number(array $data, array $keys): ?float {
    foreach ($keys as $key) {
        if (array_key_exists($key, $data)) {
            $value = panel_number_or_null($data[$key]);
            if ($value !== null) {
                return $value;
            }
        }
    }
    return null;
}

function panel_read_server_status(): array {
    foreach (panel_server_status_candidates() as $file) {
        if (!is_file($file)) {
            continue;
        }
        $raw = @file_get_contents($file);
        if (!is_string($raw) || trim($raw) === '') {
            continue;
        }

        $data = json_decode($raw, true);
        if (!is_array($data)) {
            $data = [];
            foreach (preg_split('/\R+/', $raw) ?: [] as $line) {
                if (preg_match('/^\s*([A-Za-z0-9_.-]+)\s*[:=]\s*(.*?)\s*$/', $line, $m)) {
                    $data[strtolower($m[1])] = $m[2];
                }
            }
        }

        if (!$data) {
            continue;
        }

        $lower = [];
        foreach ($data as $key => $value) {
            $lower[strtolower((string)$key)] = $value;
        }

        return [
            'online' => panel_first_number($lower, ['online', 'players', 'player', 'playeronline', 'player_online', 'players_online']),
            'thread' => panel_first_number($lower, ['thread', 'threads', 'threadcount', 'thread_count']),
            'session' => panel_first_number($lower, ['session', 'sessions', 'sessioncount', 'session_count']),
            'cpu' => [
                panel_first_number($lower, ['cpu', 'cpu1', 'processcpu', 'process_cpu', 'processcpuload']),
                panel_first_number($lower, ['cpu2', 'systemcpu', 'system_cpu', 'systemcpuload']),
            ],
            'ram' => [
                panel_first_number($lower, ['ram', 'ram1', 'usedram', 'used_ram', 'memory', 'memoryused']),
                panel_first_number($lower, ['ram2', 'maxram', 'max_ram', 'totalram', 'total_ram', 'memorymax']),
            ],
            'updated_at' => (int)(panel_first_number($lower, ['generatedat', 'generated_at', 'updatedat', 'updated_at', 'time', 'timestamp']) ?? filemtime($file)),
            'source' => basename($file),
        ];
    }

    return [
        'online' => null,
        'thread' => null,
        'session' => null,
        'cpu' => [null, null],
        'ram' => [null, null],
        'updated_at' => null,
        'source' => null,
    ];
}

function panel_public_server_status(bool $isAdmin): array {
    $status = panel_read_server_status();
    if (!$isAdmin && $status['online'] !== null && (int)$status['online'] >= 10) {
        if ($status['online'] !== null) $status['online'] += 45;
        if ($status['thread'] !== null) $status['thread'] += 100;
        if ($status['session'] !== null) $status['session'] += 45;
    }
    $status['is_admin_view'] = $isAdmin;
    return $status;
}

function append_panel_command(string $command): bool {
    $written = false;
    foreach (panel_command_targets() as $commandFile) {
        $dir = dirname($commandFile);
        if (!is_dir($dir)) {
            @mkdir($dir, 0777, true);
        }
        if (!is_dir($dir) || !is_writable($dir)) {
            continue;
        }
        if (file_put_contents($commandFile, $command . PHP_EOL, FILE_APPEND | LOCK_EX) !== false) {
            $written = true;
        }
    }
    return $written;
}

