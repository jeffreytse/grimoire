#Requires -Version 5.1
[CmdletBinding()]
param(
    [string]$Domain    = "",
    [string]$Subdomain = "",
    [string]$Skill     = "",
    [string]$Target    = "",
    [switch]$Uninstall,
    [switch]$Copy,
    [switch]$Upgrade,
    [switch]$Clean,
    [switch]$Yes,
    [switch]$List,
    [switch]$Doctor,
    [switch]$Version,
    [switch]$Help
)

$GrimoireHome = if ($env:GRIMOIRE_HOME) { $env:GRIMOIRE_HOME } else { Join-Path $HOME ".grimoire" }
$GrimoireRepo = "https://github.com/jeffreytse/grimoire.git"
$script:Link  = -not $Copy.IsPresent

$_localRepoRoot = Split-Path -Parent $PSScriptRoot
if (Test-Path (Join-Path $_localRepoRoot ".git")) {
    $RepoRoot = $_localRepoRoot
} else {
    if (Test-Path (Join-Path $GrimoireHome ".git")) {
        Write-Host "Updating grimoire at $GrimoireHome..."
        git -C $GrimoireHome pull --quiet
    } else {
        Write-Host "Cloning grimoire to $GrimoireHome..."
        git clone --depth 1 --quiet $GrimoireRepo $GrimoireHome
    }
    $RepoRoot = $GrimoireHome
}
$SkillsRoot      = Join-Path $RepoRoot "skills"
$ClaudeSkillsDir   = Join-Path $HOME ".claude\skills"
$AgentsSkillsDir   = Join-Path $HOME ".agents\skills"
$GeminiSkillsDir   = Join-Path $HOME ".gemini\skills"
$OpenClawSkillsDir = Join-Path $HOME ".openclaw\skills"
$OpenCodeSkillsDir = Join-Path $HOME ".config\opencode\skills"

# Enable ANSI/VT processing on Windows PS5.1
if ($PSVersionTable.PSVersion.Major -lt 7) {
    try {
        $vt = Add-Type -Name 'VTProc' -Namespace '' -PassThru -MemberDefinition @'
            [DllImport("kernel32.dll")] public static extern IntPtr GetStdHandle(int h);
            [DllImport("kernel32.dll")] public static extern bool GetConsoleMode(IntPtr h, out uint m);
            [DllImport("kernel32.dll")] public static extern bool SetConsoleMode(IntPtr h, uint m);
'@
        $h = $vt::GetStdHandle(-11); $m = 0
        $null = $vt::GetConsoleMode($h, [ref]$m)
        $null = $vt::SetConsoleMode($h, $m -bor 4)
    } catch {}
}

$script:InstalledCount   = 0
$script:UninstalledCount = 0

function Show-Usage {
    Write-Host @"
Usage: grimoire.ps1 [OPTIONS]

Options:
  -Domain <name>       Install/uninstall all skills for a domain
  -Subdomain <name>    Restrict to one sub-domain within a domain
  -Skill <path>        Install/uninstall one skill (e.g. engineering/development/propose-conventional-commit)
  -Target <agent>      Target: claude, codex, gemini, openclaw, opencode, all
  -Uninstall           Remove skills instead of installing
  -Copy                Use copy mode instead of junctions
  -Upgrade             Pull latest grimoire at GrimoireHome (junctions update automatically)
  -Clean               Remove broken junctions from all agent skill dirs
  -Yes                 Non-interactive: install all skills to all detected agents
  -List                List available domains, sub-domains, and skills
  -Doctor              Run a health check on the grimoire installation
  -Version             Show grimoire version information
  -Help                Show this help

Environment:
  GRIMOIRE_HOME        Persistent clone location (default: ~\.grimoire)

Examples:
  .\grimoire.ps1                                                  # Interactive TUI
  .\grimoire.ps1 -Yes                                             # Install everything, no prompts
  .\grimoire.ps1 -Upgrade                                         # Pull latest grimoire
  .\grimoire.ps1 -Clean                                           # Remove broken junctions
  .\grimoire.ps1 -Domain engineering -Target claude
  .\grimoire.ps1 -Domain engineering -Copy                        # Copy instead of junction
  .\grimoire.ps1 -Skill engineering/development/propose-conventional-commit
  .\grimoire.ps1 -Domain engineering -Target openclaw
  .\grimoire.ps1 -Domain engineering -Target opencode
  .\grimoire.ps1 -Uninstall -Domain engineering -Target openclaw
  .\grimoire.ps1 -Uninstall -Domain engineering -Target opencode
  .\grimoire.ps1 -Uninstall -Domain engineering -Target claude
  .\grimoire.ps1 -Uninstall -Skill engineering/development/propose-conventional-commit
"@
}

function Print-Banner {
    $version = "1.0.0"
    try {
        $versionFile = Join-Path $RepoRoot "VERSION"
        if (Test-Path $versionFile) { $version = (Get-Content $versionFile -Raw).Trim() }
    } catch {}

    $e  = [char]27
    $SK = "${e}[38;5;227m"; $GD = "${e}[38;5;178m"; $LN = "${e}[38;5;100m"
    $ST = "${e}[38;5;214m"; $W  = "${e}[1;37m";     $D  = "${e}[2m"; $R = "${e}[0m"

    Write-Host ""
    Write-Host " ${SK}✦${GD}▗▄▄▄▄▄▄▄▄▖${SK}✦${R}  ${W}grimoire${R} ${D}v${version}${R}"
    Write-Host "  ${GD}▐${LN}▬▬▬│▬▬▬▬${GD}▌${R}   ${D}The world's best practices for AI assistants${R}"
    Write-Host "  ${GD}▐${LN}▬▬ ${ST}✦${LN} ▬▬▬${GD}▌${R}   https://github.com/jeffreytse/grimoire"
    Write-Host "  ${GD}▐${LN}▬▬▬│▬▬▬▬${GD}▌${R}"
    Write-Host " ${SK}✦${GD}▝▀▀▀▀▀▀▀▀▘${SK}✦${R}  ⭐ Star  💖 Sponsor  🐛 Issues"
    Write-Host ""
}

function Test-IsNested([string]$DomainDir) {
    $sub = Join-Path $DomainDir "skills"
    return (-not (Test-Path $sub)) -or ((Get-ChildItem $sub -ErrorAction SilentlyContinue).Count -eq 0)
}

