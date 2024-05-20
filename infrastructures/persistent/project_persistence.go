package persistent

import (
	"log/slog"

	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	db  *gorm.DB
	slg *slog.Logger
}

func NewProjectRepo(db *gorm.DB, slg *slog.Logger) *ProjectRepo {
	return &ProjectRepo{
		db:  db,
		slg: slg,
	}
}

var _ repositories.ProjectRepository = &ProjectRepo{}

func (r *ProjectRepo) Create(data any) error {
	return r.db.Create(data).Error
}

func (r *ProjectRepo) Update(data any) error {
	return r.db.Save(data).Error
}

func (r *ProjectRepo) Delete(id uuid.UUID) error {
	var data any

	return r.db.Delete(data, id).Error
}
