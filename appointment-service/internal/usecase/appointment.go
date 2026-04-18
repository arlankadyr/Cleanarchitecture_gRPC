package usecase

import (
	"errors"
	"fmt"
	"time"

	"ap2-assignment2/appointment-service/internal/client"
	"ap2-assignment2/appointment-service/internal/model"
	"ap2-assignment2/appointment-service/internal/repository"

	"github.com/google/uuid"
)

var ErrDoctorNotFound = errors.New("doctor does not exist")

var ErrDoctorUnavailable = errors.New("doctor service unavailable")

type AppointmentUseCase struct {
	repo         repository.AppointmentRepository
	doctorClient client.DoctorClient
}

func NewAppointmentUseCase(repo repository.AppointmentRepository, dc client.DoctorClient) *AppointmentUseCase {
	return &AppointmentUseCase{repo: repo, doctorClient: dc}
}

type CreateInput struct {
	Title       string
	Description string
	DoctorID    string
}

func (uc *AppointmentUseCase) Create(input CreateInput) (*model.Appointment, error) {
	if input.Title == "" {
		return nil, errors.New("title is required")
	}
	if input.DoctorID == "" {
		return nil, errors.New("doctor_id is required")
	}

	exists, err := uc.doctorClient.DoctorExists(input.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDoctorUnavailable, err)
	}
	if !exists {
		return nil, fmt.Errorf("%w: id=%s", ErrDoctorNotFound, input.DoctorID)
	}

	now := time.Now()
	a := &model.Appointment{
		ID:          uuid.NewString(),
		Title:       input.Title,
		Description: input.Description,
		DoctorID:    input.DoctorID,
		Status:      model.StatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *AppointmentUseCase) GetByID(id string) (*model.Appointment, error) {
	return uc.repo.GetByID(id)
}

func (uc *AppointmentUseCase) GetAll() ([]*model.Appointment, error) {
	return uc.repo.GetAll()
}

func (uc *AppointmentUseCase) UpdateStatus(id string, newStatus model.Status) (*model.Appointment, error) {
	if !newStatus.IsValid() {
		return nil, fmt.Errorf("invalid status %q: must be one of new, in_progress, done", newStatus)
	}

	a, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if a.Status == model.StatusDone && newStatus == model.StatusNew {
		return nil, errors.New("cannot transition status from done back to new")
	}

	a.Status = newStatus
	a.UpdatedAt = time.Now()
	if err := uc.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}
