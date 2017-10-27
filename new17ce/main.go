package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
)

var url = "ws://admin.17ce.com:9002/router_manage?ts=%d&key=%s&r=%d"

var Username *string
var Dns *string
var IpAddress *string
// test git

const version = "3.0.10"

type Response struct {
	Act    string
	Contry string
	NodeId int
	Type   int
	UserId int
}

/*
样例输出
986947 107251 2
*/

func main() {
	Username = flag.String("-u", "2818889842@qq.com", "Username")
	Dns = flag.String("-dns", "127.0.0.1", "DNS Server Address")
	IpAddress = flag.String("-a", "192.168.1.1", "Machine IP Address")
	flag.Parse()

	// generate ts r and key
	ts := (int)(time.Now().Unix())
	r := (int)(rand.Int63n(time.Now().Unix()))
	key := md5.New()
	key.Write([]byte(reverseString(fmt.Sprintf("%s%s%s", strconv.Itoa(r), "8e2d642abac4bfbZxETNk0DL1EjN3RWC", strconv.Itoa(ts)))))
	finalKey := hex.EncodeToString(key.Sum(nil))
	finalUrl := fmt.Sprintf(url, ts, finalKey, r)
	origin := "http" + finalUrl[2:]

	Ws, err := websocket.Dial(finalUrl, "", origin)
	if nil != err {
		log.Fatal(err)
		return
	}

	data := make(map[string]string)
	data["Act"] = "Login"
	data["DnsIp"] = *Dns
	data["LocalIp"] = *IpAddress
	data["UUID"] = uuid.NewV4().String()[2:]
	data["Username"] = *Username
	data["Version"] = version
	jsonData, _ := json.Marshal(data)
	Ws.Write(jsonData)

	time.Sleep(time.Second)

	msg := make([]byte, 1024)
	n, err := Ws.Read(msg)

	// test
	// fmt.Println(string(msg))

	var v = &Response{}
	json.Unmarshal(msg[0:n], v)

	// write to console
	fmt.Println(v.UserId, v.NodeId, v.Type, v.Contry)

}

func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

/*
{
   "Act" : "LoginRt",
   "Coutry" : "China",
   "NodeId" : 107251,
   "Type" : 2,
   "UserId" : 986947
}
*/
