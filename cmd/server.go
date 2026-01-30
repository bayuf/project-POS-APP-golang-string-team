package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/wire"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

func APIserver(app *wire.App, config *utils.Configuration, logger *zap.Logger) {
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: app.Route.Handler(),
	}

	go func() {
		// service connections
		fmt.Println(config.AppName, "is running in port", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error serve server", zap.Error(err))
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	close(app.Stop)
	app.WG.Wait()

	log.Println("Shutdown", config.AppName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
