// Package json2yaml implements a converter from JSON to YAML.
package json2yaml

import (
	"encoding/json"
	"io"
)

// Convert reads JSON from r and writes YAML to w.
func Convert(w io.Writer, r io.Reader) error {
	dec := json.NewDecoder(r)
	stack := []byte{'.'}
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				if stack[len(stack)-1] == '.' {
					err = nil
				} else {
					err = io.ErrUnexpectedEOF
				}
			}
			return err
		}
		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				stack = append(stack, '{')
			case '}':
				if stack[len(stack)-1] == '{' {
					w.Write([]byte("{}\n"))
				}
				stack = stack[:len(stack)-1]
			case '[':
				stack = append(stack, '[')
			case ']':
				if stack[len(stack)-1] == '[' {
					w.Write([]byte("[]\n"))
				}
				stack = stack[:len(stack)-1]
			}
		}
	}
}
