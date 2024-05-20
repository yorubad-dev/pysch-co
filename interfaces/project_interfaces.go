package interfaces

import (
	"log/slog"
	"net/http"

	"github.com/KingDaemonX/ddd-template/domain/repository/applications"
	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/auth"
	"github.com/KingDaemonX/ddd-template/domain/repository/interfaces/response"
	"github.com/gin-gonic/gin"
)

type Project struct {
	pa  applications.ProjectAppInterface
	slg *slog.Logger
	tk  auth.TokenInterface
	rd  auth.RedisInterface
}

func NewProject(slg *slog.Logger, pa applications.ProjectAppInterface, tk auth.TokenInterface, rd auth.RedisInterface) Project {
	return Project{
		pa:  pa,
		slg: slg,
		tk:  tk,
		rd:  rd,
	}
}

// write serializer method, pass to database and interact with the main file
func (pi *Project) Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := response.NewResponse(http.StatusOK, "server up", "Hello World from your favourite health checker")
		pi.slg.Info("health check", "method", "(interfaces.Project).Health")
		c.JSON(http.StatusOK, &response)
	}
}
