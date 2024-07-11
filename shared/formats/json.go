package formats

import (
	"encoding/json"
	"errors"
	"log"
)

func EncodeJson(data any) []byte {
	rawJson, err := json.Marshal(data)
	if err != nil {
		log.Default().Panic(err)
	}
	return rawJson
}

func DecodeJson[T any](data []byte) (*T, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	var result T
	err := json.Unmarshal(data, &result)
	return &result, err
}
