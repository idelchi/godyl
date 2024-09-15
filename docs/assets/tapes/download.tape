Output gifs/download.gif

Set FontSize 16
Set Width 1600
Set Height 300

# Setup
Hide
Type "rm -rf .bin-*; godyl cache rm; clear"
Enter
Show

Type@50ms 'godyl download "https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz" --version v3.17.3'

Sleep 500ms

Enter

Wait />/

Sleep 1s

Type "find . -name 'helm'"
Sleep 500ms
Enter

Sleep 1s

Hide
Type "rm -rf .bin-*; godyl cache rm; clear"
Enter
