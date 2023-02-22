build:
		go build -o thomas -ldflags "-s -w"

build-release:
		GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o thomas