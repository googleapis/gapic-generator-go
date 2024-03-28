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
	"sort"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/cloud/extendedops"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// customOp represents a custom operation type for long running operations.
type customOp struct {
	message       *descriptorpb.DescriptorProto
	handles       []*descriptorpb.ServiceDescriptorProto
	pollingParams map[*descriptorpb.ServiceDescriptorProto][]string
}

// isCustomOp determines if the given method should return a custom operation wrapper.
func (g *generator) isCustomOp(m *descriptorpb.MethodDescriptorProto, info *httpInfo) bool {
	return g.opts.diregapic && // Generator in DIREGAPIC mode.
		g.aux.customOp != nil && // API Defines a custom operation.
		m.GetOutputType() == g.customOpProtoName() && // Method returns the custom operation.
		m.GetName() != "Wait" && // Method is not a Wait (uses POST).
		info != nil && // Must have google.api.http.
		info.verb != "get" // Method is not a GET (polling methods).
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
func (g *generator) customOpInit(rspVar, reqVar, opVar string, req *descriptorpb.DescriptorProto, s *descriptorpb.ServiceDescriptorProto) {
	h := handleName(s.GetName(), g.opts.pkgName)
	opName := g.aux.customOp.message.GetName()
	pt := g.pt.Printf

	// Collect all of the fields marked with google.cloud.operation_request_field
	// and map the getter method to the polling request's field name, while also
	// collecting these field name keys for sorted access.
	keys := []string{}
	paramToGetter := map[string]string{}
	for _, f := range req.GetField() {
		param, ok := operationRequestField(f)
		if !ok {
			continue
		}
		// Only include those operation_request_fields that are also tracked as
		// polling params for the operation service.
		param = lowerFirst(snakeToCamel(param))
		if params := g.aux.customOp.pollingParams[s]; strContains(params, param) {
			keys = append(keys, param)
			paramToGetter[param] = fmt.Sprintf("%s%s", reqVar, fieldGetter(f.GetName()))
		}
	}
	sort.Strings(keys)

	pt("%s := &%s{", opVar, opName)
	pt("  &%s{", h)
	pt("    c: c.operationClient,")
	pt("    proto: %s,", rspVar)
	for _, param := range keys {
		pt("    %s: %s,", param, paramToGetter[param])
	}
	pt("  },")
	pt("}")
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

	// Wait
	p("// Wait blocks until the operation is complete, polling regularly")
	p("// after an intial period of backing off between attempts.")
	p("func (o *%s) Wait(ctx context.Context, opts ...gax.CallOption) error {", opName)
	p("  bo := gax.Backoff{")
	p("    Initial: %s,", defaultPollInitialDelay)
	p("    Max:     %s,", defaultPollMaxDelay)
	p("  }")
	p("  for {")
	p("    if err := o.Poll(ctx, opts...); err != nil {")
	p("      return err")
	p("    }")
	p("    if o.Done() {")
	p("      return nil")
	p("    }")
	p("    if err := gax.Sleep(ctx, bo.Pause()); err != nil {")
	p("      return err")
	p("    }")
	p("  }")
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
	g.imports[pbinfo.ImportSpec{Path: "time"}] = true
	g.imports[pbinfo.ImportSpec{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}] = true
	g.imports[pbinfo.ImportSpec{Path: "github.com/googleapis/gax-go/v2/apierror"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/googleapi"}] = true

	for _, handle := range op.handles {
		pollingParams := op.pollingParams[handle]
		s := pbinfo.ReduceServName(handle.GetName(), opImp.Name)
		n := handleName(handle.GetName(), opImp.Name)

		// Look up polling method and its input.
		poll := operationPollingMethod(handle)
		pollReq := g.descInfo.Type[poll.GetInputType()].(*descriptorpb.DescriptorProto)
		pollNameField := operationResponseField(pollReq, opNameField.GetName())
		// Look up the fields for error code and error message.
		errorCodeField := operationField(op.message, extendedops.OperationResponseMapping_ERROR_CODE)
		if errorCodeField == nil {
			return fmt.Errorf("field %s not found in %T", extendedops.OperationResponseMapping_ERROR_CODE, op)
		}
		errorCode := snakeToCamel(errorCodeField.GetName())
		errorMessageField := operationField(op.message, extendedops.OperationResponseMapping_ERROR_MESSAGE)
		if errorMessageField == nil {
			return fmt.Errorf("field %s not found in %T", extendedops.OperationResponseMapping_ERROR_MESSAGE, op)
		}
		errorMessage := snakeToCamel(errorMessageField.GetName())

		// type
		p("// Implements the %s interface for %s.", handleInt, handle.GetName())
		p("type %s struct {", n)
		p("  c *%sClient", s)
		p("  proto %s", ptyp)
		for _, param := range pollingParams {
			p("  %s string", param)
		}
		p("}")
		p("")

		// Poll
		p("// Poll retrieves the latest data for the long-running operation.")
		p("func (h *%s) Poll(ctx context.Context, opts ...gax.CallOption) error {", n)
		p("  resp, err := h.c.Get(ctx, &%s.%s{", opImp.Name, upperFirst(pollReq.GetName()))
		p("    %s: h.proto%s,", snakeToCamel(pollNameField.GetName()), opNameGetter)
		for _, f := range pollingParams {
			p("    %s: h.%s,", upperFirst(f), f)
		}
		p("  }, opts...)")
		p("  if err != nil {")
		p("    return err")
		p("  }")
		p("  h.proto = resp")
		p("  if resp.%[1]s != nil && (resp.Get%[1]s() < 200 || resp.Get%[1]s() > 299) {", errorCode)
		p("  	aErr := &googleapi.Error{")
		p("  		Code: int(resp.Get%s()),", errorCode)
		if hasField(op.message, "error") {
			g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
			p(`  		Message: fmt.Sprintf("%%s: %%v", resp.Get%s(), resp.GetError()),`, errorMessage)
		} else {
			p("  		Message: resp.Get%s(),", errorMessage)
		}
		p("  	}")
		p("  	err, _ := apierror.FromError(aErr)")
		p("  	return err")
		p("  }")
		p("  return nil")
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
// to the service that owns the method(s) declaring it, as well as collects the set of
// operation services for handle generation, and maps the polling request parameters to
// that same operation service descriptor.
func (g *generator) loadCustomOpServices(servs []*descriptorpb.ServiceDescriptorProto) {
	handles := g.aux.customOp.handles
	pollingParams := map[*descriptorpb.ServiceDescriptorProto][]string{}
	for _, serv := range servs {
		for _, meth := range serv.GetMethod() {
			if opServ := g.customOpService(meth); opServ != nil {
				g.customOpServices[serv] = opServ
				if !containsService(handles, opServ) {
					handles = append(handles, opServ)
					pollingParams[opServ] = g.pollingRequestParameters(meth, opServ)
				}
				break
			}
		}
	}
	g.aux.customOp.handles = handles
	g.aux.customOp.pollingParams = pollingParams
}

// customOpService loads the ServiceDescriptorProto for the google.cloud.operation_service
// named on the given method.
func (g *generator) customOpService(m *descriptorpb.MethodDescriptorProto) *descriptorpb.ServiceDescriptorProto {
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
func (g *generator) customOpStatusCheck(st *descriptorpb.FieldDescriptorProto) string {
	ret := fmt.Sprintf("return o.Proto()%s", fieldGetter(st.GetName()))
	if st.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
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
func operationField(m *descriptorpb.DescriptorProto, target extendedops.OperationResponseMapping) *descriptorpb.FieldDescriptorProto {
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
func operationResponseField(m *descriptorpb.DescriptorProto, target string) *descriptorpb.FieldDescriptorProto {
	for _, f := range m.GetField() {
		mapping := proto.GetExtension(f.GetOptions(), extendedops.E_OperationResponseField).(string)
		if mapping == target {
			return f
		}
	}
	return nil
}

// operationRequestField is a helper for extracting the operation_request_field annotation from a field.
func operationRequestField(f *descriptorpb.FieldDescriptorProto) (string, bool) {
	mapping := proto.GetExtension(f.GetOptions(), extendedops.E_OperationRequestField).(string)
	if mapping != "" {
		return mapping, true
	}

	return "", false
}

// operationPollingMethod is a helper for finding the operation service RPC annotated with operation_polling_method.
func operationPollingMethod(s *descriptorpb.ServiceDescriptorProto) *descriptorpb.MethodDescriptorProto {
	for _, m := range s.GetMethod() {
		if proto.GetExtension(m.GetOptions(), extendedops.E_OperationPollingMethod).(bool) {
			return m
		}
	}

	return nil
}

// pollingRequestParamters collects the polling request parameters for an operation service's polling method
// based on which are annotated in an initiating RPC's request message with operation_request_field and that
// are also marked as required on the polling request message. Specifically, this weeds out the parent_id field
// of the GlobalOrganizationOperations polling params.
func (g *generator) pollingRequestParameters(m *descriptorpb.MethodDescriptorProto, opServ *descriptorpb.ServiceDescriptorProto) []string {
	var params []string
	poll := operationPollingMethod(opServ)
	if poll == nil {
		return params
	}
	pollReqName := poll.GetInputType()

	inType := g.descInfo.Type[m.GetInputType()].(*descriptorpb.DescriptorProto)

	for _, f := range inType.GetField() {

		mapping, ok := operationRequestField(f)
		if !ok {
			continue
		}
		pollField := g.lookupField(pollReqName, mapping)
		if pollField != nil && isRequired(pollField) {
			params = append(params, lowerFirst(snakeToCamel(mapping)))
		}
	}
	sort.Strings(params)

	return params
}

// handleName is a helper for constructing a operation handle name from the
// operation service name and Go package name.
func handleName(s, pkg string) string {
	s = pbinfo.ReduceServName(s, pkg)
	return lowerFirst(s + "Handle")
}
