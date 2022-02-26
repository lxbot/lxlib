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
	return &LxCommon{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (this *LxCommon) Listen(event *chan *lxtypes.Event) {
	s := bufio.NewScanner(os.Stdin)
	b := strings.Builder{}

	for {
		for s.Scan() {
			b.WriteString(s.Text())
		}
		if s.Err() == nil {
			line := b.String()
			if strings.HasSuffix("{", line) && strings.HasSuffix("}", line) {
				this.onMessage(line, event)
			}
		} else {
			ErrorLog(s.Err())
		}
		b.Reset()
	}
}

func (this *LxCommon) onMessage(line string, event *chan *lxtypes.Event) {
	d := json.NewDecoder(os.Stdin)
	var data lxtypes.StdInOutEvent

	err := d.Decode(&data)
	if err != nil {
		ErrorLog(err)
	}
	message := lxtypes.NewEvent(data.Event, data.Payload)
	message.ID = data.ID
	if err != nil {
		ErrorLog(err)
	}
	*event <- message
}

func (this *LxCommon) Send(message *lxtypes.Event) {
	m, err := ToJSON(message)
	if err != nil {
		ErrorLog(err)
	}
	this.logger.Println(m)
}

func (this *LxCommon) Close() {
	message := lxtypes.NewEvent(lxtypes.CloseEvent, nil)
	m, err := ToJSON(message)
	if err != nil {
		ErrorLog(err)
	}
	this.logger.Println(m)
}
