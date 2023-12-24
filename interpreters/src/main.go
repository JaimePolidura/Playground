package main

import (
	"bufio"
	"fmt"
	"interpreters/src/failures"
	"interpreters/src/lex"
	"interpreters/src/utils"
	"io"
	"os"
)

// usage: lox <filename>
func main() {
	args := os.Args

	if len(args) >= 2 {
		runFile(os.Args[1])
	} else if len(args) == 1 {
		showInteractivePrompt()
	} else {
		fmt.Errorf("Wrong usage: lox <filename>")
		os.Exit(-1)
	}
}

func runFile(fileName string) {
	file, err := os.Open(fileName)
	utils.Check(err)
	bytes, err := io.ReadAll(file)
	utils.Check(err)

	runCode(string(bytes))
}

func showInteractivePrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print("> ")
		text := scanner.Text()
		if err := runCode(text); err != nil {
			failures.ReportFailure(err)
		}
	}
}

func runCode(code string) error {
	lexer := lex.CreateLexer(code)

	for lexer.HasNext() {
		token, err := lexer.Next()
		if err != nil {
			return err
		}
	}
}
