package globals

import "testing"

func TestReadMyDeviceId(t *testing.T) {
	type args struct {
		storeMountPoint string
		relativePath    string
		fileName        string
	}
	tests := []struct {
		name    string
		args    args
		wantId  int64
		wantErr bool
	}{
		{
			name: "happy",
			args: args{
				storeMountPoint: "../testdata",
				relativePath:    "",
				fileName:        "deviceid",
			},
			wantId:  70000008,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := ReadMyDeviceId(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMyDeviceId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("ReadMyDeviceId() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestReadServerHostname(t *testing.T) {
	type args struct {
		storeMountPoint string
		relativePath    string
		fileName        string
	}
	tests := []struct {
		name         string
		args         args
		wantHostname string
		wantErr      bool
	}{
		{
			name: "happy",
			args: args{
				storeMountPoint: "../testdata",
				relativePath:    "",
				fileName:        "hostname",
			},
			wantHostname: "192.168.21.237",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostname, err := ReadMyServerHostname(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMyDeviceId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHostname != tt.wantHostname {
				t.Errorf("ReadMyDeviceId() gotId = %v, want %v", gotHostname, tt.wantHostname)
			}
		})
	}
}

func initGlobalsLocally(t *testing.T) {
	MyDeviceID = 70000008
	if err := ReadCompleteSiteFromPersistentStore("../testdata", "", "config.json", &MySite, &CurrentStageSchedule); err != nil {
		t.Errorf("getConfigFromServer() error = %#v", err)
	}
}

func TestValidateConfigurable(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigurable(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigurable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	initGlobalsLocally(t)
	for _, tt := range tests {
		tt.wantErr = false
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigurable(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigurable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfigured(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigured("test"); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigured() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	initGlobalsLocally(t)
	for _, tt := range tests {
		tt.wantErr = false
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigured("test"); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigured() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfigFromServer(t *testing.T) {
	type args struct {
		storeMountPoint string
		relativePath    string
		fileName        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "sad", wantErr: true, args: args{
			storeMountPoint: "./testdata",
			relativePath:    "",
			fileName:        "retreivedConfig.json",
		}},
		{name: "less sad", wantErr: true, args: args{
			storeMountPoint: "./testdata",
			relativePath:    "",
			fileName:        "retreivedConfig.json",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfigFromServer(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("GetConfigFromServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		initGlobalsLocally(t)
	}
}
