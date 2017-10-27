# Copyright 2016 Google, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

VERSION_MAJOR ?= 0
VERSION_MINOR ?= 0
VERSION_BUILD ?= 1

VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

GOOS ?= $(shell go env GOOS)
GOARCH = amd64
BUILD_DIR ?= ./out
ORG := github.com/r2d4
PROJECT := crd
REPOPATH ?= $(ORG)/$(PROJECT)
RELEASE_BUCKET ?= $(PROJECT)

SUPPORTED_PLATFORMS := darwin-$(GOARCH)
BUILD_PACKAGE = $(REPOPATH)

GO_BUILD_TAGS := ""
GO_LDFLAGS := "-X $(REPOPATH)/version.version=$(VERSION)"
GO_FILES := $(shell go list  -f '{{join .Deps "\n"}}' $(BUILD_PACKAGE) | grep $(ORG) | xargs go list -f '{{ range $$file := .GoFiles }} {{$$.Dir}}/{{$$file}}{{"\n"}}{{end}}')

$(BUILD_DIR)/$(PROJECT): $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH)
	cp $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH) $@

$(BUILD_DIR)/$(PROJECT)-%-$(GOARCH): $(GO_FILES) $(BUILD_DIR) main.go
	GOOS=$* GOARCH=$(GOARCH) CGO_ENABLED=0 go build -tags $(GO_BUILD_TAGS) -ldflags $(GO_LDFLAGS) -o $@ $(BUILD_PACKAGE)

%.sha256: %
	shasum -a 256 $< &> $@

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: code-gen
code-gen: $(wildcard pkg/apis/**/*)
	./hack/update-codegen.sh

.PHONY: run
run: $(BUILD_DIR)/$(PROJECT)
	$(BUILD_DIR)/$(PROJECT) -kubeconfig $(HOME)/.kube/config

.PHONY: test
test: $(BUILD_DIR)/$(PROJECT)
	@ ./test.sh

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
