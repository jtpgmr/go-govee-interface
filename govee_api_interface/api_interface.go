package goveeapiinterface

import (
	"bytes"
	"encoding/json"
	"errors"
	utils "govee-smart-tech-interface/utils"
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
	switch cmd.Name {
        case TURN:
            name := strings.ToLower(string(cmd.Name))
            if name != "on" && name != "off" {
                return errors.New("Invalid value for 'turn' command. Valid values are: 'on', 'off'")
            }

		case BRIGHTNESS:
            brightnessLevel, isInt := cmd.Value.(int)
            
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

type CreateGoveeAPIRequestProps struct {
    method HttpMethod
    endpoint GoveeEndpoints
    successMessage string
    deviceName *string
    body *ControlDevicesRequestBody
    queryParams    *map[string]string
}

// GetDevices retrieves devices from the Govee API
func (gAPI *GoveeAPIInterface) CreateGoveeAPIRequest(props *CreateGoveeAPIRequestProps) string {
    method := props.method
    endpointUrl := gAPI.BaseURL + string(props.endpoint)

    // checks for query params and adds them to endpointUrl if applicable
    if props.queryParams != nil {
        utils.ApplyQueryParams(&endpointUrl, *props.queryParams)
    }

    
    // Create a new request
    req, err := http.NewRequest(string(method), endpointUrl, nil)
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

    body := props.body
    
    if method == PUT && body != nil {
        var buf bytes.Buffer
        // shorthand error handling is used when only evaluating the err that is returned
        if err := json.NewEncoder(&buf).Encode(body); err != nil {
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
        responseJSON, _ := json.Marshal(response) 
        return string(responseJSON)
    }
    
    var jsonData []byte
    // fmt.Println(responseData.Data)

    // if props.deviceName != nil {
    //     var filteredDevices []GoveeDevice
    //     for _, device := range responseData.Data.Devices {
    //         if device.DeviceName == *props.deviceName {
    //             filteredDevices = append(filteredDevices, device)
    //         }
    //     }

    //     if filteredDevices == nil {
    //         response := NewGoveeAPIResponse(FailureResponse("Error searching for specified device", 404, fmt.Errorf("Device with the name %s was not found", *props.deviceName)))
    //         responseJSON, _ := json.Marshal(response) 
    //         return string(responseJSON)
    //     }

    //     responseData.Data.Devices = filteredDevices
    // } 

    // Marshal the response data back to JSON
    jsonData, err = json.Marshal(responseData)
    if err != nil {
        response := NewGoveeAPIResponse(FailureResponse("Error marshalling JSON response", 500, err))
        responseJSON, _ := json.Marshal(response) 
        return string(responseJSON)
    }

    // Return the JSON data as a string
    return string(jsonData)
}

type GoveeEndpoints string
const (
    // endpoints for LPS
    GetDevices    GoveeEndpoints = "/v1/devices"
    GetDeviceState   GoveeEndpoints = GetDevices + "/state"
    UpdateDeviceSettings    GoveeEndpoints = GetDevices + "/control"
)

// gets devices associated with API key
func (gAPI *GoveeAPIInterface) GetDevices(deviceName *string) string {
    var successMessage string 
    if deviceName != nil {
        successMessage = "Details for device named" + " " + *deviceName + " " + "successfully retrieved"
    } else {
        successMessage = "Details for all devices retrieved"
    }

    inputProps := CreateGoveeAPIRequestProps{
        method: GET,
        endpoint: GetDevices,
        successMessage: successMessage,
        deviceName: deviceName,
    }

    return gAPI.CreateGoveeAPIRequest(&inputProps)
}

// query the state of the device
// need to make a db connection which obtains the model if not already defined
func (gAPI *GoveeAPIInterface) GetDeviceState(macAddress string, model string) string {
    inputProps := CreateGoveeAPIRequestProps{
        method: GET,
        endpoint: GetDeviceState,
        queryParams: &map[string]string{
            "device": macAddress,
            "model": model,
        },
    }

    return gAPI.CreateGoveeAPIRequest(&inputProps)
}

// query the state of the device
func (gAPI *GoveeAPIInterface) UpdateDeviceSettings(body ControlDevicesRequestBody) string {
    inputProps := CreateGoveeAPIRequestProps{
        method: PUT,
        endpoint: UpdateDeviceSettings,
        body: &body,
    }

    return gAPI.CreateGoveeAPIRequest(&inputProps)
}