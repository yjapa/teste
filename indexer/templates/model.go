package templates

import (
	"bytes"
	"encoding/json"
)

// Array type rename type []interface{}
type Array []interface{}

// Object rename type map[string]interface{}
type Object map[string]interface{}

// ToBuffer convert an Object to a *bytes.Buffer
func (o *Object) ToBuffer() (*bytes.Buffer, error) {
	objectBytes, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	buff := &bytes.Buffer{}
	_, err = buff.Write(objectBytes)

	if err != nil {
		return nil, err
	}

	return buff, nil
}
