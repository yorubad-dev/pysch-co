package main

import (
	"log/slog"
	"os"

	"github.com/KingDaemonX/ddd-template/domain/repository/entry"
	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/auth"
	"github.com/KingDaemonX/ddd-template/domain/repository/infrastructures/persistent"
	"github.com/KingDaemonX/ddd-template/domain/repository/interfaces"
	"github.com/KingDaemonX/ddd-template/domain/repository/interfaces/middlewares"
	"github.com/joho/godotenv"
)

var slg *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func init() {
	if err := godotenv.Load(); err != nil {
		slg.Error("loading .env file failed", "error context", err, "function", "main.init")
		return
	}
}

func main() {
	slg.Info("initiating server up", "function", "main.main")

	if err := initiateServerUp(); err != nil {
		slg.Error("error initiating server up", "error context", err, "function", "main.main")
		return
	}

	slg.Info("server up and running", "port", os.Getenv("SERVER_PORT"), "function", "main.main")
}

func initiateServerUp() error {
	allow_db_setup := os.Getenv("ALLOW_DATABASE_SETUP") == "true"

	service, err := persistent.NewRepository(slg)
	if err != nil {
		slg.Error("error initiating service persistent", "error context", err, "function", "main.initiateServerUp")
		return err
	}

	if allow_db_setup {
		if err := service.Automigrate(); err != nil {
			slg.Error("❎ Error Occur While Migrating Database Schema", "error context", err, "function", "main.initiateServerUp")
			return err
		}
	}

	token := auth.NewToken(slg)
	redis := auth.NewRedis(slg)

	project := interfaces.NewProject(slg, service.Project, token, redis)

	server := entry.NewServer(slg, middlewares.NewMiddleware(slg), &project)

	server.Routes()

	return server.Router.Run(":" + os.Getenv("SERVER_PORT"))
}
