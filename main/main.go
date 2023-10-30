package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"plugin"
	"sync"

	"github.com/lion187chen/NCan/ncandrv"
	"github.com/lion187chen/socketcan-go/canframe"
	"github.com/nats-io/nats.go"
)

// nats pub "ncan.rx" "{\"id\":32,\"data\":\"Af8DBAAA\",\"is_extended\":true}"

var nclient *NatsWraper
var toCanChan chan canframe.Frame

func main() {
	var driver string
	var server string
	var user string
	var passwd string
	var name string
	var rsubj string
	var tsubj string
	var config string

	flag.StringVar(&driver, "driver", "wsUCanA.so", "CAN Driver.")
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
	flag.Parse()

	/*plu, err := loadDriver(driver)
	if err != nil {
		panic(err)
	}*/

	nclient = new(NatsWraper).Init(server, user, passwd, name, onSubj)
	if config == "" {
		tsubj = name + "." + tsubj
		rsubj = name + "." + rsubj
	}
	nclient.Connect()
	nclient.Subscribe(rsubj)

	toCanChan = make(chan canframe.Frame, 16)

	var wg sync.WaitGroup
	// wg.Add(1)
	// go toNats(&wg, plu)
	wg.Add(1)
	go toCan(&wg)

	var frame canframe.Frame = canframe.Frame{
		ID:         0x20,
		Data:       []byte{0x01, 0xFF, 0x03, 0x4, 0x00, 0x00},
		IsExtended: true,
		IsRemote:   false,
		IsError:    false,
	}
	nclient.Publish(tsubj, frame)
	wg.Wait()
}

func onSubj(nm *nats.Msg, subj string, data []byte) {
	fmt.Println(subj)
	var frame canframe.Frame
	json.Unmarshal(data, &frame)

	toCanChan <- frame
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

func toNats(wg *sync.WaitGroup, can ncandrv.NCanDrvIf) {
	var nexit bool = true
	for nexit {
		frame := <-can.GetReadChannel()
		fmt.Println(frame)
	}
	wg.Done()
}

func toCan(wg *sync.WaitGroup) {
	var nexit bool = true
	for nexit {
		frame := <-toCanChan
		fmt.Println(frame)
	}
	wg.Done()
}
