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

	goveeInterface := govee.NewGoveeAPI(apiKey)

	// var dev string
	// dev = ""

	// // Retrieve devices
	// devices := govee.GetDevices(&dev)

	// fmt.Println(devices)
	
	// Retrieve devices
	// devices := govee.GetDeviceState("", "")

	// fmt.Println(devices)

	body := govee.ControlDevicesRequestBody{
		Device: "",
		Model: "",
		Cmd: govee.ControlDevicesCmd{
			// Name: "turn",
			// Value: "on",
			Name: "color",
			Value: map[string]int{
				"r": 100,
				"g": 100,
				"b": 100,
			},
		},
	}

	// Retrieve devices
	devices := goveeInterface.UpdateDeviceSettings(body)

	fmt.Println(devices)

	// Print the response body
	// utils.PrintObject(devices)
}