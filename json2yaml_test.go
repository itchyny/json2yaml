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
			name: "null",
			src:  "null",
			want: join([]string{"null"}),
		},
		{
			name: "boolean",
			src:  "false true",
			want: join([]string{"false", "true"}),
		},
		{
			name: "number",
			src:  "0 128 -320 3.14 -6.63e-34",
			want: join([]string{"0", "128", "-320", "3.14", "-6.63e-34"}),
		},
		{
			name: "string",
			src:  `"" "foo" "null" "hello, world" "\b\f\n\r\n\t" "１２３４５"`,
			want: join([]string{`""`, `foo`, `"null"`, `hello, world`, `"\u0008\u000c\n\r\n\t"`, `１２３４５`}),
		},
		{
			name: "quote booleans",
			src:  `"true" "False" "YES" "y" "no" "n" "oN" "Off" "truer" "oon" "f"`,
			want: join([]string{`"true"`, `"False"`, `"YES"`, `"y"`, `"no"`, `"n"`, `"oN"`, `"Off"`, `truer`, `oon`, `f`}),
		},
		{
			name: "quote integers",
			src: `"0" "+42" "128" "900" "-1_234_567_890" "+ 1" "11:22" "+1:2" "-3:4" "0:1:02:1:0" "12:50" "12:60"
				"0b1" "0b11_00" "0b" "0b2" "0664" "0_1_2_3" "0_" "0678" "0o1_0" "0O0" "0o"
				"0x0" "0x09af" "0xFE_FF" "0x" "0xfg" "0x_F_"`,
			want: join([]string{
				`"0"`, `"+42"`, `"128"`, `"900"`, `"-1_234_567_890"`, `+ 1`, `"11:22"`, `"+1:2"`, `"-3:4"`, `"0:1:02:1:0"`, `"12:50"`, `12:60`,
				`"0b1"`, `"0b11_00"`, `0b`, `0b2`, `"0664"`, `"0_1_2_3"`, `"0_"`, `0678`, `"0o1_0"`, `"0O0"`, `0o`,
				`"0x0"`, `"0x09af"`, `"0xFE_FF"`, `0x`, `0xfg`, `"0x_F_"`,
			}),
		},
		{
			name: "quote floating point numbers",
			src: `"0.1" "3.14156" "-42.195" "-.3" "+6." "-+1" "1E+9" "6.63e-34" "1e2"
				"1_2.3_4e56" "120:30:40.56" ".inf" "+.inf" "-.inf" ".infr" ".nan" "+.nan" "-.nan" ".nan."`,
			want: join([]string{
				`"0.1"`, `"3.14156"`, `"-42.195"`, `"-.3"`, `"+6."`, `-+1`, `"1E+9"`, `"6.63e-34"`, `"1e2"`,
				`"1_2.3_4e56"`, `"120:30:40.56"`, `".inf"`, `"+.inf"`, `"-.inf"`, `.infr`, `".nan"`, `+.nan`, `-.nan`, `.nan.`,
			}),
		},
		{
			name: "quote date time",
			src: `"2022-08-04" "1000-1-1" "9999-12-31" "1999-99-99" "999-9-9" "2000-08" "2000-08-" "2000-"
				"2022-01-01T12:13:14" "2022-02-02 12:13:14.567" "2022-03-03   1:2:3" "2022-03-04 15:16:17." "2022-03-04 15:16:"
				"2000-12-31T01:02:03-09:00" "2000-12-31t01:02:03Z" "2000-12-31 01:02:03 +7" "2222-22-22  22:22:22  +22:22"`,
			want: join([]string{`"2022-08-04"`, `"1000-1-1"`, `"9999-12-31"`, `"1999-99-99"`, `999-9-9`, `2000-08`, `2000-08-`, `2000-`,
				`"2022-01-01T12:13:14"`, `"2022-02-02 12:13:14.567"`, `"2022-03-03   1:2:3"`, `"2022-03-04 15:16:17."`, `"2022-03-04 15:16:"`,
				`"2000-12-31T01:02:03-09:00"`, `"2000-12-31t01:02:03Z"`, `"2000-12-31 01:02:03 +7"`, `"2222-22-22  22:22:22  +22:22"`}),
		},
		{
			name: "quote indicators",
			src: `"!" "\"" "#" "$" "%" "&" "'" "(" ")" "*" "+" ","
				"-" "--" "---" "----" "--- -" "- ---" "- --- -" "-- --" "?-" "-?" "?---" "---?" "--- ?"
				"." "/" ":" ";" "<" "=" ">" "?" "[" "\\" "]" "^" "_" "{" "|" "}" "~"
				"%TAG" "!!str" "!<>" "&anchor" "*anchor" "https://example.com/?q=text#fragment"
				"- ." ". -" "-." ".-" "? ." ". ?" "?." ".?" ": ." ". :" "?:" ":?" ". ? :." "[]" "{}"
				". #" "# ." ".#." ". #." ".# ." ". # ." ". ! \" $ % & ' ( ) * + , - / ; < = > ? [ \\ ] ^ _ { | } ~"`,
			want: join([]string{`"!"`, `"\""`, `"#"`, `$`, `"%"`, `"\u0026"`, `"'"`, `(`, `)`, `"*"`, `+`, `","`,
				`"-"`, `--`, `"---"`, `----`, `"--- -"`, `"- ---"`, `"- --- -"`, `-- --`, `?-`, `-?`, `?---`, `---?`, `"--- ?"`,
				`.`, `/`, `":"`, `;`, `<`, `=`, `"\u003e"`, `"?"`, `"["`, `\`, `"]"`, `^`, `_`, `"{"`, `"|"`, `"}"`, `"~"`,
				`"%TAG"`, `"!!str"`, `"!\u003c\u003e"`, `"\u0026anchor"`, `"*anchor"`, `https://example.com/?q=text#fragment`,
				`"- ."`, `. -`, `-.`, `.-`, `"? ."`, `. ?`, `?.`, `.?`, `": ."`, `". :"`, `"?:"`, `:?`, `. ? :.`, `"[]"`, `"{}"`,
				`". #"`, `"# ."`, `.#.`, `". #."`, `.# .`, `". # ."`, `. ! " $ % & ' ( ) * + , - / ; < = > ? [ \ ] ^ _ { | } ~`}),
		},
		{
			name: "quote white spaces",
			src:  `" " "\t" " ." ". " "\t." ".\t" ". ." ".\t."`,
			want: join([]string{`" "`, `"\t"`, `" ."`, `". "`, `"\t."`, `".\t"`, `. .`, `".\t."`}),
		},
		{
			name: "quote special characters",
			src:  "\"\\n\" \"\x7F\" \"\uFDCF\" \"\uFDD0\" \"\uFDEF\" \"\uFEFE\" \"\uFEFF\" \"\uFFFD\" \"\uFFFE\" \"\uFFFF\"",
			want: join([]string{"\"\\n\"", "\"\x7F\"", "\uFDCF", "\"\uFDD0\"", "\"\uFDEF\"", "\uFEFE", "\"\uFEFF\"", "\uFFFD", "\"\uFFFE\"", "\"\uFFFF\""}),
		},
		{
			name: "empty object",
			src:  "{}",
			want: `{}
`,
		},
		{
			name: "simple object",
			src:  `{"foo": 128, "bar": null, "baz": false}`,
			want: `foo: 128
bar: null
baz: false
`,
		},
		{
			name: "nested object",
			src: `{
				"foo": {"bar": {"baz": 128, "bar": null}, "baz": 0},
				"bar": {"foo": {}, "bar": {"bar": {}}, "baz": {}},
				"baz": {}
			}`,
			want: `foo:
  bar:
    baz: 128
    bar: null
  baz: 0
bar:
  foo: {}
  bar:
    bar: {}
  baz: {}
baz: {}
`,
		},
		{
			name: "multiple objects",
			src:  `{}{"foo":128}{}`,
			want: join([]string{"{}", "foo: 128", "{}"}),
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
- foo bar baz
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
			want: `foo:
  - 0
  - bar: []
    foo: {}
  - - foo:
        - foo: []
  - - - - {}
bar:
  - {}
`,
		},
		{
			name: "multiple arrays",
			src:  `[][{"foo":128}][]`,
			want: join([]string{"[]", "- foo: 128", "[]"}),
		},
		{
			name: "deeply nested object",
			src:  `{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{"x":{}}}}}}}}}}}}}}}}}}}}}`,
			want: `x:
  x:
    x:
      x:
        x:
          x:
            x:
              x:
                x:
                  x:
                    x:
                      x:
                        x:
                          x:
                            x:
                              x:
                                x:
                                  x:
                                    x:
                                      x: {}
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
				if got, want := diff(sb.String(), tc.want); got != want {
					t.Fatalf("should write\n  %q\nbut got\n  %q\nwhen source is\n  %q", want, got, tc.src)
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

func join(xs []string) string {
	var sb strings.Builder
	n := 5*(len(xs)-1) + 1
	for _, x := range xs {
		n += len(x)
	}
	sb.Grow(n)
	for i, x := range xs {
		if i > 0 {
			sb.WriteString("---\n")
		}
		sb.WriteString(x)
		sb.WriteString("\n")
	}
	return sb.String()
}

func diff(xs, ys string) (string, string) {
	if xs == ys {
		return "", ""
	}
	for {
		i := strings.IndexByte(xs, '\n')
		j := strings.IndexByte(ys, '\n')
		if i < 0 || j < 0 || xs[:i] != ys[:j] {
			break
		}
		xs, ys = xs[i+1:], ys[j+1:]
	}
	for {
		i := strings.LastIndexByte(xs, '\n')
		j := strings.LastIndexByte(ys, '\n')
		if i < 0 || j < 0 || xs[i:] != ys[j:] {
			break
		}
		xs, ys = xs[:i], ys[:j]
	}
	return xs, ys
}
