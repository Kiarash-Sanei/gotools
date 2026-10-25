package wc

import (
	"bufio"
	"fmt"
	"gotools/common"
	"io"
	"os"
	"strings"
)

func main() {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tokens, err := common.ParseTokens(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	lineOption, wordOption, byteOption := false, false, false
	var flags []string
	var fileNames []string
	for _, token := range tokens {
		if strings.HasPrefix(token, "-") {
			flags = append(flags, token)
		} else {
			fileNames = append(fileNames, token)
		}
	}

	for _, flag := range flags {
		switch flag {
		case "-l":
			lineOption = true
		case "-w":
			wordOption = true
		case "-c":
			byteOption = true
		default:
			fmt.Fprintf(os.Stderr, "unknown option: %s\n", flag)
			os.Exit(2)
		}
	}
	if len(flags) == 0 {
		lineOption, wordOption, byteOption = true, true, true
	}

	exitCode := 0
	for _, fileName := range fileNames {
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
