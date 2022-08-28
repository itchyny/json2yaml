# json2yaml
[![CI Status](https://github.com/itchyny/json2yaml/workflows/CI/badge.svg)](https://github.com/itchyny/json2yaml/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/itchyny/json2yaml)](https://goreportcard.com/report/github.com/itchyny/json2yaml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/itchyny/json2yaml/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/itchyny/json2yaml/all.svg)](https://github.com/itchyny/json2yaml/releases)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/itchyny/json2yaml)](https://pkg.go.dev/github.com/itchyny/json2yaml)

This is an implementation of JSON to YAML converter written in Go language.
This tool efficiently converts each JSON tokens in streaming fashion,
so it avoids loading the entire JSON on the memory.

## Usage as a command line tool
```bash
json2yaml file.json ...
json2yaml <file.json >output.yaml
```

You can combine with other command line tools.
```bash
gh api /orgs/github/repos | json2yaml | less
```

## Usage as a library
You can use the converter as a Go library.
[`json2yaml.Convert(io.Writer, io.Reader) error`](https://pkg.go.dev/github.com/itchyny/json2yaml#Convert) is exported.

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/itchyny/json2yaml"
)

func main() {
	input := strings.NewReader(`{"Hello": "world!"}`)
	var output strings.Builder
	if err := json2yaml.Convert(&output, input); err != nil {
		log.Fatalln(err)
	}
	fmt.Print(output.String()) // outputs Hello: world!
}
```

## Installation
### Homebrew
```sh
brew install itchyny/tap/json2yaml
```

### Build from source
```bash
go install github.com/itchyny/json2yaml/cmd/json2yaml@latest
```

## Bug Tracker
Report bug at [Issues・itchyny/json2yaml - GitHub](https://github.com/itchyny/json2yaml/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
