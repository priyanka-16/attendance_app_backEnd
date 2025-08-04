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
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/attendance"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/school"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/school_grade"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/school_grade_section"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/user"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/user_otp"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/user_student"
	"github.com/priyanka-16/attendance-app-backEnd/internal/http/handlers/user_teacher"
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
	router.HandleFunc("POST /api/users", user.New(storage))
	router.HandleFunc("GET /api/users/{id}", user.GetById(storage))
	router.HandleFunc("GET /api/users", user.GetUsersList(storage))
	router.HandleFunc("POST /api/students", user_student.New(storage))
	router.HandleFunc("GET /api/students/{id}", user_student.GetById(storage))
	router.HandleFunc("GET /api/students", user_student.GetUserStudentsList(storage))
	router.HandleFunc("POST /api/teachers", user_teacher.New(storage))
	router.HandleFunc("GET /api/teachers/{id}", user_teacher.GetById(storage))
	router.HandleFunc("GET /api/teachers", user_teacher.GetUserTeachersList(storage))
	router.HandleFunc("POST /api/otps", user_otp.New(storage))
	router.HandleFunc("GET /api/otps/{id}", user_otp.GetById(storage))
	router.HandleFunc("GET /api/otps", user_otp.GetUserOTPList(storage))
	router.HandleFunc("POST /api/attendance", attendance.New(storage))
	router.HandleFunc("GET /api/attendance/{id}", attendance.GetById(storage))
	router.HandleFunc("GET /api/attendance", attendance.GetAttendancesList(storage))
	router.HandleFunc("POST /api/schools", school.New(storage))
	router.HandleFunc("GET /api/schools/{id}", school.GetById(storage))
	router.HandleFunc("GET /api/schools", school.GetSchoolsList(storage))
	router.HandleFunc("POST /api/grades", school_grade.New(storage))
	router.HandleFunc("GET /api/grades/{id}", school_grade.GetById(storage))
	router.HandleFunc("GET /api/grades", school_grade.GetSchoolGradesList(storage))
	router.HandleFunc("POST /api/sections", school_grade_section.New(storage))
	router.HandleFunc("GET /api/sections/{id}", school_grade_section.GetById(storage))
	router.HandleFunc("GET /api/sections", school_grade_section.GetSchoolGradeSectionsList(storage))
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
