package usecase

import (
	"errors"
	"fmt"

	"ap2-assignment2/doctor-service/internal/model"
	"ap2-assignment2/doctor-service/internal/repository"

	"github.com/google/uuid"
)

var ErrEmailAlreadyExists = errors.New("email already in use")

type DoctorUseCase struct {
	repo repository.DoctorRepository
}

func NewDoctorUseCase(repo repository.DoctorRepository) *DoctorUseCase {
	return &DoctorUseCase{repo: repo}
}

type CreateInput struct {
	FullName       string
	Specialization string
	Email          string
}

func (uc *DoctorUseCase) Create(input CreateInput) (*model.Doctor, error) {
	if input.FullName == "" {
		return nil, errors.New("full_name is required")
	}
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if uc.repo.ExistsByEmail(input.Email) {
		return nil, fmt.Errorf("%w: %s", ErrEmailAlreadyExists, input.Email)
	}
	doctor := &model.Doctor{
		ID:             uuid.NewString(),
		FullName:       input.FullName,
		Specialization: input.Specialization,
		Email:          input.Email,
	}
	if err := uc.repo.Create(doctor); err != nil {
		return nil, err
	}
	return doctor, nil
}

func (uc *DoctorUseCase) GetByID(id string) (*model.Doctor, error) {
	return uc.repo.GetByID(id)
}

func (uc *DoctorUseCase) GetAll() ([]*model.Doctor, error) {
	return uc.repo.GetAll()
}
