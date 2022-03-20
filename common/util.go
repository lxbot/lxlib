package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lxbot/lxlib/v2/lxtypes"
)

func ToMap(i interface{}) (lxtypes.M, error) {
	TraceLog("lxlib.common.ToMap()", "start")
	defer TraceLog("lxlib.common.ToMap()", "end")

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
	TraceLog("lxlib.common.ToJSON()", "start")
	defer TraceLog("lxlib.common.ToJSON()", "end")

	j, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func FromJSON(i json.RawMessage) (lxtypes.M, error) {
	TraceLog("lxlib.common.FromJSON()", "start")
	defer TraceLog("lxlib.common.FromJSON()", "end")

	var r lxtypes.M
	if err := json.Unmarshal(i, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func FatalLog(a ...interface{}) {
	t := make([]interface{}, 0)
	t = append(t, "[FATAL]")
	t = append(t, a...)
	fmt.Fprintln(os.Stderr, t...)
	os.Exit(1)
}

func ErrorLog(a ...interface{}) {
	t := make([]interface{}, 0)
	t = append(t, "[ERROR]")
	t = append(t, a...)
	fmt.Fprintln(os.Stderr, t...)
}

func WarnLog(a ...interface{}) {
	t := make([]interface{}, 0)
	t = append(t, "[WARN ]")
	t = append(t, a...)
	fmt.Fprintln(os.Stderr, t...)
}

func InfoLog(a ...interface{}) {
	t := make([]interface{}, 0)
	t = append(t, "[Info ]")
	t = append(t, a...)
	fmt.Fprintln(os.Stderr, t...)
}

var isTraceChecked bool
var isTrace bool

func checkIsTrace() bool {
	if !isTraceChecked {
		if d, ok := os.LookupEnv("LXLIB_TRACE"); ok && d != "0" {
			isTrace = true
		}
	}
	return isTrace
}

func TraceLog(a ...interface{}) {
	checkIsTrace()
	if isTrace {
		t := make([]interface{}, 0)
		t = append(t, "[TRACE]")
		t = append(t, a...)
		fmt.Fprintln(os.Stderr, t...)
	}
}

var isDebugChecked bool
var isDebug bool

func checkIsDebug() bool {
	if !isDebugChecked {
		if d, ok := os.LookupEnv("LXLIB_DEBUG"); ok && d != "0" {
			isDebug = true
		}
		if checkIsTrace() {
			isDebug = true
		}
	}
	return isDebug
}

func DebugLog(a ...interface{}) {
	checkIsDebug()
	if isDebug {
		t := make([]interface{}, 0)
		t = append(t, "[DEBUG]")
		t = append(t, a...)
		fmt.Fprintln(os.Stderr, t...)
	}
}
