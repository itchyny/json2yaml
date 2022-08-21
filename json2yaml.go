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
	indent := -2
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
				indent += 2
				continue
			case '}':
				if stack[len(stack)-1] == '{' {
					if stack[len(stack)-2] == ':' {
						w.Write([]byte(" "))
					}
					w.Write([]byte("{}\n"))
				}
				stack = stack[:len(stack)-1]
				indent -= 2
			case '[':
				stack = append(stack, '[')
				continue
			case ']':
				if stack[len(stack)-1] == '[' {
					w.Write([]byte("[]\n"))
				}
				stack = stack[:len(stack)-1]
			}
		} else {
			switch stack[len(stack)-1] {
			case '{':
				if stack[len(stack)-2] == ':' {
					w.Write([]byte("\n"))
					writeIndent(w, indent)
				}
				fallthrough
			case ',':
				writeValue(w, token)
				w.Write([]byte(":"))
				stack[len(stack)-1] = ':'
				continue
			case ':':
				w.Write([]byte(" "))
				writeValue(w, token)
				w.Write([]byte("\n"))
			case '[':
				stack[len(stack)-1] = '-'
				fallthrough
			case '-':
				w.Write([]byte("- "))
				writeValue(w, token)
				w.Write([]byte("\n"))
			}
		}
		if dec.More() {
			if stack[len(stack)-1] == ':' {
				writeIndent(w, indent)
				stack[len(stack)-1] = ','
			}
		}
	}
}

func writeValue(w io.Writer, v any) {
	bs, _ := json.Marshal(v)
	w.Write(bs)
}

func writeIndent(w io.Writer, count int) {
	if n := count; n > 0 {
		const spaces = "                                "
		for n > len(spaces) {
			w.Write([]byte(spaces))
			n -= len(spaces)
		}
		w.Write([]byte(spaces)[:n])
	}
}
