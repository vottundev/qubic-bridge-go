package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/vottundev/vottun-qubic-bridge-go/grpc/proto"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var server *grpc.Server

type grpcServer struct {
	pb.UnimplementedRemoteServiceServer
}

func (s *grpcServer) ProcessQubicOrder(ctx context.Context, in *pb.ProcessOrderRequest) (*pb.ProcessOrderResponse, error) {
	p, _ := peer.FromContext(ctx)

	log.Infof("Request received from addr: %s", p.Addr.String())
	return &pb.ProcessOrderResponse{Message: true}, nil
}

func StartGrpcServer(port uint16) error {
	listener, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Errorf("Failed starting gRPC server: %+v", err)
		return errors.New("ERROR_GRPC_SERVER_START")
	}

	server = grpc.NewServer()
	pb.RegisterRemoteServiceServer(server, &grpcServer{})
	log.Infof("gRPC server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Errorf("Failed to serve gRPC: %+v", err)
		return errors.New("ERROR_GRPC_SERVER_SERVE")

	}
	return nil
}

func StopGrpcServer() {
	if server != nil {
		server.Stop()
	}
}
