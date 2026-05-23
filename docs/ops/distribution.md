# Distribution Strategy

To provide the best user experience, we aim to publish `deutsch-tui` to major package managers.

## Recommended Tool: GoReleaser

[GoReleaser](https://goreleaser.com/) is the industry standard for shipping Go binaries. It can automate:
- Cross-compiling for all platforms.
- Creating GitHub Releases with changelogs.
- Generating Homebrew Formulae.
- Generating Scoop/WinGet manifests.
- Generating `.deb` and `.rpm` packages.

### Steps to Implement GoReleaser

1.  **Initialize GoReleaser:**
    ```bash
    goreleaser init
    ```
2.  **Configure `.goreleaser.yaml`:**
    Define your builds, homebrew taps, and nfpms (for deb/rpm).
3.  **Update GitHub Action:**
    Replace the manual build steps in `.github/workflows/release.yml` with the GoReleaser action:
    ```yaml
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v5
      with:
        distribution: goreleaser
        version: latest
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    ```

## Target Package Managers

### 🍺 Homebrew (macOS/Linux)
- **Goal:** `brew install tungnguyenlam/tap/deutsch-tui`
- **Method:** Create a "tap" repository (e.g., `homebrew-tap`) and configure GoReleaser to push the formula there.

### 🪟 WinGet (Windows)
- **Goal:** `winget install deutsch-tui`
- **Method:** Submit a manifest to the [Microsoft WinGet Packages](https://github.com/microsoft/winget-pkgs) repository.

### 🐧 APT (Debian/Ubuntu)
- **Goal:** `apt install deutsch-tui`
- **Method:** Use a service like [Cloudsmith](https://cloudsmith.io/) or [Gemfury](https://fury.io/) to host a private/public APT repository, or submit to Debian if the project reaches high maturity.
