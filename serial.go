package gy61

import "github.com/tarm/serial"

type serialPort struct {
	port *serial.Port
}

func NewSerial(c *serial.Config) (Serial, error) {
	s := &serialPort{}
	err := s.OpenSerial(c)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *serialPort) OpenSerial(c *serial.Config) error {
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	s.port = port
	return nil
}

func (s *serialPort) Write(b []byte) (int, error) {
	return s.port.Write(b)
}

func (s *serialPort) Read(b []byte) (int, error) {
	return s.port.Read(b)
}

func (s *serialPort) Close() error {
	return s.port.Close()
}
