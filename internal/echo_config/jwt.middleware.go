package echo_config

import (
	"api-produtos/internal/state"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    state.CHAVE_PUBLICA_JWT,
		SigningMethod: "RS256",
	})
}
