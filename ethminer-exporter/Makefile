.PHONY: build docker-build test clean all fmt

all: build docker-build

fmt:
	go fmt ./.../

test:
	go test -cover -count=1 -v ./.../

build/_output/bin/ethxporter:
		mkdir -p build/_output/bin/
		go build -a -tags netgo -mod=vendor -o build/_output/bin/ethxporter main.go

build: build/_output/bin/ethxporter

docker-build:
	docker build . -f build/Dockerfile -t ethxporter

clean:
	rm -r build/_output/bin/ethxporter