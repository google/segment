// Copyright (c) 2018, Google Inc. All rights reserved.
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

package segment

import (
	"math"
	"reflect"
	"testing"
)

// Please keep the order of test functions the same as
// the order of methods/functions in segment.go.

func TestSegmentString(t *testing.T) {
	testCases := []struct {
		input Segment
		want  string
	}{
		{
			input: Segment{11, 13},
			want:  "[start: 11, end: 13]",
		},
		{
			input: Segment{313, 313},
			want:  "[start: 313, end: 313]",
		},
	}

	for _, test := range testCases {
		if got := test.input.String(); got != test.want {
			t.Errorf("s.String() = %s, should be %s", got, test.want)
		}
	}
}

func TestSegmentsString(t *testing.T) {
	testCases := []struct {
		input Segments
		want  string
	}{
		{
			input: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 5},
			},
			want: "[start: 2, end: 3], [start: 1, end: 2], [start: 4, end: 5]",
		},
	}

	for _, test := range testCases {
		if got := test.input.String(); got != test.want {
			t.Errorf("s.String() = %s, should be %s", got, test.want)
		}
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		start, end int64
		want       Segment
		wanterr    bool
	}{
		{
			start:   0,
			end:     1,
			want:    Segment{0, 1},
			wanterr: false,
		},
		{
			start:   0,
			end:     -1,
			want:    Segment{},
			wanterr: true,
		},
	}

	for _, test := range testCases {
		got, goterr := New(test.start, test.end)
		if got != test.want || (goterr != nil) != test.wanterr {
			t.Errorf("New(%d, %d) = %s, should be %s; got error? %t, want error? %t",
				test.start, test.end, got, test.want, goterr != nil, test.wanterr)
		}
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		s, want    Segment
		start, end int64
		wanterr    bool
	}{
		{
			s:       Segment{1, 2},
			start:   0,
			end:     1,
			want:    Segment{0, 1},
			wanterr: false,
		},
		{
			s:       Segment{1, 2},
			start:   0,
			end:     -1,
			want:    Segment{1, 2},
			wanterr: true,
		},
	}

	for _, test := range testCases {
		sCopy := Segment{test.s.start, test.s.end}
		goterr := test.s.Update(test.start, test.end)
		if test.s != test.want || (goterr != nil) != test.wanterr {
			t.Errorf("%s.Update(%d, %d) = %s, should be %s; got error? %t, want error? %t",
				sCopy, test.start, test.end, test.s, test.want, goterr != nil, test.wanterr)
		}
	}
}

func TestUpdateStart(t *testing.T) {
	testCases := []struct {
		s, want Segment
		start   int64
		wanterr bool
	}{
		{
			s:       Segment{1, 2},
			start:   0,
			want:    Segment{0, 2},
			wanterr: false,
		},
		{
			s:       Segment{1, 2},
			start:   3,
			want:    Segment{1, 2},
			wanterr: true,
		},
	}

	for _, test := range testCases {
		sCopy := Segment{test.s.start, test.s.end}
		goterr := test.s.UpdateStart(test.start)
		if test.s != test.want || (goterr != nil) != test.wanterr {
			t.Errorf("%s.UpdateStart(%d) = %s, should be %s; got error? %t, want error? %t",
				sCopy, test.start, test.s, test.want, goterr != nil, test.wanterr)
		}
	}
}

func TestUpdateEnd(t *testing.T) {
	testCases := []struct {
		s, want Segment
		end     int64
		wanterr bool
	}{
		{
			s:       Segment{1, 2},
			end:     0,
			want:    Segment{1, 2},
			wanterr: true,
		},
		{
			s:       Segment{1, 2},
			end:     3,
			want:    Segment{1, 3},
			wanterr: false,
		},
	}

	for _, test := range testCases {
		sCopy := Segment{test.s.start, test.s.end}
		goterr := test.s.UpdateEnd(test.end)
		if test.s != test.want || (goterr != nil) != test.wanterr {
			t.Errorf("%s.UpdateEnd(%d) = %s, should be %s; got error? %t, want error? %t",
				sCopy, test.end, test.s, test.want, goterr != nil, test.wanterr)
		}
	}
}

