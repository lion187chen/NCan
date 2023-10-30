package main

import (
	"flag"

	"github.com/lion187chen/NCan/ncandrv"
	"github.com/lion187chen/socketcan-go/canframe"
)

type wsUCanA struct{}

func (my *wsUCanA) New() (ncandrv.NCanDrvIf, error) {
	var port string
	flag.StringVar(&port, "port", "COM6", "WaveShare USB-CAN-A's virtual serial port name.")
	return nil, nil
}

func (my *wsUCanA) WriteFrame(frame *canframe.Frame) {}

func (my *wsUCanA) GetReadChannel() <-chan canframe.Frame {
	return nil
}
