package rest

import (
	"fmt"

	"cleverreach.com/crtools/crconfig"
)

var (
	// BasePath is added before every path given in NewClient()
	BasePath string

	// Debug sets the output accordingly
	Debug bool

	// IgnoreCertificate makes the client ignoring any certificate warnings
	IgnoreCertificate bool

	// Logger can be set to get debug information.
	// it has to implement ClientLogger which should match most loggers.
	// Default logger uses fmt.Println
	Logger ClientLogger = &defaultLogger{}
)

type (
	// ClientLogger must be implemented to be used as logger for rest client
	ClientLogger interface {
		Warnln(...interface{})
		Debugln(...interface{})
	}
	// M is short for map[string]interface{}
	M map[string]interface{}

	defaultLogger struct{}
)

func (l *defaultLogger) Debugln(args ...interface{}) {
	if crconfig.GetBool("DEBUG", Debug) {
		fmt.Println(args...)
	}
}
func (l *defaultLogger) Warnln(args ...interface{}) {
	fmt.Println(args...)
}
