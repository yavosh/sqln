package grpc

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/yavosh/sqln"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
	pb "github.com/yavosh/sqln/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is a gRPC server.
type Server struct {
	pb.UnimplementedQueryServer
	server *grpc.Server
	port   int
	log    sqln.Logger
	db     *sql.DB
}

// NewServer create the http server with dependencies
func NewServer(port int, db *sql.DB) *Server {
	s := &Server{
		server: grpc.NewServer(),
		port:   port,
	}

	// Inject dependencies
	s.log = log.New(os.Stdout, "grpc ", log.LstdFlags)
	s.db = db

	pb.RegisterQueryServer(s.server, s)
	// Enable server reflection service.
	reflection.Register(s.server) // Ability to query server to describe service contracts.

	return s
}

// Start will start the grpc server
func (s *Server) Start() error {
	log.Printf("Starting grpc server @ %d ", s.port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("grpc serve: %w", err)
	}

	go func() {
		if err := s.server.Serve(ln); err != nil {
			log.Fatalf("Error starting http server %v", err)
		}
	}()

	log.Printf("Started grpc server @ %d ", s.port)
	return nil
}

// Stop will stop the grpc server
func (s *Server) Stop() error {
	if s.server == nil {
		return errors.New("can't stop, server not running")
	}
	log.Printf("Stopping grpc server")

	s.server.GracefulStop()
	log.Printf("Stopped grpc server")
	return nil
}
