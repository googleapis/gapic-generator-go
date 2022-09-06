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
	"testing"

	"github.com/google/go-cmp/cmp"
	showcase "github.com/googleapis/gapic-showcase/client"
	showcasepb "github.com/googleapis/gapic-showcase/server/genproto"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Clients are initialized in main_test.go.
var (
	identity     *showcase.IdentityClient
	identityREST *showcase.IdentityClient
)

func TestUserCRUD(t *testing.T) {
	defer check(t)
	ctx := context.Background()

	for _, client := range map[string]*showcase.IdentityClient{"grpc": identity, "rest": identityREST} {
		create := &showcasepb.CreateUserRequest{
			User: &showcasepb.User{
				DisplayName: "Jane Doe",
				Email:       "janedoe@example.com",
				Nickname:    proto.String("Doe"),
				HeightFeet:  proto.Float64(6.2),
			},
		}

		usr, err := client.CreateUser(ctx, create)
		if err != nil {
			t.Fatal(err)
		}

		want := create.GetUser()
		if usr.GetName() == "" {
			t.Errorf("CreateUser().Name was unexpectedly empty")
		}
		if usr.GetDisplayName() != want.GetDisplayName() {
			t.Errorf("CreateUser().DisplayName = %q, want = %q", usr.GetDisplayName(), want.GetDisplayName())
		}
		if usr.GetEmail() != want.GetEmail() {
			t.Errorf("CreateUser().Email = %q, want = %q", usr.GetEmail(), want.GetEmail())
		}
		if usr.GetCreateTime() == nil {
			t.Errorf("CreateUser().CreateTime was unexpectedly empty")
		}
		if usr.GetUpdateTime() == nil {
			t.Errorf("CreateUser().UpdateTime was unexpectedly empty")
		}
		if usr.GetNickname() != want.GetNickname() {
			t.Errorf("CreateUser().Nickname = %q, want = %q", usr.GetNickname(), want.GetNickname())
		}
		if usr.GetHeightFeet() != want.GetHeightFeet() {
			t.Errorf("CreateUser().HeightFeet = %f, want = %f", usr.GetHeightFeet(), want.GetHeightFeet())
		}
		if usr.Age != nil {
			t.Errorf("CreateUser().Age was unexpectedly set to: %d", usr.GetAge())
		}
		if usr.EnableNotifications != nil {
			t.Errorf("CreateUser().EnableNotifications was unexpectedly set to: %v", usr.GetEnableNotifications())
		}

		list := &showcasepb.ListUsersRequest{
			PageSize: 5,
		}

		iter := client.ListUsers(context.Background(), list)

		if max := iter.PageInfo().MaxSize; max != int(list.PageSize) {
			t.Errorf("PageInfo().MaxSize = %d, want %d", max, list.PageSize)
		}

		listed, err := iter.Next()
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(listed, usr, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("ListUsers() got=-, want=+:%s", diff)
		}

		get := &showcasepb.GetUserRequest{
			Name: usr.GetName(),
		}

		got, err := client.GetUser(ctx, get)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(got, usr, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("GetUser() got=-, want=+:%s", diff)
		}

		update := &showcasepb.UpdateUserRequest{
			User: &showcasepb.User{
				Name:                got.GetName(),
				DisplayName:         got.GetDisplayName(),
				Email:               "janedoe@jane.com",
				HeightFeet:          proto.Float64(6.0),
				EnableNotifications: proto.Bool(true),
			},
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email", "height_feet", "enable_notifications"}},
		}

		updated, err := client.UpdateUser(ctx, update)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(updated, usr, cmp.Comparer(proto.Equal)); diff == "" {
			t.Errorf("UpdateUser() users were the same, update failed")
		}
		if updated.GetEmail() == usr.GetEmail() {
			t.Errorf("UpdateUser().Email was not updated as expected")
		}
		if updated.GetNickname() != usr.GetNickname() {
			t.Errorf("UpdateUser().Nickname = %q, want = %q", updated.GetNickname(), usr.GetNickname())
		}
		if updated.GetHeightFeet() == usr.GetHeightFeet() {
			t.Errorf("UpdateUser().HeightFeet was not updated as expected")
		}
		if updated.EnableNotifications == nil || !updated.GetEnableNotifications() {
			t.Errorf("UpdateUser().EnableNotifications was not updated as expected")
		}
		if updated.Age != nil {
			t.Errorf("UpdateUser().Age was unexpectedly updated")
		}

		err = client.DeleteUser(ctx, &showcasepb.DeleteUserRequest{
			Name: usr.GetName(),
		})
		if err != nil {
			t.Fatal(err)
		}

		iter = client.ListUsers(ctx, &showcasepb.ListUsersRequest{})

		_, err = iter.Next()
		if err != iterator.Done {
			t.Errorf("ListUsers() = %q, want %q", err, iterator.Done)
		}
	}
}
