#!/bin/bash
set -e

# Configuration
REPO="tungnguyenlam/language-learning-tui"
BINARY_NAME="deutsch-tui"
INSTALL_DIR="${HOME}/.local/bin"

# Ensure install directory exists
mkdir -p "${INSTALL_DIR}"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "${OS}" in
    linux)  PLATFORM="linux" ;;
    darwin) PLATFORM="macos" ;;
    *)      echo "Unsupported OS: ${OS}"; exit 1 ;;
esac

case "${ARCH}" in
    x86_64)  CPU="amd64" ;;
    arm64|aarch64) CPU="arm64" ;;
    *)       echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
esac

# Note: The current release workflow only builds specific combinations.
# Adjusting to match the release naming convention in .github/workflows/release.yml
# - linux-amd64
# - macos-arm64
# - windows-amd64.exe (not handled by this bash script)

REMOTE_BINARY="${BINARY_NAME}-${PLATFORM}-${CPU}"

# Handle Windows (if someone runs this in Git Bash/WSL)
if [[ "${OS}" == *"mingw"* || "${OS}" == *"cygwin"* ]]; then
    echo "For Windows, please download the .exe manually or use WinGet/Scoop."
    exit 1
fi

echo "Downloading ${REMOTE_BINARY}..."

# Get latest release tag
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not find latest release for ${REPO}"
    exit 1
fi

URL="https://github.com/$(echo ${REPO} | sed 's/ /%20/g')/releases/download/${LATEST_RELEASE}/${REMOTE_BINARY}"

# Download
curl -L -o "${BINARY_NAME}" "${URL}"
chmod +x "${BINARY_NAME}"

# Install
mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
echo "Successfully installed to ${INSTALL_DIR}/${BINARY_NAME}"

# Path check
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo ""
    echo "⚠️  Note: ${INSTALL_DIR} is not in your PATH."
    echo "To run ${BINARY_NAME} from anywhere, add it to your shell config (e.g., ~/.bashrc or ~/.zshrc):"
    echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
    echo ""
fi

echo "You can now run '${INSTALL_DIR}/${BINARY_NAME}'"
