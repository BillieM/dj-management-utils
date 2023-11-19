# build all binaries
build: build-cli build-gui

# build all cli binaries
build-cli: build-cli-windows build-cli-linux build-cli-macos

# cli windows binaries
build-cli-windows: build-cli-windows-amd64 build-cli-windows-386

build-cli-windows-amd64:
	GOOS=windows GOARCH=amd64 go build ./cmd/cli -o bin/seren-cli-windows-amd64.exe

build-cli-windows-386:
	GOOS=windows GOARCH=386 go build ./cmd/cli -o bin/seren-cli-windows-x386.exe

# cli linux binaries
build-cli-linux: build-cli-linux-amd64 build-cli-linux-arm

build-cli-linux-386:
	GOOS=linux GOARCH=386 go build ./cmd/cli -o bin/seren-cli-linux-x386

build-cli-linux-amd64:
	GOOS=linux GOARCH=amd64 go build ./cmd/cli -o bin/seren-cli-linux-amd64

build-cli-linux-arm64:
	GOOS=linux GOARCH=arm64 go build ./cmd/cli -o bin/seren-cli-linux-arm64

# cli macos binaries
build-cli-macos: build-cli-macos-amd64 build-cli-macos-arm64

build-cli-macos-amd64:
	GOOS=linux GOARCH=amd64 go build ./cmd/cli -o bin/seren-cli-linux-amd64

build-cli-macos-arm64:
	GOOS=linux GOARCH=arm64 go build ./cmd/cli -o bin/seren-cli-linux-arm64

# build all gui binaries
build-gui: build-gui-windows build-gui-linux build-gui-macos

# gui windows binaries
build-gui-windows: build-gui-windows-amd64 build-gui-windows-386

build-gui-windows-amd64:
	GOOS=windows GOARCH=amd64 go build ./cmd/gui -o bin/seren-gui-windows-amd64.exe

build-gui-windows-386:
	GOOS=windows GOARCH=386 go build ./cmd/gui -o bin/seren-gui-windows-x386.exe

# gui linux binaries
build-gui-linux: build-gui-linux-amd64 build-gui-linux-arm

build-gui-linux-386:
	GOOS=linux GOARCH=386 go build ./cmd/gui -o bin/seren-gui-linux-x386

build-gui-linux-amd64:
	GOOS=linux GOARCH=amd64 go build ./cmd/gui -o bin/seren-gui-linux-amd64

build-gui-linux-arm64:
	GOOS=linux GOARCH=arm64 go build ./cmd/gui -o bin/seren-gui-linux-arm64

# gui macos binaries
build-gui-macos: build-gui-macos-amd64 build-gui-macos-arm64

build-gui-macos-amd64:
	GOOS=linux GOARCH=amd64 go build ./cmd/gui -o bin/seren-gui-linux-amd64

build-gui-macos-arm64:
	GOOS=linux GOARCH=arm64 go build ./cmd/gui -o bin/seren-gui-linux-arm64

test:
	go test ./...