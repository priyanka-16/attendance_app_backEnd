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
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/attendanceRecord"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/school"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/student"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/teacher"
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
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudentsList(storage))
	router.HandleFunc("POST /api/teachers", teacher.New(storage))
	router.HandleFunc("GET /api/teachers/{id}", teacher.GetById(storage))
	router.HandleFunc("GET /api/teachers", teacher.GetTeachersList(storage))
	router.HandleFunc("POST /api/attendanceRecords", attendanceRecord.New(storage))
	router.HandleFunc("GET /api/attendanceRecords/{id}", attendanceRecord.GetById(storage))
	router.HandleFunc("GET /api/attendanceRecords", attendanceRecord.GetAttendanceRecordsList(storage))
	router.HandleFunc("POST /api/schools", school.New(storage))
	router.HandleFunc("GET /api/schools/{id}", school.GetById(storage))
	router.HandleFunc("GET /api/schools", school.GetSchoolsList(storage))
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
