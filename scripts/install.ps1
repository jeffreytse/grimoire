$ErrorActionPreference = 'Stop'

# ── Helpers ───────────────────────────────────────────────────────────────────
function Write-Info    { Write-Host "  " -NoNewline; Write-Host "→" -ForegroundColor Cyan    -NoNewline; Write-Host "  $args" }
function Write-Success { Write-Host "  " -NoNewline; Write-Host "✓" -ForegroundColor Green   -NoNewline; Write-Host "  $args" }
function Write-Warn    { Write-Host "  " -NoNewline; Write-Host "!" -ForegroundColor Yellow  -NoNewline; Write-Host "  $args" }
function Write-Fail    { Write-Host "  " -NoNewline; Write-Host "✗" -ForegroundColor Red     -NoNewline; Write-Host "  $args"; exit 1 }

# OSC 8 hyperlink (clickable in Windows Terminal / modern consoles)
function OSC8($url, $text) { return "`e]8;;$url`a$text`e]8;;`a" }

# ── Banner ────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host " " -NoNewline
Write-Host "✦" -ForegroundColor Yellow -NoNewline
Write-Host "▗▄▄▄▄▄▄▄▄▖" -ForegroundColor DarkYellow -NoNewline
Write-Host "✦" -ForegroundColor Yellow -NoNewline
Write-Host "  " -NoNewline
Write-Host "grimoire" -ForegroundColor White -NoNewline
Write-Host "  installer" -ForegroundColor DarkGray

Write-Host "  " -NoNewline
Write-Host "▐" -ForegroundColor DarkYellow -NoNewline
Write-Host "▬▬▬│▬▬▬▬" -ForegroundColor DarkGreen -NoNewline
Write-Host "▌" -ForegroundColor DarkYellow -NoNewline
Write-Host "   The world's best practices for AI assistants" -ForegroundColor DarkGray

Write-Host "  " -NoNewline
Write-Host "▐" -ForegroundColor DarkYellow -NoNewline
Write-Host "▬▬ ✦ ▬▬▬" -ForegroundColor DarkGreen -NoNewline
Write-Host "▌" -ForegroundColor DarkYellow -NoNewline
Write-Host "   https://github.com/jeffreytse/grimoire" -ForegroundColor Cyan

Write-Host "  " -NoNewline
Write-Host "▐" -ForegroundColor DarkYellow -NoNewline
Write-Host "▬▬▬│▬▬▬▬" -ForegroundColor DarkGreen -NoNewline
Write-Host "▌" -ForegroundColor DarkYellow

Write-Host " " -NoNewline
Write-Host "✦" -ForegroundColor Yellow -NoNewline
Write-Host "▝▀▀▀▀▀▀▀▀▘" -ForegroundColor DarkYellow -NoNewline
Write-Host "✦" -ForegroundColor Yellow -NoNewline
Write-Host "  ⭐ $(OSC8 'https://github.com/jeffreytse/grimoire' 'Star')  💖 $(OSC8 'https://github.com/sponsors/jeffreytse' 'Sponsor')  🐛 $(OSC8 'https://github.com/jeffreytse/grimoire/issues' 'Issues')"

Write-Host ""

# ── Platform detection ────────────────────────────────────────────────────────
if (-not [Environment]::Is64BitOperatingSystem) {
    Write-Fail "32-bit Windows is not supported"
}
$arch = 'amd64'
Write-Info "Platform: windows/$arch"

# ── Install destination ───────────────────────────────────────────────────────
$destDir = if ($env:GRIMOIRE_INSTALL_DIR) { $env:GRIMOIRE_INSTALL_DIR } else { "$env:USERPROFILE\bin" }
$dest    = "$destDir\grimoire.exe"

if (-not (Test-Path $destDir)) {
    New-Item -ItemType Directory -Force -Path $destDir | Out-Null
}

# Existing install — show current version
if (Get-Command grimoire -ErrorAction SilentlyContinue) {
    try {
        $current = & grimoire --version 2>$null
        Write-Info "Replacing existing install ($current)"
    } catch {
        Write-Info "Replacing existing install"
    }
}

# ── Download ──────────────────────────────────────────────────────────────────
$bin = "grimoire-windows-$arch.exe"
$url = "https://github.com/jeffreytse/grimoire/releases/latest/download/$bin"
$tmp = [System.IO.Path]::GetTempFileName() + ".exe"

Write-Info "Downloading $bin"
try {
    $progressPref = $ProgressPreference
    $ProgressPreference = 'SilentlyContinue'  # hide the ugly default progress bar
    Invoke-WebRequest -Uri $url -OutFile $tmp -UseBasicParsing
    $ProgressPreference = $progressPref
} catch {
    Write-Fail "Download failed: $_`n     Check https://github.com/jeffreytse/grimoire/releases"
}

# ── Install ───────────────────────────────────────────────────────────────────
Move-Item -Force $tmp $dest

$version = ""
try { $version = (& $dest --version 2>$null) -replace '^grimoire version ', '' } catch {}
$versionLabel = if ($version) { " ($version)" } else { "" }

Write-Success "Installed grimoire$versionLabel  →  $dest"

# PATH — add to user PATH if missing
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath -notlike "*$destDir*") {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$destDir", 'User')
    Write-Warn "Added $destDir to user PATH — restart your terminal to apply"
} else {
    Write-Info "$destDir is already in your PATH"
}

# ── Next steps ────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  Get started:" -ForegroundColor White
Write-Host ""
Write-Host "    " -NoNewline; Write-Host "grimoire wizard" -ForegroundColor Cyan -NoNewline; Write-Host "     # interactive setup" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  Or manually:" -ForegroundColor DarkGray
Write-Host ""
Write-Host "    " -NoNewline; Write-Host "grimoire update" -ForegroundColor Cyan -NoNewline;  Write-Host "     # fetch the official skill library" -ForegroundColor DarkGray
Write-Host "    " -NoNewline; Write-Host "grimoire install" -ForegroundColor Cyan -NoNewline; Write-Host "    # install skills to all detected AI agents" -ForegroundColor DarkGray
Write-Host "    " -NoNewline; Write-Host "grimoire check" -ForegroundColor Cyan -NoNewline;   Write-Host "      # run a compliance check on your project" -ForegroundColor DarkGray
Write-Host ""
