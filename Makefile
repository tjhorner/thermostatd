.PHONY: dist dist-win dist-macos dist-linux ensure-dist-dir build install uninstall

GOBUILD=go build -ldflags="-s -w"
INSTALLPATH=/usr/local/bin

PROJECT_NAME=thermostatd

ensure-dist-dir:
	@- mkdir -p dist

dist-win: ensure-dist-dir
	# Build for Windows x64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/$(PROJECT_NAME)-windows-amd64.exe .

dist-macos: ensure-dist-dir
	# Build for macOS x64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dist/$(PROJECT_NAME)-darwin-amd64 .

dist-linux: ensure-dist-dir
	# Build for Linux x64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/$(PROJECT_NAME)-linux-amd64 .

dist-arm: ensure-dist-dir
	# Build for Linux ARMv7
	GOOS=linux GOARCH=arm $(GOBUILD) -o dist/$(PROJECT_NAME)-linux-arm .

run-remote: dist-arm
	scp dist/thermostatd-linux-arm pi@192.168.1.49:/home/pi/thermostatd
	ssh -t pi@192.168.1.49 /home/pi/thermostatd

deploy-remote: dist-arm
	ssh -t pi@192.168.1.49 sudo systemctl stop thermostatd
	scp dist/thermostatd-linux-arm pi@192.168.1.49:/home/pi/thermostatd
	ssh -t pi@192.168.1.49 sudo systemctl start thermostatd

dist: dist-win dist-macos dist-linux dist-arm

build:
	@- mkdir -p bin
	$(GOBUILD) -o bin/$(PROJECT_NAME) .
	@- chmod +x bin/$(PROJECT_NAME)

install: build
	mv bin/$(PROJECT_NAME) $(INSTALLPATH)/$(PROJECT_NAME)
	@- rm -rf bin
	@echo "$(PROJECT_NAME) was installed to $(INSTALLPATH)/$(PROJECT_NAME). Run make uninstall to get rid of it, or just remove the binary yourself."

uninstall:
	rm $(INSTALLPATH)/$(PROJECT_NAME)

run:
	@- go run .