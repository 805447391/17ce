package main

import (
	"fmt"
)

func main() {
	fmt.Println("17ce client by 805447391@qq.com")
	ins := &Core{
		DnsIp:   "127.0.0.1",
		Email:   "2818889842@qq.com",
		LocalIp: "192.168.1.1",
		UUID:    "2818889842@qq.comb3a27f-d7c7-49f3-86cb-72a275c72f7f",
	}
	ins.Init()
	ins.Login()
	ins.Daemon()

}

func CE(err error) {
	if nil != err {
		panic(err)
	}
}
