// Package json2yaml implements a converter from JSON to YAML.
package json2yaml

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"
)

// Convert reads JSON from r and writes YAML to w.
func Convert(w io.Writer, r io.Reader) error {
	return (&converter{w, []byte{'.'}, 0}).convert(r)
}

type converter struct {
	w      io.Writer
	stack  []byte
	indent int
}

func (c *converter) convert(r io.Reader) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				if len(c.stack) == 1 {
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
				if len(c.stack) > 1 {
					c.indent += 2
				}
				c.stack = append(c.stack, byte(delim))
				if dec.More() {
					if c.stack[len(c.stack)-2] == ':' {
						if _, err := c.w.Write([]byte("\n")); err != nil {
							return err
						}
						if err := c.writeIndent(); err != nil {
							return err
						}
					}
					if c.stack[len(c.stack)-1] == '[' {
						if _, err := c.w.Write([]byte("- ")); err != nil {
							return err
						}
					}
				} else {
					if c.stack[len(c.stack)-2] == ':' {
						if _, err := c.w.Write([]byte(" ")); err != nil {
							return err
						}
					}
					if c.stack[len(c.stack)-1] == '{' {
						if _, err := c.w.Write([]byte("{}\n")); err != nil {
							return err
						}
					} else {
						if _, err := c.w.Write([]byte("[]\n")); err != nil {
							return err
						}
					}
				}
				continue
			case '}', ']':
				c.stack = c.stack[:len(c.stack)-1]
				c.indent -= 2
			}
		} else {
			switch c.stack[len(c.stack)-1] {
			case '{':
				if err := c.writeValue(token); err != nil {
					return err
				}
				if _, err := c.w.Write([]byte(":")); err != nil {
					return err
				}
				c.stack[len(c.stack)-1] = ':'
				continue
			case ':':
				if _, err := c.w.Write([]byte(" ")); err != nil {
					return err
				}
				fallthrough
			default:
				if err := c.writeValue(token); err != nil {
					return err
				}
				if _, err := c.w.Write([]byte("\n")); err != nil {
					return err
				}
			}
		}
		if dec.More() {
			if err := c.writeIndent(); err != nil {
				return err
			}
			switch c.stack[len(c.stack)-1] {
			case ':':
				c.stack[len(c.stack)-1] = '{'
			case '[':
				if _, err := c.w.Write([]byte("- ")); err != nil {
					return err
				}
			case '.':
				if _, err := c.w.Write([]byte("---\n")); err != nil {
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
		// sequence entry, document separator, mapping key, newlines
		`|(?:-(?:--)?|\?|\n+)(?:[ \t]|$)` +
		`)` +
		// mapping value, comment, trailing white space
		`|:(?:[ \t]|$)|[ \t](?:#|\n|$)` +
		// C0 control codes - '\n', DEL
		"|[\u0000-\u0009\u000B-\u001F\u007F]" +
		// BOM, noncharacters
		"|[\uFEFF\uFDD0-\uFDEF\uFFFE\uFFFF]",
)

func (c *converter) writeValue(v any) error {
	switch v := v.(type) {
	default:
		_, err := c.w.Write([]byte("null"))
		return err
	case bool:
		if v {
			_, err := c.w.Write([]byte("true"))
			return err
		} else {
			_, err := c.w.Write([]byte("false"))
			return err
		}
	case json.Number:
		_, err := c.w.Write([]byte(v))
		return err
	case string:
		return c.writeString(v)
	}
}

func (c *converter) writeString(v string) error {
	switch {
	default:
		_, err := c.w.Write([]byte(v))
		return err
	case quoteStringPattern.MatchString(v):
		bs, _ := json.Marshal(v)
		_, err := c.w.Write(bs)
		return err
	case strings.ContainsRune(v, '\n'):
		return c.writeBlockStyleString(v)
	}
}

func (c *converter) writeBlockStyleString(v string) error {
	if c.stack[len(c.stack)-1] == '{' {
		if _, err := c.w.Write([]byte("? ")); err != nil {
			return err
		}
	}
	if _, err := c.w.Write([]byte("|")); err != nil {
		return err
	}
	if !strings.HasSuffix(v, "\n") {
		if _, err := c.w.Write([]byte("-")); err != nil {
			return err
		}
	} else if strings.HasSuffix(v, "\n\n") {
		if _, err := c.w.Write([]byte("+")); err != nil {
			return err
		}
	}
	c.indent += 2
	for s := ""; v != ""; {
		s, v, _ = strings.Cut(v, "\n")
		if _, err := c.w.Write([]byte("\n")); err != nil {
			return err
		}
		if s != "" {
			if err := c.writeIndent(); err != nil {
				return err
			}
			if _, err := c.w.Write([]byte(s)); err != nil {
				return err
			}
		}
	}
	c.indent -= 2
	if c.stack[len(c.stack)-1] == '{' {
		if _, err := c.w.Write([]byte("\n")); err != nil {
			return err
		}
		if err := c.writeIndent(); err != nil {
			return err
		}
	}
	return nil
}

func (c *converter) writeIndent() error {
	if n := c.indent; n > 0 {
		const spaces = "                                "
		for n > len(spaces) {
			if _, err := c.w.Write([]byte(spaces)); err != nil {
				return err
			}
			n -= len(spaces)
		}
		if _, err := c.w.Write([]byte(spaces)[:n]); err != nil {
			return err
		}
	}
	return nil
}