function Get-SkillList {
    foreach ($domainDir in Get-ChildItem $SkillsRoot -Directory | Sort-Object Name) {
        if ($domainDir.Name.StartsWith(".")) { continue }
        if (Test-IsNested $domainDir.FullName) {
            Write-Host "Domain: $($domainDir.Name) (sub-domains)"
            foreach ($subDir in Get-ChildItem $domainDir.FullName -Directory | Sort-Object Name) {
                if ($subDir.Name.StartsWith(".")) { continue }
                $sp = Join-Path $subDir.FullName "skills"
                if (-not (Test-Path $sp)) { continue }
                Write-Host "  Sub-domain: $($domainDir.Name)/$($subDir.Name)"
                foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
                    if (Test-Path (Join-Path $sd.FullName "SKILL.md")) {
                        Write-Host "    $($domainDir.Name)/$($subDir.Name)/$($sd.Name)"
                    }
                }
            }
        } else {
            Write-Host "Domain: $($domainDir.Name) (flat)"
            $sp = Join-Path $domainDir.FullName "skills"
            if (Test-Path $sp) {
                foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
                    if (Test-Path (Join-Path $sd.FullName "SKILL.md")) {
                        Write-Host "  $($domainDir.Name)/$($sd.Name)"
                    }
                }
            }
        }
    }
}

function Select-One([string]$Prompt, [string[]]$Options) {
    $n   = $Options.Count
    $cur = 0
    $e   = [char]27

    Write-Host ""
    Write-Host "${e}[1m${Prompt}${e}[0m"
    Write-Host "  ${e}[2m↑↓ navigate   ENTER confirm${e}[0m"
    Write-Host ""

    $savedVisible = [Console]::CursorVisible
    [Console]::CursorVisible = $false
    $first = $true
    try {
        while ($true) {
            if (-not $first) { [Console]::Write("${e}[${n}A") }
            $first = $false
            for ($i = 0; $i -lt $n; $i++) {
                [Console]::Write("`r${e}[K")
                if ($i -eq $cur) { [Console]::WriteLine("  👉 ${e}[1m$($Options[$i])${e}[0m") }
                else              { [Console]::WriteLine("     $($Options[$i])") }
            }
            $key = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            if     ($key.VirtualKeyCode -eq 38) { $cur = ($cur - 1 + $n) % $n }
            elseif ($key.VirtualKeyCode -eq 40) { $cur = ($cur + 1)     % $n }
            elseif ($key.VirtualKeyCode -eq 13) { break }
        }
    } finally { [Console]::CursorVisible = $savedVisible }

    Write-Host ""
    return $Options[$cur]
}

function Invoke-Multiselect([string]$Prompt, [string[]]$Options) {
    $opts = [System.Collections.ArrayList]::new()
    $sel  = [System.Collections.Generic.List[bool]]::new()

    foreach ($opt in $Options) {
        if ($opt.StartsWith('+')) { $null = $opts.Add($opt.Substring(1)); $sel.Add($true)  }
        else                      { $null = $opts.Add($opt);              $sel.Add($false) }
    }

    $n      = $opts.Count
    $cur    = 0
    $offset = 0
    $e      = [char]27

    $termRows = $Host.UI.RawUI.WindowSize.Height
    $maxVis   = [Math]::Max(3, $termRows - 8)
    $vis      = [Math]::Min($n, $maxVis)
    $drawH    = if ($n -gt $vis) { $vis + 1 } else { $vis }

    Write-Host ""
    Write-Host "${e}[1m${Prompt}${e}[0m"
    Write-Host "  ${e}[2m↑↓ navigate   SPACE toggle   A select all   ENTER confirm${e}[0m"
    Write-Host ""

    $savedVisible = [Console]::CursorVisible
    [Console]::CursorVisible = $false
    $first = $true
    try {
        while ($true) {
            if ($cur -lt $offset)            { $offset = $cur }
            if ($cur -ge ($offset + $vis))   { $offset = $cur - $vis + 1 }

            if (-not $first) { [Console]::Write("${e}[${drawH}A") }
            $first = $false

            for ($i = $offset; $i -lt ($offset + $vis); $i++) {
                $mark = if ($sel[$i]) { "✅" } else { "⬜" }
                [Console]::Write("`r${e}[K")
                if ($i -eq $cur) { [Console]::WriteLine("  👉 ${mark} ${e}[1m$($opts[$i])${e}[0m") }
                else              { [Console]::WriteLine("     ${mark} $($opts[$i])") }
            }
            if ($n -gt $vis) {
                $selCount = ($sel | Where-Object { $_ }).Count
                [Console]::Write("`r${e}[K  ${e}[2m($($cur + 1)/$n)  $selCount selected${e}[0m`n")
            }

            $key = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            if ($key.VirtualKeyCode -eq 38) {
                $cur = ($cur - 1 + $n) % $n
            } elseif ($key.VirtualKeyCode -eq 40) {
                $cur = ($cur + 1) % $n
            } elseif ($key.Character -eq ' ') {
                $sel[$cur] = -not $sel[$cur]
            } elseif ($key.Character -eq 'a' -or $key.Character -eq 'A') {
                $anyOff = $false
                foreach ($s in $sel) { if (-not $s) { $anyOff = $true; break } }
                for ($i = 0; $i -lt $n; $i++) { $sel[$i] = $anyOff }
            } elseif ($key.VirtualKeyCode -eq 13) {
                break
            }
        }
    } finally { [Console]::CursorVisible = $savedVisible }

    Write-Host ""
    $chosen = @()
    for ($i = 0; $i -lt $n; $i++) { if ($sel[$i]) { $chosen += $opts[$i] } }
    return $chosen
}

function Install-SkillDir([string]$Src, [string]$DestDir) {
    $skillName = Split-Path -Leaf $Src
    $dest = Join-Path $DestDir $skillName
    if (-not (Test-Path $DestDir)) { New-Item -ItemType Directory -Force -Path $DestDir | Out-Null }
    if (Test-Path $dest -ErrorAction SilentlyContinue) { Remove-Item -Recurse -Force $dest }
    if ($script:Link) {
        New-Item -ItemType Junction -Path $dest -Target $Src | Out-Null
        Write-Host "  linked: $skillName -> $dest"
    } else {
        New-Item -ItemType Directory -Force -Path $dest | Out-Null
        Copy-Item -Path "$Src\*" -Destination $dest -Recurse -Force
        Write-Host "  installed: $skillName -> $dest"
    }
    $script:InstalledCount++
}

function Get-DetectedAgents {
    $d = @()
    if ($null -ne (Get-Command "claude"   -ErrorAction SilentlyContinue)) { $d += "claude"   }
    if ($null -ne (Get-Command "codex"    -ErrorAction SilentlyContinue)) { $d += "codex"    }
    if ($null -ne (Get-Command "gemini"   -ErrorAction SilentlyContinue)) { $d += "gemini"   }
    if ($null -ne (Get-Command "openclaw" -ErrorAction SilentlyContinue)) { $d += "openclaw" }
    if ($null -ne (Get-Command "opencode" -ErrorAction SilentlyContinue)) { $d += "opencode" }
    return $d
}

