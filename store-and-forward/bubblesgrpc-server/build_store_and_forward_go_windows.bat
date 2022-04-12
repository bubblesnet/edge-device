set GOPATH=%GOPATH%;c:\Users\rodley\documents\go;c:\Users\rodley\go
echo %GOPATH%
set GOOS=windows
set GOARCH=amd64
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
go build -ldflags="-X 'main.BubblesnetBuildNumberString=201' -X 'lmain.BubblesnetVersionMajorString=2' -X 'main.BubblesnetVersionMinorString=1' -X 'main.BubblesnetVersionPatchString=1'  -X 'main.BubblesnetGitHash=%GITHASH%' -X main.BubblesnetBuildTimestamp=%TIMESTAMP%" -o build ./...
