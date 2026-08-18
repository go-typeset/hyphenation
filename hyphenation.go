// Copyright (c) the go-typeset/hyphenation authors.
// SPDX-License-Identifier: BSD-3-Clause

package hyphenation

import "strings"

// Package hyphenation answers one question: where may a hyphen fall in a word?
//
// It is Liang's algorithm (1983). A pattern carries priority digits BETWEEN
// letters; for a word, the highest value at each inter-letter position decides
// whether a break is allowed there (odd = allowed), subject to the minimum number
// of letters required before the first hyphen and after the last. The pattern
// format is the one used by the freely available files covering some seventy
// languages.

// Hyphenator holds the loaded patterns and the min-affix limits.
type Hyphenator struct {
	pat  map[string][]int
	lmin int // minimum letters before the first hyphen
	rmin int // minimum letters after the last hyphen
}

func New() *Hyphenator {
	return &Hyphenator{pat: map[string][]int{}, lmin: 2, rmin: 3}
}

// addPattern parses one Liang pattern (e.g. "a1bc3d" or ".ach4") into its letter
// key and inter-letter value array (length = letters+1).
func (h *Hyphenator) AddPattern(p string) {
	var letters strings.Builder
	var vals []int
	pending := 0
	haveDigit := false
	for _, r := range p {
		if r >= '0' && r <= '9' {
			pending = int(r - '0')
			haveDigit = true
			continue
		}
		vals = append(vals, pending)
		letters.WriteRune(r)
		pending, haveDigit = 0, false
	}
	vals = append(vals, pending)
	_ = haveDigit
	h.pat[letters.String()] = vals
}

// points returns the break positions in word: each value t means a hyphen is
// allowed after the first t letters (so between word[t-1] and word[t]).
func (h *Hyphenator) Points(word string) []int {
	w := "." + strings.ToLower(word) + "."
	n := len(w)
	val := make([]int, n+1)
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			v, ok := h.pat[w[i:j]]
			if !ok {
				continue
			}
			for k := 0; k < len(v); k++ {
				if i+k < len(val) && v[k] > val[i+k] {
					val[i+k] = v[k]
				}
			}
		}
	}
	// The augmented word has a leading '.', so original letter t sits at w[t+1];
	// the break after t letters uses val[t+1].
	L := len([]rune(word))
	var pts []int
	for t := h.lmin; t <= L-h.rmin; t++ {
		if val[t+1]%2 == 1 {
			pts = append(pts, t)
		}
	}
	return pts
}

// SetMins sets how many letters must remain before the first hyphen and after
// the last. The defaults are 2 and 3 — the long-standing English convention, and
// the reason a five-letter word is never hyphenated; other languages, and other
// house styles, want otherwise.
func (h *Hyphenator) SetMins(left, right int) {
	if left < 1 {
		left = 1
	}
	if right < 1 {
		right = 1
	}
	h.lmin, h.rmin = left, right
}

// Mins reports the current minimum number of letters before the first hyphen and
// after the last.
func (h *Hyphenator) Mins() (left, right int) { return h.lmin, h.rmin }
