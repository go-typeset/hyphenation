package hyphenation

import (
	"reflect"
	"testing"
)

// ⛔ A BREAK POSITION IS A COUNT OF LETTERS, AND HALF OF Points COUNTED BYTES.
//
// Points builds its value array over the augmented word by BYTE index —
// n := len(w), w[i:j], val[i+k] — and then reads it back by RUNE index:
// L := len([]rune(word)), and val[t+1] for a break after t letters. The two
// agree for as long as every letter is one byte, which is every English word
// and no French, German, Czech, Polish or Vietnamese one.
//
// Nothing fails loudly when they disagree: val[t+1] simply reads a different
// slot, so the word is hyphenated in the wrong place or not at all. A
// bibliography is where this lands, and a bibliography is exactly where a
// typesetting engine's readers look closely.
//
// The witness is two words with the SAME letter structure, differing only in
// whether one letter takes two bytes. The break must fall at the same letter
// count in both; it is the same word as far as Liang's algorithm is concerned.
func TestABreakPositionCountsLettersNotBytes(t *testing.T) {
	patterns := []string{"a1b"}

	mk := func() *Hyphenator {
		h := New()
		h.SetMins(1, 1) // so a four-letter word can break at all
		for _, p := range patterns {
			h.AddPattern(p)
		}
		return h
	}

	// The control: all-ASCII, one byte per letter. "xaby" breaks between a and
	// b, which is after two letters.
	ascii := mk().Points("xaby")
	if !reflect.DeepEqual(ascii, []int{2}) {
		t.Fatalf("control Points(%q) = %v, want [2] — the control is broken, so the "+
			"comparison below means nothing", "xaby", ascii)
	}

	// The same word with its first letter written as a two-byte rune. Same
	// letters in the same order as far as the pattern is concerned, so the same
	// break.
	accented := mk().Points("éaby")
	if !reflect.DeepEqual(accented, []int{2}) {
		t.Errorf("Points(%q) = %v, want [2] — the same structure as %q, which gives %v. "+
			"The value array is indexed by BYTE and read back by RUNE, so a two-byte "+
			"letter shifts every position after it.", "éaby", accented, "xaby", ascii)
	}
}

// And the same defect seen from the other side: a word that is ALL multi-byte
// letters, where the mismatch is largest.
func TestAWordOfMultiByteLettersIsNotShiftedOut(t *testing.T) {
	h := New()
	h.SetMins(1, 1)
	h.AddPattern("é1è")

	// ".ééèè." — a break between the second é and the first è is after two
	// letters. Every letter is two bytes, so a byte-indexed read looks four
	// slots along instead of two.
	got := h.Points("ééèè")
	if !reflect.DeepEqual(got, []int{2}) {
		t.Errorf("Points(%q) = %v, want [2]", "ééèè", got)
	}
}
