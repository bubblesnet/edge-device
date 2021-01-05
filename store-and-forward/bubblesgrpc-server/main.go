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
	"context"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"net"
	"runtime"
	"strings"
	"time"
)

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

type state struct {
	SampleTimestamp int64 `json:"sample_timestamp"`
	SampleTimestampS string `json:"sample_timestamp_s"`
	TempF float32 `json:"tempF,omitempty"`
	Humidity float32 `json:"humidity,omitempty"`
	Light float32 `json:"light,omitempty"`
	DistanceIn float32 `json:"distance_in,omitempty"`
	Pressure float32 `json:"pressure,omitempty"`
	Ph float32 `json:"ph,omitempty"`
	Humidifier bool `json:"humidifier"`
	Heater bool `json:"heater"`
	HeaterPad bool `json:"heater_pad"`
	GrowLightVeg bool `json:"grow_light_veg"`
}

var currentstate = state {}

// StoreAndForward implements bubblesgrpc.StoreAndForward
func (s *server) StoreAndForward(_ context.Context, in *pb.SensorRequest) (*pb.SensorReply, error) {
//	log.Printf("Received: sequence %v - %s", in.GetSequence(), in.GetData())
	go func() {
		_ = addRecord(messageBucketName, in.GetData())
	}()
	_ = parseMessage(in.GetData())
	return &pb.SensorReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Message: ""}, nil
}

// StoreAndForward implements bubblesgrpc.GetState
func (s *server) GetState(_ context.Context, in *pb.GetStateRequest) (*pb.GetStateReply, error) {
//	log.Printf("GetState Received: sequence %v - %s", in.GetSequence(), in.GetData())
	if in.GetSequence() %5 == 0 {
		return &pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "ERROR" }, nil
	} else {
		return &pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", TempF: currentstate.TempF, Humidity: currentstate.Humidity}, nil
	}
}

// StoreAndForward implements bubblesgrpc.GetState
func (s *server) GetRecordList(_ context.Context, in *pb.GetRecordListRequest) (*pb.GetRecordListReply, error) {
	log.Debug( fmt.Sprintf("GetRecordList Received: sequence %v - %s", in.GetSequence(), in.GetData()))
	_,_ = getStateAsJson(stateBucketName,2020,2,15)
	return &pb.GetRecordListReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Data: csvx}, nil
}

func parseMessage(message string) (err error) {
	var partialstate = state{ SampleTimestamp: 0, TempF: -77.7, Light: -77.7, Humidity: -77.7, DistanceIn: -77.7, Ph: -77.7, Pressure: -77.7 }
	err = json.Unmarshal([]byte(message), &partialstate)
	if err != nil {
		log.Errorf("unmarshal message error %v", err)
		return err
	}
	log.Debug(fmt.Sprintf("message %s", message))
	if strings.Contains(message, "tempF") && partialstate.TempF != -77.7 {
			currentstate.TempF = partialstate.TempF
	}
	if strings.Contains(message, "humidity") && partialstate.Humidity != -77.7 {
			currentstate.Humidity = partialstate.Humidity
	}
	if strings.Contains(message, "light")  && partialstate.Light != -77.7{
		currentstate.Light = partialstate.Light
	}
	if strings.Contains(message, "distance_in")  && partialstate.DistanceIn != -77.7 {
		currentstate.DistanceIn = partialstate.DistanceIn
	}
	if strings.Contains(message, "ph")  && partialstate.Ph != -77.7 {
		currentstate.Ph = partialstate.Ph
	}
	if strings.Contains(message, "pressure")  && partialstate.Pressure != -77.7 {
		currentstate.Pressure = partialstate.Pressure
	}
	if strings.Contains(message, "sample_timestamp")  && partialstate.SampleTimestamp != 0 {
		currentstate.SampleTimestamp = partialstate.SampleTimestamp
		nsec := (partialstate.SampleTimestamp%1000) * 1000000
		t := time.Unix(partialstate.SampleTimestamp/1000,nsec)
		currentstate.SampleTimestampS = t.String()
		log.Debug(fmt.Sprintf("Setting sampletimestamps to %s", currentstate.SampleTimestampS))
	}
	if strings.Contains(message, "humidifier") {
		currentstate.Humidifier = partialstate.Humidifier
	}
	if strings.Contains(message, "heater") {
		currentstate.Heater = partialstate.Heater
	}
	if strings.Contains(message, "heater_pad") {
		currentstate.HeaterPad = partialstate.HeaterPad
	}
	if strings.Contains(message, "sample_timestamp") {
		currentstate.GrowLightVeg = partialstate.GrowLightVeg
	}

	log.Debug(fmt.Sprintf("currentstate %v", currentstate))

	return nil
}

func saveStateDaemon( bucketName string, onceOnly bool ) error {
	for ;; {
		log.Debug(fmt.Sprintf("Saving state to writeableDb"))
		if currentstate.SampleTimestampS == "" {
			currentstate.SampleTimestamp = time.Now().UnixNano()
			currentstate.SampleTimestampS = time.Now().Format(time.ANSIC)
		}
		stringState, err := json.Marshal(currentstate)
		if err != nil {
			log.Error(err)
		} else {
			log.Debug(fmt.Sprintf("Saving state %s to writeableDb", string(stringState)))
			_ = addRecord(bucketName, string(stringState))
		}
		if onceOnly {
			break
		}
		time.Sleep(time.Minute)
	}
	return nil
}

func forwardMessages(bucketName string, oneOnly bool) (err error) {
	log.Debugf("forwardMessages %s", bucketName)
	for ;;  {
		var forwarded []string

		_ = writeableDb.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			log.Debugf("getting b")
			b := tx.Bucket([]byte(bucketName))
			log.Debugf("b = %v", b)
			c := b.Cursor()

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				// for k, v := c.First(); k != nil; k, v = c.Next() {
				//				log.Debug(fmt.Sprintf("forwarding key=%s, value=%s\n", k, v)
				forwarded = append(forwarded, string(k))
				time.Sleep(250 * time.Millisecond)
			}
			return nil
		})

		for i := 0; i < len(forwarded); i++ {
			err := deleteFromBucket(bucketName, []byte(forwarded[i]))
			if err != nil {
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

var config Config
var stageSchedule StageSchedule

func main() {
	log.ConfigureLogging("fatal,error,warn,info,debug,", ".")
	storeMountPoint := "/config"
	if  runtime.GOOS == "windows"{
		storeMountPoint = "."
		databaseFilename = "./messages.db"
	}
	_ = ReadFromPersistentStore(storeMountPoint, "", "config.json",&config,&stageSchedule)

	fmt.Printf("config = %v", config)
	fmt.Printf("stageSchedule = %v", stageSchedule)
	initDb(databaseFilename)

//	clearDatabase(stateBucketName)
//	deletePriorTo(stateBucketName, 1581483579497)

	go func() {
		_ = forwardMessages(messageBucketName, false)
	}()
	go func() {
		_ = saveStateDaemon(stateBucketName, false)
	}()

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
