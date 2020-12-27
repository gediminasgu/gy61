//go:generate mockgen -source=types.go -package=gy61 -destination types_mock.go

package gy61

import "github.com/tarm/serial"

type Serial interface {
	OpenSerial(c *serial.Config) error
	Write(b []byte) (int, error)
	Read(b []byte) (int, error)
	Close() error
}
