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
	"flag"
	"log"
	"os"
	"testing"

	showcase "github.com/googleapis/gapic-showcase/client"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func init() {
	// These "leaks" are created by the client connection not being closed at the
	// end of an individual test. This is not an issue for us, because the client
	// connection is shared across tests and closed by the TestMain. We are more
	// concerned with leaking contexts used in tests by GAPICs.
	//
	// TODO(noahdietz): explore refactoring tests to using individual connections.
	registerIgnoreGoroutine("google.golang.org/grpc/internal/transport.newHTTP2Client")
	registerIgnoreGoroutine("google.golang.org/grpc.newCCBalancerWrapper")
	registerIgnoreGoroutine("google.golang.org/grpc.(*addrConn).connect")
}

const showcaseSemver = "0.15.0"

func TestMain(m *testing.M) {
	flag.Parse()

	conn, err := grpc.Dial("localhost:7469", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	opt := option.WithGRPCConn(conn)
	ctx := context.Background()

	echo, err = showcase.NewEchoClient(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}
	defer echo.Close()

	// The custom endpoint bypasses https.
	echoREST, err = showcase.NewEchoRESTClient(ctx, option.WithEndpoint("http://localhost:7469"), option.WithoutAuthentication())
	if err != nil {
		log.Fatal(err)
	}
	defer echoREST.Close()

	identity, err = showcase.NewIdentityClient(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}
	defer identity.Close()

	sequenceClient, err = showcase.NewSequenceClient(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}
	defer sequenceClient.Close()

	// The custom endpoint bypasses https.
	complianceClient, err = showcase.NewComplianceRESTClient(ctx, option.WithEndpoint("http://localhost:7469"), option.WithoutAuthentication())

	if err != nil {
		log.Fatal(err)
	}
	defer complianceClient.Close()

	os.Exit(m.Run())
}
