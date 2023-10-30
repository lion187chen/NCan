package ncandrv

import "github.com/lion187chen/socketcan-go/canframe"

// New() (NCanDrvIf, error)

type NCanDrvIf interface {
	Delete() error
	Open(name string) error
	Close() error
	WriteFrame(frame *canframe.Frame)
	GetReadChannel() <-chan canframe.Frame
}
