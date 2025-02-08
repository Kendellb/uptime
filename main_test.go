package main

import (
	"os/exec"
	"strings"
	"testing"
)

// Mock function to simulate `uptime` output
func mockUptimeOutput() string {
	return " 10:24:56 up 3 days, 4:15,  2 users,  load average: 0.20, 0.15, 0.10"
}

func TestGetUptime(t *testing.T) {
	// Simulate extracting uptime from mock output
	mockOutput := mockUptimeOutput()
	uptimeIndex := strings.Index(mockOutput, "up ")
	if uptimeIndex == -1 {
		t.Fatalf("Failed to find 'up ' in mock uptime output")
	}

	uptimeStr := mockOutput[uptimeIndex+3:]
	parts := strings.Split(uptimeStr, ",")
	uptime := strings.Join(parts[:min(len(parts), 2)], ",")

	expected := "3 days, 4:15"
	if uptime != expected {
		t.Errorf("Expected uptime '%s', but got '%s'", expected, uptime)
	}
}

func TestGetUptimeCommand(t *testing.T) {
	_, err := exec.LookPath("uptime")
	if err != nil {
		t.Skip("Skipping test: 'uptime' command not found")
	}

	uptime, err := getUptime()
	if err != nil {
		t.Fatalf("Error running getUptime: %v", err)
	}

	if uptime == "Unknown uptime format" || uptime == "" {
		t.Errorf("Unexpected uptime output: '%s'", uptime)
	}
}