func TestStarts(t *testing.T) {
	testCases := []struct {
		input Segments
		want  []int64
	}{
		{
			input: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 5},
			},
			want: []int64{2, 1, 4},
		},
		{
			input: Segments{
				Segment{0, 1},
				Segment{-1, 5},
				Segment{4, 6},
			},
			want: []int64{0, -1, 4},
		},
		{
			input: Segments{
				Segment{3, 3},
				Segment{4, 5},
			},
			want: []int64{3, 4},
		},
		{
			input: Segments{},
			want:  nil,
		},
		{
			input: Segments{
				Segment{0, 0},
			},
			want: []int64{0},
		},
	}

	for _, test := range testCases {
		if got := test.input.Starts(); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s.Starts() = %v, should be %v", test.input, got, test.want)
		}
	}
}

func TestEnds(t *testing.T) {
	testCases := []struct {
		input Segments
		want  []int64
	}{
		{
			input: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 5},
			},
			want: []int64{3, 2, 5},
		},
		{
			input: Segments{
				Segment{0, 1},
				Segment{-1, 5},
				Segment{4, 6},
			},
			want: []int64{1, 5, 6},
		},
		{
			input: Segments{
				Segment{3, 3},
				Segment{4, 5},
			},
			want: []int64{3, 5},
		},
		{
			input: Segments{},
			want:  nil,
		},
		{
			input: Segments{
				Segment{0, 0},
			},
			want: []int64{0},
		},
	}

	for _, test := range testCases {
		if got := test.input.Ends(); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s.Ends() = %v, should be %v", test.input, got, test.want)
		}
	}
}

func TestDelta(t *testing.T) {
	testCases := []struct {
		input Segment
		want  int64
	}{
		{
			input: Segment{11, 13},
			want:  2,
		},
		{
			input: Segment{313, 313},
			want:  0,
		},
	}

	for _, test := range testCases {
		if got := test.input.Delta(); got != test.want {
			t.Errorf("%s.Delta() = %d, should be %d", test.input, got, test.want)
		}
	}
}

func TestIsDeltaPositive(t *testing.T) {
	testCases := []struct {
		input Segment
		want  bool
	}{
		{
			input: Segment{1, 2},
			want:  true,
		},
		{
			input: Segment{1, 1},
			want:  false,
		},
	}

	for _, test := range testCases {
		if got := test.input.IsDeltaPositive(); got != test.want {
			t.Errorf("%s.IsDeltaPositive() = %t, should be %t", test.input, got, test.want)
		}
	}
}

func TestSumDeltas(t *testing.T) {
	testCases := []struct {
		input Segments
		want  int64
	}{
		{
			input: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 5},
			},
			want: 1 + 1 + 1,
		},
		{
			input: Segments{
				Segment{0, 1},
				Segment{-1, 5},
				Segment{4, 6},
			},
			want: 1 + 6 + 2,
		},
		{
			input: Segments{
				Segment{3, 3},
				Segment{4, 5},
			},
			want: 0 + 1,
		},
	}

	for _, test := range testCases {
		if got := test.input.SumDeltas(); got != test.want {
			t.Errorf("%s.SumDeltas() = %d, should be %d", test.input, got, test.want)
		}
	}
}

