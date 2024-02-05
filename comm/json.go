package comm

import (
	stdjson "encoding/json"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary // 比标准库快,容错性强,默认使用此库

// DefaultJSONUnmarshal 默认json包选择
func DefaultJSONUnmarshal(data []byte, v interface{}) (err error) {
	return json.Unmarshal(data, v)
}

// DefaultJSONmarshal .
func DefaultJSONmarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// STDJSONUnmarshal .
func STDJSONUnmarshal(data []byte, v interface{}) (err error) {
	return stdjson.Unmarshal(data, v)
}

// STDJSONmarshal .
func STDJSONmarshal(v interface{}) ([]byte, error) {
	return stdjson.Marshal(v)
}
