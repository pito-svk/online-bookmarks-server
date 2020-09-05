package mocks

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(ioutil.Discard)

	return logger
}
