$ErrorActionPreference = "Stop"
$Root = Join-Path ([IO.Path]::GetTempPath()) ("courier-installer-test-" + [Guid]::NewGuid().ToString("N"))
$Original = @{
    COURIER_VERSION = $env:COURIER_VERSION
    COURIER_RELEASE_BASE_URL = $env:COURIER_RELEASE_BASE_URL
    COURIER_INSTALL_DIR = $env:COURIER_INSTALL_DIR
    COURIER_INSTALL_NO_PATH_UPDATE = $env:COURIER_INSTALL_NO_PATH_UPDATE
}

try {
    $Fixture = Join-Path $Root "releases\v1.2.3"
    $Install = Join-Path $Root "install"
    New-Item -ItemType Directory -Force -Path $Fixture | Out-Null
    $Architecture = [Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
    $TargetArch = if ($Architecture -eq "x64") { "amd64" } elseif ($Architecture -eq "arm64") { "arm64" } else { throw "unsupported test architecture" }
    $Asset = "courier_1.2.3_windows_$TargetArch.exe"
    $AssetPath = Join-Path $Fixture $Asset
    [IO.File]::WriteAllBytes($AssetPath, [Text.Encoding]::UTF8.GetBytes("standalone-binary"))
    $Digest = (Get-FileHash -Algorithm SHA256 -Path $AssetPath).Hash.ToLowerInvariant()
    [IO.File]::WriteAllText((Join-Path $Fixture "checksums.txt"), "$Digest  $Asset`n")

    $env:COURIER_VERSION = "1.2.3"
    $ReleasePath = (Join-Path $Root "releases").Replace('\', '/').TrimStart('/')
    $env:COURIER_RELEASE_BASE_URL = "file:///$ReleasePath"
    $env:COURIER_INSTALL_DIR = $Install
    $env:COURIER_INSTALL_NO_PATH_UPDATE = "1"
    & (Join-Path $PSScriptRoot "..\install.ps1")
    $Installed = [Text.Encoding]::UTF8.GetString([IO.File]::ReadAllBytes((Join-Path $Install "courier.exe")))
    if ($Installed -ne "standalone-binary") {
        throw "PowerShell installer wrote unexpected content"
    }
}
finally {
    foreach ($Name in $Original.Keys) {
        if ($null -eq $Original[$Name]) {
            Remove-Item -Path "Env:$Name" -ErrorAction SilentlyContinue
        }
        else {
            Set-Item -Path "Env:$Name" -Value $Original[$Name]
        }
    }
    Remove-Item -LiteralPath $Root -Recurse -Force -ErrorAction SilentlyContinue
}
