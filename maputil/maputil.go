package maputil

import "encoding/json"

func MapToStruct(m any, s any) error {
	b, e := json.Marshal(m)
	if e != nil {
		return e
	}

	e = json.Unmarshal(b, &s)

	if e != nil {
		return e
	}

	return nil
}
