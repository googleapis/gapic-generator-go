module showcase

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.4
	github.com/googleapis/gapic-showcase v0.12.0
	github.com/googleapis/gax-go/v2 v2.0.5
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20201201144952-b05cb90ed32e
	google.golang.org/grpc v1.33.2
)

replace github.com/googleapis/gapic-showcase => ./gen/github.com/googleapis/gapic-showcase
