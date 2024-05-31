package log

import (
	"go.uber.org/zap"
	"sync"
)

var (
	once           sync.Once
	loggerInstance *zap.Logger
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error
		loggerInstance, err = zap.NewDevelopment()
		if err != nil {
			panic(err) // this will end the program if no zap logger can be created
		}
	})

	return loggerInstance
}