func TestSegmentsWithPredicate(t *testing.T) {
	testCases := []struct {
		predicateDescription string
		ss                   Segments
		pred                 func(s Segment) bool
		want                 Segments
	}{
		{
			predicateDescription: "segment delta is positive",
			ss: Segments{
				Segment{2, 2},
				Segment{1, 2}},
			pred: func(s Segment) bool { return s.Delta() > 0 },
			want: Segments{
				Segment{1, 2}},
		},
		{
			predicateDescription: "segment delta is positive",
			ss: Segments{
				Segment{2, 2},
				Segment{1, 1}},
			pred: func(s Segment) bool { return s.Delta() > 0 },
			want: nil,
		},
	}

	for _, test := range testCases {
		if got := SegmentsWithPredicate(test.ss, test.pred); !reflect.DeepEqual(got, test.want) {
			t.Errorf("SegmentsWithPredicate(%s, predicate: %s) = %s, want %s",
				test.ss, test.predicateDescription, got, test.want)
		}
	}
}

func TestRemoveOverlaps(t *testing.T) {
	testCases := []struct {
		input, want Segments
	}{
		{
			input: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 5},
			},
			want: Segments{
				Segment{1, 3},
				Segment{4, 5},
			},
		},
		{
			input: Segments{
				Segment{0, 1},
				Segment{-1, 5},
				Segment{4, 6},
			},
			want: Segments{
				Segment{-1, 6},
			},
		},
		{
			input: Segments{
				Segment{3, 3},
				Segment{4, 5},
			},
			want: Segments{
				Segment{3, 3},
				Segment{4, 5},
			},
		},
		{
			input: Segments{
				Segment{int64(math.MinInt64), 3},
				Segment{4, 5},
			},
			want: Segments{
				Segment{4, 5},
			},
		},
	}

	for _, test := range testCases {
		if got := RemoveOverlaps(test.input); !reflect.DeepEqual(got, test.want) {
			t.Errorf("RemoveOverlaps(%s) is %s, should be %s", test.input, got, test.want)
		}
	}
}

func TestUnionWithTwoInputs(t *testing.T) {
	testCases := []struct {
		description string
		x, y, want  Segments
	}{
		{
			description: "x is empty",
			x:           Segments{},
			y: Segments{
				Segment{1, 3},
				Segment{2, 4},
			},
			want: Segments{
				Segment{1, 4},
			},
		},
		{
			description: "no overlap between x and y",
			x: Segments{
				Segment{1, 5},
				Segment{2, 10},
			},
			y: Segments{
				Segment{-5, -1},
				Segment{-10, -2},
			},
			want: Segments{
				Segment{-10, -1},
				Segment{1, 10},
			},
		},
		{
			description: "some overlap between x and y",
			x: Segments{
				Segment{1, 5},
				Segment{2, 10},
			},
			y: Segments{
				Segment{3, 12},
				Segment{14, 15},
			},
			want: Segments{
				Segment{1, 12},
				Segment{14, 15},
			},
		},
	}

	for _, test := range testCases {
		if got := Union(test.x, test.y); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: Union(%s, %s) = %s, want %s",
				test.description, test.x, test.y, got, test.want)
		}
	}
}

func TestUnionWithVaryingNumberOfInputs(t *testing.T) {

	// Test Union with nil input.
	if got := Union(nil); got != nil {
		t.Errorf("Union(nil) = %s, want nil", got)
	}

	s1 := Segments{
		Segment{1, 3},
		Segment{2, 4},
	}
	s2 := Segments{
		Segment{-1, 2},
		Segment{1, 2},
	}
	s3 := Segments{
		Segment{8, 10},
		Segment{-2, 0},
	}

	// Test Union with single slice input.
	// Results should be the same as RemoveOverlaps.
	if got, want := Union(s1), (Segments{Segment{1, 4}}); !reflect.DeepEqual(got, want) {
		t.Errorf("Union(%s) = %s, want %s", s1, got, want)
	}

	// Test Union with three slices as input.
	// Union(x, y, z) should be equivalent to Union(Union(x, y), z).
	if got, want := Union(s1, s2, s3), (Segments{Segment{-2, 4}, Segment{8, 10}}); !reflect.DeepEqual(got, want) {
		t.Errorf("Union(%s, %s, %s) = %s, want %s", s1, s2, s3, got, want)
	}
}

