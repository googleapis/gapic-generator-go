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

func TestMain(m *testing.M) {
	flag.Parse()

	conn, err := grpc.Dial("localhost:7469", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	opt := option.WithGRPCConn(conn)
	ctx := context.Background()

	echo, err = showcase.NewEchoClient(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}

	identity, err = showcase.NewIdentityClient(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
