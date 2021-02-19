/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

// Predefinded constants for the JavaScript hooks, they must correspond to the
// error codes used by gRPC, see:
// https://github.com/grpc/grpc-go/blob/master/codes/codes.go
const (
	grpcJSCodes string = `OK = 0;
			CANCELED = 1;
			UNKNOWN = 2;
			INVALIDARGUMENT = 3;
			DEADLINEEXCEEDED = 4;
			NOTFOUND = 5;
			ALREADYEXISTS = 6;
			PERMISSIONDENIED = 7;
			RESOURCEEXHAUSTED = 8;
			FAILEDPRECONDITION = 9;
			ABORTED = 10;
			OUTOFRANGE = 11;
			UNIMPLEMENTED = 12;
			INTERNAL = 13;
			UNAVAILABLE = 14;
			DATALOSS = 15;
			UNAUTHENTICATED = 16`
)
