Output gifs/download.gif

Set Theme "CLRS"
Set FontSize 16
Set Width 1600
Set Height 400
Set TypingSpeed 25ms

# Setup
Hide
Type "rm -rf ./bin; godyl cache rm; export GODYL_DOWNLOAD_OUTPUT=bin; clear"
Enter
Show

Type 'godyl -v download "https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz" --version v3.17.3'

Sleep 500ms

Enter

Wait@10s

Sleep 3s

Hide
Type "rm -rf bin; godyl cache rm; clear"
Enter
