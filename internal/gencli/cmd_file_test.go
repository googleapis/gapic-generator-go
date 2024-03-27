// Copyright 2018 Google LLC
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

package gencli

import (
	"path/filepath"
	"testing"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestCommandFile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Generating the command_file panicked: %v", r)
		}
	}()

	// unary
	createTodoCmd := &Command{
		Service:           "Todo",
		Method:            "CreateTodo",
		MethodCmd:         "create-todo",
		ShortDesc:         "creates a todo entry",
		LongDesc:          "creates a todo entry for the user to track",
		InputMessage:      "todopb.Todo",
		InputMessageVar:   "CreateTodoInput",
		OutputMessageType: "todopb.Todo",
		Imports: map[string]*pbinfo.ImportSpec{
			"todopb": &pbinfo.ImportSpec{Name: "todopb", Path: "github.com/googleapis/todo/generated"},
		},
		Flags: []*Flag{
			&Flag{
				Name:      "task",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Required:  true,
				FieldName: "Task",
				VarName:   "CreateTodoInput",
				Usage:     "task to complete",
			},
			&Flag{
				Name:      "done",
				FieldName: "Done",
				VarName:   "CreateTodoInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_BOOL,
				Usage:     "task completion status",
				Optional:  true,
			},
			&Flag{
				Name:          "priority",
				FieldName:     "Priority",
				Type:          descriptorpb.FieldDescriptorProto_TYPE_ENUM,
				Usage:         "importance of the task",
				Message:       "Priority",
				MessageImport: pbinfo.ImportSpec{Name: "todopb"},
				VarName:       "CreateTodoInputPriority",
				Optional:      true,
			},
		},
		HasEnums:    true,
		HasOptional: true,
	}

	// LRO
	startTodoCmd := &Command{
		Service:           "Todo",
		Method:            "StartTodo",
		MethodCmd:         "start-todo",
		ShortDesc:         "starts a todo",
		LongDesc:          "starts a todo that has not been completed yet",
		InputMessage:      "todopb.StartTodoRequest",
		InputMessageVar:   "StartTodoInput",
		OutputMessageType: ".google.longrunning.Operation",
		Imports: map[string]*pbinfo.ImportSpec{
			"todopb": &pbinfo.ImportSpec{Name: "todopb", Path: "github.com/googleapis/todo/generated"},
		},
		Flags: []*Flag{
			&Flag{
				Name:      "id",
				FieldName: "Id",
				VarName:   "StartTodoInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_INT32,
				Required:  true,
				Usage:     "task to start",
			},
		},
		IsLRO: true,
	}

	// client streaming
	copyTodosCmd := &Command{
		Service:           "Todo",
		Method:            "CopyTodos",
		MethodCmd:         "copy-todos",
		ShortDesc:         "stream several todos to create",
		InputMessage:      "todopb.Todo",
		InputMessageVar:   "CopyTodosInput",
		OutputMessageType: ".google.protobuf.Empty",
		Imports: map[string]*pbinfo.ImportSpec{
			"todopb": &pbinfo.ImportSpec{Name: "todopb", Path: "github.com/googleapis/todo/generated"},
		},
		ClientStreaming: true,
	}

	// server streaming
	watchTodosCmd := &Command{
		Service:           "Todo",
		Method:            "WatchTodo",
		MethodCmd:         "watch-todo",
		ShortDesc:         "watch todo",
		LongDesc:          "watch todo for changes, like completion",
		InputMessage:      "todopb.WatchTodoRequest",
		InputMessageVar:   "WatchTodoInput",
		OutputMessageType: "todopb.Todo",
		Imports: map[string]*pbinfo.ImportSpec{
			"todopb": &pbinfo.ImportSpec{Name: "todopb", Path: "github.com/googleapis/todo/generated"},
		},
		Flags: []*Flag{
			&Flag{
				Name:      "id",
				FieldName: "Id",
				VarName:   "WatchTodoInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_INT32,
				Required:  true,
				Usage:     "task to watch",
			},
		},
		ServerStreaming: true,
	}

	// bi-directional streaming
	manageTodosCmd := &Command{
		Service:           "Todo",
		Method:            "ManageTodos",
		MethodCmd:         "manage-todos",
		ShortDesc:         "manage todos live",
		LongDesc:          "manage todos live by creating and updating todos as they change",
		InputMessage:      "todopb.Todo",
		InputMessageVar:   "ManageTodosInput",
		OutputMessageType: "todopb.Todo",
		Imports: map[string]*pbinfo.ImportSpec{
			"todopb": &pbinfo.ImportSpec{Name: "todopb", Path: "github.com/googleapis/todo/generated"},
		},
		ServerStreaming: true,
		ClientStreaming: true,
	}

	for _, tst := range []struct {
		g                *gcli
		cmd              *Command
		name, goldenPath string
	}{
		{
			g: &gcli{
				format: true,
			},
			cmd:        createTodoCmd,
			name:       "create-todo",
			goldenPath: filepath.Join("testdata", "create-todo.want"),
		},
		{
			g: &gcli{
				format: true,
			},
			cmd:        startTodoCmd,
			name:       "start-todo",
			goldenPath: filepath.Join("testdata", "start-todo.want"),
		},
		{
			g: &gcli{
				format: true,
			},
			cmd:        copyTodosCmd,
			name:       "copy-todos",
			goldenPath: filepath.Join("testdata", "copy-todos.want"),
		},
		{
			g: &gcli{
				format: true,
			},
			cmd:        watchTodosCmd,
			name:       "watch-todo",
			goldenPath: filepath.Join("testdata", "watch-todo.want"),
		},
		{
			g: &gcli{
				format: true,
			},
			cmd:        manageTodosCmd,
			name:       "manage-todos",
			goldenPath: filepath.Join("testdata", "manage-todos.want"),
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			tst.g.genCommandFile(tst.cmd)
			if tst.g.response.GetError() != "" {
				t.Errorf("Error generating the command file %s: %s", tst.name, tst.g.response.GetError())
				return
			}

			file := tst.g.response.File[0]

			if file.GetName() != tst.name+".go" {
				t.Errorf("(%+v).genCommands() = %s, want %s", tst.g, file.GetName(), tst.name+".go")
			}
			txtdiff.Diff(t, file.GetContent(), tst.goldenPath)
		})
	}
}
