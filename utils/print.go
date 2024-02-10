package utils

import (
	"encoding/json"
	"fmt"
)

func PrintObject(obj interface{}) {
    // Convert the object to JSON
    objBytes, err := json.MarshalIndent(obj, "", "    ")
    if err != nil {
        fmt.Println("Error marshaling object to JSON: ", err)
        return
    }

    // Print the JSON-like representation
    fmt.Println(string(objBytes))
}