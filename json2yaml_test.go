package json2yaml_test

import (
	"strings"
	"testing"

	"github.com/itchyny/json2yaml"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		name string
		src  string
		want string
		err  string
	}{
		{
			name: "empty object",
			src:  "{}",
			want: `{}
`,
		},
		{
			name: "simple object",
			src:  `{"foo": 128, "bar": null, "baz": false}`,
			want: `"foo": 128
"bar": null
"baz": false
`,
		},
		{
			name: "nested object",
			src: `{
				"foo": {"bar": {"baz": 128, "bar": null}, "baz": 0},
				"bar": {"foo": {}, "bar": {"bar": {}}, "baz": {}},
				"baz": {}
			}`,
			want: `"foo":
  "bar":
    "baz": 128
    "bar": null
  "baz": 0
"bar":
  "foo": {}
  "bar":
    "bar": {}
  "baz": {}
"baz": {}
`,
		},
		{
			name: "unclosed object",
			src:  "{",
			err:  "unexpected EOF",
		},
		{
			name: "empty array",
			src:  "[]",
			want: `[]
`,
		},
		{
			name: "simple array",
			src:  `[null,false,true,-128,12345678901234567890,"foo bar baz"]`,
			want: `- null
- false
- true
- -128
- 12345678901234567890
- "foo bar baz"
`,
		},
		{
			name: "nested array",
			src:  "[0,[1],[2,3],[4,[5,[6,[],7],[]],[8]],[],9]",
			want: `- 0
- - 1
- - 2
  - 3
- - 4
  - - 5
    - - 6
      - []
      - 7
    - []
  - - 8
- []
- 9
`,
		},
		{
			name: "nested object and array",
			src:  `{"foo":[0,{"bar":[],"foo":{}},[{"foo":[{"foo":[]}]}],[[[{}]]]],"bar":[{}]}`,
			want: `"foo":
  - 0
  - "bar": []
    "foo": {}
  - - "foo":
        - "foo": []
  - - - - {}
"bar":
  - {}
`,
		},
		{
			name: "unclosed array",
			src:  "[",
			err:  "unexpected EOF",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			err := json2yaml.Convert(&sb, strings.NewReader(tc.src))
			if tc.err == "" {
				if err != nil {
					t.Fatalf("should not raise an error but got: %s", err)
				}
				if got := sb.String(); got != tc.want {
					t.Fatalf("should write\n  %q\nbut got\n  %q\nwhen source is\n  %q", tc.want, got, tc.src)
				}
			} else {
				if err == nil {
					t.Fatalf("should raise an error %s but got no error", tc.err)
				}
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("should raise an error %s but got error %s", tc.err, err)
				}
			}
		})
	}
}
