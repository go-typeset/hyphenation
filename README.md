# hyphenation

[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26.4%2B-00ADD8)](https://go.dev/dl/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-1a7f37)](#tests)

**Where may a hyphen fall in a word — pure Go, no cgo, no dependencies.**

It answers that one question, for any language you have patterns for.

## How it works

The method is Liang's (1983). A pattern carries priority digits *between*
letters: `hy3ph` says a break after `hy` scores 3. Every pattern matching any
substring of the word contributes (a `.` in a pattern anchors it to a word
boundary), the highest value at each position wins, and an **odd** value means a
break is allowed there.

Two limits then forbid breaks too close to either end — by default 2 letters
before the first hyphen and 3 after the last, which is why a five-letter word is
never hyphenated. Both are settable.

## Patterns

The library reads the pattern format used by the freely available pattern files
covering some seventy languages, distributed as `hyph-*.tex` (for example
`hyph-en-gb`, `hyph-fr`, `hyph-de-1996`). Feed it the pattern lines; nothing else
about those files is needed.

## Use

```go
import "github.com/go-typeset/hyphenation"

h := hyphenation.New()
for _, p := range patterns {   // "hy3ph", "he2n", "hena4", …
    h.AddPattern(p)
}

h.Points("hyphenation")   // → [2 6]: hy-phen-ation
h.SetMins(2, 2)           // fewer letters required after the last hyphen
```

Each value `t` means a hyphen may follow the first `t` letters.

Pair it with [linebreak](https://github.com/go-typeset/linebreak) to turn those
positions into candidate breakpoints.

## Tests

`go test ./...` — 100% statement coverage, run on six 64-bit architectures
(amd64, arm64, riscv64, loong64, ppc64le, s390x), three operating systems, and
both wasm targets.

## Licence

BSD-3-Clause.
