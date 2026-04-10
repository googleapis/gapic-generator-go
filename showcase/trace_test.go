package showcase

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/auth"
	"github.com/google/go-cmp/cmp"
	showcase "github.com/googleapis/gapic-showcase/client"
	gax "github.com/googleapis/gax-go/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type dummyTokenProvider struct{}

func (d dummyTokenProvider) Token(ctx context.Context) (*auth.Token, error) {
	return &auth.Token{Value: "dummy-token"}, nil
}

func setupTracingTest(t *testing.T, enableTracing bool, transport string) (*observabilityFixture, []option.ClientOption) {
	// Reset feature cache just in case something else evaluated it
	gax.TestOnlyResetIsFeatureEnabled()
	t.Cleanup(gax.TestOnlyResetIsFeatureEnabled)
	
	if enableTracing {
		os.Setenv("GOOGLE_SDK_GO_EXPERIMENTAL_TRACING", "true")
	} else {
		os.Setenv("GOOGLE_SDK_GO_EXPERIMENTAL_TRACING", "false")
	}
	t.Cleanup(func() { os.Unsetenv("GOOGLE_SDK_GO_EXPERIMENTAL_TRACING") })

	fix := setupObservabilityFixture(t)
	oldTP := otel.GetTracerProvider()
	t.Cleanup(func() { otel.SetTracerProvider(oldTP) })
	otel.SetTracerProvider(fix.provider)

	var clientOpts []option.ClientOption
	if transport == "grpc" {
		clientOpts = []option.ClientOption{
			option.WithEndpoint("127.0.0.1:7469"),
			option.WithAuthCredentials(auth.NewCredentials(&auth.CredentialsOptions{
				TokenProvider: dummyTokenProvider{},
			})),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		}
	} else {
		clientOpts = []option.ClientOption{
			option.WithEndpoint("http://127.0.0.1:7469"),
			option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "dummy-token"})),
		}
	}

	return fix, clientOpts
}

func verifyInMemorySpan(t *testing.T, fix *observabilityFixture, expectedName string, traceID trace.TraceID, wantAttrs map[string]any, unexpectedAttrs []string) {
	t.Helper()

	// Force flush the provider to ensure traces are exported
	ctxFlush, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fix.provider.ForceFlush(ctxFlush); err != nil {
		t.Fatalf("failed to flush provider: %v", err)
	}

	// Give a little time for the export to arrive
	time.Sleep(100 * time.Millisecond)

	spans := fix.traceServer.GetCapturedSpans()
	if len(spans) == 0 {
		t.Fatalf("expected to receive trace exports, got none")
	}

	var gotSpan *CapturedSpan
	for _, s := range spans {
		if string(s.TraceID) == string(traceID[:]) && s.Name == expectedName {
			gotSpan = &s
			break
		}
	}

	if gotSpan == nil {
		t.Fatalf("did not find the expected client span")
	}

	if wantAttrs != nil {
		if _, ok := gotSpan.Attributes["gcp.client.version"]; ok {
			gotSpan.Attributes["gcp.client.version"] = "DYNAMIC"
		}
		if _, ok := gotSpan.Attributes["gcp.resource.destination.id"]; ok {
			gotSpan.Attributes["gcp.resource.destination.id"] = "DYNAMIC"
		}
		if _, ok := gotSpan.Attributes["url.full"]; ok {
			gotSpan.Attributes["url.full"] = "DYNAMIC"
		}
		if _, ok := gotSpan.Attributes["exception.message"]; ok {
			// ignore exception message as it contains arbitrary text sometimes
			gotSpan.Attributes["exception.message"] = "DYNAMIC"
		}

		// Keep only the attributes we expect for diffing
		filteredGot := make(map[string]any)
		for k, v := range gotSpan.Attributes {
			if _, expected := wantAttrs[k]; expected {
				filteredGot[k] = v
			}
		}

		if diff := cmp.Diff(wantAttrs, filteredGot); diff != "" {
			t.Errorf("Client span attributes mismatch (-want +got):\n%s", diff)
		}
	}
	
	for _, attr := range unexpectedAttrs {
		if _, ok := gotSpan.Attributes[attr]; ok {
			t.Errorf("expected attribute %q to be NOT SET, but it was present", attr)
		}
	}
}

