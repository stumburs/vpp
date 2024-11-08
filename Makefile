# Temporary solution that will most likely stay permanent

linux:
	go build -ldflags="-s -w" cmd/vpp/vpp.go
	strip vpp
	upx vpp
windows:
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" cmd/vpp/vpp.go
	strip vpp.exe
	upx vpp.exe

all: linux windows

clean:
	rm vpp vpp.exe