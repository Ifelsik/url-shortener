package logger

import (
	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/sirupsen/logrus"
)

type LogrusLogWrap struct {
	log *logrus.Entry
}

func initDefault() *LogrusLogWrap {
	l := logrus.New()
	l.SetFormatter(
		&logrus.TextFormatter{
			ForceColors:            true,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		},
	)
	l.ReportCaller = true

	return &LogrusLogWrap{
		log: logrus.NewEntry(l),
	}
}

func NewLogrusLogWrap(conf *logger.LoggerConfig) *LogrusLogWrap {
	// default logger
	if conf == nil {
		return initDefault()
	}

	// TODO: add logger parameters from configuration
	// Temporary default is used
	return initDefault()
}

func (l *LogrusLogWrap) WithFields(fields logger.LoggerFields) *LogrusLogWrap {
	return &LogrusLogWrap{log: l.log.WithFields(logrus.Fields(fields))}
}

func (l *LogrusLogWrap) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *LogrusLogWrap) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *LogrusLogWrap) Warningf(format string, args ...any) {
	l.log.Warningf(format, args...)
}

func (l *LogrusLogWrap) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *LogrusLogWrap) Fatalf(format string, args ...any) {
	l.log.Fatalf(format, args...)
}
