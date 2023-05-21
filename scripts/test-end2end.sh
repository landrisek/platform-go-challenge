export GO111MODULE="on"
export GOPATH="$HOME/go"

go test -v -count=1 -tags=end2end ./internal/controller