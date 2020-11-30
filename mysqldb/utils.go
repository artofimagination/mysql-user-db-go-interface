package mysqldb

import (
	"encoding/json"
)

func ConvertToJSONRaw(references interface{}) (*json.RawMessage, error) {
	refBytes, err := json.Marshal(&references)
	if err != nil {
		return nil, err
	}

	refRaw := json.RawMessage(refBytes)
	return &refRaw, nil
}
