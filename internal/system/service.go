package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// ServiceStatus represents the current state of a systemd service.
type ServiceStatus struct {
	Name    string
	Active  bool
	Enabled bool
}

// IsServiceActive checks if a given systemd service is currently running.
func IsServiceActive(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
	err := cmd.Run()
	return err == nil
}

// IsServiceEnabled checks if a given systemd service is set to start on boot.
func IsServiceEnabled(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-enabled", "--quiet", serviceName)
	err := cmd.Run()
	return err == nil
}

// GetServiceStatus returns the full status (active/enabled) of a service.
func GetServiceStatus(serviceName string) ServiceStatus {
	return ServiceStatus{
		Name:    serviceName,
		Active:  IsServiceActive(serviceName),
		Enabled: IsServiceEnabled(serviceName),
	}
}

// StartService starts a systemd service using sudo.
func StartService(serviceName string) error {
	cmd := exec.Command("sudo", "systemctl", "start", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start %s: %v\n%s", serviceName, err, strings.TrimSpace(string(output)))
	}
	return nil
}

// StopService stops a systemd service using sudo.
func StopService(serviceName string) error {
	cmd := exec.Command("sudo", "systemctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop %s: %v\n%s", serviceName, err, strings.TrimSpace(string(output)))
	}
	return nil
}

// EnableService enables a systemd service to start on boot using sudo.
func EnableService(serviceName string) error {
	cmd := exec.Command("sudo", "systemctl", "enable", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable %s: %v\n%s", serviceName, err, strings.TrimSpace(string(output)))
	}
	return nil
}

// DisableService disables a systemd service from starting on boot using sudo.
func DisableService(serviceName string) error {
	cmd := exec.Command("sudo", "systemctl", "disable", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable %s: %v\n%s", serviceName, err, strings.TrimSpace(string(output)))
	}
	return nil
}
