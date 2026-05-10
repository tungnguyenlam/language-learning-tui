# deutsch-tui demo — requires https://github.com/charmbracelet/vhs , ffmpeg, and ttyd on PATH
# From repository root: vhs demo.tape  →  writes assets/demo.gif (120×32 terminal, 24 fps)

Output assets/demo.gif

Set Width 120
Set Height 32
Set Framerate 24

Require go

Set Shell "bash"
Set TypingSpeed 30ms

Hide
Type "rm -rf /tmp/deutsch-tui-vhs && mkdir -p /tmp/deutsch-tui-vhs"
Enter
Sleep 400ms
Show

Type "go run ./cmd/deutsch-tui -data-dir /tmp/deutsch-tui-vhs"
Enter

Sleep 4s

Type "2"
Sleep 900ms
Type "c"
Sleep 3s

Enter
Sleep 2s

Space
Sleep 1500ms
Type "g"
Sleep 1500ms

Escape
Sleep 1s

Type "8"
Sleep 3s

Type "/"
Sleep 400ms
Type "der"
Sleep 600ms
Enter
Sleep 2s

Type "q"
Sleep 1500ms
