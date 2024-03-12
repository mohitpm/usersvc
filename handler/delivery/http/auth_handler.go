package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	SignIn(ctx echo.Context) error
}

type authHandler struct {
}

func NewAuthHandler() AuthHandler {
	return &authHandler{}
}

func (a *authHandler) SignIn(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello, world")
}
