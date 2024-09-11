package utils

import "encoding/json"

func ConvertToByte(data interface{}) []byte {
	dataByte, err := json.Marshal(data)

	if err != nil {
		return []byte("")
	}

	return dataByte
}
