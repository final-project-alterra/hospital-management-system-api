package main

import (
	"github.com/final-project-alterra/hospital-management-system-api/config"
)

func main() {
	path := ".env"
	config.LoadENV(path)
	config.ConnectDB()
}
