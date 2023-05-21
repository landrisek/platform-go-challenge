export GO111MODULE="on"
export GOPATH="$HOME/go"

go test -count=1 -tags=unit ./internal/controller -v
go test -count=1 -tags=unit ./internal/models -v
go test -count=1 -tags=unit ./internal/repository -v
go test -count=1 -tags=unit ./internal/sagas -v