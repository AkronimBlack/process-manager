package shared

import (
	"encoding/json"
	"log"
)

func ToJsonByte(value interface{}) []byte {
	v, err := json.Marshal(value)
	if err != nil {
		log.Println(err.Error())
	}
	return v
}

func ToJsonPrettyByte(value interface{}) []byte {
	v, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		log.Println(err.Error())
	}
	return v
}

func ToJsonString(value interface{}) string {
	return string(ToJsonByte(value))
}

func ToJsonPrettyString(value interface{}) string {
	return string(ToJsonPrettyByte(value))
}
