package main

import (
	pb "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc"
	log "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/lawg"
	"context"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

const badTestMessage = "dkdlkdkdkdk"
const emptyTestMessage = "{}"

func initTests(t *testing.T) (err error){
	log.ConfigureTestLogging("fatal,error,warn,info,debug,", ".", t)
	storeMountPoint := "/config"
	fmt.Printf("GOOS = %s GOARCH = %s", runtime.GOOS, runtime.GOARCH)
	if  runtime.GOOS == "windows" || runtime.GOOS == "darwin" || (runtime.GOARCH != "arm" && runtime.GOARCH != "arm64") {
		storeMountPoint = "."
		databaseFilename = "./testmessages.db"
	} else {
		fmt.Printf("WTF!!!")
	}
	if MyDeviceID, err =  ReadMyDeviceId(storeMountPoint, "", "deviceid"); err != nil {
		fmt.Printf("ReadMyDeviceId error %v", err )
		return err
	}

	if err := ReadFromPersistentStore(storeMountPoint, "", "config.json",&MySite,&stageSchedule); err != nil {
		fmt.Printf("Read config error %v", err )
		return err
	}

	t.Logf("MySite = %v", MySite)
	t.Logf("stageSchedule = %v", stageSchedule)
	initDb(databaseFilename)
	t.Logf("returned from initDb")
	return nil
}

func Test_forwardMessages(t *testing.T) {
	if err := initTests(t); err != nil {
		t.Errorf("error %v", err)
		return
	}
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
//			if err := saveStateDaemon(tt.args.bucketName, tt.args.onceOnly); (err != nil) != tt.wantErr {
//				t.Errorf("saveStateDaemon() error = %v, wantErr %v", err, tt.wantErr)
//			}
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

