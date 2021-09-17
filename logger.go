package rome

import (
	"sync"

	"github.com/amermelao/rome/data"
)

type GeneralLogger interface {
	Panic(v ...interface{})
	Fatal(v ...interface{})
	Error(v ...interface{})
	Warning(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
	Trace(v ...interface{})
}

type sideRoad struct {
	way  chan data.Message
	data data.Messages
}

type CentralLogger struct {
	road    chan data.Messages
	logging GeneralLogger
	wg      sync.WaitGroup
}

func NewCentrCentralLogger(log GeneralLogger) *CentralLogger {
	rome := &CentralLogger{
		road:    make(chan data.Messages),
		logging: log,
	}
	go rome.handleChanel()

	return rome
}

func doLogging(logger GeneralLogger, lvl, msg string) {
	switch lvl {
	case "panic":
		logger.Panic(msg)
	case "fatal":
		logger.Fatal(msg)
	case "error":
		logger.Error(msg)
	case "warning":
		logger.Warning(msg)
	case "info":
		logger.Info(msg)
	case "debug":
		logger.Debug(msg)
	case "trace":
		logger.Trace(msg)
	}
}

func (logger *CentralLogger) handleChanel() {
	for value := range logger.road {
		for _, msg := range value {
			doLogging(logger.logging, msg.Level, msg.Content)
		}
		logger.wg.Done()
	}
}

func (logger *CentralLogger) Log(msg data.Messages) {
	logger.wg.Add(1)
	go func() { logger.road <- msg }()
}

func (logger *CentralLogger) Close() {
	logger.wg.Wait()
	close(logger.road)
}
