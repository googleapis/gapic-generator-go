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
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	showcase "github.com/googleapis/gapic-showcase/client"
	genprotopb "github.com/googleapis/gapic-showcase/server/genproto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var complianceClient *showcase.ComplianceClient

// TestComplianceSuite ensures the REST test suite that we require GAPIC generators to pass works
// correctly. GAPIC generators should generate GAPICs for the Showcase API and issue the unary calls
// defined in the test suite using the GAPIC surface. The generators' test should follow the
// high-level logic below, as described in the comments.
func TestComplianceSuite(t *testing.T) {
	defer check(t) // ???
	type method func(ctx context.Context, req *genprotopb.RepeatRequest, opts ...gax.CallOption) (*genprotopb.RepeatResponse, error)

	// Set handlers for each test case. When GAPIC generator tests do this, they should have
	// each of their handlers invoking the correct GAPIC library method for the Showcase API.
	restRPCs := map[string]method{
		"Compliance.RepeatDataBody":                 complianceClient.RepeatDataBody,
		"Compliance.RepeatDataBodyInfo":             complianceClient.RepeatDataBodyInfo,
		"Compliance.RepeatDataQuery":                complianceClient.RepeatDataQuery,
		"Compliance.RepeatDataSimplePath":           complianceClient.RepeatDataSimplePath,
		"Compliance.RepeatDataPathResource":         complianceClient.RepeatDataPathResource,
		"Compliance.RepeatDataPathTrailingResource": complianceClient.RepeatDataPathTrailingResource,
	}

	_ = restRPCs
	/*
		        see if suite file exists in pwd; if not, copy it  from
		         `go env GOMODCACHE`/github.com/googleapis/gapic-showcase@v0.15.0/schema/google/showcase/v1beta1/compliance_suite.json
		(or .../server/services/compliance_suite.json since that symlink does not seem to exist
		)
		       then open the file
	*/

	suite, err := getComplianceSuite()
	if err != nil {
		t.Fatalf(err.Error())
	}
	_ = suite
	//	t.Errorf("*** path is: %s", reflect.TypeOf(showcase.ComplianceClient{}).PkgPath())

	// exec, _ := os.Executable()
	// t.Errorf("*** executable is at: %s", exec)

}

func getComplianceSuite() (*genprotopb.ComplianceSuite, error) {
	filePath, err := getComplianceSuiteFile()
	if err != nil {
		return nil, err
	}

	complianceSuiteJSON, err := ioutil.ReadFile(filePath)
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
// optherwise. If such a filename cannot be found in either place, this returns an error.
func getComplianceSuiteFile() (string, error) {
	const fileName = "compliance_suite.json"
	pathInShowcase := path.Join("github.com", "googleapis", "gapic-showcase@v"+showcaseSemver, "server", "services", fileName)

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

	filePath = path.Join(strings.TrimSpace(string(output)), pathInShowcase)
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("could not determine location of %q; output: %s", fileName, err)
	}

	return filePath, nil
}
