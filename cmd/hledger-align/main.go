package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var amountRegex = regexp.MustCompile(`\s{2,}(\S+)$`)

func main() {
	// Command-line arguments
	inputFile := flag.String("f", "", "Input file (leave empty for STDIN)")
	outputFile := flag.String("o", "", "Output file (leave empty for STDOUT)")
	flag.Parse()

	// Read input
	var lines []string
	if *inputFile != "" {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		lines = strings.Split(string(data), "\n")
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from STDIN: %v\n", err)
			os.Exit(1)
		}
	}

	// Find the longest description
	maxDescLength := 0
	for _, line := range lines {
		if strings.HasPrefix(line, " ") { // Indented lines are transaction lines
			matches := amountRegex.FindStringSubmatch(line)
			if matches != nil {
				descLength := len(line) - len(matches[0])
				if descLength > maxDescLength {
					maxDescLength = descLength
				}
			}
		}
	}

	// Align amounts
	var output strings.Builder
	for _, line := range lines {
		if strings.HasPrefix(line, " ") {
			matches := amountRegex.FindStringSubmatch(line)
			if matches != nil {
				amount := matches[1]
				desc := line[:len(line)-len(matches[0])]
				output.WriteString(fmt.Sprintf("%-*s  %s\n", maxDescLength, desc, amount))
				continue
			}
		}
		output.WriteString(line + "\n")
	}

	// Write to file or STDOUT
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output.String()), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(output.String())
	}
}
