.PHONY : test

image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
	docker build -t gcr.io/gapic-images/gapic-generator-go .
	rm protoc-gen-go_gapic

test-go-cli:
	go test github.com/googleapis/gapic-generator-go/internal/gencli
	./cmd/protoc-gen-go_cli/test.sh

test-gapic:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic

golden:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic -update_golden

test:
	go test -mod=mod ./...
	go install ./cmd/protoc-gen-go_gapic
	cd showcase && ./showcase.bash && cd .. && ./test.sh

install:
	go install ./cmd/protoc-gen-go_gapic
	go install ./cmd/protoc-gen-go_cli

update-bazel-repos:
	bazelisk run //:gazelle -- update-repos -from_file=go.mod -prune -to_macro=repositories.bzl%com_googleapis_gapic_generator_go_repositories
	sed -i ''  "s/    \"go_repository\",//g" repositories.bzl
	bazelisk run //:gazelle -- update-repos -from_file=showcase/go.mod -to_macro=repositories.bzl%com_googleapis_gapic_generator_go_repositories
	sed -i ''  "s/    \"go_repository\",//g" repositories.bzl

gazelle:
	bazelisk run //:gazelle

clean:
	rm -rf testdata
	rm -rf cmd/protoc-gen-go_cli/testprotos
	rm -rf cmd/protoc-gen-go_cli/testdata	
	rm -rf showcase/gen
	rm -f showcase/gapic-showcase
	rm -f showcase/showcase_grpc_service_config.json
	rm -f showcase/compliance_suite.json
	rm -f showcase/showcase_v1beta1.yaml
	git restore showcase/go.mod showcase/go.sum
