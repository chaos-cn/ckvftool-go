#! /bin/sh

# 【darwin/arm64】
echo "start build darwin/arm64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build  -o ../release/ckvftool-mac-arm64 ../main.go

# 【darwin/amd64】
#echo "start build darwin/amd64 ..."
#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -o ../release/ckvftool-mac-amd64 ../main.go

# 【linux/amd64】
#echo "start build linux/amd64 ..."
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ../release/ckvftool-linux-amd64 ../main.go

# 【windows/amd64】
#echo "start build windows/amd64 ..."
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -o ../release/ckvftool-amd64.exe ../main.go

echo "Congratulations,all build success!!!"
