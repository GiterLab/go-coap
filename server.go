// Package coap provides a CoAP client and server.
package coap

import (
	"fmt"
	"log"
	"net"
	"time"
)

const maxPktLen = 1500

// Handler is a type that handles CoAP messages.
type Handler interface {
	// Handle the message and optionally return a response message.
	ServeCOAP(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message
}

type funcHandler func(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message

func (f funcHandler) ServeCOAP(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message {
	return f(l, a, m)
}

// FuncHandler builds a handler from a function.
func FuncHandler(f func(l *net.UDPConn, a *net.UDPAddr, m *Message) *Message) Handler {
	return funcHandler(f)
}

func handlePacket(l *net.UDPConn, data []byte, u *net.UDPAddr,
	rh Handler) {

	defer func() {
		data = nil

		// recover panic
		if err := recover(); err != nil {
			if debugEnable {
				TraceError("[coap] handle packet panic: %s", err)
			}
		}
	}()

	if debugEnable {
		tracePrintOut := true
		// health monitor for aliyun
		// Request:  RUOK
		// do not print out log for health monitor
		if healthMonitorEnable {
			if len(data) == 4 {
				if data[0] == 'R' && data[1] == 'U' && data[2] == 'O' && data[3] == 'K' {
					tracePrintOut = false
				}
			}
		}

		if tracePrintOut {
			TraceInfo("[coap] Remote: %v, Recv: %d, Bytes: %s", u, len(data), fmt.Sprintf("% X", data))
		}
	}

	// health monitor for aliyun
	// Request:  RUOK
	// Response: IMOK
	if healthMonitorEnable {
		if len(data) == 4 {
			if data[0] == 'R' && data[1] == 'U' && data[2] == 'O' && data[3] == 'K' {
				// Response IMOK
				l.WriteToUDP([]byte("IMOK"), u)
				return
			}
		}
	}

	msg, err := ParseMessage(data)
	if err != nil {
		log.Printf("Error parsing %v", err)
		return
	}

	rv := rh.ServeCOAP(l, u, &msg)
	if rv != nil {
		Transmit(l, u, *rv)
	}
}

// Transmit a message.
func Transmit(l *net.UDPConn, a *net.UDPAddr, m Message) error {
	d, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	if a == nil {
		_, err = l.Write(d)
	} else {
		_, err = l.WriteTo(d, a)
	}
	return err
}

// Receive a message.
func Receive(l *net.UDPConn, buf []byte) (Message, error) {
	l.SetReadDeadline(time.Now().Add(ResponseTimeout))

	nr, _, err := l.ReadFromUDP(buf)
	if err != nil {
		return Message{}, err
	}
	return ParseMessage(buf[:nr])
}

// ListenAndServe binds to the given address and serve requests forever.
func ListenAndServe(n, addr string, rh Handler) error {
	uaddr, err := net.ResolveUDPAddr(n, addr)
	if err != nil {
		return err
	}

	l, err := net.ListenUDP(n, uaddr)
	if err != nil {
		return err
	}

	return Serve(l, rh)
}

// Serve processes incoming UDP packets on the given listener, and processes
// these requests forever (or until the listener is closed).
func Serve(listener *net.UDPConn, rh Handler) error {
	buf := make([]byte, maxPktLen)
	for {
		nr, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			if debugEnable {
				TraceInfo("[coap] Serve ReadFromUDP error: %s", err)
			}
			continue
		}
		tmp := make([]byte, nr)
		copy(tmp, buf)
		go handlePacket(listener, tmp, addr, rh)
	}
}
