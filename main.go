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

    uptimeIndex := strings.Index(string(output), "up ")
	if uptimeIndex == -1 {
		return "Unknown uptime format", nil
	}

	// Extract the uptime portion after "up "
	uptimeStr := string(output[uptimeIndex+3:])
	parts := strings.Split(uptimeStr, ",")

	// Extract at most first two components (to get both days and hours)
	uptime := strings.Join(parts[:min(len(parts), 2)], ",")
	return strings.TrimSpace(uptime), nil
}

func min(a,b int) int {
    if a < b {
        return a
    }
    return b
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

