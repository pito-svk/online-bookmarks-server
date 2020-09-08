package repository

import "github.com/sirupsen/logrus"

type LoggerImpl struct {
	*logrus.Logger
}

func (logger *LoggerImpl) LogJSONMap(data map[string]interface{}, logLevel logrus.Level, args ...interface{}) {
	fields := logrus.Fields(logrus.Fields(data))
	entry := logger.WithFields(fields)

	entry.Log(logLevel, fields, args)
}

func (logger *LoggerImpl) Trace(args ...interface{}) {
	if len(args) > 0 {
		if jsonMap, ok := args[0].(map[string]interface{}); ok {
			logger.LogJSONMap(jsonMap, logrus.TraceLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.TraceLevel, args...)
}

func (logger *LoggerImpl) Info(args ...interface{}) {
	if len(args) > 0 {
		if jsonMap, ok := args[0].(map[string]interface{}); ok {
			logger.LogJSONMap(jsonMap, logrus.InfoLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.InfoLevel, args...)
}

func (logger *LoggerImpl) Warn(args ...interface{}) {
	if len(args) > 0 {
		if jsonMap, ok := args[0].(map[string]interface{}); ok {
			logger.LogJSONMap(jsonMap, logrus.WarnLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.WarnLevel, args...)
}

func (logger *LoggerImpl) Error(args ...interface{}) {
	if len(args) > 0 {
		if jsonMap, ok := args[0].(map[string]interface{}); ok {
			logger.LogJSONMap(jsonMap, logrus.ErrorLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.ErrorLevel, args...)
}
