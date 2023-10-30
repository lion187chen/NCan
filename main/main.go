package main

import (
	"plugin"

	"github.com/lion187chen/NCan/ncandrv"
)

func main() {}

func loadPlugin(name string) (ncandrv.NCanDrvIf, error) {
	pfile, err := plugin.Open(name)
	if err != nil {
		panic(err)
	}

	nfun, err := pfile.Lookup("New")
	if err != nil {
		panic(err)
	}

	plu, err := nfun.(func() (ncandrv.NCanDrvIf, error))()
	return plu, err
}
