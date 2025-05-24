Output gifs/install.gif

Set Theme "CLRS"
Set FontSize 16
Set Width 1600
Set Height 600
Set TypingSpeed 50ms

# Setup
Hide
Type "rm -rf ./bin; godyl cache rm; export GODYL_INSTALL_OUTPUT=bin; export GODYL_INSTALL_TOOLS=/home/user/tapes/demo.yml; clear"
Enter
Show

Type "godyl -v install --tags '!go'"

Sleep 500ms

Enter

Wait@10s

Sleep 3s

Hide
Type "rm -rf bin; godyl cache rm; clear"
Enter
