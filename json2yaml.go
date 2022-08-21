// Package json2yaml implements a converter from JSON to YAML.
package json2yaml

import "io"

// Convert reads JSON from r and writes YAML to w.
func Convert(w io.Writer, r io.Reader) error {
	io.Copy(w, r)
	return nil
}
