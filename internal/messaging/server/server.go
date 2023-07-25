package server

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/reflection"

	"github.com/hiroara/carbo/internal/messaging/pb"
)

type Server struct {
	pb.UnimplementedCommunicatorServer
	listener  net.Listener
	buffer    chan []byte
	token     string
	batch     []*pb.Message
	completed chan struct{}
	lock      *sync.Mutex
}

func New(lis net.Listener, buffer int) *Server {
	buf := make(chan []byte, buffer)
	return &Server{listener: lis, buffer: buf, completed: make(chan struct{}), lock: &sync.Mutex{}}
}

func (s *Server) FillBatch(ctx context.Context, req *pb.FillBatchRequest) (*pb.FillBatchResponse, error) {
	if s.token != req.Token {
		return nil, status.Error(codes.InvalidArgument, "Request token doesn't match.")
	}

	limit := sanitizeLimit(req.Limit)

	ok := s.fillBatch(ctx, limit)

	return &pb.FillBatchResponse{End: !ok}, nil
}

func (s *Server) fillBatch(ctx context.Context, limit int) bool {
	s.lock.Lock() // Lock the server until its batch is fulfilled

	bs, ok := <-s.buffer
	if !ok {
		defer s.lock.Unlock()
		// Buffer has already been closed.
		// No need to read the batch anymore.
		s.shutdown()
		return false
	}

	go func() {
		defer s.lock.Unlock()

		msgs := make([]*pb.Message, 0, limit)
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
	}()

	return true
}

func (s *Server) GetBatch(ctx context.Context, req *pb.GetBatchRequest) (*pb.GetBatchResponse, error) {
	s.lock.Lock() // Wait until filling batch is completed
	defer s.lock.Unlock()
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

func (s *Server) shutdown() {
	s.token = ""
	s.batch = make([]*pb.Message, 0)
	s.completed <- struct{}{}
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
	reflection.Register(srv)
	return srv.Serve(s.listener)
}

func sanitizeLimit(limit int32) int {
	lim := int(limit)
	if lim <= 0 {
		lim = 8
	}
	return lim
}
