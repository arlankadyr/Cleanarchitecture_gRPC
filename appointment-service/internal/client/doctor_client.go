package client

import (
	"context"
	"fmt"
	"time"

	doctorpb "ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type DoctorClient interface {
	DoctorExists(doctorID string) (bool, error)
}

type GRPCDoctorClient struct {
	stub doctorpb.DoctorServiceClient
}

func NewGRPCDoctorClient(addr string) (*GRPCDoctorClient, *grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext( 
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Doctor Service at %s: %w", addr, err)
	}

	return &GRPCDoctorClient{stub: doctorpb.NewDoctorServiceClient(conn)}, conn, nil
}

func (c *GRPCDoctorClient) DoctorExists(doctorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.stub.GetDoctor(ctx, &doctorpb.GetDoctorRequest{Id: doctorID})
	if err == nil {
		return true, nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return false, fmt.Errorf("doctor service unavailable: %w", err)
	}

	switch st.Code() {
	case codes.NotFound:
		return false, nil
	case codes.Unavailable:
		return false, fmt.Errorf("doctor service unavailable: %s", st.Message())
	default:
		return false, fmt.Errorf("doctor service error (%s): %s", st.Code(), st.Message())
	}
}