func TestObservability_Tracing_Success(t *testing.T) {
	transports := []string{"grpc", "rest"}
	for _, transport := range transports {
		t.Run(transport, func(t *testing.T) {
			fix, clientOpts := setupTracingTest(t, true, transport)
			ctx := context.Background()
			
			var seqClient interface {
				Close() error
			}
			var err error
			
			if transport == "grpc" {
				seqClient, err = showcase.NewSequenceClient(ctx, clientOpts...)
			} else {
				seqClient, err = showcase.NewSequenceRESTClient(ctx, clientOpts...)
			}
			if err != nil {
				t.Fatalf("failed to create sequence client: %v", err)
			}
			t.Cleanup(func() { seqClient.Close() })

			ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")
			
			if transport == "grpc" {
				_ = runTracingSuccessScenario(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			} else {
				_ = runTracingSuccessScenarioREST(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			}
			span.End()
			traceID := span.SpanContext().TraceID()

			var wantAttrs map[string]any
			var unexpectedAttrs []string
			var expectedName string

			if transport == "grpc" {
				expectedName = "google.showcase.v1beta1.SequenceService/AttemptSequence"
				wantAttrs = map[string]any{
					"gcp.client.artifact":         "github.com/googleapis/gapic-showcase/client",
					"gcp.client.repo":             "googleapis/google-cloud-go",
					"gcp.client.service":          "showcase",
					"gcp.client.version":          "DYNAMIC",
					"gcp.resource.destination.id": "DYNAMIC",
					"rpc.method":                  "google.showcase.v1beta1.SequenceService/AttemptSequence",
					"rpc.response.status_code": "OK",
					"rpc.system.name":             "grpc",
					"server.address":           "127.0.0.1",
					"server.port":                 int64(7469),
					"url.domain":                  "showcase.googleapis.com",
				}
				unexpectedAttrs = []string{"status.message", "error.type", "gcp.grpc.resend_count"}
			} else {
				expectedName = "POST /v1beta1/{name=sequences/*}"
				wantAttrs = map[string]any{
					"gcp.client.artifact":         "github.com/googleapis/gapic-showcase/client",
					"gcp.client.repo":             "googleapis/google-cloud-go",
					"gcp.client.service":          "showcase",
					"gcp.client.version":          "DYNAMIC",
					"gcp.resource.destination.id": "DYNAMIC",
					"http.request.method":         "POST",
					"http.response.status_code":   int64(200),
					"server.address":              "127.0.0.1",
					"server.port":                 int64(7469),
					"url.domain":                  "showcase.googleapis.com",
					"url.full":                    "DYNAMIC",
					"url.template":                "/v1beta1/{name=sequences/*}",
				}
				unexpectedAttrs = []string{"status.message", "error.type", "exception.type", "http.request.resend_count"}
			}

			verifyInMemorySpan(t, fix, expectedName, traceID, wantAttrs, unexpectedAttrs)
		})
	}
}

func TestObservability_Tracing_Failure(t *testing.T) {
	transports := []string{"grpc", "rest"}
	for _, transport := range transports {
		t.Run(transport, func(t *testing.T) {
			fix, clientOpts := setupTracingTest(t, true, transport)
			ctx := context.Background()
			
			var seqClient interface {
				Close() error
			}
			var err error
			
			if transport == "grpc" {
				seqClient, err = showcase.NewSequenceClient(ctx, clientOpts...)
			} else {
				seqClient, err = showcase.NewSequenceRESTClient(ctx, clientOpts...)
			}
			if err != nil {
				t.Fatalf("failed to create sequence client: %v", err)
			}
			t.Cleanup(func() { seqClient.Close() })

			ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")
			if transport == "grpc" {
				_ = runTracingServerFailureScenario(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			} else {
				_ = runTracingServerFailureScenarioREST(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			}
			span.End()
			traceID := span.SpanContext().TraceID()

			var wantAttrs map[string]any
			var unexpectedAttrs []string
			var expectedName string

			if transport == "grpc" {
				expectedName = "google.showcase.v1beta1.SequenceService/AttemptSequence"
				wantAttrs = map[string]any{
					"error.type":               "NOT_FOUND",
					"exception.type":           "*status.Error",
					"gcp.client.artifact":      "github.com/googleapis/gapic-showcase/client",
					"gcp.client.repo":          "googleapis/google-cloud-go",
					"gcp.client.service":       "showcase",
					"gcp.client.version":       "DYNAMIC",
					"gcp.resource.destination.id": "DYNAMIC",
					"rpc.method":               "google.showcase.v1beta1.SequenceService/AttemptSequence",
					"rpc.response.status_code": "NOT_FOUND",
					"rpc.system.name":             "grpc",
					"server.address":           "127.0.0.1",
					"server.port":              int64(7469),
					"status.message":           "not found",
					"url.domain":               "showcase.googleapis.com",
				}
				unexpectedAttrs = []string{}
			} else {
				expectedName = "POST /v1beta1/{name=sequences/*}"
				wantAttrs = map[string]any{
					"error.type":               "404",
					"gcp.client.artifact":      "github.com/googleapis/gapic-showcase/client",
					"gcp.client.repo":          "googleapis/google-cloud-go",
					"gcp.client.service":       "showcase",
					"gcp.client.version":       "DYNAMIC",
					"gcp.resource.destination.id": "DYNAMIC",
					"http.request.method":      "POST",
					"http.response.status_code": int64(404),
					"server.address":              "127.0.0.1",
					"server.port":              int64(7469),
					"status.message":           "not found",
					"url.domain":               "showcase.googleapis.com",
					"url.full":                 "DYNAMIC",
					"url.template":             "/v1beta1/{name=sequences/*}",
				}
				unexpectedAttrs = []string{}
			}

			verifyInMemorySpan(t, fix, expectedName, traceID, wantAttrs, unexpectedAttrs)
		})
	}
}

func TestObservability_Tracing_ClientFailure(t *testing.T) {
	transports := []string{"grpc", "rest"}
	for _, transport := range transports {
		t.Run(transport, func(t *testing.T) {
			fix, clientOpts := setupTracingTest(t, true, transport)
			ctx := context.Background()
			
			var seqClient interface {
				Close() error
			}
			var err error
			
			if transport == "grpc" {
				seqClient, err = showcase.NewSequenceClient(ctx, clientOpts...)
			} else {
				seqClient, err = showcase.NewSequenceRESTClient(ctx, clientOpts...)
			}
			if err != nil {
				t.Fatalf("failed to create sequence client: %v", err)
			}
			t.Cleanup(func() { seqClient.Close() })

			ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")
			if transport == "grpc" {
				_ = runTracingClientFailureScenario(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			} else {
				_ = runTracingClientFailureScenarioREST(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			}
			span.End()
			traceID := span.SpanContext().TraceID()

			var wantAttrs map[string]any
			var unexpectedAttrs []string
			var expectedName string

			if transport == "grpc" {
				expectedName = "google.showcase.v1beta1.SequenceService/AttemptSequence"
				wantAttrs = map[string]any{
					"error.type":               "CLIENT_TIMEOUT",
					"exception.type":           "*status.Error",
					"gcp.client.artifact":      "github.com/googleapis/gapic-showcase/client",
					"gcp.client.repo":          "googleapis/google-cloud-go",
					"gcp.client.service":       "showcase",
					"gcp.client.version":       "DYNAMIC",
					"gcp.resource.destination.id": "DYNAMIC",
					"rpc.method":               "google.showcase.v1beta1.SequenceService/AttemptSequence",
					"rpc.system.name":             "grpc",
					"server.address":              "127.0.0.1",
					"server.port":              int64(7469),
					"status.message":           "context deadline exceeded",
					"url.domain":               "showcase.googleapis.com",
				}
				unexpectedAttrs = []string{"rpc.response.status_code"}
			} else {
				expectedName = "POST /v1beta1/{name=sequences/*}"
				wantAttrs = map[string]any{
					"error.type":               "context.deadlineExceededError",
					"exception.type":           "context.deadlineExceededError",
					"gcp.client.artifact":      "github.com/googleapis/gapic-showcase/client",
					"gcp.client.repo":          "googleapis/google-cloud-go",
					"gcp.client.service":       "showcase",
					"gcp.client.version":       "DYNAMIC",
					"gcp.resource.destination.id": "DYNAMIC",
					"http.request.method":         "POST",
					"server.address":              "127.0.0.1",
					"server.port":              int64(7469),
					"url.domain":               "showcase.googleapis.com",
					"url.full":                 "DYNAMIC",
					"url.template":             "/v1beta1/{name=sequences/*}",
				}
				unexpectedAttrs = []string{"http.response.status_code"}
			}

			verifyInMemorySpan(t, fix, expectedName, traceID, wantAttrs, unexpectedAttrs)
		})
	}
}

func TestObservability_Tracing_Disablement(t *testing.T) {
	transports := []string{"grpc", "rest"}
	for _, transport := range transports {
		t.Run(transport, func(t *testing.T) {
			fix, clientOpts := setupTracingTest(t, false, transport)
			ctx := context.Background()
			
			var seqClient interface {
				Close() error
			}
			var err error
			if transport == "grpc" {
				seqClient, err = showcase.NewSequenceClient(ctx, clientOpts...)
			} else {
				seqClient, err = showcase.NewSequenceRESTClient(ctx, clientOpts...)
			}
			if err != nil {
				t.Fatalf("failed to create sequence client: %v", err)
			}
			t.Cleanup(func() { seqClient.Close() })
			
			ctxSpan, span := otel.Tracer("test-tracer").Start(context.Background(), "APP")
			if transport == "grpc" {
				runTracingDisablementScenario(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			} else {
				runTracingDisablementScenarioREST(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			}
			span.End()

			ctxFlush, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := fix.provider.ForceFlush(ctxFlush); err != nil {
				t.Fatalf("failed to flush provider: %v", err)
			}

			time.Sleep(100 * time.Millisecond)

			spans := fix.traceServer.GetCapturedSpans()
			
			for _, s := range spans {
				if _, ok := s.Attributes["gcp.client.artifact"]; ok {
					t.Errorf("found gcp.client.artifact attribute, but tracing telemetry should be disabled")
				}
			}
		})
	}
}

func TestObservability_Tracing_Retry(t *testing.T) {
	transports := []string{"grpc", "rest"}
	for _, transport := range transports {
		t.Run(transport, func(t *testing.T) {
			fix, clientOpts := setupTracingTest(t, true, transport)
			ctx := context.Background()
			
			var seqClient interface {
				Close() error
			}
			var err error
			
			if transport == "grpc" {
				seqClient, err = showcase.NewSequenceClient(ctx, clientOpts...)
			} else {
				seqClient, err = showcase.NewSequenceRESTClient(ctx, clientOpts...)
			}
			if err != nil {
				t.Fatalf("failed to create sequence client: %v", err)
			}
			t.Cleanup(func() { seqClient.Close() })

			ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")
			if transport == "grpc" {
				_ = runTracingRetryScenario(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			} else {
				_ = runTracingRetryScenarioREST(ctxSpan, t, seqClient.(*showcase.SequenceClient))
			}
			span.End()
			traceID := span.SpanContext().TraceID()

			ctxFlush, cancelFlush := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFlush()
			if err := fix.provider.ForceFlush(ctxFlush); err != nil {
				t.Fatalf("failed to flush provider: %v", err)
			}

			time.Sleep(100 * time.Millisecond)

			spans := fix.traceServer.GetCapturedSpans()
			var attemptSpans []CapturedSpan
			expectedName := "google.showcase.v1beta1.SequenceService/AttemptSequence"
			if transport == "rest" {
				expectedName = "POST /v1beta1/{name=sequences/*}"
			}
			for _, s := range spans {
				if string(s.TraceID) == string(traceID[:]) && s.Name == expectedName {
					attemptSpans = append(attemptSpans, s)
				}
			}

			if len(attemptSpans) != 4 {
				t.Errorf("expected 4 attempt spans (3 failures + 1 success), got %d", len(attemptSpans))
			}
			
			// Verify last span has correct attributes
			if len(attemptSpans) > 0 {
				lastSpan := attemptSpans[len(attemptSpans)-1]
				if transport == "rest" {
					if resend, ok := lastSpan.Attributes["http.request.resend_count"]; !ok || resend.(int64) != 3 {
						t.Errorf("expected http.request.resend_count to be 3, got %v", resend)
					}
				} else if transport == "grpc" {
					if resend, ok := lastSpan.Attributes["gcp.grpc.resend_count"]; !ok || resend.(int64) != 3 {
						t.Errorf("expected gcp.grpc.resend_count to be 3, got %v", resend)
					}
				}
			}
		})
	}
}
