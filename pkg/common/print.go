package common

import (
	"encoding/json"
)

func PrettifyJson(i interface{}, indent bool) string {
	var str []byte
	if indent {
		str, _ = json.MarshalIndent(i, "", "    ")
	} else {
		str, _ = json.Marshal(i)
	}

	return string(str)
}
