package coap

var debugEnable bool
var healthMonitorEnable bool

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
