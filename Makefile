.PHONY : check-license test

test-go-cli:
	go test github.com/googleapis/gapic-generator-go/internal/gencli
	./cmd/protoc-gen-go_cli/test.sh

test-gapic:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic

golden:
	go test github.com/googleapis/gapic-generator-go/internal/gengapic -update_golden

test:
	go test ./...
	go install ./cmd/protoc-gen-go_gapic
	cd showcase; ./showcase.bash; cd ..
	./test.sh

install:
	go install ./cmd/protoc-gen-go_gapic
	go install ./cmd/protoc-gen-go_cli

update-bazel-repos:
	bazel run //:gazelle -- update-repos -from_file=go.mod -prune -to_macro=repositories.bzl%com_googleapis_gapic_generator_go_repositories
	sed -i ''  's/    "go_repository",//g' repositories.bzl

clean:
	rm -rf testdata
	rm -rf cmd/protoc-gen-go_cli/testprotos
	rm -rf cmd/protoc-gen-go_cli/testdata	
	rm -rf showcase/gen
	rm -f showcase/gapic-showcase
	rm -f showcase/showcase_grpc_service_config.json
	cd showcase; go mod edit -dropreplace github.com/googleapis/gapic-showcase
