package grpc

import (
	"context"
	"errors"

	"ap2-assignment2/appointment-service/internal/model"
	"ap2-assignment2/appointment-service/internal/usecase"
	pb "ap2-assignment2/appointment-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppointmentHandler struct {
	pb.UnimplementedAppointmentServiceServer
	uc *usecase.AppointmentUseCase
}

func NewAppointmentHandler(uc *usecase.AppointmentUseCase) *AppointmentHandler {
	return &AppointmentHandler{uc: uc}
}

func domainToProto(a *model.Appointment) *pb.AppointmentResponse {
	return &pb.AppointmentResponse{
		Id:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		DoctorId:    a.DoctorID,
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *AppointmentHandler) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.AppointmentResponse, error) {
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.DoctorId == "" {
		return nil, status.Error(codes.InvalidArgument, "doctor_id is required")
	}

	a, err := h.uc.Create(usecase.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		DoctorID:    req.DoctorId,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrDoctorUnavailable):
			return nil, status.Errorf(codes.Unavailable, "doctor service unreachable: %v", err)
		case errors.Is(err, usecase.ErrDoctorNotFound):
			return nil, status.Errorf(codes.FailedPrecondition, "doctor does not exist: doctor_id=%s", req.DoctorId)
		default:
			return nil, status.Errorf(codes.Internal, "failed to create appointment: %v", err)
		}
	}

	return domainToProto(a), nil
}

func (h *AppointmentHandler) GetAppointment(ctx context.Context, req *pb.GetAppointmentRequest) (*pb.AppointmentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	a, err := h.uc.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "appointment not found: id=%s", req.Id)
	}

	return domainToProto(a), nil
}

func (h *AppointmentHandler) ListAppointments(ctx context.Context, req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appointments, err := h.uc.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list appointments: %v", err)
	}

	resp := make([]*pb.AppointmentResponse, 0, len(appointments))
	for _, a := range appointments {
		resp = append(resp, domainToProto(a))
	}

	return &pb.ListAppointmentsResponse{Appointments: resp}, nil
}

func (h *AppointmentHandler) UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.AppointmentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	a, err := h.uc.UpdateStatus(req.Id, model.Status(req.Status))
	if err != nil {
		switch {
		case err.Error() == "appointment not found":
			return nil, status.Errorf(codes.NotFound, "appointment not found: id=%s", req.Id)
		default:
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
	}

	return domainToProto(a), nil
}
