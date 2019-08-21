module showcase_integration

require (
	cloud.google.com/go/showcase v0.0.0
	github.com/golang/protobuf v1.3.2
	github.com/googleapis/gapic-showcase v0.3.0
	google.golang.org/api v0.8.0
	google.golang.org/genproto v0.0.0-20190802192310-fa694d86fc64
	google.golang.org/grpc v1.22.1
)

replace cloud.google.com/go/showcase => ./gen/cloud.google.com/go/showcase
