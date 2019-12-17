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

// Package segment provides mathematical operators on generic start-end line segments.
package segment

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

//////// CORE TYPES ////////

// Segment is a generic "start, end" line segment structure.
type Segment struct {
	start, end int64
}

// Segments is a slice of type Segment objects.
type Segments []Segment

//////// PRINT AS STRING ////////

// String returns the values of a segment in a string.
func (s Segment) String() string {
	return fmt.Sprintf("[start: %d, end: %d]", s.start, s.end)
}

// String returns the values of an array of Segments in a string.
func (ss Segments) String() string {
	var output []string
	for _, s := range ss {
		output = append(output, s.String())
	}
	return strings.Join(output, ", ")
}

//////// CREATE/UPDATE SEGMENT VALUES ////////

// New creates a Segment struct from a start and an end., If end < start,
// the Segment is not well-defined, so an error is returned, and the output segment is nil.
func New(start, end int64) (Segment, error) {
	if end < start {
		return Segment{}, fmt.Errorf("end < start: nil segment returned")
	}
	return Segment{start, end}, nil
}

// Update sets the values of a Segment in place.
// The values are only set if start <= end. If not, an error is returned.
func (s *Segment) Update(start, end int64) error {
	if start <= end {
		s.start, s.end = start, end
		return nil
	}
	return fmt.Errorf("end < start: segment not updated")
}

// UpdateStart sets the start value of a Segment, if the start is less than the existing end.
// If this is not the case, the start is not updated, and an error is returned.
func (s *Segment) UpdateStart(start int64) error {
	if start <= s.end {
		s.start = start
		return nil
	}
	return fmt.Errorf("new start > existing end: segment start not updated")
}

// UpdateEnd sets the end value of a Segment, if the end is greater than the existing start.
// If this is not the case, the end is not updated, and an error is returned.
func (s *Segment) UpdateEnd(end int64) error {
	if end >= s.start {
		s.end = end
		return nil
	}
	return fmt.Errorf("new end < existing start: segment end not updated")
}

// LinearTransform performs a linear transformation on a segment.
// If the multiplier is negative, the resulting segment will not be well-defined,
// so the linear transform is not performed on the segment.
// The coefficients of the linear transformation are float64, so the output
// will be rounded as per package :float.
func (s *Segment) LinearTransform(a, b float64) error {
	if a < 0 {
		return fmt.Errorf("a < 0: linear transform not performed on segment")
	}
	start, end := float64(s.start)*a+b, float64(s.end)*a+b
	s.Update(int64(math.Round(start)), int64(math.Round(end)))
	return nil
}

// LinearTransform performs a linear transformation on a slice of segments.
// If the multiplier is negative, the resulting segments will not be well-defined,
// so the linear transform is not performed on ANY of the segments.
// The coefficients of the linear transformation are float64, so all segments
// will be rounded as per package :float.
func (ss Segments) LinearTransform(a, b float64) error {
	if a < 0 {
		return fmt.Errorf("a < 0: linear transform not performed on any segment")
	}
	for i := range ss {
		ss[i].LinearTransform(a, b)
	}
	return nil
}

//////// EXTRACT SEGMENT VALUES/CHARACTERISTICS ////////

// Start returns the start of a segment.
func (s Segment) Start() int64 {
	return s.start
}

// End returns the end of a segment.
func (s Segment) End() int64 {
	return s.end
}

// Starts returns the starts of Segments in a slice.
func (ss Segments) Starts() []int64 {
	var output []int64
	for _, s := range ss {
		output = append(output, s.start)
	}
	return output
}

// Ends returns the ends of Segments in a slice.
func (ss Segments) Ends() []int64 {
	var output []int64
	for _, s := range ss {
		output = append(output, s.end)
	}
	return output
}

// Delta returns the length of the segment.
func (s Segment) Delta() int64 {
	return s.end - s.start
}

// IsDeltaPositive reports whether a segment has positive delta.
func (s Segment) IsDeltaPositive() bool {
	return s.Delta() > 0
}

// SumDeltas returns the sum of Deltas in an array of Segments.
func (ss Segments) SumDeltas() int64 {
	var output int64
	for _, s := range ss {
		output += s.Delta()
	}
	return output
}

// SegmentsWithPredicate returns a subset of segments that meet a predicate function.
func SegmentsWithPredicate(ss Segments, pred func(s Segment) bool) Segments {
	var output Segments
	for _, s := range ss {
		if pred(s) {
			output = append(output, s)
		}
	}
	return output
}

//////// SET OPERATIONS ////////

// RemoveOverlaps takes out overlapping areas in a slice of segments.
// Note segments where start > end are discarded.
// This function also sorts segments by Start, so the output will be ordered.
func RemoveOverlaps(ss Segments) Segments {
	// Sort the segments by their starts.
	// In order to not sort this in-place, we make a copy of ss.
	ssSorted := append(Segments{}, ss...)
	sort.Slice(ssSorted, func(i, j int) bool { return ssSorted[i].start < ssSorted[j].start })
	rightMost := int64(math.MinInt64)
	var output Segments

	for _, s := range ssSorted {
		// Do we need to start a new segment?
		if rightMost < s.start {
			output = append(output, Segment{s.start, s.end})
			rightMost = s.end
			// Ensuring the length of the output is positive means that segments
			// starting at math.MinInt64 are skipped.
		} else if n := len(output); n > 0 && rightMost < s.end {
			// Do we need to update the end of the existing last segment?
			output[n-1].end = s.end
			rightMost = s.end
		}
	}
	return output
}

