// Copyright 2021 Google LLC
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

package gengapic

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/cloud/extendedops"
	"google.golang.org/protobuf/proto"
)

// customOp represents a custom operation type for long running operations.
type customOp struct {
	message *descriptor.DescriptorProto
	handles []*descriptor.ServiceDescriptorProto
}

// isCustomOp determines if the given method should return a custom operation wrapper.
func (g *generator) isCustomOp(m *descriptor.MethodDescriptorProto, info *httpInfo) bool {
	return g.opts.diregapic && // Generator in DIREGAPIC mode.
		g.aux.customOp != nil && // API Defines a custom operation.
		m.GetOutputType() == g.customOpProtoName() && // Method returns the custom operation.
		info.verb != "get" && // Method is not a GET (polling methods).
		m.GetName() != "Wait" // Method is not a Wait (uses POST).
}

// customOpProtoName builds the fully-qualified proto name for the custom
// operation message type.
func (g *generator) customOpProtoName() string {
	f := g.descInfo.ParentFile[g.aux.customOp.message]
	return fmt.Sprintf(".%s.%s", f.GetPackage(), g.aux.customOp.message.GetName())
}

// customOpPointerType builds a string containing the Go code for a pointer to
// the custom operation type.
func (g *generator) customOpPointerType() (string, error) {
	op := g.aux.customOp
	if op == nil {
		return "", nil
	}

	opName, imp, err := g.descInfo.NameSpec(op.message)
	if err != nil {
		return "", err
	}

	s := fmt.Sprintf("*%s.%s", imp.Name, opName)

	return s, nil
}

// customOpInit builds a string containing the Go code for initializing the
// operation wrapper type with the Go identifier for a variable that is the
// proto-defined operation type.
func (g *generator) customOpInit(h, p string) string {
	opName := g.aux.customOp.message.GetName()

	s := fmt.Sprintf("&%s{&%s{c: c.operationClient, proto: %s}}", opName, h, p)

	return s
}

// customOperationType generates the custom operation wrapper type and operation
// service handle implementations using the generators current printer. This
// should only be called once per package.
func (g *generator) customOperationType() error {
	op := g.aux.customOp
	if op == nil {
		return nil
	}
	opName := op.message.GetName()
	handleInt := lowerFirst(opName + "Handle")

	ptyp, err := g.customOpPointerType()
	if err != nil {
		return err
	}
	_, opImp, err := g.descInfo.NameSpec(op.message)
	if err != nil {
		return err
	}
	g.imports[opImp] = true

	statusField := operationField(op.message, extendedops.OperationResponseMapping_STATUS)
	if statusField == nil {
		return fmt.Errorf("operation message %s is missing an annotated status field", op.message.GetName())
	}

	opNameField := operationField(op.message, extendedops.OperationResponseMapping_NAME)
	if opNameField == nil {
		return fmt.Errorf("operation message %s is missing an annotated name field", op.message.GetName())
	}

	opNameGetter := fieldGetter(opNameField.GetName())

	p := g.printf

	p("// %s represents a long-running operation for this API.", opName)
	p("type %s struct {", opName)
	p("  %s", handleInt)
	p("}")
	p("")

	// Done
	p("// Done reports whether the long-running operation has completed.")
	p("func (o *%s) Done() bool {", opName)
	p(g.customOpStatusCheck(statusField))
	p("}")
	p("")

	// Name
	p("// Name returns the name of the long-running operation.")
	p("// The name is assigned by the server and is unique within the service from which the operation is created.")
	p("func (o *%s) Name() string {", opName)
	p("  return o.Proto()%s", opNameGetter)
	p("}")
	p("")

	p("type %s interface {", handleInt)
	p("  // Poll retrieves the operation.")
	p("  Poll(ctx context.Context, opts ...gax.CallOption) error")
	p("")
	p("  // Proto returns the long-running operation message.")
	p("  Proto() %s", ptyp)
	p("}")
	p("")
	g.imports[pbinfo.ImportSpec{Path: "context"}] = true
	g.imports[pbinfo.ImportSpec{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}] = true

	for _, handle := range op.handles {
		s := pbinfo.ReduceServName(handle.GetName(), opImp.Name)
		n := lowerFirst(s + "Handle")

		// Look up polling method and its input.
		var get *descriptor.MethodDescriptorProto
		for _, m := range handle.GetMethod() {
			if m.GetName() == "Get" {
				get = m
				break
			}
		}
		getInput := g.descInfo.Type[get.GetInputType()]
		inNameField := operationResponseField(getInput.(*descriptor.DescriptorProto), opNameField.GetName())

		// type
		p("// Implements the %s interface for %s.", handleInt, handle.GetName())
		p("type %s struct {", n)
		p("  c *%sClient", s)
		p("  proto %s", ptyp)
		p("}")
		p("")

		// Poll
		p("// Poll retrieves the latest data for the long-running operation.")
		p("func (h *%s) Poll(ctx context.Context, opts ...gax.CallOption) error {", n)
		p("  resp, err := h.c.Get(ctx, &%s.%s{%s: h.proto%s}, opts...)", opImp.Name, upperFirst(getInput.GetName()), upperFirst(inNameField.GetName()), opNameGetter)
		p("  if err != nil {")
		p("    return err")
		p("  }")
		p("  h.proto = resp")
		p("  return nil")
		p("}")
		p("")

		// Wait
		p("// Wait blocks until the operation is complete, polling regularly")
		p("// after an intial period of backing off between attempts.")
		p("func (h *%s) Wait(ctx context.Context, opts ...gax.CallOption) error {", n)
		p("	 bo := gax.Backoff{")
		p("    Initial: 1 * time.Second,")
		p("    Max:     time.Minute,")
		p("	 }")
		p("	 for {")
		p("    if err := h.Poll(ctx, opts...); err != nil {")
		p("      return err")
		p("    }")
		p("    if h.Proto().Done() {")
		p("      return nil")
		p("    }")
		p("    if err := gax.Sleep(ctx, bo.Pause()); err != nil {")
		p("      return err")
		p("    }")
		p("  }")
		p("}")
		p("")

		// Proto
		p("// Proto returns the raw type this wraps.")
		p("func (h *%s) Proto() %s {", n, ptyp)
		p("  return h.proto")
		p("}")
		p("")
	}

	return nil
}

