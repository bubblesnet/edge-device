# export GOPATH=/Users/rodley/go
# env GOOS=linux GOARCH=arm GOARM=7 go build -o build ./...

export GOPATH=/Users/rodley/go
echo $GOPATH
set GOOS=linux
set GOARCH=arm
set GOARM=7
set CGO_ENABLED="1"
set GITHASH=""
set TIMESTAMP=""


set TIMESTAMP=%TIMESTAMP: =_%
env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-X 'main.BubblesnetBuildNumberString=201' -X 'main.BubblesnetVersionMajorString=2' -X 'main.BubblesnetVersionMinorString=1' -X 'main.BubblesnetVersionPatchString=1'  -X 'main.BubblesnetGitHash=$GITHASH' -X main.BubblesnetBuildTimestamp=$TIMESTAMP" -o build ./...
