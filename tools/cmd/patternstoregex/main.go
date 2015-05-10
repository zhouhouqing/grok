package main

/*
 This tool translate a set of grok patterns files into a golang regex (RE2) form
*/
import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/gemsi/grok"
)

var (
	input  = flag.String("input", "./patterns", "input patterns file path")
	output = flag.String("output", "./patterns.go", "output file name; default srcdir/patterns.go")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tpatternstoregex [flags] -input [directory]\n")
	fmt.Fprintf(os.Stderr, "\tpatternstoregex [flags[ -input files...\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttp://godoc.org/github.com/gemsi/grok/tools/cmd/patternstoregex\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	var src []byte
	src = []byte(fmt.Sprintf("// generated by patternstoregex ; DO NOT EDIT\n\npackage grok\n\n// denormalized patterns from patterns/ files\nvar defaultCapturePatterns = map[string]string{\n"))

	log.SetFlags(0)
	log.SetPrefix("patternstoregex: ")
	flag.Usage = Usage
	flag.Parse()

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	if !isDirectory(*input) {
		log.Fatalln(`input file "` + *input + ` is not available`)
	}

	g := grok.New(grok.NODEFAULTPATTERNS)
	err := g.AddPatternsFromPath(*input)
	if err != nil {
		log.Fatalf("error : %s", err)
	}
	patterns := g.Patterns()
	// To store the keys in slice in sorted order
	var keys []string
	for k := range patterns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		src = []byte(fmt.Sprintf("%s\t`%s`: `%s`,\n", string(src), k, patterns[k]))
	}
	src = []byte(fmt.Sprintf("%s}%s", string(src), "\nvar namedCapturePatterns = map[string]string{\n"))

	// NC
	g = grok.New(grok.NODEFAULTPATTERNS, grok.NAMEDCAPTURE)
	err = g.AddPatternsFromPath(*input)
	if err != nil {
		log.Fatalf("error : %s", err)
	}
	patterns = g.Patterns()
	keys = []string{}
	// To store the keys in slice in sorted order
	for k := range patterns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		src = []byte(fmt.Sprintf("%s\t`%s`: `%s`,\n", string(src), k, patterns[k]))
	}

	src = []byte(fmt.Sprintf("%s%s", string(src), "	}"))
	//
	//
	//

	// Write to file.
	outputName := *output
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
