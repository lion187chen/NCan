package main

import (
	"flag"
	"fmt"
	"plugin"

	"github.com/lion187chen/NCan/ncandrv"
)

// nats pub "ncan.rx" "{\"id\":32,\"data\":\"Af8DBAAA\",\"is_extended\":true}"
// cansend can0 00000123#12345678
// candump can0

func main() {
	var driver string
	var server string
	var user string
	var passwd string
	var name string
	var rsubj string
	var tsubj string
	var config string

	flag.StringVar(&driver, "driver", "../wsUCanA/wsUCanA.so", "CAN Driver.")
	flag.StringVar(&server, "server", "192.168.2.15:6666", "NATS server and port.")
	flag.StringVar(&user, "user", "ncan", "NATS login user name.")
	flag.StringVar(&passwd, "passwd", "000000", "NATS login password.")
	flag.StringVar(&name, "name", "ncan", "Application' name.")
	flag.StringVar(&tsubj, "tsubj", "tx", "Transmit through the \"<name>.<tsubj>\" subject.")
	flag.StringVar(&rsubj, "rsubj", "rx", "Receive through the \"<name>.<rsubj>\" subject.")
	flag.StringVar(&config, "config", "", "Use config as a config file. Default: not use.")

	flag.Usage = func() {
		fmt.Println("NCan version v0.0.1")
		flag.PrintDefaults()
	}

	drv, err := loadDriver(driver)
	if err != nil {
		panic(err)
	}
	flag.Parse()

	drv.Open("")

	postman := new(postman).init(server, user, passwd, name, tsubj, rsubj, drv)
	postman.joint()
	postman.del()
}

func loadDriver(name string) (ncandrv.NCanDrvIf, error) {
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
