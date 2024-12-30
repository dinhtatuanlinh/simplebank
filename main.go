package main

import (
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"simplebank/api"
	"simplebank/db"
	sqlc "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/util"

	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	db.RunMigration(&config.MigrationUrl, &config.DBSource)

	store := sqlc.NewStore(conn)
	go runGinServer(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store sqlc.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot listen:", err)
	}

	log.Printf("start gPRC server at %s", config.GRPCServerAddress)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGinServer(config util.Config, store sqlc.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
