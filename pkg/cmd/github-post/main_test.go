// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFailures(t *testing.T) {
	testCases := []struct {
		pkgEnv      string
		fileName    string
		packageName string
		testName    string
		message     string
		author      string
	}{
		{
			pkgEnv:      "github.com/cockroachdb/cockroach/pkg/storage",
			fileName:    "stress-failure.json",
			packageName: "github.com/cockroachdb/cockroach/pkg/storage",
			testName:    "TestReplicateQueueRebalance",
			message:     "replicate_queue_test.go:88: condition failed to evaluate within 45s: not balanced: [10 1 10 1 8]",
			author:      "petermattis@gmail.com",
		},
		{
			pkgEnv:      "github.com/cockroachdb/cockroach/pkg/storage",
			fileName:    "stress-fatal.json",
			packageName: "github.com/cockroachdb/cockroach/pkg/storage",
			testName:    "TestGossipHandlesReplacedNode",
			message:     "F180711 20:13:15.826193 83 storage/replica.go:1877  [n?,s1,r1/1:/M{in-ax}] on-disk and in-memory state diverged:",
			author:      "alexdwanerobinson@gmail.com",
		},
		{
			pkgEnv:      "github.com/cockroachdb/cockroach/pkg/storage",
			fileName:    "stress-unknown.json",
			packageName: "github.com/cockroachdb/cockroach/pkg/storage",
			testName:    "(unknown)",
			message:     "make: *** [bin/.submodules-initialized] Error 1",
			author:      "",
		},
		{
			pkgEnv:      "github.com/cockroachdb/cockroach/pkg/util/json",
			fileName:    "stress-subtests.json",
			packageName: "github.com/cockroachdb/cockroach/pkg/util/json",
			testName:    "TestPretty",
			message: `=== RUN   TestPretty/["hello",_["world"]]
    --- FAIL: TestPretty/["hello",_["world"]] (0.00s)
    	json_test.go:1656: injected failure`,
			author: "justin@cockroachlabs.com",
		},
	}
	for _, c := range testCases {
		t.Run(c.fileName, func(t *testing.T) {
			if err := os.Setenv("PKG", c.pkgEnv); err != nil {
				t.Fatal(err)
			}

			file, err := os.Open(filepath.Join("testdata", c.fileName))
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			f := func(_ context.Context, packageName, testName, testMessage, author string) error {
				if c.packageName != packageName {
					t.Fatalf("expected %s, but got %s", c.packageName, packageName)
				}
				if c.testName != testName {
					t.Fatalf("expected %s, but got %s", c.testName, testName)
				}
				if c.author != author {
					t.Fatalf("expected %s, but got %s", c.author, author)
				}
				if !strings.Contains(testMessage, c.message) {
					t.Fatalf("expected message containing %q, but got %s", c.message, testMessage)
				}
				return nil
			}

			if err := listFailures(context.Background(), file, f); err != nil {
				t.Fatal(err)
			}
		})
	}
}
