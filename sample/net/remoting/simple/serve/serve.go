package main

import (
	"fmt"
	"log"

	"github.com/boltmq/boltmq/net/core"
	"github.com/boltmq/boltmq/net/remoting"
	"github.com/boltmq/common/protocol"
	"github.com/boltmq/common/protocol/header/namesrv"
)

var (
	remotingServer remoting.RemotingServer
)

type GetTopicStatisInfoProcessor struct {
}

func (processor *GetTopicStatisInfoProcessor) ProcessRequest(ctx core.Context,
	request *protocol.RemotingCommand) (*protocol.RemotingCommand, error) {
	fmt.Printf("GetTopicStatisInfoProcessor %d %d\n", request.Code, request.Opaque)

	topicStatisInfoRequestHeader := &namesrv.GetTopicStatisInfoRequestHeader{}
	err := request.DecodeCommandCustomHeader(topicStatisInfoRequestHeader)
	if err != nil {
		return nil, err
	}
	fmt.Printf("DeprotocolCommandCustomHeader %v\n", topicStatisInfoRequestHeader)

	response := protocol.CreateResponseCommand(protocol.SUCCESS, "success")
	response.Opaque = request.Opaque

	return response, nil
}

type OtherProcessor struct {
}

func (processor *OtherProcessor) ProcessRequest(ctx core.Context,
	request *protocol.RemotingCommand) (*protocol.RemotingCommand, error) {
	fmt.Printf("OtherProcessor %d %d\n", request.Code, request.Opaque)

	response := protocol.CreateResponseCommand(protocol.SUCCESS, "success")
	response.Opaque = request.Opaque

	return response, nil
}

type ServerContextEventListener struct {
}

func (listener *ServerContextEventListener) OnContextActive(ctx core.Context) {
	log.Printf("one connection active: localAddr[%s] remoteAddr[%s]\n", ctx.LocalAddr(), ctx.RemoteAddr())
}

func (listener *ServerContextEventListener) OnContextConnect(ctx core.Context) {
	log.Printf("one connection create: localAddr[%s] remoteAddr[%s]\n", ctx.LocalAddr(), ctx.RemoteAddr())
}

func (listener *ServerContextEventListener) OnContextClosed(ctx core.Context) {
	log.Printf("one connection close: localAddr[%s] remoteAddr[%s]\n", ctx.LocalAddr(), ctx.RemoteAddr())
}

func (listener *ServerContextEventListener) OnContextError(ctx core.Context, err error) {
	log.Printf("one connection error: localAddr[%s] remoteAddr[%s]\n", ctx.LocalAddr(), ctx.RemoteAddr())
}

func (listener *ServerContextEventListener) OnContextIdle(ctx core.Context) {
	log.Printf("one connection idle: localAddr[%s] remoteAddr[%s]\n", ctx.LocalAddr(), ctx.RemoteAddr())
}

func main() {
	initServer()
	remotingServer.Start()
}

func initServer() {
	remotingServer = remoting.NewNMRemotingServer("0.0.0.0", 10911)
	remotingServer.RegisterProcessor(protocol.HEART_BEAT, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.SEND_MESSAGE_V2, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.GET_TOPIC_STATS_INFO, &GetTopicStatisInfoProcessor{})
	remotingServer.RegisterProcessor(protocol.GET_MAX_OFFSET, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.QUERY_CONSUMER_OFFSET, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.PULL_MESSAGE, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.UPDATE_CONSUMER_OFFSET, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.GET_CONSUMER_LIST_BY_GROUP, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.GET_ROUTEINTO_BY_TOPIC, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.UPDATE_AND_CREATE_TOPIC, &OtherProcessor{})
	remotingServer.RegisterProcessor(protocol.GET_KV_CONFIG, &OtherProcessor{})
	remotingServer.SetContextEventListener(&ServerContextEventListener{})
}
