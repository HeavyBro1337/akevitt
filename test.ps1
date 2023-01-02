$env:GOOS="linux"
go env
go build -o debug.bin .\cmd\
wsl hostname  -I
wsl ./debug.bin
pause