package controller

import (
	"api-produtos/model"
	"api-produtos/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AutenticacaoController struct {
	autenticacaoService *service.AutenticacaoService
}

func InstanciarAutenticacaoController(grupoV1 *echo.Group, autenticacaoService *service.AutenticacaoService) *AutenticacaoController {
	controller := &AutenticacaoController{
		autenticacaoService: autenticacaoService,
	}

	grupoAutenticacao := grupoV1.Group("/auth")

	grupoAutenticacao.POST("/login", controller.logar)

	return controller
}

func (c *AutenticacaoController) logar(contexto echo.Context) error {
	var (
		dadosLogin *model.LoginUsuario
	)

	if err := contexto.Bind(&dadosLogin); err != nil {
		return contexto.JSON(http.StatusBadRequest, "Dados inválidos")
	}

	usuarioAutenticado, erro := c.autenticacaoService.Logar(contexto.Request().Context(), dadosLogin.Email, dadosLogin.Senha)
	if erro != nil {
		return contexto.JSON(http.StatusUnauthorized, erro.Error())
	}

	return contexto.JSON(http.StatusOK, usuarioAutenticado)
}

func (c *AutenticacaoController) tratarErros(contexto echo.Context, erro error) error {
	switch {
	case errors.Is(erro, model.UsuarioNaoEncontrado):
		return contexto.JSON(http.StatusNotFound, erro.Error())
	case errors.Is(erro, model.SenhaIncorreta):
		return contexto.JSON(http.StatusUnauthorized, erro.Error())
	case errors.Is(erro, model.UsuarioDesativado):
		return contexto.JSON(http.StatusForbidden, erro.Error())
	default:
		return contexto.JSON(http.StatusInternalServerError, "Ocorreu um erro interno no servidor!")
	}
}
