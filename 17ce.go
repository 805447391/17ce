package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

const version = "3.0.10"

var url = "ws://admin.17ce.com:9002/router_manage?ts=1508402592&key=c1155f5fce2ff3a8b5fe62b8a173d185&r=1112346968"

// var _url = &url.URL{Scheme: "ws", Host: "admin.17ce.com:9002", Path: "/router_manage", RawQuery: "ts=1508384758&key=694e49b404e59e50aef7a3d4b4f384d7&r=1407773123"}
var origin = "http" + url[2:]

type Core struct {
	Ws *websocket.Conn
	// account information
	Email   string
	LocalIp string
	DnsIp   string
	UUID    string

	// private
	Message    chan []byte
	DaemonFlag bool
}

func (c *Core) Login() {
	data := make(map[string]string)
	data["Act"] = "Login"
	data["DnsIp"] = c.DnsIp
	data["LocalIp"] = c.LocalIp
	data["UUID"] = c.UUID
	data["Username"] = c.Email
	data["Version"] = version
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))
	c.Ws.Write(jsonData)
}

func (c *Core) Init() {
	fmt.Println("init...")

	ws, err := websocket.Dial(url, "", origin)
	c.Ws = ws
	CE(err)
	fmt.Println("init ok")
	// start goroutine
	// go c.Daemon()
}

func (c *Core) SendPing() {
	c.Ws.Write("{\"Act\": \"Ping\"}")
}

func (c *Core) ReInit() {
	if nil == c.Ws {
		fmt.Println("Reinit")
		c.Init()
	}
}

func (c *Core) OnMessage(msg []byte) {
	// str := string(msg)
	fmt.Println(string(msg))
}

func (c *Core) Finish() {
	c.Ws.Close()
}

func (c *Core) Daemon() {
	if c.DaemonFlag {
		return
	} else {
		c.DaemonFlag = true
	}
	var msgPool bytes.Buffer
	for {
		// msgPool.Reset()
		if nil != c.Ws {
			for {
				msg := make([]byte, 1024*8)
				n, err := c.Ws.Read(msg)
				fmt.Println(n)
				msgPool.Write(msg[0:n])
				if n < 1024*8 {
					break
				}

				if nil != err {
					fmt.Println(err)
					c.Ws.Close()
					c.Ws = nil
					continue
				}
			}
			// c.Message <- msgPool
			c.OnMessage(msgPool.Bytes())
		} else {
			c.ReInit()
		}
	}
}
