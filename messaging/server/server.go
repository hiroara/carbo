package server

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/hiroara/carbo/messaging/message"
	"github.com/hiroara/carbo/pb"
)

type Server struct {
	pb.UnimplementedCommunicatorServer
	listener  net.Listener
	buffer    chan message.Message
	completed chan struct{}
}

func New(lis net.Listener, buffer int) *Server {
	buf := make(chan message.Message, buffer)
	return &Server{listener: lis, buffer: buf, completed: make(chan struct{})}
}

func (s *Server) BatchPull(ctx context.Context, req *pb.BatchPullRequest) (*pb.BatchPullResponse, error) {
	limit := sanitizeLimit(req.Limit)
	closed := true
	msgs := make([]*pb.Message, 0, limit)
	for el := range s.buffer {
		dat, err := el.MarshalBinary()
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, &pb.Message{Value: dat})
		if len(msgs) == limit {
			closed = false
			break
		}
	}
	if closed {
		s.completed <- struct{}{}
	}
	return &pb.BatchPullResponse{Messages: msgs, Closed: closed}, nil
}

func (s *Server) Feed(msg message.Message) {
	s.buffer <- msg
}

func (s *Server) Run(ctx context.Context) error {
	srv := grpc.NewServer()

	go func() {
		<-ctx.Done() // Ensure all inputs has been fed
		close(s.buffer)
		<-s.completed // Wait until every message has been consumed
		srv.GracefulStop()
	}()

	pb.RegisterCommunicatorServer(srv, s)
	return srv.Serve(s.listener)
}

func sanitizeLimit(limit int32) int {
	lim := int(limit)
	if lim <= 0 {
		lim = 8
	}
	return lim
}
