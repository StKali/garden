package rpc

import (
	"fmt"
	"net"

	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/pb"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var setting = util.GetSetting()

type Server struct {
	store db.Querier
	maker token.Maker
	pb.UnimplementedAuthServiceServer
}

func NewServer(store db.Querier, maker token.Maker) (*Server, error) {
	server := &Server{store: store, maker: maker}
	return server, nil
}

func (s *Server) Start(address string) {
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", address)
	util.CheckError(fmt.Sprintf("failed to listen address: %s, err: ", address), err)
	err = grpcServer.Serve(listener)
	util.CheckError("failed to create grpc server, err: %s", err)
}
