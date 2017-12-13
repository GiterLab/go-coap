package coap

var debug bool

func init() {
	debug = false
}

// Debug swich debug mode
func Debug(mode bool) {
	debug = mode
}