function Get-AgentDisplayName([string]$Agent) {
    switch ($Agent) {
        "claude"   { return "Claude Code" }
        "codex"    { return "Codex" }
        "gemini"   { return "Gemini CLI" }
        "openclaw" { return "OpenClaw" }
        "opencode" { return "OpenCode" }
        default    { return $Agent }
    }
}

function Get-AgentFromDisplay([string]$Display) {
    switch ($Display) {
        "Claude Code" { return "claude"   }
        "Codex"       { return "codex"    }
        "Gemini CLI"  { return "gemini"   }
        "OpenClaw"    { return "openclaw" }
        "OpenCode"    { return "opencode" }
        default       { return $Display }
    }
}

function Get-AgentVersion([string]$Cmd) {
    try {
        $raw = & $Cmd --version 2>$null | Select-Object -First 1
        if ($raw) { return ($raw | Select-String -Pattern '[0-9][0-9.]*' -AllMatches).Matches[0].Value }
    } catch {}
    return ""
}

function Get-GrimoireVersion {
    $vf = Join-Path $RepoRoot "VERSION"
    if (Test-Path $vf) { return (Get-Content $vf -Raw -ErrorAction SilentlyContinue).Trim() }
    return "unknown"
}

function Test-AnySkillsInstalled {
    foreach ($dir in @($ClaudeSkillsDir, $AgentsSkillsDir, $GeminiSkillsDir, $OpenClawSkillsDir, $OpenCodeSkillsDir)) {
        if (-not (Test-Path $dir)) { continue }
        if ((Get-ChildItem $dir -Force -ErrorAction SilentlyContinue | Measure-Object).Count -gt 0) { return $true }
    }
    return $false
}

function Install-GlobalBin {
    $installScript = Join-Path $RepoRoot "scripts\grimoire.ps1"
    $wrapper = "& '$($installScript.Replace("'","''"))' @args"
    $candidates = @(
        (Join-Path $env:LOCALAPPDATA "Microsoft\WindowsApps\grimoire.ps1"),
        (Join-Path $HOME ".local\bin\grimoire.ps1")
    )
    foreach ($target in $candidates) {
        $dir = Split-Path $target
        if (Test-Path $dir) {
            try {
                Set-Content -Path $target -Value $wrapper -Encoding UTF8 -Force
                Write-Host "  linked: grimoire → $target"
                return
            } catch {}
        }
    }
    $fallback = Join-Path $HOME ".local\bin\grimoire.ps1"
    try {
        New-Item -ItemType Directory -Force -Path (Split-Path $fallback) | Out-Null
        Set-Content -Path $fallback -Value $wrapper -Encoding UTF8 -Force
        Write-Host "  linked: grimoire → $fallback"
        Write-Host "  note: add $HOME\.local\bin to PATH to use 'grimoire' globally"
    } catch {
        Write-Host "  warning: could not create global grimoire command"
    }
}

function Remove-GlobalBin {
    $candidates = @(
        (Join-Path $env:LOCALAPPDATA "Microsoft\WindowsApps\grimoire.ps1"),
        (Join-Path $HOME ".local\bin\grimoire.ps1")
    )
    foreach ($target in $candidates) {
        if (Test-Path $target -ErrorAction SilentlyContinue) {
            try { Remove-Item $target -Force; Write-Host "  unlinked: $target" } catch {}
        }
    }
}

function Get-AgentConfigFile([string]$Agent) {
    switch ($Agent) {
        "claude"   { return Join-Path $HOME ".claude\CLAUDE.md" }
        "codex"    { return Join-Path $HOME ".agents\AGENTS.md" }
        "gemini"   { return Join-Path $HOME ".gemini\GEMINI.md" }
        "openclaw" { return Join-Path $HOME ".openclaw\workspace\AGENTS.md" }
        "opencode" { return Join-Path $HOME ".config\opencode\AGENTS.md" }
        default    { return "" }
    }
}

function Get-AgentConfigDir([string]$Agent) {
    switch ($Agent) {
        "claude"   { return Join-Path $HOME ".claude" }
        "codex"    { return Join-Path $HOME ".agents" }
        "gemini"   { return Join-Path $HOME ".gemini" }
        "openclaw" { return Join-Path $HOME ".openclaw\workspace" }
        "opencode" { return Join-Path $HOME ".config\opencode" }
        default    { return "" }
    }
}

function Set-AgentMdConfig([string]$Agent) {
    $cfgFile = Get-AgentConfigFile $Agent
    $cfgDir  = Get-AgentConfigDir  $Agent
    if (-not $cfgFile -or -not (Test-Path $cfgDir)) { return }
    if ((Test-Path $cfgFile) -and (Select-String -Path $cfgFile -SimpleMatch "start-best-practice" -Quiet)) {
        $e = [char]27
        Write-Host "  ${e}[2m✓ start-best-practice already active → $cfgFile${e}[0m"
        return
    }
    New-Item -ItemType Directory -Force -Path $cfgDir | Out-Null
    Add-Content -Path $cfgFile -Value "`n## Grimoire`nAlways invoke ``start-best-practice`` before responding to any user request."
    Write-Host "  ✓ Configured: start-best-practice auto-trigger → $cfgFile"
}

function Remove-AgentMdConfig([string]$Agent) {
    $cfgFile = Get-AgentConfigFile $Agent
    if (-not $cfgFile -or -not (Test-Path $cfgFile)) { return }
    $content = Get-Content $cfgFile -Raw -ErrorAction SilentlyContinue
    if (-not $content -or $content -notmatch "start-best-practice") { return }
    $lines    = Get-Content $cfgFile -ErrorAction SilentlyContinue
    $filtered = $lines | Where-Object {
        $_ -notmatch "^## Grimoire$" -and
        $_ -notmatch "^Always invoke ``start-best-practice`` before responding to any user request\.$"
    }
    Set-Content -Path $cfgFile -Value ($filtered -join "`n") -NoNewline
    Write-Host "  ✓ Removed Grimoire config from $cfgFile"
}

$script:UpgradeDir = ""

