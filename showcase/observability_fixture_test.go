package showcase

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	pblog "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	olog "go.opentelemetry.io/proto/otlp/logs/v1"
	pbmetric "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1common "go.opentelemetry.io/proto/otlp/common/v1"
	"google.golang.org/grpc"
)

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

type mockMetricServer struct {
	pbmetric.UnimplementedMetricsServiceServer
	mu       sync.Mutex
	requests []*pbmetric.ExportMetricsServiceRequest
}

func (s *mockMetricServer) Export(ctx context.Context, req *pbmetric.ExportMetricsServiceRequest) (*pbmetric.ExportMetricsServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, req)
	return &pbmetric.ExportMetricsServiceResponse{}, nil
}

type CapturedMetric struct {
	Name       string
	Scope      string
	Attributes map[string]any
	DataPoints []float64 // Simplifying for histograms
}

func (s *mockMetricServer) getRequests() []*pbmetric.ExportMetricsServiceRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	reqs := make([]*pbmetric.ExportMetricsServiceRequest, len(s.requests))
	copy(reqs, s.requests)
	return reqs
}

func (s *mockMetricServer) GetCapturedMetrics() []CapturedMetric {
	reqs := s.getRequests()
	var metrics []CapturedMetric
	for _, req := range reqs {
		for _, rm := range req.ResourceMetrics {
			for _, sm := range rm.ScopeMetrics {
				for _, m := range sm.Metrics {
					if hist := m.GetHistogram(); hist != nil {
						for _, dp := range hist.DataPoints {
							cmDp := CapturedMetric{
								Name:  m.Name,
								Scope: sm.Scope.Name,
							}
							var dps []float64
							if dp.Sum != nil {
								dps = append(dps, *dp.Sum)
							} else {
								dps = append(dps, 0.0)
							}
							cmDp.DataPoints = dps

							var attrsMap = make(map[string]any)
							// Extract Scope Attributes
							for _, kv := range sm.Scope.Attributes {
								if kv.Value != nil {
									switch v := kv.Value.Value.(type) {
									case *v1common.AnyValue_StringValue:
										attrsMap[kv.Key] = v.StringValue
									case *v1common.AnyValue_IntValue:
										attrsMap[kv.Key] = v.IntValue
									}
								}
							}
							for _, kv := range dp.Attributes {
								if kv.Value != nil {
									switch v := kv.Value.Value.(type) {
									case *v1common.AnyValue_StringValue:
										attrsMap[kv.Key] = v.StringValue
									case *v1common.AnyValue_IntValue:
										attrsMap[kv.Key] = v.IntValue
									}
								}
							}
							cmDp.Attributes = attrsMap
							metrics = append(metrics, cmDp)
						}
					}
				}
			}
		}
	}
	return metrics
}

type mockLogServer struct {
	pblog.UnimplementedLogsServiceServer
	mu       sync.Mutex
	requests []*pblog.ExportLogsServiceRequest
}

func (s *mockLogServer) Export(ctx context.Context, req *pblog.ExportLogsServiceRequest) (*pblog.ExportLogsServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, req)
	return &pblog.ExportLogsServiceResponse{}, nil
}

type CapturedLog struct {
	Body       string
	Scope      string
	Severity   olog.SeverityNumber
	TraceID    []byte
	Attributes map[string]any
}

func (s *mockLogServer) getRequests() []*pblog.ExportLogsServiceRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	reqs := make([]*pblog.ExportLogsServiceRequest, len(s.requests))
	copy(reqs, s.requests)
	return reqs
}

func (s *mockLogServer) GetCapturedLogs() []CapturedLog {
	reqs := s.getRequests()
	var logs []CapturedLog
	for _, req := range reqs {
		for _, rl := range req.ResourceLogs {
			for _, sl := range rl.ScopeLogs {
				for _, l := range sl.LogRecords {
					cl := CapturedLog{
						Scope:    sl.Scope.Name,
						Severity: l.SeverityNumber,
						TraceID:  l.TraceId,
					}
					if l.Body != nil {
						cl.Body = l.Body.GetStringValue()
					}

					attrs := make(map[string]any)
					for _, kv := range l.Attributes {
						if kv.Value != nil {
							switch v := kv.Value.Value.(type) {
							case *v1common.AnyValue_StringValue:
								attrs[kv.Key] = v.StringValue
							case *v1common.AnyValue_IntValue:
								attrs[kv.Key] = v.IntValue
							case *v1common.AnyValue_BoolValue:
								attrs[kv.Key] = v.BoolValue
							}
						}
					}
					cl.Attributes = attrs
					logs = append(logs, cl)
				}
			}
		}
	}
	return logs
}

type observabilityFixture struct {
	grpcServer     *grpc.Server
	traceServer    *mockTraceServer
	metricServer   *mockMetricServer
	logServer      *mockLogServer
	provider       *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	loggerProvider *sdklog.LoggerProvider
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
	metricServer := &mockMetricServer{}
	logServer := &mockLogServer{}
	pb.RegisterTraceServiceServer(grpcServer, traceServer)
	pbmetric.RegisterMetricsServiceServer(grpcServer, metricServer)
	pblog.RegisterLogsServiceServer(grpcServer, logServer)

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

	metricExp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(lis.Addr().String()),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to create metric exporter: %v", err)
	}

	logExp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(lis.Addr().String()),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to create log exporter: %v", err)
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

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp, sdkmetric.WithInterval(10*time.Hour))),
		sdkmetric.WithResource(res),
	)

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		sdklog.WithResource(res),
	)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			t.Logf("Failed to shutdown tracer provider: %v", err)
		}
		if err := mp.Shutdown(ctx); err != nil {
			t.Logf("Failed to shutdown meter provider: %v", err)
		}
		if err := lp.Shutdown(ctx); err != nil {
			t.Logf("Failed to shutdown logger provider: %v", err)
		}
	})

	return &observabilityFixture{
		grpcServer:     grpcServer,
		traceServer:    traceServer,
		metricServer:   metricServer,
		logServer:      logServer,
		provider:       tp,
		meterProvider:  mp,
		loggerProvider: lp,
	}
}
