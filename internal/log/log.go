package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = &logrus.Logger{
	Out:       os.Stdout,
	Level:     logrus.DebugLevel,
	Formatter: &logrus.TextFormatter{},
}
