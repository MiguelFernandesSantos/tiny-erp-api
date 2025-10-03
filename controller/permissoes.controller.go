package controller

import (
	"api-produtos/internal/echo_config"
	"api-produtos/service"

	"github.com/labstack/echo/v4"
)

type PermissoesController struct {
	permissoesService *service.PermissoesService
}

func InstanciarPermissoesController(grupoV1 *echo.Group, permissoesService *service.PermissoesService) *PermissoesController {
	controller := &PermissoesController{
		permissoesService: permissoesService,
	}

	grupoPermissoes := grupoV1.Group("/permissoes")
	grupoPermissoes.Use(echo_config.JWTMiddleware())

	// TODO end-point de permissoes

	return controller
}
