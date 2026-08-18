package hyphenation

import (
	"reflect"
	"testing"
)

// Liang's algorithm: a pattern carries priority digits BETWEEN letters, a word is
// scanned for every pattern that matches any of its substrings (with "." marking a
// word boundary), and the MAXIMUM value at each inter-letter position decides —
// odd means a break is allowed there.
//
// The patterns below are the real TeX English ones for these words.
func TestPointsFollowsLiangsRule(t *testing.T) {
	h := New()
	for _, p := range []string{"hy3ph", "he2n", "hena4", "hen5at", "1na", "n2at", "1tio", "2io", "o2n"} {
		h.AddPattern(p)
	}
	cases := []struct {
		word string
		want []int
	}{
		// "hyphenation" → hy-phen-ation: the odd values win at 2 and 6.
		{"hyphenation", []int{2, 6}},
		// A word no pattern touches has no break point at all.
		{"xyz", nil},
	}
	for _, c := range cases {
		if got := h.Points(c.word); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Points(%q) = %v, want %v", c.word, got, c.want)
		}
	}
}

// A pattern is stored under its LETTERS, so adding the same letters twice keeps
// only the last values — a documented simplification of Liang, who keeps every
// pattern and takes the maximum. It matters only when two patterns share their
// letters exactly, which TeX pattern files do not do.
func TestSameLettersKeepTheLastValues(t *testing.T) {
	h := New()
	h.AddPattern("hy3ph")
	h.AddPattern("hy2ph")
	if got := h.Points("hyphenation"); len(got) != 0 {
		t.Errorf("the second (even) pattern did not replace the first: %v", got)
	}
}

// A "." in a pattern anchors it to a word boundary: the word is scanned as
// ".word.", so a leading-dot pattern can only match at the start.
func TestBoundaryAnchoredPatterns(t *testing.T) {
	h := New()
	h.AddPattern(".ab1cdef")
	if got := h.Points("abcdefgh"); !reflect.DeepEqual(got, []int{2}) {
		t.Errorf("Points(abcdefgh) = %v, want [2]", got)
	}
	if got := h.Points("xabcdefgh"); len(got) != 0 {
		t.Errorf("the boundary pattern matched in mid-word: %v", got)
	}
}

// Degenerate input: no patterns, an empty word, and a one-letter word all yield
// nothing rather than an out-of-range break.
func TestDegenerateInput(t *testing.T) {
	empty := New()
	for _, w := range []string{"", "a", "hyphenation"} {
		if got := empty.Points(w); len(got) != 0 {
			t.Errorf("with no patterns, Points(%q) = %v, want none", w, got)
		}
	}
	h := New()
	h.AddPattern("1a")
	for _, w := range []string{"", "a"} {
		if got := h.Points(w); len(got) != 0 {
			t.Errorf("Points(%q) = %v, want none (a break needs a letter on each side)", w, got)
		}
	}
}

// \lefthyphenmin / \righthyphenmin bound where a hyphen may fall: with TeX's
// defaults (2 and 3) a word shorter than six letters has no legal break at all,
// whatever the patterns say.
func TestHyphenMinsBoundTheBreaks(t *testing.T) {
	h := New()
	h.AddPattern("a1b")
	if got := h.Points("xaby"); len(got) != 0 {
		t.Errorf("a 4-letter word broke despite lmin=2, rmin=3: %v", got)
	}
	if got := h.Points("xxabyyy"); !reflect.DeepEqual(got, []int{3}) {
		t.Errorf("Points(xxabyyy) = %v, want [3]", got)
	}
}

// \lefthyphenmin / \righthyphenmin are settable, and a value below 1 is clamped:
// a hyphen with no letter on one side of it is not a hyphenation.
func TestSetMins(t *testing.T) {
	h := New()
	if l, r := h.Mins(); l != 2 || r != 3 {
		t.Errorf("defaults = %d/%d, want TeX's 2/3", l, r)
	}
	h.AddPattern("a1b")
	h.SetMins(1, 1)
	if got := h.Points("xaby"); !reflect.DeepEqual(got, []int{2}) {
		t.Errorf("with mins 1/1, Points(xaby) = %v, want [2]", got)
	}
	h.SetMins(0, -5)
	if l, r := h.Mins(); l != 1 || r != 1 {
		t.Errorf("mins below 1 were not clamped: %d/%d", l, r)
	}
}
