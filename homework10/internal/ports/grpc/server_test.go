package grpc

import (
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	mocks "homework10/internal/ports/grpc/app_mocks"
	"homework10/internal/ports/grpc/pb"
	"log"
	"net"
	"testing"
)

type grpcServerTestSuite struct {
	suite.Suite

	app    *mocks.App
	server Server

	grpcServer *grpc.Server
	conn       *grpc.ClientConn
	client     pb.AdServiceClient
}

func grpcServerTestSuiteInit(s *grpcServerTestSuite) {
	// for testing grpc methods
	s.app = new(mocks.App)
	s.server = Server{
		AdServiceServer: nil,
		App:             s.app,
	}

	// for testing NewGRPCServer
	s.grpcServer = NewGRPCServer(s.app)
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		if serveErr := s.grpcServer.Serve(listener); serveErr != nil {
			log.Fatalf("failed to serve: %v", serveErr)
		}
	}()
	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	s.conn = conn
	s.client = pb.NewAdServiceClient(conn)
}

func (s *grpcServerTestSuite) SetupSuite() {
	grpcServerTestSuiteInit(s)
}

func (s *grpcServerTestSuite) TearDownSuite() {
	s.grpcServer.Stop()
	err := s.conn.Close()
	if err != nil {
		log.Fatalf("failde to close conn: %v", err)
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(grpcServerTestSuite))
}
