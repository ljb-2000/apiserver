package jsonx

import (
	"encoding/json"
)

//-----------------------
//ToJson struct to json
func ToJson(v ...interface{}) (string, error) {
	vbyte, err := json.MarshalIndent(v, " ", "    ")

	if err != nil {
		return "", err
	}

	return string(vbyte), nil
}
