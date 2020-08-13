// Copyright 2019 Google LLC
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
	"io"
	"strings"
	"testing"
	"time"

	durationpb "github.com/golang/protobuf/ptypes/duration"
	showcase "github.com/googleapis/gapic-showcase/client"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/iterator"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var echo *showcase.EchoClient

func TestEcho(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.EchoRequest{
		Response: &showcasepb.EchoRequest_Content{
			Content: content,
		},
	}
	resp, err := echo.Echo(context.Background(), req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.GetContent() != req.GetContent() {
		t.Errorf("Echo() = %q, want %q", resp.GetContent(), content)
	}
}

func TestEcho_error(t *testing.T) {
	defer check(t)
	val := codes.Canceled
	req := &showcasepb.EchoRequest{
		Response: &showcasepb.EchoRequest_Error{
			Error: &spb.Status{Code: int32(val)},
		},
	}
	resp, err := echo.Echo(context.Background(), req)

	if resp != nil {
		t.Errorf("Echo() = %v, wanted error %d", resp, val)
	}
	status, _ := status.FromError(err)
	if status.Code() != val {
		t.Errorf("Echo() errors with %d, want %d", status.Code(), val)
	}
}

func TestExpand(t *testing.T) {
	defer check(t)
	content := "The rain in Spain stays mainly on the plain!"
	req := &showcasepb.ExpandRequest{Content: content}
	s, err := echo.Expand(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	resps := []string{}
	for {
		resp, err := s.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		resps = append(resps, resp.GetContent())
	}
	got := strings.Join(resps, " ")
	if content != got {
		t.Errorf("Expand() = %q, want %q", got, content)
	}
}

func TestCollect(t *testing.T) {
	defer check(t)
	content := "The rain in Spain stays mainly on the plain!"
	s, err := echo.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, str := range strings.Split(content, " ") {
		s.Send(&showcasepb.EchoRequest{
			Response: &showcasepb.EchoRequest_Content{Content: str}})
	}

	resp, err := s.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}
	if content != resp.GetContent() {
		t.Errorf("Collect() = %q, want %q", resp.GetContent(), content)
	}
}

func TestChat(t *testing.T) {
	defer check(t)
	content := "The rain in Spain stays mainly on the plain!"
	s, err := echo.Chat(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, str := range strings.Split(content, " ") {
		s.Send(&showcasepb.EchoRequest{
			Response: &showcasepb.EchoRequest_Content{Content: str}})
	}
	s.CloseSend()
	resps := []string{}
	for {
		resp, err := s.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		resps = append(resps, resp.GetContent())
	}
	got := strings.Join(resps, " ")
	if content != got {
		t.Errorf("Chat() = %q, want %q", got, content)
	}
}

func TestWait(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.WaitRequest{
		End: &showcasepb.WaitRequest_Ttl{
			Ttl: &durationpb.Duration{Nanos: 100},
		},
		Response: &showcasepb.WaitRequest_Success{
			Success: &showcasepb.WaitResponse{Content: content},
		},
	}
	op, err := echo.Wait(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := op.Wait(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetContent() != content {
		t.Errorf("Wait() = %q, want %q", resp.GetContent(), content)
	}
}

func TestWait_timeout(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.WaitRequest{
		End: &showcasepb.WaitRequest_Ttl{
			Ttl: &durationpb.Duration{Seconds: 1},
		},
		Response: &showcasepb.WaitRequest_Success{
			Success: &showcasepb.WaitResponse{Content: content},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	op, err := echo.Wait(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := op.Wait(ctx)
	if err == nil {
		t.Errorf("Wait() = %+v, want error", resp)
	}
}

func TestPagination(t *testing.T) {
	defer check(t)
	str := "foo bar biz baz"
	expected := strings.Split(str, " ")
	req := &showcasepb.PagedExpandRequest{Content: str, PageSize: 2}
	iter := echo.PagedExpand(context.Background(), req)

	ndx := 0
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if resp.GetContent() != expected[ndx] {
			t.Errorf("Chat() = %s, want %s", resp.GetContent(), expected[ndx])
		}
		ndx++
	}
}

func TestPaginationWithToken(t *testing.T) {
	defer check(t)
	str := "ab cd ef gh ij kl"
	expected := strings.Split(str, " ")[1:]
	req := &showcasepb.PagedExpandRequest{Content: str, PageSize: 2, PageToken: "1"}
	iter := echo.PagedExpand(context.Background(), req)

	ndx := 0
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		if ndx >= len(expected) {
			t.Errorf("Received more items than expected")
		} else if resp.GetContent() != expected[ndx] {
			t.Errorf("Chat() = %s, want %s", resp.GetContent(), expected[ndx])
		}
		ndx++
	}
}

func TestBlock(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.BlockRequest{
		ResponseDelay: &durationpb.Duration{Nanos: 1000},
		Response: &showcasepb.BlockRequest_Success{
			Success: &showcasepb.BlockResponse{Content: content},
		},
	}
	resp, err := echo.Block(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetContent() != content {
		t.Errorf("Block() = %q, want %q", resp.GetContent(), content)
	}
}

func TestBlock_timeout(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.BlockRequest{
		ResponseDelay: &durationpb.Duration{Seconds: 1},
		Response: &showcasepb.BlockRequest_Success{
			Success: &showcasepb.BlockResponse{Content: content},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	want := status.New(codes.DeadlineExceeded, "context deadline exceeded")
	resp, err := echo.Block(ctx, req)
	if err == nil {
		t.Errorf("Block() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want.Code() {
		t.Errorf("Block() got %+v, want %+v", err, want)
	}
}
