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
package camera

import (
	"bytes"
	"net/http"
	"testing"
)

func Test_newfileUploadRequest(t *testing.T) {

	uri := "http://nowhere.com"
	body := new(bytes.Buffer)
	req, _ := http.NewRequest("POST", uri, body)

	myMap := make(map[string]string)

	type args struct {
		uri       string
		params    map[string]string
		paramName string
		path      string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{name: "uninitialized",
			args: args{
				uri:       "",
				params:    myMap,
				paramName: "",
				path:      "/",
			},
			want:    req,
			wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newfileUploadRequest(tt.args.uri, tt.args.params, tt.args.paramName, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("newfileUploadRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//			if !reflect.DeepEqual(got, tt.want) {
			//				t.Errorf("newfileUploadRequest() got = %#v, want %#v", got, tt.want)
			//			}
		})
	}
}
