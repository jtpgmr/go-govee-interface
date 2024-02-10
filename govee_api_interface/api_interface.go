package goveeapiinterface

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var baseUrl string = "https://developer-api.govee.com"

// NewGoveeAPI creates a new GoveeAPI instance
func NewGoveeAPI(apiKey string) *GoveeAPIInterface {
    return &GoveeAPIInterface{
        BaseURL: baseUrl,
        APIKey:  apiKey,
    }
}

func validateCmd(cmd ControlDevicesCmd, value string) error {
	switch cmd.name {
        case TURN:
            name := strings.ToLower(string(cmd.name))
            if name != "on" && name != "off" {
                return errors.New("Invalid value for 'turn' command. Valid values are: 'on', 'off'")
            }

		case BRIGHTNESS:
            brightnessLevel, isInt := cmd.value.(int)
            
            if !isInt {
                return errors.New("Value for brightness must be an integer")
            }

            if brightnessLevel < 0 || brightnessLevel > 100 {
                return errors.New("Value for brightness must be in the range: 0 - 100")
            }

        case COLOR:

        case COLORTEM:
		default:
			return errors.New("Invalid command. Valid options are: turn, brightness, color, colorTem")
	}
	return errors.New("Value not found")
}

// Instance to send request
var client = &http.Client{}

// GetDevices retrieves devices from the Govee API
func (gAPI *GoveeAPIInterface) CreateGoveeAPIRequest(method HttpMethod, endPoint string, successMessage string, body ...*ControlDevicesRequestBody) string {
    // Create a new request
    req, err := http.NewRequest(string(method), gAPI.BaseURL + endPoint, nil)
    if err != nil {
        response := NewGoveeAPIResponse(FailureResponse("Error creating request", 500, err))
        responseJSON, _ := json.Marshal(response) 
        if err != nil {
            // If there's an error marshaling the response, return a generic error message
            return `{"message": "Error creating request"}`
        }
        return string(responseJSON)
    }

    // Add headers to the request
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Govee-API-Key", gAPI.APIKey)
    
    if method == PUT && len(body) > 0 && body[0] != nil {
        var buf bytes.Buffer
        // shorthand error handling is used when only evaluating the err that is returned
        if err := json.NewEncoder(&buf).Encode(body[0]); err != nil {
            response := NewGoveeAPIResponse(FailureResponse("Error encoding body", 500, err))
            responseJSON, _ := json.Marshal(response) // Ignore error here
            if err != nil {
                // If there's an error marshaling the response, return a generic error message
                return `{"message": "Error creating request"}`
            }
            return string(responseJSON)
        }
        
        req.Body = io.NopCloser(&buf)
    }

	// send request
    res, err := client.Do(req)
    if err != nil {
        response := NewGoveeAPIResponse(FailureResponse("Failure to send request", 500, err))
        responseJSON, _ := json.Marshal(response) // Ignore error here
        if err != nil {
            // If there's an error marshaling the response, return a generic error message
            return `{"message": "Error creating request"}`
        }
        return string(responseJSON)
    }
    defer res.Body.Close()

    // Read the response body
    resBody, err := io.ReadAll(res.Body)
    if err != nil {
        response := NewGoveeAPIResponse(FailureResponse("Error reading response", 500, err))
        responseJSON, _ := json.Marshal(response) 
        if err != nil {
            // If there's an error marshaling the response, return a generic error message
            return `{"message": "Error creating request"}`
        }
        return string(responseJSON)
    }

    // Unmarshal the response body into GoveeResponseDetails
    var responseData GoveeResponseDetails
    if err := json.Unmarshal(resBody, &responseData); err != nil {
        response := NewGoveeAPIResponse(FailureResponse("Error unmarshalling response", 500, err))
        responseJSON, _ := json.Marshal(response) // Ignore error here
        return string(responseJSON)
    }

    // Marshal the response data back to JSON
    jsonData, err := json.Marshal(responseData)
    if err != nil {
        response := NewGoveeAPIResponse(FailureResponse("Error marshalling JSON response", 500, err))
        responseJSON, _ := json.Marshal(response) // Ignore error here
        return string(responseJSON)
    }

    // Return the JSON data as a string
    return string(jsonData)
}

// gets devices associated with API key
func (gAPI *GoveeAPIInterface) GetDevices(deviceName *string) string {
    var successMessage string 
    if deviceName != nil {
        successMessage = "Details for device named " + *deviceName + "successfully retrieved"
    } else {
        successMessage = "Details for all devices retrieved"
    }
    return gAPI.CreateGoveeAPIRequest(GET, "/v1/devices", successMessage)
}

// func (gAPI *GoveeAPIInterface) GetDevices() (string, []byte, error) {
//     return gAPI.CreateGoveeAPIRequest(GET, "/v1/devices")
// }