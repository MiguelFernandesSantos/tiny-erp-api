package controller

import (
	"api-produtos/internal/echo_config"
	"api-produtos/model"
	"api-produtos/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UsuarioController struct {
	usuarioService *service.UsuarioService
}

func InstanciarUsuarioController(grupoV1 *echo.Group, usuarioService *service.UsuarioService) *UsuarioController {
	controller := &UsuarioController{
		usuarioService: usuarioService,
	}

	grupoUsuarios := grupoV1.Group("/usuario")
	grupoUsuarios.Use(echo_config.JWTMiddleware())

	grupoUsuarios.GET("", controller.obterUsuariosPaginado)
	grupoUsuarios.POST("", controller.adicionarUsuario)
	grupoUsuarios.GET("/:numeroUsuario", controller.obterUsuarioEspecifico)
	grupoUsuarios.DELETE("/:numeroUsuario", controller.desativarUsuarioEspecifico)

	return controller
}

func (c *UsuarioController) obterUsuariosPaginado(contexto echo.Context) error {
	var (
		pesquisa            = contexto.QueryParam("pesquisa")
		ultimoUsuarioString = contexto.QueryParam("ultimoUsuario")
		ultimoUsuario       *int64
	)

	if ultimoUsuarioString != "" {
		numero, erro := strconv.ParseInt(ultimoUsuarioString, 10, 64)
		if erro != nil {
			return contexto.JSON(http.StatusBadRequest, "Número do último usuário inválido")
		}
		ultimoUsuario = &numero
	}

	usuarios, erro := c.usuarioService.ObterUsuariosPaginado(contexto.Request().Context(), pesquisa, ultimoUsuario)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, usuarios)
}

func (c *UsuarioController) adicionarUsuario(contexto echo.Context) error {
	var (
		usuario *model.Usuario
	)

	if err := contexto.Bind(&usuario); err != nil {
		return contexto.JSON(http.StatusBadRequest, "Dados do usuário inválidos")
	}

	if err := usuario.Validar(); err != nil {
		return contexto.JSON(http.StatusBadRequest, err.Error())
	}

	erro := c.usuarioService.AdicionarUsuario(contexto.Request().Context(), usuario)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.NoContent(http.StatusCreated)
}

func (c *UsuarioController) obterUsuarioEspecifico(contexto echo.Context) error {
	var (
		numeroUsuarioString = contexto.Param("numeroUsuario")
	)

	if numeroUsuarioString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número do usuário é obrigatório")
	}

	numeroUsuario, erro := strconv.ParseInt(numeroUsuarioString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número do usuário inválido")
	}

	usuario, erro := c.usuarioService.ObterUsuarioEspecifico(contexto.Request().Context(), &numeroUsuario)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.JSON(http.StatusOK, usuario)
}

func (c *UsuarioController) desativarUsuarioEspecifico(contexto echo.Context) error {
	var (
		numeroUsuarioString = contexto.Param("numeroUsuario")
	)

	if numeroUsuarioString == "" {
		return contexto.JSON(http.StatusBadRequest, "Número do usuário é obrigatório")
	}

	numeroUsuario, erro := strconv.ParseInt(numeroUsuarioString, 10, 64)
	if erro != nil {
		return contexto.JSON(http.StatusBadRequest, "Número do usuário inválido")
	}

	erro = c.usuarioService.DesativarUsuarioEspecifico(contexto.Request().Context(), &numeroUsuario)
	if erro != nil {
		return c.tratarErros(contexto, erro)
	}

	return contexto.NoContent(http.StatusNoContent)
}

func (c *UsuarioController) tratarErros(contexto echo.Context, erro error) error {
	switch {
	case errors.Is(erro, model.UsuarioNaoEncontrado):
		return contexto.JSON(http.StatusNotFound, erro.Error())
	case errors.Is(erro, model.UsuarioJaExiste):
		return contexto.JSON(http.StatusConflict, erro.Error())
	default:
		return contexto.JSON(http.StatusInternalServerError, "Ocorreu um erro interno no servidor!")
	}
}
