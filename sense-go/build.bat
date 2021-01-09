set GOPATH=%GOPATH%;c:\Users\rodley\documents\go\src;c:\Users\rodley\go

echo %GOPATH%
set GOOS=linux
set GOARCH=arm
set GOARM=7
set CGO_ENABLED="1"
go build -o build ./...

rem copy sense_go ..\..\..\bubbles2\sense-go