func TestSimpleIntersection(t *testing.T) {
	testCases := []struct {
		s, t    Segment
		wantSeg Segment
		wantOk  bool
	}{
		{
			s:       Segment{2, 30},
			t:       Segment{20, 40},
			wantSeg: Segment{20, 30},
			wantOk:  true,
		},
		{
			s:       Segment{20, 40},
			t:       Segment{2, 30},
			wantSeg: Segment{20, 30},
			wantOk:  true,
		},
		{
			s:       Segment{-20, 40},
			t:       Segment{2, 30},
			wantSeg: Segment{2, 30},
			wantOk:  true,
		},
		{
			s:       Segment{-20, 40},
			t:       Segment{-200, 300},
			wantSeg: Segment{-20, 40},
			wantOk:  true,
		},
		{
			s:      Segment{-20, 10},
			t:      Segment{200, 300},
			wantOk: false,
		},
		{
			s:       Segment{0, 1},
			t:       Segment{1, 2},
			wantSeg: Segment{1, 1},
			wantOk:  true,
		},
	}

	for _, test := range testCases {
		if gotSeg, gotOk := SimpleIntersection(test.s, test.t); gotOk != test.wantOk || (gotOk && !reflect.DeepEqual(gotSeg, test.wantSeg)) {
			t.Errorf("SimpleIntersection(%s, %s) is %s, %t, expected %s, %t", test.s, test.t, gotSeg, gotOk, test.wantSeg, test.wantOk)
		}
	}
}

func TestIntersect(t *testing.T) {
	testCases := []struct {
		description string
		x, y, want  Segments
	}{
		{
			description: "x is empty",
			x:           Segments{},
			y:           Segments{Segment{1, 5}},
			want:        nil,
		},
		{
			description: "y is empty",
			x:           Segments{Segment{1, 5}},
			y:           Segments{},
			want:        nil,
		},
		{
			description: "no overlap between x and y",
			x: Segments{
				Segment{1, 5},
				Segment{2, 10},
			},
			y: Segments{
				Segment{-5, -1},
				Segment{-10, -2},
			},
			want: nil,
		},
		{
			description: "some overlap between x and y",
			x: Segments{
				Segment{1, 5},
				Segment{2, 10},
				Segment{12, 16},
			},
			y: Segments{
				Segment{3, 7},
				Segment{11, 15},
			},
			want: Segments{
				Segment{3, 7},
				Segment{12, 15},
			},
		},
		{
			description: "infinitesimal point overlap between x and y",
			x: Segments{
				Segment{1, 5},
				Segment{2, 10},
				Segment{12, 16},
			},
			y: Segments{
				Segment{3, 7},
				Segment{16, 17},
			},
			want: Segments{
				Segment{3, 7},
				Segment{16, 16},
			},
		},
	}

	for _, test := range testCases {
		if got := Intersect(test.x, test.y); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: Intersect(%s, %s) = %s, want %s",
				test.description, test.x, test.y, got, test.want)
		}
	}
}

func TestGetOverlaps(t *testing.T) {
	testCases := []struct {
		s, want Segments
	}{
		{
			s: Segments{
				Segment{2, 30},
				Segment{40, 50},
				Segment{60, 80},
			},
			want: nil,
		},
		{
			s: Segments{
				Segment{2, 30},
				Segment{10, 50},
				Segment{60, 80},
			},
			want: Segments{
				Segment{10, 30},
			},
		},
		{
			s: Segments{
				Segment{2, 30},
				Segment{10, 50},
				Segment{50, 80},
			},
			want: Segments{
				Segment{10, 30},
				Segment{50, 50},
			},
		},
	}

	for _, test := range testCases {
		if got := GetOverlaps(test.s); !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOverlaps(%s) is %s, expected is %s", test.s, got, test.want)
		}
	}
}