// Union finds the overlap between slices of segments.
func Union(ss ...Segments) Segments {
	var tt Segments
	for _, s := range ss {
		tt = append(tt, s...)
	}
	return RemoveOverlaps(tt)
}

// SimpleIntersection returns the intersection between segment s and segment t
// (and a bool indicating whether there is an intersection).
func SimpleIntersection(s, t Segment) (Segment, bool) {
	switch {
	case s.start <= t.start && s.end <= t.end && s.end >= t.start:
		return Segment{t.start, s.end}, true
	case s.start >= t.start && s.end >= t.end && s.start <= t.end:
		return Segment{s.start, t.end}, true
	case s.IsSubSegment(t):
		return Segment{s.start, s.end}, true
	case t.IsSubSegment(s):
		return Segment{t.start, t.end}, true
	}
	return Segment{}, false
}

// Intersect returns the segments where two slices of segments overlap.
func Intersect(ss, tt Segments) Segments {
	var output Segments
	newS, newT := RemoveOverlaps(ss), RemoveOverlaps(tt)
	sLen, tLen := len(newS), len(newT)
	for i, j := 0, 0; i < sLen && j < tLen; {
		if intersect, ok := SimpleIntersection(newS[i], newT[j]); ok {
			output = append(output, intersect)
		}

		if delta := newS[i].End() - newT[j].End(); delta == 0 {
			// If the two segments have the same end time, no remaining segments can
			// intersect with either segment, so advance both iterators.
			i++
			j++
		} else if delta > 0 {
			// If the segment from newS ends after the segment from newT, no other
			// remaining segment from newS can intersect with the newT segment, so
			// advance newT to the next segment.
			j++
		} else {
			// The segment from newT must end after the segment from newS, so see
			// the above comment.
			i++
		}
	}
	return output
}

// GetOverlaps returns segments from the intersection between any pair of segments.
// Note: the output segments do not overlap by design.
func GetOverlaps(ss Segments) Segments {
	var output Segments
	for i, s := range ss {
		for _, t := range ss[:i] {
			if intersect, ok := SimpleIntersection(s, t); ok {
				output = append(output, intersect)
			}
		}
	}
	return RemoveOverlaps(output)
}

// Complement takes as input a slice of segments, and a superset segment.
// It returns all segments in the superset that are not in the slice.
// More precisely, Complement(superset, ss) == tt if tt is the slice of segments
// with smallest length such that Union(ss, tt) == superset.
func Complement(superset Segment, ss Segments) Segments {
	// If the superset is not well-defined, return the empty slice.
	if !superset.IsDeltaPositive() {
		return nil
	}
	output := Segments{Segment{superset.start, superset.end}}

	for _, s := range RemoveOverlaps(ss) {
		// If the superset is a subset of s, return the nil segment.
		if superset.IsSubSegment(s) {
			return nil
		}
		// If there is no intersection between the superset and s, continue.
		if _, ok := SimpleIntersection(superset, s); !ok {
			continue
		}
		if s.start <= superset.start {
			output[len(output)-1].start = s.end
			continue
		}
		if s.end >= superset.end {
			output[len(output)-1].end = s.start
			continue
		}
		// Now, we know s must be a strict subsegment of superset.
		// We also are guaranteed that the segment to append has positive delta.
		output[len(output)-1].end = s.start
		output = append(output, Segment{s.end, superset.end})
	}
	return output
}

// SetDiff returns the difference in segments between two sets of segments.
// It is a generalization of the function Complement().
// If SetDiff(a, b Segments) == c Segments, then c is the slice of smallest length
// such that Union(c, b) == a.
func SetDiff(ss, tt Segments) Segments {
	var output Segments
	for _, s := range ss {
		output = append(output, Complement(s, tt)...)
	}
	return RemoveOverlaps(output)
}

// IsIntersectionEmpty returns whether the intersection between
// a segment and a slice of segments is empty.
func (s Segment) IsIntersectionEmpty(tt Segments) bool {
	return len(Intersect(Segments{s}, tt)) == 0
}

// IsSubSegment reports whether segment s is a sub-segment of segment t.
func (s Segment) IsSubSegment(t Segment) bool {
	return s.start >= t.start && s.end <= t.end
}

// IsPointInSegments returns true if and only if the point is contained
// in any of the segments in a slice of segments.
func IsPointInSegments(p int64, ss Segments) bool {
	for _, s := range ss {
		if IsPointInSegment(p, s) {
			return true
		}
	}
	return false
}

// IsPointInSegment returns true if and only if the point is contained in the segment.
func IsPointInSegment(p int64, s Segment) bool {
	return s.start <= p && p <= s.end
}
