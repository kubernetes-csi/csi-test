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

# This repository builds one command, csi-sanity, but it has its own build
# rules no image to be published.
CMDS=
all: build build-sanity

include release-tools/build.make

# We have to exclude generic testing of the csi-sanity command because
# the test binary only works in combination with a CSI driver.
# Instead we test with the special ./hack/e2e.sh and the
# csi-driver-host-path that we build from source.
TEST_GO_FILTER_CMD+=| grep -v /cmd/csi-sanity
.PHONY: test-sanity
test: test-sanity
test-sanity: bin/hostpathplugin
	@ echo; echo "### $@:"
	if [ $$(id -u) = 0 ]; then \
		./hack/e2e.sh; \
	else \
		sudo ./hack/e2e.sh; \
	fi

build-sanity:
	$(MAKE) -C cmd/csi-sanity all


TEST_HOSTPATH_VERSION=v1.14.1
TEST_HOSTPATH_SOURCE=bin/hostpath-source
TEST_HOSTPATH_REPO=https://github.com/kubernetes-csi/csi-driver-host-path.git
bin/hostpathplugin:
	mkdir -p $(@D)
	if ! [ -d $(TEST_HOSTPATH_SOURCE) ]; then \
		mkdir -p $(dir $(TEST_HOSTPATH_SOURCE)) && \
		git clone $(TEST_HOSTPATH_REPO) $(TEST_HOSTPATH_SOURCE); \
	fi
	cd $(TEST_HOSTPATH_SOURCE) && git checkout $(TEST_HOSTPATH_VERSION)
	make -C $(TEST_HOSTPATH_SOURCE) build
	ln -fs $(abspath $(TEST_HOSTPATH_SOURCE))/bin/hostpathplugin $@
