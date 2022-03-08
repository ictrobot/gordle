NAME=gordle
ENV=CGO_ENABLED=0
FLAGS=-trimpath -ldflags="-s -w"

build:
	env $(ENV) go build $(FLAGS)

install:
	env $(ENV) go install $(FLAGS)

buildAll: build linux windows darwin

linux:
	env GOOS=linux   GOARCH=386   $(ENV) go build $(FLAGS) -o bin/$(NAME)-linux-386
	env GOOS=linux   GOARCH=amd64 $(ENV) go build $(FLAGS) -o bin/$(NAME)-linux-amd64
	env GOOS=linux   GOARCH=arm   $(ENV) go build $(FLAGS) -o bin/$(NAME)-linux-arm
	env GOOS=linux   GOARCH=arm64 $(ENV) go build $(FLAGS) -o bin/$(NAME)-linux-arm64

windows:
	env GOOS=windows GOARCH=386   $(ENV) go build $(FLAGS) -o bin/$(NAME)-windows-386.exe
	env GOOS=windows GOARCH=amd64 $(ENV) go build $(FLAGS) -o bin/$(NAME)-windows-amd64.exe
	env GOOS=windows GOARCH=arm   $(ENV) go build $(FLAGS) -o bin/$(NAME)-windows-arm.exe
	env GOOS=windows GOARCH=arm64 $(ENV) go build $(FLAGS) -o bin/$(NAME)-windows-arm64.exe

darwin:
	env GOOS=darwin  GOARCH=amd64 $(ENV) go build $(FLAGS) -o bin/$(NAME)-darwin-amd64
	env GOOS=darwin  GOARCH=arm64 $(ENV) go build $(FLAGS) -o bin/$(NAME)-darwin-arm64

clean:
	rm -rf bin
	go clean