func TestComplement(t *testing.T) {
	testCases := []struct {
		description string
		ss          Segments
		superset    Segment
		want        Segments
	}{
		{
			description: "set of segments is empty",
			ss:          nil,
			superset:    Segment{0, 100},
			want: Segments{
				Segment{0, 100},
			},
		},
		{
			description: "segments only contains a point segment",
			ss: Segments{
				Segment{2, 2},
			},
			superset: Segment{0, 3},
			want: Segments{
				Segment{0, 2},
				Segment{2, 3},
			},
		},
		{
			description: "superset is empty",
			ss: Segments{
				Segment{1, 2},
				Segment{3, 4},
			},
			superset: Segment{},
			want:     nil,
		},
		{
			description: "superset is a true superset of the segments",
			ss: Segments{
				Segment{1, 3},
				Segment{0, 2},
				Segment{4, 6},
			},
			superset: Segment{-10, 6},
			want: Segments{
				Segment{-10, 0},
				Segment{3, 4},
			},
		},
		{
			description: "superset.end lies within a segment",
			ss: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
			superset: Segment{-10, 5},
			want: Segments{
				Segment{-10, 0},
				Segment{2, 4},
			},
		},
		{
			description: "superset.end lies between segments",
			ss: Segments{
				Segment{4, 6},
				Segment{0, 2},
			},
			superset: Segment{-10, 3},
			want: Segments{
				Segment{-10, 0},
				Segment{2, 3},
			},
		},
		{
			description: "superset lies completely below all segments",
			ss: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
			superset: Segment{-10, -1},
			want: Segments{
				Segment{-10, -1},
			},
		},
		{
			description: "superset.start lies within a segment",
			ss: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
			superset: Segment{1, 10},
			want: Segments{
				Segment{2, 4},
				Segment{6, 10},
			},
		},
		{
			description: "superset.start lies between segments",
			ss: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
			superset: Segment{3, 10},
			want: Segments{
				Segment{3, 4},
				Segment{6, 10},
			},
		},
		{
			description: "superset lies completely above the segments",
			ss: Segments{
				Segment{4, 6},
				Segment{0, 2},
			},
			superset: Segment{8, 10},
			want: Segments{
				Segment{8, 10},
			},
		},
		{
			description: "superset is a sub-segment of RemoveOverlaps(ss)",
			ss: Segments{
				Segment{1, 3},
				Segment{0, 2},
				Segment{4, 6},
			},
			superset: Segment{0, 3},
			want:     nil,
		},
		{
			description: "real-life session example part 1",
			ss: Segments{
				Segment{0, 30000},
				Segment{36571515, 36901489},
			},
			superset: Segment{0, 29347},
			want:     nil,
		},
		{
			description: "real-life session example part 2",
			ss: Segments{
				Segment{0, 30000},
				Segment{36571515, 36901489},
			},
			superset: Segment{36569394, 36596094},
			want: Segments{
				Segment{36569394, 36571515},
			},
		},
	}

	for _, test := range testCases {
		if got := Complement(test.superset, test.ss); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: Complement(superset = %s, ss = %s) = %s, want %s",
				test.description, test.superset, test.ss, got, test.want)
		}
	}
}

