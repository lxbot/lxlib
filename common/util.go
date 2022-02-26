package common

import (
	"encoding/json"
	"fmt"
	"os"

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

func ToJSON(i interface{}) (string, error) {
	j, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func FromJSON(i json.RawMessage) (lxtypes.M, error) {
	var r lxtypes.M
	if err := json.Unmarshal(i, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func ErrorLog(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}
