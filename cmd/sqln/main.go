package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/peterbourgon/ff/v3"
	"github.com/yavosh/sqln/grpc"
)

func main() {
	var (
		grpcPort   int
		dataSource string
	)

	fs := flag.NewFlagSet("server", flag.ExitOnError)
	fs.IntVar(&grpcPort, "grpc-port", 5051, "listen port for the grpc server")
	fs.StringVar(&dataSource, "data-source", ":memory:", "data source for sqlite")

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatalf("flag set: %v", err)
	}

	log.Printf("Using sqlite datasource %q", dataSource)
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpcPort, db)
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Error starting server %v", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := grpcServer.Stop(); err != nil {
		log.Fatalf("Error stopping server %v", err)
	}
}