func TestSetDiff(t *testing.T) {
	testCases := []struct {
		description string
		x, y, want  Segments
	}{
		{
			description: "y is empty",
			x: Segments{
				Segment{0, 2},
			},
			y: Segments{},
			want: Segments{
				Segment{0, 2},
			},
		},
		{
			description: "x is empty",
			x:           Segments{},
			y: Segments{
				Segment{0, 2},
			},
			want: nil,
		},
		{
			description: "there is some overlap between x and y",
			x: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
			y: Segments{
				Segment{1, 3},
				Segment{3, 5},
			},
			want: Segments{
				Segment{0, 1},
				Segment{5, 6},
			},
		},
		{
			description: "there is no overlap between x and y",
			x: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
			y: Segments{
				Segment{-6, -4},
				Segment{-2, 0},
			},
			want: Segments{
				Segment{0, 2},
				Segment{4, 6},
			},
		},
		{
			description: "there is some overlap within x",
			x: Segments{
				Segment{0, 4},
				Segment{2, 6},
			},
			y: Segments{
				Segment{1, 3},
				Segment{3, 5},
			},
			want: Segments{
				Segment{0, 1},
				Segment{5, 6},
			},
		},
		{
			description: "there is some overlap within y",
			x: Segments{
				Segment{0, 2},
				Segment{4, 5},
			},
			y: Segments{
				Segment{1, 5},
				Segment{3, 7},
			},
			want: Segments{
				Segment{0, 1},
			},
		},
		{
			description: "real-life session example",
			x: Segments{
				Segment{0, 29347},
				Segment{36569394, 36596094},
			},
			y: Segments{
				Segment{0, 30000},
				Segment{36571515, 36901489},
			},
			want: Segments{
				Segment{36569394, 36571515},
			},
		},
	}

	for _, test := range testCases {
		if got := SetDiff(test.x, test.y); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: SetDiff(%s, %s) = %s, want %s",
				test.description, test.x, test.y, got, test.want)
		}
	}
}

func TestIsIntersectionEmpty(t *testing.T) {
	testCases := []struct {
		s    Segment
		tt   Segments
		want bool
	}{
		{
			s: Segment{11, 13},
			tt: Segments{
				Segment{12, 14},
				Segment{20, 25},
			},
			want: false,
		},
		{
			s: Segment{11, 13},
			tt: Segments{
				Segment{120, 140},
				Segment{20, 25},
			},
			want: true,
		},
		{
			s:    Segment{11, 13},
			tt:   Segments{},
			want: true,
		},
		{
			s: Segment{1, 1},
			tt: Segments{
				Segment{0, 2},
			},
			want: false,
		},
	}

	for _, test := range testCases {
		if got := test.s.IsIntersectionEmpty(test.tt); got != test.want {
			t.Errorf("%s.IsIntersectionEmpty(%s) = %t, should be %t", test.s, test.tt, got, test.want)
		}
	}
}

func TestIsSubSegment(t *testing.T) {
	testCases := []struct {
		s, t Segment
		want bool
	}{
		{
			s:    Segment{2, 30},
			t:    Segment{20, 40},
			want: false,
		},
		{
			s:    Segment{4, 6},
			t:    Segment{3, 10},
			want: true,
		},
		{
			s:    Segment{-10, 60},
			t:    Segment{3, 10},
			want: false,
		},
	}

	for _, test := range testCases {
		if got := test.s.IsSubSegment(test.t); got != test.want {
			t.Errorf("%s.IsSubSegment(%s) is %t, expected is %t", test.s, test.t, got, test.want)
		}
	}
}

func TestIsPointInSegment(t *testing.T) {
	testCases := []struct {
		s    Segment
		p    int64
		want bool
	}{
		{
			s:    Segment{2, 30},
			p:    11,
			want: true},
		{
			s:    Segment{-10, 30},
			p:    41,
			want: false},
	}

	for _, test := range testCases {
		if got := IsPointInSegment(test.p, test.s); got != test.want {
			t.Errorf("IsPointInSegment(%d, %s) is %t, expected is %t", test.p, test.s, got, test.want)
		}
	}
}

func TestIsPointInSegments(t *testing.T) {
	testCases := []struct {
		ss   Segments
		p    int64
		want bool
	}{
		{
			ss: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 15}},
			p:    11,
			want: true},
		{
			ss: Segments{
				Segment{2, 3},
				Segment{1, 2},
				Segment{4, 5}},
			p:    41,
			want: false,
		},
	}

	for _, test := range testCases {
		if got := IsPointInSegments(test.p, test.ss); got != test.want {
			t.Errorf("IsPointInSegments(%s, %d) = %t, want %t", test.ss, test.p, got, test.want)
		}
	}
}
