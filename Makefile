.PHONY : check-license test

image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
	docker build -t gcr.io/gapic-images/gapic-generator-go .
	rm protoc-gen-go_gapic

check-license:
	find -name '*.go' -not -name '*.pb.go' | xargs go run utils/license.go --

test-go-cli:
	go test github.com/googleapis/gapic-generator-go/internal/gencli
	./cmd/protoc-gen-go_cli/test.sh

test-gapic:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic

test:
	go test ./...
	./utils/showcase.bash

clean:
	rm -rf testdata
	rm -rf cmd/protoc-gen-go_cli/testprotos
	rm -rf cmd/protoc-gen-go_cli/testdata	
	rm -rf showcase/gen
	rm -f showcase/gapic-showcase
