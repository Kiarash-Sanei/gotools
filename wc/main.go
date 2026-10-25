package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func (counts *Counts) Add(other Counts) {
	counts.Lines += other.Lines
	counts.Words += other.Words
	counts.Bytes += other.Bytes
}

func count(r io.Reader) (Counts, error) {
	reader := bufio.NewReader(r)
	inWord := false
	var counts Counts

	for {
		character, size, readErr := reader.ReadRune()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return counts, nil
			}
			return counts, readErr
		}

		counts.Bytes += size
		if character == '\n' {
			counts.Lines++
		}
		if unicode.IsSpace(character) {
			inWord = false
		} else if !inWord {
			counts.Words++
			inWord = true
		}
	}
}

func formatCounts(counts Counts, lineOption, wordOption, byteOption bool) string {
	values := make([]string, 0, 3)
	if lineOption {
		values = append(values, fmt.Sprint(counts.Lines))
	}
	if wordOption {
		values = append(values, fmt.Sprint(counts.Words))
	}
	if byteOption {
		values = append(values, fmt.Sprint(counts.Bytes))
	}
	return strings.Join(values, " ")
}

func run(args []string, stdout, stderr io.Writer, stdin io.Reader) int {
	lineOption, wordOption, byteOption := false, false, false
	flags := flag.NewFlagSet("wc", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.BoolVar(&lineOption, "l", false, "print line count")
	flags.BoolVar(&wordOption, "w", false, "print word count")
	flags.BoolVar(&byteOption, "c", false, "print byte count")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	optionProvided := false
	flags.Visit(func(_ *flag.Flag) {
		optionProvided = true
	})
	if !optionProvided {
		lineOption, wordOption, byteOption = true, true, true
	}

	fileNames := flags.Args()
	if len(fileNames) == 0 {
		counts, err := count(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "cannot read stdin: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, formatCounts(counts, lineOption, wordOption, byteOption))
		return 0
	}

	exitCode := 0
	var total Counts
	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(stderr, "cannot read %s: %v\n", fileName, err)
			exitCode = 1
			continue
		}
		counts, countErr := count(file)
		closeErr := file.Close()
		if countErr != nil {
			fmt.Fprintf(stderr, "cannot read %s: %v\n", fileName, countErr)
			exitCode = 1
			continue
		}
		if closeErr != nil {
			fmt.Fprintf(stderr, "cannot close %s: %v\n", fileName, closeErr)
			exitCode = 1
			continue
		}

		fmt.Fprintln(stdout, formatCounts(counts, lineOption, wordOption, byteOption), fileName)
		total.Add(counts)
	}
	if len(fileNames) > 1 {
		fmt.Fprintln(stdout, formatCounts(total, lineOption, wordOption, byteOption), "total")
	}
	return exitCode
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, os.Stdin))
}
