package utils

import "encoding/json"

func JsonEncode(input interface{}) string {
	raw, _ := json.Marshal(input)

	return string(raw)
}