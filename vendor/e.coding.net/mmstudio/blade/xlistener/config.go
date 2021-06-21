package xlistener

import (
	"fmt"
	"os"
	"time"
)

//go:generate optionGen
func ConfOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"BacklogAccept":    1024,
		"TimeoutCanRead":   time.Duration(time.Second * time.Duration(30)),
		"EnableHandshake":  false,
		"HandshakeTimeout": time.Duration(time.Second * time.Duration(10)),
		"Debugf": func(format string, v ...interface{}) {
			_, _ = fmt.Fprintf(os.Stdout, format, v...)
		},
		"Warningf": func(format string, v ...interface{}) {
			_, _ = fmt.Fprintf(os.Stderr, format, v...)
		},
	}
}
