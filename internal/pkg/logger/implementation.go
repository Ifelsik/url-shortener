package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

type LogrusLogWrap struct {
	log  *logrus.Entry
	conf *LoggerConfig
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
	// TODO: implement configuration for logger
	l.SetLevel(logrus.DebugLevel)
	// Need to set custom report caller
	// because usage of default Report caller
	// will show information about wrapping function
	l.ReportCaller = false

	return &LogrusLogWrap{
		log: logrus.NewEntry(l),
		conf: &LoggerConfig{
			ShowCaller: true,
		},
	}
}

func NewLogrusLogWrap(conf *LoggerConfig) *LogrusLogWrap {
	// default logger
	if conf == nil {
		return initDefault()
	}

	// TODO: add logger parameters from configuration
	// Temporary default is used
	return initDefault()
}

func (l *LogrusLogWrap) WithFields(fields LoggerFields) Logger {
	return &LogrusLogWrap{
		log:  l.log.WithFields(logrus.Fields(fields)),
		conf: l.conf,
	}
}

func (l *LogrusLogWrap) Infof(format string, args ...any) {
	logger := l.log
	if l.conf.ShowCaller {
		logger = l.withCaller()
	}
	logger.Infof(format, args...)
}

func (l *LogrusLogWrap) Debugf(format string, args ...any) {
	logger := l.log
	if l.conf.ShowCaller {
		logger = l.withCaller()
	}
	logger.Debugf(format, args...)
}

func (l *LogrusLogWrap) Warningf(format string, args ...any) {
	logger := l.log
	if l.conf.ShowCaller {
		logger = l.withCaller()
	}
	logger.Warningf(format, args...)
}

func (l *LogrusLogWrap) Errorf(format string, args ...any) {
	logger := l.log
	if l.conf.ShowCaller {
		logger = l.withCaller()
	}
	logger.Errorf(format, args...)
}

func (l *LogrusLogWrap) Fatalf(format string, args ...any) {
	logger := l.log
	if l.conf.ShowCaller {
		logger = l.withCaller()
	}
	logger.Fatalf(format, args...)
}

func (l *LogrusLogWrap) withCaller() *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(2); ok {
		funcName := runtime.FuncForPC(pc).Name()

		return l.log.WithFields(logrus.Fields{
			"func":  funcName,
			"place": fmt.Sprintf("%s:%d", file, line),
		})
	}
	l.log.Warningln("LogrusLogWrap: unable to get caller info")
	
	return l.log
}
