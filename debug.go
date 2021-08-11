package coap

import (
	"log"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var debugEnable bool
var healthMonitorEnable bool

type TraceFunc func(format string, level int, v ...interface{})

var UserTrace TraceFunc = nil

func init() {
	debugEnable = false
	healthMonitorEnable = false
}

// Debug Enable debug
func Debug(enable bool) {
	debugEnable = enable
}

// HealthMonitor Enable health monitor
func HealthMonitor(enable bool) {
	healthMonitorEnable = enable
}

// SetUserDebug 配置其他日志输出
func SetUserDebug(f TraceFunc) {
	UserTrace = f
}

// TraceInfo 调试信息日志
func TraceInfo(format string, v ...interface{}) {
	if debugEnable {
		if UserTrace != nil {
			UserTrace(format, LevelInformational, v...)
		} else {
			log.Printf(format, v...)
		}
	}
}

// TraceError 错误日志
func TraceError(format string, v ...interface{}) {
	if debugEnable {
		if UserTrace != nil {
			UserTrace(format, LevelError, v...)
		} else {
			log.Printf(format, v...)
		}
	}
}
