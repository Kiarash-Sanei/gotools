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

func count(r io.Reader) (lineCount, wordCount, byteCount int, err error) {
	reader := bufio.NewReader(r)
	inWord := false

	for {
		character, size, readErr := reader.ReadRune()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return lineCount, wordCount, byteCount, nil
			}
			return lineCount, wordCount, byteCount, readErr
		}

		byteCount += size
		if character == '\n' {
			lineCount++
		}
		if unicode.IsSpace(character) {
			inWord = false
		} else if !inWord {
			wordCount++
			inWord = true
		}
	}
}

func run(args []string, stdout, stderr io.Writer) int {
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

	exitCode := 0
	for _, fileName := range flags.Args() {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(stderr, "cannot read %s: %v\n", fileName, err)
			exitCode = 1
			continue
		}
		lineCount, wordCount, byteCount, countErr := count(file)
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

		counts := make([]string, 0, 3)
		if lineOption {
			counts = append(counts, fmt.Sprint(lineCount))
		}
		if wordOption {
			counts = append(counts, fmt.Sprint(wordCount))
		}
		if byteOption {
			counts = append(counts, fmt.Sprint(byteCount))
		}
		fmt.Fprintln(stdout, strings.Join(counts, " "), fileName)
	}
	return exitCode
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
