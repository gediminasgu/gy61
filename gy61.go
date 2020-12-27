package gy61

import (
	"encoding/binary"
)

type GY61DataReceivedFn func(x, y, z float32)
type ErrorFn func(error)

type GY61 struct {
	serialPort       Serial
	gyroReceivedFn   GY61DataReceivedFn
	accReceivedFn    GY61DataReceivedFn
	anglesReceivedFn GY61DataReceivedFn
	onError          ErrorFn
}

// NewGY61 creates a new instance for gyro data reading
func NewGY61(serialPort Serial, gyroReceivedFn, accReceivedFn, anglesReceivedFn GY61DataReceivedFn, onError ErrorFn) *GY61 {
	return &GY61{
		serialPort:       serialPort,
		gyroReceivedFn:   gyroReceivedFn,
		accReceivedFn:    accReceivedFn,
		anglesReceivedFn: anglesReceivedFn,
		onError:          onError,
	}
}

// ReadAsync starts goroutine to read data in the background
func (s *GY61) ReadAsync() {
	go s.read()
}

func (s *GY61) read() {
	hd := make([]byte, 1)
	buf := make([]byte, 9)

	for {
		n, err := s.serialPort.Read(hd)
		//fmt.Printf("%s: SERIAL1: %d; %s\n", time.Now().Format("15:04:05.000"), n, hex.EncodeToString(hd[:n]))
		if err != nil && s.onError != nil {
			s.onError(err)
			continue
		}

		if n == 1 && hd[0] == 0x55 {
			n, err = s.serialPort.Read(hd)
			if err != nil && s.onError != nil {
				s.onError(err)
				continue
			}
			if n == 1 && hd[0] >= 0x51 && hd[0] <= 0x53 {
				//fmt.Printf("%s: SERIAL2: %d; %s\n", time.Now().Format("15:04:05.000"), n, hex.EncodeToString(hd[:n]))
				i := 0
				for {
					n, err = s.serialPort.Read(buf[i:9])
					if err != nil && s.onError != nil {
						s.onError(err)
						continue
					}
					i += n
					if i >= 9 {
						break
					}
				}
				//fmt.Printf("%s: GYRO Buf: %d; %s\n", time.Now().Format("15:04:05.000"), n, hex.EncodeToString(buf))
				switch hd[0] {
				case 0x51:
					if s.gyroReceivedFn != nil {
						s.gyroReceivedFn(s.parseG(buf))
					}
				case 0x52:
					if s.accReceivedFn != nil {
						s.accReceivedFn(s.parseV(buf))
					}
				case 0x53:
					if s.anglesReceivedFn != nil {
						s.anglesReceivedFn(s.parseR(buf))
					}
				}
			}
		}
	}
}

func (s *GY61) parseG(buf []byte) (float32, float32, float32) {
	x := float32(binary.LittleEndian.Uint16(buf[0:2])) / 32768.0 * 16
	if x > 16 {
		x = -(32 - x)
	}
	y := float32(binary.LittleEndian.Uint16(buf[2:4])) / 32768.0 * 16
	if y > 16 {
		y = -(32 - y)
	}
	z := float32(binary.LittleEndian.Uint16(buf[4:6])) / 32768.0 * 16
	if z > 16 {
		z = -(32 - z)
	}
	return x, y, z
}

func (s *GY61) parseV(buf []byte) (float32, float32, float32) {
	x := float32(binary.LittleEndian.Uint16(buf[0:2])) / 32768.0 * 2000
	if x > 2000 {
		x = -(4000 - x)
	}
	y := float32(binary.LittleEndian.Uint16(buf[2:4])) / 32768.0 * 2000
	if y > 2000 {
		y = -(4000 - y)
	}
	z := float32(binary.LittleEndian.Uint16(buf[4:6])) / 32768.0 * 2000
	if z > 2000 {
		z = -(4000 - z)
	}
	return x, y, z
}

func (s *GY61) parseR(buf []byte) (float32, float32, float32) {
	x := float32(binary.LittleEndian.Uint16(buf[0:2])) / 32768.0 * 180
	if x > 180 {
		x = -(360 - x)
	}
	y := float32(binary.LittleEndian.Uint16(buf[2:4])) / 32768.0 * 180
	if y > 180 {
		y = -(360 - y)
	}
	z := float32(binary.LittleEndian.Uint16(buf[4:6])) / 32768.0 * 180
	if z > 180 {
		z = -(360 - z)
	}
	return x, y, z
}
