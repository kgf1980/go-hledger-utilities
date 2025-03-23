package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Transaction struct {
	Date        *time.Time
	Description *string
	Body        string
}

func main() {
	inputFile := flag.String("f", "", "Input file (leave blank for STDIN)")
	outputFile := flag.String("o", "", "Output file (leave blank for STDOUT)")
	sourceAccount := flag.String("s", "", "Source Account")
	targetAccount := flag.String("t", "", "Target Account")
	description := flag.String("d", "", "Description to match on")
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
				transactionDescription := ""
				if len(line) > 10 {
					transactionDescription = line[11:]
				}
				currentTransaction = Transaction{Date: &date, Body: line, Description: &transactionDescription}
			} else {
				currentTransaction.Body += "\n" + line
			}
		}
	}
	if currentTransaction.Date != nil {
		transactions = append(transactions, currentTransaction)
	}

	var output strings.Builder
	for _, transaction := range transactions {
		if transaction.Description != nil && strings.Contains(strings.ToLower(*transaction.Description), strings.ToLower(*description)) {
			output.WriteString(strings.Replace(transaction.Body, *sourceAccount, *targetAccount, -1) + "\n\n")
		} else {
			output.WriteString(transaction.Body + "\n\n")
		}
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
