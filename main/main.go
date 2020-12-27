package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gediminasgu/gy61"

	"github.com/tarm/serial"
)

const (
	baudRate = 115200
)

var (
	s        gy61.Serial
	portName string
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("please provide serial port name")
	}
	portName = os.Args[1]

	var err error
	s, err = gy61.NewSerial(&serial.Config{
		Name: portName,
		Baud: baudRate,
	})
	if err != nil {
		log.Fatalf("error creating serial: %v", err)
	}

	kalman := gy61.NewGY61(s, nil, nil, anglesReceived, onError)
	kalman.ReadAsync()

	time.Sleep(time.Hour)
}

func anglesReceived(x, y, z float32) {
	fmt.Printf("%7.2f %7.2f %7.2f\n", x, y, z)
}

func onError(err error) {
	s.Close()
	time.Sleep(time.Second)
	err = s.OpenSerial(&serial.Config{
		Name: portName,
		Baud: baudRate,
	})
	if err != nil {
		fmt.Printf("failed to open serial: %v\n", err)
	}
}
