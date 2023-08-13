.PHONY: run
run: build
	build/signaling.exe

.PHONY: build
build: fmt
	cd cmd/signaling && go build -o ../../build/

.PHONY: fmt
fmt:
	gofmt -w -s .