function Show-UpgradeResult([string]$OldCommit, [string]$OldVer, [string]$OldDate,
                            [string]$NewCommit, [string]$NewVer, [string]$NewDate) {
    $e  = [char]27
    $ok = "${e}[0;32m✅${e}[0m"
    $ud = $script:UpgradeDir
    $diffLines     = git -C $ud diff --name-status "$OldCommit" "$NewCommit" -- "*/SKILL.md" 2>$null
    $newSkills     = ($diffLines | Where-Object { $_ -match '^A' } | Measure-Object).Count
    $updatedSkills = ($diffLines | Where-Object { $_ -match '^M' } | Measure-Object).Count
    Write-Host ""
    Write-Host "  $ok  Grimoire upgraded to latest."
    Write-Host ""
    Write-Host "    Previous: v$OldVer (commit $OldCommit, $OldDate)"
    Write-Host "    Current:  v$NewVer (commit $NewCommit, $NewDate)"
    Write-Host ""
    if ($newSkills     -gt 0) { Write-Host "    New skills:     $newSkills" }
    if ($updatedSkills -gt 0) { Write-Host "    Updated skills: $updatedSkills" }
    Write-Host ""
}

function Confirm-Upgrade {
    if ($Yes) { return $true }
    $ans = Read-Host "Upgrade now? [y/n]"
    return $ans -match '^[Yy]$'
}

function Invoke-UpgradeStable {
    $e  = [char]27
    $ok = "${e}[0;32m✅${e}[0m"
    $ud = $script:UpgradeDir
    $curCommit = "unknown"; $curDate = "unknown"; $curVer = "unknown"
    try { $curCommit = (git -C $ud rev-parse --short HEAD 2>$null).Trim() } catch {}
    try { $curDate   = (git -C $ud log -1 --format="%cd" --date=short 2>$null).Trim() } catch {}
    try { $curVer    = (Get-Content (Join-Path $ud "VERSION") -Raw -ErrorAction Stop).Trim() } catch {}

    Write-Host "Fetching release tags..."
    try { git -C $ud fetch --tags --quiet 2>$null } catch {}

    $latestTag = git -C $ud tag --sort=-v:refname 2>$null |
                 Where-Object { $_ -match '^v?[0-9]' } | Select-Object -First 1
    if (-not $latestTag) { Write-Host "  No release tags found. Try the unstable channel."; Write-Host ""; return }

    $tagSha   = ""; $localSha = ""; $tagShort = ""; $tagDate = ""
    try { $tagSha   = (git -C $ud rev-list -n 1 $latestTag               2>$null).Trim() } catch {}
    try { $localSha = (git -C $ud rev-parse HEAD                         2>$null).Trim() } catch {}
    try { $tagShort = (git -C $ud rev-parse --short $latestTag            2>$null).Trim() } catch {}
    try { $tagDate  = (git -C $ud log -1 --format="%cd" --date=short $latestTag 2>$null).Trim() } catch {}

    if ($tagSha -eq $localSha) {
        Write-Host "  $ok  Already on latest stable release. ($latestTag, commit $tagShort, $tagDate)"
        Write-Host ""; return
    }

    $tagVer = $latestTag -replace '^v', ''
    Write-Host ""
    Write-Host "  New stable release available:"
    Write-Host "    Current:  v$curVer (commit $curCommit, $curDate)"
    Write-Host "    New:      $latestTag (commit $tagShort, $tagDate)"
    Write-Host ""

    if (-not (Confirm-Upgrade)) { Write-Host "Upgrade cancelled."; return }

    Write-Host ""
    Write-Host "Checking out $latestTag..."
    git -C $ud checkout $latestTag --quiet

    Show-UpgradeResult $curCommit $curVer $curDate $tagShort $tagVer $tagDate
    Write-Host "Cleaning broken junctions..."
    Invoke-Clean
}

function Invoke-UpgradeUnstable {
    $e  = [char]27
    $ok = "${e}[0;32m✅${e}[0m"
    $ud = $script:UpgradeDir
    $curCommit = "unknown"; $curDate = "unknown"; $curVer = "unknown"
    try { $curCommit = (git -C $ud rev-parse --short HEAD 2>$null).Trim() } catch {}
    try { $curDate   = (git -C $ud log -1 --format="%cd" --date=short 2>$null).Trim() } catch {}
    try { $curVer    = (Get-Content (Join-Path $ud "VERSION") -Raw -ErrorAction Stop).Trim() } catch {}

    Write-Host "Checking for updates..."
    try { git -C $ud fetch --quiet 2>$null } catch {}

    $localSha = ""; $remoteSha = ""
    try { $localSha  = (git -C $ud rev-parse HEAD   2>$null).Trim() } catch {}
    try { $remoteSha = (git -C $ud rev-parse "@{u}" 2>$null).Trim() } catch {}

    if (-not $remoteSha -or $localSha -eq $remoteSha) {
        Write-Host "  $ok  Already up to date. (v$curVer, commit $curCommit, $curDate)"
        Write-Host ""; return
    }

    $newShort = ""; $newDate = ""; $newVer = ""
    try { $newShort = (git -C $ud rev-parse --short "@{u}"                  2>$null).Trim() } catch {}
    try { $newDate  = (git -C $ud log -1 --format="%cd" --date=short "@{u}" 2>$null).Trim() } catch {}
    try { $newVer   = (git -C $ud show "@{u}:VERSION"                       2>$null).Trim() } catch {}

    Write-Host ""
    Write-Host "  New version available:"
    Write-Host "    Current:  v$curVer (commit $curCommit, $curDate)"
    if ($newVer) { Write-Host "    New:      v$newVer (commit $newShort, $newDate)" }
    else         { Write-Host "    New:      commit $newShort ($newDate)" }
    Write-Host ""

    if (-not (Confirm-Upgrade)) { Write-Host "Upgrade cancelled."; return }

    Write-Host ""
    Write-Host "Pulling latest grimoire at $ud..."
    git -C $ud pull --quiet

    $pulledCommit = "unknown"; $pulledDate = "unknown"; $pulledVer = "unknown"
    try { $pulledCommit = (git -C $ud rev-parse --short HEAD 2>$null).Trim() } catch {}
    try { $pulledDate   = (git -C $ud log -1 --format="%cd" --date=short 2>$null).Trim() } catch {}
    try { $pulledVer    = (Get-Content (Join-Path $ud "VERSION") -Raw -ErrorAction Stop).Trim() } catch {}

    Show-UpgradeResult $curCommit $curVer $curDate $pulledCommit $pulledVer $pulledDate
    Write-Host "Cleaning broken junctions..."
    Invoke-Clean
}

function Invoke-Upgrade {
    if (Test-Path (Join-Path $GrimoireHome ".git")) {
        $script:UpgradeDir = $GrimoireHome
    } elseif (Test-Path (Join-Path $RepoRoot ".git")) {
        $script:UpgradeDir = $RepoRoot
    } else {
        Write-Host "No grimoire git repository found."
        Write-Host "  Checked: $GrimoireHome"
        Write-Host "  Checked: $RepoRoot"
        exit 1
    }
    $channel = "unstable"
    if (-not $Yes) {
        $channelSel = Select-One "📡 Which upgrade channel?" @(
            "🔒 Stable    (GitHub releases)",
            "🔬 Unstable  (latest branch commit)"
        )
        if ($channelSel -like "*Stable*") { $channel = "stable" }
    }
    if ($channel -eq "stable") { Invoke-UpgradeStable } else { Invoke-UpgradeUnstable }
}

