.PHONY : check-license test

image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
	docker build -t gcr.io/gapic-images/gapic-generator-go .
	rm protoc-gen-go_gapic

check-license:
	find -name '*.go' -not -name '*.pb.go' | xargs go run utils/license.go --

test-gcli:
	go test github.com/googleapis/gapic-generator-go/internal/gencli
	./cmd/protoc-gen-gcli/test.sh

test-gapic:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic

test:
	go test ./...
	./utils/showcase.bash

clean:
	rm -rf showcase-testdir
	rm -rf testdata
	rm -rf cmd/protoc-gen-gcli/testdata	