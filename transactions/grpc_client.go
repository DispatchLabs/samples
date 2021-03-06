package main

import (
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"golang.org/x/net/context"

	"github.com/dispatchlabs/disgo/dapos/proto"
	"github.com/dispatchlabs/disgo/commons/types"
	"time"
	"github.com/dispatchlabs/disgo/commons/utils"
	"fmt"
)

func SendGrpcTransaction(tx *types.Transaction, grpcEndpoint *types.Endpoint, address string) (*types.Gossip, error) {
	node := types.Node{GrpcEndpoint: grpcEndpoint, Type: types.TypeDelegate, Address: address}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.GrpcEndpoint.Host, node.GrpcEndpoint.Port), grpc.WithInsecure())
	if err != nil {
		utils.Fatal(fmt.Sprintf("cannot dial seed [host=%s, port=%d]",  node.GrpcEndpoint.Host,  node.GrpcEndpoint.Port), err)
		log.Fatalf("cannot dial server: %v", err)
	}
	client := proto.NewDAPoSGrpcClient(conn)
	gossip := types.NewGossip(*tx)

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	// Remote gossip.
	response, err := client.GossipGrpc(contextWithTimeout, &proto.Request{Payload: gossip.String()})
	if err != nil {
		return nil, err
	}

	remoteGossip, err := types.ToGossipFromJson([]byte(response.Payload))
	if err != nil {
		return nil, err
	}

	return remoteGossip, err
}
