package repository

import "github.com/sirupsen/logrus"

type LoggerImpl struct {
	*logrus.Logger
}

func (logger *LoggerImpl) Trace(args ...interface{}) {
	if len(args) > 0 {
		if mapData, ok := args[0].(map[string]interface{}); ok {
			entry := logger.WithFields(logrus.Fields(mapData))

			entry.Log(logrus.TraceLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.TraceLevel, args...)
}

func (logger *LoggerImpl) Info(args ...interface{}) {
	if len(args) > 0 {
		if mapData, ok := args[0].(map[string]interface{}); ok {
			entry := logger.WithFields(logrus.Fields(mapData))

			entry.Log(logrus.InfoLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.InfoLevel, args...)
}

func (logger *LoggerImpl) Warn(args ...interface{}) {
	if len(args) > 0 {
		if mapData, ok := args[0].(map[string]interface{}); ok {
			entry := logger.WithFields(logrus.Fields(mapData))

			entry.Log(logrus.WarnLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.WarnLevel, args...)
}

func (logger *LoggerImpl) Error(args ...interface{}) {
	if len(args) > 0 {
		if mapData, ok := args[0].(map[string]interface{}); ok {
			entry := logger.WithFields(logrus.Fields(mapData))

			entry.Log(logrus.ErrorLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.ErrorLevel, args...)
}
