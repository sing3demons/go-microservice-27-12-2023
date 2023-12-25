package utils

import (
	"fmt"
	"os"
)

func HostName(name, id string) string {
	host := os.Getenv("API_GATEWAY_URL")
	return host + "/" + (name) + "/" + id
}

func Href(host, name, id string) string {
	return fmt.Sprintf("%s/%s/%s", host, name, id)
}
