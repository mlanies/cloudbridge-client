export default {
  async fetch(request) {
    const script = `
Write-Host ""
Write-Host ""
Write-Host "╔════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║          2GC INSTALLER             ║" -ForegroundColor Cyan
Write-Host "║      ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌           ║" -ForegroundColor DarkGray
Write-Host "║ Ваша безопасность — наш приоритет  ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""
Write-Host "    Безопасность — это не роскошь," -ForegroundColor Yellow
Write-Host "            это стандарт." -ForegroundColor Yellow
Write-Host ""
Write-Host " © 2GC, 2025 | https://2gc.ru"
Write-Host "" 

Write-Host ""
if ($ExecutionContext.SessionState.LanguageMode.value__ -ne 0) {
    Write-Host "PowerShell is not running in Full Language Mode. Please use стандартный PowerShell." -ForegroundColor Red
    exit 1
}

Write-Host "Выберите, что хотите зарегистрировать:"
Write-Host " 1 - Cloudflare tunnel (Zero Trust)"
Write-Host " 2 - CloudBridge Client"
$choice = Read-Host "Введите цифру (1/2)"

$token = Read-Host "Введите ваш 2GC Token"
if ([string]::IsNullOrWhiteSpace($token)) {
    Write-Host "Токен не введён. Выход."
    exit 1
}

function Install-Cloudflared {
    Write-Host "\\n[Cloudflared] Проверка существующей установки..."
    $existingService = Get-Service -Name "Cloudflared" -ErrorAction SilentlyContinue
    if ($existingService) {
        Write-Host "[Cloudflared] Уже установлен." -ForegroundColor Yellow
        $uninst = Read-Host "Удалить старый сервис? (y/n)"
        if ($uninst -eq 'y') {
            & "C:\\Program Files (x86)\\cloudflared\\cloudflared.exe" service uninstall
            Start-Sleep -Seconds 2
        } else {
            Write-Host "Отмена установки Cloudflared."
            return
        }
    }
    $msiUrl = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.msi"
    $msiPath = "$env:TEMP\\cloudflared_$([guid]::NewGuid().Guid).msi"
    Write-Host "[Cloudflared] Скачивание MSI..."
    try {
        Invoke-WebRequest -Uri $msiUrl -OutFile $msiPath -UseBasicParsing
    } catch {
        Write-Host "Ошибка скачивания MSI." -ForegroundColor Red
        return
    }
    if (-not (Test-Path $msiPath)) {
        Write-Host "MSI не найден после скачивания." -ForegroundColor Red
        return
    }
    Write-Host "[Cloudflared] Установка..."
    Start-Process msiexec.exe -ArgumentList "/i \`"$msiPath\`" /qn" -Wait
    $cloudflaredPath = "C:\\Program Files (x86)\\cloudflared\\cloudflared.exe"
    if (-not (Test-Path $cloudflaredPath)) {
        Write-Host "cloudflared.exe не найден!" -ForegroundColor Red
        return
    }
    Write-Host "[Cloudflared] Регистрация токена..."
    & "$cloudflaredPath" service install $token
    Remove-Item $msiPath -Force
    Write-Host "[Cloudflared] Установка и регистрация завершены!" -ForegroundColor Green
}

function Install-CloudBridgeClient {
    Write-Host "\\n[CloudBridge Client] Проверка существующей установки..."
    $clientPath = "C:\\Program Files\\CloudBridgeClient\\cloudbridge-client.exe"
    if (Test-Path $clientPath) {
        Write-Host "[CloudBridge Client] Уже установлен." -ForegroundColor Yellow
        $uninst = Read-Host "Удалить старую версию? (y/n)"
        if ($uninst -eq 'y') {
            & $clientPath service uninstall
            Start-Sleep -Seconds 2
        } else {
            Write-Host "Отмена установки CloudBridge Client."
            return
        }
    }
    $url = "https://github.com/mlanies/cloudbridge-client/releases/latest/download/cloudbridge-client-windows-amd64.exe"
    $dst = "$env:TEMP\\cloudbridge-client.exe"
    Write-Host "[CloudBridge Client] Скачивание дистрибутива..."
    try {
        Invoke-WebRequest -Uri $url -OutFile $dst -UseBasicParsing
    } catch {
        Write-Host "Ошибка скачивания CloudBridge Client." -ForegroundColor Red
        return
    }
    if (-not (Test-Path $dst)) {
        Write-Host "Файл не найден после скачивания." -ForegroundColor Red
        return
    }
    Write-Host "[CloudBridge Client] Установка и регистрация..."
    Start-Process -FilePath $dst -ArgumentList "service install $token" -Wait
    Remove-Item $dst -Force
    Write-Host "[CloudBridge Client] Установка и регистрация завершены!" -ForegroundColor Green
}

switch ($choice) {
    '1' { Install-Cloudflared }
    '2' { Install-CloudBridgeClient }
    default { Write-Host "Некорректный выбор. Выход." -ForegroundColor Red }
}

Write-Host "\\nГотово!"
pause
    `.trim();

    return new Response(script, {
      headers: { "Content-Type": "text/plain; charset=utf-8" }
    });
  }
}; 