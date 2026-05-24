# Aho-Corasick

[![Go Reference](https://pkg.go.dev/badge/github.com/ClarkThan/ahocorasick.svg)](https://pkg.go.dev/github.com/ClarkThan/ahocorasick)
[![Go Report Card](https://goreportcard.com/badge/github.com/ClarkThan/ahocorasick)](https://goreportcard.com/report/github.com/ClarkThan/ahocorasick)

Package `ahocorasick` provides a pure Go implementation of the [Aho-Corasick](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm) multiple pattern string matching algorithm.

**Efficiently find all occurrences of many patterns in a text in O(n + m + z) time** — where n is the text length, m is the total pattern length, and z is the number of matches. Build once, search many times.

## Features

- **Multiple pattern search** — find all dictionary words in a text in a single pass
- **Unicode support** — works with any UTF-8 text, including CJK characters
- **No external dependencies** — pure Go standard library only
- **Buffer-reuse API** — `SearchAppend` / `SearchIndexedAppend` for zero-allocation hot loops

## Install

```bash
go get github.com/ClarkThan/ahocorasick
```

**Go version**: 1.18+

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ClarkThan/ahocorasick"
)

func main() {
    m := ahocorasick.NewMatcher()
    m.BuildWithPatterns([]string{"she", "he", "hers"})

    // Search returns matched pattern strings
    fmt.Println(m.Search("shers"))
    // Output: [she he hers]

    // SearchIndexed returns start positions and lengths
    fmt.Println(m.SearchIndexed("shers"))
    // Output: [{0 3} {1 2} {3 4}]

    // Match checks whether any pattern exists
    fmt.Println(m.Match("shers"))  // true
    fmt.Println(m.Match("foo"))    // false
}
```

## API

### Lifecycle

```go
m := ahocorasick.NewMatcher()        // create
m.AddPattern("foo")                   // add patterns one by one
m.AddPattern("bar")
m.Build()                             // build trie and fail pointers

m.BuildWithPatterns([]string{...})    // or add + build in one call

m.Search(text)                        // search (must Build first)
m.SearchIndexed(text)
m.Match(text)
```

### Methods

| Method | Returns | Description |
|---|---|---|
| `NewMatcher()` | `*Matcher` | Create a new matcher |
| `AddPattern(pattern)` | — | Add a pattern; sets `ready = false` |
| `Build()` | — | Build fail pointers after adding patterns |
| `BuildWithPatterns(patterns)` | — | Convenience: `AddPattern` + `Build` |
| `Search(s)` | `[]string` | Return matched pattern strings |
| `SearchIndexed(s)` | `[]Hit` | Return `Hit{Start, Len}` for each match |
| `Match(s)` | `bool` | Return true if any pattern matches |
| `SearchAppend(s, buf)` | `[]string` | Append matches to caller-owned buffer |
| `SearchIndexedAppend(s, buf)` | `[]Hit` | Append hits to caller-owned buffer |

### Buffer Reuse (Hot Loop Optimization)

In tight loops, reuse a pre-allocated buffer to eliminate allocations:

```go
m.BuildWithPatterns([]string{"foo", "bar", "baz"})

// pre-allocate once
buf := make([]string, 0, 64)

texts := []string{"foo bar", "bar baz", "baz foo"}
for _, t := range texts {
    buf = m.SearchAppend(t, buf[:0])  // reuse capacity, zero allocs
    fmt.Println(buf)
}
```

Same pattern applies to `SearchIndexedAppend` with `[]Hit`.

## Performance

All benchmarks on Apple M1 Pro with 10 patterns (~120 char text).

| Scenario | Time (ns/op) | Memory (B/op) | Allocs/op |
|---|---|---|---|
| `Search` single call | 3689 | 1968 | 31 |
| `SearchAppend` with pre-allocated buffer | 3475 (−6%) | 960 (−51%) | 25 (−19%) |
| `Search` 10 calls, fresh alloc each | 9354 | 4240 | 97 |
| `SearchAppend`, buffer reused across 10 calls | 8401 (−10%) | 2224 (−48%) | 61 (−37%) |
| `SearchIndexed` single call | 3077 | 1776 | 7 |
| `SearchIndexedAppend` with pre-allocated buffer | 2527 (−18%) | 768 (−57%) | 1 (−86%) |
| `SearchIndexed` 10 calls, fresh alloc each | 7646 | 3824 | 46 |
| `SearchIndexedAppend`, buffer reused across 10 calls | 6658 (−13%) | 1808 (−53%) | 10 (−78%) |

> Use `SearchAppend` / `SearchIndexedAppend` with a reused buffer for the best performance in loop scenarios.

## Why Aho-Corasick?

Traditional approaches to multi-pattern search have trade-offs:

- **Naïve loop**: O(k·n) — scan text once per pattern, slow with many patterns
- **Regex**: depends on engine, often backtracking, unpredictable
- **Aho-Corasick**: O(n + m + z) — scans text **once**, total time grows linearly with input size regardless of pattern count

Choose this library when you have a fixed set of patterns and need to search them against many texts.
