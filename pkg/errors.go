package pkg

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type Err struct {
	Err       error    `json:"-"`
	Msg       string   `json:"msg"`
	PublicMsg string   `json:"public_msg"`
	Paths     []string `json:"paths"`
}

func (e Err) Error() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e Err) Is(target error) bool {
	return e.Err == target
}

func Error(e error, args ...string) error {
	var pubMsg string
	if len(args) > 0 {
		pubMsg = args[0]
	}

	if pubMsg == "" {
		pubMsg = "something went wrong, contact support"
	}

	_, filename, line, _ := runtime.Caller(1)
	here := fmt.Sprintf("%s:%v", filename, line)
	here = here[strings.Index(here, "/coin-commerce/"):]

	newE, ok := e.(Err)
	if !ok {
		var msg string
		if e == nil {
			msg = "no error"
		} else {
			msg = e.Error()
		}

		return Err{
			Err:       e,
			Msg:       msg,
			PublicMsg: pubMsg,
			Paths:     []string{here},
		}
	}

	newE.Paths = append(newE.Paths, here)
	return e
}
