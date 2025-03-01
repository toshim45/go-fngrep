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

// Function to extract lines inside brackets and print them
func extractLinesInsideBrackets(line string) []string {
	var result []string
	// Regular expression to match content inside brackets
	re := regexp.MustCompile(`\{(.*?)\}`)
	matches := re.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		// match[1] contains the content inside the brackets
		result = append(result, match[1])
	}
	return result
}

// Function to extract content inside curly brackets, including multi-line curly brackets
func extractMultiLineContentInsideCurlyBrackets(scanner *bufio.Scanner) string {
	var content strings.Builder
	inBrackets := false

	for scanner.Scan() {
		line := scanner.Text()

		// If we encounter an opening curly bracket, start collecting content
		if strings.Contains(line, "{") {
			inBrackets = true
			// Start collecting content, potentially including the opening bracket
			content.WriteString(line[strings.Index(line, "{")+1:])
		} else if inBrackets && strings.Contains(line, "}") {
			// If we encounter the closing curly bracket, stop collecting content
			// content.WriteString(line[:strings.Index(line, "}")])
			inBrackets = false
			break
		} else if inBrackets {
			// Collect content when inside curly brackets, including new lines
			content.WriteString(line + "\n")
		}
	}

	result := content.String()

	if strings.HasSuffix(result, "\n") {
		return result[:len(result)-1]
	}

	return result
}

func run(prefix string, filename string) {
	// Open a sample file or use standard input (e.g., `os.Stdin` for user input)
	file, err := os.Open(filename) // Replace with your file name or os.Stdin
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the specified prefix
		if containsPrefix(line, prefix) && strings.Contains(line, "{") {
			// Check and print lines inside brackets
			fmt.Println(line)
			contentInsideBrackets := extractMultiLineContentInsideCurlyBrackets(scanner)
			if contentInsideBrackets != "" {
				fmt.Println(contentInsideBrackets)
			}
			fmt.Println("}")
		}
	}

	// Check for errors reading the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

var opts struct {
	Prefix     string `short:"p" long:"prefix" description:"Function prefix to grep"`
	Positional struct {
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

	run(opts.Prefix, opts.Positional.Filename)
}
