package main

import (
	"fmt"
	"net/http"

	"github.com/KapDmitry/WB_L2/develop/dev11/internal/event"
	"github.com/KapDmitry/WB_L2/develop/dev11/internal/handlers"
	"github.com/KapDmitry/WB_L2/develop/dev11/internal/logger"
	"github.com/KapDmitry/WB_L2/develop/dev11/internal/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := logger.NewCustomLogger(zapConfig)
	if err != nil {
		panic(err.Error())
	}

	handler := &handlers.Handler{
		Repo:      event.NewInMemoryRepo(),
		Validator: validator.New(),
		Logger:    logger,
	}
	http.Handle("/create_event", middleware.LoggingMiddleware(http.HandlerFunc(handler.CreateEventHandler)))
	http.Handle("/update_event", middleware.LoggingMiddleware(http.HandlerFunc(handler.UpdateEventHandler)))
	http.Handle("/delete_event", middleware.LoggingMiddleware(http.HandlerFunc(handler.DeleteEventHandler)))
	http.Handle("/events_for_day", middleware.LoggingMiddleware(http.HandlerFunc(handler.EventsForHandler)))
	http.Handle("/events_for_week", middleware.LoggingMiddleware(http.HandlerFunc(handler.EventsForHandler)))
	http.Handle("/events_for_month", middleware.LoggingMiddleware(http.HandlerFunc(handler.EventsForHandler)))

	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("server didn't start")
	}
}
