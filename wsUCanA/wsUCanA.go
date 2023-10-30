package main

import (
	"flag"

	"github.com/lion187chen/NCan/ncandrv"
	"github.com/lion187chen/socketcan-go/canframe"
	wsucana "github.com/lion187chen/waveshare-usb-can-a-go"
)

type wsUCanA struct {
	port   string
	rate   uint
	ext    bool
	repeat bool
	ucan   *wsucana.UsbCanA
}

func New() (ncandrv.NCanDrvIf, error) {
	my := new(wsUCanA)
	flag.StringVar(&my.port, "port", "/dev/ttyUSB0", "WaveShare USB-CAN-A's virtual serial port name.")
	flag.UintVar(&my.rate, "rate", 100, "CAN bit rate 5,10,20,50,100,125,200,250,400,500,800,1000.")
	flag.BoolVar(&my.ext, "ext", true, "Use extended frame.")
	flag.BoolVar(&my.repeat, "repeat", true, "Auto repeat.")
	my.ucan = new(wsucana.UsbCanA)
	return my, nil
}

func (my *wsUCanA) Delete() error {
	return nil
}

func (my *wsUCanA) Open(name string) error {
	err := my.ucan.Open(my.port, 16)
	if err != nil {
		return err
	}

	var rate wsucana.BiterateType
	var ext wsucana.CanFrameType = wsucana.FRAME_CFG_CAN_FRAME_STD
	var repeat wsucana.RepeatType = wsucana.FRAME_CFG_REPEAT_NO
	switch my.rate {
	case 5:
		rate = wsucana.FRAME_CFG_BIT_RATE_5K
	case 10:
		rate = wsucana.FRAME_CFG_BIT_RATE_10K
	case 20:
		rate = wsucana.FRAME_CFG_BIT_RATE_20K
	case 50:
		rate = wsucana.FRAME_CFG_BIT_RATE_50K
	case 125:
		rate = wsucana.FRAME_CFG_BIT_RATE_125K
	case 200:
		rate = wsucana.FRAME_CFG_BIT_RATE_200K
	case 250:
		rate = wsucana.FRAME_CFG_BIT_RATE_250K
	case 400:
		rate = wsucana.FRAME_CFG_BIT_RATE_400K
	case 500:
		rate = wsucana.FRAME_CFG_BIT_RATE_500K
	case 800:
		rate = wsucana.FRAME_CFG_BIT_RATE_800K
	case 1000:
		rate = wsucana.FRAME_CFG_BIT_RATE_1M
	default:
		rate = wsucana.FRAME_CFG_BIT_RATE_100K
	}
	if my.ext {
		ext = wsucana.FRAME_CFG_CAN_FRAME_EXT
	}
	if my.repeat {
		repeat = wsucana.FRAME_CFG_REPEAT_AUTO
	}
	my.ucan.Config(rate, ext, wsucana.FRAME_CFG_WRK_MOD_NORMAL, repeat)
	return nil
}

func (my *wsUCanA) Close() error {
	my.ucan.Close()
	return nil
}

func (my *wsUCanA) WriteFrame(frame *canframe.Frame) {
	my.ucan.WriteFrame(frame)
}

func (my *wsUCanA) GetReadChannel() <-chan canframe.Frame {
	return my.ucan.GetReadChannel()
}
