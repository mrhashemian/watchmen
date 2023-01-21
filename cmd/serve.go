package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"watchmen/config"
	"watchmen/controller"
	"watchmen/database"
	logging "watchmen/log"
	"watchmen/monitor"
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
	linkRepo := repository.NewLinkRepository(baseAPIDB)

	mnt := monitor.NewMonitor(linkRepo, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		monitor.Run(ctx, mnt, config.C.Server.WorkerDuration)
	}()

	e := echo.New()
	e.Use(logging.EchoMiddleware())
	e.Use(middleware.Recover())

	e.GET("/", indexHandler)

	v1 := e.Group("/v1")
	v1Users := v1.Group("/users")

	v1Users.POST("/sign-up", controller.SignUp(userRepo))
	v1Users.POST("/login", controller.Login(userRepo))
	v1Users.GET("/:id/alert", controller.GetAlert(userRepo), echojwt.JWT(config.C.User.JWTSecret))

	v1Links := v1Users.Group("/:id/links", echojwt.JWT(config.C.User.JWTSecret))
	v1Links.POST("", controller.AddLink(linkRepo))
	v1Links.GET("", controller.GetLink(linkRepo))
	v1Links.GET("/:link_id", controller.RetrieveLink(linkRepo))

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

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("error in shutdown: %v", err)
	}
}

func indexHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, asciiArt)
}
