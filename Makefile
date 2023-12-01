CGO_ENABLED=0
ifeq ($(OS),Windows_NT)
	GOOS=windows
	EXT=.exe
else
	GOOS=linux
endif

ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
	GOARCH=amd64
else
ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
	GOARCH=amd64
endif
ifeq ($(PROCESSOR_ARCHITECTURE),x86)
	GOARCH=386
endif
endif

build:
	CGO_ENABLED=0 go build -o nginxlog-exporter main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nginxlog-exporter main.go

build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o nginxlog-exporter.exe main.go