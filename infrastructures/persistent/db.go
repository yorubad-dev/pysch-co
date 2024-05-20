package persistent

import (
	"log/slog"
	"os"

	entity "github.com/KingDaemonX/ddd-template/domain/repository/domains/entities"
	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB

const (
	connecting_to_db string = "🌀 Attempting To Connect Application Database"
	db_conn_err      string = "🚨 Error Occur While Connecting To Database"
	conn_success     string = "😎 Database Connected SuccessFully"
)

type Repositories struct {
	Project repositories.ProjectRepository
	db      *gorm.DB
}

func NewRepository(slg *slog.Logger) (*Repositories, error) {
	dsn := os.Getenv("DATABASE_URL")

	slg.Info(connecting_to_db, "method", "persistent.NewRepository")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slg.Error(db_conn_err, "error context", err, "method", "persistent.NewRepository")
		return nil, err
	}

	slg.Info(conn_success, "method", "persistent.NewRepository")

	return &Repositories{
		Project: NewProjectRepo(db, slg),
		db:      db,
	}, nil

}

func (r *Repositories) Automigrate() error {
	return r.db.AutoMigrate(&entity.Project{})
}
