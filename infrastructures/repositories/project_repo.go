package repositories

import "github.com/google/uuid"

type Repository interface {
	Create(data any) error
	Update(data any) error
	Delete(id uuid.UUID) error
}

type ProjectRepository interface {
	Repository
}
