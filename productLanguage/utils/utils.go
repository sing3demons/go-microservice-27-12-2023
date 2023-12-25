package utils

import (
	"os"
)

func HostName(name, id string) string {
	host := os.Getenv("API_GATEWAY_URL")
	return host + "/" + name + "/" + id
}
