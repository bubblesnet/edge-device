package grpc

import (
//	pb "bubblesgrpc.pb"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/log"
	"google.golang.org/grpc"
	"bubblesnet/edge-device/sense-go/globals"

	//	"sense-go"
	"time"
)

var SequenceNumber int32 = 0

const address = "store-and-forward:50051"

var conn *grpc.ClientConn

func GetSequenceNumber() int32 {
	if SequenceNumber >= 10000 {
		SequenceNumber = 1
	} else {
		SequenceNumber = SequenceNumber + 1
	}
	return SequenceNumber
}

func SendStoreAndForwardMessageWithRetries(sequence int32, data string, numretries int) error {
	for i := 0; i < numretries; i++ {
		err := SendStoreAndForwardMessage(sequence, data)
		if err == nil {
			if i > 0 {
				log.Warn(fmt.Sprintf("Succeeded on retry %d", i))
			}
			return err
		} else {
			//			log.Debug(fmt.Sprintf("SendStoreAndForwardMessageWithRetries failed on retry %d\n", i ))
		}
	}
	log.Error(fmt.Sprintf("error SendStoreAndForwardMessageWithRetries failed all retries - %s", data))
	return errors.New("SendStoreAndForwardMessageWithRetries failed all retries")
}

func SendStoreAndForwardMessage(sequence int32, data string) error {
	//	log.Debug(fmt.Sprintf("SendStoreAndForwardMessage sequence %d %s\n", sequence, data))
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error(fmt.Sprintf("did not connect: %v", err))
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewSensorStoreAndForwardClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err1 := c.StoreAndForward(ctx, &SensorRequest{Sequence: sequence, TypeId: "sensor", Data: data})
	if err1 != nil {
		log.Error("SendStoreAndForwardMessage failed: %v", err1)
		return err1
	} else {
//		log.Debug(fmt.Sprintf("SendStoreAndForwardMessage Received ack for sequence: %d message: %s\n", r.GetSequence(), r.GetMessage()))
	}
	return nil
}

func SendGetStateMessage(sequence int32, data string) {
//	log.Debug(fmt.Sprintf("SendGetStateMessage sequence %d %s\n", sequence, data))
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error(fmt.Sprintf("did not connect: %v", err))
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewSensorStoreAndForwardClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetState(ctx, &GetStateRequest{Sequence: sequence, TypeId: "localState", Data: data})
	if err != nil {
		log.Error(fmt.Sprintf("Received NACK for sequence: %d tempF: %f", r.GetSequence()))
		log.Fatalf("GetState call failed: %v", err)
	} else {
//		log.Debug(fmt.Sprintf("SendGetStateMessage Received ack for sequence: %d tempF: %f\n", r.GetSequence(), r.GetTempF() ))
		if r.GetTempF() != 0.0 {
			globals.ExternalCurrentState.TempF = r.GetTempF()
		}
		if r.GetHumidity() != 0.0 {
			globals.ExternalCurrentState.Humidity = r.GetHumidity()
		}
	}
}
