package server

import (
	"context"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hiroara/carbo/pb"
)

type Server struct {
	pb.UnimplementedCommunicatorServer
	listener  net.Listener
	buffer    chan []byte
	token     string
	batch     []*pb.Message
	completed chan struct{}
}

func New(lis net.Listener, buffer int) *Server {
	buf := make(chan []byte, buffer)
	return &Server{listener: lis, buffer: buf, completed: make(chan struct{})}
}

func (s *Server) FillBatch(ctx context.Context, req *pb.FillBatchRequest) (*pb.FillBatchResponse, error) {
	if s.token != req.Token {
		return nil, status.Error(codes.InvalidArgument, "Request token doesn't match.")
	}

	limit := sanitizeLimit(req.Limit)

	msgs := make([]*pb.Message, 0, limit)

	bs, ok := <-s.buffer
	if !ok {
		s.token = ""
		s.batch = msgs

		// Buffer has already been closed.
		// No need to read the batch anymore.
		s.completed <- struct{}{}

		return &pb.FillBatchResponse{End: true}, nil
	}

	for {
		msgs = append(msgs, &pb.Message{Value: bs})
		if len(msgs) == limit {
			break
		}

		if bs, ok = <-s.buffer; !ok {
			break
		}
	}

	s.token = strconv.FormatInt(time.Now().UnixNano(), 10)
	s.batch = msgs

	return &pb.FillBatchResponse{End: false}, nil
}

func (s *Server) GetBatch(ctx context.Context, req *pb.GetBatchRequest) (*pb.GetBatchResponse, error) {
	return &pb.GetBatchResponse{Token: s.token, Messages: s.batch}, nil
}

func (s *Server) Feed(ctx context.Context, bs []byte) {
	select {
	case s.buffer <- bs:
	case <-ctx.Done():
	}
}

func (s *Server) Close() error {
	close(s.buffer)
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	srv := grpc.NewServer()

	go func() {
		select {
		case <-ctx.Done():
		case <-s.completed:
		}
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
