package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	f1Ptr := flag.String("p", "", "Path to the project file")
	f2Ptr := flag.String("f", "", "Path to the input file")

	flag.Parse()

	if *f1Ptr == "" {
		fmt.Println("Usage: go run main.go -p <project file> [-f <input file>]")
		os.Exit(1)
	}

	f1Exists, f2Exists, stdinExists := checkExistence(*f1Ptr, *f2Ptr)

	if !f1Exists && !f2Exists && !stdinExists {

		f1, err := os.Create(*f1Ptr)
		if err != nil {
			fmt.Println("Error creating file:", err)
			os.Exit(1)
		}
		f1.Close()
		fmt.Printf("Created new file: %s\n", *f1Ptr)
		return
	}

	var f2Scanner *bufio.Scanner
	if f2Exists {
		f2, err := os.Open(*f2Ptr)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer f2.Close()
		f2Scanner = bufio.NewScanner(f2)
	}

	f1, err := os.OpenFile(*f1Ptr, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f1.Close()

	uniqueLines := make(map[string]struct{})
	duplicateCount := 0

	scanner := bufio.NewScanner(f1)
	for scanner.Scan() {
		line := scanner.Text()
		uniqueLines[line] = struct{}{}
	}

	if f2Scanner == nil && stdinExists {
		stdinScanner := bufio.NewScanner(os.Stdin)
		for stdinScanner.Scan() {
			line := stdinScanner.Text()
			_, exists := uniqueLines[line]
			if !exists {
				fmt.Fprintln(f1, line)
				uniqueLines[line] = struct{}{}
			} else {
				duplicateCount++
			}
		}
	}

	if f2Scanner != nil {
		for f2Scanner.Scan() {
			line := f2Scanner.Text()
			_, exists := uniqueLines[line]
			if !exists {
				fmt.Fprintln(f1, line)
				uniqueLines[line] = struct{}{}
			} else {
				duplicateCount++
			}
		}
	}

	f1.Truncate(0)
	f1.Seek(0, 0)

	// Rewrite the content of f1 with unique lines
	writer := bufio.NewWriter(f1)
	for line := range uniqueLines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	f1TotalLineCount := len(uniqueLines)

	fmt.Printf("Number of duplicate lines removed from the input file: %d\n", duplicateCount)
	fmt.Printf("Total line count of the project file: %d\n", f1TotalLineCount)
}

func checkExistence(f1Path, f2Path string) (f1Exists, f2Exists, stdinExists bool) {
	_, f1Err := os.Stat(f1Path)
	_, f2Err := os.Stat(f2Path)

	stdinInfo, _ := os.Stdin.Stat()
	stdinExists = (stdinInfo.Mode() & os.ModeCharDevice) == 0

	return !os.IsNotExist(f1Err), !os.IsNotExist(f2Err), stdinExists
}
