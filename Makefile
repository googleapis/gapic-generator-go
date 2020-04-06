.PHONY : check-license test

image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
	docker build -t gcr.io/gapic-images/gapic-generator-go .
	rm protoc-gen-go_gapic

test-go-cli:
	go test github.com/googleapis/gapic-generator-go/internal/gencli

test-gapic:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic

golden:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic -update_golden

test:
	go test ./...
	go install ./cmd/protoc-gen-go_gapic
	cd showcase; ./showcase.bash; cd ..

install:
	go install ./cmd/protoc-gen-go_gapic
	go install ./cmd/protoc-gen-go_cli

clean:
	rm -rf testdata
	rm -rf cmd/protoc-gen-go_cli/testprotos
	rm -rf cmd/protoc-gen-go_cli/testdata	
	rm -rf showcase/gen
	rm -f showcase/gapic-showcase
	rm -f showcase/showcase_grpc_service_config.json
	cd showcase; go mod edit -dropreplace github.com/googleapis/gapic-showcase
