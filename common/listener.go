package common

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/lxbot/lxlib/v2/lxtypes"
)

type LxCommon struct {
	logger *log.Logger
}

func NewLxCommon() *LxCommon {
	TraceLog("lxlib.common.NewLxCommon()", "start")
	defer TraceLog("lxlib.common.NewLxCommon()", "end")

	return &LxCommon{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (this *LxCommon) Listen(event *chan *lxtypes.Event) {
	TraceLog("lxlib.common.Listen()", "start")
	defer TraceLog("lxlib.common.Listen()", "end")

	s := bufio.NewScanner(os.Stdin)
	for {
		for s.Scan() {
			line := s.Text()
			TraceLog("lxlib.common.Listen()", "scanned:", line)

			if !strings.HasSuffix("{", line) || !strings.HasSuffix("}", line) {
				TraceLog("lxlib.common.Listen()", "stdin seems json. fire onMessage()")
				this.onMessage(s.Bytes(), event)
			} else {
				TraceLog("lxlib.common.Listen()", "malformed event?", "skip line")
			}
		}
		if s.Err() != nil {
			ErrorLog(s.Err())
		}
	}
}

func (this *LxCommon) onMessage(line []byte, event *chan *lxtypes.Event) {
	TraceLog("lxlib.common.onMessage()", "start")
	defer TraceLog("lxlib.common.onMessage()", "end")

	var data lxtypes.StdInOutEvent
	err := json.Unmarshal(line, &data)
	if err != nil {
		ErrorLog(err)
	}
	message := lxtypes.NewEvent(data.Event, data.Payload)
	message.ID = data.ID
	TraceLog("lxlib.common.onMessage()", "event:", message.Event)
	*event <- message
}

func (this *LxCommon) Send(message *lxtypes.Event) {
	TraceLog("lxlib.common.Send()", "start")
	defer TraceLog("lxlib.common.Send()", "end")

	m, err := ToJSON(message)
	if err != nil {
		ErrorLog(err)
		return
	}
	TraceLog("lxlib.common.Send()", "json:", m)
	this.logger.Println(m)
}

func (this *LxCommon) Close() {
	TraceLog("lxlib.common.Close()", "start")
	defer TraceLog("lxlib.common.Close()", "end")

	message := lxtypes.NewEvent(lxtypes.CloseEvent, nil)
	m, err := ToJSON(message)
	if err != nil {
		ErrorLog(err)
		return
	}
	TraceLog("lxlib.common.Close()", "json:", m)
	this.logger.Println(m)
}
