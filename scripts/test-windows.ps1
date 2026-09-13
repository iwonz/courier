$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$Repository = Resolve-Path (Join-Path $PSScriptRoot "..")
$TemporaryDirectory = Join-Path ([IO.Path]::GetTempPath()) ("courier-runtime-" + [Guid]::NewGuid().ToString("N"))
$Coverage = Join-Path $TemporaryDirectory "coverage.out"
$Binary = Join-Path $TemporaryDirectory "courier.exe"
New-Item -ItemType Directory -Force -Path $TemporaryDirectory | Out-Null

Push-Location $Repository
try {
    & go test ./... -covermode=atomic "-coverprofile=$Coverage"
    if ($LASTEXITCODE -ne 0) {
        throw "Windows Go tests failed with exit code $LASTEXITCODE"
    }

    $CoverageOutput = & go tool cover "-func=$Coverage"
    if ($LASTEXITCODE -ne 0) {
        throw "Windows coverage analysis failed with exit code $LASTEXITCODE"
    }
    $TotalLine = $CoverageOutput | Select-String '^total:' | Select-Object -Last 1
    if ($null -eq $TotalLine) {
        throw "Windows coverage output does not contain a total"
    }
    $Total = ($TotalLine.ToString() -split '\s+')[-1]
    if ($Total -ne "100.0%") {
        throw "Windows statement coverage is $Total, expected 100.0%"
    }

    & go build -trimpath -o $Binary ./cmd/courier
    if ($LASTEXITCODE -ne 0) {
        throw "Windows Courier build failed with exit code $LASTEXITCODE"
    }
    $Help = & $Binary help
    if ($LASTEXITCODE -ne 0 -or -not ($Help -join "`n").Contains("Safely transfer files and directories")) {
        throw "Compiled Windows Courier help check failed"
    }
    $Version = & $Binary version
    if ($LASTEXITCODE -ne 0 -or -not ($Version -join "`n").Contains("courier dev")) {
        throw "Compiled Windows Courier version check failed"
    }
}
finally {
    Pop-Location
    Remove-Item -LiteralPath $TemporaryDirectory -Recurse -Force -ErrorAction SilentlyContinue
}
