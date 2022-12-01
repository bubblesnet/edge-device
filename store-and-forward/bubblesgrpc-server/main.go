/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

// /go:generate protoc -I ../bubblesgrpc --go_out=plugins=grpc:../bubblesgrpc ../bubblesgrpc/bubblesgrpc.proto

// Package main implements a server for SensorStoreAndForward service.
package main

// copyright and license inspection - no issues 4/13/22

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
	"strconv"
	"strings"
	"time"
)

var BubblesnetVersionMajorString string
var BubblesnetVersionMinorString = ""
var BubblesnetVersionPatchString = ""
var BubblesnetBuildNumberString = ""
var BubblesnetBuildTimestamp = ""
var BubblesnetGitHash = ""

type MessageHeader struct {
	DeviceId          int64  `json:"deviceid"`
	StationId         int64  `json:"stationid"`
	SiteId            int64  `json:"siteid"`
	ContainerName     string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	EventTimestamp    int64  `json:"event_timestamp"`
	MessageType       string `json:"message_type"`
}

const (
	port              = ":50051"
	messageBucketName = "MessageBucket"
	stateBucketName   = "StateBucket"
	modeReadwrite     = 0600
)

var databaseFilename = "/data/messages.db"

var writeableDb *bolt.DB

// server is used to implement bubblesgrpc.GreeterServer.
type server struct {
	pb.UnimplementedSensorStoreAndForwardServer
}

// StoreAndForward implements bubblesgrpc.StoreAndForward
func (s *server) StoreAndForward(_ context.Context, in *pb.SensorRequest) (*pb.SensorReply, error) {
	log.Debugf("Received: sequence %v - %s", in.GetSequence(), in.GetData())
	// This asynchronicity is a joint where data can leak out and get lost if either of these methods fail
	// They should at least be independently asynchronous
	go func() {
		_ = addRecord(messageBucketName, in.GetData(), in.GetSequence())
		parseMessageForCurrentState(in.GetData())
	}()
	return &pb.SensorReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Message: ""}, nil
}

// parseMessageForCurrentState process a message for any persistent state we might want to keep.
func parseMessageForCurrentState(message string) {
	genericMessage := GenericSensorMessage{}
	err := json.Unmarshal([]byte(message), &genericMessage)
	if err != nil {
		return
	}
	switch genericMessage.MessageType {
	case message_type_measurement:
		switch genericMessage.MeasurementName {
		case Measurement_name_plant_height:
			if genericMessage.FloatValue != ExternalCurrentState.PlantHeightIn {
				log.Infof("plant_height changed from %f to %f", ExternalCurrentState.PlantHeightIn, genericMessage.FloatValue)
			}
			ExternalCurrentState.PlantHeightIn = genericMessage.FloatValue
			break

		case Measurement_name_temp_water:
			if genericMessage.FloatValue != ExternalCurrentState.WaterTemp {
				log.Infof("temp_water changed from %f to %f", ExternalCurrentState.WaterTemp, genericMessage.FloatValue)
			}
			ExternalCurrentState.WaterTemp = genericMessage.FloatValue
			break

		case measurement_name_temp_air_middle:
			if genericMessage.FloatValue != ExternalCurrentState.TempAirMiddle {
				log.Infof("temp_air_middle changed from %f to %f", ExternalCurrentState.TempAirMiddle, genericMessage.FloatValue)
			}
			ExternalCurrentState.TempAirMiddle = genericMessage.FloatValue
			break
		case measurement_name_humidity_internal:
			if genericMessage.FloatValue != ExternalCurrentState.HumidityInternal {
				log.Infof("humidity_internal changed from %f to %f", ExternalCurrentState.HumidityInternal, genericMessage.FloatValue)
			}
			ExternalCurrentState.HumidityInternal = genericMessage.FloatValue
			break
		case Measurement_name_light_internal:
			if genericMessage.FloatValue != ExternalCurrentState.LightInternal {
				log.Infof("light_internal changed from %f to %f", ExternalCurrentState.LightInternal, genericMessage.FloatValue)
			}
			ExternalCurrentState.LightInternal = genericMessage.FloatValue
			break
		case Measurement_name_pressure_internal:
			if genericMessage.FloatValue != ExternalCurrentState.PressureInternal {
				log.Infof("pressure_internal changed from %f to %f", ExternalCurrentState.PressureInternal, genericMessage.FloatValue)
			}
			ExternalCurrentState.PressureInternal = genericMessage.FloatValue
			break
		case "":
			log.Warnf("Empty state message sent from %s to store-and-forward %#v", genericMessage.ContainerName, genericMessage)
			break
		default:
			log.Infof("Unused from %s GenericMessage.SensorName/MeasurementName = %s/%s value %f", genericMessage.ContainerName, genericMessage.SensorName, genericMessage.MeasurementName, genericMessage.FloatValue)
			break
		}
		break
	default:
		log.Warnf("Non-measurement message sent from %s to store-and-forward %#v", genericMessage.ContainerName, genericMessage)
		break
	}
}

