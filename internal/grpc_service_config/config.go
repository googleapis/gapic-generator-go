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
	"io"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	duration "google.golang.org/protobuf/types/known/durationpb"
)

// Config represents parsed mapping of the gRPC ServiceConfig contents
// with methods for increased accessibility.
type Config struct {
	policies  map[string]*MethodConfig_RetryPolicy
	timeouts  map[string]*duration.Duration
	reqLimits map[string]int
	resLimits map[string]int
}

// New traverses the given gRPC ServiceConfig into more accessible constructs
// mapped by the names to the specific config values. Use the accessors on the
// resulting Config to retrieve values for a Service or Method.
func New(in io.Reader) (Config, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return Config{}, err
	}

	c := ServiceConfig{}
	err = protojson.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}

	policies := map[string]*MethodConfig_RetryPolicy{}
	timeouts := map[string]*duration.Duration{}
	reqLimits := map[string]int{}
	resLimits := map[string]int{}

	// gather retry policies from MethodConfigs
	for _, mc := range c.GetMethodConfig() {
		for _, name := range mc.GetName() {
			n := name.GetService()

			// individual method config
			if name.GetMethod() != "" {
				n = n + "." + name.GetMethod()
			}

			policies[n] = mc.GetRetryPolicy()

			if maxReq := mc.GetMaxRequestMessageBytes(); maxReq != nil {
				reqLimits[n] = int(maxReq.GetValue())
			}

			if maxRes := mc.GetMaxResponseMessageBytes(); maxRes != nil {
				resLimits[n] = int(maxRes.GetValue())
			}

			if timeout := mc.GetTimeout(); timeout != nil {
				timeouts[n] = timeout
			}
		}
	}

	return Config{
		policies:  policies,
		timeouts:  timeouts,
		reqLimits: reqLimits,
		resLimits: resLimits,
	}, nil
}

// RetryPolicy returns the retryPolicy and a presence flag for the
// given fully-qualified Service and simple Method names. A config assignment
// for a specific Method takes precendence over a Service-level assignment.
func (c Config) RetryPolicy(s, m string) (*MethodConfig_RetryPolicy, bool) {
	// Favor the policy defined for a fully-qualified Method name.
	policy, ok := c.policies[s+"."+m]
	if ok {
		return policy, ok
	}

	// Fallback on the policy defined for an entire Service
	policy, ok = c.policies[s]
	return policy, ok
}

// Timeout returns the timeout in milliseconds and a presence flag for the given
// fully-qualified Service and simple Method names. A config assignment for the
// specific Method takes precendence over a Service-level assignment.
func (c Config) Timeout(s, m string) (int64, bool) {
	// Favor the timeout defined for a fully-qualified Method name.
	timeout, ok := c.timeouts[s+"."+m]
	if ok {
		return ToMillis(timeout), ok
	}

	// Fallback on the timeout defined for an entire Service
	timeout, ok = c.timeouts[s]
	if ok {
		return ToMillis(timeout), ok
	}
	return 0, false
}

// ToMillis returns the given Duration as milliseconds.
func ToMillis(d *duration.Duration) int64 {
	if err := d.CheckValid(); err != nil {
		log.Panic("Error converting durationpb to Duration ", err)
	}
	return d.AsDuration().Milliseconds()
}

// RequestLimit returns the request limit in bytes and a presence flag for the
// given fully-qualified Service and simple Method names. A config assignment
// for a specific Method takes precendence over a Service-level assignment.
func (c Config) RequestLimit(s, m string) (int, bool) {
	// Favor the limit defined for a fully-qualified Method name.
	lim, ok := c.reqLimits[s+"."+m]
	if ok {
		return lim, ok
	}

	// Fallback on the limit defined for an entire Service
	lim, ok = c.reqLimits[s]
	return lim, ok
}

// ResponseLimit returns the response limit in bytes and a presence flag for the
// given fully-qualified Service and simple Method names. A config assignment
// for a specific Method takes precendence over a Service-level assignment.
func (c Config) ResponseLimit(s, m string) (int, bool) {
	// Favor the limit defined for a fully-qualified Method name.
	lim, ok := c.resLimits[s+"."+m]
	if ok {
		return lim, ok
	}

	// Fallback on the limit defined for an entire Service
	lim, ok = c.resLimits[s]
	return lim, ok
}
