package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	lineOption, wordOption, byteOption := false, false, false
	flag.BoolVar(&lineOption, "l", false, "print line count")
	flag.BoolVar(&wordOption, "w", false, "print word count")
	flag.BoolVar(&byteOption, "c", false, "print byte count")
	flag.Parse()

	optionProvided := false
	flag.Visit(func(_ *flag.Flag) {
		optionProvided = true
	})
	if !optionProvided {
		lineOption, wordOption, byteOption = true, true, true
	}

	exitCode := 0
	for _, fileName := range flag.Args() {
		content, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", fileName, err)
			exitCode = 1
			continue
		}
		text := string(content)

		lineCount := strings.Count(text, "\n")
		wordCount := len(strings.Fields(text))
		byteCount := len(text)

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
		fmt.Println(strings.Join(counts, " "), fileName)
	}
	os.Exit(exitCode)
}
