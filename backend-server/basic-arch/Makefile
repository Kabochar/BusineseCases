.PHONY: run
run: build
	build/basic-arch.exe

.PHONY: build
build: fmt
	cd cmd/basic-arch && go build -o ../../build/

.PHONY: fmt
fmt:
	gofmt -w -s .