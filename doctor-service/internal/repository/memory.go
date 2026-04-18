package repository

import (
	"errors"
	"sync"

	"ap2-assignment2/doctor-service/internal/model"
)

type DoctorRepository interface {
	Create(doctor *model.Doctor) error
	GetByID(id string) (*model.Doctor, error)
	GetAll() ([]*model.Doctor, error)
	ExistsByEmail(email string) bool
}

type InMemoryDoctorRepository struct {
	mu      sync.RWMutex
	doctors map[string]*model.Doctor
}

func NewInMemoryDoctorRepository() *InMemoryDoctorRepository {
	return &InMemoryDoctorRepository{
		doctors: make(map[string]*model.Doctor),
	}
}

func (r *InMemoryDoctorRepository) Create(doctor *model.Doctor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.doctors[doctor.ID] = doctor
	return nil
}

func (r *InMemoryDoctorRepository) GetByID(id string) (*model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.doctors[id]
	if !ok {
		return nil, errors.New("doctor not found")
	}
	return d, nil
}

func (r *InMemoryDoctorRepository) GetAll() ([]*model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*model.Doctor, 0, len(r.doctors))
	for _, d := range r.doctors {
		result = append(result, d)
	}
	return result, nil
}

func (r *InMemoryDoctorRepository) ExistsByEmail(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, d := range r.doctors {
		if d.Email == email {
			return true
		}
	}
	return false
}
