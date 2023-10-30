package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lion187chen/NCan/ncandrv"
	"github.com/lion187chen/socketcan-go/canframe"
	"github.com/nats-io/nats.go"
)

type postman struct {
	nclient   *NatsWraper
	tsub      string
	rsub      *nats.Subscription
	toCanChan chan *nats.Msg
	wg        sync.WaitGroup
}

func (my *postman) init(server, user, passwd, myName, tsubj, rsubj string, can ncandrv.NCanDrvIf) *postman {
	my.nclient = new(NatsWraper).Init(server, user, passwd, myName)

	my.nclient.Connect()
	my.toCanChan = make(chan *nats.Msg, 16)
	// if config == "" {
	my.tsub = myName + "." + tsubj
	my.rsub = my.nclient.Subscribe(myName+"."+rsubj, my.toCanChan)
	// }

	my.wg.Add(1)
	go my.toNats(can)
	my.wg.Add(1)
	go my.toCan(can)
	return my
}

func (my *postman) toNats(can ncandrv.NCanDrvIf) {
	var nexit bool = true
	fmt.Println("toNats is running.")
	for nexit {
		frame := <-can.GetReadChannel()
		fmt.Println("toNats:", frame)
		my.nclient.Publish(my.tsub, frame)
	}
	my.wg.Done()
}

func (my *postman) toCan(can ncandrv.NCanDrvIf) {
	var nexit bool = true
	fmt.Println("toCan is running.")
	for nexit {
		msg := <-my.toCanChan
		var frame canframe.Frame
		err := json.Unmarshal(msg.Data, &frame)
		if err != nil {
			continue
		}

		fmt.Println("toCan:", frame)
		can.WriteFrame(&frame)
	}
	my.wg.Done()
}

func (my *postman) joint() {
	my.wg.Wait()
}

func (my *postman) del() {
	my.rsub.Unsubscribe()
	// TODO: close all channel in my other packages.
	close(my.toCanChan)
}
