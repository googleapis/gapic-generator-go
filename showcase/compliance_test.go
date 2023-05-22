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

package showcase

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	showcase "github.com/googleapis/gapic-showcase/client"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Client is initialized in main_test.go.
var complianceClient *showcase.ComplianceClient

// method is the function type that implements any of the Compliance.RepeatData* RPCs, which all have the same signature.
type method func(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error)

// wipMethod is a stub of type `method` for use in TestComplianceSuite() as a stand-in for the RPC
// types listed in the compliance suite file that are known to not be fully implemented in the Go
// generator yet. This allows the test suite to pass during development, verifying that the
// implemented methods work correctly but not erroring on methods we know are not yet implemented.
//
// TODO: Remove the need for wipMethod by implementing all features as per design docs. Once
// everything is in place, this method and its references should not be needed for all tests to
// pass.
func wipMethod(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error) {
	return &genprotopb.RepeatResponse{Request: req}, nil
}

// TestComplianceSuite is used to ensure the REST transport for the GAPIC client emitted by this
// generator for the Showcase API works correctly. It depends on complianceClient having been
// already initialized to be REST client.
func TestComplianceSuite(t *testing.T) {
	defer check(t)

	// restRPCs is a map of the RPC names as they appear in the compliance suite file to the
	// actual Go methods that implement them.
	restRPCs := map[string]method{
		"Compliance.RepeatDataBody":                 complianceClient.RepeatDataBody,
		"Compliance.RepeatDataBodyInfo":             complianceClient.RepeatDataBodyInfo,
		"Compliance.RepeatDataQuery":                complianceClient.RepeatDataQuery,
		"Compliance.RepeatDataSimplePath":           complianceClient.RepeatDataSimplePath,
		"Compliance.RepeatDataPathResource":         wipMethod, // TODO: replace with complianceClient.RepeatDataPathResource,
		"Compliance.RepeatDataPathTrailingResource": wipMethod, // TODO: replace with complianceClient.RepeatDataPathTrailingResource,
		"Compliance.RepeatDataBodyPut":              complianceClient.RepeatDataBodyPut,
		"Compliance.RepeatDataBodyPatch":            complianceClient.RepeatDataBodyPatch,
	}

	suite, err := getComplianceSuite()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, group := range suite.GetGroup() {
		rpcsToTest := group.GetRpcs()
		for requestIdx, requestProto := range group.GetRequests() {
			for rpcIdx, rpcName := range rpcsToTest {
				errorPrefix := fmt.Sprintf("[request %d/%q: rpc %q/%d/%q]",
					requestIdx, requestProto.GetName(), group.Name, rpcIdx, rpcName)

				// Ensure that we issue only the RPCs the test suite is expecting.
				method, ok := restRPCs[rpcName]
				if !ok {
					t.Errorf("%s could not find client library method for this RPC", errorPrefix)
					continue
				}

				response, err := method(context.Background(), requestProto)
				if err != nil {
					t.Errorf("%s error: %s", errorPrefix, err)
				}
				// Check for expected response.
				if diff := cmp.Diff(response.GetRequest().GetInfo(), requestProto.GetInfo(), cmp.Comparer(proto.Equal)); diff != "" {
					t.Errorf("%s unexpected response: got=-, want=+:%s\n------------------------------\n",
						errorPrefix, diff)
				}
			}
		}
	}
}

// getComplianceSuite returns the ComplianceSuite read and parsed from the appropriate location
// (current directory if available, otherwise the module path).
func getComplianceSuite() (*genprotopb.ComplianceSuite, error) {
	filePath, err := getComplianceSuiteFile()
	if err != nil {
		return nil, err
	}

	complianceSuiteJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open suite file %q", filePath)
	}

	suite := &genprotopb.ComplianceSuite{}
	if err := protojson.Unmarshal(complianceSuiteJSON, suite); err != nil {
		return nil, fmt.Errorf("error unmarshalling from json %s:\n   file: %s\n   input was: %s", err, filePath, complianceSuiteJSON)
	}
	return suite, nil
}

// getComplianceSuiteFile returns the path to the compliance_suite.json file at the current
// directory, if a file with that name is found there, or in the Showcase module's path
// otherwise. If such a filename cannot be found in either place, this returns an error.
func getComplianceSuiteFile() (string, error) {
	const fileName = "compliance_suite.json"

	// Return the file in the current directory, if it exists there.
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine current directory: %s", err)
	}

	filePath := path.Join(currentDir, fileName)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	// Return the file in the installed module.
	goCmd := exec.Command("go", "env", "GOMODCACHE")
	output, err := goCmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not determine GOMODCACHE: %s", err)
	}
	modCache := strings.TrimSpace(string(output))

	// Look up version of gapic-showcase in the go.mod.
	versionCmd := exec.Command("go", "list", "-m", "-f", "'{{ .Version }}'", "github.com/googleapis/gapic-showcase")
	versionOut, err := versionCmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not determine version of gapic-showcase dependency: %s", err)
	}
	version := strings.TrimSpace(string(versionOut))
	version = strings.TrimFunc(version, func(r rune) bool { return r == '\'' })

	pathInShowcase := path.Join("github.com", "googleapis", "gapic-showcase@"+version, "server", "services", fileName)
	filePath = path.Join(modCache, pathInShowcase)
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("could not determine location of %q; %s", fileName, err)
	}

	return filePath, nil
}
