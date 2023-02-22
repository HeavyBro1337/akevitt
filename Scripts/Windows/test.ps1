$env:GOOS="linux"
echo "Switched to linux. Building..."
go build -o .\debug.bin "..\..\"
wsl echo "Launching"
wsl ./debug.bin
pause