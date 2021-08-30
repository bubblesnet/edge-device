/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

///go:generate protoc -I ../bubblesgrpc --go_out=plugins=grpc:../bubblesgrpc ../bubblesgrpc/bubblesgrpc.proto

// Package main implements a server for SensorStoreAndForward service.
package main

import (
	pb "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc"
	log "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/lawg"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)


var BubblesnetVersionMajorString string
var BubblesnetVersionMinorString=""
var BubblesnetVersionPatchString=""
var BubblesnetBuildNumberString=""
var BubblesnetBuildTimestamp=""
var BubblesnetGitHash=""

const (
	port              = ":50051"
	messageBucketName = "MessageBucket"
	stateBucketName   = "StateBucket"
	modeReadwrite     = 0600
)
var databaseFilename  = "/data/messages.db"

var writeableDb *bolt.DB

// server is used to implement bubblesgrpc.GreeterServer.
type server struct {
	pb.UnimplementedSensorStoreAndForwardServer
}

// StoreAndForward implements bubblesgrpc.StoreAndForward
func (s *server) StoreAndForward(_ context.Context, in *pb.SensorRequest) (*pb.SensorReply, error) {
//	log.Debugf("Received: sequence %v - %s", in.GetSequence(), in.GetData())
	go func() {
		_ = addRecord(messageBucketName, in.GetData(), in.GetSequence())
		parseMessageForCurrentState(in.GetData())
	}()
	return &pb.SensorReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Message: ""}, nil
}

func parseMessageForCurrentState( message string ) {
	genericMessage := GenericSensorMessage{}
	log.Infof("parsing %s", message)
	err := json.Unmarshal([]byte(message), &genericMessage)
	if err != nil {
		return
	}
	log.Infof("GenericMessage.SensorName/MeasurementName = %s/%s", genericMessage.SensorName, genericMessage.MeasurementName)
	switch genericMessage.MeasurementName {
	case "temp_air_middle":
		ExternalCurrentState.TempF = genericMessage.FloatValue
		break
	case "humidity_internal":
		ExternalCurrentState.Humidity = genericMessage.FloatValue
		break
	}
}

// StoreAndForward implements bubblesgrpc.GetState
func (s *server) GetState(_ context.Context, in *pb.GetStateRequest) (*pb.GetStateReply, error) {
//	log.Debugf("GetState Received: sequence %v - %s", in.GetSequence(), in.GetData())
//	if in.GetSequence() %5 == 0 {
//		return &pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "ERROR" }, nil
//	} else {
		ret := pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", TempF: float32(ExternalCurrentState.TempF), Humidity: float32(ExternalCurrentState.Humidity)}
//		log.Infof("GetState returning %v", ret)
		return &ret, nil
		//	}
}

//
func (s *server) GetRecordList(_ context.Context, in *pb.GetRecordListRequest) (*pb.GetRecordListReply, error) {
	log.Debug( fmt.Sprintf("GetRecordList Received: sequence %v - %s", in.GetSequence(), in.GetData()))
	jsn,_ := getStateAsJson(stateBucketName,2020,2,15)
	return &pb.GetRecordListReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Data: jsn}, nil
}

func forwardMessages(bucketName string, oneOnly bool) (err error) {
	log.Debugf("forwardMessages %s", bucketName)
	for ;;  {
		var forwarded []string

		_ = writeableDb.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
//			log.Debugf("getting b")
			b := tx.Bucket([]byte(bucketName))
//			log.Debugf("b = %v", b)
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				// for k, v := c.First(); k != nil; k, v = c.Next() {
//				log.Debugf("forwarding key=%s, value=%s from %s\n", k, string(v), bucketName)
				_ = postIt(v)
				forwarded = append(forwarded, string(k))
				time.Sleep(250 * time.Millisecond)
			}
			return nil
		})

		for i := 0; i < len(forwarded); i++ {
			if err := deleteFromBucket(bucketName, []byte(forwarded[i])); err != nil {
				log.Errorf( "delete from bucket failed %v", err)
			}
		}
		// delete the forwarded messages
		forwarded = forwarded[:0]
		if oneOnly {
			break
		}
		time.Sleep(3*time.Second)
	}
	return err
}

func postIt(message []byte) (err error){
	url := fmt.Sprintf("http://%s:%d/api/measurement/%8.8d/%8.8d", MySite.ControllerHostName, MySite.ControllerAPIPort, MySite.UserID, MyDeviceID)
	log.Infof("Sending to %s", url)
	resp, err := http.Post(url,
		"application/json", bytes.NewBuffer(message))
	if err != nil {
		log.Errorf("post error %v", err)
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("readall error %v", err)
		return err
	}
//	log.Debugf("response %s", string(body))
	return nil
}

var stageSchedule StageSchedule
func handleVersioningFromLoader() (err error ) {
	BubblesnetBuildTimestamp = strings.ReplaceAll(BubblesnetBuildTimestamp, "'", "")
	BubblesnetBuildTimestamp = strings.ReplaceAll(BubblesnetBuildTimestamp, "_", " ")
	return nil
}

func main() {
	log.ConfigureLogging("fatal,error,warn,info,debug,", ".")

	if err := handleVersioningFromLoader(); err != nil {
		log.Errorf("handleVersioningFromLoader %+v", err )
		os.Exit(222)
	}
	fmt.Printf("Bubblesnet %s.%s.%s build %s timestamp %s githash %s\n", BubblesnetVersionMajorString,
		BubblesnetVersionMinorString, BubblesnetVersionPatchString, BubblesnetBuildNumberString,
		BubblesnetBuildTimestamp, BubblesnetGitHash)

	storeMountPoint := "/config"

	fmt.Printf("GOOS = %s, GOARCH=%s\n", runtime.GOOS, runtime.GOARCH)
	if  runtime.GOOS == "windows" || runtime.GOOS == "darwin" || (runtime.GOARCH != "arm" && runtime.GOARCH != "arm64") {
		storeMountPoint = "./testdata"
		databaseFilename = "./messages.db"
	}
	var err error
	MyDeviceID, err = ReadMyDeviceId("/config","", "deviceid")
	if err != nil {
		fmt.Printf("error read device %v\n", err)
		return
	}
	fmt.Printf("Read deviceid %d\n", MyDeviceID)

	WaitForConfigFile(storeMountPoint, "", "config.json")
	err = ReadFromPersistentStore(storeMountPoint, "", "config.json", &MySite, &stageSchedule)

	fmt.Printf("MySite = %v", MySite)
	fmt.Printf("stageSchedule = %v", stageSchedule)
	initDb(databaseFilename)

//	clearDatabase(stateBucketName)
//	deletePriorTo(stateBucketName, 1581483579497)

	go func() {
		_ = forwardMessages(messageBucketName, false)
	}()
/*	go func() {
		_ = forwardMessages(stateBucketName, false)
	}()
	go func() {
		_ = saveStateDaemon(stateBucketName, false)
	}()

 */

//	go StartApiServer()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSensorStoreAndForwardServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
