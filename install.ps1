[CmdletBinding()]
param(
    [string]$Version = $env:COURIER_VERSION,
    [string]$InstallDir = $env:COURIER_INSTALL_DIR
)

$ErrorActionPreference = "Stop"
$Repository = if ($env:COURIER_REPOSITORY) { $env:COURIER_REPOSITORY } else { "iwonz/courier" }
$ReleaseApi = if ($env:COURIER_RELEASE_API_URL) { $env:COURIER_RELEASE_API_URL } else { "https://api.github.com/repos/$Repository/releases/latest" }
$ReleaseBase = if ($env:COURIER_RELEASE_BASE_URL) { $env:COURIER_RELEASE_BASE_URL } else { "https://github.com/$Repository/releases/download" }
if (-not $InstallDir) {
    $InstallDir = Join-Path ([Environment]::GetFolderPath("LocalApplicationData")) "Programs\Courier"
}

$TemporaryDir = Join-Path ([IO.Path]::GetTempPath()) ("courier-install-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $TemporaryDir | Out-Null

function Receive-CourierFile {
    param([string]$Uri, [string]$Destination)
    $Parsed = [Uri]$Uri
    if ($Parsed.IsFile) {
        Copy-Item -LiteralPath $Parsed.LocalPath -Destination $Destination
    }
    else {
        Invoke-WebRequest -Uri $Parsed -OutFile $Destination -Headers @{ "User-Agent" = "courier-installer" }
    }
}

function Receive-CourierJson {
    param([string]$Uri)
    $Parsed = [Uri]$Uri
    if ($Parsed.IsFile) {
        return Get-Content -LiteralPath $Parsed.LocalPath -Raw | ConvertFrom-Json
    }
    return Invoke-RestMethod -Uri $Parsed -Headers @{ "User-Agent" = "courier-installer" }
}

try {
    if (-not $Version -or $Version -eq "latest") {
        $Release = Receive-CourierJson -Uri $ReleaseApi
        $Version = [string]$Release.tag_name
    }
    if ($Version -notmatch '^v?\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$') {
        throw "release version must be vMAJOR.MINOR.PATCH"
    }
    $Tag = if ($Version.StartsWith("v")) { $Version } else { "v$Version" }
    $VersionNumber = $Tag.Substring(1)

    $Architecture = [Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
    switch ($Architecture) {
        "x64" { $TargetArch = "amd64" }
        "arm64" { $TargetArch = "arm64" }
        default { throw "unsupported architecture: $Architecture" }
    }

    $Asset = "courier_${VersionNumber}_windows_${TargetArch}.exe"
    $BinaryPath = Join-Path $TemporaryDir $Asset
    $ChecksumPath = Join-Path $TemporaryDir "checksums.txt"
    Receive-CourierFile -Uri "$ReleaseBase/$Tag/$Asset" -Destination $BinaryPath
    Receive-CourierFile -Uri "$ReleaseBase/$Tag/checksums.txt" -Destination $ChecksumPath

    $Expected = $null
    foreach ($Line in [IO.File]::ReadLines($ChecksumPath)) {
        if ($Line -match '^([0-9A-Fa-f]{64})\s+\*?(.+)$' -and $Matches[2] -eq $Asset) {
            $Expected = $Matches[1].ToLowerInvariant()
            break
        }
    }
    if (-not $Expected) {
        throw "checksums.txt has no entry for $Asset"
    }
    $Actual = (Get-FileHash -Algorithm SHA256 -Path $BinaryPath).Hash.ToLowerInvariant()
    if ($Actual -ne $Expected) {
        throw "checksum mismatch for $Asset"
    }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $Staged = Join-Path $InstallDir (".courier-install-" + [Guid]::NewGuid().ToString("N") + ".exe")
    try {
        Copy-Item -LiteralPath $BinaryPath -Destination $Staged
        Move-Item -LiteralPath $Staged -Destination (Join-Path $InstallDir "courier.exe") -Force
    }
    finally {
        Remove-Item -LiteralPath $Staged -Force -ErrorAction SilentlyContinue
    }

    if ($env:COURIER_INSTALL_NO_PATH_UPDATE -ne "1") {
        $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $PathEntries = @($UserPath -split ';' | Where-Object { $_ })
        if (-not ($PathEntries | Where-Object { $_.TrimEnd('\') -ieq $InstallDir.TrimEnd('\') })) {
            $NewPath = (@($PathEntries) + $InstallDir) -join ';'
            [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        }
        if (-not (($env:Path -split ';') | Where-Object { $_.TrimEnd('\') -ieq $InstallDir.TrimEnd('\') })) {
            $env:Path = "$InstallDir;$env:Path"
        }
    }
    Write-Output "Installed courier $VersionNumber to $(Join-Path $InstallDir 'courier.exe')"
}
finally {
    Remove-Item -LiteralPath $TemporaryDir -Recurse -Force -ErrorAction SilentlyContinue
}
