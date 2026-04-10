package showcase

import (
	"context"
	"testing"
	"time"

	showcase "github.com/googleapis/gapic-showcase/client"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func runTracingSuccessScenario(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.OK, "OK").Proto()},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	err = seqClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err != nil {
		t.Fatalf("AttemptSequence RPC failed: %v", err)
	}

	return seq
}

func runTracingServerFailureScenario(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.NotFound, "not found").Proto()},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	err = seqClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	return seq
}

func runTracingClientFailureScenario(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{
			Status: status.New(codes.OK, "OK").Proto(),
			Delay:  durationpb.New(1 * time.Second),
		},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")

	timeoutCtx, cancelTimeout := context.WithTimeout(ctxSpan, 1*time.Millisecond)
	defer cancelTimeout()

	err = seqClient.AttemptSequence(timeoutCtx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	span.End()

	return seq
}

func runTracingRetryScenario(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.Unavailable, "Unavailable").Proto()},
		{Status: status.New(codes.Unavailable, "Unavailable").Proto()},
		{Status: status.New(codes.Unavailable, "Unavailable").Proto()},
		{Status: status.New(codes.OK, "OK").Proto()},
	}

	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")

	retryCtx, cancel := context.WithTimeout(ctxSpan, 5*time.Second)
	defer cancel()

	bo := gax.Backoff{
		Initial:    10 * time.Millisecond,
		Max:        100 * time.Millisecond,
		Multiplier: 2.00,
	}
	retryOpt := gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, bo)
	})

	opts := append(seqClient.CallOptions.AttemptSequence, retryOpt)
	err = seqClient.AttemptSequence(retryCtx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, opts...)
	if err != nil {
		t.Fatalf("AttemptSequence failed: %v", err)
	}
	span.End()

	return seq
}

func runTracingDisablementScenario(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.OK, "OK").Proto()},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	err = seqClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err != nil {
		t.Fatalf("AttemptSequence RPC failed: %v", err)
	}
}
func runTracingSuccessScenarioREST(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.OK, "OK").Proto()},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	err = seqClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err != nil {
		t.Fatalf("AttemptSequence RPC failed: %v", err)
	}

	return seq
}

func runTracingServerFailureScenarioREST(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.NotFound, "not found").Proto()},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	err = seqClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	return seq
}

func runTracingClientFailureScenarioREST(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{
			Status: status.New(codes.OK, "OK").Proto(),
			Delay:  durationpb.New(1 * time.Second),
		},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")

	timeoutCtx, cancelTimeout := context.WithTimeout(ctxSpan, 1*time.Millisecond)
	defer cancelTimeout()

	err = seqClient.AttemptSequence(timeoutCtx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	span.End()

	return seq
}

func runTracingRetryScenarioREST(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) *showcasepb.Sequence {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.Unavailable, "Unavailable").Proto()},
		{Status: status.New(codes.Unavailable, "Unavailable").Proto()},
		{Status: status.New(codes.Unavailable, "Unavailable").Proto()},
		{Status: status.New(codes.OK, "OK").Proto()},
	}

	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	ctxSpan, span := otel.Tracer("test-tracer").Start(ctx, "APP")

	retryCtx, cancel := context.WithTimeout(ctxSpan, 5*time.Second)
	defer cancel()

	bo := gax.Backoff{
		Initial:    10 * time.Millisecond,
		Max:        100 * time.Millisecond,
		Multiplier: 2.00,
	}
	retryOpt := gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, bo)
	})

	opts := append(seqClient.CallOptions.AttemptSequence, retryOpt)
	err = seqClient.AttemptSequence(retryCtx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, opts...)
	if err != nil {
		t.Fatalf("AttemptSequence failed: %v", err)
	}
	span.End()

	return seq
}

func runTracingDisablementScenarioREST(ctx context.Context, t *testing.T, seqClient *showcase.SequenceClient) {
	responses := []*showcasepb.Sequence_Response{
		{Status: status.New(codes.OK, "OK").Proto()},
	}
	seq, err := seqClient.CreateSequence(ctx, &showcasepb.CreateSequenceRequest{
		Sequence: &showcasepb.Sequence{Responses: responses},
	})
	if err != nil {
		t.Fatalf("CreateSequence failed: %v", err)
	}

	err = seqClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, seqClient.CallOptions.AttemptSequence...)
	if err != nil {
		t.Fatalf("AttemptSequence RPC failed: %v", err)
	}
}
