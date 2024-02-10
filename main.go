package main

import (
	"fmt"
	"os"

	govee "govee-smart-tech-interface/govee_api_interface"
	// utils "govee-smart-tech-interface/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	apiKey := os.Getenv("GOVEE_API_KEY")

	govee := govee.NewGoveeAPI(apiKey)

	var dev *string
	// Retrieve devices
	devices := govee.GetDevices(dev)

	fmt.Println(devices)

	// Print the response body
	// utils.PrintObject(devices)
}