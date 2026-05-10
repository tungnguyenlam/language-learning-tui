# deutsch-tui · charmbracelet/vhs (needs ffmpeg + ttyd on PATH)
# Regenerate demo: vhs demo.tape   →   assets/demo.gif
# Install VHS: brew install vhs   OR   go install github.com/charmbracelet/vhs@latest

Output assets/demo.gif

Set Shell "bash"
Set FontFamily "JetBrains Mono"
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Framerate 24
Set PlaybackSpeed 0.8
Set TypingSpeed 35ms

Require go

Hide
Type "rm -rf /tmp/deutsch-tui-vhs && mkdir -p /tmp/deutsch-tui-vhs"
Enter
Sleep 400ms
Show

Type "go run ./cmd/deutsch-tui -data-dir /tmp/deutsch-tui-vhs"
Sleep 500ms
Enter

Sleep 3s

Type "2"
Sleep 1s

Down
Sleep 500ms
Enter
Sleep 1s

Type "3"
Sleep 1s

Sleep 2s

Space
Sleep 1s

Type "g"
Sleep 500ms

Type "8"
Sleep 1s

Type "/"
Sleep 500ms
Type "vocabulary"
Sleep 1s
Enter
Sleep 1s

Type "q"
Sleep 500ms
