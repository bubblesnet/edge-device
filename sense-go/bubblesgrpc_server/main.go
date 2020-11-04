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

//go:generate protoc -I ../bubblesgrpc --go_out=plugins=grpc:../bubblesgrpc ../bubblesgrpc/bubblesgrpc.proto

// Package main implements a server for SensorStoreAndForward service.
package main

import (
	pb "../../bubblesgrpc/bubblesgrpc"
	bu "bitbucket.org/jrodley/balena-utils-go"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"google.golang.org/grpc"
	"net"
	"strings"
	"time"
	//	pb "bubblesgrpc"
	bolt "go.etcd.io/bbolt"
)

const (
	port              = ":50051"
	messageBucketName = "MessageBucket"
	stateBucketName   = "StateBucket"
	databaseFilename = "/data/messages.db"
	mode_readwrite = 0600
)

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
	Ph float32 `json:"pH,omitempty"`
	Humidifier bool `json:"humidifier"`
	Heater bool `json:"heater"`
	HeaterPad bool `json:"heater_pad"`
	GrowLightVeg bool `json:"grow_light_veg"`
}

var currentstate = state {}

// StoreAndForward implements bubblesgrpc.StoreAndForward
func (s *server) StoreAndForward(ctx context.Context, in *pb.SensorRequest) (*pb.SensorReply, error) {
//	log.Printf("Received: sequence %v - %s", in.GetSequence(), in.GetData())
	go addRecord(messageBucketName, in.GetData())
	parseMessage(in.GetData())
	return &pb.SensorReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Message: ""}, nil
}

// StoreAndForward implements bubblesgrpc.GetState
func (s *server) GetState(ctx context.Context, in *pb.GetStateRequest) (*pb.GetStateReply, error) {
//	log.Printf("GetState Received: sequence %v - %s", in.GetSequence(), in.GetData())
	if in.GetSequence() %5 == 0 {
		return &pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "ERROR" }, nil
	} else {
		return &pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", TempF: currentstate.TempF, Humidity: currentstate.Humidity}, nil
	}
}

// StoreAndForward implements bubblesgrpc.GetState
func (s *server) GetRecordList(ctx context.Context, in *pb.GetRecordListRequest) (*pb.GetRecordListReply, error) {
	log.Debug( fmt.Sprintf("GetRecordList Received: sequence %v - %s", in.GetSequence(), in.GetData()))
	getStateAsJson(stateBucketName,2020,2,15)
	return &pb.GetRecordListReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Data: csvx}, nil
}

func parseMessage(message string) error {
	var partialstate = state{ SampleTimestamp: 0, TempF: -77.7, Light: -77.7, Humidity: -77.7, DistanceIn: -77.7, Ph: -77.7, Pressure: -77.7 }
	_ = json.Unmarshal([]byte(message), &partialstate)
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

func saveState( bucketName string ) error {
	for ; 1 != 2; {
		log.Debug(fmt.Sprintf("Saving state to writeableDb"))
		if currentstate.SampleTimestampS == nil {
			currentstate.SampleTimestamp = time.Now().UnixNano()
			currentstate.SampleTimestampS = t.String()
		}
		stringState, err := json.Marshal(currentstate)
		if err != nil {
			log.Error(err)
		} else {
			log.Debug(fmt.Sprintf("Saving state %s to writeableDb", string(stringState)))
			addRecord(stateBucketName, string(stringState))
		}
		time.Sleep(time.Minute)
	}
	return nil
}

func forwardMessages(bucketName string) error {
	for ; 1 != 2;  {
		var forwarded []string

		writeableDb.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(bucketName))

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
				log.Error( fmt.Sprintf("delete from bucket failed %v", err))
			}
		}
		// delete the forwarded messages
		forwarded = forwarded[:0]

		time.Sleep(3*time.Second)
	}
	return nil
}

var config bu.Config
var stageSchedule bu.StageSchedule

func main() {

	bu.ReadFromPersistentStore("/config", "", "config.json",&config,&stageSchedule)

	log.Info(fmt.Sprintf("config = %v", config))
	log.Info(fmt.Sprintf("stageSchedule = %v", stageSchedule))
	initDb()

//	clearDatabase(stateBucketName)
//	deletePriorTo(stateBucketName, 1581483579497)

	go forwardMessages(messageBucketName)
	go saveState(stateBucketName)

	go StartApiServer()

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
