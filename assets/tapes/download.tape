Output assets/gifs/download.gif

Set Theme "CLRS"
Set FontSize 16
Set Width 1600
Set Height 500
Set TypingSpeed 25ms

Env GODYL_TOOLS "/home/user/assets/tapes/demo.yml"
Env GODYL_DOWNLOAD_OUTPUT "/home/user/.local/bin"
Env GODYL_VERBOSE "1"
Env GODYL_ENV_FILE "/home/user/assets/tokens.env"

Type 'godyl download derailed/k9s'

Sleep 500ms

Enter

Wait@10s

Sleep 3s

Hide
Type "rm -rf bin; godyl cache rm; clear"
Enter
