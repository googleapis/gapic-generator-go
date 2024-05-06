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
	"errors"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	"github.com/google/go-cmp/cmp"
	showcase "github.com/googleapis/gapic-showcase/client"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	locationpb "google.golang.org/genproto/googleapis/cloud/location"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	metadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Clients are initialized in main_test.go.
var (
	echo     *showcase.EchoClient
	echoREST *showcase.EchoClient
)

func TestEcho(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.EchoRequest{
		Response: &showcasepb.EchoRequest_Content{
			Content: content,
		},
	}
	for typ, client := range map[string]*showcase.EchoClient{"grpc": echo, "rest": echoREST} {
		resp, err := client.Echo(context.Background(), req)

		if err != nil {
			t.Fatal(err)
		}
		if resp.GetContent() != req.GetContent() {
			t.Errorf("%s Echo() = %q, want %q", typ, resp.GetContent(), content)
		}
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
	for typ, client := range map[string]*showcase.EchoClient{"grpc": echo, "rest": echoREST} {
		resp, err := client.Echo(context.Background(), req)
		if resp != nil || err == nil {
			t.Errorf("%s Echo() = %v, wanted error %d", typ, resp, val)
		}
		if typ == "grpc" {
			status, _ := status.FromError(err)
			if status.Code() != val {
				t.Errorf("%s Echo() errors with %d, want %d", typ, status.Code(), val)
			}
		} else {
			want := 499
			gerr := &googleapi.Error{}
			if !errors.As(err, &gerr) {
				t.Errorf("%s Echo() returned unexpected error type: %v", typ, err)
			} else if gerr.Code != want {
				t.Errorf("%s Echo() errors with %d, want %d", typ, gerr.Code, want)
			}
		}

	}

}

// Test dynamic routing header generation. We cannot guarantee the order that headers are sent, so we check that the header sent contains the correct elements as opposed to checking
// the header itself.
func TestEchoHeader(t *testing.T) {
	defer check(t)
	for _, tst := range []struct {
		req  *showcasepb.EchoRequest
		want []string
	}{
		{
			req:  &showcasepb.EchoRequest{OtherHeader: "projects/123/instances/456"},
			want: []string{"baz=projects%2F123%2Finstances%2F456", "qux=projects%2F123"},
		},
		{
			req:  &showcasepb.EchoRequest{OtherHeader: "instances/456"},
			want: []string{"baz=instances%2F456"},
		},
		{
			req:  &showcasepb.EchoRequest{Header: "potato"},
			want: []string{"header=potato", "routing_id=potato"},
		},
		{
			req:  &showcasepb.EchoRequest{Header: "projects/123/instances/456"},
			want: []string{"header=projects%2F123%2Finstances%2F456", "routing_id=projects%2F123%2Finstances%2F456", "super_id=projects%2F123", "table_name=projects%2F123%2Finstances%2F456", "instance_id=instances%2F456"},
		},
		{
			req: &showcasepb.EchoRequest{
				Header:      "regions/123/zones/456",
				OtherHeader: "projects/123/instances/456",
			},
			want: []string{"baz=projects%2F123%2Finstances%2F456", "qux=projects%2F123", "table_name=regions%2F123%2Fzones%2F456", "header=regions%2F123%2Fzones%2F456", "routing_id=regions%2F123%2Fzones%2F456"},
		},
	} {
		mdForHeaders := metadata.New(map[string]string{})
		echo.Echo(context.Background(), tst.req, gax.WithGRPCOptions(grpc.Header(&mdForHeaders)))
		got := mdForHeaders.Get("x-goog-request-params")
		got = strings.Split(got[0], "&")
		sort.Strings(got)
		sort.Strings(tst.want)
		if !reflect.DeepEqual(got, tst.want) {
			t.Errorf("Echo() header = %q, want %q", got, tst.want)
		}
	}
}

// Test gRPC trailers for required headers.
func TestEchoTrailers(t *testing.T) {
	defer check(t)
	req := &showcasepb.EchoRequest{Header: "potato"}
	mdForTrailers := metadata.New(map[string]string{})
	echo.Echo(context.Background(), req, gax.WithGRPCOptions(grpc.Trailer(&mdForTrailers)))
	apiVersion := mdForTrailers.Get("x-goog-api-version")
	if len(apiVersion) < 1 || apiVersion[0] == "" {
		t.Error("Echo() x-goog-api-version header not present or empty")
	}
}

func TestEchoHeaderREST(t *testing.T) {
	defer check(t)
	for _, tst := range []struct {
		req  *showcasepb.EchoRequest
		want []string
	}{
		{
			req:  &showcasepb.EchoRequest{OtherHeader: "projects/123/instances/456"},
			want: []string{"baz=projects%2F123%2Finstances%2F456", "qux=projects%2F123"},
		},
		{
			req:  &showcasepb.EchoRequest{OtherHeader: "instances/456"},
			want: []string{"baz=instances%2F456"},
		},
		{
			req:  &showcasepb.EchoRequest{Header: "potato"},
			want: []string{"header=potato", "routing_id=potato"},
		},
		{
			req:  &showcasepb.EchoRequest{Header: "projects/123/instances/456"},
			want: []string{"header=projects%2F123%2Finstances%2F456", "routing_id=projects%2F123%2Finstances%2F456", "super_id=projects%2F123", "table_name=projects%2F123%2Finstances%2F456", "instance_id=instances%2F456"},
		},
		{
			req: &showcasepb.EchoRequest{
				Header:      "regions/123/zones/456",
				OtherHeader: "projects/123/instances/456",
			},
			want: []string{"baz=projects%2F123%2Finstances%2F456", "qux=projects%2F123", "table_name=regions%2F123%2Fzones%2F456", "header=regions%2F123%2Fzones%2F456", "routing_id=regions%2F123%2Fzones%2F456"},
		},
	} {
		// Wrap the default RoundTripper with our own that asserts on the response
		// headers expected by the test.
		wrapped := &http.Client{}
		wrapped.Transport = headerChecker{rt: wrapped.Transport, want: tst.want, t: t}
		echoWrapped, err := showcase.NewEchoRESTClient(context.Background(),
			option.WithEndpoint("http://localhost:7469"),
			option.WithoutAuthentication())
		if err != nil {
			t.Fatal(err)
		}
		echoWrapped.Echo(context.Background(), tst.req)
		echoWrapped.Close()
	}
}

func TestXGoogeHeaders(t *testing.T) {
	// Inspect the private property `xGoogHeaders` of the transport-specific
	// client implementation that is populated on creation of the client.
	w := reflect.ValueOf(*echo)
	x := w.FieldByName("internalClient")
	y := x.Elem().Elem()
	info := y.FieldByName("xGoogHeaders")

	var goVersion string
	var apiVersion string
	vals := make([]string, 0)
	for i := 0; i < info.Len(); i++ {
		key := info.Index(i)
		// Only check for the client info set by the generated setGoogleClientInfo()
		if key.String() == "x-goog-api-client" {
			vals = append(vals, info.Index(i+1).String())
		}
		if key.String() == "x-goog-api-version" {
			apiVersion = info.Index(i + 1).String()
		}

	}

	if apiVersion == "" {
		t.Error("expected x-goog-api-version to be present and non-empty")
	}

	for i := 0; goVersion == "" || i < len(vals); i++ {
		v := vals[i]
		split := strings.Split(v, " ")
		for _, s := range split {
			// For now, we only want to check that the Go version is being
			// properly populated.
			if strings.HasPrefix(s, "gl-go/") {
				goVersion = s
				break
			}
		}
	}

	if goVersion == "" {
		t.Errorf("expected Go version pair to be populated, but wasn't: %v", info)
	} else if strings.Contains(goVersion, "UNKNOWN") {
		t.Errorf("expected Go version pair to not be UNKNOWN: %q", goVersion)
	}
}

type headerChecker struct {
	rt   http.RoundTripper
	want []string
	t    *testing.T
}

func (hc headerChecker) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := hc.rt.RoundTrip(r)

	header := resp.Header
	got := header["x-goog-request-params"]
	got = strings.Split(got[0], "&")
	sort.Strings(got)
	sort.Strings(hc.want)
	if diff := cmp.Diff(got, hc.want); diff != "" {
		hc.t.Errorf("got(-),want(+):\n%s", diff)
	}

	return resp, err
}

