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

// copyright and license inspection - no issues 4/13/22

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

func SendPictureTakenEvent(PictureFilename string, PictureDatetimeMillis int64) {
	dm := messaging.NewPictureTakenMessage(PictureFilename, PictureDatetimeMillis)
	bytearray, err := json.Marshal(dm)
	message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: globals.Grpc_message_typeid_picture, Data: string(bytearray)}
	_, err = globals.Client.StoreAndForward(context.Background(), &message)
	if err != nil {
		log.Errorf("SendPictureTakenEvent ERROR %#v", err)
	} else {
		//				log.Debugf("%#v", sensor_reply)
	}
}

func uploadFile(name string) (err error) {
	//	log.Infof("uploadFile %s", name)
	path, _ := os.Getwd()
	path += "/" + name
	extraParams := map[string]string{
		"title":       "picture",
		"author":      "JR",
		"description": "uploaded picture",
	}
	url := fmt.Sprintf("http://%s:%d/api/video/%8.8d/%8.8d/upload", globals.MySite.ControllerAPIHostName,
		globals.MySite.ControllerAPIPort, globals.MySite.UserID, globals.MyDevice.DeviceID)
	log.Infof("Uploading to api at %s", url)
	request, err := newfileUploadRequest(url, extraParams, "filename", name)
	if err != nil {
		log.Errorf("uploadFile 1 fatal %#v", err)
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("uploadFile %s 2 fatal %#v", name, err)
		return err
	} else {
		var bodyContent []byte
		//		log.Infof("File upload returned %d", resp.StatusCode)
		//		log.Infof("AND %#v", resp.Header)
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		//		log.Infof("AND %#v", bodyContent)
	}
	return nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(w *multipart.Writer, fieldname string, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	//	h.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundarynvWX2EwdIfT37B1G")
	return w.CreatePart(h)
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName string, path string) (*http.Request, error) {
	//	log.Infof("newfileUploadRequest %s", uri)
	file, err := os.Open(path)
	if err != nil {
		log.Errorf("newfileUploadRequest 1 failed %#v", err)
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("newfileUploadRequest 2 failed %#v", err)
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		log.Errorf("newfileUploadRequest 3 failed %#v", err)
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	//	log.Infof("CreateFormFile %s %s", paramName, fi.Name())
	part, err := createFormFile(writer, paramName, fi.Name())
	if err != nil {
		log.Errorf("newfileUploadRequest 4 failed %#v", err)
		return nil, err
	}

	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		log.Errorf("newfileUploadRequest 5 failed %#v", err)
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}
