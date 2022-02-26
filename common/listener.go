package common

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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
	d := json.NewDecoder(os.Stdin)
	var data map[string]interface{}

	for {
		err := d.Decode(&data)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
		message := lxtypes.NewEvent(lxtypes.IncomingMessageEvent, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		*event <- message
	}
}

func (this *LxCommon) Send(message *lxtypes.Event) {
	m, err := ToMap(message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	this.logger.Println(m)
}

func (this *LxCommon) Close() {
	message := lxtypes.NewEvent(lxtypes.CloseEvent, nil)
	m, err := ToMap(message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	this.logger.Println(m)
}
