$RepoOwner = "rasel9t6"
$RepoName = "the-go-engineer"
$ConfigPath = "scripts/audit-issues-config.json"
$EnvPath = ".env"

# 1. Load Token from .env
if (Test-Path $EnvPath) {
    $EnvContent = Get-Content $EnvPath
    foreach ($line in $EnvContent) {
        if ($line -like "GITHUB_TOKEN=*") {
            $Global:GithubToken = $line.Split("=")[1].Trim()
        }
    }
}

if (-not $Global:GithubToken) {
    Write-Error "GITHUB_TOKEN not found in .env"
    exit 1
}

# 2. Load Config
if (-not (Test-Path $ConfigPath)) {
    Write-Error "Config file not found: $ConfigPath"
    exit 1
}

$Config = Get-Content $ConfigPath | ConvertFrom-Json
$Headers = @{
    "Authorization" = "token $Global:GithubToken"
    "Accept"        = "application/vnd.github.v3+json"
}

Write-Host "Starting GitHub Automation for $($Config.owner)/$($Config.repo)"
Write-Host "---------------------------------------------------"

# 3. Create Labels
Write-Host "Step 1: Creating/Verifying Labels ($($Config.labels.Count))..."
foreach ($label in $Config.labels) {
    $uri = "https://api.github.com/repos/$($Config.owner)/$($Config.repo)/labels"
    $body = @{
        name        = $label.name
        color       = $label.color
        description = $label.description
    } | ConvertTo-Json

    try {
        $resp = Invoke-RestMethod -Uri $uri -Method Post -Headers $Headers -Body $body -ContentType "application/json"
        Write-Host "  [OK] Created Label: $($label.name)"
    } catch {
        if ($_.Exception.Response.StatusCode -eq 422) {
            Write-Host "  [SKIP] Label already exists: $($label.name)"
        } else {
            Write-Host "  [FAIL] Failed Label: $($label.name) ($($_.Exception.Message))"
        }
    }
}

# 4. Create Issues
Write-Host "`nStep 2: Creating Issues ($($Config.issues.Count))..."
$count = 1
foreach ($issue in $Config.issues) {
    $uri = "https://api.github.com/repos/$($Config.owner)/$($Config.repo)/issues"
    $body = @{
        title     = $issue.title
        body      = $issue.body
        labels    = $issue.labels
        assignees = $issue.assignees
    } | ConvertTo-Json

    try {
        $resp = Invoke-RestMethod -Uri $uri -Method Post -Headers $Headers -Body $body -ContentType "application/json"
        Write-Host "  [OK] Created Issue #$($count): $($issue.title)"
    } catch {
        Write-Host "  [FAIL] Failed Issue #$($count): $($issue.title) ($($_.Exception.Message))"
    }
    $count++
}

Write-Host "---------------------------------------------------"
Write-Host "✨ Task complete!"
