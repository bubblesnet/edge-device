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
