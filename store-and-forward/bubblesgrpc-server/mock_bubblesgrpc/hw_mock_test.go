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
	req := &bubblesgrpc.SensorRequest{TypeId: "unit_test"}
	mockSensorStoreAndForwardClient.EXPECT().StoreAndForward(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&bubblesgrpc.SensorReply{Message: "Mocked Interface"}, nil)
	doStoreAndForward(t, mockSensorStoreAndForwardClient)
}

func doStoreAndForward(t *testing.T, client bubblesgrpc.SensorStoreAndForwardClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.StoreAndForward(ctx, &bubblesgrpc.SensorRequest{TypeId: "unit_test"})
	if err != nil || r.Message != "Mocked Interface" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.Message)
}
