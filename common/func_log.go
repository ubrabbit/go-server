package common

import (
	golog "github.com/davyxu/golog"
)

var (
	g_LoggerInfo  *golog.Logger = nil
	g_LoggerWarn  *golog.Logger = nil
	g_LoggerError *golog.Logger = nil
)

func NewLog(name string, filepath string) *golog.Logger {
	obj := golog.New(name)
	if name != "" {
		golog.SetOutputLogger(name, filepath)
	}
	return obj
}

func LogFatal(v ...interface{}) {
	if g_LoggerError == nil {
		g_LoggerError = golog.New("error")
	}
	g_LoggerError.Fatalln(v...)
}

func LogError(v ...interface{}) {
	if g_LoggerError == nil {
		g_LoggerError = golog.New("error")
	}
	g_LoggerError.Errorln(v...)
}

func LogInfo(v ...interface{}) {
	if g_LoggerInfo == nil {
		g_LoggerInfo = golog.New("info")
	}
	g_LoggerInfo.Infoln(v...)
}

func LogWarning(v ...interface{}) {
	if g_LoggerWarn == nil {
		g_LoggerWarn = golog.New("warning")
	}
	g_LoggerWarn.Warnln(v...)
}
