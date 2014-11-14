package main

import (
	"log"
	"time"

	"thriftAndGob/trans"
)

var (
	logData = `42.62.41.64 - - [26/Jan/2014:06:59:59 +0800] "GET /mshopapi/index.php/v1/authorize/sso?client_id=180888088813&callback=http%3A%2F%2Fm.demo.com%2Findex.html%23ac%3Daccount%26op%3Dindex HTTP/1.0" 302 0 "-" "Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166" "10.100.2.11/127.0.0.1:9999" "0.001/0.001/302/888888888"`
)

func main() {
	pdr, err := trans.NewProducer(":9102", false)
	if err != nil {
		log.Printf("NewProducer: %v\n", err)
	}
	defer pdr.Close()
	var count = 0
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println(count)
			count = 0
		default:

			var buf []*trans.LogEntry
			for i := 0; i < 200; i++ {
				entry := &trans.LogEntry{
					Hostname: "",
					Message:  logData,
				}
				buf = append(buf, entry)
			}
			err = pdr.SendArray(buf)
			if err != nil {
				log.Printf("SendArray error: %v\n", err)
			}
			count += 200
			time.Sleep(time.Millisecond)
		}
	}
}
