package utils

import (
	"log"
	"net/url"
)

func ApplyQueryParams(endpointUrl *string, params map[string]string) {
    // Encode query parameters directly into the URL
    if len(params) > 0 {
        // Parse the URL to handle existing query parameters
        u, err := url.Parse(*endpointUrl)
        if err != nil {
			log.Printf("Warning: Failed to parse URL: %v. URL will not update ", err)
            return
        }

        // Add new query parameters to the existing ones
        query := u.Query()
        for key, value := range params {
            query.Set(key, value)
        }
        u.RawQuery = query.Encode()

        // Update the endpointUrl's value in memory to the url string
        *endpointUrl = u.String()
    }
}