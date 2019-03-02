module showcase_integration

require (
	cloud.google.com/go/showcase v0.0.0
	github.com/golang/protobuf v1.3.0
	github.com/googleapis/gapic-showcase v0.0.12
	github.com/googleapis/gax-go v0.0.0-20181219185031-c8a15bac9b9f // indirect
	github.com/googleapis/gax-go/v2 v2.0.3 // indirect
	google.golang.org/api v0.1.0
	google.golang.org/genproto v0.0.0-20190219182410-082222b4a5c5
	google.golang.org/grpc v1.18.0
)

replace cloud.google.com/go/showcase => ./gen/cloud.google.com/go/showcase
