package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

//var (
//	once           sync.Once
//	loggerInstance *zap.Logger
//)
//
//func GetLogger() *zap.Logger {
//	once.Do(func() {
//		var err error
//		loggerInstance, err = zap.NewDevelopment()
//		if err != nil {
//			panic(err) // this will end the program if no zap logger can be created
//		}
//	})
//
//	return loggerInstance
//}

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
