package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/grpc/pb"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/application/usecase"
	"github.com/ajvideira/imersao-fullstack-fullcycle/codepix/infrastructure/repository"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//StartGrpcServer create a new gRPC server
func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	pixUseCase := usecase.PixKeyUseCase{PixKeyRepository: &pixRepository}
	pixService := PixGrpcService {PixUseCase: pixUseCase}
	pb.RegisterPixKeyServiceServer(grpcServer, &pixService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start gRPC server", err)
	}
	log.Printf("gRPC server has been started on port %d", port)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot serve gRPC server", err)
	}

}