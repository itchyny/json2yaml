// Package json2yaml implements a converter from JSON to YAML.
package json2yaml

import (
	"encoding/json"
	"io"
	"regexp"
)

// Convert reads JSON from r and writes YAML to w.
func Convert(w io.Writer, r io.Reader) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
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
			case '{', '[':
				stack = append(stack, byte(delim))
				indent += 2
				if dec.More() {
					if stack[len(stack)-2] == ':' {
						w.Write([]byte("\n"))
						writeIndent(w, indent)
					}
					if stack[len(stack)-1] == '[' {
						w.Write([]byte("- "))
					}
				} else {
					if stack[len(stack)-2] == ':' {
						w.Write([]byte(" "))
					}
					if stack[len(stack)-1] == '{' {
						w.Write([]byte("{}\n"))
					} else {
						w.Write([]byte("[]\n"))
					}
				}
				continue
			case '}', ']':
				stack = stack[:len(stack)-1]
				indent -= 2
			}
		} else {
			switch stack[len(stack)-1] {
			case '{':
				writeValue(w, token)
				w.Write([]byte(":"))
				stack[len(stack)-1] = ':'
				continue
			case ':':
				w.Write([]byte(" "))
				fallthrough
			default:
				writeValue(w, token)
				w.Write([]byte("\n"))
			}
		}
		if dec.More() {
			writeIndent(w, indent)
			switch stack[len(stack)-1] {
			case ':':
				stack[len(stack)-1] = '{'
			case '[':
				w.Write([]byte("- "))
			case '.':
				w.Write([]byte("---\n"))
			}
		}
	}
}

// This pattern matches more than the specifications,
// but it is okay to quote for parsers just in case.
var quoteStringPattern = regexp.MustCompile(
	`^(?:` +
		`(?i:` +
		// tag:yaml.org,2002:null
		`|~|null` +
		// tag:yaml.org,2002:bool
		`|true|false|y(?:es)?|no?|on|off` +
		// tag:yaml.org,2002:int, tag:yaml.org,2002:float
		`|[-+]?(?:0(?:b[01_]+|o?[0-7_]+|x[0-9a-f_]+)` + // base 2, 8, 16
		`|(?:(?:0|[1-9][0-9_]*)(?::[0-5]?[0-9])*(?:\.[0-9_]*)?` +
		`|\.[0-9_]+)(?:E[-+]?[0-9]+)?` + // base 10, 60
		`|\.inf)|\.nan` + // infinities, not-a-number
		// tag:yaml.org,2002:timestamp
		`|\d\d\d\d-\d\d?-\d\d?` + // date
		`(?:(?:T|\s+)\d\d?:\d\d?:\d\d?(?:\.\d*)?` + // time
		`(?:\s*(?:Z|[-+]\d\d?(?::\d\d?)?))?)?` + // time zone
		`)$` +
		// c-indicator - '-' - '?' - ':', leading white space
		"|[,\\[\\]{}#&*!|>'\"%@` \\t]" +
		// sequence entry, document separator, mapping key
		`|(?:-(?:--)?|\?)(?:\s|$)` +
		`)` +
		// mapping value, comment, trailing white space
		`|:(?:\s|$)|\s(?:#|$)` +
		// C0 control codes, DEL
		"|[\u0000-\u001F\u007F]" +
		// BOM, noncharacters
		"|[\uFEFF\uFDD0-\uFDEF\uFFFE\uFFFF]",
)

func writeValue(w io.Writer, v any) {
	switch v := v.(type) {
	case nil:
		w.Write([]byte("null"))
	case bool:
		if v {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("false"))
		}
	case json.Number:
		w.Write([]byte(v))
	case string:
		if quoteStringPattern.MatchString(v) {
			bs, _ := json.Marshal(v)
			w.Write(bs)
		} else {
			w.Write([]byte(v))
		}
	}
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