// Chat, Collect, and Expand are streaming methods and don't have interesting REST semantics
func TestExpand(t *testing.T) {
	defer check(t)
	content := "The rain in Spain stays mainly on the plain!"
	req := &showcasepb.ExpandRequest{Content: content}

	for typ, client := range map[string]*showcase.EchoClient{"grpc": echo, "rest": echoREST} {
		s, err := client.Expand(context.Background(), req)
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
			t.Errorf("%s Expand() = %q, want %q", typ, got, content)
		}
	}
}

// Chat, Collect, and Expand are streaming methods and don't have interesting REST semantics
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

// Chat, Collect, and Expand are streaming methods and don't have interesting REST semantics
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
			Ttl: &durationpb.Duration{Seconds: 2},
		},
		Response: &showcasepb.WaitRequest_Success{
			Success: &showcasepb.WaitResponse{Content: content},
		},
	}
	for typ, client := range map[string]*showcase.EchoClient{"grpc": echo, "rest": echoREST} {
		op, err := client.Wait(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := op.Wait(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if resp.GetContent() != content {
			t.Errorf("%s Wait() = %q, want %q", typ, resp.GetContent(), content)
		}
	}
}

// TODO(dovs): add REST testing
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
	for typ, client := range map[string]*showcase.EchoClient{"grpc": echo, "rest": echoREST} {
		iter := client.PagedExpand(context.Background(), req)

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
				t.Errorf("%s Chat() = %s, want %s", typ, resp.GetContent(), expected[ndx])
			}
			ndx++
		}
	}
}

