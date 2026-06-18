$ErrorActionPreference = 'Stop'

$arch = if ([Environment]::Is64BitOperatingSystem) { 'amd64' } else {
    Write-Error "error: unsupported architecture (32-bit not supported)"
    exit 1
}

$bin     = "grimoire-windows-$arch.exe"
$url     = "https://github.com/jeffreytse/grimoire/releases/latest/download/$bin"
$destDir = "$env:USERPROFILE\bin"
$dest    = "$destDir\grimoire.exe"

if (-not (Test-Path $destDir)) {
    New-Item -ItemType Directory -Force -Path $destDir | Out-Null
}

Write-Host "Downloading grimoire (windows/$arch)..."
Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath -notlike "*$destDir*") {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$destDir", 'User')
    Write-Host "Added $destDir to user PATH (restart terminal to apply)"
}

Write-Host "Installed: $dest"
Write-Host ""
Write-Host "Next:"
Write-Host "  grimoire update    # clone skill library -> ~\.grimoire"
Write-Host "  grimoire install   # install to all detected AI agents"
