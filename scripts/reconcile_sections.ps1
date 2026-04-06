$dirs10 = "10-web-and-database"
$dirs11 = "11-concurrency"

# Update Section 10 files
Get-ChildItem -Path $dirs10 -Recurse -Filter "*.go" | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    $newContent = $content -replace "Section (12|13|16|21):", "Section 10:"
    if ($content -ne $newContent) {
        $newContent | Set-Content $_.FullName
        Write-Host "Updated $($_.FullName)"
    }
}

# Update Section 11 files
Get-ChildItem -Path $dirs11 -Recurse -Filter "*.go" | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    $newContent = $content -replace "Section (9|11|15|17|24):", "Section 11:"
    if ($content -ne $newContent) {
        $newContent | Set-Content $_.FullName
        Write-Host "Updated $($_.FullName)"
    }
}
