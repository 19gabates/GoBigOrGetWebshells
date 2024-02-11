package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	configFile := "linux_config.txt"
	executableName := "LinuxWebsite"

	directoryLineNumber := 2

	// Create the list variables
	directoryList, err := readLineFromFile(configFile, directoryLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate through the list of directories
	for _, dir := range directoryList {
		// Check if the directory exists, create it if not
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Error creating directory %s: %v\n", dir, err)
				continue
			}
			fmt.Printf("Directory created: %s\n", dir)
		}

		// Copy the Exe File to the directory
		destFilePath := filepath.Join(dir, executableName)

		// Open the source file
		sourceFile, err := os.Open(executableName)
		if err != nil {
			fmt.Printf("Error opening source file %s: %v\n", executableName, err)
			continue
		}
		defer sourceFile.Close()

		// Create the destination file
		destFile, err := os.Create(destFilePath)
		if err != nil {
			fmt.Printf("Error creating destination file %s: %v\n", destFilePath, err)
			continue
		}
		defer destFile.Close()

		// Set the necessary file permissions (adjust as needed)
		perm := os.FileMode(0755) // for example, set to 0755 for read, write, execute permissions for owner, and read, execute permissions for group and others
		err = destFile.Chmod(perm)
		if err != nil {
			fmt.Printf("Error setting permissions for %s: %v\n", destFilePath, err)
			continue
		}

		// Copy the contents of the source file to the destination file
		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			fmt.Printf("Error copying file to %s: %v\n", destFilePath, err)
			continue
		}

		fmt.Printf("File copied to: %s\n", destFilePath)
	}
}

func renameLinux() {
	configFile := "linux_config.txt"
	renameBinaryLineNumber := 8
	directoryLineNumber := 2

	// Read the Rename Binary line from the config file
	renameBinaryList, err := readLineFromFile(configFile, renameBinaryLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read the Directories line from the config file
	directoryList, err := readLineFromFile(configFile, directoryLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract the new binary names
	newBinaryNameList := strings.Split(renameBinaryList[0], ",")

	// Check if the lengths match
	if len(newBinaryNameList) != len(directoryList) {
		fmt.Println("Error: Mismatch in the number of directories and new binary names in the config file.")
		return
	}

	// Iterate through the directories and rename the file to the specified names
	for i, dir := range directoryList {
		oldFilePath := filepath.Join(dir, "LinuxWebsite")

		// Check if the index is within the valid range
		if i < len(newBinaryNameList) {
			newFilePath := filepath.Join(dir, newBinaryNameList[i])

			err := os.Rename(oldFilePath, newFilePath)
			if err != nil {
				fmt.Printf("Error renaming file %s: %v\n", "LinuxWebsite", err)
			} else {
				fmt.Printf("File renamed: %s -> %s\n", "LinuxWebsite", newBinaryNameList[i])
			}
		} else {
			fmt.Println("Error: Index out of range for newBinaryNameList.")
		}
	}
}

func linuxStart() {
	configFile := "linux_config.txt"
	portLineNumber := 4
	renameBinaryLineNumber := 8
	directoryLineNumber := 2

	// Read the Port line from the config file
	portList, err := readLineFromFile(configFile, portLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read the Rename Binary line from the config file
	renameBinaryList, err := readLineFromFile(configFile, renameBinaryLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read the Directories line from the config file
	directoryList, err := readLineFromFile(configFile, directoryLineNumber)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Extract the port, new binary names, and directories
	ports := strings.Split(portList[0], ",")
	newBinaryNameList := strings.Split(renameBinaryList[0], ",")

	// Iterate through the directories and start the websites
	for i, dir := range directoryList {
		// Construct the file path for the renamed binary
		filePath := filepath.Join(dir, newBinaryNameList[i])

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File not found in directory %s\n", dir)
			continue
		}

		// Start the website in the background with the specified port
		cmd := exec.Command(filePath, ports[i])
		cmd.Dir = dir

		err := cmd.Start()
		if err != nil {
			fmt.Printf("Error starting website in directory %s: %v\n", dir, err)
			continue
		}

		fmt.Printf("Website started in directory %s with port %s (PID: %d)\n", dir, ports[i], cmd.Process.Pid)

		// Save the PID to use later (appends to the file)
		savePIDToFile(cmd.Process.Pid, i)
	}
}

func savePIDToFile(pid int, index int) {
	// Save the PID to a file (appends to the file)
	pidFilePath := fmt.Sprintf("linux_pid_file.txt")
	pidFile, err := os.OpenFile(pidFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening PID file %s: %v\n", pidFilePath, err)
		return
	}
	defer pidFile.Close()

	_, err = pidFile.WriteString(fmt.Sprintf("%d\n", pid))
	if err != nil {
		fmt.Printf("Error writing PID to file %s: %v\n", pidFilePath, err)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: <name> <option>")
		fmt.Println("Linux Options:")
		fmt.Println("  -ld   Distribute the exe across the directories named in linux_config.txt")
		fmt.Println("  -ln   Rename files on Linux")
		fmt.Println("  -ls   Start Bigger Websites V2 on Linux")
		fmt.Println("  -lrm  Remove Bigger Websites files on Linux")
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
	case "-ld":
		fmt.Println("Running Distrubution on Linux")
		linuxDistrobution()
	case "-ln":
		fmt.Println("Renaming files on Linux")
		renameLinux()
	case "-ls":
		fmt.Println("Starting Bigger Websites V2 on Linux")
		linuxStart()
	case "-lh":
		fmt.Println("Starting Bigger Websites V2 on Linux")
	case "-lrm":
		fmt.Println("Starting Bigger Websites V2 on Linux")
	case "option2":
		fmt.Println(executableName)
	default:
		fmt.Println("Invalid option. Supported options: option1, option2")
	}
}