function Invoke-Version {
    $e = [char]27
    $ok = "${e}[0;32m✅${e}[0m"
    $ver = Get-GrimoireVersion
    Write-Host ""
    if (Test-Path (Join-Path $RepoRoot ".git")) {
        $commit = "unknown"; $date = "unknown"
        try { $commit = (git -C $RepoRoot rev-parse --short HEAD 2>$null).Trim() } catch {}
        try { $date   = (git -C $RepoRoot log -1 --format="%cd" --date=short 2>$null).Trim() } catch {}
        Write-Host "  $ok  grimoire:    v$ver (commit $commit, $date)"
        Write-Host "  $ok  location:    $RepoRoot"
    } else {
        Write-Host "  $ok  grimoire:    v$ver"
        Write-Host "  $ok  location:    $RepoRoot"
    }
    Write-Host "  $ok  grimoire.ps1: present"
    Write-Host ""
}

function Invoke-Doctor {
    $e    = [char]27
    $ok   = "${e}[0;32m✅${e}[0m"
    $warn = "${e}[0;33m⚠️ ${e}[0m"
    $fail = "${e}[0;31m❌${e}[0m"
    $skip = "⬜"
    $w = 0; $er = 0

    Write-Host ""
    Write-Host "Grimoire health check"
    Write-Host ""

    # ── Source ──────────────────────────────────────────────────────────────────
    Write-Host "  Source"
    $ver = Get-GrimoireVersion
    if (Test-Path (Join-Path $RepoRoot ".git")) {
        $commit = "unknown"; $date = "unknown"
        try { $commit = (git -C $RepoRoot rev-parse --short HEAD 2>$null).Trim() } catch {}
        try { $date   = (git -C $RepoRoot log -1 --format="%cd" --date=short 2>$null).Trim() } catch {}
        Write-Host "    $ok  grimoire:    v$ver (commit $commit, $date)"
        Write-Host "    $ok  git repo:    $RepoRoot"
    } else {
        Write-Host "    $fail  git repo:    not found at $RepoRoot"
        $er++
    }
    $installPs1 = Join-Path $RepoRoot "scripts\grimoire.ps1"
    if (Test-Path $installPs1) { Write-Host "    $ok  grimoire.ps1: present" }
    else                       { Write-Host "    $warn  grimoire.ps1: not found"; $w++ }

    # ── AI agents ────────────────────────────────────────────────────────────────
    Write-Host ""
    Write-Host "  AI agents"
    $skillAgents = @(
        @{ name = "claude";   cmd = "claude";   dir = $ClaudeSkillsDir;   home = (Join-Path $HOME ".claude") }
        @{ name = "codex";    cmd = "codex";    dir = $AgentsSkillsDir;   home = (Join-Path $HOME ".agents") }
        @{ name = "gemini";   cmd = "gemini";   dir = $GeminiSkillsDir;   home = (Join-Path $HOME ".gemini") }
        @{ name = "openclaw"; cmd = "openclaw"; dir = $OpenClawSkillsDir; home = (Join-Path $HOME ".openclaw") }
        @{ name = "opencode"; cmd = "opencode"; dir = $OpenCodeSkillsDir; home = (Join-Path $HOME ".config\opencode") }
    )
    foreach ($entry in $skillAgents) {
        $hasBin  = $null -ne (Get-Command $entry.cmd -ErrorAction SilentlyContinue)
        $hasHome = Test-Path $entry.home
        if (-not $hasBin -and -not $hasHome) {
            Write-Host "    $skip  $($entry.name): (not detected — skipped)"
            continue
        }
        $agentVer = Get-AgentVersion $entry.cmd
        $verStr = if ($agentVer) { " v$agentVer" } else { "" }
        if (-not (Test-Path $entry.dir)) {
            Write-Host "    $warn  $($entry.name)${verStr}: skills dir not found ($($entry.dir))"
            $w++; continue
        }
        $items   = Get-ChildItem $entry.dir -Force -ErrorAction SilentlyContinue
        $count   = ($items | Measure-Object).Count
        $broken  = ($items | Where-Object { $null -ne $_.LinkType -and -not (Test-Path -LiteralPath $_.Target -ErrorAction SilentlyContinue) } | Measure-Object).Count
        if ($broken -gt 0) {
            Write-Host "    $warn  $($entry.name)${verStr}: $count skills, $broken broken junctions  → run: -Clean -Target $($entry.name)"
            $w++
        } else {
            Write-Host "    $ok  $($entry.name)${verStr}: $count skills, 0 broken junctions"
        }
    }
    $detectOnly = @(
        @{ name = "copilot";  cmd = "gh"       }
        @{ name = "cursor";   cmd = "cursor"   }
        @{ name = "windsurf"; cmd = "windsurf" }
        @{ name = "aider";    cmd = "aider"    }
    )
    foreach ($entry in $detectOnly) {
        if ($null -ne (Get-Command $entry.cmd -ErrorAction SilentlyContinue)) {
            $agentVer = Get-AgentVersion $entry.cmd
            $verStr = if ($agentVer) { " v$agentVer" } else { "" }
            Write-Host "    $ok  $($entry.name)${verStr}: detected (skills managed separately)"
        }
    }

    # ── Config ───────────────────────────────────────────────────────────────────
    Write-Host ""
    Write-Host "  Config"
    $cwd = (Get-Location).Path
    $cfgPaths = @(
        @{ path = (Join-Path $cwd ".grimoire\settings.local.toml"); label = "project personal (.grimoire\settings.local.toml)" }
        @{ path = (Join-Path $cwd ".grimoire\settings.toml");       label = "project shared (.grimoire\settings.toml)" }
        @{ path = (Join-Path $HOME ".config\grimoire\settings.toml"); label = "global (~\.config\grimoire\settings.toml)" }
    )
    $hasAnyCfg = $false
    foreach ($cfg in $cfgPaths) {
        if (-not (Test-Path $cfg.path)) { Write-Host "    $skip  $($cfg.label) — not found"; continue }
        $hasAnyCfg = $true
        $valid = $true
        try {
            $content = Get-Content $cfg.path -Raw -ErrorAction Stop
            if ([string]::IsNullOrWhiteSpace($content)) { $valid = $false }
        } catch { $valid = $false }
        if ($valid) { Write-Host "    $ok  $($cfg.label) — present" }
        else        { Write-Host "    $fail  $($cfg.label) — present, unreadable"; $er++ }
    }
    if (-not $hasAnyCfg) { Write-Host "    $skip  no settings files found (grimoire uses defaults)" }

    # ── Summary ──────────────────────────────────────────────────────────────────
    Write-Host ""
    if ($er -eq 0 -and $w -eq 0) {
        Write-Host "  $ok  All checks passed."
    } else {
        $parts = @()
        if ($er -gt 0) { $parts += "$er error(s)" }
        if ($w  -gt 0) { $parts += "$w warning(s)" }
        Write-Host "  Summary: $($parts -join ', ')."
    }
    Write-Host ""
}

