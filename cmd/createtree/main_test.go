// Copyright 2017 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/trillian"
	"github.com/google/trillian/testonly"
	"github.com/google/trillian/testonly/flagsaver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

// defaultTree reflects all flag defaults with the addition of a valid private key.
var defaultTree = &trillian.Tree{
	TreeState:       trillian.TreeState_ACTIVE,
	TreeType:        trillian.TreeType_LOG,
	MaxRootDuration: durationpb.New(0 * time.Millisecond),
}

type testCase struct {
	desc        string
	setFlags    func()
	validateErr error
	createErr   error
	initErr     error
	wantErr     bool
	wantTree    *trillian.Tree
}

func TestCreateTree(t *testing.T) {
	nonDefaultTree := proto.Clone(defaultTree).(*trillian.Tree)
	nonDefaultTree.TreeType = trillian.TreeType_LOG
	nonDefaultTree.DisplayName = "Llamas Log"
	nonDefaultTree.Description = "For all your digital llama needs!"

	runTest(t, []*testCase{
		{
			desc: "validOpts",
			// runTest sets mandatory options, so no need to provide a setFlags func.
			wantTree: defaultTree,
		},
		{
			desc: "nonDefaultOpts",
			setFlags: func() {
				*treeType = nonDefaultTree.TreeType.String()
				*displayName = nonDefaultTree.DisplayName
				*description = nonDefaultTree.Description
			},
			wantTree: nonDefaultTree,
		},
		{
			desc: "mandatoryOptsNotSet",
			// Undo the flags set by runTest, so that mandatory options are no longer set.
			setFlags:    flagsaver.Save().MustRestore,
			validateErr: errAdminAddrNotSet,
			wantErr:     true,
		},
		{
			desc:        "emptyAddr",
			setFlags:    func() { *adminServerAddr = "" },
			validateErr: errAdminAddrNotSet,
			wantErr:     true,
		},
		{
			desc:        "invalidEnumOpts",
			setFlags:    func() { *treeType = "LLAMA!" },
			validateErr: errors.New("unknown TreeType"),
			wantErr:     true,
		},
		{
			desc:      "createErr",
			createErr: status.Errorf(codes.Unavailable, "create tree failed"),
			wantErr:   true,
		},
		{
			desc: "logInitErr",
			setFlags: func() {
				nonDefaultTree.TreeType = trillian.TreeType_LOG
				*treeType = nonDefaultTree.TreeType.String()
			},
			wantTree: defaultTree,
			initErr:  status.Errorf(codes.Unavailable, "log init failed"),
			wantErr:  true,
		},
	})
}

// runTest executes the createtree command against a fake TrillianAdminServer
// for each of the provided tests, and checks that the tree in the request is
// as expected, or an expected error occurs.
// Prior to each test case, it:
// 1. Resets all flags to their original values.
// 2. Sets the adminServerAddr flag to point to the fake server.
// 3. Calls the test's setFlags func (if provided) to allow it to change flags specific to the test.
func runTest(t *testing.T, tests []*testCase) {
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			// Note: Restore() must be called after the flag-reading bits are
			// stopped, otherwise there might be a data race.
			defer flagsaver.Save().MustRestore()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s, stopFakeServer, err := testonly.NewMockServer(ctrl)
			if err != nil {
				t.Fatalf("Error starting fake server: %v", err)
			}
			defer stopFakeServer()
			*adminServerAddr = s.Addr
			if tc.setFlags != nil {
				tc.setFlags()
			}

			call := s.Admin.EXPECT().CreateTree(gomock.Any(), gomock.Any()).Return(tc.wantTree, tc.createErr)
			expectCalls(call, tc.createErr, tc.validateErr)
			switch *treeType {
			case "LOG":
				call := s.Log.EXPECT().InitLog(gomock.Any(), gomock.Any()).Return(&trillian.InitLogResponse{}, tc.initErr)
				expectCalls(call, tc.initErr, tc.validateErr, tc.createErr)
				call = s.Log.EXPECT().GetLatestSignedLogRoot(gomock.Any(), gomock.Any()).Return(&trillian.GetLatestSignedLogRootResponse{}, nil)
				expectCalls(call, nil, tc.validateErr, tc.createErr, tc.initErr)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_, err = createTree(ctx)
			if hasErr := err != nil; hasErr != tc.wantErr {
				t.Errorf("createTree() '%v', wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

// expectCalls returns the minimum number of times a function is expected to be called
// given the return error for the function (err), and all previous errors in the function's
// code path.
func expectCalls(call *gomock.Call, err error, prevErr ...error) *gomock.Call {
	// If a function prior to this function errored,
	// we do not expect this function to be called.
	for _, e := range prevErr {
		if e != nil {
			return call.Times(0)
		}
	}
	// If this function errors, it will be retried multiple times.
	if err != nil {
		return call.MinTimes(2)
	}
	// If this function succeeds it should only be called once.
	return call.Times(1)
}
