package goczmq

/*
#cgo !windows pkg-config: libczmq
#cgo windows CFLAGS: -I/usr/local/include
#cgo windows LDFLAGS: -L/usr/local/lib -lczmq
#include "czmq.h"

zactor_t *Beacon_new () { zactor_t *beacon = zactor_new(zbeacon, NULL); return beacon; }
*/
import "C"

import (
	"strconv"
	"unsafe"
)

type Beacon struct {
	zactor_t *C.struct__zactor_t
}

// NewBeacon creates a new Beacon instance.
func NewBeacon() *Beacon {
	z := &Beacon{}
	z.zactor_t = C.Beacon_new()
	return z
}

// Verbose sets the beacon to log information to stdout.
func (z *Beacon) Verbose() error {
	rc := C.zstr_send(unsafe.Pointer(z.zactor_t), C.CString("VERBOSE"))
	if rc == -1 {
		return ErrActorCmd
	}

	return nil
}

// Configure accepts a port number and configures the beacon, returning an address
func (z *Beacon) Configure(port int) (string, error) {
	rc := C.zstr_sendm(unsafe.Pointer(z.zactor_t), C.CString("CONFIGURE"))
	if rc == -1 {
		return "", ErrActorCmd
	}

	rc = C.zstr_send(unsafe.Pointer(z.zactor_t), C.CString(strconv.Itoa(port)))
	if rc == -1 {
		return "", ErrActorCmd
	}

	Chostname := C.zstr_recv(unsafe.Pointer(z.zactor_t))
	hostname := C.GoString(Chostname)
	return hostname, nil
}

// Publish publishes an announcement at an interval
func (z *Beacon) Publish(announcement string, interval int) error {
	rc := C.zstr_sendm(unsafe.Pointer(z.zactor_t), C.CString("PUBLISH"))
	if rc == -1 {
		return ErrActorCmd
	}

	rc = C.zstr_sendm(unsafe.Pointer(z.zactor_t), C.CString(announcement))
	if rc == -1 {
		return ErrActorCmd
	}

	rc = C.zstr_send(unsafe.Pointer(z.zactor_t), C.CString(strconv.Itoa(interval)))
	if rc == -1 {
		return ErrActorCmd
	}

	return nil
}

// Subscribe subscribes to beacons matching the filter
func (z *Beacon) Subscribe(filter string) error {
	rc := C.zstr_sendm(unsafe.Pointer(z.zactor_t), C.CString("SUBSCRIBE"))
	if rc == -1 {
		return ErrActorCmd
	}

	rc = C.zstr_send(unsafe.Pointer(z.zactor_t), C.CString(filter))
	if rc == -1 {
		return ErrActorCmd
	}

	return nil
}

// Recv waits for the specific timeout in milliseconds to receive a beacon
func (z *Beacon) Recv(timeout int) string {
	C.zsock_set_rcvtimeo(unsafe.Pointer(z.zactor_t), C.int(timeout))
	msg := C.zstr_recv(unsafe.Pointer(z.zactor_t))
	return C.GoString(msg)
}

// Destroy destroys the beacon.
func (z *Beacon) Destroy() {
	C.zactor_destroy(&z.zactor_t)
}