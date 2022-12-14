// json2yaml - convert JSON to YAML
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/itchyny/json2yaml"
)

const name = "json2yaml"

const version = "0.1.4"

var revision = "HEAD"

func main() {
	os.Exit(run(os.Args[1:]))
}

const (
	exitCodeOK = iota
	exitCodeErr
)

func run(args []string) (exitCode int) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fs.SetOutput(os.Stdout)
		fmt.Printf(`%[1]s - convert JSON to YAML

Version: %s (rev: %s/%s)

Synopsis:
  %% %[1]s file ...

Options:
`, name, version, revision, runtime.Version())
		fs.PrintDefaults()
	}
	var showVersion bool
	fs.BoolVar(&showVersion, "version", false, "print version")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitCodeOK
		}
		return exitCodeErr
	}
	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return exitCodeOK
	}
	if args = fs.Args(); len(args) == 0 {
		if err := convert("-"); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
			exitCode = exitCodeErr
		}
	} else {
		for i, arg := range args {
			if i > 0 {
				fmt.Fprintln(os.Stdout, "---")
			}
			if err := convert(arg); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
				exitCode = exitCodeErr
			}
		}
	}
	return
}

func convert(name string) (err error) {
	if name == "-" {
		if err := json2yaml.Convert(os.Stdout, os.Stdin); err != nil {
			return fmt.Errorf("<stdin>: %w", err)
		}
		return nil
	}
	f, err := os.Open(filepath.Clean(name))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()
	if err := json2yaml.Convert(os.Stdout, f); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}
