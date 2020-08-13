module showcase

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.1
	github.com/googleapis/gapic-showcase v0.12.0
	github.com/googleapis/gax-go/v2 v2.0.5
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200808173500-a06252235341
	google.golang.org/grpc v1.31.0
)

replace github.com/googleapis/gapic-showcase => ./gen/github.com/googleapis/gapic-showcase
