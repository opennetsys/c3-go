package rpc

import (
	"errors"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"

	context "golang.org/x/net/context"

	"github.com/c3systems/c3-go/core/p2p"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	"github.com/c3systems/c3-go/node/store"
	pb "github.com/c3systems/c3-go/rpc/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// ErrMethodNotSupported ...
	ErrMethodNotSupported = errors.New("method not supported")
)

// TODO: move to config
const port = ":5005"

// RPC ...
type RPC struct {
	mempool store.Interface
	p2p     *p2p.Service
}

// Config ...
type Config struct {
	Mempool store.Interface
	P2P     *p2p.Service
}

// Server ...
type Server struct {
	service *RPC
}

// New ...
func New(cfg *Config) *RPC {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	svc := &RPC{
		mempool: cfg.Mempool,
		p2p:     cfg.P2P,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterC3ServiceServer(grpcServer, &Server{
		service: svc,
	})
	reflection.Register(grpcServer)
	grpcServer.Serve(listen)

	log.Printf("[rpc] server running on port %s", port)

	return svc
}

// Send ...
func (s *Server) Send(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	method := strings.ToLower(r.Method)
	result, err := s.handleRequest(method, r)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.Response{
		Jsonrpc: r.Jsonrpc,
		Id:      r.Id,
		Result:  result,
	}, nil
}

// handleRequest ...
func (s *Server) handleRequest(method string, r *pb.Request) (*any.Any, error) {
	switch method {
	case "c3_ping":
		return ptypes.MarshalAny(s.service.ping())
	case "c3_latestblock":
		return ptypes.MarshalAny(s.service.latestBlock())
	case "c3_getblock":
		result, err := s.service.getBlock(r.Params)
		if err != nil {
			return ptypes.MarshalAny(&pb.ErrorResponse{
				Code:    400,
				Message: err.Error(),
			})
		}
		return ptypes.MarshalAny(result)
	case "c3_getstateblock":
		return ptypes.MarshalAny(s.service.getStateblock(r.Params))
	default:
		return nil, ErrMethodNotSupported
	}
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
