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

package gengapic

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/duration"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func (g *generator) clientOptions(serv *descriptor.ServiceDescriptorProto, servName string) error {
	p := g.printf
	sFQN := g.descInfo.ParentFile[serv].GetPackage() + "." + serv.GetName()

	// CallOptions struct
	{
		p("// %[1]sCallOptions contains the retry settings for each method of %[1]sClient.", servName)
		p("type %sCallOptions struct {", servName)
		for _, m := range serv.Method {
			p("%s []gax.CallOption", *m.Name)
		}
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{"gax", "github.com/googleapis/gax-go/v2"}] = true
	}

	// defaultClientOptions
	{
		var host string
		if eHost, err := proto.GetExtension(serv.Options, annotations.E_DefaultHost); err == nil {
			host = *eHost.(*string)
		} else if g.serviceConfig != nil {
			// TODO(ndietz) remove this once default_host annotation is acepted
			host = g.serviceConfig.Name
		} else {
			fqn := g.descInfo.ParentFile[serv].GetPackage() + "." + serv.GetName()
			return fmt.Errorf("service %q is missing option google.api.default_host", fqn)
		}

		if !strings.Contains(host, ":") {
			host += ":443"
		}

		p("func default%sClientOptions() []option.ClientOption {", servName)
		p("  return []option.ClientOption{")
		p("    option.WithEndpoint(%q),", host)
		p("    option.WithGRPCDialOption(grpc.WithDisableServiceConfig()),")
		p("    option.WithScopes(DefaultAuthScopes()...),")
		p("  }")
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option"}] = true
	}

	// defaultCallOptions
	{
		policies := map[string]*conf.MethodConfig_RetryPolicy{}

		// gather retry policies from gRPC ServiceConfig
		for _, mc := range g.grpcConf.GetMethodConfig() {
			for _, name := range mc.GetName() {
				base := name.GetService()

				// skip the Name entry if it's not the current service
				if base != sFQN {
					continue
				}

				// individual method config, overwrites service-level config
				if name.GetMethod() != "" {
					base = base + "." + name.GetMethod()
					policies[base] = mc.GetRetryPolicy()
					continue
				}

				// service-level config, apply to all *unset* methods
				for _, m := range serv.GetMethod() {
					// build fully-qualified name
					fqn := base + "." + m.GetName()
					if _, ok := policies[fqn]; !ok {
						policies[fqn] = mc.GetRetryPolicy()
					}
				}
			}
		}

		if len(policies) > 0 {
			g.imports[pbinfo.ImportSpec{Path: "time"}] = true
			g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc/codes"}] = true
		}

		// read retry params from gRPC ServiceConfig
		p("func default%[1]sCallOptions() *%[1]sCallOptions {", servName)
		p("  return &%sCallOptions{", servName)
		for _, m := range serv.GetMethod() {
			mFQN := sFQN + "." + m.GetName()
			p("%s: []gax.CallOption{", m.GetName())
			if rp, ok := policies[mFQN]; ok && rp != nil && len(rp.GetRetryableStatusCodes()) > 0 {
				p("gax.WithRetry(func() gax.Retryer {")
				p("  return gax.OnCodes([]codes.Code{")
				for _, c := range rp.GetRetryableStatusCodes() {
					p("    codes.%s,", snakeToCamel(c.String()))
				}
				p("	 }, gax.Backoff{")
				// this ignores max_attempts
				p("		Initial:    %d * time.Millisecond,", durationToMillis(rp.GetInitialBackoff()))
				p("		Max:        %d * time.Millisecond,", durationToMillis(rp.GetMaxBackoff()))
				p("		Multiplier: %.2f,", rp.GetBackoffMultiplier())
				p("	 })")
				p("}),")
			}
			p("},")
		}
		p("  }")
		p("}")
		p("")
	}

	return nil
}

func durationToMillis(d *duration.Duration) int64 {
	return d.GetSeconds()*1000 + int64(d.GetNanos()/1000000)
}

