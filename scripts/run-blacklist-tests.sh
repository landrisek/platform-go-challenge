export GO111MODULE="on"
export GOPATH="$HOME/go"

# run both unit and integration tests for blacklist microservice
go test -v -run "TestBlacklist|TestHandleConnection" ./internal/controller
go test -v -run TestFindBlacklist ./internal/repository
