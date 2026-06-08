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
$ClaudeSkillsDir = Join-Path $HOME ".claude\skills"
$AgentsSkillsDir = Join-Path $HOME ".agents\skills"
$GeminiSkillsDir = Join-Path $HOME ".gemini\skills"

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
Usage: install.ps1 [OPTIONS]

Options:
  -Domain <name>       Install/uninstall all skills for a domain
  -Subdomain <name>    Restrict to one sub-domain within a domain
  -Skill <path>        Install/uninstall one skill (e.g. engineering/development/propose-conventional-commit)
  -Target <agent>      Target: claude, codex, gemini, all
  -Uninstall           Remove skills instead of installing
  -Copy                Use copy mode instead of junctions
  -Upgrade             Pull latest grimoire at GrimoireHome (junctions update automatically)
  -Clean               Remove broken junctions from all agent skill dirs
  -Yes                 Non-interactive: install all skills to all detected agents
  -List                List available domains, sub-domains, and skills
  -Help                Show this help

Environment:
  GRIMOIRE_HOME        Persistent clone location (default: ~\.grimoire)

Examples:
  .\install.ps1                                                   # Interactive TUI
  .\install.ps1 -Yes                                              # Install everything, no prompts
  .\install.ps1 -Upgrade                                          # Pull latest grimoire
  .\install.ps1 -Clean                                            # Remove broken junctions
  .\install.ps1 -Domain engineering -Target claude
  .\install.ps1 -Domain engineering -Copy                         # Copy instead of junction
  .\install.ps1 -Skill engineering/development/propose-conventional-commit
  .\install.ps1 -Uninstall -Domain engineering -Target claude
  .\install.ps1 -Uninstall -Skill engineering/development/propose-conventional-commit
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
    if (Test-Path (Join-Path $HOME ".claude"))  { $d += "claude" }
    if (Test-Path (Join-Path $HOME ".agents"))  { $d += "codex" }
    if (Test-Path (Join-Path $HOME ".gemini"))  { $d += "gemini" }
    return $d
}

function Invoke-Install([string]$Src, [string]$TargetAgent) {
    switch ($TargetAgent) {
        "claude" { Install-SkillDir $Src $ClaudeSkillsDir }
        "codex"  { Install-SkillDir $Src $AgentsSkillsDir }
        "gemini" { Install-SkillDir $Src $GeminiSkillsDir }
        "all" {
            Install-SkillDir $Src $ClaudeSkillsDir
            Install-SkillDir $Src $AgentsSkillsDir
            Install-SkillDir $Src $GeminiSkillsDir
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
    $dirs = @($ClaudeSkillsDir, $AgentsSkillsDir, $GeminiSkillsDir)
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
        "claude" { Uninstall-SkillDir $SkillName $ClaudeSkillsDir }
        "codex"  { Uninstall-SkillDir $SkillName $AgentsSkillsDir }
        "gemini" { Uninstall-SkillDir $SkillName $GeminiSkillsDir }
        "all" {
            Uninstall-SkillDir $SkillName $ClaudeSkillsDir
            Uninstall-SkillDir $SkillName $AgentsSkillsDir
            Uninstall-SkillDir $SkillName $GeminiSkillsDir
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
if ($Help) { Show-Usage; exit 0 }
if ($List) { Get-SkillList; exit 0 }

if ($Upgrade) {
    if (-not (Test-Path (Join-Path $GrimoireHome ".git"))) {
        Write-Host "No grimoire clone found at $GrimoireHome. Run install first."
        exit 1
    }
    Write-Host "Pulling latest grimoire at $GrimoireHome..."
    git -C $GrimoireHome pull
    Write-Host "Cleaning broken junctions..."
    Invoke-Clean
    Write-Host "Done."
    exit 0
}

if ($Clean -and -not $Uninstall -and -not $Domain -and -not $Skill) {
    Write-Host "Cleaning broken junctions..."
    Invoke-Clean
    exit 0
}

$detected    = @()
$isInteractive = (-not $Domain -and -not $Skill -and -not $Target -and -not $Yes -and -not $Uninstall)

if ($isInteractive) {
    Print-Banner

    $tuiMode = Select-One "⚙️  What would you like to do?" @("📥 Install", "🗑️  Uninstall")

    $detected = Get-DetectedAgents
    if ($detected.Count -eq 0) {
        Write-Host "No agents detected. Defaulting to Claude Code."
        $agentList = @("claude")
    } else {
        if ($tuiMode -like "*Uninstall*") { $agentList = Invoke-Multiselect "🤖 Which agents to uninstall from?" $detected }
        else                              { $agentList = Invoke-Multiselect "🤖 Which agents to install to?"    $detected }
        if ($agentList.Count -eq 0) { Write-Host "No agents selected. Exiting."; exit 0 }
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
}

# ── Footer ────────────────────────────────────────────────────────────────────
$e = [char]27
$Bold = "${e}[1m"; $Dim = "${e}[2m"; $Gold = "${e}[38;5;178m"; $Cyan = "${e}[0;36m"; $Reset = "${e}[0m"
Write-Host ""
Write-Host "${Bold}💡 Also available via marketplace:${Reset}"
Write-Host "   ${Gold}🖱️  Cursor${Reset}   ${Cyan}/add-plugin grimoire${Reset}  ${Dim}(in Agent chat)${Reset}"
Write-Host "   ${Gold}🐙 Copilot${Reset}  ${Cyan}copilot plugin install grimoire${Reset}"
Write-Host "   ${Gold}📖 Docs${Reset}     ${Dim}https://github.com/jeffreytse/grimoire#-install${Reset}"
Write-Host ""
