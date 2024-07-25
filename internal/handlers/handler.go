package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/back1ng1/messaggio-test/internal/entity"
	"gitlab.com/back1ng1/messaggio-test/internal/usecase"
)

type handler struct {
	uc usecase.Usecase
}

func New(uc usecase.Usecase) *echo.Echo {
	h := handler{uc: uc}
	e := echo.New()

	e.POST("/message", h.storeMessage)
	e.GET("/stats", h.getStats)

	return e
}

func (h handler) storeMessage(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	message, ok := m["message"]
	if !ok {
		return c.JSON(http.StatusUnprocessableEntity, "Incorrect message given")
	}

	msg, err := h.uc.StoreMessage(entity.Message{
		Message: message.(string),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to store message: %v", err))
	}

	return c.JSON(http.StatusOK, msg)
}

func (h handler) getStats(c echo.Context) error {
	return c.String(http.StatusOK, "TODO Statistics!!!")
}
