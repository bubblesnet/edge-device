set GOPATH=%GOPATH%;c:\Users\rodley\documents\go;c:\Users\rodley\go
echo %GOPATH%

set API_HOST = 192.168.23.237
set NO_FAN_WITH_HEATER = false
set SLEEP_ON_EXIT_FOR_DEBUGGING = 60
set ACTIVEMQ_HOST = 192.168.23.237
set ACTIVEMQ_PORT  = 61611
set API_PORT = 4001
set USERID = 90000009
set DEVICEID = 70000008

set GOOS=windows
set GOARCH=amd64
set GOARM=7
set CGO_ENABLED="1"
set GITHASH=""
set TIMESTAMP=""
for /f "tokens=*" %%a in ('git rev-parse HEAD') do (
    set GITHASH=%%a
)
for /f "tokens=*" %%a in ('date /t') do (
    set TIMESTAMP='%%a
)
for /f "tokens=*" %%a in ('time /t') do (
    set TIMESTAMP=%TIMESTAMP% %%a'
)
set TIMESTAMP=%TIMESTAMP: =_%
go build -ldflags="-X 'main.BubblesnetBuildNumberString=201' -X 'main.BubblesnetVersionMajorString=2' -X 'main.BubblesnetVersionMinorString=1' -X 'main.BubblesnetVersionPatchString=1'  -X 'main.BubblesnetGitHash=%GITHASH%' -X main.BubblesnetBuildTimestamp=%TIMESTAMP%" -o build ./...
go test ./...