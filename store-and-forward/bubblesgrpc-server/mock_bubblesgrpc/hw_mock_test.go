/*
 *
 * Copyright 2017 gRPC authors.
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

package mock_bubblesgrpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	bubblesgrpc "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/bubblesgrpc"
	hwmock "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/mock_bubblesgrpc"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestAcknowledgeStoreAndForward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSensorStoreAndForwardClient := hwmock.NewMockSensorStoreAndForwardClient(ctrl)
	req := &bubblesgrpc.SensorRequest{Name: "unit_test"}
	mockSensorStoreAndForwardClient.EXPECT().StoreAndForward(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&bubblesgrpc.SensorReply{Message: "Mocked Interface"}, nil)
	testStoreAndForward(t, mockSensorStoreAndForwardClient)
}

func testStoreAndForward(t *testing.T, client bubblesgrpc.SensorStoreAndForwardClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.StoreAndForward(ctx, &bubblesgrpc.SensorRequest{Name: "unit_test"})
	if err != nil || r.Message != "Mocked Interface" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.Message)
}
