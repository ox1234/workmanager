package workmanager

import (
	"io"

	"github.com/sirupsen/logrus"
)

var log = func() *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logrus.NewEntry(logger)
}()

func SetLogger(l *logrus.Entry) {
	log = l
}
