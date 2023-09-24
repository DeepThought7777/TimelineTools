package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var errInvalidArguments = errors.New("not enough arguments specified")

func main() {
	inputFolderFile := "inputFolders.csv"
	outputFolder := "R:\\TIMELINE"
	extensionsFile := "extensions.csv"

	extensions, err := readExtensions(extensionsFile)
	if err != nil {
		fmt.Printf("Error while reading the extensions file: %v\n", err)
		os.Exit(2)
	}

	inputFolders, err := readInputFolders(inputFolderFile)
	if err != nil {
		fmt.Printf("Error while reading the input folders file: %v\n", err)
		os.Exit(2)
	}

	for _, folder := range inputFolders {
		fmt.Printf("FOLDER:  %s", folder)
		err = filepath.Walk(folder, walkFn(outputFolder, extensions))
		if err != nil {
			fmt.Printf("Error while walking the input folder: %v\n", err)
			os.Exit(3)
		}
	}
}

func readInputFolders(filePath string) ([]string, error) {
	var inputFolders []string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(line) > 0 {
			inputFolders = append(inputFolders, line[0])
		}
	}

	return inputFolders, nil
}

func readExtensions(filePath string) ([]string, error) {
	var extensions []string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(line) > 0 {
			extensions = append(extensions, line[0])
		}
	}

	return extensions, nil
}

func hasValidExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}

func walkFn(outputFolder string, extensions []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Printf("Processing folder: %s\n", path)
			return nil
		}

		if err != nil {
			return err
		}

		if !hasValidExtension(path, extensions) {
			return nil
		}

		if strings.Contains(path, "$Recycle.Bin") {
			return nil
		}

		if strings.Contains(path, "$") {
			return nil
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		h := sha1.New()
		if _, err := io.Copy(h, src); err != nil {
			return err
		}

		hash := fmt.Sprintf("%x", h.Sum(nil))
		fileName := filepath.Base(path)
		ext := filepath.Ext(fileName)
		base := fileName[:len(fileName)-len(ext)]
		newFileName := base + "_" + hash[:6] + ext

		// Get the last modification time of the file
		modTime := info.ModTime().Format("2006_01_02")

		dstPath := filepath.Join(outputFolder, modTime, newFileName)
		dstDir := filepath.Dir(dstPath)
		if _, err := os.Stat(dstDir); os.IsNotExist(err) {
			if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
				return err
			}
		}

		src.Seek(0, 0)
		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	}
}
