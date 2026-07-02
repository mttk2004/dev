package hosts

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// HostEntry represents a host routing entry in /etc/hosts.
type HostEntry struct {
	IP     string
	Domain string
}

// GetLocalDomains parses /etc/hosts and returns custom domains mapped to 127.0.0.1 or ::1.
func GetLocalDomains() ([]HostEntry, error) {
	file, err := os.Open("/etc/hosts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []HostEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		ip := fields[0]
		// Only collect custom local mappings, ignore localhost
		if ip == "127.0.0.1" || ip == "::1" {
			for _, domain := range fields[1:] {
				if domain == "localhost" || domain == "localhost.localdomain" || domain == "ip6-localhost" || domain == "ip6-loopback" {
					continue
				}
				entries = append(entries, HostEntry{
					IP:     ip,
					Domain: domain,
				})
			}
		}
	}

	return entries, scanner.Err()
}

// AddLocalDomain adds a new local domain mapped to 127.0.0.1 in /etc/hosts.
func AddLocalDomain(domain string) error {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return fmt.Errorf("domain name cannot be empty")
	}

	existing, err := GetLocalDomains()
	if err == nil {
		for _, e := range existing {
			if e.Domain == domain {
				return fmt.Errorf("domain %s is already configured", domain)
			}
		}
	}

	entry := fmt.Sprintf("\n127.0.0.1  %s", domain)
	cmd := exec.Command("sudo", "tee", "-a", "/etc/hosts")
	cmd.Stdin = strings.NewReader(entry)
	return cmd.Run()
}

// RemoveLocalDomain removes a local domain entry from /etc/hosts using sed.
func RemoveLocalDomain(domain string) error {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return fmt.Errorf("domain name cannot be empty")
	}

	pattern := fmt.Sprintf("/[[:space:]]%s[[:space:]]*$/d", domain)
	cmd := exec.Command("sudo", "sed", "-i", pattern, "/etc/hosts")
	return cmd.Run()
}
