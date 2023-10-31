package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"plugin"

	"github.com/lion187chen/NCan/ncandrv"
)

// nats pub "ncan0" "{\"id\":32,\"data\":\"Af8DBAAA\",\"is_extended\":true}"
// nats sub "ncan1"
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
	var cfile string

	flag.StringVar(&driver, "driver", "./wsUCanA.so", "CAN Driver.")
	flag.StringVar(&server, "server", "127.0.0.1:6666", "NATS server and port.")
	flag.StringVar(&user, "user", "ncan", "NATS login user name.")
	flag.StringVar(&passwd, "passwd", "000000", "NATS login password.")
	flag.StringVar(&name, "name", "ncan", "Application' name.")
	flag.StringVar(&tsubj, "tsubj", "ncan1", "Transmit through the \"<name>.<tsubj>\" subject.")
	flag.StringVar(&rsubj, "rsubj", "ncan0", "Receive through the \"<name>.<rsubj>\" subject.")
	flag.StringVar(&cfile, "config", "./config.json", "Use config as a config file. Default: not use.")

	flag.Usage = func() {
		fmt.Println("NCan version v0.0.1")
		flag.PrintDefaults()
	}

	drv, err := loadDriver(driver)
	if err != nil {
		panic(err)
	}
	flag.Parse()

	var cfg *config
	var driverCfg []byte
	if cfile != "" {
		cfg = new(config).init(cfile)

		if cfg.Driver != nil {
			driverCfg, err = json.Marshal(cfg.Driver)
			if err != nil {
				panic(err)
			}
			if string(driverCfg) == "{}" {
				driverCfg = nil
			}
		}
	}

	err = drv.Open("", driverCfg)
	if err != nil {
		panic(err)
	}

	var post *postman
	if cfile == "" {
		post = new(postman).initWithParam(server, user, passwd, name, tsubj, rsubj, drv)
	} else {
		post = new(postman).initWithConfig(cfg, drv)
	}

	post.joint()
	post.del()
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