func (g *generator) clientInit(serv *descriptor.ServiceDescriptorProto, servName string) error {
	p := g.printf

	var hasLRO bool
	for _, m := range serv.Method {
		if *m.OutputType == lroType {
			hasLRO = true
			break
		}
	}

	imp, err := g.descInfo.ImportSpec(serv)
	if err != nil {
		return err
	}

	// client struct
	{
		p("// %sClient is a client for interacting with %s.", servName, g.apiName)
		p("//")
		p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
		p("type %sClient struct {", servName)

		p("// The connection to the service.")
		p("conn *grpc.ClientConn")
		p("")

		p("// The gRPC API client.")
		p("%s %s.%sClient", grpcClientField(servName), imp.Name, serv.GetName())
		p("")

		if hasLRO {
			p("// LROClient is used internally to handle longrunning operations.")
			p("// It is exposed so that its CallOptions can be modified if required.")
			p("// Users should not Close this client.")
			p("LROClient *lroauto.OperationsClient")
			p("")

			g.imports[pbinfo.ImportSpec{Name: "lroauto", Path: "cloud.google.com/go/longrunning/autogen"}] = true
		}

		p("// The call options for this service.")
		p("CallOptions *%sCallOptions", servName)
		p("")

		p("// The x-goog-* metadata to be sent with each request.")
		p("xGoogMetadata metadata.MD")
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc/metadata"}] = true
	}

	// Client constructor
	{
		clientName := camelToSnake(serv.GetName())
		clientName = strings.Replace(clientName, "_", " ", -1)

		p("// New%sClient creates a new %s client.", servName, clientName)
		p("//")
		g.comment(g.comments[serv])
		p("func New%[1]sClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
		p("  conn, err := transport.DialGRPC(ctx, append(default%sClientOptions(), opts...)...)", servName)
		p("  if err != nil {")
		p("    return nil, err")
		p("  }")
		p("  c := &%sClient{", servName)
		p("    conn:        conn,")
		p("    CallOptions: default%sCallOptions(),", servName)
		p("")
		p("    %s: %s.New%sClient(conn),", grpcClientField(servName), imp.Name, serv.GetName())
		p("  }")
		p("  c.setGoogleClientInfo()")
		p("")

		if hasLRO {
			p("  c.LROClient, err = lroauto.NewOperationsClient(ctx, option.WithGRPCConn(conn))")
			p("  if err != nil {")
			p("    // This error \"should not happen\", since we are just reusing old connection")
			p("    // and never actually need to dial.")
			p("    // If this does happen, we could leak conn. However, we cannot close conn:")
			p("    // If the user invoked the function with option.WithGRPCConn,")
			p("    // we would close a connection that's still in use.")
			p("    // TODO(pongad): investigate error conditions.")
			p("    return nil, err")
			p("  }")
		}

		p("  return c, nil")
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/transport"}] = true
		g.imports[pbinfo.ImportSpec{Path: "context"}] = true
	}

	// Connection()
	{
		p("// Connection returns the client's connection to the API service.")
		p("func (c *%sClient) Connection() *grpc.ClientConn {", servName)
		p("  return c.conn")
		p("}")
		p("")
	}

	// Close()
	{
		p("// Close closes the connection to the API service. The user should invoke this when")
		p("// the client is no longer required.")
		p("func (c *%sClient) Close() error {", servName)
		p("  return c.conn.Close()")
		p("}")
		p("")
	}

	// setGoogleClientInfo
	{
		p("// setGoogleClientInfo sets the name and version of the application in")
		p("// the `x-goog-api-client` header passed on each request. Intended for")
		p("// use by Google-written clients.")
		p("func (c *%sClient) setGoogleClientInfo(keyval ...string) {", servName)
		p(`  kv := append([]string{"gl-go", versionGo()}, keyval...)`)
		p(`  kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)`)
		p(`  c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))`)
		p("}")
		p("")
	}

	return nil
}
