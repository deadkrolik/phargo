# PHARGO

PHAR-files reader and parser written in golang

## Installation

1. Download and install:

```sh
$ go get -u github.com/deadkrolik/phargo
```

2. Import and use it:

```go
package main

import (
    "log"

    "github.com/deadkrolik/phargo"
)

func main() {
    r := phargo.NewReader()
    file, err := r.Parse("file.phar")
    log.Println(file, err)
}
```

## Running the tests

Just run the command:

```sh
$ go test
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
