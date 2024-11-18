package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/dto"
	pb "github.com/vottundev/vottun-qubic-bridge-go/grpc/proto"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcClientConnection       *grpc.ClientConn
	grpcClientContext          context.Context
	grpcClientConnectionCancel context.CancelFunc
	serviceConnection          pb.RemoteServiceClient
)

func StartGrpcClientConnection(port uint16) error {
	var err error

	grpcClientConnection, err = grpc.NewClient(fmt.Sprintf("%s:%d", config.Config.Grpc.Server, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Failed connecting to gRpc bridge server: %+v", err)
		return errors.New("ERROR_GRPC_CONNECTION")
	}

	serviceConnection = pb.NewRemoteServiceClient(grpcClientConnection)
	grpcClientContext, grpcClientConnectionCancel = context.WithTimeout(context.Background(), 10*time.Second)

	log.Infof("client set at %s:%d", config.Config.Grpc.Server, port)

	return nil
}

func ProcessQubicOrder(order *dto.OrderReceivedDTO) {

	grpcOrder := &pb.ProcessOrderRequest{}
	grpcOrder.OrderID = order.OrderID
	grpcOrder.Amount = order.Amount
	grpcOrder.DestinationAccount = order.DestinationAccount
	grpcOrder.Memo = order.Memo
	grpcOrder.OriginAccount = order.OriginAccount
	grpcOrder.SourceChain = order.SourceChain

	result, err := serviceConnection.ProcessQubicOrder(grpcClientContext, grpcOrder)
	if err != nil {
		log.Errorf("Failed processing order: %+v", err)
	} else {
		log.Infof("Response from gRPC %+v", result.GetMessage())
	}
}
func StopGrprClientConnection() {
	grpcClientConnectionCancel()
	if grpcClientConnection != nil {
		grpcClientConnection.Close()
	}
}