func TestPaginationWithToken(t *testing.T) {
	defer check(t)
	str := "ab cd ef gh ij kl"
	expected := strings.Split(str, " ")[1:]
	req := &showcasepb.PagedExpandRequest{Content: str, PageSize: 2, PageToken: "1"}
	for typ, client := range map[string]*showcase.EchoClient{"grpc": echo, "rest": echoREST} {
		iter := client.PagedExpand(context.Background(), req)

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
				t.Errorf("%s Chat() eceived more items than expected", typ)
			} else if resp.GetContent() != expected[ndx] {
				t.Errorf("%s Chat() = %s, want %s", typ, resp.GetContent(), expected[ndx])
			}
			ndx++
		}
	}
}

// TODO(dovs): the server side is not finished for PagedExpandLegacy.
// Add tests for that when it's been merged.

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

func TestBlock_default_timeout(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.BlockRequest{
		ResponseDelay: &durationpb.Duration{Seconds: 6},
		Response: &showcasepb.BlockRequest_Success{
			Success: &showcasepb.BlockResponse{Content: content},
		},
	}

	want := status.New(codes.DeadlineExceeded, "context deadline exceeded")
	resp, err := echo.Block(context.Background(), req)
	if err == nil {
		t.Errorf("Block() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want.Code() {
		t.Errorf("Block() got %+v, want %+v", err, want)
	}
}

func TestBlock_override_default_timeout(t *testing.T) {
	defer check(t)
	content := "hello world!"
	req := &showcasepb.BlockRequest{
		ResponseDelay: &durationpb.Duration{Seconds: 6},
		Response: &showcasepb.BlockRequest_Success{
			Success: &showcasepb.BlockResponse{Content: content},
		},
	}

	resp, err := echo.Block(context.Background(), req, gax.WithTimeout(10*time.Second))
	if err != nil {
		t.Error(err)
	}
	if resp.GetContent() != content {
		t.Errorf("Block() = %q, want %q", resp.GetContent(), content)
	}
}

func TestGetLocation(t *testing.T) {
	defer check(t)
	want := &locationpb.Location{
		Name:        "projects/showcase/location/us-central1",
		DisplayName: "us-central1",
	}
	req := &locationpb.GetLocationRequest{
		Name: want.GetName(),
	}

	got, err := echo.GetLocation(context.Background(), req)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("GetLocation() got(-),want(+):\n%s", diff)
	}
}

