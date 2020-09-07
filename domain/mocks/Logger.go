package mocks

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"peterparada.com/online-bookmarks/domain"
)

func NewLogger() domain.Logger {
	logger := logrus.New()

	logger.SetOutput(ioutil.Discard)

	return logger
}
