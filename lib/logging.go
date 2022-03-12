package lib

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logging struct {
	Level string
}

var once sync.Once
var instance Logging

func GetLogger() Logging {
	once.Do(func() {
		var lvl string
		levels := []string{"INFO", "DEBUG", "WARN", "ERROR"}
		if StringInSlice(levels, os.Getenv("GOG_LOG_LEVEL")) {
			lvl = os.Getenv("GOG_LOG_LEVEL")
		} else {
			lvl = "INFO"
		}

		instance = Logging{Level: lvl}
	})

	return instance
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