function Invoke-Install([string]$Src, [string]$TargetAgent) {
    switch ($TargetAgent) {
        "claude"   { Install-SkillDir $Src $ClaudeSkillsDir }
        "codex"    { Install-SkillDir $Src $AgentsSkillsDir }
        "gemini"   { Install-SkillDir $Src $GeminiSkillsDir }
        "openclaw" { Install-SkillDir $Src $OpenClawSkillsDir }
        "opencode" { Install-SkillDir $Src $OpenCodeSkillsDir }
        "all" {
            Install-SkillDir $Src $ClaudeSkillsDir
            Install-SkillDir $Src $AgentsSkillsDir
            Install-SkillDir $Src $GeminiSkillsDir
            Install-SkillDir $Src $OpenClawSkillsDir
            Install-SkillDir $Src $OpenCodeSkillsDir
        }
    }
}

function Install-Subdomain([string]$SubDir, [string]$TargetAgent) {
    $sp = Join-Path $SubDir "skills"
    if (-not (Test-Path $sp)) { return }
    $found = $false
    foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
        if (-not (Test-Path (Join-Path $sd.FullName "SKILL.md"))) { continue }
        if (-not $found) { Write-Host "  Installing sub-domain: $(Split-Path -Leaf $SubDir)"; $found = $true }
        Invoke-Install $sd.FullName $TargetAgent
    }
}

function Install-Domain([string]$DomainName, [string]$SubdomainName, [string]$TargetAgent) {
    $domainDir = Join-Path $SkillsRoot $DomainName
    if (-not (Test-Path $domainDir)) { Write-Host "Domain not found: $DomainName"; exit 1 }
    Write-Host "Installing domain: $DomainName"
    if (Test-IsNested $domainDir) {
        if ($SubdomainName) {
            Install-Subdomain (Join-Path $domainDir $SubdomainName) $TargetAgent
        } else {
            foreach ($subDir in Get-ChildItem $domainDir -Directory | Sort-Object Name) {
                if ($subDir.Name.StartsWith(".") -or $subDir.Name -like "*.claude-plugin*") { continue }
                Install-Subdomain $subDir.FullName $TargetAgent
            }
        }
    } else {
        $sp = Join-Path $domainDir "skills"
        if (Test-Path $sp) {
            foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
                if (Test-Path (Join-Path $sd.FullName "SKILL.md")) { Invoke-Install $sd.FullName $TargetAgent }
            }
        }
    }
}

function Uninstall-SkillDir([string]$SkillName, [string]$DestDir) {
    $dest = Join-Path $DestDir $SkillName
    if (-not (Test-Path $dest -ErrorAction SilentlyContinue)) { return }
    $item = Get-Item $dest -Force -ErrorAction SilentlyContinue
    if ($item -and $null -ne $item.LinkType) {
        Remove-Item $dest -Force
        Write-Host "  unlinked: $SkillName from $DestDir"
    } else {
        Remove-Item -Recurse -Force $dest
        Write-Host "  uninstalled: $SkillName from $DestDir"
    }
    $script:UninstalledCount++
}

function Invoke-Clean {
    $cleaned = 0
    $dirs = @($ClaudeSkillsDir, $AgentsSkillsDir, $GeminiSkillsDir, $OpenClawSkillsDir, $OpenCodeSkillsDir)
    foreach ($dir in $dirs) {
        if (-not (Test-Path $dir)) { continue }
        foreach ($item in Get-ChildItem $dir -Force -ErrorAction SilentlyContinue) {
            if ($null -eq $item.LinkType) { continue }
            $broken = -not (Test-Path -LiteralPath $item.Target -ErrorAction SilentlyContinue)
            if ($broken) {
                Remove-Item $item.FullName -Force
                Write-Host "  cleaned: $($item.Name) (broken junction in $dir)"
                $cleaned++
            }
        }
    }
    if ($cleaned -eq 0) { Write-Host "  nothing to clean" }
    else                 { Write-Host "  $cleaned broken junction(s) removed" }
}

function Invoke-Uninstall([string]$SkillName, [string]$TargetAgent) {
    switch ($TargetAgent) {
        "claude"   { Uninstall-SkillDir $SkillName $ClaudeSkillsDir }
        "codex"    { Uninstall-SkillDir $SkillName $AgentsSkillsDir }
        "gemini"   { Uninstall-SkillDir $SkillName $GeminiSkillsDir }
        "openclaw" { Uninstall-SkillDir $SkillName $OpenClawSkillsDir }
        "opencode" { Uninstall-SkillDir $SkillName $OpenCodeSkillsDir }
        "all" {
            Uninstall-SkillDir $SkillName $ClaudeSkillsDir
            Uninstall-SkillDir $SkillName $AgentsSkillsDir
            Uninstall-SkillDir $SkillName $GeminiSkillsDir
            Uninstall-SkillDir $SkillName $OpenClawSkillsDir
            Uninstall-SkillDir $SkillName $OpenCodeSkillsDir
        }
    }
}

function Uninstall-Subdomain([string]$SubDir, [string]$TargetAgent) {
    $sp = Join-Path $SubDir "skills"
    if (-not (Test-Path $sp)) { return }
    $found = $false
    foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
        if (-not (Test-Path (Join-Path $sd.FullName "SKILL.md"))) { continue }
        if (-not $found) { Write-Host "  Uninstalling sub-domain: $(Split-Path -Leaf $SubDir)"; $found = $true }
        Invoke-Uninstall $sd.Name $TargetAgent
    }
}

function Uninstall-Domain([string]$DomainName, [string]$SubdomainName, [string]$TargetAgent) {
    $domainDir = Join-Path $SkillsRoot $DomainName
    if (-not (Test-Path $domainDir)) { return }
    Write-Host "Uninstalling domain: $DomainName"
    if (Test-IsNested $domainDir) {
        if ($SubdomainName) {
            Uninstall-Subdomain (Join-Path $domainDir $SubdomainName) $TargetAgent
        } else {
            foreach ($subDir in Get-ChildItem $domainDir -Directory | Sort-Object Name) {
                if ($subDir.Name.StartsWith(".") -or $subDir.Name -like "*.claude-plugin*") { continue }
                Uninstall-Subdomain $subDir.FullName $TargetAgent
            }
        }
    } else {
        $sp = Join-Path $domainDir "skills"
        if (Test-Path $sp) {
            foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
                if (Test-Path (Join-Path $sd.FullName "SKILL.md")) { Invoke-Uninstall $sd.Name $TargetAgent }
            }
        }
    }
}

