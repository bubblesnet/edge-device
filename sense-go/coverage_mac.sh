echo $GOPATH
set SAVE_GOOS=$GOOS
set SAVE_GOARCH=$GOARCH
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED="1"
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
set GOOS=$SAVE_GOOS
set GOARCH=$SAVE_GOARCH
