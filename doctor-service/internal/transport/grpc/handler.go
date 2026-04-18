package grpc

import (
	"context"
	"errors"

	"ap2-assignment2/doctor-service/internal/usecase"
	pb "ap2-assignment2/doctor-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DoctorHandler struct {
	pb.UnimplementedDoctorServiceServer
	uc *usecase.DoctorUseCase
}

func NewDoctorHandler(uc *usecase.DoctorUseCase) *DoctorHandler {
	return &DoctorHandler{uc: uc}
}

func (h *DoctorHandler) CreateDoctor(ctx context.Context, req *pb.CreateDoctorRequest) (*pb.DoctorResponse, error) {
	if req.FullName == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	doctor, err := h.uc.Create(usecase.CreateInput{
		FullName:       req.FullName,
		Specialization: req.Specialization,
		Email:          req.Email,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "email already in use: %s", req.Email)
		}
		return nil, status.Errorf(codes.Internal, "failed to create doctor: %v", err)
	}

	return &pb.DoctorResponse{
		Id:             doctor.ID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		Email:          doctor.Email,
	}, nil
}

func (h *DoctorHandler) GetDoctor(ctx context.Context, req *pb.GetDoctorRequest) (*pb.DoctorResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	doctor, err := h.uc.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "doctor not found: id=%s", req.Id)
	}

	return &pb.DoctorResponse{
		Id:             doctor.ID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		Email:          doctor.Email,
	}, nil
}

func (h *DoctorHandler) ListDoctors(ctx context.Context, req *pb.ListDoctorsRequest) (*pb.ListDoctorsResponse, error) {
	doctors, err := h.uc.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list doctors: %v", err)
	}

	resp := make([]*pb.DoctorResponse, 0, len(doctors))
	for _, d := range doctors {
		resp = append(resp, &pb.DoctorResponse{
			Id:             d.ID,
			FullName:       d.FullName,
			Specialization: d.Specialization,
			Email:          d.Email,
		})
	}

	return &pb.ListDoctorsResponse{Doctors: resp}, nil
}
