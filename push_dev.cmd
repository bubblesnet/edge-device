cd sense-go
powershell .\build_sense_go.bat
cd ..\store-and-forward\bubblesgrpc-server
powershell .\build_store_and_forward_go.bat
cd ..\..

balena push bubblesnet3_dev

