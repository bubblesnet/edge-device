echo $GOPATH
export CGO_ENABLED="1"
export CODECOV_TOKEN='bd6757f7-5f19-40b6-81f3-68547d5b9177'
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
