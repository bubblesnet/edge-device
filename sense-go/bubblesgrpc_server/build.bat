
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o bubblesgrpc_server ./...
copy bubblesgrpc_server ..\..\..\..\bubbles2\bubblesgrpc-server\