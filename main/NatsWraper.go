package main

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Publish()：单纯的发布 Data 信息。
// PublishRequest()：带 Reply 的 Publish()。
// PublishMsg()：带 Header 的 PublishRequest()。
// Request()：同步版本的 PublishRequest()。
// RequestMsg()：高级版本的带 Header 的 Request()。
// Subscribe()：订阅感兴趣信息，收到信息后将调用回调函数。
// SubscribeSync()：同步版本的 Subscribe()。
// Respond()：通过 Reply 主题 Publish()。
// RespondMsg()：通过 Reply 主题 PublishMsg()。

// NATS 客户端.
type NatsWraper struct {
	// requests gtasks.Queue
	NConn *nats.Conn
	NatsState

	server string
	user   string
	passwd string
	myName string

	subjHandler OnSubjectHandler
}

type NatsState struct {
	connected bool
	sync.RWMutex
}

type OnSubjectHandler func(nm *nats.Msg, subj string, data []byte)

func (n *NatsState) Write(state bool) {
	n.RWMutex.Lock()
	n.connected = state
	defer n.RWMutex.Unlock()
}

func (n *NatsState) Read() bool {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()
	return n.connected
}

// NatsWraper 构造方法.
// 返回 NatsWraper 对象指针.
func (n *NatsWraper) Init(server, user, passwd, myName string, hdl OnSubjectHandler) *NatsWraper {
	n.NatsState.Write(false)

	if myName == "" {
		n.myName = "nclient"
	} else {
		n.myName = myName
	}
	if server == "" {
		n.server = "localhost:9119"
	} else {
		n.server = server
	}
	n.user = user
	n.passwd = passwd
	n.subjHandler = hdl

	return n
}

func (n *NatsWraper) Connect() {
	var e error

	for {
		// AllowReconnect:     true,
		// MaxReconnect:       DefaultMaxReconnect,
		// ReconnectWait:      DefaultReconnectWait,
		// ReconnectJitter:    DefaultReconnectJitter,
		// ReconnectJitterTLS: DefaultReconnectJitterTLS,
		// Timeout:            DefaultTimeout,
		// PingInterval:       DefaultPingInterval,
		// MaxPingsOut:        DefaultMaxPingOut,
		// SubChanLen:         DefaultMaxChanLen,
		// ReconnectBufSize:   1MB,
		// DrainTimeout:       DefaultDrainTimeout,
		n.NConn, e = nats.Connect(
			n.server,
			nats.UserInfo(n.user, n.passwd),
			nats.ClosedHandler(n.onClose),
			nats.DisconnectHandler(n.onDisconnect),
			nats.ErrorHandler(n.onErr),
			nats.ReconnectHandler(n.onReconnect),
			nats.ReconnectBufSize(1024*1024),
			nats.Name(n.myName))
		if e == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	n.NatsState.Write(true)
}

// 关闭连接.
func (n *NatsWraper) Close() {
	if !n.NatsState.Read() {
		return
	}
	n.NatsState.Write(false)
	n.NConn.Close()
}

func (n *NatsWraper) IsConnected() bool {
	return n.NatsState.Read()
}

// 参数 subj 为接收使用的主题
func (n *NatsWraper) Subscribe(subj string) {
	_, e := n.NConn.Subscribe(subj, n.onSubject)
	if e != nil {
		panic(e)
	}
}

func (n *NatsWraper) onSubject(m *nats.Msg) {
	n.subjHandler(m, m.Subject, m.Data)
}

// 参数 subj 为发布信息使用的主题.
// 参数 data 为发布内容
func (n *NatsWraper) Publish(sub string, data interface{}) error {
	if !n.NatsState.Read() {
		return errors.New("wrong nats state")
	}

	msg, e := json.Marshal(data)
	if e != nil {
		return e
	}

	e = n.NConn.Publish(sub, msg)
	return e
}

func (n *NatsWraper) onErr(*nats.Conn, *nats.Subscription, error) {
	// Do nothing
}

// 重连超时则断开连接
func (n *NatsWraper) onClose(nc *nats.Conn) {
	n.NatsState.Write(false)
	n.Connect()
}

// 响应断开连接事件.
func (n *NatsWraper) onDisconnect(nc *nats.Conn) {
	n.NatsState.Write(false)
}

// 响应重连接事件.
func (n *NatsWraper) onReconnect(nc *nats.Conn) {
	n.NatsState.Write(true)
}
