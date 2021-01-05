package main

import (
	pb "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc"
	log "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/lawg"
	"context"
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const badTestMessage = "dkdlkdkdkdk"
const emptyTestMessage = "{}"

var goodState = state{
	SampleTimestamp:  1609808763,
	SampleTimestampS: "Tue Jan 05 01:06:03 2021 UTC",
	TempF:            82.0,
	Humidity:         67.7,
	Light:            10053,
	DistanceIn:       12,
	Pressure:         1015,
	Ph:               5.9,
	Humidifier:       true,
	Heater:           true,
	HeaterPad:        true,
	GrowLightVeg:     true,
}

func initTests(t *testing.T) {
	log.ConfigureTestLogging("fatal,error,warn,info,debug,", ".", t)
	storeMountPoint := "/config"
	if  runtime.GOOS == "windows"{
		storeMountPoint = "."
		databaseFilename = "./testmessages.db"
	} else if runtime.GOOS == "darwin" {
		storeMountPoint = "."
		databaseFilename = "./testmessages.db"
	}
	_ = ReadFromPersistentStore(storeMountPoint, "", "config.json",&config,&stageSchedule)

	t.Logf("config = %v", config)
	t.Logf("stageSchedule = %v", stageSchedule)
	initDb(databaseFilename)
	t.Logf("returned from initDb")
}

func Test_forwardMessages(t *testing.T) {
	initTests(t)
	type args struct {
		bucketName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "forwardMessages", args: args{bucketName: "StateBucket"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := forwardMessages(tt.args.bucketName, true); (err != nil) != tt.wantErr {
				t.Errorf("forwardMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Logf("done")
}


func Test_parseMessage(t *testing.T) {
	bytearray, _ := json.Marshal(goodState)
	goodTestMessage := string(bytearray)

	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "good parseMessage", args: args{message: goodTestMessage}, wantErr: false},
		{name: "no data, empty json parseMessage", args: args{message: emptyTestMessage}, wantErr: false},
		{name: "not json parseMessage", args: args{message: badTestMessage}, wantErr: true},
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
		onceOnly bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		 {name: "good", args: args{bucketName: "StateBucket", onceOnly: true},  wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveStateDaemon(tt.args.bucketName, tt.args.onceOnly); (err != nil) != tt.wantErr {
				t.Errorf("saveStateDaemon() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_GetRecordList(t *testing.T) {
	type fields struct {
		UnimplementedSensorStoreAndForwardServer pb.UnimplementedSensorStoreAndForwardServer
	}
	type args struct {
		ctx context.Context
		in  *pb.GetRecordListRequest
	}

	reply := pb.GetRecordListReply{}
	reply.Result = "OK"
	bytearray, _ := json.Marshal(goodState)
	reply.Data = string(bytearray)

	ffields := fields{
		UnimplementedSensorStoreAndForwardServer: pb.UnimplementedSensorStoreAndForwardServer{},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetRecordListReply
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "good", fields: ffields, args: args{}, want: &reply, wantErr: false},
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
			got.Data = strings.ReplaceAll(got.Data, "\n","")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecordList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_GetState(t *testing.T) {
	type fields struct {
		UnimplementedSensorStoreAndForwardServer pb.UnimplementedSensorStoreAndForwardServer
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
		UnimplementedSensorStoreAndForwardServer pb.UnimplementedSensorStoreAndForwardServer
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

