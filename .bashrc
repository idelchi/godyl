source .env
alias g="godyl --log debug"
alias gr="go run ./cmd/godyl -t trntv/sshed --log=debug"
alias gr2="go run ./cmd/godyl -t jqlang/jq --log=debug"
alias gg="godyl -t trntv/sshed --log=debug"
alias ggg="godyl --dry --log debug -c cmd/godyl/config.yml"


echo "done!"
