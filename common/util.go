package common

import (
	"encoding/json"

	"github.com/lxbot/lxlib/v2/lxtypes"
)

func ToMap(i interface{}) (lxtypes.M, error) {
	t, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var r lxtypes.M
	if err := json.Unmarshal(t, &r); err != nil {
		return nil, err
	}
	return r, nil
}
