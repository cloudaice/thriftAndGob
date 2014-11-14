package main

import (
	"log"
	"time"

	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
	"thriftAndGob/translate"
)

const (
	HOSTPORT = "127.0.0.1:9102"
)

var (
	logData = `42.62.41.64 - - [26/Jan/2014:06:59:59 +0800] "GET /mshopapi/index.php/v1/authorize/sso?client_id=180188831013&callback=http%3A%2F%2Fm.demo.com%2Findex.html%23ac%3Daccount%26op%3Dindex HTTP/1.0" 302 0 "-" "Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166" "10.100.2.11/127.0.0.1:9999" "0.001/0.001/302/888888888"`
)

func main() {
	tSocket, err := thrift.NewTSocket(HOSTPORT)
	if err != nil {
		log.Printf("NewTSocket error: %v\n", err)
	}

	transport := thrift.NewTFramedTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := translate.NewProxyTransClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		log.Printf("transport.Open error: %v\n", err)
	}
	defer transport.Close()

	var count = 0
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println(count)
			count = 0
		default:
			var buf []*translate.LogEntry
			for i := 0; i < 200; i++ {
				entry := &translate.LogEntry{
					Hostname: "",
					Message:  logData,
				}
				buf = append(buf, entry)
			}
			code, err := client.Log(buf)
			if err != nil {
				log.Printf("Log error: %v, code: %d\n", err, code)
			}
			count += 200
			time.Sleep(time.Millisecond)
		}
	}
}
