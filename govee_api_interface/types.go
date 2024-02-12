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
type GetDevicesData struct {
    Devices []GoveeDevice `json:"devices,omitempty"`
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
    Name SupportedCommand   `json:"name"`
    Value interface{}   `json:"value"`
}

type ControlDevicesRequestBody struct {
    Device string   `json:"device"`
    Model string    `json:"model"`
    Cmd ControlDevicesCmd   `json:"cmd"`
}

type GoveeResponseDetails struct {
    Message string   `json:"message"`
    StatusCode    int      `json:"code"`
    HasErrors bool  `json:"hasErrors"`
    Data       interface{}   `json:"data,omitempty"` // Omit if nil
    Err        *string           `json:"error,omitempty"` // Omit if nil
}

type ResponseDetailsFunc func(*GoveeResponseDetails)

func SucessfulResponse(message string, statusCode int, data interface{}) ResponseDetailsFunc {
    return func(r *GoveeResponseDetails) {
        r.Message = message
        r.StatusCode = statusCode
        r.Data = &data
        r.HasErrors = false
        r.Err = nil
    }
}

func FailureResponse(message string, statusCode int, err error) ResponseDetailsFunc {
    return func(r *GoveeResponseDetails) {
        r.Message = message
        r.StatusCode = statusCode
        r.HasErrors = true
        errMsg := err.Error()
        r.Err = &errMsg
        r.Data = nil
    }
}

func NewGoveeAPIResponse(functions ...ResponseDetailsFunc) *GoveeResponseDetails {
    response := &GoveeResponseDetails{}

    for _, function := range functions {
        function(response)
    }

    return response
}