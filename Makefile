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

archive: all
	tar -czvf vpp_VERSION_linux_amd64.tar.gz vpp
	tar -czvf vpp_VERSION_windows_amd64.tar.gz vpp.exe

clean:
	rm vpp vpp.exe