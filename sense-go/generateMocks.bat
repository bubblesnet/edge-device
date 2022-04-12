rem mockgen -build_flags="-mod=vendor" -source /users/rodley/documents/go/src/bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc SensorStoreAndForwardClient -destination mock_bubblesgrpc/hw_mock.go

mockgen -source="bubblesgrpc/bubblesgrpc.pb.go" -self_package bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc > mock_bubblesgrpc/hw_mock.go
