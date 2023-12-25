package utils

import (
	"fmt"
	"os"
)

func HostName(name, id string) string {
	host := os.Getenv("PRODUCT_SERVICE_URL")
	return host + "/" + (name) + "/" + id
}

func Href(host, name, id string) string {
	return fmt.Sprintf("%s/%s/%s", host, name, id)
}