// loadCustomOpServices maps the service declared as a google.cloud.operation_service
// to the service that owns the method(s) declaring it.
func (g *generator) loadCustomOpServices(servs []*descriptor.ServiceDescriptorProto) {
	handles := g.aux.customOp.handles
	for _, serv := range servs {
		for _, meth := range serv.GetMethod() {
			if opServ := g.customOpService(meth); opServ != nil {
				g.customOpServices[serv] = opServ
				if !containsService(handles, opServ) {
					handles = append(handles, opServ)
				}
				break
			}
		}
	}
	g.aux.customOp.handles = handles
}

// customOpService loads the ServiceDescriptorProto for the google.cloud.operation_service
// named on the given method.
func (g *generator) customOpService(m *descriptor.MethodDescriptorProto) *descriptor.ServiceDescriptorProto {
	opServName := proto.GetExtension(m.GetOptions(), extendedops.E_OperationService).(string)
	if opServName == "" {
		return nil
	}

	file := g.descInfo.ParentFile[m]
	fqn := fmt.Sprintf(".%s.%s", file.GetPackage(), opServName)

	return g.descInfo.Serv[fqn]
}

// customOpStatusCheck constructs a return statement that checks if the operation's Status
// field indicates it is done.
func (g *generator) customOpStatusCheck(st *descriptor.FieldDescriptorProto) string {
	ret := fmt.Sprintf("return o.Proto()%s", fieldGetter(st.GetName()))
	if st.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM {
		done := g.customOpStatusEnumDone()
		ret = fmt.Sprintf("%s == %s", ret, done)
	}

	return ret
}

// customOpStatusEnumDone constructs the Go name of the operation's status enum
// DONE value.
func (g *generator) customOpStatusEnumDone() string {
	op := g.aux.customOp.message

	// Ignore the error here, it would failed much earlier if the
	// operation type was unresolvable.
	_, imp, _ := g.descInfo.NameSpec(op)

	// Ignore the nil case here, it would failed earlier if the
	// status field was not present.
	statusField := operationField(op, extendedops.OperationResponseMapping_STATUS)
	statusEnum := g.descInfo.Type[statusField.GetTypeName()]

	enum := fmt.Sprintf("%s_DONE", g.nestedName(g.descInfo.ParentElement[statusEnum]))

	s := fmt.Sprintf("%s.%s", imp.Name, enum)

	return s
}

// operationField is a helper for loading the target google.cloud.operation_field annotation value
// if present on a field in the given message.
func operationField(m *descriptor.DescriptorProto, target extendedops.OperationResponseMapping) *descriptor.FieldDescriptorProto {
	for _, f := range m.GetField() {
		mapping := proto.GetExtension(f.GetOptions(), extendedops.E_OperationField).(extendedops.OperationResponseMapping)
		if mapping == target {
			return f
		}
	}
	return nil
}

// operationResponseField is a helper for finding the message field that declares the target field name
// in the google.cloud.operation_response_field annotation.
func operationResponseField(m *descriptor.DescriptorProto, target string) *descriptor.FieldDescriptorProto {
	for _, f := range m.GetField() {
		mapping := proto.GetExtension(f.GetOptions(), extendedops.E_OperationResponseField).(string)
		if mapping == target {
			return f
		}
	}
	return nil
}

// handleName is a helper for constructing a operation handle name from the
// operation service name and Go package name.
func handleName(s, pkg string) string {
	s = pbinfo.ReduceServName(s, pkg)
	return lowerFirst(s + "Handle")
}
