package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ftfmtavares/online-companies-finder/csv"
	"github.com/ftfmtavares/online-companies-finder/query"
)

func main() {
	inputFile := flag.String("input", "", "Path to input CSV file")
	inputColumn := flag.String("input-column", "", "Name of the input column to process")
	outputFile := flag.String("output", "", "Path to output CSV file")
	outputColumn := flag.String("output-column", "", "Name of the output column to write results to")
	additionalSearch := flag.String("additional-search-params", "", "Additional search parameters")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" || *inputColumn == "" || *outputColumn == "" {
		fmt.Println("Usage: go run main.go -input <input.csv> -input-column <column_name> -output <output.csv> -output-column <column_name>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	records, err := csv.ReadCSV(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input CSV: %v\n", err)
		os.Exit(1)
	}

	inputColumnIndex := -1
	outputColumnIndex := -1
	for i := range records[0] {
		if records[0][i] == *inputColumn {
			inputColumnIndex = i
		}

		if records[0][i] == *outputColumn {
			outputColumnIndex = i
		}
	}

	if inputColumnIndex == -1 {
		fmt.Printf("Input column '%s' not found in CSV header\n", *inputColumn)
		os.Exit(1)
	}

	if outputColumnIndex == -1 {
		records[0] = append(records[0], *outputColumn)
	}

	for i := 1; i < len(records); i++ {
		fmt.Printf("searching for %s\n", records[i][inputColumnIndex])

		result, err := query.DuckDuckGoFirstResult(*additionalSearch + " " + records[i][inputColumnIndex])
		if err != nil {
			fmt.Printf("Error searching for %s: %v\n", records[i][inputColumnIndex], err)
		}

		fmt.Printf("found %s\n", result)
		time.Sleep(time.Duration((rand.Intn(5) + 5) * 1000 * 1000 * 1000)) // sleep 10-20 seconds

		if outputColumnIndex == -1 {
			records[i] = append(records[i], result)
			continue
		}

		records[i][outputColumnIndex] = result
	}

	fmt.Println("CSV processing completed successfully.")

	err = csv.WriteCSV(*outputFile, records)
	if err != nil {
		fmt.Printf("Error writing output CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Output CSV written successfully.")
}
