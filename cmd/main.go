package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/maximfedotov74/cloud-api/docs"
	"github.com/maximfedotov74/cloud-api/internal/cfg"
	"github.com/maximfedotov74/cloud-api/internal/handler"
	"github.com/maximfedotov74/cloud-api/internal/mw"
	"github.com/maximfedotov74/cloud-api/internal/repository"
	"github.com/maximfedotov74/cloud-api/internal/service"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/file"
	"github.com/maximfedotov74/cloud-api/internal/shared/jwt"
	"github.com/maximfedotov74/cloud-api/internal/shared/mail"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"os/signal"
	"syscall"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	runServer(ctx)
}

func initDeps(r *http.ServeMux, cfg *cfg.Config, dbClient *pgxpool.Pool, fileClient *file.FileClient, cron *gocron.Scheduler) {

	jwtService := jwt.NewJwtService(jwt.JwtConfig{
		RefreshTokenExp:    cfg.RefreshTokenExp,
		AccessTokenExp:     cfg.AccessTokenExp,
		RefreshTokenSecret: cfg.RefreshTokenSecret,
		AccessTokenSecret:  cfg.AccessTokenSecret,
	})

	mailService := mail.NewMailService(mail.MailConfig{SmtpKey: cfg.SmtpKey, SenderEmail: cfg.SmtpMail, SmtpHost: cfg.SmtpHost, SmtpPort: cfg.SmtpPort, AppLink: cfg.AppLink})

	sessionRepository := repository.NewSessionRepository(dbClient)

	roleRepository := repository.NewRoleRepository(dbClient)
	userRepository := repository.NewUserRepository(dbClient, roleRepository)
	folderRepository := repository.NewFolderRepository(dbClient)
	fileRepository := repository.NewFileRepository(dbClient)

	userService := service.NewUserService(userRepository, sessionRepository, jwtService, mailService, dbClient, fileClient)
	authService := service.NewAuthService(userService, sessionRepository, mailService, jwtService)
	folderService := service.NewFolderService(folderRepository)
	fileService := service.NewFileService(fileRepository, folderRepository, dbClient, fileClient)

	authMW := mw.NewAuthMW(userService, sessionRepository, jwtService)
	roleMw := mw.NewRoleMW(roleRepository)

	userHandler := handler.NewUserHandler(userService, r)
	authHandler := handler.NewAuthHandler(authService, r)
	helloHandler := handler.NewHelloHandler(cfg, r, authMW, roleMw)
	folderHandler := handler.NewFolderHandler(folderService, r, authMW)
	fileHandler := handler.NewFileHandler(fileService, r, authMW)

	authHandler.StartHandlers()
	userHandler.StartHandlers()
	folderHandler.StartHandlers()
	fileHandler.StartHandlers()
	helloHandler.StartHandlers()
}

func runServer(ctx context.Context) {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	config := cfg.MustLoadConfig()

	db := db.NewPostgresConnection(config.DatabaseUrl)

	fileClient := file.New(config.MinioApiUrl, config.MinioUser, config.MinioPassword)

	cron := gocron.NewScheduler(time.UTC)

	cron.StartAsync()
	log.Println("Scheduler service started successfully!")

	initDeps(mux, config, db, fileClient, cron)

	r := mw.ApplyLogger(mw.ApplyHeaders(mux))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Cannot start http server, because: %v", err)
		}
	}()

	log.Printf("Swagger Api docs working on : %s", "/swagger")
	log.Printf("Server started on PORT: %d", config.Port)
	<-ctx.Done()

	log.Println("Gracefully shutting down...")
	log.Println("Cleaning")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	db.Close()
	log.Println("Db connection closed")
	cron.Stop()
	log.Println("Scheduler service stopped")

	select {
	case <-shutdownCtx.Done():
		log.Fatalf("server shutdown: %v", ctx.Err())
	default:
		log.Println("Server shutdown successfully")
	}
}
