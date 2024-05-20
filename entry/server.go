package entry

import (
	"log/slog"

	"github.com/KingDaemonX/ddd-template/domain/repository/interfaces"
	"github.com/KingDaemonX/ddd-template/domain/repository/interfaces/middlewares"
	"github.com/gin-gonic/gin"
)

type server struct {
	Router *gin.Engine
	mw     *middlewares.Middleware
	ph     *interfaces.Project
}

func NewServer(slg *slog.Logger, mw *middlewares.Middleware, ph *interfaces.Project) *server {
	return &server{
		Router: gin.Default(),
		mw:     mw,
		ph:     ph,
	}
}
