// Copyright 2019 Google LLC. All Rights Reserved.
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

package monitoring_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/trillian/monitoring"
)

func TestPercentileBucketsInvalid(t *testing.T) {
	for _, inc := range []int64{0, -1, -50, 300, 40000000} {
		t.Run(fmt.Sprintf("increment %d", inc), func(t *testing.T) {
			if got := monitoring.PercentileBuckets(inc); got != nil {
				t.Errorf("PercentileBuckets: got: %v for invalid case, want: nil", got)
			}
		})
	}
}

func TestPercentileBuckets(t *testing.T) {
	for _, inc := range []int64{1, 2, 10, 25, 46, 97} {
		t.Run(fmt.Sprintf("increment %d", inc), func(t *testing.T) {
			buckets := monitoring.PercentileBuckets(inc)
			// The number of buckets expected to be created is fixed.
			if got, want := len(buckets), int(100/inc); math.Abs(float64(got-want)) > 1 {
				t.Errorf("PercentileBuckets(): got len: %d, want: %d", got, want)
			}
			// The first bucket should always be close to 0%.
			if buckets[0] < 0 || buckets[0] > 0.0001 {
				t.Errorf("PercentileBuckets(): got first bucket: %v, want: ~0.0", buckets[0])
			}
			// The last bucket should be on the way towards 100%. It doesn't make a
			// lot of sense to create an extremely coarse grained distribution but
			// it's not actually wrong so no reason to reject it.
			if got, want := math.Abs(buckets[len(buckets)-1]-75.0), 25.0; got > want {
				t.Errorf("PercentileBuckets(): got last bucket diff: %v, want: <%v", got, want)
			}
			// Percentile buckets should increase monotonically.
			for i := 0; i < len(buckets)-1; i++ {
				if buckets[i] > buckets[i+1] {
					t.Errorf("PercentileBuckets(): buckets out of order at index: %d", i)
				}
			}
		})
	}
}

func TestLatencyBuckets(t *testing.T) {
	// Just do some probes on the result to make sure it looks sensible.
	buckets := monitoring.LatencyBuckets()
	checkExpBuckets(t, buckets, 0.04, 1.07, 300)
	// Highest bucket should be about 282 days (allow some leeway).
	expected := 282 * 24 * 3600.0 // 282 days.
	precision := 17 * 3600.0      // A bit less than one day.
	if got := math.Abs(buckets[len(buckets)-1] - expected); got > precision {
		t.Errorf("LatencyBuckets(): got last bucket diff: %v, want: <%v", got, precision)
	}
}

func TestExpBuckets(t *testing.T) {
	for _, tc := range []struct {
		base  float64
		mult  float64
		count uint
	}{
		{base: 1.0, mult: 2.0, count: 20},
		{base: 0.04, mult: 1.07, count: 300},
	} {
		t.Run("", func(t *testing.T) {
			buckets := monitoring.ExpBuckets(tc.base, tc.mult, tc.count)
			checkExpBuckets(t, buckets, tc.base, tc.mult, tc.count)
		})
	}
}

func checkExpBuckets(t *testing.T, buckets []float64, base, mult float64, count uint) {
	t.Helper()
	if got, want := len(buckets), int(count); got != want {
		t.Errorf("unexpected length %d, want %d", got, want)
	}
	// Bucket thresholds should increase monotonically.
	for i := 0; i < len(buckets)-1; i++ {
		if buckets[i] >= buckets[i+1] {
			t.Fatalf("out of order at index %d", i)
		}
	}
	// Lowest bucket should be equal to base.
	if got, want := buckets[0], base; math.Abs(got-want) > 0.001 {
		t.Errorf("got first bucket: %v, want: ~%v", got, want)
	}
	// Highest bucket should be about 86400 sec = 1 day (allow some leeway)/
	last := base * math.Pow(mult, float64(count-1))
	if got, want := buckets[len(buckets)-1], last; math.Abs(got-want) > 0.001 {
		t.Errorf("got last bucket: %v, want: ~%v", got, want)
	}
}
