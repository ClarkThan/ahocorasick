# Aho-Corasick

Package `ahocorasick` is a Golang implementation of the Aho-Corasick multiple pattern string matching algorithm.


## Install

```
go get github.com/ClarkThan/ahocorasick
```

## Usage

```go 
package main

import (
	"fmt"

    "github.com/ClarkThan/ahocorasick"
)

func main() {
	m := ahocorasick.NewMatcher()
    m.BuildWithPatterns([]string{"she", "he"})
    ret1 := m.Search("shers")  // ["she", "he"]
	fmt.Println(ret1)
    
    indexedHit := m.SearchIndexed("shers")  // [{0, 3}, {1, 2}]
    fmt.Println(indexedHit)
	
    m.AddPattern("her")
    m.AddPattern("hers")
    m.Build()  // as far as new pattern added, you must Build again
    fmt.Println(m.Search("shershe"))  // ["she", "he", "her", "hers", "she", "he"]
}
```