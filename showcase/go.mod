module showcase_integration

require (
	cloud.google.com/go/showcase v0.0.0
	github.com/golang/protobuf v1.3.2
	github.com/googleapis/gapic-showcase v0.1.0
	google.golang.org/api v0.2.0
	google.golang.org/genproto v0.0.0-20190321212433-e79c0c59cdb5
	google.golang.org/grpc v1.19.1
)

replace cloud.google.com/go/showcase => ./gen/cloud.google.com/go/showcase
