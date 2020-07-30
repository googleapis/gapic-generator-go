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
package grpc_service_config

import (
	"testing"

	duration "github.com/golang/protobuf/ptypes/duration"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
	code "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
)

func Test_Parse(t *testing.T) {
	in := ServiceConfig{
		MethodConfig: []*MethodConfig{
			{
				Name: []*MethodConfig_Name{
					{
						Service: "bar.FooService",
						Method:  "Zip",
					},
				},
				MaxRequestMessageBytes:  &wrappers.UInt32Value{Value: 123456},
				MaxResponseMessageBytes: &wrappers.UInt32Value{Value: 123456},
				RetryOrHedgingPolicy: &MethodConfig_RetryPolicy_{
					RetryPolicy: &MethodConfig_RetryPolicy{
						InitialBackoff:    &duration.Duration{Nanos: 100000000},
						MaxBackoff:        &duration.Duration{Seconds: 60},
						BackoffMultiplier: 1.3,
						RetryableStatusCodes: []code.Code{
							code.Code_UNKNOWN,
						},
					},
				},
				Timeout: &duration.Duration{Seconds: 30},
			},
			{
				Name: []*MethodConfig_Name{
					{
						Service: "bar.FooService",
					},
				},
				MaxRequestMessageBytes:  &wrappers.UInt32Value{Value: 654321},
				MaxResponseMessageBytes: &wrappers.UInt32Value{Value: 654321},
				RetryOrHedgingPolicy: &MethodConfig_RetryPolicy_{
					RetryPolicy: &MethodConfig_RetryPolicy{
						InitialBackoff:    &duration.Duration{Nanos: 10000000},
						MaxBackoff:        &duration.Duration{Seconds: 7},
						BackoffMultiplier: 1.1,
						RetryableStatusCodes: []code.Code{
							code.Code_UNKNOWN,
						},
					},
				},
				Timeout: &duration.Duration{Seconds: 60},
			},
		},
	}

	want := Config{
		policies: map[string]*MethodConfig_RetryPolicy{
			"bar.FooService": {
				InitialBackoff:    &duration.Duration{Nanos: 10000000},
				MaxBackoff:        &duration.Duration{Seconds: 7},
				BackoffMultiplier: 1.1,
				RetryableStatusCodes: []code.Code{
					code.Code_UNKNOWN,
				},
			},
			"bar.FooService.Zip": {
				InitialBackoff:    &duration.Duration{Nanos: 100000000},
				MaxBackoff:        &duration.Duration{Seconds: 60},
				BackoffMultiplier: 1.3,
				RetryableStatusCodes: []code.Code{
					code.Code_UNKNOWN,
				},
			},
		},
		timeouts: map[string]*duration.Duration{
			"bar.FooService":     {Seconds: 60},
			"bar.FooService.Zip": {Seconds: 30},
		},
		reqLimits: map[string]int{
			"bar.FooService":     654321,
			"bar.FooService.Zip": 123456,
		},
		resLimits: map[string]int{
			"bar.FooService":     654321,
			"bar.FooService.Zip": 123456,
		},
	}

	got := New(in)
	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal), cmp.AllowUnexported(Config{})); diff != "" {
		t.Errorf("%s: %s", t.Name(), diff)
	}
}

func Test_GetTimeout(t *testing.T) {
	s := "bar.FooService"
	m := "Zip"
	mFQN := s + "." + m
	c := Config{
		timeouts: map[string]*duration.Duration{
			s:    {Seconds: 60},
			mFQN: {Seconds: 30},
		},
	}

	want := ToMillis(c.timeouts[s])
	if got, ok := c.GetTimeout(s, "dne"); !ok || got != want {
		t.Errorf("%s: expected %d got %d", t.Name(), want, got)
	}

	want = ToMillis(c.timeouts[mFQN])
	if got, ok := c.GetTimeout(s, m); !ok || got != want {
		t.Errorf("%s: expected %d got %d", t.Name(), want, got)
	}

	if got, ok := c.GetTimeout("dne", "dne"); ok {
		t.Errorf("%s: expected !ok got %d", t.Name(), got)
	}
}

func Test_GetRetryPolicy(t *testing.T) {
	s := "bar.FooService"
	m := "Zip"
	mFQN := s + "." + m
	c := Config{
		policies: map[string]*MethodConfig_RetryPolicy{
			s: {
				InitialBackoff:    &duration.Duration{Nanos: 10000000},
				MaxBackoff:        &duration.Duration{Seconds: 7},
				BackoffMultiplier: 1.1,
				RetryableStatusCodes: []code.Code{
					code.Code_UNKNOWN,
				},
			},
			mFQN: {
				InitialBackoff:    &duration.Duration{Nanos: 100000000},
				MaxBackoff:        &duration.Duration{Seconds: 60},
				BackoffMultiplier: 1.3,
				RetryableStatusCodes: []code.Code{
					code.Code_UNKNOWN,
				},
			},
		},
	}

	want := c.policies[s]
	if got, ok := c.GetRetryPolicy(s, "dne"); !ok || !cmp.Equal(got, want, cmp.Comparer(proto.Equal)) {
		t.Errorf("%s: expected %v got %v", t.Name(), want, got)
	}

	want = c.policies[mFQN]
	if got, ok := c.GetRetryPolicy(s, m); !ok || !cmp.Equal(got, want, cmp.Comparer(proto.Equal)) {
		t.Errorf("%s: expected %v got %v", t.Name(), want, got)
	}

	if got, ok := c.GetRetryPolicy("dne", "dne"); ok {
		t.Errorf("%s: expected !ok got %v", t.Name(), got)
	}
}

func Test_GetRequestLimit(t *testing.T) {
	s := "bar.FooService"
	m := "Zip"
	mFQN := s + "." + m
	c := Config{
		reqLimits: map[string]int{
			s:    654321,
			mFQN: 123456,
		},
	}

	want := c.reqLimits[s]
	if got, ok := c.GetRequestLimit(s, "dne"); !ok || got != want {
		t.Errorf("%s: expected %d got %d", t.Name(), want, got)
	}

	want = c.reqLimits[mFQN]
	if got, ok := c.GetRequestLimit(s, m); !ok || got != want {
		t.Errorf("%s: expected %d got %d", t.Name(), want, got)
	}

	if got, ok := c.GetRequestLimit("dne", "dne"); ok {
		t.Errorf("%s: expected !ok got %d", t.Name(), got)
	}
}

func Test_GetResponseLimit(t *testing.T) {
	s := "bar.FooService"
	m := "Zip"
	mFQN := s + "." + m
	c := Config{
		resLimits: map[string]int{
			s:    654321,
			mFQN: 123456,
		},
	}

	want := c.resLimits[s]
	if got, ok := c.GetResponseLimit(s, "dne"); !ok || got != want {
		t.Errorf("%s: expected %d got %d", t.Name(), want, got)
	}

	want = c.resLimits[mFQN]
	if got, ok := c.GetResponseLimit(s, m); !ok || got != want {
		t.Errorf("%s: expected %d got %d", t.Name(), want, got)
	}

	if got, ok := c.GetResponseLimit("dne", "dne"); ok {
		t.Errorf("%s: expected !ok got %d", t.Name(), got)
	}
}
