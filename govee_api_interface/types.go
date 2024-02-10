package goveeapiinterface

// GoveeAPI represents the Govee API client
type GoveeAPIInterface struct {
    BaseURL string
    APIKey  string
}

type SupportedCommand string
const (
    TURN    SupportedCommand = "turn"
    BRIGHTNESS    SupportedCommand = "brightness"
    COLOR    SupportedCommand = "color"
    COLORTEM SupportedCommand = "colorTem"
)

// Device represents the structure of each device object in the "devices" array
type GoveeDevice struct {
    Device        string            `json:"device"`
    Model         string            `json:"model"`
    DeviceName    string            `json:"deviceName"`
    Controllable  bool              `json:"controllable"`
    Retrievable   bool              `json:"retrievable"`
    SupportCmds   []SupportedCommand          `json:"supportCmds"`
    Properties    map[string]Property `json:"properties"`
}

// Data represents the structure of the "data" field in the response
type GoveeDevices struct {
    Devices []GoveeDevice `json:"devices"`
}

// Range represents the structure of a range object
type Range struct {
    Min int `json:"min"`
    Max int `json:"max"`
}

// Property represents the structure of a property object
type Property struct {
    Range Range `json:"range"`
}


type HttpMethod string
const (
    GET    HttpMethod = "GET"
    POST   HttpMethod = "POST"
    PUT    HttpMethod = "PUT"
    DELETE HttpMethod = "DELETE"
)

type ControlDevicesCmd struct {
    name SupportedCommand
    value interface{}
}

type ControlDevicesRequestBody struct {
    device string
    model string
    cmd ControlDevicesCmd
}

type GoveeResponseDetails struct {
    Message string   `json:"message"`
    StatusCode    int      `json:"code"`
    Data    GoveeDevices     `json:"data"`
    HasErrors bool  `json:"hasErrors"`
    Err       error `json:"error"`
}

type ResponseDetailsFunc func(*GoveeResponseDetails)

func SucessfulResponse(message string, statusCode int, data GoveeDevices) ResponseDetailsFunc {
    return func(r *GoveeResponseDetails) {
        r.Message = message
        r.StatusCode = statusCode
        r.Data = data
        r.HasErrors = false
        r.Err = nil
    }
}

func FailureResponse(message string, statusCode int, err error) ResponseDetailsFunc {
    return func(r *GoveeResponseDetails) {
        r.Message = message
        r.StatusCode = statusCode
        r.HasErrors = true
        r.Err = err
    }
}
func NewGoveeAPIResponse(functions ...ResponseDetailsFunc) *GoveeResponseDetails {
    response := &GoveeResponseDetails{}

    for _, function := range functions {
        function(response)
    }

    return response
}