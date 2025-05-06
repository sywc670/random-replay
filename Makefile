BUILD=go build 

default: build

.PHONY: build
build: 
	@echo "Building"
	GOOS=windows GOARCH=amd64 $(BUILD)

.PHONY: buildgui
buildgui: 
	@echo "Building"
	GOOS=windows GOARCH=amd64 $(BUILD) -ldflags="-H windowsgui"