// GetState StoreAndForward implements bubblesgrpc.GetState
func (s *server) GetState(_ context.Context, in *pb.GetStateRequest) (*pb.GetStateReply, error) {
	//	log.Debugf("GetState Received: sequence %v - %s", in.GetSequence(), in.GetData())
	//	if in.GetSequence() %5 == 0 {
	//		return &pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "ERROR" }, nil
	//	} else {

	ret := pb.GetStateReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(),
		Result: "OK", TempAirMiddle: float32(ExternalCurrentState.TempAirMiddle),
		HumidityInternal: float32(ExternalCurrentState.HumidityInternal),
		TempWater:        float32(ExternalCurrentState.WaterTemp), LightInternal: float32(ExternalCurrentState.LightInternal)}
	//	log.Infof("GetState returning %#v", ret)
	return &ret, nil
	//	}
}

//
func (s *server) GetRecordList(_ context.Context, in *pb.GetRecordListRequest) (*pb.GetRecordListReply, error) {
	log.Debug(fmt.Sprintf("GetRecordList Received: sequence %v - %s", in.GetSequence(), in.GetData()))
	jsn, _ := getStateAsJson(stateBucketName, 2020, 2, 15)
	return &pb.GetRecordListReply{Sequence: in.GetSequence(), TypeId: in.GetTypeId(), Result: "OK", Data: jsn}, nil
}

func forwardMessages(bucketName string, oneOnly bool) (err error) {
	log.Debugf("forwardMessages %s", bucketName)
	for {
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
				log.Errorf("delete frm bucket failed %v", err)
			}
		}
		// delete the forwarded messages
		forwarded = forwarded[:0]
		if oneOnly {
			break
		}
		time.Sleep(3 * time.Second)
	}
	return err
}

