// Copyright 2020 Google LLC
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
	"net/http"
	"testing"
	"time"

	showcase "github.com/googleapis/gapic-showcase/client"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Clients are initialized in main_test.go.
var (
	sequenceClient     *showcase.SequenceClient
	sequenceRESTClient *showcase.SequenceClient
)

func Test_Sequence_Empty(t *testing.T) {
	defer check(t)
	seq, err := sequenceClient.CreateSequence(context.Background(), &showcasepb.CreateSequenceRequest{})
	if err != nil {
		t.Errorf("CreateSequence(empty): unexpected err %+v", err)
	}

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = sequenceClient.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()})
	if err != nil {
		t.Errorf("AttemptSequence(empty): unexpected err %+v", err)
	}

	r := seq.GetName() + "/sequenceReport"
	report, err := sequenceClient.GetSequenceReport(context.Background(), &showcasepb.GetSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("GetSequenceReport(empty): unexpected err %+v", err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != 1 {
		t.Errorf("%s: expected number of attempts to be 1 but was %d", t.Name(), len(attempts))
	}

	a := attempts[0]
	d, _ := ctx.Deadline()
	if diff := d.Sub(a.GetAttemptDeadline().AsTime()); diff > time.Millisecond {
		t.Errorf("%s: difference between server and client deadline was more than 1ms: %v", t.Name(), diff.Milliseconds())
	}
}

func Test_Sequence_Empty_DefaultDeadline(t *testing.T) {
	defer check(t)
	seq, err := sequenceClient.CreateSequence(context.Background(), &showcasepb.CreateSequenceRequest{})
	if err != nil {
		t.Errorf("%s: unexpected err %+v", t.Name(), err)
	}

	start := time.Now()
	err = sequenceClient.AttemptSequence(context.Background(), &showcasepb.AttemptSequenceRequest{Name: seq.GetName()})
	if err != nil {
		t.Errorf("%s: unexpected err %+v", t.Name(), err)
	}

	r := seq.GetName() + "/sequenceReport"
	report, err := sequenceClient.GetSequenceReport(context.Background(), &showcasepb.GetSequenceReportRequest{Name: r})
	if err != nil {
		t.Errorf("%s: unexpected err %+v", t.Name(), err)
	}

	attempts := report.GetAttempts()
	if len(attempts) != 1 {
		t.Errorf("%s: expected number of attempts to be 1 but was %d", t.Name(), len(attempts))
	}

	a := attempts[0]
	// Ensure that the default deadline of ~10s was set.
	d := start.Add(time.Second * 10)
	if diff := d.Sub(a.GetAttemptDeadline().AsTime()); diff > time.Millisecond {
		t.Errorf("%s: difference between server and client deadline was more than 1ms: %v", t.Name(), diff.Milliseconds())
	}
}

func Test_Sequence_Retry(t *testing.T) {
	defer check(t)
	responses := []*showcasepb.Sequence_Response{
		{
			Status: status.New(codes.Unavailable, "Unavailable").Proto(),
			Delay:  durationpb.New(100 * time.Millisecond),
		},
		{
			Status: status.New(codes.Unavailable, "Unavailable").Proto(),
			Delay:  durationpb.New(100 * time.Millisecond),
		},
		{
			Status: status.New(codes.Unavailable, "Unavailable").Proto(),
			Delay:  durationpb.New(100 * time.Millisecond),
		},
		{
			Status: status.New(codes.Unavailable, "Unavailable").Proto(),
			Delay:  durationpb.New(100 * time.Millisecond),
		},
		{
			Status: status.New(codes.OK, "OK").Proto(),
		},
	}

	for typ, client := range map[string]*showcase.SequenceClient{"grpc": sequenceClient, "rest": sequenceRESTClient} {
		seq, err := client.CreateSequence(context.Background(), &showcasepb.CreateSequenceRequest{
			Sequence: &showcasepb.Sequence{Responses: responses},
		})
		if err != nil {
			t.Errorf("CreateSequence(retry): unexpected err %+v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var retryOpt gax.CallOption
		bo := gax.Backoff{
			Initial:    100 * time.Millisecond,
			Max:        3000 * time.Millisecond,
			Multiplier: 2.00,
		}
		if typ == "grpc" {
			retryOpt = gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.Unavailable,
				}, bo)
			})
		} else {
			retryOpt = gax.WithRetry(func() gax.Retryer {
				return gax.OnHTTPCodes(bo, http.StatusServiceUnavailable)
			})
		}
		err = client.AttemptSequence(ctx, &showcasepb.AttemptSequenceRequest{Name: seq.GetName()}, retryOpt)
		if err != nil {
			t.Errorf("%s: unexpected AttemptSequence error %v", t.Name(), err)
		}

		r := seq.GetName() + "/sequenceReport"
		report, err := client.GetSequenceReport(context.Background(), &showcasepb.GetSequenceReportRequest{Name: r})
		if err != nil {
			t.Errorf("GetSequenceReport(retry): unexpected err %+v", err)
		}

		attempts := report.GetAttempts()
		if len(attempts) != len(responses) {
			t.Errorf("%s: expected number of attempts to be %d but was %d", t.Name(), len(responses), len(attempts))
		}

		d, _ := ctx.Deadline()
		for n, a := range attempts {
			if got, want := a.GetAttemptNumber(), int32(n); got != want {
				t.Errorf("%s: want attempt #%d but got attempt #%d", t.Name(), want, got)
			}

			// Ensure that the server-perceived attempt deadline is the same or only
			// slightly before the client-specified deadline - there seems to be
			// flaky nanosecond differences.
			//
			// TODO(noahdietz): I think there is a bug in showcase REST server because the context deadline
			// isn't being conveyed to the handler via the request.
			if diff := d.Sub(a.GetAttemptDeadline().AsTime()); typ == "grpc" && diff > time.Millisecond {
				t.Errorf("%s: difference between server and client deadline was more than 1ms: %v", t.Name(), diff.Milliseconds())
			}

			if got, want := a.GetStatus().GetCode(), responses[n].GetStatus().GetCode(); got != want {
				t.Errorf("%s: want response %v but got %v", t.Name(), want, got)
			}

			if n > 0 {
				cur, prev := a.GetAttemptDelay().AsDuration(), attempts[n-1].GetAttemptDelay().AsDuration()

				// gax.Backoff uses full jitter, so delay is not garaunteed to be monotonically increasing.
				// Thus, we can only check that it is not the same between attempts.
				//
				// TODO(noahdietz) investigate gax.Backoff jitter for signs of predictability in pseudo-random numbers.
				if cur.Milliseconds() == prev.Milliseconds() {
					t.Errorf("%s: want attempt(%d) delay: %v to differ from previous(%d): %v", t.Name(), n, cur, prev, n-1)
				}
			}
		}
	}
}
