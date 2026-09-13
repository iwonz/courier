$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$Repository = Resolve-Path (Join-Path $PSScriptRoot "..")
$Coverage = Join-Path $Repository "coverage-windows.out"

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
}
finally {
    Pop-Location
    Remove-Item -LiteralPath $Coverage -Force -ErrorAction SilentlyContinue
}
