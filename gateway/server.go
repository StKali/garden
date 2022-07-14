package gateway

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/pb"
	"github.com/stkali/garden/rpc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	store db.Querier
	maker token.Maker
}

func NewServer(store db.Querier, maker token.Maker) (*Server, error) {
	server := Server{
		store: store,
		maker: maker,
	}
	return &server, nil
}

func (s *Server) Start(address string) {
	server, err := rpc.NewServer(s.store, s.maker)
	util.CheckError("failed to create rpc server", err)

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterAuthServiceHandlerServer(context.Background(), grpcMux, server)
	util.CheckError("failed to register service handle", err)

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", address)
	util.CheckError("failed to create listener", err)

	err = http.Serve(listener, mux)
	util.CheckError("failed to start grpc-gateway server", err)
}
