package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Function to check for lines containing a given prefix
func containsPrefix(line, prefix string) bool {
	return strings.HasPrefix(line, prefix)
}

// Function to check for lines containing a given regex pattern
func matchPattern(line string, rp *regexp.Regexp) bool {
	return rp.MatchString(line)
}

// Function to extract content inside curly brackets, including multi-line curly brackets
func extractMultiLineContentInsideCurlyBrackets(scanner *bufio.Scanner) string {
	return extractMultiLineContentInsideXBrackets(scanner, "{", "}")
}

// Function to extract content inside round brackets/parentheses, including multi-line round brackets/parentheses
func extractMultiLineContentInsideParentheses(scanner *bufio.Scanner) string {
	return extractMultiLineContentInsideXBrackets(scanner, "(", ")")
}

func extractMultiLineContentInsideXBrackets(scanner *bufio.Scanner, openBracket, closeBracket string) string {
	var content strings.Builder
	inBrackets := false

	for scanner.Scan() {
		line := scanner.Text()

		// If we encounter an opening curly bracket, start collecting content
		if strings.Contains(line, openBracket) {
			inBrackets = true
			// Start collecting content, potentially including the opening bracket
			content.WriteString(line[strings.Index(line, openBracket)+1:])
		} else if inBrackets && strings.Contains(line, closeBracket) {
			// If we encounter the closing curly bracket, stop collecting content
			content.WriteString(line)
			inBrackets = false
			break
		} else if inBrackets {
			// Collect content when inside curly brackets, including new lines
			content.WriteString(line + "\n")
		}
	}

	return content.String()
}

func run(prefix string, regexPattern string, filename string) {
	// Open a sample file or use standard input (e.g., `os.Stdin` for user input)
	file, err := os.Open(filename) // Replace with your file name or os.Stdin
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var rp *regexp.Regexp
	if regexPattern != "" {
		rp = regexp.MustCompile(regexPattern)
	}

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the specified prefix
		if prefix != "" && containsPrefix(line, prefix) && strings.Contains(line, "{") {
			// Check and print lines inside brackets
			fmt.Println(line)
			contentInsideBrackets := extractMultiLineContentInsideCurlyBrackets(scanner)
			if contentInsideBrackets != "" {
				fmt.Println(contentInsideBrackets)
			}
		} else if regexPattern != "" && matchPattern(line, rp) {
			fmt.Println(line)
			contentInsideBracket := extractMultiLineContentInsideParentheses(scanner)
			if contentInsideBracket != "" {
				fmt.Println(contentInsideBracket)
			}
		}
	}

	// Check for errors reading the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

var opts struct {
	Prefix       string `short:"p" long:"prefix" description:"Function prefix to grep"`
	RegexPattern string `short:"r" long:"regexp" description:"Function line regex to grep"`
	Positional   struct {
		CmdName  string
		Filename string
	} `positional-args:"yes" required:"yes"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		if err != nil {
			// Check the specific error type
			if flagsErr, ok := err.(*flags.Error); ok {
				// If it's a help request, print the help message and exit gracefully
				if flagsErr.Type == flags.ErrHelp {
					fmt.Println(err)
					return
				}
			}
			// For other errors, print the error message and exit with a non-zero status
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	run(opts.Prefix, opts.RegexPattern, opts.Positional.Filename)
}
