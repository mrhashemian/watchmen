package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	logging "watchmen/log"

	"watchmen/config"
	"watchmen/controller"
	"watchmen/database"
	"watchmen/repository"
)

var serveCMD = &cobra.Command{
	Use:   "serve",
	Short: "serve API",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	baseAPIDB := database.InitBaseAPIDB()
	defer database.CloseDB(baseAPIDB)

	userRepo := repository.NewUserRepository(baseAPIDB)

	e := echo.New()
	e.Use(logging.EchoMiddleware())

	e.GET("/", indexHandler)

	v1 := e.Group("/v1")
	v1Users := v1.Group("/users")

	v1Users.POST("/sign-up", controller.SignUp(userRepo))

	go func() {
		if err := e.Start(config.C.Server.Address); err != nil {
			e.Logger.Fatal(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGTERM,
		syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("error in shutdown: %v", err)
	}
}

func indexHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, asciiArt)
}
