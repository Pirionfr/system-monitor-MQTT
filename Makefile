# Variables
BUILD_DIR 		:= build
FORMAT_PATHS 	:= .

PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64 arm

all:init format build

init:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)

format:
	gofmt -w -s $(FORMAT_PATHS)
	goimports -w $(FORMAT_PATHS)

dev: format build

build: init
	go build -o $(BUILD_DIR)/smm

dist:
	@for GOOS in $(PLATFORMS) ; do \
		for GOARCH in $(ARCHITECTURES) ; do \
			GOOS=$${GOOS} GOARCH=$${GOARCH} go build -o $(BUILD_DIR)/$${GOOS}/$${GOARCH}/smm; \
		done \
	done


install: build
	cp -v $(BUILD_DIR)/smm /usr/bin/smm
	mkdir /etc/smm
	cp -R default/sensors /etc/smm/sensors
	cp default/config.yml /etc/smm/
	cp default/smm.service /usr/lib/systemd/system/