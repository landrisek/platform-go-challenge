export GO111MODULE="on"
export GOPATH="$HOME/go"

go test -v ./internal/controller -run TestServer