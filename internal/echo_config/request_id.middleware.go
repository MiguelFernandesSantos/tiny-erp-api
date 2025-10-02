package echo_config

import (
	"api-produtos/internal"
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := c.Request().Header.Get(echo.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		contexto := context.WithValue(c.Request().Context(), internal.ChaveContextoUnicoKey, reqID)
		c.SetRequest(c.Request().WithContext(contexto))

		c.Response().Header().Set(echo.HeaderXRequestID, reqID)

		return next(c)
	}
}
