package configuration

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

type Unmarshaler func(data []byte, v any) error

func ProtoJSONUnmarshaler(data []byte, v any) (err error) {
	if m, ok := v.(proto.Message); ok {
		return protojson.Unmarshal(data, m)
	}
	return json.Unmarshal(data, v)
}

func YAMLUnmarshaler(data []byte, v any) (err error) {
	return yaml.Unmarshal(data, v)
}

func JSONUnmarshaler(data []byte, v any) (err error) {
	return json.Unmarshal(data, v)
}

func ExpandUnmarshaler(unmarshaler Unmarshaler, mapping func(int, string) ([]byte, error)) Unmarshaler {
	return func(data []byte, v any) (err error) {
		data, err = ExpandBytesError(data, mapping)
		if err != nil {
			return err
		}
		return unmarshaler(data, v)
	}
}
