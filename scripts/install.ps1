#Requires -Version 5.1
[CmdletBinding()]
param(
    [string]$Domain    = "",
    [string]$Subdomain = "",
    [string]$Skill     = "",
    [string]$Target    = "auto",
    [switch]$List,
    [switch]$Help
)

$RepoRoot        = Split-Path -Parent $PSScriptRoot
$SkillsRoot      = Join-Path $RepoRoot "skills"
$ClaudeSkillsDir = Join-Path $HOME ".claude\skills"
$AgentsSkillsDir = Join-Path $HOME ".agents\skills"
$GeminiSkillsDir = Join-Path $HOME ".gemini\skills"

function Show-Usage {
    Write-Host @"
Usage: install.ps1 [OPTIONS]

Options:
  -Domain <name>       Install all skills for a domain (e.g. engineering, photography)
  -Subdomain <name>    Restrict to one sub-domain within a domain (e.g. development)
  -Skill <path>        Install one skill (e.g. engineering/development/propose-conventional-commit)
  -Target <agent>      Target: auto (default), claude, codex, gemini, all
  -List                List available domains, sub-domains, and skills
  -Help                Show this help

Examples:
  .\install.ps1
  .\install.ps1 -Domain engineering
  .\install.ps1 -Domain engineering -Subdomain development
  .\install.ps1 -Skill engineering/development/propose-conventional-commit
  .\install.ps1 -Domain engineering -Target all
"@
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
            foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
                if (Test-Path (Join-Path $sd.FullName "SKILL.md")) {
                    Write-Host "  $($domainDir.Name)/$($sd.Name)"
                }
            }
        }
    }
}

function Install-SkillDir([string]$Src, [string]$DestDir) {
    $skillName = Split-Path -Leaf $Src
    $dest = Join-Path $DestDir $skillName
    New-Item -ItemType Directory -Force -Path $dest | Out-Null
    Copy-Item -Path "$Src\*" -Destination $dest -Recurse -Force
    Write-Host "  installed: $skillName -> $dest"
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
        "auto" {
            if (Test-Path (Join-Path $HOME ".claude"))  { Install-SkillDir $Src $ClaudeSkillsDir }
            if (Test-Path (Join-Path $HOME ".agents"))  { Install-SkillDir $Src $AgentsSkillsDir }
            if (Test-Path (Join-Path $HOME ".gemini"))  { Install-SkillDir $Src $GeminiSkillsDir }
        }
    }
}

function Install-Subdomain([string]$SubDir, [string]$TargetAgent) {
    $sp = Join-Path $SubDir "skills"
    if (-not (Test-Path $sp)) { return }
    $found = $false
    foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
        if (-not (Test-Path (Join-Path $sd.FullName "SKILL.md"))) { continue }
        if (-not $found) {
            Write-Host "  Installing sub-domain: $(Split-Path -Leaf $SubDir)"
            $found = $true
        }
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
        foreach ($sd in Get-ChildItem $sp -Directory | Sort-Object Name) {
            if (Test-Path (Join-Path $sd.FullName "SKILL.md")) {
                Invoke-Install $sd.FullName $TargetAgent
            }
        }
    }
}

# Main
if ($Help) { Show-Usage; exit 0 }
if ($List) { Get-SkillList; exit 0 }

if ($Target -eq "auto") {
    $detected = Get-DetectedAgents
    if ($detected.Count -eq 0) {
        Write-Host "No agents detected (~\.claude, ~\.agents, ~\.gemini). Defaulting to Claude Code."
        $Target = "claude"
    } else {
        Write-Host "Detected agents: $($detected -join ', ')"
    }
}

if ($Skill) {
    $parts = $Skill -split "/"
    if ($parts.Count -eq 3) {
        $skillPath = Join-Path $SkillsRoot "$($parts[0])\$($parts[1])\skills\$($parts[2])"
    } elseif ($parts.Count -eq 2) {
        $skillPath = Join-Path $SkillsRoot "$($parts[0])\skills\$($parts[1])"
    } else {
        Write-Host "Invalid skill path: $Skill"; exit 1
    }
    if (-not (Test-Path (Join-Path $skillPath "SKILL.md"))) {
        Write-Host "Skill not found: $Skill"; exit 1
    }
    Write-Host "Installing skill: $Skill"
    Invoke-Install $skillPath $Target
} elseif ($Domain) {
    Install-Domain $Domain $Subdomain $Target
} else {
    foreach ($domainDir in Get-ChildItem $SkillsRoot -Directory | Sort-Object Name) {
        if ($domainDir.Name.StartsWith(".")) { continue }
        Install-Domain $domainDir.Name "" $Target
    }
}

Write-Host "Done."
Write-Host ""
Write-Host "Note: Cursor and GitHub Copilot CLI require marketplace install (not handled by this script):"
Write-Host "  Cursor:  /add-plugin grimoire  (in Agent chat)"
Write-Host "  Copilot: copilot plugin install grimoire"
Write-Host "  See: https://github.com/jeffreytse/grimoire#-install"
