package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

var uid int
var nid int
var ntype int

var events chan string

var Ws *websocket.Conn

var url = "ws://admin.17ce.com:9002/router_manage?ts=%d&key=%s&r=%d"

func main() {
	uid = *flag.Int("-uid", 986947, "User ID")
	nid = *flag.Int("-nid", 107317, "Node ID")
	ntype = *flag.Int("-t", 2, "Node Type")
	flag.Parse()
	fmt.Println("UID:", uid)
	fmt.Println("NID:", nid)
	fmt.Println("Type:", ntype)

	// generate ts r and key
	ts := (int)(time.Now().Unix())
	r := (int)(rand.Int63n(time.Now().Unix()))
	key := md5.New()
	key.Write([]byte(reverseString(fmt.Sprintf("%s%s%s", strconv.Itoa(r), "8e2d642abac4bfbZxETNk0DL1EjN3RWC", strconv.Itoa(ts)))))
	finalKey := hex.EncodeToString(key.Sum(nil))
	finalUrl := fmt.Sprintf(url, ts, finalKey, r)
	origin := "http" + finalUrl[2:]

	// init
	var err error
	Ws, err = websocket.Dial(finalUrl, "", origin)
	if nil != err {
		fmt.Println(err)
		return
	}

	// ping
	go Ping()
	go GetTask()

	for {
		msg := make([]byte, 1024)
		n, err := Ws.Read(msg)
		if nil != err {
			fmt.Println(err)
			break
		}
		fmt.Println("Read", n, "bytes :", string(msg))
	}

}

func Ping() {
	for {
		_, err := Ws.Write([]byte("{\"Act\": \"Ping\"}"))
		if nil != err {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 30)
	}
}

func GetTask() {
	for {
		str := "{\"Act\": \"GetTask\", \"UserId\": " + strconv.Itoa(uid) + ", \"NodeId\": " + strconv.Itoa(nid) + "}"
		_, err := Ws.Write([]byte(str))
		if nil != err {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 30)
	}
}

func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}
