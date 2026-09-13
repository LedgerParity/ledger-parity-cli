param(
    [string]$Go = 'go',
    [string]$Destination = 'dist'
)
$ErrorActionPreference = 'Stop'
$repo = Split-Path $PSScriptRoot -Parent
$dest = [IO.Path]::GetFullPath((Join-Path $repo $Destination))
if (Test-Path -LiteralPath $dest) { throw 'Release destination must be new to avoid mixing artifacts.' }
$versionLine = Select-String -LiteralPath (Join-Path $repo 'cmd/ledger-parity/main.go') -Pattern '^const Version = "([^"]+)"$'
$version = $versionLine.Matches[0].Groups[1].Value
if ($version -notmatch '^\d+\.\d+\.\d+-preview$') { throw 'Expected an explicit preview version.' }
$oldGOOS = $env:GOOS
$oldGOARCH = $env:GOARCH
$oldCGO = $env:CGO_ENABLED
Push-Location $repo
try {
    New-Item -ItemType Directory -Path $dest | Out-Null
    foreach ($target in @('windows/amd64', 'linux/amd64', 'linux/arm64', 'darwin/amd64', 'darwin/arm64')) {
        $env:GOOS, $env:GOARCH = $target.Split('/')
        $env:CGO_ENABLED = '0'
        $name = "ledger-parity-$version-$($env:GOOS)-$($env:GOARCH)"
        $stage = Join-Path $dest $name
        New-Item -ItemType Directory -Path $stage | Out-Null
        $exe = if ($env:GOOS -eq 'windows') { 'ledger-parity.exe' } else { 'ledger-parity' }
        & $Go build -trimpath -o (Join-Path $stage $exe) ./cmd/ledger-parity
        if ($LASTEXITCODE -ne 0) { throw "Build failed for $target" }
        foreach ($item in @('LICENSE', 'README.md', 'RELEASE_NOTES.md', 'examples', 'dashboard', 'docs')) {
            Copy-Item -LiteralPath (Join-Path $repo $item) -Destination $stage -Recurse
        }
        Compress-Archive -Path $stage -DestinationPath (Join-Path $dest "$name.zip")
    }
    $checksums = Get-ChildItem -LiteralPath $dest -Filter '*.zip' | Sort-Object Name | ForEach-Object {
        (Get-FileHash -LiteralPath $_.FullName -Algorithm SHA256).Hash.ToLowerInvariant() + '  ' + $_.Name
    }
    [IO.File]::WriteAllLines((Join-Path $dest 'SHA256SUMS'), $checksums, [Text.UTF8Encoding]::new($false))
    Write-Output "Built $version packages in $dest"
} finally {
    $env:GOOS = $oldGOOS
    $env:GOARCH = $oldGOARCH
    $env:CGO_ENABLED = $oldCGO
    Pop-Location
}
