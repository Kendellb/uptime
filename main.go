package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
    "log"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>System Uptime</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <h1>System Uptime</h1>
    <p>{{.}}</p>
</body>
</html>`

func getUptime() (string, error) {
	cmd := exec.Command("uptime") // Get full uptime output
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fields := strings.Fields(string(output))
	if len(fields) < 3 {
		return "Unknown uptime format", nil
	}

	// Extract only the first part of uptime (e.g., "up 3 days, 4:15")
	uptimeStartIndex := strings.Index(string(output), "up ")
	if uptimeStartIndex == -1 {
		return "Unknown uptime format", nil
	}

	uptime := strings.SplitN(string(output[uptimeStartIndex+3:]), ",", 2)[0] // Extract the first part only
	return strings.TrimSpace(uptime), nil
}

func uptimeHandler(w http.ResponseWriter, r *http.Request) {
	uptime, err := getUptime()
	if err != nil {
		http.Error(w, "Failed to get uptime", http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.New("uptime").Parse(htmlTemplate))
	tmpl.Execute(w, uptime)
}

func main() {
    // Serve static files from the "static" directory
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static", fs))

    // Handle uptime endpoint
    http.HandleFunc("/", uptimeHandler)

    fmt.Println("Server running on http://localhost:8080")
    
    // Start server and log errors
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err) // Log any errors that occur
    }
}

