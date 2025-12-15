package configuration

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Unmarshaler func(data []byte, v any) error

func UnmarshalJSON(data []byte, v any) (err error) {
	if m, ok := v.(proto.Message); ok {
		return protojson.Unmarshal(data, m)
	}
	return json.Unmarshal(data, v)
}
