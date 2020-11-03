export GOPATH=/Users/rodley/go/src/sense-go
env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED="1" go build -o sense_go main.go atlasezo_driver.go
cp sense_go ../../../bubbles2/sense-go
cd ../../../bubbles2
git commit -a -m "deploy by script"
balena push rpi3
balena push balenafin
cd ../go/src/sense-go