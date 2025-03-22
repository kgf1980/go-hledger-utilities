package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Transaction struct {
	Date *time.Time
	Body string
}

func main() {
	inputFile := flag.String("i", "", "Input file (leave blank for STDIN)")
	outputFile := flag.String("o", "", "Output file (leave blank for STDOUT)")
	flag.Parse()

	var lines []string
	if *inputFile != "" {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
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

	var transactions []Transaction
	var currentTransaction Transaction
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			if currentTransaction.Date != nil {
				transactions = append(transactions, currentTransaction)
				currentTransaction = Transaction{}
			}
		} else {
			date, err := time.Parse("2006-01-02", line[:10])
			if err == nil {
				if currentTransaction.Date != nil {
					transactions = append(transactions, currentTransaction)
				}
				currentTransaction = Transaction{Date: &date, Body: line}
			} else {
				currentTransaction.Body += "\n" + line
			}
		}
	}
	if currentTransaction.Date != nil {
		transactions = append(transactions, currentTransaction)
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.Before(*transactions[j].Date)
	})

	var output strings.Builder
	for _, transaction := range transactions {
		output.WriteString(transaction.Body + "\n\n")
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
