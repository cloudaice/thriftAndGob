package main

import (
	"log"
	"runtime"

	"thriftAndGob/translate"
	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

const (
	NetworkAddr = "0.0.0.0:9102"
)

type proxyTrans struct{}

func (pt *proxyTrans) Log(msgs []*translate.LogEntry) (translate.ResultCode, error) {
	_ = msgs
	return translate.ResultCode_OK, nil
}

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		log.Printf("NewTServerSocket error: %v\n", err)
		return
	}

	processor := translate.NewProxyTransProcessor(&proxyTrans{})
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	log.Println(server.Serve())
}
