regommend
=========

Recommendation engine for Go

## Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

To install regommend, simply run:

    go get github.com/muesli/regommend

To compile it from source:

    git clone git://github.com/muesli/regommend.git
    cd regommend && go build && go test -v

## Example
```go
package main

import (
	"github.com/muesli/regommend"
	"fmt"
)

func main() {
	// Accessing a new regommend table for the first time will create it.
	books := regommend.Table("books")
}
```

## Development
API docs can be found [here](http://godoc.org/github.com/muesli/regommend).
