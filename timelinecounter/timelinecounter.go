package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var outfilename string = "file_count.csv"

func walkFn(counts map[string]map[string]int) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.Contains(path, outfilename) {
			return nil
		}

		if err != nil {
			return err
		}

		folderName := filepath.Base(filepath.Dir(path))
		fileType := strings.TrimPrefix(filepath.Ext(path), ".")

		if _, ok := counts[folderName]; !ok {
			counts[folderName] = make(map[string]int)
		}

		counts[folderName][fileType]++

		return nil
	}
}

var errInvalidArguments = errors.New("not enough arguments specified")

func processArguments() (string, error) {
	if len(os.Args) != 2 {
		fmt.Println("Taking default \\TIMELINE")
		return "\\TIMELINE", nil
	}

	return os.Args[1], nil
}

type record [3]string

func (r record) Less(other record) bool {
	if r[0] != other[0] {
		return r[0] > other[0]
	}
	return r[1] < other[1]
}

func main() {
	outputFolder, err := processArguments()
	if err != nil {
		fmt.Printf("Arguments are invalid: %v\n", err)
		os.Exit(1)
	}

	csvFile, err := os.Create(filepath.Join(outputFolder, outfilename))
	if err != nil {
		fmt.Printf("Error while creating the CSV file: %v\n", err)
		os.Exit(2)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(csvFile))
	defer csvWriter.Flush()
	csvWriter.Comma = ';'

	header := []string{"Folder", "File Type", "Count"}
	if err := csvWriter.Write(header); err != nil {
		fmt.Printf("Error while writing the header to the CSV file: %v\n", err)
		os.Exit(3)
	}

	counts := make(map[string]map[string]int)

	err = filepath.Walk(outputFolder, walkFn(counts))
	if err != nil && err != io.EOF {
		fmt.Printf("Error while walking the output folder: %v\n", err)
		os.Exit(4)
	}

	records := make([]record, 0)

	for folderName, fileCounts := range counts {
		for fileType, count := range fileCounts {
			records = append(records, record{folderName, fileType, fmt.Sprintf("%d", count)})
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Less(records[j])
	})

	for _, r := range records {
		if err := csvWriter.Write(r[:]); err != nil {
			fmt.Printf("Error while writing a record to the CSV file: %v\n", err)
			os.Exit(5)
		}
	}
}
