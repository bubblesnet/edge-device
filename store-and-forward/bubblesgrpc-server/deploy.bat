set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o bubblesgrpc_server main.go
copy bubblesgrpc_server ..\..\..\..\bubbles2\bubblesgrpc-server\
cd ..\..\..\..\bubbles2
git commit -a -m "deploy by script"
git push rpi3
cd ..\go\src\bubblesgrpc\bubblesgrpc_server