module showcase_integration

require (
	cloud.google.com/go/showcase v0.0.0
	github.com/golang/protobuf v1.3.2
	github.com/googleapis/gapic-showcase v0.5.0
	google.golang.org/api v0.9.0
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.23.0
)

replace cloud.google.com/go/showcase => ./gen/cloud.google.com/go/showcase

go 1.13
