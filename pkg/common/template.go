package common

import (
	"bytes"
	"text/template"
)

func ParserTemplate(filePath string, obj interface{}) ([]byte, error) {

	t, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, obj)
	return buf.Bytes(), err
}
