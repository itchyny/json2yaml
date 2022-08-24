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
						if _, err := w.Write([]byte("\n")); err != nil {
							return err
						}
						if err := writeIndent(w, indent); err != nil {
							return err
						}
					}
					if stack[len(stack)-1] == '[' {
						if _, err := w.Write([]byte("- ")); err != nil {
							return err
						}
					}
				} else {
					if stack[len(stack)-2] == ':' {
						if _, err := w.Write([]byte(" ")); err != nil {
							return err
						}
					}
					if stack[len(stack)-1] == '{' {
						if _, err := w.Write([]byte("{}\n")); err != nil {
							return err
						}
					} else {
						if _, err := w.Write([]byte("[]\n")); err != nil {
							return err
						}
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
				if err := writeValue(w, token); err != nil {
					return err
				}
				if _, err := w.Write([]byte(":")); err != nil {
					return err
				}
				stack[len(stack)-1] = ':'
				continue
			case ':':
				if _, err := w.Write([]byte(" ")); err != nil {
					return err
				}
				fallthrough
			default:
				if err := writeValue(w, token); err != nil {
					return err
				}
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
			}
		}
		if dec.More() {
			if err := writeIndent(w, indent); err != nil {
				return err
			}
			switch stack[len(stack)-1] {
			case ':':
				stack[len(stack)-1] = '{'
			case '[':
				if _, err := w.Write([]byte("- ")); err != nil {
					return err
				}
			case '.':
				if _, err := w.Write([]byte("---\n")); err != nil {
					return err
				}
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
		`|[-+]?(?:0(?:b[01_]+|o[0-7_]+|x[0-9a-f_]+)` + // base 2, 8, 16
		`|(?:[0-9][0-9_]*(?::[0-5]?[0-9])*(?:\.[0-9_]*)?` +
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

func writeValue(w io.Writer, v any) error {
	switch v := v.(type) {
	case nil:
		if _, err := w.Write([]byte("null")); err != nil {
			return err
		}
	case bool:
		if v {
			if _, err := w.Write([]byte("true")); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte("false")); err != nil {
				return err
			}
		}
	case json.Number:
		if _, err := w.Write([]byte(v)); err != nil {
			return err
		}
	case string:
		if quoteStringPattern.MatchString(v) {
			bs, _ := json.Marshal(v)
			if _, err := w.Write(bs); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte(v)); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeIndent(w io.Writer, count int) error {
	if n := count; n > 0 {
		const spaces = "                                "
		for n > len(spaces) {
			if _, err := w.Write([]byte(spaces)); err != nil {
				return err
			}
			n -= len(spaces)
		}
		if _, err := w.Write([]byte(spaces)[:n]); err != nil {
			return err
		}
	}
	return nil
}
