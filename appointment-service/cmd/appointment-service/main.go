package main

import (
	"log"
	"net"
	"os"

	"ap2-assignment2/appointment-service/internal/client"
	"ap2-assignment2/appointment-service/internal/repository"
	grpchandler "ap2-assignment2/appointment-service/internal/transport/grpc"
	"ap2-assignment2/appointment-service/internal/usecase"
	pb "ap2-assignment2/appointment-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	doctorAddr := os.Getenv("DOCTOR_SERVICE_ADDR")
	if doctorAddr == "" {
		doctorAddr = "localhost:50051"
	}

	doctorClient, conn, err := client.NewGRPCDoctorClient(doctorAddr)
	if err != nil {
		log.Fatalf("cannot connect to Doctor Service at %s: %v", doctorAddr, err)
	}
	defer conn.Close()

	log.Printf("Connected to Doctor Service at %s", doctorAddr)

	repo := repository.NewInMemoryAppointmentRepository()
	uc := usecase.NewAppointmentUseCase(repo, doctorClient)
	handler := grpchandler.NewAppointmentHandler(uc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on :50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAppointmentServiceServer(grpcServer, handler)
	
	reflection.Register(grpcServer)

	log.Println("Appointment Service listening on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Appointment Service terminated: %v", err)
	}
}
