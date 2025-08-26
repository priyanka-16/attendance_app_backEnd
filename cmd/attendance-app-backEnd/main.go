package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/priyanka-16/attendance-app-backEnd/internal/config"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/auth"
	"github.com/priyanka-16/attendance-app-backEnd/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()
	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized", slog.String("env", cfg.Env))
	//setup router
	router := http.NewServeMux()
	// Auth
	router.HandleFunc("POST /api/auth/request-otp", auth.RequestOTP(storage))
	router.HandleFunc("POST /api/auth/login-otp", auth.LoginWithOTP(storage))
	// // Profile
	// router.HandleFunc("GET /api/profile", profile.Get(storage))    // requires auth middleware
	// router.HandleFunc("PUT /api/profile", profile.Update(storage)) // requires auth middleware
	// // Home/Dashboard
	// router.HandleFunc("GET /api/home", home.Get(storage)) // requires auth middleware
	// // Attendance
	// router.HandleFunc("GET /api/attendance/calendar", attendance.Calendar(storage))
	// router.HandleFunc("GET /api/attendance/class/{sectionId}/date/{date}", attendance.GetClassAttendance(storage))
	// router.HandleFunc("POST /api/attendance/class/{sectionId}/date/{date}", attendance.SubmitClassAttendance(storage))
	// router.HandleFunc("GET /api/attendance/me", attendance.GetMyAttendance(storage))
	// router.HandleFunc("GET /api/attendance/summary/class/{sectionId}", attendance.GetClassSummary(storage))
	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))
	fmt.Printf("server started %s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
