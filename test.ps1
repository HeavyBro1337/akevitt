$env:GOOS="linux"
go env
go build -o debug.bin .\cmd\
wsl echo "Launching"
wsl ./debug.bin
pause