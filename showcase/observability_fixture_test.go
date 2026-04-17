// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package showcase

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1common "go.opentelemetry.io/proto/otlp/common/v1"
	"google.golang.org/grpc"
)

// mockTraceServer implements an OTLP service that receives traces in-memory
// for assertions in tests.
type mockTraceServer struct {
	pb.UnimplementedTraceServiceServer
	mu       sync.Mutex
	requests []*pb.ExportTraceServiceRequest
}

func (s *mockTraceServer) Export(ctx context.Context, req *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, req)
	return &pb.ExportTraceServiceResponse{}, nil
}

// CapturedSpan represents a simplified OpenTelemetry span captured by the mock trace
// server for testing assertions.
type CapturedSpan struct {
	Name       string
	Scope      string
	TraceID    []byte
	Attributes map[string]any
}

func (s *mockTraceServer) getRequests() []*pb.ExportTraceServiceRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	reqs := make([]*pb.ExportTraceServiceRequest, len(s.requests))
	copy(reqs, s.requests)
	return reqs
}

func (s *mockTraceServer) GetCapturedSpans() []CapturedSpan {
	reqs := s.getRequests()
	var spans []CapturedSpan
	for _, req := range reqs {
		for _, rs := range req.ResourceSpans {
			for _, ss := range rs.ScopeSpans {
				for _, s := range ss.Spans {
					attrs := make(map[string]any)
					// Extract Resource attributes
					if rs.Resource != nil {
						for _, kv := range rs.Resource.Attributes {
							if kv.Value != nil {
								switch v := kv.Value.Value.(type) {
								case *v1common.AnyValue_StringValue:
									attrs[kv.Key] = v.StringValue
								case *v1common.AnyValue_IntValue:
									attrs[kv.Key] = v.IntValue
								case *v1common.AnyValue_BoolValue:
									attrs[kv.Key] = v.BoolValue
								case *v1common.AnyValue_DoubleValue:
									attrs[kv.Key] = v.DoubleValue
								default:
									attrs[kv.Key] = kv.Value.String()
								}
							}
						}
					}
					// Extract Span attributes
					for _, kv := range s.Attributes {
						if kv.Value != nil {
							switch v := kv.Value.Value.(type) {
							case *v1common.AnyValue_StringValue:
								attrs[kv.Key] = v.StringValue
							case *v1common.AnyValue_IntValue:
								attrs[kv.Key] = v.IntValue
							case *v1common.AnyValue_BoolValue:
								attrs[kv.Key] = v.BoolValue
							case *v1common.AnyValue_DoubleValue:
								attrs[kv.Key] = v.DoubleValue
							default:
								attrs[kv.Key] = kv.Value.String()
							}
						}
					}
					spans = append(spans, CapturedSpan{
						Name:       s.Name,
						Scope:      ss.Scope.Name,
						TraceID:    s.TraceId,
						Attributes: attrs,
					})
				}
			}
		}
	}
	return spans
}

// observabilityFixture encapsulates an in-memory OTLP gRPC server and the
// OpenTelemetry provider configurations for integration testing.
type observabilityFixture struct {
	grpcServer  *grpc.Server
	traceServer *mockTraceServer
	provider    *sdktrace.TracerProvider
}

// setupObservabilityFixture creates an in-memory OTLP trace server and configures the OTel SDK to export to it.
func setupObservabilityFixture(t *testing.T) *observabilityFixture {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	traceServer := &mockTraceServer{}
	pb.RegisterTraceServiceServer(grpcServer, traceServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("grpc server serve err: %v", err)
		}
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
	})

	ctx := context.Background()
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(lis.Addr().String()),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-app"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			t.Logf("Failed to shutdown tracer provider: %v", err)
		}
	})

	return &observabilityFixture{
		grpcServer:  grpcServer,
		traceServer: traceServer,
		provider:    tp,
	}
}
