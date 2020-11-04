set GOPATH=%GOPATH%;c:\Users\rodley\documents\go\src;c:\Users\rodley\go
echo %GOPATH%
set GOOS=linux
set GOARCH=arm
set GOARM=7
set CGO_ENABLED="1"
go build -o sense_go main.go atlasezo_driver.go adc.go grpc.go powerstrip.go
copy sense_go ..\..\..\bubbles2\sense-go
cd ..\..\..\bubbles2
git commit -a -m "deploy by script"
git push rpi3
git push rpi4
git push balenafin
cd ..\go\src\sense-go