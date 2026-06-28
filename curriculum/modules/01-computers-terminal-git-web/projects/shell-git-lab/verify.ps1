param(
    [string]$ProjectDir = "project"
)

$ProjectRoot = Resolve-Path $ProjectDir -ErrorAction Stop
$pass = 0
$fail = 0
$total = 14

function Check {
    param([string]$Name, [scriptblock]$Condition)
    $script:pass++
    try {
        $result = & $Condition
        if ($result) { Write-Host "  PASS  $Name" -ForegroundColor Green }
        else { throw "assertion failed" }
    } catch {
        $script:pass--
        $script:fail++
        Write-Host "  FAIL  $Name" -ForegroundColor Red
    }
}

Write-Host "`nShell + Git Lab Verification`n" -ForegroundColor Cyan

# --- Shell checks ---
Check "project directory exists" { Test-Path $ProjectRoot }
Check "src subdirectory exists" { Test-Path (Join-Path $ProjectRoot "src") }
Check "docs subdirectory exists" { Test-Path (Join-Path $ProjectRoot "docs") }
Check "tests subdirectory exists" { Test-Path (Join-Path $ProjectRoot "tests") }
Check "README.md exists" { Test-Path (Join-Path $ProjectRoot "README.md") }
Check "src/main.go exists" { Test-Path (Join-Path $ProjectRoot "src" "main.go") }
Check "docs/design.md exists" { Test-Path (Join-Path $ProjectRoot "docs" "design.md") }

# --- Git checks ---
Push-Location $ProjectRoot
try {
    Check "git repository initialized" { Test-Path ".git" }
    Check "at least 3 commits exist" { (git log --oneline 2>$null).Count -ge 3 }
    Check "feature branch merged" {
        $branches = git branch --merged 2>$null
        $branches -match "feature"
    }
    Check "git log shows merge commit or fast-forward" {
        $log = git log --oneline --graph 2>$null
        $log.Count -ge 3
    }

    # Check README.md was an initial commit (exists at root)
    Check "README.md in git tree" {
        $null -ne (git ls-files "README.md" 2>$null)
    }

    # Check merge conflict marker resolved (no conflict markers in any tracked file)
    Check "no unresolved merge conflicts" {
        $files = git grep -l ">>>>>>>" -- ':/*.md' ':/*.go' 2>$null
        -not $files
    }
}
finally { Pop-Location }

# --- Summary ---
Write-Host "`nResults: $pass/$total passed, $fail failed" -ForegroundColor $(if ($fail -eq 0) { "Green" } else { "Red" })
if ($fail -eq 0) { Write-Host "All checks passed!" -ForegroundColor Green }
else { Write-Host "Some checks failed. Review the output above." -ForegroundColor Yellow }
