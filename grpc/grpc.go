package grpc

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"
)

func ToRpcStruct(data any) *structpb.Struct {
	b, e := json.Marshal(data)
	if e != nil {
		return nil
	}
	var m map[string]any
	e = json.Unmarshal(b, &m)

	if e != nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}

	return s
}
