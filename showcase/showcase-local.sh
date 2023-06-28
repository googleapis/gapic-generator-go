# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# To run the gapic-generator-go Showcase tests using a local version of Showcase
# (eg in development), set the following variables before calling the
# test-with-local-showcase function below.

GO_GENERATOR="$(pwd)/.."
API_COMMON_PROTOS="" # path to local version of api-common-protos
LOCAL_SHOWCASE_REPO="" # path to local version of gapic-showcase

test-with-local-showcase(){

  [[ -n "${GO_GENERATOR}" ]] || { echo '${GO_GENERATOR} must be set' ; return 1 ; }
  [[ -n "${API_COMMON_PROTOS}" ]] || { echo '${API_COMMON_PROTOS} must be set' ; return 1 ; }
  [[ -n "${LOCAL_SHOWCASE_REPO}" ]] || { echo '${LOCAL_SHOWCASE_REPO} must be set' ; return 1 ; }

  GO_GENERATOR="$(realpath ${GO_GENERATOR})"
  API_COMMON_PROTOS="$(realpath ${API_COMMON_PROTOS})"
  LOCAL_SHOWCASE_REPO="$(realpath ${LOCAL_SHOWCASE_REPO})"
  SHOWCASE_SCHEMA="${LOCAL_SHOWCASE_REPO}/schema/google/showcase/v1beta1"


  cat <<EOF
Running with:
  GO_GENERATOR:        ${GO_GENERATOR}
  API_COMMON_PROTOS:   ${API_COMMON_PROTOS}
  LOCAL_SHOWCASE_REPO: ${LOCAL_SHOWCASE_REPO}
EOF


  pushd "${GO_GENERATOR}" >& /dev/null
  go install ./cmd/protoc-gen-go_gapic

  cd showcase

  GAPIC_OUT_DIR="./gen"
  rm -rf $GAPIC_OUT_DIR
  mkdir -p $GAPIC_OUT_DIR
  protoc --experimental_allow_proto3_optional -I $API_COMMON_PROTOS  \
    --go_gapic_out $GAPIC_OUT_DIR \
    --go_gapic_opt 'transport=rest+grpc' \
    --go_gapic_opt 'rest-numeric-enums' \
    --go_gapic_opt='go-gapic-package=github.com/googleapis/gapic-showcase/client;client' \
    --go_gapic_opt=grpc-service-config=${SHOWCASE_SCHEMA}/showcase_grpc_service_config.json \
    --go_gapic_opt=api-service-config=${SHOWCASE_SCHEMA}/showcase_v1beta1.yaml \
    $PLUGIN --proto_path=$SHOWCASE_SCHEMA $SHOWCASE_SCHEMA/*.proto

  go mod edit -replace=github.com/googleapis/gapic-showcase=./gen/github.com/googleapis/gapic-showcase

  pushd ./gen/github.com/googleapis/gapic-showcase  >& /dev/null
  go mod init gapic-showcase
  mkdir server
  ln -s ${LOCAL_SHOWCASE_REPO}/server/genproto server/
  popd  >& /dev/null

  # ensure your Showcase server is running

  cp ${LOCAL_SHOWCASE_REPO}/server/services/compliance_suite.json .
  go test -mod=mod -count=1 ./...

  popd  >& /dev/null
}

test-with-local-showcase
