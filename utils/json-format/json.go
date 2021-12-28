package jsonformat

import "encoding/json"

func JSON(v interface{}) string {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "Failed marshalling to JSON"
	}
	return string(result)
}
