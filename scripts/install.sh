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

# Path check and auto-configuration
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo ""
    echo "⚠️  Note: ${INSTALL_DIR} is not in your PATH."
    
    SHELL_NAME=$(basename "$SHELL")
    CONFIG_FILE=""
    EXPORT_CMD="export PATH=\"\$HOME/.local/bin:\$PATH\""

    case "$SHELL_NAME" in
        "zsh") CONFIG_FILE="${HOME}/.zshrc" ;;
        "bash") CONFIG_FILE="${HOME}/.bashrc" ;;
        "fish") 
            CONFIG_FILE="${HOME}/.config/fish/config.fish"
            EXPORT_CMD="fish_add_path \$HOME/.local/bin"
            ;;
        *) CONFIG_FILE="${HOME}/.profile" ;;
    esac

    if ! grep -q "\.local/bin" "${CONFIG_FILE}" 2>/dev/null; then
        echo "Auto-configuring ${SHELL_NAME} by updating ${CONFIG_FILE}..."
        
        if [ "$SHELL_NAME" = "fish" ]; then
            mkdir -p "${HOME}/.config/fish"
        fi

        echo "" >> "${CONFIG_FILE}"
        echo "# Added by deutsch-tui installer" >> "${CONFIG_FILE}"
        echo "${EXPORT_CMD}" >> "${CONFIG_FILE}"
        
        echo "✅ Added to ${CONFIG_FILE}."
        echo "To use '${BINARY_NAME}' immediately, run: source ${CONFIG_FILE} (or restart your terminal)"
    else
        echo "ℹ️  It looks like ${INSTALL_DIR} is already in ${CONFIG_FILE}, but hasn't been loaded in this session."
        echo "To use '${BINARY_NAME}' immediately, run: source ${CONFIG_FILE} (or restart your terminal)"
    fi
else
    echo "You can now run '${BINARY_NAME}'"
fi
