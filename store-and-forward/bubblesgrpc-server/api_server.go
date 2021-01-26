package main

import (
	pb "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc"
	log "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/lawg"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
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
/*
func convertJsonStateArrayStringToCsv( states []state ) (string, error) {
	ret := "SampleTimestamp,SampleTimestampS,DistanceIn,Ph,TempF,Humidity,Pressure,Light,GrowLightVeg,Heater,HeaterPad,Humidifier\n"
	for i := 0; i < len(states); i++ {
		ret = ret + fmt.Sprintf("%d,%s,%f,%f,%f,%f,%f,%f,%t,%t,%t,%t\n",
			states[i].SampleTimestamp,
			states[i].SampleTimestampS,
			states[i].DistanceIn,
			states[i].Ph,
			states[i].TempF,
			states[i].Humidity,
			states[i].Pressure,
			states[i].Light,
			states[i].GrowLightVeg,
			states[i].Heater,
			states[i].HeaterPad,
			states[i].Humidifier)
	}
	return ret, nil
}
*/

func getContentDisposition( format string ) string {
	t := time.Now()
	filename := fmt.Sprintf("%d-%02d-%02dT%02d_%02d_%02d-00_00.%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), format)
	content_dispostion := fmt.Sprintf("attachment; filename=%s", filename)
	return content_dispostion
}
/*
func StartApiServer() {
	log.Infof("StartApiServer"))
	csvHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debugf("Received API request %v", req ))
		ret, err := requestStateList()
		if err == nil {
			w.Header().Set("Content-Type", "text/csv; charset=utf-8")
			w.Header().Set("Content-Disposition", getContentDisposition("csv"))
			var states []state
			_ = json.Unmarshal([]byte(ret), &states)
			s, _ := convertJsonStateArrayStringToCsv(states)
			io.WriteString(w, s)
		} else {
			io.WriteString(w,"ERROR\n")
		}
	}
	http.HandleFunc("/api/state/csv", csvHandler)

	jsonHandler := func(w http.ResponseWriter, req *http.Request) {
		ret, err := requestStateList()
		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
//			w.Header().Set("Content-Disposition", getContentDisposition("json"))
			io.WriteString(w, ret)
		} else {
			io.WriteString(w,"ERROR\n")
		}
	}
	http.HandleFunc("/api/state/json", jsonHandler)
	log.Fatalf("%v", http.ListenAndServe(":8080", nil))
}
*/

type getRecordsRequest struct {
	BucketName string `json:"bucket_name"`
	Year int `json:"year"`
	Month int `json:"month"`
	Day int `json:"day"`
}

func requestStateList() (string, error) {
	log.Infof("requestStateList"))
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSensorStoreAndForwardClient(conn)

	var getRecsReq = getRecordsRequest {
		BucketName: stateBucketName,
		Year:       2020,
		Month:      2,
		Day:        15,
	}
	bytearray, err := json.Marshal(getRecsReq)
	if err != nil {
		log.Errorf("requestStateList error %v", err)
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pbr := pb.GetRecordListRequest{Sequence: GetSequenceNumber(), TypeId: "database", Data: string(bytearray)}
	log.Debugf("c.GetRecordList"))
	r, err1 := c.GetRecordList(ctx, &pbr)
	if err1 != nil {
		log.Errorf("GetRecordListRequest failed: %v", err1)
		return "", err1
	} else {
		log.Debugf("GetRecordListRequest Received ack for sequence: %d message: TOO LONG YOU DOPE!", r.GetSequence()))
	}
	return r.GetData(), nil

}