func postIt(message []byte) (err error) {
	var messageHeader MessageHeader
	apiName := "measurement"
	if err1 := json.Unmarshal(message, &messageHeader); err1 != nil {
		log.Errorf("error unmarshalling message pre post #+v", err1)
	} else {
		log.Infof("message type %s", messageHeader.MessageType)
		if messageHeader.MessageType != "measurement" {
			log.Infof("non-measurement message %s", string(message))
		}
	}
	url := fmt.Sprintf("http://%s:%d/api/%s/%8.8d/%8.8d", MySite.ControllerAPIHostName, MySite.ControllerAPIPort, apiName, MySite.UserID, MyDeviceID)
	//	log.Infof("Sending to %s", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(message))
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

func handleVersioningFromLoader() (err error) {
	BubblesnetBuildTimestamp = strings.ReplaceAll(BubblesnetBuildTimestamp, "'", "")
	BubblesnetBuildTimestamp = strings.ReplaceAll(BubblesnetBuildTimestamp, "_", " ")
	return nil
}

func SleepBeforeExit() {
	snaptime := os.Getenv(ENV_SLEEP_ON_EXIT_FOR_DEBUGGING)
	naptime, err := strconv.ParseInt(snaptime, 10, 32)
	if err != nil {
		log.Errorf("SLEEP_ON_EXIT_FOR_DEBUGGING %s conversion error %#v", snaptime, err)
		naptime = 60
	}
	fmt.Printf("Exiting because of bad configuration - sleeping for %d seconds to allow intervention\n", naptime)
	time.Sleep(time.Duration(naptime) * time.Second)
}

func main() {
	log.ConfigureLogging("fatal,error,warn,info,debug,", ".")

	if err := handleVersioningFromLoader(); err != nil {
		log.Errorf("handleVersioningFromLoader %+v", err)
		SleepBeforeExit()
		os.Exit(222)
	}
	fmt.Printf("Bubblesnet %s.%s.%s build %s timestamp %s githash %s\n", BubblesnetVersionMajorString,
		BubblesnetVersionMinorString, BubblesnetVersionPatchString, BubblesnetBuildNumberString,
		BubblesnetBuildTimestamp, BubblesnetGitHash)

	storeMountPoint := "/config"

	fmt.Printf("GOOS = %s, GOARCH=%s\n", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || (runtime.GOARCH != "arm" && runtime.GOARCH != "arm64") {
		storeMountPoint = "./testdata"
		databaseFilename = "./messages.db"
	}
	var err error
	MyDeviceID, err = ReadMyDeviceId()
	fmt.Printf("Read deviceid %d\n", MyDeviceID)
	if err != nil {
		fmt.Printf("error read device %v\n", err)
		return
	}

	WaitForConfigFile(storeMountPoint, "", "config.json")
	err = ReadCompleteSiteFromPersistentStore(storeMountPoint, "", "config.json", &MySite, &stageSchedule)
	var nilerr error
	MySite.ControllerAPIHostName, _ = os.Getenv(ENV_API_HOST), nilerr
	MySite.ControllerActiveMQHostName, _ = os.Getenv(ENV_ACTIVEMQ_HOST), nilerr
	MySite.ControllerAPIPort, _ = strconv.Atoi(os.Getenv(ENV_API_PORT))
	MySite.ControllerActiveMQPort, _ = strconv.Atoi(os.Getenv(ENV_ACTIVEMQ_PORT))
	MySite.UserID, _ = strconv.ParseInt(os.Getenv(ENV_USERID), 10, 64)
	d := EdgeDevice{DeviceID: MyDeviceID}
	MyDevice = &d

	fmt.Printf("MySite = %v", MySite)
	fmt.Printf("stageSchedule = %v", stageSchedule)
	initDb(databaseFilename)

	//	clearDatabase(stateBucketName)
	//	deletePriorTo(stateBucketName, 1581483579497)

	go func() {
		_ = forwardMessages(messageBucketName, false)
	}()

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

func WaitForConfigFile(storeMountPoint string, relativePath string, fileName string) {
	fmt.Printf("WaitForConfigFile %s %s %s\n", storeMountPoint, relativePath, fileName)
	for index := 0; index <= 60; index++ {
		if exists, err := ConfigFileExists(storeMountPoint, "", "config.json"); (exists == true) && (err == nil) {
			fmt.Printf("apparently config.json exists\n")
			return
		}
		if index == 60 {
			fmt.Printf("waited too long for file %s to be downloaded. Probably no connection.  Exiting\n", fileName)
			SleepBeforeExit()
			os.Exit(1)
		}
		fmt.Printf("Sleeping 60 seconds waiting for someone to bring us a /config/config.json\n")
		time.Sleep(60 * time.Second)
	}
}

func ConfigFileExists(storeMountPoint string, relativePath string, fileName string) (exists bool, err error) {
	fmt.Printf("ConfigFileExists\n")
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	if relativePath == "" {
		fullpath = storeMountPoint + "/" + fileName
	}
	if _, err := os.Stat(fullpath); err != nil {
		return false, err

	}
	return true, nil
}
