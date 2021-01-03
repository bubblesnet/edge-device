package main

import (
	"bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc"
	"context"
	"reflect"
	"testing"
)

func Test_forwardMessages(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := forwardMessages(tt.args.bucketName); (err != nil) != tt.wantErr {
				t.Errorf("forwardMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseMessage(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("parseMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_saveState(t *testing.T) {
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveState(tt.args.bucketName); (err != nil) != tt.wantErr {
				t.Errorf("saveState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_GetRecordList(t *testing.T) {
	type fields struct {
		UnimplementedSensorStoreAndForwardServer bubblesgrpc.UnimplementedSensorStoreAndForwardServer
	}
	type args struct {
		ctx context.Context
		in  *pb.GetRecordListRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetRecordListReply
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				UnimplementedSensorStoreAndForwardServer: tt.fields.UnimplementedSensorStoreAndForwardServer,
			}
			got, err := s.GetRecordList(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecordList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_GetState(t *testing.T) {
	type fields struct {
		UnimplementedSensorStoreAndForwardServer bubblesgrpc.UnimplementedSensorStoreAndForwardServer
	}
	type args struct {
		ctx context.Context
		in  *pb.GetStateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetStateReply
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				UnimplementedSensorStoreAndForwardServer: tt.fields.UnimplementedSensorStoreAndForwardServer,
			}
			got, err := s.GetState(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetState() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_StoreAndForward(t *testing.T) {
	type fields struct {
		UnimplementedSensorStoreAndForwardServer bubblesgrpc.UnimplementedSensorStoreAndForwardServer
	}
	type args struct {
		ctx context.Context
		in  *pb.SensorRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.SensorReply
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				UnimplementedSensorStoreAndForwardServer: tt.fields.UnimplementedSensorStoreAndForwardServer,
			}
			got, err := s.StoreAndForward(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreAndForward() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreAndForward() got = %v, want %v", got, tt.want)
			}
		})
	}
}