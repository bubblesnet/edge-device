REM protoc -I bubblesgrpc/ bubblesgrpc/bubblesgrpc.proto --go_out=plugins=grpc:bubblesgrpc
python -m grpc_tools.protoc -I../../store-and-forward/bubblesgrpc-server/bubblesgrpc --python_out=./grpc --pyi_out=. --grpc_python_out=./grpc ../../store-and-forward/bubblesgrpc-server/bubblesgrpc/bubblesgrpc.proto