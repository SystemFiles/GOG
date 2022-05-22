package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"sykesdev.ca/gog/internal/common"
)

type Logging struct {
	Level string
}

var (
	once sync.Once
	instance Logging
)
var SeverityLevels = []string{"INFO", "DEBUG", "WARN", "ERROR"}

func Instance() *Logging {
	once.Do(func() {
		var lvl string
		if common.StringInSlice(SeverityLevels, os.Getenv("GOG_LOG_LEVEL")) {
			lvl = os.Getenv("GOG_LOG_LEVEL")
		} else {
			lvl = "INFO"
		}

		instance = Logging{Level: lvl}
	})

	return &instance
}

func (l *Logging) Setup(level string) {
	if common.StringInSlice(SeverityLevels, strings.ToUpper(level)) {
		l.Level = strings.ToUpper(level)
	}
}

func (l Logging) Info(message string) {
	fmt.Printf("%v-%v-%v %v:%v:%v [INFO] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		message)
}

func (l Logging) Infof(fmtMessage string, arguments ...interface{}) {
	fmt.Printf("%v-%v-%v %v:%v:%v [INFO] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		fmt.Sprintf(fmtMessage, arguments...))
}

func (l Logging) Warn(message string) {
	if l.Level == "WARN" || l.Level == "DEBUG" || l.Level == "INFO" {
		fmt.Printf("%v-%v-%v %v:%v:%v [WARN] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		message)
	}
}

func (l Logging) Warnf(fmtMessage string, arguments ...interface{}) {
	if l.Level == "WARN" || l.Level == "DEBUG" || l.Level == "INFO" {
		fmt.Printf("%v-%v-%v %v:%v:%v [WARN] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		fmt.Sprintf(fmtMessage, arguments...))
	}
}

func (l Logging) Debug(message string) {
	if l.Level == "DEBUG" {
		fmt.Printf("%v-%v-%v %v:%v:%v [DEBUG] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		message)
	}
}

func (l Logging) Debugf(fmtMessage string, arguments ...interface{}) {
	if l.Level == "DEBUG" {
		fmt.Printf("%v-%v-%v %v:%v:%v [DEBUG] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		fmt.Sprintf(fmtMessage, arguments...))
	}
}

func (l Logging) Error(message string) {
	fmt.Printf("%v-%v-%v %v:%v:%v [ERROR] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		message)
}

func (l Logging) Errorf(fmtMessage string, arguments ...interface{}) {
	fmt.Printf("%v-%v-%v %v:%v:%v [ERROR] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		fmt.Sprintf(fmtMessage, arguments...))
}

func (l Logging) Fatal(message string) {
	fmt.Printf("%v-%v-%v %v:%v:%v [FATAL] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		message)
	os.Exit(1)
}

func (l Logging) Fatalf(fmtMessage string, arguments ...interface{}) {
	fmt.Printf("%v-%v-%v %v:%v:%v [FATAL] - %v\n",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		fmt.Sprintf(fmtMessage, arguments...))
	os.Exit(1)
}