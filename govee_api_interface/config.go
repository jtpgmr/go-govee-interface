package goveeapiinterface

// DeviceConfig represents the valid configurations for a device model
type DeviceConfig struct {
    SupportCmds []SupportedCommand `json:"supportCmds"`
    Range       Range    `json:"range"`
}

// Configurations holds valid configurations for each device model
var ModelConfigurations = map[string]DeviceConfig {
    "H6159": {
		Range: Range{ Min: 2000, Max: 9000 },
        SupportCmds: []SupportedCommand{TURN, BRIGHTNESS, COLOR, COLORTEM},
    },
    "H5081": {
        SupportCmds: []SupportedCommand{TURN},
        // No Range specified for H5081
    },
	"H6008": {
		Range: Range{ Min: 2700, Max: 6500 },
        SupportCmds: []SupportedCommand{TURN, BRIGHTNESS, COLOR, COLORTEM},
	},
}

// isValidSupportCmd checks if a support command is valid
func isValidSupportCmd(cmd SupportedCommand, validCmds []SupportedCommand) bool {
    for _, validCmd := range validCmds {
        if cmd == validCmd {
            return true
        }
    }
    return false
}


// // Validate validates the GoveeResponse object
// func (g *GoveeResponseDetails) Validate() error {
//     for _, device := range g.Data.Devices {
//         // Check if the model exists in the Configurations map
//         config, ok := ModelConfigurations[device.Model]
//         if !ok {
//             return fmt.Errorf("invalid model: %s", device.Model)
//         }

//         // Validate SupportCmds
//         for _, cmd := range device.SupportCmds {
//             if !isValidSupportCmd(cmd, config.SupportCmds) {
//                 return fmt.Errorf("invalid support command: %s", cmd)
//             }
//         }

//         // Validate Range
//         if config.Range.Min != 0 && config.Range.Max != 0 && (device.Properties["colorTem"].Range.Min < config.Range.Min || device.Properties["colorTem"].Range.Max > config.Range.Max) {
//             return fmt.Errorf("invalid range for model %s: %+v", device.Model, device.Properties["colorTem"].Range)
//         }
//     }

//     return nil
// }


