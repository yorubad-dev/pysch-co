package applications

import (
	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/repositories"
	"github.com/google/uuid"
)

type ProjectApp struct {
	pr repositories.ProjectRepository
}

type ProjectAppInterface interface {
	// implement thesame thing as the repository because it serves as a passage between the interface and domain
	Create(data any) error
	Update(data any) error
	Delete(id uuid.UUID) error
}

// create passage method here
var _ repositories.ProjectRepository = &ProjectApp{}

func (pa *ProjectApp) Create(data any) error {
	return pa.pr.Create(data)
}
func (pa *ProjectApp) Update(data any) error {
	return pa.pr.Update(data)
}
func (pa *ProjectApp) Delete(id uuid.UUID) error {
	return pa.pr.Delete(id)
}
