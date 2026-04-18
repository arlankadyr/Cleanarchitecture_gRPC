package main

import (
	"log"
	"net"

	"ap2-assignment2/doctor-service/internal/repository"
	grpchandler "ap2-assignment2/doctor-service/internal/transport/grpc"
	"ap2-assignment2/doctor-service/internal/usecase"
	pb "ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	repo := repository.NewInMemoryDoctorRepository()
	uc := usecase.NewDoctorUseCase(repo)
	handler := grpchandler.NewDoctorHandler(uc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on :50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDoctorServiceServer(grpcServer, handler)
	
	reflection.Register(grpcServer)

	log.Println("Doctor Service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Doctor Service terminated: %v", err)
	}
}
