package ncandrv

import "github.com/lion187chen/socketcan-go/canframe"

type NCanDrvIf interface {
	New() (NCanDrvIf, error)
	WriteFrame(frame *canframe.Frame)
	GetReadChannel() <-chan canframe.Frame
}
