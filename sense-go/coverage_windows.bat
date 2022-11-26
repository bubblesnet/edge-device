echo %GOPATH%
set SAVE_GOOS=%GOOS%
set SAVE_GOARCH=%GOARCH%
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED="1"
set CODECOV_TOKEN='bd6757f7-5f19-40b6-81f3-68547d5b9177'

set API_HOST = 192.168.23.237
set NO_FAN_WITH_HEATER = false
set SLEEP_ON_EXIT_FOR_DEBUGGING = 60
set ACTIVEMQ_HOST = 192.168.23.237
set ACTIVEMQ_PORT  = 61611
set API_PORT = 4001
set USERID = 90000009
set DEVICEID = 70000008

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
set GOOS=%SAVE_GOOS%
set GOARCH=%SAVE_GOARCH%



