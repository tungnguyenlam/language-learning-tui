$ErrorActionPreference = "Stop"

$Repo = "tungnguyenlam/language-learning-tui"
$BinaryName = "deutsch-tui.exe"
$InstallDir = "$env:USERPROFILE\.local\bin"

# Create install directory
if (-not (Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Host "Fetching latest release from $Repo..."

# Get latest release tag
$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
try {
    $Release = Invoke-RestMethod -Uri $ApiUrl
    $LatestTag = $Release.tag_name
} catch {
    Write-Error "Failed to fetch latest release from GitHub API."
    exit 1
}

if (-not $LatestTag) {
    Write-Error "Could not find the latest release."
    exit 1
}

$RemoteBinary = "deutsch-tui-windows-amd64.exe"
$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestTag/$RemoteBinary"

Write-Host "Downloading $RemoteBinary ($LatestTag)..."

$Destination = Join-Path -Path $InstallDir -ChildPath $BinaryName
Invoke-WebRequest -Uri $DownloadUrl -OutFile $Destination

Write-Host "Successfully installed to $Destination"

# Add to PATH if necessary
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notmatch [regex]::Escape($InstallDir)) {
    Write-Host "Adding $InstallDir to your User PATH environment variable..."
    $NewPath = "$InstallDir;$UserPath"
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
    
    # Update current session PATH so it works immediately
    $env:PATH = "$InstallDir;$env:PATH"
    
    Write-Host "✅ Added to PATH."
} else {
    Write-Host "ℹ️ $InstallDir is already in your PATH."
}

Write-Host ""
Write-Host "You can now run 'deutsch-tui' from your terminal."
Write-Host "(Note: You may need to restart your terminal if the command is not recognized)"
