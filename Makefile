# Copyright 2018 The Kubernetes Authors.
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

# This repository builds two commands, mock-driver and csi-sanity,
# but csi-sanity has its own build rules and only mock-driver gets
# published as a container image.
CMDS=mock-driver
all: build build-sanity

include release-tools/build.make

# We have to exclude generic testing of the csi-sanity command because
# the test binary only works in combination with a CSI driver.
# Instead we test with the special ./hack/e2e.sh.
TEST_GO_FILTER_CMD+=| grep -v /cmd/csi-sanity
.PHONY: test-sanity
test: test-sanity
test-sanity:
	@ echo; echo "### $@:"
	./hack/e2e.sh

build-sanity:
	$(MAKE) -C cmd/csi-sanity all
