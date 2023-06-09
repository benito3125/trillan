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

// Package matchers contains additional gomock matchers.
package matchers

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

// AtLeast returns a matcher that requires a number >= n.
func AtLeast(n int) gomock.Matcher {
	return &atLeastMatcher{n}
}

type atLeastMatcher struct {
	num int
}

// Matches tests whether a supplied value, which must be of an int type is
// at least the value the AtLeast matcher expects. If so then it returns true.
func (m atLeastMatcher) Matches(x interface{}) bool {
	if x, ok := x.(int); ok {
		return x >= m.num
	}
	return false
}

func (m atLeastMatcher) String() string {
	return fmt.Sprintf("at least %v", m.num)
}