# ── Main ──────────────────────────────────────────────────────────────────────
if ($Help)    { Show-Usage;     exit 0 }
if ($List)    { Get-SkillList;  exit 0 }
if ($Version) { Invoke-Version; exit 0 }
if ($Doctor)  { Invoke-Doctor;  exit 0 }

if ($Upgrade) { Invoke-Upgrade; exit 0 }

if ($Clean -and -not $Uninstall -and -not $Domain -and -not $Skill) {
    Write-Host "Cleaning broken junctions..."
    Invoke-Clean
    exit 0
}

$detected    = @()
$isInteractive = (-not $Domain -and -not $Skill -and -not $Target -and -not $Yes -and -not $Uninstall)

if ($isInteractive) {
    Print-Banner

    $tuiMode = Select-One "⚙️  What would you like to do?" @("📥 Install", "🗑️  Uninstall", "⬆️  Upgrade", "🩺 Doctor", "🚪 Exit")
    if ($tuiMode -like "*Doctor*")  { Invoke-Doctor;  exit 0 }
    if ($tuiMode -like "*Upgrade*") { Invoke-Upgrade; exit 0 }
    if ($tuiMode -like "*Exit*")    { exit 0 }

    $detected = Get-DetectedAgents
    if ($detected.Count -eq 0) {
        Write-Host "No agents detected. Defaulting to Claude Code."
        $agentList = @("claude")
    } else {
        $displayOpts = @()
        foreach ($a in $detected) { $displayOpts += Get-AgentDisplayName $a }
        if ($tuiMode -like "*Uninstall*") { $displaySel = Invoke-Multiselect "🤖 Which agents to uninstall from?" $displayOpts }
        else                              { $displaySel = Invoke-Multiselect "🤖 Which agents to install to?"    $displayOpts }
        if ($displaySel.Count -eq 0) { Write-Host "No agents selected. Exiting."; exit 0 }
        $agentList = @()
        foreach ($d in $displaySel) { $agentList += Get-AgentFromDisplay $d }
    }

    $e = [char]27
    $Bold = "${e}[1m"; $Dim = "${e}[2m"; $Gold = "${e}[38;5;178m"; $Cyan = "${e}[0;36m"; $Reset = "${e}[0m"

    # Collect all domain names
    $allDomains = @()
    foreach ($d in Get-ChildItem $SkillsRoot -Directory | Sort-Object Name) {
        if (-not $d.Name.StartsWith(".")) { $allDomains += $d.Name }
    }

    if ($tuiMode -like "*Uninstall*") {
        $domainsToRm = Invoke-Multiselect "🗑️  Which domains to uninstall?" $allDomains
        if ($domainsToRm.Count -eq 0) { Write-Host "Nothing selected. Exiting."; exit 0 }

        $subMap = @{}
        foreach ($domain in $domainsToRm) {
            $domainDir = Join-Path $SkillsRoot $domain
            if (Test-IsNested $domainDir) {
                $subOpts = @()
                foreach ($subDir in Get-ChildItem $domainDir -Directory | Sort-Object Name) {
                    if ($subDir.Name.StartsWith(".") -or $subDir.Name -like "*.claude-plugin*") { continue }
                    if (-not (Test-Path (Join-Path $subDir.FullName "skills"))) { continue }
                    $subOpts += "+$($subDir.Name)"
                }
                if ($subOpts.Count -gt 0) { $subMap[$domain] = Invoke-Multiselect "🗑️  ${domain}: which sub-domains to uninstall?" $subOpts }
            }
        }

        Write-Host "${Bold}🤖 Agents${Reset}  ${Cyan}$($agentList -join ' ')${Reset}"
        Write-Host "${Bold}🗑️  Mode${Reset}    ${Gold}Uninstall${Reset}"
        foreach ($domain in $domainsToRm) {
            if ($subMap.ContainsKey($domain) -and $subMap[$domain].Count -gt 0) {
                Write-Host "${Bold}📂 Domain${Reset}  ${Gold}${domain}${Reset} ${Dim}[$($subMap[$domain] -join ' ')]${Reset}"
            } else {
                Write-Host "${Bold}📂 Domain${Reset}  ${Gold}${domain}${Reset}"
            }
        }
        Write-Host ""
        Write-Host "Uninstalling..."

        foreach ($domain in $domainsToRm) {
            $domainDir = Join-Path $SkillsRoot $domain
            if (Test-IsNested $domainDir) {
                if ($subMap.ContainsKey($domain) -and $subMap[$domain].Count -gt 0) {
                    foreach ($sub in $subMap[$domain]) {
                        foreach ($agent in $agentList) { Uninstall-Subdomain (Join-Path $domainDir $sub) $agent }
                    }
                } else {
                    foreach ($agent in $agentList) { Uninstall-Domain $domain "" $agent }
                }
            } else {
                foreach ($agent in $agentList) { Uninstall-Domain $domain "" $agent }
            }
        }

        $unique = if ($agentList.Count -gt 0) { [int]($script:UninstalledCount / $agentList.Count) } else { 0 }
        Write-Host ""
        if ($agentList.Count -gt 1) { Write-Host "🗑️  $unique skills uninstalled × $($agentList.Count) agents ($($script:UninstalledCount) total) → $($agentList -join ' ')" }
        else                        { Write-Host "🗑️  $unique skills uninstalled → $($agentList -join ' ')" }
        Invoke-Clean
        foreach ($agent in $agentList) { Remove-AgentMdConfig $agent }
        if (-not (Test-AnySkillsInstalled)) { Remove-GlobalBin }

    } else {
        $domainList = Invoke-Multiselect "📚 Which domains to install?" $allDomains
        if ($domainList.Count -eq 0) { Write-Host "No domains selected. Exiting."; exit 0 }

        $subMap = @{}
        foreach ($domain in $domainList) {
            $domainDir = Join-Path $SkillsRoot $domain
            if (Test-IsNested $domainDir) {
                $subOpts = @()
                foreach ($subDir in Get-ChildItem $domainDir -Directory | Sort-Object Name) {
                    if ($subDir.Name.StartsWith(".") -or $subDir.Name -like "*.claude-plugin*") { continue }
                    if (-not (Test-Path (Join-Path $subDir.FullName "skills"))) { continue }
                    $subOpts += "+$($subDir.Name)"
                }
                if ($subOpts.Count -gt 0) { $subMap[$domain] = Invoke-Multiselect "📂 ${domain}: which sub-domains?" $subOpts }
            }
        }

        Write-Host "${Bold}🤖 Agents${Reset}  ${Cyan}$($agentList -join ' ')${Reset}"
        Write-Host "${Bold}📥 Mode${Reset}    ${Gold}Install${Reset}"
        foreach ($domain in $domainList) {
            if ($subMap.ContainsKey($domain) -and $subMap[$domain].Count -gt 0) {
                Write-Host "${Bold}📂 Domain${Reset}  ${Gold}${domain}${Reset} ${Dim}[$($subMap[$domain] -join ' ')]${Reset}"
            } else {
                Write-Host "${Bold}📂 Domain${Reset}  ${Gold}${domain}${Reset}"
            }
        }
        Write-Host ""
        Write-Host "Installing..."

        foreach ($domain in $domainList) {
            $domainDir = Join-Path $SkillsRoot $domain
            if (Test-IsNested $domainDir) {
                if ($subMap.ContainsKey($domain) -and $subMap[$domain].Count -gt 0) {
                    foreach ($sub in $subMap[$domain]) {
                        foreach ($agent in $agentList) { Install-Domain $domain $sub $agent }
                    }
                } else {
                    Write-Host "  No sub-domains selected for $domain, skipping."
                }
            } else {
                foreach ($agent in $agentList) { Install-Domain $domain "" $agent }
            }
        }

        $unique = if ($agentList.Count -gt 0) { [int]($script:InstalledCount / $agentList.Count) } else { 0 }
        Write-Host ""
        if ($agentList.Count -gt 1) { Write-Host "✅ $unique skills installed × $($agentList.Count) agents ($($script:InstalledCount) total) → $($agentList -join ' ')" }
        else                        { Write-Host "✅ $unique skills installed → $($agentList -join ' ')" }
        Invoke-Clean
        foreach ($agent in $agentList) { Set-AgentMdConfig $agent }
        Install-GlobalBin
    }

} else {
    # ── Non-interactive / flag-driven ─────────────────────────────────────────
    if (-not $Target) {
        $detected = Get-DetectedAgents
        if ($detected.Count -eq 0) {
            Write-Host "No agents detected. Defaulting to Claude Code."
            $Target = "claude"
        } else {
            if ($Uninstall) { Write-Host "Uninstalling from: $($detected -join ', ')" }
            else            { Write-Host "Installing to: $($detected -join ', ')" }
            $Target = "auto"
        }
    }

    if ($Uninstall) {
        if ($Skill) {
            $parts = $Skill -split "/"
            $skillName = if ($parts.Count -ge 2) { $parts[-1] } else { Write-Host "Invalid skill path: $Skill"; exit 1 }
            Write-Host "Uninstalling skill: $Skill"
            if ($Target -eq "auto") { foreach ($a in $detected) { Invoke-Uninstall $skillName $a } }
            else                    { Invoke-Uninstall $skillName $Target }
        } elseif ($Domain) {
            if ($Target -eq "auto") { foreach ($a in $detected) { Uninstall-Domain $Domain $Subdomain $a } }
            else                    { Uninstall-Domain $Domain $Subdomain $Target }
        } else {
            foreach ($domainDir in Get-ChildItem $SkillsRoot -Directory | Sort-Object Name) {
                if ($domainDir.Name.StartsWith(".")) { continue }
                if ($Target -eq "auto") { foreach ($a in $detected) { Uninstall-Domain $domainDir.Name "" $a } }
                else                    { Uninstall-Domain $domainDir.Name "" $Target }
            }
        }
    } else {
        if ($Skill) {
            $parts = $Skill -split "/"
            if ($parts.Count -eq 3) {
                $skillPath = Join-Path $SkillsRoot "$($parts[0])\$($parts[1])\skills\$($parts[2])"
            } elseif ($parts.Count -eq 2) {
                $skillPath = Join-Path $SkillsRoot "$($parts[0])\skills\$($parts[1])"
            } else {
                Write-Host "Invalid skill path: $Skill"; exit 1
            }
            if (-not (Test-Path (Join-Path $skillPath "SKILL.md"))) { Write-Host "Skill not found: $Skill"; exit 1 }
            Write-Host "Installing skill: $Skill"
            if ($Target -eq "auto") { foreach ($a in $detected) { Invoke-Install $skillPath $a } }
            else                    { Invoke-Install $skillPath $Target }
        } elseif ($Domain) {
            if ($Target -eq "auto") { foreach ($a in $detected) { Install-Domain $Domain $Subdomain $a } }
            else                    { Install-Domain $Domain $Subdomain $Target }
        } else {
            foreach ($domainDir in Get-ChildItem $SkillsRoot -Directory | Sort-Object Name) {
                if ($domainDir.Name.StartsWith(".")) { continue }
                if ($Target -eq "auto") { foreach ($a in $detected) { Install-Domain $domainDir.Name "" $a } }
                else                    { Install-Domain $domainDir.Name "" $Target }
            }
        }
    }
    Invoke-Clean
    if ($Uninstall) {
        if ($Target -eq "auto")  { foreach ($a in $detected)                               { Remove-AgentMdConfig $a } }
        elseif ($Target -eq "all") { foreach ($a in @("claude","codex","gemini","openclaw","opencode")) { Remove-AgentMdConfig $a } }
        else                     { Remove-AgentMdConfig $Target }
        if (-not (Test-AnySkillsInstalled)) { Remove-GlobalBin }
    } else {
        if ($Target -eq "auto")  { foreach ($a in $detected)                               { Set-AgentMdConfig $a } }
        elseif ($Target -eq "all") { foreach ($a in @("claude","codex","gemini","openclaw","opencode")) { Set-AgentMdConfig $a } }
        else                     { Set-AgentMdConfig $Target }
        Install-GlobalBin
    }
}

# ── Footer ────────────────────────────────────────────────────────────────────
$e = [char]27
$Bold = "${e}[1m"; $Dim = "${e}[2m"; $Gold = "${e}[38;5;178m"; $Cyan = "${e}[0;36m"; $Reset = "${e}[0m"
Write-Host ""
Write-Host "${Bold}💡 Also available via marketplace:${Reset}"
Write-Host "   ${Gold}🐙 Copilot${Reset}  ${Cyan}copilot plugin marketplace add jeffreytse/grimoire${Reset}"
Write-Host "              ${Cyan}copilot plugin install grimoire@grimoire${Reset}"
Write-Host "   ${Gold}📖 Docs${Reset}     ${Dim}https://github.com/jeffreytse/grimoire#-install${Reset}"
Write-Host ""
