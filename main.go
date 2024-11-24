package main

import (
	"errors"
	"log/slog"
	"net/http"
	"receiptprocessor/db"
	"receiptprocessor/handler"
	"receiptprocessor/structs"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &structs.CustomValidator{Validator: validator.New()}

	e.GET("/receipts/:id/points", handler.GetPointsV1)
	e.POST("/receipts/process", handler.ProcessReceiptsV1)

	db.InitDB()

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
