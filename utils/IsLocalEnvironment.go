package utils

import (
	"net"
	"os"
)

func IsLocalEnvironment() bool {
	// Check if the PORT environment variable is set (common in production)
	if os.Getenv("PORT") != "" {
		return false
	}

	// Check if running on localhost (hostname/IP check)
	host, err := os.Hostname()
	if err == nil {
		addrs, err := net.LookupHost(host)
		if err == nil {
			for _, addr := range addrs {
				if addr == "127.0.0.1" || addr == "::1" { // Localhost
					return true
				}
			}
		}
	}

	// Check if the .env file exists (assumes local development)
	if _, err := os.Stat(".env"); err == nil {
		return true
	}

	return false
}
