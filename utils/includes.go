package utils

import (
	"errors"
)

func includes(arr interface{}, val interface{}) error {
	switch arr.(type) {
		case []int:
			for _, element := range arr.([]int) {
				if element == val.(int) {
					return nil
				}
			}
		case []string:
			for _, element := range arr.([]string) {
				if element == val.(string) {
					return nil
				}
			}
		case []bool:
			for _, element := range arr.([]bool) {
				if element == val.(bool) {
					return nil
				}
			}
		default:
			return errors.New("Unsupported slice type. Valid options are: int, string, bool")
	}
	return errors.New("Value not found")
}