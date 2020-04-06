module showcase

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/google/go-cmp v0.4.0
	github.com/googleapis/gapic-showcase v0.7.0
	google.golang.org/api v0.17.0
	google.golang.org/genproto v0.0.0-20200205142000-a86caf926a67
	google.golang.org/grpc v1.27.1
)

replace github.com/googleapis/gapic-showcase => ./gen/github.com/googleapis/gapic-showcase