func TestListLocations(t *testing.T) {
	defer check(t)
	req := &locationpb.ListLocationsRequest{
		Name: "projects/showcase",
	}
	want := []*locationpb.Location{
		{
			Name:        req.GetName() + "/locations/us-north",
			DisplayName: "us-north",
		},
		{
			Name:        req.GetName() + "/locations/us-south",
			DisplayName: "us-south",
		},
		{
			Name:        req.GetName() + "/locations/us-east",
			DisplayName: "us-east",
		},
		{
			Name:        req.GetName() + "/locations/us-west",
			DisplayName: "us-west",
		},
	}

	iter := echo.ListLocations(context.Background(), req)

	got := []*locationpb.Location{}
	for loc, err := iter.Next(); err == nil; loc, err = iter.Next() {
		got = append(got, loc)
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("ListLocations got(-),want(+):\n%s", diff)
	}
}

func TestIamPolicy(t *testing.T) {
	defer check(t)
	want := &iampb.Policy{
		Bindings: []*iampb.Binding{
			{
				Role:    "foo.editor",
				Members: []string{"allUsers"},
			},
		},
	}
	set := &iampb.SetIamPolicyRequest{
		Resource: "projects/showcase/location/us-central1",
		Policy:   want,
	}

	got, err := echo.SetIamPolicy(context.Background(), set)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("TestIamPolicy() got(-),want(+):\n%s", diff)
	}

	get := &iampb.GetIamPolicyRequest{
		Resource: set.GetResource(),
	}

	got, err = echo.GetIamPolicy(context.Background(), get)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("TestIamPolicy() got(-),want(+):\n%s", diff)
	}

	test := &iampb.TestIamPermissionsRequest{
		Resource:    set.GetResource(),
		Permissions: []string{"foo.create"},
	}
	_, err = echo.TestIamPermissions(context.Background(), test)
	if err != nil {
		t.Error(err)
	}
}

func TestGetIamPolicy_doesNotExist(t *testing.T) {
	defer check(t)
	want := codes.NotFound
	req := &iampb.GetIamPolicyRequest{
		Resource: "projects/foo/location/bar",
	}

	resp, err := echo.GetIamPolicy(context.Background(), req)
	if err == nil {
		t.Errorf("GetIamPolicy() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want {
		t.Errorf("GetIamPolicy() got %+v, want %+v", err, want)
	}
}

func TestGetIamPolicy_missingResource(t *testing.T) {
	defer check(t)
	want := codes.InvalidArgument

	resp, err := echo.GetIamPolicy(context.Background(), &iampb.GetIamPolicyRequest{})
	if err == nil {
		t.Errorf("GetIamPolicy() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want {
		t.Errorf("GetIamPolicy() got %+v, want %+v", err, want)
	}
}

func TestSetIamPolicy_missingResource(t *testing.T) {
	defer check(t)
	want := codes.InvalidArgument

	resp, err := echo.SetIamPolicy(context.Background(), &iampb.SetIamPolicyRequest{})
	if err == nil {
		t.Errorf("SetIamPolicy() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want {
		t.Errorf("SetIamPolicy() got %+v, want %+v", err, want)
	}
}

func TestSetIamPolicy_missingPolicy(t *testing.T) {
	defer check(t)
	want := codes.InvalidArgument
	req := &iampb.SetIamPolicyRequest{
		Resource: "projects/showcase/location/us-central1",
	}

	resp, err := echo.SetIamPolicy(context.Background(), req)
	if err == nil {
		t.Errorf("SetIamPolicy() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want {
		t.Errorf("SetIamPolicy() got %+v, want %+v", err, want)
	}
}

func TestTestIamPermissions_doesNotExist(t *testing.T) {
	defer check(t)
	want := codes.NotFound
	req := &iampb.TestIamPermissionsRequest{
		Resource: "projects/foo/location/bar",
	}

	resp, err := echo.TestIamPermissions(context.Background(), req)
	if err == nil {
		t.Errorf("TestIamPermissions() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want {
		t.Errorf("TestIamPermissions() got %+v, want %+v", err, want)
	}
}

func TestTestIamPermissions_missingResource(t *testing.T) {
	defer check(t)
	want := codes.InvalidArgument

	resp, err := echo.TestIamPermissions(context.Background(), &iampb.TestIamPermissionsRequest{})
	if err == nil {
		t.Errorf("TestIamPermissions() got %+v, want %+v", resp, want)
	} else if got, ok := status.FromError(err); !ok || got.Code() != want {
		t.Errorf("TestIamPermissions() got %+v, want %+v", err, want)
	}
}
